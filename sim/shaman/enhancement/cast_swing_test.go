package enhancement

import (
	"testing"
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/shaman"
)

// A hard-cast Lightning Bolt restarts the swing timer when it completes, even when the next swing
// was not due until after it: beta logs 2701 and 2712 (foreverlogs.gg) show the next swing one
// weapon speed after every completed hard cast. An instant bolt (5 Maelstrom stacks) leaves it alone.
func TestHardCastRestartsTheSwingTimer(t *testing.T) {
	nextSwingAfterBolt := func(maelstrom int32) (next, speed time.Duration) {
		sim := core.NewSim(&proto.RaidSimRequest{
			SimOptions: &proto.SimOptions{RandomSeed: 1},
			Raid: core.SinglePlayerRaidProto(&proto.Player{
				Name: "enh", Class: proto.Class_ClassShaman, Race: proto.Race_RaceOrc, TalentsString: DefaultTalents,
				Equipment: &proto.EquipmentSpec{Items: []*proto.ItemSpec{
					{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {},
					{Id: 17182}, // Sulfuras, Hand of Ragnaros
				}},
				Spec:     &proto.Player_EnhancementShaman{EnhancementShaman: &proto.EnhancementShaman{Options: &proto.EnhancementShaman_Options{ClassOptions: &proto.ShamanOptions{}}}},
				Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
			}, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
			Encounter: core.MakeSingleTargetEncounter(0),
		}, simsignals.CreateSignals())
		sim.Reset()
		sham := sim.Raid.Parties[0].Players[0].(shaman.ShamanAgent).GetShaman()
		sham.AutoAttacks.EnableAutoSwing(sim)
		sham.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime) // a swing just landed
		if maelstrom > 0 {
			mw := sham.GetAura("Maelstrom Weapon")
			mw.Activate(sim)
			mw.SetStacks(sim, maelstrom)
		}
		if !sham.GetSpell(core.ActionID{SpellID: 403}).Cast(sim, sham.CurrentTarget) { // rank 1, 1.5 sec
			t.Fatal("Lightning Bolt did not cast")
		}
		return sham.AutoAttacks.NextAttackAt() - sim.CurrentTime, sham.AutoAttacks.MainhandSwingSpeed()
	}

	if got, speed := nextSwingAfterBolt(0); got != 1500*time.Millisecond+speed {
		t.Errorf("hard cast: next swing %s from the cast start, want %s (cast end + a full swing)", got, 1500*time.Millisecond+speed)
	}
	if got, speed := nextSwingAfterBolt(5); got != speed {
		t.Errorf("instant bolt: next swing %s from the cast, want %s (untouched)", got, speed)
	}
}
