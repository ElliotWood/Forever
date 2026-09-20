package druid

import (
	"github.com/wowsims/forever/sim/common/shared"
)

var insectSwarmRank = spellData.InsectSwarm.HighestRank()
var insectSwarmTick = insectSwarmRank.Periodic.(shared.SpellDataPeriodic)

// TODO: To be implemented.
func (druid *Druid) registerInsectSwarmSpell() {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// druid.InsectSwarm = druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
	// 	ActionID:       core.ActionID{SpellID: insectSwarmRank.SpellID},
	// 	SpellSchool:    insectSwarmRank.SpellSchool,
	// 	DefenseType:    insectSwarmRank.DefenseType,
	// 	ProcMask:       core.ProcMaskSpellDamage,
	// 	ClassSpellMask: DruidSpellInsectSwarm,
	// 	Flags:          core.SpellFlagAPL | core.SpellFlagBinary,
	//
	// 	DamageMultiplier: 1,
	// 	ThreatMultiplier: 1,
	// 	MaxRange:         insectSwarmRank.MaxRange,
	//
	// 	ManaCost: core.ManaCostOptions{
	// 		FlatCost: insectSwarmRank.Cost,
	// 	},
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			GCD: insectSwarmRank.GCD,
	// 		},
	// 	},
	//
	// 	Dot: core.DotConfig{
	// 		Aura: core.Aura{
	// 			Label: "Insect Swarm",
	// 		},
	//
	// 		NumberOfTicks:       insectSwarmTick.NumberOfTicks,
	// 		TickLength:          insectSwarmTick.TickLength,
	// 		AffectedByCastSpeed: false,
	// 		BonusCoefficient:    insectSwarmTick.Coef,
	//
	// 		OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
	// 			dot.Snapshot(target, insectSwarmTick.Tick)
	// 		},
	// 		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
	// 			dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
	// 		},
	// 	},
	//
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
	//
	// 		if result.Landed() {
	// 			spell.Dot(target).Apply(sim)
	// 		}
	//
	// 		spell.DealOutcome(sim, result)
	// 	},
	// })
}
