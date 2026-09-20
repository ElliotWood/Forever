package hunter

var multiShotRank = spellData.MultiShot.HighestRank()

// TODO: To be implemented.
func (hunter *Hunter) registerMultiShotSpell() {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// hunter.MultiShot = hunter.RegisterRangedSpell(core.SpellConfig{
	// 	ActionID:       core.ActionID{SpellID: multiShotRank.SpellID},
	// 	SpellSchool:    multiShotRank.SpellSchool,
	// 	DefenseType:    multiShotRank.DefenseType,
	// 	ProcMask:       core.ProcMaskRangedSpecial,
	// 	ClassSpellMask: HunterSpellMultiShot,
	// 	Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
	//
	// 	MissileSpeed: multiShotRank.MissileSpeed,
	//
	// 	ManaCost: core.ManaCostOptions{
	// 		FlatCost: multiShotRank.Cost,
	// 	},
	//
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			CastTime: time.Millisecond * 500,
	// 		},
	// 		CD: core.Cooldown{
	// 			Timer:    hunter.NewTimer(),
	// 			Duration: multiShotRank.Cooldown,
	// 		},
	// 	},
	//
	// 	BonusCoefficient: multiShotRank.Direct.BonusCoefficient(),
	//
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		baseDamage := spell.RangedAttackPower(target)*0.2 +
	// 			hunter.AutoAttacks.Ranged().BaseDamage(sim) +
	// 			hunter.talonOfAlarBonus() +
	// 			multiShotRank.Direct.Damage(sim)
	//
	// 		spell.CalcAoeDamage(sim, baseDamage, spell.OutcomeRangedHitAndCrit)
	//
	// 		spell.WaitTravelTime(sim, func(sim *core.Simulation) {
	// 			spell.DealBatchedAoeDamage(sim)
	// 		})
	// 	},
	// })
}
