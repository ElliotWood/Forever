package warrior

import (
	"github.com/wowsims/forever/sim/core"
)

// This spell works but the Sim never kills a target
func (warrior *Warrior) registerVictoryRush() {
	victoryRushRank := spellData.VictoryRush.HighestRank()
	victoryRushAPCoef := victoryRushRank.Effects[2].Fraction()
	victoryRushHealPercent := victoryRushRank.Effects[1].Fraction()
	victoriousRank := spellData.VictoryRushTriggered.HighestRank()

	actionID := core.ActionID{SpellID: victoryRushRank.SpellID}
	healthMetrics := warrior.NewHealthMetrics(actionID)

	victoriousAura := warrior.RegisterAura(core.Aura{
		Label:    "Victorious",
		ActionID: core.ActionID{SpellID: victoriousRank.SpellID},
		Duration: victoriousRank.Duration,
	})

	warrior.VictoryRush = warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskVictoryRush,
		MaxRange:       victoryRushRank.MaxRange,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: victoryRushRank.GCD,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: victoryRushRank.Cooldown,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return victoriousAura.IsActive()
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := victoryRushRank.Direct.Damage(sim) + victoryRushAPCoef*spell.MeleeAttackPower(target)
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			warrior.GainHealth(sim, warrior.MaxHealth()*victoryRushHealPercent, healthMetrics)
			victoriousAura.Deactivate(sim)
		},

		RelatedSelfBuff: victoriousAura,
	})
}
