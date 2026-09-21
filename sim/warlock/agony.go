package warlock

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

// Forever renamed Curse of Agony to Bane of Agony and moved it onto the bane slot, so it holds the
// target beside a curse. It ramps: the snapshot pays half the tick, and the other half is added
// back every four ticks, leaving the last four at 150%. Amplify Curse adds half again and is spent
// on the snapshot.
func (warlock *Warlock) registerCurseOfAgony() {
	rank := spellData.BaneOfAgony.HighestRank()
	tick := rank.Periodic.(shared.SpellDataPeriodic)
	amplify := 1 + spellData.AmplifyCurse.EffectAt(0).FractionAt(1)

	rampStep := 0.0

	warlock.CurseOfAgony = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.SpellID},
		SpellSchool:    rank.SpellSchool,
		DefenseType:    rank.DefenseType,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellCurseOfAgony,

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
				Label: "Bane of Agony",
				Tag:   "Affliction",
			},
			NumberOfTicks:    tick.NumberOfTicks,
			TickLength:       tick.TickLength,
			BonusCoefficient: tick.Coef,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				base := tick.Tick
				if warlock.AmplifyCurseAura.IsActive() {
					base *= amplify
					warlock.AmplifyCurseAura.Deactivate(sim)
				}

				rampStep = base * 0.5
				dot.Snapshot(target, rampStep)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, shared.PeriodicTickOutcome(rank, dot))
				if dot.TickCount()%4 == 0 {
					dot.SnapshotBaseDamage += rampStep
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)
			if result.Landed() {
				dot := spell.Dot(target)
				warlock.takeBaneSlot(sim, target, dot.Aura)
				dot.Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},

		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			// The ramp averages out to the plain tick over the full duration.
			dot := spell.Dot(target)
			if useSnapshot {
				return dot.CalcSnapshotDamage(sim, target, spell.OutcomeExpectedMagicAlwaysHit)
			}
			return spell.CalcPeriodicDamage(sim, target, tick.Tick, spell.OutcomeExpectedMagicAlwaysHit)
		},
	})
}
