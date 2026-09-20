package mage

import (
	"math"
	"time"

	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/stats"
)

func (mage *Mage) ApplyTalents() {

	//------- ARCANE --------
	mage.registerArcaneSubtlety()
	mage.registerArcaneFocus()
	//mage.registerImprovedArcaneMissiles()

	mage.registerWandSpecialization()
	mage.registerMagicAbsorption()
	mage.registerArcaneConcentration()

	// mage.registerMagicAttunement()
	mage.registerArcaneImpact()
	//mage.registerArcaneFortitude()

	// mage.registerImprovedManaShield()
	mage.registerImprovedCounterspell()
	mage.registerArcaneMeditation()

	//mage.registerImprovedBlink()
	mage.registerArcaneMind()

	// mage.registerPrismaticCloak()
	mage.registerArcaneInstability()

	//-------  FIRE  --------
	mage.registerImprovedFireball()
	mage.registerImpact()

	mage.registerIgnite()
	mage.registerFlameThrowing()

	mage.registerIncineration()
	mage.registerImprovedFlamestrike()
	mage.registerBurningSoul()

	// mage.registerMoltenShields()
	mage.registerMasterOfElements()

	mage.registerCriticalMass()

	// mage.registerBlazingSpeed()
	mage.registerFirePower()

	//------- FROST --------
	mage.registerFrostWarding()
	mage.registerImprovedFrostbolt()
	mage.registerElementalPrecision()

	mage.registerIceShards()
	mage.registerFrostbite()
	mage.registerImprovedFrostNova()
	mage.registerPermafrost()

	mage.registerPiercingIce()
	mage.registerImprovedBlizzard()

	mage.registerArcticReach()
	mage.registerFrostChanneling()
	mage.registerShatter()

	mage.registerImprovedConeOfCold()

	mage.registerWinterChill()

	mage.registerIceBarrier()

	// Forever additions, not yet implemented.
	mage.registerArcaneBlastTalent()
	mage.registerArcaneGeometry()
	mage.registerArcaneResilience()
	mage.registerArcaneShielding()
	mage.registerImprovedChanneling()
	mage.registerMissileBarrage()
	mage.registerPyroblastTalent()
	mage.registerHotStreak()
	mage.registerImprovedFireWard()
	mage.registerWakeOfFire()
	mage.registerFingersOfFrost()
	mage.registerIceBlock()
	mage.registerIceLance()
}

func (mage *Mage) registerArcaneSubtlety() {
	if mage.Talents.ArcaneSubtlety == 0 {
		return
	}

	//all spells resist 5 & arcance spells threat 20% per rank
	mage.AddStaticMod(core.SpellModConfig{
		School:     core.SpellSchoolArcane,
		FloatValue: -.20 * float64(mage.Talents.ArcaneSubtlety),
		Kind:       core.SpellMod_ThreatMultiplier_Pct,
	})
}

func (mage *Mage) registerArcaneFocus() {
	if mage.Talents.ArcaneFocus == 0 {
		return
	}

	mage.PseudoStats.SchoolBonusHitChance[stats.SchoolIndexArcane] += spellData.ArcaneFocus.ValueAt(mage.Talents.ArcaneFocus)
}

func (mage *Mage) registerArcaneConcentration() {
	if mage.Talents.ArcaneConcentration == 0 {
		return
	}

	// TODO: Forever drops Arcane Potency; the Clearcasting crit bonus it fed is pinned to
	// 0 until we know whether the effect moved onto another talent.
	bonusCrit := 0.0
	var proccedAt time.Duration
	var proccedSpell *core.Spell

	mage.ClearCasting = mage.RegisterAura(core.Aura{
		Label:    "Clearcasting",
		ActionID: core.ActionID{SpellID: 12536},
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			mage.AddStatDynamic(sim, stats.SpellCritRating, bonusCrit)
			aura.Unit.PseudoStats.SpellCostPercentModifier -= 100
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			mage.AddStatDynamic(sim, stats.SpellCritRating, -bonusCrit)
			aura.Unit.PseudoStats.SpellCostPercentModifier += 100
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.ClassSpellMask&MageSpellsAllDamaging == 0 {
				return
			}

			if spell.DefaultCast.Cost == 0 {
				return
			}

			if proccedAt == sim.CurrentTime && proccedSpell == spell {
				// Means this is another hit from the same cast that procced CC.
				return
			}

			aura.Deactivate(sim)
		},
	})

	mage.RegisterAura(core.Aura{
		Label:    "Arcane Concentration",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ClassSpellMask&MageSpellsAllDamaging == 0 {
				return
			}

			if !result.Landed() {
				return
			}

			// Forever states a flat SpellAuraOptions.ProcChance of 100 on the talent spell and puts the
			// real per-rank chance on the effect, so ProcChanceAt would read 100% at every rank.
			procChance := spellData.ArcaneConcentration.FractionAt(mage.Talents.ArcaneConcentration)
			if sim.Proc(procChance, "Arcane Concentration") {
				proccedAt = sim.CurrentTime
				proccedSpell = spell
				mage.ClearCasting.Activate(sim)
			}
		},
	})
}

func (mage *Mage) registerArcaneImpact() {
	if mage.Talents.ArcaneImpact == 0 {
		return
	}

	mage.AddStaticMod(core.SpellModConfig{
		ClassMask:  MageSpellArcaneBlast | MageSpellArcaneExplosion,
		FloatValue: spellData.ArcaneImpact.ValueAt(mage.Talents.ArcaneImpact),
		Kind:       core.SpellMod_BonusCrit_Percent,
	})
}

func (mage *Mage) registerArcaneMeditation() {
	if mage.Talents.ArcaneMeditation == 0 {
		return
	}

	mage.PseudoStats.SpiritRegenRateCasting += float64(mage.Talents.ArcaneMeditation) * 0.1
	mage.UpdateManaRegenRates()
}

func (mage *Mage) registerArcaneMind() {
	if mage.Talents.ArcaneMind == 0 {
		return
	}

	mage.MultiplyStat(stats.Intellect, 1+(float64(mage.Talents.ArcaneMind)*.03))
}

func (mage *Mage) registerArcaneInstability() {
	if mage.Talents.ArcaneInstability == 0 {
		return
	}

	// TODO: To be implemented. Forever's regenerated aura enum no longer carries the crit
	// chance aura this talent's rank data used (A_MOD_SPELL_CRIT_CHANCE_SCHOOL is gone from
	// the auto-generated table), so only the damage-done portion below is applied.

	mage.AddStaticMod(core.SpellModConfig{
		ClassMask:  MageSpellsAll,
		FloatValue: spellData.ArcaneInstability.Effect(shared.A_MOD_DAMAGE_PERCENT_DONE, 126).FractionAt(mage.Talents.ArcaneInstability),
		Kind:       core.SpellMod_DamageDone_Pct,
	})

}

// ------ FIRE TALENTS ------

func (mage *Mage) registerImprovedFireball() {
	if mage.Talents.ImprovedFireball == 0 {
		return
	}

	mage.AddStaticMod(core.SpellModConfig{
		ClassMask: MageSpellFireball,
		TimeValue: time.Millisecond * time.Duration(spellData.ImprovedFireball.Effect(shared.A_ADD_FLAT_MODIFIER, shared.SPELLMOD_CASTING_TIME).ValueAt(mage.Talents.ImprovedFireball)),
		Kind:      core.SpellMod_CastTime_Flat,
	})
}

func (mage *Mage) registerIgnite() {
	if mage.Talents.Ignite == 0 {
		return
	}
	igniteDamageMultiplier := float64(mage.Talents.Ignite) * .08

	igniteSpell := mage.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 12846},
		SpellSchool:      core.SpellSchoolFire,
		ProcMask:         core.ProcMaskSpellDamage,
		ClassSpellMask:   MageSpellIgnite,
		Flags:            core.SpellFlagIgnoreModifiers | core.SpellFlagNoSpellMods | core.SpellFlagNoOnCastComplete | core.SpellFlagIgnoreResists | core.SpellFlagProc,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:     "Ignite",
				Tag:       "IgniteDot",
				MaxStacks: math.MaxInt32,
			},
			NumberOfTicks: 2,
			TickLength:    2 * time.Second,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.Spell.CalcAndDealPeriodicDamage(sim, target, dot.SnapshotBaseDamage, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).Apply(sim)
		},
	})

	refreshIgnite := func(sim *core.Simulation, target *core.Unit, damagePerTick float64) {
		dot := igniteSpell.Dot(target)
		igniteSpell.Cast(sim, target)
		dot.SnapshotBaseDamage = damagePerTick
		dot.Aura.SetStacks(sim, int32(dot.SnapshotBaseDamage))
	}

	procTrigger := core.ProcTrigger{
		Name:               "Ignite Talent",
		CanProcFromProcs:   true, // 11119, 11120, 12846-12848 carry the bit.
		Callback:           core.CallbackOnSpellHitDealt,
		ProcMask:           core.ProcMaskSpellDamage,
		ClassSpellMask:     FireSpellIgnitable,
		Outcome:            core.OutcomeCrit,
		TriggerImmediately: true,
		Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
			target := result.Target
			dot := igniteSpell.Dot(target)
			outstandingDamage := dot.OutstandingDmg()
			if dot.RemainingTicks() <= 0 {
				outstandingDamage = 0
			}

			newDamage := result.Damage * igniteDamageMultiplier
			totalDamage := outstandingDamage + newDamage
			damagePerTick := totalDamage / float64(dot.BaseTickCount)

			refreshIgnite(sim, target, damagePerTick)
		},
	}
	igniteSpell.Unit.MakeProcTriggerAura(procTrigger)
	mage.Ignite = igniteSpell

	// This is needed because we want to listen for the spell "cast" event that refreshes the Dot
	mage.Ignite.Flags ^= core.SpellFlagNoOnCastComplete
}

func (mage *Mage) registerIncineration() {
	if mage.Talents.Incineration == 0 {
		return
	}

	mage.AddStaticMod(core.SpellModConfig{
		ClassMask:  MageSpellFireBlast | MageSpellScorch,
		FloatValue: spellData.Incineration.ValueAt(mage.Talents.Incineration),
		Kind:       core.SpellMod_BonusCrit_Percent,
	})
}

func (mage *Mage) registerImprovedFlamestrike() {
	if mage.Talents.ImprovedFlamestrike == 0 {
		return
	}

	mage.AddStaticMod(core.SpellModConfig{
		ClassMask:  MageSpellFlamestrike,
		FloatValue: spellData.ImprovedFlamestrike.ValueAt(mage.Talents.ImprovedFlamestrike),
		Kind:       core.SpellMod_BonusCrit_Percent,
	})
}

func (mage *Mage) registerBurningSoul() {
	if mage.Talents.BurningSoul == 0 {
		return
	}

	mage.AddStaticMod(core.SpellModConfig{
		School:     core.SpellSchoolFire,
		FloatValue: -.05 * float64(mage.Talents.BurningSoul),
		Kind:       core.SpellMod_ThreatMultiplier_Pct,
	})
}

func (mage *Mage) registerMasterOfElements() {
	if mage.Talents.MasterOfElements == 0 {
		return
	}

	refundCoeff := spellData.MasterOfElements.FractionAt(mage.Talents.MasterOfElements)
	manaMetrics := mage.NewManaMetrics(core.ActionID{SpellID: 29076})

	mage.MakeProcTriggerAura(core.ProcTrigger{
		Name:           "Master of Elements",
		Duration:       core.NeverExpires,
		ClassSpellMask: MageSpellFire | MageSpellFrost,
		Outcome:        core.OutcomeCrit,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.CurCast.Cost == 0 {
				return
			}
			mage.AddMana(sim, spell.DefaultCast.Cost*refundCoeff, manaMetrics)
		},
	})
}

func (mage *Mage) registerCriticalMass() {
	if mage.Talents.CriticalMass == 0 {
		return
	}

	mage.AddStaticMod(core.SpellModConfig{
		SpellFlag:  core.SpellFlag(core.SpellSchoolFire),
		FloatValue: spellData.CriticalMass.ValueAt(mage.Talents.CriticalMass),
		Kind:       core.SpellMod_BonusCrit_Percent,
	})
}

func (mage *Mage) registerFirePower() {
	if mage.Talents.FirePower == 0 {
		return
	}

	mage.AddStaticMod(core.SpellModConfig{
		School:     core.SpellSchoolFire,
		FloatValue: spellData.FirePower.Effect(shared.A_ADD_PCT_MODIFIER, shared.SPELLMOD_DAMAGE).FractionAt(mage.Talents.FirePower),
		Kind:       core.SpellMod_DamageDone_Flat,
	})
}

// ------ FROST TALENTS ------

func (mage *Mage) registerImprovedFrostbolt() {
	if mage.Talents.ImprovedFrostbolt == 0 {
		return
	}

	mage.AddStaticMod(core.SpellModConfig{
		ClassMask: MageSpellFrostbolt,
		TimeValue: time.Millisecond * time.Duration(spellData.ImprovedFrostbolt.Effect(shared.A_ADD_FLAT_MODIFIER, shared.SPELLMOD_CASTING_TIME).ValueAt(mage.Talents.ImprovedFrostbolt)),
		Kind:      core.SpellMod_CastTime_Flat,
	})
}

func (mage *Mage) registerElementalPrecision() {
	if mage.Talents.ElementalPrecision == 0 {
		return
	}
	percent := spellData.ElementalPrecision.Effect(shared.A_ADD_FLAT_MODIFIER, shared.SPELLMOD_RESIST_MISS_CHANCE).ValueAt(mage.Talents.ElementalPrecision)
	mage.AddStaticMod(core.SpellModConfig{
		School:     core.SpellSchoolFrostfire,
		FloatValue: -percent / 100,
		Kind:       core.SpellMod_PowerCost_Pct,
	})

	// Bug: Gives 2% hit per point instead of 1% to frost spells.
	// https://www.warcraftlogs.com/reports/kwd3V8MA9FgrRYhf/?boss=-3&difficulty=0&type=damage-done&source=1&target=2
	mage.PseudoStats.SchoolBonusHitChance[stats.SchoolIndexFrost] += (percent * 2)
	mage.PseudoStats.SchoolBonusHitChance[stats.SchoolIndexFire] += percent
}

func (mage *Mage) registerIceShards() {
	if mage.Talents.IceShards == 0 {
		return
	}

	mage.AddStaticMod(core.SpellModConfig{
		School:     core.SpellSchoolFrost,
		FloatValue: spellData.IceShards.FractionAt(mage.Talents.IceShards),
		Kind:       core.SpellMod_CritMultiplier_Flat,
	})
}

func (mage *Mage) registerImprovedFrostNova() {
	if mage.Talents.ImprovedFrostNova == 0 {
		return
	}

	mage.AddStaticMod(core.SpellModConfig{
		ClassMask: MageSpellFrostNova,
		TimeValue: time.Second * time.Duration(-2*mage.Talents.ImprovedFrostNova),
		Kind:      core.SpellMod_CastTime_Flat,
	})
}

func (mage *Mage) registerPiercingIce() {
	if mage.Talents.PiercingIce == 0 {
		return
	}

	mage.AddStaticMod(core.SpellModConfig{
		ClassMask:  MageSpellFrost,
		FloatValue: spellData.PiercingIce.Effect(shared.A_ADD_PCT_MODIFIER, shared.SPELLMOD_DAMAGE).FractionAt(mage.Talents.PiercingIce),
		Kind:       core.SpellMod_DamageDone_Flat,
	})
}

func (mage *Mage) registerFrostChanneling() {
	if mage.Talents.FrostChanneling == 0 {
		return
	}

	mage.AddStaticMod(core.SpellModConfig{
		ClassMask:  MageSpellFrost,
		FloatValue: -.05 * float64(mage.Talents.FrostChanneling),
		Kind:       core.SpellMod_PowerCost_Pct_Add,
	})

	threatMod := []float64{.04, .07, .1}
	mage.AddStaticMod(core.SpellModConfig{
		School:     core.SpellSchoolFrost,
		FloatValue: -threatMod[mage.Talents.FrostChanneling-1],
		Kind:       core.SpellMod_ThreatMultiplier_Pct,
	})
}

func (mage *Mage) registerImprovedConeOfCold() {
	if mage.Talents.ImprovedConeOfCold == 0 {
		return
	}

	mage.AddStaticMod(core.SpellModConfig{
		ClassMask:  MageSpellConeOfCold,
		FloatValue: .15 + (.10 * (float64(mage.Talents.ImprovedConeOfCold) - 1)),
		Kind:       core.SpellMod_DamageDone_Flat,
	})
}

func (mage *Mage) registerWinterChill() {
	if mage.Talents.WintersChill == 0 {
		return
	}

	// Forever states a flat SpellAuraOptions.ProcChance of 100 on the talent spell and puts the
	// real per-rank chance on the effect, so ProcChanceAt would read 100% at every rank.
	// Effect 0 is the stack count (1..5); effect 1 is the chance (20..100).
	procChance := spellData.WintersChill.EffectAt(1).FractionAt(mage.Talents.WintersChill)

	wcAuras := mage.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return core.WintersChillAura(target, 0)
	})

	mage.Env.RegisterPreFinalizeEffect(func() {
		for _, spell := range mage.GetSpellsMatchingSchool(core.SpellSchoolFrost) {
			spell.RelatedAuraArrays.Append(wcAuras)
		}
	})

	mage.MakeProcTriggerAura(core.ProcTrigger{
		Name:           "Winters Chill Talent",
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		ClassSpellMask: MageSpellFrost,
		ProcChance:     procChance,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			aura := wcAuras.Get(result.Target)
			aura.Activate(sim)
			aura.AddStack(sim)
		},
	})
}

// registerWandSpecialization implements Wand Specialization, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (mage *Mage) registerWandSpecialization() {
	if mage.Talents.WandSpecialization == 0 {
		return
	}

	panic("To be implemented")
}

// registerMagicAbsorption implements Magic Absorption, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (mage *Mage) registerMagicAbsorption() {
	if mage.Talents.MagicAbsorption == 0 {
		return
	}

	panic("To be implemented")
}

// registerImprovedCounterspell implements Improved Counterspell, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (mage *Mage) registerImprovedCounterspell() {
	if mage.Talents.ImprovedCounterspell == 0 {
		return
	}

	panic("To be implemented")
}

// registerImpact implements Impact, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (mage *Mage) registerImpact() {
	if mage.Talents.Impact == 0 {
		return
	}

	panic("To be implemented")
}

// registerFlameThrowing implements Flame Throwing, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (mage *Mage) registerFlameThrowing() {
	if mage.Talents.FlameThrowing == 0 {
		return
	}

	panic("To be implemented")
}

// registerFrostWarding implements Frost Warding, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (mage *Mage) registerFrostWarding() {
	if mage.Talents.FrostWarding == 0 {
		return
	}

	panic("To be implemented")
}

// registerFrostbite implements Frostbite, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (mage *Mage) registerFrostbite() {
	if mage.Talents.Frostbite == 0 {
		return
	}

	panic("To be implemented")
}

// registerPermafrost implements Permafrost, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (mage *Mage) registerPermafrost() {
	if mage.Talents.Permafrost == 0 {
		return
	}

	panic("To be implemented")
}

// registerImprovedBlizzard implements Improved Blizzard, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (mage *Mage) registerImprovedBlizzard() {
	if mage.Talents.ImprovedBlizzard == 0 {
		return
	}

	panic("To be implemented")
}

// registerArcticReach implements Arctic Reach, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (mage *Mage) registerArcticReach() {
	if mage.Talents.ArcticReach == 0 {
		return
	}

	panic("To be implemented")
}

// registerShatter implements Shatter, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (mage *Mage) registerShatter() {
	if mage.Talents.Shatter == 0 {
		return
	}

	panic("To be implemented")
}

// registerIceBarrier implements Ice Barrier, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (mage *Mage) registerIceBarrier() {
	if !mage.Talents.IceBarrier {
		return
	}

	panic("To be implemented")
}

// registerArcaneGeometry implements Arcane Geometry, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (mage *Mage) registerArcaneGeometry() {
	if mage.Talents.ArcaneGeometry == 0 {
		return
	}

	panic("To be implemented")
}

// registerArcaneResilience implements Arcane Resilience, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (mage *Mage) registerArcaneResilience() {
	if mage.Talents.ArcaneResilience == 0 {
		return
	}

	panic("To be implemented")
}

// registerArcaneShielding implements Arcane Shielding, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (mage *Mage) registerArcaneShielding() {
	if mage.Talents.ArcaneShielding == 0 {
		return
	}

	panic("To be implemented")
}

// registerImprovedChanneling implements Improved Channeling, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (mage *Mage) registerImprovedChanneling() {
	if mage.Talents.ImprovedChanneling == 0 {
		return
	}

	panic("To be implemented")
}

// registerMissileBarrage implements Missile Barrage, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (mage *Mage) registerMissileBarrage() {
	if !mage.Talents.MissileBarrage {
		return
	}

	panic("To be implemented")
}

// registerHotStreak implements Hot Streak, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (mage *Mage) registerHotStreak() {
	if !mage.Talents.HotStreak {
		return
	}

	panic("To be implemented")
}

// registerImprovedFireWard implements Improved Fire Ward, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (mage *Mage) registerImprovedFireWard() {
	if mage.Talents.ImprovedFireWard == 0 {
		return
	}

	panic("To be implemented")
}

// registerWakeOfFire implements Wake of Fire, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (mage *Mage) registerWakeOfFire() {
	if mage.Talents.WakeOfFire == 0 {
		return
	}

	panic("To be implemented")
}

// registerFingersOfFrost implements Fingers of Frost, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (mage *Mage) registerFingersOfFrost() {
	if mage.Talents.FingersOfFrost == 0 {
		return
	}

	panic("To be implemented")
}

// registerIceBlock implements Ice Block, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (mage *Mage) registerIceBlock() {
	if !mage.Talents.IceBlock {
		return
	}

	panic("To be implemented")
}

// registerIceLance implements Ice Lance, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (mage *Mage) registerIceLance() {
	if !mage.Talents.IceLance {
		return
	}

	panic("To be implemented")
}

// registerArcaneBlastTalent implements the Forever talent gate for Arcane Blast.
//
// TODO: Forever makes Arcane Blast a talent; the baseline registration in
// arcane_blast.go (registerArcaneBlastSpell, called unconditionally from
// registerSpells) is untouched and should be gated on this field.
func (mage *Mage) registerArcaneBlastTalent() {
	if !mage.Talents.ArcaneBlast {
		return
	}

	panic("To be implemented")
}

// registerPyroblastTalent implements the Forever talent gate for Pyroblast.
//
// TODO: Forever makes Pyroblast a talent; the baseline registration in
// pyroblast.go (registerPyroblastSpell, called unconditionally from
// registerSpells) is untouched and should be gated on this field.
func (mage *Mage) registerPyroblastTalent() {
	if !mage.Talents.Pyroblast {
		return
	}

	panic("To be implemented")
}
