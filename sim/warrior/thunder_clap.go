package warrior

import (
	"github.com/wowsims/forever/sim/core"
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

	warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: thunderClapRank.ID},
		SpellSchool: thunderClapRank.SpellSchool(),
		// Thunder Clap is Physical but Magic in SpellCategories: it rolls on the spell hit table
		// (logs show full resists next to armor mitigation) and crits on spell crit chance for
		// 1.5x. Warriors have no base spell crit, so logs without Totem of Wrath show none
		// (0 of 799 landed hits from 6 prot warriors on fresh.warcraftlogs.com, 2026-09-14).
		DefenseType:    thunderClapRank.DefenseTypeCore(),
		ProcMask:       core.ProcMaskRangedSpecial,
		Flags:          core.SpellFlagAPL | core.SpellFlagBinary,
		ClassSpellMask: SpellMaskThunderClap,
		ClassFlags:     SpellFlagsThunderClap,

		RageCost: core.RageCostOptions{
			Cost: rageCost(thunderClapRank),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: thunderClapRank.GCD(),
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: cooldownOf(thunderClapRank),
			},
		},

		DamageMultiplier: 1,
		// TODO: Manual review needed -- the client states no threat coefficient; 1 until measured in game.
		ThreatMultiplier: 1,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			// Thunder Clap (11581) is usable in Battle and Defensive Stance.
			return warrior.StanceMatches(BattleStance | DefensiveStance)
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			results := spell.CalcCleaveDamage(sim, target, int32(thunderClapRank.MaxTargets), thunderClapBaseDamage, spell.OutcomeMagicHitAndCrit)
			warrior.CastNormalizedSweepingStrikesAttack(results, sim)

			for _, result := range results {
				if result.Landed() {
					auras.Get(result.Target).Activate(sim)
				}
				spell.DealDamage(sim, result)
			}
		},

		RelatedAuraArrays: auras.ToMap(),
	})
}
