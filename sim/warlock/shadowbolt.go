package warlock

import (
	"github.com/wowsims/forever/sim/core"
)

var shadowBoltRank = spellData.ShadowBolt.BySpellID(27209)
var shadowBoltCoeff = shadowBoltRank.Direct.BonusCoefficient()

func (warlock *Warlock) registerShadowBolt() {

	warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: shadowBoltRank.SpellID},
		SpellSchool:    shadowBoltRank.SpellSchool,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellShadowBolt,
		MissileSpeed:   shadowBoltRank.MissileSpeed,

		ManaCost: core.ManaCostOptions{FlatCost: shadowBoltRank.Cost},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      shadowBoltRank.GCD,
				CastTime: shadowBoltRank.CastTime,
			},
		},

		DamageMultiplierAdditive: 1,
		DefenseType:              shadowBoltRank.DefenseType,
		ThreatMultiplier:         1,
		BonusCoefficient:         shadowBoltCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dmgRoll := shadowBoltRank.Direct.Damage(sim)
			result := spell.CalcDamage(sim, target, dmgRoll, spell.OutcomeMagicHitAndCrit)
			existingAura := target.GetAurasWithTag("ImprovedShadowBolt")

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
				if len(existingAura) == 0 || existingAura[0].Duration != core.NeverExpires {
					if result.Landed() && result.Outcome.Matches(core.OutcomeCrit) && warlock.Talents.ImprovedShadowBolt > 0 {
						if !warlock.ImpShadowboltAura.IsActive() {
							warlock.ImpShadowboltAura.Activate(sim)
						}
						warlock.ImpShadowboltAura.SetStacks(sim, 4)
					}
				}
			})
		},
	})
}
