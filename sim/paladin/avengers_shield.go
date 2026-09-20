package paladin

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

// Package-level state the commented-out implementations used:
// var AvengersShieldRankMap = spellData.AvengersShield

func (paladin *Paladin) getAvengersShieldTimer() *core.Timer {
	if paladin.avengersShieldTimer == nil {
		paladin.avengersShieldTimer = paladin.NewTimer()
	}
	return paladin.avengersShieldTimer
}

// Avenger's Shield (Talent)
// https://www.wowhead.com/forever/spell=31935
//
// Hurls a holy shield at the enemy, dealing Holy damage, dazing them and
// then jumping to additional nearby enemies. Affects 3 total targets.
//
// TODO: uncalled -- Forever drops the Avenger's Shield talent; re-gate before wiring
// back into registerTalentSpells.
// TODO: To be implemented. Spells of this name exist in the Forever client, but none of them
// has a class ability row -- no SkillLineAbility entry in a CategoryID 7 skill line -- so the
// generator has no class spell to build a ladder from.
func (paladin *Paladin) registerAvengersShield(rankConfig shared.SpellData) {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// paladin.RegisterSpell(core.SpellConfig{
	// 	ActionID:       core.ActionID{SpellID: rankConfig.SpellID},
	// 	SpellSchool:    core.SpellSchoolHoly,
	// 	DefenseType:    core.DefenseTypeRanged,
	// 	ProcMask:       core.ProcMaskRangedSpecial,
	// 	Flags:          core.SpellFlagAPL | core.SpellFlagBinary,
	// 	ClassSpellMask: SpellMaskAvengersShield,
	// 	Rank:           rankConfig.Rank,
	//
	// 	DamageMultiplier: 1,
	// 	ThreatMultiplier: 1,
	//
	// 	MaxRange:     rankConfig.MaxRange,
	// 	MissileSpeed: rankConfig.MissileSpeed,
	//
	// 	ManaCost: core.ManaCostOptions{
	// 		FlatCost: rankConfig.Cost,
	// 	},
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			GCD:      rankConfig.GCD,
	// 			CastTime: rankConfig.CastTime,
	// 		},
	// 		CD: core.Cooldown{
	// 			Timer:    paladin.getAvengersShieldTimer(),
	// 			Duration: rankConfig.Cooldown,
	// 		},
	// 		ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
	// 			castTime := paladin.ApplyCastSpeedForSpell(cast.CastTime, spell)
	// 			paladin.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime+castTime)
	// 		},
	// 	},
	//
	// 	BonusCoefficient: rankConfig.Direct.BonusCoefficient(),
	//
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		damage := rankConfig.Direct.Damage(sim)
	// 		results := spell.CalcCleaveDamage(sim, target, 3, damage, spell.OutcomeRangedHitAndCrit)
	// 		spell.WaitTravelTime(sim, func(sim *core.Simulation) {
	// 			for _, result := range results {
	// 				spell.DealDamage(sim, result)
	// 			}
	// 		})
	// 	},
	// })
}
