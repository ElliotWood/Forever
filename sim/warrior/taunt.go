package warrior

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/spelldata"
)

var tauntRank = spellData.Taunt.Highest()

func (warrior *Warrior) registerTaunt() {
	config := spelldata.SpellConfig(&warrior.Unit, tauntRank, spelldata.Flags(core.SpellFlagAPL))
	config.ClassSpellMask = SpellMaskTaunt
	config.ProcMask = core.ProcMaskEmpty
	config.ThreatMultiplier = 1
	config.Cast.DefaultCast.NonEmpty = true

	// Spell 355's ShapeshiftMask is Defensive Stance only.
	config.ExtraCastCondition = func(sim *core.Simulation, target *core.Unit) bool {
		return warrior.StanceMatches(DefensiveStance)
	}

	// TODO: taunt sets the caster's threat to the highest on the target, which the sim has no
	// threat table to do; the cast is modelled and the threat is not.
	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
	}

	warrior.Taunt = warrior.RegisterSpell(config)
}
