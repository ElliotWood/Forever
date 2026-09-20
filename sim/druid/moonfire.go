package druid

import (
	"github.com/wowsims/forever/sim/common/shared"
)

var moonfireRank = spellData.Moonfire.HighestRank()
var moonfireTick = moonfireRank.Periodic.(shared.SpellDataPeriodic)

func (druid *Druid) registerMoonfireSpell() {
	druid.registerMoonfireImpactSpell()
	druid.registerMoonfireDoTSpell()
}

// TODO: To be implemented.
func (druid *Druid) registerMoonfireDoTSpell() {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// druid.Moonfire.RelatedDotSpell = druid.Unit.RegisterSpell(core.SpellConfig{
	// 	ActionID:       core.ActionID{SpellID: moonfireRank.SpellID}.WithTag(1),
	// 	SpellSchool:    moonfireRank.SpellSchool,
	// 	DefenseType:    moonfireRank.DefenseType,
	// 	ProcMask:       core.ProcMaskSpellDamage,
	// 	ClassSpellMask: DruidSpellMoonfireDoT,
	// 	Flags:          core.SpellFlagPassiveSpell,
	//
	// 	DamageMultiplier: 1,
	// 	ThreatMultiplier: 1,
	//
	// 	Dot: core.DotConfig{
	// 		Aura: core.Aura{
	// 			Label: "Moonfire",
	// 		},
	// 		NumberOfTicks:       moonfireTick.NumberOfTicks,
	// 		TickLength:          moonfireTick.TickLength,
	// 		AffectedByCastSpeed: false,
	// 		BonusCoefficient:    moonfireTick.Coef,
	//
	// 		OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
	// 			dot.Snapshot(target, moonfireTick.Tick)
	// 		},
	// 		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
	// 			dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
	// 		},
	// 	},
	//
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		result := spell.CalcOutcome(sim, target, spell.OutcomeAlwaysHitNoHitCounter)
	//
	// 		spell.Dot(target).Apply(sim)
	// 		spell.DealOutcome(sim, result)
	// 	},
	// })
}

// TODO: To be implemented.
func (druid *Druid) registerMoonfireImpactSpell() {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// druid.Moonfire = druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
	// 	ActionID:       core.ActionID{SpellID: moonfireRank.SpellID},
	// 	SpellSchool:    moonfireRank.SpellSchool,
	// 	DefenseType:    moonfireRank.DefenseType,
	// 	ProcMask:       core.ProcMaskSpellDamage,
	// 	ClassSpellMask: DruidSpellMoonfire,
	// 	Flags:          core.SpellFlagAPL,
	//
	// 	ManaCost: core.ManaCostOptions{
	// 		FlatCost: moonfireRank.Cost,
	// 	},
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			GCD: moonfireRank.GCD,
	// 		},
	// 	},
	//
	// 	BonusCoefficient: moonfireRank.Direct.BonusCoefficient(),
	// 	DamageMultiplier: 1,
	// 	ThreatMultiplier: 1,
	// 	MaxRange:         moonfireRank.MaxRange,
	//
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		baseDamage := moonfireRank.Direct.Damage(sim)
	// 		result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
	//
	// 		if result.Landed() {
	// 			druid.Moonfire.RelatedDotSpell.Cast(sim, target)
	// 		}
	//
	// 		spell.DealDamage(sim, result)
	// 	},
	// })
}
