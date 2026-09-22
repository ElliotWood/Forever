package warlock

import (
	"github.com/wowsims/forever/sim/core"
)

var hellfireRank = spellData.Hellfire.Highest()

// The dump's Periodic role for this family sits at effect position 2: position 1 is the
// A_PERIODIC_TRIGGER_SPELL that fires the Hellfire Effect spell each tick, not the damage tick.
var hellfireTick = hellfireRank.EffectN(2)
var hellFireCoeff = hellfireTick.Coeff()

// TODO: To be implemented. Port the TBC Hellfire implementation below; not yet verified against the Forever client.
func (warlock *Warlock) registerHellfire() *core.Spell {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// hellfireActionID := core.ActionID{SpellID: hellfireRank.ID}
	// tickLength := hellfireTick.Period()
	//
	// manaCost := hellfireRank.Cost()
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
	// 			GCD: hellfireRank.GCD(),
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
	// 		TickLength:           tickLength,
	// 		NumberOfTicks:        int32(hellfireRank.Duration() / tickLength),
	// 		HasteReducesDuration: true,
	// 		AffectedByCastSpeed:  true,
	// 		BonusCoefficient:     hellFireCoeff,
	//
	// 		OnTick: func(sim *core.Simulation, _ *core.Unit, dot *core.Dot) {
	// 			// Rolled once: the warlock burns exactly what it deals.
	// 			tickDamage := hellfireTick.Average(core.CharacterLevel)
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
