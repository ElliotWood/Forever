package elemental

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
)

// Since the 16 Oct build a Water Shield globe also goes when the shaman's own spell crits, so a
// shaman nobody hits still gets mana from it: one globe per 3.5 sec at most, three in all. Rank 1
// Lightning Bolt so a naked shaman keeps casting for the whole fight.
func TestWaterShieldRestoresManaOnSpellCrits(t *testing.T) {
	cast := func(id int32) *proto.APLAction {
		return &proto.APLAction{Action: &proto.APLAction_CastSpell{
			CastSpell: &proto.APLActionCastSpell{SpellId: &proto.ActionID{RawId: &proto.ActionID_SpellId{SpellId: id}}},
		}}
	}
	player := &proto.Player{
		Name: "ele", Class: proto.Class_ClassShaman, Race: proto.Race_RaceOrc, TalentsString: DefaultTalents,
		Equipment: &proto.EquipmentSpec{},
		Spec:      &proto.Player_ElementalShaman{ElementalShaman: &proto.ElementalShaman{Options: &proto.ElementalShaman_Options{ClassOptions: &proto.ShamanOptions{}}}},
		Rotation: &proto.APLRotation{
			Type:           proto.APLRotation_TypeAPL,
			PrepullActions: []*proto.APLPrepullAction{{Action: cast(408510), DoAtValue: &proto.APLValue{Value: &proto.APLValue_Const{Const: &proto.APLValueConst{Val: "-2s"}}}}},
			PriorityList:   []*proto.APLListItem{{Action: cast(403)}},
		},
	}
	encounter := core.MakeSingleTargetEncounter(0)
	encounter.Duration = 300
	res := core.RunRaidSim(&proto.RaidSimRequest{
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  encounter,
		SimOptions: &proto.SimOptions{Iterations: 20, RandomSeed: 1},
	})
	if res.Error != nil {
		t.Fatal(res.Error.Message)
	}
	globes := 0.0
	for _, r := range res.RaidMetrics.Parties[0].Players[0].Resources {
		if r.Id.GetSpellId() == 408510 {
			globes = float64(r.Events) / 20
		}
	}
	if globes < 2.5 || globes > 3 {
		t.Errorf("Water Shield globes spent on crits in 5 min: %v a fight, want 2.5-3", globes)
	}
}
