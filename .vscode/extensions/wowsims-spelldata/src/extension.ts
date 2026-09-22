import { spawn } from 'node:child_process';
import { existsSync, readdirSync, readFileSync } from 'node:fs';
import { basename, dirname, join } from 'node:path';
import * as vscode from 'vscode';

import { familyFieldAt, identifierAt, ladderDeclarations, spellIdAt } from './spellIds';

interface EffectLine {
	human: string;
	literal: string;
}

interface SpellRow {
	id: number;
	name: string;
	rank: string;
	title: string;
	ladder: string[];
	header: string[];
	effects: EffectLine[];
	proc: string;
	refs: string[];
	wowhead: string;
}

interface FamilyRow {
	family: string;
	ranks: { id: number; name: string; rank: string; accessor: string }[];
	highest: SpellRow;
}

const LANGUAGES = ['go', 'json', 'typescript', 'typescriptreact'];

export function activate(context: vscode.ExtensionContext): void {
	const log = vscode.window.createOutputChannel('WoWSims Spelldata');
	// One promise per repository and question, so the hovers that arrive while the first `go run` is
	// still compiling wait on that run rather than starting their own.
	const cache = new Map<string, Promise<unknown>>();
	// The ladder picks each class folder declares, read once per folder and dropped when a file in it
	// is saved: an identifier the map does not know is not worth a process.
	const declarations = new Map<string, Map<string, string>>();
	let warnedMissingGo = false;

	function ask<T>(root: string, key: string, args: string[]): Promise<T | undefined> {
		const full = `${root}\u0000${key}`;
		let pending = cache.get(full) as Promise<T | undefined> | undefined;
		if (pending === undefined) {
			pending = readTool<T>(root, args, log, () => {
				if (!warnedMissingGo) {
					warnedMissingGo = true;
					void vscode.window.showWarningMessage('WoWSims Spelldata: the go binary was not found. Set wowsims-spelldata.goBinary to its path.');
				}
			});
			cache.set(full, pending);
		}
		return pending.then(row => {
			// A failure is not an answer: a tree that did not compile a moment ago compiles once it is
			// fixed, and the next hover has to ask again.
			if (row === undefined) {
				cache.delete(full);
			}
			return row;
		});
	}

	const provider: vscode.HoverProvider = {
		async provideHover(document, position) {
			const root = moduleRoot(document.uri.fsPath);
			if (root === undefined) {
				return undefined;
			}
			const line = document.lineAt(position.line).text;

			const id = spellIdAt(line, position.character);
			if (id !== undefined) {
				const row = await ask<SpellRow>(root, `id\u0000${id}`, ['-json', String(id)]);
				return row === undefined ? undefined : new vscode.Hover(spellCard(row, `**${row.title}**`));
			}

			// The rest are Go names, and reading a folder for them only makes sense in a class package.
			if (document.languageId !== 'go') {
				return undefined;
			}
			const folder = dirname(document.uri.fsPath);
			const pkg = basename(folder);

			const field = familyFieldAt(line, position.character);
			if (field !== undefined) {
				const family = await ask<FamilyRow>(root, `family\u0000${pkg}/${field}`, ['-family', `${pkg}/${field}`, '-json']);
				return family === undefined ? undefined : new vscode.Hover(familyCard(family));
			}

			const name = identifierAt(line, position.character);
			if (name === undefined) {
				return undefined;
			}
			const expr = declaredLadderPicks(folder, declarations, log).get(name);
			if (expr === undefined) {
				return undefined;
			}
			const row = await ask<SpellRow>(root, `expr\u0000${pkg}\u0000${expr}`, ['-expr', expr, '-package', pkg, '-json']);
			return row === undefined ? undefined : new vscode.Hover(spellCard(row, `\`${name}\` = ${row.title}`));
		},
	};

	context.subscriptions.push(
		log,
		vscode.languages.registerHoverProvider(LANGUAGES, provider),
		vscode.workspace.onDidSaveTextDocument(document => declarations.delete(dirname(document.uri.fsPath))),
	);
}

export function deactivate(): void {}

// The repository the file belongs to, which is the folder holding go.mod: the tool is run from there
// and reads the store out of it, so a file outside one gets no hover.
function moduleRoot(filePath: string): string | undefined {
	let dir = dirname(filePath);
	for (;;) {
		if (existsSync(join(dir, 'go.mod'))) {
			return dir;
		}
		const parent = dirname(dir);
		if (parent === dir) {
			return undefined;
		}
		dir = parent;
	}
}

function declaredLadderPicks(folder: string, cache: Map<string, Map<string, string>>, log: vscode.OutputChannel): Map<string, string> {
	let found = cache.get(folder);
	if (found !== undefined) {
		return found;
	}

	found = new Map<string, string>();
	try {
		for (const entry of readdirSync(folder)) {
			if (!entry.endsWith('.go')) {
				continue;
			}
			for (const [name, expr] of ladderDeclarations(readFileSync(join(folder, entry), 'utf8'))) {
				found.set(name, expr);
			}
		}
	} catch (error) {
		log.appendLine(`${folder}: ${String(error)}`);
	}
	cache.set(folder, found);
	return found;
}

function readTool<T>(root: string, args: string[], log: vscode.OutputChannel, onMissingGo: () => void): Promise<T | undefined> {
	const goBinary = vscode.workspace.getConfiguration('wowsims-spelldata').get<string>('goBinary', 'go');
	const run = ['run', './tools/spelldata', ...args];

	return new Promise(resolve => {
		const child = spawn(goBinary, run, { cwd: root });
		let stdout = '';
		let stderr = '';
		child.stdout.on('data', chunk => (stdout += chunk));
		child.stderr.on('data', chunk => (stderr += chunk));

		child.on('error', error => {
			log.appendLine(`${goBinary} ${run.join(' ')}: ${error.message}`);
			if ((error as NodeJS.ErrnoException).code === 'ENOENT') {
				onMissingGo();
			}
			resolve(undefined);
		});

		child.on('close', code => {
			if (code !== 0) {
				log.appendLine(`${run.join(' ')}: the tool exited ${code}: ${stderr.trim()}`);
				resolve(undefined);
				return;
			}
			try {
				resolve(JSON.parse(stdout) as T);
			} catch (error) {
				log.appendLine(`${run.join(' ')}: ${String(error)}`);
				resolve(undefined);
			}
		});
	});
}

function spellCard(row: SpellRow, heading: string, md = new vscode.MarkdownString()): vscode.MarkdownString {
	md.appendMarkdown(`${heading}\n\n`);
	for (const call of row.ladder) {
		md.appendMarkdown(`\`${call}\`  \n`);
	}
	for (const line of row.header) {
		md.appendMarkdown(`${line}  \n`);
	}
	if (row.proc !== '') {
		md.appendMarkdown(`proc      ${row.proc}  \n`);
	}
	if (row.refs.length > 0) {
		md.appendMarkdown(`refs      ${row.refs.join(', ')}  \n`);
	}
	row.effects.forEach((effect, index) => {
		md.appendMarkdown(`\neffect ${index + 1}${effect.human === '' ? '' : ` — ${effect.human}`}\n`);
		md.appendCodeblock(effect.literal);
	});
	md.appendMarkdown(`\n[Wowhead](${row.wowhead})`);
	return md;
}

// The whole ladder: every rank with the call that reaches it, then the highest rank in full.
function familyCard(family: FamilyRow): vscode.MarkdownString {
	const md = new vscode.MarkdownString();
	md.appendMarkdown(`**${family.family}**\n\n`);
	md.appendMarkdown('| id | name | rank | call |\n| --- | --- | --- | --- |\n');
	for (const rank of family.ranks) {
		md.appendMarkdown(`| ${rank.id} | ${rank.name} | ${rank.rank} | \`${rank.accessor}\` |\n`);
	}
	md.appendMarkdown('\n');
	return spellCard(family.highest, `**${family.highest.title}**`, md);
}
