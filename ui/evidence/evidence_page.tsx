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

// What would actually move a number, worst first, with what to send for each.
//
// Hand-written on purpose. The manifest knows which abilities are unsettled, but not which
// unsettled thing matters - downranking is one line in no JSON file and moves every caster
// on the site, while a hunter pet's attack speed is a rounding error. A list generated from
// the manifest would rank those the same and quietly waste the first person who offers to
// help. The counts inside it come from the manifest, so those cannot go stale.
type Need = {
	title: string;
	why: string;
	/** What to send, in the words of someone who has the game open. */
	send: string;
	/** Pre-fills the search box below, so the row in question is one click away. */
	find?: string;
};

const NEEDS: Array<Need> = [
	{
		title: 'Downranking: does a low rank still hit for full?',
		why: 'The client carries the full coefficient on low ranks where Classic Era carried a reduced one. Read as written, a rank 4 Lightning Bolt does most of a rank 10 for a quarter of the mana, which would rewrite every caster rotation on this site. No table anywhere says whether Forever kept the penalty.',
		send: 'A DamageMeter.bin from a session where you deliberately spammed a low rank of a direct damage spell - rank 1-4 Lightning Bolt, Fireball, Shadow Bolt. Twenty casts is plenty. The meter records the biggest hit, and that alone settles it.',
		find: 'Lightning Bolt',
	},
	{
		title: 'Hunter, nearly everything',
		why: 'Hunter carries more guesses than the rest of the game put together. Volley, Serpent Sting, Arcane Shot and the pet abilities all run on numbers the client does not settle, and Summon Hawk models one hawk where the tooltip can be read as allowing two.',
		send: 'A DamageMeter.bin from any hunter at any level, or screenshots of those tooltips out of your spellbook. Either one. The tooltips are worth as much as the damage here, because half of what is wrong is about what a number applies to rather than what it is.',
		find: 'hunter',
	},
	{
		title: 'Sky Elf abilities that appear in no file at all',
		why: 'Three spell ids turn up in the beta client with no home: 1259231 Infusion of Wind, 1259652 Shock, 1248802 Wind Spike. They are not registered anywhere in this sim because nobody knows whether they are racials, a quest reward or cut content.',
		send: 'A screenshot of a Sky Elf spellbook, or of the racials pane on the character screen. One picture ends this.',
	},
	{
		title: 'Hotfixes, from anyone, any day',
		why: 'Blizzard tunes after the build ships and none of it reaches a datamining site. It exists only in the cache your own client downloads it into, so a number here can go stale with nothing to indicate it has.',
		send: 'DBCache.bin, as it is. It carries no character name, account, realm or Battle.net tag, so there is nothing to strip. Sending the same file again next week is useful - it is the change that matters.',
	},
	{
		title: 'Any tooltip that disagrees with this sim',
		why: 'Five bugs so far passed a value check and were wrong about what the value applied to. Improved Seals scaled half of what it should while every number matched. A rank curve says what a talent’s numbers are, never what they do.',
		send: 'A screenshot, cropped to the tooltip. If it contradicts what is written on a row below, that row is wrong and it takes one picture to prove it.',
	},
];

const tierOf = (s: SpellSource): Tier => s.source as Tier;

/** A row's searchable text, so filtering never has to walk the DOM. */
const haystack = (id: number, s: SpellSource) => `${id} ${s.ability} ${s.file} ${s.note ?? ''} ${(s.assumptions ?? []).join(' ')}`.toLowerCase();

// 996 rows means 996 icon lookups, and an icon the bundled database has never heard of goes
// out to Wowhead for it. Firing those on load would be a thousand requests for the forty
// rows anyone can actually see, so icons fill when they scroll into view.
//
// The href waits for the same moment, which is not an optimisation of the same kind: the
// site loads Wowhead's tooltip script with colorLinks on, and that script walks every link
// to a Wowhead URL and fetches it. A page that hands it 996 links at once is asking for 996
// requests before anyone has scrolled, and most of them 404 because Wowhead has never heard
// of an ability Forever invented.
const iconsInView = new IntersectionObserver(
	(entries, self) => {
		for (const entry of entries) {
			if (!entry.isIntersecting) continue;
			self.unobserve(entry.target);
			const elem = entry.target as HTMLAnchorElement;
			const actionId = ActionId.fromSpellId(Number(elem.dataset.spellId));
			// The link needs no lookup, so it must not wait behind one. Going through
			// fillAndSet held the href until the icon request came back, and for an id Wowhead
			// has never heard of that request is a slow 404 - which left rows with a working
			// icon lookup pending and no link at all in the meantime.
			actionId.setWowheadHref(elem);
			actionId.fill().then(filled => filled.setBackground(elem));
		}
	},
	{ rootMargin: '200px' },
);

export class EvidencePage {
	private readonly rows: Array<{ elem: HTMLElement; icon: HTMLAnchorElement; tiers: Set<Tier>; text: string }> = [];
	private readonly count: HTMLElement;
	private readonly search: HTMLInputElement;
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
					<section className="evidence-wanted" id="most-wanted">
						<h2 className="evidence-wanted-title">Most wanted</h2>
						<p className="evidence-wanted-lede">
							Worst first. Each of these can be closed by one person with the game open, and the top one moves every caster on the site.{' '}
							<a className="evidence-wanted-link" href={`${SITE_BASE}scrub/`}>
								Send a file or a screenshot
							</a>{' '}
							&mdash; no account, no form.
						</p>
						<ol className="evidence-wanted-list">
							{NEEDS.map((need, i) => (
								<li className="evidence-want">
									<h3 className="evidence-want-title">
										<span className="evidence-want-rank">{String(i + 1)}</span>
										{need.title}
									</h3>
									<p className="evidence-want-why">{need.why}</p>
									<p className="evidence-want-send">
										<strong>What settles it:</strong> {need.send}
									</p>
									{need.find ? this.findButton(need.find) : <></>}
								</li>
							))}
						</ol>
					</section>

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
		this.search = parent.querySelector('.evidence-search') as HTMLInputElement;
		// Observed only once the tree is in the document. Observing while the rows were still
		// being built produced one callback per icon saying "not intersecting" - correct, since
		// a detached element intersects nothing - and no second one when they were attached, so
		// every icon below the first screenful stayed blank and unlinked.
		for (const row of this.rows) iconsInView.observe(row.icon);
		this.apply();
	}

	/** Drops a most-wanted item straight into the list below it. */
	private findButton(find: string): Element {
		const button = (
			<button className="evidence-find" type="button">
				<i className="fas fa-search" />
				<span>Show these rows</span>
			</button>
		) as HTMLButtonElement;

		button.addEventListener('click', () => {
			this.search.value = find;
			this.setQuery(find);
			this.search.scrollIntoView({ behavior: 'smooth', block: 'start' });
		});
		return button;
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

		const icon = (<a className="evidence-icon" target="_blank" rel="noreferrer" dataset={{ spellId: String(id) }} />) as HTMLAnchorElement;

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

		this.rows.push({ elem, icon, tiers, text: haystack(id, s) });
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
