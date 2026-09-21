package druid

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

var hurricaneRank = spellData.Hurricane.HighestRank()

func (druid *Druid) registerHurricaneSpell() {
	hurricaneTick := hurricaneRank.Periodic.(shared.SpellDataPeriodic)

	druid.Hurricane = druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
		ActionID:       core.ActionID{SpellID: hurricaneRank.SpellID},
		SpellSchool:    hurricaneRank.SpellSchool,
		DefenseType:    hurricaneRank.DefenseType,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagChanneled | core.SpellFlagAPL,
		ClassSpellMask: DruidSpellHurricane,
		MaxRange:       hurricaneRank.MaxRange,
		Rank:           hurricaneRank.Rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: hurricaneRank.Cost,
		},
		// Forever states no cooldown on Hurricane (the client rows carry none), so the spell is
		// registered without one rather than with an invented duration.
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: hurricaneRank.GCD,
			},
		},
		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "Hurricane (Aura)",
			},
			NumberOfTicks:       hurricaneTick.NumberOfTicks,
			TickLength:          hurricaneTick.TickLength,
			AffectedByCastSpeed: true,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				druid.Hurricane.RelatedDotSpell.Cast(sim, target)
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			spell.AOEDot().Apply(sim)
		},
	})

	druid.Hurricane.RelatedDotSpell = druid.Unit.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: hurricaneTick.SpellID},
		SpellSchool:    hurricaneRank.SpellSchool,
		DefenseType:    hurricaneRank.DefenseType,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: DruidSpellHurricane,
		// The tick is its own client row that the channel triggers, a proc rather than a cast.
		Flags: core.SpellFlagProc,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: hurricaneTick.Coef,

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			spell.CalcAndDealAoeDamage(sim, hurricaneTick.Tick, spell.OutcomeMagicHit)
		},
	})
}
