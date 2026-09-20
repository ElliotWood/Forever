package warlock

// Package-level state the commented-out implementations used:
// var doomRank = spellData.CurseOfDoom.BySpellID(30910)
// var doomTick = doomRank.Periodic.(shared.SpellDataPeriodic)
// var doomCoeff = doomTick.Coef

// TODO: To be implemented. Spells of this name exist in the Forever client, but none of them
// has a class ability row -- no SkillLineAbility entry in a CategoryID 7 skill line -- so the
// generator has no class spell to build a ladder from.
func (warlock *Warlock) registerCurseOfDoom() {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	//
	// calculateBaseDamage := func() float64 {
	// 	damageMultiplier := core.TernaryFloat64(warlock.AmplifyCurseAura != nil && warlock.AmplifyCurseAura.IsActive(), 1.5, 1.0)
	// 	return doomTick.Tick * damageMultiplier
	// }
	//
	// warlock.CurseOfDoom = warlock.RegisterSpell(core.SpellConfig{
	// 	ActionID:       core.ActionID{SpellID: doomRank.SpellID},
	// 	SpellSchool:    doomRank.SpellSchool,
	// 	DefenseType:    doomRank.DefenseType,
	// 	ProcMask:       core.ProcMaskSpellDamage,
	// 	Flags:          core.SpellFlagAPL,
	// 	ClassSpellMask: WarlockSpellCurseOfDoom,
	//
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			GCD: doomRank.GCD,
	// 		},
	// 		CD: core.Cooldown{
	// 			Timer:    warlock.NewTimer(),
	// 			Duration: doomRank.Cooldown,
	// 		},
	// 	},
	//
	// 	ThreatMultiplier: 1,
	// 	DamageMultiplier: 1,
	// 	BonusCoefficient: doomCoeff,
	//
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
	// 		if result.Landed() {
	// 			warlock.DeactivateOtherCurses(sim, spell, target)
	// 			spell.Dot(target).Apply(sim)
	// 		}
	// 		spell.DealOutcome(sim, result)
	// 	},
	//
	// 	Dot: core.DotConfig{
	// 		Aura: core.Aura{
	// 			Label: "Doom",
	// 			Tag:   "Affliction",
	// 		},
	// 		NumberOfTicks:            doomTick.NumberOfTicks,
	// 		TickLength:               doomTick.TickLength,
	// 		BonusCoefficient:         doomCoeff,
	// 		PeriodicDamageMultiplier: 1,
	//
	// 		OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
	// 			dot.Snapshot(target, calculateBaseDamage())
	// 		},
	// 		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
	// 			dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
	// 		},
	// 	},
	//
	// 	ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
	// 		dot := spell.Dot(target)
	// 		if useSnapshot {
	// 			return dot.CalcSnapshotDamage(sim, target, dot.OutcomeTick)
	// 		} else {
	// 			return spell.CalcPeriodicDamage(sim, target, calculateBaseDamage(), spell.OutcomeExpectedMagicHit)
	// 		}
	// 	},
	// })
}
