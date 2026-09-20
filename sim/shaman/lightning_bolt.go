package shaman

import (
	"time"

	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

var LightningBoltRankMap = spellData.LightningBolt

// TODO: To be implemented. Port the TBC Lightning Bolt Spell implementation below; not yet verified against the Forever client.
func (shaman *Shaman) registerLightningBoltSpell() {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// shaman.LightningBoltOverloads = make(map[int32]*core.Spell, len(LightningBoltRankMap))
	// LightningBoltRankMap.RegisterAll(func(config shared.SpellData) {
	// 	shaman.RegisterSpell(shaman.newLightningBoltSpellConfig(config, false))
	// 	shaman.LightningBoltOverloads[config.Rank] = shaman.RegisterSpell(shaman.newLightningBoltSpellConfig(config, true))
	// })
}

func (shaman *Shaman) newLightningBoltSpellConfig(config shared.SpellData, isElementalOverload bool) core.SpellConfig {
	shamConfig := ShamSpellConfig{
		ActionID:            core.ActionID{SpellID: config.SpellID},
		Rank:                config.Rank,
		IsElementalOverload: isElementalOverload,
		BaseFlatCost:        config.Cost,
		BonusCoefficient:    config.Direct.BonusCoefficient(),
		BaseCastTime:        time.Millisecond * 2500,
	}
	spellConfig := shaman.newElectricSpellConfig(shamConfig)

	spellConfig.ClassSpellMask = core.TernaryInt64(isElementalOverload, SpellMaskLightningBoltOverload, SpellMaskLightningBolt)
	spellConfig.MissileSpeed = 20

	spellConfig.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		baseDamage := config.Direct.Damage(sim)
		result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

		spell.WaitTravelTime(sim, func(sim *core.Simulation) {
			if !isElementalOverload && result.Landed() && sim.Proc(shaman.GetOverloadChance(), "Lightning Bolt Elemental Overload") {
				shaman.LightningBoltOverloads[config.Rank].Cast(sim, target)
			}

			spell.DealDamage(sim, result)
		})
	}

	return spellConfig
}
