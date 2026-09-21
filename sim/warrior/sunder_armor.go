package warrior

import (
	"github.com/wowsims/forever/sim/core"
)

// The client supplies Sunder Armor's flat threat per rank: 405/608/810/1013 for ranks 2-5.
//
// TODO: rank 1 reads a flat threat of 1, which looks like placeholder data next to the
// rest of the ladder. Harmless while this pins the highest rank, but worth confirming.
var sunderArmorRank = spellData.SunderArmor.HighestRank()

func (warrior *Warrior) registerSunderArmor() {
	actionId := core.ActionID{SpellID: sunderArmorRank.SpellID}

	warrior.SunderArmorAuras = warrior.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return core.SunderArmorAura(target)
	})

	warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionId,
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskSunderArmor,
		MaxRange:       core.MaxMeleeRange,

		RageCost: core.RageCostOptions{
			Cost:   sunderArmorRank.Cost,
			Refund: sunderArmorRank.MissRefund(),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: sunderArmorRank.GCD,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.CanApplySunderAura(target)
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		FlatThreatBonus:  sunderArmorRank.FlatThreatBonus,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHit)

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
	})
}

func (warrior *Warrior) CanApplySunderAura(target *core.Unit) bool {
	return warrior.SunderArmorAuras.Get(target).IsActive() || !warrior.SunderArmorAuras.Get(target).ExclusiveEffects[0].Category.AnyActive()
}
