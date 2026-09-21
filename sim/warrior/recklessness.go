package warrior

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/stats"
)

var recklessnessRank = spellData.Recklessness.HighestRank()

func (warrior *Warrior) registerRecklessness() {
	actionID := core.ActionID{SpellID: recklessnessRank.SpellID}

	aura := warrior.RegisterAura(core.Aura{
		Label:    "Recklessness",
		ActionID: actionID,
		Duration: recklessnessRank.Duration,
	}).AttachStatsBuff(
		stats.Stats{
			stats.PhysicalCritPercent: recklessnessRank.Effect(shared.A_MOD_CRIT_PCT, 0).Value,
			stats.SpellCritPercent:    recklessnessRank.Effect(shared.A_MOD_CRIT_PCT, 0).Value,
		},
	).AttachMultiplicativePseudoStatBuff(
		&warrior.PseudoStats.DamageTakenMultiplier,
		recklessnessRank.Effect(shared.A_MOD_DAMAGE_PERCENT_TAKEN, 127).Multiplier(),
	).
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
