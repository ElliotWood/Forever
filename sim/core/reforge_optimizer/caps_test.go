package reforgeoptimizer

import (
	"testing"

	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

func TestBuildDebuffUnitStats(t *testing.T) {
	critStats := []proto.PseudoStat{
		proto.PseudoStat_PseudoStatMeleeCritPercent,
		proto.PseudoStat_PseudoStatRangedCritPercent,
		proto.PseudoStat_PseudoStatSpellCritPercent,
	}

	cases := []struct {
		name    string
		debuffs *proto.Debuffs
		want    float64
	}{
		{name: "seal of the crusader", debuffs: &proto.Debuffs{ImprovedSealOfTheCrusader: true}, want: 3},
		{name: "no seal of the crusader", debuffs: &proto.Debuffs{}, want: 0},
		{name: "another debuff", debuffs: &proto.Debuffs{FaerieFire: true}, want: 0},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			unitStats := buildDebuffUnitStats(&proto.Raid{Debuffs: c.debuffs})

			for _, pseudoStat := range critStats {
				unitStat := stats.UnitStatFromPseudoStat(pseudoStat)
				if got := unitStats.PseudoStats[unitStat.PseudoStatIdx()]; got != c.want {
					t.Errorf("%s = %v, want %v", pseudoStat, got, c.want)
				}
			}
		})
	}
}
