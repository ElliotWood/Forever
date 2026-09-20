package druid

import (
	"github.com/wowsims/forever/sim/core"
)

func init() {
	// Idol of the Moon
	core.NewItemEffect(23197, func(agent core.Agent) {
		character := agent.GetCharacter()
		aura := core.MakePermanent(character.RegisterAura(core.Aura{
			Label: "Improved Moonfire",
		}).AttachSpellMod(core.SpellModConfig{
			ClassMask:  DruidSpellMoonfire,
			Kind:       core.SpellMod_BaseDamage_Flat,
			FloatValue: 33.0,
		}))

		// TODO: this registers the swap proc against 32330, not 23197. That item is not in
		// the Forever database either, so the registration is dead whichever was intended.
		character.ItemSwap.RegisterProc(32330, aura)
	})

	// Wolfshead Helm (8345): When shapeshifting into Cat form the Druid gains 20 energy,
	// when shapeshifting into Bear form the Druid gains 5 rage.
	core.NewItemEffect(8345, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()
		core.MakePermanent(druid.RegisterAura(core.Aura{
			Label:    "Wolfshead Helm",
			ActionID: core.ActionID{SpellID: 17768},
			OnGain: func(_ *core.Aura, _ *core.Simulation) {
				druid.WolfsheadEnergyBonus += 20
				druid.WolfsheadRageBonus += 5
			},
			OnExpire: func(_ *core.Aura, _ *core.Simulation) {
				druid.WolfsheadEnergyBonus -= 20
				druid.WolfsheadRageBonus -= 5
			},
		}))
	})
}
