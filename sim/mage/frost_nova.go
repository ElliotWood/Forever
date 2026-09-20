package mage

var frostNovaRank = spellData.FrostNova.HighestRank()

// TODO: To be implemented. TBC body below needs no porting; kept commented until this class's port is reviewed.
func (mage *Mage) registerFrostNovaSpell() {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	//
	// frostNovaCoefficient := 0.18799999356 // Per https://wago.tools/db2/SpellEffect?build=2.5.5.65295&filter%5BSpellID%5D=exact%253A122 Field "EffetBonusCoefficient"
	//
	// mage.RegisterSpell(core.SpellConfig{
	// 	ActionID:       core.ActionID{SpellID: frostNovaRank.SpellID},
	// 	SpellSchool:    frostNovaRank.SpellSchool,
	// 	DefenseType:    frostNovaRank.DefenseType,
	// 	ProcMask:       core.ProcMaskSpellDamage,
	// 	Flags:          core.SpellFlagAPL | core.SpellFlagBinary,
	// 	ClassSpellMask: MageSpellFrostNova,
	//
	// 	ManaCost: core.ManaCostOptions{
	// 		FlatCost: frostNovaRank.Cost,
	// 	},
	//
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			GCD: frostNovaRank.GCD,
	// 		},
	// 		CD: core.Cooldown{
	// 			Timer:    mage.NewTimer(),
	// 			Duration: frostNovaRank.Cooldown,
	// 		},
	// 	},
	//
	// 	DamageMultiplier: 1,
	// 	BonusCoefficient: frostNovaCoefficient,
	// 	ThreatMultiplier: 1,
	//
	// 	ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
	// 		baseDamage := frostNovaRank.Direct.Damage(sim)
	// 		spell.CalcAndDealAoeDamage(sim, baseDamage, spell.OutcomeMagicHitAndCrit)
	// 	},
	// })
}
