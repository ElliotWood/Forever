package warrior

import (
	"time"

	"github.com/wowsims/forever/sim/core"
)

// TODO: Manual review needed -- spell 1680 states 4 targets, a 25 rage cost and a 10 second cooldown.
const whirlwindMaxTargets int32 = 4

func (warrior *Warrior) registerWhirlwind() {
	actionID := core.ActionID{SpellID: 1680}

	whirlwindOH := warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID.WithTag(2),
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeOHSpecial,
		ClassSpellMask: SpellMaskWhirlwindOh,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagPassiveSpell | core.SpellFlagNoOnCastComplete,

		DamageMultiplier: 1,
		ThreatMultiplier: 1.25,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := warrior.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower(target))
			spell.CalcCleaveDamage(sim, target, whirlwindMaxTargets, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			spell.DealBatchedAoeDamage(sim)
		},
	})

	warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID.WithTag(1),
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		ClassSpellMask: SpellMaskWhirlwind,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		RageCost: core.RageCostOptions{
			Cost: 25,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Second * 10,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1.25,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.StanceMatches(BerserkerStance)
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := warrior.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower(target))
			results := spell.CalcCleaveDamage(sim, target, whirlwindMaxTargets, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			warrior.CastNormalizedSweepingStrikesAttack(results, sim)
			spell.DealBatchedAoeDamage(sim)

			// Raging Blows (1310315) adds the off-hand strike.
			if warrior.HasOHWeapon() && warrior.Talents.RagingBlows {
				whirlwindOH.Cast(sim, target)
			}
		},
	})
}
