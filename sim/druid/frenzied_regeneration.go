package druid

import (
	"time"

	"github.com/wowsims/classic/sim/core"
)

const FrenziedRegenerationRanks = 3

var FrenziedRegenerationSpellId = [FrenziedRegenerationRanks + 1]int32{0, 22842, 22895, 22896}
var FrenziedRegenerationHealthPerRage = [FrenziedRegenerationRanks + 1]float64{0, 10, 15, 20}
var FrenziedRegenerationLevel = [FrenziedRegenerationRanks + 1]int{0, 36, 46, 56}

// Converts up to 10 Rage per second into health for 10 sec.
func (druid *Druid) registerFrenziedRegenerationCD() {
	rank := map[int32]int{
		40: 1,
		50: 2,
		60: 3,
	}[druid.Level]
	if rank == 0 {
		return
	}

	// Beta client 1.60.1.69893: one rank (22842), and each point of Rage heals 1% of maximum health instead of a flat
	// 10 / 15 / 20.
	forever := druid.Env.IsForever()
	if forever {
		rank = 1
	}

	actionID := core.ActionID{SpellID: FrenziedRegenerationSpellId[rank]}
	healthPerRage := FrenziedRegenerationHealthPerRage[rank]
	healthMetrics := druid.NewHealthMetrics(actionID)
	rageMetrics := druid.NewRageMetrics(actionID)

	druid.FrenziedRegenerationAura = druid.RegisterAura(core.Aura{
		Label:    "Frenzied Regeneration",
		ActionID: actionID,
		Duration: time.Second * 10,
	})

	druid.FrenziedRegeneration = druid.RegisterSpell(Bear, core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagHelpful | core.SpellFlagAPL,

		Rank:          rank,
		RequiredLevel: FrenziedRegenerationLevel[rank],

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Minute * 3,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				NumTicks: 10,
				Period:   time.Second,
				OnAction: func(sim *core.Simulation) {
					if !druid.FrenziedRegenerationAura.IsActive() {
						return
					}

					rageDumped := min(druid.CurrentRage(), 10.0)
					druid.SpendRage(sim, rageDumped, rageMetrics)
					if forever {
						healthPerRage = 0.01 * druid.MaxHealth()
					}
					druid.GainHealth(sim, rageDumped*healthPerRage*druid.PseudoStats.HealingTakenMultiplier, healthMetrics)
				},
			})

			druid.FrenziedRegenerationAura.Activate(sim)
		},
	})

	druid.AddMajorCooldown(core.MajorCooldown{
		Spell: druid.FrenziedRegeneration.Spell,
		Type:  core.CooldownTypeSurvival,
	})
}
