package priest

import (
	"time"

	"github.com/wowsims/forever/sim/core"
)

// Penance is the Discipline talent the Smite build goes deep for: three Holy bolts over the channel,
// on a 12 second cooldown. Only the level 60 rank is registered, the one the rotation casts.
//
// The client's rank 3 bolt (180) is larger than rank 4's (131); the table is taken as it is. The
// bolts land at 0/1/2 sec in game, the sim spreads them evenly over the channel.
const PenanceTicks = 3

func (priest *Priest) registerPenanceSpell() {
	rank := spellData.Penance.HighestRank()
	bolt := spellData.PenanceTriggered.BySpellID(1316993)
	// Each bolt is its own direct Holy hit (1316993, School Damage), so every one can crit.
	boltCrits := bolt
	boltCrits.PeriodicCanCrit = true

	priest.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.SpellID},
		SpellSchool:    core.SpellSchoolHoly,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL | core.SpellFlagChanneled,
		ClassSpellMask: PriestSpellPenance,
		Rank:           rank.Rank,
		MaxRange:       rank.MaxRange,

		ManaCost: core.ManaCostOptions{
			FlatCost: rank.Cost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: rank.GCD,
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: rank.Cooldown,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Penance",
			},
			NumberOfTicks:       PenanceTicks,
			TickLength:          time.Second * 2 / PenanceTicks,
			AffectedByCastSpeed: false,
			BonusCoefficient:    bolt.Direct.BonusCoefficient(),

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.Snapshot(target, bolt.Direct.Damage(sim))
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, priestTickOutcome(boltCrits, dot))
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},

		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			if useSnapshot {
				return spell.Dot(target).CalcSnapshotDamage(sim, target, spell.OutcomeExpectedMagicHit)
			}
			return spell.CalcPeriodicDamage(sim, target, bolt.Direct.Damage(sim), spell.OutcomeExpectedMagicHit)
		},
	})
}
