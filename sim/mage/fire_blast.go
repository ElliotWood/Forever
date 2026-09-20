package mage

// TODO: To be implemented. Fire Blast exists in Forever with a full seven-rank ladder
// (spellData.FireBlast, 2136/2137/2138/8412/8413/10197/10199). The implementation below is
// the TBC one already ported onto that ladder -- HighestRank() in place of the level-70
// rank it pinned -- so bringing it back is uncommenting it. It stays commented until
// it has been reviewed.
//
// Imports and rank pin the implementation needs, kept with it:
// import (
// 	"github.com/wowsims/forever/sim/core"
// )
//
// var fireBlastRank = spellData.FireBlast.HighestRank()

func (mage *Mage) registerFireBlastSpell() {
	panic("To be implemented")

	//
	// mage.FireBlast = mage.RegisterSpell(core.SpellConfig{
	// 	ActionID:       core.ActionID{SpellID: fireBlastRank.SpellID},
	// 	SpellSchool:    fireBlastRank.SpellSchool,
	// 	DefenseType:    fireBlastRank.DefenseType,
	// 	ProcMask:       core.ProcMaskSpellDamage,
	// 	Flags:          core.SpellFlagAPL,
	// 	ClassSpellMask: MageSpellFireBlast,
	//
	// 	ManaCost: core.ManaCostOptions{
	// 		FlatCost: fireBlastRank.Cost,
	// 	},
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			GCD: fireBlastRank.GCD,
	// 		},
	// 		CD: core.Cooldown{
	// 			Timer:    mage.NewTimer(),
	// 			Duration: fireBlastRank.Cooldown,
	// 		},
	// 	},
	//
	// 	DamageMultiplier: 1,
	// 	BonusCoefficient: fireBlastRank.Direct.BonusCoefficient(),
	//
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		baseDamage := fireBlastRank.Direct.Damage(sim)
	// 		spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
	// 	},
	// })
}
