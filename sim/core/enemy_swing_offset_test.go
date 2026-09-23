package core

import (
	"slices"
	"testing"
	"time"

	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
)

func TestEnemiesOpenAtTheirOwnPointInTheSwingTimer(t *testing.T) {
	targets := make([]*proto.Target, 5)
	for i := range targets {
		targets[i] = FreshDefaultTargetConfig()
	}
	firstSwings := func() []time.Duration {
		sim := NewSim(&proto.RaidSimRequest{
			SimOptions: &proto.SimOptions{RandomSeed: 100},
			Raid: &proto.Raid{
				Parties: []*proto.Party{{
					Players: []*proto.Player{{
						Name:      "Warrior",
						Class:     proto.Class_ClassWarrior,
						Buffs:     &proto.IndividualBuffs{},
						Spec:      &proto.Player_ProtectionWarrior{},
						Equipment: &proto.EquipmentSpec{},
						Rotation:  &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
					}},
					Buffs: &proto.PartyBuffs{},
				}},
			},
			Encounter: &proto.Encounter{Targets: targets, Duration: 180},
		}, simsignals.CreateSignals())
		sim.Reseed(7)
		sim.Reset()
		var swings []time.Duration
		for _, target := range sim.Encounter.AllTargets {
			aa := &target.AutoAttacks
			if aa.mh.swingAt < 0 || aa.mh.swingAt >= aa.MainhandSwingSpeed() {
				t.Fatalf("%s opens at %v, want within its %v swing", target.Label, aa.mh.swingAt, aa.MainhandSwingSpeed())
			}
			if aa.mh.previousSwing != aa.mh.swingAt-aa.MainhandSwingSpeed() {
				t.Fatalf("%s's previous swing is %v, want a swing before %v", target.Label, aa.mh.previousSwing, aa.mh.swingAt)
			}
			swings = append(swings, aa.mh.swingAt)
		}
		return swings
	}

	swings := firstSwings()
	if slices.Min(swings) == slices.Max(swings) {
		t.Fatalf("every enemy opens at %v", swings[0])
	}
	if again := firstSwings(); !slices.Equal(swings, again) {
		t.Fatalf("the same seed opened at %v, then at %v", swings, again)
	}
}
