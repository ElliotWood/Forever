package paladin

import (
	"time"

	"github.com/wowsims/classic/sim/core/proto"

	"github.com/wowsims/classic/sim/core"
)

func (paladin *Paladin) registerHammerOfWrath() {
	// Beta client 1.60.1.69893, with each rank's growth to its max level folded in: rank 3 504-566
	// (Classic's client reads 504-556) -> 473-523. Cost, cast, cooldown and the 0.429 are unchanged.
	ranks := []struct {
		level     int32
		minDamage float64
		maxDamage float64
		manaCost  float64
	}{
		{level: 44, manaCost: 295, minDamage: 285, maxDamage: 315},
		{level: 52, manaCost: 360, minDamage: 382, maxDamage: 421},
		{level: 60, manaCost: 425, minDamage: 473, maxDamage: 523},
	}

	cd := core.Cooldown{
		Timer:    paladin.NewTimer(),
		Duration: time.Second * 6,
	}

	// Instrument of Law: 0.5 sec a rank, confirmed by the beta client's talent data.
	castTime := time.Second - time.Millisecond*500*time.Duration(paladin.Talents.InstrumentOfLaw)

	for i, rank := range ranks {
		rank := rank
		spellID := []int32{24275, 24274, 24239}[i]
		if paladin.Level < rank.level {
			break
		}

		paladin.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: spellID},
			SpellSchool: core.SpellSchoolHoly,
			DefenseType: core.DefenseTypeRanged,
			ProcMask:    core.ProcMaskRangedSpecial, // TODO to be tested
			Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
			CastType:    proto.CastType_CastTypeRanged,

			Rank:          i + 1,
			RequiredLevel: int(rank.level),
			SpellCode:     SpellCode_PaladinHammerOfWrath,

			ManaCost: core.ManaCostOptions{
				FlatCost:   rank.manaCost,
				Multiplier: paladin.holyConduit(),
			},
			Cast: core.CastConfig{
				DefaultCast: core.Cast{
					GCD:      time.Second,
					CastTime: castTime,
				},
				IgnoreHaste: true,
				CD:          cd,
			},

			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			BonusCoefficient: 0.429,

			ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
				return sim.IsExecutePhase20()
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				damage := sim.Roll(rank.minDamage, rank.maxDamage)
				spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeRangedHitAndCrit)
			},
		})
	}
}
