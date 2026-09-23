package druid

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/stats"
)

var barkskinRank = spellData.Barkskin.HighestRank()

// Barkskin: 20% less Physical damage taken for 15 sec, no cost. The client states the reduction on
// A_MOD_DAMAGE_PERCENT_TAKEN with school mask 1 (Physical).
func (druid *Druid) registerBarkskin() {
	actionID := core.ActionID{SpellID: barkskinRank.SpellID}
	multiplier := 1 + barkskinRank.Effect(shared.A_MOD_DAMAGE_PERCENT_TAKEN, 1).Value/100

	barkskinAura := druid.RegisterAura(core.Aura{
		Label:    "Barkskin",
		ActionID: actionID,
		Duration: barkskinRank.Duration,
	}).AttachMultiplicativePseudoStatBuff(&druid.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexPhysical], multiplier)

	druid.Barkskin = druid.RegisterSpell(Any, core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: barkskinRank.SpellSchool,
		DefenseType: barkskinRank.DefenseType,
		Flags:       core.SpellFlagAPL,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: barkskinRank.Cooldown,
			},
			DefaultCast: core.Cast{
				GCD: barkskinRank.GCD,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			barkskinAura.Activate(sim)
			if sim.CurrentTime > 0 {
				druid.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime)
			}
		},

		RelatedSelfBuff: barkskinAura,
	})

	druid.AddMajorCooldown(core.MajorCooldown{
		Spell: druid.Barkskin.Spell,
		Type:  core.CooldownTypeSurvival,
		// Manual only, as upstream's Barkskin: on cooldown it spends what a tank keeps for real damage
		// (master never auto-uses survival cooldowns either).
		ShouldActivate: func(_ *core.Simulation, _ *core.Character) bool { return false },
	})
}
