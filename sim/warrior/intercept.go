package warrior

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/spelldata"
)

var interceptRank = spellData.Intercept.Highest()

// The damage sits on the stun the charge triggers, not on Intercept itself.
var interceptStunDamage = spellData.InterceptTriggered.Highest().DamageEffect().Average(core.CharacterLevel)

func (warrior *Warrior) registerIntercept() {
	actionID := core.ActionID{SpellID: interceptRank.ID}
	chargeMinRange := float64(interceptRank.MinRange)

	var spell *core.Spell
	var interceptTarget *core.Unit

	config := spelldata.SpellConfig(&warrior.Unit, interceptRank, spelldata.Flags(core.SpellFlagAPL))
	config.ClassSpellMask = SpellMaskIntercept

	aura := warrior.RegisterAura(core.Aura{
		Label:    "Intercept",
		ActionID: actionID,
		Duration: config.Cast.CD.Duration,
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

	config.ExtraCastCondition = func(sim *core.Simulation, target *core.Unit) bool {
		return warrior.StanceMatches(BerserkerStance)
	}

	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		interceptTarget = target
		aura.Duration = spell.CD.Duration
		aura.Activate(sim)
		warrior.MoveTo(chargeMinRange-3.5, sim) // movement aura is discretized in 1 yard intervals, so need to overshoot to guarantee melee range
	}

	spell = warrior.RegisterSpell(config)
}
