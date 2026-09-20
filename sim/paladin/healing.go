package paladin

import (
	"github.com/wowsims/forever/sim/common/shared"
)

func (paladin *Paladin) registerHealingSpells() {
	HolyLightRankMap.RegisterAll(paladin.registerHolyLight)
	FlashOfLightRankMap.RegisterAll(paladin.registerFlashOfLight)
	LayOnHandsRankMap.RegisterAll(paladin.registerLayOnHands)
}

var HolyLightRankMap = spellData.HolyLight

// TODO: To be implemented. TBC body below needs no porting; kept commented until this class's port is reviewed.
//
// Holy Light
// https://www.wowhead.com/forever/spell=27136
//
// Heals a friendly target for a large amount.
func (paladin *Paladin) registerHolyLight(rankConfig shared.SpellData) {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// spellID := rankConfig.SpellID
	// cost := rankConfig.Cost
	// coefficient := rankConfig.Heal.BonusCoefficient()
	//
	// paladin.RegisterSpell(core.SpellConfig{
	// 	ActionID:       core.ActionID{SpellID: spellID},
	// 	SpellSchool:    core.SpellSchoolHoly,
	// 	DefenseType:    core.DefenseTypeMagic,
	// 	ProcMask:       core.ProcMaskSpellHealing,
	// 	Flags:          core.SpellFlagAPL | core.SpellFlagHelpful,
	// 	Rank:           rankConfig.Rank,
	// 	ClassSpellMask: SpellMaskHolyLight,
	//
	// 	DamageMultiplier: 1,
	// 	ThreatMultiplier: 1,
	//
	// 	MaxRange: rankConfig.MaxRange,
	//
	// 	ManaCost: core.ManaCostOptions{
	// 		FlatCost: cost,
	// 	},
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			GCD:      rankConfig.GCD,
	// 			CastTime: rankConfig.CastTime,
	// 		},
	// 	},
	//
	// 	BonusCoefficient: coefficient,
	//
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		spell.CalcAndDealHealing(sim, target, rankConfig.Heal.Damage(sim), spell.OutcomeHealingCrit)
	// 	},
	// })
}

var FlashOfLightRankMap = spellData.FlashOfLight

// TODO: To be implemented. TBC body below needs no porting; kept commented until this class's port is reviewed.
//
// Flash of Light
// https://www.wowhead.com/forever/spell=27137
//
// Heals a friendly target for a small amount.
func (paladin *Paladin) registerFlashOfLight(rankConfig shared.SpellData) {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// spellID := rankConfig.SpellID
	// cost := rankConfig.Cost
	// coefficient := rankConfig.Heal.BonusCoefficient()
	//
	// paladin.RegisterSpell(core.SpellConfig{
	// 	ActionID:       core.ActionID{SpellID: spellID},
	// 	SpellSchool:    core.SpellSchoolHoly,
	// 	DefenseType:    core.DefenseTypeMagic,
	// 	ProcMask:       core.ProcMaskSpellHealing,
	// 	Flags:          core.SpellFlagAPL | core.SpellFlagHelpful,
	// 	Rank:           rankConfig.Rank,
	// 	ClassSpellMask: SpellMaskFlashOfLight,
	//
	// 	DamageMultiplier: 1,
	// 	ThreatMultiplier: 1,
	//
	// 	MaxRange: rankConfig.MaxRange,
	//
	// 	ManaCost: core.ManaCostOptions{
	// 		FlatCost: cost,
	// 	},
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			GCD:      rankConfig.GCD,
	// 			CastTime: rankConfig.CastTime,
	// 		},
	// 	},
	//
	// 	BonusCoefficient: coefficient,
	//
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		spell.CalcAndDealHealing(sim, target, rankConfig.Heal.Damage(sim), spell.OutcomeHealingCrit)
	// 	},
	// })
}

var LayOnHandsRankMap = spellData.LayOnHands

// TODO: To be implemented. TBC body below needs no porting; kept commented until this class's port is reviewed.
//
// Lay on Hands
// https://www.wowhead.com/forever/spell=27154
//
// Heals a friendly target for an amount equal to the Paladin's maximum health
// and restores mana to the target. Causes Forbearance for 1 min.
func (paladin *Paladin) registerLayOnHands(rankConfig shared.SpellData) {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// spellID := rankConfig.SpellID
	// manaRestore := shared.SpellDataMin(rankConfig.Energize)
	//
	// cd := core.Cooldown{
	// 	Timer:    paladin.NewTimer(),
	// 	Duration: time.Hour,
	// }
	//
	// manaMetrics := paladin.NewManaMetrics(core.ActionID{SpellID: spellID})
	//
	// paladin.RegisterSpell(core.SpellConfig{
	// 	ActionID:       core.ActionID{SpellID: spellID},
	// 	SpellSchool:    core.SpellSchoolHoly,
	// 	DefenseType:    core.DefenseTypeMagic,
	// 	ProcMask:       core.ProcMaskSpellHealing,
	// 	Flags:          core.SpellFlagAPL | core.SpellFlagHelpful,
	// 	Rank:           rankConfig.Rank,
	// 	ClassSpellMask: SpellMaskLayOnHands,
	//
	// 	DamageMultiplier: 1,
	// 	ThreatMultiplier: 1,
	//
	// 	MaxRange: rankConfig.MaxRange,
	//
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			GCD: core.GCDDefault,
	// 		},
	// 		CD: cd,
	// 	},
	//
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		// Drain all of the caster's mana
	// 		spell.Unit.AddMana(sim, -spell.Unit.CurrentMana(), manaMetrics)
	//
	// 		// Restore mana and health to the target
	// 		target.AddMana(sim, manaRestore, manaMetrics)
	// 		spell.CalcAndDealHealing(sim, target, spell.Unit.MaxHealth(), spell.OutcomeHealingCrit)
	// 	},
	// })
}
