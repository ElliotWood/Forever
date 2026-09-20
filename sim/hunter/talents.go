package hunter

import (
	"time"

	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/stats"
)

func (hunter *Hunter) ApplyTalents() {
	// Beast Mastery
	hunter.registerEnduranceTraining()
	hunter.registerFocusedFire()
	hunter.registerUnleashedFury()
	hunter.registerFerocity()
	// Bestial Discipline handled in pet.go
	hunter.registerFrenzy()
	hunter.registerBestialWrath()

	// Marksmanship
	// Improved Hunter's Mark handled in hunters_mark.go
	hunter.registerEfficiency()
	hunter.registerImprovedArcaneShot()
	hunter.registerAimedShot()
	hunter.registerRapidKilling()
	hunter.registerImprovedStings()
	hunter.registerMortalShots()
	hunter.registerBarrage()
	hunter.registerRangedWeaponSpecialization()
	hunter.registerCarefulAim()
	// Trueshot Aura handled as a group buff in hunter.go

	// Survival
	hunter.registerHawkEye()
	hunter.registerSavageStrikes()
	hunter.registerSurvivalist()
	hunter.registerSurefooted()
	hunter.registerResourcefulness()
	hunter.registerLightningReflexes()
	// Readiness: Forever drops the talent; see registerReadiness

	// Forever additions, not yet implemented.
	hunter.registerDeadlyAspects()
	hunter.registerImprovedAspectOfTheMonkey()
	hunter.registerPathfinding()
	hunter.registerImprovedRevivePet()
	hunter.registerBestialSwiftness()
	hunter.registerImprovedMendPet()
	hunter.registerSummonHawk()
	hunter.registerSpiritBond()
	hunter.registerIntimidation()
	hunter.registerImprovedConcussiveShot()
	hunter.registerLethalAttacks()
	hunter.registerLoneWolf()
	hunter.registerImprovedSerpentSting()
	hunter.registerRapidRecuperation()
	hunter.registerScatterShot()
	hunter.registerSniperShot()
	hunter.registerImprovedTracking()
	hunter.registerDeflection()
	hunter.registerEntrapment()
	hunter.registerImprovedWingClip()
	hunter.registerCleverTraps()
	hunter.registerDeterrence()
	hunter.registerSurvivalTactics()
	hunter.registerPredatorsEdge()
	hunter.registerCounterattack()
	hunter.registerExposePrey()
	hunter.registerSurvivalistsDiscipline()
	hunter.registerStriderKick()
	hunter.registerLaceratingStrikes()

	if hunter.Pet != nil {
		hunter.Pet.ApplyTalents()
	}
}

func (hunter *Hunter) registerEnduranceTraining() {
	if hunter.Pet == nil || hunter.Talents.EnduranceTraining == 0 {
		return
	}

	hunter.Pet.StatDependencyManager.EnableDynamicStatDep(
		hunter.Pet.NewDynamicMultiplyStat(stats.Health, spellData.EnduranceTraining.Effect(shared.A_ADD_FLAT_MODIFIER, shared.SPELLMOD_ALL_EFFECTS).MultiplierAt(hunter.Talents.EnduranceTraining)),
	)

	hunter.StatDependencyManager.EnableDynamicStatDep(
		hunter.NewDynamicMultiplyStat(stats.Health, spellData.EnduranceTraining.Effect(shared.A_MOD_INCREASE_HEALTH_PERCENT, 0).MultiplierAt(hunter.Talents.EnduranceTraining)),
	)
}

func (hunter *Hunter) registerFocusedFire() {
	if hunter.Pet == nil || hunter.Talents.FocusedFire == 0 {
		return
	}

	hunter.PseudoStats.DamageDealtMultiplier *= spellData.FocusedFire.Effect(shared.A_NONE, 0).MultiplierAt(hunter.Talents.FocusedFire)
	hunter.Pet.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_BonusCrit_Percent,
		ClassMask:  HunterSpellKillCommandPet,
		FloatValue: spellData.FocusedFire.Effect(shared.A_ADD_FLAT_MODIFIER, shared.SPELLMOD_CRITICAL_CHANCE).ValueAt(hunter.Talents.FocusedFire),
	})
}

func (hunter *Hunter) registerUnleashedFury() {
	if hunter.Pet == nil || hunter.Talents.UnleashedFury == 0 {
		return
	}

	hunter.Pet.PseudoStats.DamageDealtMultiplier *= spellData.UnleashedFury.MultiplierAt(hunter.Talents.UnleashedFury)
}

func (hunter *Hunter) registerFerocity() {
	if hunter.Pet == nil || hunter.Talents.Ferocity == 0 {
		return
	}

	hunter.Pet.AddStats(stats.Stats{
		stats.PhysicalCritPercent: spellData.Ferocity.ValueAt(hunter.Talents.Ferocity),
		stats.SpellCritPercent:    spellData.Ferocity.ValueAt(hunter.Talents.Ferocity),
	})
}

func (hunter *Hunter) registerFrenzy() {
	if hunter.Pet == nil || hunter.Talents.Frenzy == 0 {
		return
	}

	frenzy := hunter.Pet.RegisterAura(core.Aura{
		Label:    "Frenzy Effect",
		ActionID: core.ActionID{SpellID: 19615},
		Duration: time.Second * 8,
	}).AttachMultiplyMeleeSpeed(1.3)

	hunter.Pet.MakeProcTriggerAura(core.ProcTrigger{
		Name:       "Frenzy",
		Callback:   core.CallbackOnSpellHitDealt,
		Outcome:    core.OutcomeCrit,
		ProcChance: spellData.Frenzy.FractionAt(hunter.Talents.Frenzy),

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			frenzy.Activate(sim)
		},
	})
}

func (hunter *Hunter) registerBestialWrath() {
	if hunter.Pet == nil || !hunter.Talents.BestialWrath {
		return
	}

	actionID := core.ActionID{SpellID: 19574}

	hunter.Pet.BestialWrathAura = hunter.Pet.RegisterAura(core.Aura{
		Label:    "Bestial Wrath",
		ActionID: actionID,
		Duration: time.Second * 18,
	}).AttachMultiplicativePseudoStatBuff(
		&hunter.Pet.PseudoStats.DamageDealtMultiplier, 1.5,
	)

	hunter.BestialWrath = hunter.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ClassSpellMask: HunterSpellBestialWrath,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 10,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Minute * 2,
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hunter.GCD.IsReady(sim)
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			hunter.Pet.BestialWrathAura.Activate(sim)
		},
	})

	hunter.AddMajorCooldown(core.MajorCooldown{
		Spell: hunter.BestialWrath,
		Type:  core.CooldownTypeDPS,
	})
}

func (hunter *Hunter) registerEfficiency() {
	if hunter.Talents.Efficiency == 0 {
		return
	}

	hunter.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_PowerCost_Pct_Add,
		ClassMask:  HunterSpellsShotsAndStings,
		FloatValue: -0.02 * float64(hunter.Talents.Efficiency),
	})
}

func (hunter *Hunter) registerImprovedArcaneShot() {
	if hunter.Talents.ImprovedArcaneShot == 0 {
		return
	}

	hunter.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Flat,
		ClassMask: HunterSpellArcaneShot,
		TimeValue: -core.DurationFromSeconds(0.2 * float64(hunter.Talents.ImprovedArcaneShot)),
	})
}

var aimedShotRank = spellData.AimedShot.HighestRank()

func (hunter *Hunter) registerAimedShot() {
	hunter.AimedShot = hunter.RegisterRangedSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: aimedShotRank.SpellID},
		SpellSchool:    aimedShotRank.SpellSchool,
		DefenseType:    aimedShotRank.DefenseType,
		ClassSpellMask: HunterSpellAimedShot,
		ProcMask:       core.ProcMaskRangedSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			FlatCost: aimedShotRank.Cost,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				CastTime: time.Millisecond * 3000,
			},
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: aimedShotRank.Cooldown,
			},
		},

		BonusCoefficient: aimedShotRank.Direct.BonusCoefficient(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0.2*spell.RangedAttackPower(target) +
				hunter.AutoAttacks.Ranged().BaseDamage(sim) +
				hunter.talonOfAlarBonus() +
				aimedShotRank.Direct.Damage(sim)

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}

func (hunter *Hunter) registerRapidKilling() {
	if hunter.Talents.RapidKilling == 0 {
		return
	}

	hunter.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Flat,
		ClassMask: HunterSpellRapidFire,
		TimeValue: -core.DurationFromSeconds(60 * float64(hunter.Talents.RapidKilling)),
	})
}

func (hunter *Hunter) registerImprovedStings() {
	if hunter.Talents.ImprovedStings == 0 {
		return
	}

	hunter.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		ClassMask:  HunterSpellSerpentSting,
		FloatValue: spellData.ImprovedStings.Effect(shared.A_ADD_PCT_MODIFIER, shared.SPELLMOD_DOT).FractionAt(hunter.Talents.ImprovedStings),
	})
}

func (hunter *Hunter) registerMortalShots() {
	if hunter.Talents.MortalShots == 0 {
		return
	}

	hunter.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_CritMultiplier_Flat,
		ProcMask:   core.ProcMaskRanged,
		FloatValue: spellData.MortalShots.FractionAt(hunter.Talents.MortalShots),
	})
}

func (hunter *Hunter) registerBarrage() {
	if hunter.Talents.Barrage == 0 {
		return
	}

	hunter.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		ClassMask:  HunterSpellMultiShot | HunterSpellVolley,
		FloatValue: spellData.Barrage.Effect(shared.A_ADD_PCT_MODIFIER, shared.SPELLMOD_DAMAGE).FractionAt(hunter.Talents.Barrage),
	})
}

func (hunter *Hunter) registerRangedWeaponSpecialization() {
	if hunter.Talents.RangedWeaponSpecialization == 0 {
		return
	}

	hunter.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ProcMask:   core.ProcMaskRanged,
		FloatValue: spellData.RangedWeaponSpecialization.FractionAt(hunter.Talents.RangedWeaponSpecialization),
	})
}

func (hunter *Hunter) registerCarefulAim() {
	if hunter.Talents.CarefulAim == 0 {
		return
	}

	hunter.AddStatDependency(stats.Intellect, stats.RangedAttackPower, 0.15*float64(hunter.Talents.CarefulAim))
}

func (hunter *Hunter) registerHawkEye() {
	if hunter.Talents.HawkEye == 0 {
		return
	}

	bonusRange := float64(hunter.Talents.HawkEye) * 2
	ranged := hunter.AutoAttacks.Ranged()

	if ranged != nil {
		ranged.MaxRange += bonusRange
	}

	hunter.AddStaticMod(core.SpellModConfig{
		Kind:     core.SpellMod_Custom,
		ProcMask: core.ProcMaskRanged,
		ApplyCustom: func(mod *core.SpellMod, spell *core.Spell) {
			if spell.MaxRange > 0 {
				spell.MaxRange += bonusRange
			}
		},
		RemoveCustom: func(mod *core.SpellMod, spell *core.Spell) {
			if spell.MaxRange > 0 {
				spell.MaxRange -= bonusRange
			}
		},
	})
}

func (hunter *Hunter) registerSavageStrikes() {
	if hunter.Talents.SavageStrikes == 0 {
		return
	}

	hunter.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_BonusCrit_Percent,
		ClassMask:  HunterSpellRaptorStrike,
		FloatValue: spellData.SavageStrikes.ValueAt(hunter.Talents.SavageStrikes),
	})
}

func (hunter *Hunter) registerSurvivalist() {
	if hunter.Talents.Survivalist == 0 {
		return
	}

	hunter.MultiplyStat(stats.Health, spellData.Survivalist.MultiplierAt(hunter.Talents.Survivalist))
}

func (hunter *Hunter) registerSurefooted() {
	if hunter.Talents.Surefooted == 0 {
		return
	}

	hunter.AddStat(stats.PhysicalHitPercent, float64(hunter.Talents.Surefooted))
}

func (hunter *Hunter) registerResourcefulness() {
	if hunter.Talents.Resourcefulness == 0 {
		return
	}

	hunter.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_PowerCost_Pct_Add,
		ClassMask:  HunterSpellRaptorStrike,
		FloatValue: -0.2 * float64(hunter.Talents.Resourcefulness),
	})
}

func (hunter *Hunter) registerLightningReflexes() {
	if hunter.Talents.LightningReflexes == 0 {
		return
	}

	hunter.MultiplyStat(stats.Agility, spellData.LightningReflexes.MultiplierAt(hunter.Talents.LightningReflexes))
}

// TODO: uncalled -- Forever drops the Readiness talent; re-gate before wiring back
// into ApplyTalents.
func (hunter *Hunter) registerReadiness() {
	hunter.Readiness = hunter.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 23989},
		SpellSchool:    core.SpellSchoolPhysical,
		ClassSpellMask: HunterSpellReadiness,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Minute * 5,
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !hunter.RapidFire.IsReady(sim) ||
				!hunter.MultiShot.IsReady(sim) ||
				!hunter.ArcaneShot.IsReady(sim) ||
				!hunter.KillCommand.IsReady(sim) ||
				!hunter.RaptorStrike.IsReady(sim)
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			hunter.RapidFire.CD.Reset()
			hunter.MultiShot.CD.Reset()
			hunter.ArcaneShot.CD.Reset()
			hunter.KillCommand.CD.Reset()
			hunter.RaptorStrike.CD.Reset()
		},
	})

	hunter.AddMajorCooldown(core.MajorCooldown{
		Spell: hunter.Readiness,
		Type:  core.CooldownTypeDPS,

		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			return !hunter.RapidFire.IsReady(sim)
		},
	})
}

// registerDeadlyAspects implements Deadly Aspects, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerDeadlyAspects() {
	if hunter.Talents.DeadlyAspects == 0 {
		return
	}

	panic("To be implemented")
}

// registerImprovedAspectOfTheMonkey implements Improved Aspect of the Monkey, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerImprovedAspectOfTheMonkey() {
	if hunter.Talents.ImprovedAspectOfTheMonkey == 0 {
		return
	}

	panic("To be implemented")
}

// registerPathfinding implements Pathfinding, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerPathfinding() {
	if hunter.Talents.Pathfinding == 0 {
		return
	}

	panic("To be implemented")
}

// registerImprovedRevivePet implements Improved Revive Pet, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerImprovedRevivePet() {
	if hunter.Talents.ImprovedRevivePet == 0 {
		return
	}

	panic("To be implemented")
}

// registerBestialSwiftness implements Bestial Swiftness, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerBestialSwiftness() {
	if !hunter.Talents.BestialSwiftness {
		return
	}

	panic("To be implemented")
}

// registerImprovedMendPet implements Improved Mend Pet, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerImprovedMendPet() {
	if hunter.Talents.ImprovedMendPet == 0 {
		return
	}

	panic("To be implemented")
}

// registerSummonHawk implements Summon Hawk, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerSummonHawk() {
	if !hunter.Talents.SummonHawk {
		return
	}

	panic("To be implemented")
}

// registerSpiritBond implements Spirit Bond, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerSpiritBond() {
	if hunter.Talents.SpiritBond == 0 {
		return
	}

	panic("To be implemented")
}

// registerIntimidation implements Intimidation, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerIntimidation() {
	if !hunter.Talents.Intimidation {
		return
	}

	panic("To be implemented")
}

// registerImprovedConcussiveShot implements Improved Concussive Shot, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerImprovedConcussiveShot() {
	if hunter.Talents.ImprovedConcussiveShot == 0 {
		return
	}

	panic("To be implemented")
}

// registerLethalAttacks implements Lethal Attacks, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerLethalAttacks() {
	if hunter.Talents.LethalAttacks == 0 {
		return
	}

	panic("To be implemented")
}

// registerLoneWolf implements Lone Wolf, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerLoneWolf() {
	if !hunter.Talents.LoneWolf {
		return
	}

	panic("To be implemented")
}

// registerImprovedSerpentSting implements Improved Serpent Sting, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerImprovedSerpentSting() {
	if hunter.Talents.ImprovedSerpentSting == 0 {
		return
	}

	panic("To be implemented")
}

// registerRapidRecuperation implements Rapid Recuperation, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerRapidRecuperation() {
	if hunter.Talents.RapidRecuperation == 0 {
		return
	}

	panic("To be implemented")
}

// registerScatterShot implements Scatter Shot, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerScatterShot() {
	if !hunter.Talents.ScatterShot {
		return
	}

	panic("To be implemented")
}

// registerSniperShot implements Sniper Shot, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerSniperShot() {
	if !hunter.Talents.SniperShot {
		return
	}

	panic("To be implemented")
}

// registerImprovedTracking implements Improved Tracking, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerImprovedTracking() {
	if hunter.Talents.ImprovedTracking == 0 {
		return
	}

	panic("To be implemented")
}

// registerDeflection implements Deflection, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerDeflection() {
	if hunter.Talents.Deflection == 0 {
		return
	}

	panic("To be implemented")
}

// registerEntrapment implements Entrapment, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerEntrapment() {
	if hunter.Talents.Entrapment == 0 {
		return
	}

	panic("To be implemented")
}

// registerImprovedWingClip implements Improved Wing Clip, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerImprovedWingClip() {
	if hunter.Talents.ImprovedWingClip == 0 {
		return
	}

	panic("To be implemented")
}

// registerCleverTraps implements Clever Traps, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerCleverTraps() {
	if hunter.Talents.CleverTraps == 0 {
		return
	}

	panic("To be implemented")
}

// registerDeterrence implements Deterrence, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerDeterrence() {
	if !hunter.Talents.Deterrence {
		return
	}

	panic("To be implemented")
}

// registerSurvivalTactics implements Survival Tactics, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerSurvivalTactics() {
	if hunter.Talents.SurvivalTactics == 0 {
		return
	}

	panic("To be implemented")
}

// registerPredatorsEdge implements Predator's Edge, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerPredatorsEdge() {
	if hunter.Talents.PredatorsEdge == 0 {
		return
	}

	panic("To be implemented")
}

// registerCounterattack implements Counterattack, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerCounterattack() {
	if !hunter.Talents.Counterattack {
		return
	}

	panic("To be implemented")
}

// registerExposePrey implements Expose Prey, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerExposePrey() {
	if hunter.Talents.ExposePrey == 0 {
		return
	}

	panic("To be implemented")
}

// registerSurvivalistsDiscipline implements Survivalist's Discipline, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerSurvivalistsDiscipline() {
	if hunter.Talents.SurvivalistsDiscipline == 0 {
		return
	}

	panic("To be implemented")
}

// registerStriderKick implements Strider Kick, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerStriderKick() {
	if !hunter.Talents.StriderKick {
		return
	}

	panic("To be implemented")
}

// registerLaceratingStrikes implements Lacerating Strikes, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerLaceratingStrikes() {
	if !hunter.Talents.LaceratingStrikes {
		return
	}

	panic("To be implemented")
}
