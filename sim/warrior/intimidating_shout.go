package warrior

import (
	"time"

	"github.com/wowsims/forever/sim/core"
)

// TODO: Manual review needed -- spell 5246 states a 25 rage cost, a 3 minute cooldown, an
// 8 second fear on up to 5 targets and a 10 yard range.
const (
	intimidatingShoutRageCost int32 = 25
	intimidatingShoutCooldown       = time.Minute * 3
	intimidatingShoutMaxRange       = 10.0
)

func (warrior *Warrior) registerIntimidatingShout() {
	warrior.IntimidatingShout = warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 5246},
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskIntimidatingShout,
		MaxRange:       intimidatingShoutMaxRange,

		RageCost: core.RageCostOptions{
			Cost: intimidatingShoutRageCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: intimidatingShoutCooldown,
			},
		},

		ThreatMultiplier: 1,

		// TODO: the fear is not modelled; encounter targets take no crowd control in the sim.
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
		},
	})
}
