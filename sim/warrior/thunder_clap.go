package warrior

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/spelldata"
)

var thunderClapRank = spellData.ThunderClap.Highest()

var thunderClapBaseDamage = thunderClapRank.DamageEffect().Average(core.CharacterLevel)
var thunderClapSlow = -thunderClapRank.EffectN(2).Percent()

func (warrior *Warrior) registerThunderClap() {
	auras := warrior.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return core.ThunderClapAura(target).ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
			slow := thunderClapSlow * (1 + warrior.thunderClapEffectBonus)
			aura.ExclusiveEffects[0].SetPriority(sim, 1/(1-slow))
		})
	})

	// Thunder Clap is Physical but Magic in SpellCategories: it rolls on the spell hit table
	// (logs show full resists next to armor mitigation) and crits on spell crit chance for
	// 1.5x. Warriors have no base spell crit, so logs without Totem of Wrath show none
	// (0 of 799 landed hits from 6 prot warriors on fresh.warcraftlogs.com, 2026-09-14).
	config := spelldata.SpellConfig(&warrior.Unit, thunderClapRank,
		spelldata.Flags(core.SpellFlagAPL|core.SpellFlagBinary))
	config.ClassSpellMask = SpellMaskThunderClap
	config.ProcMask = core.ProcMaskRangedSpecial
	config.DamageMultiplier = 1
	// TODO: Manual review needed -- the client states no threat coefficient; 1 until measured in game.
	config.ThreatMultiplier = 1

	config.ExtraCastCondition = func(sim *core.Simulation, target *core.Unit) bool {
		// Thunder Clap (11581) is usable in Battle and Defensive Stance.
		return warrior.StanceMatches(BattleStance | DefensiveStance)
	}

	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		results := spell.CalcCleaveDamage(sim, target, int32(thunderClapRank.MaxTargets), thunderClapBaseDamage, spell.OutcomeMagicHitAndCrit)
		warrior.CastNormalizedSweepingStrikesAttack(results, sim)

		for _, result := range results {
			if result.Landed() {
				auras.Get(result.Target).Activate(sim)
			}
			spell.DealDamage(sim, result)
		}
	}

	config.RelatedAuraArrays = auras.ToMap()

	warrior.RegisterSpell(config)
}
