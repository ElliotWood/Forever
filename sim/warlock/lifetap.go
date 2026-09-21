package warlock

import (
	"github.com/wowsims/forever/sim/core"
)

// Life Tap turns health into mana one for one, so the restore is the damage the spell rolls against
// the warlock. The row carries the 424 on its energize effect and no coefficient of its own; the
// 0.8 is ours. Improved Life Tap rides on the talent as a SpellMod, Demonic Energies hands the pet
// a share of the restore (the talent's second effect, 50% per point).
func (warlock *Warlock) registerLifeTap() {
	rank := spellData.LifeTap.HighestRank()
	actionID := core.ActionID{SpellID: rank.SpellID}
	baseDamage := spellData.LifeTap.EffectAt(0).ValueAt(rank.Rank)
	petManaShare := spellData.DemonicEnergies.EffectAt(1).FractionAt(warlock.Talents.DemonicEnergies)

	manaMetrics := warlock.NewManaMetrics(actionID)
	petManaMetrics := make(map[*WarlockPet]*core.ResourceMetrics, len(warlock.BasePets))
	for _, pet := range warlock.BasePets {
		petManaMetrics[pet] = pet.NewManaMetrics(actionID)
	}

	warlock.LifeTap = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    rank.SpellSchool,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL | core.SpellFlagBinary,
		ClassSpellMask: WarlockSpellLifeTap,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplierAdditive: 1,
		DamageMultiplier:         1,
		ThreatMultiplier:         1,
		BonusCoefficient:         0.8,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcDamage(sim, &warlock.Unit, baseDamage, spell.OutcomeAlwaysHit)
			warlock.RemoveHealth(sim, result.Damage)

			warlock.AddMana(sim, result.Damage, manaMetrics)

			if petManaShare > 0 && warlock.ActivePet != nil {
				warlock.ActivePet.AddMana(sim, result.Damage*petManaShare, petManaMetrics[warlock.ActivePet])
			}
		},
	})
}
