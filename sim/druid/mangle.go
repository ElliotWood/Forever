package druid

import (
	"github.com/wowsims/forever/sim/core"
)

var mangleRank = spellData.Mangle.HighestRank()

// Forever ships ONE Mangle - 407995 and 1238069/1238070/1238073 on the Feral Combat line, all with
// ShapeshiftMask [144,0], which is Bear and Dire Bear only. The TBC Cat/Bear split is gone with it,
// so this is the only Mangle registrar and the name still says "Bear" because that is the form it
// is restricted to. Its effects are weapon damage and a flat bonus only (1238073): none of TBC's
// bleed debuff.
func (druid *Druid) registerMangleBearSpell() {
	if !druid.Talents.Mangle {
		return
	}

	druid.MangleAuras = druid.NewEnemyAuraArray(core.MangleAura)

	druid.MangleBear = druid.RegisterSpell(Bear, core.SpellConfig{
		ActionID:       core.ActionID{SpellID: mangleRank.SpellID},
		SpellSchool:    mangleRank.SpellSchool,
		DefenseType:    mangleRank.DefenseType,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		ClassSpellMask: DruidSpellMangleBear,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		Rank:           mangleRank.Rank,

		RageCost: core.RageCostOptions{
			Cost:   mangleRank.Cost,
			Refund: mangleRank.MissRefund(),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: mangleRank.GCD,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: mangleRank.Cooldown,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1.5,
		MaxRange:         core.MaxMeleeRange,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := mangleRank.Direct.Damage(sim) + spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower(target))
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if !result.Landed() {
				spell.IssueRefund(sim)
			}

			// Berserk removes Mangle's cooldown (client 417141).
			if druid.BerserkAura.IsActive() {
				spell.CD.Reset()
			}
		},
	})
}
