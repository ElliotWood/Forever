package paladin

import (
	"time"

	"github.com/wowsims/classic/sim/core"
)

// Templar's Bulwark is new in Forever and borrows Sacred Shield's spell id, the paladin absorb
// the tooltip describes.
// TODO: assumed baseline, beta will confirm - the tooltip carries no cooldown, so it shares the
// 5 minutes of the two Forbearance abilities Sacred Duty shortens alongside it.
func (paladin *Paladin) registerTemplarsBulwark() {
	if !paladin.Talents.TemplarsBulwark {
		return
	}

	actionID := core.ActionID{SpellID: 53601}

	// The sim has no absorb model. A shield worth the paladin's whole health pool is far more
	// than a tank takes in 8 seconds, so it is modelled as damage taken dropping to nothing.
	const damageTaken = 0.01

	bulwarkAura := paladin.RegisterAura(core.Aura{
		Label:    "Templar's Bulwark",
		ActionID: actionID,
		Duration: time.Second * 8,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			paladin.PseudoStats.DamageTakenMultiplier *= damageTaken
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			paladin.PseudoStats.DamageTakenMultiplier /= damageTaken
		},
	})

	// TODO: Both ranks of Sacred Duty read 30 sec, the second is assumed to scale linearly.
	cooldown := time.Minute*5 - time.Second*30*time.Duration(paladin.Talents.SacredDuty)

	bulwark := paladin.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL | SpellFlag_Forbearance,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: cooldown,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			bulwarkAura.Activate(sim)
		},
	})

	paladin.AddMajorCooldown(core.MajorCooldown{
		Spell: bulwark,
		Type:  core.CooldownTypeSurvival,
	})
}
