package database

import (
	"testing"

	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
	"github.com/wowsims/forever/tools/database/dbc"
)

func enchantByName(t *testing.T, effectID int, name string) *proto.UIEnchant {
	t.Helper()
	for _, enchant := range dbc.GetDBC().Enchants {
		if enchant.EffectId == effectID && enchant.Name == name {
			return enchant.ToProto()
		}
	}
	t.Fatalf("enchant %d %q is not in the enchant inputs", effectID, name)
	return nil
}

func TestEnchantStats(t *testing.T) {
	inRepositoryRoot(t)

	for _, tc := range []struct {
		effectID int
		name     string
		want     stats.Stats
	}{
		{7655, "Enchant Bracer - Spell Power", stats.Stats{stats.SpellDamage: 12, stats.HealingPower: 12}},
		{8211, "Enchant 2H Weapon - Mighty Spell Power", stats.Stats{stats.SpellDamage: 55, stats.HealingPower: 55}},
		{8482, "Mystic Medium Armor Kit", stats.Stats{stats.BonusArmor: 16, stats.SpellDamage: 2, stats.HealingPower: 2}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := stats.FromProtoArray(enchantByName(t, tc.effectID, tc.name).Stats)
			if got != tc.want {
				t.Errorf("stats %v, want %v", got.FlatString(), tc.want.FlatString())
			}
		})
	}
}

// Healing items state spell power (45) and their extra healing (41) in separate slots, and the
// two add up: 231622 Field Marshal's Satin Hood.
func TestItemSpellPowerAddsToHealing(t *testing.T) {
	inRepositoryRoot(t)

	item, ok := dbc.GetDBC().Items[231622]
	if !ok {
		t.Fatal("item 231622 is not in the item inputs")
	}
	got := item.GetStats(item.ItemLevel)
	if got[stats.SpellDamage] != 21 || got[stats.HealingPower] != 62 {
		t.Errorf("spell damage %v and healing %v, want 21 and 62", got[stats.SpellDamage], got[stats.HealingPower])
	}
}
