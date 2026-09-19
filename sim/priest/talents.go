package priest

import (
	"time"

	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/stats"
)

func (priest *Priest) applyPowerInfusion() {
	if !priest.Talents.PowerInfusion {
		return
	}

	piAura := core.PowerInfusionAura(priest.GetCharacter(), 0)

	piSpell := priest.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 10060},
		SpellSchool: core.SpellSchoolHoly,
		Flags:       core.SpellFlagHelpful,
		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 16,
		},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: core.PowerInfusionCD,
			},
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, _ *core.Spell) {
			piAura.Activate(sim)
		},
	})

	priest.AddMajorCooldown(core.MajorCooldown{
		Spell:    piSpell,
		Priority: core.CooldownPriorityBloodlust,
		Type:     core.CooldownTypeMana,
	})
}

// Spirit of Redemption's passive half: +5% total Spirit. The on-death form is not modelled.
func (priest *Priest) applySpiritOfRedemption() {
	if !priest.Talents.SpiritOfRedemption {
		return
	}
	priest.MultiplyStat(stats.Spirit, 1.05)
}

func (priest *Priest) applyMentalStrength() {
	if priest.Talents.MentalStrength == 0 {
		return
	}
	// +2% mana per rank
	priest.MultiplyStat(stats.Mana, spellData.MentalStrength.MultiplierAt(priest.Talents.MentalStrength))
}

func (priest *Priest) applySpiritualGuidance() {
	if priest.Talents.SpiritualGuidance == 0 {
		return
	}
	// 5% of Spirit added to spell damage per rank
	coeff := spellData.SpiritualGuidance.Effect(shared.A_MOD_SPELL_DAMAGE_OF_STAT_PERCENT, 126).FractionAt(priest.Talents.SpiritualGuidance)
	priest.AddStatDependency(stats.Spirit, stats.SpellDamage, coeff) // Only scaling damage for now since no healing sim....yet!
}

func (priest *Priest) applyDivineFury() {
	if priest.Talents.DivineFury == 0 {
		return
	}
	// -0.1s per rank
	priest.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_CastTime_Flat,
		TimeValue: time.Millisecond * time.Duration(-100*priest.Talents.DivineFury),
		ClassMask: PriestSpellSmite | PriestSpellHolyFire,
	})
}

func (priest *Priest) applySearingLight() {
	if priest.Talents.SearingLight == 0 {
		return
	}
	// +5% damage per rank
	priest.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: spellData.SearingLight.Effect(shared.A_ADD_PCT_MODIFIER, shared.SPELLMOD_DAMAGE).FractionAt(priest.Talents.SearingLight),
		ClassMask:  PriestSpellSmite | PriestSpellHolyFire,
	})
}

func (priest *Priest) applySilentResolve() {
	if priest.Talents.SilentResolve == 0 {
		return
	}
	// -4% threat per rank for discipline and holy spells
	threatReduction := []float64{0, -0.04, -0.08, -0.12, -0.16, -0.20}[priest.Talents.SilentResolve]
	priest.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_ThreatMultiplier_Pct,
		FloatValue: threatReduction,
		ClassMask:  PriestHolySpells,
	})
}

func (priest *Priest) applyHolyNova() {
	if !priest.Talents.HolyNova {
		return
	}
	HolyNovaRankMap.RegisterAll(priest.registerHolyNovaSpell)
}

func (priest *Priest) applyMindFlay() {
	if !priest.Talents.MindFlay {
		return
	}
	MindFlayRankMap.RegisterAll(priest.registerMindFlaySpell)
}

func (priest *Priest) applyImprovedMindBlast() {
	if priest.Talents.ImprovedMindBlast == 0 {
		return
	}

	priest.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: time.Millisecond * time.Duration(-500*priest.Talents.ImprovedMindBlast),
		ClassMask: PriestSpellMindBlast,
	})
}

func (priest *Priest) applyInnerFocus() {
	if !priest.Talents.InnerFocus {
		return
	}

	critMod := priest.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_BonusCrit_Percent,
		FloatValue: 25.0,
		ClassMask:  PriestSpellsAll,
	})

	var innerFocusSpell *core.Spell
	priest.InnerFocusAura = priest.RegisterAura(core.Aura{
		Label:    "Inner Focus",
		ActionID: core.ActionID{SpellID: 14751},
		Duration: time.Hour,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SpellCostPercentModifier -= 100
			critMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SpellCostPercentModifier += 100
			critMod.Deactivate()
			innerFocusSpell.CD.Use(sim)
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if !spell.Matches(PriestSpellsAll) {
				return
			}
			aura.Deactivate(sim)
		},
	})

	innerFocusSpell = priest.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 14751},
		DefenseType:    core.DefenseTypeMagic,
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,
		ClassSpellMask: PriestSpellFlagNone,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: time.Second * 180,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			priest.InnerFocusAura.Activate(sim)
		},
		RelatedSelfBuff: priest.InnerFocusAura,
	})

	priest.AddMajorCooldown(core.MajorCooldown{
		Spell: innerFocusSpell,
		Type:  core.CooldownTypeMana,
	})
}

func (priest *Priest) applyMeditation() {
	if priest.Talents.Meditation == 0 {
		return
	}

	priest.PseudoStats.SpiritRegenRateCasting += spellData.Meditation.FractionAt(priest.Talents.Meditation)
	priest.UpdateManaRegenRates()
}

func (priest *Priest) applyMentalAgility() {
	if priest.Talents.MentalAgility == 0 {
		return
	}

	priest.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_PowerCost_Pct_Add,
		FloatValue: -0.02 * float64(priest.Talents.MentalAgility),
		ClassMask:  PriestSpellInstant,
	})
}

func (priest *Priest) applyDarkness() {
	if priest.Talents.Darkness == 0 {
		return
	}

	priest.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: spellData.Darkness.Effect(shared.A_ADD_PCT_MODIFIER, shared.SPELLMOD_DOT).FractionAt(priest.Talents.Darkness),
		ClassMask:  PriestShadowSpells,
	})
}

func (priest *Priest) applyShadowFocus() {
	if priest.Talents.ShadowFocus == 0 {
		return
	}

	priest.PseudoStats.SchoolBonusHitChance[stats.SchoolIndexShadow] += spellData.ShadowFocus.ValueAt(priest.Talents.ShadowFocus)

}

func (priest *Priest) applyImprovedShadowWordPain() {
	if priest.Talents.ImprovedShadowWordPain == 0 {
		return
	}

	priest.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_DotNumberOfTicks_Flat,
		IntValue:  int32(priest.Talents.ImprovedShadowWordPain),
		ClassMask: PriestSpellShadowWordPain,
	})
}

func (priest *Priest) applyShadowAffinity() {
	if priest.Talents.ShadowAffinity == 0 {
		return
	}

	threatReduction := []float64{0, -0.08, -0.16, -0.25}[priest.Talents.ShadowAffinity]

	priest.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_ThreatMultiplier_Pct,
		FloatValue: threatReduction,
		ClassMask:  PriestShadowSpells,
	})
}

func (priest *Priest) applyShadowWeaving() {
	if priest.Talents.ShadowWeaving == 0 {
		return
	}

	swAuras := priest.NewEnemyAuraArray(core.ShadowWeavingAura)

	priest.MakeProcTriggerAura(core.ProcTrigger{
		Name:             "Shadow Weaving Trigger",
		CanProcFromProcs: true, // 15257, 15331-15334 carry the bit.
		ClassSpellMask:   PriestShadowSpells,
		Callback:         core.CallbackOnSpellHitDealt,
		Outcome:          core.OutcomeLanded,
		ProcChance:       spellData.ShadowWeaving.FractionAt(priest.Talents.ShadowWeaving),
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			swAuras.Get(result.Target).Activate(sim)
			swAuras.Get(result.Target).AddStack(sim)
		},
	})
}

func (priest *Priest) applyShadowform() {
	if !priest.Talents.Shadowform {
		return
	}

	shadowformAura := priest.RegisterAura(core.Aura{
		Label:    "Shadowform",
		ActionID: core.ActionID{SpellID: 15473},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			if priest.SelfBuffs.PreShadowform {
				aura.Activate(sim)
			}
		},
		// Casting any holy-school spell breaks Shadowform.
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.SpellSchool.Matches(core.SpellSchoolHoly) {
				aura.Deactivate(sim)
			}
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: 0.15,
		ClassMask:  PriestShadowSpells,
	}).AttachMultiplicativePseudoStatBuff(
		&priest.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexPhysical], 0.85,
	)

	priest.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 15473},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: PriestSpellShadowform,
		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 32,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			shadowformAura.Activate(sim)
		},
	})
}

func (priest *Priest) applyVampiricEmbrace() {
	if !priest.Talents.VampiricEmbrace {
		return
	}

	// TODO: Forever drops Improved Vampiric Embrace; base heal percent only until we know
	// whether the bonus moved onto another talent.
	healPct := 0.15
	healthMetrics := priest.NewHealthMetrics(core.ActionID{SpellID: 15286})

	veDebuffAuras := priest.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		aura := target.RegisterAura(core.Aura{
			Label:    "Vampiric Embrace",
			ActionID: core.ActionID{SpellID: 15286},
			Duration: time.Second * 60,
		})
		aura.AttachProcTriggerCallback(target, core.ProcTrigger{
			Name:               "Vampiric Embrace Proc",
			Callback:           core.CallbackOnSpellHitTaken | core.CallbackOnPeriodicDamageTaken,
			ClassSpellMask:     PriestShadowSpells,
			RequireDamageDealt: true,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				priest.GainHealth(sim, result.Damage*healPct, healthMetrics)
			},
		})
		return aura
	})

	priest.VampiricEmbrace = priest.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 15286},
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: PriestSpellVampiricEmbrace,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 2,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: time.Second * 10,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			veDebuffAuras.Get(target).Activate(sim)
		},

		RelatedAuraArrays: veDebuffAuras.ToMap(),
	})
}

// applyPowerInLight implements Power in Light, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applyPowerInLight() {
	if priest.Talents.PowerInLight == 0 {
		return
	}

	panic("To be implemented")
}

// applyWandSpecialization implements Wand Specialization, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applyWandSpecialization() {
	if priest.Talents.WandSpecialization == 0 {
		return
	}

	panic("To be implemented")
}

// applyTwinDisciplines implements Twin Disciplines, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applyTwinDisciplines() {
	if priest.Talents.TwinDisciplines == 0 {
		return
	}

	panic("To be implemented")
}

// applyHolyPrecision implements Holy Precision, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applyHolyPrecision() {
	if priest.Talents.HolyPrecision == 0 {
		return
	}

	panic("To be implemented")
}

// applyImprovedPowerWordShield implements Improved Power Word: Shield, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applyImprovedPowerWordShield() {
	if priest.Talents.ImprovedPowerWordShield == 0 {
		return
	}

	panic("To be implemented")
}

// applyMartyrdom implements Martyrdom, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applyMartyrdom() {
	if priest.Talents.Martyrdom == 0 {
		return
	}

	panic("To be implemented")
}

// applyImprovedInnerFire implements Improved Inner Fire, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applyImprovedInnerFire() {
	if priest.Talents.ImprovedInnerFire == 0 {
		return
	}

	panic("To be implemented")
}

// applySoulWarding implements Soul Warding, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applySoulWarding() {
	if !priest.Talents.SoulWarding {
		return
	}

	panic("To be implemented")
}

// applyImprovedManaBurn implements Improved Mana Burn, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applyImprovedManaBurn() {
	if priest.Talents.ImprovedManaBurn == 0 {
		return
	}

	panic("To be implemented")
}

// applyPenance implements Penance, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applyPenance() {
	if !priest.Talents.Penance {
		return
	}

	panic("To be implemented")
}

// applyRenewedHope implements Renewed Hope, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applyRenewedHope() {
	if priest.Talents.RenewedHope == 0 {
		return
	}

	panic("To be implemented")
}

// applyDivineAegis implements Divine Aegis, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applyDivineAegis() {
	if priest.Talents.DivineAegis == 0 {
		return
	}

	panic("To be implemented")
}

// applyTwilightFocus implements Twilight Focus, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applyTwilightFocus() {
	if priest.Talents.TwilightFocus == 0 {
		return
	}

	panic("To be implemented")
}

// applyImprovedRenew implements Improved Renew, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applyImprovedRenew() {
	if priest.Talents.ImprovedRenew == 0 {
		return
	}

	panic("To be implemented")
}

// applyHolySpecialization implements Holy Specialization, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applyHolySpecialization() {
	if priest.Talents.HolySpecialization == 0 {
		return
	}

	panic("To be implemented")
}

// applySpellWarding implements Spell Warding, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applySpellWarding() {
	if priest.Talents.SpellWarding == 0 {
		return
	}

	panic("To be implemented")
}

// applyBlessedRecovery implements Blessed Recovery, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applyBlessedRecovery() {
	if priest.Talents.BlessedRecovery == 0 {
		return
	}

	panic("To be implemented")
}

// applyInspiration implements Inspiration, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applyInspiration() {
	if priest.Talents.Inspiration == 0 {
		return
	}

	panic("To be implemented")
}

// applyHolyReach implements Holy Reach, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applyHolyReach() {
	if priest.Talents.HolyReach == 0 {
		return
	}

	panic("To be implemented")
}

// applyImprovedHealing implements Improved Healing, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applyImprovedHealing() {
	if priest.Talents.ImprovedHealing == 0 {
		return
	}

	panic("To be implemented")
}

// applyBindingHeal implements Binding Heal, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applyBindingHeal() {
	if !priest.Talents.BindingHeal {
		return
	}

	panic("To be implemented")
}

// applyLitanyOfLight implements Litany of Light, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applyLitanyOfLight() {
	if priest.Talents.LitanyOfLight == 0 {
		return
	}

	panic("To be implemented")
}

// applySpiritualHealing implements Spiritual Healing, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applySpiritualHealing() {
	if priest.Talents.SpiritualHealing == 0 {
		return
	}

	panic("To be implemented")
}

// applyPrayerOfMending implements Prayer of Mending, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applyPrayerOfMending() {
	if !priest.Talents.PrayerOfMending {
		return
	}

	panic("To be implemented")
}

// applyBlackout implements Blackout, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applyBlackout() {
	if priest.Talents.Blackout == 0 {
		return
	}

	panic("To be implemented")
}

// applySpiritTap implements Spirit Tap, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applySpiritTap() {
	if priest.Talents.SpiritTap == 0 {
		return
	}

	panic("To be implemented")
}

// applyShadowReach implements Shadow Reach, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applyShadowReach() {
	if priest.Talents.ShadowReach == 0 {
		return
	}

	panic("To be implemented")
}

// applyImprovedPsychicScream implements Improved Psychic Scream, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applyImprovedPsychicScream() {
	if priest.Talents.ImprovedPsychicScream == 0 {
		return
	}

	panic("To be implemented")
}

// applyImprovedMindFlay implements Improved Mind Flay, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applyImprovedMindFlay() {
	if priest.Talents.ImprovedMindFlay == 0 {
		return
	}

	panic("To be implemented")
}

// applyImprovedFade implements Improved Fade, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applyImprovedFade() {
	if priest.Talents.ImprovedFade == 0 {
		return
	}

	panic("To be implemented")
}

// applySilence implements Silence, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applySilence() {
	if !priest.Talents.Silence {
		return
	}

	panic("To be implemented")
}

// applyDevouringContagion implements Devouring Contagion, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applyDevouringContagion() {
	if priest.Talents.DevouringContagion == 0 {
		return
	}

	panic("To be implemented")
}

// applyEarlyDemise implements Early Demise, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applyEarlyDemise() {
	if priest.Talents.EarlyDemise == 0 {
		return
	}

	panic("To be implemented")
}
