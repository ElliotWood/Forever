package warrior

import (
	"time"

	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/stats"
)

// TODO: Manual review needed -- this was modelled during the Forever port, not carried
// over unchanged, so its numbers and shape want checking against the client.
func (war *Warrior) registerFuryTalents() {
	// Tier 1
	// Booming Voice implemented in shouts.go
	war.registerCruelty()

	// Tier 2
	// Improved Demoralizing Shout implemented in demoralizing_shout.go
	war.registerUnbridledWrath()

	// Tier 3
	// Improved Cleave implemented in heroic_strike_cleave.go
	// Piercing Howl not implemented
	// Blood Craze not implemented
	// Commanding Presence implemented in shouts.go

	// Tier 4
	war.registerDualWieldSpecialization()
	war.registerImprovedExecute()
	war.registerEnrage()

	// Tier 5
	war.registerImprovedSlam()
	war.registerSweepingStrikes()

	// Tier 6
	war.registerImprovedBerserkerRage()
	war.registerFlurry()

	// Tier 7
	war.registerPrecision()
	war.registerBloodthirst()

	// Tier 8

	// Tier 9

	// Forever additions, not yet implemented.
	war.registerIronWill()
	war.registerPiercingHowl()
	war.registerBloodCraze()
	war.registerBoundlessRage()
	war.registerRagingBlows()
}

func (war *Warrior) registerCruelty() {
	if war.Talents.Cruelty == 0 {
		return
	}

	war.AddStat(stats.PhysicalCritPercent, spellData.Cruelty.ValueAt(war.Talents.Cruelty))
}

func (war *Warrior) registerUnbridledWrath() {
	if war.Talents.UnbridledWrath == 0 {
		return
	}

	rageMetrics := war.NewRageMetrics(core.ActionID{SpellID: 13002})

	war.MakeProcTriggerAura(core.ProcTrigger{
		Name:               "Unbridled Wrath",
		DPM:                war.NewStaticLegacyPPMManager(3*float64(war.Talents.UnbridledWrath), core.ProcMaskMeleeWhiteHit),
		RequireDamageDealt: true,
		Outcome:            core.OutcomeLanded,
		Callback:           core.CallbackOnSpellHitDealt,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			war.AddRage(sim, 1, rageMetrics)
		},
	})
}

// TODO: Manual review needed -- this was modelled during the Forever port, not carried
// over unchanged, so its numbers and shape want checking against the client.
func (war *Warrior) registerDualWieldSpecialization() {
	if war.Talents.DualWieldSpecialization == 0 {
		return
	}

	war.AddStaticMod(core.SpellModConfig{
		ProcMask:   core.ProcMaskMeleeOH,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: spellData.DualWieldSpecialization.Effect(shared.A_MOD_OFFHAND_DAMAGE_PCT, 0).FractionAt(war.Talents.DualWieldSpecialization),
	})
}

func (war *Warrior) registerImprovedExecute() {
	if war.Talents.ImprovedExecute == 0 {
		return
	}

	rageCostReduction := []int32{0, 2, 5}[war.Talents.ImprovedExecute]

	war.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskExecute,
		Kind:      core.SpellMod_PowerCost_Flat,
		IntValue:  -rageCostReduction,
	})
}

func (war *Warrior) registerEnrage() {
	if war.Talents.Enrage == 0 {
		return
	}

	war.EnrageAura = war.GetOrRegisterAura(core.Aura{
		Label:     "Enrage",
		ActionID:  core.ActionID{SpellID: 13048},
		Duration:  time.Second * 12,
		MaxStacks: 12,
	}).AttachSpellMod(core.SpellModConfig{
		School:     core.SpellSchoolPhysical,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 0.05 * float64(war.Talents.Enrage),
	}).AttachProcTrigger(core.ProcTrigger{
		Name:               "Enrage - Spend",
		TriggerImmediately: true,
		ProcMask:           core.ProcMaskMelee,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			war.EnrageAura.RemoveStack(sim)
		},
	})

	war.EnrageAura.NewExclusiveEffect("Enrage", true, core.ExclusiveEffect{Priority: 5 * float64(war.Talents.Enrage)})

	war.MakeProcTriggerAura(core.ProcTrigger{
		Name:     "Enrage - Trigger",
		ProcMask: core.ProcMaskMelee,
		Outcome:  core.OutcomeCrit | core.OutcomeSuppressedCrit,
		Callback: core.CallbackOnSpellHitTaken,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			war.EnrageAura.Activate(sim)
			war.EnrageAura.SetStacks(sim, 12)
		},
	})
}

func (war *Warrior) registerImprovedSlam() {
	if war.Talents.ImprovedSlam == 0 {
		return
	}

	war.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskSlam,
		Kind:      core.SpellMod_CastTime_Flat,
		TimeValue: -time.Millisecond * time.Duration(500*war.Talents.ImprovedSlam),
	})
}

func (war *Warrior) registerSweepingStrikes() {
	if !war.Talents.SweepingStrikes {
		return
	}

	actionID := core.ActionID{SpellID: 12723}

	var copyDamage float64
	hitSpell := war.RegisterSpell(core.SpellConfig{
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

	war.SweepingStrikesNormalizedAttack = war.RegisterSpell(core.SpellConfig{
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

	war.SweepingStrikesAura = war.MakeProcTriggerAura(core.ProcTrigger{
		Name:               "Sweeping Strikes",
		ActionID:           actionID,
		MetricsActionID:    actionID,
		Duration:           time.Second * 10,
		Callback:           core.CallbackOnSpellHitDealt,
		ProcMask:           core.ProcMaskMelee,
		Outcome:            core.OutcomeLanded,
		TriggerImmediately: true,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if war.Env.ActiveTargetCount() < 2 || war.SweepingStrikesAura.GetStacks() == 0 || result.PostOutcomeDamage <= 0 || !spell.ProcMask.Matches(core.ProcMaskMelee) {
				return
			}

			if spell.Matches(SpellMaskSweepingStrikesHit | SpellMaskSweepingStrikesNormalizedHit | SpellMaskThunderClap | SpellMaskWhirlwind | SpellMaskWhirlwindOh) {
				return
			}

			nextTarget := war.Env.NextActiveTargetUnit(result.Target)
			if spell.Matches(SpellMaskExecute) && sim.IsExecutePhase20() {
				war.SweepingStrikesNormalizedAttack.Cast(sim, nextTarget)
			} else {
				copyDamage = result.Damage / result.ArmorAndResistanceMultiplier
				hitSpell.Cast(sim, nextTarget)
			}

			war.SweepingStrikesAura.RemoveStack(sim)
		},
	})
	war.SweepingStrikesAura.MaxStacks = 10

	ssCD := war.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: SpellMaskSweepingStrikes,
		SpellSchool:    core.SpellSchoolPhysical,

		RageCost: core.RageCostOptions{
			Cost: 30,
		},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    war.NewTimer(),
				Duration: time.Second * 30,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return war.StanceMatches(BattleStance|BerserkerStance) || sim.ActiveTargetCount() > 1
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			spell.RelatedSelfBuff.Activate(sim)
			war.SweepingStrikesAura.SetStacks(sim, 10)
		},

		RelatedSelfBuff: war.SweepingStrikesAura,
	})

	war.AddMajorCooldown(core.MajorCooldown{
		Spell: ssCD,
		Type:  core.CooldownTypeDPS,
	})
}

func (war *Warrior) registerImprovedBerserkerRage() {
	if war.Talents.ImprovedBerserkerRage == 0 {
		return
	}

	core.MakePermanent(war.RegisterAura(core.Aura{
		Label:    "Improved Berserker Rage",
		ActionID: core.ActionID{SpellID: 20500}.WithTag(war.Talents.ImprovedBerserkerRage),
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			war.BerserkerRageRageGain += 5 * float64(war.Talents.ImprovedBerserkerRage)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			war.BerserkerRageRageGain -= 5 * float64(war.Talents.ImprovedBerserkerRage)
		},
	}))
}

func (war *Warrior) registerFlurry() {
	if war.Talents.Flurry == 0 {
		return
	}

	flurryAura := war.RegisterAura(core.Aura{
		Label:     "Flurry",
		ActionID:  core.ActionID{SpellID: 12970},
		Duration:  15 * time.Second,
		MaxStacks: 3,
	}).AttachMultiplyMeleeSpeed(1 + 0.05*float64(war.Talents.Flurry))

	war.MakeProcTriggerAura(core.ProcTrigger{
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
func (war *Warrior) registerPrecision() {
	if war.Talents.Precision == 0 {
		return
	}

	// The spell-hit effect carries the same ladder; only melee hit is taken here.
	war.AddStat(stats.PhysicalHitPercent, spellData.Precision.Effect(shared.A_MOD_HIT_CHANCE, 0).ValueAt(war.Talents.Precision))
}

func (war *Warrior) registerBloodthirst() {
	if !war.Talents.Bloodthirst {
		return
	}

	actionID := core.ActionID{SpellID: 30335}

	war.RegisterSpell(core.SpellConfig{
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
				Timer:    war.NewTimer(),
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
func (war *Warrior) registerIronWill() {
	if war.Talents.IronWill == 0 {
		return
	}
}

// registerPiercingHowl implements Piercing Howl, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (war *Warrior) registerPiercingHowl() {
	if !war.Talents.PiercingHowl {
		return
	}
}

// registerBloodCraze implements Blood Craze, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (war *Warrior) registerBloodCraze() {
	if war.Talents.BloodCraze == 0 {
		return
	}
}

// registerBoundlessRage implements Boundless Rage, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (war *Warrior) registerBoundlessRage() {
	if war.Talents.BoundlessRage == 0 {
		return
	}
}

// registerRagingBlows implements Raging Blows, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (war *Warrior) registerRagingBlows() {
	if !war.Talents.RagingBlows {
		return
	}
}
