package druid

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

var rakeRank = spellData.Rake.HighestRank()
var rakeTick = rakeRank.Periodic.(shared.SpellDataPeriodic)

// Forever's Rake no longer scales with attack power: the client states a flat hit and a flat tick,
// and carries no BonusCoefficientFromAP on either.
func (druid *Druid) registerRakeSpell() {
	druid.Rake = druid.RegisterSpell(Cat, core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rakeRank.SpellID},
		SpellSchool:    rakeRank.SpellSchool,
		DefenseType:    rakeRank.DefenseType,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		ClassSpellMask: DruidSpellRake,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		Rank:           rakeRank.Rank,

		EnergyCost: core.EnergyCostOptions{
			Cost:   rakeRank.Cost,
			Refund: rakeRank.MissRefund(),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: rakeRank.GCD,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		MaxRange:         core.MaxMeleeRange,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:    "Rake",
				Duration: rakeRank.Duration,
			},
			NumberOfTicks: rakeTick.NumberOfTicks,
			TickLength:    rakeTick.TickLength,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.SnapshotPhysical(target, rakeTick.Tick)
				druid.UpdateBleedPower(druid.Rake, sim, target, true, true)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, shared.PeriodicTickOutcome(rakeRank, dot))
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealDamage(sim, target, rakeRank.Direct.Damage(sim), spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if result.Landed() {
				druid.AddComboPoints(sim, 1, spell.ComboPointMetrics())
				spell.Dot(target).Apply(sim)
			} else {
				spell.IssueRefund(sim)
			}
		},

		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			if useSnapshot {
				dot := spell.Dot(target)
				return dot.CalcSnapshotDamage(sim, target, dot.OutcomeTick)
			}
			ticks := spell.CalcPeriodicDamage(sim, target, rakeTick.Tick, spell.OutcomeExpectedMagicAlwaysHit)
			attackTable := spell.Unit.AttackTables[target.UnitIndex]
			critChance := spell.PhysicalCritChance(attackTable)
			ticks.Damage *= 1 + critChance*(spell.CritDamageMultiplier(attackTable)-1)
			return ticks
		},
	})

	druid.Rake.ShortName = "Rake"
}

func (druid *Druid) CurrentRakeCost() float64 {
	return druid.Rake.Cost.GetCurrentCost()
}
