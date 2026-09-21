package warrior

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/stats"
)

var bloodrageRank = spellData.Bloodrage.HighestRank()
var bloodrageOverTime = spellData.BloodrageTriggered.HighestRank().Energize.(shared.SpellDataPeriodic)

func (warrior *Warrior) registerBloodrage() {
	actionID := core.ActionID{SpellID: bloodrageRank.SpellID}
	rageMetrics := warrior.NewRageMetrics(actionID)
	healthCost := warrior.GetBaseStats()[stats.Health] * bloodrageRank.PowerCostPct / 100
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
				NumTicks: int(bloodrageOverTime.NumberOfTicks),
				Period:   bloodrageOverTime.TickLength,
				OnAction: func(sim *core.Simulation) {
					warrior.AddRage(sim, bloodrageOverTime.Tick/10*improvedBloodrage, rageMetrics)
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
