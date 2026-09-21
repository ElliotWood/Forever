package mage

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

func (mage *Mage) registerFireballSpell() {
	fireballRank := spellData.Fireball.HighestRank()
	fireballTick := fireballRank.Periodic.(shared.SpellDataPeriodic)

	mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: fireballRank.SpellID},
		SpellSchool:    fireballRank.SpellSchool,
		DefenseType:    fireballRank.DefenseType,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MageSpellFireball,
		MissileSpeed:   fireballRank.MissileSpeed,

		ManaCost: core.ManaCostOptions{
			FlatCost: fireballRank.Cost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      fireballRank.GCD,
				CastTime: fireballRank.CastTime,
			},
		},

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "FireballDoT",
			},
			NumberOfTicks:    fireballTick.NumberOfTicks,
			TickLength:       fireballTick.TickLength,
			BonusCoefficient: fireballTick.Coef,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.Snapshot(target, fireballTick.Tick)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, shared.PeriodicTickOutcome(fireballRank, dot))
			},
		},

		DamageMultiplier: 1,
		BonusCoefficient: fireballRank.Direct.BonusCoefficient(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcDamage(sim, target, fireballRank.Direct.Damage(sim), spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
				if result.Landed() {
					spell.Dot(target).Apply(sim)
				}
			})
		},
	})
}
