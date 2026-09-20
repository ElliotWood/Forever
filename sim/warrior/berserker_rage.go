package warrior

import (
	"time"

	"github.com/wowsims/forever/sim/core"
)

func (warrior *Warrior) registerBerserkerRage() {
	actionID := core.ActionID{SpellID: 18499}
	rageMetrics := warrior.NewRageMetrics(actionID)

	aura := warrior.RegisterAura(core.Aura{
		Label:    "Berserker Rage",
		ActionID: actionID,
		// TODO: Manual review needed -- spell 18499 states a 10 second duration.
		Duration: time.Second * 10,
	}).
		// Grants immunity to Fear, Sap and Incapacitate effects.
		AttachFearImmunity()

	spell := warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: SpellMaskBerserkerRage,
		Flags:          core.SpellFlagAPL | core.SpellFlagCastWhileIncapacitated,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer: warrior.NewTimer(),
				// TODO: Manual review needed -- spell 18499 states a 30 second cooldown.
				Duration: time.Second * 30,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.StanceMatches(BerserkerStance)
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			if warrior.BerserkerRageRageGain > 0 {
				warrior.AddRage(sim, warrior.BerserkerRageRageGain, rageMetrics)
			}
			aura.Activate(sim)
		},
		RelatedSelfBuff: aura,
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeSurvival,
		ShouldActivate: func(s *core.Simulation, c *core.Character) bool {
			return warrior.BerserkerRageRageGain > 0 && warrior.CurrentRage()+warrior.BerserkerRageRageGain <= warrior.MaximumRage()
		},
	})
}
