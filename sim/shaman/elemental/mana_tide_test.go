package elemental

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
)

func hasManaTide(t *testing.T, talents string) bool {
	player := core.WithSpec(&proto.Player{
		Race:          proto.Race_RaceTroll,
		Class:         proto.Class_ClassShaman,
		Equipment:     &proto.EquipmentSpec{},
		Consumables:   &proto.ConsumesSpec{},
		TalentsString: talents,
		Rotation:      &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
	}, &proto.Player_ElementalShaman{ElementalShaman: &proto.ElementalShaman{Options: &proto.ElementalShaman_Options{ClassOptions: &proto.ShamanOptions{}}}})
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 100},
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()
	return sim.Raid.Parties[0].Players[0].GetCharacter().HasAura("Mana Tide Totem (External)")
}

// The talented totem is the shaman's own, not only a party buff set by hand.
func TestManaTideTotemTalentDropsTheTotem(t *testing.T) {
	if hasManaTide(t, DefaultTalents) {
		t.Errorf("Mana Tide Totem without the talent")
	}
	if !hasManaTide(t, DefaultTalents+"001") {
		t.Errorf("no Mana Tide Totem with the talent")
	}
}

// Elemental Reach (28999) adds 3/6 yards to Lightning Bolt, Chain Lightning and Lava Burst range and
// 8/15 to Flame Shock's; it was an empty stub and the three nature/fire bolts carried no range.
func TestElementalReachExtendsRange(t *testing.T) {
	reach := "3505331000021351--553302001" // DefaultTalents with Elemental Reach 2/2 for two Convection
	for _, c := range []struct {
		talents string
		spell   int32
		want    float64
	}{
		{DefaultTalents, 15208, 30}, // Lightning Bolt
		{DefaultTalents, 29228, 20}, // Flame Shock
		{reach, 15208, 36},          // Lightning Bolt
		{reach, 10605, 36},          // Chain Lightning
		{reach, 1326047, 36},        // Lava Burst
		{reach, 29228, 35},          // Flame Shock
		{reach, 10473, 20},          // Frost Shock, not reached
	} {
		player := core.WithSpec(&proto.Player{
			Race: proto.Race_RaceTroll, Class: proto.Class_ClassShaman, Equipment: &proto.EquipmentSpec{},
			Consumables: &proto.ConsumesSpec{}, TalentsString: c.talents, Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}, &proto.Player_ElementalShaman{ElementalShaman: &proto.ElementalShaman{Options: &proto.ElementalShaman_Options{ClassOptions: &proto.ShamanOptions{}}}})
		sim := core.NewSim(&proto.RaidSimRequest{
			SimOptions: &proto.SimOptions{RandomSeed: 100},
			Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
			Encounter:  core.MakeSingleTargetEncounter(0),
		}, simsignals.CreateSignals())
		sim.Reset()
		if got := sim.Raid.Parties[0].Players[0].GetCharacter().GetSpell(core.ActionID{SpellID: c.spell}).MaxRange; got != c.want {
			t.Errorf("%q: spell %d reaches %.1f yd, want %.1f", c.talents, c.spell, got, c.want)
		}
	}
}

// Call of Flame 3/3 (in DefaultTalents) adds 9 sec to Flame Shock's 12, three more 3 second ticks.
func TestCallOfFlameExtendsFlameShock(t *testing.T) {
	for _, c := range []struct {
		talents string
		want    int32
	}{
		{"", 4},
		{DefaultTalents, 7},
	} {
		player := core.WithSpec(&proto.Player{
			Race: proto.Race_RaceTroll, Class: proto.Class_ClassShaman, Equipment: &proto.EquipmentSpec{},
			Consumables: &proto.ConsumesSpec{}, TalentsString: c.talents, Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}, &proto.Player_ElementalShaman{ElementalShaman: &proto.ElementalShaman{Options: &proto.ElementalShaman_Options{ClassOptions: &proto.ShamanOptions{}}}})
		sim := core.NewSim(&proto.RaidSimRequest{
			SimOptions: &proto.SimOptions{RandomSeed: 100},
			Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
			Encounter:  core.MakeSingleTargetEncounter(0),
		}, simsignals.CreateSignals())
		sim.Reset()
		dot := sim.Raid.Parties[0].Players[0].GetCharacter().GetSpell(core.ActionID{SpellID: 29228, Tag: 1}).Dot(sim.Encounter.AllTargetUnits[0])
		if dot.BaseTickCount != c.want {
			t.Errorf("%q: Flame Shock ticks %d times, want %d", c.talents, dot.BaseTickCount, c.want)
		}
	}
}

// Trial by Fire 2/2 adds 20 sec to Magma Totem's 20, ten more 2 second pulses.
func TestTrialByFireExtendsMagmaTotem(t *testing.T) {
	for _, c := range []struct {
		talents string
		want    int32
	}{
		{DefaultTalents, 10},
		{"3505331020001351--553302001", 20}, // DefaultTalents with Trial by Fire 2/2 for two Convection
	} {
		player := core.WithSpec(&proto.Player{
			Race: proto.Race_RaceTroll, Class: proto.Class_ClassShaman, Equipment: &proto.EquipmentSpec{},
			Consumables: &proto.ConsumesSpec{}, TalentsString: c.talents, Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}, &proto.Player_ElementalShaman{ElementalShaman: &proto.ElementalShaman{Options: &proto.ElementalShaman_Options{ClassOptions: &proto.ShamanOptions{}}}})
		sim := core.NewSim(&proto.RaidSimRequest{
			SimOptions: &proto.SimOptions{RandomSeed: 100},
			Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
			Encounter:  core.MakeSingleTargetEncounter(0),
		}, simsignals.CreateSignals())
		sim.Reset()
		if got := sim.Raid.Parties[0].Players[0].GetCharacter().GetSpell(core.ActionID{SpellID: 10587}).AOEDot().BaseTickCount; got != c.want {
			t.Errorf("%q: Magma Totem pulses %d times, want %d", c.talents, got, c.want)
		}
	}
}
