// Every ability the sim registers, and what is actually known about its numbers.
//
// The manifest in ui/core/spells has always held this, and sim/spell_sources_test.go has
// always enforced it, but the only way to read it was to clone the repository. On screen an
// ability nobody has settled looked exactly like one lifted from the client. The icons now
// carry a dot for the unsettled ones; this page is where that dot leads.

import { SITE_BASE, SITE_REPO_URL } from '../core/constants/other';
import { ActionId } from '../core/proto_utils/action_id';
import { allSpellSources, SpellSource } from '../core/spells/index';

type Tier = 'measured' | 'forever' | 'classic' | 'assumed';

const TIERS: Array<{ key: Tier; label: string; blurb: string }> = [
	{ key: 'measured', label: 'Seen in game', blurb: "Watched happen on a running server, through the client's own damage meter." },
	{ key: 'forever', label: 'Read from the client', blurb: "Taken out of the beta client's own data tables." },
	{ key: 'classic', label: 'Unchanged from Classic', blurb: 'Identical to Classic Era, so the long-established value stands.' },
	{ key: 'assumed', label: 'Still a guess', blurb: 'At least one number here is unconfirmed, and the row says which.' },
];

const tierOf = (s: SpellSource): Tier => s.source as Tier;

/** A row's searchable text, so filtering never has to walk the DOM. */
const haystack = (id: number, s: SpellSource) => `${id} ${s.ability} ${s.file} ${s.note ?? ''} ${(s.assumptions ?? []).join(' ')}`.toLowerCase();

// 996 rows means 996 icon lookups, and an icon the bundled database has never heard of goes
// out to Wowhead for it. Firing those on load would be a thousand requests for the forty
// rows anyone can actually see, so icons fill when they scroll into view.
const iconsInView = new IntersectionObserver(
	(entries, self) => {
		for (const entry of entries) {
			if (!entry.isIntersecting) continue;
			self.unobserve(entry.target);
			const elem = entry.target as HTMLAnchorElement;
			ActionId.fromSpellId(Number(elem.dataset.spellId)).fillAndSet(elem, false, true);
		}
	},
	{ rootMargin: '200px' },
);

export class EvidencePage {
	private readonly rows: Array<{ elem: HTMLElement; tiers: Set<Tier>; text: string }> = [];
	private readonly count: HTMLElement;
	private query = '';
	private tier: Tier | 'all' = 'all';

	constructor(parent: HTMLElement) {
		const entries = allSpellSources()
			.filter(([, s]) => s.source !== 'unreviewed')
			.sort(([, a], [, b]) => a.ability.localeCompare(b.ability));

		const tally = new Map<Tier, number>();
		for (const [, s] of entries) tally.set(tierOf(s), (tally.get(tierOf(s)) ?? 0) + 1);
		const measured = entries.filter(([, s]) => !!s.measured).length;

		parent.appendChild(
			<div id="evidence-page">
				<header className="evidence-header">
					<div className="container evidence-header-container">
						<a href={SITE_BASE} className="evidence-home-link">
							<img className="forever-logo" src={`${SITE_BASE}assets/img/forever_logo.png`} alt="World of Warcraft: Forever" />
						</a>
						<div className="evidence-title-block">
							<h1 className="evidence-title">Where every number came from</h1>
							<p className="evidence-subtitle">
								Every ability this sim runs, and how much is actually known about it. {String(entries.length)} of them: {String(measured)} have
								been watched happen on a running server, {String(tally.get('assumed') ?? 0)} still carry a guess, and the rest are read straight
								out of the beta client.
							</p>
						</div>
					</div>
				</header>

				<main className="container evidence-content">
					<div className="evidence-controls">
						<input
							className="evidence-search"
							type="search"
							placeholder="Filter by name, spell id, file or reason"
							oninput={(e: Event) => this.setQuery((e.target as HTMLInputElement).value)}
						/>
						<div className="evidence-tiers">
							{this.tierButton('all', 'Everything', entries.length)}
							{TIERS.map(t => this.tierButton(t.key, t.label, t.key === 'measured' ? measured : tally.get(t.key) ?? 0))}
						</div>
					</div>

					<ul className="evidence-legend">
						{TIERS.map(t => (
							<li>
								<span className={`evidence-chip evidence-chip-${t.key}`}>{t.label}</span>
								<span>{t.blurb}</span>
							</li>
						))}
					</ul>

					<p className="evidence-count" />

					<ul className="evidence-rows">{entries.map(([id, s]) => this.row(id, s))}</ul>

					<p className="evidence-foot">
						Kept honest by <code>sim/spell_sources_test.go</code>, which will not let an ability be registered without saying where its numbers came
						from. The entries themselves live in{' '}
						<a href={`${SITE_REPO_URL}/tree/master/ui/core/spells`} target="_blank" rel="noreferrer">
							ui/core/spells
						</a>
						. If you can move a row up a tier, <a href={`${SITE_BASE}scrub/`}>send the beta&apos;s own numbers</a>.
					</p>
				</main>
			</div>,
		);

		this.count = parent.querySelector('.evidence-count') as HTMLElement;
		this.apply();
	}

	private tierButton(tier: Tier | 'all', label: string, count: number): Element {
		const button = (
			<button className={`evidence-tier${tier === 'all' ? ' evidence-tier-on' : ''}`} type="button">
				<span>{label}</span>
				<span className="evidence-tier-count">{String(count)}</span>
			</button>
		) as HTMLButtonElement;

		button.addEventListener('click', () => {
			this.tier = tier;
			for (const other of document.querySelectorAll('.evidence-tier')) other.classList.remove('evidence-tier-on');
			button.classList.add('evidence-tier-on');
			this.apply();
		});
		return button;
	}

	private row(id: number, s: SpellSource): Element {
		const tier = tierOf(s);
		const tiers = new Set<Tier>([tier]);
		if (s.measured) tiers.add('measured');

		// The href needs no lookup, so it works before the icon has loaded and still works for
		// the Forever-only ids Wowhead will never have a page for.
		const icon = (<a className="evidence-icon" target="_blank" rel="noreferrer" dataset={{ spellId: String(id) }} />) as HTMLAnchorElement;
		ActionId.fromSpellId(id).setWowheadHref(icon);
		iconsInView.observe(icon);

		const elem = (
			<li className="evidence-row" dataset={{ spellSource: s.source }}>
				{icon}
				<div className="evidence-name">
					<span className="evidence-ability">{s.ability}</span>
					<span className="evidence-id">{String(id)}</span>
				</div>
				<div className="evidence-tags">
					<span className={`evidence-chip evidence-chip-${tier}`}>{TIERS.find(t => t.key === tier)!.label}</span>
					{s.measured ? <span className="evidence-chip evidence-chip-measured">Seen {s.measured.date}</span> : <></>}
				</div>
				<div className="evidence-detail">
					{s.tooltip ? <p className="evidence-tooltip">{s.tooltip}</p> : <></>}
					{s.measured ? (
						<p className="evidence-measured">
							Client says <strong>{s.measured.client}</strong>; the game&apos;s own meter recorded a largest hit of{' '}
							<strong>{String(s.measured.biggest)}</strong>, via {s.measured.how}.
						</p>
					) : (
						<></>
					)}
					{(s.assumptions ?? []).map(a => (
						<p className="evidence-assumption">{a}</p>
					))}
					{s.note ? <p className="evidence-note">{s.note}</p> : <></>}
					<p className="evidence-file">
						<code>{s.file}</code>
					</p>
				</div>
			</li>
		) as HTMLElement;

		this.rows.push({ elem, tiers, text: haystack(id, s) });
		return elem;
	}

	private setQuery(value: string) {
		this.query = value.trim().toLowerCase();
		this.apply();
	}

	private apply() {
		let shown = 0;
		for (const row of this.rows) {
			const visible = (this.tier === 'all' || row.tiers.has(this.tier)) && (!this.query || row.text.includes(this.query));
			row.elem.classList.toggle('evidence-hidden', !visible);
			shown += visible ? 1 : 0;
		}
		this.count.textContent = `Showing ${shown} of ${this.rows.length}`;
	}
}
