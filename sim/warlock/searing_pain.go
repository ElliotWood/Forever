package warlock

var searingPainRank = spellData.SearingPain.HighestRank()
var searingPainCoeff = searingPainRank.Direct.BonusCoefficient()

// TODO: To be implemented. Port the TBC Searing Pain implementation below; not yet verified against the Forever client.
func (warlock *Warlock) registerSearingPain() {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	//
	// warlock.Shadowburn = warlock.RegisterSpell(core.SpellConfig{
	// 	ActionID:       core.ActionID{SpellID: searingPainRank.SpellID},
	// 	SpellSchool:    searingPainRank.SpellSchool,
	// 	ProcMask:       core.ProcMaskSpellDamage,
	// 	Flags:          core.SpellFlagAPL,
	// 	ClassSpellMask: WarlockSpellSearingPain,
	// 	MaxRange:       searingPainRank.MaxRange,
	//
	// 	ManaCost: core.ManaCostOptions{FlatCost: searingPainRank.Cost},
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			GCD:      searingPainRank.GCD,
	// 			CastTime: searingPainRank.CastTime,
	// 		},
	// 	},
	//
	// 	DamageMultiplier: 1,
	// 	DefenseType:      searingPainRank.DefenseType,
	// 	ThreatMultiplier: 2,
	// 	BonusCoefficient: searingPainCoeff,
	//
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		dmgRoll := searingPainRank.Direct.Damage(sim)
	// 		spell.CalcAndDealDamage(sim, target, dmgRoll, spell.OutcomeMagicHitAndCrit)
	// 	},
	// })
}
