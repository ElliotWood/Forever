package mage

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

// TODO: was Ranks(7, 6); Forever's Flamestrike tops out at rank 6, so only that remains.
var FlameStrikeRankMap = spellData.Flamestrike.Ranks(6)

// Flamestrike
// https://www.wowhead.com/forever/spell=10216
//
// Calls down a pillar of fire, burning all enemies within the area for X Fire damage and an additional
// Y Fire damage over 8 sec. Only one Flamestrike can be active per Mage at a time.
func (mage *Mage) registerFlamestrike(rankConfig shared.SpellData) {
	actionID := core.ActionID{SpellID: rankConfig.SpellID}
	tick := rankConfig.Periodic.(shared.SpellDataPeriodic)

	spell := mage.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    rankConfig.SpellSchool,
		DefenseType:    rankConfig.DefenseType,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MageSpellFlamestrike,
		Rank:           rankConfig.Rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: rankConfig.Cost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      rankConfig.GCD,
				CastTime: rankConfig.CastTime,
			},
		},

		DamageMultiplier: 1,
		BonusCoefficient: rankConfig.Direct.BonusCoefficient(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				ActionID: actionID,
				Label:    "Flamestrike" + mage.Label + " " + rankConfig.GetRankLabel(),
			},
			NumberOfTicks:    tick.NumberOfTicks,
			TickLength:       tick.TickLength,
			BonusCoefficient: tick.Coef,
			OnTick: func(sim *core.Simulation, _ *core.Unit, dot *core.Dot) {
				dot.Spell.CalcAndDealPeriodicAoeDamage(sim, tick.Tick, dot.OutcomeTickMagicHit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			spell.CalcAndDealAoeDamage(sim, rankConfig.Direct.Damage(sim), spell.OutcomeMagicHitAndCrit)
			spell.AOEDot().Apply(sim)
		},
	})

	mage.Flamestrike = append(mage.Flamestrike, spell)
}
