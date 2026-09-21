package warrior

import (
	"github.com/wowsims/forever/sim/core"
)

var intimidatingShoutRank = spellData.IntimidatingShout.Highest()

func (warrior *Warrior) registerIntimidatingShout() {
	warrior.IntimidatingShout = warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: intimidatingShoutRank.ID},
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskIntimidatingShout,
		MaxRange:       float64(intimidatingShoutRank.MaxRange),

		RageCost: core.RageCostOptions{
			Cost: rageCost(intimidatingShoutRank),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: intimidatingShoutRank.GCD(),
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: cooldownOf(intimidatingShoutRank),
			},
		},

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
		},
	})
}
