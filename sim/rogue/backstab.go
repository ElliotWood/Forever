package rogue

import (
	"github.com/wowsims/forever/sim/core"
)

var backstabRank = spellData.Backstab.BySpellID(26863)

func (rogue *Rogue) registerBackstabSpell() {
	baseDamage, _ := backstabRank.Direct.Range()
	weaponDamage := 1.5

	rogue.Backstab = rogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: backstabRank.SpellID},
		SpellSchool:    backstabRank.SpellSchool,
		DefenseType:    backstabRank.DefenseType,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | SpellFlagBuilder | core.SpellFlagAPL,
		ClassSpellMask: RogueSpellBackstab,

		EnergyCost: core.EnergyCostOptions{
			Cost:   backstabRank.Cost,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: backstabRank.GCD,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !rogue.PseudoStats.InFrontOfTarget && rogue.HasDagger(core.MainHand)
		},

		DamageMultiplierAdditive: weaponDamage,
		DamageMultiplier:         1,
		ThreatMultiplier:         1,

		BonusCoefficient: backstabRank.Direct.BonusCoefficient(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			baseDamage := baseDamage +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower(target))

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if result.Landed() {
				rogue.AddComboPoints(sim, 1, spell.ComboPointMetrics())
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}
