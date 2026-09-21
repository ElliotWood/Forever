package shaman

import (
	"time"

	"github.com/wowsims/classic/sim/core"
)

const ChainLightningRanks = 4
const ChainLightningTargetCount = int32(3)

// Forever beta client values, scaled to level 60 the same way as Lightning Bolt's. Rank 3's 0.517 is what the
// client stores, against 0.571 on the ranks either side; it reads like a transposed digit but is kept as written.
var ChainLightningSpellId = [ChainLightningRanks + 1]int32{0, 421, 930, 2860, 10605}
var ChainLightningBaseDamage = [ChainLightningRanks + 1][]float64{{0}, {85, 97}, {97, 109}, {109, 122}, {119, 134}}
var ChainLightningSpellCoef = [ChainLightningRanks + 1]float64{0, .571, .571, .517, .571}
var ChainLightningManaCost = [ChainLightningRanks + 1]float64{0, 225, 305, 390, 485}
var ChainLightningLevel = [ChainLightningRanks + 1]int{0, 32, 40, 48, 56}

func (shaman *Shaman) registerChainLightningSpell() {
	shaman.ChainLightning = make([]*core.Spell, ChainLightningRanks+1)
	shaman.ChainLightningOverload = make([]*core.Spell, ChainLightningRanks+1)

	cdTimer := shaman.NewTimer()

	for rank := 1; rank <= ChainLightningRanks; rank++ {
		config := shaman.newChainLightningSpellConfig(rank, cdTimer)

		if config.RequiredLevel <= int(shaman.Level) {
			// The overload gets a config of its own so that the two casts share no bounce results.
			shaman.ChainLightningOverload[rank] = shaman.registerOverloadSpell(shaman.newChainLightningSpellConfig(rank, cdTimer))
			shaman.ChainLightning[rank] = shaman.RegisterSpell(config)
		}
	}
}

func (shaman *Shaman) newChainLightningSpellConfig(rank int, cdTimer *core.Timer) core.SpellConfig {
	spellId := ChainLightningSpellId[rank]
	baseDamageLow := ChainLightningBaseDamage[rank][0]
	baseDamageHigh := ChainLightningBaseDamage[rank][1]
	spellCoeff := ChainLightningSpellCoef[rank]
	manaCost := ChainLightningManaCost[rank]
	level := ChainLightningLevel[rank]

	cooldown := time.Second * 6
	castTime := time.Millisecond * 2000

	shaman.ChainLightningBounceCoefficient = .70 // 30% reduction per bounce
	targetCount := ChainLightningTargetCount

	spell := shaman.newElectricSpellConfig(
		core.ActionID{SpellID: spellId},
		manaCost,
		castTime,
	)

	spell.SpellCode = SpellCode_ShamanChainLightning
	spell.ClassSpellMask = SpellMaskChainLightning
	spell.RequiredLevel = level
	spell.Rank = rank
	spell.BonusCoefficient = spellCoeff
	spell.Cast.CD = core.Cooldown{
		Timer:    cdTimer,
		Duration: cooldown,
	}

	results := make([]*core.SpellResult, min(targetCount, shaman.Env.GetNumTargets()))

	spell.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		primaryTarget := target
		origMult := spell.DamageMultiplier
		for hitIndex := range results {
			baseDamage := sim.Roll(baseDamageLow, baseDamageHigh)
			results[hitIndex] = spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			target = sim.Environment.NextTargetUnit(target)
			spell.DamageMultiplier *= shaman.ChainLightningBounceCoefficient
		}

		for _, result := range results {
			spell.DealDamage(sim, result)
		}

		spell.DamageMultiplier = origMult

		shaman.tryLightningOverload(sim, primaryTarget, spell, shaman.ChainLightningOverload[rank])
	}

	return spell
}
