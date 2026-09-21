package druid

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

// Forever's Starfire tops out at rank 7, so the max-rank entry moves down from TBC's 8 rather than
// naming a rank the table does not hold. Rank 6 stays registered for downranking.
var StarfireRankMap = spellData.Starfire.Ranks(6, 7)

func (druid *Druid) registerStarfireSpell(rankConfig shared.SpellData) {
	spell := druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rankConfig.SpellID},
		SpellSchool:    rankConfig.SpellSchool,
		DefenseType:    rankConfig.DefenseType,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: DruidSpellStarfire,
		Flags:          core.SpellFlagAPL,
		Rank:           rankConfig.Rank,
		MaxRange:       rankConfig.MaxRange,

		ManaCost: core.ManaCostOptions{
			FlatCost: rankConfig.Cost,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      rankConfig.GCD,
				CastTime: rankConfig.CastTime,
			},
		},

		BonusCoefficient: rankConfig.Direct.BonusCoefficient(),
		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, rankConfig.Direct.Damage(sim), spell.OutcomeMagicHitAndCrit)
		},
	})

	druid.Starfire = append(druid.Starfire, spell)
}
