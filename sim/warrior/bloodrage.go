package warrior

import (
	"time"

	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/stats"
)

func (warrior *Warrior) registerBloodrage() {
	bloodrageRank := spellData.Bloodrage.HighestRank()
	bloodrageTriggered := spellData.BloodrageTriggered.HighestRank()

	actionID := core.ActionID{SpellID: bloodrageRank.SpellID}
	rageMetrics := warrior.NewRageMetrics(actionID)
	// The client costs 20% of base health (SpellPower.PowerCostPct); the generator on forever-next
	// does not carry the column yet.
	healthCost := warrior.GetBaseStats()[stats.Health] * 20 / 100
	improvedBloodrage := spellData.ImprovedBloodrage.MultiplierAt(warrior.Talents.ImprovedBloodrage)
	instantRage := spellData.Bloodrage.EffectAt(0).TenthsAt(1) * improvedBloodrage
	// 29131's periodic energize: 1 rage a second for its 10 sec. The generator files it as a flat
	// Energize of 100 rather than a periodic, so the tick is read off the effect and the schedule off
	// the duration.
	ragePerTick := bloodrageTriggered.Effect(shared.A_PERIODIC_ENERGIZE, 1).Tenths() * improvedBloodrage
	tickLength := time.Second
	numTicks := int(bloodrageTriggered.Duration / tickLength)

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
				NumTicks: numTicks,
				Period:   tickLength,
				OnAction: func(sim *core.Simulation) {
					warrior.AddRage(sim, ragePerTick, rageMetrics)
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
