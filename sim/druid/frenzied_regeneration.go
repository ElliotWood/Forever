package druid

// Package-level state the commented-out implementations used:
// var frenziedRegenerationRank = spellData.FrenziedRegeneration.BySpellID(26999)
// var frenziedRegenerationTick = frenziedRegenerationRank.Periodic.(shared.SpellDataPeriodic)

// TODO: To be implemented. The Forever client ships this as a single unranked class spell:
// it has a SkillLineAbility row but no "Rank N" subtext, so no ladder can be built for it.
func (druid *Druid) registerFrenziedRegenerationSpell() {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// actionID := core.ActionID{SpellID: frenziedRegenerationRank.SpellID}
	// rageMetrics := druid.NewRageMetrics(actionID)
	//
	// druid.FrenziedRegenerationAura = druid.RegisterAura(core.Aura{
	// 	Label:    "Frenzied Regeneration",
	// 	ActionID: actionID,
	// 	Duration: 10 * time.Second,
	// })
	//
	// // Deactivate when leaving Bear Form.
	// druid.BearFormAura.ApplyOnExpire(func(_ *core.Aura, sim *core.Simulation) {
	// 	druid.FrenziedRegenerationAura.Deactivate(sim)
	// })
	//
	// druid.FrenziedRegeneration = druid.RegisterSpell(Bear, core.SpellConfig{
	// 	ActionID:         actionID,
	// 	SpellSchool:      core.SpellSchoolPhysical,
	// 	ProcMask:         core.ProcMaskEmpty,
	// 	ClassSpellMask:   DruidSpellFrenziedRegeneration,
	// 	Flags:            core.SpellFlagAPL,
	// 	DamageMultiplier: 1,
	//
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			GCD: frenziedRegenerationRank.GCD,
	// 		},
	// 		CD: core.Cooldown{
	// 			Timer:    druid.NewTimer(),
	// 			Duration: frenziedRegenerationRank.Cooldown,
	// 		},
	// 		IgnoreHaste: true,
	// 	},
	//
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		druid.FrenziedRegenerationAura.Activate(sim)
	// 		// Converts up to 10 rage per second into 25 health per rage, for 10 sec.
	// 		core.StartPeriodicAction(sim, core.PeriodicActionOptions{
	// 			Period:   frenziedRegenerationTick.TickLength,
	// 			NumTicks: int(frenziedRegenerationTick.NumberOfTicks),
	// 			Priority: core.ActionPriorityDOT,
	// 			OnAction: func(sim *core.Simulation) {
	// 				rage := min(druid.CurrentRage(), 10)
	// 				if rage > 0 {
	// 					druid.SpendRage(sim, rage, rageMetrics)
	// 					spell.CalcAndDealPeriodicHealing(sim, &druid.Unit, rage*frenziedRegenerationTick.Tick, spell.OutcomeHealing)
	// 				}
	// 			},
	// 		})
	// 	},
	// })
}
