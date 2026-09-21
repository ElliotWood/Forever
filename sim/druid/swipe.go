package druid

import (
	"github.com/wowsims/forever/sim/core"
)

var swipeRank = spellData.Swipe.HighestRank()

func (druid *Druid) registerSwipeBearSpell() {
	druid.Swipe = druid.RegisterSpell(Bear, core.SpellConfig{
		ActionID:       core.ActionID{SpellID: swipeRank.SpellID},
		SpellSchool:    swipeRank.SpellSchool,
		DefenseType:    swipeRank.DefenseType,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		ClassSpellMask: DruidSpellSwipe,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		RageCost: core.RageCostOptions{
			Cost:   swipeRank.Cost,
			Refund: swipeRank.MissRefund(),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: swipeRank.GCD,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		// Season of Discovery's "Modifies Threat +101%", which the client does not carry.
		ThreatMultiplier: 2,
		MaxRange:         core.MaxMeleeRange,

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			numHits := min(3, len(druid.Env.Encounter.AllTargetUnits))
			for i := 0; i < numHits; i++ {
				aoeTarget := druid.Env.Encounter.AllTargetUnits[i]
				spell.CalcAndDealDamage(sim, aoeTarget, swipeRank.Direct.Damage(sim), spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			}
		},
	})
}
