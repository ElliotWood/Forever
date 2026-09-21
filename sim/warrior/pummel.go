package warrior

import (
	"github.com/wowsims/forever/sim/core"
)

var pummelRank = spellData.Pummel.BySpellID(6554)
var pummelBaseDamage, _ = pummelRank.Direct.Range()

func (warrior *Warrior) registerPummel() {
	warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: pummelRank.SpellID},
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskPummel,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		SpellSchool:    pummelRank.SpellSchool,
		DefenseType:    pummelRank.DefenseType,
		MaxRange:       core.MaxMeleeRange,

		RageCost: core.RageCostOptions{
			Cost: pummelRank.Cost,
			// TODO: Manual review needed -- the 80% rage refund on a miss is the sim's convention; the client states none.
			Refund: 0.8,
		},

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: pummelRank.Cooldown,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.StanceMatches(BerserkerStance)
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealDamage(sim, target, pummelBaseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
}
