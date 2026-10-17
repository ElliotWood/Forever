package elemental

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
)

// Elemental Mastery (16166) takes Chain Lightning's 6 s cooldown away for 15 s, so a 15 s fight of
// Elemental Mastery then Chain Lightning casts it back to back (2 s casts) instead of every 6 s.
func TestElementalMasteryRemovesChainLightningCooldown(t *testing.T) {
	casts := func(talents string) float64 {
		cast := func(id int32) *proto.APLListItem {
			return &proto.APLListItem{Action: &proto.APLAction{Action: &proto.APLAction_CastSpell{
				CastSpell: &proto.APLActionCastSpell{SpellId: &proto.ActionID{RawId: &proto.ActionID_SpellId{SpellId: id}}},
			}}}
		}
		player := &proto.Player{
			Name: "ele", Class: proto.Class_ClassShaman, Race: proto.Race_RaceOrc, TalentsString: talents,
			Equipment: &proto.EquipmentSpec{},
			Spec:      &proto.Player_ElementalShaman{ElementalShaman: &proto.ElementalShaman{Options: &proto.ElementalShaman_Options{ClassOptions: &proto.ShamanOptions{}}}},
			Rotation:  &proto.APLRotation{Type: proto.APLRotation_TypeAPL, PriorityList: []*proto.APLListItem{cast(16166), cast(10605)}},
		}
		encounter := core.MakeSingleTargetEncounter(0)
		encounter.Duration = 15
		res := core.RunRaidSim(&proto.RaidSimRequest{
			Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
			Encounter:  encounter,
			SimOptions: &proto.SimOptions{Iterations: 1, RandomSeed: 1},
		})
		if res.Error != nil {
			t.Fatal(res.Error.Message)
		}
		for _, a := range res.RaidMetrics.Parties[0].Players[0].Actions {
			if a.Id.GetSpellId() == 10605 && a.Id.Tag == 0 { // tags 1-6 are its overloads
				n := int32(0)
				for _, tgt := range a.Targets {
					n += tgt.Casts
				}
				return float64(n)
			}
		}
		return 0
	}

	without, with := casts("5505331000001350--553302001"), casts(DefaultTalents)
	if without > 3 || with < 6 {
		t.Errorf("Chain Lightning casts in 15 s: %v without Elemental Mastery (want <= 3), %v with (want >= 6)", without, with)
	}
}
