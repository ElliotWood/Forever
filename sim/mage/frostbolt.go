package mage

import (
	"github.com/wowsims/forever/sim/core"
)

func (mage *Mage) registerFrostboltSpell() {
	frostboltRank := spellData.Frostbolt.HighestRank()

	mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: frostboltRank.SpellID},
		SpellSchool:    frostboltRank.SpellSchool,
		DefenseType:    frostboltRank.DefenseType,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL | core.SpellFlagBinary,
		ClassSpellMask: MageSpellFrostbolt,
		MissileSpeed:   frostboltRank.MissileSpeed,

		ManaCost: core.ManaCostOptions{
			FlatCost: frostboltRank.Cost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      frostboltRank.GCD,
				CastTime: frostboltRank.CastTime,
			},
		},

		DamageMultiplier: 1,
		BonusCoefficient: frostboltRank.Direct.BonusCoefficient(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcDamage(sim, target, frostboltRank.Direct.Damage(sim), spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
