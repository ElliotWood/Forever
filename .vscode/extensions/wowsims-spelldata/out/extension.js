"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.activate = activate;
exports.deactivate = deactivate;
const node_child_process_1 = require("node:child_process");
const node_fs_1 = require("node:fs");
const node_path_1 = require("node:path");
const vscode = require("vscode");
const spellIds_1 = require("./spellIds");
const LANGUAGES = ['go', 'json', 'typescript', 'typescriptreact'];
function activate(context) {
    const log = vscode.window.createOutputChannel('WoWSims Spelldata');
    // One promise per repository and question, so the hovers that arrive while the first `go run` is
    // still compiling wait on that run rather than starting their own.
    const cache = new Map();
    // The chains each class folder declares, read once per folder and dropped when a file in it is
    // saved: an identifier the map does not know is not worth a process.
    const declarations = new Map();
    let warnedMissingGo = false;
    function ask(root, key, args) {
        const full = `${root}\u0000${key}`;
        let pending = cache.get(full);
        if (pending === undefined) {
            pending = readTool(root, args, log, () => {
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
    const provider = {
        async provideHover(document, position) {
            const root = moduleRoot(document.uri.fsPath);
            if (root === undefined) {
                return undefined;
            }
            const line = document.lineAt(position.line).text;
            const id = (0, spellIds_1.spellIdAt)(line, position.character);
            if (id !== undefined) {
                const row = await ask(root, `id\u0000${id}`, ['-json', String(id)]);
                return row === undefined ? undefined : new vscode.Hover(spellCard(row, `**${row.title}**`));
            }
            // The rest are Go names, and reading a folder for them only makes sense in a class package.
            if (document.languageId !== 'go') {
                return undefined;
            }
            const folder = (0, node_path_1.dirname)(document.uri.fsPath);
            const pkg = (0, node_path_1.basename)(folder);
            const field = (0, spellIds_1.familyFieldAt)(line, position.character);
            if (field !== undefined) {
                const family = await ask(root, `family\u0000${pkg}/${field}`, ['-family', `${pkg}/${field}`, '-json']);
                return family === undefined ? undefined : new vscode.Hover(familyCard(family));
            }
            const hover = (0, spellIds_1.chainHoverAt)(line, position.character, declaredChains(folder, declarations, log));
            if (hover === undefined) {
                return undefined;
            }
            const expr = hover.expr;
            const row = await ask(root, `expr\u0000${pkg}\u0000${expr}`, ['-expr', expr, '-package', pkg, '-json']);
            return row === undefined ? undefined : new vscode.Hover(exprCard(row, hover));
        },
    };
    context.subscriptions.push(log, vscode.languages.registerHoverProvider(LANGUAGES, provider), vscode.workspace.onDidSaveTextDocument(document => declarations.delete((0, node_path_1.dirname)(document.uri.fsPath))));
}
function deactivate() { }
// The repository the file belongs to, which is the folder holding go.mod: the tool is run from there
// and reads the store out of it, so a file outside one gets no hover.
function moduleRoot(filePath) {
    let dir = (0, node_path_1.dirname)(filePath);
    for (;;) {
        if ((0, node_fs_1.existsSync)((0, node_path_1.join)(dir, 'go.mod'))) {
            return dir;
        }
        const parent = (0, node_path_1.dirname)(dir);
        if (parent === dir) {
            return undefined;
        }
        dir = parent;
    }
}
function declaredChains(folder, cache, log) {
    let found = cache.get(folder);
    if (found !== undefined) {
        return found;
    }
    found = new Map();
    try {
        for (const entry of (0, node_fs_1.readdirSync)(folder)) {
            if (!entry.endsWith('.go')) {
                continue;
            }
            for (const [name, expr] of (0, spellIds_1.chainDeclarations)((0, node_fs_1.readFileSync)((0, node_path_1.join)(folder, entry), 'utf8'))) {
                found.set(name, expr);
            }
        }
    }
    catch (error) {
        log.appendLine(`${folder}: ${String(error)}`);
    }
    cache.set(folder, found);
    return found;
}
function readTool(root, args, log, onMissingGo) {
    const goBinary = vscode.workspace.getConfiguration('wowsims-spelldata').get('goBinary', 'go');
    const run = ['run', './tools/spelldata', ...args];
    return new Promise(resolve => {
        const child = (0, node_child_process_1.spawn)(goBinary, run, { cwd: root });
        let stdout = '';
        let stderr = '';
        child.stdout.on('data', chunk => (stdout += chunk));
        child.stderr.on('data', chunk => (stderr += chunk));
        child.on('error', error => {
            log.appendLine(`${goBinary} ${run.join(' ')}: ${error.message}`);
            if (error.code === 'ENOENT') {
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
                resolve(JSON.parse(stdout));
            }
            catch (error) {
                log.appendLine(`${run.join(' ')}: ${String(error)}`);
                resolve(undefined);
            }
        });
    });
}
function spellCard(row, heading, md = new vscode.MarkdownString(), read = 0) {
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
        const label = index + 1 === read ? `effect ${index + 1} (read)` : `effect ${index + 1}`;
        md.appendMarkdown(`\n${label}${effect.human === '' ? '' : ` — ${effect.human}`}\n`);
        md.appendCodeblock(effect.literal);
    });
    md.appendMarkdown(`\n[Wowhead](${row.wowhead})`);
    return md;
}
// What a chain answered. A pick is the row it names; an accessor call that stopped on an effect is that
// effect and the accessors that read something off it; a value is the number, the chain it was read
// through and, where the cursor was on the accessor itself, what the store says that accessor answers.
function exprCard(row, hover) {
    if (row.kind === 'spell') {
        return spellCard(row, `\`${hover.label}\` = ${row.title}`);
    }
    const md = new vscode.MarkdownString();
    const called = lastCall(row.trail);
    if (row.kind === 'effect') {
        md.appendMarkdown(`\`${called}\` = **${row.value}** of ${row.title}\n\n`);
        md.appendMarkdown(`\`${row.trail}\`\n\n`);
        const effect = row.effects[row.read_effect - 1];
        if (effect !== undefined) {
            md.appendMarkdown(`${effect.human}\n`);
            md.appendCodeblock(effect.literal);
        }
        for (const accessor of row.accessors) {
            md.appendMarkdown(`\`${accessor}\`  \n`);
        }
        md.appendMarkdown(`\n[Wowhead](${row.wowhead})`);
        return md;
    }
    md.appendMarkdown(`\`${hover.segment ? called : hover.label}\` = **${row.value}**\n\n`);
    md.appendMarkdown(`\`${row.trail}\`\n\n`);
    if (hover.segment && row.doc !== '') {
        md.appendMarkdown(`${row.doc}\n\n`);
    }
    return spellCard(row, `**${row.title}**`, md, row.read_effect);
}
// The accessor the chain ends on, as the tool substituted it: `Average(60)` where the file wrote
// `Average(core.CharacterLevel)`.
function lastCall(trail) {
    const match = /\.([A-Za-z_]\w*\([^()]*\))$/.exec(trail);
    return match === null ? trail : match[1];
}
// The whole ladder: every rank with the call that reaches it, then the highest rank in full.
function familyCard(family) {
    const md = new vscode.MarkdownString();
    md.appendMarkdown(`**${family.family}**\n\n`);
    md.appendMarkdown('| id | name | rank | call |\n| --- | --- | --- | --- |\n');
    for (const rank of family.ranks) {
        md.appendMarkdown(`| ${rank.id} | ${rank.name} | ${rank.rank} | \`${rank.accessor}\` |\n`);
    }
    md.appendMarkdown('\n');
    return spellCard(family.highest, `**${family.highest.title}**`, md);
}
//# sourceMappingURL=extension.js.map