package druid

import (
	"github.com/wowsims/forever/sim/core"
)

var faerieFireRank = spellData.FaerieFire.BySpellID(26993)
var faerieFireFeralRank = spellData.FaerieFireFeral.BySpellID(27011)

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
func (druid *Druid) registerFaerieFireFeralSpell() {
	druid.FaerieFireAuras = druid.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		// TODO: Forever drops Improved Faerie Fire; untalented (0 points) until we know
		// whether the effect moved onto another talent.
		return core.FaerieFireAura(target, 0)
	})

	druid.FaerieFireFeral = druid.RegisterSpell(Cat|Bear, core.SpellConfig{
		ClassSpellMask: DruidSpellFaerieFireFeral,
		ActionID:       core.ActionID{SpellID: faerieFireFeralRank.SpellID},
		SpellSchool:    faerieFireFeralRank.SpellSchool,
		DefenseType:    faerieFireFeralRank.DefenseType,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: faerieFireFeralRank.GCD,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: faerieFireFeralRank.Cooldown,
			},
		},

		ThreatMultiplier: 1,
		FlatThreatBonus:  132,
		MaxRange:         faerieFireFeralRank.MaxRange,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)

			if result.Landed() {
				druid.FaerieFireAuras.Get(target).Activate(sim)
			}
		},

		RelatedAuraArrays: druid.FaerieFireAuras.ToMap(),
	})
}
