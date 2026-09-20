package shaman

import (
	"github.com/wowsims/forever/sim/core"
)

func init() {

	// 	aura := core.MakePermanent(character.RegisterAura(core.Aura{
	// 		Label: "Increased Shock Damage",
	// 	}).AttachSpellMod(core.SpellModConfig{
	// 		Kind:       core.SpellMod_BaseDamage_Flat,
	// 		FloatValue: 30.0,
	// 		ClassMask:  SpellMaskShock,
	// 	}))

	// 	character.ItemSwap.RegisterProc(22395, aura)
	// })

	// Totem of the Storm
	core.NewItemEffect(23199, func(agent core.Agent) {
		character := agent.GetCharacter()

		aura := core.MakePermanent(character.RegisterAura(core.Aura{
			Label: "Increased Lightning Damage",
		}).AttachSpellMod(core.SpellModConfig{
			Kind:       core.SpellMod_BaseDamage_Flat,
			FloatValue: 33.0,
			ClassMask:  SpellMaskLightningBolt | SpellMaskChainLightning | SpellMaskOverload,
		}))

		character.ItemSwap.RegisterProc(23199, aura)
	})

}
