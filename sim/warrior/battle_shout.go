package warrior

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/proto"
)

var battleShoutRank = spellData.BattleShout.Highest()
var battleShoutAttackPower = battleShoutRank.Effect(dbcenums.A_MOD_ATTACK_POWER, 0).Average(core.CharacterLevel)

func (warrior *Warrior) battleShoutValue() float64 {
	return battleShoutAttackPower + core.TernaryFloat64(warrior.HasBsT2, core.BattleShoutWrathBonus, 0)
}

func (warrior *Warrior) registerBattleShout() {
	auras := warrior.NewAllyAuraArray(func(unit *core.Unit) *core.Aura {
		aura := core.BattleShoutAura(unit, true, battleShoutAttackPower, battleShoutRank.Duration())
		aura.BuildPhase = core.Ternary(warrior.DefaultShout == proto.WarriorShout_WarriorShoutBattle, core.CharacterBuildPhaseBuffs, core.CharacterBuildPhaseNone)
		return aura.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
			aura.ExclusiveEffects[0].SetPriority(sim, warrior.battleShoutValue())
		})
	})

	warrior.BattleShout = warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: battleShoutRank.ID},
		ClassSpellMask: SpellMaskBattleShout,
		ClassFlags:     SpellFlagsBattleShout,
		SpellSchool:    battleShoutRank.SpellSchool(),
		Flags:          core.SpellFlagAPL | core.SpellFlagHelpful,
		ProcMask:       core.ProcMaskEmpty,

		RageCost: core.RageCostOptions{
			Cost: rageCost(battleShoutRank),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: battleShoutRank.GCD(),
			},
			IgnoreHaste: true,
		},

		ThreatMultiplier: 1,
		// TODO: Manual review needed -- spell 25289 carries no threat effect; none is modelled until
		// measured in game.
		FlatThreatBonus: 0,

		ExtraCastCondition: func(sim *core.Simulation, _ *core.Unit) bool {
			aura := auras.Get(&warrior.Unit)
			return !aura.IsActive() || aura.ExclusiveEffects[0].Priority <= warrior.battleShoutValue()
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
			auras.ActivateAllPlayers(sim)
		},

		RelatedAuraArrays: auras.ToMap(),
	})
}
