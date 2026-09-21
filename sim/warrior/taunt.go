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

		// Spell 355's ShapeshiftMask is Defensive Stance only.
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.StanceMatches(DefensiveStance)
		},

		// TODO: taunt sets the caster's threat to the highest on the target, which the sim has no
		// threat table to do; the cast is modelled and the threat is not.
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
		},
	})
}
