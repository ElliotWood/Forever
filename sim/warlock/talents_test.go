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

// The Imp has no melee, so its Firebolt spends the brand, and the branded hit is Fire (1293698);
// the other demons' is Shadow (1293697).
func TestDemonicBrandImpSpendsWithFirebolt(t *testing.T) {
	for _, tc := range []struct {
		summon proto.WarlockOptions_Summon
		school core.SpellSchool
	}{{proto.WarlockOptions_Imp, core.SpellSchoolFire}, {proto.WarlockOptions_Succubus, core.SpellSchoolShadow}} {
		player := core.WithSpec(&proto.Player{
			Race:          proto.Race_RaceOrc,
			Class:         proto.Class_ClassWarlock,
			Equipment:     &proto.EquipmentSpec{},
			Consumables:   &proto.ConsumesSpec{},
			TalentsString: "-00000000000003",
			Rotation:      &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}, &proto.Player_Warlock{Warlock: &proto.Warlock{Options: &proto.Warlock_Options{ClassOptions: &proto.WarlockOptions{
			Summon: tc.summon,
		}}}})
		sim := core.NewSim(&proto.RaidSimRequest{
			SimOptions: &proto.SimOptions{RandomSeed: 100},
			Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
			Encounter:  core.MakeSingleTargetEncounter(0),
		}, simsignals.CreateSignals())
		sim.Reset()

		pet := sim.Raid.Parties[0].Players[0].(WarlockAgent).GetWarlock().ActivePet
		brand := pet.DemonicBrandAura
		brand.Activate(sim)
		brand.SetStacks(sim, brand.MaxStacks)
		for _, spell := range pet.AutoCastAbilities {
			spell.SkipCastAndApplyEffects(sim, pet.CurrentTarget)
		}

		hit := pet.GetSpell(brand.ActionID)
		if hit.SpellSchool != tc.school {
			t.Errorf("%v: branded hit is %v, want %v", tc.summon, hit.SpellSchool, tc.school)
		}
		if hit.SpellMetrics[0].Casts == 0 || brand.GetStacks() != brand.MaxStacks-hit.SpellMetrics[0].Casts {
			t.Errorf("%v: pet spell spent %d of %d charges", tc.summon, brand.MaxStacks-brand.GetStacks(), brand.MaxStacks)
		}
	}
}
