package warrior

import (
	"time"

	"github.com/wowsims/forever/sim/core"
)

// TODO: Manual review needed -- spell 676 states a 20 rage cost, a 60 second cooldown and a
// 10 second disarm.
const (
	disarmRageCost int32 = 20
	disarmCooldown       = time.Second * 60
	disarmDuration       = time.Second * 10
)

func (warrior *Warrior) registerDisarm() {
	actionID := core.ActionID{SpellID: 676}

	// TODO: core has no disarm effect, so the aura only tracks the debuff's uptime.
	auras := warrior.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    "Disarm-" + warrior.Label,
			ActionID: actionID,
			Duration: disarmDuration,
		})
	})

	warrior.Disarm = warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskDisarm,
		MaxRange:       core.MaxMeleeRange,

		RageCost: core.RageCostOptions{
			Cost:   disarmRageCost,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: disarmCooldown,
			},
		},

		ThreatMultiplier: 1,

		// Spell 676's ShapeshiftMask is Defensive Stance only.
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.StanceMatches(DefensiveStance)
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMeleeSpecialHit)

			if result.Landed() {
				auras.Get(target).Activate(sim)
			} else {
				spell.IssueRefund(sim)
			}
		},

		RelatedAuraArrays: auras.ToMap(),
	})
}
