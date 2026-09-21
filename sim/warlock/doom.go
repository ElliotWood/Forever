package warlock

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

// Forever renamed Curse of Doom to Bane of Doom (603): one tick of 1742 after a minute, with a 4.0
// spell power coefficient, on the bane slot.
func (warlock *Warlock) registerCurseOfDoom() {
	rank := spellData.BaneOfDoom.HighestRank()
	tick := rank.Periodic.(shared.SpellDataPeriodic)
	amplify := 1 + spellData.AmplifyCurse.EffectAt(0).FractionAt(1)

	baseDamage := func(sim *core.Simulation) float64 {
		if warlock.AmplifyCurseAura.IsActive() {
			warlock.AmplifyCurseAura.Deactivate(sim)
			return tick.Tick * amplify
		}
		return tick.Tick
	}

	warlock.CurseOfDoom = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.SpellID},
		SpellSchool:    rank.SpellSchool,
		DefenseType:    rank.DefenseType,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellCurseOfDoom,

		ManaCost: core.ManaCostOptions{FlatCost: rank.Cost},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: rank.GCD,
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: rank.Cooldown,
			},
		},

		DamageMultiplierAdditive: 1,
		DamageMultiplier:         1,
		ThreatMultiplier:         1,
		BonusCoefficient:         tick.Coef,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Bane of Doom",
				Tag:   "Affliction",
			},
			NumberOfTicks:    tick.NumberOfTicks,
			TickLength:       tick.TickLength,
			BonusCoefficient: tick.Coef,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.Snapshot(target, baseDamage(sim))
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, shared.PeriodicTickOutcome(rank, dot))
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
			dot := spell.Dot(target)
			if useSnapshot {
				return dot.CalcSnapshotDamage(sim, target, spell.OutcomeExpectedMagicAlwaysHit)
			}
			return spell.CalcPeriodicDamage(sim, target, tick.Tick, spell.OutcomeExpectedMagicAlwaysHit)
		},
	})
}
