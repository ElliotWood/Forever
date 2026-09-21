package rogue

import (
	"time"

	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/stats"
)

// Hemorrhage has no rank subtext, so the generator gives it a single row.
var hemorrhageRank = spellData.Hemorrhage.BySpellID(16511)

func (rogue *Rogue) registerSubtletyTalents() {
	// Tier 1
	rogue.registerMasterOfDeception()
	rogue.registerOpportunity()

	// Tier 2
	rogue.registerSetup()
	rogue.registerCamouflage()

	// Tier 3
	rogue.registerInitiative()
	rogue.registerGhostlyStrike()
	rogue.registerImprovedAmbush()
	rogue.registerImprovedDistract()

	// Tier 4
	rogue.registerElusiveness()
	rogue.registerSerratedBlades()
	rogue.registerDirtyTricks()

	// Tier 5
	rogue.registerHeightenedSenses()
	rogue.registerPreparation()
	rogue.registerDirtyDeeds()
	rogue.registerHemorrhage()

	// Tier 6
	rogue.registerQuietus()
	rogue.registerCutthroat()

	// Tier 7
	rogue.registerPremeditation()
	rogue.registerThousandCuts()
}

func (rogue *Rogue) registerOpportunity() {
	if rogue.Talents.Opportunity == 0 {
		return
	}

	rogue.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		ClassMask:  RogueSpellBackstab | RogueSpellMutilate | RogueSpellMutilateHit | RogueSpellAmbush | RogueSpellGarrote,
		FloatValue: spellData.Opportunity.Effect(shared.A_ADD_PCT_MODIFIER, shared.SPELLMOD_DAMAGE).FractionAt(rogue.Talents.Opportunity),
	})
}

func (rogue *Rogue) registerInitiative() {
	if rogue.Talents.Initiative == 0 {
		return
	}

	initMetrics := rogue.NewComboPointMetrics(core.ActionID{SpellID: spellData.InitiativeTriggered.HighestRank().SpellID})

	rogue.MakeProcTriggerAura(core.ProcTrigger{
		Name:     "Initiative Trigger",
		ActionID: core.ActionID{SpellID: spellData.Initiative.HighestRank().SpellID},
		// The beta rounds rank 2 up to 67% rather than doubling rank 1's 33%; ProcChanceAt would
		// read the flat 100 the talent spell carries.
		ProcChance:     spellData.Initiative.FractionAt(rogue.Talents.Initiative),
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		ClassSpellMask: RogueSpellGarrote | RogueSpellAmbush,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			rogue.AddComboPoints(sim, 1, initMetrics)
		},
	})
}

func (rogue *Rogue) registerGhostlyStrike() {
	if !rogue.Talents.GhostlyStrike {
		return
	}

	ghostlyStrikeRank := spellData.GhostlyStrike.HighestRank()
	actionID := core.ActionID{SpellID: ghostlyStrikeRank.SpellID}

	// Effect 0 is the plain weapon share, effect 3 the larger one a dagger gets.
	weaponDamage := spellData.GhostlyStrike.EffectAt(0).ValueAt(1) / 100
	if rogue.HasDagger(core.MainHand) {
		weaponDamage = spellData.GhostlyStrike.EffectAt(3).ValueAt(1) / 100
	}

	dodgeAura := rogue.RegisterAura(core.Aura{
		Label:    "Ghostly Strike Buff",
		ActionID: actionID,
		Duration: ghostlyStrikeRank.Duration,
	}).AttachStatBuff(stats.DodgeRating, ghostlyStrikeRank.Effect(shared.A_MOD_DODGE_PERCENT, 0).Value*core.DodgeRatingPerDodgePercent)

	rogue.GhostlyStrike = rogue.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: RogueSpellGhostlyStrike,
		SpellSchool:    ghostlyStrikeRank.SpellSchool,
		DefenseType:    ghostlyStrikeRank.DefenseType,
		Flags:          core.SpellFlagAPL | core.SpellFlagMeleeMetrics | SpellFlagBuilder,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		MaxRange:       core.MaxMeleeRange,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: ghostlyStrikeRank.GCD,
			},
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: ghostlyStrikeRank.Cooldown,
			},
			IgnoreHaste: true,
		},
		// The table reads zero for the cost, a generator gap; the client charges 40.
		EnergyCost: core.EnergyCostOptions{
			Cost:   40,
			Refund: ghostlyStrikeRank.MissRefund(),
		},

		DamageMultiplier:         weaponDamage,
		DamageMultiplierAdditive: 1,
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)

			damage := rogue.MHWeaponDamage(sim, spell.MeleeAttackPower(target))
			result := spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			dodgeAura.Activate(sim)

			if result.Landed() {
				rogue.AddComboPoints(sim, 1, spell.ComboPointMetrics())
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}

func (rogue *Rogue) registerImprovedAmbush() {
	if rogue.Talents.ImprovedAmbush == 0 {
		return
	}

	rogue.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_BonusCrit_Percent,
		ClassMask:  RogueSpellAmbush,
		FloatValue: spellData.ImprovedAmbush.ValueAt(rogue.Talents.ImprovedAmbush),
	})
}

func (rogue *Rogue) registerElusiveness() {
	if rogue.Talents.Elusiveness == 0 {
		return
	}

	rogue.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Flat,
		ClassMask: RogueSpellVanish,
		TimeValue: time.Duration(spellData.Elusiveness.Effect(shared.A_ADD_FLAT_MODIFIER, shared.SPELLMOD_COOLDOWN).ValueAt(rogue.Talents.Elusiveness)) * time.Millisecond,
	})
}

// Serrated Blades ignores a share of the target's Armor rather than a flat amount, and raises
// the rogue's own Rupture.
func (rogue *Rogue) registerSerratedBlades() {
	if rogue.Talents.SerratedBlades == 0 {
		return
	}

	rogue.addArmorIgnore(spellData.SerratedBlades.EffectAt(0).ValueAt(rogue.Talents.SerratedBlades) / 100)
	rogue.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		ClassMask:  RogueSpellRupture,
		FloatValue: spellData.SerratedBlades.Effect(shared.A_ADD_PCT_MODIFIER, shared.SPELLMOD_DOT).FractionAt(rogue.Talents.SerratedBlades),
	})
}

func (rogue *Rogue) registerPreparation() {
	if !rogue.Talents.Preparation {
		return
	}

	preparationRank := spellData.Preparation.HighestRank()

	rogue.Preparation = rogue.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: preparationRank.SpellID},
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: RogueSpellPreparation,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: preparationRank.GCD,
			},
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: preparationRank.Cooldown,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, affected := range []*core.Spell{rogue.ColdBlood, rogue.Shadowstep, rogue.Premeditation, rogue.Vanish} {
				if affected != nil {
					affected.CD.Reset()
				}
			}
		},
	})

	rogue.AddMajorCooldown(core.MajorCooldown{
		Spell: rogue.Preparation,
		Type:  core.CooldownTypeDPS,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			return rogue.Vanish != nil && !rogue.Vanish.CD.IsReady(sim)
		},
	})
}

// Forever states only an energy discount on Dirty Deeds; the Classic damage bonus below 35%
// health is gone. It also drops Garrote's positional requirement, handled in garrote.go.
func (rogue *Rogue) registerDirtyDeeds() {
	if rogue.Talents.DirtyDeeds == 0 {
		return
	}

	rogue.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_PowerCost_Flat,
		ClassMask: RogueSpellGarrote,
		IntValue:  int32(spellData.DirtyDeeds.Effect(shared.A_ADD_FLAT_MODIFIER, shared.SPELLMOD_COST).ValueAt(rogue.Talents.DirtyDeeds)),
	})
}

// Hemorrhage no longer weakens the target for the whole raid, it makes the rogue's own Rupture
// hit harder.
const HemorrhageRuptureMultiplier = 1.15

func (rogue *Rogue) registerHemorrhage() {
	if !rogue.Talents.Hemorrhage {
		return
	}

	actionID := core.ActionID{SpellID: hemorrhageRank.SpellID}

	rogue.HemorrhageAuras = rogue.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    "Hemorrhage-" + rogue.Label,
			ActionID: actionID,
			Duration: hemorrhageRank.Duration,
		})
	})

	// The plain weapon share is effect 3 (100%); effect 4 is the larger one a dagger gets.
	weaponDamage := spellData.Hemorrhage.EffectAt(3).ValueAt(1) / 100
	if rogue.HasDagger(core.MainHand) {
		weaponDamage = spellData.Hemorrhage.EffectAt(4).ValueAt(1) / 100
	}

	rogue.Hemorrhage = rogue.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: RogueSpellHemorrhage,
		SpellSchool:    hemorrhageRank.SpellSchool,
		DefenseType:    hemorrhageRank.DefenseType,
		Flags:          core.SpellFlagAPL | core.SpellFlagMeleeMetrics | SpellFlagBuilder,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		MaxRange:       core.MaxMeleeRange,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: hemorrhageRank.GCD,
			},
			IgnoreHaste: true,
		},
		EnergyCost: core.EnergyCostOptions{
			Cost:   hemorrhageRank.Cost,
			Refund: hemorrhageRank.MissRefund(),
		},

		DamageMultiplier:         weaponDamage,
		DamageMultiplierAdditive: 1,
		ThreatMultiplier:         1,

		BonusCoefficient: hemorrhageRank.Direct.BonusCoefficient(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)

			// The beta client moved Hemorrhage from plain to normalized weapon damage.
			damage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower(target))
			result := spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if result.Landed() {
				rogue.AddComboPoints(sim, 1, spell.ComboPointMetrics())
				rogue.HemorrhageAuras.Get(target).Activate(sim)
			} else {
				spell.IssueRefund(sim)
			}
		},

		RelatedAuraArrays: rogue.HemorrhageAuras.ToMap(),
	})
}

// Hemorrhage is optional, so the debuff has to be looked up defensively.
func (rogue *Rogue) isHemorrhaging(target *core.Unit) bool {
	return rogue.HemorrhageAuras != nil && rogue.HemorrhageAuras.Get(target).IsActive()
}

func (rogue *Rogue) registerPremeditation() {
	if !rogue.Talents.Premeditation {
		return
	}

	premeditationRank := spellData.Premeditation.HighestRank()
	actionID := core.ActionID{SpellID: premeditationRank.SpellID}
	comboMetrics := rogue.NewComboPointMetrics(actionID)
	points, _ := premeditationRank.Energize.Range()

	rogue.Premeditation = rogue.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagAPL | core.SpellFlagNoOnCastComplete,
		ClassSpellMask: RogueSpellPremeditation,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: 0,
				GCD:  0,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: premeditationRank.Cooldown,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return rogue.IsStealthed()
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			rogue.AddComboPoints(sim, int32(points), comboMetrics)
		},
	})

	rogue.AddMajorCooldown(core.MajorCooldown{
		Spell:              rogue.Premeditation,
		Type:               core.CooldownTypeDPS,
		Priority:           core.CooldownPriorityLow,
		AllowSpellQueueing: true,
	})
}

// Quietus, new in Forever: the rogue's strikes hit harder once the target is in execute range.
func (rogue *Rogue) registerQuietus() {
	if rogue.Talents.Quietus == 0 {
		return
	}

	quietusAura := rogue.GetOrRegisterAura(core.Aura{
		Label:    "Quietus",
		ActionID: core.ActionID{SpellID: spellData.Quietus.HighestRank().SpellID},
		Duration: core.NeverExpires,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		ClassMask:  RogueSpellStrikes,
		FloatValue: spellData.Quietus.FractionAt(rogue.Talents.Quietus),
	})

	rogue.RegisterResetEffect(func(sim *core.Simulation) {
		quietusAura.Deactivate(sim)
		sim.RegisterExecutePhaseCallback(func(sim *core.Simulation, isExecute int32) {
			if isExecute == 35 {
				quietusAura.Activate(sim)
			}
		})
	})
}

// Cutthroat, new in Forever: a Backstab can let the next Ambush be used outside of Stealth.
func (rogue *Rogue) registerCutthroat() {
	if rogue.Talents.Cutthroat == 0 {
		return
	}

	triggered := spellData.CutthroatTriggered.HighestRank()

	rogue.CutthroatAura = rogue.RegisterAura(core.Aura{
		Label:    "Cutthroat",
		ActionID: core.ActionID{SpellID: triggered.SpellID},
		Duration: triggered.Duration,
	})

	rogue.MakeProcTriggerAura(core.ProcTrigger{
		Name:           "Cutthroat Trigger",
		ActionID:       core.ActionID{SpellID: spellData.Cutthroat.HighestRank().SpellID},
		ProcChance:     spellData.Cutthroat.FractionAt(rogue.Talents.Cutthroat),
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		ClassSpellMask: RogueSpellBackstab,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			rogue.CutthroatAura.Activate(sim)
		},
	})
}

// Thousand Cuts, new in Forever: Rupture's ticks discount the next Hemorrhage or Backstab.
// The client ships no ranked spell for it, so the numbers are our Forever sim's.
func (rogue *Rogue) registerThousandCuts() {
	if !rogue.Talents.ThousandCuts {
		return
	}

	rogue.ThousandCutsAura = rogue.RegisterAura(core.Aura{
		Label:     "Thousand Cuts",
		ActionID:  core.ActionID{SpellID: 1310714},
		Duration:  time.Second * 10,
		MaxStacks: 5,
		OnApplyEffects: func(aura *core.Aura, sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if spell.Matches(RogueSpellBackstab | RogueSpellHemorrhage) {
				aura.Deactivate(sim)
			}
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_PowerCost_Flat,
		ClassMask: RogueSpellBackstab | RogueSpellHemorrhage,
		IntValue:  -3,
	})

	rogue.MakeProcTriggerAura(core.ProcTrigger{
		Name:           "Thousand Cuts Trigger",
		Callback:       core.CallbackOnPeriodicDamageDealt,
		ClassSpellMask: RogueSpellRupture,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			rogue.ThousandCutsAura.Activate(sim)
			rogue.ThousandCutsAura.AddStack(sim)
		},
	})
}

// registerCamouflage implements Camouflage.
//
// TODO: Not modelled. Stealth movement speed and a shorter Vanish cooldown.
func (rogue *Rogue) registerCamouflage() {}

// registerDirtyTricks implements Dirty Tricks.
//
// TODO: Not modelled. Discounts Sap and Blind, neither of which our sim casts.
func (rogue *Rogue) registerDirtyTricks() {}

// registerHeightenedSenses implements Heightened Senses.
//
// TODO: Not modelled. Stealth detection and a lower chance to be hit by spells.
func (rogue *Rogue) registerHeightenedSenses() {}

// registerImprovedDistract implements Improved Distract.
//
// TODO: Not modelled. Distract radius and cost.
func (rogue *Rogue) registerImprovedDistract() {}

// registerMasterOfDeception implements Master of Deception.
//
// TODO: Not modelled. Stealth detection only.
func (rogue *Rogue) registerMasterOfDeception() {}

// registerSetup implements Setup.
//
// TODO: Not modelled. Gives a combo point when the rogue dodges, and the rogue is not the one
// being attacked in a DPS sim.
func (rogue *Rogue) registerSetup() {}
