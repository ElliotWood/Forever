package enhancement

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
)

// Flametongue Totem (16387) adds 16389's fire hit to each landed main-hand auto. It no longer stacks
// with Windfury Totem or Flametongue Weapon (Forever beta development notes), and a main-hand
// Flametongue Weapon disables it (tooltips 8024/16342).
func TestFlametongueTotem(t *testing.T) {
	run := func(mh proto.ShamanImbue, party *proto.PartyBuffs) (autos, hits int32) {
		player := &proto.Player{
			Name: "enh", Class: proto.Class_ClassShaman, Race: proto.Race_RaceOrc, TalentsString: DefaultTalents,
			Equipment: &proto.EquipmentSpec{Items: []*proto.ItemSpec{
				{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {},
				{Id: 17182}, // Sulfuras, Hand of Ragnaros
			}},
			Spec: &proto.Player_EnhancementShaman{EnhancementShaman: &proto.EnhancementShaman{Options: &proto.EnhancementShaman_Options{
				ClassOptions: &proto.ShamanOptions{ImbueMh: mh},
			}}},
			Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL, PriorityList: []*proto.APLListItem{{Action: &proto.APLAction{Action: &proto.APLAction_CastSpell{
				CastSpell: &proto.APLActionCastSpell{SpellId: &proto.ActionID{RawId: &proto.ActionID_SpellId{SpellId: 16387}}},
			}}}}},
		}
		raid := &proto.Raid{Parties: []*proto.Party{{Players: []*proto.Player{player}, Buffs: party}}, Buffs: &proto.RaidBuffs{}, Debuffs: &proto.Debuffs{}, NumActiveParties: 1}
		res := core.RunRaidSim(&proto.RaidSimRequest{Raid: raid, Encounter: core.MakeSingleTargetEncounter(0), SimOptions: &proto.SimOptions{Iterations: 5, RandomSeed: 1}})
		if res.Error != nil {
			t.Fatal(res.Error.Message)
		}
		for _, a := range res.RaidMetrics.Parties[0].Players[0].Actions {
			for _, tgt := range a.Targets {
				if a.Id.GetOtherId() == proto.OtherAction_OtherActionAttack && a.Id.Tag == 1 {
					autos += tgt.Hits + tgt.Crits + tgt.Glances + tgt.Blocks
				}
				if a.Id.GetSpellId() == 16389 {
					hits += tgt.Hits + tgt.Crits + tgt.Misses
				}
			}
		}
		return
	}

	if autos, hits := run(proto.ShamanImbue_NoImbue, &proto.PartyBuffs{}); hits == 0 || hits != autos {
		t.Errorf("alone: %d totem hits for %d landed main-hand autos, want one each", hits, autos)
	}
	if autos, hits := run(proto.ShamanImbue_FrostbrandWeapon, &proto.PartyBuffs{}); hits == 0 || hits != autos {
		t.Errorf("Frostbrand main hand: %d totem hits for %d landed main-hand autos, want one each", hits, autos)
	}
	if _, hits := run(proto.ShamanImbue_NoImbue, &proto.PartyBuffs{WindfuryTotem: true}); hits != 0 {
		t.Errorf("with the party's Windfury Totem: %d totem hits, want 0", hits)
	}
	if _, hits := run(proto.ShamanImbue_FlametongueWeapon, &proto.PartyBuffs{}); hits != 0 {
		t.Errorf("Flametongue main hand: %d totem hits, want 0", hits)
	}
}
