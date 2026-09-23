package warrior

import (
	"github.com/wowsims/forever/sim/core"
)

func (warrior *Warrior) registerHeroicStrike() {
	heroicStrikeRank := spellData.HeroicStrike.HighestRank()
	heroicStrikeBaseDamage, _ := heroicStrikeRank.Direct.Range()

	warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: heroicStrikeRank.SpellID},
		SpellSchool:    heroicStrikeRank.SpellSchool,
		DefenseType:    heroicStrikeRank.DefenseType,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskHeroicStrike,
		MaxRange:       core.MaxMeleeRange,

		RageCost: core.RageCostOptions{
			Cost:   heroicStrikeRank.Cost,
			Refund: heroicStrikeRank.MissRefund(),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		// Not in the client table (no E_THREAT effect), so the Classic rank 9 value stays ours.
		FlatThreatBonus: 173,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := heroicStrikeBaseDamage + warrior.MHWeaponDamage(sim, spell.MeleeAttackPower(target))
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
}

func (warrior *Warrior) registerCleave() {
	cleaveRank := spellData.Cleave.HighestRank()
	cleaveBaseDamage, _ := cleaveRank.Direct.Range()

	const maxTargets int32 = 2

	warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: cleaveRank.SpellID},
		SpellSchool:    cleaveRank.SpellSchool,
		DefenseType:    cleaveRank.DefenseType,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskCleave,
		MaxRange:       core.MaxMeleeRange,

		RageCost: core.RageCostOptions{
			Cost:   cleaveRank.Cost,
			Refund: cleaveRank.MissRefund(),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		// Not in the client table (no E_THREAT effect), so the Classic rank 5 value stays ours.
		FlatThreatBonus: 100,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := cleaveBaseDamage + warrior.MHWeaponDamage(sim, spell.MeleeAttackPower(target))
			results := spell.CalcCleaveDamage(sim, target, maxTargets, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			spell.DealBatchedAoeDamage(sim)
			if !results[0].Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
}
