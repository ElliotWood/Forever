package paladin

import (
	"github.com/wowsims/forever/sim/core"
)

func (paladin *Paladin) registerAuras() {
	paladin.registerDevotionAura()
	paladin.registerRetributionAura()
	paladin.registerConcentrationAura()
	paladin.registerFireResistanceAura()
	paladin.registerFrostResistanceAura()
	paladin.registerShadowResistanceAura()
}

// TODO: To be implemented. TBC body below needs no porting; kept commented until this class's port is reviewed.
//
// registerAuraSpell wires up a castable paladin aura spell that toggles the
// given self-cast aura (which must already be in the PaladinAuraCategory).
func (paladin *Paladin) registerAuraSpell(aura *core.Aura, classSpellMask int64) {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// paladin.RegisterSpell(core.SpellConfig{
	// 	ActionID:       aura.ActionID,
	// 	SpellSchool:    core.SpellSchoolHoly,
	// 	DefenseType:    core.DefenseTypeMagic,
	// 	ProcMask:       core.ProcMaskEmpty,
	// 	Flags:          core.SpellFlagAPL | core.SpellFlagHelpful,
	// 	ClassSpellMask: classSpellMask,
	//
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			GCD: core.GCDDefault,
	// 		},
	// 	},
	//
	// 	ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
	// 		aura.Activate(sim)
	// 	},
	// })
}

// TODO: To be implemented. TBC body below already accounts for Forever dropping Improved Devotion Aura (untalented, per the TODO inside); kept commented until this class's port is reviewed.
//
// Devotion Aura
// https://www.wowhead.com/forever/spell=27149
//
// Gives 861 additional armor to party members within 30 yards.
// Improved Devotion Aura talent increases the armor bonus by up to 40%.
// Players may only have one Aura on them per Paladin at any one time.
func (paladin *Paladin) registerDevotionAura() {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// // TODO: Forever drops Improved Devotion Aura; untalented (0 points) until we know
	// // whether the effect moved onto another talent.
	// aura := core.DevotionAuraBuff(&paladin.Character, true, 0)
	// paladin.registerAuraSpell(aura, SpellMaskDevotionAura)
}

// TODO: To be implemented. TBC body below already accounts for Forever dropping Improved Retribution Aura (untalented, per the TODO inside); kept commented until this class's port is reviewed.
//
// Retribution Aura
// https://www.wowhead.com/forever/spell=27150
//
// Causes 26 Holy damage to any creature that strikes a party member within 30 yards.
// Improved Retribution Aura talent increases damage by up to 50%.
// Players may only have one Aura on them per Paladin at any one time.
func (paladin *Paladin) registerRetributionAura() {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// // TODO: Forever drops Improved Retribution Aura; untalented (0 points) until we know
	// // whether the effect moved onto another talent.
	// aura := core.RetributionAuraBuff(&paladin.Character, true, 0)
	// paladin.registerAuraSpell(aura, SpellMaskRetributionAura)
}

// TODO: To be implemented. TBC body below needs no porting; kept commented until this class's port is reviewed.
//
// registerSelfCastAura creates a self-cast paladin aura with the standard
// PaladinAuraCategory exclusivity and returns the aura for further configuration.
func (paladin *Paladin) registerSelfCastAura(label string, actionID core.ActionID) *core.Aura {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// aura := paladin.RegisterAura(core.Aura{
	// 	Label:    label + " (Player)" + paladin.Label,
	// 	ActionID: actionID,
	// 	Duration: core.NeverExpires,
	// })
	// aura.NewExclusiveEffect(core.PaladinAuraCategory, true, core.ExclusiveEffect{})
	// return aura
}

// TODO: To be implemented. TBC body below already accounts for Forever dropping Improved Concentration Aura (untalented, per the TODO inside); kept commented until this class's port is reviewed.
//
// Concentration Aura
// https://www.wowhead.com/forever/spell=19746
//
// All party members within 30 yards lose 35% less casting or channeling time
// when damaged. Players may only have one Aura on them per Paladin at any one time.
func (paladin *Paladin) registerConcentrationAura() {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// // TODO: Forever drops Improved Concentration Aura; untalented (0 points) until we
	// // know whether the effect moved onto another talent.
	// aura := core.ConcentrationAura(&paladin.Character, true, 0)
	// paladin.registerAuraSpell(aura, SpellMaskConcentrationAura)
}

// TODO: To be implemented. TBC body below needs no porting; kept commented until this class's port is reviewed.
//
// Fire Resistance Aura
// https://www.wowhead.com/forever/spell=27153
//
// Gives 70 fire resistance to party members within 30 yards.
// Players may only have one Aura on them per Paladin at any one time.
func (paladin *Paladin) registerFireResistanceAura() {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// aura := core.FireResistanceAura(&paladin.Character, true)
	// paladin.registerAuraSpell(aura, SpellMaskFireResistanceAura)
}

// TODO: To be implemented. TBC body below needs no porting; kept commented until this class's port is reviewed.
//
// Frost Resistance Aura
// https://www.wowhead.com/forever/spell=27152
//
// Gives 70 frost resistance to party members within 30 yards.
// Players may only have one Aura on them per Paladin at any one time.
func (paladin *Paladin) registerFrostResistanceAura() {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// aura := core.FrostResistanceAura(&paladin.Character, true)
	// paladin.registerAuraSpell(aura, SpellMaskFrostResistanceAura)
}

// TODO: To be implemented. TBC body below needs no porting; kept commented until this class's port is reviewed.
//
// Shadow Resistance Aura
// https://www.wowhead.com/forever/spell=27151
//
// Gives 70 shadow resistance to party members within 30 yards.
// Players may only have one Aura on them per Paladin at any one time.
func (paladin *Paladin) registerShadowResistanceAura() {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// aura := core.ShadowResistanceAura(&paladin.Character, true)
	// paladin.registerAuraSpell(aura, SpellMaskShadowResistanceAura)
}
