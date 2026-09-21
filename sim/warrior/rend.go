package warrior

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

var rendRank = spellData.Rend.HighestRank()

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
			Refund: 0.8,
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
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				// TODO: Manual review needed -- spell 11574 states no weapon or attack power scaling,
				// so the 0.00743 per point of average weapon damage is hand-supplied.
				dot.SnapshotBaseDamage = tick.Tick + warrior.AutoAttacks.MH().CalculateAverageWeaponDamage(dot.Spell.MeleeAttackPower(target))*0.00743
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex], true)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
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
