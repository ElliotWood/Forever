package warlock

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
)

// The Demonic Pact 2/31/18 preset: Demonic Sacrifice, Soul Link and Demonic Pact.
const pactTalents = "113-0005003221220311351-0550005"

func sacrificeAura(t *testing.T, talents string, pact proto.WarlockOptions_Summon) *core.Aura {
	player := core.WithSpec(&proto.Player{
		Race:          proto.Race_RaceOrc,
		Class:         proto.Class_ClassWarlock,
		Equipment:     &proto.EquipmentSpec{},
		Consumables:   &proto.ConsumesSpec{},
		TalentsString: talents,
		Rotation:      &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
	}, &proto.Player_Warlock{Warlock: &proto.Warlock{Options: &proto.Warlock_Options{ClassOptions: &proto.WarlockOptions{
		Summon:        proto.WarlockOptions_Succubus,
		PactSacrifice: pact,
	}}}})
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 100},
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()
	return sim.Raid.Parties[0].Players[0].(WarlockAgent).GetWarlock().GetAura("Demonic Sacrifice")
}

// Demonic Pact keeps a sacrificed Imp's Shadow buff with the Succubus out, but not a sacrificed
// Succubus (resummoning the sacrificed demon cancels it), and nothing without the talent.
func TestDemonicPactKeepsTheSacrifice(t *testing.T) {
	if aura := sacrificeAura(t, pactTalents, proto.WarlockOptions_Imp); aura == nil || aura.ActionID.SpellID != 18789 {
		t.Errorf("Pact + sacrificed Imp: got %v, want Touch of Shadow (18789)", aura)
	}
	if aura := sacrificeAura(t, pactTalents, proto.WarlockOptions_Succubus); aura != nil {
		t.Errorf("Pact + sacrificed Succubus with the Succubus out: got %v, want none", aura)
	}
	if aura := sacrificeAura(t, "113-0005003221220311350-0550005", proto.WarlockOptions_Imp); aura != nil {
		t.Errorf("no Demonic Pact: got %v, want none", aura)
	}
}
