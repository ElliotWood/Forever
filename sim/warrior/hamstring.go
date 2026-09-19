package warrior

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

var hamstringRank = shared.WithSpellDataFlatThreat(spellData.Hamstring, 167.5).BySpellID(25212)
var hamstringBaseDamage, _ = hamstringRank.Direct.Range()

func (war *Warrior) registerHamstring() {
	war.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: hamstringRank.SpellID},
		SpellSchool:    hamstringRank.SpellSchool,
		DefenseType:    hamstringRank.DefenseType,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskHamstring,
		MaxRange:       core.MaxMeleeRange,

		RageCost: core.RageCostOptions{
			Cost:   hamstringRank.Cost,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: hamstringRank.GCD,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1.25,
		FlatThreatBonus:  hamstringRank.FlatThreatBonus,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return war.StanceMatches(BerserkerStance)
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealDamage(sim, target, hamstringBaseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
}
