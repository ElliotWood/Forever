package warlock

import (
	"github.com/wowsims/forever/sim/core"
)

var deathCoilRank = spellData.DeathCoil.HighestRank()

func (warlock *Warlock) registerDeathCoil() {

	warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: deathCoilRank.SpellID},
		SpellSchool:    deathCoilRank.SpellSchool,
		DefenseType:    deathCoilRank.DefenseType,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellDeathCoil,
		MissileSpeed:   deathCoilRank.MissileSpeed,
		MaxRange:       deathCoilRank.MaxRange,

		ManaCost: core.ManaCostOptions{FlatCost: deathCoilRank.Cost},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: deathCoilRank.GCD,
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: deathCoilRank.Cooldown,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 0.214,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, 526, spell.OutcomeMagicHit)
		},
	})
}
