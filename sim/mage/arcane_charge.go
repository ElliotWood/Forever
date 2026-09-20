package mage

// TODO: To be implemented. TBC body below needs no porting; kept commented until this class's port is reviewed.
func (mage *Mage) registerArcaneCharges() {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// powerCostIncrease := 0.75
	// abCostMod := mage.AddDynamicMod(core.SpellModConfig{
	// 	ClassMask:  MageSpellArcaneBlast,
	// 	FloatValue: powerCostIncrease,
	// 	Kind:       core.SpellMod_PowerCost_Pct_Add,
	// })
	//
	// castTimeReduction := time.Millisecond * -334
	// abCastMod := mage.AddDynamicMod(core.SpellModConfig{
	// 	ClassMask: MageSpellArcaneBlast,
	// 	TimeValue: castTimeReduction,
	// 	Kind:      core.SpellMod_CastTime_Flat,
	// })
	//
	// mage.ArcaneChargesAura = core.BlockPrepull(mage.GetOrRegisterAura(core.Aura{
	// 	Label:     "Arcane Charges Aura",
	// 	ActionID:  core.ActionID{SpellID: 36032},
	// 	Duration:  time.Second * 8,
	// 	MaxStacks: 3,
	// 	OnGain: func(aura *core.Aura, sim *core.Simulation) {
	// 		abCastMod.Activate()
	// 		abCostMod.Activate()
	// 	},
	// 	OnExpire: func(aura *core.Aura, sim *core.Simulation) {
	// 		abCastMod.Deactivate()
	// 		abCostMod.Deactivate()
	// 	},
	// 	OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
	// 		stacks := float64(newStacks)
	// 		abCastMod.UpdateTimeValue(castTimeReduction * time.Duration(newStacks))
	// 		abCostMod.UpdateFloatValue(powerCostIncrease * stacks)
	// 	},
	// }))
}
