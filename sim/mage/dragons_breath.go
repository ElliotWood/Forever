package mage

import (
	"github.com/wowsims/forever/sim/core"
)

const dragonsBreathCoefficient = 0.1930000037

var dragonsBreathRank = spellData.DragonsBreath.BySpellID(33043)

func (mage *Mage) registerDragonsBreathSpell() {
	if !mage.Talents.DragonsBreath {
		return
	}

	mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: dragonsBreathRank.SpellID},
		SpellSchool:    dragonsBreathRank.SpellSchool,
		DefenseType:    dragonsBreathRank.DefenseType,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MageSpellDragonsBreath,

		ManaCost: core.ManaCostOptions{
			FlatCost: dragonsBreathRank.Cost,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: dragonsBreathRank.GCD,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: dragonsBreathRank.Cooldown,
			},
		},

		DamageMultiplier: 1,
		BonusCoefficient: dragonsBreathRank.Direct.BonusCoefficient(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			baseDamage := dragonsBreathRank.Direct.Damage(sim)
			spell.CalcAndDealAoeDamage(sim, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}
