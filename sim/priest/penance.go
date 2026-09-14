package priest

import (
	"time"

	"github.com/wowsims/classic/sim/core"
)

const PenanceTicks = 3

func (priest *Priest) registerPenanceSpell() {
	if !priest.Talents.Penance {
		return
	}

	// The demo tooltip showed neither a mana cost nor a cooldown, both are taken from the
	// spell of the same name.
	// TODO: beta will confirm the cost, the cooldown and the level 60 damage.
	baseDamage := 93.0
	spellCoeff := 0.268

	priest.Penance = priest.RegisterSpell(core.SpellConfig{
		SpellCode:   SpellCode_PriestPenance,
		ActionID:    core.ActionID{SpellID: 47540},
		SpellSchool: core.SpellSchoolHoly,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagPriest | core.SpellFlagAPL | core.SpellFlagChanneled,

		RequiredLevel: 60,
		Rank:          1,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.16,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: time.Second * 10,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Penance",
			},

			// The bolts land at 0/1/2 sec in game, the sim spreads them evenly over the channel.
			NumberOfTicks:    PenanceTicks,
			TickLength:       time.Second * 2 / PenanceTicks,
			BonusCoefficient: spellCoeff,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, baseDamage, isRollover)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},

		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			return spell.CalcPeriodicDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicAlwaysHit)
		},
	})
}
