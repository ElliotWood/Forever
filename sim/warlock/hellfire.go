package warlock

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

func (warlock *Warlock) registerHellfire() {
	rank := spellData.Hellfire.HighestRank()
	tick := rank.Periodic.(shared.SpellDataPeriodic)

	warlock.Hellfire = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.SpellID},
		SpellSchool:    rank.SpellSchool,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagChanneled | core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellHellfire,

		ManaCost: core.ManaCostOptions{FlatCost: rank.Cost},
		Cast:     core.CastConfig{DefaultCast: core.Cast{GCD: rank.GCD}},

		DamageMultiplierAdditive: 1,
		DamageMultiplier:         1,
		ThreatMultiplier:         1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Hellfire",
			},

			IsAOE:                true,
			TickLength:           tick.TickLength,
			NumberOfTicks:        tick.NumberOfTicks,
			HasteReducesDuration: true,
			AffectedByCastSpeed:  true,
			BonusCoefficient:     tick.Coef,

			OnTick: func(sim *core.Simulation, _ *core.Unit, dot *core.Dot) {
				// Rolled once: the warlock burns exactly what it deals.
				tickDamage := tick.Damage(sim)

				resultSlice := dot.Spell.CalcPeriodicAoeDamage(sim, tickDamage, dot.Spell.OutcomeTickMagicHitNoHitCounter)
				if resultSlice[0].Damage > warlock.CurrentHealth() {
					dot.Deactivate(sim)
				}

				dot.Spell.DealBatchedPeriodicDamage(sim)
				warlock.RemoveHealth(sim, tickDamage)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.AOEDot().Apply(sim)
		},
	})
}
