package druid

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

var shredRank = spellData.Shred.HighestRank()

func (druid *Druid) registerShredSpell() {
	druid.Shred = druid.RegisterSpell(Cat, core.SpellConfig{
		ActionID:       core.ActionID{SpellID: shredRank.SpellID},
		SpellSchool:    shredRank.SpellSchool,
		DefenseType:    shredRank.DefenseType,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		ClassSpellMask: DruidSpellShred,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		EnergyCost: core.EnergyCostOptions{
			Cost:   shredRank.Cost,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: shredRank.GCD,
			},
			IgnoreHaste: true,
		},

		ExtraCastCondition: func(_ *core.Simulation, _ *core.Unit) bool {
			return !druid.PseudoStats.InFrontOfTarget && !druid.CannotShredTarget
		},

		// Weapon damage * 2.25, boosted by 30% if Mangle is active.
		//
		// TODO: the client disagrees with this multiplier. Every rank of spellData.Shred
		// states E_WEAPON_PERCENT_DAMAGE 155, not 225, so Forever appears to have retuned
		// it. Left at the TBC value rather than changed on one reading of the data.
		DamageMultiplier: 2.25,
		ThreatMultiplier: 1,
		MaxRange:         core.MaxMeleeRange,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// shredRank.Direct is the pre-multiplier flat value (180); the idol/flat
			// bonuses are historically expressed post-multiplier, so scale up and back
			// down around them to keep the result identical.
			baseDamage := shredRank.Direct.Damage(sim) + spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower(target))
			if druid.MangleAuras != nil && druid.MangleAuras.Get(target).IsActive() {
				baseDamage *= 1.3
			}

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if result.Landed() {
				druid.AddComboPoints(sim, 1, spell.ComboPointMetrics())
			} else {
				spell.IssueRefund(sim)
			}
		},

		ExpectedInitialDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			baseDamage := shared.SpellDataMin(shredRank.Direct) + spell.Unit.AutoAttacks.MH().CalculateAverageWeaponDamage(spell.MeleeAttackPower(target))
			if druid.MangleAuras != nil && druid.MangleAuras.Get(target).IsActive() {
				baseDamage *= 1.3
			}
			return spell.CalcDamage(sim, target, baseDamage, spell.OutcomeExpectedMeleeWeaponSpecialHitAndCrit)
		},
	})
}

func (druid *Druid) CurrentShredCost() float64 {
	return druid.Shred.Cost.GetCurrentCost()
}
