package buffs

import (
	"fmt"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/stats"
)

// One rank of a paladin aura: the spell the caster used and the number its row states (armor for
// Devotion, damage for Retribution, a percentage for Concentration, resistance for the three
// resistance auras). The paladin registers an aura per rank; each joins the same categories as the
// generated party-buff copy, which is the top rank.
type PaladinAuraRank struct {
	SpellID int32
	Rank    int32
	Value   float64
}

var RetributionAuraMaxRank = PaladinAuraRank{SpellID: 10301, Value: RetributionAuraValue(0)}

func paladinAuraLabel(name string, isPlayer bool, rank PaladinAuraRank) string {
	label := fmt.Sprintf("%s (%s)", name, core.Ternary(isPlayer, "Player", "External"))
	if rank.Rank > 0 {
		label += fmt.Sprintf(" Rank %d", rank.Rank)
	}
	return label
}

func paladinAuraBuff(name string, category string, isPlayer bool, rank PaladinAuraRank) core.GeneratedBuff {
	return core.GeneratedBuff{
		Label:          paladinAuraLabel(name, isPlayer, rank),
		ActionID:       core.ActionID{SpellID: rank.SpellID}.WithTag(core.TernaryInt32(isPlayer, 0, -1)),
		Duration:       core.NeverExpires,
		Category:       category,
		SharedCategory: PaladinAuraCategory,
		SingleAura:     true,
		IsPlayer:       isPlayer,
	}
}

func DevotionAuraBuff(char *core.Character, isPlayer bool, rank PaladinAuraRank) *core.Aura {
	config := paladinAuraBuff("Devotion Aura", DevotionAuraCategory, isPlayer, rank)
	config.Stats = []core.StatConfig{{Stat: stats.Armor, Amount: rank.Value}}
	return core.NewGeneratedStatAura(&char.Unit, config)
}

func ConcentrationAura(char *core.Character, isPlayer bool, rank PaladinAuraRank) *core.Aura {
	config := paladinAuraBuff("Concentration Aura", ConcentrationAuraCategory, isPlayer, rank)
	config.Pseudo = []core.PseudoConfig{{Kind: core.PseudoStatPushbackChance, Amount: -rank.Value / 100}}
	return core.NewGeneratedStatAura(&char.Unit, config)
}

func FireResistanceAura(char *core.Character, isPlayer bool, rank PaladinAuraRank) *core.Aura {
	config := paladinAuraBuff("Fire Resistance Aura", FireResistanceAuraCategory, isPlayer, rank)
	config.Stats = []core.StatConfig{{Stat: stats.FireResistance, Amount: rank.Value}}
	return core.NewGeneratedStatAura(&char.Unit, config)
}

func FrostResistanceAura(char *core.Character, isPlayer bool, rank PaladinAuraRank) *core.Aura {
	config := paladinAuraBuff("Frost Resistance Aura", FrostResistanceAuraCategory, isPlayer, rank)
	config.Stats = []core.StatConfig{{Stat: stats.FrostResistance, Amount: rank.Value}}
	return core.NewGeneratedStatAura(&char.Unit, config)
}

func ShadowResistanceAura(char *core.Character, isPlayer bool, rank PaladinAuraRank) *core.Aura {
	config := paladinAuraBuff("Shadow Resistance Aura", ShadowResistanceAuraCategory, isPlayer, rank)
	config.Stats = []core.StatConfig{{Stat: stats.ShadowResistance, Amount: rank.Value}}
	return core.NewGeneratedStatAura(&char.Unit, config)
}

// Retribution Aura scales with the casting paladin's Holy spell power in Forever even though its
// client row carries no coefficient (every rank and Thorns are the same: EffectBonusCoefficient 0,
// and the damage still moves with spell power in game). The coefficient is the 1.5 s cast-time
// floor over 3.5, the AoE divisor because the shield hits every attacker, and the 0.95 penalty for
// the aura effect. Confirmed at level 20: with 80 spell power rank 1 (base 7) hits for 17 to 18,
// mostly 18, which is the 17.86 this coefficient predicts; the 0.95² variant (0.129) would have
// shown mostly 17.
const RetributionAuraSpellPowerCoefficient = 1.5 / 3.5 / 3 * 0.95

// RetributionAuraBuff is the aura on the unit the shield protects. The self-cast variant reads
// the paladin's own Holy spell power through the proc spell; the external (party-buff) variant
// cannot see the providing paladin, so externalSpellPower stands in for it and the recipient's
// own stats stay out of the damage.
func RetributionAuraBuff(char *core.Character, isPlayer bool, rank PaladinAuraRank, externalSpellPower float64) *core.Aura {
	config := paladinAuraBuff("Retribution Aura", RetributionAuraCategory, isPlayer, rank)
	if char.HasAura(config.Label) {
		return char.GetAura(config.Label)
	}

	if isPlayer {
		return core.NewDamageShield(&char.Unit, config, core.SpellSchoolHoly, rank.Value, RetributionAuraSpellPowerCoefficient)
	}
	return core.NewDamageShield(&char.Unit, config, core.SpellSchoolHoly, rank.Value+RetributionAuraSpellPowerCoefficient*externalSpellPower, 0)
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
func JudgementOfTheCrusaderAura(target *core.Unit, rank JudgementRank) *core.Aura {
	bonus := rank.Value
	label := judgementLabel("Judgement of the Crusader", rank)
	if target.HasAura(label) {
		return target.GetAura(label)
	}

	aura := target.GetOrRegisterAura(core.Aura{
		Label:    label,
		ActionID: core.ActionID{SpellID: rank.SpellID},
		Tag:      JudgementAuraTag,
		Duration: JudgementOfLightDuration(0),
	})

	aura.NewExclusiveEffect("Judgement of the Crusader", true, core.ExclusiveEffect{
		Priority: bonus,
		OnGain: func(ee *core.ExclusiveEffect, sim *core.Simulation) {
			target.PseudoStats.SchoolBonusSpellDamage[stats.SchoolIndexHoly] += bonus
		},
		OnExpire: func(ee *core.ExclusiveEffect, sim *core.Simulation) {
			target.PseudoStats.SchoolBonusSpellDamage[stats.SchoolIndexHoly] -= bonus
		},
	})

	return aura
}

// The paladin's own Judgement of Light at one rank: the bare 40-second debuff, plus the heal it
// grants whoever strikes the target.
func JudgementOfLightRankAura(target *core.Unit, rank JudgementRank) *core.Aura {
	label := judgementLabel("Judgement of Light", rank)
	if target.HasAura(label) {
		return target.GetAura(label)
	}

	return AttachJudgementOfLightHeal(target.GetOrRegisterAura(core.Aura{
		Label:    label,
		ActionID: core.ActionID{SpellID: rank.SpellID},
		Tag:      JudgementAuraTag,
		Duration: JudgementOfLightDuration(0),
	}), rank)
}

// Whoever lands a melee hit on the judged target has a chance to be healed for the rank's amount.
func AttachJudgementOfLightHeal(aura *core.Aura, rank JudgementRank) *core.Aura {
	healthMetrics := aura.Unit.NewHealthMetrics(core.ActionID{SpellID: rank.SpellID})
	heal := rank.Value

	return aura.AttachProcTrigger(core.ProcTrigger{
		Name:     aura.Label + " - Heal",
		Callback: core.CallbackOnSpellHitTaken,
		ProcMask: core.ProcMaskMelee,
		Outcome:  core.OutcomeLanded,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if sim.Proc(judgementProcChance, "Judgement of Light - Heal") {
				spell.Unit.GainHealth(sim, heal, healthMetrics)
			}
		},
	})
}

// The paladin's own Judgement of Wisdom at one rank.
func JudgementOfWisdomRankAura(target *core.Unit, rank JudgementRank) *core.Aura {
	label := judgementLabel("Judgement of Wisdom", rank)
	if target.HasAura(label) {
		return target.GetAura(label)
	}

	return AttachJudgementOfWisdomMana(target.GetOrRegisterAura(core.Aura{
		Label:    label,
		ActionID: core.ActionID{SpellID: rank.SpellID},
		Tag:      JudgementAuraTag,
		Duration: JudgementOfWisdomDuration(0),
	}), rank)
}

// Whoever lands an attack or spell on the judged target has a chance to regain the rank's mana.
// Melee claim it returns mana on a miss as well.
func AttachJudgementOfWisdomMana(aura *core.Aura, rank JudgementRank) *core.Aura {
	actionID := core.ActionID{SpellID: rank.SpellID}
	mana := rank.Value

	return aura.AttachProcTrigger(core.ProcTrigger{
		Name:            aura.Label,
		ActionID:        actionID,
		MetricsActionID: actionID,
		ProcChance:      judgementProcChance,
		ProcMask:        core.ProcMaskDirect,
		Callback:        core.CallbackOnSpellHitTaken,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskMeleeOrRanged) && !result.Landed() {
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
