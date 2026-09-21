package mage

import (
	"github.com/wowsims/forever/sim/core"
)

func (mage *Mage) registerFrostNovaSpell() {
	frostNovaRank := spellData.FrostNova.HighestRank()

	mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: frostNovaRank.SpellID},
		SpellSchool:    frostNovaRank.SpellSchool,
		DefenseType:    frostNovaRank.DefenseType,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL | core.SpellFlagBinary,
		ClassSpellMask: MageSpellFrostNova,

		ManaCost: core.ManaCostOptions{
			FlatCost: frostNovaRank.Cost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: frostNovaRank.GCD,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: frostNovaRank.Cooldown,
			},
		},

		DamageMultiplier: 1,
		BonusCoefficient: frostNovaRank.Direct.BonusCoefficient(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			spell.CalcAndDealAoeDamage(sim, frostNovaRank.Direct.Damage(sim), spell.OutcomeMagicHitAndCrit)
		},
	})
}
