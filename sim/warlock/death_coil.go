package warlock

import (
	"github.com/wowsims/forever/sim/core"
)

// The generator files Death Coil's damage under E_HEALTH_LEECH, which it has no Direct role for, so
// the row's Direct is nil and the damage is read off the effect. The 0.214 coefficient is ours: the
// client states none for a leech.
func (warlock *Warlock) registerDeathCoil() {
	rank := spellData.DeathCoil.HighestRank()
	baseDamage := spellData.DeathCoil.EffectAt(0).ValueAt(rank.Rank)

	healingSpell := warlock.GetOrRegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: rank.SpellID}.WithTag(1),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskSpellHealing,
		Flags:       core.SpellFlagPassiveSpell | core.SpellFlagHelpful,

		DamageMultiplier: 1,
		ThreatMultiplier: 0,
	})

	warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.SpellID},
		SpellSchool:    rank.SpellSchool,
		DefenseType:    rank.DefenseType,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL | core.SpellFlagBinary,
		ClassSpellMask: WarlockSpellDeathCoil,
		MissileSpeed:   rank.MissileSpeed,
		MaxRange:       rank.MaxRange,

		ManaCost: core.ManaCostOptions{FlatCost: rank.Cost},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: rank.GCD,
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: rank.Cooldown,
			},
		},

		DamageMultiplierAdditive: 1,
		DamageMultiplier:         1,
		ThreatMultiplier:         1,
		BonusCoefficient:         0.214,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
				if result.Landed() {
					healingSpell.CalcAndDealHealing(sim, healingSpell.Unit, result.Damage, healingSpell.OutcomeHealing)
				}
			})
		},
	})
}
