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
		{923, "Enchant Bracer - Deflection", stats.Stats{stats.DefenseRating: 7}},
		{2583, "Presence of Might", stats.Stats{stats.Stamina: 10, stats.DefenseRating: 7, stats.BlockValue: 15}},
		{7633, "Presence of Valor", stats.Stats{stats.Stamina: 20, stats.DefenseRating: 7, stats.BlockValue: 15}},
		{8214, "Enchant Bracer - Superior Deflection", stats.Stats{stats.DefenseRating: 9}},
		{8719, "Wild Leather Armor Kit", stats.Stats{stats.DefenseRating: 4, stats.Stamina: 10}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := stats.FromProtoArray(enchantByName(t, tc.effectID, tc.name).Stats)
			if got != tc.want {
				t.Errorf("stats %v, want %v", got.FlatString(), tc.want.FlatString())
			}
		})
	}
}

func TestEnchantWeaponDamage(t *testing.T) {
	inRepositoryRoot(t)

	for _, tc := range []struct {
		effectID int
		name     string
		want     float64
	}{
		{963, "Enchant 2H Weapon - Greater Impact", 7},
		{1897, "Enchant Weapon - Superior Striking", 5},
		{664, "Sniper Scope", 7},
		{7944, "Poultry Precision Scope", 10},
		{7655, "Enchant Bracer - Spell Power", 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			enchant := enchantByName(t, tc.effectID, tc.name)
			if enchant.WeaponDamage != tc.want {
				t.Errorf("weapon damage %v, want %v", enchant.WeaponDamage, tc.want)
			}
			if tc.want != 0 && stats.FromProtoArray(enchant.Stats) != (stats.Stats{}) {
				t.Errorf("stats %v, want none", enchant.Stats)
			}
		})
	}
}

func TestEnchantSlot(t *testing.T) {
	inRepositoryRoot(t)

	for _, tc := range []struct {
		effectID    int
		name        string
		itemType    proto.ItemType
		enchantType proto.EnchantType
	}{
		{7663, "Enchant Shield - Excellent Stamina", proto.ItemType_ItemTypeWeapon, proto.EnchantType_EnchantTypeShield},
		{7664, "Enchant Shield - Critical Strike", proto.ItemType_ItemTypeWeapon, proto.EnchantType_EnchantTypeShield},
		{255, "Enchant Shield - Lesser Spirit", proto.ItemType_ItemTypeWeapon, proto.EnchantType_EnchantTypeShield},
		{848, "Enchant Shield - Lesser Protection", proto.ItemType_ItemTypeWeapon, proto.EnchantType_EnchantTypeShield},
		{848, "Enchant Cloak - Defense", proto.ItemType_ItemTypeBack, proto.EnchantType_EnchantTypeNormal},
		{929, "Enchant Shield - Stamina", proto.ItemType_ItemTypeWeapon, proto.EnchantType_EnchantTypeShield},
		{7659, "Enchant Off-Hand - Superior Intellect", proto.ItemType_ItemTypeWeapon, proto.EnchantType_EnchantTypeOffHand},
		{7660, "Enchant Off-Hand - Excellent Spirit", proto.ItemType_ItemTypeWeapon, proto.EnchantType_EnchantTypeOffHand},
		{7661, "Enchant Off-Hand - Wisdom", proto.ItemType_ItemTypeWeapon, proto.EnchantType_EnchantTypeOffHand},
	} {
		t.Run(tc.name, func(t *testing.T) {
			enchant := enchantByName(t, tc.effectID, tc.name)
			if enchant.Type != tc.itemType || enchant.EnchantType != tc.enchantType || len(enchant.ExtraTypes) != 0 {
				t.Errorf("type %v, enchant type %v, extra types %v; want %v, %v, none",
					enchant.Type, enchant.EnchantType, enchant.ExtraTypes, tc.itemType, tc.enchantType)
			}
		})
	}
}

func TestEquipSpellStats(t *testing.T) {
	inRepositoryRoot(t)
	dbc.GetDBC()

	for _, tc := range []struct {
		name    string
		spellID int
		want    stats.Stats
	}{
		{"13922 Enchant Bracer - Deflection", 13922, stats.Stats{stats.DefenseRating: 7}},
		{"24148 Presence of Might", 24148, stats.Stats{stats.DefenseRating: 7, stats.BlockValue: 15}},
		{"7778 Enchant Boots - Minor Agility", 7778, stats.Stats{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := stats.Stats{}
			added := dbc.AddEquipSpellStats(&got, tc.spellID)
			if got != tc.want {
				t.Errorf("stats %v, want %v", got.FlatString(), tc.want.FlatString())
			}
			if added != (tc.want != stats.Stats{}) {
				t.Errorf("reported added=%v", added)
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
