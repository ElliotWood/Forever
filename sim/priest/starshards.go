package priest

import (
	"fmt"

	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

// Starshards is the Night Elf racial: an Arcane channel, free, on a 30 second cooldown.
var StarshardsRankMap = spellData.Starshards

func (priest *Priest) registerStarshardsSpell(rank shared.SpellData, cdTimer *core.Timer) {
	tick := rank.Periodic.(shared.SpellDataPeriodic)

	priest.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.SpellID},
		SpellSchool:    rank.SpellSchool,
		DefenseType:    rank.DefenseType,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL | core.SpellFlagChanneled,
		ClassSpellMask: PriestSpellStarshards,
		Rank:           rank.Rank,
		MaxRange:       rank.MaxRange,

		ManaCost: core.ManaCostOptions{
			FlatCost: rank.Cost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: rank.GCD,
			},
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: rank.Cooldown,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: fmt.Sprintf("Starshards-%d", rank.Rank),
			},
			NumberOfTicks:       tick.NumberOfTicks,
			TickLength:          tick.TickLength,
			AffectedByCastSpeed: false,
			BonusCoefficient:    tick.Coef,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.Snapshot(target, tick.Tick)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, priestTickOutcome(rank, dot))
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},

		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			if useSnapshot {
				return spell.Dot(target).CalcSnapshotDamage(sim, target, spell.OutcomeExpectedMagicHit)
			}
			return spell.CalcPeriodicDamage(sim, target, tick.Tick, spell.OutcomeExpectedMagicHit)
		},
	})
}
