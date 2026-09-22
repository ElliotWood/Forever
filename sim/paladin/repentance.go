package paladin

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
)

var RepentanceRankMap = spellData.Repentance

// Repentance (talent)
// https://www.wowhead.com/forever/spell=20066
//
// Puts the enemy target in a state of meditation, incapacitating them for up to 6 sec. Any damage
// caused will awaken the target. Only works against Humanoids.
//
// The sim has no meditation to model; the cast puts the debuff on a humanoid target so Judgement
// of Command can see it as incapacitated, and nothing wakes it.
func (paladin *Paladin) registerRepentance() {
	row := RepentanceRankMap.HighestRank()

	auras := paladin.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    "Repentance",
			ActionID: core.ActionID{SpellID: row.SpellID},
			Duration: row.Duration,
			OnGain: func(_ *core.Aura, _ *core.Simulation) {
				target.PseudoStats.Stunned = true
			},
			OnExpire: func(_ *core.Aura, _ *core.Simulation) {
				target.PseudoStats.Stunned = false
			},
		})
	})

	paladin.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: row.SpellID},
		SpellSchool:    row.SpellSchool,
		DefenseType:    row.DefenseType,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskRepentance,
		MaxRange:       row.MaxRange,

		ManaCost: manaCost(row),
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: row.GCD,
			},
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: row.Cooldown,
			},
		},

		ExtraCastCondition: func(_ *core.Simulation, target *core.Unit) bool {
			return target.MobType == proto.MobType_MobTypeHumanoid
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				auras.Get(target).Activate(sim)
			}
		},

		RelatedAuraArrays: auras.ToMap(),
	})
}
