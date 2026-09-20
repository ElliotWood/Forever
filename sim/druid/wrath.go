package druid

import (
	"github.com/wowsims/forever/sim/core"
)

var wrathRank = spellData.Wrath.HighestRank()

func (druid *Druid) registerWrathSpell() {
	druid.Wrath = druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
		ActionID:       core.ActionID{SpellID: wrathRank.SpellID},
		SpellSchool:    wrathRank.SpellSchool,
		DefenseType:    wrathRank.DefenseType,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: DruidSpellWrath,
		Flags:          core.SpellFlagAPL,
		MissileSpeed:   wrathRank.MissileSpeed,

		ManaCost: core.ManaCostOptions{
			FlatCost: wrathRank.Cost,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      wrathRank.GCD,
				CastTime: wrathRank.CastTime,
			},
		},

		BonusCoefficient: wrathRank.Direct.BonusCoefficient(),
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		MaxRange:         wrathRank.MaxRange,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := wrathRank.Direct.Damage(sim)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
