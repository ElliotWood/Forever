package hunter

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/stats"
)

func (hunter *Hunter) registerSurvivalTalents() {
	// Tier 1
	hunter.registerImprovedTracking()
	hunter.registerDeflection()

	// Tier 2
	hunter.registerEntrapment()
	hunter.registerSavageStrikes()
	hunter.registerSurvivalist()
	hunter.registerImprovedWingClip()

	// Tier 3
	hunter.registerCleverTraps()
	hunter.registerSurefooted()
	hunter.registerDeterrence()

	// Tier 4
	hunter.registerSurvivalTactics()
	hunter.registerPredatorsEdge()
	hunter.registerCounterattack()

	// Tier 5
	hunter.registerResourcefulness()
	hunter.registerExposePrey()
	hunter.registerSurvivalistsDiscipline()
	hunter.registerStriderKick()

	// Tier 6
	hunter.registerLightningReflexes()

	// Tier 7
	hunter.registerLaceratingStrikes()
}

func (hunter *Hunter) registerSavageStrikes() {
	if hunter.Talents.SavageStrikes == 0 {
		return
	}

	hunter.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_BonusCrit_Percent,
		ClassMask:  HunterSpellRaptorStrike,
		FloatValue: spellData.SavageStrikes.ValueAt(hunter.Talents.SavageStrikes),
	})
}

func (hunter *Hunter) registerSurvivalist() {
	if hunter.Talents.Survivalist == 0 {
		return
	}

	hunter.MultiplyStat(stats.Health, spellData.Survivalist.MultiplierAt(hunter.Talents.Survivalist))
}

func (hunter *Hunter) registerSurefooted() {
	if hunter.Talents.Surefooted == 0 {
		return
	}

	hunter.AddStat(stats.PhysicalHitPercent, float64(hunter.Talents.Surefooted))
}

func (hunter *Hunter) registerResourcefulness() {
	if hunter.Talents.Resourcefulness == 0 {
		return
	}

	hunter.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_PowerCost_Pct_Add,
		ClassMask:  HunterSpellRaptorStrike,
		FloatValue: -0.2 * float64(hunter.Talents.Resourcefulness),
	})
}

func (hunter *Hunter) registerLightningReflexes() {
	if hunter.Talents.LightningReflexes == 0 {
		return
	}

	hunter.MultiplyStat(stats.Agility, spellData.LightningReflexes.MultiplierAt(hunter.Talents.LightningReflexes))
}

// registerImprovedTracking implements Improved Tracking, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerImprovedTracking() {
	if hunter.Talents.ImprovedTracking == 0 {
		return
	}

	panic("To be implemented")
}

// registerDeflection implements Deflection, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerDeflection() {
	if hunter.Talents.Deflection == 0 {
		return
	}

	panic("To be implemented")
}

// registerEntrapment implements Entrapment, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerEntrapment() {
	if hunter.Talents.Entrapment == 0 {
		return
	}

	panic("To be implemented")
}

// registerImprovedWingClip implements Improved Wing Clip, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerImprovedWingClip() {
	if hunter.Talents.ImprovedWingClip == 0 {
		return
	}

	panic("To be implemented")
}

// registerCleverTraps implements Clever Traps, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerCleverTraps() {
	if hunter.Talents.CleverTraps == 0 {
		return
	}

	panic("To be implemented")
}

// registerDeterrence implements Deterrence, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerDeterrence() {
	if !hunter.Talents.Deterrence {
		return
	}

	panic("To be implemented")
}

// registerSurvivalTactics implements Survival Tactics, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerSurvivalTactics() {
	if hunter.Talents.SurvivalTactics == 0 {
		return
	}

	panic("To be implemented")
}

// registerPredatorsEdge implements Predator's Edge, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerPredatorsEdge() {
	if hunter.Talents.PredatorsEdge == 0 {
		return
	}

	panic("To be implemented")
}

// registerCounterattack implements Counterattack, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerCounterattack() {
	if !hunter.Talents.Counterattack {
		return
	}

	panic("To be implemented")
}

// registerExposePrey implements Expose Prey, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerExposePrey() {
	if hunter.Talents.ExposePrey == 0 {
		return
	}

	panic("To be implemented")
}

// registerSurvivalistsDiscipline implements Survivalist's Discipline, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerSurvivalistsDiscipline() {
	if hunter.Talents.SurvivalistsDiscipline == 0 {
		return
	}

	panic("To be implemented")
}

// registerStriderKick implements Strider Kick, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerStriderKick() {
	if !hunter.Talents.StriderKick {
		return
	}

	panic("To be implemented")
}

// registerLaceratingStrikes implements Lacerating Strikes, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerLaceratingStrikes() {
	if !hunter.Talents.LaceratingStrikes {
		return
	}

	panic("To be implemented")
}
