package warrior

import (
	"time"

	"github.com/wowsims/forever/sim/core"
)

func (warrior *Warrior) registerBerserkerRage() {
	actionID := core.ActionID{SpellID: 18499}
	rageMetrics := warrior.NewRageMetrics(actionID)
	// Improved Berserker Rage (20500) adds rage on the cast; both of its effects are dummies, so
	// the rage one is named by its index. Its second effect, shedding movement impairment, has
	// nothing to act on in the sim.
	rageGain := spellData.ImprovedBerserkerRage.EffectAt(0).TenthsAt(warrior.Talents.ImprovedBerserkerRage)

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
			if rageGain > 0 {
				warrior.AddRage(sim, rageGain, rageMetrics)
			}
			aura.Activate(sim)
		},
		RelatedSelfBuff: aura,
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeSurvival,
		ShouldActivate: func(s *core.Simulation, c *core.Character) bool {
			return rageGain > 0 && warrior.CurrentRage()+rageGain <= warrior.MaximumRage()
		},
	})
}
