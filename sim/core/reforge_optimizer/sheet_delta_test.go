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
	unshielded := protectionWarriorRaid("preraid")
	unshielded.Parties[0].Players[0].Equipment.Items[proto.ItemSlot_ItemSlotOffHand] = &proto.ItemSpec{}

	block := proto.PseudoStat_PseudoStatBlockPercent
	dodge := proto.PseudoStat_PseudoStatDodgePercent
	parry := proto.PseudoStat_PseudoStatParryPercent
	critTaken := proto.PseudoStat_PseudoStatReducedCritTakenPercent
	for _, character := range []struct {
		name               string
		raid               *proto.Raid
		canBlock, canParry bool
	}{
		{"warrior with a shield", protectionWarriorRaid("preraid"), true, true},
		{"warrior without a shield", unshielded, false, true},
	} {
		t.Run(character.name, func(t *testing.T) {
			baseResult, sdm, pseudoStats := computeReforgeStatsAndDeps(&proto.ComputeStatsRequest{Raid: protopkg.Clone(character.raid).(*proto.Raid)})
			if baseResult.ErrorResult != "" {
				t.Fatalf("ComputeStats: %s", baseResult.ErrorResult)
			}
			if pseudoStats.CanBlock != character.canBlock || pseudoStats.CanParry != character.canParry {
				t.Fatalf("can block %v and parry %v, want %v and %v", pseudoStats.CanBlock, pseudoStats.CanParry, character.canBlock, character.canParry)
			}
			o := &reforgeOptimizer{
				statDeps:  sdm,
				canBlock:  pseudoStats.CanBlock,
				canParry:  pseudoStats.CanParry,
				baseStats: protoToCoreUnitStats(baseResult.RaidStats.Parties[0].Players[0].FinalStats),
			}

			for _, test := range []struct {
				stat  stats.Stat
				moves []proto.PseudoStat
			}{
				{stats.BlockRating, []proto.PseudoStat{block}},
				{stats.DefenseRating, []proto.PseudoStat{block, dodge, parry, critTaken}},
				{stats.ResilienceRating, []proto.PseudoStat{critTaken}},
				{stats.DodgeRating, []proto.PseudoStat{dodge}},
				{stats.ParryRating, []proto.PseudoStat{parry}},
			} {
				var delta stats.Stats
				delta[test.stat] = 25

				bonusRaid := protopkg.Clone(character.raid).(*proto.Raid)
				bonusRaid.Parties[0].Players[0].BonusStats = &proto.UnitStats{Stats: delta[:stats.ProtoStatsLen]}
				sheetDelta := subtractUnitStats(sheetStats(t, bonusRaid), o.baseStats)
				capCoeffs := o.resolveCapCoeffs(delta)

				for _, pseudoStat := range []proto.PseudoStat{block, dodge, parry, critTaken} {
					moves := slices.Contains(test.moves, pseudoStat) &&
						(pseudoStat != block || character.canBlock) &&
						(pseudoStat != parry || character.canParry)
					want := getUnitStat(sheetDelta, stats.UnitStatFromPseudoStat(pseudoStat))
					if (want != 0) != moves {
						t.Errorf("%s moves the sheet's %s by %v", test.stat.StatName(), pseudoStat, want)
					}
					if got := capCoeffs[pseudoStatCoeffKey(pseudoStat)]; math.Abs(got-want) > 1e-9 {
						t.Errorf("%s moves the cap space's %s by %v, the sheet's by %v", test.stat.StatName(), pseudoStat, got, want)
					}
				}
			}
		})
	}
}
