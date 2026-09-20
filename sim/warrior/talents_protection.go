package warrior

import (
	"time"

	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

func (warrior *Warrior) registerProtectionTalents() {
	// Tier 1
	warrior.registerShieldSpecialization()
	warrior.registerAnticipation()

	// Tier 2
	// Improved Bloodrage: bloodrage.go
	warrior.registerToughness()
	warrior.registerImprovedThunderClap()

	// Tier 3
	warrior.registerLastStand()
	warrior.registerMasterOfDefense()
	warrior.registerImprovedRevenge()
	warrior.registerDefiance()

	// Tier 4
	warrior.registerImprovedSunderArmor()
	warrior.registerImprovedDisarm()
	warrior.registerVanguard()

	// Tier 5
	warrior.registerImprovedShieldWall()
	warrior.registerConcussionBlow()
	warrior.registerImprovedShieldBash()
	warrior.registerBastion()

	// Tier 6
	warrior.registerFocusedRage()

	// Tier 7
	warrior.registerShieldSlam()
}

func (warrior *Warrior) registerDefiance() {
	if warrior.Talents.Defiance == 0 {
		return
	}

	// Spell 12792 carries one effect, the threat modifier applied below (A_MOD_THREAT, +5/10/15%).
	// TODO: the client also requires a shield equipped, which needs a condition in stances.go.
	warrior.OnSpellRegistered(func(spell *core.Spell) {
		if !spell.Matches(SpellMaskDefensiveStance) {
			return
		}
		spell.RelatedSelfBuff.
			AttachMultiplicativePseudoStatBuff(&warrior.PseudoStats.ThreatMultiplier, spellData.Defiance.Effect(shared.A_MOD_THREAT, 127).MultiplierAt(warrior.Talents.Defiance))
	})
}

func (warrior *Warrior) registerAnticipation() {
	if warrior.Talents.Anticipation == 0 {
		return
	}

	warrior.AddStat(stats.DefenseRating, spellData.Anticipation.ValueAt(warrior.Talents.Anticipation)*core.DefenseRatingPerDefenseLevel)
}

func (warrior *Warrior) registerShieldSpecialization() {
	if warrior.Talents.ShieldSpecialization == 0 {
		return
	}

	warrior.AddStat(stats.BlockPercent, spellData.ShieldSpecialization.Effect(shared.A_MOD_BLOCK_PERCENT, 0).FractionAt(warrior.Talents.ShieldSpecialization))

	rageMetrics := warrior.NewRageMetrics(core.ActionID{SpellID: 1310318})

	warrior.MakeProcTriggerAura(core.ProcTrigger{
		Name: "Shield Specialization",
		// Effect 0 is the block bonus; the tooltip states the chance as $m2%, so effect 1's
		// ladder is the chance and the 100 in the proc chance column is noise.
		ProcChance:         spellData.ShieldSpecialization.EffectAt(1).FractionAt(warrior.Talents.ShieldSpecialization),
		TriggerImmediately: true,
		Outcome:            core.OutcomeBlock,
		Callback:           core.CallbackOnSpellHitTaken,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// TODO: Manual review needed -- spell 1310318 generates 5 Rage on a block ($m1/10).
			warrior.AddRage(sim, 5, rageMetrics)
		},
	})
}

func (warrior *Warrior) registerToughness() {
	if warrior.Talents.Toughness == 0 {
		return
	}

	// The client states the ladder twice, once on base armor and once on bonus armor; the sim's
	// single Armor stat takes one multiplier.
	warrior.MultiplyStat(stats.Armor, spellData.Toughness.Effect(shared.A_MOD_BASE_RESISTANCE_PCT, 1).MultiplierAt(warrior.Talents.Toughness))
}

func (warrior *Warrior) registerLastStand() {
	if !warrior.Talents.LastStand {
		return
	}

	actionID := core.ActionID{SpellID: 12975}
	healthMetrics := warrior.NewHealthMetrics(actionID)

	var bonusHealth float64
	aura := warrior.RegisterAura(core.Aura{
		Label:    "Last Stand",
		ActionID: actionID,
		// TODO: Manual review needed -- spell 12976 lasts 20 seconds.
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			// TODO: Manual review needed -- spell 12976 grants 30% of maximum health.
			bonusHealth = warrior.MaxHealth() * 0.3
			warrior.AddStatsDynamic(sim, stats.Stats{stats.Health: bonusHealth})
			warrior.GainHealth(sim, bonusHealth, healthMetrics)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.AddStatsDynamic(sim, stats.Stats{stats.Health: -bonusHealth})
		},
	})

	spell := warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: SpellMaskLastStand,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer: warrior.NewTimer(),
				// TODO: Manual review needed -- spell 12975 has a 3 minute cooldown.
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			aura.Activate(sim)
		},

		RelatedSelfBuff: aura,
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeSurvival,
		BuffAura: &core.StatBuffAura{
			Aura:            aura,
			BuffedStatTypes: []stats.Stat{stats.Health},
		},
	})
}

func (warrior *Warrior) registerImprovedSunderArmor() {
	if warrior.Talents.ImprovedSunderArmor == 0 {
		return
	}

	warrior.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskSunderArmor,
		Kind:      core.SpellMod_PowerCost_Flat,
		// The client states the rage reduction on a 0-1000 bar, so the ladder is tenths.
		IntValue: int32(spellData.ImprovedSunderArmor.ValueAt(warrior.Talents.ImprovedSunderArmor) / 10),
	})
}

func (warrior *Warrior) registerImprovedShieldWall() {
	if warrior.Talents.ImprovedShieldWall == 0 {
		return
	}

	warrior.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskShieldWall,
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: time.Duration(spellData.ImprovedShieldWall.ValueAt(warrior.Talents.ImprovedShieldWall)) * time.Millisecond,
	})
}

func (warrior *Warrior) registerConcussionBlow() {
	if !warrior.Talents.ConcussionBlow {
		return
	}

	warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 12809},
		ClassSpellMask: SpellMaskConcussionBlow,
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		MaxRange:       core.MaxMeleeRange,

		RageCost: core.RageCostOptions{
			// TODO: Manual review needed -- spell 12809 costs 10 Rage.
			Cost:   10,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer: warrior.NewTimer(),
				// TODO: Manual review needed -- spell 12809 has a 45 second cooldown.
				Duration: time.Second * 45,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMeleeSpecialHit)

			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
}

// TODO: Manual review needed -- spell 23922 states only "a very high amount of threat", so the
// flat threat is hand-supplied.
var shieldSlamRank = shared.WithSpellDataFlatThreat(spellData.ShieldSlam, 305).HighestRank()

func (warrior *Warrior) registerShieldSlam() {
	if !warrior.Talents.ShieldSlam {
		return
	}

	warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: shieldSlamRank.SpellID},
		ClassSpellMask: SpellMaskShieldSlam,
		SpellSchool:    shieldSlamRank.SpellSchool,
		DefenseType:    shieldSlamRank.DefenseType,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		MaxRange:       core.MaxMeleeRange,

		RageCost: core.RageCostOptions{
			Cost:   shieldSlamRank.Cost,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: shieldSlamRank.GCD,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: shieldSlamRank.Cooldown,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.PseudoStats.CanBlock
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		FlatThreatBonus:  shieldSlamRank.FlatThreatBonus,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := shieldSlamRank.Direct.Damage(sim) + warrior.BlockDamageReduction()
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
}

func (warrior *Warrior) registerFocusedRage() {
	if warrior.Talents.FocusedRage == 0 {
		return
	}

	warrior.AddStaticMod(core.SpellModConfig{
		ClassMask: WarriorSpellsAll ^ (SpellMaskDeathWish | SpellMaskBattleShout),
		Kind:      core.SpellMod_PowerCost_Flat,
		// The client states the rage reduction on a 0-1000 bar, so the ladder is tenths.
		IntValue: int32(spellData.FocusedRage.ValueAt(warrior.Talents.FocusedRage) / 10),
	})
}

func (warrior *Warrior) registerMasterOfDefense() {
	if warrior.Talents.MasterOfDefense == 0 {
		return
	}

	rageMetrics := warrior.NewRageMetrics(core.ActionID{SpellID: 23602})

	warrior.MakeProcTriggerAura(core.ProcTrigger{
		Name: "Master of Defense",
		// The tooltip states the chance as $m1%, so the talent's ladder is the chance and the 100
		// in the proc chance column is noise.
		ProcChance:         spellData.MasterOfDefense.FractionAt(warrior.Talents.MasterOfDefense),
		TriggerImmediately: true,
		Outcome:            core.OutcomeDodge | core.OutcomeParry,
		Callback:           core.CallbackOnSpellHitTaken,
		ExtraCondition: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) bool {
			return warrior.PseudoStats.CanBlock
		},
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// TODO: Manual review needed -- spell 23602 generates 5 Rage on a dodge or parry ($m1/10).
			warrior.AddRage(sim, 5, rageMetrics)
		},
	})
}

func (warrior *Warrior) registerImprovedRevenge() {
	if warrior.Talents.ImprovedRevenge == 0 {
		return
	}

	warrior.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskRevenge,
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: spellData.ImprovedRevenge.FractionAt(warrior.Talents.ImprovedRevenge),
	})
}

func (warrior *Warrior) registerImprovedDisarm() {
	if warrior.Talents.ImprovedDisarm == 0 {
		return
	}

	warrior.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskDisarm,
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: time.Duration(spellData.ImprovedDisarm.ValueAt(warrior.Talents.ImprovedDisarm)) * time.Millisecond,
	})
}

// TODO: Vanguard makes Charge usable in Defensive Stance -- spell 1310317 overrides each Charge
// rank with 1240287/1240288/1240289, which carry the same rage and cooldown. That needs the stance
// condition in charge.go.
func (warrior *Warrior) registerVanguard() {
	if !warrior.Talents.Vanguard {
		return
	}
}

func (warrior *Warrior) registerImprovedShieldBash() {
	if warrior.Talents.ImprovedShieldBash == 0 {
		return
	}

	// TODO: nothing in the sim reads a silence on an enemy, so the aura only shows up in metrics.
	silenceAuras := warrior.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    "Shield Bash - Silence",
			ActionID: core.ActionID{SpellID: 18498},
			// TODO: Manual review needed -- spell 18498 silences for 3 seconds.
			Duration: time.Second * 3,
		})
	})

	warrior.MakeProcTriggerAura(core.ProcTrigger{
		Name: "Improved Shield Bash",
		// The tooltip states the chance as $m1%, so the talent's ladder is the chance and the 100
		// in the proc chance column is noise.
		ProcChance:         spellData.ImprovedShieldBash.FractionAt(warrior.Talents.ImprovedShieldBash),
		TriggerImmediately: true,
		ClassSpellMask:     SpellMaskShieldBash,
		Outcome:            core.OutcomeLanded,
		Callback:           core.CallbackOnSpellHitDealt,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			silenceAuras.Get(result.Target).Activate(sim)
		},
	})
}

func (warrior *Warrior) registerBastion() {
	if warrior.Talents.Bastion == 0 {
		return
	}

	// The client applies the bonus to the physical school (A_MOD_DAMAGE_PERCENT_DONE, mask 1).
	damageMod := warrior.AddDynamicMod(core.SpellModConfig{
		School:     core.SpellSchoolPhysical,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: spellData.Bastion.FractionAt(warrior.Talents.Bastion),
	})

	if warrior.PseudoStats.CanBlock {
		damageMod.Activate()
	}

	warrior.RegisterItemSwapCallback([]proto.ItemSlot{proto.ItemSlot_ItemSlotOffHand}, func(sim *core.Simulation, slot proto.ItemSlot) {
		if warrior.PseudoStats.CanBlock {
			damageMod.Activate()
		} else {
			damageMod.Deactivate()
		}
	})
}

func (warrior *Warrior) registerImprovedThunderClap() {
	if warrior.Talents.ImprovedThunderClap == 0 {
		return
	}

	// Slowing effect implemented in core/debuffs.go

	warrior.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskThunderClap,
		Kind:      core.SpellMod_PowerCost_Flat,
		// The client states the rage reduction on a 0-1000 bar, so the ladder is tenths.
		IntValue: int32(spellData.ImprovedThunderClap.ValueAt(warrior.Talents.ImprovedThunderClap) / 10),
	})
}
