package shaman

import (
	"time"

	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/stats"
)

func (shaman *Shaman) registerElementalTalents() {
	// Tier 1
	shaman.applyConvection()
	shaman.applyConcussion()

	// Tier 2
	shaman.applyElementalWarding()
	shaman.applyReverberation()
	shaman.applyCallOfFlame()
	shaman.applyElementalDevastation()

	// Tier 3
	shaman.applyElementalFocus()
	shaman.applyElementalFury()

	// Tier 4
	shaman.applyImprovedFireNova()
	shaman.applyEyeOfTheStorm()
	shaman.applyCallOfThunder()

	// Tier 5
	shaman.applyElementalReach()
	shaman.applyLightningOverload()
	shaman.applyEarthbound()

	// Tier 6
	shaman.applyElementalAlacrity()

	// Tier 7
	shaman.applyLavaBurst()
}

func (shaman *Shaman) applyCallOfFlame() {
	if shaman.Talents.CallOfFlame == 0 {
		return
	}

	shaman.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: spellData.CallOfFlame.Effect(shared.A_ADD_PCT_MODIFIER, shared.SPELLMOD_DAMAGE).FractionAt(shaman.Talents.CallOfFlame),
		ClassMask:  SpellMaskFireTotem | SpellMaskFlameShock | SpellMaskFireNova | SpellMaskLavaBurst,
	})
}

func (shaman *Shaman) applyCallOfThunder() {
	if !shaman.Talents.CallOfThunder {
		return
	}

	shaman.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_BonusCrit_Percent,
		FloatValue: spellData.CallOfThunder.Effect(shared.A_ADD_FLAT_MODIFIER, shared.SPELLMOD_CRITICAL_CHANCE).ValueAt(1),
		ClassMask:  SpellMaskLightningBolt | SpellMaskChainLightning | SpellMaskOverload,
	})
}

func (shaman *Shaman) applyConcussion() {
	if shaman.Talents.Concussion == 0 {
		return
	}

	shaman.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: spellData.Concussion.Effect(shared.A_ADD_PCT_MODIFIER, shared.SPELLMOD_DAMAGE).FractionAt(shaman.Talents.Concussion),
		ClassMask:  SpellMaskLightningBolt | SpellMaskChainLightning | SpellMaskOverload | SpellMaskShock,
	})
}

func (shaman *Shaman) applyConvection() {
	if shaman.Talents.Convection == 0 {
		return
	}

	shaman.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_PowerCost_Pct_Add,
		FloatValue: spellData.Convection.Effect(shared.A_ADD_PCT_MODIFIER, shared.SPELLMOD_COST).FractionAt(shaman.Talents.Convection),
		ClassMask:  SpellMaskLightningBolt | SpellMaskChainLightning | SpellMaskOverload | SpellMaskShock | SpellMaskLavaBurst,
	})
}

func (shaman *Shaman) applyElementalDevastation() {
	if shaman.Talents.ElementalDevastation == 0 {
		return
	}

	critBuffAura := shaman.RegisterAura(core.Aura{
		Label:    "Elemental Devastation",
		ActionID: core.ActionID{SpellID: 29178},
		Duration: time.Second * 10,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_BonusCrit_Percent,
		FloatValue: spellData.ElementalDevastation.Effect(shared.A_DUMMY, 0).ValueAt(shaman.Talents.ElementalDevastation),
		ProcMask:   core.ProcMaskMelee,
	})
	shaman.MakeProcTriggerAura(core.ProcTrigger{
		Name:             "Elemental Devastation Trigger",
		CanProcFromProcs: true, // 29179/29180/30160 carry the bit.
		Callback:         core.CallbackOnSpellHitDealt,
		ProcMask:         core.ProcMaskSpellDamage,
		Outcome:          core.OutcomeCrit,
		Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
			critBuffAura.Activate(sim)
		},
	})
}

func (shaman *Shaman) applyElementalFocus() {
	if !shaman.Talents.ElementalFocus {
		return
	}

	var triggeringSpell *core.Spell
	var triggerTime time.Duration

	canConsumeSpells := SpellMaskLightningBolt | SpellMaskChainLightning | SpellMaskLavaBurst | (SpellMaskShock & ^SpellMaskFlameShockDot)

	clearcasting := spellData.ElementalFocusTriggered.HighestRank()
	maxStacks := clearcasting.ProcCharges

	clearcastingAura := shaman.RegisterAura(core.Aura{
		Label:     "Clearcasting",
		ActionID:  core.ActionID{SpellID: 16246},
		Duration:  clearcasting.Duration,
		MaxStacks: maxStacks,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if !spell.Matches(canConsumeSpells) {
				return
			}
			if spell == triggeringSpell && sim.CurrentTime == triggerTime {
				return
			}
			aura.RemoveStack(sim)
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_PowerCost_Pct_Add,
		ClassMask:  canConsumeSpells,
		FloatValue: clearcasting.Effect(shared.A_ADD_PCT_MODIFIER, shared.SPELLMOD_COST).Value / 100,
	})

	shaman.MakeProcTriggerAura(core.ProcTrigger{
		Name:               "Elemental Focus",
		ProcChance:         spellData.ElementalFocus.ProcChanceAt(1),
		Callback:           core.CallbackOnSpellHitDealt,
		CanProcFromProcs:   true, // 16164 carries the bit: Lightning Overload crits count.
		Outcome:            core.OutcomeLanded,
		TriggerImmediately: true,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.SpellSchool.Matches(core.SpellSchoolElemental) {
				return
			}
			triggeringSpell = spell
			triggerTime = sim.CurrentTime
			clearcastingAura.Activate(sim)
			clearcastingAura.SetStacks(sim, maxStacks)
		},
	})
}

func (shaman *Shaman) applyElementalFury() {
	if shaman.Talents.ElementalFury == 0 {
		return
	}

	// The talent's class mask (16089) also covers Flametongue Attack (bit 21) and Frostbrand
	// Attack (bit 24): shamans' Flametongue Weapon hits crit for 2.0x in logs while the same
	// attack granted by Flametongue Totem crits for 1.5x on other players.
	shaman.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_CritMultiplier_Flat,
		FloatValue: spellData.ElementalFury.Effect(shared.A_ADD_PCT_MODIFIER, shared.SPELLMOD_CRIT_DAMAGE_BONUS).FractionAt(shaman.Talents.ElementalFury),
		ClassMask:  SpellMaskFireTotem | SpellMaskFire | SpellMaskNature | SpellMaskFrost | SpellMaskFlametongueWeapon | SpellMaskFrostbrandWeapon,
	})
}
func (shaman *Shaman) applyLightningOverload() {
	if shaman.Talents.LightningOverload == 0 {
		return
	}
	// In shaman.go -> GetOverloadChance()
}

func (shaman *Shaman) applyReverberation() {
	if shaman.Talents.Reverberation == 0 {
		return
	}

	shaman.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: time.Millisecond * time.Duration(spellData.Reverberation.Effect(shared.A_ADD_FLAT_MODIFIER, shared.SPELLMOD_COOLDOWN).ValueAt(shaman.Talents.Reverberation)),
		ClassMask: SpellMaskShock,
	})
}

// applyElementalAlacrity implements Elemental Alacrity, new in Forever: a flat cast-time cut.
func (shaman *Shaman) applyElementalAlacrity() {
	if shaman.Talents.ElementalAlacrity == 0 {
		return
	}

	shaman.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_CastTime_Flat,
		TimeValue: time.Millisecond * time.Duration(spellData.ElementalAlacrity.Effect(shared.A_ADD_FLAT_MODIFIER, shared.SPELLMOD_CASTING_TIME).ValueAt(shaman.Talents.ElementalAlacrity)),
		ClassMask: SpellMaskLightningBolt | SpellMaskChainLightning | SpellMaskLavaBurst,
	})
}

// applyElementalReach is Forever's spell-range talent (28999, +3/6 yards). Range costs a sim nothing,
// so it changes no number here.
func (shaman *Shaman) applyElementalReach() {
	if shaman.Talents.ElementalReach == 0 {
		return
	}
}

// applyElementalWarding implements Elemental Warding, new in Forever: less fire, frost and nature
// damage taken (the client states one modifier over the three-school mask 28).
func (shaman *Shaman) applyElementalWarding() {
	if shaman.Talents.ElementalWarding == 0 {
		return
	}

	multiplier := spellData.ElementalWarding.Effect(shared.A_MOD_DAMAGE_PERCENT_TAKEN, 28).MultiplierAt(shaman.Talents.ElementalWarding)
	for _, school := range []stats.SchoolIndex{stats.SchoolIndexFire, stats.SchoolIndexFrost, stats.SchoolIndexNature} {
		shaman.PseudoStats.SchoolDamageTakenMultiplier[school] *= multiplier
	}
}

// applyEyeOfTheStorm is Forever's pushback-resistance talent (29062, 23/47/70%). The sim models no
// spell pushback, so it changes no number here.
func (shaman *Shaman) applyEyeOfTheStorm() {
	if shaman.Talents.EyeOfTheStorm == 0 {
		return
	}
}

// applyImprovedFireNova implements Improved Fire Nova, new in Forever: more Fire Nova damage and a
// shorter cooldown on it.
func (shaman *Shaman) applyImprovedFireNova() {
	if shaman.Talents.ImprovedFireNova == 0 {
		return
	}

	shaman.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: spellData.ImprovedFireNova.Effect(shared.A_ADD_PCT_MODIFIER, shared.SPELLMOD_DAMAGE).FractionAt(shaman.Talents.ImprovedFireNova),
		ClassMask:  SpellMaskFireNova,
	})
	shaman.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: time.Millisecond * time.Duration(spellData.ImprovedFireNova.Effect(shared.A_ADD_FLAT_MODIFIER, shared.SPELLMOD_COOLDOWN).ValueAt(shaman.Talents.ImprovedFireNova)),
		ClassMask: SpellMaskFireNova,
	})
}

// applyEarthbound implements Earthbound, new in Forever.
//
// TODO: To be implemented. The generated tables carry no Earthbound row, so there is nothing to read
// its effect from yet.
func (shaman *Shaman) applyEarthbound() {
	if !shaman.Talents.Earthbound {
		return
	}
}

// applyLavaBurst grants Lava Burst, new in Forever. See lava_burst.go.
func (shaman *Shaman) applyLavaBurst() {
	shaman.registerLavaBurstSpell()
}
