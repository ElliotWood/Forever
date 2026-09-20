package warlock

import (
	"time"

	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

func (warlock *Warlock) registerDemonologyTalents() {
	// Tier 1
	warlock.applyImprovedHealthFunnel()
	warlock.appyImprovedImp()
	warlock.applyDemonicEmbrace()
	warlock.applyUnholyPower()

	// Tier 2
	// Demonic Aegis implemented in armors.go
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
