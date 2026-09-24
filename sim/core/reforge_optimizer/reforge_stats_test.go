package reforgeoptimizer

import (
	"math"
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

func pseudoStatsProto(values map[proto.PseudoStat]float64) *proto.UnitStats {
	pseudoStats := make([]float64, stats.PseudoStatsLen)
	for pseudoStat, value := range values {
		pseudoStats[pseudoStat] = value
	}
	return &proto.UnitStats{Stats: make([]float64, stats.ProtoStatsLen), PseudoStats: pseudoStats}
}

func TestReforgeWeightsKeepPercentWeightsPerPercent(t *testing.T) {
	weights := protoToCoreUnitStats(pseudoStatsProto(map[proto.PseudoStat]float64{
		proto.PseudoStat_PseudoStatBlockPercent:     5,
		proto.PseudoStat_PseudoStatMeleeHitPercent:  2,
		proto.PseudoStat_PseudoStatRangedHitPercent: 3,
	}))

	if got := computeCoeffScore(map[string]float64{pseudoStatCoeffKey(proto.PseudoStat_PseudoStatBlockPercent): 1}, weights); got != 5 {
		t.Errorf("one Block%% scores %v, want the 5 per percent the EP values state", got)
	}
	if got := getUnitStat(weights, stats.UnitStatFromPseudoStat(proto.PseudoStat_PseudoStatRangedHitPercent)); got != 3 {
		t.Errorf("ranged hit weight %v, want 3", got)
	}
	for statIdx := int(stats.ProtoStatsLen); statIdx < int(stats.SimStatsLen); statIdx++ {
		if weights.Stats[statIdx] != 0 {
			t.Errorf("%s carries %v; the optimizer reads percent weights from PseudoStats only", stats.Stat(statIdx).StatName(), weights.Stats[statIdx])
		}
	}
}

func TestReforgeHitCapMatchesTheBackEndHitPercent(t *testing.T) {
	sdm := stats.NewStatDependencyManager()
	sdm.AddStatDependency(stats.MeleeHitRating, stats.PhysicalHitPercent, 1/core.PhysicalHitRatingPerHitPercent)
	sdm.FinalizeStatDeps()

	baseStats := protoToCoreUnitStats(pseudoStatsProto(map[proto.PseudoStat]float64{
		proto.PseudoStat_PseudoStatMeleeHitPercent: 7,
		proto.PseudoStat_PseudoStatBlockPercent:    25,
	}))
	caps := computeStatCapsDelta(baseStats, protoToCoreUnitStats(pseudoStatsProto(map[proto.PseudoStat]float64{
		proto.PseudoStat_PseudoStatMeleeHitPercent: 9,
		proto.PseudoStat_PseudoStatBlockPercent:    30,
	})))

	meleeHit := stats.UnitStatFromPseudoStat(proto.PseudoStat_PseudoStatMeleeHitPercent)
	if gap := getUnitStat(caps, meleeHit); gap != 2 {
		t.Fatalf("9%% hit cap leaves a gap of %v over 7%% hit, want 2", gap)
	}
	if gap := getUnitStat(caps, stats.UnitStatFromPseudoStat(proto.PseudoStat_PseudoStatBlockPercent)); gap != 5 {
		t.Errorf("30%% block cap leaves a gap of %v over 25%% block, want 5", gap)
	}

	delta := core.NewUnitStats()
	delta.Stats[stats.MeleeHitRating] = 2 * core.PhysicalHitRatingPerHitPercent
	resolved := resolveStatDelta(&sdm, baseStats, delta)
	if got := getUnitStat(resolved, meleeHit); math.Abs(got-2) > 1e-9 {
		t.Errorf("rating for 2%% hit fills %v of the hit cap gap, want 2", got)
	}
}

func TestReforgeBlockDeltaLandsInPercent(t *testing.T) {
	sdm := stats.NewStatDependencyManager()
	sdm.AddStatDependency(stats.Strength, stats.BlockPercent, 0.005)
	sdm.FinalizeStatDeps()

	delta := core.NewUnitStats()
	delta.Stats[stats.Strength] = 10
	resolved := resolveStatDelta(&sdm, core.NewUnitStats(), delta)
	if got := resolved.Stats[stats.BlockPercent]; math.Abs(got-0.05) > 1e-12 {
		t.Fatalf("back-end BlockPercent delta %v, want the 0.05 probability", got)
	}
	if got := getUnitStat(resolved, stats.UnitStatFromPseudoStat(proto.PseudoStat_PseudoStatBlockPercent)); math.Abs(got-5) > 1e-9 {
		t.Errorf("Block%% delta %v, want 5 percent", got)
	}
}

func TestReforgeDodgeAndParryRatingReachTheirPercentCaps(t *testing.T) {
	sdm := stats.NewStatDependencyManager()
	sdm.AddStatDependency(stats.DodgeRating, stats.DodgePercent, 1/core.DodgeRatingPerDodgePercent)
	sdm.AddStatDependency(stats.ParryRating, stats.ParryPercent, 1/core.ParryRatingPerParryPercent)
	sdm.FinalizeStatDeps()

	delta := core.NewUnitStats()
	delta.Stats[stats.DodgeRating] = 2 * core.DodgeRatingPerDodgePercent
	delta.Stats[stats.ParryRating] = 3 * core.ParryRatingPerParryPercent
	resolved := resolveStatDelta(&sdm, core.NewUnitStats(), delta)
	if got := getUnitStat(resolved, stats.UnitStatFromPseudoStat(proto.PseudoStat_PseudoStatDodgePercent)); math.Abs(got-2) > 1e-9 {
		t.Errorf("rating for 2%% dodge moves Dodge%% by %v, want 2", got)
	}
	if got := getUnitStat(resolved, stats.UnitStatFromPseudoStat(proto.PseudoStat_PseudoStatParryPercent)); math.Abs(got-3) > 1e-9 {
		t.Errorf("rating for 3%% parry moves Parry%% by %v, want 3", got)
	}
}

func TestReforgeDefenseAndResilienceReachTheCritTakenCap(t *testing.T) {
	sdm := stats.NewStatDependencyManager()
	sdm.FinalizeStatDeps()

	delta := core.NewUnitStats()
	delta.Stats[stats.DefenseRating] = 25 * core.DefenseRatingPerDefenseLevel
	delta.Stats[stats.ResilienceRating] = 2 * core.ResilienceRatingPerCritReductionChance
	resolved := resolveStatDelta(&sdm, core.NewUnitStats(), delta)
	if got := getUnitStat(resolved, stats.UnitStatFromPseudoStat(proto.PseudoStat_PseudoStatReducedCritTakenPercent)); math.Abs(got-3) > 1e-9 {
		t.Errorf("25 defense and the resilience for 2%% move crit taken by %v, want 3", got)
	}
}
