package druid

import (
	"fmt"
	"time"

	"github.com/wowsims/classic/sim/core"
)

const MoonfireRanks = 10

var MoonfireSpellId = [MoonfireRanks + 1]int32{0, 8921, 8924, 8925, 8926, 8927, 8928, 8929, 9833, 9834, 9835}

// Beta client 1.60.1.69893: every rank carries the full .15 direct and .13 per tick coefficients, and the damage is
// lower from rank 2 up (rank 10 195-228 plus 96 a tick -> 128-151 plus 60 a tick). The dot total is the client's tick
// times the tick count; the coefficient stays per tick, as the client stores it.
var MoonfiresSpellCoeff = [MoonfireRanks + 1]float64{0, .15, .15, .15, .15, .15, .15, .15, .15, .15, .15}
var MoonfiresSellDotCoeff = [MoonfireRanks + 1]float64{0, .13, .13, .13, .13, .13, .13, .13, .13, .13, .13}
var MoonfireBaseDamage = [MoonfireRanks + 1][]float64{{0}, {9, 12}, {15, 19}, {25, 29}, {36, 43}, {50, 59}, {62, 72}, {75, 88}, {96, 113}, {112, 131}, {128, 151}}
var MoonfireBaseDotDamage = [MoonfireRanks + 1]float64{0, 12, 24, 36, 52, 76, 92, 124, 156, 196, 240}
var MoonfireDotTicks = [MoonfireRanks + 1]int32{0, 3, 4, 4, 4, 4, 4, 4, 4, 4, 4}
var MoonfireManaCost = [MoonfireRanks + 1]float64{0, 25, 50, 75, 105, 150, 190, 235, 280, 325, 375}
var MoonfireLevel = [MoonfireRanks + 1]int{0, 4, 10, 16, 22, 28, 34, 40, 46, 52, 58}

func (druid *Druid) registerMoonfireSpell() {
	druid.Moonfire = make([]*DruidSpell, 0)

	for rank := 1; rank <= MoonfireRanks; rank++ {
		config := druid.getMoonfireBaseConfig(rank)

		if config.RequiredLevel <= int(druid.Level) {
			druid.Moonfire = append(druid.Moonfire, druid.RegisterSpell(Humanoid|Moonkin, config))
		}
	}
}

func (druid *Druid) getMoonfireBaseConfig(rank int) core.SpellConfig {
	ticks := MoonfireDotTicks[rank]
	tickLength := time.Second * 3

	spellId := MoonfireSpellId[rank]
	spellCoeff := MoonfiresSpellCoeff[rank]
	spellDotCoeff := MoonfiresSellDotCoeff[rank]
	baseDamageLow := MoonfireBaseDamage[rank][0]
	baseDamageHigh := MoonfireBaseDamage[rank][1]
	baseDotDamage := (MoonfireBaseDotDamage[rank] / float64(ticks))
	manaCost := MoonfireManaCost[rank]
	level := MoonfireLevel[rank]

	return core.SpellConfig{
		ActionID:       core.ActionID{SpellID: spellId},
		SpellCode:      SpellCode_DruidMoonfire,
		ClassSpellMask: SpellMaskMoonfire,
		SpellSchool:    core.SpellSchoolArcane,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL | core.SpellFlagResetAttackSwing,

		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: 0,
			},
		},
		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:    fmt.Sprintf("Moonfire (Rank %d)", rank),
				ActionID: core.ActionID{SpellID: spellId},
			},
			NumberOfTicks:    ticks,
			TickLength:       tickLength,
			BonusCoefficient: spellDotCoeff,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, baseDotDamage, isRollover)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		BonusCoefficient: spellCoeff,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageLow, baseDamageHigh)
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			if result.Landed() {
				dot := spell.Dot(target)
				dot.Apply(sim)
			}
		},
	}
}
