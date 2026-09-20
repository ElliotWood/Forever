package warlock

import (
	"time"

	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

func (warlock *Warlock) applyAfflictionTalents() {
	warlock.applySuppression()
	warlock.applyImprovedCorruption()
	warlock.registerAmplifyCurse()
	warlock.applyNightfall()
	warlock.applyEmpoweredCorruption()
	warlock.applyShadowMastery()

	// Forever additions, not yet implemented.
	warlock.applySoulHarvesting()
	warlock.applyImprovedDrains()
	warlock.applyImprovedBaneOfAgony()
	warlock.applyFelConcentration()
	warlock.applyPandemic()
	warlock.applyMalevolence()
	warlock.applyCurseOfExhaustion()
	warlock.applySiphonLife()
	warlock.applyWrack()
}

func (warlock *Warlock) applyDemonologyTalents() {
	warlock.appyImprovedImp()
	warlock.applyDemonicEmbrace()
	warlock.applyImprovedSayaad()
	warlock.applyUnholyPower()
	warlock.applyDemonicSacrifice()
	warlock.applyMasterDemonologist()
	warlock.applySoulLink()
	warlock.applyDemonicKnowledge()

	// Forever additions, not yet implemented.
	warlock.applyImprovedHealthFunnel()
	warlock.applyImprovedVoidwalker()
	warlock.applyFelVitality()
	warlock.applyDemonicEnergies()
	warlock.applyMasterSummoner()
	warlock.applyDecimation()
	warlock.applyFelDomination()
	warlock.applyDemonicBrand()
	warlock.applyImprovedFelhunter()
	warlock.applyDemonicPact()
}

func (warlock *Warlock) applyDestructionTalents() {
	warlock.applyCataclysm()
	warlock.applyBane()
	warlock.applyShadowburn()
	warlock.applyImprovedShadowBolt()
	warlock.applyDestructiveReach()
	warlock.applyRuin()
	warlock.applyConflagrate()
	warlock.applyShadowAndFlame()

	// Forever additions, not yet implemented.
	warlock.applyMoltenSkin()
	warlock.applyAftermath()
	warlock.applyIntensity()
	warlock.applyAgonizingFlames()
	warlock.applyPyroclasm()
	warlock.applyBaneOfHavoc()
	warlock.applyFireAndBrimstone()
	warlock.applyIncinerate()
}

/*
Affliction
Skipping the following (for now)
- Soul Siphon -> included in drain_life.go
- Improved Life Tap -> included in lifetap.go
- Empowered Corruption -> included in corruption.go
- Siphon Life -> implemented in siphon_life.go
- Fel Concentration
- Grim Reach
- Shadow Embrace -> implemented in corruption.go, curseOfAgony.go, siphon_life.go, and seed.go
- Curse of Weakness
- Curse of Exhaustion
- Dark Pact
- Improved Howl of Terror
*/
func (warlock *Warlock) applySuppression() {
	if warlock.Talents.Suppression == 0 {
		return
	}

	warlock.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_BonusHit_Percent,
		FloatValue: spellData.Suppression.Effect(shared.A_MOD_SPELL_HIT_CHANCE, 0).ValueAt(warlock.Talents.Suppression),
		ClassMask:  WarlockAfflictionSpells,
	})
}

func (warlock *Warlock) applyImprovedCorruption() {
	if warlock.Talents.ImprovedCorruption == 0 {
		return
	}

	warlock.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_CastTime_Flat,
		TimeValue: time.Millisecond * (-400 * time.Duration(warlock.Talents.ImprovedCorruption)),
		ClassMask: WarlockSpellCorruption,
	})
}

func (warlock *Warlock) registerAmplifyCurse() {
	if !warlock.Talents.AmplifyCurse {
		return
	}

	actionID := core.ActionID{SpellID: 18288}

	warlock.AmplifyCurseAura = warlock.GetOrRegisterAura(core.Aura{
		Label:    "Amplify Curse",
		Tag:      "Affliction",
		ActionID: actionID,
		Duration: time.Second * 30,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Matches(WarlockSpellCurseOfAgony | WarlockSpellCurseOfDoom) {
				warlock.AmplifyCurseAura.Deactivate(sim)
			}
		},
	})

	warlock.AmplifyCurse = warlock.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Minute * 3,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			warlock.AmplifyCurseAura.Activate(sim)
		},
	})

	warlock.AddMajorCooldown(core.MajorCooldown{
		Spell: warlock.AmplifyCurse,
		Type:  core.CooldownTypeDPS,
	})
}

func (warlock *Warlock) applyNightfall() {
	if warlock.Talents.Nightfall == 0 {
		return
	}

	warlock.NightfallProcAura = warlock.MakeProcTriggerAura(core.ProcTrigger{
		Name:            "Shadow Trance",
		MetricsActionID: core.ActionID{SpellID: 17941},
		ClassSpellMask:  WarlockSpellShadowBolt,
		Callback:        core.CallbackOnCastComplete,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.CurCast.CastTime != 0 {
				return
			}
			warlock.NightfallProcAura.Deactivate(sim)
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_CastTime_Pct,
		FloatValue: -1.0,
		ClassMask:  WarlockSpellShadowBolt,
	})

	warlock.MakeProcTriggerAura(core.ProcTrigger{
		Name:           "Nightfall",
		ClassSpellMask: WarlockSpellCorruption | WarlockSpellDrainLife,
		// Forever puts the real per-rank chance on the effect; ProcChanceAt reads a flat 100%.
		ProcChance: spellData.Nightfall.FractionAt(warlock.Talents.Nightfall),
		Callback:   core.CallbackOnPeriodicDamageDealt,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			warlock.NightfallProcAura.Activate(sim)
		},
	})
}

func (warlock *Warlock) applyEmpoweredCorruption() {
	if warlock.Talents.ImprovedCorruption == 0 {
		return
	}

	// TODO: this read Talents.EmpoweredCorruption, a talent Forever drops entirely, while
	// the registrar's guard is ImprovedCorruption (kept) -- looks like a long-standing
	// copy-paste bug. EmpoweredCorruption's spellData ladder is unrelated to
	// ImprovedCorruption's rank count, so the bonus coefficient is unknown and pinned to 0
	// until confirmed against the Forever tooltip.
	warlock.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DotBonusCoeffecient_Flat,
		FloatValue: 0,
		ClassMask:  WarlockSpellCorruption,
	})
}

func (warlock *Warlock) applyShadowMastery() {
	if warlock.Talents.ShadowMastery == 0 {
		return
	}

	warlock.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: spellData.ShadowMastery.Effect(shared.A_ADD_PCT_MODIFIER, shared.SPELLMOD_DAMAGE).FractionAt(warlock.Talents.ShadowMastery),
		ClassMask:  WarlockShadowDamage,
	})
}

/*
Demonology
Skipping the following:
  - Improved Healthstone
  - Improved Health Funnel
  - Improved Voidwalker
  - Fel Domination
  - Demonic Aegis -> implemented in armors.go
  - Mana Feed -> applied in lifetap.go
*/
func (warlock *Warlock) appyImprovedImp() {
	if warlock.Talents.ImprovedImp == 0 || warlock.Options.SacrificeSummon {
		return
	}

	warlock.Imp.AddStaticMod(core.SpellModConfig{
		Kind: core.SpellMod_DamageDone_Flat,
		// SPELLMOD_ALL_EFFECTS carries the same ladder and also covers Blood Pact, which is
		// buffed elsewhere; this mod is the Firebolt damage half.
		FloatValue: spellData.ImprovedImp.Effect(shared.A_ADD_PCT_MODIFIER, shared.SPELLMOD_DAMAGE).FractionAt(warlock.Talents.ImprovedImp),
		ClassMask:  WarlockSpellImpFireBolt,
	})
}

func (warlock *Warlock) applyDemonicEmbrace() {
	if warlock.Talents.DemonicEmbrace == 0 {
		return
	}

	warlock.MultiplyStat(stats.Stamina, 1.0+(0.03)*float64(warlock.Talents.DemonicEmbrace))
	warlock.MultiplyStat(stats.Spirit, 1.0-(0.01)*float64(warlock.Talents.DemonicEmbrace))
}

func (warlock *Warlock) applyImprovedSayaad() {
	if warlock.Talents.ImprovedSayaad == 0 || warlock.Options.SacrificeSummon {
		return
	}

	//This might not actually increase the damage, find a source to prove this
	warlock.Succubus.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: 0.1 * float64(warlock.Talents.ImprovedSayaad),
		ClassMask:  WarlockSpellSuccubusLashOfPain,
	})
}

func (warlock *Warlock) applyUnholyPower() {
	if warlock.Talents.UnholyPower == 0 || warlock.Options.SacrificeSummon {
		return
	}

	for _, pet := range warlock.Pets {
		if pet != &warlock.Imp.Pet {
			pet.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= spellData.UnholyPower.MultiplierAt(warlock.Talents.UnholyPower)
		}
	}

	warlock.Imp.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: spellData.UnholyPower.FractionAt(warlock.Talents.UnholyPower),
		ClassMask:  WarlockSpellImpFireBolt,
	})
}

func (warlock *Warlock) applyDemonicSacrifice() {
	if !warlock.Talents.DemonicSacrifice || warlock.Options.SacrificeSummon == false {
		return
	}

	switch warlock.Options.Summon {
	case proto.WarlockOptions_Succubus:
		warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= 1.15
	case proto.WarlockOptions_Imp:
		warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] *= 1.15
	case proto.WarlockOptions_Felguard:
		warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= 1.10
		warlock.applyDemonicSacrificeManaRegen(core.ActionID{SpellID: 18788}, 0.02)
	case proto.WarlockOptions_Felhunter:
		warlock.applyDemonicSacrificeManaRegen(core.ActionID{SpellID: 18792}, 0.03)
	}
}

// Demonic Sacrifice restores a percentage of maximum mana every 4 seconds, which is
// independent of the sim's regular 2 second mana ticks.
func (warlock *Warlock) applyDemonicSacrificeManaRegen(actionID core.ActionID, manaPercent float64) {
	manaMetrics := warlock.NewManaMetrics(actionID)

	warlock.RegisterResetEffect(func(sim *core.Simulation) {
		core.StartPeriodicAction(sim, core.PeriodicActionOptions{
			Period:   time.Second * 4,
			Priority: core.ActionPriorityRegen,
			OnAction: func(sim *core.Simulation) {
				warlock.AddMana(sim, warlock.MaxMana()*manaPercent, manaMetrics)
			},
		})
	})
}

func (warlock *Warlock) applyMasterDemonologist() {
	if warlock.Talents.MasterDemonologist == 0 || warlock.Options.SacrificeSummon == true {
		return
	}
	points := float64(warlock.Talents.MasterDemonologist)

	switch warlock.Options.Summon {

	case proto.WarlockOptions_Imp:
		warlock.MasterDemonologistAura = warlock.NewTemporaryStatsAura("Master Demonologist", core.ActionID{SpellID: (23825 + int32(points))}, stats.Stats{}, core.NeverExpires).Aura
		warlock.MasterDemonologistAura.AttachMultiplicativePseudoStatBuff(&warlock.PseudoStats.ThreatMultiplier, 1.0-0.04*points)
		for _, pet := range warlock.Pets {
			if pet == &warlock.Imp.Pet {
				pet.PseudoStats.ThreatMultiplier *= 1.0 - 0.04*points
			}
		}
	case proto.WarlockOptions_Succubus:
		warlock.MasterDemonologistAura = warlock.NewTemporaryStatsAura("Master Demonologist", core.ActionID{SpellID: (23832 + int32(points))}, stats.Stats{}, core.NeverExpires).Aura
		warlock.MasterDemonologistAura.AttachMultiplicativePseudoStatBuff(&warlock.PseudoStats.DamageDealtMultiplier, 1.0+0.02*points)
		for _, pet := range warlock.Pets {
			if pet == &warlock.Succubus.Pet {
				pet.PseudoStats.DamageDealtMultiplier *= 1.0 + 0.02*points
			}
		}
	case proto.WarlockOptions_Felguard:
		resistsBonus := 0.10 * points * 70
		warlock.MasterDemonologistAura = warlock.NewTemporaryStatsAura("Master Demonologist", core.ActionID{SpellID: (35701 + int32(points))}, stats.Stats{
			stats.ArcaneResistance: resistsBonus,
			stats.FireResistance:   resistsBonus,
			stats.FrostResistance:  resistsBonus,
			stats.NatureResistance: resistsBonus,
			stats.ShadowResistance: resistsBonus,
		}, core.NeverExpires).Aura
		warlock.MasterDemonologistAura.AttachMultiplicativePseudoStatBuff(&warlock.PseudoStats.DamageDealtMultiplier, 1.0+0.01*points)

		for _, pet := range warlock.Pets {
			if pet == &warlock.Felguard.Pet {
				pet.PseudoStats.DamageDealtMultiplier *= 1.0 + 0.01*points
			}
		}
	case proto.WarlockOptions_Voidwalker:
		warlock.PseudoStats.BonusPhysicalDamageTaken *= 1.0 - 0.02*points
		warlock.MasterDemonologistAura = warlock.NewTemporaryStatsAura("Master Demonologist", core.ActionID{SpellID: (23840 + int32(points))}, stats.Stats{}, core.NeverExpires).Aura
		warlock.MasterDemonologistAura.AttachMultiplicativePseudoStatBuff(&warlock.PseudoStats.BonusPhysicalDamageTaken, 1.0-0.02*points)
		for _, pet := range warlock.Pets {
			if pet == &warlock.Voidwalker.Pet {
				pet.PseudoStats.BonusPhysicalDamageTaken *= 1.0 - 0.02*points
			}
		}
	case proto.WarlockOptions_Felhunter:
		resistsBonus := 0.20 * points * 70
		warlock.MasterDemonologistAura = warlock.NewTemporaryStatsAura("Master Demonologist", core.ActionID{SpellID: (23836 + int32(points))}, stats.Stats{}, core.NeverExpires).Aura
		warlock.MasterDemonologistAura.AttachStatsBuff(stats.Stats{
			stats.ArcaneResistance: resistsBonus,
			stats.FireResistance:   resistsBonus,
			stats.FrostResistance:  resistsBonus,
			stats.NatureResistance: resistsBonus,
			stats.ShadowResistance: resistsBonus,
		})
		for _, pet := range warlock.Pets {
			if pet == &warlock.Felhunter.Pet {
				pet.NewTemporaryStatsAura("Master Demonologist", core.ActionID{SpellID: (23836 + int32(points))}, stats.Stats{
					stats.ArcaneResistance: resistsBonus,
					stats.FireResistance:   resistsBonus,
					stats.FrostResistance:  resistsBonus,
					stats.NatureResistance: resistsBonus,
					stats.ShadowResistance: resistsBonus,
				}, core.NeverExpires)
			}
		}
	}

}

func (warlock *Warlock) applySoulLink() {
	if !warlock.Talents.SoulLink {
		return
	}

	// TODO Add if/while pet is alive
	warlock.PseudoStats.DamageTakenMultiplier *= 0.80
	warlock.PseudoStats.DamageDealtMultiplier *= 1.05

	for _, pet := range warlock.Pets {
		pet.PseudoStats.DamageDealtMultiplier *= 1.05
	}
}

func (warlock *Warlock) applyDemonicKnowledge() {
	if warlock.Talents.DemonicKnowledge == 0 {
		return
	}

	warlock.DemonicKnowledgeAura = warlock.RegisterAura(core.Aura{
		Label:    "Demonic Knowledge",
		Duration: core.NeverExpires,
	})
}

func (warlock *Warlock) updateDemonicKnowledge(sim *core.Simulation) {
	if warlock.DemonicKnowledgeBonus != 0 {
		warlock.AddStatDynamic(sim, stats.SpellDamage, -warlock.DemonicKnowledgeBonus)
	}

	if warlock.ActivePet == nil {
		warlock.DemonicKnowledgeBonus = 0
		return
	}

	coeff := spellData.DemonicKnowledge.FractionAt(warlock.Talents.DemonicKnowledge)
	bonus := coeff * (warlock.ActivePet.GetStat(stats.Stamina) + warlock.ActivePet.GetStat(stats.Intellect))

	warlock.DemonicKnowledgeBonus = bonus
	warlock.AddStatDynamic(sim, stats.SpellDamage, bonus)
}

/*
Destruction
Skipped Talents:
  - Aftermath
*/
func (warlock *Warlock) applyImprovedShadowBolt() {
	if warlock.Talents.ImprovedShadowBolt == 0 {
		return
	}
	warlock.ImpShadowboltAura = core.ImprovedShadowBoltAura(warlock.CurrentTarget, 0, warlock.Talents.ImprovedShadowBolt)
}

func (warlock *Warlock) applyCataclysm() {
	if warlock.Talents.Cataclysm == 0 {
		return
	}

	warlock.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_PowerCost_Pct_Add,
		FloatValue: -0.01 * float64(warlock.Talents.Cataclysm),
		ClassMask:  WarlockDestructionSpells,
	})
}

func (warlock *Warlock) applyBane() {
	if warlock.Talents.Bane == 0 {
		return
	}

	warlock.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_CastTime_Flat,
		TimeValue: time.Millisecond * time.Duration(-100*warlock.Talents.Bane),
		ClassMask: WarlockSpellShadowBolt | WarlockSpellImmolate,
	})

	warlock.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_CastTime_Flat,
		TimeValue: time.Millisecond * time.Duration(-400*warlock.Talents.Bane),
		ClassMask: WarlockSpellSoulFire,
	})
}

func (warlock *Warlock) applyShadowburn() {
	if !warlock.Talents.Shadowburn {
		return
	}

	warlock.registerShadowBurn()
}

func (warlock *Warlock) applyDestructiveReach() {
	if warlock.Talents.DestructiveReach == 0 {
		return
	}

	warlock.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_ThreatMultiplier_Pct,
		FloatValue: -0.05 * float64(warlock.Talents.DestructiveReach),
		ClassMask:  WarlockDestructionSpells,
	})

}

func (warlock *Warlock) applyRuin() {
	if warlock.Talents.Ruin == 0 {
		return
	}

	// TODO: Forever expands Ruin from 1 rank to 5; the per-rank crit-multiplier bonus is
	// unconfirmed, so it is pinned to 0 until the Forever tooltip is known.
	warlock.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_CritMultiplier_Flat,
		FloatValue: 0,
		ClassMask:  WarlockDestructionSpells,
	})
}

func (warlock *Warlock) applyConflagrate() {
	if !warlock.Talents.Conflagrate {
		return
	}

	warlock.registerConflagrate()
}

func (warlock *Warlock) applyShadowAndFlame() {
	if warlock.Talents.ShadowAndFlame == 0 {
		return
	}

	warlock.AddStaticMod(core.SpellModConfig{
		Kind: core.SpellMod_BonusCoeffecient_Flat,
		// Four dummy effects, of which only the first (4% per rank) matches the coefficient
		// bonus this talent has always granted; the other three (20/2/2 per rank) are
		// unidentified.
		FloatValue: spellData.ShadowAndFlame.EffectAt(0).FractionAt(warlock.Talents.ShadowAndFlame),
		ClassMask:  WarlockSpellShadowBolt | WarlockSpellIncinerate,
	})
}

// applySoulHarvesting implements Soul Harvesting, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (warlock *Warlock) applySoulHarvesting() {
	if warlock.Talents.SoulHarvesting == 0 {
		return
	}

	panic("To be implemented")
}

// applyImprovedDrains implements Improved Drains, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (warlock *Warlock) applyImprovedDrains() {
	if warlock.Talents.ImprovedDrains == 0 {
		return
	}

	panic("To be implemented")
}

// applyImprovedBaneOfAgony implements Improved Bane of Agony, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (warlock *Warlock) applyImprovedBaneOfAgony() {
	if warlock.Talents.ImprovedBaneOfAgony == 0 {
		return
	}

	panic("To be implemented")
}

// applyFelConcentration implements Fel Concentration, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (warlock *Warlock) applyFelConcentration() {
	if warlock.Talents.FelConcentration == 0 {
		return
	}

	panic("To be implemented")
}

// applyPandemic implements Pandemic, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (warlock *Warlock) applyPandemic() {
	if warlock.Talents.Pandemic == 0 {
		return
	}

	panic("To be implemented")
}

// applyMalevolence implements Malevolence, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (warlock *Warlock) applyMalevolence() {
	if warlock.Talents.Malevolence == 0 {
		return
	}

	panic("To be implemented")
}

// applyCurseOfExhaustion implements Curse of Exhaustion, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (warlock *Warlock) applyCurseOfExhaustion() {
	if !warlock.Talents.CurseOfExhaustion {
		return
	}

	panic("To be implemented")
}

// applySiphonLife implements Siphon Life, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (warlock *Warlock) applySiphonLife() {
	if !warlock.Talents.SiphonLife {
		return
	}

	panic("To be implemented")
}

// applyWrack implements Wrack, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (warlock *Warlock) applyWrack() {
	if !warlock.Talents.Wrack {
		return
	}

	panic("To be implemented")
}

// applyImprovedHealthFunnel implements Improved Health Funnel, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (warlock *Warlock) applyImprovedHealthFunnel() {
	if warlock.Talents.ImprovedHealthFunnel == 0 {
		return
	}

	panic("To be implemented")
}

// applyImprovedVoidwalker implements Improved Voidwalker, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (warlock *Warlock) applyImprovedVoidwalker() {
	if warlock.Talents.ImprovedVoidwalker == 0 {
		return
	}

	panic("To be implemented")
}

// applyFelVitality implements Fel Vitality, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (warlock *Warlock) applyFelVitality() {
	if warlock.Talents.FelVitality == 0 {
		return
	}

	panic("To be implemented")
}

// applyDemonicEnergies implements Demonic Energies, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (warlock *Warlock) applyDemonicEnergies() {
	if warlock.Talents.DemonicEnergies == 0 {
		return
	}

	panic("To be implemented")
}

// applyMasterSummoner implements Master Summoner, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (warlock *Warlock) applyMasterSummoner() {
	if warlock.Talents.MasterSummoner == 0 {
		return
	}

	panic("To be implemented")
}

// applyDecimation implements Decimation, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (warlock *Warlock) applyDecimation() {
	if warlock.Talents.Decimation == 0 {
		return
	}

	panic("To be implemented")
}

// applyFelDomination implements Fel Domination, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (warlock *Warlock) applyFelDomination() {
	if !warlock.Talents.FelDomination {
		return
	}

	panic("To be implemented")
}

// applyDemonicBrand implements Demonic Brand, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (warlock *Warlock) applyDemonicBrand() {
	if warlock.Talents.DemonicBrand == 0 {
		return
	}

	panic("To be implemented")
}

// applyImprovedFelhunter implements Improved Felhunter, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (warlock *Warlock) applyImprovedFelhunter() {
	if warlock.Talents.ImprovedFelhunter == 0 {
		return
	}

	panic("To be implemented")
}

// applyDemonicPact implements Demonic Pact, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (warlock *Warlock) applyDemonicPact() {
	if !warlock.Talents.DemonicPact {
		return
	}

	panic("To be implemented")
}

// applyMoltenSkin implements Molten Skin, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (warlock *Warlock) applyMoltenSkin() {
	if warlock.Talents.MoltenSkin == 0 {
		return
	}

	panic("To be implemented")
}

// applyAftermath implements Aftermath, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (warlock *Warlock) applyAftermath() {
	if warlock.Talents.Aftermath == 0 {
		return
	}

	panic("To be implemented")
}

// applyIntensity implements Intensity, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (warlock *Warlock) applyIntensity() {
	if warlock.Talents.Intensity == 0 {
		return
	}

	panic("To be implemented")
}

// applyAgonizingFlames implements Agonizing Flames, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (warlock *Warlock) applyAgonizingFlames() {
	if warlock.Talents.AgonizingFlames == 0 {
		return
	}

	panic("To be implemented")
}

// applyPyroclasm implements Pyroclasm, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (warlock *Warlock) applyPyroclasm() {
	if warlock.Talents.Pyroclasm == 0 {
		return
	}

	panic("To be implemented")
}

// applyBaneOfHavoc implements Bane of Havoc, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (warlock *Warlock) applyBaneOfHavoc() {
	if !warlock.Talents.BaneOfHavoc {
		return
	}

	panic("To be implemented")
}

// applyFireAndBrimstone implements Fire and Brimstone, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (warlock *Warlock) applyFireAndBrimstone() {
	if warlock.Talents.FireAndBrimstone == 0 {
		return
	}

	panic("To be implemented")
}

// applyIncinerate implements Incinerate, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (warlock *Warlock) applyIncinerate() {
	if !warlock.Talents.Incinerate {
		return
	}

	panic("To be implemented")
}
