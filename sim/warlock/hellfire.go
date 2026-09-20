package warlock

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

var hellfireRank = spellData.Hellfire.HighestRank()
var hellfireTick = hellfireRank.Periodic.(shared.SpellDataPeriodic)
var hellFireCoeff = hellfireTick.Coef

// TODO: To be implemented. Port the TBC Hellfire implementation below; not yet verified against the Forever client.
func (warlock *Warlock) registerHellfire() *core.Spell {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// hellfireActionID := core.ActionID{SpellID: hellfireRank.SpellID}
	//
	// manaCost := hellfireRank.Cost
	// warlock.Hellfire = warlock.RegisterSpell(core.SpellConfig{
	// 	ActionID:         hellfireActionID,
	// 	SpellSchool:      core.SpellSchoolFire,
	// 	DefenseType:      core.DefenseTypeMagic,
	// 	Flags:            core.SpellFlagChanneled | core.SpellFlagAPL,
	// 	ProcMask:         core.ProcMaskSpellDamage,
	// 	ClassSpellMask:   WarlockSpellHellfire,
	// 	ThreatMultiplier: 1,
	// 	DamageMultiplier: 1,
	//
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			GCD: hellfireRank.GCD,
	// 		},
	// 	},
	// 	ManaCost: core.ManaCostOptions{FlatCost: manaCost},
	//
	// 	Dot: core.DotConfig{
	// 		Aura: core.Aura{
	// 			Label: "Hellfire",
	// 		},
	//
	// 		IsAOE:                true,
	// 		TickLength:           hellfireTick.TickLength,
	// 		NumberOfTicks:        hellfireTick.NumberOfTicks,
	// 		HasteReducesDuration: true,
	// 		AffectedByCastSpeed:  true,
	// 		BonusCoefficient:     hellFireCoeff,
	//
	// 		OnTick: func(sim *core.Simulation, _ *core.Unit, dot *core.Dot) {
	// 			// Rolled once: the warlock burns exactly what it deals.
	// 			tickDamage := hellfireTick.Damage(sim)
	//
	// 			resultSlice := dot.Spell.CalcPeriodicAoeDamage(sim, tickDamage, dot.Spell.OutcomeTickMagicHitNoHitCounter)
	// 			if resultSlice[0].Damage > warlock.CurrentHealth() {
	// 				dot.Deactivate(sim)
	// 			}
	//
	// 			dot.Spell.DealBatchedPeriodicDamage(sim)
	// 			warlock.RemoveHealth(sim, tickDamage)
	//
	// 		},
	// 	},
	//
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		spell.AOEDot().Apply(sim)
	// 	},
	// })
	//
	// return warlock.Hellfire
}
