package hunter

import (
	"github.com/wowsims/forever/sim/core"
)

// Aimed Shot is no longer a talent. The beta client (1.60.1.69893) keeps every rank on the hunter
// with a 2 sec cast, down from 3, and a much smaller flat bonus (rank 6 600 -> 166); mana costs and
// the 6 sec cooldown are Classic's. Forever puts that cooldown on the Multi-Shot timer.
func (hunter *Hunter) registerAimedShotSpell(timer *core.Timer) {
	rank := spellData.AimedShot.HighestRank()
	flatBonus := rank.Direct.Damage

	hunter.AimedShot = hunter.RegisterRangedSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.SpellID},
		SpellSchool:    rank.SpellSchool,
		DefenseType:    rank.DefenseType,
		ClassSpellMask: HunterSpellAimedShot,
		ProcMask:       core.ProcMaskRangedSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		MissileSpeed:   rank.MissileSpeed,

		ManaCost: core.ManaCostOptions{
			FlatCost: rank.Cost,
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
			baseDamage := hunter.AutoAttacks.Ranged().CalculateNormalizedWeaponDamage(sim, spell.RangedAttackPower(target)) +
				flatBonus(sim)

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
