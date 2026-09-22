package core

import (
	"testing"

	"github.com/wowsims/forever/sim/core/proto"
)

// A fixed chance bound to a weapon enchant rolls on the hits of the hand carrying the enchant, and
// moves with it when an item swap hands the enchanted weapon to the other hand.
func TestEnchantProcChanceFollowsTheEnchantedWeapon(t *testing.T) {
	const enchantedID, plainID, enchantID int32 = 990501, 990502, 990503

	oneHander := func(id int32) *proto.SimItem {
		return &proto.SimItem{
			Id:             id,
			Name:           "Test Sword",
			Type:           proto.ItemType_ItemTypeWeapon,
			WeaponType:     proto.WeaponType_WeaponTypeSword,
			HandType:       proto.HandType_HandTypeOneHand,
			WeaponSpeed:    2.6,
			ScalingOptions: map[int32]*proto.ScalingItemProperties{0: {WeaponDamageMin: 50, WeaponDamageMax: 90}},
		}
	}
	addToDatabase(&proto.SimDatabase{
		Items:    []*proto.SimItem{oneHander(enchantedID), oneHander(plainID)},
		Enchants: []*proto.SimEnchant{{EffectId: enchantID, Name: "Test Fiery Blaze"}},
	})

	hands := func(mainHand, offHand *proto.ItemSpec) []*proto.ItemSpec {
		items := make([]*proto.ItemSpec, proto.ItemSlot_ItemSlotOffHand+1)
		for i := range items {
			items[i] = &proto.ItemSpec{}
		}
		items[proto.ItemSlot_ItemSlotMainHand] = mainHand
		items[proto.ItemSlot_ItemSlotOffHand] = offHand
		return items
	}
	enchanted := &proto.ItemSpec{Id: enchantedID, Enchant: enchantID}
	plain := &proto.ItemSpec{Id: plainID}

	agent := &FakeAgent{Character: NewCharacter(&Party{}, 0, &proto.Player{
		Name:           "Swap Tester",
		Class:          proto.Class_ClassWarrior,
		Race:           proto.Race_RaceHuman,
		Spec:           &proto.Player_ProtectionWarrior{ProtectionWarrior: &proto.ProtectionWarrior{}},
		Equipment:      &proto.EquipmentSpec{Items: hands(enchanted, plain)},
		EnableItemSwap: true,
		ItemSwap:       &proto.ItemSwap{Items: hands(plain, enchanted)},
	})}
	character := &agent.Character
	character.ItemSwap.initialize(character)
	character.EnableAutoAttacks(agent, AutoAttackOptions{
		MainHand:       character.WeaponFromMainHand(),
		OffHand:        character.WeaponFromOffHand(),
		AutoSwingMelee: true,
	})

	dpm := character.NewDynamicLegacyProcForEnchant(enchantID, 0, 0.15)

	check := func(when string, mainHand, offHand float64) {
		t.Helper()
		if got := dpm.Chance(ProcMaskMeleeMHAuto, nil); got != mainHand {
			t.Errorf("%s: main hand chance = %v, want %v", when, got, mainHand)
		}
		if got := dpm.Chance(ProcMaskMeleeOHAuto, nil); got != offHand {
			t.Errorf("%s: off hand chance = %v, want %v", when, got, offHand)
		}
	}
	check("enchant on the main hand", 0.15, 0)

	// SwapItems' per-slot step. The rest of it moves stats and activates auras, which needs a running
	// sim, and a prepull swap exchanges the items without one.
	for _, slot := range character.ItemSwap.slots {
		character.ItemSwap.swapItem(nil, slot, true, false)
		for _, onSwap := range character.ItemSwap.onItemSwapCallbacks[slot] {
			onSwap(nil, slot)
		}
	}
	if character.OffHand().Enchant.EffectID != enchantID {
		t.Fatal("the swap did not move the enchanted weapon to the off hand")
	}
	check("enchant swapped to the off hand", 0, 0.15)
}
