package druid

import (
	"time"

	"github.com/wowsims/classic/sim/core"
)

// Berserk widens Mangle to a cleave.
const MangleBerserkTargets = 3

func (druid *Druid) registerMangleCatSpell() {
	if !druid.Talents.Mangle {
		return
	}

	// TODO: Only the tooltip was seen, the Energy cost is taken from the Classic Mangle (Cat).
	flatDamageBonus := 26.0
	results := make([]*core.SpellResult, min(MangleBerserkTargets, druid.Env.GetNumTargets()))

	druid.MangleCat = druid.RegisterSpell(Cat, core.SpellConfig{
		SpellCode:   SpellCode_DruidMangle,
		ActionID:    core.ActionID{SpellID: 33876},
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL | SpellFlagBuilder,

		EnergyCost: core.EnergyCostOptions{
			Cost:   45 - float64(druid.Talents.Ferocity),
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			numHits := 1
			if druid.BerserkAura.IsActive() {
				numHits = len(results)
			}

			for idx := 0; idx < numHits; idx++ {
				baseDamage := flatDamageBonus + spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower(target))
				results[idx] = spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
				target = sim.Environment.NextTargetUnit(target)
			}

			for idx := 0; idx < numHits; idx++ {
				spell.DealDamage(sim, results[idx])
			}

			if results[0].Landed() {
				druid.AddComboPoints(sim, 1, results[0].Target, spell.ComboPointMetrics())
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}
