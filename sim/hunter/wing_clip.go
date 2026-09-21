package hunter

import (
	"github.com/wowsims/forever/sim/core"
)

func (hunter *Hunter) registerWingClipSpell() {
	rank := spellData.WingClip.HighestRank()
	// Effect 0 is the snare; the damage is effect 1.
	baseDamage := spellData.WingClip.EffectAt(1).ValueAt(rank.Rank)

	hunter.WingClip = hunter.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.SpellID},
		SpellSchool:    rank.SpellSchool,
		DefenseType:    rank.DefenseType,
		ClassSpellMask: HunterSpellWingClip,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL | core.SpellFlagBinary,
		MaxRange:       rank.MaxRange,

		ManaCost: core.ManaCostOptions{
			FlatCost: rank.Cost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: rank.GCD,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
		},
	})
}
