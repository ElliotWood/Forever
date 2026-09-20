package warrior

import (
	"time"

	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/stats"
)

func (war *Warrior) registerProtectionTalents() {
	// Tier 1
	// Improved Bloodrage implemented in bloodrage.go
	war.registerAnticipation()

	// Tier 2
	war.registerShieldSpecialization()
	war.registerToughness()

	// Tier 3
	war.registerLastStand()
	// Improved Revenge not implemented
	war.registerDefiance()

	// Tier 4
	war.registerImprovedSunderArmor()
	// Improved Disarm not implemented
	// Improved Taunt not implemented

	// Tier 5
	war.registerImprovedShieldWall()
	war.registerConcussionBlow()
	// Improved Shield Bash not implemented

	// Tier 6

	// Tier 7
	war.registerShieldSlam()
	war.registerFocusedRage()

	// Tier 8

	// Tier 9

	// Forever additions, not yet implemented.
	war.registerMasterOfDefense()
	war.registerImprovedRevenge()
	war.registerImprovedDisarm()
	war.registerVanguard()
	war.registerImprovedShieldBash()
	war.registerBastion()
}

func (war *Warrior) registerDefiance() {
	if war.Talents.Defiance == 0 {
		return
	}

	war.AddStat(stats.ExpertiseRating, spellData.Defiance.Effect(shared.A_ADD_FLAT_MODIFIER, shared.SPELLMOD_ALL_EFFECTS).ValueAt(war.Talents.Defiance)*core.ExpertisePerQuarterPercentReduction)
	war.OnSpellRegistered(func(spell *core.Spell) {
		if !spell.Matches(SpellMaskDefensiveStance) {
			return
		}
		spell.RelatedSelfBuff.
			AttachMultiplicativePseudoStatBuff(&war.PseudoStats.ThreatMultiplier, spellData.Defiance.Effect(shared.A_MOD_THREAT, 127).MultiplierAt(war.Talents.Defiance))
	})
}

func (war *Warrior) registerAnticipation() {
	if war.Talents.Anticipation == 0 {
		return
	}

	war.AddStat(stats.DefenseRating, spellData.Anticipation.ValueAt(war.Talents.Anticipation)*core.DefenseRatingPerDefenseLevel)
}

func (war *Warrior) registerShieldSpecialization() {
	if war.Talents.ShieldSpecialization == 0 {
		return
	}

	war.AddStat(stats.BlockPercent, spellData.ShieldSpecialization.Effect(shared.A_MOD_BLOCK_PERCENT, 0).FractionAt(war.Talents.ShieldSpecialization))

	rageMetrics := war.NewRageMetrics(core.ActionID{SpellID: 23602})

	war.MakeProcTriggerAura(core.ProcTrigger{
		Name: "Shield Specialization",
		// Effect 0 is the block bonus; effect 1 is the proc chance. ProcChanceAt reads a flat 100%.
		ProcChance:         spellData.ShieldSpecialization.EffectAt(1).FractionAt(war.Talents.ShieldSpecialization),
		TriggerImmediately: true,
		Outcome:            core.OutcomeBlock,
		Callback:           core.CallbackOnSpellHitTaken,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			war.AddRage(sim, 1, rageMetrics)
		},
	})
}

func (war *Warrior) registerToughness() {
	if war.Talents.Toughness == 0 {
		return
	}

	war.MultiplyStat(stats.Armor, spellData.Toughness.MultiplierAt(war.Talents.Toughness))
}

func (war *Warrior) registerLastStand() {
	if !war.Talents.LastStand {
		return
	}

	actionID := core.ActionID{SpellID: 12975}
	healthMetrics := war.NewHealthMetrics(actionID)

	var bonusHealth float64
	aura := war.RegisterAura(core.Aura{
		Label:    "Last Stand",
		ActionID: actionID,
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			bonusHealth = war.MaxHealth() * 0.3
			war.AddStatsDynamic(sim, stats.Stats{stats.Health: bonusHealth})
			war.GainHealth(sim, bonusHealth, healthMetrics)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			war.AddStatsDynamic(sim, stats.Stats{stats.Health: -bonusHealth})
		},
	})

	spell := war.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: SpellMaskLastStand,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    war.NewTimer(),
				Duration: time.Minute * 8,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			aura.Activate(sim)
		},

		RelatedSelfBuff: aura,
	})

	war.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeSurvival,
		BuffAura: &core.StatBuffAura{
			Aura:            aura,
			BuffedStatTypes: []stats.Stat{stats.Health},
		},
	})
}

func (war *Warrior) registerImprovedSunderArmor() {
	if war.Talents.ImprovedSunderArmor == 0 {
		return
	}

	// Retained rage when swapping stances implemented in stances.go
	war.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskSunderArmor | SpellMaskDevastate,
		Kind:      core.SpellMod_PowerCost_Flat,
		// TODO: this read war.Talents.TacticalMastery, which looks like a long-standing
		// copy-paste bug -- the registrar guards ImprovedSunderArmor. Forever drops
		// Tactical Mastery entirely, so it now scales off its own talent; the per-rank
		// rage reduction needs confirming against the Forever tooltip.
		IntValue: -war.Talents.ImprovedSunderArmor,
	})
}

func (war *Warrior) registerImprovedShieldWall() {
	if war.Talents.ImprovedShieldWall == 0 {
		return
	}

	duration := []time.Duration{0, 3, 5}[war.Talents.ImprovedShieldWall]

	war.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskShieldWall,
		Kind:      core.SpellMod_BuffDuration_Flat,
		TimeValue: time.Second * duration,
	})
}

func (war *Warrior) registerConcussionBlow() {
	if !war.Talents.ConcussionBlow {
		return
	}

	war.RegisterSpell(core.SpellConfig{
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
				Timer:    war.NewTimer(),
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

var shieldSlamRank = spellData.ShieldSlam.BySpellID(30356)
var devastateRank = spellData.Devastate.BySpellID(30022)

func (war *Warrior) registerShieldSlam() {
	if !war.Talents.ShieldSlam {
		return
	}

	war.RegisterSpell(core.SpellConfig{
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
				Timer:    war.NewTimer(),
				Duration: shieldSlamRank.Cooldown,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return war.PseudoStats.CanBlock
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		FlatThreatBonus:  305,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := shieldSlamRank.Direct.Damage(sim) + war.BlockDamageReduction()
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
}

func (war *Warrior) registerFocusedRage() {
	if war.Talents.FocusedRage == 0 {
		return
	}

	war.AddStaticMod(core.SpellModConfig{
		ClassMask: WarriorSpellsAll ^ (SpellMaskRampage | SpellMaskDeathWish | SpellMaskBattleShout | SpellMaskCommandingShout),
		Kind:      core.SpellMod_PowerCost_Flat,
		IntValue:  -war.Talents.FocusedRage,
	})
}

// registerMasterOfDefense implements Master of Defense, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (war *Warrior) registerMasterOfDefense() {
	if war.Talents.MasterOfDefense == 0 {
		return
	}

	panic("To be implemented")
}

// registerImprovedRevenge implements Improved Revenge, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (war *Warrior) registerImprovedRevenge() {
	if war.Talents.ImprovedRevenge == 0 {
		return
	}

	panic("To be implemented")
}

// registerImprovedDisarm implements Improved Disarm, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (war *Warrior) registerImprovedDisarm() {
	if war.Talents.ImprovedDisarm == 0 {
		return
	}

	panic("To be implemented")
}

// registerVanguard implements Vanguard, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (war *Warrior) registerVanguard() {
	if !war.Talents.Vanguard {
		return
	}

	panic("To be implemented")
}

// registerImprovedShieldBash implements Improved Shield Bash, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (war *Warrior) registerImprovedShieldBash() {
	if war.Talents.ImprovedShieldBash == 0 {
		return
	}

	panic("To be implemented")
}

// registerBastion implements Bastion, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (war *Warrior) registerBastion() {
	if war.Talents.Bastion == 0 {
		return
	}

	panic("To be implemented")
}
