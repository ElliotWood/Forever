package paladin

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

func (paladin *Paladin) getHolyShockTimer() *core.Timer {
	if paladin.holyShockTimer == nil {
		paladin.holyShockTimer = paladin.NewTimer()
	}
	return paladin.holyShockTimer
}

var HolyShockRankMap = spellData.HolyShock

// Holy Shock
// https://www.wowhead.com/forever/spell=20473
//
// Blasts the target with Holy energy, causing X to Y Holy damage to an enemy,
// or X*1.267 to Y*1.267 healing to an ally.
func (paladin *Paladin) registerHolyShock(rankConfig shared.SpellData) {
	spellID := rankConfig.SpellID
	cost := rankConfig.Cost
	coefficient := rankConfig.Direct.BonusCoefficient()

	// Holy Shock heals for 1.267x the damage component of the spell.
	healingCoeff := 1.267

	paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellID},
		SpellSchool: core.SpellSchoolHoly,
		// The cast (20473 .. 33072) is a dummy; the damage (25912 .. 33073) and heal (25914 .. 33074)
		// spells it triggers are Magic in SpellCategories.
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskHolyShock,
		Rank:           rankConfig.Rank,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		MaxRange: rankConfig.MaxRange,

		ManaCost: core.ManaCostOptions{
			FlatCost: cost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: rankConfig.GCD,
			},
			CD: core.Cooldown{
				Timer:    paladin.getHolyShockTimer(),
				Duration: rankConfig.Cooldown,
			},
		},

		BonusCoefficient: coefficient,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if spell.Unit.IsOpponent(target) {
				damage := rankConfig.Direct.Damage(sim)
				spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)
			} else {
				// Temporarily configure the spell as a healing spell
				spell.Flags |= core.SpellFlagHelpful
				originalProcMask := spell.ProcMask
				spell.ProcMask = core.ProcMaskSpellHealing

				// TODO: Use healing power instead of holy power for healing calculations
				healing := rankConfig.Direct.Damage(sim) * healingCoeff
				spell.CalcAndDealHealing(sim, target, healing, spell.OutcomeHealingCrit)

				// Reset the spell to its original configuration
				spell.Flags &= ^core.SpellFlagHelpful
				spell.ProcMask = originalProcMask
			}
		},
	})
}
