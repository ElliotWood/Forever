package mage

import (
	"github.com/wowsims/forever/sim/core"
)

var scorchRank = spellData.Scorch.HighestRank()

func (mage *Mage) registerScorchSpell() {

	procChance := []float64{0, 0.33, 0.66, 1}[mage.Talents.ImprovedScorch]

	mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: scorchRank.SpellID},
		SpellSchool:    scorchRank.SpellSchool,
		DefenseType:    scorchRank.DefenseType,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MageSpellScorch,

		ManaCost: core.ManaCostOptions{
			FlatCost: scorchRank.Cost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      scorchRank.GCD,
				CastTime: scorchRank.CastTime,
			},
		},

		DamageMultiplierAdditive: 1,
		BonusCoefficient:         scorchRank.Direct.BonusCoefficient(),
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := scorchRank.Direct.Damage(sim)
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if result.Landed() && mage.Talents.ImprovedScorch > 0 {
				if sim.Proc(procChance, "Improved Scorch") {
					aura := mage.ImprovedScorchAuras.Get(target)
					aura.Activate(sim)
					aura.AddStack(sim)
				}
			}
		},

		RelatedAuraArrays: mage.ImprovedScorchAuras.ToMap(),
	})
}
