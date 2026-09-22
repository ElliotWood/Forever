package hunter

import (
	"time"

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
	// Rapid Killing: rapid_fire.go
	hunter.registerImprovedArcaneShot()
	hunter.registerLoneWolf()

	// Tier 4
	// Trueshot Aura handled as a party buff in hunter.go
	hunter.registerMortalShots()

	// Tier 5
	hunter.registerRapidRecuperation()
	hunter.registerBarrage()
	hunter.registerScatterShot()

	// Tier 6
	hunter.registerRangedWeaponSpecialization()

	// Tier 7
	// Sniper Shot: sniper_shot.go
}

func (hunter *Hunter) registerEfficiency() {
	if hunter.Talents.Efficiency == 0 {
		return
	}

	// The tooltip reads "Shots, Stings and melee abilities".
	hunter.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_PowerCost_Pct_Add,
		ClassMask:  HunterSpellsShotsAndStings | HunterSpellsMelee,
		FloatValue: spellData.Efficiency.Effect(shared.A_ADD_PCT_MODIFIER, shared.SPELLMOD_COST).FractionAt(hunter.Talents.Efficiency),
	})
}

func (hunter *Hunter) registerImprovedArcaneShot() {
	if hunter.Talents.ImprovedArcaneShot == 0 {
		return
	}

	hunter.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Flat,
		ClassMask: HunterSpellArcaneShot,
		TimeValue: time.Millisecond * time.Duration(spellData.ImprovedArcaneShot.
			Effect(shared.A_ADD_FLAT_MODIFIER, shared.SPELLMOD_COOLDOWN).
			ValueAt(hunter.Talents.ImprovedArcaneShot)),
	})
}

func (hunter *Hunter) registerImprovedStings() {
	if hunter.Talents.ImprovedStings == 0 {
		return
	}

	// The beta showed rank 1 at 6%; the client's rank curve gives 6/13/20, not the linear 6/12/18
	// that was assumed for the ranks nobody saw.
	hunter.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DotDamageDone_Pct,
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
		ClassMask:  HunterSpellMultiShot | HunterSpellAimedShot | HunterSpellVolley,
		FloatValue: spellData.Barrage.EffectAt(0).FractionAt(hunter.Talents.Barrage),
	})
}

func (hunter *Hunter) registerRangedWeaponSpecialization() {
	if hunter.Talents.RangedWeaponSpecialization == 0 {
		return
	}

	// Serpent Sting is a sting, not a ranged weapon attack, and is left out on master. It carries
	// the ranged special proc mask, so the shots are named instead: Auto Shot has no class mask.
	value := spellData.RangedWeaponSpecialization.FractionAt(hunter.Talents.RangedWeaponSpecialization)
	hunter.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ProcMask:   core.ProcMaskRangedAuto,
		FloatValue: value,
	})
	hunter.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  HunterSpellsAll &^ HunterSpellSerpentSting,
		ProcMask:   core.ProcMaskRangedSpecial,
		FloatValue: value,
	})
}

// Generator gap: spell 1223984 carries no effect rows at all, so the 0.2 attack power per intellect
// a rank stays from our client-verified sim. The tooltip only says Attack Power, but melee and
// ranged attack power are separate stats here and every other attack power buff feeds both.
func (hunter *Hunter) registerCarefulAim() {
	if hunter.Talents.CarefulAim == 0 {
		return
	}

	apPerInt := 0.2 * float64(hunter.Talents.CarefulAim)
	hunter.AddStatDependency(stats.Intellect, stats.AttackPower, apPerInt)
	hunter.AddStatDependency(stats.Intellect, stats.RangedAttackPower, apPerInt)
}

func (hunter *Hunter) registerHawkEye() {
	if hunter.Talents.HawkEye == 0 {
		return
	}

	bonusRange := spellData.HawkEye.Effect(shared.A_ADD_FLAT_MODIFIER, shared.SPELLMOD_RANGE).ValueAt(hunter.Talents.HawkEye)

	if ranged := hunter.AutoAttacks.Ranged(); ranged != nil {
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

func (hunter *Hunter) registerLethalAttacks() {
	if hunter.Talents.LethalAttacks == 0 {
		return
	}

	crit := spellData.LethalAttacks.ValueAt(hunter.Talents.LethalAttacks)
	hunter.AddStat(stats.PhysicalCritPercent, crit)
	hunter.AddStat(stats.SpellCritPercent, crit)
}

// Generator gap: Lone Wolf (415370) has no generated row, so the 20% from our client-verified sim
// stands.
func (hunter *Hunter) registerLoneWolf() {
	if !hunter.Talents.LoneWolf || hunter.Pet != nil {
		return
	}

	hunter.PseudoStats.DamageDealtMultiplier *= 1.2
}

// Only the Serpent Sting half is modelled: nothing dies mid fight to hand out the kill half.
func (hunter *Hunter) registerRapidRecuperation() {
	if hunter.Talents.RapidRecuperation == 0 {
		return
	}

	buff := spellData.RapidRecuperationTriggered.HighestRank()
	regen := spellData.RapidRecuperation.EffectAt(0).FractionAt(hunter.Talents.RapidRecuperation)

	procAura := hunter.RegisterAura(core.Aura{
		Label:    "Rapid Recuperation",
		ActionID: core.ActionID{SpellID: buff.SpellID},
		Duration: buff.Duration,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			hunter.PseudoStats.SpiritRegenRateCasting += regen
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			hunter.PseudoStats.SpiritRegenRateCasting -= regen
		},
	})

	hunter.MakeProcTriggerAura(core.ProcTrigger{
		Name:           "Rapid Recuperation Trigger",
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: HunterSpellSerpentSting,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() {
				procAura.Activate(sim)
			}
		},
	})
}

// registerImprovedConcussiveShot implements Improved Concussive Shot, new in Forever.
//
// TODO: a stun chance on Concussive Shot; bosses are immune and the sim does not cast it.
func (hunter *Hunter) registerImprovedConcussiveShot() {
	if hunter.Talents.ImprovedConcussiveShot == 0 {
		return
	}
}

// registerScatterShot implements Scatter Shot, new in Forever.
//
// TODO: a 4 sec disorient; bosses are immune.
func (hunter *Hunter) registerScatterShot() {
	if !hunter.Talents.ScatterShot {
		return
	}
}
