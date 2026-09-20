package rogue

import (
	"github.com/wowsims/forever/sim/core"
)

var exposeArmorRank = spellData.ExposeArmor.HighestRank()

func (rogue *Rogue) registerExposeArmorSpell() {
	rogue.ExposeArmorAuras = rogue.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return core.ExposeArmorAura(target, rogue.ComboPoints, rogue.Talents.ImprovedExposeArmor)
	})

	rogue.ExposeArmor = rogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: exposeArmorRank.SpellID},
		SpellSchool:    exposeArmorRank.SpellSchool,
		DefenseType:    exposeArmorRank.DefenseType,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		MetricSplits:   6,
		ClassSpellMask: RogueSpellExposeArmor,

		EnergyCost: core.EnergyCostOptions{
			Cost: exposeArmorRank.Cost,
			// TODO: Forever drops Quick Recovery; no energy refund until we know whether the
			// effect moved onto another talent.
			Refund:        0,
			RefundMetrics: rogue.EnergyRefundMetrics,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: exposeArmorRank.GCD,
			},
			IgnoreHaste: true,
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				spell.SetMetricsSplit(rogue.ComboPoints())
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return rogue.ComboPoints() > 0
		},

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if rogue.CanApplyExposeArmorAura(target) {
				rogue.BreakStealth(sim)
				result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
				if result.Landed() {
					rogue.ExposeArmorAuras.Get(target).Activate(sim)
					rogue.ApplyFinisher(sim, spell)
				} else {
					spell.IssueRefund(sim)
				}
				spell.DealOutcome(sim, result)
			}
		},

		RelatedAuraArrays: rogue.ExposeArmorAuras.ToMap(),
	})
}

func (rogue *Rogue) GetExposeArmorValue() float64 {
	return 410.0 * float64(rogue.ComboPoints()) * spellData.ImprovedExposeArmor.MultiplierAt(rogue.Talents.ImprovedExposeArmor)
}

func (rogue *Rogue) CanApplyExposeArmorAura(target *core.Unit) bool {
	return !rogue.ExposeArmorAuras.Get(target).IsActive() || rogue.ExposeArmorAuras.Get(target).ExclusiveEffects[0].Priority <= rogue.GetExposeArmorValue()
}
