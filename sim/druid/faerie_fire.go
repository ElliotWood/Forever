package druid

import (
	"github.com/wowsims/forever/sim/core"
)

var faerieFireRank = spellData.FaerieFire.HighestRank()

func (druid *Druid) registerFaerieFireSpell() {
	auras := druid.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		// TODO: Forever drops Improved Faerie Fire; untalented (0 points) until we know
		// whether the effect moved onto another talent.
		return core.FaerieFireAura(target, 0)
	})

	druid.FaerieFire = druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
		ClassSpellMask: DruidSpellFaerieFire,
		ActionID:       core.ActionID{SpellID: faerieFireRank.SpellID},
		SpellSchool:    faerieFireRank.SpellSchool,
		DefenseType:    faerieFireRank.DefenseType,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			FlatCost: faerieFireRank.Cost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: faerieFireRank.GCD,
			},
		},

		ThreatMultiplier: 1,
		FlatThreatBonus:  132,
		MaxRange:         faerieFireRank.MaxRange,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)

			if result.Landed() {
				auras.Get(target).Activate(sim)
			}
		},

		RelatedAuraArrays: auras.ToMap(),
	})
}

// TODO: uncalled -- Forever drops the Faerie Fire (Feral) talent; re-gate before wiring
// back into RegisterFeralCatSpells/RegisterFeralTankSpells.
// TODO: To be implemented. The Forever client ships no rank ladder the generator can
// read for this ability -- it survives as a single spell with no "Rank N" subtext and
// no ranked SkillLineAbility row -- so there is no data to build the spell from.
func (druid *Druid) registerFaerieFireFeralSpell() {
	panic("To be implemented")
}
