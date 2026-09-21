package warrior

import (
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/stats"
)

var bloodrageRank = spellData.Bloodrage.HighestRank()

func (warrior *Warrior) registerBloodrage() {
	actionID := core.ActionID{SpellID: bloodrageRank.SpellID}
	rageMetrics := warrior.NewRageMetrics(actionID)
	// TODO: Manual review needed -- SpellPower states PowerCostPct 20 on 2687, which the generator does not carry yet.
	healthCost := warrior.GetBaseStats()[stats.Health] * 0.20
	improvedBloodrage := spellData.ImprovedBloodrage.MultiplierAt(warrior.Talents.ImprovedBloodrage)
	instantRage := spellData.Bloodrage.EffectAt(0).TenthsAt(1) * improvedBloodrage

	spell := warrior.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: bloodrageRank.Cooldown,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			warrior.AddRage(sim, instantRage, rageMetrics)
			warrior.RemoveHealth(sim, healthCost)

			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				// TODO: Manual review needed -- the tooltip states 10 rage over 10 s; the 1 s period is not in the client.
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
