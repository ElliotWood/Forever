package priest

import (
	"github.com/wowsims/classic/sim/core"
)

func (priest *Priest) registerHolyNovaSpell() {
	if !priest.Talents.HolyNova {
		return
	}

	// Only the rank 1 tooltip was shown in the demo, so only rank 1 is registered.
	// TODO: beta will confirm the higher ranks and the mana cost.
	baseDamage := []float64{30, 35}
	baseHealing := []float64{64, 73}

	partyPlayers := priest.Env.Raid.GetPlayerParty(&priest.Unit).Players

	healSpell := priest.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 23455},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskSpellHealing,
		Flags:       core.SpellFlagHelpful | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

		DamageMultiplier: 1,
		ThreatMultiplier: 0,
		BonusCoefficient: 0.286,

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			for _, player := range partyPlayers {
				spell.CalcAndDealHealing(sim, &player.GetCharacter().Unit, sim.Roll(baseHealing[0], baseHealing[1]), spell.OutcomeHealingCrit)
			}
		},
	})

	priest.HolyNova = priest.RegisterSpell(core.SpellConfig{
		SpellCode:   SpellCode_PriestHolyNova,
		ActionID:    core.ActionID{SpellID: 15237},
		SpellSchool: core.SpellSchoolHoly,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagPriest | core.SpellFlagAPL,

		RequiredLevel: 20,
		Rank:          1,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.22,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 0,
		BonusCoefficient: 0.143,

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, sim.Roll(baseDamage[0], baseDamage[1]), spell.OutcomeMagicHitAndCrit)
			}

			healSpell.Cast(sim, &priest.Unit)
		},
	})
}
