package warrior

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/spelldata"
)

var chargeRank = spellData.Charge.ByID(11578)

func (warrior *Warrior) registerCharge() {
	actionID := core.ActionID{SpellID: chargeRank.ID}
	metrics := warrior.NewRageMetrics(actionID)

	chargeRage := chargeRank.EnergizeEffect().Tenths()
	if warrior.Talents.ImprovedCharge > 0 {
		chargeRage += spellData.ImprovedCharge.TenthsAt(warrior.Talents.ImprovedCharge)
	}

	config := spelldata.SpellConfig(&warrior.Unit, chargeRank, spelldata.Flags(core.SpellFlagAPL))

	aura := warrior.RegisterAura(core.Aura{
		Label:    "Charge",
		ActionID: actionID,
		Duration: config.Cast.CD.Duration,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			// TODO: Manual review needed -- the run speed and the overshoot below are the sim's movement model.
			warrior.MultiplyMovementSpeed(sim, 3.0)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.MultiplyMovementSpeed(sim, 1.0/3.0)
		},
	})

	warrior.RegisterMovementCallback(func(sim *core.Simulation, position float64, kind core.MovementUpdateType) {
		if kind == core.MovementEnd && aura.IsActive() {
			aura.Deactivate(sim)
		}
	})

	config.ExtraCastCondition = func(sim *core.Simulation, target *core.Unit) bool {
		return sim.CurrentTime < 0 && (warrior.StanceMatches(BattleStance) || (warrior.Talents.Vanguard && warrior.StanceMatches(DefensiveStance)))
	}

	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		aura.Duration = spell.CD.Duration
		aura.Activate(sim)
		warrior.AddRage(sim, chargeRage, metrics)
		warrior.MoveTo(spell.MinRange-3.5, sim) // movement aura is discretized in 1 yard intervals, so need to overshoot to guarantee melee range
	}

	warrior.RegisterSpell(config)
}
