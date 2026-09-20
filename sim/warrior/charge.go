package warrior

import (
	"time"

	"github.com/wowsims/forever/sim/core"
)

var chargeRank = spellData.Charge.BySpellID(11578)

func (warrior *Warrior) registerCharge() {
	// Charge (11578) energizes 15 rage; Improved Charge adds its ladder.
	chargeRage := 15 + spellData.ImprovedCharge.TenthsAt(warrior.Talents.ImprovedCharge)
	actionID := core.ActionID{SpellID: chargeRank.SpellID}
	metrics := warrior.NewRageMetrics(actionID)

	chargeMinRange := chargeRank.MinRange

	// TODO: Manual review needed -- spell 11578 states no duration; the 15 seconds covers the
	// sim's movement to the target and is hand-supplied.
	aura := warrior.RegisterAura(core.Aura{
		Label:    "Charge",
		ActionID: actionID,
		Duration: 15 * time.Second,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
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

	warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskCharge,
		MinRange:       chargeMinRange,
		MaxRange:       chargeRank.MaxRange,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: chargeRank.Cooldown,
			},
			IgnoreHaste: true,
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			// Vanguard (1310317) replaces Charge with a copy usable in Defensive Stance.
			return sim.CurrentTime < 0 && (warrior.StanceMatches(BattleStance) || (warrior.Talents.Vanguard && warrior.StanceMatches(DefensiveStance)))
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			aura.Activate(sim)
			warrior.AddRage(sim, chargeRage, metrics)
			warrior.MoveTo(chargeMinRange-3.5, sim) // movement aura is discretized in 1 yard intervals, so need to overshoot to guarantee melee range
		},
	})
}
