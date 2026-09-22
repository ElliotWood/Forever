package druid

var thornsRank = spellData.Thorns.HighestRank()

// TODO: To be implemented.
// Self-cast Thorns. The druid's own copy of the generated shield deals the
// client's 22 nature damage; the raid buff registers the external copy before
// Initialize runs, and the two auras bid in ThornsCategory, which holds one of
// them at a time.
func (druid *Druid) registerThornsSpell() {
	panic("To be implemented")

	// The body the port needs:
	// thornsAura := druid.GetAura("Thorns (Player)")
	// if thornsAura == nil {
	// 	thornsAura = core.ThornsAura(&druid.Unit, true, 0)
	// }
	//
	// druid.RegisterSpell(Humanoid, core.SpellConfig{
	// 	ActionID:       core.ActionID{SpellID: thornsRank.SpellID},
	// 	SpellSchool:    thornsRank.SpellSchool,
	// 	DefenseType:    thornsRank.DefenseType,
	// 	Flags:          core.SpellFlagAPL | core.SpellFlagHelpful,
	// 	ClassSpellMask: DruidSpellThorns,
	// 	ProcMask:       core.ProcMaskEmpty,
	// 	MaxRange:       thornsRank.MaxRange,
	//
	// 	ManaCost: core.ManaCostOptions{
	// 		FlatCost: thornsRank.Cost,
	// 	},
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			GCD: thornsRank.GCD,
	// 		},
	// 	},
	//
	// 	ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
	// 		thornsAura.Activate(sim)
	// 	},
	//
	// 	RelatedSelfBuff: thornsAura,
	// })
}
