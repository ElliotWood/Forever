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
	// Booming Voice: shouts.go
	warrior.registerCruelty()

	// Tier 2
	warrior.registerIronWill()
	warrior.registerUnbridledWrath()

	// Tier 3
	// Improved Cleave: heroic_strike_cleave.go
	warrior.registerPiercingHowl()
	warrior.registerBloodCraze()
	warrior.registerBoundlessRage()

	// Tier 4
	warrior.registerDualWieldSpecialization()
	warrior.registerRagingBlows()
	warrior.registerEnrage()
	warrior.registerImprovedExecute()

	// Tier 5
	warrior.registerPrecision()
	warrior.registerDeathWish()
	warrior.registerImprovedIntercept()

	// Tier 6
	warrior.registerImprovedBerserkerRage()
	warrior.registerFlurry()

	// Tier 7
	warrior.registerBloodthirst()
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

// TODO: registerIronWill models nothing yet; spellData.IronWill carries the ranks.
func (warrior *Warrior) registerIronWill() {
	if warrior.Talents.IronWill == 0 {
		return
	}
}

// TODO: registerPiercingHowl models nothing yet; the client has no ladder for it, so the tooltip is the source.
func (warrior *Warrior) registerPiercingHowl() {
	if !warrior.Talents.PiercingHowl {
		return
	}
}

// TODO: registerBloodCraze models nothing yet; spellData.BloodCraze carries the ranks.
func (warrior *Warrior) registerBloodCraze() {
	if warrior.Talents.BloodCraze == 0 {
		return
	}
}

// TODO: registerBoundlessRage models nothing yet; spellData.BoundlessRage carries the ranks.
func (warrior *Warrior) registerBoundlessRage() {
	if warrior.Talents.BoundlessRage == 0 {
		return
	}
}

// TODO: registerRagingBlows models nothing yet; the client has no ladder for it, so the tooltip is the source.
func (warrior *Warrior) registerRagingBlows() {
	if !warrior.Talents.RagingBlows {
		return
	}
}

func (warrior *Warrior) registerDeathWish() {
	if !warrior.Talents.DeathWish {
		return
	}

	actionID := core.ActionID{SpellID: 12292}

	deathWishAura := warrior.RegisterAura(core.Aura{
		Label:    "Death Wish",
		ActionID: actionID,
		Duration: time.Second * 30,
	}).
		AttachMultiplicativePseudoStatBuff(
			&warrior.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical], 1.2,
		).
		AttachMultiplicativePseudoStatBuff(
			&warrior.PseudoStats.DamageTakenMultiplier, 1.05,
		).
		// Grants immunity to Fear effects.
		AttachFearImmunity()

	deathWishSpell := warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: SpellMaskDeathWish,
		Flags:          core.SpellFlagCastWhileIncapacitated,

		RageCost: core.RageCostOptions{
			Cost: 10,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			deathWishAura.Activate(sim)
			warrior.WaitUntil(sim, sim.CurrentTime+core.GCDDefault)
		},

		RelatedSelfBuff: deathWishAura,
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: deathWishSpell,
		Type:  core.CooldownTypeDPS,
	})
}

func (warrior *Warrior) registerImprovedIntercept() {
	if warrior.Talents.ImprovedIntercept == 0 {
		return
	}

	warrior.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskIntercept,
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: -time.Second * time.Duration(5*warrior.Talents.ImprovedIntercept),
	})
}
