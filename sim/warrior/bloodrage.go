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
	healthCost := warrior.GetBaseStats()[stats.Health] * 0.16
	// Improved Bloodrage (12301) raises every rage amount Bloodrage generates by 25% per rank.
	improvedBloodrage := spellData.ImprovedBloodrage.MultiplierAt(warrior.Talents.ImprovedBloodrage)
	// Bloodrage (2687) energizes 100 on the cast, which is 10 rage; its other effect triggers the
	// rage over time, which is a spell of its own with no table.
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
