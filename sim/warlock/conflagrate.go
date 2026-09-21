package warlock

import (
	"github.com/wowsims/forever/sim/core"
)

// Conflagrate burns the Immolate on the target. Shadow and Flame's second effect (20% per point)
// is the chance it survives, so at 5/5 Conflagrate stops consuming Immolate altogether.
func (warlock *Warlock) registerConflagrate() {
	rank := spellData.Conflagrate.HighestRank()
	keepImmolateChance := spellData.ShadowAndFlame.EffectAt(1).FractionAt(warlock.Talents.ShadowAndFlame)

	warlock.Conflagrate = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.SpellID},
		SpellSchool:    rank.SpellSchool,
		DefenseType:    rank.DefenseType,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellConflagrate,

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
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warlock.Immolate.Dot(target).IsActive()
		},

		DamageMultiplierAdditive: 1,
		DamageMultiplier:         1,
		ThreatMultiplier:         1,
		BonusCoefficient:         rank.Direct.BonusCoefficient(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, rank.Direct.Damage(sim), spell.OutcomeMagicHitAndCrit)

			dot := warlock.Immolate.Dot(target)
			if dot.IsActive() && !sim.Proc(keepImmolateChance, "Shadow and Flame") {
				dot.Deactivate(sim)
			}
		},
	})
}
