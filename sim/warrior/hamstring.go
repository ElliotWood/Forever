package warrior

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

// TODO: Ingame research needed if this adds flat threat
var hamstringRank = shared.WithSpellDataFlatThreat(spellData.Hamstring, 0).HighestRank()
var hamstringBaseDamage, _ = hamstringRank.Direct.Range()

func (warrior *Warrior) registerHamstring() {
	warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: hamstringRank.SpellID},
		SpellSchool:    hamstringRank.SpellSchool,
		DefenseType:    hamstringRank.DefenseType,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskHamstring,
		MaxRange:       core.MaxMeleeRange,

		RageCost: core.RageCostOptions{
			Cost: hamstringRank.Cost,
			// TODO: Manual review needed -- the 80% rage refund on a miss is the sim's convention; the client states none.
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: hamstringRank.GCD,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		// TODO: Manual review needed -- the threat coefficient is not in the client.
		ThreatMultiplier: 1.25,
		FlatThreatBonus:  hamstringRank.FlatThreatBonus,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.StanceMatches(BattleStance | BerserkerStance)
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealDamage(sim, target, hamstringBaseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
}
