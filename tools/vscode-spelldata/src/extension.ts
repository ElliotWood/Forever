import { existsSync } from 'node:fs';
import { join } from 'node:path';
import * as vscode from 'vscode';
import { LanguageClient, LanguageClientOptions, ServerOptions, TransportKind } from 'vscode-languageclient/node';

const LANGUAGES = ['go', 'json', 'typescript', 'typescriptreact'];

let client: LanguageClient | undefined;

export async function activate(context: vscode.ExtensionContext): Promise<void> {
	const root = vscode.workspace.workspaceFolders?.map(folder => folder.uri.fsPath).find(path => existsSync(join(path, 'go.mod')));
	if (root === undefined) {
		return;
	}

	const settings = vscode.workspace.getConfiguration('wowsims-spelldata');
	const serverOptions: ServerOptions = {
		command: settings.get<string>('goBinary', 'go'),
		args: ['run', './tools/spelldata', '-lsp'],
		options: { cwd: root },
		transport: TransportKind.stdio,
	};
	const outputChannel = vscode.window.createOutputChannel('WoWSims Spelldata', { log: true });
	const clientOptions: LanguageClientOptions = {
		documentSelector: LANGUAGES.map(language => ({ scheme: 'file', language })),
		initializationOptions: { trace: settings.get<string>('trace', 'on') },
		outputChannel,
	};

	client = new LanguageClient('wowsims-spelldata', 'WoWSims Spelldata', serverOptions, clientOptions);
	context.subscriptions.push(outputChannel);
	await client.start();
}

export function deactivate(): Promise<void> | undefined {
	return client?.stop();
}
