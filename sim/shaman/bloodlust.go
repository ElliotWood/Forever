package shaman

import (
	"github.com/wowsims/forever/sim/core"
)

// Package-level state the commented-out implementations used:
// var bloodlustRank = spellData.Bloodlust.BySpellID(2825)

func (shaman *Shaman) BloodlustActionID() core.ActionID {
	return core.ActionID{
		SpellID: 2825,
		Tag:     shaman.Index,
	}
}

// TODO: To be implemented. Spells of this name exist in the Forever client, but none of them
// has a class ability row -- no SkillLineAbility entry in a CategoryID 7 skill line -- so the
// generator has no class spell to build a ladder from.
func (shaman *Shaman) registerBloodlustCD() {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// actionID := shaman.BloodlustActionID()
	//
	// blAuras := []*core.Aura{}
	// for _, party := range shaman.Env.Raid.Parties {
	// 	for _, partyMember := range party.Players {
	// 		blAuras = append(blAuras, core.BloodlustAura(partyMember.GetCharacter(), actionID.Tag))
	// 	}
	// }
	//
	// spell := shaman.RegisterSpell(core.SpellConfig{
	// 	ActionID:       actionID,
	// 	Flags:          core.SpellFlagAPL | SpellFlagInstant,
	// 	ClassSpellMask: SpellMaskBloodlust,
	//
	// 	ManaCost: core.ManaCostOptions{
	// 		FlatCost: bloodlustRank.Cost,
	// 	},
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			GCD: bloodlustRank.GCD,
	// 		},
	// 		CD: core.Cooldown{
	// 			Timer:    shaman.NewTimer(),
	// 			Duration: core.BloodlustCD,
	// 		},
	// 	},
	//
	// 	ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
	// 		for _, blAura := range blAuras {
	// 			target := blAura.Unit
	// 			// Only activate bloodlust on units without sated.
	// 			if !target.HasActiveAura(core.SatedAuraLabel) {
	// 				blAura.Activate(sim)
	// 			}
	// 		}
	// 	},
	// })
	//
	// shaman.AddMajorCooldown(core.MajorCooldown{
	// 	Spell:    spell,
	// 	Type:     core.CooldownTypeDPS,
	// 	Priority: core.CooldownPriorityBloodlust,
	// 	ShouldActivate: func(_ *core.Simulation, _ *core.Character) bool {
	// 		// Only cast if there is a player missing Sated.
	// 		for _, playerUnit := range shaman.Env.Raid.AllPlayerUnits {
	// 			if !playerUnit.HasActiveAura(core.SatedAuraLabel) {
	// 				return true
	// 			}
	// 		}
	// 		return false
	// 	},
	// })
}
