package mage

import (
	"github.com/wowsims/forever/sim/core"
)

var coneOfColdRank = spellData.ConeOfCold.HighestRank()

func (mage *Mage) registerConeOfColdSpell() {

	coneOfColdCoefficient := 0.1930000037 // Per https://wago.tools/db2/SpellEffect?build=2.5.5.65295&filter%5BSpellID%5D=exact%253A120 Field "EffetBonusCoefficient"

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
		BonusCoefficient: coneOfColdCoefficient,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			baseDamage := coneOfColdRank.Direct.Damage(sim)
			spell.CalcAndDealAoeDamage(sim, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}
