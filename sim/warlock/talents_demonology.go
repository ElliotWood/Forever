package warlock

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

func (warlock *Warlock) registerDemonologyTalents() {
	// Tier 1
	warlock.applyImprovedHealthFunnel()
	warlock.applyImprovedImp()
	warlock.applyDemonicEmbrace()
	warlock.applyUnholyPower()

	// Tier 2
	// Demonic Aegis: armors.go
	warlock.applyImprovedVoidwalker()
	warlock.applyFelVitality()
	warlock.applyDemonicEnergies()

	// Tier 3
	warlock.applyImprovedSayaad()
	warlock.applyDemonicSacrifice()
	warlock.applyMasterSummoner()

	// Tier 4
	warlock.applyDecimation()
	warlock.applyFelDomination()
	warlock.applyDemonicBrand()

	// Tier 5
	warlock.applyImprovedFelhunter()
	warlock.applySoulLink()
	warlock.applyDemonicKnowledge()

	// Tier 6
	warlock.applyMasterDemonologist()

	// Tier 7
	warlock.applyDemonicPact()
}

// The Firebolt half of 18694; its first effect carries the same ladder for Blood Pact, which the
// raid buff handles.
func (warlock *Warlock) applyImprovedImp() {
	if warlock.Talents.ImprovedImp == 0 || warlock.Options.SacrificeSummon {
		return
	}

	warlock.Imp.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: spellData.ImprovedImp.EffectAt(2).FractionAt(warlock.Talents.ImprovedImp),
		ClassMask:  WarlockSpellImpFireBolt,
	})
}

// Forever drops Classic's spirit penalty: 18697 only raises stamina.
func (warlock *Warlock) applyDemonicEmbrace() {
	if warlock.Talents.DemonicEmbrace == 0 {
		return
	}

	warlock.MultiplyStat(stats.Stamina, spellData.DemonicEmbrace.EffectAt(1).MultiplierAt(warlock.Talents.DemonicEmbrace))
}

// 2% more pet damage a point (18769).
func (warlock *Warlock) applyUnholyPower() {
	if warlock.Talents.UnholyPower == 0 || warlock.Options.SacrificeSummon {
		return
	}

	multiplier := spellData.UnholyPower.MultiplierAt(warlock.Talents.UnholyPower)
	for _, pet := range warlock.BasePets {
		pet.PseudoStats.DamageDealtMultiplier *= multiplier
	}
}

// 5% more mana for the warlock and 5% more health and mana for the demon, a point (18731).
func (warlock *Warlock) applyFelVitality() {
	if warlock.Talents.FelVitality == 0 {
		return
	}

	multiplier := spellData.FelVitality.EffectAt(1).MultiplierAt(warlock.Talents.FelVitality)
	warlock.MultiplyStat(stats.Mana, multiplier)
	for _, pet := range warlock.BasePets {
		pet.MultiplyStat(stats.Health, multiplier)
		pet.MultiplyStat(stats.Mana, multiplier)
	}
}

// 10% more Lash of Pain damage a point (18754, second effect).
func (warlock *Warlock) applyImprovedSayaad() {
	if warlock.Talents.ImprovedSayaad == 0 || warlock.Options.SacrificeSummon {
		return
	}

	warlock.Succubus.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: spellData.ImprovedSayaad.EffectAt(2).FractionAt(warlock.Talents.ImprovedSayaad),
		ClassMask:  WarlockSpellSuccubusLashOfPain,
	})
}

// Each demon leaves behind the opposing aspect, and Forever's pairing is the reverse of Classic's:
// the Imp leaves Shadow damage (18789, school mask 32), the Succubus Fire (18791, mask 4), the
// Voidwalker mana (18792) and the Felhunter health (18790). The demon is sacrificed before the pull,
// so the buff is simply permanent and no pet is ever summoned.
func (warlock *Warlock) applyDemonicSacrifice() {
	if !warlock.Talents.DemonicSacrifice || !warlock.Options.SacrificeSummon {
		return
	}

	var spellID int32
	var school stats.SchoolIndex
	switch warlock.Options.Summon {
	case proto.WarlockOptions_Imp:
		spellID, school = 18789, stats.SchoolIndexShadow
	case proto.WarlockOptions_Succubus:
		spellID, school = 18791, stats.SchoolIndexFire
	default:
		// The Voidwalker's mana and the Felhunter's health are regeneration, not damage; they are
		// left out until the sim needs them.
		return
	}

	row := spellData.DemonicSacrificeTriggered.ByID(spellID)
	multiplier := 1 + row.EffectN(1).Percent()

	core.MakePermanent(warlock.RegisterAura(core.Aura{
		Label:    "Demonic Sacrifice",
		ActionID: core.ActionID{SpellID: spellID},
		Duration: row.Duration(),
	}).AttachMultiplicativePseudoStatBuff(&warlock.PseudoStats.SchoolDamageDealtMultiplier[school], multiplier))
}

// Shadow Bolt and Searing Pain hit 3% harder a point below 35% health, and Soul Fire casts 20% a
// point faster and comes off cooldown 45% a point sooner (440870 / 440873).
func (warlock *Warlock) applyDecimation() {
	if warlock.Talents.Decimation == 0 {
		return
	}

	points := warlock.Talents.Decimation
	triggered := spellData.DecimationTriggered.Highest()

	warlock.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_Cooldown_Multiplier,
		FloatValue: 1 + spellData.Decimation.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_COOLDOWN)).FractionAt(points),
		ClassMask:  WarlockSpellSoulFire,
	})

	warlock.DecimationAura = warlock.RegisterAura(core.Aura{
		Label:    "Decimation",
		ActionID: core.ActionID{SpellID: triggered.ID},
		Duration: triggered.Duration(),
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: spellData.Decimation.EffectAt(4).FractionAt(points),
		ClassMask:  WarlockSpellShadowBolt | WarlockSpellSearingPain,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_CastTime_Pct,
		FloatValue: spellData.Decimation.EffectAt(1).FractionAt(points),
		ClassMask:  WarlockSpellSoulFire,
	})

	warlock.MakeProcTriggerAura(core.ProcTrigger{
		Name:               "Decimation Trigger",
		Callback:           core.CallbackOnSpellHitDealt,
		ClassSpellMask:     WarlockSpellShadowBolt | WarlockSpellSearingPain,
		Outcome:            core.OutcomeLanded,
		TriggerImmediately: true,
		ExtraCondition: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) bool {
			return sim.IsExecutePhase35()
		},
		Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
			warlock.DecimationAura.Activate(sim)
		},
	})
}

// Searing Pain sheds 17/33/50% of its threat and arms the demon with 2/4/6 branded attacks
// (1293695 / 1293696).
//
// TODO: the client writes the pet hit as a $<minDam> to $<maxDam> formula the exported tables do not
// carry, so the 39 to 42 from the BlizzCon tooltip is kept.
func (warlock *Warlock) applyDemonicBrand() {
	if warlock.Talents.DemonicBrand == 0 {
		return
	}

	points := warlock.Talents.DemonicBrand
	triggered := spellData.DemonicBrandTriggered.Highest()
	actionID := core.ActionID{SpellID: triggered.ID}

	warlock.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_ThreatMultiplier_Pct,
		FloatValue: spellData.DemonicBrand.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_THREAT)).FractionAt(points),
		ClassMask:  WarlockSpellSearingPain,
	})

	if warlock.Options.SacrificeSummon {
		return
	}

	charges := int32(spellData.DemonicBrand.Effect(dbcenums.A_ADD_FLAT_MODIFIER, int32(dbcenums.SPELLMOD_CHARGES)).ValueAt(points))

	for _, pet := range warlock.BasePets {
		brandSpell := pet.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolShadow,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagPassiveSpell | core.SpellFlagNoOnCastComplete,

			DamageMultiplier: 1,
			ThreatMultiplier: 3,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, sim.Roll(39, 42), spell.OutcomeMagicHit)
			},
		})

		pet.DemonicBrandAura = pet.RegisterAura(core.Aura{
			Label:     "Demonic Brand",
			ActionID:  actionID,
			Duration:  triggered.Duration(),
			MaxStacks: charges,
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if result.Landed() && spell.ProcMask.Matches(core.ProcMaskMelee) {
					brandSpell.Cast(sim, result.Target)
					aura.RemoveStack(sim)
				}
			},
		})
	}

	warlock.MakeProcTriggerAura(core.ProcTrigger{
		Name:               "Demonic Brand Trigger",
		Callback:           core.CallbackOnSpellHitDealt,
		ClassSpellMask:     WarlockSpellSearingPain,
		Outcome:            core.OutcomeLanded,
		TriggerImmediately: true,
		Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
			if warlock.ActivePet == nil {
				return
			}
			brandAura := warlock.ActivePet.DemonicBrandAura
			brandAura.Activate(sim)
			brandAura.SetStacks(sim, brandAura.MaxStacks)
		},
	})
}

// 3% more damage dealt and 30% of the damage taken split with the demon (25228).
func (warlock *Warlock) applySoulLink() {
	if !warlock.Talents.SoulLink || warlock.Options.SacrificeSummon {
		return
	}

	row := spellData.SoulLinkTriggered.ByID(25228)
	damageDealt := 1 + row.EffectN(1).Percent()
	damageTaken := 1 - row.EffectN(2).Percent()

	config := func(unit *core.Unit) core.Aura {
		return core.Aura{
			Label:    "Soul Link",
			ActionID: core.ActionID{SpellID: 19028},
			Duration: core.NeverExpires,
			OnGain: func(aura *core.Aura, _ *core.Simulation) {
				aura.Unit.PseudoStats.DamageDealtMultiplier *= damageDealt
				aura.Unit.PseudoStats.DamageTakenMultiplier *= damageTaken
			},
			OnExpire: func(aura *core.Aura, _ *core.Simulation) {
				aura.Unit.PseudoStats.DamageDealtMultiplier /= damageDealt
				aura.Unit.PseudoStats.DamageTakenMultiplier /= damageTaken
			},
		}
	}

	warlock.SoulLinkAura = core.MakePermanent(warlock.RegisterAura(config(&warlock.Unit)))
	for _, pet := range warlock.BasePets {
		pet.SoulLinkAura = core.MakePermanent(pet.RegisterAura(config(&pet.Unit)))
	}
}

// The demon lends the warlock 33/67/100% of its level in spell power while it is out (412732).
func (warlock *Warlock) applyDemonicKnowledge() {
	if warlock.Talents.DemonicKnowledge == 0 || warlock.Options.SacrificeSummon {
		return
	}

	bonus := spellData.DemonicKnowledge.FractionAt(warlock.Talents.DemonicKnowledge) * float64(core.CharacterLevel)

	core.MakePermanent(warlock.RegisterAura(core.Aura{
		Label:    "Demonic Knowledge",
		ActionID: core.ActionID{SpellID: 412732},
		Duration: core.NeverExpires,
	}).AttachStatBuff(stats.SpellDamage, bonus))
}

// 2% a point, on the school the demon out matches (23785): Fire for the Imp, Shadow for the
// Succubus, damage taken for the Voidwalker and the Felhunter.
func (warlock *Warlock) applyMasterDemonologist() {
	if warlock.Talents.MasterDemonologist == 0 || warlock.Options.SacrificeSummon {
		return
	}

	fraction := spellData.MasterDemonologist.EffectAt(1).FractionAt(warlock.Talents.MasterDemonologist)

	var buff *core.Aura
	switch warlock.Options.Summon {
	case proto.WarlockOptions_Imp:
		buff = warlock.RegisterAura(core.Aura{
			Label:    "Master Demonologist (Imp)",
			ActionID: core.ActionID{SpellID: 23785, Tag: 1},
			Duration: core.NeverExpires,
		}).AttachMultiplicativePseudoStatBuff(&warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire], 1+fraction)
	case proto.WarlockOptions_Succubus:
		buff = warlock.RegisterAura(core.Aura{
			Label:    "Master Demonologist (Succubus)",
			ActionID: core.ActionID{SpellID: 23785, Tag: 3},
			Duration: core.NeverExpires,
		}).AttachMultiplicativePseudoStatBuff(&warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow], 1+fraction)
	default:
		// The Voidwalker's and the Felhunter's halves only cut damage taken.
		return
	}

	warlock.MasterDemonologistAura = core.MakePermanent(buff)
}

// applyImprovedHealthFunnel implements Improved Health Funnel, new in Forever.
//
// Health Funnel is not modelled, so its cost, threat and healing bonuses have nothing to act on.
func (warlock *Warlock) applyImprovedHealthFunnel() {
	if warlock.Talents.ImprovedHealthFunnel == 0 {
		return
	}
}

// applyImprovedVoidwalker implements Improved Voidwalker, new in Forever.
//
// 10% a point on the Voidwalker's Torment and Sacrifice, neither of which is modelled.
func (warlock *Warlock) applyImprovedVoidwalker() {
	if warlock.Talents.ImprovedVoidwalker == 0 {
		return
	}
}

// applyDemonicEnergies implements Demonic Energies, new in Forever.
//
// The pet's share of Life Tap is in lifetap.go.
// TODO: the talent's first effect (8/15) is unidentified.
func (warlock *Warlock) applyDemonicEnergies() {
	if warlock.Talents.DemonicEnergies == 0 {
		return
	}
}

// applyMasterSummoner implements Master Summoner, new in Forever.
//
// The summon spells are not modelled - the demon is out from the start - so the cast time and cost
// cuts have nothing to act on.
func (warlock *Warlock) applyMasterSummoner() {
	if warlock.Talents.MasterSummoner == 0 {
		return
	}
}

// applyFelDomination implements Fel Domination, new in Forever.
//
// A summon cooldown; nothing to model while the demon never has to be resummoned.
func (warlock *Warlock) applyFelDomination() {
	if !warlock.Talents.FelDomination {
		return
	}
}

// applyImprovedFelhunter implements Improved Felhunter, new in Forever.
//
// 10% a point on the Felhunter's abilities, none of which are modelled.
func (warlock *Warlock) applyImprovedFelhunter() {
	if warlock.Talents.ImprovedFelhunter == 0 {
		return
	}
}

// applyDemonicPact implements Demonic Pact, new in Forever.
//
// It keeps a Demonic Sacrifice buff alive while another demon is out, which only matters once the
// sim summons and sacrifices during the fight.
func (warlock *Warlock) applyDemonicPact() {
	if !warlock.Talents.DemonicPact {
		return
	}
}
