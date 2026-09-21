package rogue

import (
	"slices"
	"time"

	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

func (rogue *Rogue) registerAssassinationTalents() {
	// Tier 1
	rogue.registerImprovedEviscerate()
	rogue.registerRemorselessAttacks()
	rogue.registerMalice()

	// Tier 2
	// Ruthlessness implemented in ApplyFinisher
	rogue.registerMurder()
	rogue.registerPuncturingWounds()

	// Tier 3
	// Relentless Strikes implemented in ApplyFinisher
	rogue.registerImprovedExposeArmor()
	rogue.registerLethality()

	// Tier 4
	rogue.registerVilePoisons()
	// Improved Poisons implemented in poisons.go

	// Tier 5
	rogue.registerColdBlood()
	rogue.registerImprovedKidneyShot()

	// Tier 6
	rogue.registerSealFate()

	// Tier 7
	// Vigor implemented in rogue.go
	rogue.registerVenom()

	// Tier 9
	rogue.registerMutilate()
}

func (rogue *Rogue) registerImprovedEviscerate() {
	if rogue.Talents.ImprovedEviscerate == 0 {
		return
	}

	rogue.AddStaticMod(core.SpellModConfig{
		ClassMask:  RogueSpellEviscerate,
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: spellData.ImprovedEviscerate.FractionAt(rogue.Talents.ImprovedEviscerate),
	})
}

func (rogue *Rogue) registerMalice() {
	if rogue.Talents.Malice == 0 {
		return
	}

	// Malice covers Poisons in Forever, and those roll against the spell hit table.
	crit := spellData.Malice.ValueAt(rogue.Talents.Malice)
	rogue.AddStat(stats.PhysicalCritPercent, crit)
	rogue.AddStat(stats.SpellCritPercent, crit)
}

func (rogue *Rogue) registerMurder() {
	if rogue.Talents.Murder == 0 {
		return
	}

	multiplier := spellData.Murder.MultiplierAt(rogue.Talents.Murder)
	rogue.Env.RegisterPostFinalizeEffect(func() {
		for _, at := range rogue.AttackTables {
			if slices.Contains([]proto.MobType{proto.MobType_MobTypeHumanoid, proto.MobType_MobTypeGiant, proto.MobType_MobTypeBeast, proto.MobType_MobTypeDragonkin}, at.Defender.MobType) {
				at.DamageDealtMultiplier *= multiplier
				at.CritMultiplier *= multiplier
			}
		}
	})
}

func (rogue *Rogue) registerPuncturingWounds() {
	if rogue.Talents.PuncturingWounds == 0 {
		return
	}

	// Effect 1 is the combo point trigger, handled in backstab.go. The two crit modifiers share
	// an aura and misc pair, so each has to be named by index: 0 is Backstab, 2 is Mutilate.
	rogue.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_BonusCrit_Percent,
		ClassMask:  RogueSpellBackstab,
		FloatValue: spellData.PuncturingWounds.EffectAt(0).ValueAt(rogue.Talents.PuncturingWounds),
	})
	rogue.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_BonusCrit_Percent,
		ClassMask:  RogueSpellMutilate | RogueSpellMutilateHit,
		FloatValue: spellData.PuncturingWounds.EffectAt(2).ValueAt(rogue.Talents.PuncturingWounds),
	})
}

// Forever repurposes Improved Expose Armor: the energy discount is a SpellMod, and the combo
// points it hands back on a full spend live in expose_armor.go.
func (rogue *Rogue) registerImprovedExposeArmor() {
	if rogue.Talents.ImprovedExposeArmor == 0 {
		return
	}

	rogue.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_PowerCost_Flat,
		ClassMask: RogueSpellExposeArmor,
		IntValue:  int32(spellData.ImprovedExposeArmor.Effect(shared.A_ADD_FLAT_MODIFIER, shared.SPELLMOD_COST).ValueAt(rogue.Talents.ImprovedExposeArmor)),
	})

	// A dummy aura so the APL can ask whether the talent is taken.
	core.MakePermanent(rogue.RegisterAura(core.Aura{
		Label:    "Improved Expose Armor",
		ActionID: core.ActionID{SpellID: spellData.ImprovedExposeArmor.HighestRank().SpellID},
	}))
}

func (rogue *Rogue) registerLethality() {
	if rogue.Talents.Lethality == 0 {
		return
	}

	rogue.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_CritMultiplier_Flat,
		ClassMask:  RogueSpellLethality,
		FloatValue: spellData.Lethality.FractionAt(rogue.Talents.Lethality),
	})
}

func (rogue *Rogue) registerVilePoisons() {
	if rogue.Talents.VilePoisons == 0 {
		return
	}

	rogue.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		ClassMask:  RogueSpellPoisons,
		FloatValue: spellData.VilePoisons.Effect(shared.A_ADD_PCT_MODIFIER, shared.SPELLMOD_DAMAGE).FractionAt(rogue.Talents.VilePoisons),
	})
}

func (rogue *Rogue) registerColdBlood() {
	if !rogue.Talents.ColdBlood {
		return
	}

	coldBloodRank := spellData.ColdBlood.HighestRank()
	actionID := core.ActionID{SpellID: coldBloodRank.SpellID}

	cbAura := rogue.GetOrRegisterAura(core.Aura{
		Label:    "Cold Blood",
		ActionID: actionID,
		Duration: core.NeverExpires,

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(RogueSpellActives) {
				aura.Deactivate(sim)
			}
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_BonusCrit_Percent,
		ClassMask:  RogueSpellActives,
		FloatValue: coldBloodRank.Effect(shared.A_ADD_FLAT_MODIFIER, shared.SPELLMOD_CRITICAL_CHANCE).Value,
	})

	rogue.ColdBlood = rogue.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: RogueSpellColdBlood,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: coldBloodRank.Cooldown,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			cbAura.Activate(sim)
		},
	})

	rogue.AddMajorCooldown(core.MajorCooldown{
		Spell: rogue.ColdBlood,
		Type:  core.CooldownTypeDPS,
	})
}

func (rogue *Rogue) registerSealFate() {
	if rogue.Talents.SealFate == 0 {
		return
	}

	sfMetrics := rogue.NewComboPointMetrics(core.ActionID{SpellID: 14195})

	rogue.MakeProcTriggerAura(core.ProcTrigger{
		Name:     "Seal Fate Trigger",
		ActionID: core.ActionID{SpellID: spellData.SealFate.HighestRank().SpellID},
		// Forever puts the real per-rank chance on the effect; ProcChanceAt reads a flat 100%.
		ProcChance: spellData.SealFate.FractionAt(rogue.Talents.SealFate),
		Callback:   core.CallbackOnSpellHitDealt,
		Outcome:    core.OutcomeCrit,
		SpellFlags: SpellFlagBuilder,
		ICD:        time.Millisecond * 500,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			rogue.AddComboPoints(sim, 1, sfMetrics)
		},
	})
}

// Was 34413; Forever reworked Mutilate onto an entirely new set of spell ids, so this
// follows the highest rank the client actually ships.
var MutilateSpellID int32 = spellData.Mutilate.HighestRank().SpellID

var mutilateRank = spellData.Mutilate.BySpellID(MutilateSpellID)

func (rogue *Rogue) registerMutilate() {
	if !rogue.Talents.Mutilate {
		return
	}

	rogue.MutilateMH = rogue.newMutilateHitSpell(true)
	rogue.MutilateOH = rogue.newMutilateHitSpell(false)

	rogue.Mutilate = rogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: MutilateSpellID, Tag: 0},
		SpellSchool:    mutilateRank.SpellSchool,
		DefenseType:    mutilateRank.DefenseType,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: RogueSpellMutilate,
		MaxRange:       core.MaxMeleeRange,

		EnergyCost: core.EnergyCostOptions{
			Cost:   mutilateRank.Cost,
			Refund: mutilateRank.MissRefund(),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: mutilateRank.GCD,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			// Every rank, and both hand strikes, require a Dagger in the beta client.
			return rogue.HasDagger(core.MainHand) && rogue.HasDagger(core.OffHand)
		},

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHit) // Miss/Dodge/Parry/Hit
			if result.Landed() {
				rogue.AddComboPoints(sim, mutilateComboPoints(), spell.ComboPointMetrics())
				rogue.MutilateOH.Cast(sim, target)
				rogue.MutilateMH.Cast(sim, target)
			} else {
				spell.IssueRefund(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})
}

// The client energizes off the rank's own effect, 2 combo points.
func mutilateComboPoints() int32 {
	min, _ := mutilateRank.Energize.Range()
	return int32(min)
}

// The client splits Mutilate into a parent and two triggered hit spells, and gen_spelldata's
// rank order for the hits does not follow the parent's (rank 1's hits are generated last), so
// the flat damage is keyed on the parent spell id instead: 23/33/48/67 for ranks 1-4.
var mutilateFlatDamage = map[int32]float64{
	1310707: 23,
	399956:  33,
	1241582: 48,
	1241584: 67,
}

func (rogue *Rogue) newMutilateHitSpell(isMH bool) *core.Spell {
	actionID := core.ActionID{SpellID: MutilateSpellID, Tag: 1}
	procMask := core.ProcMaskMeleeMHSpecial
	if !isMH {
		actionID = core.ActionID{SpellID: MutilateSpellID, Tag: 2}
		procMask = core.ProcMaskMeleeOHSpecial
	}

	mutBaseDamage := mutilateFlatDamage[MutilateSpellID]
	// Every hit rank states the same 75% weapon share.
	weaponDamage := spellData.MutilateTriggered.EffectAt(1).ValueAt(1) / 100
	// Mutilate hits harder while one of the rogue's lingering poisons is on the target; the
	// client states the bonus as a percentage on the parent rank's dummy effect.
	poisonBonus := spellData.Mutilate.EffectAt(3).ValueAt(mutilateRank.Rank) / 100

	return rogue.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    mutilateRank.SpellSchool,
		DefenseType:    mutilateRank.DefenseType,
		ProcMask:       procMask,
		Flags:          core.SpellFlagMeleeMetrics | SpellFlagBuilder,
		ClassSpellMask: RogueSpellMutilateHit,

		DamageMultiplier:         weaponDamage,
		DamageMultiplierAdditive: 1,
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := mutBaseDamage
			if isMH {
				baseDamage += spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower(target))
			} else {
				baseDamage += spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower(target))
			}

			oldMultiplier := spell.DamageMultiplier
			if rogue.isPoisoned(target) {
				spell.DamageMultiplier *= 1 + poisonBonus
			}

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			spell.DamageMultiplier = oldMultiplier
		},
	})
}

// registerImprovedKidneyShot implements Improved Kidney Shot, new in Forever.
//
// TODO: Not modelled. The talent adds 5/10% damage taken while Kidney Shot holds the target, and
// Kidney Shot itself is not an ability our sim casts.
func (rogue *Rogue) registerImprovedKidneyShot() {}

// registerRemorselessAttacks implements Remorseless Attacks, new in Forever.
//
// TODO: Not modelled. The client gives the next ability after a killing blow +40% crit, and a
// single-target sim never scores one.
func (rogue *Rogue) registerRemorselessAttacks() {}

// Venom is a finisher that raises poison damage and application chance for a combo point scaled
// duration, on the same 9-21 second ladder as Slice and Dice.
func (rogue *Rogue) registerVenom() {
	if !rogue.Talents.Venom {
		return
	}

	venomRank := spellData.Venom.HighestRank()
	actionID := core.ActionID{SpellID: venomRank.SpellID}

	damageBonus := venomRank.Effect(shared.A_ADD_PCT_MODIFIER, shared.SPELLMOD_DAMAGE).Value / 100
	chanceBonus := venomRank.Effect(shared.A_ADD_FLAT_MODIFIER, shared.SPELLMOD_CHANCE_OF_SUCCESS).Value / 100

	rogue.VenomAura = rogue.RegisterAura(core.Aura{
		Label:    "Venom",
		ActionID: actionID,
		// Overridden on cast; a non-zero default keeps an APL prepull from crashing.
		Duration: rogue.sliceAndDiceDurations[5],
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			rogue.additivePoisonBonusChance += chanceBonus
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.additivePoisonBonusChance -= chanceBonus
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		ClassMask:  RogueSpellPoisons,
		FloatValue: damageBonus,
	})

	rogue.Venom = rogue.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          SpellFlagFinisher | core.SpellFlagAPL,
		MetricSplits:   6,
		ClassSpellMask: RogueSpellVenom,

		EnergyCost: core.EnergyCostOptions{
			Cost: venomRank.Cost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: venomRank.GCD,
			},
			IgnoreHaste: true,
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				spell.SetMetricsSplit(rogue.ComboPoints())
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return rogue.ComboPoints() > 0
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			comboPoints := rogue.ComboPoints()
			rogue.ApplyFinisher(sim, spell)
			spell.RelatedSelfBuff.Deactivate(sim)
			spell.RelatedSelfBuff.Duration = rogue.sliceAndDiceDurations[comboPoints]
			spell.RelatedSelfBuff.Activate(sim)
		},

		RelatedSelfBuff: rogue.VenomAura,
	})
}
