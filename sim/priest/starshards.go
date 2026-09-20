package priest

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

// Starshards - Night Elf Racial
// Arcane school DoT, 0 mana cost, 30s cooldown, 15s duration
var StarshardsRankMap = spellData.Starshards

// TODO: To be implemented. Starshards already has a full Forever rank ladder (spellData.Starshards); the
// TBC body needs review before it's uncommented.
func (priest *Priest) registerStarshardsSpell(rank shared.SpellData, cdTimer *core.Timer) {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// tick := rank.Periodic.(shared.SpellDataPeriodic)
	//
	// priest.RegisterSpell(core.SpellConfig{
	// 	ActionID:       core.ActionID{SpellID: rank.SpellID},
	// 	SpellSchool:    core.SpellSchoolArcane,
	// 	DefenseType:    core.DefenseTypeMagic,
	// 	ProcMask:       core.ProcMaskSpellDamage,
	// 	Flags:          core.SpellFlagAPL,
	// 	ClassSpellMask: PriestSpellStarshards,
	// 	Rank:           rank.Rank,
	// 	MaxRange:       rank.MaxRange,
	//
	// 	ManaCost: core.ManaCostOptions{
	// 		FlatCost: 0,
	// 	},
	//
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			GCD: rank.GCD,
	// 		},
	// 		CD: core.Cooldown{
	// 			Timer:    cdTimer,
	// 			Duration: rank.Cooldown,
	// 		},
	// 	},
	//
	// 	DamageMultiplier:         1,
	// 	DamageMultiplierAdditive: 1,
	// 	ThreatMultiplier:         1,
	//
	// 	Dot: core.DotConfig{
	// 		Aura: core.Aura{
	// 			Label: fmt.Sprintf("Starshards-%d", rank.Rank),
	// 		},
	// 		NumberOfTicks:       tick.NumberOfTicks,
	// 		TickLength:          tick.TickLength,
	// 		AffectedByCastSpeed: false,
	// 		BonusCoefficient:    tick.Coef,
	//
	// 		OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
	// 			dot.Snapshot(target, tick.Tick)
	// 		},
	// 		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
	// 			dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
	// 		},
	// 	},
	//
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)
	// 		if result.Landed() {
	// 			spell.Dot(target).Apply(sim)
	// 		}
	// 		spell.DealOutcome(sim, result)
	// 	},
	//
	// 	ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
	// 		if useSnapshot {
	// 			dot := spell.Dot(target)
	// 			return dot.CalcSnapshotDamage(sim, target, spell.OutcomeExpectedMagicHit)
	// 		}
	// 		return spell.CalcPeriodicDamage(sim, target, tick.Tick, spell.OutcomeExpectedMagicHit)
	// 	},
	// })
}
