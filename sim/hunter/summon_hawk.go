package hunter

import (
	"time"

	"github.com/wowsims/classic/sim/core"
)

// Summon Hawk shares its cooldown with Arcane Shot, and Ferocity and Unleashed Fury buff hawks
// the same way they buff pets.
//
// The beta client (1293241, 1293525-1293527) gives the dive bomb, 32/47/85/108 plus 5% of ranged
// attack power, the mana cost and the 18 sec hawk (1293248). The hawk that stays is a guardian whose
// swings the client does not describe, so the assault is modelled as the rank's dive bomb base damage
// every 3 sec. Only one hawk at a time is modelled, not the two the tooltip allows.
func (hunter *Hunter) registerSummonHawkSpell(timer *core.Timer) {
	if !hunter.Talents.SummonHawk {
		return
	}

	rank := 1
	switch {
	case hunter.Level >= 60:
		rank = 4
	case hunter.Level >= 48:
		rank = 3
	case hunter.Level >= 36:
		rank = 2
	}
	spellId := [5]int32{0, 1293241, 1293525, 1293526, 1293527}[rank]
	baseDamage := [5]float64{0, 32, 47, 85, 108}[rank]
	manaCost := [5]float64{0, 80, 105, 135, 190}[rank]

	hunter.SummonHawk = hunter.RegisterSpell(core.SpellConfig{
		SpellCode:      SpellCode_HunterSummonHawk,
		ClassSpellMask: SpellMaskSummonHawk,
		ActionID:       core.ActionID{SpellID: spellId},
		Rank:           rank,
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true, // Hunter GCD is locked at 1.5s
			CD: core.Cooldown{
				Timer:    timer,
				Duration: time.Second * 6,
			},
		},

		BonusCritRating:  2 * float64(hunter.Talents.Ferocity) * core.CritRatingPerCritChance,
		DamageMultiplier: 1 + 0.03*float64(hunter.Talents.UnleashedFury),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Summon Hawk" + hunter.Label,
			},
			NumberOfTicks: 6,
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, baseDamage, isRollover)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := baseDamage + 0.05*spell.RangedAttackPower(target, false)
			result := spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
		},
	})
}
