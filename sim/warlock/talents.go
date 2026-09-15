package warlock

import (
	"slices"
	"time"

	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
	"github.com/wowsims/classic/sim/core/stats"
)

func (warlock *Warlock) ApplyTalents() {
	warlock.applyWeaponImbue()

	// Affliction
	warlock.applySuppression()
	warlock.applyMalediction()
	warlock.applyPandemic()
	warlock.applyMalevolence()
	warlock.applyNightfall()
	warlock.applyShadowMastery()

	// Demonology
	warlock.applyDemonicEmbrace()
	warlock.applyUnholyPower()
	warlock.applyFelVitality()
	warlock.registerFelDominationCD()
	warlock.applyMasterSummoner()
	warlock.applyDecimation()
	warlock.applyDemonicBrand()
	warlock.applyDemonicKnowledge()
	warlock.applyMasterDemonologist()
	warlock.applyDemonicSacrifice()
	warlock.applySoulLink()

	// Destruction
	warlock.applyImprovedShadowBolt()
	warlock.applyCataclysm()
	warlock.applyBane()
	warlock.applyRuin()
	warlock.applyAgonizingFlames()
	warlock.applyFireAndBrimstone()
	warlock.applyShadowAndFlame()
}

func (warlock *Warlock) applyWeaponImbue() {
	if warlock.GetCharacter().Equipment.OffHand().Type != proto.ItemType_ItemTypeUnknown {
		return
	}

	level := warlock.Level
	if warlock.Options.WeaponImbue == proto.WarlockOptions_Firestone {
		warlock.applyFirestone()
	}
	if warlock.Options.WeaponImbue == proto.WarlockOptions_Spellstone {
		if level >= 55 {
			warlock.AddStat(stats.SpellCrit, 1*core.SpellCritRatingPerCritChance)
		}
	}
}

func (warlock *Warlock) applyFirestone() {
	level := warlock.Level

	damageMin := 0.0
	damageMax := 0.0

	// TODO: Test for spell scaling
	spellCoeff := 0.0
	spellId := int32(0)

	// TODO: Test PPM
	ppm := warlock.AutoAttacks.NewPPMManager(8, core.ProcMaskMelee)

	firestoneMulti := 1.0

	if level >= 56 {
		warlock.AddStat(stats.FirePower, 21*firestoneMulti)
		damageMin = 80.0
		damageMax = 120.0
		spellId = 17949
	} else if level >= 46 {
		warlock.AddStat(stats.FirePower, 17*firestoneMulti)
		damageMin = 60.0
		damageMax = 90.0
		spellId = 17947
	} else if level >= 36 {
		warlock.AddStat(stats.FirePower, 14*firestoneMulti)
		damageMin = 40.0
		damageMax = 60.0
		spellId = 17945
	} else if level >= 28 {
		warlock.AddStat(stats.FirePower, 10*firestoneMulti)
		damageMin = 25.0
		damageMax = 35.0
		spellId = 758
	}

	if level >= 28 && warlock.Consumes.MainHandImbue == proto.WeaponImbue_WeaponImbueUnknown {
		fireProcSpell := warlock.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: spellId},
			SpellSchool: core.SpellSchoolFire,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier:         firestoneMulti,
			ThreatMultiplier:         1,
			DamageMultiplierAdditive: 1,
			BonusCoefficient:         spellCoeff,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				baseDamage := sim.Roll(damageMin, damageMax)

				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicCrit)
			},
		})

		core.MakePermanent(warlock.GetOrRegisterAura(core.Aura{
			Label: "Firestone Proc",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if !spell.ProcMask.Matches(core.ProcMaskMelee) {
					return
				}

				if !ppm.Proc(sim, core.ProcMaskMelee, "Firestone Proc") {
					return
				}

				fireProcSpell.Cast(sim, result.Target)
			},
		}))
	}
}

///////////////////////////////////////////////////////////////////////////
//                            Affliction
///////////////////////////////////////////////////////////////////////////

func (warlock *Warlock) applySuppression() {
	if warlock.Talents.Suppression == 0 {
		return
	}

	points := float64(warlock.Talents.Suppression)
	warlock.AddStat(stats.SpellHit, points*core.SpellHitRatingPerHitChance)
	warlock.AddStat(stats.MeleeHit, points*core.MeleeHitRatingPerHitChance)
	warlock.PseudoStats.ThreatMultiplier *= 1 - 0.04*points
}

func (warlock *Warlock) applyMalediction() {
	if warlock.Talents.Malediction == 0 {
		return
	}

	points := float64(warlock.Talents.Malediction)
	warlock.OnSpellRegistered(func(spell *core.Spell) {
		if isWarlockSpell(spell) && (len(spell.Dots()) > 0 || spell.AOEDot() != nil) {
			spell.PeriodicDamageMultiplierAdditive += 0.01 * points
		}
	})
}

// Drain Life and Drain Soul deal 2% more per other Affliction effect on the target (up to 3), tripled on targets below 20%
func (warlock *Warlock) improvedDrainsMultiplier(sim *core.Simulation, target *core.Unit) float64 {
	if warlock.Talents.ImprovedDrains == 0 {
		return 1
	}

	effects := 0
	for _, spell := range warlock.DoTSpells {
		if spell.Flags.Matches(WarlockFlagAffliction) && spell.Dot(target).IsActive() {
			effects++
		}
	}
	for _, spell := range warlock.DebuffSpells {
		if spell.Flags.Matches(WarlockFlagAffliction) && spell.RelatedAuras[0].Get(target).IsActive() {
			effects++
		}
	}

	bonus := 0.02 * float64(warlock.Talents.ImprovedDrains) * float64(min(effects, 3))
	if sim.IsExecutePhase20() {
		bonus *= 3
	}

	return 1 + bonus
}

func (warlock *Warlock) drainTickLength(baseTickLength time.Duration) time.Duration {
	return time.Duration(float64(baseTickLength) / (1 + 0.17*float64(warlock.Talents.SoulSiphon)))
}

func (warlock *Warlock) applyPandemic() {
	if warlock.Talents.Pandemic == 0 {
		return
	}

	affectedSpellCodes := []int32{SpellCode_WarlockCorruption, SpellCode_WarlockBaneOfAgony, SpellCode_WarlockBaneOfDoom, SpellCode_WarlockDrainSoul, SpellCode_WarlockDrainLife, SpellCode_WarlockSiphonLife, SpellCode_WarlockDrainHope}
	bonus := 0.33 * float64(warlock.Talents.Pandemic)
	warlock.OnSpellRegistered(func(spell *core.Spell) {
		if slices.Contains(affectedSpellCodes, spell.SpellCode) {
			spell.CritDamageBonus += bonus
		}
	})
}

func (warlock *Warlock) applyMalevolence() {
	if warlock.Talents.Malevolence == 0 {
		return
	}

	points := float64(warlock.Talents.Malevolence)
	warlock.OnSpellRegistered(func(spell *core.Spell) {
		if spell.SpellSchool.Matches(core.SpellSchoolShadow) && isWarlockSpell(spell) {
			spell.BonusCritRating += points * core.SpellCritRatingPerCritChance
		}
	})
}

func (warlock *Warlock) applyNightfall() {
	if warlock.Talents.Nightfall <= 0 {
		return
	}

	shadowTranceAura := warlock.RegisterAura(core.Aura{
		Label:    "Nightfall Shadow Trance",
		ActionID: core.ActionID{SpellID: 17941},
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range warlock.ShadowBolt {
				spell.CastTimeMultiplier -= 1
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range warlock.ShadowBolt {
				spell.CastTimeMultiplier += 1
			}
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			// Check if the shadowbolt was instant cast and not a normal one
			if spell.SpellCode == SpellCode_WarlockShadowBolt && spell.CurCast.CastTime == 0 {
				aura.Deactivate(sim)
			}
		},
	})

	affectedSpellCodes := []int32{SpellCode_WarlockCorruption, SpellCode_WarlockDrainLife, SpellCode_WarlockDrainSoul}
	procChance := 0.02 * float64(warlock.Talents.Nightfall)

	core.MakePermanent(warlock.RegisterAura(core.Aura{
		Label: "Nightfall Hidden Aura",
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if slices.Contains(affectedSpellCodes, spell.SpellCode) && sim.Proc(procChance, "Nightfall") {
				shadowTranceAura.Activate(sim)
			}
		},
	}))
}

func (warlock *Warlock) applyShadowMastery() {
	if warlock.Talents.ShadowMastery == 0 {
		return
	}

	// These spells have their base damage modded instead
	// Apply Aura: Modifies Spell Effectiveness (8)
	excludedSpellCodes := []int32{SpellCode_WarlockBaneOfAgony, SpellCode_WarlockDeathCoil, SpellCode_WarlockDrainLife, SpellCode_WarlockDrainSoul}

	warlock.OnSpellRegistered(func(spell *core.Spell) {
		// Shadow Mastery applies a base damage modifier to all dots / channeled spells instead
		if spell.SpellSchool.Matches(core.SpellSchoolShadow) && isWarlockSpell(spell) && !slices.Contains(excludedSpellCodes, spell.SpellCode) {
			spell.DamageMultiplierAdditive += warlock.shadowMasteryBonus()
		}
	})
}

func (warlock *Warlock) shadowMasteryBonus() float64 {
	return .01 * float64(warlock.Talents.ShadowMastery)
}

///////////////////////////////////////////////////////////////////////////
//                            Demonology Talents
///////////////////////////////////////////////////////////////////////////

func (warlock *Warlock) applyDemonicEmbrace() {
	if warlock.Talents.DemonicEmbrace == 0 {
		return
	}

	warlock.MultiplyStat(stats.Stamina, 1+.03*float64(warlock.Talents.DemonicEmbrace))
}

func (warlock *Warlock) applyUnholyPower() {
	if warlock.Talents.UnholyPower == 0 {
		return
	}

	multiplier := 1 + 0.02*float64(warlock.Talents.UnholyPower)
	for _, pet := range warlock.BasePets {
		pet.PseudoStats.DamageDealtMultiplier *= multiplier
	}
}

func (warlock *Warlock) applyFelVitality() {
	if warlock.Talents.FelVitality == 0 {
		return
	}

	multiplier := 1 + 0.05*float64(warlock.Talents.FelVitality)
	warlock.MultiplyStat(stats.Mana, multiplier)
	for _, pet := range warlock.BasePets {
		pet.MultiplyStat(stats.Health, multiplier)
		pet.MultiplyStat(stats.Mana, multiplier)
	}
}

func (warlock *Warlock) applyMasterSummoner() {
	if warlock.Talents.MasterSummoner == 0 {
		return
	}

	castTimeReduction := time.Second * 2 * time.Duration(warlock.Talents.MasterSummoner)
	costReduction := 20 * warlock.Talents.MasterSummoner

	// Use an aura because the summon spells aren't registered by this point
	warlock.RegisterAura(core.Aura{
		Label:    "Master Summoner Hidden Aura",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range warlock.SummonDemonSpells {
				spell.DefaultCast.CastTime -= castTimeReduction
				spell.Cost.Multiplier -= costReduction
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range warlock.SummonDemonSpells {
				spell.DefaultCast.CastTime += castTimeReduction
				spell.Cost.Multiplier += costReduction
			}
		},
	})
}

func (warlock *Warlock) applyDecimation() {
	if warlock.Talents.Decimation == 0 {
		return
	}

	// TODO: Both ranks of Decimation read the same numbers, so the per point scaling is a guess
	points := float64(warlock.Talents.Decimation)
	damageBonus := 0.03 * points
	castTimeReduction := 0.2 * points

	// The damage half of the tooltip is about the two spells that trigger it, not about
	// everything the warlock casts; only the Soul Fire half carries the ten second window.
	decimationSpells := func() []*core.Spell {
		return append(append([]*core.Spell{}, warlock.ShadowBolt...), warlock.SearingPain...)
	}

	warlock.DecimationAura = warlock.RegisterAura(core.Aura{
		Label:    "Decimation",
		ActionID: core.ActionID{SpellID: 63165},
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range decimationSpells() {
				spell.DamageMultiplierAdditive += damageBonus
			}
			for _, spell := range warlock.SoulFire {
				spell.CastTimeMultiplier -= castTimeReduction
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range decimationSpells() {
				spell.DamageMultiplierAdditive -= damageBonus
			}
			for _, spell := range warlock.SoulFire {
				spell.CastTimeMultiplier += castTimeReduction
			}
		},
	})

	affectedSpellCodes := []int32{SpellCode_WarlockShadowBolt, SpellCode_WarlockSearingPain}
	core.MakePermanent(warlock.RegisterAura(core.Aura{
		Label: "Decimation Trigger",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && sim.IsExecutePhase35() && slices.Contains(affectedSpellCodes, spell.SpellCode) {
				warlock.DecimationAura.Activate(sim)
			}
		},
	}))
}

func (warlock *Warlock) applyDemonicBrand() {
	if warlock.Talents.DemonicBrand == 0 {
		return
	}

	points := float64(warlock.Talents.DemonicBrand)
	actionID := core.ActionID{SpellID: 18821}

	warlock.OnSpellRegistered(func(spell *core.Spell) {
		if spell.SpellCode == SpellCode_WarlockSearingPain {
			spell.ThreatMultiplier *= 1 - 0.17*points
		}
	})

	for _, pet := range warlock.BasePets {
		brandSpell := pet.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolShadow,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagPassiveSpell | core.SpellFlagNoOnCastComplete,

			DamageMultiplier: 1,
			ThreatMultiplier: 3,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				// TODO: Level 60 values
				spell.CalcAndDealDamage(sim, target, sim.Roll(39, 42), spell.OutcomeMagicHit)
			},
		})

		pet.DemonicBrandAura = pet.RegisterAura(core.Aura{
			Label:     "Demonic Brand",
			ActionID:  actionID,
			Duration:  time.Second * 10,
			MaxStacks: 2,
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if result.Landed() && spell.ProcMask.Matches(core.ProcMaskMelee) {
					brandSpell.Cast(sim, result.Target)
					aura.RemoveStack(sim)
				}
			},
		})
	}

	core.MakePermanent(warlock.RegisterAura(core.Aura{
		Label: "Demonic Brand Trigger",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && spell.SpellCode == SpellCode_WarlockSearingPain && warlock.ActivePet != nil {
				brandAura := warlock.ActivePet.DemonicBrandAura
				brandAura.Activate(sim)
				brandAura.SetStacks(sim, brandAura.MaxStacks)
			}
		},
	}))
}

func (warlock *Warlock) applyDemonicKnowledge() {
	if warlock.Talents.DemonicKnowledge == 0 {
		return
	}

	// The tooltip only pays the bonus out while a demon is active, so it rides on the pet
	// rather than sitting on the character sheet.
	// TODO: Beta will show whether the 33% of level is per rank or the full value
	bonus := 0.33 * float64(warlock.Talents.DemonicKnowledge) * float64(warlock.Level)

	demonicKnowledgeAura := warlock.RegisterAura(core.Aura{
		Label:    "Demonic Knowledge",
		ActionID: core.ActionID{SpellID: 35696},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warlock.AddStatDynamic(sim, stats.SpellPower, bonus)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warlock.AddStatDynamic(sim, stats.SpellPower, -bonus)
		},
	})

	for _, pet := range warlock.BasePets {
		pet.ApplyOnPetEnable(func(sim *core.Simulation) {
			demonicKnowledgeAura.Activate(sim)
		})
		pet.ApplyOnPetDisable(func(sim *core.Simulation, isSacrifice bool) {
			demonicKnowledgeAura.Deactivate(sim)
		})
	}
}

func (warlock *Warlock) applyMasterDemonologist() {
	if warlock.Talents.MasterDemonologist == 0 {
		return
	}

	points := float64(warlock.Talents.MasterDemonologist)
	damageDealtMultiplier := 1 + 0.02*points
	damageTakenMultiplier := 1 - 0.02*points

	impConfig := core.Aura{
		Label:    "Master Demonologist (Imp)",
		ActionID: core.ActionID{SpellID: 23825, Tag: 1},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] *= damageDealtMultiplier
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] /= damageDealtMultiplier
		},
	}

	voidwalkerConfig := core.Aura{
		Label:    "Master Demonologist (Voidwalker)",
		ActionID: core.ActionID{SpellID: 23825, Tag: 2},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageTakenMultiplier *= damageTakenMultiplier
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageTakenMultiplier /= damageTakenMultiplier
		},
	}

	succubusConfig := core.Aura{
		Label:    "Master Demonologist (Succubus)",
		ActionID: core.ActionID{SpellID: 23825, Tag: 3},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= damageDealtMultiplier
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] /= damageDealtMultiplier
		},
	}

	felhunterConfig := core.Aura{
		Label:    "Master Demonologist (Felhunter)",
		ActionID: core.ActionID{SpellID: 23825, Tag: 4},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for school := stats.SchoolIndexArcane; school <= stats.SchoolIndexShadow; school++ {
				aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[school] *= damageTakenMultiplier
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for school := stats.SchoolIndexArcane; school <= stats.SchoolIndexShadow; school++ {
				aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[school] /= damageTakenMultiplier
			}
		},
	}

	for _, pet := range warlock.BasePets {
		pet.ApplyOnPetEnable(func(sim *core.Simulation) {
			if warlock.MasterDemonologistAura != nil {
				warlock.MasterDemonologistAura.Deactivate(sim)
			}
		})
	}

	warlockImpAura := warlock.RegisterAura(impConfig)
	impAura := warlock.Imp.RegisterAura(impConfig)
	warlock.Imp.ApplyOnPetEnable(func(sim *core.Simulation) {
		impAura.Activate(sim)
		warlock.MasterDemonologistAura = warlockImpAura
	})
	warlock.Imp.ApplyOnPetDisable(func(sim *core.Simulation, isSacrifice bool) {
		impAura.Deactivate(sim)
	})

	warlockVoidwalkerAura := warlock.RegisterAura(voidwalkerConfig)
	voidwalkerAura := warlock.Voidwalker.RegisterAura(voidwalkerConfig)
	warlock.Voidwalker.ApplyOnPetEnable(func(sim *core.Simulation) {
		voidwalkerAura.Activate(sim)
		warlock.MasterDemonologistAura = warlockVoidwalkerAura
	})
	warlock.Voidwalker.ApplyOnPetDisable(func(sim *core.Simulation, isSacrifice bool) {
		voidwalkerAura.Deactivate(sim)
	})

	warlockSuccubusAura := warlock.RegisterAura(succubusConfig)
	succubusAura := warlock.Succubus.RegisterAura(succubusConfig)
	warlock.Succubus.ApplyOnPetEnable(func(sim *core.Simulation) {
		succubusAura.Activate(sim)
		warlock.MasterDemonologistAura = warlockSuccubusAura
	})
	warlock.Succubus.ApplyOnPetDisable(func(sim *core.Simulation, isSacrifice bool) {
		succubusAura.Deactivate(sim)
	})

	warlockFelhunterAura := warlock.RegisterAura(felhunterConfig)
	felhunterAura := warlock.Felhunter.RegisterAura(felhunterConfig)
	warlock.Felhunter.ApplyOnPetEnable(func(sim *core.Simulation) {
		felhunterAura.Activate(sim)
		warlock.MasterDemonologistAura = warlockFelhunterAura
	})
	warlock.Felhunter.ApplyOnPetDisable(func(sim *core.Simulation, isSacrifice bool) {
		felhunterAura.Deactivate(sim)
	})

	for _, pet := range warlock.BasePets {
		pet.ApplyOnPetEnable(func(sim *core.Simulation) {
			warlock.MasterDemonologistAura.Activate(sim)
		})

		pet.ApplyOnPetDisable(func(sim *core.Simulation, isSacrifice bool) {
			if warlock.MasterDemonologistAura != nil {
				warlock.MasterDemonologistAura.Deactivate(sim)
				warlock.MasterDemonologistAura = nil
			}
		})
	}
}

func (warlock *Warlock) applySoulLink() {
	if !warlock.Talents.SoulLink {
		return
	}

	actionID := core.ActionID{SpellID: 19028}
	soulLinkConfig := core.Aura{
		Label:    "Soul Link Aura",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageTakenMultiplier /= 1.3
			aura.Unit.PseudoStats.DamageDealtMultiplier *= 1.03
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier /= 1.03
			aura.Unit.PseudoStats.DamageTakenMultiplier *= 1.3
		},
	}

	warlock.SoulLinkAura = warlock.RegisterAura(soulLinkConfig)
	for _, pet := range warlock.BasePets {
		pet.SoulLinkAura = pet.RegisterAura(soulLinkConfig)

		oldOnPetDisable := pet.OnPetDisable
		pet.OnPetDisable = func(sim *core.Simulation, isSacrifice bool) {
			oldOnPetDisable(sim, isSacrifice)
			warlock.SoulLinkAura.Deactivate(sim)
			pet.SoulLinkAura.Deactivate(sim)
		}
	}

	warlock.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warlock.ActivePet != nil
		},

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.2,
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			warlock.SoulLinkAura.Activate(sim)
			warlock.ActivePet.SoulLinkAura.Activate(sim)
		},
	})
}

func (warlock *Warlock) applyDemonicSacrifice() {
	if !warlock.Talents.DemonicSacrifice {
		return
	}

	duration := time.Hour * 2

	impAura := warlock.GetOrRegisterAura(core.Aura{
		Label:    "Touch of Shadow",
		ActionID: core.ActionID{SpellID: 18791},
		Duration: duration,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= 1.15
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] /= 1.15
		},
	})

	succubusAura := warlock.GetOrRegisterAura(core.Aura{
		Label:    "Burning Wish",
		ActionID: core.ActionID{SpellID: 18789},
		Duration: duration,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] *= 1.15
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] /= 1.15
		},
	})

	var vwPa *core.PendingAction
	manaMetric := warlock.NewManaMetrics(core.ActionID{SpellID: 18792})
	voidwalkerAura := warlock.GetOrRegisterAura(core.Aura{
		Label:    "Fel Energy",
		ActionID: core.ActionID{SpellID: 18792},
		Duration: duration,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			vwPa = core.NewPeriodicAction(sim, core.PeriodicActionOptions{
				Period: time.Second * 4,
				OnAction: func(s *core.Simulation) {
					warlock.AddMana(sim, warlock.MaxMana()*0.02, manaMetric)
				},
			})
			sim.AddPendingAction(vwPa)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			vwPa.Cancel(sim)
		},
	})

	var fhPa *core.PendingAction
	healthMetric := warlock.NewHealthMetrics(core.ActionID{SpellID: 18790})
	felhunterAura := warlock.GetOrRegisterAura(core.Aura{
		Label:    "Fel Stamina",
		ActionID: core.ActionID{SpellID: 18790},
		Duration: duration,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			fhPa = core.NewPeriodicAction(sim, core.PeriodicActionOptions{
				Period: time.Second * 4,
				OnAction: func(s *core.Simulation) {
					warlock.GainHealth(sim, warlock.MaxHealth()*0.03, healthMetric)
				},
			})
			sim.AddPendingAction(fhPa)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			fhPa.Cancel(sim)
		},
	})

	warlock.DemonicSacrificeAuras = map[*WarlockPet]*core.Aura{
		warlock.Felhunter:  felhunterAura,
		warlock.Imp:        impAura,
		warlock.Succubus:   succubusAura,
		warlock.Voidwalker: voidwalkerAura,
	}
	for _, pet := range warlock.BasePets {
		pet.ApplyOnPetEnable(func(sim *core.Simulation) {
			// With Demonic Pact only resummoning the sacrificed demon ends the effect
			if warlock.Talents.DemonicPact && pet != warlock.SacrificedPet {
				return
			}

			for _, dsAura := range warlock.DemonicSacrificeAuras {
				dsAura.Deactivate(sim)
			}
			warlock.SacrificedPet = nil
		})
	}

	warlock.GetOrRegisterSpell(core.SpellConfig{
		SpellCode:   SpellCode_WarlockDemonicSacrifice,
		ActionID:    core.ActionID{SpellID: 18788},
		SpellSchool: core.SpellSchoolShadow,
		Flags:       core.SpellFlagAPL,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warlock.ActivePet != nil
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, dsAura := range warlock.DemonicSacrificeAuras {
				dsAura.Deactivate(sim)
			}
			warlock.DemonicSacrificeAuras[warlock.ActivePet].Activate(sim)

			warlock.SacrificedPet = warlock.ActivePet
			warlock.changeActivePet(sim, nil, true)
		},
	})
}

///////////////////////////////////////////////////////////////////////////
//                            Destruction Talents
///////////////////////////////////////////////////////////////////////////

func (warlock *Warlock) applyImprovedShadowBolt() {
	if warlock.Talents.ImprovedShadowBolt == 0 {
		return
	}

	damageMultiplier := 1 + 0.04*float64(warlock.Talents.ImprovedShadowBolt)

	warlock.ImprovedShadowBoltAuras = warlock.NewEnemyAuraArray(func(unit *core.Unit) *core.Aura {
		return unit.RegisterAura(core.Aura{
			Label:    "Improved Shadow Bolt-" + warlock.Label,
			ActionID: core.ActionID{SpellID: 17800},
			Duration: time.Second * 12,
		})
	})

	// Only the warlock's own shadow damage benefits, so this can't go through the target's school multiplier
	for _, target := range warlock.Env.Encounter.TargetUnits {
		target.AddDynamicDamageTakenModifier(func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Unit == &warlock.Unit && spell.SpellSchool.Matches(core.SpellSchoolShadow) && warlock.ImprovedShadowBoltAuras.Get(result.Target).IsActive() {
				result.Damage *= damageMultiplier
			}
		})
	}

	core.MakePermanent(warlock.RegisterAura(core.Aura{
		Label: "ISB Trigger",
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range warlock.ShadowBolt {
				spell.RelatedAuras = []core.AuraArray{warlock.ImprovedShadowBoltAuras}
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.DidCrit() && spell.SpellCode == SpellCode_WarlockShadowBolt {
				warlock.ImprovedShadowBoltAuras.Get(result.Target).Activate(sim)
			}
		},
	}))
}

func (warlock *Warlock) applyCataclysm() {
	if warlock.Talents.Cataclysm == 0 {
		return
	}

	warlock.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Flags.Matches(WarlockFlagDestruction) && spell.Cost != nil {
			spell.Cost.Multiplier -= 3 * warlock.Talents.Cataclysm
		}
	})
}

func (warlock *Warlock) applyBane() {
	if warlock.Talents.Bane == 0 {
		return
	}

	points := time.Duration(warlock.Talents.Bane)
	warlock.OnSpellRegistered(func(spell *core.Spell) {
		if spell.SpellCode == SpellCode_WarlockShadowBolt || spell.SpellCode == SpellCode_WarlockImmolate || spell.SpellCode == SpellCode_WarlockIncinerate {
			spell.DefaultCast.CastTime -= time.Millisecond * 100 * points
		} else if spell.SpellCode == SpellCode_WarlockSoulFire {
			spell.DefaultCast.CastTime -= time.Millisecond * 400 * points
		}
	})
}

func (warlock *Warlock) applyRuin() {
	if warlock.Talents.Ruin == 0 {
		return
	}

	bonus := 0.2 * float64(warlock.Talents.Ruin)
	warlock.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Flags.Matches(WarlockFlagDestruction) {
			spell.CritDamageBonus += bonus
		}
	})
}

func (warlock *Warlock) applyAgonizingFlames() {
	if warlock.Talents.AgonizingFlames == 0 {
		return
	}

	points := float64(warlock.Talents.AgonizingFlames)
	warlock.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Flags.Matches(WarlockFlagDestruction) {
			spell.DamageMultiplierAdditive += 0.03 * points
		}
		if spell.SpellCode == SpellCode_WarlockSearingPain {
			spell.BonusCritRating += 3 * points * core.SpellCritRatingPerCritChance
		}
	})
}

func (warlock *Warlock) applyFireAndBrimstone() {
	if warlock.Talents.FireAndBrimstone == 0 {
		return
	}

	points := float64(warlock.Talents.FireAndBrimstone)
	warlock.OnSpellRegistered(func(spell *core.Spell) {
		if spell.SpellCode == SpellCode_WarlockConflagrate {
			spell.BonusCritRating += 8 * points * core.SpellCritRatingPerCritChance
		}
	})
}

func (warlock *Warlock) applyShadowAndFlame() {
	if warlock.Talents.ShadowAndFlame == 0 {
		return
	}

	multiplier := 1 + 0.02*float64(warlock.Talents.ShadowAndFlame)

	shadowAura := warlock.RegisterAura(core.Aura{
		Label:    "Shadow and Flame (Shadow)",
		ActionID: core.ActionID{SpellID: 30288, Tag: 1},
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= multiplier
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] /= multiplier
		},
	})

	fireAura := warlock.RegisterAura(core.Aura{
		Label:    "Shadow and Flame (Fire)",
		ActionID: core.ActionID{SpellID: 30288, Tag: 2},
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] *= multiplier
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] /= multiplier
		},
	})

	core.MakePermanent(warlock.RegisterAura(core.Aura{
		Label: "Shadow and Flame Trigger",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}

			switch spell.SpellCode {
			case SpellCode_WarlockConflagrate:
				shadowAura.Activate(sim)
			case SpellCode_WarlockShadowburn:
				fireAura.Activate(sim)
			}
		},
	}))
}
