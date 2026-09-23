package database

import (
	"testing"

	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/tools/database/dbc"
)

// Biznicks 247x128 Accurascope is enchant 2523, granted by 22779. Searing Totem's 3599 carries the
// same number in the misc value of an effect that is not E_ENCHANT_ITEM, and must not take its place.
func TestEnchantGrantEffectsReadOnlyEnchantItemEffects(t *testing.T) {
	effects := map[int]dbc.SpellEffect{
		1: {SpellID: 3599, EffectType: dbcenums.E_SUMMON, EffectMiscValues: []int{2523, 0}},
		2: {SpellID: 22779, EffectType: dbcenums.E_ENCHANT_ITEM, EffectMiscValues: []int{2523, 0}},
		3: {SpellID: 7418, EffectType: dbcenums.E_ENCHANT_ITEM, EffectMiscValues: []int{41, 0}},
		4: {SpellID: 3000, EffectType: dbcenums.E_ENCHANT_ITEM, EffectMiscValues: []int{41, 0}},
	}

	grants := enchantGrantEffects(effects)

	for enchantID, want := range map[int]int{2523: 22779, 41: 3000} {
		if got := grants[enchantID]; got == nil || got.SpellID != want {
			t.Errorf("enchant %d is granted by %v, want spell %d", enchantID, got, want)
		}
	}
	if len(grants) != 2 {
		t.Errorf("%d grants, want 2: %v", len(grants), grants)
	}
}
