package warrior

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
)

// Battlegear of Wrath's three pieces (23563) add 30 attack power to Battle Shout.
const battleShoutWrathBonus = 30.0

func (warrior *Warrior) registerBattleShout() {
	battleShoutRank := spellData.BattleShout.HighestRank()
	// Rank 7 (25289): 139 attack power for 3 min. core.BattleShoutAura still holds TBC's 306 for
	// 2 min (the raid buff path, which their #39 rewrites), so the warrior's own shout overrides the
	// value and the duration with its row.
	baseAttackPower := battleShoutRank.Effect(shared.A_MOD_ATTACK_POWER, 0).Value
	attackPower := func() float64 {
		return baseAttackPower + core.TernaryFloat64(warrior.HasBsT2, battleShoutWrathBonus, 0)
	}

	auras := warrior.NewAllyAuraArray(func(unit *core.Unit) *core.Aura {
		agent := warrior.Env.Raid.GetPlayerFromUnit(unit)
		if agent == nil {
			return nil
		}
		aura := core.BattleShoutAura(agent.GetCharacter(), true, 0, 1, false, false)
		aura.ActionID = core.ActionID{SpellID: battleShoutRank.SpellID}
		aura.Duration = battleShoutRank.Duration
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
		// TODO: spell 25289 carries no threat effect; none is modelled until measured in game.
		FlatThreatBonus: battleShoutRank.FlatThreatBonus,

		// Battle Shout is a single-aura exclusive category: a stronger one from another source (the
		// party buff) blocks ours, so casting would only burn rage every GCD.
		ExtraCastCondition: func(sim *core.Simulation, _ *core.Unit) bool {
			active := selfAura.ExclusiveEffects[0].Category.GetActiveEffect()
			return active == nil || active.Priority <= attackPower()
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
			// The exclusive check runs before OnGain sets the value, so give it the real one first.
			for _, aura := range auras {
				if aura != nil {
					aura.ExclusiveEffects[0].SetPriority(sim, attackPower())
				}
			}
			auras.ActivateAllPlayers(sim)
		},

		RelatedAuraArrays: auras.ToMap(),
	})
}
