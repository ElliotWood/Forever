package warrior

import (
	"time"

	"github.com/wowsims/forever/sim/core"
)

// TODO: Manual review needed -- spell 402927 states ${1+$AP*$m3/100} damage with $m3 of 15,
// heals 10% of maximum health and has a 30 second cooldown; spell 402975 states 20 seconds.
const (
	victoryRushBaseDamage        = 1.0
	victoryRushAPCoef            = 0.15
	victoryRushHealPercent       = 0.10
	victoryRushCooldown          = time.Second * 30
	victoriousDuration           = time.Second * 20
	victoriousSpellID      int32 = 402975
)

func (warrior *Warrior) registerVictoryRush() {
	actionID := core.ActionID{SpellID: 402927}
	healthMetrics := warrior.NewHealthMetrics(actionID)

	// TODO: spell 402974 grants this on a kill, which the sim never simulates, so nothing
	// activates it and Victory Rush stays uncastable.
	victoriousAura := warrior.RegisterAura(core.Aura{
		Label:    "Victorious",
		ActionID: core.ActionID{SpellID: victoriousSpellID},
		Duration: victoriousDuration,
	})

	warrior.VictoryRush = warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskVictoryRush,
		MaxRange:       core.MaxMeleeRange,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: victoryRushCooldown,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return victoriousAura.IsActive()
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := victoryRushBaseDamage + victoryRushAPCoef*spell.MeleeAttackPower(target)
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			warrior.GainHealth(sim, warrior.MaxHealth()*victoryRushHealPercent, healthMetrics)
			victoriousAura.Deactivate(sim)
		},

		RelatedSelfBuff: victoriousAura,
	})
}
