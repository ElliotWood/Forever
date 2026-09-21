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
	// Defiance: stances.go

	// Tier 4
	warrior.registerImprovedSunderArmor()
	warrior.registerImprovedDisarm()
	// Vanguard: charge.go

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

func (warrior *Warrior) registerAnticipation() {
	if warrior.Talents.Anticipation == 0 {
		return
	}

	warrior.AddStat(stats.DefenseRating, spellData.Anticipation.ValueAt(warrior.Talents.Anticipation)*core.DefenseRatingPerDefenseLevel)
}

var shieldSpecializationEnergize = spellData.ShieldSpecializationTriggered.HighestRank()

func (warrior *Warrior) registerShieldSpecialization() {
	if warrior.Talents.ShieldSpecialization == 0 {
		return
	}

	warrior.AddStat(stats.BlockPercent, spellData.ShieldSpecialization.Effect(shared.A_MOD_BLOCK_PERCENT, 0).FractionAt(warrior.Talents.ShieldSpecialization))

	// Effect 0 is the block bonus; the tooltip states the chance as $m2%, so effect 1's ladder is
	// the chance and the 100 in the proc chance column is noise.
	warrior.registerRageOnAvoid("Shield Specialization", shieldSpecializationEnergize.SpellID,
		shared.SpellDataMin(shieldSpecializationEnergize.Energize)/10,
		spellData.ShieldSpecialization.EffectAt(1).FractionAt(warrior.Talents.ShieldSpecialization), core.OutcomeBlock, nil)
}

// A chance to gain rage when an incoming attack is blocked, dodged or parried. The energize the
// triggered spell states is on the client's 0-1000 rage bar, so the call site divides it by ten.
func (warrior *Warrior) registerRageOnAvoid(name string, spellID int32, rage float64, chance float64, outcome core.HitOutcome, extra func() bool) {
	rageMetrics := warrior.NewRageMetrics(core.ActionID{SpellID: spellID})
	trigger := core.ProcTrigger{
		Name:               name,
		ProcChance:         chance,
		TriggerImmediately: true,
		Outcome:            outcome,
		Callback:           core.CallbackOnSpellHitTaken,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			warrior.AddRage(sim, rage, rageMetrics)
		},
	}
	if extra != nil {
		trigger.ExtraCondition = func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) bool { return extra() }
	}
	warrior.MakeProcTriggerAura(trigger)
}

func (warrior *Warrior) registerToughness() {
	if warrior.Talents.Toughness == 0 {
		return
	}

	// The client states the ladder twice, once on base armor and once on bonus armor; the sim's
	// single Armor stat takes one multiplier.
	// The tooltip states armor from items, which is the equipment share of the stat.
	warrior.ApplyEquipScaling(stats.Armor, spellData.Toughness.Effect(shared.A_MOD_BASE_RESISTANCE_PCT, 1).MultiplierAt(warrior.Talents.Toughness))
}

var lastStandRank = spellData.LastStand.HighestRank()
var lastStandBuff = spellData.LastStandTriggered.HighestRank()

func (warrior *Warrior) registerLastStand() {
	if !warrior.Talents.LastStand {
		return
	}

	actionID := core.ActionID{SpellID: lastStandRank.SpellID}
	healthMetrics := warrior.NewHealthMetrics(actionID)

	var bonusHealth float64
	aura := warrior.RegisterAura(core.Aura{
		Label:    "Last Stand",
		ActionID: actionID,
		Duration: lastStandBuff.Duration,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			bonusHealth = warrior.MaxHealth() * lastStandBuff.Effect(shared.A_MOD_MAX_HEALTH, 0).Fraction()
			warrior.UpdateMaxHealth(sim, bonusHealth, healthMetrics)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.UpdateMaxHealth(sim, -bonusHealth, healthMetrics)
		},
	})

	spell := warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: SpellMaskLastStand,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: lastStandRank.Cooldown,
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
		IntValue:  int32(spellData.ImprovedSunderArmor.TenthsAt(warrior.Talents.ImprovedSunderArmor)),
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

var concussionBlowRank = spellData.ConcussionBlow.HighestRank()

func (warrior *Warrior) registerConcussionBlow() {
	if !warrior.Talents.ConcussionBlow {
		return
	}

	warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: concussionBlowRank.SpellID},
		ClassSpellMask: SpellMaskConcussionBlow,
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		MaxRange:       core.MaxMeleeRange,

		RageCost: core.RageCostOptions{
			Cost:   concussionBlowRank.Cost,
			Refund: concussionBlowRank.MissRefund(),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: concussionBlowRank.Cooldown,
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

// TODO: Manual review needed -- spell 23922 states only "a very high amount of threat"; none is
// modelled until measured in game.
var shieldSlamRank = spellData.ShieldSlam.HighestRank()

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
			Refund: shieldSlamRank.MissRefund(),
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
		ClassMask: SpellMaskOffensiveAbilities,
		Kind:      core.SpellMod_PowerCost_Flat,
		IntValue:  int32(spellData.FocusedRage.TenthsAt(warrior.Talents.FocusedRage)),
	})
}

var masterOfDefenseEnergize = spellData.MasterOfDefenseTriggered.HighestRank()

func (warrior *Warrior) registerMasterOfDefense() {
	if warrior.Talents.MasterOfDefense == 0 {
		return
	}

	// The tooltip states the chance as $m1%, so the talent's ladder is the chance and the 100 in
	// the proc chance column is noise; a shield has to be equipped.
	warrior.registerRageOnAvoid("Master of Defense", masterOfDefenseEnergize.SpellID,
		shared.SpellDataMin(masterOfDefenseEnergize.Energize)/10,
		spellData.MasterOfDefense.FractionAt(warrior.Talents.MasterOfDefense), core.OutcomeDodge|core.OutcomeParry,
		func() bool { return warrior.PseudoStats.CanBlock })
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

var improvedShieldBashSilence = spellData.ImprovedShieldBashTriggered.HighestRank()

func (warrior *Warrior) registerImprovedShieldBash() {
	if warrior.Talents.ImprovedShieldBash == 0 {
		return
	}

	// TODO: nothing in the sim reads a silence on an enemy, so the aura only shows up in metrics.
	silenceAuras := warrior.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    "Shield Bash - Silence",
			ActionID: core.ActionID{SpellID: improvedShieldBashSilence.SpellID},
			Duration: improvedShieldBashSilence.Duration,
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
		IntValue:  int32(spellData.ImprovedThunderClap.TenthsAt(warrior.Talents.ImprovedThunderClap)),
	})
}
