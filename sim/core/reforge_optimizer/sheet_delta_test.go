//go:build with_db

package reforgeoptimizer

import (
	"math"
	"slices"
	"testing"

	"github.com/wowsims/forever/sim"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
	protopkg "google.golang.org/protobuf/proto"
)

func protectionWarriorRaid(gearSet string) *proto.Raid {
	return &proto.Raid{Parties: []*proto.Party{{Players: []*proto.Player{{
		Class:     proto.Class_ClassWarrior,
		Race:      proto.Race_RaceOrc,
		Equipment: core.GetGearSet("../../../ui/specs/warrior/protection/gear_sets", gearSet).GearSet,
		Spec: &proto.Player_ProtectionWarrior{ProtectionWarrior: &proto.ProtectionWarrior{
			Options: &proto.ProtectionWarrior_Options{ClassOptions: &proto.WarriorOptions{}},
		}},
	}}}}}
}

func sheetStats(t *testing.T, raid *proto.Raid) core.UnitStats {
	t.Helper()
	result := computeReforgeStats(&proto.ComputeStatsRequest{Raid: raid})
	if result.ErrorResult != "" {
		t.Fatalf("ComputeStats: %s", result.ErrorResult)
	}
	return protoToCoreUnitStats(result.RaidStats.Parties[0].Players[0].FinalStats)
}

func TestCapSpaceDeltaMatchesTheSheet(t *testing.T) {
	sim.RegisterAll()
	raid := protectionWarriorRaid("preraid")
	baseResult, sdm := computeReforgeStatsAndDeps(&proto.ComputeStatsRequest{Raid: protopkg.Clone(raid).(*proto.Raid)})
	if baseResult.ErrorResult != "" {
		t.Fatalf("ComputeStats: %s", baseResult.ErrorResult)
	}
	baseStats := protoToCoreUnitStats(baseResult.RaidStats.Parties[0].Players[0].FinalStats)

	block := proto.PseudoStat_PseudoStatBlockPercent
	dodge := proto.PseudoStat_PseudoStatDodgePercent
	parry := proto.PseudoStat_PseudoStatParryPercent
	for _, test := range []struct {
		stat  stats.Stat
		moves []proto.PseudoStat
	}{
		{stats.BlockRating, []proto.PseudoStat{block}},
		{stats.DefenseRating, []proto.PseudoStat{block, dodge, parry}},
		{stats.DodgeRating, []proto.PseudoStat{dodge}},
		{stats.ParryRating, []proto.PseudoStat{parry}},
	} {
		t.Run(test.stat.StatName(), func(t *testing.T) {
			delta := core.NewUnitStats()
			delta.Stats[test.stat] = 25

			bonusRaid := protopkg.Clone(raid).(*proto.Raid)
			bonusRaid.Parties[0].Players[0].BonusStats = &proto.UnitStats{Stats: delta.Stats[:stats.ProtoStatsLen]}
			sheetDelta := subtractUnitStats(sheetStats(t, bonusRaid), baseStats)
			resolved := resolveStatDelta(sdm, baseStats, delta)

			for _, pseudoStat := range []proto.PseudoStat{block, dodge, parry} {
				unitStat := stats.UnitStatFromPseudoStat(pseudoStat)
				want := getUnitStat(sheetDelta, unitStat)
				if (want != 0) != slices.Contains(test.moves, pseudoStat) {
					t.Errorf("%s: the sheet moves by %v", pseudoStat, want)
				}
				if got := getUnitStat(resolved, unitStat); math.Abs(got-want) > 1e-9 {
					t.Errorf("%s: cap space moves by %v, the sheet by %v", pseudoStat, got, want)
				}
			}
		})
	}
}
