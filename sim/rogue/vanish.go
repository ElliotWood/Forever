package rogue

var vanishRank = spellData.Vanish.BySpellID(1856)

// TODO: To be implemented. Vanish already resolves against Forever data
// (spellData.Vanish.BySpellID(1856)); the TBC body needs review before it's uncommented.
func (rogue *Rogue) registerVanishSpell() {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// rogue.Vanish = rogue.RegisterSpell(core.SpellConfig{
	// 	ActionID:       core.ActionID{SpellID: vanishRank.SpellID},
	// 	SpellSchool:    vanishRank.SpellSchool,
	// 	Flags:          core.SpellFlagAPL,
	// 	ClassSpellMask: RogueSpellVanish,
	//
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			GCD: 0,
	// 		},
	// 		IgnoreHaste: true,
	// 		CD: core.Cooldown{
	// 			Timer:    rogue.NewTimer(),
	// 			Duration: vanishRank.Cooldown,
	// 		},
	// 	},
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		// Pause auto attacks
	// 		rogue.AutoAttacks.CancelAutoSwing(sim)
	// 		// Apply stealth
	// 		rogue.StealthAura.Activate(sim)
	// 	},
	// })
}
