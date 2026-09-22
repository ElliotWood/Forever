package warrior

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
)

func (warrior *Warrior) registerBattleShout() {
	battleShoutRank := spellData.BattleShout.HighestRank()
	baseAttackPower := battleShoutRank.Effect(shared.A_MOD_ATTACK_POWER, 0).Value
	attackPower := func() float64 {
		return baseAttackPower + core.TernaryFloat64(warrior.HasBsT2, core.BattleShoutWrathBonus, 0)
	}

	auras := warrior.NewAllyAuraArray(func(unit *core.Unit) *core.Aura {
		aura := core.BattleShoutAura(unit, true, baseAttackPower, battleShoutRank.Duration)
		aura.BuildPhase = core.Ternary(warrior.DefaultShout == proto.WarriorShout_WarriorShoutBattle, core.CharacterBuildPhaseBuffs, core.CharacterBuildPhaseNone)
		return aura.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
			if ee := aura.ExclusiveEffects[0]; ee.Priority != attackPower() {
				ee.SetPriority(sim, attackPower())
			}
		})
	})
	selfAura := auras.Get(&warrior.Unit)

	warrior.BattleShout = warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: battleShoutRank.SpellID},
		ClassSpellMask: SpellMaskBattleShout,
		SpellSchool:    battleShoutRank.SpellSchool,
		Flags:          core.SpellFlagAPL | core.SpellFlagHelpful,
		ProcMask:       core.ProcMaskEmpty,

		RageCost: core.RageCostOptions{
			Cost: battleShoutRank.Cost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: battleShoutRank.GCD,
			},
			IgnoreHaste: true,
		},

		ThreatMultiplier: 1,
		// TODO: Manual review needed -- spell 25289 carries no threat effect; none is modelled until
		// measured in game.
		FlatThreatBonus: battleShoutRank.FlatThreatBonus,

		ExtraCastCondition: func(sim *core.Simulation, _ *core.Unit) bool {
			return !selfAura.IsActive() || selfAura.ExclusiveEffects[0].Priority <= attackPower()
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
			auras.ActivateAllPlayers(sim)
		},

		RelatedAuraArrays: auras.ToMap(),
	})
}
