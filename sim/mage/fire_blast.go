package mage

// TODO: To be implemented. spellData.FireBlast holds the seven trainer ranks, 2136 to 10199.
// The client also carries 400616 to 400623, the copies the Season of Discovery rune passive
// Overheat (400615) swaps onto the action bar. Overheat is an Engrave grant with no place in
// Forever, and the generator drops its stand-ins.
func (mage *Mage) registerFireBlastSpell() {
	panic("To be implemented")

	//
	// mage.FireBlast = mage.RegisterSpell(core.SpellConfig{
	// 	ActionID:       core.ActionID{SpellID: fireBlastRank.SpellID},
	// 	SpellSchool:    fireBlastRank.SpellSchool,
	// 	DefenseType:    fireBlastRank.DefenseType,
	// 	ProcMask:       core.ProcMaskSpellDamage,
	// 	Flags:          core.SpellFlagAPL,
	// 	ClassSpellMask: MageSpellFireBlast,
	//
	// 	ManaCost: core.ManaCostOptions{
	// 		FlatCost: fireBlastRank.Cost,
	// 	},
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			GCD: fireBlastRank.GCD,
	// 		},
	// 		CD: core.Cooldown{
	// 			Timer:    mage.NewTimer(),
	// 			Duration: fireBlastRank.Cooldown,
	// 		},
	// 	},
	//
	// 	DamageMultiplier: 1,
	// 	BonusCoefficient: fireBlastRank.Direct.BonusCoefficient(),
	//
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		baseDamage := fireBlastRank.Direct.Damage(sim)
	// 		spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
	// 	},
	// })
}
