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

func (war *Warrior) registerSunderArmor() {
	actionId := core.ActionID{SpellID: sunderArmorRank.SpellID}

	war.SunderArmorAuras = war.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return core.SunderArmorAura(target, true, 0)
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
				return war.CanApplySunderAura(target)
			},

			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			FlatThreatBonus:  sunderArmorRank.FlatThreatBonus,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				result := spell.CalcOutcome(sim, target, shared.GetOutcome(spell, outcome))

				if result.Landed() {
					aura := war.SunderArmorAuras.Get(target)
					aura.Activate(sim)
					aura.AddStack(sim)
				} else if spell.Cost != nil {
					spell.IssueRefund(sim)
				}

				spell.DealOutcome(sim, result)
			},

			RelatedAuraArrays: war.SunderArmorAuras.ToMap(),
		}
	}

	war.RegisterSpell(getSunderArmorConfig(core.SpellConfig{
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

	war.SunderArmorDevastate = war.RegisterSpell(getSunderArmorConfig(core.SpellConfig{
		ActionID: actionId.WithTag(1),
	}, shared.OutcomeAlwaysHit))
}

func (warrior *Warrior) CanApplySunderAura(target *core.Unit) bool {
	return warrior.SunderArmorAuras.Get(target).IsActive() || !warrior.SunderArmorAuras.Get(target).ExclusiveEffects[0].Category.AnyActive()
}
