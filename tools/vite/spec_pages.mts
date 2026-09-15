import fs from 'fs';
import { IncomingMessage, ServerResponse } from 'http';
import path from 'path';
import { Connect, normalizePath, PluginOption, ViteDevServer } from 'vite';

/**
 * Generates one HTML page per sim from `ui/index_template.html`, so that no per-page
 * `index.html` has to exist in the source tree. The makefile used to sed them into it, and
 * overwrote any hand-written page on the way.
 *
 * A directory `ui/<page>` is a page if it has both an `index.ts` entry point and a stylesheet
 * at `ui/scss/sims/<page>/index.scss` -- exactly the two things the template references.
 */

export type SpecPage = {
	name: string;
	title: string;
	/** '<page>/index.html', the page's location relative to the site root. */
	outPath: string;
};

const TEMPLATE_NAME = 'index_template.html';

// Pages whose directory name doesn't spell their title. Everything else falls back to
// 'WoW Forever <Dir Name> Simulator'.
// Pages kept off the homepage and out of search. The damage table is a self-check for
// the sim, not a claim about Forever, and a table built from tooltips inherits every
// tooltip that is wrong; it stays at its URL for anyone checking the sim.
const unlistedPages = new Set(['dps_rankings']);

const pageTitles: Record<string, string> = {
	bis: 'Best in Slot - Forever Sim (unofficial)',
	changelog: 'Changelog - Forever Sim (unofficial)',
	dps_rankings: 'Damage comparison - Forever Sim (unofficial)',
	stat_weights: 'Stat Weights - Forever Sim (unofficial)',
};

const titleFor = (name: string) =>
	pageTitles[name] ??
	`${name
		.split('_')
		.map(word => word.charAt(0).toUpperCase() + word.slice(1))
		.join(' ')} - Forever Sim (unofficial)`;

export function discoverSpecPages(uiRoot: string): SpecPage[] {
	const pages: SpecPage[] = [];

	for (const dir of fs.readdirSync(uiRoot, { withFileTypes: true })) {
		if (!dir.isDirectory()) continue;

		const name = dir.name;
		const hasEntry = fs.existsSync(path.join(uiRoot, name, 'index.ts'));
		const hasStyles = fs.existsSync(path.join(uiRoot, 'scss', 'sims', name, 'index.scss'));
		if (hasEntry && hasStyles) pages.push({ name, title: titleFor(name), outPath: `${name}/index.html` });
	}

	return pages.sort((a, b) => a.outPath.localeCompare(b.outPath));
}

/**
 * Fills in the template's placeholders and rewrites its depth-relative references to
 * root-absolute ones, so the page renders identically no matter where it is served from. Vite
 * resolves root-absolute specifiers against its root in both dev and build. Already-absolute
 * URLs and CDN URLs are left alone.
 */
export function renderSpecPage(template: string, page: SpecPage, siteBase: string): string {
	return (
		template
			.replaceAll('@@TITLE@@', page.title)
			.replaceAll('@@ROBOTS@@', unlistedPages.has(page.name) ? '<meta name="robots" content="noindex" />' : '')
			.replaceAll('@@SPEC@@', page.name)
			.replaceAll('@@BASE@@', siteBase)
			// '../scss/...' and '../index.ts' are relative to ui/, i.e. the vite root.
			.replace(/(["'])\.\.\//g, '$1/')
			// './index.ts' is this page's own entry point.
			.replace(/(["'])\.\//g, `$1/${page.name}/`)
	);
}

export function specPages(uiRoot: string, siteBase: string): PluginOption {
	const templatePath = path.join(uiRoot, TEMPLATE_NAME);
	const readTemplate = () => fs.readFileSync(templatePath, 'utf-8');

	/**
	 * Absolute id of a page -> the page it renders. The ids point at `ui/<page>/index.html`,
	 * which is where the pages used to be generated on disk: rollup names an HTML output after
	 * the input's path relative to the vite root, so this is what puts each page at its final
	 * URL. The files themselves never exist -- the `load` hook below renders them on demand.
	 */
	const pageIds = new Map<string, SpecPage>();

	return {
		name: 'spec-pages',
		// Beat vite:resolve and vite:load-fallback to the (nonexistent) page files.
		enforce: 'pre',

		config(_userConfig, env) {
			if (env.command !== 'build') return;

			const input: Record<string, string> = {};
			for (const page of discoverSpecPages(uiRoot)) {
				const id = normalizePath(path.join(uiRoot, page.outPath));
				pageIds.set(id, page);
				// Keyed the same way as the landing page in vite.config.mts.
				input[`${path.basename(uiRoot)}/${page.outPath}`] = id;
			}

			return { build: { rollupOptions: { input } } };
		},

		resolveId(source) {
			const id = normalizePath(source);
			return pageIds.has(id) ? id : null;
		},

		load(id) {
			const page = pageIds.get(normalizePath(id));
			return page ? renderSpecPage(readTemplate(), page, siteBase) : null;
		},

		configureServer(server) {
			server.middlewares.use(serveSpecPage(server, readTemplate, discoverSpecPages(uiRoot), siteBase));
		},
	} satisfies PluginOption;
}

function serveSpecPage(server: ViteDevServer, readTemplate: () => string, pages: SpecPage[], siteBase: string): Connect.NextHandleFunction {
	const byUrl = new Map(pages.map(page => [page.name, page]));

	return (req: IncomingMessage, res: ServerResponse<IncomingMessage>, next: Connect.NextFunction) => {
		const base = server.config.base;
		const parsed = new URL(req.url!, 'http://localhost');
		if (!parsed.pathname.startsWith(base)) return next();

		const relativePath = parsed.pathname.slice(base.length).replace(/(^|\/)index\.html$/, '$1');
		const page = byUrl.get(relativePath.replace(/\/$/, ''));
		if (!page) return next();

		// '/classic/<page>' -> '/classic/<page>/', so relative URLs in the page keep working.
		if (!relativePath.endsWith('/')) {
			res.writeHead(301, { Location: `${base}${relativePath}/${parsed.search}` });
			res.end();
			return;
		}

		// transformIndexHtml keys its inline-script proxy modules off a root-relative url.
		server
			.transformIndexHtml(`/${page.outPath}`, renderSpecPage(readTemplate(), page, siteBase), (req as Connect.IncomingMessage).originalUrl)
			.then(html => {
				res.writeHead(200, { 'Content-Type': 'text/html' });
				res.end(html);
			})
			.catch(next);
	};
}
