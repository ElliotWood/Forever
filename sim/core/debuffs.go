package core

import (
	"fmt"
	"time"

	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

// applyRaidDebuffEffects applies all raid-level debuffs based on the provided Debuffs proto.
func applyDebuffEffects(target *Unit, targetIdx int, debuffs *proto.Debuffs, raid *proto.Raid) {
	applyGeneratedDebuffs(target, debuffs, raid)

	if debuffs.BloodFrenzy {
		MakePermanent(BloodFrenzyAura(target, 2))
	}

	if debuffs.ImprovedSealOfTheCrusader {
		MakePermanent(ImprovedSealOfTheCrusaderAura(target, -1, 0, 0.0, Ternary(debuffs.JocRetribution_2Pt4, 1.15, 1.0)))
	}

	if debuffs.Mangle {
		MakePermanent(MangleAura(target))
	}

	if debuffs.Misery {
		MakePermanent(MiseryAura(target, 5))
	}

	if debuffs.Screech {
		MakePermanent(ScreechAura(target))
	}

	if debuffs.ShadowEmbrace {
		MakePermanent(ShadowEmbraceAura(target, 5))
	}
}

func ScheduledAura(aura *Aura, options PeriodicActionOptions) {
	aura.OnReset = func(aura *Aura, sim *Simulation) {
		aura.Duration = NeverExpires
		StartPeriodicAction(sim, options)
	}
}

// Physical and Armor Related Debuffs
func BloodFrenzyAura(target *Unit, points int32) *Aura {
	return damageTakenDebuff(target, 0,
		"Blood Frenzy",
		29859,
		[]stats.SchoolIndex{stats.SchoolIndexPhysical},
		1+0.02*float64(points),
		NeverExpires,
	)
}

func SlowAura(target *Unit) *Aura {
	return castSlowReductionAura(target, "Slow", 31589, 1.5, time.Second*15)
}

func castSlowReductionAura(target *Unit, label string, spellID int32, multiplier float64, duration time.Duration) *Aura {
	aura := target.GetOrRegisterAura(Aura{Label: label, ActionID: ActionID{SpellID: spellID}, Duration: duration})
	aura.NewExclusiveEffect("CastSpdReduction", false, ExclusiveEffect{
		Priority: multiplier,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.MultiplyCastSpeed(sim, 1/multiplier)
			ee.Aura.Unit.MultiplyRangedSpeed(sim, 1/multiplier)
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.MultiplyCastSpeed(sim, multiplier)
			ee.Aura.Unit.MultiplyRangedSpeed(sim, multiplier)
		},
	})
	return aura
}

// points is number of talent points in improved seal of the crusader
//
// flatBonus is used when the character has a flat bonus to the holy damage taken
//
// percentBonus is used when the character has a percent bonus to the holy damage taken
func ImprovedSealOfTheCrusaderAura(target *Unit, casterIndex, points int32, flatBonus, percentBonus float64) *Aura {
	holySpellDamageBonus := 219.0*percentBonus + flatBonus //assumed Max Rank Seal Of Crusader (Rank 7)

	auraLabel := fmt.Sprintf("Improved Seal of the Crusader (%s)", Ternary(casterIndex == -1, "External", "Self"))

	actionID := ActionID{SpellID: 20337}
	if casterIndex != 0 {
		actionID = actionID.WithTag(casterIndex)
	}

	aura := target.GetOrRegisterAura(Aura{
		Label:    auraLabel,
		ActionID: actionID,
		Duration: time.Second * 20,
	})

	aura.NewExclusiveEffect("Improved Seal of the Crusader", true, ExclusiveEffect{
		Priority: holySpellDamageBonus + float64(casterIndex),
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			target.AddReducedCritTakenPercent(float64(-1 * points))
			target.PseudoStats.SchoolBonusSpellDamage[stats.SchoolIndexHoly] += holySpellDamageBonus
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			target.AddReducedCritTakenPercent(float64(1 * points))
			target.PseudoStats.SchoolBonusSpellDamage[stats.SchoolIndexHoly] -= holySpellDamageBonus
		},
	})

	return aura
}

func MangleAura(target *Unit) *Aura {
	multiplier := 1.3

	aura := target.GetOrRegisterAura(Aura{
		Label:    "Mangle",
		ActionID: ActionID{SpellID: 33876},
		Duration: time.Second * 12,
	})

	aura.NewExclusiveEffect("Mangle", true, ExclusiveEffect{
		Priority: multiplier,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.PeriodicPhysicalDamageTakenMultiplier *= ee.Priority
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.PeriodicPhysicalDamageTakenMultiplier /= ee.Priority
		},
	})

	return aura
}

func MiseryAura(target *Unit, ranks int32) *Aura {
	multiplier := 1.0 + 0.01*float64(ranks)
	schools := []stats.SchoolIndex{
		stats.SchoolIndexArcane, stats.SchoolIndexFire, stats.SchoolIndexFrost,
		stats.SchoolIndexHoly, stats.SchoolIndexNature, stats.SchoolIndexShadow,
	}
	aura := target.GetOrRegisterAura(Aura{
		Label:    "Misery",
		ActionID: ActionID{SpellID: 33195},
		Duration: NeverExpires,
	})
	effect := aura.NewExclusiveEffect("MiseryBonus", true, ExclusiveEffect{
		Priority: multiplier,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			for _, school := range schools {
				target.PseudoStats.SchoolDamageTakenMultiplier[school] *= ee.Priority
			}
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			for _, school := range schools {
				target.PseudoStats.SchoolDamageTakenMultiplier[school] /= ee.Priority
			}
		},
	})
	if effect.Priority < multiplier {
		effect.Priority = multiplier
	}
	return aura
}

func ScreechAura(target *Unit) *Aura {
	return statsDebuff(target, 0, "Screech", 27051, stats.Stats{stats.AttackPower: -210}, time.Second*4)
}

func ShadowEmbraceAura(target *Unit, ranks int32) *Aura {
	return damageDealtDebuff(target, "Shadow Embrace", 32394, []stats.SchoolIndex{stats.SchoolIndexPhysical}, 1.0-(.01*float64(ranks)), NeverExpires)
}

func StormstrikeAura(target *Unit, uptime float64) *Aura {
	multiplier := 1.20
	hasAura := target.HasAura("Stormstrike")
	aura := damageTakenDebuff(target, 0, "Stormstrike", 17364, []stats.SchoolIndex{stats.SchoolIndexNature}, multiplier, time.Second*12)

	if !hasAura {
		ApplyFixedUptimeAura(aura, uptime, aura.Duration, 1)
	}

	return aura
}

func AtkSpeedReductionEffect(aura *Aura, speedMultiplier float64) *ExclusiveEffect {
	return aura.NewExclusiveEffect("AtkSpdReduction", false, ExclusiveEffect{
		// How far from 1 the applied factor is, which is the scale every member
		// of the category bids on: a 20% slow outbids a 10% one. The factor is
		// 1/speedMultiplier, so a helper stating 1.2 is a 16.67% slow.
		Priority: 1 - 1/speedMultiplier,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.MultiplyAttackSpeed(sim, 1/speedMultiplier)
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.MultiplyAttackSpeed(sim, speedMultiplier)
		},
	})
}

func damageTakenDebuff(target *Unit, casterIndex int32, label string, spellID int32, schools []stats.SchoolIndex, multiplier float64, duration time.Duration) *Aura {
	actionID := ActionID{SpellID: spellID}
	if casterIndex != 0 {
		actionID = actionID.WithTag(casterIndex)
	}

	return target.GetOrRegisterAura(Aura{
		Label:    label,
		ActionID: actionID,
		Duration: duration,
		OnGain: func(aura *Aura, sim *Simulation) {
			for _, school := range schools {
				target.PseudoStats.SchoolDamageTakenMultiplier[school] *= multiplier
			}
		},

		OnExpire: func(aura *Aura, sim *Simulation) {
			for _, school := range schools {
				target.PseudoStats.SchoolDamageTakenMultiplier[school] /= multiplier
			}
		},
	})
}

func damageDealtDebuff(target *Unit, label string, spellID int32, schools []stats.SchoolIndex, multiplier float64, duration time.Duration) *Aura {
	return target.GetOrRegisterAura(Aura{
		Label:    label,
		ActionID: ActionID{SpellID: spellID},
		Duration: duration,

		OnGain: func(aura *Aura, sim *Simulation) {
			for _, school := range schools {
				target.PseudoStats.SchoolDamageDealtMultiplier[school] *= multiplier
			}
		},

		OnExpire: func(aura *Aura, sim *Simulation) {
			for _, school := range schools {
				target.PseudoStats.SchoolDamageDealtMultiplier[school] /= multiplier
			}
		},
	})
}

func statsDebuff(target *Unit, casterIndex int32, label string, spellID int32, stats stats.Stats, duration time.Duration) *Aura {
	if duration == 0 {
		duration = time.Second * 30
	}

	actionID := ActionID{SpellID: spellID}
	if casterIndex != 0 {
		actionID = actionID.WithTag(casterIndex)
	}

	aura := target.GetAuraByID(actionID)
	if aura != nil {
		return aura
	}

	return target.RegisterAura(Aura{
		Label:    label,
		ActionID: actionID,
		Duration: duration,
	}).AttachStatsBuff(stats)
}
