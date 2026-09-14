package warlock

import (
	"time"

	"github.com/wowsims/classic/sim/core"
)

func (warlock *Warlock) registerIncinerateSpell() {
	if !warlock.Talents.Incinerate {
		return
	}

	// Level 60 values aren't datamined yet, scaled up from the demo tooltip
	baseDamage := []float64{380, 440}
	spellCoeff := 0.714
	manaCost := 300.0
	castTime := time.Millisecond * 2500

	warlock.Incinerate = warlock.RegisterSpell(core.SpellConfig{
		SpellCode:     SpellCode_WarlockIncinerate,
		ActionID:      core.ActionID{SpellID: 29722},
		SpellSchool:   core.SpellSchoolFire,
		DefenseType:   core.DefenseTypeMagic,
		ProcMask:      core.ProcMaskSpellDamage,
		Flags:         core.SpellFlagAPL | core.SpellFlagResetAttackSwing | WarlockFlagDestruction,
		RequiredLevel: 60,
		Rank:          1,
		MissileSpeed:  24,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: castTime,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: spellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := sim.Roll(baseDamage[0], baseDamage[1])
			if warlock.getActiveImmolateSpell(target) != nil {
				damage *= 1.25
			}

			result := spell.CalcDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
