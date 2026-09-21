package paladin

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

func (paladin *Paladin) getConsecrationTimer() *core.Timer {
	if paladin.consecrationTimer == nil {
		paladin.consecrationTimer = paladin.NewTimer()
	}
	return paladin.consecrationTimer
}

var ConsecrationRankMap = spellData.Consecration

// Consecration
// https://www.wowhead.com/forever/spell=26573
//
// Consecrates the land beneath the Paladin, doing X Holy damage over 8 sec to enemies who enter the
// area. The first N enemies who enter the area will take an additional Y damage over 8 sec.
// TODO: To be implemented.
func (paladin *Paladin) registerConsecration(rankConfig shared.SpellData) {
	panic("To be implemented")

	// The ported implementation, kept until this class is done:
	// tick := rankConfig.Periodic.(shared.SpellDataPeriodic)
	//
	// // The extra damage the first few targets take, which is the only part of the spell the client
	// // gives a spell power coefficient: the tick everyone takes has none.
	// bonus := rankConfig.SecondaryPeriodic.(shared.SpellDataPeriodic)
	// bonusTargets := int(rankConfig.Effect(shared.A_PERIODIC_DUMMY, 0).Value)
	//
	// spellID := rankConfig.SpellID
	// cost := rankConfig.Cost
	//
	// // The bonus scales on its own coefficient, so it is added to the base damage here rather than
	// // through the dot's, which is the base tick's.
	// dealTick := func(sim *core.Simulation, dot *core.Dot) {
	// 	for i, target := range sim.Encounter.ActiveTargetUnits {
	// 		damage := tick.Tick
	// 		if i < bonusTargets {
	// 			damage += bonus.Tick + bonus.Coef*dot.Spell.BonusDamage(dot.Spell.Unit.AttackTables[target.UnitIndex])
	// 		}
	// 		dot.Spell.CalcAndDealPeriodicDamage(sim, target, damage, dot.OutcomeTickMagicHit)
	// 	}
	// }
	//
	// paladin.RegisterSpell(core.SpellConfig{
	// 	ActionID:       core.ActionID{SpellID: spellID},
	// 	SpellSchool:    core.SpellSchoolHoly,
	// 	DefenseType:    core.DefenseTypeMagic,
	// 	ProcMask:       core.ProcMaskSpellDamage,
	// 	Flags:          core.SpellFlagAPL,
	// 	ClassSpellMask: SpellMaskConsecration,
	// 	Rank:           rankConfig.Rank,
	//
	// 	DamageMultiplier: 1,
	// 	ThreatMultiplier: 1,
	//
	// 	MaxRange: 8,
	//
	// 	ManaCost: core.ManaCostOptions{
	// 		FlatCost: cost,
	// 	},
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			GCD: rankConfig.GCD,
	// 		},
	// 		CD: core.Cooldown{
	// 			Timer:    paladin.getConsecrationTimer(),
	// 			Duration: rankConfig.Cooldown,
	// 		},
	// 	},
	//
	// 	Dot: core.DotConfig{
	// 		IsAOE: true,
	// 		Aura: core.Aura{
	// 			ActionID: core.ActionID{SpellID: spellID},
	// 			Label:    "Consecration" + paladin.Label + " " + rankConfig.GetRankLabel(),
	// 		},
	// 		NumberOfTicks:    tick.NumberOfTicks - 1, // the sim adds an immediate tick below
	// 		TickLength:       tick.TickLength,
	// 		BonusCoefficient: tick.Coef,
	// 		OnTick: func(sim *core.Simulation, _ *core.Unit, dot *core.Dot) {
	// 			dealTick(sim, dot)
	// 		},
	// 	},
	//
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		// Consecration does one hit check on cast but the ground effect will still be applied
	// 		// meaning it's only needed to proc things like Eye of Magtheridon (procs on resist)
	// 		spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)
	//
	// 		dot := spell.AOEDot()
	// 		dot.Apply(sim)
	// 		dealTick(sim, dot)
	// 	},
	// })
}
