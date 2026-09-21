package warlock

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

// Wrack is the Forever Affliction capstone (1316697): a six second shadow channel that also makes
// the warlock's other shadow dots on the target tick 10% harder while it runs - the second effect
// on the client's row.
func (warlock *Warlock) registerWrack() {
	if !warlock.Talents.Wrack {
		return
	}

	rank := spellData.Wrack.HighestRank()
	tick := rank.Periodic.(shared.SpellDataPeriodic)
	dotBonus := 1 + spellData.Wrack.EffectAt(1).FractionAt(rank.Rank)

	warlock.Wrack = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.SpellID},
		SpellSchool:    rank.SpellSchool,
		DefenseType:    rank.DefenseType,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagChanneled | core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellWrack,

		ManaCost: core.ManaCostOptions{FlatCost: rank.Cost},
		Cast:     core.CastConfig{DefaultCast: core.Cast{GCD: rank.GCD}},

		DamageMultiplierAdditive: 1,
		DamageMultiplier:         1,
		ThreatMultiplier:         1,
		BonusCoefficient:         tick.Coef,

		Dot: core.DotConfig{
			Aura:                 core.Aura{Label: "Wrack"},
			NumberOfTicks:        tick.NumberOfTicks,
			TickLength:           tick.TickLength,
			AffectedByCastSpeed:  true,
			HasteReducesDuration: true,
			BonusCoefficient:     tick.Coef,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.Snapshot(target, tick.Tick)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				result := dot.CalcSnapshotDamage(sim, target, shared.PeriodicTickOutcome(rank, dot))
				result.Damage *= warlock.soulSiphonMultiplier(target)
				dot.Spell.DealPeriodicDamage(sim, result)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})

	for _, target := range warlock.Env.Encounter.AllTargetUnits {
		target.AddDynamicDamageTakenModifier(func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult, isPeriodic bool) {
			if spell.Unit != &warlock.Unit || spell == warlock.Wrack {
				return
			}

			if isPeriodic && spell.SpellSchool.Matches(core.SpellSchoolShadow) && warlock.Wrack.Dot(result.Target).IsActive() {
				result.Damage *= dotBonus
			}
		})
	}
}
