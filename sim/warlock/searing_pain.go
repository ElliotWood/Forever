package warlock

import (
	"github.com/wowsims/forever/sim/core"
)

func (warlock *Warlock) registerSearingPain() {
	rank := spellData.SearingPain.HighestRank()

	warlock.SearingPain = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.SpellID},
		SpellSchool:    rank.SpellSchool,
		DefenseType:    rank.DefenseType,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellSearingPain,
		MaxRange:       rank.MaxRange,

		ManaCost: core.ManaCostOptions{FlatCost: rank.Cost},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      rank.GCD,
				CastTime: rank.CastTime,
			},
		},

		DamageMultiplierAdditive: 1,
		DamageMultiplier:         1,
		ThreatMultiplier:         2,
		BonusCoefficient:         rank.Direct.BonusCoefficient(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, rank.Direct.Damage(sim), spell.OutcomeMagicHitAndCrit)
		},
	})
}
