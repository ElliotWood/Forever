package warrior

import (
	"github.com/wowsims/forever/sim/core"
)

var hamstringRank = spellData.Hamstring.Highest()
var hamstringBaseDamage = hamstringRank.DamageEffect().Average(core.CharacterLevel)

func (warrior *Warrior) registerHamstring() {
	warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: hamstringRank.ID},
		SpellSchool:    hamstringRank.SpellSchool(),
		DefenseType:    hamstringRank.DefenseTypeCore(),
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskHamstring,
		ClassFlags:     SpellFlagsHamstring,
		MaxRange:       core.MaxMeleeRange,

		RageCost: core.RageCostOptions{
			Cost:   rageCost(hamstringRank),
			Refund: hamstringRank.MissRefund(),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: hamstringRank.GCD(),
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		// TODO: Manual review needed -- the client states no threat coefficient; 1 until measured in game.
		ThreatMultiplier: 1,
		// TODO: Ingame research needed if this adds flat threat
		FlatThreatBonus: 0,

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
