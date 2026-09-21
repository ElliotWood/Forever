package rogue

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

var garroteRank = spellData.Garrote.HighestRank()

func (rogue *Rogue) registerGarrote() {
	tick := garroteRank.Periodic.(shared.SpellDataPeriodic)

	rogue.Garrote = rogue.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: garroteRank.SpellID},
		SpellSchool:    garroteRank.SpellSchool,
		DefenseType:    garroteRank.DefenseType,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | SpellFlagBuilder | core.SpellFlagAPL,
		ClassSpellMask: RogueSpellGarrote,
		MaxRange:       core.MaxMeleeRange,

		EnergyCost: core.EnergyCostOptions{
			Cost:   garroteRank.Cost,
			Refund: garroteRank.MissRefund(),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: garroteRank.GCD,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			if !rogue.IsStealthed() {
				return false
			}
			// Dirty Deeds drops the positional requirement.
			return rogue.Talents.DirtyDeeds > 0 || !rogue.PseudoStats.InFrontOfTarget
		},

		DamageMultiplierAdditive: 1,
		DamageMultiplier:         1,
		ThreatMultiplier:         1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Garrote",
				Tag:   RogueBleedTag,
			},
			NumberOfTicks: tick.NumberOfTicks,
			TickLength:    tick.TickLength,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.SnapshotPhysical(target, tick.Tick+dot.Spell.MeleeAttackPower(target)*0.03)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, shared.PeriodicTickOutcome(garroteRank, dot))
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialNoBlockDodgeParryNoCrit)
			if result.Landed() {
				rogue.AddComboPoints(sim, 1, spell.ComboPointMetrics())
				spell.Dot(target).Apply(sim)
			} else {
				spell.IssueRefund(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})
}
