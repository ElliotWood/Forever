package druid

import (
	"time"

	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/stats"
)

func (druid *Druid) ApplyTalents() {
	// Balance
	druid.applyImprovedMoonfire()
	// Brambles: applied as Thorns aura points in thorns.go
	druid.applyInsectSwarm()
	druid.applyNaturesReach()
	druid.applyVengeance()
	druid.applyNaturesGrace()
	druid.applyMoonglow()
	druid.applyMoonfury()
	// Omen of Clarity: Forever drops the talent; see omen_of_clarity.go

	// Feral
	druid.applyFuror()
	druid.applyFerocity()
	druid.applyFeralInstincts()
	druid.applyPredatoryInstincts()
	druid.applyThickHide()
	druid.applyFeralSwiftness()
	druid.applySharpenedClaws()
	druid.applyShreddingAttacks()
	druid.applyPredatoryStrikes()
	druid.applyPrimalFury()
	druid.applySavageFury()
	druid.applyHeartOfTheWild()

	// Restoration
	druid.applyNaturalist()
	druid.applyNaturalShapeshifter()
	druid.applySubtlety()
	druid.applyLivingSpirit()

	// Forever additions, not yet implemented.
	druid.applyImprovedWrath()
	druid.applyGenesis()
	druid.applyNaturesMajesty()
	druid.applyImprovedEntanglingRoots()
	druid.applyNaturesSplendor()
	druid.applyImprovedStarfire()
	druid.applyOvergrowth()
	druid.applyEclipse()
	druid.applyBrutalImpact()
	druid.applyFeralCharge()
	druid.applyKingOfTheJungle()
	druid.applyNaturalReaction()
	druid.applyRendAndTear()
	druid.applyBerserk()
	druid.applyNaturesFocus()
	druid.applyReflection()
	druid.applyGiftOfNature()
	druid.applyGiftOfTheEarthmother()
	druid.applyTranquilSpirit()
	druid.applyImprovedRejuvenation()
	druid.applySwiftmend()
	druid.applyNaturesSwiftness()
	druid.applyImprovedTranquility()
	druid.applyImprovedRegrowth()
	druid.applyWildGrowth()
}

// applyThickHide increases armor contribution from items by 4/7/10% (ranks 1/2/3).
// Applies to both Armor and BonusArmor equip stats
func (druid *Druid) applyThickHide() {
	if druid.Talents.ThickHide == 0 {
		return
	}

	bonusByRank := [3]float64{0.04, 0.07, 0.10}
	bonus := bonusByRank[druid.Talents.ThickHide-1]
	druid.ApplyEquipScaling(stats.Armor, 1.0+bonus)
	druid.ApplyEquipScaling(stats.BonusArmor, 1.0+bonus)
}

// Predatory Instincts: +2% melee critical strike damage per rank while in Cat or Bear form.
// The client aura (33859 and its ranks) is school-masked to Physical, so it covers every
// physical attack in form (abilities and auto attacks, Ravage included) and nothing else.
//
// TODO: To be implemented. Forever's regenerated aura enum no longer carries the crit damage
// aura this talent's rank data used (A_MOD_CRIT_DAMAGE_BONUS is gone from the auto-generated
// table), and this talent has no other effect to fall back on, so it is fully disabled.
func (druid *Druid) applyPredatoryInstincts() {
	if druid.Talents.PredatoryInstincts == 0 {
		return
	}
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// if druid.Talents.PredatoryInstincts == 0 {
	// 	return
	// }
	//
	// druid.AddStaticMod(core.SpellModConfig{
	// 	Kind:       core.SpellMod_CritMultiplier_Pct,
	// 	School:     core.SpellSchoolPhysical,
	// 	FloatValue: spellData.PredatoryInstincts.Effect(shared.A_MOD_CRIT_DAMAGE_BONUS, 1).FractionAt(druid.Talents.PredatoryInstincts),
	// })
}

func (druid *Druid) applyMoonfury() {
	if druid.Talents.Moonfury == 0 {
		return
	}

	// Forever states Moonfury as +2% damage per rank to the Arcane|Nature schools (mask 72)
	// rather than as a spell modifier; the class mask keeps it on the TBC spell list.
	druid.AddStaticMod(core.SpellModConfig{
		ClassMask:  DruidSpellWrath | DruidSpellStarfire | DruidSpellMoonfire,
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: spellData.Moonfury.Effect(shared.A_MOD_DAMAGE_PERCENT_DONE, 72).FractionAt(druid.Talents.Moonfury),
	})
}

func (druid *Druid) applyMoonglow() {
	if druid.Talents.Moonglow == 0 {
		return
	}

	druid.AddStaticMod(core.SpellModConfig{
		ClassMask:  DruidSpellMoonfire | DruidSpellStarfire | DruidSpellWrath | DruidSpellHealingTouch | DruidSpellRegrowth | DruidSpellRejuvenation,
		FloatValue: -0.03 * float64(druid.Talents.Moonglow),
		Kind:       core.SpellMod_PowerCost_Pct_Add,
	})
}

func (druid *Druid) applyNaturesGrace() {
	if !druid.Talents.NaturesGrace {
		return
	}

	lastProcAt := time.Duration(-1)

	aura := druid.RegisterAura(core.Aura{
		Label:    "Nature's Grace",
		ActionID: core.ActionID{SpellID: 16886},
		Duration: time.Second * 15,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			lastProcAt = -1
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.CurCast.CastTime == 0 {
				return
			}
			// Only consume if the aura was already active when this cast started;
			// a cast in flight when the proc landed did not benefit from it.
			if aura.TimeActive(sim) < spell.CurCast.CastTime {
				return
			}
			// A proc that landed during this cast re-arms the buff for the next
			// cast instead of being consumed by this one.
			if lastProcAt > sim.CurrentTime-spell.CurCast.CastTime {
				return
			}

			aura.Deactivate(sim)
		},
	}).AttachSpellMod(core.SpellModConfig{
		ClassMask: DruidSpellStarfire | DruidSpellWrath,
		Kind:      core.SpellMod_CastTime_Flat,
		TimeValue: time.Millisecond * -500,
	})

	druid.MakeProcTriggerAura(core.ProcTrigger{
		Name:               "Nature's Grace Trigger",
		Callback:           core.CallbackOnSpellHitDealt,
		ClassSpellMask:     DruidSpellWrath | DruidSpellStarfire,
		Outcome:            core.OutcomeCrit,
		TriggerImmediately: true,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			lastProcAt = sim.CurrentTime
			aura.Activate(sim)
		},
	})
}

func (druid *Druid) applyVengeance() {
	if druid.Talents.Vengeance == 0 {
		return
	}

	druid.AddStaticMod(core.SpellModConfig{
		ClassMask:  DruidSpellWrath | DruidSpellStarfire | DruidSpellMoonfire,
		Kind:       core.SpellMod_CritMultiplier_Flat,
		FloatValue: spellData.Vengeance.FractionAt(druid.Talents.Vengeance),
	})
}

func (druid *Druid) applyNaturesReach() {
	if druid.Talents.NaturesReach == 0 {
		return
	}

	// druid.AddStaticMod(core.SpellModConfig{
	// 	ClassMask:  DruidSpellsBalance | DruidSpellFearieFireFeral,
	// 	Kind:       ****BONUS RANGE**** most likely irrelevant for sim
	// 	FloatValue: 10.0 * float64(druid.Talents.NaturesReach),
	// })
}

func (druid *Druid) applyInsectSwarm() {
	if !druid.Talents.InsectSwarm {
		return
	}

	druid.registerInsectSwarmSpell()
}

func (druid *Druid) applyImprovedMoonfire() {
	if druid.Talents.ImprovedMoonfire == 0 {
		return
	}

	// 5% per point damage increase to Moonfire and its DoT
	druid.AddStaticMod(core.SpellModConfig{
		ClassMask:  DruidSpellMoonfire,
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: spellData.ImprovedMoonfire.Effect(shared.A_ADD_PCT_MODIFIER, shared.SPELLMOD_DAMAGE).FractionAt(druid.Talents.ImprovedMoonfire),
	})

	// 5% per point chance to crit with Moonfire
	druid.AddStaticMod(core.SpellModConfig{
		ClassMask:  DruidSpellMoonfire,
		Kind:       core.SpellMod_BonusCrit_Percent,
		FloatValue: spellData.ImprovedMoonfire.Effect(shared.A_ADD_FLAT_MODIFIER, shared.SPELLMOD_CRITICAL_CHANCE).ValueAt(druid.Talents.ImprovedMoonfire),
	})
}

func (druid *Druid) applyNaturalShapeshifter() {
	if druid.Talents.NaturalShapeshifter == 0 {
		return
	}

	druid.AddStaticMod(core.SpellModConfig{
		ClassMask:  DruidSpellCatForm | DruidSpellBearForm,
		Kind:       core.SpellMod_PowerCost_Pct_Add,
		FloatValue: -0.1 * float64(druid.Talents.NaturalShapeshifter),
	})
}

func (druid *Druid) applyNaturalist() {
	if druid.Talents.Naturalist == 0 {
		return
	}

	// Forever states the damage bonus against every school (mask 127) instead of physical
	// only; the sim keeps it on physical, which is all a feral druid deals.
	druid.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= spellData.Naturalist.Effect(shared.A_MOD_DAMAGE_PERCENT_DONE, 127).MultiplierAt(druid.Talents.Naturalist)
}

func (druid *Druid) applyHeartOfTheWild() {
	if druid.Talents.HeartOfTheWild == 0 {
		return
	}

	// +2% Intellect per rank (all forms, always active). The stat aura carries misc 0 in the
	// Forever data, as every stat-percent effect does, so the stat is the call site's choice.
	// The Cat/Bear form-specific bonuses are handled dynamically in RegisterCatFormAura /
	// RegisterBearFormAura.
	druid.MultiplyStat(stats.Intellect, spellData.HeartOfTheWild.Effect(shared.A_MOD_TOTAL_STAT_PERCENTAGE, 0).MultiplierAt(druid.Talents.HeartOfTheWild))
}

func (druid *Druid) applySharpenedClaws() {
	if druid.Talents.SharpenedClaws == 0 {
		return
	}
	bonus := stats.Stats{
		stats.PhysicalCritPercent: float64(druid.Talents.SharpenedClaws) * 2.0,
	}
	druid.CatFormAura.AttachStatsBuff(bonus)
	druid.BearFormAura.AttachStatsBuff(bonus)
}

func (druid *Druid) applyFeralSwiftness() {
	if druid.Talents.FeralSwiftness == 0 {
		return
	}
	bonus := stats.Stats{
		stats.DodgeRating: core.DodgeRatingPerDodgePercent * 2.0 * float64(druid.Talents.FeralSwiftness),
	}
	druid.CatFormAura.AttachStatsBuff(bonus)
	druid.BearFormAura.AttachStatsBuff(bonus)
}

func (druid *Druid) applyPredatoryStrikes() {
	if druid.Talents.PredatoryStrikes == 0 {
		return
	}
	bonus := stats.Stats{
		stats.AttackPower: float64(druid.Talents.PredatoryStrikes) * 0.5 * core.CharacterLevel,
	}
	druid.CatFormAura.AttachStatsBuff(bonus)
	druid.BearFormAura.AttachStatsBuff(bonus)
}

// applyFuror gives a 20% chance per rank to gain 40 energy when shifting into
// Cat Form, or 10 rage when shifting into Bear Form.
func (druid *Druid) applyFuror() {
	if druid.Talents.Furor == 0 {
		return
	}

	// Both dummy effects carry the same ladder, one per form, so either answers the chance.
	druid.FurorProcChance = spellData.Furor.EffectAt(0).FractionAt(druid.Talents.Furor)
}

func (druid *Druid) applyFerocity() {
	if druid.Talents.Ferocity == 0 {
		return
	}

	druid.AddStaticMod(core.SpellModConfig{
		ClassMask: DruidSpellRake,
		Kind:      core.SpellMod_PowerCost_Flat,
		IntValue:  -druid.Talents.Ferocity,
	})

	druid.AddStaticMod(core.SpellModConfig{
		ClassMask: DruidSpellMangleBear | DruidSpellMaul | DruidSpellSwipe,
		Kind:      core.SpellMod_PowerCost_Flat,
		IntValue:  -druid.Talents.Ferocity,
	})
}

func (druid *Druid) applySavageFury() {
	if druid.Talents.SavageFury == 0 {
		return
	}

	druid.AddStaticMod(core.SpellModConfig{
		ClassMask:  DruidSpellRake,
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: spellData.SavageFury.Effect(shared.A_ADD_PCT_MODIFIER, shared.SPELLMOD_DAMAGE).FractionAt(druid.Talents.SavageFury),
	})
}

func (druid *Druid) applyShreddingAttacks() {
	if druid.Talents.ShreddingAttacks == 0 {
		return
	}

	druid.AddStaticMod(core.SpellModConfig{
		ClassMask: DruidSpellShred,
		Kind:      core.SpellMod_PowerCost_Flat,
		IntValue:  -9 * druid.Talents.ShreddingAttacks,
	})

	druid.AddStaticMod(core.SpellModConfig{
		ClassMask: DruidSpellLacerate,
		Kind:      core.SpellMod_PowerCost_Flat,
		IntValue:  -druid.Talents.ShreddingAttacks,
	})
}

func (druid *Druid) applyPrimalFury() {
	if druid.Talents.PrimalFury == 0 {
		return
	}

	procChance := 0.5 * float64(druid.Talents.PrimalFury)
	actionID := core.ActionID{SpellID: 37117}
	rageMetrics := druid.NewRageMetrics(actionID)
	cpMetrics := druid.NewComboPointMetrics(actionID)

	// Cat form: +1 combo point on builder crits.
	druid.MakeProcTriggerAura(core.ProcTrigger{
		Name:           "Primal Fury (Cat)",
		ActionID:       actionID,
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: DruidSpellBuilder,
		Outcome:        core.OutcomeCrit,
		ProcChance:     procChance,
		ExtraCondition: func(_ *core.Simulation, _ *core.Spell, _ *core.SpellResult) bool {
			return druid.InForm(Cat)
		},
		Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
			druid.AddComboPoints(sim, 1, cpMetrics)
		},
	})

	// Bear form: +5 rage on auto-attack crits.
	druid.MakeProcTriggerAura(core.ProcTrigger{
		Name:       "Primal Fury (Bear)",
		ActionID:   actionID,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeCrit,
		ProcChance: procChance,
		ExtraCondition: func(_ *core.Simulation, _ *core.Spell, _ *core.SpellResult) bool {
			return druid.InForm(Bear)
		},
		Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
			druid.AddRage(sim, 5, rageMetrics)
		},
	})
}

func (druid *Druid) applyFeralInstincts() {
	if druid.Talents.FeralInstinct == 0 {
		return
	}

	// TODO: Forever repurposes Feral Instinct: the Dire Bear threat bonus is gone and the
	// spell now carries SPELLMOD_EFFECT1 +5/10/15 and SPELLMOD_DAMAGE +10/20/30% against an
	// unknown spell, so the threat multiplier is pinned to the untalented 1.0.
	threatMultiplier := 1.0
	druid.BearFormAura.AttachMultiplicativePseudoStatBuff(&druid.PseudoStats.ThreatMultiplier, threatMultiplier)
}

func (druid *Druid) applySubtlety() {
	if druid.Talents.Subtlety == 0 {
		return
	}

	// Reduces the threat caused by healing/damage spells by 4/8/12/16/20% per rank.
	druid.AddStaticMod(core.SpellModConfig{
		ClassMask:  DruidHealingSpells | DruidDamagingSpells,
		Kind:       core.SpellMod_ThreatMultiplier_Pct,
		FloatValue: -0.04 * float64(druid.Talents.Subtlety),
	})
}

func (druid *Druid) applyLivingSpirit() {
	if druid.Talents.LivingSpirit == 0 {
		return
	}

	// Increases total Spirit by 5/10/15% per rank.
	druid.MultiplyStat(stats.Spirit, spellData.LivingSpirit.MultiplierAt(druid.Talents.LivingSpirit))
}

// applyImprovedWrath implements Improved Wrath, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (druid *Druid) applyImprovedWrath() {
	if druid.Talents.ImprovedWrath == 0 {
		return
	}

	panic("To be implemented")
}

// applyGenesis implements Genesis, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (druid *Druid) applyGenesis() {
	if druid.Talents.Genesis == 0 {
		return
	}

	panic("To be implemented")
}

// applyNaturesMajesty implements Nature's Majesty, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (druid *Druid) applyNaturesMajesty() {
	if druid.Talents.NaturesMajesty == 0 {
		return
	}

	panic("To be implemented")
}

// applyImprovedEntanglingRoots implements Improved Entangling Roots, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (druid *Druid) applyImprovedEntanglingRoots() {
	if druid.Talents.ImprovedEntanglingRoots == 0 {
		return
	}

	panic("To be implemented")
}

// applyNaturesSplendor implements Nature's Splendor, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (druid *Druid) applyNaturesSplendor() {
	if !druid.Talents.NaturesSplendor {
		return
	}

	panic("To be implemented")
}

// applyImprovedStarfire implements Improved Starfire, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (druid *Druid) applyImprovedStarfire() {
	if druid.Talents.ImprovedStarfire == 0 {
		return
	}

	panic("To be implemented")
}

// applyOvergrowth implements Overgrowth, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (druid *Druid) applyOvergrowth() {
	if druid.Talents.Overgrowth == 0 {
		return
	}

	panic("To be implemented")
}

// applyEclipse implements Eclipse, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (druid *Druid) applyEclipse() {
	if druid.Talents.Eclipse == 0 {
		return
	}

	panic("To be implemented")
}

// applyBrutalImpact implements Brutal Impact, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (druid *Druid) applyBrutalImpact() {
	if druid.Talents.BrutalImpact == 0 {
		return
	}

	panic("To be implemented")
}

// applyFeralCharge implements Feral Charge, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (druid *Druid) applyFeralCharge() {
	if !druid.Talents.FeralCharge {
		return
	}

	panic("To be implemented")
}

// applyKingOfTheJungle implements King of the Jungle, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (druid *Druid) applyKingOfTheJungle() {
	if druid.Talents.KingOfTheJungle == 0 {
		return
	}

	panic("To be implemented")
}

// applyNaturalReaction implements Natural Reaction, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (druid *Druid) applyNaturalReaction() {
	if druid.Talents.NaturalReaction == 0 {
		return
	}

	panic("To be implemented")
}

// applyRendAndTear implements Rend and Tear, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (druid *Druid) applyRendAndTear() {
	if druid.Talents.RendAndTear == 0 {
		return
	}

	panic("To be implemented")
}

// applyBerserk implements Berserk, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (druid *Druid) applyBerserk() {
	if !druid.Talents.Berserk {
		return
	}

	panic("To be implemented")
}

// applyNaturesFocus implements Nature's Focus, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (druid *Druid) applyNaturesFocus() {
	if druid.Talents.NaturesFocus == 0 {
		return
	}

	panic("To be implemented")
}

// applyReflection implements Reflection, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (druid *Druid) applyReflection() {
	if druid.Talents.Reflection == 0 {
		return
	}

	panic("To be implemented")
}

// applyGiftOfNature implements Gift of Nature, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (druid *Druid) applyGiftOfNature() {
	if druid.Talents.GiftOfNature == 0 {
		return
	}

	panic("To be implemented")
}

// applyGiftOfTheEarthmother implements Gift of the Earthmother, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (druid *Druid) applyGiftOfTheEarthmother() {
	if !druid.Talents.GiftOfTheEarthmother {
		return
	}

	panic("To be implemented")
}

// applyTranquilSpirit implements Tranquil Spirit, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (druid *Druid) applyTranquilSpirit() {
	if druid.Talents.TranquilSpirit == 0 {
		return
	}

	panic("To be implemented")
}

// applyImprovedRejuvenation implements Improved Rejuvenation, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (druid *Druid) applyImprovedRejuvenation() {
	if druid.Talents.ImprovedRejuvenation == 0 {
		return
	}

	panic("To be implemented")
}

// applySwiftmend implements Swiftmend, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (druid *Druid) applySwiftmend() {
	if !druid.Talents.Swiftmend {
		return
	}

	panic("To be implemented")
}

// applyNaturesSwiftness implements Nature's Swiftness, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (druid *Druid) applyNaturesSwiftness() {
	if !druid.Talents.NaturesSwiftness {
		return
	}

	panic("To be implemented")
}

// applyImprovedTranquility implements Improved Tranquility, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (druid *Druid) applyImprovedTranquility() {
	if druid.Talents.ImprovedTranquility == 0 {
		return
	}

	panic("To be implemented")
}

// applyImprovedRegrowth implements Improved Regrowth, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (druid *Druid) applyImprovedRegrowth() {
	if druid.Talents.ImprovedRegrowth == 0 {
		return
	}

	panic("To be implemented")
}

// applyWildGrowth implements Wild Growth, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (druid *Druid) applyWildGrowth() {
	if !druid.Talents.WildGrowth {
		return
	}

	panic("To be implemented")
}
