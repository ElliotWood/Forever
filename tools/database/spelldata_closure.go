package database

import (
	"database/sql"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"sort"
	"strconv"
	"strings"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/tools/database/dbc"
)

// Where the store's hand-kept extra ids live. Parsed rather than imported: package database is what
// gen_spelldata runs, and importing the store would let a generated file that does not compile stop
// the generator that rewrites it.
const extraIDsPath = "sim/core/spelldata/extra_ids.go"

// The spells the sim can reach without anything naming them first: everything the nine class files
// are built from, the store's own extra ids, and every spell an item, an enchant or a set bonus
// casts. What those spells trigger is reached from here by reachableSpells.
func storeRoots(db *sql.DB, t *spellTables, ladderIDs []int32) ([]int32, error) {
	roots := map[int32]bool{}
	for _, id := range ladderIDs {
		roots[id] = true
	}

	extras, err := parseExtraIDs()
	if err != nil {
		return nil, err
	}
	for _, id := range extras {
		roots[id] = true
	}

	for _, source := range []struct {
		name  string
		query string
	}{
		{"item effects", itemEffectSpellQuery()},
		{"enchant effects", enchantSpellQuery},
		{"set bonuses", `SELECT DISTINCT SpellID FROM ItemSetSpell WHERE SpellID > 0`},
	} {
		if err := eachRow(db, source.query, func(rows *sql.Rows) error {
			var id int32
			if err := rows.Scan(&id); err != nil {
				return err
			}
			roots[id] = true
			return nil
		}); err != nil {
			return nil, fmt.Errorf("%s: %w", source.name, err)
		}
	}

	// An id no SpellName row carries is not a spell in this build: ItemEffect keeps rows for spells
	// the client dropped, and following one would put an empty row in the store.
	named := map[int32]bool{}
	for id := range roots {
		if _, ok := t.names[id]; ok {
			named[id] = true
		}
	}
	return sortedIDs(named), nil
}

// The gear and consumables gen_db ships, which is where the item procs the sim registers come from.
// The gear half is the predicate gen_db itself passes to LoadAndWriteRawItems; the consumable half
// states the classes LoadAndWriteConsumables selects, whose query cannot be reused as it stands
// because its own filter reads a subquery alias. Both halves let the allowlists through, as gen_db
// does: Hand of Justice and the raid consumables are shipped by id rather than by predicate.
func itemEffectSpellQuery() string {
	return `
		SELECT DISTINCT ie.SpellID
		FROM ItemEffect ie
		JOIN ItemXItemEffect ixie ON ixie.ItemEffectID = ie.ID
		JOIN Item i ON i.ID = ixie.ItemID
		LEFT JOIN ItemSparse s ON s.ID = i.ID
		WHERE ie.SpellID > 0 AND (
			(` + SimItemFilter(core.CharacterLevel) + `)
			OR (((i.ClassID = 0 AND i.SubclassID IS NOT 0 AND i.SubclassID IS NOT 8 AND i.SubclassID IS NOT 6)
			     OR (i.ClassID = 7 AND i.SubclassID = 2))
			    AND s.RequiredLevel >= 50 AND s.Display_lang != ''
			    AND s.Display_lang NOT LIKE '%Test%' AND s.Display_lang NOT LIKE 'QA%')
			OR i.ID IN (` + allowListedItemIDs() + `))`
}

// The items gen_db ships by id: ItemAllowList bypasses the gear filter and ConsumableAllowList is
// added to the consumable query. Hand of Justice and Skullflame Shield have no ItemSparse row in
// this build, which is why the join to it above is a left one: an allowlisted item is taken by id,
// whatever the sparse table says about it.
func allowListedItemIDs() string {
	ids := make([]int32, 0, len(ItemAllowList)+len(ConsumableAllowList))
	for id := range ItemAllowList {
		ids = append(ids, id)
	}
	ids = append(ids, ConsumableAllowList...)
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })

	list := make([]string, len(ids))
	for i, id := range ids {
		list[i] = strconv.Itoa(int(id))
	}
	return strings.Join(list, ", ")
}

// An enchant states the spell it applies in the EffectArg of the effect that applies it, which is
// what LoadAndWriteRawEnchants reads as its spell id: effect 1 and effect 3 are the two that name a
// spell.
const enchantSpellQuery = `
	SELECT DISTINCT EffectArg_0 FROM SpellItemEnchantment WHERE Effect_0 IN (1, 3) AND EffectArg_0 > 0
	UNION
	SELECT DISTINCT EffectArg_1 FROM SpellItemEnchantment WHERE Effect_1 IN (1, 3) AND EffectArg_1 > 0
	UNION
	SELECT DISTINCT EffectArg_2 FROM SpellItemEnchantment WHERE Effect_2 IN (1, 3) AND EffectArg_2 > 0`

// Every spell the roots reach: what an effect triggers, what an actionbar override swaps in, what a
// hand link names, and what a tooltip's $<id> token points at. A spell the sim registers reads its
// numbers off those, so the store carries them all rather than the roots alone.
func reachableSpells(t *spellTables, roots []int32) []int32 {
	seen := map[int32]bool{}
	queue := make([]int32, 0, len(roots))
	for _, id := range roots {
		if !seen[id] {
			seen[id] = true
			queue = append(queue, id)
		}
	}

	for len(queue) > 0 {
		id := queue[0]
		queue = queue[1:]

		for _, next := range spellEdges(t, id) {
			if seen[next] {
				continue
			}
			if _, named := t.names[next]; !named {
				continue
			}
			seen[next] = true
			queue = append(queue, next)
		}
	}
	return sortedIDs(seen)
}

// A_OVERRIDE_ACTIONBAR_SPELLS states the spell it swaps in as its base points, the way
// overrideReplacements reads it.
func spellEdges(t *spellTables, id int32) []int32 {
	var next []int32
	for _, e := range t.effects[id] {
		if e.TriggerID > 0 {
			next = append(next, e.TriggerID)
		}
		if e.Aura == int32(dbc.A_OVERRIDE_ACTIONBAR_SPELLS) && e.BasePoints > 0 {
			next = append(next, int32(e.BasePoints))
		}
	}
	next = append(next, handTriggers[id]...)
	return append(next, t.referencedIDs(id)...)
}

// The store's ExtraIDs, read out of its source: a Go program cannot ask a package for the values of
// a variable it does not import.
func parseExtraIDs() ([]int32, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, extraIDsPath, nil, 0)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %w", extraIDsPath, err)
	}

	var ids []int32
	found := false
	ast.Inspect(file, func(n ast.Node) bool {
		vs, ok := n.(*ast.ValueSpec)
		if !ok || len(vs.Names) != 1 || vs.Names[0].Name != "ExtraIDs" || len(vs.Values) != 1 {
			return true
		}
		lit, ok := vs.Values[0].(*ast.CompositeLit)
		if !ok {
			return true
		}
		found = true
		for _, elt := range lit.Elts {
			if id, ok := constInt(elt); ok {
				ids = append(ids, id)
			}
		}
		return false
	})
	if !found {
		return nil, fmt.Errorf("%s declares no ExtraIDs slice", extraIDsPath)
	}
	return ids, nil
}

func parseSpellID(s string) int32 {
	id, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return int32(id)
}
