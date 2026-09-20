package shaman

// Package-level state the commented-out implementations used:
// var fireElementalTotemRank = spellData.FireElementalTotem.BySpellID(2894)

// TODO: To be implemented, or removed. The Forever client does not ship this spell at all --
// no SpellName row carries the name -- so there is nothing to build a registrar from. The
// body below is kept commented as the record of the TBC implementation.
func (shaman *Shaman) registerFireElementalTotem() {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	//
	// actionID := core.ActionID{SpellID: fireElementalTotemRank.SpellID}
	//
	// totalDuration := time.Second * 120
	//
	// fireElementalAura := shaman.RegisterAura(core.Aura{
	// 	Label:    "Fire Elemental Totem",
	// 	ActionID: actionID,
	// 	Duration: totalDuration,
	// 	OnExpire: func(aura *core.Aura, sim *core.Simulation) {
	// 		shaman.FireElemental.Disable(sim)
	// 	},
	// })
	//
	// shaman.FireElementalTotem = shaman.RegisterSpell(core.SpellConfig{
	// 	ActionID:       actionID,
	// 	DefenseType:    core.DefenseTypeMagic,
	// 	Flags:          core.SpellFlagAPL | SpellFlagInstant,
	// 	ClassSpellMask: SpellMaskFireElementalTotem,
	// 	ManaCost: core.ManaCostOptions{
	// 		FlatCost: fireElementalTotemRank.Cost,
	// 	},
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			GCD: fireElementalTotemRank.GCD,
	// 		},
	// 		CD: core.Cooldown{
	// 			Timer:    shaman.NewTimer(),
	// 			Duration: fireElementalTotemRank.Cooldown,
	// 		},
	// 		SharedCD: core.Cooldown{
	// 			Timer:    shaman.GetOrInitTimer(&shaman.ElementalSharedCDTimer),
	// 			Duration: time.Minute * 1,
	// 		},
	// 	},
	//
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, _ *core.Spell) {
	// 		shaman.cancelFireTotems(sim)
	// 		shaman.TotemExpirations[FireTotem] = sim.CurrentTime + fireElementalAura.Duration
	//
	// 		shaman.FireElemental.Disable(sim)
	// 		shaman.FireElemental.EnableWithTimeout(sim, shaman.FireElemental, fireElementalAura.Duration)
	//
	// 		// Add a dummy aura to show in metrics
	// 		fireElementalAura.Activate(sim)
	// 	},
	// 	RelatedSelfBuff: fireElementalAura,
	// })
	//
	// shaman.AddMajorCooldown(core.MajorCooldown{
	// 	Spell: shaman.FireElementalTotem,
	// 	Type:  core.CooldownTypeDPS,
	// 	ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
	// 		// Fele should only be cast by manual APL intervention
	// 		return false
	// 	},
	// })
}
