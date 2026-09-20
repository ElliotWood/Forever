package hunter

var arcaneShotRank = spellData.ArcaneShot.HighestRank()

// TODO: To be implemented.
func (hunter *Hunter) registerArcaneShotSpell() {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// hunter.ArcaneShot = hunter.RegisterRangedSpell(core.SpellConfig{
	// 	ActionID:       core.ActionID{SpellID: arcaneShotRank.SpellID},
	// 	SpellSchool:    arcaneShotRank.SpellSchool,
	// 	DefenseType:    arcaneShotRank.DefenseType,
	// 	ClassSpellMask: HunterSpellArcaneShot,
	// 	ProcMask:       core.ProcMaskRangedSpecial,
	// 	Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
	//
	// 	ManaCost: core.ManaCostOptions{
	// 		FlatCost: arcaneShotRank.Cost,
	// 	},
	//
	// 	Cast: core.CastConfig{
	// 		CD: core.Cooldown{
	// 			Timer:    hunter.NewTimer(),
	// 			Duration: arcaneShotRank.Cooldown,
	// 		},
	// 	},
	//
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		baseDamage := spell.RangedAttackPower(target)*0.15 +
	// 			hunter.talonOfAlarBonus() +
	// 			arcaneShotRank.Direct.Damage(sim)
	//
	// 		result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)
	//
	// 		spell.WaitTravelTime(sim, func(sim *core.Simulation) {
	// 			spell.DealDamage(sim, result)
	// 		})
	// 	},
	// })
}
