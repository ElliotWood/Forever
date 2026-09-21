package warrior

import (
	"github.com/wowsims/forever/sim/core"
)

var thunderClapRank = spellData.ThunderClap.HighestRank()

var thunderClapBaseDamage, _ = thunderClapRank.Direct.Range()

func (warrior *Warrior) registerThunderClap() {
	auras := warrior.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return core.ThunderClapAura(target, warrior.Talents.ImprovedThunderClap)
	})

	warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: thunderClapRank.SpellID},
		SpellSchool: thunderClapRank.SpellSchool,
		// Thunder Clap is Physical but Magic in SpellCategories: it rolls on the spell hit table
		// (logs show full resists next to armor mitigation) and crits on spell crit chance for
		// 1.5x. Warriors have no base spell crit, so logs without Totem of Wrath show none
		// (0 of 799 landed hits from 6 prot warriors on fresh.warcraftlogs.com, 2026-09-14).
		DefenseType:    thunderClapRank.DefenseType,
		ProcMask:       core.ProcMaskRangedSpecial,
		Flags:          core.SpellFlagAPL | core.SpellFlagBinary,
		ClassSpellMask: SpellMaskThunderClap,

		RageCost: core.RageCostOptions{
			Cost: thunderClapRank.Cost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: thunderClapRank.GCD,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: thunderClapRank.Cooldown,
			},
		},

		DamageMultiplier: 1,
		// TODO: Manual review needed -- the threat coefficient is not in the client.
		ThreatMultiplier: 1.75,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			// Thunder Clap (11581) is usable in Battle and Defensive Stance.
			return warrior.StanceMatches(BattleStance | DefensiveStance)
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			results := spell.CalcCleaveDamage(sim, target, thunderClapRank.MaxTargets, thunderClapBaseDamage, spell.OutcomeMagicHitAndCrit)
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
