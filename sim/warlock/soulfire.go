package warlock

import (
	"github.com/wowsims/forever/sim/core"
)

func (warlock *Warlock) registerSoulfire() {
	rank := spellData.SoulFire.HighestRank()

	// Bane's cast time cut and Decimation's cooldown cut ride on the talents as SpellMods.
	warlock.Soulfire = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.SpellID},
		SpellSchool:    rank.SpellSchool,
		DefenseType:    rank.DefenseType,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellSoulFire,
		MissileSpeed:   rank.MissileSpeed,

		ManaCost: core.ManaCostOptions{FlatCost: rank.Cost},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      rank.GCD,
				CastTime: rank.CastTime,
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: rank.Cooldown,
			},
		},

		DamageMultiplierAdditive: 1,
		DamageMultiplier:         1,
		ThreatMultiplier:         1,
		BonusCoefficient:         rank.Direct.BonusCoefficient(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcDamage(sim, target, rank.Direct.Damage(sim), spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
