package warlock

import (
	"github.com/wowsims/forever/sim/core"
)

var conflagrateRank = spellData.Conflagrate.HighestRank()
var conflagrateCoeff = conflagrateRank.Direct.BonusCoefficient()

func (warlock *Warlock) registerConflagrate() {

	if !warlock.Talents.Conflagrate {
		return
	}

	warlock.Conflagrate = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: conflagrateRank.SpellID},
		SpellSchool:    conflagrateRank.SpellSchool,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellConflagrate,

		ManaCost: core.ManaCostOptions{FlatCost: conflagrateRank.Cost},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: conflagrateRank.GCD,
			},
			CD: core.Cooldown{
				Duration: conflagrateRank.Cooldown,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warlock.Immolate.Dot(target).IsActive()
		},

		DamageMultiplier: 1.0,
		DefenseType:      conflagrateRank.DefenseType,
		ThreatMultiplier: 1,
		BonusCoefficient: conflagrateCoeff,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			//tie this to landed/hit
			dmgRoll := conflagrateRank.Direct.Damage(sim)
			result := spell.CalcAndDealDamage(sim, target, dmgRoll, spell.OutcomeMagicHitAndCrit)

			if result.Landed() || result.DidResist() {
				warlock.Immolate.Dot(target).Deactivate(sim)
			}
		},
	})
}
