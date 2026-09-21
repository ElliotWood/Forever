package hunter

import (
	"time"

	"github.com/wowsims/forever/sim/core"
)

// Summon Hawk shares its cooldown with Arcane Shot, and Ferocity and Unleashed Fury buff hawks the
// same way they buff pets.
//
// The beta client (1293241, 1293525-1293527) gives the dive bomb, 32/47/85/108 plus 5% of ranged
// attack power, the mana cost and the 18 sec hawk (1293248). The hawk that stays is a guardian whose
// swings the client does not describe, so the assault is modelled as the rank's dive bomb base
// damage every 3 sec. Only one hawk at a time is modelled, not the two the tooltip allows.
func (hunter *Hunter) registerSummonHawkSpell(timer *core.Timer) {
	if !hunter.Talents.SummonHawk {
		return
	}

	rank := spellData.SummonHawk.HighestRank()
	baseDamage := rank.Direct.Damage

	hunter.SummonHawk = hunter.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.SpellID},
		SpellSchool:    rank.SpellSchool,
		DefenseType:    core.DefenseTypeMelee,
		ClassSpellMask: HunterSpellSummonHawk,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		MaxRange:       rank.MaxRange,

		ManaCost: core.ManaCostOptions{
			FlatCost: rank.Cost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: rank.GCD,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    timer,
				Duration: rank.Cooldown,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Summon Hawk" + hunter.Label,
			},
			NumberOfTicks: 6,
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.Snapshot(target, baseDamage(sim))
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := baseDamage(sim) + 0.05*spell.RangedAttackPower(target)
			result := spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
		},
	})
}
