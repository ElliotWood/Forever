package rogue

import (
	"github.com/wowsims/forever/sim/core"
)

var eviscerateRank = spellData.Eviscerate.BySpellID(26865)

func (rogue *Rogue) registerEviscerate() {
	flatDamage, flatDamageMax := eviscerateRank.Direct.Range()
	comboDamageBonus := 185.0 + rogue.DeathmantleBonus
	damageVariance := flatDamageMax - flatDamage

	rogue.Eviscerate = rogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: eviscerateRank.SpellID},
		SpellSchool:    eviscerateRank.SpellSchool,
		DefenseType:    eviscerateRank.DefenseType,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | SpellFlagFinisher | core.SpellFlagAPL,
		MetricSplits:   6,
		ClassSpellMask: RogueSpellEviscerate,

		EnergyCost: core.EnergyCostOptions{
			Cost: eviscerateRank.Cost,
			// TODO: Forever drops Quick Recovery; no energy refund until we know whether the
			// effect moved onto another talent.
			Refund:        0,
			RefundMetrics: rogue.EnergyRefundMetrics,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: eviscerateRank.GCD,
			},
			IgnoreHaste: true,
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				spell.SetMetricsSplit(rogue.ComboPoints())
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return rogue.ComboPoints() > 0
		},

		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		ThreatMultiplier:         1,

		BonusCoefficient: eviscerateRank.Direct.BonusCoefficient(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)

			comboPoints := float64(rogue.ComboPoints())
			flatBaseDamage := flatDamage + comboDamageBonus*float64(comboPoints)

			baseDamage := sim.Roll(flatBaseDamage, flatBaseDamage+damageVariance) + 0.03*float64(comboPoints)*spell.MeleeAttackPower(target)

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				rogue.ApplyFinisher(sim, spell)
			} else {
				spell.IssueRefund(sim)
			}

			spell.DealDamage(sim, result)
		},
	})
}
