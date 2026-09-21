package warrior

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/spelldata"
)

// The client supplies Sunder Armor's flat threat per rank: 405/608/810/1013 for ranks 2-5.
//
// TODO: rank 1 reads a flat threat of 1, which looks like placeholder data next to the
// rest of the ladder. Harmless while this pins the highest rank, but worth confirming.
var sunderArmorRank = spellData.SunderArmor.Highest()

func (warrior *Warrior) registerSunderArmor() {
	warrior.SunderArmorAuras = warrior.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return core.SunderArmorAura(target)
	})

	config := spelldata.SpellConfig(&warrior.Unit, sunderArmorRank, spelldata.Melee(core.ProcMaskMeleeMHSpecial))
	config.FlatThreatBonus = sunderArmorRank.FindEffect(dbcenums.E_THREAT, 0, 0).Average(core.CharacterLevel)

	config.ExtraCastCondition = func(sim *core.Simulation, target *core.Unit) bool {
		return warrior.CanApplySunderAura(target)
	}

	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHit)

		if result.Landed() {
			aura := warrior.SunderArmorAuras.Get(target)
			aura.Activate(sim)
			aura.AddStack(sim)
		} else if spell.Cost != nil {
			spell.IssueRefund(sim)
		}

		spell.DealOutcome(sim, result)
	}

	config.RelatedAuraArrays = warrior.SunderArmorAuras.ToMap()

	warrior.RegisterSpell(config)
}

func (warrior *Warrior) CanApplySunderAura(target *core.Unit) bool {
	return warrior.SunderArmorAuras.Get(target).IsActive() || !warrior.SunderArmorAuras.Get(target).ExclusiveEffects[0].Category.AnyActive()
}
