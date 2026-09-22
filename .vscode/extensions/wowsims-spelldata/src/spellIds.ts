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
		const match = matchCovering(pattern, lineText, column);
		if (match !== undefined) {
			return Number(match[1]);
		}
	}
	return undefined;
}

// A class file's ladder, as `spellData.Execute` names it. The cursor on the prefix or on the field
// answers the field; on the accessor that follows it does not, since that is a rank rather than the
// family.
const FAMILY_PATTERN = /\bspellData\.([A-Za-z_]\w*)/g;

export function familyFieldAt(lineText: string, column: number): string | undefined {
	return matchCovering(FAMILY_PATTERN, lineText, column)?.[1];
}

// Neither a name a declaration can carry nor one worth asking the tool about.
const GO_KEYWORDS = new Set(
	(
		'break case chan const continue default defer else fallthrough for func go goto if import ' +
		'interface map package range return select struct switch type var'
	).split(' '),
);

export function identifierAt(lineText: string, column: number): string | undefined {
	const match = matchCovering(/[A-Za-z_]\w*/g, lineText, column);
	if (match === undefined || GO_KEYWORDS.has(match[0])) {
		return undefined;
	}
	return match[0];
}

// A name bound to one rank of a ladder: `var executeRank = spellData.Execute.Highest()` at package
// level, or the same inside a function with `:=`. The value is the call, which is what -expr takes.
const DECLARATION_PATTERN = /\b([A-Za-z_]\w*)\s*(?::=|=)\s*(?:spellData\.)?([A-Za-z_]\w*)\.(Highest\(\)|Rank\(\d+\)|ByID\(\d+\))/g;

export function ladderDeclarations(text: string): Map<string, string> {
	const found = new Map<string, string>();
	const scan = new RegExp(DECLARATION_PATTERN.source, 'g');

	// A stubbed class file keeps its ladder picks as comments, and a commented name is not one a hover
	// should spend a process on.
	for (const line of text.split('\n')) {
		if (line.trimStart().startsWith('//')) {
			continue;
		}
		scan.lastIndex = 0;
		let match: RegExpExecArray | null;
		while ((match = scan.exec(line)) !== null) {
			found.set(match[1], `spellData.${match[2]}.${match[3]}`);
		}
	}
	return found;
}

// The first match of a pattern whose span covers the column, the column at either end included.
function matchCovering(pattern: RegExp, lineText: string, column: number): RegExpExecArray | undefined {
	const scan = new RegExp(pattern.source, 'g');
	let match: RegExpExecArray | null;
	while ((match = scan.exec(lineText)) !== null) {
		if (column >= match.index && column <= match.index + match[0].length) {
			return match;
		}
	}
	return undefined;
}
