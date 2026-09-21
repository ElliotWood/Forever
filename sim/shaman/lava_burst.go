package shaman

import (
	"github.com/wowsims/forever/sim/core"
)

var lavaBurstRank = spellData.LavaBurst.HighestRank()

// Lava Burst is new in Forever. Its second effect is the bonus it gains against a target already
// burning with Flame Shock.
func (shaman *Shaman) registerLavaBurstSpell() {
	if !shaman.Talents.LavaBurst {
		return
	}

	flameShockBonus := spellData.LavaBurst.EffectAt(1).MultiplierAt(lavaBurstRank.Rank)

	shaman.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: lavaBurstRank.SpellID},
		SpellSchool:    lavaBurstRank.SpellSchool,
		DefenseType:    lavaBurstRank.DefenseType,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL | SpellFlagShamanSpell | SpellFlagFocusable,
		ClassSpellMask: SpellMaskLavaBurst,
		MissileSpeed:   lavaBurstRank.MissileSpeed,

		ManaCost: core.ManaCostOptions{
			FlatCost: lavaBurstRank.Cost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      lavaBurstRank.GCD,
				CastTime: lavaBurstRank.CastTime,
			},
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: lavaBurstRank.Cooldown,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: lavaBurstRank.Direct.BonusCoefficient(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if shaman.FlameShock.RelatedDotSpell.Dot(target).IsActive() {
				spell.DamageMultiplier *= flameShockBonus
				defer func() { spell.DamageMultiplier /= flameShockBonus }()
			}

			result := spell.CalcDamage(sim, target, lavaBurstRank.Direct.Damage(sim), spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
