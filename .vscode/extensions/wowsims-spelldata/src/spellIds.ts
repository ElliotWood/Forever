// Every way hand-written code names a spell id: the store's accessors and a ladder's ByID in Go,
// core.ActionID's SpellID, an APL file's "spellId" and the TS ActionId.fromSpellId. Each pattern
// captures the id in group 1.
export const SPELL_ID_PATTERNS: readonly RegExp[] = [
	/\bMustFind\(\s*(\d+)\s*\)/g,
	/\bFind\(\s*(\d+)\s*\)/g,
	/\.ByID\(\s*(\d+)\s*\)/g,
	/\bSpellID:\s*(\d+)/g,
	/"spellId":\s*(\d+)/g,
	/\bfromSpellId\(\s*(\d+)\s*\)/g,
	/\bspellId:\s*(\d+)/g,
];

// The id the cursor sits in, anywhere within a match - on the accessor's name as readily as on the
// digits. The first pattern that covers the column wins.
export function spellIdAt(lineText: string, column: number): number | undefined {
	for (const pattern of SPELL_ID_PATTERNS) {
		const scan = new RegExp(pattern.source, 'g');
		let match: RegExpExecArray | null;
		while ((match = scan.exec(lineText)) !== null) {
			const start = match.index;
			const end = start + match[0].length;
			if (column >= start && column <= end) {
				return Number(match[1]);
			}
		}
	}
	return undefined;
}
