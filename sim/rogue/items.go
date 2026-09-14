package rogue

import (
	"time"

	"github.com/wowsims/classic/sim/core"
)

const (
	RenatakisCharmofTrickery = 19954
)

func init() {
	core.AddEffectsToTest = false

	// https://www.wowhead.com/classic/item=19954/renatakis-charm-of-trickery
	// Use: Instantly increases your energy by 60. (3 Min Cooldown)
	core.NewItemEffect(RenatakisCharmofTrickery, func(agent core.Agent) {
		rogue := agent.(RogueAgent).GetRogue()
		energyMetrics := rogue.NewEnergyMetrics(core.ActionID{SpellID: 24532})

		spell := rogue.RegisterSpell(core.SpellConfig{
			ActionID: core.ActionID{ItemID: RenatakisCharmofTrickery},
			ProcMask: core.ProcMaskEmpty,
			Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    rogue.NewTimer(),
					Duration: time.Second * 180,
				},
				SharedCD: core.Cooldown{
					Timer:    rogue.GetOffensiveTrinketCD(),
					Duration: time.Second * 10,
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				rogue.AddEnergy(sim, 60, energyMetrics)
			},
		})

		rogue.AddMajorCooldown(core.MajorCooldown{
			Type:  core.CooldownTypeDPS,
			Spell: spell,
			ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
				// Make sure we have plenty of room so we dont energy cap right after using.
				return rogue.CurrentEnergy() <= 40
			},
		})
	})

	core.AddEffectsToTest = true
}
