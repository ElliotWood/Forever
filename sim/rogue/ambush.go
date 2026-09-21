package rogue

import (
	"github.com/wowsims/forever/sim/core"
)

var ambushRank = spellData.Ambush.HighestRank()

func (rogue *Rogue) registerAmbushSpell() {
	baseDamage, _ := ambushRank.Direct.Range()
	// The client states the weapon share as a percentage on effect 1 (250, where TBC had 275).
	// Effects 0 and 2 share its aura and misc pair, so the effect has to be named by index.
	weaponDamage := spellData.Ambush.EffectAt(1).ValueAt(ambushRank.Rank) / 100

	rogue.Ambush = rogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: ambushRank.SpellID},
		SpellSchool:    ambushRank.SpellSchool,
		DefenseType:    ambushRank.DefenseType,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | SpellFlagBuilder | core.SpellFlagAPL,
		ClassSpellMask: RogueSpellAmbush,
		MaxRange:       core.MaxMeleeRange,

		EnergyCost: core.EnergyCostOptions{
			Cost:   ambushRank.Cost,
			Refund: ambushRank.MissRefund(),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: ambushRank.GCD,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			if !rogue.HasDagger(core.MainHand) {
				return false
			}
			// Cutthroat lets the Stealth requirement slide for a short while after a Backstab.
			return rogue.IsStealthed() || (rogue.CutthroatAura != nil && rogue.CutthroatAura.IsActive())
		},

		DamageMultiplier:         weaponDamage,
		DamageMultiplierAdditive: 1,
		ThreatMultiplier:         1,

		BonusCoefficient: ambushRank.Direct.BonusCoefficient(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			if rogue.CutthroatAura != nil {
				rogue.CutthroatAura.Deactivate(sim)
			}

			damage := baseDamage + spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower(target))
			result := spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMeleeSpecialNoBlockDodgeParry)

			if result.Landed() {
				rogue.AddComboPoints(sim, 1, spell.ComboPointMetrics())
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}
