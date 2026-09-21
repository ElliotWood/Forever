package warrior

import (
	"github.com/wowsims/forever/sim/core"
)

var berserkerRageRank = spellData.BerserkerRage.HighestRank()

func (warrior *Warrior) registerBerserkerRage() {
	actionID := core.ActionID{SpellID: berserkerRageRank.SpellID}
	rageMetrics := warrior.NewRageMetrics(actionID)
	rageGain := spellData.ImprovedBerserkerRage.EffectAt(0).TenthsAt(warrior.Talents.ImprovedBerserkerRage)

	// TODO: Ingame test needed -- the client states no amount for the extra rage a hit taken
	// generates; doubled here.
	aura := warrior.RegisterAura(core.Aura{
		Label:    "Berserker Rage",
		ActionID: actionID,
		Duration: berserkerRageRank.Duration,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.MultiplyDamageTakenRageGen(2)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.MultiplyDamageTakenRageGen(0.5)
		},
	}).
		AttachFearImmunity()

	spell := warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: SpellMaskBerserkerRage,
		Flags:          core.SpellFlagAPL | core.SpellFlagCastWhileIncapacitated,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: berserkerRageRank.GCD,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: berserkerRageRank.Cooldown,
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
