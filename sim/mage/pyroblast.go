package mage

var pyroblastRank = spellData.Pyroblast.HighestRank()

// TODO: To be implemented. TBC body below needs no porting; kept commented until this class's port is reviewed.
func (mage *Mage) registerPyroblastSpell() {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// actionID := core.ActionID{SpellID: pyroblastRank.SpellID}
	//
	// pyroblastDotCoefficient := 0.05000000075
	// pyroblastTick := pyroblastRank.Periodic.(shared.SpellDataPeriodic)
	//
	// mage.Pyroblast = mage.RegisterSpell(core.SpellConfig{
	// 	ActionID:       actionID,
	// 	SpellSchool:    core.SpellSchoolFire,
	// 	DefenseType:    core.DefenseTypeMagic,
	// 	ProcMask:       core.ProcMaskSpellDamage,
	// 	Flags:          core.SpellFlagAPL,
	// 	ClassSpellMask: MageSpellPyroblast,
	// 	MissileSpeed:   pyroblastRank.MissileSpeed,
	//
	// 	ManaCost: core.ManaCostOptions{
	// 		FlatCost: pyroblastRank.Cost,
	// 	},
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			GCD:      pyroblastRank.GCD,
	// 			CastTime: pyroblastRank.CastTime,
	// 		},
	// 	},
	//
	// 	DamageMultiplier: 1,
	// 	BonusCoefficient: pyroblastRank.Direct.BonusCoefficient(),
	// 	ThreatMultiplier: 1,
	//
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		baseDamage := pyroblastRank.Direct.Damage(sim)
	// 		result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
	//
	// 		spell.WaitTravelTime(sim, func(s *core.Simulation) {
	// 			spell.DealDamage(sim, result)
	// 			if result.Landed() {
	// 				spell.RelatedDotSpell.Cast(sim, target)
	// 			}
	// 		})
	// 	},
	// })
	//
	// mage.Pyroblast.RelatedDotSpell = mage.RegisterSpell(core.SpellConfig{
	// 	ActionID:       actionID.WithTag(1),
	// 	SpellSchool:    core.SpellSchoolFire,
	// 	DefenseType:    core.DefenseTypeMagic,
	// 	ProcMask:       core.ProcMaskSpellDamage,
	// 	ClassSpellMask: MageSpellPyroblastDot,
	// 	Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,
	//
	// 	DamageMultiplier: 1,
	// 	ThreatMultiplier: 1,
	//
	// 	Dot: core.DotConfig{
	// 		Aura: core.Aura{
	// 			Label: "PyroblastDoT",
	// 		},
	// 		NumberOfTicks:    pyroblastTick.NumberOfTicks,
	// 		TickLength:       pyroblastTick.TickLength,
	// 		BonusCoefficient: pyroblastDotCoefficient,
	// 		OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
	// 			dot.Snapshot(target, pyroblastTick.Tick)
	// 		},
	// 		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
	// 			dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
	// 		},
	// 	},
	//
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		spell.Dot(target).Apply(sim)
	// 	},
	// })
}
