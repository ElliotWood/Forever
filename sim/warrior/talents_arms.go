package warrior

import (
	"time"

	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
)

// TODO: Manual review needed -- this was modelled during the Forever port, not carried
// over unchanged, so its numbers and shape want checking against the client.
func (warrior *Warrior) registerArmsTalents() {
	// Tier 1
	warrior.registerImprovedHeroicStrike()
	warrior.registerDeflection()
	warrior.registerImprovedRend()

	// Tier 2
	warrior.registerImprovedCharge()
	warrior.registerImprovedTacticalMastery()
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

	warrior.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskHeroicStrike,
		Kind:      core.SpellMod_PowerCost_Flat,
		IntValue:  -warrior.Talents.ImprovedHeroicStrike,
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

func (warrior *Warrior) registerImprovedCharge() {
	if warrior.Talents.ImprovedCharge == 0 {
		return
	}

	warrior.ChargeRageGain += 3.0 + float64(warrior.Talents.ImprovedCharge)
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

	warrior.DeepWounds = warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 12867},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskEmpty,
		ClassSpellMask: SpellMaskDeepWounds,
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagIgnoreResists | core.SpellFlagProc, // 12867 lacks Not a Proc.

		// Deep Wounds (12867) has no SpellCategories row in the client DB. It's a bleed that
		// snapshots on proc; the application uses OutcomeAlwaysHitNoHitCounter and the DoT ticks
		// with OutcomeTick, so it never rolls a crit and DefenseType is intentionally left unset.
		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "DeepWounds",
			},
			NumberOfTicks: 6,
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				baseDamage := warrior.AutoAttacks.MH().CalculateAverageWeaponDamage(dot.Spell.MeleeAttackPower(target))
				dot.SnapshotPhysical(target, baseDamage/float64(dot.HastedTickCount())*0.2*float64(warrior.Talents.DeepWounds))
			},

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHitNoHitCounter)
			spell.Dot(target).Deactivate(sim)
			spell.Dot(target).Apply(sim)
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

// TODO: registerImprovedTacticalMastery models nothing yet; spellData.ImprovedTacticalMastery carries the ranks.
func (warrior *Warrior) registerImprovedTacticalMastery() {
	if warrior.Talents.ImprovedTacticalMastery == 0 {
		return
	}
}

// TODO: registerSpearingStrike models nothing yet; the client has no ladder for it, so the tooltip is the source.
func (warrior *Warrior) registerSpearingStrike() {
	if !warrior.Talents.SpearingStrike {
		return
	}
}

// TODO: registerBloodthrill models nothing yet; spellData.Bloodthrill carries the ranks.
func (warrior *Warrior) registerBloodthrill() {
	if warrior.Talents.Bloodthrill == 0 {
		return
	}
}

// TODO: registerWeaponmaster models nothing yet; spellData.Weaponmaster carries the ranks.
func (warrior *Warrior) registerWeaponmaster() {
	if warrior.Talents.Weaponmaster == 0 {
		return
	}
}

// TODO: registerImprovedHamstring models nothing yet; spellData.ImprovedHamstring carries the ranks.
func (warrior *Warrior) registerImprovedHamstring() {
	if warrior.Talents.ImprovedHamstring == 0 {
		return
	}
}

func (warrior *Warrior) registerImprovedSlam() {
	if warrior.Talents.ImprovedSlam == 0 {
		return
	}

	warrior.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskSlam,
		Kind:      core.SpellMod_CastTime_Flat,
		TimeValue: -time.Millisecond * time.Duration(500*warrior.Talents.ImprovedSlam),
	})
}

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
		Duration:           time.Second * 10,
		Callback:           core.CallbackOnSpellHitDealt,
		ProcMask:           core.ProcMaskMelee,
		Outcome:            core.OutcomeLanded,
		TriggerImmediately: true,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if warrior.Env.ActiveTargetCount() < 2 || warrior.SweepingStrikesAura.GetStacks() == 0 || result.PostOutcomeDamage <= 0 || !spell.ProcMask.Matches(core.ProcMaskMelee) {
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
	warrior.SweepingStrikesAura.MaxStacks = 10

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
			return warrior.StanceMatches(BattleStance|BerserkerStance) || sim.ActiveTargetCount() > 1
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			spell.RelatedSelfBuff.Activate(sim)
			warrior.SweepingStrikesAura.SetStacks(sim, 10)
		},

		RelatedSelfBuff: warrior.SweepingStrikesAura,
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: ssCD,
		Type:  core.CooldownTypeDPS,
	})
}
