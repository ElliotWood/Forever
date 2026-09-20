package warrior

import (
	"time"

	"github.com/wowsims/forever/sim/core"
)

var victoryRushRank = spellData.VictoryRush.HighestRank()

// Spell 402927 states ${1+$AP*$m3/100}: the dummy at effect index 2 is the attack power
// coefficient as a percentage, and the heal at index 1 is a percentage of maximum health.
var victoryRushAPCoef = victoryRushRank.Effects[2].Value / 100
var victoryRushHealPercent = victoryRushRank.Effects[1].Value / 100

// TODO: Manual review needed -- the window Victory Rush has to be used in is spell 402975's 20
// seconds; it is not a ranked row, so it carries no table.
const (
	victoriousDuration       = time.Second * 20
	victoriousSpellID  int32 = 402975
)

func (warrior *Warrior) registerVictoryRush() {
	actionID := core.ActionID{SpellID: victoryRushRank.SpellID}
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
