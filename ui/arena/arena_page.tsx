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
};

const builds = results.builds as Array<Build>;

/**
 * Best build per spec, on launch gear only.
 *
 * Not a detail. Specs ship different numbers of gear sets - balance druid has a phase 2 BiS
 * on file and smite priest has nothing but launch - so ranking each spec's best run across
 * all of them ranks how far ahead somebody wrote its gear, not the spec. Every spec has a
 * launch set, so that is the one tier they can all be compared on. The full arena below
 * still shows every set; it just does not pretend the resulting order means one thing.
 */
const bestPerSpec = (): Array<Build> => {
	const best = new Map<string, Build>();
	for (const build of builds.filter(build => build.gear.includes('launch'))) {
		if (!best.has(build.spec)) best.set(build.spec, build);
	}
	return [...best.values()].sort((a, b) => b.dps - a.dps);
};

const specName = (spec: string) => (SPECS[spec] !== undefined ? specNames[SPECS[spec]] : spec);

export class ArenaPage {
	private readonly body: HTMLElement;
	private readonly count: HTMLElement;
	private showAll = false;

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
								against the same target, with the same buffs and the same consumables. Nothing is simulated in your browser.
							</p>
						</div>
					</div>
				</header>

				<main className="container arena-content">
					<div className="arena-controls">
						{this.toggle()}
						<p className="arena-count" />
					</div>

					<table className="metrics-table arena-table">
						<thead>
							<tr className="metrics-table-header-row">
								<th className="metrics-table-header-cell arena-rank-cell">#</th>
								<th className="metrics-table-header-cell arena-build-cell">Build</th>
								<th className="metrics-table-header-cell arena-setup-cell">Gear and rotation</th>
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
							<strong>The leaderboard is launch gear only.</strong> Specs ship different numbers of gear sets - balance druid has a phase 2 BiS on
							file, smite priest has nothing but launch - so taking each spec's best run across all of them would rank how far ahead somebody
							wrote its gear rather than the spec. Launch is the one tier every spec has. Show every build to see the rest; that view is an
							inventory of what has been simulated, not a like-for-like table.
						</li>
						<li>
							<strong>Rests on a guess</strong> is what the build's damage is made of, not a verdict on it. Each ability is weighted by its share
							of that build's damage and looked up in the <a href={`${SITE_BASE}evidence/`}>evidence manifest</a>. A build ten DPS ahead means
							something different if a quarter of it is unconfirmed. Hover the bar for the breakdown.
						</li>
						<li>
							<strong>Builds come from what is already here</strong> - the community talent builds on each spec's page, its gear sets and its
							rotations. Nothing searches for a better build than the ones people have written down, so a spec with one rotation on file gets one
							rotation ranked. That is a gap in the data, not a finding about the spec.
						</li>
						<li>
							Tank specs are measured on damage alone and healing specs are absent, because damage is the only axis this table has. A protection
							paladin at the bottom is not a bad tank.
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

	private render() {
		const shown = this.showAll ? builds : bestPerSpec();
		const top = shown[0]?.dps || 1;
		this.count.textContent = this.showAll
			? `All ${builds.length} builds, best first - mixed gear tiers, so not a like-for-like ranking`
			: `The best launch-gear build of each of ${shown.length} specs. ${builds.length} builds were run in total.`;
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
					</span>
				</td>
				<td className="arena-setup-cell">
					<span className="arena-setup">{build.gear}</span>
					<span className="arena-setup arena-setup-quiet">{build.rotation || 'default rotation'}</span>
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
