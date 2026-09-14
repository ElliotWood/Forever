package warlock

import (
	"time"

	"github.com/wowsims/classic/sim/core"
)

func (warlock *Warlock) registerDrainHopeSpell() {
	if !warlock.Talents.DrainHope {
		return
	}

	// Level 60 values aren't datamined yet, scaled up from the demo tooltip
	numTicks := int32(6)
	tickLength := time.Second
	baseDamage := 90.0
	spellCoeff := 0.1
	manaCost := 320.0

	warlock.DrainHope = warlock.RegisterSpell(core.SpellConfig{
		SpellCode:   SpellCode_WarlockDrainHope,
		ActionID:    core.ActionID{SpellID: 11704},
		SpellSchool: core.SpellSchoolShadow,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL | core.SpellFlagChanneled | core.SpellFlagResetAttackSwing | WarlockFlagAffliction,

		RequiredLevel: 60,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "DrainHope-" + warlock.Label,
			},
			NumberOfTicks:    numTicks,
			TickLength:       tickLength,
			BonusCoefficient: spellCoeff,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, baseDamage, isRollover)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)
			if result.Landed() {
				dot := spell.Dot(target)
				dot.Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})

	// Other shadow dots on the target tick for 10% more while it's being drained
	for _, target := range warlock.Env.Encounter.TargetUnits {
		target.AddDynamicDamageTakenModifier(func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Unit != &warlock.Unit || spell == warlock.DrainHope {
				return
			}

			if spell.SpellSchool.Matches(core.SpellSchoolShadow) && len(spell.Dots()) > 0 && warlock.DrainHope.Dot(result.Target).IsActive() {
				result.Damage *= 1.1
			}
		})
	}
}
