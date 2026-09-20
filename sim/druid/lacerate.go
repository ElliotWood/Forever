package druid

import (
	"time"

	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

var lacerateRank = shared.WithSpellDataFlatThreat(spellData.Lacerate, 267).HighestRank()
var lacerateTick = lacerateRank.Periodic.(shared.SpellDataPeriodic)

func (druid *Druid) registerLacerateSpell() {
	tickDamageBase := lacerateTick.Tick

	druid.Lacerate = druid.RegisterSpell(Bear, core.SpellConfig{
		ActionID:       core.ActionID{SpellID: lacerateRank.SpellID},
		SpellSchool:    lacerateRank.SpellSchool,
		DefenseType:    lacerateRank.DefenseType,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		ClassSpellMask: DruidSpellLacerate,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		RageCost: core.RageCostOptions{
			Cost:   lacerateRank.Cost,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: lacerateRank.GCD,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 0.5,
		FlatThreatBonus:  lacerateRank.FlatThreatBonus,
		MaxRange:         core.MaxMeleeRange,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:     "Lacerate",
				MaxStacks: 5,
				Duration:  time.Second * 15,
			},
			NumberOfTicks: lacerateTick.NumberOfTicks,
			TickLength:    lacerateTick.TickLength,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				perStack := tickDamageBase + druid.LacerateTickBonus + 0.01*dot.Spell.MeleeAttackPower(target)
				dot.SnapshotPhysical(target, perStack*float64(dot.Aura.GetStacks()))
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := tickDamageBase + druid.LacerateTickBonus + 0.01*spell.MeleeAttackPower(target)
			if druid.MangleAuras != nil && druid.MangleAuras.Get(target).IsActive() {
				baseDamage *= 1.3
			}
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if result.Landed() {
				dot := spell.Dot(target)
				if dot.IsActive() {
					dot.Refresh(sim)
					dot.AddStack(sim)
					dot.TakeSnapshot(sim)
				} else {
					dot.Apply(sim)
					dot.SetStacks(sim, 1)
					dot.TakeSnapshot(sim)
				}
			} else {
				spell.IssueRefund(sim)
			}
		},
	})

	druid.Lacerate.ShortName = "Lacerate"
}
