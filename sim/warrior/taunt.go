package warrior

import (
	"time"

	"github.com/wowsims/forever/sim/core"
)

// TODO: Manual review needed -- spell 355 states no rage cost, an 8 second cooldown, a 3 second
// taunt and no global cooldown.
const tauntCooldown = time.Second * 8

func (warrior *Warrior) registerTaunt() {
	warrior.Taunt = warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 355},
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskTaunt,
		MaxRange:       core.MaxMeleeRange,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: tauntCooldown,
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
