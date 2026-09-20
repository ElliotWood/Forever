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
			Refund: 0.8,
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
		ThreatMultiplier: 1,

		// Spell 20560's ShapeshiftMask is Battle Stance only.
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.StanceMatches(BattleStance)
		},

		// TODO: the taunt half of spell 20560 is not modelled; the sim has no threat table to
		// force the target onto.
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealDamage(sim, target, mockingBlowBaseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
}
