package rogue

import (
	"github.com/wowsims/forever/sim/core"
)

var vanishRank = spellData.Vanish.BySpellID(1856)

func (rogue *Rogue) registerVanishSpell() {
	rogue.Vanish = rogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: vanishRank.SpellID},
		SpellSchool:    vanishRank.SpellSchool,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: RogueSpellVanish,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: 0,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: vanishRank.Cooldown,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// Pause auto attacks
			rogue.AutoAttacks.CancelAutoSwing(sim)
			// Apply stealth
			rogue.StealthAura.Activate(sim)
		},
	})
}
