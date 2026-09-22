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
    // One promise per repository and id, so the hovers that arrive while the first `go run` is still
    // compiling wait on that run rather than starting their own.
    const cache = new Map();
    let warnedMissingGo = false;
    const provider = {
        async provideHover(document, position) {
            const line = document.lineAt(position.line).text;
            const id = (0, spellIds_1.spellIdAt)(line, position.character);
            if (id === undefined) {
                return undefined;
            }
            const root = moduleRoot(document.uri.fsPath);
            if (root === undefined) {
                return undefined;
            }
            const key = `${root}\u0000${id}`;
            let pending = cache.get(key);
            if (pending === undefined) {
                pending = readSpell(root, id, log, () => {
                    if (!warnedMissingGo) {
                        warnedMissingGo = true;
                        void vscode.window.showWarningMessage('WoWSims Spelldata: the go binary was not found. Set wowsims-spelldata.goBinary to its path.');
                    }
                });
                cache.set(key, pending);
            }
            const row = await pending;
            if (row === undefined) {
                // A failure is not an answer: a tree that did not compile a moment ago compiles once it
                // is fixed, and the next hover has to ask again.
                cache.delete(key);
                return undefined;
            }
            return new vscode.Hover(render(row));
        },
    };
    context.subscriptions.push(log, vscode.languages.registerHoverProvider(LANGUAGES, provider));
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
function readSpell(root, id, log, onMissingGo) {
    const goBinary = vscode.workspace.getConfiguration('wowsims-spelldata').get('goBinary', 'go');
    return new Promise(resolve => {
        const child = (0, node_child_process_1.spawn)(goBinary, ['run', './tools/spelldata', '-json', String(id)], { cwd: root });
        let stdout = '';
        let stderr = '';
        child.stdout.on('data', chunk => (stdout += chunk));
        child.stderr.on('data', chunk => (stderr += chunk));
        child.on('error', error => {
            log.appendLine(`${goBinary} run ./tools/spelldata -json ${id}: ${error.message}`);
            if (error.code === 'ENOENT') {
                onMissingGo();
            }
            resolve(undefined);
        });
        child.on('close', code => {
            if (code !== 0) {
                log.appendLine(`spell ${id}: the tool exited ${code}: ${stderr.trim()}`);
                resolve(undefined);
                return;
            }
            try {
                resolve(JSON.parse(stdout));
            }
            catch (error) {
                log.appendLine(`spell ${id}: ${String(error)}`);
                resolve(undefined);
            }
        });
    });
}
function render(row) {
    const md = new vscode.MarkdownString();
    md.appendMarkdown(`**${row.title}**\n\n`);
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
//# sourceMappingURL=extension.js.map