package mage

import (
	"github.com/wowsims/forever/sim/core"
)

func (mage *Mage) registerBlastWaveSpell() {
	if !mage.Talents.BlastWave {
		return
	}

	blastWaveRank := spellData.BlastWave.HighestRank()

	mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: blastWaveRank.SpellID},
		Flags:          core.SpellFlagAPL | core.SpellFlagBinary,
		SpellSchool:    blastWaveRank.SpellSchool,
		DefenseType:    blastWaveRank.DefenseType,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: MageSpellBlastWave,

		BonusCoefficient: blastWaveRank.Direct.BonusCoefficient(),
		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ManaCost: core.ManaCostOptions{
			FlatCost: blastWaveRank.Cost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: blastWaveRank.GCD,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: blastWaveRank.Cooldown,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealAoeDamage(sim, blastWaveRank.Direct.Damage(sim), spell.OutcomeMagicHitAndCrit)
		},
	})
}
