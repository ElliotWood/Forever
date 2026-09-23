package core

import (
	"math"
	"testing"

	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/core/stats"
)

func TestEnchantSpeedFollowsTheEnchantedItem(t *testing.T) {
	const counterweightAxeID, plainAxeID, glovesID, counterweightID, minorHasteID int32 = 990601, 990602, 990603, 990604, 990605

	twoHander := func(id int32) *proto.SimItem {
		return &proto.SimItem{
			Id:             id,
			Name:           "Test Speed Axe",
			Type:           proto.ItemType_ItemTypeWeapon,
			WeaponType:     proto.WeaponType_WeaponTypeAxe,
			HandType:       proto.HandType_HandTypeTwoHand,
			WeaponSpeed:    3.5,
			ScalingOptions: map[int32]*proto.ScalingItemProperties{0: {WeaponDamageMin: 100, WeaponDamageMax: 150}},
		}
	}
	haste := func(melee, ranged, spell float64) []float64 {
		pseudoStats := make([]float64, stats.PseudoStatsLen)
		pseudoStats[proto.PseudoStat_PseudoStatMeleeHastePercent] = melee
		pseudoStats[proto.PseudoStat_PseudoStatRangedHastePercent] = ranged
		pseudoStats[proto.PseudoStat_PseudoStatSpellHastePercent] = spell
		return pseudoStats
	}
	addToDatabase(&proto.SimDatabase{
		Items: []*proto.SimItem{
			twoHander(counterweightAxeID),
			twoHander(plainAxeID),
			{Id: glovesID, Name: "Test Speed Gloves", Type: proto.ItemType_ItemTypeHands, ScalingOptions: map[int32]*proto.ScalingItemProperties{0: {}}},
		},
		Enchants: []*proto.SimEnchant{
			{EffectId: counterweightID, Name: "Test Counterweight", PseudoStats: haste(3, 0, 0)},
			{EffectId: minorHasteID, Name: "Test Minor Haste", PseudoStats: haste(1, 1, 1)},
		},
	})

	withMainHand := func(mainHand *proto.ItemSpec, hands *proto.ItemSpec) []*proto.ItemSpec {
		items := make([]*proto.ItemSpec, proto.ItemSlot_ItemSlotMainHand+1)
		for i := range items {
			items[i] = &proto.ItemSpec{}
		}
		items[proto.ItemSlot_ItemSlotHands] = hands
		items[proto.ItemSlot_ItemSlotMainHand] = mainHand
		return items
	}
	counterweightAxe := &proto.ItemSpec{Id: counterweightAxeID, Enchant: counterweightID}
	plainAxe := &proto.ItemSpec{Id: plainAxeID}
	encounter := &proto.Encounter{
		Targets:  []*proto.Target{{Name: "target", Level: 63, MobType: proto.MobType_MobTypeHumanoid}},
		Duration: 180,
	}
	raid := func(mainHand, swapMainHand *proto.ItemSpec) *proto.Raid {
		return &proto.Raid{Parties: []*proto.Party{{
			Players: []*proto.Player{{
				Name:           "Speed Tester",
				Class:          proto.Class_ClassWarrior,
				Race:           proto.Race_RaceHuman,
				Buffs:          &proto.IndividualBuffs{},
				Spec:           &proto.Player_DpsWarrior{},
				Equipment:      &proto.EquipmentSpec{Items: withMainHand(mainHand, &proto.ItemSpec{Id: glovesID, Enchant: minorHasteID})},
				EnableItemSwap: true,
				ItemSwap:       &proto.ItemSwap{Items: withMainHand(swapMainHand, &proto.ItemSpec{})},
				Rotation:       &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
			}},
			Buffs: &proto.PartyBuffs{},
		}}}
	}
	newSim := func(mainHand, swapMainHand *proto.ItemSpec) (*Simulation, *Character) {
		sim := NewSim(&proto.RaidSimRequest{Raid: raid(mainHand, swapMainHand), Encounter: encounter, SimOptions: &proto.SimOptions{RandomSeed: 1}}, simsignals.CreateSignals())
		sim.Reset()
		return sim, sim.Raid.Parties[0].Players[0].GetCharacter()
	}

	check := func(when string, unit *Unit, melee float64) {
		t.Helper()
		for _, speed := range []struct {
			name      string
			got, want float64
		}{
			{"melee", unit.TotalMeleeHasteMultiplier(), melee * 1.01},
			{"ranged", unit.TotalRangedHasteMultiplier(), 1.01},
			{"cast", unit.TotalSpellHasteMultiplier(), 1.01},
		} {
			if math.Abs(speed.got-speed.want) > 1e-9 {
				t.Errorf("%s: %s speed multiplier %v, want %v", when, speed.name, speed.got, speed.want)
			}
		}
	}

	gearStats := ComputeStats(&proto.ComputeStatsRequest{Raid: raid(counterweightAxe, plainAxe), Encounter: encounter}).RaidStats.Parties[0].Players[0].GearStats
	if got, want := gearStats.PseudoStats[proto.PseudoStat_PseudoStatMeleeHastePercent], (1.03*1.01-1)*100; math.Abs(got-want) > 1e-9 {
		t.Errorf("gear melee haste %v%%, want %v%%", got, want)
	}

	sim, character := newSim(counterweightAxe, plainAxe)
	check("counterweight equipped", &character.Unit, 1.03)
	character.ItemSwap.SwapItems(sim, proto.APLActionItemSwap_Swap1, false)
	check("counterweight swapped off", &character.Unit, 1)
	sim.Cleanup()
	sim.Reset()
	check("counterweight equipped again on the next iteration", &character.Unit, 1.03)

	sim, character = newSim(plainAxe, counterweightAxe)
	check("counterweight in the swap set", &character.Unit, 1)
	character.ItemSwap.SwapItems(sim, proto.APLActionItemSwap_Swap1, false)
	check("counterweight swapped on", &character.Unit, 1.03)
	character.ItemSwap.SwapItems(sim, proto.APLActionItemSwap_Main, false)
	check("counterweight swapped back off", &character.Unit, 1)
}
