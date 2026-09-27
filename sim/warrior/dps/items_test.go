package dps

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
)

// Diamond Flask's Strength comes from finishing its 5 sec channel, so cast at the pull it holds
// for a full minute from the 5 sec mark.
func TestDiamondFlask(t *testing.T) {
	equipment := weaponsOnly(15240, 15238)
	equipment.Items[proto.ItemSlot_ItemSlotTrinket1] = &proto.ItemSpec{Id: 20130}
	rotation := core.GetAplRotation("../../../ui/specs/warrior/dps/apls", "dps_reck").Rotation
	rotation.PriorityList = append([]*proto.APLListItem{{Action: &proto.APLAction{Action: &proto.APLAction_CastSpell{
		CastSpell: &proto.APLActionCastSpell{SpellId: &proto.ActionID{RawId: &proto.ActionID_ItemId{ItemId: 20130}}},
	}}}}, rotation.PriorityList...)

	player := core.WithSpec(&proto.Player{
		Race:          proto.Race_RaceOrc,
		Class:         proto.Class_ClassWarrior,
		Equipment:     equipment,
		Consumables:   &proto.ConsumesSpec{},
		TalentsString: DpsTalents,
		Rotation:      rotation,
	}, DefaultOptions)
	raid := core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{})
	result := core.RunRaidSim(&proto.RaidSimRequest{Raid: raid, Encounter: core.MakeSingleTargetEncounter(0),
		SimOptions: &proto.SimOptions{Iterations: 5, RandomSeed: 101}})
	if result.Error != nil {
		t.Fatal(result.Error.Message)
	}
	for _, aura := range result.RaidMetrics.Parties[0].Players[0].Auras {
		if aura.Id.GetSpellId() == 1318070 {
			if aura.UptimeSecondsAvg != 60 {
				t.Errorf("Diamond Flask's Strength up %v sec, want 60", aura.UptimeSecondsAvg)
			}
			return
		}
	}
	t.Error("Diamond Flask never granted its Strength")
}

// Shadowstrike, The Cruel Hand of Timmy and Skullforge Reaver drain life: a health leech (the first
// two) or a damage aura (Skullforge Brand). Each proc has to land its damage on the target.
func TestLifeDrainWeaponProcs(t *testing.T) {
	for _, weapon := range []struct {
		itemID, spellID int32
	}{{17074, 21170}, {13401, 17505}, {13361, 17484}} {
		player := core.WithSpec(&proto.Player{
			Race:          proto.Race_RaceOrc,
			Class:         proto.Class_ClassWarrior,
			Equipment:     weaponsOnly(weapon.itemID, 0),
			Consumables:   &proto.ConsumesSpec{},
			TalentsString: DpsTalents,
			Rotation:      core.GetAplRotation("../../../ui/specs/warrior/dps/apls", "dps_reck").Rotation,
		}, DefaultOptions)
		raid := core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{})
		result := core.RunRaidSim(&proto.RaidSimRequest{Raid: raid, Encounter: core.MakeSingleTargetEncounter(0),
			SimOptions: &proto.SimOptions{Iterations: 5, RandomSeed: 101}})
		if result.Error != nil {
			t.Fatal(result.Error.Message)
		}
		damage := 0.0
		for _, action := range result.RaidMetrics.Parties[0].Players[0].Actions {
			if action.Id.GetSpellId() == weapon.spellID {
				for _, target := range action.Targets {
					damage += target.Damage
				}
			}
		}
		if damage <= 0 {
			t.Errorf("item %d: proc %d dealt no damage", weapon.itemID, weapon.spellID)
		}
	}
}
