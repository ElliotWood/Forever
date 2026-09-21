package warrior

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/spelldata"
)

var victoryRushRank = spellData.VictoryRush.Highest()

// Spell 402927 states ${1+$AP*$m3/100}: the dummy at effect index 2 is the attack power
// coefficient as a percentage, and the heal at index 1 is a percentage of maximum health.
var victoryRushAPCoef = victoryRushRank.EffectN(3).Percent()
var victoryRushHealPercent = victoryRushRank.EffectN(2).Percent()

// The window Victory Rush has to be used in is the Victorious buff the kill grants.
var victoriousRank = spellData.VictoryRushTriggered.Highest()

func (warrior *Warrior) registerVictoryRush() {
	healthMetrics := warrior.NewHealthMetrics(core.ActionID{SpellID: victoryRushRank.ID})

	// TODO: spell 402974 grants this on a kill, which the sim never simulates, so nothing
	// activates it and Victory Rush stays uncastable.
	victoriousAura := warrior.RegisterAura(spelldata.AuraConfig(victoriousRank))

	config := spelldata.SpellConfig(&warrior.Unit, victoryRushRank, spelldata.Melee(core.ProcMaskMeleeMHSpecial))
	config.ClassSpellMask = SpellMaskVictoryRush

	config.ExtraCastCondition = func(sim *core.Simulation, target *core.Unit) bool {
		return victoriousAura.IsActive()
	}

	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		baseDamage := victoryRushRank.DamageEffect().Average(core.CharacterLevel) + victoryRushAPCoef*spell.MeleeAttackPower(target)
		spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
		warrior.GainHealth(sim, warrior.MaxHealth()*victoryRushHealPercent, healthMetrics)
		victoriousAura.Deactivate(sim)
	}

	config.RelatedSelfBuff = victoriousAura

	warrior.VictoryRush = warrior.RegisterSpell(config)
}
