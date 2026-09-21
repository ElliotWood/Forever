package warrior

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/spelldata"
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

	config := spelldata.SpellConfig(&warrior.Unit, battleShoutRank, spelldata.Flags(core.SpellFlagAPL))
	config.ClassSpellMask = SpellMaskBattleShout
	config.ProcMask = core.ProcMaskEmpty
	config.ThreatMultiplier = 1
	// TODO: Manual review needed -- spell 25289 carries no threat effect; none is modelled until
	// measured in game.
	config.FlatThreatBonus = 0

	config.ExtraCastCondition = func(sim *core.Simulation, _ *core.Unit) bool {
		aura := auras.Get(&warrior.Unit)
		return !aura.IsActive() || aura.ExclusiveEffects[0].Priority <= warrior.battleShoutValue()
	}

	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
		auras.ActivateAllPlayers(sim)
	}

	config.RelatedAuraArrays = auras.ToMap()

	warrior.BattleShout = warrior.RegisterSpell(config)
}
