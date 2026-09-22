package warrior

import (
	"github.com/wowsims/forever/sim/core"
)

func (warrior *Warrior) registerIntimidatingShout() {
	intimidatingShoutRank := spellData.IntimidatingShout.HighestRank()

	warrior.IntimidatingShout = warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: intimidatingShoutRank.SpellID},
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskIntimidatingShout,
		MaxRange:       intimidatingShoutRank.MaxRange,

		RageCost: core.RageCostOptions{
			Cost: intimidatingShoutRank.Cost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: intimidatingShoutRank.GCD,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: intimidatingShoutRank.Cooldown,
			},
		},

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
		},
	})
}
