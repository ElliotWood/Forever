package mage

import (
	"github.com/wowsims/classic/sim/core"
)

// Damage bonus against a target the mage counts as frozen, i.e. one held by Fingers of Frost.
const IceLanceFrozenMultiplier = 4.0

func (mage *Mage) registerIceLanceSpell() {
	if !mage.Talents.IceLance {
		return
	}

	// Level 60 values aren't datamined yet. The demo tooltip reads 28 to 33 for what is a rank 1,
	// so these keep Ice Lance at roughly a fifth of the level 60 Frostbolt it is cast alongside.
	baseDamage := []float64{115, 133}
	spellCoeff := .143
	manaCost := 160.0

	mage.IceLance = mage.RegisterSpell(core.SpellConfig{
		SpellCode:    SpellCode_MageIceLance,
		ActionID:     core.ActionID{SpellID: 30455},
		SpellSchool:  core.SpellSchoolFrost,
		DefenseType:  core.DefenseTypeMagic,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        SpellFlagMage | core.SpellFlagBinary | core.SpellFlagAPL,
		MissileSpeed: 38,

		RequiredLevel: 60,
		Rank:          1,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: spellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := sim.Roll(baseDamage[0], baseDamage[1])
			result := spell.CalcDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)

			// The tooltip reads as a bonus on the whole hit rather than on the base roll, so spell
			// power gets multiplied too.
			if mage.IsTargetFrozen() {
				result.Damage *= IceLanceFrozenMultiplier
			}

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
