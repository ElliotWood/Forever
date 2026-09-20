package warlock

// Package-level state the commented-out implementations used:
// var shadowFuryRank = spellData.Shadowfury.BySpellID(30414)
// var shadowFuryCoeff = shadowFuryRank.Direct.BonusCoefficient()

// TODO: uncalled -- Forever drops the Shadowfury talent that gated this spell, so
// nothing calls this any more. Kept rather than deleted because "the talent is gone"
// and "the spell is gone" are not the same claim, and the client data does not
// distinguish them. Re-gate before wiring it back up.
// TODO: To be implemented, or removed. The Forever client does not ship this spell at all --
// no SpellName row carries the name -- so there is nothing to build a registrar from. The
// body below is kept commented as the record of the TBC implementation.
func (warlock *Warlock) registerShadowfury() {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	//
	// warlock.Shadowfury = warlock.RegisterSpell(core.SpellConfig{
	// 	ActionID:       core.ActionID{SpellID: shadowFuryRank.SpellID},
	// 	SpellSchool:    shadowFuryRank.SpellSchool,
	// 	ProcMask:       core.ProcMaskSpellDamage,
	// 	Flags:          core.SpellFlagAPL,
	// 	ClassSpellMask: WarlockSpellShadowFury,
	//
	// 	ManaCost: core.ManaCostOptions{FlatCost: shadowFuryRank.Cost},
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			GCD: shadowFuryRank.GCD,
	// 			// Below core's 1s floor, so it has to be named as the floor too - see GCDTime.
	// 			GCDMin:   shadowFuryRank.GCD,
	// 			CastTime: shadowFuryRank.CastTime,
	// 		},
	// 		CD: core.Cooldown{
	// 			Timer:    warlock.NewTimer(),
	// 			Duration: shadowFuryRank.Cooldown,
	// 		},
	// 	},
	//
	// 	DamageMultiplier: 1,
	// 	DefenseType:      shadowFuryRank.DefenseType,
	// 	ThreatMultiplier: 1,
	// 	BonusCoefficient: shadowFuryCoeff,
	//
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		dmgRoll := shadowFuryRank.Direct.Damage(sim)
	// 		result := spell.CalcDamage(sim, target, dmgRoll, spell.OutcomeMagicHitAndCrit)
	// 		spell.WaitTravelTime(sim, func(sim *core.Simulation) {
	// 			spell.DealDamage(sim, result)
	// 		})
	// 	},
	// })
}
