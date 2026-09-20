package warrior

import (
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/stats"
)

func (warrior *Warrior) registerBloodrage() {
	actionID := core.ActionID{SpellID: 2687}
	rageMetrics := warrior.NewRageMetrics(actionID)
	healthCost := warrior.GetBaseStats()[stats.Health] * 0.16
	// Improved Bloodrage (12301) raises every rage amount Bloodrage generates by 25% per rank.
	improvedBloodrage := spellData.ImprovedBloodrage.MultiplierAt(warrior.Talents.ImprovedBloodrage)
	instantRage := 10.0 * improvedBloodrage

	spell := warrior.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Minute,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			warrior.AddRage(sim, instantRage, rageMetrics)
			warrior.RemoveHealth(sim, healthCost)

			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				NumTicks: 10,
				Period:   time.Second * 1,
				OnAction: func(sim *core.Simulation) {
					warrior.AddRage(sim, improvedBloodrage, rageMetrics)
				},
			})
		},
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			return warrior.CurrentRage() < 70
		},
	})
}
