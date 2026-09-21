package paladin

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
	// "github.com/wowsims/forever/sim/core/proto" -- used only by the commented-out implementation below
)

var ExorcismRankMap = spellData.Exorcism

func (paladin *Paladin) getExorcismTimer() *core.Timer {
	if paladin.exorcismTimer == nil {
		paladin.exorcismTimer = paladin.NewTimer()
	}
	return paladin.exorcismTimer
}

// Exorcism
// https://www.wowhead.com/forever/spell=10314
//
// Causes X to Y Holy damage to an Undead or Demon target.
//
// TODO: To be implemented. spellData.Exorcism holds the six trainer ranks, 879 to 10314.
// The client carries a second ladder, 415068 to 415073, that the Season of Discovery
// passive Exorcist (415076) swaps onto the action bar so the spell can hit any target.
// The generator drops those stand-ins, and nothing in Forever teaches Exorcist, so the
// trainer ranks are the ones to build from.
func (paladin *Paladin) registerExorcism(rankConfig shared.SpellData) {
	panic("To be implemented")

	// spellID := rankConfig.SpellID
	// cost := rankConfig.Cost
	// coefficient := rankConfig.Direct.BonusCoefficient()
	//
	// paladin.RegisterSpell(core.SpellConfig{
	// 	ActionID:       core.ActionID{SpellID: spellID},
	// 	SpellSchool:    core.SpellSchoolHoly,
	// 	DefenseType:    core.DefenseTypeMagic,
	// 	ProcMask:       core.ProcMaskSpellDamage,
	// 	Flags:          core.SpellFlagAPL,
	// 	Rank:           rankConfig.Rank,
	// 	ClassSpellMask: SpellMaskExorcism,
	//
	// 	DamageMultiplier: 1,
	// 	ThreatMultiplier: 1,
	//
	// 	MaxRange: rankConfig.MaxRange,
	//
	// 	ManaCost: core.ManaCostOptions{
	// 		FlatCost: cost,
	// 	},
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			GCD: rankConfig.GCD,
	// 		},
	// 		CD: core.Cooldown{
	// 			Timer:    paladin.getExorcismTimer(),
	// 			Duration: rankConfig.Cooldown,
	// 		},
	// 	},
	//
	// 	BonusCoefficient: coefficient,
	//
	// 	ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
	// 		return target.MobType == proto.MobType_MobTypeUndead || target.MobType == proto.MobType_MobTypeDemon
	// 	},
	//
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		spell.CalcAndDealDamage(sim, target, rankConfig.Direct.Damage(sim), spell.OutcomeMagicHitAndCrit)
	// 	},
	// })
}
