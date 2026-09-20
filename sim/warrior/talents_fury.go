package warrior

import (
	"time"

	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/stats"
)

// TODO: Manual review needed -- this was modelled during the Forever port, not carried
// over unchanged, so its numbers and shape want checking against the client.
func (warrior *Warrior) registerFuryTalents() {
	// Tier 1
	// Booming Voice implemented in shouts.go
	warrior.registerCruelty()

	// Tier 2
	// Improved Demoralizing Shout implemented in demoralizing_shout.go
	warrior.registerUnbridledWrath()

	// Tier 3
	// Improved Cleave implemented in heroic_strike_cleave.go
	// Piercing Howl not implemented
	// Blood Craze not implemented
	// Commanding Presence implemented in shouts.go

	// Tier 4
	warrior.registerDualWieldSpecialization()
	warrior.registerImprovedExecute()
	warrior.registerEnrage()

	// Tier 5
	warrior.registerImprovedSlam()
	warrior.registerSweepingStrikes()

	// Tier 6
	warrior.registerImprovedBerserkerRage()
	warrior.registerFlurry()

	// Tier 7
	warrior.registerPrecision()
	warrior.registerBloodthirst()

	// Tier 8

	// Tier 9

	// Forever additions, not yet implemented.
	warrior.registerIronWill()
	warrior.registerPiercingHowl()
	warrior.registerBloodCraze()
	warrior.registerBoundlessRage()
	warrior.registerRagingBlows()
}

func (warrior *Warrior) registerCruelty() {
	if warrior.Talents.Cruelty == 0 {
		return
	}

	warrior.AddStat(stats.PhysicalCritPercent, spellData.Cruelty.ValueAt(warrior.Talents.Cruelty))
}

func (warrior *Warrior) registerUnbridledWrath() {
	if warrior.Talents.UnbridledWrath == 0 {
		return
	}

	rageMetrics := warrior.NewRageMetrics(core.ActionID{SpellID: 13002})

	warrior.MakeProcTriggerAura(core.ProcTrigger{
		Name:               "Unbridled Wrath",
		DPM:                warrior.NewStaticLegacyPPMManager(3*float64(warrior.Talents.UnbridledWrath), core.ProcMaskMeleeWhiteHit),
		RequireDamageDealt: true,
		Outcome:            core.OutcomeLanded,
		Callback:           core.CallbackOnSpellHitDealt,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			warrior.AddRage(sim, 1, rageMetrics)
		},
	})
}

// TODO: Manual review needed -- this was modelled during the Forever port, not carried
// over unchanged, so its numbers and shape want checking against the client.
func (warrior *Warrior) registerDualWieldSpecialization() {
	if warrior.Talents.DualWieldSpecialization == 0 {
		return
	}

	warrior.AddStaticMod(core.SpellModConfig{
		ProcMask:   core.ProcMaskMeleeOH,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: spellData.DualWieldSpecialization.Effect(shared.A_MOD_OFFHAND_DAMAGE_PCT, 0).FractionAt(warrior.Talents.DualWieldSpecialization),
	})
}

func (warrior *Warrior) registerImprovedExecute() {
	if warrior.Talents.ImprovedExecute == 0 {
		return
	}

	rageCostReduction := []int32{0, 2, 5}[warrior.Talents.ImprovedExecute]

	warrior.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskExecute,
		Kind:      core.SpellMod_PowerCost_Flat,
		IntValue:  -rageCostReduction,
	})
}

func (warrior *Warrior) registerEnrage() {
	if warrior.Talents.Enrage == 0 {
		return
	}

	warrior.EnrageAura = warrior.GetOrRegisterAura(core.Aura{
		Label:     "Enrage",
		ActionID:  core.ActionID{SpellID: 13048},
		Duration:  time.Second * 12,
		MaxStacks: 12,
	}).AttachSpellMod(core.SpellModConfig{
		School:     core.SpellSchoolPhysical,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 0.05 * float64(warrior.Talents.Enrage),
	}).AttachProcTrigger(core.ProcTrigger{
		Name:               "Enrage - Spend",
		TriggerImmediately: true,
		ProcMask:           core.ProcMaskMelee,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			warrior.EnrageAura.RemoveStack(sim)
		},
	})

	warrior.EnrageAura.NewExclusiveEffect("Enrage", true, core.ExclusiveEffect{Priority: 5 * float64(warrior.Talents.Enrage)})

	warrior.MakeProcTriggerAura(core.ProcTrigger{
		Name:     "Enrage - Trigger",
		ProcMask: core.ProcMaskMelee,
		Outcome:  core.OutcomeCrit | core.OutcomeSuppressedCrit,
		Callback: core.CallbackOnSpellHitTaken,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			warrior.EnrageAura.Activate(sim)
			warrior.EnrageAura.SetStacks(sim, 12)
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

func (warrior *Warrior) registerImprovedBerserkerRage() {
	if warrior.Talents.ImprovedBerserkerRage == 0 {
		return
	}

	core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label:    "Improved Berserker Rage",
		ActionID: core.ActionID{SpellID: 20500}.WithTag(warrior.Talents.ImprovedBerserkerRage),
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.BerserkerRageRageGain += 5 * float64(warrior.Talents.ImprovedBerserkerRage)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.BerserkerRageRageGain -= 5 * float64(warrior.Talents.ImprovedBerserkerRage)
		},
	}))
}

func (warrior *Warrior) registerFlurry() {
	if warrior.Talents.Flurry == 0 {
		return
	}

	flurryAura := warrior.RegisterAura(core.Aura{
		Label:     "Flurry",
		ActionID:  core.ActionID{SpellID: 12970},
		Duration:  15 * time.Second,
		MaxStacks: 3,
	}).AttachMultiplyMeleeSpeed(1 + 0.05*float64(warrior.Talents.Flurry))

	warrior.MakeProcTriggerAura(core.ProcTrigger{
		Name:               "Flurry - Trigger",
		ActionID:           core.ActionID{SpellID: 12974},
		ProcMask:           core.ProcMaskMelee,
		CanProcFromProcs:   true, // 12319, 12971-12974 carry the bit.
		TriggerImmediately: true,
		Callback:           core.CallbackOnSpellHitDealt,
		Outcome:            core.OutcomeLanded,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(SpellMaskWhirlwindOh) {
				return
			}

			if result.Outcome.Matches(core.OutcomeCrit) {
				flurryAura.Activate(sim)
				flurryAura.SetStacks(sim, 3)
				return
			}

			if flurryAura.IsActive() && spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
				flurryAura.RemoveStack(sim)
			}
		},
	})
}

// TODO: Manual review needed -- this was modelled during the Forever port, not carried
// over unchanged, so its numbers and shape want checking against the client.
func (warrior *Warrior) registerPrecision() {
	if warrior.Talents.Precision == 0 {
		return
	}

	// The spell-hit effect carries the same ladder; only melee hit is taken here.
	warrior.AddStat(stats.PhysicalHitPercent, spellData.Precision.Effect(shared.A_MOD_HIT_CHANCE, 0).ValueAt(warrior.Talents.Precision))
}

func (warrior *Warrior) registerBloodthirst() {
	if !warrior.Talents.Bloodthirst {
		return
	}

	actionID := core.ActionID{SpellID: 30335}

	warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskBloodthirst,
		MaxRange:       core.MaxMeleeRange,

		RageCost: core.RageCostOptions{
			Cost:   30,
			Refund: 0.8,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.MeleeAttackPower(target) * 0.45
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
}

// registerIronWill implements Iron Will, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (warrior *Warrior) registerIronWill() {
	if warrior.Talents.IronWill == 0 {
		return
	}
}

// registerPiercingHowl implements Piercing Howl, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (warrior *Warrior) registerPiercingHowl() {
	if !warrior.Talents.PiercingHowl {
		return
	}
}

// registerBloodCraze implements Blood Craze, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (warrior *Warrior) registerBloodCraze() {
	if warrior.Talents.BloodCraze == 0 {
		return
	}
}

// registerBoundlessRage implements Boundless Rage, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (warrior *Warrior) registerBoundlessRage() {
	if warrior.Talents.BoundlessRage == 0 {
		return
	}
}

// registerRagingBlows implements Raging Blows, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (warrior *Warrior) registerRagingBlows() {
	if !warrior.Talents.RagingBlows {
		return
	}
}
