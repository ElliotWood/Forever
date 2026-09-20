package warrior

import (
	"time"

	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

func (warrior *Warrior) registerArmsTalents() {
	// Tier 1
	warrior.registerImprovedHeroicStrike()
	warrior.registerDeflection()
	warrior.registerImprovedRend()

	// Tier 2
	// Improved Charge: charge.go
	// Improved Tactical Mastery: stances.go
	warrior.registerImprovedOverpower()

	// Tier 3
	warrior.registerAngerManagement()
	warrior.registerDeepWounds()

	// Tier 4
	warrior.registerSpearingStrike()
	warrior.registerTwoHandedWeaponSpecialization()
	warrior.registerImpale()

	// Tier 5
	warrior.registerBloodthrill()
	warrior.registerSweepingStrikes()
	warrior.registerWeaponmaster()

	// Tier 6
	warrior.registerImprovedSlam()
	warrior.registerImprovedHamstring()

	// Tier 7
	warrior.registerMortalStrike()
}

/*
 * Arms
 */
func (warrior *Warrior) registerImprovedHeroicStrike() {
	if warrior.Talents.ImprovedHeroicStrike == 0 {
		return
	}

	// The client states the cost on the 0-1000 rage bar, so the ladder is -10/-20/-30.
	warrior.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskHeroicStrike,
		Kind:      core.SpellMod_PowerCost_Flat,
		IntValue:  int32(spellData.ImprovedHeroicStrike.TenthsAt(warrior.Talents.ImprovedHeroicStrike)),
	})
}
func (warrior *Warrior) registerDeflection() {
	if warrior.Talents.Deflection == 0 {
		return
	}

	warrior.PseudoStats.BaseParryChance += spellData.Deflection.FractionAt(warrior.Talents.Deflection)
}

func (warrior *Warrior) registerImprovedRend() {
	if warrior.Talents.ImprovedRend == 0 {
		return
	}

	warrior.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskRend,
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: spellData.ImprovedRend.FractionAt(warrior.Talents.ImprovedRend),
	})
}

func (warrior *Warrior) registerImprovedOverpower() {
	if warrior.Talents.ImprovedOverpower == 0 {
		return
	}

	core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label:    "Improved Overpower",
		ActionID: core.ActionID{SpellID: 12963}.WithTag(warrior.Talents.ImprovedOverpower),
	})).AttachSpellMod(core.SpellModConfig{
		ClassMask:  SpellMaskOverpower,
		Kind:       core.SpellMod_BonusCrit_Percent,
		FloatValue: spellData.ImprovedOverpower.ValueAt(warrior.Talents.ImprovedOverpower),
	})
}

func (warrior *Warrior) registerAngerManagement() {
	if !warrior.Talents.AngerManagement {
		return
	}

	rageMetrics := warrior.NewRageMetrics(core.ActionID{SpellID: 12296})

	// TODO: Manual review needed -- Anger Management has no generated table; the client states
	// 1 Rage every 3 seconds in combat (12296).
	warrior.RegisterResetEffect(func(sim *core.Simulation) {
		core.StartPeriodicAction(sim, core.PeriodicActionOptions{
			Period: time.Second * 3,
			OnAction: func(sim *core.Simulation) {
				if sim.CurrentTime > 0 {
					warrior.AddRage(sim, 1, rageMetrics)
				}
			},
		})
	})
}

func (warrior *Warrior) registerDeepWounds() {
	if warrior.Talents.DeepWounds == 0 {
		return
	}

	share := spellData.DeepWounds.FractionAt(warrior.Talents.DeepWounds)
	warrior.DeepWounds = warrior.RegisterSpell(core.SpellConfig{
		// The bleed the talent (12834) reaches through 12162.
		ActionID:       core.ActionID{SpellID: 412609},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskEmpty,
		ClassSpellMask: SpellMaskDeepWounds,
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagIgnoreResists | core.SpellFlagProc, // 12162 and 412609 lack Not a Proc.

		// 12162 and 412609 state DefenseType 0. It's a bleed that snapshots on proc; the
		// application uses OutcomeAlwaysHitNoHitCounter and the DoT ticks with OutcomeTick, so it
		// never rolls a crit and DefenseType is intentionally left unset.
		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "DeepWounds",
			},
			// TODO: Manual review needed -- 412609 has no generated table; the client states a
			// 3 second period over a 12 second duration.
			NumberOfTicks: 4,
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				baseDamage := warrior.AutoAttacks.MH().CalculateAverageWeaponDamage(dot.Spell.MeleeAttackPower(target))
				dot.SnapshotPhysical(target, baseDamage/float64(dot.HastedTickCount())*share)
			},

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHitNoHitCounter)
			dot := spell.Dot(target)
			dot.Deactivate(sim)
			dot.Apply(sim)
		},
	})

	warrior.MakeProcTriggerAura(core.ProcTrigger{
		Name:               "Deep Wounds - Trigger",
		TriggerImmediately: true,
		ProcMaskExclude:    core.ProcMaskEmpty,
		Outcome:            core.OutcomeCrit,
		Callback:           core.CallbackOnSpellHitDealt,
		ExtraCondition: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) bool {
			return spell.SpellSchool.Matches(core.SpellSchoolPhysical)
		},
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			warrior.DeepWounds.Cast(sim, result.Target)
		},
	})

}

func (warrior *Warrior) registerTwoHandedWeaponSpecialization() {
	if warrior.Talents.TwoHandedWeaponSpecialization == 0 {
		return
	}

	weaponMod := warrior.AddDynamicMod(core.SpellModConfig{
		ClassMask:  SpellMaskDirectDamageSpells,
		School:     core.SpellSchoolPhysical,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: spellData.TwoHandedWeaponSpecialization.Effect(shared.A_MOD_DAMAGE_PERCENT_DONE, 1).FractionAt(warrior.Talents.TwoHandedWeaponSpecialization),
	})

	if warrior.GetMainHandType() == proto.HandType_HandTypeTwoHand {
		weaponMod.Activate()
	}

	warrior.RegisterItemSwapCallback(core.AllMeleeWeaponSlots(), func(sim *core.Simulation, slot proto.ItemSlot) {
		if warrior.GetMainHandType() == proto.HandType_HandTypeTwoHand {
			weaponMod.Activate()
		} else {
			weaponMod.Deactivate()
		}
	})
}

func (warrior *Warrior) registerImpale() {
	if warrior.Talents.Impale == 0 {
		return
	}

	warrior.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskDamageSpells,
		Kind:       core.SpellMod_CritMultiplier_Flat,
		FloatValue: spellData.Impale.FractionAt(warrior.Talents.Impale),
	})
}

var mortalStrikeRank = spellData.MortalStrike.HighestRank()
var mortalStrikeBaseDamage, _ = mortalStrikeRank.Direct.Range()

func (warrior *Warrior) registerMortalStrike() {
	if !warrior.Talents.MortalStrike {
		return
	}

	warrior.MortalStrike = warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: mortalStrikeRank.SpellID},
		SpellSchool:    mortalStrikeRank.SpellSchool,
		DefenseType:    mortalStrikeRank.DefenseType,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagAPL | core.SpellFlagMeleeMetrics,
		ClassSpellMask: SpellMaskMortalStrike,
		MaxRange:       core.MaxMeleeRange,

		RageCost: core.RageCostOptions{
			Cost:   mortalStrikeRank.Cost,
			Refund: 0.8,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: mortalStrikeRank.GCD,
			},
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: mortalStrikeRank.Cooldown,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := mortalStrikeBaseDamage + spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower(target))
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
}

// TODO: Manual review needed -- Spearing Strike has no generated table. The client states 15 Rage,
// a 20 second cooldown, a 1.5 second global cooldown, melee range, 40% of normalized weapon damage
// and triple that against Giants and Dragonkin (1310222).
func (warrior *Warrior) registerSpearingStrike() {
	if !warrior.Talents.SpearingStrike {
		return
	}

	warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 1310222},
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagAPL | core.SpellFlagMeleeMetrics,
		ClassSpellMask: SpellMaskSpearingStrike,
		MaxRange:       core.MaxMeleeRange,

		RageCost: core.RageCostOptions{
			Cost:   15,
			Refund: 0.8,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Second * 20,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0.4 * spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower(target))
			if target.MobType == proto.MobType_MobTypeGiant || target.MobType == proto.MobType_MobTypeDragonkin {
				baseDamage *= 3
			}

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
}

func (warrior *Warrior) registerBloodthrill() {
	if warrior.Talents.Bloodthrill == 0 {
		return
	}

	// The proc (1289681) makes Overpower usable for 6 s; the cast consumes it like a dodge would.
	warrior.MakeProcTriggerAura(core.ProcTrigger{
		Name:       "Bloodthrill - Trigger",
		ActionID:   core.ActionID{SpellID: 1289682},
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeWhiteHit,
		Outcome:    core.OutcomeLanded,
		ProcChance: spellData.Bloodthrill.FractionAt(warrior.Talents.Bloodthrill),
		ExtraCondition: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) bool {
			return warrior.Rend.Dot(result.Target).IsActive()
		},
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			warrior.OverpowerAura.Activate(sim)
			warrior.OverpowerAura.UpdateExpires(sim.CurrentTime + time.Second*6)
		},
	})
}

func (warrior *Warrior) registerWeaponmaster() {
	if warrior.Talents.Weaponmaster == 0 {
		return
	}

	rank := warrior.Talents.Weaponmaster
	actionID := core.ActionID{SpellID: 1290261}

	// The three branches share A_DUMMY and misc 0, so each is named by its effect index. Which
	// branch applies follows the main hand, re-read on a weapon swap.
	mainHandIs := func(weaponTypes ...proto.WeaponType) bool {
		return warrior.GetProcMaskForTypes(weaponTypes...).Matches(core.ProcMaskMeleeMH)
	}
	var critOn, armorIgnoreOn, swordOn bool
	readMainHand := func() {
		critOn = mainHandIs(proto.WeaponType_WeaponTypeAxe, proto.WeaponType_WeaponTypePolearm)
		armorIgnoreOn = mainHandIs(proto.WeaponType_WeaponTypeMace, proto.WeaponType_WeaponTypeStaff)
		swordOn = mainHandIs(proto.WeaponType_WeaponTypeSword)
	}
	readMainHand()

	critAura := warrior.RegisterAura(core.Aura{
		Label:    "Weaponmaster (Axe/Polearm)",
		ActionID: actionID.WithTag(1),
		Duration: core.NeverExpires,
	}).AttachStatBuff(stats.PhysicalCritPercent, spellData.Weaponmaster.EffectAt(0).ValueAt(rank))
	if critOn {
		core.MakePermanent(critAura)
	}

	// The attack tables exist only once the environment is built, so the factor is written on reset.
	armorIgnore := spellData.Weaponmaster.EffectAt(1).FractionAt(rank)
	applyArmorIgnore := func() {
		for _, attackTable := range warrior.AttackTables {
			attackTable.ArmorIgnoreFactor = core.TernaryFloat64(armorIgnoreOn, armorIgnore, 0)
		}
	}
	warrior.RegisterResetEffect(func(sim *core.Simulation) {
		applyArmorIgnore()
	})

	var extraAttack *core.Spell
	warrior.MakeProcTriggerAura(core.ProcTrigger{
		Name:               "Weaponmaster (Sword)",
		ActionID:           actionID.WithTag(3),
		Callback:           core.CallbackOnSpellHitDealt,
		ProcMask:           core.ProcMaskMelee,
		Outcome:            core.OutcomeLanded,
		ProcChance:         spellData.Weaponmaster.EffectAt(2).FractionAt(rank),
		TriggerImmediately: true,
		ExtraCondition: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) bool {
			// An extra attack does not give another one.
			return swordOn && spell != extraAttack
		},
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			warrior.AutoAttacks.MaybeReplaceMHSwing(sim, extraAttack).Cast(sim, result.Target)
		},
	}).ApplyOnInit(func(aura *core.Aura, sim *core.Simulation) {
		config := *warrior.AutoAttacks.MHConfig()
		config.ActionID = config.ActionID.WithTag(actionID.SpellID)
		extraAttack = warrior.GetOrRegisterSpell(config)
	})

	warrior.RegisterItemSwapCallback(core.AllMeleeWeaponSlots(), func(sim *core.Simulation, slot proto.ItemSlot) {
		readMainHand()
		if critOn {
			critAura.Activate(sim)
		} else {
			critAura.Deactivate(sim)
		}
		applyArmorIgnore()
	})
}

func (warrior *Warrior) registerImprovedHamstring() {
	if warrior.Talents.ImprovedHamstring == 0 {
		return
	}

	// TODO: Manual review needed -- 23694 has no generated table; the client states a 5 second
	// immobilize, which a stationary sim target does not feel.
	immobilizeAuras := warrior.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    "Improved Hamstring-" + warrior.Label,
			ActionID: core.ActionID{SpellID: 23694},
			Duration: time.Second * 5,
		})
	})

	warrior.MakeProcTriggerAura(core.ProcTrigger{
		Name:           "Improved Hamstring - Trigger",
		ActionID:       core.ActionID{SpellID: 12289},
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: SpellMaskHamstring,
		Outcome:        core.OutcomeLanded,
		ProcChance:     spellData.ImprovedHamstring.FractionAt(warrior.Talents.ImprovedHamstring),
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			immobilizeAuras.Get(result.Target).Activate(sim)
		},
	})
}

func (warrior *Warrior) registerImprovedSlam() {
	if warrior.Talents.ImprovedSlam == 0 {
		return
	}

	warrior.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskSlam,
		Kind:      core.SpellMod_CastTime_Flat,
		TimeValue: time.Millisecond * time.Duration(spellData.ImprovedSlam.Effect(shared.A_ADD_FLAT_MODIFIER, shared.SPELLMOD_CASTING_TIME).ValueAt(warrior.Talents.ImprovedSlam)),
	})

	warrior.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskSlam,
		Kind:      core.SpellMod_GlobalCooldown_Flat,
		TimeValue: time.Millisecond * time.Duration(spellData.ImprovedSlam.Effect(shared.A_ADD_FLAT_MODIFIER, shared.SPELLMOD_GLOBAL_COOLDOWN).ValueAt(warrior.Talents.ImprovedSlam)),
	})
}

// TODO: Manual review needed -- Sweeping Strikes has no generated table; the client states 30 Rage,
// a 30 second cooldown, 5 charges over 20 seconds and Battle Stance only (12292).
const sweepingStrikesCharges = 5

func (warrior *Warrior) registerSweepingStrikes() {
	if !warrior.Talents.SweepingStrikes {
		return
	}

	actionID := core.ActionID{SpellID: 12723}

	var copyDamage float64
	hitSpell := warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: SpellMaskSweepingStrikesHit,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeSpecial,
		Flags:          core.SpellFlagIgnoreModifiers | core.SpellFlagMeleeMetrics | core.SpellFlagPassiveSpell | core.SpellFlagNoOnCastComplete,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, copyDamage, spell.OutcomeAlwaysHit)
		},
	})

	warrior.SweepingStrikesNormalizedAttack = warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID.WithTag(1), // Real SpellID: 26654
		ClassSpellMask: SpellMaskSweepingStrikesNormalizedHit,
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagPassiveSpell | core.SpellFlagNoOnCastComplete,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower(target))
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeAlwaysHit)
		},
	})

	warrior.SweepingStrikesAura = warrior.MakeProcTriggerAura(core.ProcTrigger{
		Name:               "Sweeping Strikes",
		ActionID:           actionID,
		MetricsActionID:    actionID,
		Duration:           time.Second * 20,
		Callback:           core.CallbackOnSpellHitDealt,
		ProcMask:           core.ProcMaskMelee,
		Outcome:            core.OutcomeLanded,
		TriggerImmediately: true,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if warrior.Env.ActiveTargetCount() < 2 || warrior.SweepingStrikesAura.GetStacks() == 0 || result.PostOutcomeDamage <= 0 {
				return
			}

			if spell.Matches(SpellMaskSweepingStrikesHit | SpellMaskSweepingStrikesNormalizedHit | SpellMaskThunderClap | SpellMaskWhirlwind | SpellMaskWhirlwindOh) {
				return
			}

			nextTarget := warrior.Env.NextActiveTargetUnit(result.Target)
			if spell.Matches(SpellMaskExecute) && sim.IsExecutePhase20() {
				warrior.SweepingStrikesNormalizedAttack.Cast(sim, nextTarget)
			} else {
				copyDamage = result.Damage / result.ArmorAndResistanceMultiplier
				hitSpell.Cast(sim, nextTarget)
			}

			warrior.SweepingStrikesAura.RemoveStack(sim)
		},
	})
	warrior.SweepingStrikesAura.MaxStacks = sweepingStrikesCharges

	ssCD := warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: SpellMaskSweepingStrikes,
		SpellSchool:    core.SpellSchoolPhysical,

		RageCost: core.RageCostOptions{
			Cost: 30,
		},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Second * 30,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			// Sweeping Strikes (12292) is usable in Battle Stance only.
			return warrior.StanceMatches(BattleStance)
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			spell.RelatedSelfBuff.Activate(sim)
			warrior.SweepingStrikesAura.SetStacks(sim, sweepingStrikesCharges)
		},

		RelatedSelfBuff: warrior.SweepingStrikesAura,
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: ssCD,
		Type:  core.CooldownTypeDPS,
	})
}
