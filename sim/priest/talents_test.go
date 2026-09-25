package priest

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
)

func penanceCost(t *testing.T, talents string) float64 {
	player := core.WithSpec(&proto.Player{
		Race:          proto.Race_RaceUndead,
		Class:         proto.Class_ClassPriest,
		Equipment:     &proto.EquipmentSpec{},
		Consumables:   &proto.ConsumesSpec{},
		TalentsString: talents,
		Rotation:      &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
	}, &proto.Player_DpsPriest{DpsPriest: &proto.DpsPriest{Options: &proto.DpsPriest_Options{ClassOptions: &proto.PriestOptions{}}}})
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 100},
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()
	priest := sim.Raid.Parties[0].Players[0].(PriestAgent).GetPriest()
	penance := priest.GetSpell(core.ActionID{SpellID: spellData.Penance.Highest().ID})
	if penance == nil {
		t.Fatalf("%s: no Penance", talents)
	}
	return penance.Cost.GetCurrentCost()
}

// Improved Healing 3/3 takes 15% off Penance.
func TestImprovedHealingDiscountsPenance(t *testing.T) {
	base := penanceCost(t, "504020031305001-13505100202-50002")
	discounted := penanceCost(t, "504020031305001-13505100032-50002")
	if !core.WithinToleranceFloat64(base*0.85, discounted, 1e-6) {
		t.Errorf("Penance costs %v with Improved Healing 3, %v without; want 15%% off", discounted, base)
	}
}
