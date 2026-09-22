const assert = require('node:assert/strict');
const { test } = require('node:test');

const { chainDeclarations, chainHoverAt, familyFieldAt, identifierAt, resolveChain, spellIdAt } = require('../out/spellIds.js');

// The lines sim/warrior/execute.go states, which is the case the hover was built for.
const EXECUTE_GO = [
	'var executeRank = spellData.Execute.Highest()',
	'',
	'// TODO: The dummy effect carries the base damage, and no finder picks it out: the row states no',
	'// damage effect, and both of its effects share the aura and misc values Effect() selects on.',
	'var executeBaseDamage = executeRank.EffectN(1).Average(core.CharacterLevel)',
	'',
	"// The tooltip's $*10;F1: the dummy's chain amplitude, times 10, per extra point of rage.",
	'var executeDamagePerRage = float64(executeRank.EffectN(1).ChainAmp) * 10',
].join('\n');

const DERIVED_LINE = 'var executeBaseDamage = executeRank.EffectN(1).Average(core.CharacterLevel)';

test('the id under the cursor, in each of the shapes code names one', () => {
	const cases = [
		['\tvar rend = spelldata.MustFind(11574)', 'MustFind(11574'],
		['\tspelldata.Find(116)', 'Find(116'],
		['\tspellData.Rend.ByID(11574)', 'ByID(11574'],
		['\tActionID: core.ActionID{SpellID: 11574},', 'SpellID: 11574'],
		['\t\t\t"spellId": 11574', '"spellId": 11574'],
		['\tActionId.fromSpellId(23563)', 'fromSpellId(23563'],
		['\t\tspellId: 23563,', 'spellId: 23563'],
	];

	for (const [line, needle] of cases) {
		const column = line.indexOf(needle) + needle.length;
		const want = Number(needle.match(/(\d+)/)[1]);
		assert.equal(spellIdAt(line, column), want, line);
	}
});

test('a line with no id, and a column outside the match', () => {
	assert.equal(spellIdAt('\tspell.ApplyEffects = nil', 12), undefined);
	assert.equal(spellIdAt('\tspelldata.Find(116) // the rank 1', 30), undefined);
});

test('the ladder family under the cursor', () => {
	const line = '\tvar executeRank = spellData.Execute.Highest()';
	assert.equal(familyFieldAt(line, line.indexOf('Execute.') + 3), 'Execute');
	assert.equal(familyFieldAt(line, line.indexOf('spellData')), 'Execute');
	// The accessor names a rank, not the family, and the line's own id match answers there instead.
	assert.equal(familyFieldAt(line, line.indexOf('Highest') + 3), undefined);
	assert.equal(familyFieldAt('\tspellData.Rend.ByID(11574)', 12), 'Rend');
	assert.equal(familyFieldAt('\tspelldata.Find(116)', 12), undefined);
});

test('the Go identifier under the cursor, keywords aside', () => {
	assert.equal(identifierAt(DERIVED_LINE, DERIVED_LINE.indexOf('executeRank') + 4), 'executeRank');
	assert.equal(identifierAt(DERIVED_LINE, DERIVED_LINE.indexOf('EffectN') + 2), 'EffectN');
	assert.equal(identifierAt('\treturn spell', 4), undefined);
	assert.equal(identifierAt('\t\t\t\t', 2), undefined);
});

test('the chains a class file declares', () => {
	const found = chainDeclarations(EXECUTE_GO);
	assert.equal(found.get('executeRank'), 'spellData.Execute.Highest()');
	assert.equal(found.get('executeBaseDamage'), 'executeRank.EffectN(1).Average(core.CharacterLevel)');
	// A conversion is not a chain: the tool would have to run the package to read it.
	assert.equal(found.get('executeDamagePerRage'), undefined);
	assert.equal(found.size, 2);
});

test('the chains a class file does not declare', () => {
	const file = [
		'\tsecondRank := spellData.Execute.Rank(2)',
		'\tnamed := spellData.Execute.ByID(20658)',
		'\tbare := Cruelty.Rank(3)',
		'\tif executeRank == spellData.Execute.Highest() {',
		'// var stubbedRank = spellData.MangleBear.ByID(33987)',
		'\tconfig := spelldata.SpellConfig(&warrior.Unit, executeRank, spelldata.Melee(cost))',
		'\tmaxRage := warrior.MaximumRage() - spell.Cost.GetCurrentCost()',
		'\tticks := executeRank.EffectN(1).Average(level)',
	].join('\n');

	const found = chainDeclarations(file);
	assert.equal(found.get('secondRank'), 'spellData.Execute.Rank(2)');
	assert.equal(found.get('named'), 'spellData.Execute.ByID(20658)');
	assert.equal(found.get('bare'), 'Cruelty.Rank(3)');
	assert.equal(found.get('executeRank'), undefined);
	assert.equal(found.get('stubbedRank'), undefined);
	assert.equal(found.get('config'), undefined);
	assert.equal(found.get('maxRage'), undefined);
	// An argument that is a variable rather than a literal is not one the tool substitutes.
	assert.equal(found.get('ticks'), undefined);
	assert.equal(found.size, 3);
});

test('a chain resolves to the pick at its head', () => {
	const declared = chainDeclarations(EXECUTE_GO);

	assert.equal(resolveChain('executeBaseDamage', declared), 'spellData.Execute.Highest().EffectN(1).Average(core.CharacterLevel)');
	assert.equal(resolveChain('executeRank', declared), 'spellData.Execute.Highest()');
	assert.equal(resolveChain('executeRank.EffectN(1)', declared), 'spellData.Execute.Highest().EffectN(1)');
	// A ladder as a class file names it inside its own package.
	assert.equal(resolveChain('Cruelty.Rank(3)', declared), 'spellData.Cruelty.Rank(3)');
	assert.equal(resolveChain('spellData.Rend.Highest().EffectN(1).Period()', declared), 'spellData.Rend.Highest().EffectN(1).Period()');

	// A head no declaration binds and no pick follows, a family with no pick, and an argument that is
	// not a literal.
	assert.equal(resolveChain('warrior.MaximumRage()', declared), undefined);
	assert.equal(resolveChain('spellData.Execute', declared), undefined);
	assert.equal(resolveChain('executeRank.EffectN(n)', declared), undefined);
});

test('the substitution stops at four names and at a cycle', () => {
	const deep = new Map([
		['a', 'spellData.Execute.Highest()'],
		['b', 'a.EffectN(1)'],
		['c', 'b.Average(60)'],
		['d', 'c.Min(60)'],
		['e', 'd.Min(60)'],
	]);
	assert.equal(resolveChain('c', deep), 'spellData.Execute.Highest().EffectN(1).Average(60)');
	assert.equal(resolveChain('e', deep), undefined);

	const cyclic = new Map([
		['loop', 'other.EffectN(1)'],
		['other', 'loop.EffectN(2)'],
		['self', 'self.EffectN(1)'],
	]);
	assert.equal(resolveChain('loop', cyclic), undefined);
	assert.equal(resolveChain('self', cyclic), undefined);
});

test('what the cursor is on, along the declaration of a derived value', () => {
	const declared = chainDeclarations(EXECUTE_GO);
	const at = column => chainHoverAt(DERIVED_LINE, column, declared);

	assert.deepEqual(at(DERIVED_LINE.indexOf('executeBaseDamage') + 4), {
		label: 'executeBaseDamage',
		expr: 'spellData.Execute.Highest().EffectN(1).Average(core.CharacterLevel)',
		segment: false,
	});
	assert.deepEqual(at(DERIVED_LINE.indexOf('executeRank.') + 4), {
		label: 'executeRank',
		expr: 'spellData.Execute.Highest()',
		segment: true,
	});
	assert.deepEqual(at(DERIVED_LINE.indexOf('EffectN') + 2), {
		label: 'EffectN(1)',
		expr: 'spellData.Execute.Highest().EffectN(1)',
		segment: true,
	});
	assert.deepEqual(at(DERIVED_LINE.indexOf('Average') + 2), {
		label: 'Average(core.CharacterLevel)',
		expr: 'spellData.Execute.Highest().EffectN(1).Average(core.CharacterLevel)',
		segment: true,
	});
	// The name an argument states is inside the call it belongs to, not a hover of its own.
	assert.deepEqual(at(DERIVED_LINE.indexOf('core.CharacterLevel') + 2), {
		label: 'Average(core.CharacterLevel)',
		expr: 'spellData.Execute.Highest().EffectN(1).Average(core.CharacterLevel)',
		segment: true,
	});
});

test('what the cursor is on elsewhere in the package', () => {
	const declared = chainDeclarations(EXECUTE_GO);
	const use = '\t\tbaseDamage := executeBaseDamage + executeDamagePerRage*extraRage';

	assert.deepEqual(chainHoverAt(use, use.indexOf('executeBaseDamage') + 4, declared), {
		label: 'executeBaseDamage',
		expr: 'spellData.Execute.Highest().EffectN(1).Average(core.CharacterLevel)',
		segment: false,
	});
	// A bare accessor name outside a declaration, and a name no declaration binds, stay silent: a chain
	// standing on a name nothing binds reaches no ladder pick.
	const call = '\t\tamount := spell.EffectN(1).Average(60)';
	assert.equal(chainHoverAt(call, call.indexOf('EffectN') + 2, declared), undefined);
	assert.equal(chainHoverAt(use, use.indexOf('extraRage') + 2, declared), undefined);
});
