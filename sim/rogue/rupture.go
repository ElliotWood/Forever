package rogue

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

// Was the TBC rank-6 id 26867, which the level squish removed. Derived from the table so
// it follows the data instead of naming a rank that may not exist.
var RuptureSpellID = spellData.Rupture.HighestRank().SpellID

var ruptureRank = spellData.Rupture.BySpellID(RuptureSpellID)

func (rogue *Rogue) registerRupture() {
	tick := ruptureRank.Periodic.(shared.SpellDataPeriodic)

	// The beta client cut the per combo point step with the tick (rank 6: 60 + 8 -> 35 + 4.73).
	// The table carries the 35; the step sits on a dummy effect the generator reads as 0.
	const damagePerComboPoint = 4.73

	rogue.Rupture = rogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: ruptureRank.SpellID},
		SpellSchool:    ruptureRank.SpellSchool,
		DefenseType:    ruptureRank.DefenseType,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | SpellFlagFinisher | core.SpellFlagAPL,
		MetricSplits:   6,
		ClassSpellMask: RogueSpellRupture,
		MaxRange:       core.MaxMeleeRange,

		EnergyCost: core.EnergyCostOptions{
			Cost:          ruptureRank.Cost,
			Refund:        0,
			RefundMetrics: rogue.EnergyRefundMetrics,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: ruptureRank.GCD,
			},
			IgnoreHaste: true,
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				spell.SetMetricsSplit(rogue.ComboPoints())
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return rogue.ComboPoints() > 0
		},

		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		ThreatMultiplier:         1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Rupture",
				Tag:   RogueBleedTag,
			},
			NumberOfTicks: 0, // Set dynamically
			TickLength:    tick.TickLength,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				damage := rogue.ruptureDamage(target, rogue.ComboPoints(), tick.Tick, damagePerComboPoint)
				if rogue.isHemorrhaging(target) {
					damage *= HemorrhageRuptureMultiplier
				}
				dot.SnapshotPhysical(target, damage)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, shared.PeriodicTickOutcome(ruptureRank, dot))
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
			if result.Landed() {
				dot := spell.Dot(target)
				dot.BaseTickCount = tick.NumberOfTicks + rogue.ComboPoints()
				dot.Apply(sim)
				rogue.ApplyFinisher(sim, spell)
			} else {
				spell.IssueRefund(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})
}

func (rogue *Rogue) ruptureDamage(target *core.Unit, comboPoints int32, baseDamage float64, damagePerComboPoint float64) float64 {
	return baseDamage +
		damagePerComboPoint*float64(comboPoints) +
		[]float64{0, 0.01, 0.02, 0.03, 0.03, 0.03}[comboPoints]*rogue.Rupture.MeleeAttackPower(target)
}
