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

	if debuffs.JudgementOfTheCrusader {
		MakePermanent(JudgementOfTheCrusaderAura(target, JudgementOfTheCrusaderMaxRank))
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

// One rank of a judgement debuff: the spell the target shows and the number its row states. The
// paladin registers a rank per row; the debuff panel applies the max rank, carried by the *MaxRank
// values with Rank 0.
type JudgementRank struct {
	SpellID int32
	Rank    int32
	Value   float64
}

var (
	JudgementOfTheCrusaderMaxRank = JudgementRank{SpellID: 20303, Value: 161}
	JudgementOfLightMaxRank       = JudgementRank{SpellID: 20346, Value: 61}
	JudgementOfWisdomMaxRank      = JudgementRank{SpellID: 20355, Value: 59}
)

// Every judgement debuff a paladin puts up carries the tag, so an effect that refreshes "all
// Judgement effects on the target" can find them.
const JudgementAuraTag = "JudgementAura"

// The client says the judgements proc on a chance without stating it; the sim uses 50% until
// in-game testing says otherwise.
const judgementProcChance = 0.5

func judgementLabel(name string, rank JudgementRank) string {
	if rank.Rank > 0 {
		return fmt.Sprintf("%s Rank %d", name, rank.Rank)
	}
	return name
}

// Judgement of the Crusader raises the Holy damage the target takes by a flat amount. Every rank
// and every paladin share one exclusive category, so the strongest active one is the one that
// counts.
func JudgementOfTheCrusaderAura(target *Unit, rank JudgementRank) *Aura {
	bonus := rank.Value
	label := judgementLabel("Judgement of the Crusader", rank)
	if target.HasAura(label) {
		return target.GetAura(label)
	}

	aura := target.GetOrRegisterAura(Aura{
		Label:    label,
		ActionID: ActionID{SpellID: rank.SpellID},
		Tag:      JudgementAuraTag,
		Duration: JudgementOfLightDuration(0),
	})

	aura.NewExclusiveEffect("Judgement of the Crusader", true, ExclusiveEffect{
		Priority: bonus,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			target.PseudoStats.SchoolBonusSpellDamage[stats.SchoolIndexHoly] += bonus
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			target.PseudoStats.SchoolBonusSpellDamage[stats.SchoolIndexHoly] -= bonus
		},
	})

	return aura
}

// The paladin's own Judgement of Light at one rank: the bare 40-second debuff, plus the heal it
// grants whoever strikes the target.
func JudgementOfLightRankAura(target *Unit, rank JudgementRank) *Aura {
	label := judgementLabel("Judgement of Light", rank)
	if target.HasAura(label) {
		return target.GetAura(label)
	}

	return AttachJudgementOfLightHeal(target.GetOrRegisterAura(Aura{
		Label:    label,
		ActionID: ActionID{SpellID: rank.SpellID},
		Tag:      JudgementAuraTag,
		Duration: JudgementOfLightDuration(0),
	}), rank)
}

// Whoever lands a melee hit on the judged target has a chance to be healed for the rank's amount.
func AttachJudgementOfLightHeal(aura *Aura, rank JudgementRank) *Aura {
	healthMetrics := aura.Unit.NewHealthMetrics(ActionID{SpellID: rank.SpellID})
	heal := rank.Value

	return aura.AttachProcTrigger(ProcTrigger{
		Name:     aura.Label + " - Heal",
		Callback: CallbackOnSpellHitTaken,
		ProcMask: ProcMaskMelee,
		Outcome:  OutcomeLanded,
		Handler: func(sim *Simulation, spell *Spell, result *SpellResult) {
			if sim.Proc(judgementProcChance, "Judgement of Light - Heal") {
				spell.Unit.GainHealth(sim, heal, healthMetrics)
			}
		},
	})
}

// The paladin's own Judgement of Wisdom at one rank.
func JudgementOfWisdomRankAura(target *Unit, rank JudgementRank) *Aura {
	label := judgementLabel("Judgement of Wisdom", rank)
	if target.HasAura(label) {
		return target.GetAura(label)
	}

	return AttachJudgementOfWisdomMana(target.GetOrRegisterAura(Aura{
		Label:    label,
		ActionID: ActionID{SpellID: rank.SpellID},
		Tag:      JudgementAuraTag,
		Duration: JudgementOfWisdomDuration(0),
	}), rank)
}

// Whoever lands an attack or spell on the judged target has a chance to regain the rank's mana.
// Melee claim it returns mana on a miss as well.
func AttachJudgementOfWisdomMana(aura *Aura, rank JudgementRank) *Aura {
	actionID := ActionID{SpellID: rank.SpellID}
	mana := rank.Value

	return aura.AttachProcTrigger(ProcTrigger{
		Name:            aura.Label,
		ActionID:        actionID,
		MetricsActionID: actionID,
		ProcChance:      judgementProcChance,
		ProcMask:        ProcMaskDirect,
		Callback:        CallbackOnSpellHitTaken,
		Handler: func(sim *Simulation, spell *Spell, result *SpellResult) {
			if !spell.ProcMask.Matches(ProcMaskMeleeOrRanged) && !result.Landed() {
				return
			}

			unit := spell.Unit
			if !unit.HasManaBar() {
				return
			}
			if unit.JowManaMetrics == nil {
				unit.JowManaMetrics = unit.NewManaMetrics(actionID)
			}
			unit.AddMana(sim, mana, unit.JowManaMetrics)
		},
	})
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
