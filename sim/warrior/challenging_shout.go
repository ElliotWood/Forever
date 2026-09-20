package warrior

import (
	"time"

	"github.com/wowsims/forever/sim/core"
)

// TODO: Manual review needed -- spell 1161 states a 5 rage cost, a 10 minute cooldown and a
// 6 second taunt, and states no threat amount.
const (
	challengingShoutRageCost int32 = 5
	challengingShoutCooldown       = time.Minute * 10
)

func (warrior *Warrior) registerChallengingShout() {
	warrior.ChallengingShout = warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 1161},
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskChallengingShout,

		RageCost: core.RageCostOptions{
			Cost: challengingShoutRageCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: challengingShoutCooldown,
			},
		},

		ThreatMultiplier: 1,

		// TODO: the taunt is not modelled; the sim has no threat table to force the targets onto.
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.ActiveTargetUnits {
				spell.CalcAndDealOutcome(sim, aoeTarget, spell.OutcomeAlwaysHit)
			}
		},
	})
}
