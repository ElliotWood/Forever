package warlock

import (
	"time"

	"github.com/wowsims/forever/sim/core"
)

var soulfireRank = spellData.SoulFire.BySpellID(30545)
var soulfireCoeff = soulfireRank.Direct.BonusCoefficient()

func (warlock *Warlock) registerSoulfire() {
	warlock.Soulfire = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: soulfireRank.SpellID},
		SpellSchool:    soulfireRank.SpellSchool,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellSoulFire,
		MissileSpeed:   soulfireRank.MissileSpeed,

		ManaCost: core.ManaCostOptions{FlatCost: soulfireRank.Cost},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      soulfireRank.GCD,
				CastTime: soulfireRank.CastTime - time.Duration(400*warlock.Talents.Bane),
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: soulfireRank.Cooldown,
			},
		},

		DamageMultiplier: 1,
		DefenseType:      soulfireRank.DefenseType,
		ThreatMultiplier: 1,
		BonusCoefficient: soulfireCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dmgRoll := soulfireRank.Direct.Damage(sim)
			result := spell.CalcDamage(sim, target, dmgRoll, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})

}
