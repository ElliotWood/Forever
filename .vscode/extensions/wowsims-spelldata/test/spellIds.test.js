const assert = require('node:assert/strict');
const { test } = require('node:test');

const { familyFieldAt, identifierAt, ladderDeclarations, spellIdAt } = require('../out/spellIds.js');

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
	const line = '\tvar executeBaseDamage = executeRank.EffectN(1).Average(core.CharacterLevel)';
	assert.equal(identifierAt(line, line.indexOf('executeRank') + 4), 'executeRank');
	assert.equal(identifierAt(line, line.indexOf('EffectN') + 2), 'EffectN');
	assert.equal(identifierAt(line, 2), undefined);
	assert.equal(identifierAt('\treturn spell', 4), undefined);
	assert.equal(identifierAt('\t\t\t\t', 2), undefined);
});

test('the ladder picks a class file declares', () => {
	const file = [
		'var executeRank = spellData.Execute.Highest()',
		'var whirlwindRank = spellData.Whirlwind.Highest()',
		'var executeBaseDamage = executeRank.EffectN(1).Average(core.CharacterLevel)',
		'\tsecondRank := spellData.Execute.Rank(2)',
		'\tnamed := spellData.Execute.ByID(20658)',
		'\tbare := Cruelty.Rank(3)',
		'\tif executeRank == spellData.Execute.Highest() {',
		'// var stubbedRank = spellData.MangleBear.ByID(33987)',
		'\tconfig := spelldata.SpellConfig(&warrior.Unit, executeRank, spelldata.Melee(cost))',
	].join('\n');

	const found = ladderDeclarations(file);
	assert.equal(found.get('executeRank'), 'spellData.Execute.Highest()');
	assert.equal(found.get('whirlwindRank'), 'spellData.Whirlwind.Highest()');
	assert.equal(found.get('secondRank'), 'spellData.Execute.Rank(2)');
	assert.equal(found.get('named'), 'spellData.Execute.ByID(20658)');
	assert.equal(found.get('bare'), 'spellData.Cruelty.Rank(3)');
	assert.equal(found.get('executeBaseDamage'), undefined);
	assert.equal(found.get('config'), undefined);
	assert.equal(found.get('stubbedRank'), undefined);
	assert.equal(found.size, 5);
});
