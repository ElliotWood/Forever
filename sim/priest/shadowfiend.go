package priest

// Package-level state the commented-out implementations used:
// var shadowfiendRank = spellData.Shadowfiend.BySpellID(34433)

// TODO: To be implemented. The Forever client ships this as a single unranked class spell:
// it has a SkillLineAbility row but no "Rank N" subtext, so no ladder can be built for it.
func (priest *Priest) registerShadowfiendSpell() {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// actionID := core.ActionID{SpellID: shadowfiendRank.SpellID}
	//
	// // Timeline aura
	// priest.ShadowfiendAura = priest.RegisterAura(core.Aura{
	// 	ActionID: actionID,
	// 	Label:    "Shadowfiend",
	// 	Duration: time.Second * 15,
	// })
	//
	// priest.Shadowfiend = priest.RegisterSpell(core.SpellConfig{
	// 	ActionID:       actionID,
	// 	SpellSchool:    core.SpellSchoolShadow,
	// 	ProcMask:       core.ProcMaskEmpty,
	// 	Flags:          core.SpellFlagAPL,
	// 	ClassSpellMask: PriestSpellShadowFiend,
	// 	Rank:           shadowfiendRank.Rank,
	//
	// 	ManaCost: core.ManaCostOptions{
	// 		BaseCostPercent: 6,
	// 	},
	//
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			GCD: shadowfiendRank.GCD,
	// 		},
	// 		CD: core.Cooldown{
	// 			Timer:    priest.NewTimer(),
	// 			Duration: time.Minute * 5,
	// 		},
	// 	},
	//
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		priest.ShadowfiendPet.EnableWithTimeout(sim, priest.ShadowfiendPet, spell.RelatedSelfBuff.Duration)
	// 		spell.RelatedSelfBuff.Activate(sim)
	// 	},
	//
	// 	RelatedSelfBuff: priest.ShadowfiendAura,
	// })
}
