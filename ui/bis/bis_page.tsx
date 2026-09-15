import { ref } from 'tsx-vanilla';

import { SITE_BASE } from '../core/constants/other';
import { setItemQualityCssClass } from '../core/css_utils';
import { Class, ItemRandomSuffix, Spec } from '../core/proto/common';
import { UIItem as Item } from '../core/proto/ui';
import { ActionId } from '../core/proto_utils/action_id';
import { Database } from '../core/proto_utils/database';
import { slotNames } from '../core/proto_utils/names';
import { Stats, UnitStat } from '../core/proto_utils/stats';
import { classNames, specNames, specToClass } from '../core/proto_utils/utils';
import { RankedItem, rankGear, SlotRanking } from './gear';
import { bisSpecs, loadSpecConfig } from './specs';

export const specSlug = (spec: Spec): string => specNames[spec].toLowerCase().replace(/[^a-z]+/g, '-');

const iconUrl = (item: Item): string => `https://wow.zamimg.com/images/wow/icons/large/${item.icon}.jpg`;

// Rounding first keeps an item worth nothing but its unique penalty off the page as -0.0.
const formatEP = (ep: number): string => (Math.round(ep * 10) / 10).toFixed(1);

const loadingMessage = (text: string) => (
	<p className="bis-loading">
		<i className="fa fa-spinner fa-spin" />
		{text}
	</p>
);

// The weights a spec ranks by, spelled out so a reader can see what produced the order.
// A weight can read back undefined rather than zero; see weightFor in gear.ts.
const weightsSummary = (epWeights: Stats, playerClass: Class): string =>
	UnitStat.getAll()
		.map(unitStat => ({ name: unitStat.getName(playerClass), weight: epWeights.getUnitStat(unitStat) || 0 }))
		.filter(({ name, weight }) => !!name && weight !== 0)
		.map(({ name, weight }) => `${name} ${weight.toFixed(2)}`)
		.join(' · ');

const itemCell = (ranked: RankedItem, rank: number) => {
	const nameRef = ref<HTMLAnchorElement>();
	const cell = (
		<td className="bis-item">
			<img className="bis-item-icon" src={iconUrl(ranked.item)} alt="" loading="lazy" />
			<a ref={nameRef} className="bis-item-name" href={ActionId.makeItemUrl(ranked.item.id)} target="_blank" rel="noreferrer">
				{ranked.item.name}
			</a>
			{rank > 1 && <span className="bis-item-rank">#{rank}</span>}
		</td>
	);
	setItemQualityCssClass(nameRef.value!, ranked.item.quality);
	return cell;
};

// One row for the slot's pick, plus a hidden row per item it beat. The banding is set here
// rather than with nth-child, which would count the hidden rows and stripe by accident.
const slotRows = (ranking: SlotRanking, banded: boolean): Array<HTMLElement> => {
	const best = ranking.ranked[0];
	const runnersUp = ranking.ranked.slice(1);
	const runnerRows = runnersUp.map((ranked, index) => (
		<tr className={`bis-runner-up ${banded ? 'bis-banded' : ''}`} hidden>
			<th scope="row" />
			{itemCell(ranked, index + 2)}
			<td className="bis-ilvl">{String(ranked.item.ilvl)}</td>
			<td className="bis-ep">{formatEP(ranked.ep)}</td>
		</tr>
	)) as Array<HTMLElement>;

	const toggleRef = ref<HTMLButtonElement>();
	const row = (
		<tr className={`bis-pick ${banded ? 'bis-banded' : ''}`}>
			<th scope="row" className="bis-slot">
				{slotNames.get(ranking.slot)}
			</th>
			{best ? itemCell(best, 1) : <td className="bis-item bis-item-empty">{ranking.note || 'Nothing in the launch pool fits this slot.'}</td>}
			<td className="bis-ilvl">{best ? String(best.item.ilvl) : '—'}</td>
			<td className="bis-ep">
				<span className="bis-ep-value">{best ? formatEP(best.ep) : '—'}</span>
				{!!runnersUp.length && (
					<button ref={toggleRef} className="bis-toggle" type="button" aria-expanded="false" title={`Show the ${runnersUp.length} items this beat`}>
						<i className="fas fa-chevron-down" />
					</button>
				)}
			</td>
		</tr>
	) as HTMLElement;

	if (toggleRef.value) {
		const toggle = toggleRef.value;
		toggle.addEventListener('click', () => {
			const showing = toggle.getAttribute('aria-expanded') === 'true';
			toggle.setAttribute('aria-expanded', String(!showing));
			runnerRows.forEach(runnerRow => (runnerRow.hidden = showing));
		});
	}

	return [row, ...runnerRows];
};

export class BisPage {
	private readonly resultsElem: HTMLElement;
	private readonly selectElem: HTMLSelectElement;
	private items: Array<Item> = [];
	private randomSuffix: (id: number) => ItemRandomSuffix | undefined = () => undefined;
	private spec: Spec;

	// The page's index.html is generated from the shared template like every other page,
	// so the markup lives here rather than in a hand-written file the build would overwrite.
	constructor(parentElem: HTMLElement) {
		this.spec = this.specFromHash() ?? bisSpecs[0];
		this.resultsElem = (<div className="bis-results">{loadingMessage('Loading the item database…')}</div>) as HTMLElement;

		const selectRef = ref<HTMLSelectElement>();
		const byClass = new Map<Class, Array<Spec>>();
		bisSpecs.forEach(spec => byClass.set(specToClass[spec], (byClass.get(specToClass[spec]) ?? []).concat(spec)));

		parentElem.appendChild(
			<div id="bis-page">
				<header className="bis-header">
					<div className="container bis-header-container">
						<a href={SITE_BASE} className="bis-home-link">
							<img className="forever-logo" src={`${SITE_BASE}assets/img/forever_logo.png`} alt="World of Warcraft: Forever" />
						</a>
						<div className="bis-title-block">
							<h1 className="bis-title">Best in Slot</h1>
							<p className="bis-subtitle">The highest-EP item in every slot, for each spec the sim has launched.</p>
						</div>
					</div>
				</header>
				<main className="container bis-content">
					<section className="content-block bis-provenance">
						<div className="content-block-header">
							<h2 className="content-block-title">Read this first</h2>
						</div>
						<div className="content-block-body">
							<p>
								<strong>Where the EP weights come from.</strong> Every spec ships a set of default EP weights in its sim, and its gear picker
								sorts items by them. This page uses those same weights, so a pick here is the item that sits at the top of that spec's gear
								picker for that slot. <strong>The simulator does not run on this page.</strong> Nothing here is derived live: these weights were
								worked out for each spec and then written down, so treat them as a well-kept starting point rather than an answer computed for
								your character. For weights derived from the sim against your own gear, talents and encounter, open a spec's sim and press{' '}
								<em>Stat Weights</em>.
							</p>
							<p>
								<strong>Which specs are here.</strong> Only specs the sim has launched. Restoration Druid, Restoration Shaman, Holy Paladin and
								Healing Priest are not implemented, and neither is Feral Tank Druid, so they have no EP weights and are absent rather than
								guessed at.
							</p>
							<p>
								<strong>What the item pool is.</strong> Forever launch content only: Onyxia and Molten Core, alongside the dungeon, crafted,
								reputation and PvP gear available at launch. There is no Blackwing Lair, Zul'Gurub, Ahn'Qiraj or Naxxramas gear in the database
								this ranks, so nothing from those raids can appear.
							</p>
							<p>
								<strong>What EP cannot see.</strong> An EP score reads an item's stat line and its weapon damage, and nothing else. Set bonuses,
								on-use and proc effects, weapon skill and weapon specialisation talents, resistance requirements and stat caps are all invisible
								to it. That is why the rankings favour raw stats, and why a dagger can out-rank a sword for a warrior here in a way it would not
								in a raid.
							</p>
						</div>
					</section>
					<section className="content-block bis-gear-block">
						<div className="content-block-header">
							<h2 className="content-block-title">Gear</h2>
						</div>
						<div className="content-block-body">
							<label className="bis-spec-picker">
								<span className="bis-spec-picker-label">Spec</span>
								<select ref={selectRef} className="form-select">
									{Array.from(byClass.entries()).map(([playerClass, specs]) => (
										<optgroup label={classNames[playerClass]}>
											{specs.map(spec => (
												<option value={String(spec)}>{specNames[spec]}</option>
											))}
										</optgroup>
									))}
								</select>
							</label>
							{this.resultsElem}
						</div>
					</section>
				</main>
			</div>,
		);

		this.selectElem = selectRef.value!;
		this.selectElem.value = String(this.spec);
		this.selectElem.addEventListener('change', () => this.setSpec(Number(this.selectElem.value) as Spec));

		window.addEventListener('hashchange', () => {
			const spec = this.specFromHash();
			if (spec !== undefined && spec !== this.spec) {
				this.setSpec(spec);
			}
		});

		this.load();
	}

	private specFromHash(): Spec | undefined {
		return bisSpecs.find(spec => specSlug(spec) === window.location.hash.slice(1));
	}

	private async load() {
		const db = await Database.get();
		this.items = db.getAllItems();
		this.randomSuffix = id => db.getRandomSuffixById(id);
		await this.render();
	}

	private setSpec(spec: Spec) {
		this.spec = spec;
		this.selectElem.value = String(spec);
		window.history.replaceState(null, '', `#${specSlug(spec)}`);
		this.render();
	}

	private async render() {
		const spec = this.spec;
		this.resultsElem.replaceChildren(loadingMessage(`Ranking ${specNames[spec]} gear…`));

		const config = await loadSpecConfig(spec);
		// The reader may have moved on to another spec while this one's module loaded.
		if (this.spec !== spec) {
			return;
		}

		const epWeights = config.defaults.epWeights;
		const gear = rankGear(spec, epWeights, this.items, this.randomSuffix);

		this.resultsElem.replaceChildren(
			<>
				<div className="bis-spec-header">
					<h2 className={`bis-spec-name text-${config.cssScheme}`}>{specNames[spec]}</h2>
					<p className="bis-spec-pool">Best of the {gear.poolSize.toLocaleString()} launch items this spec can wear.</p>
				</div>
				<div className="bis-table-container">
					<table className="bis-table">
						<thead>
							<tr>
								<th scope="col">Slot</th>
								<th scope="col">Item</th>
								<th scope="col" className="bis-ilvl">
									ilvl
								</th>
								<th scope="col" className="bis-ep">
									EP
								</th>
							</tr>
						</thead>
						<tbody>{gear.slots.map((ranking, index) => slotRows(ranking, index % 2 === 1)).flat()}</tbody>
					</table>
				</div>
				<p className="bis-weights">
					<span className="bis-weights-label">EP weights used</span>
					{weightsSummary(epWeights, specToClass[spec])}
				</p>
			</>,
		);
	}
}
