package warrior

import (
	"time"

	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

func (warrior *Warrior) registerFuryTalents() {
	// Tier 1
	// Booming Voice: shouts.go
	warrior.registerCruelty()

	// Tier 2
	// TODO: Iron Will (12962) shortens stuns and fears by 3% per rank; core registers them at a
	// fixed duration with no per-unit modifier to hang the talent on.
	warrior.registerUnbridledWrath()

	// Tier 3
	warrior.registerImprovedCleave()
	warrior.registerPiercingHowl()
	warrior.registerBloodCraze()
	// Boundless Rage: warrior.go, when it enables the rage bar

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

	rageMetrics := warrior.NewRageMetrics(core.ActionID{SpellID: 12964})

	// TODO: Manual review needed -- 12964 restores 10 rage tenths, which 12322 doubles for a
	// two-handed weapon.
	rageGain := func() float64 {
		if warrior.GetMainHandType() == proto.HandType_HandTypeTwoHand {
			return 2
		}
		return 1
	}

	warrior.MakeProcTriggerAura(core.ProcTrigger{
		Name:               "Unbridled Wrath",
		ProcMask:           core.ProcMaskMeleeWhiteHit,
		ProcChance:         spellData.UnbridledWrath.FractionAt(warrior.Talents.UnbridledWrath),
		RequireDamageDealt: true,
		Outcome:            core.OutcomeLanded,
		Callback:           core.CallbackOnSpellHitDealt,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			warrior.AddRage(sim, rageGain(), rageMetrics)
		},
	})
}

func (warrior *Warrior) registerDualWieldSpecialization() {
	if warrior.Talents.DualWieldSpecialization == 0 {
		return
	}

	warrior.AddStaticMod(core.SpellModConfig{
		ProcMask:   core.ProcMaskMeleeOH,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: spellData.DualWieldSpecialization.Effect(shared.A_MOD_OFFHAND_DAMAGE_PCT, 0).FractionAt(warrior.Talents.DualWieldSpecialization),
	})

	warrior.AddStaticMod(core.SpellModConfig{
		ProcMask:   core.ProcMaskMeleeOH,
		Kind:       core.SpellMod_BonusHit_Percent,
		FloatValue: spellData.DualWieldSpecialization.Effect(shared.A_MOD_HIT_CHANCE, 0).ValueAt(warrior.Talents.DualWieldSpecialization),
	})

	// TODO: The off-hand Rage generation the talent's second effect states is not modelled; the
	// rage bar halves the hit factor for off-hand swings inside sim/core/rage.go and takes no
	// per-hand multiplier.
}

func (warrior *Warrior) registerImprovedExecute() {
	if warrior.Talents.ImprovedExecute == 0 {
		return
	}

	warrior.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskExecute,
		Kind:      core.SpellMod_PowerCost_Flat,
		IntValue:  int32(spellData.ImprovedExecute.ValueAt(warrior.Talents.ImprovedExecute) / 10),
	})
}

func (warrior *Warrior) registerEnrage() {
	if warrior.Talents.Enrage == 0 {
		return
	}

	warrior.EnrageAura = warrior.GetOrRegisterAura(core.Aura{
		Label:    "Enrage",
		ActionID: core.ActionID{SpellID: 12880},
		// TODO: Manual review needed -- 12880 lasts 12 seconds and carries no charge count.
		Duration: time.Second * 12,
	}).AttachSpellMod(core.SpellModConfig{
		School:     core.SpellSchoolPhysical,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: spellData.Enrage.FractionAt(warrior.Talents.Enrage),
	})

	warrior.EnrageAura.NewExclusiveEffect("Enrage", true, core.ExclusiveEffect{Priority: spellData.Enrage.ValueAt(warrior.Talents.Enrage)})

	warrior.MakeProcTriggerAura(core.ProcTrigger{
		Name:               "Enrage - Trigger",
		Callback:           core.CallbackOnSpellHitTaken,
		Outcome:            core.OutcomeLanded,
		RequireDamageDealt: true,
		ProcChance:         spellData.Enrage.ProcChanceAt(warrior.Talents.Enrage),
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			warrior.EnrageAura.Activate(sim)
		},
	})
}

func (warrior *Warrior) registerImprovedBerserkerRage() {
	if warrior.Talents.ImprovedBerserkerRage == 0 {
		return
	}

	// Both of the talent's effects are dummies, so the rage one is named by its index.
	rageGain := spellData.ImprovedBerserkerRage.EffectAt(0).ValueAt(warrior.Talents.ImprovedBerserkerRage) / 10

	core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label:    "Improved Berserker Rage",
		ActionID: core.ActionID{SpellID: 20500}.WithTag(warrior.Talents.ImprovedBerserkerRage),
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.BerserkerRageRageGain += rageGain
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.BerserkerRageRageGain -= rageGain
		},
	}))

	// TODO: The chance to shed movement impairing effects the talent's second effect states is not
	// modelled; nothing snares the warrior in the sim.
}

func (warrior *Warrior) registerFlurry() {
	if warrior.Talents.Flurry == 0 {
		return
	}

	flurryAura := warrior.RegisterAura(core.Aura{
		Label:    "Flurry",
		ActionID: core.ActionID{SpellID: 12966},
		// TODO: Manual review needed -- 12966 lasts 15 seconds and carries 3 charges.
		Duration:  15 * time.Second,
		MaxStacks: 3,
	}).AttachMultiplyMeleeSpeed(spellData.Flurry.MultiplierAt(warrior.Talents.Flurry))

	warrior.MakeProcTriggerAura(core.ProcTrigger{
		Name:               "Flurry - Trigger",
		ActionID:           core.ActionID{SpellID: 12319},
		ProcMask:           core.ProcMaskMelee,
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

func (warrior *Warrior) registerPrecision() {
	if warrior.Talents.Precision == 0 {
		return
	}

	warrior.AddStat(stats.PhysicalHitPercent, spellData.Precision.Effect(shared.A_MOD_HIT_CHANCE, 0).ValueAt(warrior.Talents.Precision))
	warrior.AddStat(stats.SpellHitPercent, spellData.Precision.Effect(shared.A_MOD_SPELL_HIT_CHANCE, 0).ValueAt(warrior.Talents.Precision))
}

var bloodthirstRank = spellData.Bloodthirst.BySpellID(23894)

func (warrior *Warrior) registerBloodthirst() {
	if !warrior.Talents.Bloodthirst {
		return
	}

	// The attack power share sits on the second effect; the first is the flat damage added to it.
	apShare := bloodthirstRank.Effects[1].Value / 100

	warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: bloodthirstRank.SpellID},
		Rank:           bloodthirstRank.Rank,
		SpellSchool:    bloodthirstRank.SpellSchool,
		DefenseType:    bloodthirstRank.DefenseType,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskBloodthirst,
		MaxRange:       bloodthirstRank.MaxRange,

		RageCost: core.RageCostOptions{
			Cost:   bloodthirstRank.Cost,
			Refund: 0.8,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: bloodthirstRank.GCD,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: bloodthirstRank.Cooldown,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.MeleeAttackPower(target)*apShare + bloodthirstRank.Direct.Damage(sim)
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
}

// TODO: The daze itself is not modelled; the encounter's targets do not move, so the -50% movement
// speed 12323 applies for 6 seconds within 10 yards has nothing to act on.
func (warrior *Warrior) registerPiercingHowl() {
	if !warrior.Talents.PiercingHowl {
		return
	}

	warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 12323},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskNone,

		RageCost: core.RageCostOptions{
			// TODO: Manual review needed -- 12323 costs 100 rage tenths.
			Cost: 10,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				// TODO: Manual review needed -- 12323 states a 1.5 second global cooldown and no
				// cooldown of its own.
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},
	})
}

// TODO: Manual review needed -- 16487's second effect sets the share of maximum health a single hit
// has to exceed at 20%; the generator leaves that effect out for having no rank curve.
const bloodCrazeHealthThreshold = 0.2

func (warrior *Warrior) registerBloodCraze() {
	if warrior.Talents.BloodCraze == 0 {
		return
	}

	healthFraction := spellData.BloodCraze.FractionAt(warrior.Talents.BloodCraze)

	bloodCraze := warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 16488},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskSpellHealing,
		Flags:       core.SpellFlagPassiveSpell | core.SpellFlagHelpful | core.SpellFlagNoOnCastComplete,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Hot: core.DotConfig{
			Aura:     core.Aura{Label: "Blood Craze"},
			SelfOnly: true,
			// TODO: Manual review needed -- 16488 ticks every 2 seconds for 6 seconds.
			NumberOfTicks: 3,
			TickLength:    time.Second * 2,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				healPerTick := warrior.MaxHealth() * healthFraction / float64(dot.ExpectedTickCount())
				dot.Spell.CalcAndDealPeriodicHealing(sim, target, healPerTick, dot.OutcomeTick)
			},
		},
	})

	warrior.MakeProcTriggerAura(core.ProcTrigger{
		Name:               "Blood Craze - Damage Taken",
		Callback:           core.CallbackOnSpellHitTaken,
		Outcome:            core.OutcomeLanded,
		RequireDamageDealt: true,
		ExtraCondition: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) bool {
			return result.Outcome.Matches(core.OutcomeCrit) || result.Damage > warrior.MaxHealth()*bloodCrazeHealthThreshold
		},
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			bloodCraze.SelfHot().Apply(sim)
		},
	})

	warrior.MakeProcTriggerAura(core.ProcTrigger{
		Name:               "Blood Craze - Bloodthirst",
		Callback:           core.CallbackOnSpellHitDealt,
		ClassSpellMask:     SpellMaskBloodthirst,
		Outcome:            core.OutcomeLanded,
		RequireDamageDealt: true,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			bloodCraze.SelfHot().Apply(sim)
		},
	})
}

// TODO: The other half of 1310315, Whirlwind striking with the off-hand as well, is not gated on the
// talent: whirlwind.go casts its off-hand hit off every Whirlwind the warrior lands with an off-hand
// weapon equipped.
func (warrior *Warrior) registerRagingBlows() {
	if !warrior.Talents.RagingBlows {
		return
	}

	warrior.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskCleave,
		Kind:      core.SpellMod_PowerCost_Flat,
		// TODO: Manual review needed -- 1310315's second effect states -20 rage tenths on Cleave.
		IntValue: -2,
	})
}

func (warrior *Warrior) registerDeathWish() {
	if !warrior.Talents.DeathWish {
		return
	}

	actionID := core.ActionID{SpellID: 12328}

	deathWishAura := warrior.RegisterAura(core.Aura{
		Label:    "Death Wish",
		ActionID: actionID,
		// TODO: Manual review needed -- 12328 lasts 30 seconds.
		Duration: time.Second * 30,
	}).
		// TODO: Manual review needed -- 12328 states +20% Physical damage done.
		AttachMultiplicativePseudoStatBuff(
			&warrior.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical], 1.2,
		).
		// TODO: Manual review needed -- 12328 states +5% damage taken.
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
			// TODO: Manual review needed -- 12328 costs 100 rage tenths.
			Cost: 10,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer: warrior.NewTimer(),
				// TODO: Manual review needed -- 12328 states a 3 minute cooldown.
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
		TimeValue: time.Duration(spellData.ImprovedIntercept.ValueAt(warrior.Talents.ImprovedIntercept)) * time.Millisecond,
	})
}

// Improved Cleave (12329) states only a rage discount on Cleave.
func (warrior *Warrior) registerImprovedCleave() {
	if warrior.Talents.ImprovedCleave == 0 {
		return
	}
	warrior.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskCleave,
		Kind:      core.SpellMod_PowerCost_Flat,
		IntValue:  int32(spellData.ImprovedCleave.ValueAt(warrior.Talents.ImprovedCleave) / 10),
	})
}
