package warlock

import (
	"github.com/wowsims/forever/sim/core"
)

func (warlock *Warlock) registerCurseOfRecklessness() {
	rank := spellData.CurseOfRecklessness.HighestRank()

	warlock.CurseOfRecklessnessAuras = warlock.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return core.CurseOfRecklessnessAura(target, warlock.Index)
	})

	warlock.CurseOfRecklessness = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.SpellID},
		SpellSchool:    rank.SpellSchool,
		DefenseType:    rank.DefenseType,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellCurseOfRecklessness,

		ManaCost: core.ManaCostOptions{FlatCost: rank.Cost},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: rank.GCD,
			},
		},

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)
			if result.Landed() {
				aura := warlock.CurseOfRecklessnessAuras.Get(target)
				warlock.takeCurseSlot(sim, target, aura)
				aura.Activate(sim)
			}

			spell.DealOutcome(sim, result)
		},

		RelatedAuraArrays: warlock.CurseOfRecklessnessAuras.ToMap(),
	})
}
