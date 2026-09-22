package hunter

import (
	"github.com/wowsims/forever/sim/core"
)

// The beta client has one rank of Multi-Shot: no flat bonus, and a 6 sec cooldown shared with Aimed
// Shot where Classic had 10 sec alone. Ranks 2-5 are gone from the spellbook.
func (hunter *Hunter) registerMultiShotSpell(timer *core.Timer) {
	rank := spellData.MultiShot.HighestRank()
	numHits := min(3, hunter.Env.ActiveTargetCount())

	hunter.MultiShot = hunter.RegisterRangedSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.SpellID},
		SpellSchool:    rank.SpellSchool,
		DefenseType:    rank.DefenseType,
		ClassSpellMask: HunterSpellMultiShot,
		ProcMask:       core.ProcMaskRangedSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		MissileSpeed:   rank.MissileSpeed,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 13.9, // client SpellPower 170887: PowerCostPct 13.9, no flat cost
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				CastTime: rank.CastTime,
			},
			CD: core.Cooldown{
				Timer:    timer,
				Duration: rank.Cooldown,
			},
		},

		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			results := make([]*core.SpellResult, numHits)
			curTarget := target

			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				baseDamage := hunter.AutoAttacks.Ranged().CalculateNormalizedWeaponDamage(sim, spell.RangedAttackPower(target)) +
					hunter.AmmoDamageBonus

				results[hitIndex] = spell.CalcDamage(sim, curTarget, baseDamage, spell.OutcomeRangedHitAndCrit)
				curTarget = sim.Environment.NextActiveTargetUnit(curTarget)
			}

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
					spell.DealDamage(sim, results[hitIndex])
				}
			})
		},
	})
}
