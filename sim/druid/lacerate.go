package druid

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

// The client carries no threat for Lacerate, so the 3.33x multiplier is still Season of Discovery's.
var lacerateRank = spellData.Lacerate.HighestRank()
var lacerateTick = lacerateRank.Periodic.(shared.SpellDataPeriodic)

// The hit is a share of weapon damage per stack, which the client states on the rank's dummy
// effect (10 at every rank).
var lacerateWeaponPctPerStack = spellData.Lacerate.EffectAt(1).FractionAt(lacerateRank.Rank)

const LacerateMaxStacks int32 = 5

// Forever's bleed no longer scales with attack power: the client states a flat tick per stack.
func (druid *Druid) registerLacerateSpell() {
	druid.Lacerate = druid.RegisterSpell(Bear, core.SpellConfig{
		ActionID:       core.ActionID{SpellID: lacerateRank.SpellID},
		SpellSchool:    lacerateRank.SpellSchool,
		DefenseType:    lacerateRank.DefenseType,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		ClassSpellMask: DruidSpellLacerate,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		Rank:           lacerateRank.Rank,

		RageCost: core.RageCostOptions{
			Cost:   lacerateRank.Cost,
			Refund: lacerateRank.MissRefund(),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: lacerateRank.GCD,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 3.33,
		MaxRange:         core.MaxMeleeRange,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:     "Lacerate",
				MaxStacks: LacerateMaxStacks,
				Duration:  lacerateRank.Duration,
			},
			NumberOfTicks: lacerateTick.NumberOfTicks,
			TickLength:    lacerateTick.TickLength,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.SnapshotPhysical(target, lacerateTick.Tick*float64(dot.Aura.GetStacks()))
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, shared.PeriodicTickOutcome(lacerateRank, dot))
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dot := spell.Dot(target)
			stacks := min(dot.Aura.GetStacks()+1, LacerateMaxStacks)
			baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower(target)) * lacerateWeaponPctPerStack * float64(stacks)
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if result.Landed() {
				if dot.IsActive() {
					dot.Refresh(sim)
					dot.AddStack(sim)
				} else {
					dot.Apply(sim)
					dot.SetStacks(sim, 1)
				}
				// Snapshot again once the stacks are in, since the damage grows with them.
				dot.TakeSnapshot(sim)
			} else {
				spell.IssueRefund(sim)
			}
		},
	})

	druid.Lacerate.ShortName = "Lacerate"
}
