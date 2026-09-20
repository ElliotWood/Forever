package warrior

import (
	"time"

	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

// TODO: Manual review needed -- this was modelled during the Forever port, not carried
// over unchanged, so its numbers and shape want checking against the client.
func (war *Warrior) registerArmsTalents() {
	// Tier 1
	war.registerImprovedHeroicStrike()
	war.registerDeflection()
	war.registerImprovedRend()

	// Tier 2
	war.registerImprovedCharge()
	// Iron Will not implemented
	war.registerImprovedThunderClap()

	// Tier 3
	war.registerImprovedOverpower()
	war.registerAngerManagement()
	war.registerDeepWounds()

	// Tier 4
	war.registerTwoHandedWeaponSpecialization()
	war.registerImpale()

	// Tier 5
	war.registerDeathWish()

	// Tier 6
	war.registerImprovedIntercept()
	// Improved Hamstring not implemented

	// Tier 7
	war.registerMortalStrike()
	// Second Wind not implemented

	// Tier 8

	// Tier 9

	// Forever additions, not yet implemented.
	war.registerImprovedTacticalMastery()
	war.registerSpearingStrike()
	war.registerBloodthrill()
	war.registerWeaponmaster()
	war.registerImprovedHamstring()
}

/*
 * Arms
 */
func (war *Warrior) registerImprovedHeroicStrike() {
	if war.Talents.ImprovedHeroicStrike == 0 {
		return
	}

	war.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskHeroicStrike,
		Kind:      core.SpellMod_PowerCost_Flat,
		IntValue:  -war.Talents.ImprovedHeroicStrike,
	})
}
func (war *Warrior) registerDeflection() {
	if war.Talents.Deflection == 0 {
		return
	}

	war.PseudoStats.BaseParryChance += spellData.Deflection.FractionAt(war.Talents.Deflection)
}

func (war *Warrior) registerImprovedRend() {
	if war.Talents.ImprovedRend == 0 {
		return
	}

	war.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskRend,
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: spellData.ImprovedRend.FractionAt(war.Talents.ImprovedRend),
	})
}

func (war *Warrior) registerImprovedCharge() {
	if war.Talents.ImprovedCharge == 0 {
		return
	}

	war.ChargeRageGain += 3.0 + float64(war.Talents.ImprovedCharge)
}

func (war *Warrior) registerImprovedThunderClap() {
	if war.Talents.ImprovedThunderClap == 0 {
		return
	}

	// Slowing effect implemented in core/debuffs.go

	rageCostReduction := []int32{0, 1, 2, 4}[war.Talents.ImprovedThunderClap]
	damageGain := []float64{0, 0.4, 0.7, 1.0}[war.Talents.ImprovedThunderClap]

	war.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskThunderClap,
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: damageGain,
	})

	war.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskThunderClap,
		Kind:      core.SpellMod_PowerCost_Flat,
		IntValue:  -rageCostReduction,
	})
}

func (war *Warrior) registerImprovedOverpower() {
	if war.Talents.ImprovedOverpower == 0 {
		return
	}

	core.MakePermanent(war.RegisterAura(core.Aura{
		Label:    "Improved Overpower",
		ActionID: core.ActionID{SpellID: 12963}.WithTag(war.Talents.ImprovedOverpower),
	})).AttachSpellMod(core.SpellModConfig{
		ClassMask:  SpellMaskOverpower,
		Kind:       core.SpellMod_BonusCrit_Percent,
		FloatValue: spellData.ImprovedOverpower.ValueAt(war.Talents.ImprovedOverpower),
	})
}

func (war *Warrior) registerAngerManagement() {
	if !war.Talents.AngerManagement {
		return
	}

	rageMetrics := war.NewRageMetrics(core.ActionID{SpellID: 12296})

	war.RegisterResetEffect(func(sim *core.Simulation) {
		core.StartPeriodicAction(sim, core.PeriodicActionOptions{
			Period: time.Second * 3,
			OnAction: func(sim *core.Simulation) {
				if sim.CurrentTime > 0 {
					war.AddRage(sim, 1, rageMetrics)
				}
			},
		})
	})
}

func (war *Warrior) registerDeepWounds() {
	if war.Talents.DeepWounds == 0 {
		return
	}

	war.DeepWounds = war.RegisterSpell(core.SpellConfig{
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
				baseDamage := war.AutoAttacks.MH().CalculateAverageWeaponDamage(dot.Spell.MeleeAttackPower(target))
				dot.SnapshotPhysical(target, baseDamage/float64(dot.HastedTickCount())*0.2*float64(war.Talents.DeepWounds))
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

	war.MakeProcTriggerAura(core.ProcTrigger{
		Name:               "Deep Wounds - Trigger",
		TriggerImmediately: true,
		ProcMaskExclude:    core.ProcMaskEmpty,
		Outcome:            core.OutcomeCrit,
		Callback:           core.CallbackOnSpellHitDealt,
		ExtraCondition: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) bool {
			return spell.SpellSchool.Matches(core.SpellSchoolPhysical)
		},
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			war.DeepWounds.Cast(sim, result.Target)
		},
	})

}

func (war *Warrior) registerTwoHandedWeaponSpecialization() {
	if war.Talents.TwoHandedWeaponSpecialization == 0 {
		return
	}

	weaponMod := war.AddDynamicMod(core.SpellModConfig{
		ClassMask:  SpellMaskDirectDamageSpells,
		School:     core.SpellSchoolPhysical,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: spellData.TwoHandedWeaponSpecialization.Effect(shared.A_MOD_DAMAGE_PERCENT_DONE, 1).FractionAt(war.Talents.TwoHandedWeaponSpecialization),
	})

	if war.GetMainHandType() == proto.HandType_HandTypeTwoHand {
		weaponMod.Activate()
	}

	war.RegisterItemSwapCallback(core.AllMeleeWeaponSlots(), func(sim *core.Simulation, slot proto.ItemSlot) {
		if war.GetMainHandType() == proto.HandType_HandTypeTwoHand {
			weaponMod.Activate()
		} else {
			weaponMod.Deactivate()
		}
	})
}

func (war *Warrior) registerImpale() {
	if war.Talents.Impale == 0 {
		return
	}

	war.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskDamageSpells,
		Kind:       core.SpellMod_CritMultiplier_Flat,
		FloatValue: spellData.Impale.FractionAt(war.Talents.Impale),
	})
}

func (war *Warrior) registerDeathWish() {
	if !war.Talents.DeathWish {
		return
	}

	actionID := core.ActionID{SpellID: 12292}

	deathWishAura := war.RegisterAura(core.Aura{
		Label:    "Death Wish",
		ActionID: actionID,
		Duration: time.Second * 30,
	}).
		AttachMultiplicativePseudoStatBuff(
			&war.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical], 1.2,
		).
		AttachMultiplicativePseudoStatBuff(
			&war.PseudoStats.DamageTakenMultiplier, 1.05,
		).
		// Grants immunity to Fear effects.
		AttachFearImmunity()

	deathWishSpell := war.RegisterSpell(core.SpellConfig{
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
				Timer:    war.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			deathWishAura.Activate(sim)
			war.WaitUntil(sim, sim.CurrentTime+core.GCDDefault)
		},

		RelatedSelfBuff: deathWishAura,
	})

	war.AddMajorCooldown(core.MajorCooldown{
		Spell: deathWishSpell,
		Type:  core.CooldownTypeDPS,
	})
}

func (war *Warrior) registerImprovedIntercept() {
	if war.Talents.ImprovedIntercept == 0 {
		return
	}

	war.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskIntercept,
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: -time.Second * time.Duration(5*war.Talents.ImprovedIntercept),
	})
}

var mortalStrikeRank = spellData.MortalStrike.HighestRank()
var mortalStrikeBaseDamage, _ = mortalStrikeRank.Direct.Range()

func (war *Warrior) registerMortalStrike() {
	if !war.Talents.MortalStrike {
		return
	}

	war.MortalStrike = war.RegisterSpell(core.SpellConfig{
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
				Timer:    war.NewTimer(),
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

// registerImprovedTacticalMastery implements Improved Tactical Mastery, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (war *Warrior) registerImprovedTacticalMastery() {
	if war.Talents.ImprovedTacticalMastery == 0 {
		return
	}
}

// registerSpearingStrike implements Spearing Strike, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (war *Warrior) registerSpearingStrike() {
	if !war.Talents.SpearingStrike {
		return
	}
}

// registerBloodthrill implements Bloodthrill, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (war *Warrior) registerBloodthrill() {
	if war.Talents.Bloodthrill == 0 {
		return
	}
}

// registerWeaponmaster implements Weaponmaster, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (war *Warrior) registerWeaponmaster() {
	if war.Talents.Weaponmaster == 0 {
		return
	}
}

// registerImprovedHamstring implements Improved Hamstring, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (war *Warrior) registerImprovedHamstring() {
	if war.Talents.ImprovedHamstring == 0 {
		return
	}
}
