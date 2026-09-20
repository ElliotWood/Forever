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

	if debuffs.ImprovedSealOfTheCrusader {
		MakePermanent(ImprovedSealOfTheCrusaderAura(target, -1, 0, 0.0, 1.0))
	}

	if debuffs.Mangle {
		MakePermanent(MangleAura(target))
	}
}

func ScheduledAura(aura *Aura, options PeriodicActionOptions) {
	aura.OnReset = func(aura *Aura, sim *Simulation) {
		aura.Duration = NeverExpires
		StartPeriodicAction(sim, options)
	}
}

func SlowAura(target *Unit) *Aura {
	return castSlowReductionAura(target, "Slow", 31589, 1.5, time.Second*15)
}

func castSlowReductionAura(target *Unit, label string, spellID int32, multiplier float64, duration time.Duration) *Aura {
	aura := target.GetOrRegisterAura(Aura{Label: label, ActionID: ActionID{SpellID: spellID}, Duration: duration})
	aura.NewExclusiveEffect("CastSpdReduction", false, ExclusiveEffect{
		// How far from 1 the applied factor is, which is the scale every member
		// of the category bids on: a 33% slow outbids a 20% one. The factor is
		// 1/multiplier, so a caller stating 1.5 is a 33.3% slow.
		Priority: 1 - 1/multiplier,
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

func ScreechAura(target *Unit) *Aura {
	return statsDebuff(target, 0, "Screech", 27051, stats.Stats{stats.AttackPower: -210}, time.Second*4)
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
