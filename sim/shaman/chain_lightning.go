package shaman

import (
	"time"

	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

var ChainLightningRankMap = spellData.ChainLightning

// TODO: To be implemented. Port the TBC Chain Lightning Spell implementation below; not yet verified against the Forever client.
func (shaman *Shaman) registerChainLightningSpell() {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// maxHits := min(3, shaman.Env.TotalTargetCount())
	// sharedCDTimer := shaman.NewTimer()
	// shaman.ChainLightningOverloads = make(map[int32][]*core.Spell, len(ChainLightningRankMap))
	// ChainLightningRankMap.RegisterAll(func(config shared.SpellData) {
	// 	shaman.newChainLightningSpell(config, false, sharedCDTimer)
	// 	for range maxHits {
	// 		shaman.ChainLightningOverloads[config.Rank] = append(shaman.ChainLightningOverloads[config.Rank], shaman.newChainLightningSpell(config, true, nil))
	// 	}
	// })
	//
}

func (shaman *Shaman) newChainLightningSpell(config shared.SpellData, isElementalOverload bool, sharedCDTimer *core.Timer) *core.Spell {
	shamConfig := ShamSpellConfig{
		ActionID:            core.ActionID{SpellID: config.SpellID},
		Rank:                config.Rank,
		IsElementalOverload: isElementalOverload,
		BaseFlatCost:        config.Cost,
		BonusCoefficient:    config.Direct.BonusCoefficient(),
		SpellSchool:         core.SpellSchoolNature,
		Overloads:           shaman.ChainLightningOverloads,
		BounceReduction:     0.7,
		ClassSpellMask:      core.TernaryInt64(isElementalOverload, SpellMaskChainLightningOverload, SpellMaskChainLightning),
		BaseCastTime:        time.Second * 2,
	}
	spellConfig := shaman.newElectricSpellConfig(shamConfig)
	if !isElementalOverload {
		spellConfig.Cast.CD = core.Cooldown{
			Timer:    sharedCDTimer,
			Duration: time.Second * 6,
		}
	}
	maxHits := int32(3)
	maxHits = min(maxHits, shaman.Env.TotalTargetCount())

	spellConfig.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		curTarget := target

		// Damage calculation and DealDamage are in separate loops so that e.g. a spell power proc
		// can't proc on the first target and apply to the second
		numHits := min(maxHits, shaman.Env.ActiveTargetCount())
		results := make([]*core.SpellResult, numHits)
		for hitIndex := range numHits {
			baseDamage := config.Direct.Damage(sim)
			results[hitIndex] = spell.CalcDamage(sim, curTarget, baseDamage, spell.OutcomeMagicHitAndCrit)

			curTarget = sim.Environment.NextActiveTargetUnit(curTarget)
			spell.DamageMultiplier *= shamConfig.BounceReduction
		}

		for hitIndex := range numHits {
			if !isElementalOverload && results[hitIndex].Landed() && sim.Proc(shaman.GetOverloadChance()/3, "Chain Lightning Elemental Overload") {
				shamConfig.Overloads[config.Rank][hitIndex].Cast(sim, results[hitIndex].Target)
			}
			spell.DealDamage(sim, results[hitIndex])
			spell.DamageMultiplier /= shamConfig.BounceReduction
		}
	}

	return shaman.RegisterSpell(spellConfig)
}
