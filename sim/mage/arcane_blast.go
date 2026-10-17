package mage

import (
	"github.com/wowsims/forever/sim/core"
)

func (mage *Mage) registerArcaneBlastSpell() {
	if !mage.Talents.ArcaneBlast {
		return
	}

	arcaneBlastRank := spellData.ArcaneBlast.Highest()

	mage.ArcaneBlast = mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: arcaneBlastRank.ID},
		SpellSchool:    arcaneBlastRank.SpellSchool(),
		DefenseType:    arcaneBlastRank.DefenseTypeCore(),
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MageSpellArcaneBlast,
		MaxRange:       float64(arcaneBlastRank.MaxRange),

		ManaCost: arcaneBlastRank.ManaCost(),
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      arcaneBlastRank.GCD(),
				CastTime: arcaneBlastRank.CastTime(),
			},
		},

		DamageMultiplier: 1,
		BonusCoefficient: arcaneBlastRank.DamageEffect().Coeff(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, arcaneBlastRank.DamageEffect().Roll(sim, core.CharacterLevel), spell.OutcomeMagicHitAndCrit)
			// The stacks run 15 sec from the first one, so a new stack doesn't refresh them.
			if !mage.ArcaneBlastAura.IsActive() {
				mage.ArcaneBlastAura.Activate(sim)
			}
			mage.ArcaneBlastAura.AddStack(sim)
		},
	})
}
