package warlock

import (
	"github.com/wowsims/forever/sim/core"
)

func (warlock *Warlock) registerShadowBurn() {
	rank := spellData.Shadowburn.HighestRank()

	warlock.Shadowburn = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.SpellID},
		SpellSchool:    rank.SpellSchool,
		DefenseType:    rank.DefenseType,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL | core.SpellFlagBinary,
		ClassSpellMask: WarlockSpellShadowBurn,

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
		BonusCoefficient:         rank.Direct.BonusCoefficient(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, rank.Direct.Damage(sim), spell.OutcomeMagicHitAndCrit)
		},
	})
}
