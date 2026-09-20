package mage

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

var blizzardRank = spellData.Blizzard.HighestRank()

func (mage *Mage) registerBlizzardSpell() {
	blizzardActionId := core.ActionID{SpellID: blizzardRank.SpellID}
	blizzardTick := blizzardRank.Periodic.(shared.SpellDataPeriodic)

	blizzardTickSpell := mage.RegisterSpell(core.SpellConfig{
		ActionID:       blizzardActionId,
		SpellSchool:    core.SpellSchoolFrost,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: MageSpellBlizzard,

		DamageMultiplier: 1,
		BonusCoefficient: blizzardTick.Coef,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			spell.CalcAndDealAoeDamage(sim, blizzardTick.Tick, spell.OutcomeMagicHit)
		},
	})

	mage.RegisterSpell(core.SpellConfig{
		ActionID:       blizzardActionId,
		SpellSchool:    core.SpellSchoolFrost,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagChanneled | core.SpellFlagAPL,
		ClassSpellMask: MageSpellBlizzard,
		ManaCost: core.ManaCostOptions{
			FlatCost: blizzardRank.Cost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: blizzardRank.GCD,
			},
		},
		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label:    "Blizzard",
				ActionID: blizzardActionId,
			},
			NumberOfTicks:        blizzardTick.NumberOfTicks,
			TickLength:           blizzardTick.TickLength,
			AffectedByCastSpeed:  true,
			HasteReducesDuration: true,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				blizzardTickSpell.Cast(sim, target)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.AOEDot().Apply(sim)
		},
	})
}
