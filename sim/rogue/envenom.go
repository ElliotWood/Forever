package rogue

// TODO: To be implemented. Envenom pins spell 32645 directly in the TBC body; the implementation needs
// review before it's uncommented.
func (rogue *Rogue) registerEnvenom() {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// baseDamage := 180.0 + rogue.DeathmantleBonus
	// apScalingPerComboPoint := 0.03
	//
	// rogue.Envenom = rogue.RegisterSpell(core.SpellConfig{
	// 	ActionID:       core.ActionID{SpellID: 32645},
	// 	SpellSchool:    core.SpellSchoolNature,
	// 	DefenseType:    core.DefenseTypeMelee,
	// 	ProcMask:       core.ProcMaskMeleeMHSpecial, // not core.ProcMaskSpellDamage
	// 	Flags:          core.SpellFlagMeleeMetrics | SpellFlagFinisher | core.SpellFlagAPL,
	// 	MetricSplits:   6,
	// 	ClassSpellMask: RogueSpellEnvenom,
	//
	// 	EnergyCost: core.EnergyCostOptions{
	// 		Cost: 35,
	// 		// TODO: Forever drops Quick Recovery; no energy refund until we know whether the
	// 		// effect moved onto another talent.
	// 		Refund:        0,
	// 		RefundMetrics: rogue.EnergyRefundMetrics,
	// 	},
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			GCD: time.Second,
	// 		},
	// 		IgnoreHaste: true,
	// 		ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
	// 			spell.SetMetricsSplit(rogue.ComboPoints())
	// 		},
	// 	},
	// 	ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
	// 		return rogue.ComboPoints() > 0
	// 	},
	//
	// 	DamageMultiplier:         1,
	// 	DamageMultiplierAdditive: 1,
	// 	ThreatMultiplier:         1,
	//
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		rogue.BreakStealth(sim)
	// 		comboPoints := rogue.ComboPoints()
	// 		dp := rogue.DeadlyPoison.Dot(target)
	// 		consumed := min(dp.GetStacks(), comboPoints)
	//
	// 		baseDamage := baseDamage*float64(consumed) +
	// 			apScalingPerComboPoint*float64(consumed)*spell.MeleeAttackPower(target)
	//
	// 		result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
	//
	// 		if result.Landed() {
	// 			rogue.ApplyFinisher(sim, spell)
	// 			if newStacks := dp.GetStacks() - comboPoints; newStacks > 0 {
	// 				dp.SetStacks(sim, newStacks)
	// 			} else {
	// 				dp.Deactivate(sim)
	// 			}
	// 		} else {
	// 			spell.IssueRefund(sim)
	// 		}
	//
	// 		spell.DealDamage(sim, result)
	// 	},
	// })
}
