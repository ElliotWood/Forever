package hunter

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/stats"
)

func (hunter *Hunter) registerMarksmanshipTalents() {
	// Tier 1
	hunter.registerHawkEye()
	hunter.registerImprovedConcussiveShot()
	hunter.registerLethalAttacks()

	// Tier 2
	hunter.registerImprovedStings()
	hunter.registerEfficiency()
	hunter.registerCarefulAim()

	// Tier 3
	hunter.registerRapidKilling()
	hunter.registerImprovedArcaneShot()
	hunter.registerLoneWolf()

	// Tier 4
	// Trueshot Aura handled as a group buff in hunter.go
	hunter.registerMortalShots()
	hunter.registerImprovedSerpentSting()

	// Tier 5
	hunter.registerRapidRecuperation()
	hunter.registerBarrage()
	hunter.registerScatterShot()

	// Tier 6
	hunter.registerRangedWeaponSpecialization()

	// Tier 7
	hunter.registerSniperShot()
}

func (hunter *Hunter) registerEfficiency() {
	if hunter.Talents.Efficiency == 0 {
		return
	}

	hunter.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_PowerCost_Pct_Add,
		ClassMask:  HunterSpellsShotsAndStings,
		FloatValue: -0.02 * float64(hunter.Talents.Efficiency),
	})
}

func (hunter *Hunter) registerImprovedArcaneShot() {
	if hunter.Talents.ImprovedArcaneShot == 0 {
		return
	}

	hunter.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Flat,
		ClassMask: HunterSpellArcaneShot,
		TimeValue: -core.DurationFromSeconds(0.2 * float64(hunter.Talents.ImprovedArcaneShot)),
	})
}

func (hunter *Hunter) registerRapidKilling() {
	if hunter.Talents.RapidKilling == 0 {
		return
	}

	hunter.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Flat,
		ClassMask: HunterSpellRapidFire,
		TimeValue: -core.DurationFromSeconds(60 * float64(hunter.Talents.RapidKilling)),
	})
}

func (hunter *Hunter) registerImprovedStings() {
	if hunter.Talents.ImprovedStings == 0 {
		return
	}

	hunter.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		ClassMask:  HunterSpellSerpentSting,
		FloatValue: spellData.ImprovedStings.Effect(shared.A_ADD_PCT_MODIFIER, shared.SPELLMOD_DOT).FractionAt(hunter.Talents.ImprovedStings),
	})
}

func (hunter *Hunter) registerMortalShots() {
	if hunter.Talents.MortalShots == 0 {
		return
	}

	hunter.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_CritMultiplier_Flat,
		ProcMask:   core.ProcMaskRanged,
		FloatValue: spellData.MortalShots.FractionAt(hunter.Talents.MortalShots),
	})
}

func (hunter *Hunter) registerBarrage() {
	if hunter.Talents.Barrage == 0 {
		return
	}

	hunter.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		ClassMask:  HunterSpellMultiShot | HunterSpellVolley,
		FloatValue: spellData.Barrage.Effect(shared.A_ADD_PCT_MODIFIER, shared.SPELLMOD_DAMAGE).FractionAt(hunter.Talents.Barrage),
	})
}

func (hunter *Hunter) registerRangedWeaponSpecialization() {
	if hunter.Talents.RangedWeaponSpecialization == 0 {
		return
	}

	hunter.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ProcMask:   core.ProcMaskRanged,
		FloatValue: spellData.RangedWeaponSpecialization.FractionAt(hunter.Talents.RangedWeaponSpecialization),
	})
}

func (hunter *Hunter) registerCarefulAim() {
	if hunter.Talents.CarefulAim == 0 {
		return
	}

	hunter.AddStatDependency(stats.Intellect, stats.RangedAttackPower, 0.15*float64(hunter.Talents.CarefulAim))
}

func (hunter *Hunter) registerHawkEye() {
	if hunter.Talents.HawkEye == 0 {
		return
	}

	bonusRange := float64(hunter.Talents.HawkEye) * 2
	ranged := hunter.AutoAttacks.Ranged()

	if ranged != nil {
		ranged.MaxRange += bonusRange
	}

	hunter.AddStaticMod(core.SpellModConfig{
		Kind:     core.SpellMod_Custom,
		ProcMask: core.ProcMaskRanged,
		ApplyCustom: func(mod *core.SpellMod, spell *core.Spell) {
			if spell.MaxRange > 0 {
				spell.MaxRange += bonusRange
			}
		},
		RemoveCustom: func(mod *core.SpellMod, spell *core.Spell) {
			if spell.MaxRange > 0 {
				spell.MaxRange -= bonusRange
			}
		},
	})
}

// registerImprovedConcussiveShot implements Improved Concussive Shot, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerImprovedConcussiveShot() {
	if hunter.Talents.ImprovedConcussiveShot == 0 {
		return
	}

	panic("To be implemented")
}

// registerLethalAttacks implements Lethal Attacks, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerLethalAttacks() {
	if hunter.Talents.LethalAttacks == 0 {
		return
	}

	panic("To be implemented")
}

// registerLoneWolf implements Lone Wolf, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerLoneWolf() {
	if !hunter.Talents.LoneWolf {
		return
	}

	panic("To be implemented")
}

// registerImprovedSerpentSting implements Improved Serpent Sting, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerImprovedSerpentSting() {
	if hunter.Talents.ImprovedSerpentSting == 0 {
		return
	}

	panic("To be implemented")
}

// registerRapidRecuperation implements Rapid Recuperation, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerRapidRecuperation() {
	if hunter.Talents.RapidRecuperation == 0 {
		return
	}

	panic("To be implemented")
}

// registerScatterShot implements Scatter Shot, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerScatterShot() {
	if !hunter.Talents.ScatterShot {
		return
	}

	panic("To be implemented")
}

// registerSniperShot implements Sniper Shot, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (hunter *Hunter) registerSniperShot() {
	if !hunter.Talents.SniperShot {
		return
	}

	panic("To be implemented")
}
