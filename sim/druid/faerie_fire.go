package druid

import (
	"github.com/wowsims/forever/sim/core"
)

var faerieFireRank = spellData.FaerieFire.HighestRank()

// Forever has no Faerie Fire (Feral): the client keeps only the Balance line (770, 778, 9749,
// 9907), so one registration serves every form.
func (druid *Druid) registerFaerieFireSpell() {
	druid.FaerieFireAuras = druid.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		// TODO: Forever drops Improved Faerie Fire; untalented (0 points) until we know
		// whether the effect moved onto another talent.
		return core.FaerieFireAura(target, 0)
	})

	druid.FaerieFire = druid.RegisterSpell(Any, core.SpellConfig{
		ClassSpellMask: DruidSpellFaerieFire,
		ActionID:       core.ActionID{SpellID: faerieFireRank.SpellID},
		SpellSchool:    faerieFireRank.SpellSchool,
		DefenseType:    faerieFireRank.DefenseType,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		Rank:           faerieFireRank.Rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: faerieFireRank.Cost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: faerieFireRank.GCD,
			},
		},

		ThreatMultiplier: 1,
		// Two threat a level, the sim's long-standing value; the client states none.
		FlatThreatBonus: 2 * float64(core.CharacterLevel),
		MaxRange:        faerieFireRank.MaxRange,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)

			if result.Landed() {
				druid.FaerieFireAuras.Get(target).Activate(sim)
			}
		},

		RelatedAuraArrays: druid.FaerieFireAuras.ToMap(),
	})
}
