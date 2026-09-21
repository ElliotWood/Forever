package warrior

import (
	"github.com/wowsims/forever/sim/core"
)

func (warrior *Warrior) registerIntercept() {
	interceptRank := spellData.Intercept.HighestRank()
	// The damage sits on the stun the charge triggers, which the generator follows onto the row.
	interceptStunDamage, _ := interceptRank.Direct.Range()

	actionID := core.ActionID{SpellID: interceptRank.SpellID}
	chargeMinRange := interceptRank.MinRange
	interceptCD := interceptRank.Cooldown

	var spell *core.Spell
	var interceptTarget *core.Unit

	aura := warrior.registerDashAura("Intercept", actionID, interceptCD, func(sim *core.Simulation) {
		spell.CalcAndDealDamage(sim, interceptTarget, interceptStunDamage, spell.OutcomeAlwaysHit)
	})

	spell = warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskIntercept,
		MinRange:       chargeMinRange,
		MaxRange:       interceptRank.MaxRange,

		RageCost: core.RageCostOptions{
			Cost: interceptRank.Cost,
		},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: interceptCD,
			},
			IgnoreHaste: true,
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.StanceMatches(BerserkerStance)
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			interceptTarget = target
			aura.Duration = spell.CD.Duration
			aura.Activate(sim)
			warrior.MoveTo(chargeMinRange-3.5, sim) // movement aura is discretized in 1 yard intervals, so need to overshoot to guarantee melee range
		},
	})
}
