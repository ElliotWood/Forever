package warlock

import (
	"github.com/wowsims/forever/sim/core"
)

// Forever folded Curse of Shadow into Curse of the Elements: its top rank (1311680) drops every
// magic resistance and raises all magic damage taken by 10%, which is what core's aura already
// does. Malediction is no longer a bonus on the curse - the beta client makes it a flat damage
// modifier on the warlock's own spells - so no ranks are passed in.
func (warlock *Warlock) registerCurseOfElements() {
	rank := spellData.CurseOfTheElements.HighestRank()

	warlock.CurseOfElementsAuras = warlock.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return core.CurseOfElementsAura(target, warlock.Index, 0)
	})

	warlock.CurseOfElements = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.SpellID},
		SpellSchool:    rank.SpellSchool,
		DefenseType:    rank.DefenseType,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellCurseOfElements,

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
				aura := warlock.CurseOfElementsAuras.Get(target)
				warlock.takeCurseSlot(sim, target, aura)
				aura.Activate(sim)
			}

			spell.DealOutcome(sim, result)
		},

		RelatedAuraArrays: warlock.CurseOfElementsAuras.ToMap(),
	})
}
