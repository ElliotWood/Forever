package warrior

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/stats"
)

var bloodrageRank = spellData.Bloodrage.Highest()
var bloodrageOverTime = spellData.BloodrageTriggered.Highest()
var bloodrageOverTimeTick = bloodrageOverTime.PeriodicEffect()

func (warrior *Warrior) registerBloodrage() {
	actionID := core.ActionID{SpellID: bloodrageRank.ID}
	rageMetrics := warrior.NewRageMetrics(actionID)
	healthCost := warrior.GetBaseStats()[stats.Health] * float64(bloodrageRank.Powers[0].CostPct) / 100
	improvedBloodrage := spellData.ImprovedBloodrage.MultiplierAt(warrior.Talents.ImprovedBloodrage)
	instantRage := spellData.Bloodrage.EffectAt(1).TenthsAt(1) * improvedBloodrage

	spell := warrior.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: cooldownOf(bloodrageRank),
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			warrior.AddRage(sim, instantRage, rageMetrics)
			warrior.RemoveHealth(sim, healthCost)

			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				NumTicks: int(bloodrageOverTime.Duration() / bloodrageOverTimeTick.Period()),
				Period:   bloodrageOverTimeTick.Period(),
				OnAction: func(sim *core.Simulation) {
					warrior.AddRage(sim, bloodrageOverTimeTick.Tenths()*improvedBloodrage, rageMetrics)
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
