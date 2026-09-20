package mage

// The shared and core imports belong with the commented implementation.

var blizzardRank = spellData.Blizzard.HighestRank()

// TODO: To be implemented. TBC body below needs no porting; kept commented until this class's port is reviewed.
func (mage *Mage) registerBlizzardSpell() {
	panic("To be implemented")

	// The ported implementation, kept until this class is done:
	// blizzardActionId := core.ActionID{SpellID: blizzardRank.SpellID}
	// blizzardTick := blizzardRank.Periodic.(shared.SpellDataPeriodic)
	//
	// blizzardTickSpell := mage.RegisterSpell(core.SpellConfig{
	// 	ActionID:       blizzardActionId,
	// 	SpellSchool:    core.SpellSchoolFrost,
	// 	DefenseType:    core.DefenseTypeMagic,
	// 	ProcMask:       core.ProcMaskSpellDamage,
	// 	ClassSpellMask: MageSpellBlizzard,
	//
	// 	DamageMultiplier: 1,
	// 	BonusCoefficient: blizzardTick.Coef,
	// 	ThreatMultiplier: 1,
	//
	// 	ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
	// 		spell.CalcAndDealAoeDamage(sim, blizzardTick.Tick, spell.OutcomeMagicHit)
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
	// 		NumberOfTicks:        blizzardTick.NumberOfTicks,
	// 		TickLength:           blizzardTick.TickLength,
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
