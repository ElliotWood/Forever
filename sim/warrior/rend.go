package warrior

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

var rendRank = spellData.Rend.HighestRank()

// TODO: Ingame testing needed if Rend has a coef
func (warrior *Warrior) registerRend() {
	tick := rendRank.Periodic.(shared.SpellDataPeriodic)

	warrior.Rend = warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rendRank.SpellID},
		SpellSchool:    rendRank.SpellSchool,
		DefenseType:    rendRank.DefenseType,
		ClassSpellMask: SpellMaskRend,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagAPL | SpellFlagBleed,

		RageCost: core.RageCostOptions{
			Cost:   rendRank.Cost,
			Refund: rendRank.MissRefund(),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: rendRank.GCD,
			},
			IgnoreHaste: true,
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.StanceMatches(BattleStance | DefensiveStance)
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Rend",
			},
			NumberOfTicks: tick.NumberOfTicks,
			TickLength:    tick.TickLength,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.Spell.CalcAndDealPeriodicDamage(sim, target, tick.Tick, shared.PeriodicTickOutcome(rendRank, dot))
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}
