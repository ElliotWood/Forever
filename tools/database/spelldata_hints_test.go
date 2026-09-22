package database

// The proc shapes docs/spell_data.md names, read off the client database the way the store's rows
// are. Which shape a spell has is a reading of its tooltip, and a reading is what silently changes
// when a matcher is widened: Flurry states "$m1%" for its attack speed and Unbridled Wrath states
// "$m1% chance" for its rage, and only the second is a roll.
//
// Skips when tools/database/wowsims.db is absent, which is why CI is unaffected.

import (
	"os"
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/tools/database/overrides"
)

func TestProcShapeOfNamedSpells(t *testing.T) {
	DatabasePath = "wowsims.db"
	if _, err := os.Stat(DatabasePath); err != nil {
		t.Skipf("no client database at %s - run `make db` from a local WoW install to enable this gate", DatabasePath)
	}

	helper, err := NewDBHelper()
	if err != nil {
		t.Fatalf("opening %s: %v", DatabasePath, err)
	}
	defer helper.Close()

	tables, err := loadSpellTables(helper.db)
	if err != nil {
		t.Fatalf("loading the spell tables: %v", err)
	}

	for _, want := range []struct {
		id      int32
		name    string
		source  storeProcChanceSource
		effect  int8
		why     string
		hint    core.ProcHint
		hintWhy string
	}{
		{12317, "Enrage", procChanceColumn, 0, `the tooltip's "$h%" is the ProcChance column`,
			0, `"the victim of any damaging attack" states no outcome and no wording any hint reads`},
		{12322, "Unbridled Wrath", procChanceEffectN, 1, `"$m1% chance to generate Rage" is effect 1's ladder`,
			0, `"when you deal melee damage with a weapon" is the mask's own shape`},
		{12298, "Shield Specialization", procChanceEffectN, 2, `"$m2% chance to generate Rage", past the "$s1%" block value`,
			core.ProcHintOutcomeTaken, `"when you Block" is an outcome no ProcTypeMask has a bit for`},
		{12319, "Flurry", procChanceAlways, 0, `"$m1%" is the attack speed, and a crit is the condition`,
			0, `"after dealing a melee critical strike" is the crit the caller reads off the result, since the same listener has to see the white hits that spend a charge`},
		{12834, "Deep Wounds", procChanceAlways, 0, `"$m1%" is the share of weapon damage, and a crit is the condition`,
			core.ProcHintCrit, `"Your critical strikes" names the trigger`},
		{16928, "Armor Shatter", procChancePPM, 0, "the 101 sentinel, answered by an override into RPPM",
			0, "the tooltip is the armor reduction, not a trigger"},
		{1308935, "Striking", procChanceColumn, 0, "no tooltip at all, so the column is all there is",
			0, "no tooltip to read a hint off"},
		{23548, "Battlegear of Wrath, the parry grant", procChanceColumn, 0, `the tooltip's "$h%" is the ProcChance column`,
			core.ProcHintOutcomeTaken, `"after a block" is an outcome no ProcTypeMask has a bit for`},
		{23547, "Battlegear of Wrath, the parry buff", procChanceAlways, 0, "the 100 in the column is the client's no-roll",
			0, "the row ships no tooltip, so the parry it is spent on stays the caller's"},
	} {
		t.Run(want.name, func(t *testing.T) {
			rows := []storeSpell{tables.row(want.id)}
			applyTooltipHints(tables, &rows[0])
			if err := applyOverrides(rows, overridesOf(want.id)); err != nil {
				t.Fatalf("baking the overrides: %v", err)
			}
			row := rows[0]

			if row.ProcChanceSource != want.source || row.ProcChanceEffect != want.effect {
				t.Errorf("spell %d reads %s effect %d, and %s, so it is %s effect %d",
					want.id, row.ProcChanceSource, row.ProcChanceEffect, want.why, want.source, want.effect)
			}

			if row.ProcHint != want.hint {
				t.Errorf("spell %d reads hint %d, and %s, so it is %d",
					want.id, row.ProcHint, want.hintWhy, want.hint)
			}

			// An effect the chance sits on has to be one the sim can reach by that number, since
			// EffectN counts positions and the client's EffectIndex has gaps.
			if row.ProcChanceEffect != 0 {
				position := int(row.ProcChanceEffect) - 1
				if position >= len(row.Effects) {
					t.Fatalf("spell %d names effect %d of %d", want.id, row.ProcChanceEffect, len(row.Effects))
				}
				if !isProcEffect(&row.Effects[position]) {
					t.Errorf("spell %d's effect %d carries aura %d, which is no proc aura",
						want.id, row.ProcChanceEffect, row.Effects[position].Aura)
				}
			}
		})
	}
}

func overridesOf(id int32) []overrides.Override {
	var mine []overrides.Override
	for _, o := range overrides.Spells {
		if o.SpellID == id {
			mine = append(mine, o)
		}
	}
	return mine
}
