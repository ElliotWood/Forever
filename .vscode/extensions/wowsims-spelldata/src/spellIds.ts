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

// The accessors that name a rank of a ladder, which is where every chain the tool reads starts.
const PICKS = new Set(['Highest', 'Rank', 'ByID']);

// The names a class file writes into an accessor call, which the tool substitutes.
const NAMED_VALUES = new Set(['core.CharacterLevel']);

const CHAIN_HEAD = /^(?:spellData\.)?[A-Za-z_]\w*/;
const CHAIN_SEGMENT = /^\.([A-Za-z_]\w*)(?:\(([^()]*)\))?/;

// A name bound to a chain: `var executeRank = spellData.Execute.Highest()` at package level, or the
// same inside a function with `:=`. `=(?!=)` keeps a comparison out.
const DECLARATION = /^(\s*(?:var\s+)?)([A-Za-z_]\w*)(\s*(?::=|=(?!=))\s*)(\S.*)$/;

const MAX_SUBSTITUTIONS = 4;

// One `.Name(args)` of a chain, with the columns it spans so a cursor can be placed in it.
interface Segment {
	name: string;
	text: string;
	start: number;
	end: number;
}

interface Chain {
	head: string;
	headStart: number;
	headEnd: number;
	segments: Segment[];
}

// What a hover is about: the name or accessor call under the cursor, and the chain that reaches it -
// the declaration's chain cut off after that call, resolved to the ladder pick at its head, which is
// what `-expr` takes.
export interface ChainHover {
	label: string;
	expr: string;
	segment: boolean;
}

// A chain the tool can read: a name or `spellData.<Family>`, then accessor calls whose every argument
// is a number or one of the names sim/core states. Anything else - an operator, a conversion, a
// variable as an argument - is none, since the tool would have to run the package to read it.
function parseChain(text: string, offset = 0): Chain | undefined {
	const head = CHAIN_HEAD.exec(text);
	if (head === null) {
		return undefined;
	}

	const segments: Segment[] = [];
	let at = head[0].length;
	while (at < text.length) {
		const match = CHAIN_SEGMENT.exec(text.slice(at));
		if (match === null || (match[2] !== undefined && !argumentsAreLiterals(match[2]))) {
			return undefined;
		}
		segments.push({
			name: match[1],
			text: match[0].slice(1),
			start: offset + at,
			end: offset + at + match[0].length,
		});
		at += match[0].length;
	}
	return { head: head[0], headStart: offset, headEnd: offset + head[0].length, segments };
}

function argumentsAreLiterals(args: string): boolean {
	if (args.trim() === '') {
		return true;
	}
	return args.split(',').every(arg => /^-?\d+(?:\.\d+)?$/.test(arg.trim()) || NAMED_VALUES.has(arg.trim()));
}

function chainText(chain: Chain, upTo = chain.segments.length): string {
	return (
		chain.head +
		chain.segments
			.slice(0, upTo)
			.map(segment => `.${segment.text}`)
			.join('')
	);
}

function isPick(segment: Segment | undefined): boolean {
	return segment !== undefined && PICKS.has(segment.name) && segment.text.endsWith(')');
}

// Every name a class file's folder binds to a chain, so a hover on one of them knows what it reads.
export function chainDeclarations(text: string): Map<string, string> {
	const found = new Map<string, string>();
	for (const line of text.split('\n')) {
		const declared = declarationOnLine(line);
		if (declared !== undefined) {
			found.set(declared.name, chainText(declared.chain));
		}
	}
	return found;
}

interface Declaration {
	name: string;
	nameStart: number;
	nameEnd: number;
	chain: Chain;
}

function declarationOnLine(lineText: string): Declaration | undefined {
	// A stubbed class file keeps its picks as comments, and a commented name is not one a hover should
	// spend a process on.
	if (lineText.trimStart().startsWith('//')) {
		return undefined;
	}
	const match = DECLARATION.exec(lineText.split('//')[0]);
	if (match === null) {
		return undefined;
	}

	const nameStart = match[1].length;
	const valueStart = nameStart + match[2].length + match[3].length;
	const chain = parseChain(match[4].trimEnd(), valueStart);
	if (chain === undefined) {
		return undefined;
	}
	return { name: match[2], nameStart, nameEnd: nameStart + match[2].length, chain };
}

// The chain with every declared name it stands on inlined, down to the ladder pick at its head: the
// expression `-expr` takes. A head no declaration binds and no pick follows is not one the tool reads,
// and neither is a chain that stands on more than a few names or on itself.
export function resolveChain(chain: string, declarations: Map<string, string>, depth = MAX_SUBSTITUTIONS, seen = new Set<string>()): string | undefined {
	const parsed = parseChain(chain);
	if (parsed === undefined) {
		return undefined;
	}

	if (parsed.head.startsWith('spellData.')) {
		return isPick(parsed.segments[0]) ? chain : undefined;
	}

	const bound = declarations.get(parsed.head);
	if (bound === undefined) {
		// A ladder as a class file names it inside its own package, without the `spellData.` prefix.
		return isPick(parsed.segments[0]) ? `spellData.${chain}` : undefined;
	}
	if (depth <= 0 || seen.has(parsed.head)) {
		return undefined;
	}

	seen.add(parsed.head);
	const head = resolveChain(bound, declarations, depth - 1, seen);
	if (head === undefined) {
		return undefined;
	}
	return head + parsed.segments.map(segment => `.${segment.text}`).join('');
}

// What the cursor is on: the name a declaration binds, one accessor call of its chain - which reads
// the chain cut off there - or a name used elsewhere in the package. Anything else is silent.
export function chainHoverAt(lineText: string, column: number, declarations: Map<string, string>): ChainHover | undefined {
	const declared = declarationOnLine(lineText);
	if (declared !== undefined) {
		const hover = hoverInDeclaration(declared, column, declarations);
		if (hover !== undefined) {
			return hover;
		}
	}

	const name = identifierAt(lineText, column);
	if (name === undefined || !declarations.has(name)) {
		return undefined;
	}
	return resolved(name, name, false, declarations);
}

function hoverInDeclaration(declared: Declaration, column: number, declarations: Map<string, string>): ChainHover | undefined {
	if (column >= declared.nameStart && column <= declared.nameEnd) {
		return resolved(declared.name, chainText(declared.chain), false, declarations);
	}

	const found = declared.chain.segments.findIndex(segment => column >= segment.start && column <= segment.end);
	if (found >= 0) {
		return resolved(declared.chain.segments[found].text, chainText(declared.chain, found + 1), true, declarations);
	}
	if (column >= declared.chain.headStart && column <= declared.chain.headEnd) {
		return resolved(declared.chain.head, declared.chain.head, true, declarations);
	}
	return undefined;
}

function resolved(label: string, chain: string, segment: boolean, declarations: Map<string, string>): ChainHover | undefined {
	const expr = resolveChain(chain, declarations);
	return expr === undefined ? undefined : { label, expr, segment };
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
