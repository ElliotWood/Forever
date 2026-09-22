package warrior

import (
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
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

// The manifest's battle_shout row resolves to this rank too, so the warrior's own shout and the
// party's copy state one spell.
var battleShoutRank = spellData.BattleShout.BySpellID(25289)

// TODO: Manual review needed -- this was modelled during the Forever port, not carried
// over unchanged, so its numbers and shape want checking against the client.
func (warrior *Warrior) registerShouts() {
	warrior.registerDemoralizingShout()

	// A warrior that shouts builds a copy of its own; one that shouts nothing
	// gets the isPlayer=false constructor, whose aura is the party's external
	// copy, and what that copy is worth is the party's Battle Shout state's to say.
	castsOwnShout := warrior.DefaultShout != proto.WarriorShout_WarriorShoutNone

	// Three pieces of Battlegear of Wrath add a flat 30 to the shout. HasBsT2 is
	// the user saying this warrior wears them; the equipped set is not read. The
	// set is worth the 30 on the shout this warrior makes, so it raises what that
	// copy applies and what the cast is worth to the category together.
	battleShoutBase := core.BattleShoutValue(warrior.Talents.BoomingVoice)
	battleShoutValue := battleShoutBase
	shoutsWithTheSet := castsOwnShout && warrior.HasBsT2
	if shoutsWithTheSet {
		battleShoutValue += core.BattleShoutT2Bonus
	}

	battleShoutAuras := warrior.NewAllyAuraArray(func(unit *core.Unit) *core.Aura {
		// Booming Voice modifies the radius only, so the generated aura ignores the
		// points; they are passed for signature uniformity.
		aura := core.BattleShoutAura(unit, castsOwnShout, warrior.Talents.BoomingVoice)
		if shoutsWithTheSet {
			core.AddGeneratedFlatBonus(aura, stats.AttackPower, battleShoutBase, core.BattleShoutT2Bonus)
		}
		aura.BuildPhase = core.Ternary(warrior.DefaultShout == proto.WarriorShout_WarriorShoutBattle, core.CharacterBuildPhaseBuffs, core.CharacterBuildPhaseNone)
		return aura
	})

	// What holds the Battle Shout category decides whether shouting is worth a
	// global. Nothing there and the cast puts the buff up. This warrior's own
	// copy there and the cast only refreshes it, which is worth a global once it
	// is about to run out. Anything else has to be outbid first - a warrior
	// without the set would otherwise keep recasting a 139 shout the category
	// turns away while an external 169 one is up.
	battleShoutCategory := warrior.GetExclusiveEffectCategory(core.BattleShoutCategory)

	warrior.BattleShout = warrior.MakeShoutSpellHelper(ShoutHelperConfig{
		ActionID:    core.ActionID{SpellID: battleShoutRank.SpellID},
		RageCost:    battleShoutRank.Cost,
		SpellMask:   SpellMaskBattleShout,
		ThreatBonus: 69,
		ExtraCastCondition: func(sim *core.Simulation, _ *core.Unit) bool {
			active := battleShoutCategory.GetActiveEffect()
			if active == nil {
				return true
			}
			if active.Aura == battleShoutAuras.Get(&warrior.Unit) {
				return active.Aura.RemainingDuration(sim) <= ShoutExpirationThreshold
			}
			return battleShoutValue >= active.Priority
		},
		AllyAuras: battleShoutAuras,
	})
}
