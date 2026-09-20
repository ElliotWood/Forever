// The build arena: every build the sim can be given, ranked, computed ahead of time.
//
// The live rankings page runs one preset per spec in your browser and takes about forty
// seconds to do it. That is the ceiling of what a browser can be asked for, and it is a long
// way short of a leaderboard: the gear sets and rotations already in the repository multiply
// out to hundreds of builds. So those run headless in CI instead, and this page renders what
// they produced - instantly, with no sim in the browser at all.
//
// Which makes the cache question answer itself. The results file is the cache, the commit
// that produced it is the key, and the workflow that regenerates it on every push to master
// is the invalidation. Nothing has to notice the sim changed; the only thing that ever
// writes the file is a run that happened afterwards.

import { restsCell } from '../core/components/rests_cell';
import { SITE_BASE, SITE_VERSION } from '../core/constants/other';
import { Spec } from '../core/proto/common';
import { cssClassForClass, specNames, specToClass, titleIcons } from '../core/proto_utils/utils';
import { Composition } from '../core/spells/rests';
import { formatToNumber, formatToPercent } from '../core/utils';
import results from './results.json';

// The ui/ directory each spec's results are keyed by. The arena runs in Go, which knows the
// directory but not this enum, so the join happens here.
const SPECS: Record<string, Spec> = {
	balance_druid: Spec.SpecBalanceDruid,
	feral_druid: Spec.SpecFeralDruid,
	feral_tank_druid: Spec.SpecFeralTankDruid,
	elemental_shaman: Spec.SpecElementalShaman,
	enhancement_shaman: Spec.SpecEnhancementShaman,
	hunter: Spec.SpecHunter,
	mage: Spec.SpecMage,
	protection_paladin: Spec.SpecProtectionPaladin,
	retribution_paladin: Spec.SpecRetributionPaladin,
	rogue: Spec.SpecRogue,
	shadow_priest: Spec.SpecShadowPriest,
	smite_priest: Spec.SpecSmitePriest,
	tank_warrior: Spec.SpecTankWarrior,
	warlock: Spec.SpecWarlock,
	warrior: Spec.SpecWarrior,
};

type Build = {
	spec: string;
	build: string;
	talents: string;
	gear: string;
	rotation: string;
	dps: number;
	rests: Composition;
	/** Average item level of the gear set, and how many slots it fills. */
	ilvl: number;
	slots: number;
	optimised?: boolean;
	/** Only present when the build does not spend all 51 points. */
	points?: number;
};

const builds = results.builds as Array<Build>;

/**
 * Best build per spec, out of whatever the item level filter has left.
 *
 * Specs ship different numbers of gear sets - balance druid has a phase 2 BiS on file and
 * smite priest has nothing but launch - so ranking each spec's best run across all of them
 * ranks how far ahead somebody wrote its gear rather than the spec. Filtering by item level
 * is what makes the comparison mean something; the name on the gear file does not.
 */
const bestPerSpec = (from: Array<Build>): Array<Build> => {
	const best = new Map<string, Build>();
	for (const build of from) {
		if (!best.has(build.spec)) best.set(build.spec, build);
	}
	return [...best.values()].sort((a, b) => b.dps - a.dps);
};

/**
 * What the search was worth, on each row it appears.
 *
 * The leaderboard is searched builds top to bottom, which on its own hides the only question
 * anyone actually has - was the community build already right? The gain is that answer, and
 * for most specs it is not small.
 *
 * Matched on the exact gear set and rotation, never across two of them. A gain measured
 * across gear sets would be part talent search and part item level, which is the confusion
 * the item level column exists to remove.
 */
const key = (build: Build) => `${build.spec}|${build.gear}|${build.rotation}`;

const searchGains = (): Map<string, number> => {
	// Best written build per setup, so a searched row is measured against the strongest thing
	// a person put on file for that exact gear and rotation rather than whichever turned up
	// first. Keyed per row, not per spec: a searched build is re-run on every gear set the
	// spec has, and what those talents were worth is not the same number on each of them.
	const written = new Map<string, number>();
	for (const build of builds.filter(build => !build.optimised)) {
		written.set(key(build), Math.max(written.get(key(build)) ?? 0, build.dps));
	}

	const gains = new Map<string, number>();
	for (const searched of builds.filter(build => build.optimised)) {
		const base = written.get(key(searched));
		if (base) gains.set(`${key(searched)}|${searched.talents}`, searched.dps / base - 1);
	}
	return gains;
};

const gains = searchGains();

const specName = (spec: string) => (SPECS[spec] !== undefined ? specNames[SPECS[spec]] : spec);

/**
 * Item level brackets, because "same gear tier" was doing work it could not support.
 *
 * The leaderboard used to compare every spec's launch set on the grounds that launch is the
 * one tier they all have. They do all have one - and those sets run from item level 63.1 to
 * 70.0, which is most of a tier of difference sitting inside a table claiming to compare
 * specs. A balance druid's launch gear is seven levels below a feral tank's. So the gear is
 * now a number on every row and a filter above them, and "launch" is just a file name again.
 */
const BRACKETS: Array<{ key: string; label: string; holds: (ilvl: number) => boolean }> = [
	{ key: 'all', label: 'Any gear', holds: () => true },
	{ key: 'low', label: 'Under 63', holds: ilvl => ilvl < 63 },
	{ key: 'mid', label: '63 to 65', holds: ilvl => ilvl >= 63 && ilvl < 66 },
	{ key: 'high', label: '66 to 68', holds: ilvl => ilvl >= 66 && ilvl < 69 },
	{ key: 'top', label: '69 and up', holds: ilvl => ilvl >= 69 },
];

const specsIn = (bracket: (typeof BRACKETS)[number]) => new Set(builds.filter(b => bracket.holds(b.ilvl)).map(b => b.spec)).size;

export class ArenaPage {
	private readonly body: HTMLElement;
	private readonly count: HTMLElement;
	private readonly spread: HTMLElement;
	private showAll = false;
	// Defaults to the band covering the most specs, since a bracket holding three of them is
	// a narrower comparison than it looks.
	private bracket = BRACKETS.slice(1).reduce((best, b) => (specsIn(b) > specsIn(best) ? b : best), BRACKETS[1]);

	constructor(parent: HTMLElement) {
		const generated = new Date(results.generated);
		const stale = !!results.commit && !SITE_VERSION.startsWith(results.commit.slice(0, 7));

		parent.appendChild(
			<div id="arena-page">
				<header className="arena-header">
					<div className="container arena-header-container">
						<a href={SITE_BASE} className="arena-home-link">
							<img className="forever-logo" src={`${SITE_BASE}assets/img/forever_logo.png`} alt="World of Warcraft: Forever" />
						</a>
						<div className="arena-title-block">
							<h1 className="arena-title">The build arena</h1>
							<p className="arena-subtitle">
								Every talent build crossed with every gear set and every rotation this sim has on file: {String(builds.length)} builds across{' '}
								{String(new Set(builds.map(b => b.spec)).size)} specs, each run on its own at {formatToNumber(results.iterations)} iterations
								against the same target, with the same buffs and the same consumables, plus the builds a talent search found on top of those.
								Nothing is simulated in your browser.
							</p>
						</div>
					</div>
				</header>

				<main className="container arena-content">
					<section className="arena-intro">
						<div className="arena-intro-block">
							<h2 className="arena-intro-title">How a number gets onto this page</h2>
							<p>
								Every build here was simulated: a character is assembled, given a talent build, a gear set and a rotation, and run against the
								same target for {formatToNumber(results.iterations)} iterations. What comes out is the average damage per second of those runs.
								Nothing is estimated, interpolated or predicted - each row is the outcome of that build being played out five thousand times.
							</p>
							<p>
								The gear sets and rotations are files in the repository, written by people. The talent builds are the community ones from each
								spec's own page, plus whatever a search found on top of them. All of it runs headless when the sim changes, and the site ships
								the results, which is why the table is instant and why nothing is simulated in your browser.
							</p>
						</div>

						<div className="arena-intro-block">
							<h2 className="arena-intro-title">Where AI comes into it, and where it does not</h2>
							<p>
								<strong>Not into any number on this page.</strong> The damage figures come from a simulator - the{' '}
								<a href="https://github.com/wowsims/classic" target="_blank" rel="noreferrer">
									wowsims
								</a>{' '}
								engine, forked and adjusted for Forever. It is ordinary code doing arithmetic on the client's own data tables. No language model
								produces, adjusts or estimates a DPS figure, and the talent search is a hill climb that measures builds rather than reasons
								about them.
							</p>
							<p>
								<strong>Into the code, heavily.</strong> This sim's Forever changes, the talent search, this page and most of what surrounds
								them were written by an AI assistant working to one person's direction. That is worth saying plainly, because it is exactly the
								situation where confident-sounding output is cheap and being wrong is easy.
							</p>
							<p>
								So the checking is the point rather than an afterthought. <a href={`${SITE_BASE}evidence/`}>Every ability the sim registers</a>{' '}
								records where its numbers came from, and a test refuses to let one be added without that. The <strong>rests on a guess</strong>{' '}
								column on the right carries it through to here: it is how much of a build's damage depends on something nobody has confirmed.
								Where those numbers came out wrong, they came out wrong visibly - which is the only version of this that is worth publishing.
							</p>
						</div>
					</section>

					<div className="arena-controls">
						{this.toggle()}
						<div className="arena-brackets">{BRACKETS.map(bracket => this.bracketButton(bracket))}</div>
						<p className="arena-count" />
					</div>
					<p className="arena-spread" />

					<table className="metrics-table arena-table">
						<thead>
							<tr className="metrics-table-header-row">
								<th className="metrics-table-header-cell arena-rank-cell">#</th>
								<th className="metrics-table-header-cell arena-build-cell">Build</th>
								<th className="metrics-table-header-cell arena-setup-cell">Gear and rotation</th>
								<th className="metrics-table-header-cell arena-ilvl-cell">ilvl</th>
								<th className="metrics-table-header-cell arena-dps-cell">DPS</th>
								<th className="metrics-table-header-cell arena-share-cell">Share of top</th>
								<th className="metrics-table-header-cell rests-cell">Rests on a guess</th>
							</tr>
						</thead>
						<tbody className="metrics-table-body arena-body" />
					</table>

					<ul className="arena-notes">
						<li>
							<strong>Every build meets the same conditions.</strong> One target, one encounter length, one buff set, one set of consumables. That
							is what makes two numbers comparable - the live <a href={`${SITE_BASE}dps_rankings/`}>rankings page</a> achieves the same thing by
							putting everyone in one raid, which stops being possible at this count.
						</li>
						<li>
							<strong>Item level is the filter, not the file name.</strong> This table used to compare every spec on its "launch" gear set, on the
							grounds that launch is the one tier they all have. They do all have one - and those sets run from item level 63.1 to 70.0, which is
							most of a tier of difference sitting inside a table claiming to compare specs. A balance druid's launch gear is seven levels below a
							feral tank's. So gear is a number on every row now and a bracket above them, and the line under the filter says how far apart the
							rows you are looking at actually are.
						</li>
						<li>
							<strong>Rests on a guess</strong> is what the build's damage is made of, not a verdict on it. Each ability is weighted by its share
							of that build's damage and looked up in the <a href={`${SITE_BASE}evidence/`}>evidence manifest</a>. A build ten DPS ahead means
							something different if a quarter of it is unconfirmed. Hover the bar for the breakdown.
						</li>
						<li>
							<strong>Gear and rotations come from what is already here</strong> - the sets and priority lists on each spec's page. Nothing
							invents a better rotation than the ones people have written down, so a spec with one rotation on file gets one rotation ranked. That
							is a gap in the data, not a finding about the spec.
						</li>
						<li>
							<strong>Talents are searched, because they cannot be enumerated.</strong> A warrior has <strong>89,776,730,783,606,094</strong>{' '}
							builds it could actually spend - counted from the trees, enforcing rank caps, row gates and the prerequisite arrows. Legality is not
							what stands in the way: the arrows cut the count about fourfold, and what is left is two billion years of simulating at a second a
							build. Nor does throwing out the builds nobody would run. Count only the all-or-nothing ones, every talent maxed or untouched, which
							is roughly the shape of an optimal build, and a warrior still has 57,341,667 and a mage 1,261,940,421 - eighteen months and forty
							years. So each spec's best known build is improved one point at a time instead: price every point that could come out, price every
							point that could go in, make the best trade, repeat until no single move helps. Rows marked{' '}
							<span className="arena-found">found by search</span> came out of that, and hovering one shows its talent string.
						</li>
						<li>
							<strong>What that does and does not promise.</strong> It climbs from every distinct build the spec has on file rather than only its
							best one, because a climb goes to the nearest peak and one start reports the nearest peak to one build. Several starts agreeing is
							the cheapest evidence available that the peak is not merely nearby - it is still not proof that nothing higher exists somewhere else
							in the tree. It also only knows what this sim models: a talent flagged as unimplemented is worth zero here, so the search will
							happily empty it, and that is a fact about the sim rather than advice. Every build it reaches is checked against the game's own
							rules first - rank caps, row gates and prerequisites - so nothing in this table is a build you could not actually spend.
						</li>
						<li>
							<strong>A build that does not spend 51 points says so.</strong> One does: the mage Frost community build spends 49, so every Frost
							number on this site has been two points short of a character. It is left as written rather than quietly corrected - it is somebody
							else's build - but the searched row beside it shows what those two points are worth.
						</li>
						<li>
							Tank specs are measured on damage alone and healing specs are absent, because damage is the only axis this table has. A protection
							paladin at the bottom is not a bad tank - and a searched tank build is a tank build with the mitigation optimised out of it, so read
							those two rows as what the spec can do to a target dummy and nothing else.
						</li>
						<li>
							Computed from sim <code>{results.commit ? results.commit.slice(0, 7) : 'unknown'}</code> on {generated.toISOString().slice(0, 10)}.{' '}
							{stale ? (
								<strong className="arena-stale">
									The site has been rebuilt since, so these numbers are older than the sim running on the rest of the site.
								</strong>
							) : (
								'That is the sim this site is built from.'
							)}
						</li>
					</ul>
				</main>
			</div>,
		);

		this.body = parent.querySelector('.arena-body') as HTMLElement;
		this.count = parent.querySelector('.arena-count') as HTMLElement;
		this.spread = parent.querySelector('.arena-spread') as HTMLElement;
		this.render();
	}

	private toggle(): Element {
		const button = (
			<button className="arena-toggle" type="button">
				Show every build
			</button>
		) as HTMLButtonElement;
		button.addEventListener('click', () => {
			this.showAll = !this.showAll;
			button.textContent = this.showAll ? 'Show only the best of each spec' : 'Show every build';
			this.render();
		});
		return button;
	}

	private bracketButton(bracket: (typeof BRACKETS)[number]): Element {
		const specs = bracket.key === 'all' ? new Set(builds.map(b => b.spec)).size : specsIn(bracket);
		const button = (
			<button className={`arena-bracket${bracket.key === this.bracket.key ? ' arena-bracket-on' : ''}`} type="button">
				<span>{bracket.label}</span>
				<span className="arena-bracket-count">{String(specs)}</span>
			</button>
		) as HTMLButtonElement;

		button.addEventListener('click', () => {
			this.bracket = bracket;
			for (const other of document.querySelectorAll('.arena-bracket')) other.classList.remove('arena-bracket-on');
			button.classList.add('arena-bracket-on');
			this.render();
		});
		return button;
	}

	private render() {
		const inBracket = builds.filter(build => this.bracket.holds(build.ilvl));
		const shown = this.showAll ? inBracket : bestPerSpec(inBracket);
		const top = shown[0]?.dps || 1;

		this.count.textContent = this.showAll
			? `All ${shown.length} builds in this bracket, best first`
			: `The best build of each of ${new Set(shown.map(b => b.spec)).size} spec${new Set(shown.map(b => b.spec)).size === 1 ? '' : 's'} in this bracket.`;

		// The one thing a ranking table has to admit when it is not true: that the rows are not
		// wearing comparable gear. Three item levels is about a fifth of a tier, and a spec can
		// lose more than that to the gear alone.
		const levels = shown.map(build => build.ilvl).filter(ilvl => ilvl > 0);
		const low = Math.min(...levels);
		const high = Math.max(...levels);
		const wide = levels.length > 1 && high - low > 3;
		this.spread.textContent = levels.length
			? wide
				? `These rows span item level ${low.toFixed(1)} to ${high.toFixed(1)}, so some of the gap between them is gear rather than spec.`
				: `These rows span item level ${low.toFixed(1)} to ${high.toFixed(1)}.`
			: '';
		this.spread.classList.toggle('arena-spread-wide', wide);

		this.body.replaceChildren(...shown.map((build, index) => this.row(build, index + 1, top)));
	}

	private row(build: Build, rank: number, top: number): Element {
		const spec = SPECS[build.spec];
		const classColor = spec !== undefined ? cssClassForClass(specToClass[spec]) : 'white';
		const share = (build.dps / top) * 100;

		return (
			<tr className="arena-row">
				<td className="arena-rank-cell">{String(rank)}</td>
				<td className="arena-build-cell">
					{spec !== undefined ? <img className="metrics-action-icon" src={titleIcons[spec]} alt="" /> : <></>}
					<span className="arena-build-names">
						<span className="arena-spec-name">{specName(build.spec)}</span>
						<span className={`arena-talents text-${classColor}`}>{build.build}</span>
						{build.optimised ? (
							<span
								className="arena-found"
								attributes={{
									title: `${build.talents}\n\nClimbed from every distinct build this spec has on file for this gear and rotation.`,
								}}>
								found by search
								{gains.has(`${key(build)}|${build.talents}`)
									? ` ${formatToPercent(gains.get(`${key(build)}|${build.talents}`)! * 100, {
											maximumFractionDigits: 1,
											signDisplay: 'always',
									  })}`
									: ''}
							</span>
						) : (
							<></>
						)}
						{build.points ? (
							<span className="arena-short" attributes={{ title: 'This build does not spend every talent point a level 60 character has.' }}>
								{String(build.points)} of 51 points
							</span>
						) : (
							<></>
						)}
					</span>
				</td>
				<td className="arena-setup-cell">
					<span className="arena-setup">{build.gear}</span>
					<span className="arena-setup arena-setup-quiet">{build.rotation || 'default rotation'}</span>
				</td>
				<td className="arena-ilvl-cell">
					<span className="arena-ilvl">{build.ilvl ? build.ilvl.toFixed(1) : '?'}</span>
					{build.slots && build.slots < 15 ? (
						<span className="arena-slots" attributes={{ title: 'This gear set leaves slots empty, so the character is not fully equipped.' }}>
							{String(build.slots)} slots
						</span>
					) : (
						<></>
					)}
				</td>
				<td className="arena-dps-cell">{formatToNumber(build.dps, { maximumFractionDigits: 1, minimumFractionDigits: 1 })}</td>
				<td className="arena-share-cell">
					<div className="arena-share">
						<div className="arena-share-bar">
							<div className={`arena-share-bar-fill bg-${classColor}`} style={{ '--percentage': formatToPercent(share) }} />
						</div>
						<span className="arena-share-percent">{formatToPercent(share, { maximumFractionDigits: 1 })}</span>
					</div>
				</td>
				<td className="rests-cell">{restsCell(build.rests)}</td>
			</tr>
		) as Element;
	}
}
