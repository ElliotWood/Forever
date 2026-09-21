package hunter

import (
	"github.com/wowsims/forever/sim/core"
)

func (hunter *Hunter) registerMongooseBiteSpell() {
	rank := spellData.MongooseBite.HighestRank()
	baseDamage := rank.Direct.Damage

	// The aura is only a pre-requisite for Mongoose Bite: a dodge opens the window, and Expose Prey
	// opens it off any landed hit on a marked target.
	defensiveWindow := spellData.ExposePreyTriggered.ByRank(1)
	hunter.DefensiveState = hunter.RegisterAura(core.Aura{
		Label:    "Defensive State",
		ActionID: core.ActionID{SpellID: defensiveWindow.SpellID},
		Duration: defensiveWindow.Duration,

		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.DidDodge() {
				aura.Activate(sim)
			}
		},
	})

	hunter.MongooseBite = hunter.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.SpellID},
		SpellSchool:    rank.SpellSchool,
		DefenseType:    rank.DefenseType,
		ClassSpellMask: HunterSpellMongooseBite,
		ProcMask:       core.ProcMaskMeleeSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		MaxRange:       rank.MaxRange,

		ManaCost: core.ManaCostOptions{
			FlatCost: rank.Cost,
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
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hunter.DefensiveState.IsActive()
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			hunter.DefensiveState.Deactivate(sim)
			// Forever: normalized melee weapon damage plus a smaller flat amount, where Classic
			// dealt the flat amount alone.
			damage := baseDamage(sim) + hunter.AutoAttacks.MH().CalculateNormalizedWeaponDamage(sim, spell.MeleeAttackPower(target))
			result := spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if hunter.LaceratingStrikes != nil && result.Landed() {
				hunter.procLaceratingStrikes(sim, result)
			}
		},
	})
}
