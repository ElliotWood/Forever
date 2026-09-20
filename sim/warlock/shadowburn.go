package warlock

var shadowBurnRank = spellData.Shadowburn.HighestRank()
var shadowBurnCoeff = shadowBurnRank.Direct.BonusCoefficient()

// TODO: To be implemented. Port the TBC Shadow Burn implementation below; not yet verified against the Forever client.
func (warlock *Warlock) registerShadowBurn() {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	//
	// warlock.Shadowburn = warlock.RegisterSpell(core.SpellConfig{
	// 	ActionID:       core.ActionID{SpellID: shadowBurnRank.SpellID},
	// 	SpellSchool:    shadowBurnRank.SpellSchool,
	// 	ProcMask:       core.ProcMaskSpellDamage,
	// 	Flags:          core.SpellFlagAPL | core.SpellFlagBinary,
	// 	ClassSpellMask: WarlockSpellShadowBurn,
	//
	// 	ManaCost: core.ManaCostOptions{FlatCost: shadowBurnRank.Cost},
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			GCD: shadowBurnRank.GCD,
	// 		},
	// 		CD: core.Cooldown{
	// 			Timer:    warlock.NewTimer(),
	// 			Duration: shadowBurnRank.Cooldown,
	// 		},
	// 	},
	//
	// 	DamageMultiplier: 1,
	// 	DefenseType:      shadowBurnRank.DefenseType,
	// 	ThreatMultiplier: 1,
	// 	BonusCoefficient: shadowBurnCoeff,
	//
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		dmgRoll := shadowBurnRank.Direct.Damage(sim)
	// 		spell.CalcAndDealDamage(sim, target, dmgRoll, spell.OutcomeMagicHitAndCrit)
	//
	// 	},
	// })
}
