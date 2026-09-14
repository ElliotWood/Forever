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
import { Class, IndividualBuffs, Spec } from '../core/proto/common.js';
import { SimResult } from '../core/proto_utils/sim_result.js';
import { classNames, cssClassForClass, makeDefaultBlessings, naturalSpecOrder, specNames, specToClass, titleIcons } from '../core/proto_utils/utils.js';
import { Sim, SimError } from '../core/sim.js';
import { EventID, TypedEvent } from '../core/typed_event.js';
import { formatToNumber, formatToPercent } from '../core/utils.js';
import { applyBlessings, applyNewPlayerAssignments, newPlayerFromPreset, playerPresets } from '../raid/presets.js';

// Fourteen players over a two minute encounter, so this is not free: measured at about
// twenty seconds from opening the page to the table appearing, on four cores. The run has
// to finish while someone is still looking at it, and at this count the noise on a single
// spec is a few DPS - far under the gaps the ranking is showing. The picker raises it for
// anyone who wants the error smaller than that.
const DEFAULT_ITERATIONS = 1000;

const isLaunched = (spec: Spec) => simLaunchStatuses[spec].status != LaunchStatus.Unlaunched;

// One raid slot per spec. Several specs ship more than one raid preset - three mage
// builds, two rogue builds - which differ only in name and icon, so the first one stands
// for the spec the way it does in the raid picker's class menu.
const rankedSpecs = (): Array<Spec> => [...new Set(playerPresets.map(preset => preset.spec))].filter(isLaunched);

type Ranking = {
	spec: Spec;
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

// World buffs are per-player, so a raid assembled from presets has none - but every
// spec's own sim turns them on by default, and they are worth far more to the physical
// specs than to the casters. Leaving them off would quietly push melee down the table,
// so give everyone the same union of the world buffs the specs ask for. The rest of
// IndividualBuffs stays empty on purpose: blessings belong to the paladins in the raid,
// and innervates and power infusions have to be cast by somebody who is in it.
function worldBuffsFor(specDefaults: Array<IndividualSimUIConfig<any>['defaults']>): IndividualBuffs {
	const union = strongestOf(specDefaults.map(defaults => defaults.individualBuffs));
	return IndividualBuffs.create({
		rallyingCryOfTheDragonslayer: union.rallyingCryOfTheDragonslayer,
		saygesFortune: union.saygesFortune,
		spiritOfZandalar: union.spiritOfZandalar,
		songflowerSerenade: union.songflowerSerenade,
		warchiefsBlessing: union.warchiefsBlessing,
		fengusFerocity: union.fengusFerocity,
		moldarsMoxie: union.moldarsMoxie,
		slipkiksSavvy: union.slipkiksSavvy,
	});
}

export class DpsRankings extends Component {
	readonly sim: Sim;

	private readonly presets: Array<RaidSimPreset<any>>;
	private readonly runButton: HTMLButtonElement;
	private readonly statusElem: HTMLElement;
	private readonly resultsElem: HTMLElement;
	private readonly iterationsElem: HTMLElement;

	private running = false;

	constructor(parentElem: HTMLElement) {
		super(parentElem, 'dps-rankings-ui');
		this.sim = new Sim();
		this.presets = rankedSpecs().map(spec => playerPresets.find(preset => preset.spec == spec)!);

		const runButtonRef = ref<HTMLButtonElement>();
		const statusRef = ref<HTMLDivElement>();
		const resultsRef = ref<HTMLDivElement>();
		const iterationsRef = ref<HTMLSpanElement>();
		const controlsRef = ref<HTMLDivElement>();
		const socialsRef = ref<HTMLDivElement>();
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
						<h1 className="dps-rankings-title">DPS Rankings</h1>
						<p className="dps-rankings-lead fs-5">
							Every launched spec in one raid, damage only, highest first. One run, one encounter, one set of buffs, so the numbers can sit in the
							same table.
						</p>
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
		socials.appendChild(SocialLinks.buildPatreonLink());

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

		const notes = new ContentBlock(parentElem, 'dps-rankings-notes-block', {
			header: { title: 'What these numbers are, and are not' },
		});
		notes.bodyElement.appendChild(
			<ul className="dps-rankings-notes-list">
				<li>
					All {this.presets.length} specs are simulated <strong>together, in a single raid</strong>, sharing one encounter, duration and set of buffs.
					Running each spec on its own gives numbers that cannot honestly be put side by side, because each spec's own defaults differ.
				</li>
				<li>
					Each player is the <strong>raid sim's preset build</strong> for its spec - the same one the raid picker drops into a slot. Those are shared
					defaults, not per-spec optimised gear, talents or rotations, and some specs' presets are better tuned than others.
				</li>
				<li>
					Absent, because there is nothing to put in the raid: <strong>{unimplemented}</strong> have no working sim yet, and the raid sim has no
					preset build for <strong>{withoutPreset}</strong>. The raid is therefore missing its healers, and the tank specs that are present are ranked
					on damage alone.
				</li>
				<li>
					Raid buffs, party buffs, debuffs and world buffs are the strongest of what each launched spec's own sim assumes by default, given to
					everyone alike, so nobody is missing a buff it expects. Blessings come from the paladins actually in the raid; innervates and power
					infusions are off, because nobody in the raid is casting them.
				</li>
				<li>
					Every run is a fresh simulation in your browser at the iteration count below. Fewer iterations means a noisier ranking; raise it if two
					specs are close.
				</li>
			</ul>,
		);
	}

	// Mirrors the raid sim's own defaults so a spec's number here means the same thing it
	// would mean in the raid sim, then fills every slot from the launched presets.
	private buildRaid() {
		const eventID = TypedEvent.nextEventID();
		const specDefaults = this.presets.map(preset => (getSpecConfig(preset.spec) as IndividualSimUIConfig<any>).defaults);

		TypedEvent.freezeAllAndDo(() => {
			this.sim.raid.fromProto(eventID, RaidProto.create({ numActiveParties: 5 }));
			this.sim.encounter.applyDefaults(eventID);
			this.sim.applyDefaults(eventID, true, true);
			this.sim.setShowDamageMetrics(eventID, true);
			this.sim.setIterations(eventID, DEFAULT_ITERATIONS);

			this.sim.raid.setBuffs(eventID, strongestOf(specDefaults.map(defaults => defaults.raidBuffs)));
			this.sim.raid.setDebuffs(eventID, strongestOf(specDefaults.map(defaults => defaults.debuffs)));
			const partyBuffs = strongestOf(specDefaults.map(defaults => defaults.partyBuffs));
			this.sim.raid.getParties().forEach(party => party.setBuffs(eventID, partyBuffs));

			const worldBuffs = worldBuffsFor(specDefaults);
			this.presets.forEach((preset, index) => {
				const player = newPlayerFromPreset(eventID, this.sim, preset);
				player.setBuffs(eventID, worldBuffs);
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
		const bySpec = new Map<Spec, number>();
		result.getPlayers().forEach(player => bySpec.set(player.spec, player.dps.avg));

		const rankings: Array<Ranking> = this.presets.map(preset => ({ spec: preset.spec, dps: bySpec.get(preset.spec) || 0 })).sort((a, b) => b.dps - a.dps);
		const topDps = rankings[0]?.dps || 1;

		this.iterationsElem.textContent = `${formatToNumber(result.iterations)} iterations, ${formatToNumber(result.duration)}s encounter`;
		this.statusElem.replaceChildren();
		this.resultsElem.replaceChildren(
			<table className="metrics-table dps-rankings-table">
				<thead>
					<tr className="metrics-table-header-row">
						<th className="metrics-table-header-cell dps-rankings-rank-cell">#</th>
						<th className="metrics-table-header-cell dps-rankings-spec-cell">Spec</th>
						<th className="metrics-table-header-cell dps-rankings-dps-cell">DPS</th>
						<th className="metrics-table-header-cell dps-rankings-share-cell">Share of top</th>
					</tr>
				</thead>
				<tbody className="metrics-table-body">{rankings.map((ranking, index) => this.buildRow(ranking, index + 1, topDps))}</tbody>
			</table>,
		);
	}

	private buildRow(ranking: Ranking, rank: number, topDps: number): Element {
		const spec = ranking.spec;
		const classColor = cssClassForClass(specToClass[spec]);
		const share = (ranking.dps / topDps) * 100;

		return (
			<tr className="dps-rankings-row">
				<td className="dps-rankings-rank-cell">{rank}</td>
				<td className="dps-rankings-spec-cell">
					<img className="metrics-action-icon" src={titleIcons[spec]} alt="" />
					<span className="dps-rankings-spec-names">
						<span className="dps-rankings-class-name">{classNames[specToClass[spec]]}</span>
						<span className={`dps-rankings-spec-name text-${classColor}`}>{specNames[spec]}</span>
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
