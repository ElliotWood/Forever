import clsx from 'clsx';
import { ref } from 'tsx-vanilla';

import { Component } from '../core/components/component.js';
import { ContentBlock } from '../core/components/content_block.jsx';
import { NumberPicker } from '../core/components/number_picker.js';
import { SocialLinks } from '../core/components/social_links.jsx';
import { SITE_BASE, SITE_VERSION } from '../core/constants/other.js';
import { IndividualSimUIConfig, RaidSimPreset } from '../core/individual_sim_ui.js';
import { LaunchStatus, simLaunchStatuses } from '../core/launched_sims.js';
import { getSpecConfig } from '../core/player.js';
import { ErrorOutcomeType, ProgressMetrics, Raid as RaidProto } from '../core/proto/api.js';
import { Class, Spec } from '../core/proto/common.js';
import { SimResult } from '../core/proto_utils/sim_result.js';
import { MAX_PARTY_SIZE } from '../core/party.js';
import {
	cssClassForClass,
	getTalentTree,
	getTalentTreePoints,
	makeDefaultBlessings,
	naturalSpecOrder,
	specNames,
	specToClass,
	titleIcons,
} from '../core/proto_utils/utils.js';
import { MAX_NUM_PARTIES } from '../core/raid.js';
import { Sim, SimError } from '../core/sim.js';
import { EventID, TypedEvent } from '../core/typed_event.js';
import { formatToNumber, formatToPercent } from '../core/utils.js';
import { applyBlessings, applyNewPlayerAssignments, communityBuilds, newPlayerFromPreset, playerPresets, RaidBuild } from '../raid/presets.js';

// Twenty-six players over a two minute encounter, so this is not free: measured at about
// forty seconds from opening the page to the table appearing, on four cores. The run has
// to finish while someone is still looking at it, and at this count the noise on a single
// build is a few DPS - far under the gaps the ranking is showing. The picker raises it for
// anyone who wants the error smaller than that.
const DEFAULT_ITERATIONS = 1000;

const isLaunched = (spec: Spec) => simLaunchStatuses[spec].status != LaunchStatus.Unlaunched;

type Ranking = {
	build: RaidBuild;
	dps: number;
};

// Every spec's own sim starts with the raid buffs and debuffs that spec assumes somebody
// else in the raid is providing. A ranking has to hand all of them the same set, so take
// the strongest of each: nobody goes without a buff their own sim would have had, and
// nobody gets one that no launched spec brings. Tristate fields keep the better version.
function strongestOf<T extends object>(buffs: Array<T>): T {
	const merged: Record<string, boolean | number> = {};
	buffs.forEach(buff =>
		Object.entries(buff as Record<string, boolean | number>).forEach(([field, value]) => {
			if (typeof value === 'boolean') {
				merged[field] = !!merged[field] || value;
			} else {
				merged[field] = Math.max((merged[field] as number) || 0, value);
			}
		}),
	);
	return merged as T;
}

export class DpsRankings extends Component {
	readonly sim: Sim;

	private readonly builds: Array<RaidBuild>;
	private readonly runButton: HTMLButtonElement;
	private readonly statusElem: HTMLElement;
	private readonly resultsElem: HTMLElement;
	private readonly iterationsElem: HTMLElement;

	private running = false;

	constructor(parentElem: HTMLElement) {
		super(parentElem, 'dps-rankings-ui');
		this.sim = new Sim();
		this.builds = communityBuilds();

		const runButtonRef = ref<HTMLButtonElement>();
		const statusRef = ref<HTMLDivElement>();
		const resultsRef = ref<HTMLDivElement>();
		const iterationsRef = ref<HTMLSpanElement>();
		const controlsRef = ref<HTMLDivElement>();
		const socialsRef = ref<HTMLDivElement>();
		const provenanceRef = ref<HTMLDivElement>();
		const notesRef = ref<HTMLDivElement>();

		this.rootElem.appendChild(
			<>
				<div className="sim-bg" />
				<div className="container dps-rankings-container">
					<header className="dps-rankings-header">
						<a href={SITE_BASE} className="dps-rankings-home">
							<img className="dps-rankings-logo" src={`${SITE_BASE}assets/img/forever_logo.png`} alt="World of Warcraft: Forever" />
						</a>
						<div ref={socialsRef} className="dps-rankings-socials" />
					</header>
					<main className="dps-rankings-main">
						<h1 className="dps-rankings-title">Damage comparison</h1>
						<p className="dps-rankings-lead fs-5">
							Every community build of every launched spec in one raid, damage only. One run, one encounter, one set of buffs, so the numbers can
							sit in the same table. It is here to find bugs in this sim, not to say which spec is stronger in Forever.
						</p>
						<div ref={provenanceRef} className="dps-rankings-provenance" />
						<div ref={notesRef} className="dps-rankings-notes" />
						<div className="dps-rankings-controls">
							<button ref={runButtonRef} className="btn btn-primary dps-rankings-run">
								Run again
							</button>
							<div ref={controlsRef} className="dps-rankings-iterations" />
							<span ref={iterationsRef} className="dps-rankings-iterations-note" />
						</div>
						<div ref={statusRef} className="dps-rankings-status" />
						<div ref={resultsRef} className="dps-rankings-results" />
					</main>
					<footer className="dps-rankings-footer">WoW Forever sim {SITE_VERSION}</footer>
				</div>
			</>,
		);

		this.runButton = runButtonRef.value!;
		this.statusElem = statusRef.value!;
		this.resultsElem = resultsRef.value!;
		this.iterationsElem = iterationsRef.value!;

		const socials = socialsRef.value!;
		socials.appendChild(SocialLinks.buildDiscordLink());
		socials.appendChild(SocialLinks.buildGitHubLink());

		this.buildProvenance(provenanceRef.value!);
		this.buildNotes(notesRef.value!);

		new NumberPicker(controlsRef.value!, this.sim, {
			id: 'dps-rankings-iterations',
			label: 'Iterations',
			inline: true,
			positive: true,
			extraCssClasses: ['dps-rankings-iterations-picker'],
			changedEvent: (sim: Sim) => sim.iterationsChangeEmitter,
			getValue: (sim: Sim) => sim.getIterations(),
			setValue: (eventID: EventID, sim: Sim, newValue: number) => {
				sim.setIterations(eventID, newValue);
			},
		});

		this.runButton.addEventListener('click', () => this.run());

		this.setStatus('Loading sim...');
		// The encounter's default target and every preset's gear are looked up in the item
		// database, so there is no raid to assemble until that has loaded.
		this.sim.waitForInit().then(() => {
			this.buildRaid();
			this.run();
		});
	}

	// Where the numbers come from matters more than how the raid was assembled, so it goes
	// above the fold rather than into the notes below. Nothing in Forever has been measured
	// in a client: every talent here was read off a BlizzCon demo tooltip, and the values
	// the demo never showed are assumptions this sim writes down rather than hides.
	private buildProvenance(parentElem: HTMLElement) {
		parentElem.appendChild(
			<div className="dps-rankings-provenance-block">
				<h2 className="dps-rankings-provenance-title">These numbers are provisional</h2>
				<p className="dps-rankings-provenance-body">
					Forever has no beta client yet. Every talent, rank and coefficient in this sim was read off a BlizzCon 2026 demo tooltip, and the demo
					mostly showed rank 1, so the rest is assumed - linear unless a tooltip said otherwise. Those assumptions are counted and listed in the{' '}
					<a href="https://github.com/ElliotWood/Forever/blob/master/docs/forever_beta_checklist.md" target="_blank" rel="noreferrer">
						beta re-verification checklist
					</a>
					, one line each, and every one of them can move a number in this table.
				</p>
				<p className="dps-rankings-provenance-body">
					This is an unofficial fork, not the official Forever sim, and the table is a self-check: run every build under identical conditions and a
					build the sim is getting wrong stands out. That is what it has been for. It is how the feral cat was found stuck at 398 with eight gear
					slots empty, the enhancement shaman re-dropping one totem until it ran out of mana, and the retribution paladin carrying a Holy Strike
					invented at nearly three times its published damage. Each of those is now fixed and written up in the{' '}
					<a href={`${SITE_BASE}changelog/`}>changelog</a>.
				</p>
				<p className="dps-rankings-provenance-body">
					There is a second limit under the first one, and the beta will not lift it. A tooltip gives a spell's damage at one level; what a sim needs
					is how that damage is built out of attack power, spell power and weapon speed. None of those relationships was published, so where Forever
					has not changed a spell this sim assumes Classic's, and where it has, the coefficient is inferred from the one number the demo showed. The
					beta is capped at level 30, so it will not settle the level 60 ranks or the scaling either. Both are the sort of thing only the live client
					or the people building the game can answer.
				</p>
				<p className="dps-rankings-provenance-body">
					So: do not pick a main off this table, and do not quote it as a Forever balance claim. It is a place to catch the sim getting something
					obviously wrong, and nothing more than that until there is real data to check it against.
				</p>
			</div>,
		);
	}

	// The caveats belong on the page, not in a commit message: a ranking that doesn't say
	// what it left out is worse than no ranking.
	private buildNotes(parentElem: HTMLElement) {
		// One name per spec, from the preset's own tooltip rather than specNames, because
		// that calls the healing priest just 'Priest' and reads as the whole class missing.
		const unimplemented = [...new Set(playerPresets.map(preset => preset.spec))]
			.filter(spec => !isLaunched(spec))
			.map(spec => playerPresets.find(preset => preset.spec == spec)!.tooltip)
			.join(', ');
		// Launched specs the raid picker has no preset for. Hunter is one today: its raid
		// presets are commented out upstream, so there is no build to put in the raid.
		const withoutPreset = naturalSpecOrder
			.filter(spec => isLaunched(spec) && !playerPresets.some(preset => preset.spec == spec))
			.map(spec => specNames[spec])
			.join(', ');

		const specCount = new Set(this.builds.map(build => build.preset.spec)).size;

		const notes = new ContentBlock(parentElem, 'dps-rankings-notes-block', {
			header: { title: 'What these numbers are, and are not' },
		});
		notes.bodyElement.appendChild(
			<ul className="dps-rankings-notes-list">
				<li>
					All {this.builds.length} builds are simulated <strong>together, as one {this.builds.length}-player raid</strong>: one run, one encounter,
					one set of buffs, so every number in this table comes from the same fight. Those {this.builds.length} players span {specCount} specs.
					Running each build on its own would give numbers that cannot honestly be put side by side, because each spec's own defaults differ.
				</li>
				<li>
					Each player is a <strong>community talent build</strong> on its spec's raid preset - the gear, consumes and rotation the raid picker drops
					into a slot, with the build's talents in place of the preset's. Two builds of one spec therefore differ by talents alone; the presets are
					shared defaults, not per-build optimised gear or rotations, and some specs' presets are better tuned than others.
				</li>
				<li>
					{unimplemented || withoutPreset ? (
						<>
							Absent, because there is nothing to put in the raid:{' '}
							{unimplemented ? (
								<>
									<strong>{unimplemented}</strong> have no working sim yet
								</>
							) : (
								''
							)}
							{unimplemented && withoutPreset ? ', and ' : ''}
							{withoutPreset ? (
								<>
									the raid sim has no preset build for <strong>{withoutPreset}</strong>
								</>
							) : (
								''
							)}
							. The raid is therefore missing its healers, and the tank specs that are present are measured on damage alone.
						</>
					) : (
						'Every launched spec has a build in the raid. The raid is still missing its healers, and the tank specs that are present are measured on damage alone.'
					)}
				</li>
				<li>
					Every build wears its spec's <strong>Launch</strong> gear set: the best pre-raid gear in the launch item pool, picked by that spec's own
					stat weights, with raid drops left out. Same tier, chosen the same way, so the table compares specs and not gear.
				</li>
				<li>
					Raid buffs, party buffs and debuffs are the strongest of what each launched spec's own sim assumes by default, given to everyone alike, so
					nobody is missing a buff it expects. Blessings come from the paladins actually in the raid; innervates and power infusions are off, because
					nobody in the raid is casting them. <strong>No world buffs</strong>: they do not work inside Forever raids.
				</li>
				<li>
					Every run is a fresh simulation in your browser at the iteration count below. Fewer iterations means a noisier comparison; raise it if two
					specs are close.
				</li>
			</ul>,
		);
	}

	// Mirrors the raid sim's own defaults so a build's number here means the same thing it
	// would mean in the raid sim, then fills a slot per build. More builds than a five party
	// raid holds, so open as many parties as they need.
	private buildRaid() {
		const eventID = TypedEvent.nextEventID();
		const specs = [...new Set(this.builds.map(build => build.preset.spec))];
		const specDefaults = specs.map(spec => (getSpecConfig(spec) as IndividualSimUIConfig<any>).defaults);
		const numActiveParties = Math.min(MAX_NUM_PARTIES, Math.ceil(this.builds.length / MAX_PARTY_SIZE));

		TypedEvent.freezeAllAndDo(() => {
			this.sim.raid.fromProto(eventID, RaidProto.create({ numActiveParties }));
			this.sim.encounter.applyDefaults(eventID);
			this.sim.applyDefaults(eventID, true, true);
			this.sim.setShowDamageMetrics(eventID, true);
			this.sim.setIterations(eventID, DEFAULT_ITERATIONS);

			this.sim.raid.setBuffs(eventID, strongestOf(specDefaults.map(defaults => defaults.raidBuffs)));
			// Every spec's defaults describe the curses somebody else is keeping up, and none of
			// them asks for Curse of Shadow, so the union had Curse of Elements and not its twin.
			// This raid has four warlocks in it, so both curses are up.
			const debuffs = strongestOf(specDefaults.map(defaults => defaults.debuffs));
			if (this.builds.some(build => specToClass[build.preset.spec] == Class.ClassWarlock)) {
				debuffs.curseOfElements = true;
				debuffs.curseOfShadow = true;
			}
			this.sim.raid.setDebuffs(eventID, debuffs);
			const partyBuffs = strongestOf(specDefaults.map(defaults => defaults.partyBuffs));
			this.sim.raid.getParties().forEach(party => party.setBuffs(eventID, partyBuffs));

			this.builds.forEach((build, index) => {
				const player = newPlayerFromPreset(eventID, this.sim, build.preset);
				player.setTalentsString(eventID, build.talentsString);
				player.setName(eventID, build.name);
				this.sim.raid.setPlayer(eventID, index, player);
				applyNewPlayerAssignments(eventID, player, this.sim.raid);
			});
		});

		const blessings = makeDefaultBlessings(this.sim.raid.getClassCount(Class.ClassPaladin));
		this.sim.setModifyRaidProto(raidProto => applyBlessings(raidProto, blessings, this.sim.raid.getClassCount(Class.ClassPaladin)));
	}

	private async run() {
		if (this.running) {
			return;
		}
		this.running = true;
		this.runButton.disabled = true;
		this.setStatus('Simulating...');

		try {
			const result = await this.sim.runRaidSim(TypedEvent.nextEventID(), (progress: ProgressMetrics) => this.setProgress(progress));
			if (result instanceof SimResult) {
				this.setResults(result);
			} else if (result.type != ErrorOutcomeType.ErrorOutcomeAborted) {
				this.setStatus('The sim stopped before it finished. Try running it again.');
			}
		} catch (error) {
			console.error(error);
			this.setStatus(error instanceof SimError ? error.errorStr : 'Something went wrong running the sim. Reload the page and try again.');
		} finally {
			this.running = false;
			this.runButton.disabled = false;
		}
	}

	private setStatus(message: string) {
		this.statusElem.replaceChildren(<span className="dps-rankings-status-text">{message}</span>);
	}

	private setProgress(progress: ProgressMetrics) {
		const total = progress.totalIterations || this.sim.getIterations();
		this.statusElem.replaceChildren(
			<>
				<div className="loader dps-rankings-loader" />
				<span className="dps-rankings-status-text">
					{progress.presimRunning
						? 'Presimulations running...'
						: `${formatToNumber(progress.completedIterations)} / ${formatToNumber(total)} iterations`}
				</span>
			</>,
		);
	}

	private setResults(result: SimResult) {
		// Builds of one spec share it, so match each player back by its raid slot, not its spec.
		const rankings: Array<Ranking> = this.builds
			.map((build, index) => ({ build, dps: result.getPlayerWithRaidIndex(index)?.dps.avg || 0 }))
			.sort((a, b) => b.dps - a.dps);
		const topDps = rankings[0]?.dps || 1;

		this.iterationsElem.textContent = `${formatToNumber(result.iterations)} iterations, ${formatToNumber(result.duration)}s encounter`;
		this.statusElem.replaceChildren();
		this.resultsElem.replaceChildren(
			<table className="metrics-table dps-rankings-table">
				<thead>
					<tr className="metrics-table-header-row">
						<th className="metrics-table-header-cell dps-rankings-spec-cell">Build</th>
						<th className="metrics-table-header-cell dps-rankings-dps-cell">DPS</th>
						<th className="metrics-table-header-cell dps-rankings-share-cell">Share of top</th>
					</tr>
				</thead>
				<tbody className="metrics-table-body">{rankings.map(ranking => this.buildRow(ranking, topDps))}</tbody>
			</table>,
		);
	}

	private buildRow(ranking: Ranking, topDps: number): Element {
		const spec = ranking.build.preset.spec;
		const classColor = cssClassForClass(specToClass[spec]);
		const share = (ranking.dps / topDps) * 100;

		return (
			<tr className="dps-rankings-row">
				<td className="dps-rankings-spec-cell">
					<img className="metrics-action-icon" src={titleIcons[spec]} alt="" />
					<span className="dps-rankings-spec-names">
						<span className="dps-rankings-class-name">{specNames[spec]}</span>
						<span className={`dps-rankings-spec-name text-${classColor}`}>{ranking.build.name}</span>
					</span>
				</td>
				<td className="dps-rankings-dps-cell">{formatToNumber(ranking.dps, { maximumFractionDigits: 1, minimumFractionDigits: 1 })}</td>
				<td className="dps-rankings-share-cell">
					<div className="dps-rankings-share">
						<div className="dps-rankings-share-bar">
							<div className={clsx('dps-rankings-share-bar-fill', `bg-${classColor}`)} style={{ '--percentage': formatToPercent(share) }} />
						</div>
						<span className="dps-rankings-share-percent">{formatToPercent(share, { maximumFractionDigits: 1 })}</span>
					</div>
				</td>
			</tr>
		) as Element;
	}
}
