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

// TODO: To be implemented. The body below is the one the port needs; kept commented until this class's port is reviewed.
//
// Devotion Aura
// https://www.wowhead.com/forever/spell=10293
//
// Gives 735 additional armor to party members within 30 yards.
// Players may only have one Aura on them per Paladin at any one time.
func (paladin *Paladin) registerDevotionAura() {
	panic("To be implemented")

	// The implementation, kept for the port:
	// aura := core.DevotionAuraAura(&paladin.Unit, true, 0)
	// paladin.registerAuraSpell(aura, SpellMaskDevotionAura)
}

// TODO: To be implemented. The body below needs the spell around the aura;
// kept commented until this class's port is reviewed.
//
// Retribution Aura
// https://www.wowhead.com/forever/spell=10301
//
// Causes 30 Holy damage to any creature that strikes a party member.
// Players may only have one Aura on them per Paladin at any one time.
func (paladin *Paladin) registerRetributionAura() {
	panic("To be implemented")

	// The body the port needs:
	// aura := core.RetributionAuraAura(&paladin.Unit, true, 0)
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

	// The implementation, kept for the port:
	// aura := core.ConcentrationAuraAura(&paladin.Unit, true, 0)
	// paladin.registerAuraSpell(aura, SpellMaskConcentrationAura)
}

// TODO: To be implemented. The body below is the one the port needs; kept commented until this class's port is reviewed.
//
// Fire Resistance Aura
// https://www.wowhead.com/forever/spell=19900
//
// Gives 60 fire resistance to party members within 30 yards.
// Players may only have one Aura on them per Paladin at any one time.
func (paladin *Paladin) registerFireResistanceAura() {
	panic("To be implemented")

	// The implementation, kept for the port:
	// aura := core.FireResistanceAuraAura(&paladin.Unit, true, 0)
	// paladin.registerAuraSpell(aura, SpellMaskFireResistanceAura)
}

// TODO: To be implemented. The body below is the one the port needs; kept commented until this class's port is reviewed.
//
// Frost Resistance Aura
// https://www.wowhead.com/forever/spell=19898
//
// Gives 60 frost resistance to party members within 30 yards.
// Players may only have one Aura on them per Paladin at any one time.
func (paladin *Paladin) registerFrostResistanceAura() {
	panic("To be implemented")

	// The implementation, kept for the port:
	// aura := core.FrostResistanceAuraAura(&paladin.Unit, true, 0)
	// paladin.registerAuraSpell(aura, SpellMaskFrostResistanceAura)
}

// TODO: To be implemented. The body below is the one the port needs; kept commented until this class's port is reviewed.
//
// Shadow Resistance Aura
// https://www.wowhead.com/forever/spell=19896
//
// Gives 60 shadow resistance to party members within 30 yards.
// Players may only have one Aura on them per Paladin at any one time.
func (paladin *Paladin) registerShadowResistanceAura() {
	panic("To be implemented")

	// The implementation, kept for the port:
	// aura := core.ShadowResistanceAuraAura(&paladin.Unit, true, 0)
	// paladin.registerAuraSpell(aura, SpellMaskShadowResistanceAura)
}

// Sanctity Aura (Talent)
// https://www.wowhead.com/forever/spell=20218
//
// Increases Holy damage done by party members within 30 yards by 10%.
// Players may only have one Aura on them per Paladin at any one time.
//
// TODO: uncalled -- no SkillLineAbility row grants spell 20218 and paladin tree 1100
// has no node for it, so sim/core builds no aura for it either. Re-gate before wiring
// back into registerTalentSpells.
// TODO: To be implemented once the client has a spell to read the aura from.
func (paladin *Paladin) registerSanctityAura() {
	panic("To be implemented")
}
