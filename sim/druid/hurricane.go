package druid

var hurricaneRank = spellData.Hurricane.HighestRank()

// TODO: To be implemented.
func (druid *Druid) registerHurricaneSpell() {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// druid.Hurricane = druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
	// 	ActionID:       core.ActionID{SpellID: hurricaneRank.SpellID},
	// 	SpellSchool:    hurricaneRank.SpellSchool,
	// 	DefenseType:    hurricaneRank.DefenseType,
	// 	ProcMask:       core.ProcMaskSpellDamage,
	// 	Flags:          core.SpellFlagChanneled | core.SpellFlagAPL,
	// 	ClassSpellMask: DruidSpellHurricane,
	// 	MaxRange:       hurricaneRank.MaxRange,
	//
	// 	ManaCost: core.ManaCostOptions{
	// 		FlatCost: hurricaneRank.Cost,
	// 	},
	// 	// TODO: Forever states no cooldown on Hurricane (the client rows carry none), so the
	// 	// spell is registered without one rather than with an invented duration.
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			GCD: hurricaneRank.GCD,
	// 		},
	// 	},
	// 	Dot: core.DotConfig{
	// 		IsAOE: true,
	// 		Aura: core.Aura{
	// 			Label: "Hurricane (Aura)",
	// 		},
	// 		NumberOfTicks:       10,
	// 		TickLength:          time.Second * 1,
	// 		AffectedByCastSpeed: true,
	// 		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
	// 			druid.Hurricane.RelatedDotSpell.Cast(sim, target)
	// 		},
	// 	},
	//
	// 	ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
	// 		spell.AOEDot().Apply(sim)
	// 	},
	// })
	//
	// druid.Hurricane.RelatedDotSpell = druid.Unit.RegisterSpell(core.SpellConfig{
	// 	ActionID:       core.ActionID{SpellID: 42230},
	// 	SpellSchool:    core.SpellSchoolNature,
	// 	DefenseType:    core.DefenseTypeMagic,
	// 	ProcMask:       core.ProcMaskSpellDamage,
	// 	ClassSpellMask: DruidSpellHurricane,
	// 	// 42230 is the tick the channel triggers, a proc rather than a cast.
	// 	Flags: core.SpellFlagProc,
	//
	// 	DamageMultiplier: 1,
	// 	ThreatMultiplier: 1,
	// 	// TODO: Forever moves Hurricane's damage onto the area trigger its second effect
	// 	// creates, which the client tables do not carry, so the rank has no Direct value at
	// 	// all and the tick is pinned to no damage rather than an invented one.
	// 	BonusCoefficient: shared.SpellDataCoef(hurricaneRank.Direct),
	//
	// 	ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
	// 		spell.CalcAndDealAoeDamage(sim, 0, spell.OutcomeMagicHit)
	// 	},
	// })
}
