package warlock

import (
	"github.com/wowsims/forever/sim/common/shared"
)

var siphonLifeRank = spellData.SiphonLife.HighestRank()
var siphonLifeTick = siphonLifeRank.Periodic.(shared.SpellDataPeriodic)
var siphonLifeCoeff = siphonLifeTick.Coef

// TODO: To be implemented. Port the TBC Siphon Life Spell implementation below; not yet verified against the Forever client.
func (warlock *Warlock) registerSiphonLifeSpell() {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// actionID := core.ActionID{SpellID: siphonLifeRank.SpellID}
	// baseCost := float64(siphonLifeRank.Cost)
	//
	// healthMetrics := warlock.NewHealthMetrics(actionID)
	//
	// warlock.SiphonLife = warlock.RegisterSpell(core.SpellConfig{
	// 	ActionID:       actionID,
	// 	SpellSchool:    core.SpellSchoolShadow,
	// 	DefenseType:    core.DefenseTypeMagic,
	// 	ClassSpellMask: WarlockSpellSiphonLife,
	// 	ProcMask:       core.ProcMaskSpellDamage,
	// 	Flags:          core.SpellFlagAPL,
	// 	BaseCost:       baseCost,
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			Cost: baseCost,
	// 			GCD:  siphonLifeRank.GCD,
	// 		},
	// 	},
	// 	DamageMultiplier: 1,
	// 	BonusCoefficient: siphonLifeCoeff,
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
	//
	// 		if result.Landed() {
	// 			spell.Dot(target).Apply(sim)
	// 		}
	// 		spell.DealOutcome(sim, result)
	// 	},
	//
	// 	Dot: core.DotConfig{
	// 		Aura: core.Aura{
	// 			Label: "SiphonLife",
	// 			Tag:   "Affliction",
	// 		},
	// 		NumberOfTicks:    siphonLifeTick.NumberOfTicks,
	// 		TickLength:       siphonLifeTick.TickLength,
	// 		BonusCoefficient: siphonLifeCoeff,
	//
	// 		OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
	// 			dot.Snapshot(target, siphonLifeTick.Tick)
	// 		},
	// 		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
	// 			result := dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
	//
	// 			healthToRegain := result.Damage * (1 * warlock.PseudoStats.BonusHealingTaken)
	// 			warlock.GainHealth(sim, healthToRegain, healthMetrics)
	// 			dot.Spell.ApplyAOEThreat(healthToRegain * 0.5)
	// 		},
	// 	},
	//
	// 	ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
	// 		dot := spell.Dot(target)
	// 		if useSnapshot {
	// 			result := dot.CalcSnapshotDamage(sim, target, dot.OutcomeTick)
	// 			result.Damage /= dot.TickPeriod().Seconds()
	// 			return result
	// 		} else {
	// 			result := spell.CalcPeriodicDamage(sim, target, siphonLifeTick.Tick*float64(siphonLifeTick.NumberOfTicks), spell.OutcomeExpectedMagicHit)
	// 			result.Damage /= dot.CalcTickPeriod().Round(time.Millisecond).Seconds()
	// 			return result
	// 		}
	// 	},
	// })
}
