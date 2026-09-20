package mage

// Package-level state the commented-out implementations used:
// var dragonsBreathRank = spellData.DragonsBreath.BySpellID(33043)

const dragonsBreathCoefficient = 0.1930000037

// TODO: uncalled -- Forever drops the Dragon's Breath talent; re-gate before wiring
// back into registerSpells.
// TODO: To be implemented. Spells of this name exist in the Forever client, but none of them
// has a class ability row -- no SkillLineAbility entry in a CategoryID 7 skill line -- so the
// generator has no class spell to build a ladder from.
func (mage *Mage) registerDragonsBreathSpell() {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// mage.RegisterSpell(core.SpellConfig{
	// 	ActionID:       core.ActionID{SpellID: dragonsBreathRank.SpellID},
	// 	SpellSchool:    dragonsBreathRank.SpellSchool,
	// 	DefenseType:    dragonsBreathRank.DefenseType,
	// 	ProcMask:       core.ProcMaskSpellDamage,
	// 	Flags:          core.SpellFlagAPL,
	// 	ClassSpellMask: MageSpellDragonsBreath,
	//
	// 	ManaCost: core.ManaCostOptions{
	// 		FlatCost: dragonsBreathRank.Cost,
	// 	},
	//
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			GCD: dragonsBreathRank.GCD,
	// 		},
	// 		CD: core.Cooldown{
	// 			Timer:    mage.NewTimer(),
	// 			Duration: dragonsBreathRank.Cooldown,
	// 		},
	// 	},
	//
	// 	DamageMultiplier: 1,
	// 	BonusCoefficient: dragonsBreathRank.Direct.BonusCoefficient(),
	// 	ThreatMultiplier: 1,
	//
	// 	ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
	// 		baseDamage := dragonsBreathRank.Direct.Damage(sim)
	// 		spell.CalcAndDealAoeDamage(sim, baseDamage, spell.OutcomeMagicHitAndCrit)
	// 	},
	// })
}
