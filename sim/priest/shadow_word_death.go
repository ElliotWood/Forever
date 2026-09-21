package priest

import (
	"time"

	"github.com/wowsims/classic/sim/core"
)

// Shadow Word: Death is new in Forever: trained on the Shadow Magic line (SkillLineAbility
// AcquireMethod 0), four ranks, an instant shadow nuke on a 15 second cooldown. Every number
// below is read from the beta client 1.60.1.69893 - SpellEffect, SpellLevels, SpellPower,
// SpellCooldowns. It was missing from this sim entirely until the upstream harvest turned up
// its generated table in wowsims/forever, which was then checked against the client here.
//
// Not modelled, deliberately:
//   - The backlash (effect 2, 10% of max health when the target survives). It costs health,
//     not damage, and nothing in a DPS sim reads the priest's health.
//   - Effect 1, a script effect worth 150 that no tooltip names. It could be an execute bonus;
//     until a combat log says what it does, guessing would put an invented number in the
//     shadow priest's DPS.
const ShadowWordDeathRanks = 4

var ShadowWordDeathSpellId = [ShadowWordDeathRanks + 1]int32{0, 1309595, 1309633, 1309635, 1309636}
var ShadowWordDeathLevel = [ShadowWordDeathRanks + 1]int{0, 32, 40, 48, 56}
var ShadowWordDeathMaxLevel = [ShadowWordDeathRanks + 1]int{0, 37, 45, 53, 61}
var ShadowWordDeathBaseDamage = [ShadowWordDeathRanks + 1]float64{0, 295, 370, 403, 448}
var ShadowWordDeathPerLevel = [ShadowWordDeathRanks + 1]float64{0, 1.5, 1.9, 2.2, 2.5}
var ShadowWordDeathManaCost = [ShadowWordDeathRanks + 1]float64{0, 175, 205, 250, 340}

const ShadowWordDeathSpellCoef = 0.429

// Early Demise (1310076): +15% crit per rank on Shadow Word: Death against a target at or
// below 20% health.
const earlyDemiseCritPerRank = 15.0

func (priest *Priest) registerShadowWordDeath() {
	if !priest.Env.IsForever() {
		return
	}
	priest.ShadowWordDeath = make([]*core.Spell, ShadowWordDeathRanks+1)
	cdTimer := priest.NewTimer()

	for rank := 1; rank <= ShadowWordDeathRanks; rank++ {
		if ShadowWordDeathLevel[rank] <= int(priest.Level) {
			priest.ShadowWordDeath[rank] = priest.GetOrRegisterSpell(priest.shadowWordDeathConfig(rank, cdTimer))
		}
	}
}

func (priest *Priest) shadowWordDeathConfig(rank int, cdTimer *core.Timer) core.SpellConfig {
	level := min(int(priest.Level), ShadowWordDeathMaxLevel[rank])
	baseDamage := ShadowWordDeathBaseDamage[rank] + ShadowWordDeathPerLevel[rank]*float64(level-ShadowWordDeathLevel[rank])
	earlyDemise := earlyDemiseCritPerRank * float64(priest.Talents.EarlyDemise) * core.SpellCritRatingPerCritChance

	return core.SpellConfig{
		SpellCode:   SpellCode_PriestShadowWordDeath,
		ActionID:    core.ActionID{SpellID: ShadowWordDeathSpellId[rank]},
		SpellSchool: core.SpellSchoolShadow,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagPriest | core.SpellFlagAPL,

		RequiredLevel: ShadowWordDeathLevel[rank],
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: ShadowWordDeathManaCost[rank],
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: time.Second * 15,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: ShadowWordDeathSpellCoef,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			bonus := 0.0
			if earlyDemise > 0 && sim.IsExecutePhase20() {
				bonus = earlyDemise
			}
			spell.BonusCritRating += bonus
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.BonusCritRating -= bonus

			if result.Landed() {
				priest.AddShadowWeavingStack(sim)
			}
			spell.DealDamage(sim, result)
		},

		ExpectedInitialDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			return spell.CalcDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicHitAndCrit)
		},
	}
}
