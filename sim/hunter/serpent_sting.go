package hunter

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

func (hunter *Hunter) registerSerpentStingSpell() {
	rank := spellData.SerpentSting.HighestRank()
	tick := rank.Periodic.(shared.SpellDataPeriodic)

	// The beta client carries no spell power coefficient on Serpent Sting at all, so Classic's
	// stands: the full-duration 1.0 split across the ticks.
	spellCoeff := 1.0 / float64(tick.NumberOfTicks)

	hunter.SerpentSting = hunter.RegisterRangedSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.SpellID},
		SpellSchool:    rank.SpellSchool,
		DefenseType:    rank.DefenseType,
		ClassSpellMask: HunterSpellSerpentSting,
		ProcMask:       core.ProcMaskRangedSpecial,
		Flags:          core.SpellFlagAPL | core.SpellFlagPoison,
		MissileSpeed:   rank.MissileSpeed,

		ManaCost: core.ManaCostOptions{
			FlatCost: rank.Cost,
		},

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Serpent Sting",
				Tag:   "Sting",
			},
			NumberOfTicks:    tick.NumberOfTicks,
			TickLength:       tick.TickLength,
			BonusCoefficient: spellCoeff,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.Snapshot(target, tick.Damage(sim))
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, shared.PeriodicTickOutcome(rank, dot))
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeRangedHitNoHitCounter)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealOutcome(sim, result)

				// Serpent Sting is our only sting, so there is no other sting of ours to knock off; another
				// hunter's stays up beside it (it used to be removed, and two hunters kept undoing each other).
				if result.Landed() {
					spell.Dot(target).Apply(sim)
				}
			})
		},
	})
}
