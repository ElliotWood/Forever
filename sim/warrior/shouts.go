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

var battleShoutRank = spellData.BattleShout.BySpellID(25289)

// TODO: Manual review needed -- this was modelled during the Forever port, not carried
// over unchanged, so its numbers and shape want checking against the client.
func (warrior *Warrior) registerShouts() {
	warrior.registerDemoralizingShout()

	battleShoutAuras := warrior.NewAllyAuraArray(func(unit *core.Unit) *core.Aura {
		// Booming Voice modifies the radius only, so the generated aura ignores the
		// points; they are passed for signature uniformity.
		aura := core.BattleShoutAura(unit, warrior.DefaultShout != proto.WarriorShout_WarriorShoutNone, warrior.Talents.BoomingVoice)
		aura.BuildPhase = core.Ternary(warrior.DefaultShout == proto.WarriorShout_WarriorShoutBattle, core.CharacterBuildPhaseBuffs, core.CharacterBuildPhaseNone)
		return aura
	})

	warrior.BattleShout = warrior.MakeShoutSpellHelper(ShoutHelperConfig{
		ActionID:    core.ActionID{SpellID: battleShoutRank.SpellID},
		RageCost:    battleShoutRank.Cost,
		SpellMask:   SpellMaskBattleShout,
		ThreatBonus: 69,
		ExtraCastCondition: func(_ *core.Simulation, _ *core.Unit) bool {
			aura := battleShoutAuras.Get(&warrior.Unit)
			return !aura.IsActive() || aura.ExclusiveEffects[0].Priority <= core.BattleShoutValue(warrior.Talents.BoomingVoice)
		},
		AllyAuras: battleShoutAuras,
	})

	// TODO: To be implemented. The Forever client ships no rank ladder the generator can
	// read for this ability -- it survives as a single spell with no "Rank N" subtext and
	// no ranked SkillLineAbility row -- it has a class ability row but no "Rank N" subtext,
	// so no ladder can be built for it. Left unregistered (warrior.CommandingShout stays nil)
	// rather than panicking here, since this function also builds Battle Shout, which every
	// warrior needs working.
}
