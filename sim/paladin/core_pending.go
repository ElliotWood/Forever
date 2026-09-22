package paladin

import (
	"fmt"
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/stats"
)

// The paladin auras, judgement debuffs and Blessing of Light marker in the Forever shape their #48
// stack moves into sim/core (commit "Give the paladin buffs and debuffs their Forever shape"). The
// shared core on forever-next still builds the TBC versions for the raid-buff panel, and that
// rewrite is their #39's to land, so the self-cast ranks the paladin registers are built here from
// the paladin's own rows instead.
// ponytail: delete this file and point the callers back at core.* once #39/#48's core lands.

// One rank of a paladin aura: the spell cast and the number its row states.
type paladinAuraRank struct {
	SpellID int32
	Rank    int32
	Value   float64
}

// The self-cast aura: one per rank, exclusive with the other paladin auras and winning over the
// raid-buff panel's copy of the same aura.
func selfCastPaladinAura(char *core.Character, name, category string, rank paladinAuraRank) *core.Aura {
	aura := char.GetOrRegisterAura(core.Aura{
		Label:    fmt.Sprintf("%s (Player) Rank %d", name, rank.Rank),
		ActionID: core.ActionID{SpellID: rank.SpellID},
		Duration: core.NeverExpires,
	})
	aura.NewExclusiveEffect(category, true, core.ExclusiveEffect{Priority: 1})
	aura.NewExclusiveEffect(core.PaladinAuraCategory, true, core.ExclusiveEffect{})
	return aura
}

func devotionAuraBuff(char *core.Character, rank paladinAuraRank) *core.Aura {
	return selfCastPaladinAura(char, "Devotion Aura", core.DevotionAuraCategory, rank).AttachStatBuff(stats.Armor, rank.Value)
}

func fireResistanceAura(char *core.Character, rank paladinAuraRank) *core.Aura {
	return selfCastPaladinAura(char, "Fire Resistance Aura", core.FireResistanceAuraCategory, rank).AttachStatBuff(stats.FireResistance, rank.Value)
}

func frostResistanceAura(char *core.Character, rank paladinAuraRank) *core.Aura {
	return selfCastPaladinAura(char, "Frost Resistance Aura", core.FrostResistanceAuraCategory, rank).AttachStatBuff(stats.FrostResistance, rank.Value)
}

func shadowResistanceAura(char *core.Character, rank paladinAuraRank) *core.Aura {
	return selfCastPaladinAura(char, "Shadow Resistance Aura", core.ShadowResistanceAuraCategory, rank).AttachStatBuff(stats.ShadowResistance, rank.Value)
}

func concentrationAura(char *core.Character, rank paladinAuraRank) *core.Aura {
	return selfCastPaladinAura(char, "Concentration Aura", core.ConcentrationAuraCategory, rank).
		AttachAdditivePseudoStatBuff(&char.PseudoStats.PushbackChance, -rank.Value/100)
}

func retributionAuraBuff(char *core.Character, rank paladinAuraRank) *core.Aura {
	damage := rank.Value
	procSpell := char.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: rank.SpellID}.WithTag(2),
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagBinary | core.SpellFlagPassiveSpell,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeAlwaysHit)
		},
	})

	return selfCastPaladinAura(char, "Retribution Aura", core.RetributionAuraCategory, rank).AttachProcTrigger(core.ProcTrigger{
		Callback: core.CallbackOnSpellHitTaken,
		Outcome:  core.OutcomeLanded,
		Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
			if spell.SpellSchool.Matches(core.SpellSchoolPhysical) {
				procSpell.Cast(sim, spell.Unit)
			}
		},
	})
}

// One rank of a judgement debuff: the spell the target shows and the number its row states.
type judgementRank struct {
	SpellID int32
	Rank    int32
	Value   float64
}

// Forever's judgements last 40 sec (the TBC debuffs in core last 20).
const judgementDuration = time.Second * 40

// Every judgement debuff carries the tag, so an effect that refreshes "all Judgement effects on the
// target" can find them.
const judgementAuraTag = "JudgementAura"

// The client says the judgements proc on a chance; the 50% the TBC sim settled on stays until
// Forever testing says otherwise.
const judgementProcChance = 0.5

func judgementAura(target *core.Unit, name string, rank judgementRank) *core.Aura {
	return target.GetOrRegisterAura(core.Aura{
		Label:    fmt.Sprintf("%s Rank %d", name, rank.Rank),
		ActionID: core.ActionID{SpellID: rank.SpellID},
		Tag:      judgementAuraTag,
		Duration: judgementDuration,
	})
}

// Judgement of the Crusader raises the Holy damage the target takes by a flat amount; every rank
// and every paladin share one exclusive category, so the strongest one counts.
func judgementOfTheCrusaderAura(target *core.Unit, rank judgementRank) *core.Aura {
	label := fmt.Sprintf("Judgement of the Crusader Rank %d", rank.Rank)
	if target.HasAura(label) {
		return target.GetAura(label)
	}

	bonus := rank.Value
	aura := judgementAura(target, "Judgement of the Crusader", rank)
	aura.NewExclusiveEffect("Judgement of the Crusader", true, core.ExclusiveEffect{
		Priority: bonus,
		OnGain: func(_ *core.ExclusiveEffect, _ *core.Simulation) {
			target.PseudoStats.SchoolBonusSpellDamage[stats.SchoolIndexHoly] += bonus
		},
		OnExpire: func(_ *core.ExclusiveEffect, _ *core.Simulation) {
			target.PseudoStats.SchoolBonusSpellDamage[stats.SchoolIndexHoly] -= bonus
		},
	})
	return aura
}

// Judgement of Light heals whoever lands a melee hit on the target.
func judgementOfLightAura(target *core.Unit, rank judgementRank) *core.Aura {
	label := fmt.Sprintf("Judgement of Light Rank %d", rank.Rank)
	if target.HasAura(label) {
		return target.GetAura(label)
	}

	healthMetrics := target.NewHealthMetrics(core.ActionID{SpellID: rank.SpellID})
	heal := rank.Value
	return judgementAura(target, "Judgement of Light", rank).AttachProcTrigger(core.ProcTrigger{
		Callback:   core.CallbackOnSpellHitTaken,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeLanded,
		ProcChance: judgementProcChance,
		Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
			spell.Unit.GainHealth(sim, heal, healthMetrics)
		},
	})
}

// Judgement of Wisdom returns mana to whoever lands a direct attack or spell on the target.
func judgementOfWisdomAura(target *core.Unit, rank judgementRank) *core.Aura {
	label := fmt.Sprintf("Judgement of Wisdom Rank %d", rank.Rank)
	if target.HasAura(label) {
		return target.GetAura(label)
	}

	actionID := core.ActionID{SpellID: rank.SpellID}
	mana := rank.Value
	return judgementAura(target, "Judgement of Wisdom", rank).AttachProcTrigger(core.ProcTrigger{
		Callback:   core.CallbackOnSpellHitTaken,
		ProcMask:   core.ProcMaskDirect,
		ProcChance: judgementProcChance,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// Melee claim that wisdom can proc on misses.
			if !spell.ProcMask.Matches(core.ProcMaskMeleeOrRanged) && !result.Landed() {
				return
			}
			unit := spell.Unit
			if unit.HasManaBar() {
				if unit.JowManaMetrics == nil {
					unit.JowManaMetrics = unit.NewManaMetrics(actionID)
				}
				unit.AddMana(sim, mana, unit.JowManaMetrics)
			}
		},
	})
}

// Core has no Blessing of Light on forever-next yet, so no target carries it and the heals never
// add its bonus.
const blessingOfLightAuraLabel = "Blessing of Light"
