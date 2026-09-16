package warrior

import (
	"time"

	"github.com/wowsims/classic/sim/core"
)

// TODO: Classic Update
func (warrior *Warrior) RegisterShieldWallCD() {
	// Forever trades strength for uptime: 60% for 12 sec against Classic's 75% for 10.
	// Read off Esfand's Warrior in "Warriors Just Got a MASSIVE Rework" at 11:49.
	forever := warrior.Env.IsForever()
	duration := core.TernaryDuration(forever, time.Second*12, time.Second*10)
	//This is the inverse of the tooltip since it is a damage TAKEN coefficient
	damageTaken := core.TernaryFloat64(forever, 0.40, 0.25)

	actionID := core.ActionID{SpellID: 871}
	swAura := warrior.RegisterAura(core.Aura{
		Label:    "Shield Wall",
		ActionID: actionID,
		Duration: duration,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.DamageTakenMultiplier *= damageTaken
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.DamageTakenMultiplier /= damageTaken
		},
	})

	// Improved Shield Wall reduces the cooldown instead of extending the duration in Forever.
	// TODO: only rank 1 was shown, beta will confirm that rank 2 is another 5.5 minutes.
	// Forever halves the untalented cooldown to 15 min, read off Esfand's Warrior in
	// "Warriors Just Got a MASSIVE Rework" at 11:49.
	baseCooldown := core.TernaryDuration(forever, time.Minute*15, time.Minute*30)
	cooldownDur := baseCooldown - time.Second*330*time.Duration(warrior.Talents.ImprovedShieldWall)

	swSpell := warrior.RegisterSpell(DefensiveStance, core.SpellConfig{
		ActionID: actionID,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: 0,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: cooldownDur,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.PseudoStats.CanBlock
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			swAura.Activate(sim)
		},
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: swSpell.Spell,
		Type:  core.CooldownTypeSurvival,
	})
}
