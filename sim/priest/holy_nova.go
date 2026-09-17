package priest

import (
	"github.com/wowsims/classic/sim/core"
)

const HolyNovaRanks = 6

// Forever beta client 1.60.1.69893. Every rank is a little below Classic's, at .107 for both halves.
var HolyNovaSpellId = [HolyNovaRanks + 1]int32{0, 15237, 15430, 15431, 27799, 27800, 27801}
var HolyNovaHealSpellId = [HolyNovaRanks + 1]int32{0, 23455, 23458, 23459, 27803, 27804, 27805}
var HolyNovaBaseDamage = [HolyNovaRanks + 1][]float64{{0}, {26, 31}, {47, 56}, {73, 84}, {103, 118}, {139, 159}, {174, 200}}
var HolyNovaBaseHealing = [HolyNovaRanks + 1][]float64{{0}, {49, 58}, {80, 90}, {111, 128}, {151, 176}, {225, 260}, {288, 334}}
var HolyNovaManaCost = [HolyNovaRanks + 1]float64{0, 185, 290, 400, 520, 635, 750}
var HolyNovaLevel = [HolyNovaRanks + 1]int{0, 20, 28, 36, 44, 52, 60}

// Only the highest rank the priest knows is registered, since Searing Light's free cast is tied to one spell.
func (priest *Priest) registerHolyNovaSpell() {
	if !priest.Talents.HolyNova {
		return
	}

	rank := 1
	for rank < HolyNovaRanks && HolyNovaLevel[rank+1] <= int(priest.Level) {
		rank++
	}

	baseDamage := HolyNovaBaseDamage[rank]
	baseHealing := HolyNovaBaseHealing[rank]
	spellCoeff := 0.107

	partyPlayers := priest.Env.Raid.GetPlayerParty(&priest.Unit).Players

	healSpell := priest.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: HolyNovaHealSpellId[rank]},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskSpellHealing,
		Flags:       core.SpellFlagHelpful | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

		DamageMultiplier: 1,
		ThreatMultiplier: 0,
		BonusCoefficient: spellCoeff,

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			for _, player := range partyPlayers {
				spell.CalcAndDealHealing(sim, &player.GetCharacter().Unit, sim.Roll(baseHealing[0], baseHealing[1]), spell.OutcomeHealingCrit)
			}
		},
	})

	priest.HolyNova = priest.RegisterSpell(core.SpellConfig{
		SpellCode:   SpellCode_PriestHolyNova,
		ActionID:    core.ActionID{SpellID: HolyNovaSpellId[rank]},
		SpellSchool: core.SpellSchoolHoly,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagPriest | core.SpellFlagAPL,

		RequiredLevel: HolyNovaLevel[rank],
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: HolyNovaManaCost[rank],
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 0,
		BonusCoefficient: spellCoeff,

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, sim.Roll(baseDamage[0], baseDamage[1]), spell.OutcomeMagicHitAndCrit)
			}

			healSpell.Cast(sim, &priest.Unit)
		},
	})
}
