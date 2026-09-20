package mage

// TODO: uncalled -- Forever drops the Slow talent; re-gate before wiring back into
// registerSpells.
// TODO: To be implemented. The TBC body below is otherwise a clean port; kept commented until this class's port is reviewed.
func (mage *Mage) registerSlowSpell() {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// mage.RegisterSpell(core.SpellConfig{
	// 	ActionID:       core.ActionID{SpellID: 31589},
	// 	ClassSpellMask: MageSpellSlow,
	// 	SpellSchool:    core.SpellSchoolArcane,
	// 	DefenseType:    core.DefenseTypeMagic,
	// 	Flags:          core.SpellFlagAPL | core.SpellFlagBinary,
	//
	// 	ManaCost: core.ManaCostOptions{
	// 		BaseCostPercent: 20,
	// 	},
	//
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			GCD: core.GCDDefault,
	// 		},
	// 	},
	//
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)
	// 		if result.Landed() {
	// 			aura := mage.SlowAuras.Get(target)
	// 			aura.Activate(sim)
	// 		}
	// 	},
	// })
}
