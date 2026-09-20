package paladin

// Divine Favor
// https://www.wowhead.com/forever/spell=20216
//
// When activated, gives your next Flash of Light, Holy Light, or Holy Shock
// spell a 100% critical strike chance.

// TODO: To be implemented. TBC body below needs no porting; kept commented until this class's port is reviewed.
func (paladin *Paladin) registerDivineFavor() {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// actionID := core.ActionID{SpellID: 20216}
	//
	// var divineFavorAura *core.Aura
	// divineFavorAura = paladin.RegisterAura(core.Aura{
	// 	Label:    "Divine Favor" + paladin.Label,
	// 	ActionID: actionID,
	// 	Duration: core.NeverExpires,
	// }).AttachSpellMod(core.SpellModConfig{
	// 	Kind:       core.SpellMod_BonusCrit_Percent,
	// 	ClassMask:  SpellMaskHolyLight | SpellMaskFlashOfLight | SpellMaskHolyShock,
	// 	FloatValue: 100,
	// }).AttachProcTrigger(core.ProcTrigger{
	// 	CanProcFromProcs:   true, // 20216 carries the bit.
	// 	Callback:           core.CallbackOnCastComplete,
	// 	ClassSpellMask:     SpellMaskHolyLight | SpellMaskFlashOfLight | SpellMaskHolyShock,
	// 	TriggerImmediately: true,
	// 	Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
	// 		divineFavorAura.Deactivate(sim)
	// 	},
	// })
	//
	// divineFavor := paladin.RegisterSpell(core.SpellConfig{
	// 	ActionID:       actionID,
	// 	Flags:          core.SpellFlagAPL | core.SpellFlagHelpful,
	// 	ClassSpellMask: SpellMaskDivineFavor,
	// 	SpellSchool:    core.SpellSchoolHoly,
	// 	DefenseType:    core.DefenseTypeMagic,
	//
	// 	ManaCost: core.ManaCostOptions{
	// 		BaseCostPercent: 3,
	// 	},
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			NonEmpty: true,
	// 		},
	// 		CD: core.Cooldown{
	// 			Timer:    paladin.NewTimer(),
	// 			Duration: 2 * time.Minute,
	// 		},
	// 	},
	//
	// 	ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
	// 		spell.RelatedSelfBuff.Activate(sim)
	// 	},
	//
	// 	RelatedSelfBuff: divineFavorAura,
	// })
	//
	// paladin.AddMajorCooldown(core.MajorCooldown{
	// 	Spell: divineFavor,
	// 	Type:  core.CooldownTypeDPS,
	// })
}
