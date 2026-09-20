package paladin

// Package-level state the commented-out implementations used:
// var spiritualAttunementRank = spellData.SpiritualAttunement.BySpellID(33776)

// Spiritual Attunement (Rank 2, SpellID 33776): Whenever you are healed by another character's spell,
// you regain 10% of the amount healed as mana.
// In the sim, this is modeled as mana return from damage taken (since the healing model offsets damage).
//
// TODO: To be implemented, or removed. The Forever client does not ship this spell at all --
// no SpellName row carries the name -- so there is nothing to build a registrar from. The
// body below is kept commented as the record of the TBC implementation.
func (paladin *Paladin) RegisterSpiritualAttunement() {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// manaMetrics := paladin.NewManaMetrics(core.ActionID{SpellID: spiritualAttunementRank.SpellID})
	//
	// paladin.MakeProcTriggerAura(core.ProcTrigger{
	// 	Name:               "Spiritual Attunement",
	// 	CanProcFromProcs:   true, // 31785/33776 carry the bit: proc heals count.
	// 	ActionID:           core.ActionID{SpellID: spiritualAttunementRank.SpellID},
	// 	Callback:           core.CallbackOnSpellHitTaken,
	// 	RequireDamageDealt: true,
	// 	Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
	// 		coeff := shared.SpellDataMin(spiritualAttunementRank.Direct) / 100
	// 		// Lightbringer Armor 2pc: +10% mana from Spiritual Attunement
	// 		if paladin.T6_4pcAura.IsActive() {
	// 			coeff *= 1.1
	// 		}
	// 		paladin.AddMana(sim, result.Damage*coeff, manaMetrics)
	// 	},
	// })
}
