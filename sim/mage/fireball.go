package mage

var fireballRank = spellData.Fireball.HighestRank()

// TODO: To be implemented. TBC body below needs no porting; kept commented until this class's port is reviewed.
func (mage *Mage) registerFireballSpell() {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// fireballTick := fireballRank.Periodic.(shared.SpellDataPeriodic)
	//
	// mage.RegisterSpell(core.SpellConfig{
	// 	ActionID:       core.ActionID{SpellID: fireballRank.SpellID},
	// 	SpellSchool:    fireballRank.SpellSchool,
	// 	DefenseType:    fireballRank.DefenseType,
	// 	ProcMask:       core.ProcMaskSpellDamage,
	// 	Flags:          core.SpellFlagAPL,
	// 	ClassSpellMask: MageSpellFireball,
	// 	MissileSpeed:   fireballRank.MissileSpeed,
	//
	// 	ManaCost: core.ManaCostOptions{
	// 		FlatCost: fireballRank.Cost,
	// 	},
	//
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			GCD:      fireballRank.GCD,
	// 			CastTime: fireballRank.CastTime,
	// 		},
	// 	},
	//
	// 	Dot: core.DotConfig{
	// 		Aura: core.Aura{
	// 			Label: "FireballDoT",
	// 		},
	// 		NumberOfTicks: fireballTick.NumberOfTicks,
	// 		TickLength:    fireballTick.TickLength,
	// 		OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
	// 			dot.Snapshot(target, fireballTick.Tick)
	// 		},
	// 		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
	// 			dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
	// 		},
	// 	},
	//
	// 	DamageMultiplier: 1,
	// 	BonusCoefficient: fireballRank.Direct.BonusCoefficient(),
	// 	ThreatMultiplier: 1,
	//
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		baseDamage := fireballRank.Direct.Damage(sim)
	// 		result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
	//
	// 		spell.WaitTravelTime(sim, func(s *core.Simulation) {
	// 			spell.DealDamage(sim, result)
	// 			if result.Landed() {
	// 				spell.Dot(target).Apply(sim)
	// 			}
	// 		})
	//
	// 	},
	// })
}
