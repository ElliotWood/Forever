package warrior

import (
	"github.com/wowsims/forever/sim/core"
)

var whirlwindRank = spellData.Whirlwind.HighestRank()

func (warrior *Warrior) registerWhirlwind() {
	actionID := core.ActionID{SpellID: whirlwindRank.SpellID}

	// Raging Blows (1310315) adds the off-hand strike.
	var whirlwindOH *core.Spell
	if warrior.Talents.RagingBlows {
		whirlwindOH = warrior.RegisterSpell(core.SpellConfig{
			ActionID:       actionID.WithTag(2),
			SpellSchool:    core.SpellSchoolPhysical,
			DefenseType:    core.DefenseTypeMelee,
			ProcMask:       core.ProcMaskMeleeOHSpecial,
			ClassSpellMask: SpellMaskWhirlwindOh,
			Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagPassiveSpell | core.SpellFlagNoOnCastComplete,

			DamageMultiplier: 1,
			// TODO: Manual review needed -- the client states no threat coefficient; 1 until measured in game.
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				baseDamage := warrior.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower(target))
				spell.CalcCleaveDamage(sim, target, whirlwindRank.MaxTargets, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
				spell.DealBatchedAoeDamage(sim)
			},
		})
	}

	warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID.WithTag(1),
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		ClassSpellMask: SpellMaskWhirlwind,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		RageCost: core.RageCostOptions{
			Cost: whirlwindRank.Cost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: whirlwindRank.GCD,
			},
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: whirlwindRank.Cooldown,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		// TODO: Manual review needed -- the client states no threat coefficient; 1 until measured in game.
		ThreatMultiplier: 1,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.StanceMatches(BerserkerStance)
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := warrior.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower(target))
			results := spell.CalcCleaveDamage(sim, target, whirlwindRank.MaxTargets, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			warrior.CastNormalizedSweepingStrikesAttack(results, sim)
			spell.DealBatchedAoeDamage(sim)

			if whirlwindOH != nil && warrior.HasOHWeapon() {
				whirlwindOH.Cast(sim, target)
			}
		},
	})
}
