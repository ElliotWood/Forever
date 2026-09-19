import { SITE_BASE, SITE_REPO_URL } from '../core/constants/other';
import { ActorRecord, plausible, readRecords, scrub } from './scrub';
import { upload } from './upload';

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

const dropZone = (id: string, title: string, hint: string, onFile: (file: File | undefined) => void) => {
	const label = (
		<label className="scrub-drop" id={id}>
			<input type="file" onchange={(event: Event) => onFile((event.target as HTMLInputElement).files?.[0])} />
			<i className="fas fa-upload" />
			<span className="scrub-drop-title">{title}</span>
			<span className="scrub-drop-hint">{hint}</span>
		</label>
	) as HTMLElement;

	label.addEventListener('dragover', event => {
		event.preventDefault();
		label.classList.add('scrub-drop-over');
	});
	label.addEventListener('dragleave', () => label.classList.remove('scrub-drop-over'));
	label.addEventListener('drop', event => {
		event.preventDefault();
		label.classList.remove('scrub-drop-over');
		onFile((event as DragEvent).dataTransfer?.files?.[0]);
	});
	return label;
};

export class ScrubPage {
	private readonly hotfixResult: HTMLElement;
	private readonly meterResult: HTMLElement;
	private note = '';

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
								things the tables cannot. Drop either one below, and that is the whole job &mdash; no account, no form.
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
								Carries no character name, account, realm or Battle.net tag &mdash; it is NPC dialogue, item names and tuning rows &mdash; so it
								goes straight up.
							</p>
							{dropZone('hotfix-drop', 'Choose DBCache.bin', 'or drag it here', file => this.sendHotfix(file))}
							<div className="scrub-result scrub-result-hotfix" />
						</section>

						<section className="scrub-file">
							<h2 className="scrub-file-title">
								<code>DamageMeter.bin</code>
							</h2>
							<p className="scrub-file-sub">What the server actually paid out</p>
							<p>
								Nothing on this site has been checked against a running game, and Forever blocks addons from reading damage, so the
								client&apos;s own meter is the only measurement there will be. The client clears it between sessions, so grab it while you are
								logged in.
							</p>
							<p className="scrub-path">
								<code>World of Warcraft\_classic_beta_\Cache\</code>
							</p>
							<p className="scrub-file-note scrub-warn">
								This one holds character names, yours and everyone you grouped with. They are taken out <strong>in your browser</strong> before
								anything is sent, and you get to see what is left first.
							</p>
							{dropZone('meter-drop', 'Choose DamageMeter.bin', 'or drag it here', file => this.scrubMeter(file))}
							<div className="scrub-result scrub-result-meter" />
						</section>
					</div>

					<section className="scrub-drop-section">
						<label className="scrub-note">
							<span>Anything worth saying about it? Optional.</span>
							<input
								type="text"
								maxLength={200}
								placeholder="Fury warrior, dummy, 3 min - or a tooltip that disagrees with the sim"
								oninput={(event: Event) => (this.note = (event.target as HTMLInputElement).value)}
							/>
						</label>
						<p className="scrub-file-note">
							Files land in a private bucket, are read and thrown away, and are never committed to the repository. If you would rather send it
							yourself,{' '}
							<a href={ISSUE_URL} target="_blank" rel="noreferrer">
								open an issue
							</a>{' '}
							instead.
						</p>
					</section>
				</main>
			</div>,
		);
		this.hotfixResult = parent.querySelector('.scrub-result-hotfix') as HTMLElement;
		this.meterResult = parent.querySelector('.scrub-result-meter') as HTMLElement;
	}

	/** No names in it, so it goes as soon as it is dropped. */
	private async sendHotfix(file: File | undefined) {
		if (!file) return;
		this.hotfixResult.replaceChildren(<p className="scrub-message">Sending {file.name}...</p>);
		const bytes = new Uint8Array(await file.arrayBuffer());
		const result = await upload('dbcache', bytes, this.note);
		this.hotfixResult.replaceChildren(
			result.ok ? (
				<p className="scrub-sent">
					<i className="fas fa-check" /> Sent, thank you. Reference <code>{result.receipt.split('/').pop()!.slice(0, 8)}</code>.
				</p>
			) : (
				<p className="scrub-warn scrub-message">{result.error}</p>
			),
		);
	}

	/** Names out first, then show what is left, then send on a deliberate click. */
	private async scrubMeter(file: File | undefined) {
		if (!file) return;
		this.meterResult.replaceChildren(<p className="scrub-message">Reading {file.name}...</p>);
		try {
			const bytes = new Uint8Array(await file.arrayBuffer());
			const { scrubbed, records, namesRemoved } = scrub(bytes);
			if (!records.length) {
				this.meterResult.replaceChildren(
					<p className="scrub-warn scrub-message">No damage meter records in that file. Is it Cache/DamageMeter.bin?</p>,
				);
				return;
			}
			// Summarised from the scrubbed copy, not the original. Showing the real names back
			// would say "24 names removed" above a list of those 24 names, which reads as though
			// the scrub had not worked.
			this.meterResult.replaceChildren(this.report(file.name, scrubbed, readRecords(scrubbed), namesRemoved));
		} catch (error) {
			this.meterResult.replaceChildren(<p className="scrub-warn scrub-message">Could not read that file: {String(error)}</p>);
		}
	}

	private report(name: string, scrubbed: Uint8Array, records: ActorRecord[], namesRemoved: number): Element {
		const { byActor, bySpell } = totals(records);
		const shown = records.filter(plausible).length;
		const url = URL.createObjectURL(new Blob([scrubbed as unknown as BlobPart], { type: 'application/octet-stream' }));

		const send = (
			<button className="scrub-button" type="button">
				<i className="fas fa-paper-plane" />
				<span>Send it</span>
			</button>
		) as HTMLButtonElement;

		send.addEventListener('click', async () => {
			send.disabled = true;
			send.replaceChildren(<span>Sending...</span>);
			const result = await upload('damagemeter', scrubbed, this.note);
			send.replaceWith(
				result.ok ? (
					<p className="scrub-sent">
						<i className="fas fa-check" /> Sent, thank you. Reference <code>{result.receipt.split('/').pop()!.slice(0, 8)}</code>.
					</p>
				) : (
					<p className="scrub-warn scrub-message">{result.error}</p>
				),
			);
		});

		return (
			<div className="scrub-report">
				<p className="scrub-report-head">
					<strong>
						{String(namesRemoved)} name{namesRemoved === 1 ? '' : 's'} removed
					</strong>{' '}
					from {String(records.length)} records. Same size, every damage number untouched.
				</p>
				<p className="scrub-actions">
					{send}
					<a className="scrub-button scrub-button-quiet" href={url} download={`scrubbed-${name}`}>
						<i className="fas fa-download" />
						<span>Or download it</span>
					</a>
				</p>
				<p className="scrub-report-note">The {String(shown)} records whose numbers read cleanly are below, exactly as they would arrive.</p>
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
