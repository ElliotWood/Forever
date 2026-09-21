package hunter

import (
	"github.com/wowsims/forever/sim/core"
)

// Strider Kick from the beta client (1317257): 100% normalized melee weapon damage, 8 sec cooldown,
// 5.81% of base mana. The generated row states no cost, so the percentage is kept from our
// client-verified sim.
func (hunter *Hunter) registerStriderKickSpell() {
	if !hunter.Talents.StriderKick {
		return
	}

	rank := spellData.StriderKick.HighestRank()

	hunter.StriderKick = hunter.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.SpellID},
		SpellSchool:    rank.SpellSchool,
		DefenseType:    rank.DefenseType,
		ClassSpellMask: HunterSpellStriderKick,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		MaxRange:       rank.MaxRange,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 5.81,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: rank.GCD,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: rank.Cooldown,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := hunter.AutoAttacks.MH().CalculateNormalizedWeaponDamage(sim, spell.MeleeAttackPower(target))
			spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
		},
	})
}
