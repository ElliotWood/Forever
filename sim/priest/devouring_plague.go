package priest

import (
	"fmt"

	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

// Devouring Plague is an Undead racial in Classic. The Forever beta client teaches it to priests of
// every race (SkillLineAbility race mask -1), so priest.go registers it for all of them. The ticks
// heal the priest for what they deal, and the client leaves Periodic Can Crit off, so they never crit.
var DevouringPlagueRankMap = spellData.DevouringPlague

func (priest *Priest) registerDevouringPlagueSpell(rank shared.SpellData, cdTimer *core.Timer) {
	tick := rank.Periodic.(shared.SpellDataPeriodic)
	healthMetrics := priest.NewHealthMetrics(core.ActionID{SpellID: rank.SpellID}.WithTag(1))

	priest.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.SpellID},
		SpellSchool:    rank.SpellSchool,
		DefenseType:    rank.DefenseType,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: PriestSpellDevouringPlague,
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
				Label: fmt.Sprintf("DevouringPlague-%d", rank.Rank),
			},
			NumberOfTicks:       tick.NumberOfTicks,
			TickLength:          tick.TickLength,
			AffectedByCastSpeed: false,
			BonusCoefficient:    tick.Coef,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.Snapshot(target, tick.Tick)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				result := dot.CalcAndDealPeriodicSnapshotDamage(sim, target, priestTickOutcome(rank, dot))
				priest.GainHealth(sim, result.Damage, healthMetrics)
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
