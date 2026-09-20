package warrior

import (
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
)

const ShoutExpirationThreshold = time.Second * 3

type ShoutHelperConfig struct {
	ActionID           core.ActionID
	RageCost           int32
	SpellMask          int64
	ThreatBonus        float64
	AllyAuras          core.AuraArray
	ExtraCastCondition core.CanCastCondition
}

func (warrior *Warrior) MakeShoutSpellHelper(config ShoutHelperConfig) *core.Spell {
	return warrior.RegisterSpell(core.SpellConfig{
		ActionID:       config.ActionID,
		ClassSpellMask: config.SpellMask,
		SpellSchool:    core.SpellSchoolPhysical,
		Flags:          core.SpellFlagAPL | core.SpellFlagHelpful,
		ProcMask:       core.ProcMaskEmpty,

		RageCost: core.RageCostOptions{
			Cost: config.RageCost,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		ThreatMultiplier: 1,
		FlatThreatBonus:  config.ThreatBonus,

		ExtraCastCondition: config.ExtraCastCondition,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// Assuming full party, thus multiplying by 5
			spell.FlatThreatBonus = core.TernaryFloat64(sim.CurrentTime > 0, config.ThreatBonus*5/float64(sim.Environment.ActiveTargetCount()), 0)
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
			config.AllyAuras.ActivateAllPlayers(sim)
		},

		RelatedAuraArrays: config.AllyAuras.ToMap(),
	})
}

var battleShoutRank = spellData.BattleShout.HighestRank()

// TODO: The core Battle Shout aura still takes Booming Voice points, a Commanding Presence
// multiplier and two item flags. In the client Booming Voice (12321) widens the radius only,
// Commanding Presence does not exist, Solarian's Sapphire (30446) is not an item, and the shout
// itself states 139 attack power for 3 minutes (25289). Pending the shared shout aura rework, the
// talent and multiplier are passed as neutral.
func (warrior *Warrior) registerShouts() {
	commandingPresenceMultiplier := 1.0

	warrior.registerDemoralizingShout()

	battleShoutAuras := warrior.NewAllyAuraArray(func(unit *core.Unit) *core.Aura {
		aura := core.BattleShoutAura(
			warrior.GetCharacter(),
			warrior.DefaultShout != proto.WarriorShout_WarriorShoutNone,
			0,
			commandingPresenceMultiplier,
			warrior.HasBsSolarianSapphire,
			warrior.HasBsT2,
		)
		aura.BuildPhase = core.Ternary(warrior.DefaultShout == proto.WarriorShout_WarriorShoutBattle, core.CharacterBuildPhaseBuffs, core.CharacterBuildPhaseNone)
		return aura
	})

	warrior.BattleShout = warrior.MakeShoutSpellHelper(ShoutHelperConfig{
		ActionID:    core.ActionID{SpellID: battleShoutRank.SpellID},
		RageCost:    battleShoutRank.Cost,
		SpellMask:   SpellMaskBattleShout,
		ThreatBonus: 69,
		ExtraCastCondition: func(sim *core.Simulation, _ *core.Unit) bool {
			aura := battleShoutAuras.Get(&warrior.Unit)
			return !aura.IsActive() || aura.ExclusiveEffects[0].Priority <= core.GetBattleShoutValue(0, commandingPresenceMultiplier, warrior.HasBsSolarianSapphire, warrior.HasBsT2, sim.CurrentTime < 0)
		},
		AllyAuras: battleShoutAuras,
	})
}
