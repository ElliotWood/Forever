package mage

import (
	"time"

	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

// Evocation raises spirit regen by the row's 1500% and lets it run in full while channeling.
func (mage *Mage) registerEvocation() {
	evocationRank := spellData.Evocation.HighestRank()
	regenMultiplier := evocationRank.Effect(shared.A_MOD_POWER_REGEN_PERCENT, 0).Value / 100
	actionID := core.ActionID{SpellID: evocationRank.SpellID}

	// The row states the channel's length, not a period; ticks only mark the channel.
	tickLength := time.Millisecond * 250

	regenAura := mage.RegisterAura(core.Aura{
		Label:    "Evocation Regen",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnGain: func(_ *core.Aura, _ *core.Simulation) {
			mage.PseudoStats.SpiritRegenMultiplier += regenMultiplier
			mage.PseudoStats.ForceFullSpiritRegen = true
			mage.UpdateManaRegenRates()
		},
		OnExpire: func(_ *core.Aura, _ *core.Simulation) {
			mage.PseudoStats.SpiritRegenMultiplier -= regenMultiplier
			mage.PseudoStats.ForceFullSpiritRegen = false
			mage.UpdateManaRegenRates()
		},
	})

	evocation := mage.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagHelpful | core.SpellFlagChanneled | core.SpellFlagAPL,
		ClassSpellMask: MageSpellEvocation,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: evocationRank.GCD,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: evocationRank.Cooldown,
			},
		},

		Hot: core.DotConfig{
			SelfOnly: true,
			Aura: core.Aura{
				Label: "Evocation",
				OnGain: func(_ *core.Aura, sim *core.Simulation) {
					regenAura.Activate(sim)
				},
				OnExpire: func(_ *core.Aura, sim *core.Simulation) {
					regenAura.Deactivate(sim)
				},
			},
			NumberOfTicks: int32(evocationRank.Duration / tickLength),
			TickLength:    tickLength,
			OnTick:        func(_ *core.Simulation, _ *core.Unit, _ *core.Dot) {},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			spell.SelfHot().Apply(sim)
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: evocation,
		Type:  core.CooldownTypeMana,
		ShouldActivate: func(_ *core.Simulation, _ *core.Character) bool {
			return false // Left to the APL.
		},
	})
}
