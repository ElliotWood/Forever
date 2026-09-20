package druid

import (
	"github.com/wowsims/forever/sim/core"
)

// Package-level state the commented-out implementations used:
// var mangleCatRank = spellData.MangleCat.BySpellID(33983)
// var mangleBearRank = spellData.MangleBear.BySpellID(33987)

func (druid *Druid) registerMangleAuras() {
	if druid.MangleAuras != nil {
		return
	}
	druid.MangleAuras = druid.NewEnemyAuraArray(core.MangleAura)
}

// TODO: To be implemented, or removed. The Forever client does not ship this spell at all --
// no SpellName row carries the name -- so there is nothing to build a registrar from. The
// body below is kept commented as the record of the TBC implementation.
func (druid *Druid) registerMangleCatSpell() {
	if !druid.Talents.Mangle {
		return
	}
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// if !druid.Talents.Mangle {
	// 	return
	// }
	//
	// druid.registerMangleAuras()
	//
	// druid.MangleCat = druid.RegisterSpell(Cat, core.SpellConfig{
	// 	ActionID:       core.ActionID{SpellID: mangleCatRank.SpellID},
	// 	SpellSchool:    mangleCatRank.SpellSchool,
	// 	DefenseType:    mangleCatRank.DefenseType,
	// 	ProcMask:       core.ProcMaskMeleeMHSpecial,
	// 	ClassSpellMask: DruidSpellMangleCat,
	// 	Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
	//
	// 	EnergyCost: core.EnergyCostOptions{
	// 		Cost:   mangleCatRank.Cost,
	// 		Refund: 0.8,
	// 	},
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			GCD: mangleCatRank.GCD,
	// 		},
	// 		IgnoreHaste: true,
	// 	},
	//
	// 	DamageMultiplier: 1.6,
	// 	ThreatMultiplier: 1,
	// 	MaxRange:         core.MaxMeleeRange,
	//
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		// mangleCatRank.Direct is the pre-multiplier flat value (165); the idol
	// 		// bonus is historically expressed post-multiplier, so scale up and back
	// 		// down around it to keep the result identical.
	// 		baseDamage := (mangleCatRank.Direct.Damage(sim)*1.6+druid.IdolMangleCatBonus)/1.6 + spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower(target))
	// 		result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
	//
	// 		if result.Landed() {
	// 			druid.AddComboPoints(sim, 1, spell.ComboPointMetrics())
	// 			druid.MangleAuras.Get(target).Activate(sim)
	// 		} else {
	// 			spell.IssueRefund(sim)
	// 		}
	// 	},
	//
	// 	ExpectedInitialDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
	// 		baseDamage := (shared.SpellDataMin(mangleCatRank.Direct)*1.6+druid.IdolMangleCatBonus)/1.6 + spell.Unit.AutoAttacks.MH().CalculateAverageWeaponDamage(spell.MeleeAttackPower(target))
	// 		return spell.CalcDamage(sim, target, baseDamage, spell.OutcomeExpectedMeleeWeaponSpecialHitAndCrit)
	// 	},
	// })
}

// TODO: To be implemented, or removed. The Forever client does not ship this spell at all --
// no SpellName row carries the name -- so there is nothing to build a registrar from. The
// body below is kept commented as the record of the TBC implementation.
func (druid *Druid) registerMangleBearSpell() {
	if !druid.Talents.Mangle {
		return
	}
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// if !druid.Talents.Mangle {
	// 	return
	// }
	//
	// druid.registerMangleAuras()
	//
	// druid.MangleBear = druid.RegisterSpell(Bear, core.SpellConfig{
	// 	ActionID:       core.ActionID{SpellID: mangleBearRank.SpellID},
	// 	SpellSchool:    mangleBearRank.SpellSchool,
	// 	DefenseType:    mangleBearRank.DefenseType,
	// 	ProcMask:       core.ProcMaskMeleeMHSpecial,
	// 	ClassSpellMask: DruidSpellMangleBear,
	// 	Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
	//
	// 	RageCost: core.RageCostOptions{
	// 		Cost:   mangleBearRank.Cost,
	// 		Refund: 0.8,
	// 	},
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			GCD: mangleBearRank.GCD,
	// 		},
	// 		IgnoreHaste: true,
	// 		CD: core.Cooldown{
	// 			Timer:    druid.NewTimer(),
	// 			Duration: mangleBearRank.Cooldown,
	// 		},
	// 	},
	//
	// 	DamageMultiplier: 1.15,
	// 	ThreatMultiplier: 1.5 / 1.15,
	// 	MaxRange:         core.MaxMeleeRange,
	//
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		baseDamage := (mangleBearRank.Direct.Damage(sim)*1.15+druid.IdolMangleBearBonus)/1.15 + spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower(target))
	// 		result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
	//
	// 		if result.Landed() {
	// 			druid.MangleAuras.Get(target).Activate(sim)
	// 		} else {
	// 			spell.IssueRefund(sim)
	// 		}
	// 	},
	// })
}

func (druid *Druid) CurrentMangleCatCost() float64 {
	return druid.MangleCat.Cost.GetCurrentCost()
}
