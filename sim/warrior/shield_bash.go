package warrior

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

// TODO: Manual review needed -- spell 1672 carries no threat effect, so the 192 is hand-supplied.
var shieldBashRank = shared.WithSpellDataFlatThreat(spellData.ShieldBash, 192).HighestRank()

func (warrior *Warrior) registerShieldBash() {
	actionID := core.ActionID{SpellID: shieldBashRank.SpellID}

	warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: SpellMaskShieldBash,
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		MaxRange:       core.MaxMeleeRange,

		RageCost: core.RageCostOptions{
			Cost: shieldBashRank.Cost,
			// TODO: Manual review needed -- the 80% rage refund on a miss is the sim's convention; the client states none.
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: shieldBashRank.GCD,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: shieldBashRank.Cooldown,
			},
		},

		DamageMultiplier: 1,
		// TODO: Manual review needed -- the threat coefficient is not in the client.
		ThreatMultiplier: 1.5,
		FlatThreatBonus:  shieldBashRank.FlatThreatBonus,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.PseudoStats.CanBlock && warrior.StanceMatches(DefensiveStance|BattleStance)
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := shieldBashRank.Direct.Damage(sim)
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
}
