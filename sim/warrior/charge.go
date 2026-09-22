package warrior

import (
	"time"

	"github.com/wowsims/forever/sim/core"
)

func (warrior *Warrior) registerCharge() {
	chargeRank := spellData.Charge.BySpellID(11578)

	actionID := core.ActionID{SpellID: chargeRank.SpellID}
	metrics := warrior.NewRageMetrics(actionID)

	chargeCD := chargeRank.Cooldown
	chargeRage := chargeRank.Energize.Tenths() + spellData.ImprovedCharge.TenthsAt(warrior.Talents.ImprovedCharge)

	aura := warrior.registerDashAura("Charge", actionID, chargeCD, nil)

	warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskCharge,
		MinRange:       chargeRank.MinRange,
		MaxRange:       chargeRank.MaxRange,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: chargeCD,
			},
			IgnoreHaste: true,
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return sim.CurrentTime < 0 && (warrior.StanceMatches(BattleStance) || (warrior.Talents.Vanguard && warrior.StanceMatches(DefensiveStance)))
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			aura.Duration = spell.CD.Duration
			aura.Activate(sim)
			warrior.AddRage(sim, chargeRage, metrics)
			warrior.MoveTo(spell.MinRange-3.5, sim) // movement aura is discretized in 1 yard intervals, so need to overshoot to guarantee melee range
		},
	})
}

// TODO: Manual review needed -- the run speed and the callers' overshoot are the sim's movement model.
func (warrior *Warrior) registerDashAura(label string, actionID core.ActionID, duration time.Duration, onEnd func(sim *core.Simulation)) *core.Aura {
	aura := warrior.RegisterAura(core.Aura{
		Label:    label,
		ActionID: actionID,
		Duration: duration,
		OnGain: func(_ *core.Aura, sim *core.Simulation) {
			warrior.MultiplyMovementSpeed(sim, 3.0)
		},
		OnExpire: func(_ *core.Aura, sim *core.Simulation) {
			warrior.MultiplyMovementSpeed(sim, 1.0/3.0)
			if onEnd != nil {
				onEnd(sim)
			}
		},
	})

	warrior.RegisterMovementCallback(func(sim *core.Simulation, _ float64, kind core.MovementUpdateType) {
		if kind == core.MovementEnd && aura.IsActive() {
			aura.Deactivate(sim)
		}
	})

	return aura
}
