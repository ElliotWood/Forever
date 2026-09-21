package warlock

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

// Siphon Life is a Forever Affliction talent: a 30 sec shadow dot that heals the warlock for what
// it deals.
func (warlock *Warlock) registerSiphonLifeSpell() {
	if !warlock.Talents.SiphonLife {
		return
	}

	rank := spellData.SiphonLife.HighestRank()
	tick := rank.Periodic.(shared.SpellDataPeriodic)
	actionID := core.ActionID{SpellID: rank.SpellID}
	healthMetrics := warlock.NewHealthMetrics(actionID)

	warlock.SiphonLife = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    rank.SpellSchool,
		DefenseType:    rank.DefenseType,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL | core.SpellFlagBinary,
		ClassSpellMask: WarlockSpellSiphonLife,

		ManaCost: core.ManaCostOptions{FlatCost: rank.Cost},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: rank.GCD,
			},
		},

		DamageMultiplierAdditive: 1,
		DamageMultiplier:         1,
		ThreatMultiplier:         1,
		BonusCoefficient:         tick.Coef,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Siphon Life",
				Tag:   "Affliction",
			},
			NumberOfTicks:    tick.NumberOfTicks,
			TickLength:       tick.TickLength,
			BonusCoefficient: tick.Coef,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.Snapshot(target, tick.Tick)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				result := dot.CalcAndDealPeriodicSnapshotDamage(sim, target, shared.PeriodicTickOutcome(rank, dot))
				warlock.GainHealth(sim, result.Damage*warlock.PseudoStats.SelfHealingMultiplier, healthMetrics)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},

		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			dot := spell.Dot(target)
			if useSnapshot {
				return dot.CalcSnapshotDamage(sim, target, spell.OutcomeExpectedMagicAlwaysHit)
			}
			return spell.CalcPeriodicDamage(sim, target, tick.Tick, spell.OutcomeExpectedMagicAlwaysHit)
		},
	})
}
