package warrior

import (
	"github.com/wowsims/forever/sim/core"
)

func (warrior *Warrior) registerTaunt() {
	tauntRank := spellData.Taunt.HighestRank()

	warrior.Taunt = warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: tauntRank.SpellID},
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskTaunt,
		MaxRange:       tauntRank.MaxRange,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: tauntRank.Cooldown,
			},
		},

		ThreatMultiplier: 1,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.StanceMatches(DefensiveStance)
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
		},
	})
}
