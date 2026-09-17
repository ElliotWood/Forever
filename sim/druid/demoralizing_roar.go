package druid

import (
	"github.com/wowsims/classic/sim/core"
)

const DemoralizingRoarRanks = 5

var DemoralizingRoarSpellId = [DemoralizingRoarRanks + 1]int32{0, 99, 1735, 9490, 9747, 9898}
var DemoralizingRoarLevel = [DemoralizingRoarRanks + 1]int{0, 10, 20, 30, 40, 50}

func (druid *Druid) registerDemoralizingRoarSpell() {
	rank := map[int32]int{
		25: 2,
		40: 4,
		50: 5,
		60: 5,
	}[druid.Level]

	druid.DemoralizingRoarAuras = druid.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		// Feral Aggression is gone from the Forever tree and more than folded in: the beta client's rank 5
		// is 193 plus 1.4 a level, 204 at 60, where Classic's is 130 plus 1 a level. 204 is Classic's 138
		// with six points of the old 8%, which is how the shared aura is asked for it.
		return core.DemoralizingRoarAura(target, 6)
	})

	druid.DemoralizingRoar = druid.RegisterSpell(Bear, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: DemoralizingRoarSpellId[rank]},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagAPL,

		Rank:          rank,
		RequiredLevel: DemoralizingRoarLevel[rank],

		RageCost: core.RageCostOptions{
			Cost: 10,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		ThreatMultiplier: 1,
		FlatThreatBonus:  2 * float64(DemoralizingRoarLevel[rank]),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				result := spell.CalcAndDealOutcome(sim, aoeTarget, spell.OutcomeMagicHit)
				if result.Landed() {
					druid.DemoralizingRoarAuras.Get(aoeTarget).Activate(sim)
				}
			}
		},

		RelatedAuras: []core.AuraArray{druid.DemoralizingRoarAuras},
	})
}
