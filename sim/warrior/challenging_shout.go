package warrior

import (
	"github.com/wowsims/forever/sim/core"
)

var challengingShoutRank = spellData.ChallengingShout.HighestRank()

func (warrior *Warrior) registerChallengingShout() {
	warrior.ChallengingShout = warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: challengingShoutRank.SpellID},
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskChallengingShout,

		RageCost: core.RageCostOptions{
			Cost: challengingShoutRank.Cost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: challengingShoutRank.GCD,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: challengingShoutRank.Cooldown,
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
