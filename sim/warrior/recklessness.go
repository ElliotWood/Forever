package warrior

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

var recklessnessRank = spellData.Recklessness.HighestRank()

func (warrior *Warrior) registerRecklessness() {
	actionID := core.ActionID{SpellID: recklessnessRank.SpellID}

	aura := warrior.RegisterAura(core.Aura{
		Label:    "Recklessness",
		ActionID: actionID,
		Duration: recklessnessRank.Duration,
	}).AttachSpellMod(core.SpellModConfig{
		ProcMask:   core.ProcMaskMeleeSpecial,
		Kind:       core.SpellMod_BonusCrit_Percent,
		FloatValue: recklessnessRank.Effect(shared.A_MOD_CRIT_PCT, 0).Value,
	}).AttachMultiplicativePseudoStatBuff(
		&warrior.PseudoStats.DamageTakenMultiplier,
		1+recklessnessRank.Effect(shared.A_MOD_DAMAGE_PERCENT_TAKEN, 127).Value/100,
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
				GCD: recklessnessRank.GCD,
			},
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: recklessnessRank.Cooldown,
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
