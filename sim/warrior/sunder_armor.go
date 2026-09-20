package warrior

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

// Forever's client supplies Sunder Armor's flat threat per rank (405/608/810/1013 for
// ranks 2-5), so the hand-applied 301.5 is gone -- WithSpellDataFlatThreat panics rather
// than silently overriding client data.
//
// TODO: rank 1 reads a flat threat of 1, which looks like placeholder data next to the
// rest of the ladder. Harmless while this pins the highest rank, but worth confirming.
var sunderArmorRank = spellData.SunderArmor.HighestRank()

func (warrior *Warrior) registerSunderArmor() {
	actionId := core.ActionID{SpellID: sunderArmorRank.SpellID}

	warrior.SunderArmorAuras = warrior.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return core.SunderArmorAura(target)
	})

	getSunderArmorConfig := func(config core.SpellConfig, outcome shared.OutcomeType) core.SpellConfig {
		return core.SpellConfig{
			ActionID:       config.ActionID,
			SpellSchool:    core.SpellSchoolPhysical,
			DefenseType:    core.DefenseTypeMelee,
			ProcMask:       core.ProcMaskMeleeMHSpecial,
			Flags:          config.Flags,
			ClassSpellMask: SpellMaskSunderArmor,
			MaxRange:       core.MaxMeleeRange,

			RageCost: config.RageCost,
			Cast:     config.Cast,
			ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
				return warrior.CanApplySunderAura(target)
			},

			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			FlatThreatBonus:  sunderArmorRank.FlatThreatBonus,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				result := spell.CalcOutcome(sim, target, shared.GetOutcome(spell, outcome))

				if result.Landed() {
					aura := warrior.SunderArmorAuras.Get(target)
					aura.Activate(sim)
					aura.AddStack(sim)
				} else if spell.Cost != nil {
					spell.IssueRefund(sim)
				}

				spell.DealOutcome(sim, result)
			},

			RelatedAuraArrays: warrior.SunderArmorAuras.ToMap(),
		}
	}

	warrior.RegisterSpell(getSunderArmorConfig(core.SpellConfig{
		ActionID: actionId,
		Flags:    core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		RageCost: core.RageCostOptions{
			Cost:   sunderArmorRank.Cost,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: sunderArmorRank.GCD,
			},
			IgnoreHaste: true,
		},
	}, shared.OutcomeMeleeNoCrit))

}

func (warrior *Warrior) CanApplySunderAura(target *core.Unit) bool {
	return warrior.SunderArmorAuras.Get(target).IsActive() || !warrior.SunderArmorAuras.Get(target).ExclusiveEffects[0].Category.AnyActive()
}
