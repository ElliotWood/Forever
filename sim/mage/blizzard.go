package mage

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

// Forever's Blizzard is an area trigger that casts a damage spell (1279949 at rank 6) every second.
func (mage *Mage) registerBlizzardSpell() {
	blizzardRank := spellData.Blizzard.HighestRank()
	blizzardTick := blizzardRank.Periodic.(shared.SpellDataPeriodic)
	blizzardActionId := core.ActionID{SpellID: blizzardRank.SpellID}

	// Improved Blizzard's chill, a separate spell so Fingers of Frost can roll on it.
	var improvedBlizzard *core.Spell
	if mage.Talents.ImprovedBlizzard > 0 {
		improvedBlizzardRank := spellData.ImprovedBlizzardTriggered.HighestRank()
		improvedBlizzard = mage.RegisterSpell(core.SpellConfig{
			ActionID:       core.ActionID{SpellID: improvedBlizzardRank.SpellID},
			SpellSchool:    core.SpellSchoolFrost,
			DefenseType:    core.DefenseTypeMagic,
			ProcMask:       core.ProcMaskSpellDamageProc,
			Flags:          core.SpellFlagNoLogs | core.SpellFlagNoMetrics | core.SpellFlagNoOnCastComplete,
			ClassSpellMask: MageSpellImprovedBlizzard,
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHitNoHitCounter)
			},
		})
	}

	blizzardTickSpell := mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: blizzardTick.SpellID},
		SpellSchool:    blizzardRank.SpellSchool,
		DefenseType:    blizzardRank.DefenseType,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagNoOnCastComplete,
		ClassSpellMask: MageSpellBlizzard,

		DamageMultiplier: 1,
		BonusCoefficient: blizzardTick.Coef,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			results := spell.CalcAndDealAoeDamage(sim, blizzardTick.Tick, spell.OutcomeMagicHit)
			if improvedBlizzard == nil {
				return
			}
			for _, result := range results {
				if result.Landed() {
					improvedBlizzard.Cast(sim, result.Target)
				}
			}
		},
	})

	mage.RegisterSpell(core.SpellConfig{
		ActionID:       blizzardActionId,
		SpellSchool:    blizzardRank.SpellSchool,
		DefenseType:    blizzardRank.DefenseType,
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
			NumberOfTicks: blizzardTick.NumberOfTicks,
			TickLength:    blizzardTick.TickLength,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				blizzardTickSpell.Cast(sim, target)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.AOEDot().Apply(sim)
		},
	})
}
