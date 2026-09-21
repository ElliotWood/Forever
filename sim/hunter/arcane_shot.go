package hunter

import (
	"github.com/wowsims/forever/sim/core"
)

// The beta client carries no spell power coefficient on Arcane Shot at all, so Classic's stand.
var arcaneShotCoefficients = [9]float64{0, .204, .3, .429, .429, .429, .429, .429, .429}

func (hunter *Hunter) registerArcaneShotSpell(timer *core.Timer) {
	rank := spellData.ArcaneShot.HighestRank()
	baseDamage := rank.Direct.Damage

	hunter.ArcaneShot = hunter.RegisterRangedSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.SpellID},
		SpellSchool:    rank.SpellSchool,
		DefenseType:    rank.DefenseType,
		ClassSpellMask: HunterSpellArcaneShot,
		ProcMask:       core.ProcMaskRangedSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		MissileSpeed:   rank.MissileSpeed,

		ManaCost: core.ManaCostOptions{
			FlatCost: rank.Cost,
		},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    timer,
				Duration: rank.Cooldown,
			},
		},

		BonusCoefficient: arcaneShotCoefficients[rank.Rank],

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcDamage(sim, target, baseDamage(sim), spell.OutcomeRangedHitAndCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
