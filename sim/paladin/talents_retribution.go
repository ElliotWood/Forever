package paladin

import (
	"slices"
	"time"

	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

func (paladin *Paladin) registerRetributionTalents() {
	// Tier 1
	paladin.applyDeflection()
	paladin.applyBenediction()

	// Tier 2
	paladin.applyImprovedJudgement()
	paladin.applyHolyConduit()
	paladin.applyConviction()

	// Tier 3
	paladin.applyVindication()
	paladin.applySanctifiedJudgement()
	// Seal of Command registered in registerTalentSpells
	paladin.applyPursuitOfJustice()

	// Tier 4
	paladin.applyEyeForAnEye()
	paladin.applySacredArbiter()
	paladin.applyCrusade()

	// Tier 5
	paladin.applyTwoHandedWeaponSpecialization()
	paladin.applyVengeance()
	// Repentance not implemented

	// Tier 6
	paladin.applyChampionOfTheLight()
	paladin.applyInstrumentOfLaw()

	// Tier 7
	paladin.applyTwistOfLight()
}

// Benediction - Reduces the mana cost of your Judgement and Seal spells by 3/6/9/12/15%
func (paladin *Paladin) applyBenediction() {
	if paladin.Talents.Benediction == 0 {
		return
	}

	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskAllSeals | SpellMaskJudgement,
		Kind:       core.SpellMod_PowerCost_Pct_Add,
		FloatValue: -.03 * float64(paladin.Talents.Benediction),
	})
}

// Improved Judgement - Decreases the cooldown of your Judgement spell by 1/2 sec
func (paladin *Paladin) applyImprovedJudgement() {
	if paladin.Talents.ImprovedJudgement == 0 {
		return
	}

	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskJudgement,
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: -time.Second * time.Duration(paladin.Talents.ImprovedJudgement),
	})
}

// Deflection - Increases your Parry chance by 1/2/3/4/5%
func (paladin *Paladin) applyDeflection() {
	if paladin.Talents.Deflection == 0 {
		return
	}

	paladin.PseudoStats.BaseParryChance += spellData.Deflection.FractionAt(paladin.Talents.Deflection)
}

// Conviction - Increases your chance to get a critical strike with melee attacks by 1/2/3/4/5%
func (paladin *Paladin) applyConviction() {
	if paladin.Talents.Conviction == 0 {
		return
	}

	paladin.AddStat(stats.PhysicalCritPercent, float64(paladin.Talents.Conviction))
}

// Crusade - Increases all damage caused by 1/2/3% against Humanoids, Demons, Undead and Elementals
func (paladin *Paladin) applyCrusade() {
	if paladin.Talents.Crusade == 0 {
		return
	}

	paladin.Env.RegisterPostFinalizeEffect(func() {
		for _, at := range paladin.AttackTables {
			if slices.Contains([]proto.MobType{proto.MobType_MobTypeDemon, proto.MobType_MobTypeHumanoid, proto.MobType_MobTypeUndead, proto.MobType_MobTypeElemental}, at.Defender.MobType) {
				// Misc 36 is the creature-type mask the client states the bonus against; the
				// mob list above is the sim's own reading of it.
				at.DamageDealtMultiplier *= spellData.Crusade.Effect(shared.A_MOD_DAMAGE_DONE_VERSUS, 36).MultiplierAt(paladin.Talents.Crusade)
			}
		}
	})
}

// Two-Handed Weapon Specialization - Increases the damage you deal with two-handed melee weapons by 2/4/6%
func (paladin *Paladin) applyTwoHandedWeaponSpecialization() {
	if paladin.Talents.TwoHandedWeaponSpecialization == 0 {
		return
	}

	weaponMod := paladin.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ProcMask:   core.ProcMaskMelee,
		FloatValue: spellData.TwoHandedWeaponSpecialization.FractionAt(paladin.Talents.TwoHandedWeaponSpecialization),
	})

	if paladin.GetMainHandType() == proto.HandType_HandTypeTwoHand {
		weaponMod.Activate()
	}

	paladin.RegisterItemSwapCallback(core.AllMeleeWeaponSlots(), func(sim *core.Simulation, slot proto.ItemSlot) {
		if paladin.GetMainHandType() == proto.HandType_HandTypeTwoHand {
			weaponMod.Activate()
		} else {
			weaponMod.Deactivate()
		}
	})
}

// Vengeance - Gives you a 1/2/3/4/5% bonus to Physical and Holy damage you deal for 30 sec after dealing a critical strike from a weapon swing, spell, or ability
func (paladin *Paladin) applyVengeance() {
	if paladin.Talents.Vengeance == 0 {
		return
	}

	bonusMod := .01 * float64(paladin.Talents.Vengeance)

	dmgMod := paladin.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: bonusMod,
		School:     core.SpellSchoolHoly | core.SpellSchoolPhysical,
	})

	aura := paladin.RegisterAura(core.Aura{
		Label:     "Vengeance",
		ActionID:  core.ActionID{SpellID: 20055},
		Duration:  time.Second * 30,
		MaxStacks: 3,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			dmgMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dmgMod.Deactivate()
		},
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			stacks := float64(newStacks)
			dmgMod.UpdateFloatValue(bonusMod * stacks)
		},
	})

	paladin.MakeProcTriggerAura(core.ProcTrigger{
		Name:             "Vengeance - Trigger",
		Callback:         core.CallbackOnSpellHitDealt,
		Outcome:          core.OutcomeCrit,
		ProcChance:       1,
		CanProcFromProcs: true, // 20049/20056/20057 carry the bit: Seal of Blood crits count.
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			aura.Activate(sim)
			aura.AddStack(sim)
		},
	})
}

// Sanctified Judgement - Gives your Judgement spell a 33/66/100% chance to return 80% of the mana cost of the Judged Seal
func (paladin *Paladin) applySanctifiedJudgement() {
	if paladin.Talents.SanctifiedJudgement == 0 {
		return
	}

	procChance := []float64{0, 0.33, 0.66, 1}[paladin.Talents.SanctifiedJudgement]
	sancJudgementManaMetric := paladin.NewManaMetrics(core.ActionID{SpellID: 31930})

	paladin.MakeProcTriggerAura(core.ProcTrigger{
		Name:               "Sanctified Judgement - Trigger",
		ClassSpellMask:     SpellMaskAllJudgements,
		Callback:           core.CallbackOnSpellHitDealt,
		ProcChance:         procChance,
		TriggerImmediately: true,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if paladin.PreviousSeal.IsActive() {
				paladin.AddMana(sim, paladin.PreviousSealSpell.CurCast.Cost*.8, sancJudgementManaMetric)
			} else {
				paladin.AddMana(sim, paladin.CurrentSealSpell.CurCast.Cost*.8, sancJudgementManaMetric)
			}
		},
	})
}

// applyChampionOfTheLight implements Champion of the Light, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (paladin *Paladin) applyChampionOfTheLight() {
	if paladin.Talents.ChampionOfTheLight == 0 {
		return
	}
}

// applyEyeForAnEye implements Eye for an Eye, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (paladin *Paladin) applyEyeForAnEye() {
	if paladin.Talents.EyeForAnEye == 0 {
		return
	}
}

// applyHolyConduit implements Holy Conduit, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (paladin *Paladin) applyHolyConduit() {
	if paladin.Talents.HolyConduit == 0 {
		return
	}
}

// applyInstrumentOfLaw implements Instrument of Law, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (paladin *Paladin) applyInstrumentOfLaw() {
	if paladin.Talents.InstrumentOfLaw == 0 {
		return
	}
}

// applyPursuitOfJustice implements Pursuit of Justice, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (paladin *Paladin) applyPursuitOfJustice() {
	if paladin.Talents.PursuitOfJustice == 0 {
		return
	}
}

// applySacredArbiter implements Sacred Arbiter, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (paladin *Paladin) applySacredArbiter() {
	if !paladin.Talents.SacredArbiter {
		return
	}
}

// applyTwistOfLight implements Twist of Light, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (paladin *Paladin) applyTwistOfLight() {
	if !paladin.Talents.TwistOfLight {
		return
	}
}

// applyVindication implements Vindication, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (paladin *Paladin) applyVindication() {
	if paladin.Talents.Vindication == 0 {
		return
	}
}
