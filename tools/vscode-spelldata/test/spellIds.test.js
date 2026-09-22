const assert = require('node:assert/strict');
const { test } = require('node:test');

const { spellIdAt } = require('../out/spellIds.js');

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
