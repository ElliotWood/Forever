package dps

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
)

// The golden suites are skipped while the class is stubbed, so this is what
// checks that registerShouts still runs and that the warrior's own Battle Shout
// is the generated aura.
func TestWarriorCastsTheGeneratedBattleShout(t *testing.T) {
	player := &proto.Player{
		Name:          "Battle Shout",
		Race:          proto.Race_RaceOrc,
		Class:         proto.Class_ClassWarrior,
		Equipment:     &proto.EquipmentSpec{},
		TalentsString: DefaultFuryTalents,
		Spec:          DefaultOptions,
	}

	env, _, _ := core.NewEnvironment(
		core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		core.MakeSingleTargetEncounter(0), false, true)

	war := env.Raid.Parties[0].Players[0].(*DpsWarrior)
	if war.BattleShout == nil {
		t.Fatal("registerShouts did not register Battle Shout")
	}
	if want := (core.ActionID{SpellID: 25289}); war.BattleShout.ActionID != want {
		t.Errorf("the warrior casts %v, want %v", war.BattleShout.ActionID, want)
	}

	aura := war.GetAura("Battle Shout (Player)")
	if aura == nil {
		t.Fatal("the ally aura array holds no aura labelled \"Battle Shout (Player)\"")
	}
	if want := (core.ActionID{SpellID: 25289, Tag: 0}); aura.ActionID != want {
		t.Errorf("the warrior's own copy is %v, want %v", aura.ActionID, want)
	}
	if aura.Duration != core.BattleShoutDuration(war.Talents.BoomingVoice) {
		t.Errorf("the warrior's own copy lasts %v, want %v",
			aura.Duration, core.BattleShoutDuration(war.Talents.BoomingVoice))
	}
	if got := aura.ExclusiveEffects[0].Priority; got != core.BattleShoutValue(war.Talents.BoomingVoice) {
		t.Errorf("the warrior's own copy bids %v, want %v",
			got, core.BattleShoutValue(war.Talents.BoomingVoice))
	}
}
