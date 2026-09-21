package warrior

import (
	"github.com/wowsims/forever/sim/core"
)

var chargeRank = spellData.Charge.ByID(11578)

func (warrior *Warrior) registerCharge() {
	actionID := core.ActionID{SpellID: chargeRank.ID}
	metrics := warrior.NewRageMetrics(actionID)

	chargeCD := cooldownOf(chargeRank)
	chargeRage := chargeRank.EnergizeEffect().Tenths()
	if warrior.Talents.ImprovedCharge > 0 {
		chargeRage += spellData.ImprovedCharge.TenthsAt(warrior.Talents.ImprovedCharge)
	}

	aura := warrior.RegisterAura(core.Aura{
		Label:    "Charge",
		ActionID: actionID,
		Duration: chargeCD,
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

	warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskCharge,
		MinRange:       float64(chargeRank.MinRange),
		MaxRange:       float64(chargeRank.MaxRange),

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
