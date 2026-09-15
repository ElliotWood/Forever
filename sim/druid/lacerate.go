package druid

import (
	"time"

	"github.com/wowsims/classic/sim/core"
)

const LacerateMaxStacks int32 = 5

// Forever gives the bear Lacerate: Shredding Attacks in the Forever tree cuts its Rage
// cost, and nothing else in the tree or the spellbook has been seen. The Classic era
// never had it at 60, so the shape and the numbers here are the Season of Discovery
// Lacerate, the only level 60 tuning of the ability there is: 10 Rage, 20% weapon damage
// per stack on the hit, a bleed of 29.8 per tick per stack over 15 sec, and 3.33x threat.
// TODO: assumed from Season of Discovery, beta will confirm the cost, the damage and the threat.
func (druid *Druid) registerLacerateSpell() {
	druid.registerLacerateBleedSpell()

	results := make([]*core.SpellResult, min(MangleBerserkTargets, druid.Env.GetNumTargets()))

	druid.Lacerate = druid.RegisterSpell(Bear, core.SpellConfig{
		SpellCode:   SpellCode_DruidLacerate,
		ActionID:    core.ActionID{SpellID: 414644},
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		RageCost: core.RageCostOptions{
			Cost:   10 - float64(druid.Talents.ShreddingAttacks),
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 3.33,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			numHits := 1
			if druid.BerserkAura.IsActive() {
				numHits = len(results)
			}

			for idx := 0; idx < numHits; idx++ {
				stacks := min(druid.LacerateBleed.Dot(target).GetStacks()+1, LacerateMaxStacks)
				baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower(target)) * 0.2 * float64(stacks)
				results[idx] = spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

				if results[idx].Landed() {
					druid.LacerateBleed.Cast(sim, target)
				}
				target = sim.Environment.NextTargetUnit(target)
			}

			if !results[0].Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
}

func (druid *Druid) registerLacerateBleedSpell() {
	tickDamage := 29.8312

	druid.LacerateBleed = druid.RegisterSpell(Bear, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 414647},
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete,

		DamageMultiplier: 1,
		ThreatMultiplier: 3.33,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:     "Lacerate",
				MaxStacks: LacerateMaxStacks,
				Duration:  time.Second * 15,
			},
			NumberOfTicks: 5,
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, tickDamage*float64(dot.Aura.GetStacks()), isRollover)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dot := spell.Dot(target)
			if dot.IsActive() {
				dot.Refresh(sim)
				dot.AddStack(sim)
			} else {
				dot.Apply(sim)
				dot.SetStacks(sim, 1)
			}
			// Snapshot again once the stacks are in, since the damage grows with them.
			dot.TakeSnapshot(sim, false)
		},
	})
}
