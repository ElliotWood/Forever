package dpswarrior

import (
	"testing"

	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
)

// The Skyborne are a Forever race, so their racials must not reach a Classic sim. The race
// picker offers whatever the spec allows without filtering by ruleset, so a saved setting
// can name a Skyborne under Classic, and nothing else stops those racials from applying.
func TestSkyborneRacialsOnlyApplyUnderForever(t *testing.T) {
	windBlessed := func(ruleset proto.Ruleset) float64 {
		result := core.ComputeStats(&proto.ComputeStatsRequest{
			Raid: core.SinglePlayerRaidProto(&proto.Player{
				Class:         proto.Class_ClassWarrior,
				Race:          proto.Race_RaceSkyborneWindshaper,
				Equipment:     core.GetGearSet("../../../ui/warrior/gear_sets", "p0.bis").GearSet,
				TalentsString: P1Talents,
				Spec:          PlayerOptionsArms,
			}, nil, nil, nil),
			Ruleset: ruleset,
		})

		stats := result.RaidStats.Parties[0].Players[0].FinalStats.PseudoStats
		return stats[proto.PseudoStat_PseudoStatMeleeSpeedMultiplier]
	}

	// Wind Blessed is 1% melee, ranged and cast haste; the multiplier is the cheapest of
	// the three to read back.
	if forever := windBlessed(proto.Ruleset_RulesetForever); forever <= 1 {
		t.Errorf("Forever melee speed multiplier = %v, want the 1%% from Wind Blessed", forever)
	}
	if classic := windBlessed(proto.Ruleset_RulesetClassic); classic != 1 {
		t.Errorf("Classic melee speed multiplier = %v, want 1: the Skyborne have no racials outside Forever", classic)
	}
}
