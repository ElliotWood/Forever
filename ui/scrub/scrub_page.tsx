import { SITE_REPO_URL } from '../core/constants/other';
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

export class ScrubPage {
	private readonly result: HTMLElement;

	constructor(parent: HTMLElement) {
		parent.appendChild(
			<div className="container scrub-container">
				<h1>Send the beta's own numbers</h1>
				<p className="scrub-lede">
					This sim reads Blizzard's data tables, which say what an ability is <em>meant</em> to do. Two files on your machine say things the tables
					cannot, and between them they close the two largest gaps here. Send either, or both.
				</p>

				<div className="scrub-file">
					<h2>
						<code>DBCache.bin</code> &mdash; what Blizzard changed after the build shipped
					</h2>
					<p>
						Blizzard tunes values <em>after</em> a build ships. Those hotfixes never reach a datamining site; they exist only in the cache your
						client downloads them into, so a value this sim reads can be stale the moment it is tuned and there is no way to know from the outside.
					</p>
					<p>
						In <code>World of Warcraft\_classic_beta_\Cache\ADB\enUS\</code>. It carries no character name, account, realm or Battle.net tag &mdash;
						it is NPC dialogue, item names and tuning rows &mdash; so it can be attached as it is.
					</p>
					<a className="scrub-button" href={ISSUE_URL} target="_blank" rel="noreferrer">
						Attach it to an issue
					</a>
				</div>

				<div className="scrub-file">
					<h2>
						<code>DamageMeter.bin</code> &mdash; what the server actually paid out
					</h2>
					<p>
						Every number on this site is unconfirmed against a running game, and Forever blocks addons from reading damage, so the client's own
						meter is the only measurement there will be. This file holds it: per ability, per fight, how many hits and how much damage. It is the
						only thing that can tell us a coefficient here is wrong rather than merely unverified.
					</p>
					<p className="scrub-warn">
						<strong>It also holds character names</strong> &mdash; yours, and everyone you grouped with, who did not agree to anything. So it needs
						the names taken out first, which is what this page is for.
					</p>
					<p>
						In <code>World of Warcraft\_classic_beta_\Cache\</code>. Drop it below. It is read in your browser and <strong>never uploaded</strong>;
						the names are gone before anything leaves your machine.
					</p>

					<label className="scrub-drop">
						<input type="file" onchange={(event: Event) => this.handle((event.target as HTMLInputElement).files?.[0])} />
						<span>Choose DamageMeter.bin</span>
					</label>
				</div>

				<div className="scrub-result" />
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
				this.show(<p className="scrub-warn">No damage meter records in that file. Is it Cache/DamageMeter.bin?</p>);
				return;
			}
			// Summarised from the scrubbed copy, not the original. Showing the real names back
			// would say "24 names removed" above a list of those 24 names, which reads as though
			// the scrub had not worked.
			this.show(this.report(file.name, scrubbed, readRecords(scrubbed), namesRemoved));
		} catch (error) {
			this.show(<p className="scrub-warn">Could not read that file: {String(error)}</p>);
		}
	}

	private show(node: Element) {
		this.result.replaceChildren(node);
	}

	private report(name: string, scrubbed: Uint8Array, records: ActorRecord[], namesRemoved: number): Element {
		const { byActor, bySpell } = totals(records);
		const url = URL.createObjectURL(new Blob([scrubbed as unknown as BlobPart], { type: 'application/octet-stream' }));
		return (
			<div>
				<h2>Ready to send</h2>
				<p>
					{String(namesRemoved)} name{namesRemoved === 1 ? '' : 's'} removed from {records.length} records. The file is the same size it was and every
					damage number is untouched. Summarised below are the {records.filter(plausible).length} records whose numbers read cleanly &mdash; the rest
					are scrubbed and sent too, they just are not worth showing back to you:
				</p>
				<ul className="scrub-summary">
					{[...byActor.entries()].map(([actor, row]) => (
						<li>
							<span className="scrub-actor">{actor}</span>
							<span>{row.className}</span>
							<span>{row.damage.toLocaleString()} damage</span>
							<span>{row.hits.toLocaleString()} hits</span>
						</li>
					))}
				</ul>
				<ul className="scrub-summary">
					{[...bySpell.entries()]
						.sort((a, b) => b[1].damage - a[1].damage)
						.map(([spellId, row]) => (
							<li>
								<span className="scrub-actor">spell {String(spellId)}</span>
								<span>{row.damage.toLocaleString()} damage</span>
								<span>{row.hits.toLocaleString()} hits</span>
							</li>
						))}
				</ul>
				<p className="scrub-actions">
					<a className="scrub-button" href={url} download={`scrubbed-${name}`}>
						Download the scrubbed file
					</a>
					<a className="scrub-button scrub-button-quiet" href={ISSUE_URL} target="_blank" rel="noreferrer">
						Then attach it to an issue
					</a>
				</p>
			</div>
		);
	}
}
