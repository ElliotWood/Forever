package warrior

import (
	"time"

	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/stats"
)

// TODO: Manual review needed -- this was modelled during the Forever port, not carried
// over unchanged, so its numbers and shape want checking against the client.
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

// TODO: Manual review needed -- this was modelled during the Forever port, not carried
// over unchanged, so its numbers and shape want checking against the client.
func (warrior *Warrior) registerDefiance() {
	if warrior.Talents.Defiance == 0 {
		return
	}

	// TODO: Forever drops Defiance's expertise; the spell carries only the threat modifier
	// applied below (A_MOD_THREAT, +5/10/15%), so expertise is pinned to the untalented 0.
	expertiseBonus := 0.0
	warrior.AddStat(stats.ExpertiseRating, expertiseBonus)
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

// TODO: Manual review needed -- this was modelled during the Forever port, not carried
// over unchanged, so its numbers and shape want checking against the client.
func (warrior *Warrior) registerShieldSpecialization() {
	if warrior.Talents.ShieldSpecialization == 0 {
		return
	}

	warrior.AddStat(stats.BlockPercent, spellData.ShieldSpecialization.Effect(shared.A_MOD_BLOCK_PERCENT, 0).FractionAt(warrior.Talents.ShieldSpecialization))

	rageMetrics := warrior.NewRageMetrics(core.ActionID{SpellID: 23602})

	warrior.MakeProcTriggerAura(core.ProcTrigger{
		Name: "Shield Specialization",
		// Effect 0 is the block bonus; effect 1 is the proc chance. ProcChanceAt reads a flat 100%.
		ProcChance:         spellData.ShieldSpecialization.EffectAt(1).FractionAt(warrior.Talents.ShieldSpecialization),
		TriggerImmediately: true,
		Outcome:            core.OutcomeBlock,
		Callback:           core.CallbackOnSpellHitTaken,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			warrior.AddRage(sim, 1, rageMetrics)
		},
	})
}

// TODO: Manual review needed -- this was modelled during the Forever port, not carried
// over unchanged, so its numbers and shape want checking against the client.
func (warrior *Warrior) registerToughness() {
	if warrior.Talents.Toughness == 0 {
		return
	}

	// The bonus-armor effect carries the same ladder; this multiplies base armor only.
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
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
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
				Timer:    warrior.NewTimer(),
				Duration: time.Minute * 8,
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

// TODO: Manual review needed -- this was modelled during the Forever port, not carried
// over unchanged, so its numbers and shape want checking against the client.
func (warrior *Warrior) registerImprovedSunderArmor() {
	if warrior.Talents.ImprovedSunderArmor == 0 {
		return
	}

	// Retained rage when swapping stances implemented in stances.go
	warrior.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskSunderArmor,
		Kind:      core.SpellMod_PowerCost_Flat,
		// TODO: this read warrior.Talents.TacticalMastery, which looks like a long-standing
		// copy-paste bug -- the registrar guards ImprovedSunderArmor. Forever drops
		// Tactical Mastery entirely, so it now scales off its own talent; the per-rank
		// rage reduction needs confirming against the Forever tooltip.
		IntValue: -warrior.Talents.ImprovedSunderArmor,
	})
}

func (warrior *Warrior) registerImprovedShieldWall() {
	if warrior.Talents.ImprovedShieldWall == 0 {
		return
	}

	duration := []time.Duration{0, 3, 5}[warrior.Talents.ImprovedShieldWall]

	warrior.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskShieldWall,
		Kind:      core.SpellMod_BuffDuration_Flat,
		TimeValue: time.Second * duration,
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
			Cost:   15,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
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
		FlatThreatBonus:  305,

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
		IntValue:  -warrior.Talents.FocusedRage,
	})
}

// TODO: registerMasterOfDefense models nothing yet; spellData.MasterOfDefense carries the ranks.
func (warrior *Warrior) registerMasterOfDefense() {
	if warrior.Talents.MasterOfDefense == 0 {
		return
	}
}

// TODO: registerImprovedRevenge models nothing yet; spellData.ImprovedRevenge carries the ranks.
func (warrior *Warrior) registerImprovedRevenge() {
	if warrior.Talents.ImprovedRevenge == 0 {
		return
	}
}

// TODO: registerImprovedDisarm models nothing yet; spellData.ImprovedDisarm carries the ranks.
func (warrior *Warrior) registerImprovedDisarm() {
	if warrior.Talents.ImprovedDisarm == 0 {
		return
	}
}

// TODO: registerVanguard models nothing yet; the client has no ladder for it, so the tooltip is the source.
func (warrior *Warrior) registerVanguard() {
	if !warrior.Talents.Vanguard {
		return
	}
}

// TODO: registerImprovedShieldBash models nothing yet; spellData.ImprovedShieldBash carries the ranks.
func (warrior *Warrior) registerImprovedShieldBash() {
	if warrior.Talents.ImprovedShieldBash == 0 {
		return
	}
}

// TODO: registerBastion models nothing yet; spellData.Bastion carries the ranks.
func (warrior *Warrior) registerBastion() {
	if warrior.Talents.Bastion == 0 {
		return
	}
}

func (warrior *Warrior) registerImprovedThunderClap() {
	if warrior.Talents.ImprovedThunderClap == 0 {
		return
	}

	// Slowing effect implemented in core/debuffs.go

	rageCostReduction := []int32{0, 1, 2, 4}[warrior.Talents.ImprovedThunderClap]
	damageGain := []float64{0, 0.4, 0.7, 1.0}[warrior.Talents.ImprovedThunderClap]

	warrior.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskThunderClap,
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: damageGain,
	})

	warrior.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskThunderClap,
		Kind:      core.SpellMod_PowerCost_Flat,
		IntValue:  -rageCostReduction,
	})
}
