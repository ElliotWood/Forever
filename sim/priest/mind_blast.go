package priest

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

var MindBlastRankMap = spellData.MindBlast

// TODO: To be implemented. Mind Blast already has a full Forever rank ladder (spellData.MindBlast); the
// TBC body needs review before it's uncommented.
func (priest *Priest) registerMindBlastSpell(rank shared.SpellData, cdTimer *core.Timer) {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	//
	// priest.RegisterSpell(core.SpellConfig{
	// 	ActionID:       core.ActionID{SpellID: rank.SpellID},
	// 	SpellSchool:    core.SpellSchoolShadow,
	// 	DefenseType:    core.DefenseTypeMagic,
	// 	ProcMask:       core.ProcMaskSpellDamage,
	// 	Flags:          core.SpellFlagAPL,
	// 	ClassSpellMask: PriestSpellMindBlast,
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
	// 		CD: core.Cooldown{
	// 			Timer:    cdTimer,
	// 			Duration: rank.Cooldown,
	// 		},
	// 	},
	//
	// 	DamageMultiplier:         1,
	// 	DamageMultiplierAdditive: 1,
	// 	BonusCoefficient:         rank.Direct.BonusCoefficient(),
	// 	ThreatMultiplier:         1,
	//
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		baseDamage := rank.Direct.Damage(sim)
	// 		result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
	// 		spell.DealDamage(sim, result)
	// 	},
	// })
}
