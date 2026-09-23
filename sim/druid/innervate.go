package druid

// TODO: To be implemented.
// Returns the time to wait before the next action, or 0 if innervate is on CD
// or disabled.
func (druid *Druid) registerInnervateCD() {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// innervateTarget := druid.GetUnit(druid.SelfBuffs.InnervateTarget)
	// if innervateTarget == nil {
	// 	innervateTarget = &druid.Unit
	// }
	// innervateTargetChar := druid.Env.Raid.GetPlayerFromUnit(innervateTarget).GetCharacter()
	//
	// actionID := core.ActionID{SpellID: 29166, Tag: druid.Index}
	// var innervateSpell *DruidSpell
	//
	// innervateCD := core.InnervateCD
	//
	// amount := 0.05
	// if innervateTarget == &druid.Unit {
	// 	// TODO: Forever drops Dreamstate; untalented (self-cast base amount) until we
	// 	// know whether the effect moved onto another talent.
	// 	amount = 0.2
	// }
	//
	// var innervateAura = core.InnervateAura(innervateTargetChar, amount, actionID.Tag)
	//
	// innervateSpell = druid.RegisterSpell(core.SpellConfig{
	// 	ActionID:        actionID,
	// 	CastRequirement: spellData.Innervate.Highest().CastRequirement(),
	// 	DefenseType:     core.DefenseTypeMagic,
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			GCD: core.GCDDefault,
	// 		},
	// 		CD: core.Cooldown{
	// 			Timer:    druid.NewTimer(),
	// 			Duration: innervateCD,
	// 		},
	// 	},
	// 	ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
	// 		// If target already has another innervate, don't cast.
	// 		return !innervateTarget.HasActiveAuraWithTag(core.InnervateAuraTag)
	// 	},
	// 	ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
	// 		innervateAura.Activate(sim)
	// 	},
	// })
	//
	// druid.AddMajorCooldown(core.MajorCooldown{
	// 	Spell: innervateSpell.Spell,
	// 	Type:  core.CooldownTypeMana,
	// 	ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
	// 		// Require manual APL usage
	// 		return false
	// 	},
	// })
}
