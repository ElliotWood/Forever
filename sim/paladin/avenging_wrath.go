package paladin

// TODO: To be implemented. TBC body below needs no porting; kept commented until this class's port is reviewed.
//
// Avenging Wrath
// https://www.wowhead.com/forever/spell=31884
//
// Increases all damage caused by 30% for 20 sec.
// Causes Forebearance, preventing the use of Divine Shield,
// Divine Protection, Blessing of Protection again for 1 min.
func (paladin *Paladin) registerAvengingWrath() {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	//
	// actionID := core.ActionID{SpellID: 31884}
	// avengingWrathAura := paladin.RegisterAura(core.Aura{
	// 	Label:    "Avenging Wrath" + paladin.Label,
	// 	ActionID: actionID,
	// 	Duration: time.Second * 20,
	// }).AttachMultiplicativePseudoStatBuff(
	// 	&paladin.PseudoStats.DamageDealtMultiplier, 1.3,
	// )
	//
	// avengingWrath := paladin.RegisterSpell(core.SpellConfig{
	// 	ActionID:       actionID,
	// 	SpellSchool:    core.SpellSchoolHoly,
	// 	ProcMask:       core.ProcMaskEmpty,
	// 	Flags:          core.SpellFlagAPL,
	// 	ClassSpellMask: SpellMaskAvengingWrath,
	//
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			NonEmpty: true,
	// 		},
	// 		CD: core.Cooldown{
	// 			Timer:    paladin.NewTimer(),
	// 			Duration: time.Minute * 3,
	// 		},
	// 	},
	// 	ManaCost: core.ManaCostOptions{
	// 		BaseCostPercent: 8,
	// 	},
	//
	// 	ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
	// 		return !paladin.Forbearance.IsActive()
	// 	},
	//
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		spell.RelatedSelfBuff.Activate(sim)
	// 		paladin.Forbearance.Activate(sim)
	// 	},
	//
	// 	RelatedSelfBuff: avengingWrathAura,
	// })
	//
	// paladin.AddMajorCooldown(core.MajorCooldown{
	// 	Spell:    avengingWrath,
	// 	Priority: int32(core.CooldownTypeDPS),
	// 	Type:     core.CooldownTypeDPS,
	// })
}
