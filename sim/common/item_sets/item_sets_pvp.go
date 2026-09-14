package item_sets

import (
	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/stats"
)

///////////////////////////////////////////////////////////////////////////
//                                 Mail
///////////////////////////////////////////////////////////////////////////

var ItemSetTheHighlandersFortitude = core.NewItemSet(core.ItemSet{
	Name: "The Highlander's Fortitude",
	Bonuses: map[int32]core.ApplyEffect{
		// Increase Stamina +5
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Stamina, 5)
		},
		// +1 Crit with Spells.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.SpellCrit, 1)
		},
	},
})
