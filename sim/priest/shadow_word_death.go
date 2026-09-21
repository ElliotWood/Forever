package priest

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

// Shadow Word: Death is new in Forever: trained on the Shadow Magic line, four ranks, an instant
// shadow nuke on a 15 second cooldown. The generated table already carries the level-60 damage of
// every rank, which matches the beta client 1.60.1.69893 our sim reads.
//
// Not modelled, deliberately:
//   - The backlash (effect 2, 10% of max health when the target survives). It costs health, not
//     damage, and nothing in a DPS sim reads the priest's health.
//   - Effect 1, a script effect worth 150 that no tooltip names. It could be an execute bonus;
//     until a combat log says what it does, guessing would put an invented number in the DPS.
var ShadowWordDeathRankMap = spellData.ShadowWordDeath

func (priest *Priest) registerShadowWordDeathSpell(rank shared.SpellData, cdTimer *core.Timer) {
	// Early Demise (1310076): +15% crit per point against a target at or below 20% health.
	earlyDemise := spellData.EarlyDemise.ValueAt(priest.Talents.EarlyDemise)

	priest.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.SpellID},
		SpellSchool:    rank.SpellSchool,
		DefenseType:    rank.DefenseType,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: PriestSpellShadowWordDeath,
		Rank:           rank.Rank,
		MaxRange:       rank.MaxRange,

		ManaCost: core.ManaCostOptions{
			FlatCost: rank.Cost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: rank.GCD,
			},
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: rank.Cooldown,
			},
		},

		DamageMultiplier: 1,
		BonusCoefficient: rank.Direct.BonusCoefficient(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			bonus := 0.0
			if earlyDemise > 0 && sim.IsExecutePhase20() {
				bonus = earlyDemise
			}
			spell.BonusCritPercent += bonus
			spell.CalcAndDealDamage(sim, target, rank.Direct.Damage(sim), spell.OutcomeMagicHitAndCrit)
			spell.BonusCritPercent -= bonus
		},

		ExpectedInitialDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			return spell.CalcDamage(sim, target, rank.Direct.Damage(sim), spell.OutcomeExpectedMagicHitAndCrit)
		},
	})
}
