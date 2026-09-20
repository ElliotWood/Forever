package priest

import (
	"github.com/wowsims/forever/sim/common/shared"
)

var SmiteRankMap = spellData.Smite

// TODO: To be implemented. Smite already has a full Forever rank ladder (spellData.Smite); the TBC body
// needs review before it's uncommented.
func (priest *Priest) registerSmiteSpell(rank shared.SpellData) {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// priest.RegisterSpell(core.SpellConfig{
	// 	ActionID:       core.ActionID{SpellID: rank.SpellID},
	// 	SpellSchool:    core.SpellSchoolHoly,
	// 	DefenseType:    core.DefenseTypeMagic,
	// 	ProcMask:       core.ProcMaskSpellDamage,
	// 	Flags:          core.SpellFlagAPL,
	// 	ClassSpellMask: PriestSpellSmite,
	// 	Rank:           rank.Rank,
	// 	MaxRange:       rank.MaxRange,
	// 	ManaCost: core.ManaCostOptions{
	// 		FlatCost: rank.Cost,
	// 	},
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			GCD:      rank.GCD,
	// 			CastTime: rank.CastTime,
	// 		},
	// 	},
	//
	// 	DamageMultiplier: 1,
	// 	BonusCoefficient: rank.Direct.BonusCoefficient(),
	// 	ThreatMultiplier: 1,
	//
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		baseDamage := rank.Direct.Damage(sim)
	// 		spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
	// 	},
	// })
}
