package mage

var blizzardRank = spellData.Blizzard.HighestRank()

// TODO: To be implemented. TBC body below needs no porting; kept commented until this class's port is reviewed.
func (mage *Mage) registerBlizzardSpell() {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	//
	// blizzardActionId := core.ActionID{SpellID: blizzardRank.SpellID}
	//
	// // https://wago.tools/db2/SpellEffect?build=2.5.5.65295&filter%5BSpellID%5D=42208
	// blizzardCoefficient := 0.11900000274
	//
	// blizzardTickSpell := mage.RegisterSpell(core.SpellConfig{
	// 	ActionID:       blizzardActionId,
	// 	SpellSchool:    core.SpellSchoolFrost,
	// 	DefenseType:    core.DefenseTypeMagic,
	// 	ProcMask:       core.ProcMaskSpellDamage,
	// 	ClassSpellMask: MageSpellBlizzard,
	//
	// 	DamageMultiplier: 1,
	// 	BonusCoefficient: blizzardCoefficient,
	// 	ThreatMultiplier: 1,
	//
	// 	ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
	// 		spell.CalcAndDealAoeDamage(sim, 184, spell.OutcomeMagicHit)
	// 	},
	// })
	//
	// mage.RegisterSpell(core.SpellConfig{
	// 	ActionID:       blizzardActionId,
	// 	SpellSchool:    core.SpellSchoolFrost,
	// 	DefenseType:    core.DefenseTypeMagic,
	// 	ProcMask:       core.ProcMaskSpellDamage,
	// 	Flags:          core.SpellFlagChanneled | core.SpellFlagAPL,
	// 	ClassSpellMask: MageSpellBlizzard,
	// 	ManaCost: core.ManaCostOptions{
	// 		FlatCost: blizzardRank.Cost,
	// 	},
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			GCD: blizzardRank.GCD,
	// 		},
	// 	},
	// 	Dot: core.DotConfig{
	// 		IsAOE: true,
	// 		Aura: core.Aura{
	// 			Label:    "Blizzard",
	// 			ActionID: blizzardActionId,
	// 		},
	// 		NumberOfTicks:        8,
	// 		TickLength:           time.Second * 1,
	// 		AffectedByCastSpeed:  true,
	// 		HasteReducesDuration: true,
	// 		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
	// 			blizzardTickSpell.Cast(sim, target)
	// 		},
	// 	},
	//
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		spell.AOEDot().Apply(sim)
	// 	},
	// })
}
