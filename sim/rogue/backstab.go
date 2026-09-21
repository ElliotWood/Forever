package rogue

import (
	"github.com/wowsims/forever/sim/core"
)

var backstabRank = spellData.Backstab.HighestRank()

func (rogue *Rogue) registerBackstabSpell() {
	baseDamage, _ := backstabRank.Direct.Range()
	weaponDamage := spellData.Backstab.EffectAt(1).ValueAt(backstabRank.Rank) / 100

	// Puncturing Wounds also hands a combo point back, on effect 1 of the talent.
	extraComboPointChance := spellData.PuncturingWounds.EffectAt(1).ValueAt(rogue.Talents.PuncturingWounds) / 100
	cpMetrics := rogue.NewComboPointMetrics(core.ActionID{SpellID: spellData.PuncturingWoundsTriggered.HighestRank().SpellID})

	rogue.Backstab = rogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: backstabRank.SpellID},
		SpellSchool:    backstabRank.SpellSchool,
		DefenseType:    backstabRank.DefenseType,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | SpellFlagBuilder | core.SpellFlagAPL,
		ClassSpellMask: RogueSpellBackstab,
		MaxRange:       core.MaxMeleeRange,

		EnergyCost: core.EnergyCostOptions{
			Cost:   backstabRank.Cost,
			Refund: backstabRank.MissRefund(),
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

		DamageMultiplier:         weaponDamage,
		DamageMultiplierAdditive: 1,
		ThreatMultiplier:         1,

		BonusCoefficient: backstabRank.Direct.BonusCoefficient(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)

			damage := baseDamage + spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower(target))
			result := spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if result.Landed() {
				rogue.AddComboPoints(sim, 1, spell.ComboPointMetrics())
				if extraComboPointChance > 0 && sim.Proc(extraComboPointChance, "Puncturing Wounds") {
					rogue.AddComboPoints(sim, 1, cpMetrics)
				}
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}
