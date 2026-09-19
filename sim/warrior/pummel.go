package warrior

import (
	"github.com/wowsims/forever/sim/core"
)

var pummelRank = spellData.Pummel.BySpellID(6554)
var pummelBaseDamage, _ = pummelRank.Direct.Range()

func (war *Warrior) registerPummel() {
	war.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: pummelRank.SpellID},
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskPummel,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		SpellSchool:    pummelRank.SpellSchool,
		DefenseType:    pummelRank.DefenseType,
		MaxRange:       core.MaxMeleeRange,

		RageCost: core.RageCostOptions{
			Cost:   pummelRank.Cost,
			Refund: 0.8,
		},

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    war.NewTimer(),
				Duration: pummelRank.Cooldown,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return war.StanceMatches(BerserkerStance)
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealDamage(sim, target, pummelBaseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
}
