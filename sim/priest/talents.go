package priest

import (
	"slices"
	"time"

	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/stats"
)

func (priest *Priest) ApplyTalents() {
	// Discipline
	priest.registerInnerFocus()
	priest.applyPowerInLight()
	priest.applyTwinDisciplines()
	priest.applyHolyPrecision()
	priest.applyMentalAgility()

	if priest.Talents.SilentResolve > 0 {
		priest.OnSpellRegistered(func(spell *core.Spell) {
			if spell.Flags.Matches(SpellFlagPriest) && spell.SpellSchool.Matches(core.SpellSchoolHoly) {
				spell.ThreatMultiplier *= 1 - .1*float64(priest.Talents.SilentResolve)
			}
		})
	}

	priest.PseudoStats.SpiritRegenRateCasting = []float64{0.0, 0.17, 0.34, 0.51}[priest.Talents.Meditation]

	if priest.Talents.MentalStrength > 0 {
		priest.MultiplyStat(stats.Intellect, 1.0+0.03*float64(priest.Talents.MentalStrength))
	}

	// Holy
	priest.applyInspiration()
	priest.applyHolySpecialization()
	priest.applySearingLight()

	priest.PseudoStats.SchoolDamageTakenMultiplier.MultiplyMagicSchools(1 - 0.02*float64(priest.Talents.SpellWarding))

	// TODO: only rank 1 was shown, beta will confirm that the two halves scale at 5% and 1% per point
	if priest.Talents.SpiritualGuidance > 0 {
		priest.AddStatDependency(stats.Spirit, stats.HealingPower, 0.05*float64(priest.Talents.SpiritualGuidance))
		priest.AddStatDependency(stats.Spirit, stats.SpellDamage, 0.01*float64(priest.Talents.SpiritualGuidance))
	}

	// Shadow Magic
	priest.registerVampiricEmbraceSpell()
	priest.registerShadowform()
	priest.applySpiritTap()
	priest.applyShadowAffinity()
	priest.applyShadowFocus()
	priest.applyShadowWeaving()
	priest.applyDarkness()
}

// Smite and Penance hit harder while the target is burning from this priest's Holy Fire.
// TODO: only rank 1 was shown, beta will confirm the 2% per point
func (priest *Priest) applyPowerInLight() {
	if priest.Talents.PowerInLight == 0 {
		return
	}

	multiplier := 1 + 0.02*float64(priest.Talents.PowerInLight)
	affectedSpellCodes := []int32{SpellCode_PriestSmite, SpellCode_PriestPenance}

	for _, target := range priest.Env.Encounter.TargetUnits {
		target.AddDynamicDamageTakenModifier(func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Unit == &priest.Unit && slices.Contains(affectedSpellCodes, spell.SpellCode) && priest.hasActiveHolyFire(result.Target) {
				result.Damage *= multiplier
			}
		})
	}
}

func (priest *Priest) hasActiveHolyFire(target *core.Unit) bool {
	for _, spell := range priest.HolyFire {
		if spell != nil && spell.Dot(target).IsActive() {
			return true
		}
	}
	return false
}

func (priest *Priest) applyTwinDisciplines() {
	if priest.Talents.TwinDisciplines == 0 {
		return
	}

	points := float64(priest.Talents.TwinDisciplines)
	priest.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Flags.Matches(SpellFlagPriest) && spell.DefaultCast.CastTime == 0 {
			spell.DamageMultiplierAdditive += 0.01 * points
		}
	})
}

func (priest *Priest) applyHolyPrecision() {
	if priest.Talents.HolyPrecision == 0 {
		return
	}

	bonusHit := 6 * float64(priest.Talents.HolyPrecision) * core.SpellHitRatingPerHitChance
	priest.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Flags.Matches(SpellFlagPriest) && spell.SpellSchool.Matches(core.SpellSchoolHoly) {
			spell.BonusHitRating += bonusHit
		}
	})
}

func (priest *Priest) applyMentalAgility() {
	if priest.Talents.MentalAgility == 0 {
		return
	}

	affectedSpellCodes := []int32{SpellCode_PriestSmite, SpellCode_PriestHolyFire}
	priest.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Cost == nil || !spell.Flags.Matches(SpellFlagPriest) {
			return
		}

		if spell.DefaultCast.CastTime == 0 || slices.Contains(affectedSpellCodes, spell.SpellCode) {
			spell.Cost.Multiplier -= 3 * priest.Talents.MentalAgility
		}
	})
}

func (priest *Priest) applyHolySpecialization() {
	if priest.Talents.HolySpecialization == 0 {
		return
	}

	priest.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Flags.Matches(SpellFlagPriest) && spell.SpellSchool.Matches(core.SpellSchoolHoly) {
			spell.BonusCritRating += 1 * float64(priest.Talents.HolySpecialization) * core.CritRatingPerCritChance
		}
	})
}

func (priest *Priest) applyInspiration() {
	if priest.Talents.Inspiration == 0 {
		return
	}

	auras := make([]*core.Aura, len(priest.Env.AllUnits))
	for _, unit := range priest.Env.AllUnits {
		if !priest.IsOpponent(unit) {
			aura := core.InspirationAura(unit, priest.Talents.Inspiration)
			auras[unit.UnitIndex] = aura
		}
	}

	priest.RegisterAura(core.Aura{
		Label:    "Inspiration Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnHealDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if slices.Contains([]int32{SpellCode_PriestFlashHeal, SpellCode_PriestHeal, SpellCode_PriestGreaterHeal}, spell.SpellCode) {
				auras[result.Target.UnitIndex].Activate(sim)
			}
		},
	})
}

// Searing Light now buffs every Holy spell and lets Holy Fire ticks refund the next Holy Nova.
func (priest *Priest) applySearingLight() {
	if priest.Talents.SearingLight == 0 {
		return
	}

	points := float64(priest.Talents.SearingLight)
	priest.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexHoly] *= 1 + 0.02*points

	priest.SearingLightAura = priest.RegisterAura(core.Aura{
		Label:    "Searing Light",
		ActionID: core.ActionID{SpellID: 14909},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if priest.HolyNova != nil {
				priest.HolyNova.Cost.Multiplier -= 100
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if priest.HolyNova != nil {
				priest.HolyNova.Cost.Multiplier += 100
			}
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.SpellCode == SpellCode_PriestHolyNova {
				aura.Deactivate(sim)
			}
		},
	})

	procChance := 0.05 * points
	core.MakePermanent(priest.RegisterAura(core.Aura{
		Label: "Searing Light Trigger",
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.SpellCode == SpellCode_PriestHolyFire && sim.Proc(procChance, "Searing Light") {
				priest.SearingLightAura.Activate(sim)
			}
		},
	}))
}

func (priest *Priest) applySpiritTap() {
	if priest.Talents.SpiritTap == 0 {
		return
	}

	spellID := []int32{0, 15270, 15335, 15336, 15337, 15338}[priest.Talents.SpiritTap]
	statDep := priest.NewDynamicMultiplyStat(stats.Spirit, 2.0)

	priest.SpiritTapAura = priest.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: spellID},
		Label:    "Spirit Tap",
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			priest.EnableDynamicStatDep(sim, statDep)
			priest.PseudoStats.SpiritRegenRateCasting += 0.50
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			priest.DisableDynamicStatDep(sim, statDep)
			priest.PseudoStats.SpiritRegenRateCasting -= 0.50
		},
	})
}

func (priest *Priest) applyShadowAffinity() {
	if priest.Talents.ShadowAffinity == 0 {
		return
	}

	priest.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Flags.Matches(SpellFlagPriest) && spell.SpellSchool.Matches(core.SpellSchoolShadow) {
			spell.ThreatMultiplier *= 1 - 0.1*float64(priest.Talents.ShadowAffinity)
		}
	})
}

func (priest *Priest) applyShadowFocus() {
	if priest.Talents.ShadowFocus == 0 {
		return
	}

	bonusHit := 1 * float64(priest.Talents.ShadowFocus) * core.SpellHitRatingPerHitChance
	priest.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Flags.Matches(SpellFlagPriest) && spell.SpellSchool.Matches(core.SpellSchoolShadow) {
			spell.BonusHitRating += bonusHit
		}
	})
}

// The raid debuff version of Shadow Weaving is a Classic mechanic, in Forever it buffs the priest instead.
func (priest *Priest) applyShadowWeaving() {
	if priest.Talents.ShadowWeaving == 0 {
		return
	}

	priest.shadowWeavingProcChance = 0.33 * float64(priest.Talents.ShadowWeaving)

	priest.ShadowWeavingAura = priest.RegisterAura(core.Aura{
		Label:     "Shadow Weaving",
		ActionID:  core.ActionID{SpellID: core.ShadowWeavingSpellIDs[int(priest.Talents.ShadowWeaving)]},
		Duration:  time.Second * 15,
		MaxStacks: 5,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] /= 1 + 0.02*float64(oldStacks)
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= 1 + 0.02*float64(newStacks)
		},
	})
}

func (priest *Priest) AddShadowWeavingStack(sim *core.Simulation) {
	if priest.ShadowWeavingAura == nil || !sim.Proc(priest.shadowWeavingProcChance, "Shadow Weaving") {
		return
	}

	priest.ShadowWeavingAura.Activate(sim)
	priest.ShadowWeavingAura.AddStack(sim)
}

func (priest *Priest) applyDarkness() {
	if priest.Talents.Darkness == 0 {
		return
	}

	multiplier := 0.02 * float64(priest.Talents.Darkness)

	priest.RegisterAura(core.Aura{
		Label: "Darkness",
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			baseDamageAffectedSpells := core.FilterSlice(
				core.Flatten(
					[][]*core.Spell{
						priest.MindBlast,
						priest.DevouringPlague,
					},
				),
				func(spell *core.Spell) bool { return spell != nil },
			)

			fullDamageAffectedSpells := core.FilterSlice(
				core.Flatten(
					[][]*core.Spell{
						priest.ShadowWordPain,
					},
				),
				func(spell *core.Spell) bool { return spell != nil },
			)

			for _, spells := range priest.MindFlay {
				fullDamageAffectedSpells = append(
					fullDamageAffectedSpells,
					core.FilterSlice(spells, func(spell *core.Spell) bool { return spell != nil })...,
				)
			}

			for _, spell := range baseDamageAffectedSpells {
				spell.BaseDamageMultiplierAdditive += multiplier
			}

			for _, spell := range fullDamageAffectedSpells {
				spell.DamageMultiplierAdditive += multiplier
			}
		},
	})
}

func (priest *Priest) registerInnerFocus() {
	if !priest.Talents.InnerFocus {
		return
	}

	actionID := core.ActionID{SpellID: 14751}

	priest.InnerFocusAura = priest.RegisterAura(core.Aura{
		Label:    "Inner Focus",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range priest.Spellbook {
				if spell.Flags.Matches(SpellFlagPriest) && spell.Cost != nil {
					spell.Cost.Multiplier -= 100
					spell.BonusCritRating += 25 * core.SpellCritRatingPerCritChance
				}
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range priest.Spellbook {
				if spell.Flags.Matches(SpellFlagPriest) && spell.Cost != nil {
					spell.Cost.Multiplier += 100
					spell.BonusCritRating -= 25 * core.SpellCritRatingPerCritChance
				}
			}
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Flags.Matches(SpellFlagPriest) {
				// Remove the buff and put skill on CD
				aura.Deactivate(sim)
				priest.InnerFocus.CD.Use(sim)
				priest.UpdateMajorCooldowns()
			}
		},
	})

	priest.InnerFocus = priest.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			priest.InnerFocusAura.Activate(sim)
		},
	})

	priest.AddMajorCooldown(core.MajorCooldown{
		Spell: priest.InnerFocus,
		Type:  core.CooldownTypeDPS,
	})
}

func (priest *Priest) registerShadowform() {
	if !priest.Talents.Shadowform {
		return
	}

	actionID := core.ActionID{SpellID: 15473}

	// The mana discount is the biggest single change in the tree and only the demo tooltip backs it up.
	// TODO: beta will confirm the 50%.
	priest.ShadowformAura = priest.RegisterAura(core.Aura{
		Label:    "Shadowform",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= 1.10
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexPhysical] *= 0.85
			for _, spell := range priest.shadowformSpells() {
				spell.CritDamageBonus += 1
				if spell.Cost != nil {
					spell.Cost.Multiplier -= 50
				}
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] /= 1.10
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexPhysical] /= 0.85
			for _, spell := range priest.shadowformSpells() {
				spell.CritDamageBonus -= 1
				if spell.Cost != nil {
					spell.Cost.Multiplier += 50
				}
			}
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.SpellSchool.Matches(core.SpellSchoolHoly) {
				aura.Deactivate(sim)
			}
		},
	})

	priest.Shadowform = priest.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: 0,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			priest.ShadowformAura.Activate(sim)
		},
	})
}

func (priest *Priest) shadowformSpells() []*core.Spell {
	return core.FilterSlice(priest.Spellbook, func(spell *core.Spell) bool {
		return spell.Flags.Matches(SpellFlagPriest) && spell.SpellSchool.Matches(core.SpellSchoolShadow)
	})
}
