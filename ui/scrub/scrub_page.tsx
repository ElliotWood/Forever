import { SITE_BASE, SITE_REPO_URL } from '../core/constants/other';
import { ActorRecord, plausible, readRecords, scrub } from './scrub';

const ISSUE_URL = `${SITE_REPO_URL}/issues/new?template=hotfix_cache.md`;

const totals = (records: ActorRecord[]) => {
	const byActor = new Map<string, { className: string; damage: number; hits: number }>();
	const bySpell = new Map<number, { damage: number; hits: number }>();
	for (const record of records.filter(plausible)) {
		const actor = byActor.get(record.name) ?? { className: record.className, damage: 0, hits: 0 };
		actor.damage += record.damage;
		actor.hits += record.hits;
		byActor.set(record.name, actor);

		const spell = bySpell.get(record.spellId) ?? { damage: 0, hits: 0 };
		spell.damage += record.damage;
		spell.hits += record.hits;
		bySpell.set(record.spellId, spell);
	}
	return { byActor, bySpell };
};

const summaryRow = (label: string, cells: string[]) => (
	<li className="scrub-row">
		<span className="scrub-row-label">{label}</span>
		{cells.map(cell => (
			<span className="scrub-row-cell">{cell}</span>
		))}
	</li>
);

export class ScrubPage {
	private readonly result: HTMLElement;

	constructor(parent: HTMLElement) {
		parent.appendChild(
			<div id="scrub-page">
				<header className="scrub-header">
					<div className="container scrub-header-container">
						<a href={SITE_BASE} className="scrub-home-link">
							<img className="forever-logo" src={`${SITE_BASE}assets/img/forever_logo.png`} alt="World of Warcraft: Forever" />
						</a>
						<div className="scrub-title-block">
							<h1 className="scrub-title">Send the beta&apos;s own numbers</h1>
							<p className="scrub-subtitle">
								This sim reads Blizzard&apos;s data tables, which say what an ability is <em>meant</em> to do. Two files on your machine say
								things the tables cannot. Send either, or both.
							</p>
						</div>
					</div>
				</header>

				<main className="container scrub-content">
					<div className="scrub-files">
						<section className="scrub-file">
							<h2 className="scrub-file-title">
								<code>DBCache.bin</code>
							</h2>
							<p className="scrub-file-sub">What Blizzard changed after the build shipped</p>
							<p>
								Hotfixes never reach a datamining site. They exist only in the cache your client downloads them into, so a value this sim reads
								can be stale the moment it is tuned, and there is no way to know from outside.
							</p>
							<p className="scrub-path">
								<code>World of Warcraft\_classic_beta_\Cache\ADB\enUS\</code>
							</p>
							<p className="scrub-file-note">
								Carries no character name, account, realm or Battle.net tag. It is NPC dialogue, item names and tuning rows, so it can be
								attached as it is.
							</p>
							<a className="scrub-button" href={ISSUE_URL} target="_blank" rel="noreferrer">
								<i className="fab fa-github" />
								<span>Attach it to an issue</span>
							</a>
						</section>

						<section className="scrub-file">
							<h2 className="scrub-file-title">
								<code>DamageMeter.bin</code>
							</h2>
							<p className="scrub-file-sub">What the server actually paid out</p>
							<p>
								Nothing on this site has been checked against a running game, and Forever blocks addons from reading damage, so the
								client&apos;s own meter is the only measurement there will be. It is the only thing that can show a number here is wrong rather
								than merely unverified.
							</p>
							<p className="scrub-path">
								<code>World of Warcraft\_classic_beta_\Cache\</code>
							</p>
							<p className="scrub-file-note scrub-warn">
								It also holds character names, yours and everyone you grouped with, who did not agree to anything. So the names come out first,
								below.
							</p>
						</section>
					</div>

					<section className="scrub-drop-section">
						<h2 className="scrub-file-title">Take the names out</h2>
						<p>
							The file is read in your browser and <strong>never uploaded</strong>. The names are gone before anything leaves your machine, and
							you get a cleaned copy back to attach.
						</p>
						<label className="scrub-drop">
							<input type="file" onchange={(event: Event) => this.handle((event.target as HTMLInputElement).files?.[0])} />
							<i className="fas fa-upload" />
							<span className="scrub-drop-title">Choose DamageMeter.bin</span>
							<span className="scrub-drop-hint">or drag it here</span>
						</label>
						<div className="scrub-result" />
					</section>
				</main>
			</div>,
		);
		this.result = parent.querySelector('.scrub-result') as HTMLElement;

		const drop = parent.querySelector('.scrub-drop') as HTMLElement;
		drop.addEventListener('dragover', event => {
			event.preventDefault();
			drop.classList.add('scrub-drop-over');
		});
		drop.addEventListener('dragleave', () => drop.classList.remove('scrub-drop-over'));
		drop.addEventListener('drop', event => {
			event.preventDefault();
			drop.classList.remove('scrub-drop-over');
			this.handle(event.dataTransfer?.files?.[0]);
		});
	}

	private async handle(file: File | undefined) {
		if (!file) return;
		this.result.replaceChildren();
		try {
			const bytes = new Uint8Array(await file.arrayBuffer());
			const { scrubbed, records, namesRemoved } = scrub(bytes);
			if (!records.length) {
				this.show(<p className="scrub-warn scrub-message">No damage meter records in that file. Is it Cache/DamageMeter.bin?</p>);
				return;
			}
			// Summarised from the scrubbed copy, not the original. Showing the real names back
			// would say "24 names removed" above a list of those 24 names, which reads as though
			// the scrub had not worked.
			this.show(this.report(file.name, scrubbed, readRecords(scrubbed), namesRemoved));
		} catch (error) {
			this.show(<p className="scrub-warn scrub-message">Could not read that file: {String(error)}</p>);
		}
	}

	private show(node: Element) {
		this.result.replaceChildren(node);
	}

	private report(name: string, scrubbed: Uint8Array, records: ActorRecord[], namesRemoved: number): Element {
		const { byActor, bySpell } = totals(records);
		const shown = records.filter(plausible).length;
		const url = URL.createObjectURL(new Blob([scrubbed as unknown as BlobPart], { type: 'application/octet-stream' }));
		return (
			<div className="scrub-report">
				<p className="scrub-report-head">
					<strong>
						{String(namesRemoved)} name{namesRemoved === 1 ? '' : 's'} removed
					</strong>{' '}
					from {String(records.length)} records. Same size, every damage number untouched.
				</p>
				<p className="scrub-actions">
					<a className="scrub-button" href={url} download={`scrubbed-${name}`}>
						<i className="fas fa-download" />
						<span>Download the scrubbed file</span>
					</a>
					<a className="scrub-button scrub-button-quiet" href={ISSUE_URL} target="_blank" rel="noreferrer">
						<i className="fab fa-github" />
						<span>Then attach it to an issue</span>
					</a>
				</p>
				<p className="scrub-report-note">
					The {String(shown)} records whose numbers read cleanly are below, exactly as they will arrive. The rest are scrubbed and sent too, they are
					just not worth showing back to you.
				</p>
				<ul className="scrub-summary">
					{[...byActor.entries()].map(([actor, row]) =>
						summaryRow(actor, [row.className, `${row.damage.toLocaleString()} damage`, `${row.hits.toLocaleString()} hits`]),
					)}
				</ul>
				<ul className="scrub-summary">
					{[...bySpell.entries()]
						.sort((a, b) => b[1].damage - a[1].damage)
						.map(([spellId, row]) =>
							summaryRow(`spell ${spellId}`, [`${row.damage.toLocaleString()} damage`, `${row.hits.toLocaleString()} hits`]),
						)}
				</ul>
			</div>
		);
	}
}
