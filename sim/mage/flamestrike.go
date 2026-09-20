package mage

import (
	"github.com/wowsims/forever/sim/common/shared"
)

// TODO: was Ranks(7, 6); Forever's Flamestrike tops out at rank 6, so only that remains.
var FlameStrikeRankMap = spellData.Flamestrike.Ranks(6)

// TODO: To be implemented. TBC body below needs no porting (Forever's dropped DoT is already left out, per the TODO inside); kept commented until this class's port is reviewed.
func (mage *Mage) registerFlamestrike(rankConfig shared.SpellData) {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// // TODO: Forever moves Flamestrike's damage over time onto the area trigger its second
	// // effect creates, which the client tables do not carry, so the rank has no Periodic value
	// // and only the direct hit is registered - no DoT rather than an invented one.
	// flameStrikeCoefficient := 0.23600000143 // Per https://wago.tools/db2/SpellEffect?build=2.5.5.65295&filter%5BSpellID%5D=exact%253A2120 Field: "BonusCoefficient"
	//
	// spell := mage.RegisterSpell(core.SpellConfig{
	// 	ActionID:       core.ActionID{SpellID: rankConfig.SpellID},
	// 	SpellSchool:    core.SpellSchoolFire,
	// 	DefenseType:    core.DefenseTypeMagic,
	// 	ProcMask:       core.ProcMaskSpellDamage,
	// 	Flags:          core.SpellFlagAPL,
	// 	ClassSpellMask: MageSpellFlamestrike,
	// 	Rank:           rankConfig.Rank,
	//
	// 	ManaCost: core.ManaCostOptions{
	// 		FlatCost: rankConfig.Cost,
	// 	},
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			GCD:      rankConfig.GCD,
	// 			CastTime: rankConfig.CastTime,
	// 		},
	// 	},
	//
	// 	DamageMultiplier: 1,
	// 	BonusCoefficient: flameStrikeCoefficient,
	// 	ThreatMultiplier: 1,
	//
	// 	ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
	// 		baseDamage := rankConfig.Direct.Damage(sim)
	// 		spell.CalcAndDealAoeDamage(sim, baseDamage, spell.OutcomeMagicHitAndCrit)
	// 	},
	// })
	//
	// mage.Flamestrike = append(mage.Flamestrike, spell)
}
