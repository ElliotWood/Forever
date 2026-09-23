package paladin

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

func (paladin *Paladin) registerAuras() {
	paladin.registerDevotionAura()
	paladin.registerRetributionAura()
	paladin.registerConcentrationAura()
	paladin.registerFireResistanceAura()
	paladin.registerFrostResistanceAura()
	paladin.registerShadowResistanceAura()
}

// The rank of a paladin aura as core wants it: the spell the paladin cast, its rank, and the number
// its row states.
func auraRank(row shared.SpellData) core.PaladinAuraRank {
	return core.PaladinAuraRank{SpellID: row.SpellID, Rank: row.Rank, Value: shared.SpellDataMin(row.Direct)}
}

// The castable aura spell: instant, on the GCD, free. The aura it turns on is one of core's
// self-cast paladin auras, which already sit in PaladinAuraCategory so one cast replaces the last.
func (paladin *Paladin) registerAuraSpell(row shared.SpellData, aura *core.Aura, classMask int64) {
	paladin.RegisterSpell(core.SpellConfig{
		ActionID:       aura.ActionID,
		SpellSchool:    row.SpellSchool,
		DefenseType:    row.DefenseType,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL | core.SpellFlagHelpful,
		ClassSpellMask: classMask,
		Rank:           row.Rank,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: row.GCD,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			aura.Activate(sim)
		},

		RelatedSelfBuff: aura,
	})
}
