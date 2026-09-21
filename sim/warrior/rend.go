package warrior

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/spelldata"
)

var rendRank = spellData.Rend.Highest()

// TODO: Ingame testing needed if Rend has a coef
func (warrior *Warrior) registerRend() {
	tick := rendRank.PeriodicEffect()

	config := spelldata.SpellConfig(&warrior.Unit, rendRank,
		spelldata.Flags(core.SpellFlagNoOnCastComplete|core.SpellFlagAPL))
	config.ClassSpellMask = SpellMaskRend
	config.ProcMask = core.ProcMaskMeleeMHSpecial
	config.DamageMultiplier = 1
	config.ThreatMultiplier = 1

	config.Dot = spelldata.DotConfig(rendRank, tick)
	// The resolver's callbacks snapshot the multipliers when the dot is applied; these ticks take
	// the multipliers in force at the tick. The tick count, the period and BonusCoefficient stay
	// the row's, and this effect states no spell power coefficient - one that did would put spell
	// power on every tick.
	config.Dot.OnSnapshot = nil
	config.Dot.OnTick = func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
		dot.Spell.CalcAndDealPeriodicDamage(sim, target, tick.Average(core.CharacterLevel), rendRank.TickOutcome(dot))
	}

	config.ExtraCastCondition = func(sim *core.Simulation, target *core.Unit) bool {
		return warrior.StanceMatches(BattleStance | DefensiveStance)
	}

	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
		if result.Landed() {
			spell.Dot(target).Apply(sim)
		} else {
			spell.IssueRefund(sim)
		}
	}

	warrior.Rend = warrior.RegisterSpell(config)
}
