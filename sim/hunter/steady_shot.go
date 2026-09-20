package hunter

// Package-level state the commented-out implementations used:
// var steadyShotRank = spellData.SteadyShot.BySpellID(34120)

// TODO: To be implemented. Spells of this name exist in the Forever client, but none of them
// has a class ability row -- no SkillLineAbility entry in a CategoryID 7 skill line -- so the
// generator has no class spell to build a ladder from.
func (hunter *Hunter) registerSteadyShotSpell() {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// hunter.SteadyShot = hunter.RegisterRangedSpell(core.SpellConfig{
	// 	ActionID:       core.ActionID{SpellID: steadyShotRank.SpellID},
	// 	SpellSchool:    steadyShotRank.SpellSchool,
	// 	DefenseType:    steadyShotRank.DefenseType,
	// 	ClassSpellMask: HunterSpellSteadyShot,
	// 	ProcMask:       core.ProcMaskRangedSpecial,
	// 	Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
	//
	// 	ManaCost: core.ManaCostOptions{
	// 		FlatCost: steadyShotRank.Cost,
	// 	},
	//
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			CastTime: time.Millisecond * 1500,
	// 		},
	// 	},
	//
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		// Steady Shot isn't affected by Ammo, Scopes or Adamantite Weightstone
	// 		weaponDamage := hunter.AutoAttacks.Ranged().BaseDamage(sim) - hunter.AmmoDamageBonus
	//
	// 		if ranged := hunter.Ranged(); ranged != nil && ranged.Enchant.EffectID == 2722 {
	// 			weaponDamage -= 10
	// 		} else if ranged != nil && ranged.Enchant.EffectID == 2723 {
	// 			weaponDamage -= 12
	// 		}
	// 		if hunter.Consumables.OhImbueId == 34340 || hunter.Consumables.OhImbueId == 29453 {
	// 			weaponDamage -= 12
	// 		}
	// 		if (hunter.Consumables.MhImbueId == 34340 || hunter.Consumables.MhImbueId == 29453) && !hunter.windFuryEnabled {
	// 			weaponDamage -= 12
	// 		}
	//
	// 		baseDamage := 0.2*spell.RangedAttackPower(target) +
	// 			weaponDamage*2.8/hunter.AutoAttacks.Ranged().SwingSpeed +
	// 			hunter.talonOfAlarBonus() +
	// 			steadyShotRank.Direct.Damage(sim)
	//
	// 		result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)
	//
	// 		spell.WaitTravelTime(sim, func(sim *core.Simulation) {
	// 			spell.DealDamage(sim, result)
	// 		})
	// 	},
	// })
}
