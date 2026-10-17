package mage

import (
	"github.com/wowsims/forever/sim/core"
)

// Arcane Wrath (new 16-point Arcane talent, landed with client 70450): a 3.5 sec, 6 sec cooldown Arcane bolt
// that hits as if the mage had extra spell power equal to 5% of maximum mana (effect index 1's E_DUMMY 5,
// read as the percent; Blizzard dev notes t/2360696 post 6). Its damage effect's coefficient is 1.0.
func (mage *Mage) registerArcaneWrathSpell() {
	if !mage.Talents.ArcaneWrath {
		return
	}

	rank := spellData.ArcaneWrath.Highest()
	damage := rank.DamageEffect()
	manaAsSpellPower := rank.EffectN(2).Percent()

	mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.ID},
		SpellSchool:    rank.SpellSchool(),
		DefenseType:    rank.DefenseTypeCore(),
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MageSpellArcaneWrath,
		MaxRange:       float64(rank.MaxRange),
		MissileSpeed:   float64(rank.Speed),

		ManaCost: rank.ManaCost(),
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      rank.GCD(),
				CastTime: rank.CastTime(),
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: max(rank.Cooldown(), rank.CategoryCooldown()),
			},
		},

		DamageMultiplier: 1,
		BonusCoefficient: damage.Coeff(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := damage.Roll(sim, core.CharacterLevel) + spell.BonusCoefficient*manaAsSpellPower*mage.MaxMana()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
