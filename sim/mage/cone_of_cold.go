package mage

import (
	"github.com/wowsims/forever/sim/core"
)

func (mage *Mage) registerConeOfColdSpell() {
	coneOfColdRank := spellData.ConeOfCold.HighestRank()

	mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: coneOfColdRank.SpellID},
		SpellSchool:    coneOfColdRank.SpellSchool,
		DefenseType:    coneOfColdRank.DefenseType,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL | core.SpellFlagBinary,
		ClassSpellMask: MageSpellConeOfCold,

		ManaCost: core.ManaCostOptions{
			FlatCost: coneOfColdRank.Cost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: coneOfColdRank.GCD,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: coneOfColdRank.Cooldown,
			},
		},

		DamageMultiplier: 1,
		BonusCoefficient: coneOfColdRank.Direct.BonusCoefficient(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			spell.CalcAndDealAoeDamage(sim, coneOfColdRank.Direct.Damage(sim), spell.OutcomeMagicHitAndCrit)
		},
	})
}
