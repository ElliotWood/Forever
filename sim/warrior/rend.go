package warrior

import (
	"github.com/wowsims/forever/sim/core"
)

var rendRank = spellData.Rend.Highest()

// TODO: Ingame testing needed if Rend has a coef
func (warrior *Warrior) registerRend() {
	tick := rendRank.PeriodicEffect()

	warrior.Rend = warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rendRank.ID},
		SpellSchool:    rendRank.SpellSchool(),
		DefenseType:    rendRank.DefenseTypeCore(),
		ClassSpellMask: SpellMaskRend,
		ClassFlags:     SpellFlagsRend,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,

		RageCost: core.RageCostOptions{
			Cost:   rageCost(rendRank),
			Refund: rendRank.MissRefund(),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: rendRank.GCD(),
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
			NumberOfTicks: int32(rendRank.Duration() / tick.Period()),
			TickLength:    tick.Period(),
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.Spell.CalcAndDealPeriodicDamage(sim, target, tick.Average(core.CharacterLevel), rendRank.TickOutcome(dot))
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
