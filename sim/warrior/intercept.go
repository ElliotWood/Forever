package warrior

import (
	"github.com/wowsims/forever/sim/core"
)

var interceptRank = spellData.Intercept.HighestRank()

var interceptStunRank = spellData.InterceptTriggered.ByRank(interceptRank.Rank)
var interceptStunDamage, _ = interceptStunRank.Direct.Range()

func (warrior *Warrior) registerIntercept() {
	actionID := core.ActionID{SpellID: interceptRank.SpellID}
	chargeMinRange := interceptRank.MinRange
	interceptCD := interceptRank.Cooldown

	var spell *core.Spell
	var interceptTarget *core.Unit

	aura := warrior.RegisterAura(core.Aura{
		Label:    "Intercept",
		ActionID: actionID,
		Duration: interceptCD,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			// TODO: Manual review needed -- the run speed and the overshoot below are the sim's movement model.
			warrior.MultiplyMovementSpeed(sim, 3.0)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.MultiplyMovementSpeed(sim, 1.0/3.0)
			spell.CalcAndDealDamage(sim, interceptTarget, interceptStunDamage, spell.OutcomeAlwaysHit)
		},
	})

	warrior.RegisterMovementCallback(func(sim *core.Simulation, position float64, kind core.MovementUpdateType) {
		if kind == core.MovementEnd && aura.IsActive() {
			aura.Deactivate(sim)
		}
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
