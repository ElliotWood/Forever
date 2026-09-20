package warrior

import (
	"time"

	"github.com/wowsims/forever/sim/core"
)

func (warrior *Warrior) registerRecklessness() {
	actionID := core.ActionID{SpellID: 1719}

	aura := warrior.RegisterAura(core.Aura{
		Label:    "Recklessness",
		ActionID: actionID,
		// TODO: Manual review needed -- spell 1719 states a 15 second duration.
		Duration: time.Second * 15,
	}).AttachSpellMod(core.SpellModConfig{
		ProcMask: core.ProcMaskMeleeSpecial,
		Kind:     core.SpellMod_BonusCrit_Percent,
		// TODO: Manual review needed -- spell 1719 states 100% critical strike chance.
		FloatValue: 100,
	}).AttachMultiplicativePseudoStatBuff(
		// TODO: Manual review needed -- spell 1719 states 20% increased damage taken.
		&warrior.PseudoStats.DamageTakenMultiplier, 1.2,
	).
		// Grants immunity to Fear effects.
		AttachFearImmunity()

	spell := warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		DefenseType:    core.DefenseTypeMelee,
		Flags:          core.SpellFlagAPL | core.SpellFlagCastWhileIncapacitated,
		ClassSpellMask: SpellMaskRecklessness,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer: warrior.NewTimer(),
				// TODO: Manual review needed -- spell 1719 states a 30 minute cooldown.
				Duration: time.Minute * 30,
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.StanceMatches(BerserkerStance)
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			aura.Activate(sim)
		},

		RelatedSelfBuff: aura,
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
	})
}
