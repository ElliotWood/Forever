package mage

import (
	"github.com/wowsims/forever/sim/core"
)

func (mage *Mage) registerArcaneExplosionSpell() {
	arcaneExplosionRank := spellData.ArcaneExplosion.HighestRank()

	mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: arcaneExplosionRank.SpellID},
		SpellSchool:    arcaneExplosionRank.SpellSchool,
		DefenseType:    arcaneExplosionRank.DefenseType,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MageSpellArcaneExplosion,

		ManaCost: core.ManaCostOptions{
			FlatCost: arcaneExplosionRank.Cost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: arcaneExplosionRank.GCD,
			},
		},

		DamageMultiplier: 1,
		BonusCoefficient: arcaneExplosionRank.Direct.BonusCoefficient(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealAoeDamage(sim, arcaneExplosionRank.Direct.Damage(sim), spell.OutcomeMagicHitAndCrit)
		},
	})
}
