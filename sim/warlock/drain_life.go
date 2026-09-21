package warlock

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

// spellData.DrainLife holds the six trainer ranks, 689 to 11700. The client also carries 403677 to
// 403689, the copies the Season of Discovery rune passive Master Channeler (403668) swaps onto the
// action bar; that is an Engrave grant with no place in Forever and the generator drops it.
//
// Improved Drains rides on the talent as a dot SpellMod. Soul Siphon has to be counted per tick,
// so it stays here.
func (warlock *Warlock) registerDrainLife() {
	rank := spellData.DrainLife.HighestRank()
	tick := rank.Periodic.(shared.SpellDataPeriodic)
	healthMetric := warlock.NewHealthMetrics(core.ActionID{SpellID: rank.SpellID})

	warlock.DrainLife = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.SpellID},
		SpellSchool:    rank.SpellSchool,
		DefenseType:    rank.DefenseType,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagChanneled | core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellDrainLife,

		ManaCost: core.ManaCostOptions{FlatCost: rank.Cost},
		Cast:     core.CastConfig{DefaultCast: core.Cast{GCD: rank.GCD}},

		DamageMultiplierAdditive: 1,
		DamageMultiplier:         1,
		ThreatMultiplier:         1,
		BonusCoefficient:         tick.Coef,

		Dot: core.DotConfig{
			Aura:                 core.Aura{Label: "Drain Life"},
			NumberOfTicks:        tick.NumberOfTicks,
			TickLength:           tick.TickLength,
			AffectedByCastSpeed:  true,
			HasteReducesDuration: true,
			BonusCoefficient:     tick.Coef,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.Snapshot(target, tick.Tick)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				// Scaled here rather than through dot.PeriodicDamageMultiplier, which Improved
				// Drains' SpellMod already owns.
				result := dot.CalcSnapshotDamage(sim, target, shared.PeriodicTickOutcome(rank, dot))
				result.Damage *= warlock.soulSiphonMultiplier(target)
				dot.Spell.DealPeriodicDamage(sim, result)
				warlock.GainHealth(sim, result.Damage*warlock.PseudoStats.SelfHealingMultiplier, healthMetric)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})
}

// Soul Siphon pays 4% per point for each of the warlock's other Affliction effects on the target,
// counting up to three. The per-effect value is the talent's own dummy ladder.
//
// ponytail: counted off the "Affliction"-tagged dot auras only; a curse without a dot is missed.
// Corruption, Bane of Agony and Siphon Life already reach the cap in every build that takes it.
func (warlock *Warlock) soulSiphonMultiplier(target *core.Unit) float64 {
	if warlock.Talents.SoulSiphon == 0 {
		return 1
	}

	perEffect := spellData.SoulSiphon.FractionAt(warlock.Talents.SoulSiphon)
	return 1 + perEffect*min(warlock.AfflictionCount(target), 3)
}
