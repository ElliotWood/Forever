package warlock

import (
	"time"

	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

var uaRank = spellData.UnstableAffliction.BySpellID(30405)
var uaTick = uaRank.Periodic.(shared.SpellDataPeriodic)
var uaCoeff = uaTick.Coef

// TODO: uncalled -- Forever drops the Unstable Affliction talent that gated this spell, so
// nothing calls this any more. Kept rather than deleted because "the talent is gone"
// and "the spell is gone" are not the same claim, and the client data does not
// distinguish them. Re-gate before wiring it back up.
func (warlock *Warlock) registerUnstableAffliction() {
	warlock.UnstableAffliction = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: uaRank.SpellID},
		SpellSchool:    uaRank.SpellSchool,
		DefenseType:    uaRank.DefenseType,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellUnstableAffliction,

		ManaCost: core.ManaCostOptions{FlatCost: uaRank.Cost},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      uaRank.GCD,
				CastTime: uaRank.CastTime,
			},
		},

		DamageMultiplierAdditive: 1,
		ThreatMultiplier:         1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
		},
		BonusCoefficient: uaCoeff,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:    "Unstable Affliction",
				Tag:      "Affliction",
				ActionID: core.ActionID{SpellID: 30108},
			},
			NumberOfTicks:    uaTick.NumberOfTicks,
			TickLength:       uaTick.TickLength,
			BonusCoefficient: uaCoeff,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.Snapshot(target, uaTick.Tick)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			dot := spell.Dot(target)
			if useSnapshot {
				result := dot.CalcSnapshotDamage(sim, target, dot.OutcomeTick)
				result.Damage /= dot.TickPeriod().Seconds()
				return result
			} else {
				result := spell.CalcPeriodicDamage(sim, target, uaTick.Tick*float64(uaTick.NumberOfTicks), spell.OutcomeExpectedMagicHit)
				result.Damage /= dot.CalcTickPeriod().Round(time.Millisecond).Seconds()
				return result
			}
		},
	})
}
