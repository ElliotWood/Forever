package warrior

import (
	"github.com/wowsims/forever/sim/core"
)

func (warrior *Warrior) registerHamstring() {
	// TODO: Ingame research needed if this adds flat threat
	hamstringRank := spellData.Hamstring.HighestRank()
	hamstringBaseDamage, _ := hamstringRank.Direct.Range()

	warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: hamstringRank.SpellID},
		SpellSchool:    hamstringRank.SpellSchool,
		DefenseType:    hamstringRank.DefenseType,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskHamstring,
		MaxRange:       core.MaxMeleeRange,

		RageCost: core.RageCostOptions{
			Cost:   hamstringRank.Cost,
			Refund: hamstringRank.MissRefund(),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: hamstringRank.GCD,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		// TODO: Manual review needed -- the client states no threat coefficient; 1 until measured in game.
		ThreatMultiplier: 1,
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
