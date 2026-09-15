import { SITE_BASE, SITE_REPO_URL } from '../core/constants/other';
import { Entry, sections, Source } from './entries';
// Every pull request merged into master, kept current by the Update Changelog workflow.
import merged from './merged.json';

type MergedPullRequest = {
	number: number;
	title: string;
	mergedAt: string;
	url: string;
};
const mergedPullRequests = merged as Array<MergedPullRequest>;

const prLink = (pr: number) => (
	<a className="changelog-pr" href={`${SITE_REPO_URL}/pull/${pr}`} target="_blank" rel="noreferrer">
		#{pr}
	</a>
);

const sourceLink = (source: Source) => (
	<li>
		<a href={source.url} target="_blank" rel="noreferrer">
			{source.label}
		</a>
	</li>
);

const mergedRow = (pull: MergedPullRequest) => (
	<li className="changelog-merged-row">
		<a className="changelog-pr" href={pull.url} target="_blank" rel="noreferrer">
			#{pull.number}
		</a>
		<span className="changelog-merged-title">{pull.title}</span>
		<time className="changelog-merged-date" dateTime={pull.mergedAt}>
			{pull.mergedAt}
		</time>
	</li>
);

const entryBlock = (entry: Entry) => (
	<article className="changelog-entry">
		<h3 className="changelog-entry-title">
			{entry.title}
			<span className="changelog-entry-prs">{entry.prs.map(prLink)}</span>
		</h3>
		<dl className="changelog-entry-body">
			<dt>What changed</dt>
			<dd>{entry.changed}</dd>
			<dt>Effect on the sim</dt>
			<dd>{entry.effect}</dd>
			{entry.sources && (
				<>
					<dt>Where it came from</dt>
					<dd>
						<ul className="changelog-sources">{entry.sources.map(sourceLink)}</ul>
					</dd>
				</>
			)}
		</dl>
	</article>
);

// The page's index.html is generated from the shared template like every other page, so
// the markup lives here rather than in a hand-written file the build would overwrite.
export class ChangelogPage {
	constructor(parentElem: HTMLElement) {
		parentElem.appendChild(
			<div id="changelog-page">
				<header className="changelog-header">
					<div className="container changelog-header-container">
						<a href={SITE_BASE} className="changelog-home-link">
							<img className="forever-logo" src={`${SITE_BASE}assets/img/forever_logo.png`} alt="World of Warcraft: Forever" />
						</a>
						<div className="changelog-title-block">
							<h1 className="changelog-title">What changed for Forever</h1>
							<p className="changelog-subtitle">
								Every way this simulator differs from the Classic Era sim it was forked from, what each change does to the numbers, and where
								the Forever information came from.
							</p>
						</div>
					</div>
				</header>
				<main className="container changelog-content">
					<section className="content-block changelog-provenance">
						<div className="content-block-header">
							<h2 className="content-block-title">Read this first</h2>
						</div>
						<div className="content-block-body">
							<p>
								<strong>Sources.</strong> Forever is not out. What is modelled here was read from Blizzard's BlizzCon 2026 announcements and
								panel, Wowhead's Forever guides, the community talent calculators rebuilt from the demo's tooltips, and reports from people who
								played the demo. Each entry below links the source it came from; where the source was a demo tooltip the sim carries a{' '}
								<code>TODO</code> for the beta pass, all of them listed in{' '}
								<a href={`${SITE_REPO_URL}/blob/master/docs/forever_beta_checklist.md`} target="_blank" rel="noreferrer">
									the beta checklist
								</a>
								. The beta opens on 17 September; expect numbers to move.
							</p>
							<p>
								<strong>Where to read more.</strong> Each change links its pull request, which carries the full reasoning and the before and
								after numbers. The one-line rules are collected in{' '}
								<a href={`${SITE_REPO_URL}/blob/master/docs/forever_rules.md`} target="_blank" rel="noreferrer">
									the rules sheet
								</a>
								.
							</p>
						</div>
					</section>
					{sections.map(section => (
						<section className="content-block changelog-section">
							<div className="content-block-header">
								<h2 className="content-block-title">{section.title}</h2>
							</div>
							<div className="content-block-body">
								<p className="changelog-section-intro">{section.intro}</p>
								{section.entries.map(entryBlock)}
							</div>
						</section>
					))}
					<section className="content-block changelog-section changelog-merged">
						<div className="content-block-header">
							<h2 className="content-block-title">Every merged pull request ({mergedPullRequests.length})</h2>
						</div>
						<div className="content-block-body">
							<p className="changelog-section-intro">
								The sections above are written by hand and group the work by what it did. This list is the raw record, newest first, refreshed
								by the build every time a pull request is merged.
							</p>
							<ol className="changelog-merged-list">{mergedPullRequests.map(mergedRow)}</ol>
						</div>
					</section>
				</main>
			</div>,
		);
	}
}
