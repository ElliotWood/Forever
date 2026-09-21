package druid

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

var ripRank = spellData.Rip.HighestRank()
var ripTick = ripRank.Periodic.(shared.SpellDataPeriodic)

// The per-combo-point damage and the attack power share are not in the generated table - the client
// states the tick base only - so they stay the sim's client-read values: 25.5 a point a tick at rank
// 6, and 1% of attack power a point, which stops growing at four points.
const ripTickPerComboPoint = 25.5

func (druid *Druid) registerRipSpell() {
	druid.Rip = druid.RegisterSpell(Cat, core.SpellConfig{
		ActionID:       core.ActionID{SpellID: ripRank.SpellID},
		SpellSchool:    ripRank.SpellSchool,
		DefenseType:    ripRank.DefenseType,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		ClassSpellMask: DruidSpellRip,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		Rank:           ripRank.Rank,

		EnergyCost: core.EnergyCostOptions{
			Cost: ripRank.Cost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: ripRank.GCD,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(_ *core.Simulation, _ *core.Unit) bool {
			return druid.ComboPoints() > 0
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		MaxRange:         core.MaxMeleeRange,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Rip",
			},
			NumberOfTicks: ripTick.NumberOfTicks,
			TickLength:    ripTick.TickLength,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.SnapshotPhysical(target, ripTickDamage(float64(druid.ComboPoints()), dot.Spell.MeleeAttackPower(target)))
				druid.UpdateBleedPower(druid.Rip, sim, target, true, true)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, shared.PeriodicTickOutcome(ripRank, dot))
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHitNoHitCounter)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
				druid.SpendComboPoints(sim, spell.ComboPointMetrics())
			} else {
				spell.IssueRefund(sim)
			}
			spell.DealOutcome(sim, result)
		},

		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			if useSnapshot {
				dot := spell.Dot(target)
				return dot.CalcSnapshotDamage(sim, target, dot.OutcomeTick)
			}
			// Assume 5 CP for projections.
			tick := ripTickDamage(5, spell.MeleeAttackPower(target))
			result := spell.CalcPeriodicDamage(sim, target, tick, spell.OutcomeExpectedMagicAlwaysHit)
			attackTable := spell.Unit.AttackTables[target.UnitIndex]
			critChance := spell.PhysicalCritChance(attackTable)
			result.Damage *= 1 + critChance*(spell.CritDamageMultiplier(attackTable)-1)
			return result
		},
	})

	druid.Rip.ShortName = "Rip"
}

func ripTickDamage(comboPoints float64, attackPower float64) float64 {
	return ripTick.Tick + ripTickPerComboPoint*comboPoints + 0.01*min(comboPoints, 4)*attackPower
}

func (druid *Druid) CurrentRipCost() float64 {
	return druid.Rip.Cost.GetCurrentCost()
}
