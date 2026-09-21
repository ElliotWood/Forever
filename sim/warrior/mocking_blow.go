package warrior

import (
	"github.com/wowsims/forever/sim/core"
)

var mockingBlowRank = spellData.MockingBlow.HighestRank()
var mockingBlowBaseDamage, _ = mockingBlowRank.Direct.Range()

func (warrior *Warrior) registerMockingBlow() {
	warrior.MockingBlow = warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: mockingBlowRank.SpellID},
		SpellSchool:    mockingBlowRank.SpellSchool,
		DefenseType:    mockingBlowRank.DefenseType,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskMockingBlow,
		MaxRange:       mockingBlowRank.MaxRange,

		RageCost: core.RageCostOptions{
			Cost:   mockingBlowRank.Cost,
			Refund: mockingBlowRank.MissRefund(),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: mockingBlowRank.GCD,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: mockingBlowRank.Cooldown,
			},
		},

		DamageMultiplier: 1,
		// TODO: Test in-game
		ThreatMultiplier: 1,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.StanceMatches(BattleStance)
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealDamage(sim, target, mockingBlowBaseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
}
