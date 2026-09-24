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

func TestReforgeAvoidanceRatingWeightMovesOntoItsCappedPercent(t *testing.T) {
	for _, test := range []struct {
		rating           stats.Stat
		percent          proto.PseudoStat
		ratingPerPercent float64
	}{
		{stats.DodgeRating, proto.PseudoStat_PseudoStatDodgePercent, core.DodgeRatingPerDodgePercent},
		{stats.ParryRating, proto.PseudoStat_PseudoStatParryPercent, core.ParryRatingPerParryPercent},
		{stats.BlockRating, proto.PseudoStat_PseudoStatBlockPercent, core.BlockRatingPerBlockPercent},
	} {
		percent := stats.UnitStatFromPseudoStat(test.percent)
		weights := core.NewUnitStats()
		weights.Stats[test.rating] = 1
		caps := setUnitStat(core.NewUnitStats(), percent, 2)

		validated := checkWeights(weights, caps, nil)
		if validated.Stats[test.rating] != 0 {
			t.Errorf("%s keeps weight %v under a %s cap, want 0", test.rating.StatName(), validated.Stats[test.rating], test.percent)
		}
		if got := getUnitStat(validated, percent); got != test.ratingPerPercent {
			t.Errorf("%s weighs %v per percent, want %v", test.percent, got, test.ratingPerPercent)
		}

		coeffs := map[string]float64{}
		o := &reforgeOptimizer{player: &proto.Player{}}
		o.applyReforgeStat(coeffs, test.rating, 3*test.ratingPerPercent, validated)
		if got := coeffs[pseudoStatCoeffKey(test.percent)]; math.Abs(got-3) > 1e-9 {
			t.Errorf("the rating for 3%% counts %v toward %s, want 3", got, test.percent)
		}
	}
}

func TestReforgeDefenseWeightMovesOntoCritTakenOnly(t *testing.T) {
	critTaken := stats.UnitStatFromPseudoStat(proto.PseudoStat_PseudoStatReducedCritTakenPercent)
	dodge := stats.UnitStatFromPseudoStat(proto.PseudoStat_PseudoStatDodgePercent)
	weights := core.NewUnitStats()
	weights.Stats[stats.DefenseRating] = 1
	caps := setUnitStat(setUnitStat(core.NewUnitStats(), critTaken, 2), dodge, 2)

	validated := checkWeights(weights, caps, nil)
	if got := getUnitStat(validated, critTaken); got != core.DefenseRatingPerAvoidancePercent {
		t.Errorf("crit taken weighs %v per percent, want the 25 defense one percent costs", got)
	}
	if got := getUnitStat(validated, dodge); got != 0 {
		t.Errorf("Dodge%% takes weight %v from defense, want 0", got)
	}
}

func TestReforgeCapSpaceSkipsBlockAndParryTheCharacterLacks(t *testing.T) {
	sdm := stats.NewStatDependencyManager()
	sdm.AddStatDependency(stats.DodgeRating, stats.DodgePercent, 1/core.DodgeRatingPerDodgePercent)
	sdm.AddStatDependency(stats.ParryRating, stats.ParryPercent, 1/core.ParryRatingPerParryPercent)
	sdm.AddStatDependency(stats.BlockRating, stats.BlockPercent, 1/core.BlockRatingPerBlockPercent/100)
	sdm.FinalizeStatDeps()

	var delta stats.Stats
	delta[stats.DodgeRating] = core.DodgeRatingPerDodgePercent
	delta[stats.ParryRating] = core.ParryRatingPerParryPercent
	delta[stats.BlockRating] = core.BlockRatingPerBlockPercent

	for _, character := range []struct{ canBlock, canParry bool }{{true, true}, {false, true}, {true, false}, {false, false}} {
		o := &reforgeOptimizer{statDeps: &sdm, baseStats: core.NewUnitStats(), canBlock: character.canBlock, canParry: character.canParry}
		coeffs := o.resolveCapCoeffs(delta)
		for _, want := range []struct {
			pseudoStat proto.PseudoStat
			counts     bool
		}{
			{proto.PseudoStat_PseudoStatDodgePercent, true},
			{proto.PseudoStat_PseudoStatBlockPercent, character.canBlock},
			{proto.PseudoStat_PseudoStatParryPercent, character.canParry},
		} {
			_, counts := coeffs[pseudoStatCoeffKey(want.pseudoStat)]
			if counts != want.counts {
				t.Errorf("can block %v, parry %v: %s in the cap space is %v", character.canBlock, character.canParry, want.pseudoStat, counts)
			}
		}
	}
}
