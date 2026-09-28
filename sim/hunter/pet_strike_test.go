package hunter

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
)

// The Crab casts Pinch, the Crocolisk Dismember, the Owl Mine! and the Hyena Tendon Rip (client rows
// 1264742 / 1264933 / 1265058 / 1265042), each landing inside its rank 5 range.
func TestPetStrikes(t *testing.T) {
	for _, c := range []struct {
		pet      proto.HunterOptions_PetType
		id       int32
		min, max float64
	}{
		{proto.HunterOptions_Crab, spellData.PinchTriggered.Highest().ID, 88, 102},
		{proto.HunterOptions_Crocolisk, spellData.DismemberTriggered.Highest().ID, 50, 58},
		{proto.HunterOptions_Owl, spellData.MineTriggered.Highest().ID, 41, 47},
		// The whole bleed, 3 ticks of 20.
		{proto.HunterOptions_Hyena, spellData.TendonRipTriggered.Highest().ID, 60, 60},
	} {
		player := &proto.Player{
			Name: "bm", Class: proto.Class_ClassHunter, Race: proto.Race_RaceOrc, TalentsString: BeastMasteryTalents,
			Equipment: WeaponsOnly, DistanceFromTarget: 30,
			Spec: &proto.Player_Hunter{Hunter: &proto.Hunter{Options: &proto.Hunter_Options{ClassOptions: &proto.HunterOptions{
				Ammo: proto.HunterOptions_Doomshot, QuiverBonus: proto.HunterOptions_Speed15, PetType: c.pet,
				PetAttackSpeed: proto.HunterOptions_OneTwo, PetUptime: 1}}}},
			Rotation: core.GetAplRotation("../../ui/specs/hunter/dps/apls", "bm").Rotation,
		}
		raid := &proto.Raid{Parties: []*proto.Party{{Players: []*proto.Player{player}, Buffs: &proto.PartyBuffs{}}}, Buffs: &proto.RaidBuffs{}, Debuffs: &proto.Debuffs{}, NumActiveParties: 1}
		// A target with no armor, so a hit shows the row's own damage.
		encounter := core.MakeSingleTargetEncounter(0)
		encounter.Targets[0].Stats[proto.Stat_StatArmor] = 0
		res := core.RunRaidSim(&proto.RaidSimRequest{Raid: raid, Encounter: encounter, SimOptions: &proto.SimOptions{Iterations: 1, RandomSeed: 1}})
		if res.Error != nil {
			t.Fatal(res.Error.Message)
		}

		var hits int32
		var damage float64
		for _, pet := range res.RaidMetrics.Parties[0].Players[0].Pets {
			for _, action := range pet.Actions {
				if action.Id.GetSpellId() != c.id {
					continue
				}
				for _, target := range action.Targets {
					hits += target.Hits
					damage += target.Damage - target.CritDamage
				}
			}
		}
		if hits == 0 {
			t.Fatalf("%v: no hit from spell %d", c.pet, c.id)
		}
		avg := damage / float64(hits)
		t.Logf("%v: %d hits, %.1f a hit", c.pet, hits, avg)
		if avg < c.min || avg > c.max*2 {
			t.Fatalf("%v: spell %d averaged %.1f a hit, want the row's %.0f-%.0f before pet multipliers", c.pet, c.id, avg, c.min, c.max)
		}
	}
}
