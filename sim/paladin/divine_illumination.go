package paladin

// TODO: To be implemented. The TBC body below is otherwise a clean port; kept commented until this class's port is reviewed.
//
// Divine Illumination (Talent)
// https://www.wowhead.com/forever/spell=31842
//
// Reduces the mana cost of all spells by 50% for 15 sec.
//
// TODO: uncalled -- Forever drops the Divine Illumination talent; re-gate before wiring
// back into registerTalentSpells.
func (paladin *Paladin) registerDivineIllumination() {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// actionId := core.ActionID{SpellID: 31842}
	// divineIlluminationAura := paladin.RegisterAura(core.Aura{
	// 	Label:    "Divine Illumination" + paladin.Name,
	// 	ActionID: actionId,
	// 	Duration: time.Second * 15,
	// }).AttachSpellMod(core.SpellModConfig{
	// 	Kind:       core.SpellMod_PowerCost_Pct,
	// 	FloatValue: -0.5,
	// })
	//
	// divineIlluminationSpell := paladin.RegisterSpell(core.SpellConfig{
	// 	ActionID:       actionId,
	// 	SpellSchool:    core.SpellSchoolHoly,
	// 	DefenseType:    core.DefenseTypeMagic,
	// 	ProcMask:       core.ProcMaskEmpty,
	// 	Flags:          core.SpellFlagAPL,
	// 	ClassSpellMask: SpellMaskDivineIllumination,
	//
	// 	Cast: core.CastConfig{
	// 		CD: core.Cooldown{
	// 			Timer:    paladin.NewTimer(),
	// 			Duration: time.Minute * 3,
	// 		},
	// 	},
	//
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		spell.RelatedSelfBuff.Activate(sim)
	// 	},
	//
	// 	RelatedSelfBuff: divineIlluminationAura,
	// })
	//
	// paladin.AddMajorCooldown(core.MajorCooldown{
	// 	Spell:    divineIlluminationSpell,
	// 	Priority: core.CooldownPriorityLow,
	// 	Type:     core.CooldownTypeMana,
	// })
}
