package warlock

// TODO: To be implemented. Drain Life exists in Forever with a full six-rank ladder
// (spellData.DrainLife, 689/699/709/7651/11699/11700). The implementation below is the
// TBC one already ported onto that ladder -- HighestRank() in place of the level-70 rank
// it used to pin, and Soul Siphon reading spellData.SoulSiphon -- so bringing it back is
// uncommenting it, not rewriting it. It stays commented until it has been reviewed.
//
// Imports and rank pins the implementation needs, kept with it:
// import (
// 	"math"
//
// 	"github.com/wowsims/forever/sim/common/shared"
// 	"github.com/wowsims/forever/sim/core"
// )
//
// var drainLifeRank = spellData.DrainLife.HighestRank()
// var drainLifeTick = drainLifeRank.Periodic.(shared.SpellDataPeriodic)
// var drainLifeCoeff = drainLifeTick.Coef

func (warlock *Warlock) registerDrainLife() {
	panic("To be implemented")

	// healthMetric := warlock.NewHealthMetrics(core.ActionID{SpellID: drainLifeRank.SpellID})
	// resultSlice := make(core.SpellResultSlice, 1)
	//
	// // TODO: cappedDmgBonus/the == 2 check are TBC's 2-rank Soul Siphon (2%/4% per Affliction effect,
	// // capped at 1.24/1.60 total). spellData.SoulSiphon now ships 3 ranks at 4/8/12% with no cap effect
	// // in the table, so rank 3 currently falls through to the rank-2 cap of 1.60 rather than a sourced
	// // rank-3 number - neither the client nor spellData states what the rank-3 cap should be.
	// cappedDmgBonus := 1.24
	// if warlock.Talents.SoulSiphon == 2 {
	// 	cappedDmgBonus = 1.60
	// }
	//
	// // Read once rather than per tick: the ladder lookup scans the table and boxes the rank.
	// soulSiphonPerAffliction := spellData.SoulSiphon.Effect(shared.A_DUMMY, 0).FractionAt(warlock.Talents.SoulSiphon)
	//
	// warlock.DrainLife = warlock.RegisterSpell(core.SpellConfig{
	// 	ActionID:       core.ActionID{SpellID: drainLifeRank.SpellID},
	// 	SpellSchool:    drainLifeRank.SpellSchool,
	// 	DefenseType:    drainLifeRank.DefenseType,
	// 	ProcMask:       core.ProcMaskSpellDamage,
	// 	Flags:          core.SpellFlagChanneled | core.SpellFlagAPL,
	// 	ClassSpellMask: WarlockSpellDrainLife,
	//
	// 	ManaCost: core.ManaCostOptions{FlatCost: drainLifeRank.Cost},
	// 	Cast:     core.CastConfig{DefaultCast: core.Cast{GCD: drainLifeRank.GCD}},
	//
	// 	DamageMultiplierAdditive: 1,
	// 	ThreatMultiplier:         1,
	// 	BonusCoefficient:         drainLifeCoeff,
	//
	// 	Dot: core.DotConfig{
	// 		Aura:                 core.Aura{Label: "Drain Life"},
	// 		NumberOfTicks:        drainLifeTick.NumberOfTicks,
	// 		TickLength:           drainLifeTick.TickLength,
	// 		AffectedByCastSpeed:  true,
	// 		HasteReducesDuration: true,
	// 		BonusCoefficient:     drainLifeCoeff,
	// 		OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
	// 			dot.Snapshot(target, drainLifeTick.Tick)
	// 		},
	// 		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
	// 			dot.PeriodicDamageMultiplier = math.Max(1, math.Min(1+(soulSiphonPerAffliction*warlock.AfflictionCount(target)), cappedDmgBonus))
	// 			resultSlice[0] = dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
	// 			warlock.GainHealth(sim, resultSlice[0].Damage*warlock.PseudoStats.SelfHealingMultiplier, healthMetric)
	// 		},
	// 	},
	//
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)
	// 		if result.Landed() {
	// 			spell.Dot(target).Apply(sim)
	// 			spell.DealOutcome(sim, result)
	// 		}
	// 	},
	// })
}
