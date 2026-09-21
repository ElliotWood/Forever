package mage

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

func (mage *Mage) registerPyroblastSpell() {
	if !mage.Talents.Pyroblast {
		return
	}

	pyroblastRank := spellData.Pyroblast.HighestRank()
	pyroblastTick := pyroblastRank.Periodic.(shared.SpellDataPeriodic)

	mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: pyroblastRank.SpellID},
		SpellSchool:    pyroblastRank.SpellSchool,
		DefenseType:    pyroblastRank.DefenseType,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MageSpellPyroblast,
		MissileSpeed:   pyroblastRank.MissileSpeed,

		ManaCost: core.ManaCostOptions{
			FlatCost: pyroblastRank.Cost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      pyroblastRank.GCD,
				CastTime: pyroblastRank.CastTime,
			},
		},

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "PyroblastDoT",
			},
			NumberOfTicks:    pyroblastTick.NumberOfTicks,
			TickLength:       pyroblastTick.TickLength,
			BonusCoefficient: pyroblastTick.Coef,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.Snapshot(target, pyroblastTick.Tick)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, shared.PeriodicTickOutcome(pyroblastRank, dot))
			},
		},

		DamageMultiplier: 1,
		BonusCoefficient: pyroblastRank.Direct.BonusCoefficient(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcDamage(sim, target, pyroblastRank.Direct.Damage(sim), spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
				if result.Landed() {
					spell.Dot(target).Apply(sim)
				}
			})
		},
	})
}
