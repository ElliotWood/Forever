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

// The warrior's own shout with three pieces of Battlegear of Wrath on: the set
// is a class option rather than the equipped gear, and it raises what the
// shout bids for its category by the 30 spell 23563 states.
func TestWarriorShoutsForTheTierTwoBonus(t *testing.T) {
	bidOf := func(hasT2 bool) float64 {
		t.Helper()

		player := &proto.Player{
			Name:          "Enhanced Battle Shout",
			Race:          proto.Race_RaceOrc,
			Class:         proto.Class_ClassWarrior,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: DefaultFuryTalents,
			Spec: &proto.Player_DpsWarrior{
				DpsWarrior: &proto.DpsWarrior{
					Options: &proto.DpsWarrior_Options{
						ClassOptions: &proto.WarriorOptions{
							DefaultShout:  proto.WarriorShout_WarriorShoutBattle,
							DefaultStance: proto.WarriorStance_WarriorStanceBerserker,
							HasBsT2:       hasT2,
						},
					},
				},
			},
		}

		env, _, _ := core.NewEnvironment(
			core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
			core.MakeSingleTargetEncounter(0), false, true)

		war := env.Raid.Parties[0].Players[0].(*DpsWarrior)
		aura := war.GetAura("Battle Shout (Player)")
		if aura == nil {
			t.Fatal("the ally aura array holds no aura labelled \"Battle Shout (Player)\"")
		}
		return aura.ExclusiveEffects[0].Priority
	}

	bare, withSet := bidOf(false), bidOf(true)

	if want := core.BattleShoutValue(0); bare != want {
		t.Errorf("a warrior without the set bids %v, want the client's %v", bare, want)
	}
	if want := bare + core.BattleShoutT2Bonus; withSet != want {
		t.Errorf("a warrior wearing the set bids %v, want %v", withSet, want)
	}
}

// The ally aura array registers one "Battle Shout (Player)" per unit, which two
// warriors in the same party share, so the set has to be worth 30 once rather
// than once per warrior wearing it.
func TestTwoWarriorsWearingTheSetShoutForOneBonus(t *testing.T) {
	warriorProto := func(name string) *proto.Player {
		return &proto.Player{
			Name:          name,
			Race:          proto.Race_RaceOrc,
			Class:         proto.Class_ClassWarrior,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: DefaultFuryTalents,
			Spec: &proto.Player_DpsWarrior{
				DpsWarrior: &proto.DpsWarrior{
					Options: &proto.DpsWarrior_Options{
						ClassOptions: &proto.WarriorOptions{
							DefaultShout:  proto.WarriorShout_WarriorShoutBattle,
							DefaultStance: proto.WarriorStance_WarriorStanceBerserker,
							HasBsT2:       true,
						},
					},
				},
			},
		}
	}

	raid := &proto.Raid{
		Parties: []*proto.Party{{
			Players: []*proto.Player{warriorProto("First"), warriorProto("Second")},
			Buffs:   &proto.PartyBuffs{},
		}},
		Buffs:            &proto.RaidBuffs{},
		Debuffs:          &proto.Debuffs{},
		NumActiveParties: 1,
	}

	env, _, _ := core.NewEnvironment(raid, core.MakeSingleTargetEncounter(0), false, true)

	want := core.BattleShoutValue(0) + core.BattleShoutT2Bonus
	for _, agent := range env.Raid.Parties[0].Players {
		war := agent.(*DpsWarrior)
		aura := war.GetAura("Battle Shout (Player)")
		if aura == nil {
			t.Fatalf("%s carries no aura labelled \"Battle Shout (Player)\"", war.Label)
		}
		if got := aura.ExclusiveEffects[0].Priority; got != want {
			t.Errorf("%s reads a shout worth %v, want %v", war.Label, got, want)
		}
	}
}

// What the warrior does when something else holds the Battle Shout category.
// The buff is the unit's, so a 169 shout from the party or from the warrior
// standing next to this one is the same category, and a 139 shout the category
// would turn away is a global and some rage thrown away.
func TestWarriorShoutsOnlyWhenTheCategoryWillTakeIt(t *testing.T) {
	warriorProto := func(name string, hasT2 bool) *proto.Player {
		return &proto.Player{
			Name:          name,
			Race:          proto.Race_RaceOrc,
			Class:         proto.Class_ClassWarrior,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: DefaultFuryTalents,
			Spec: &proto.Player_DpsWarrior{
				DpsWarrior: &proto.DpsWarrior{
					Options: &proto.DpsWarrior_Options{
						ClassOptions: &proto.WarriorOptions{
							DefaultShout:  proto.WarriorShout_WarriorShoutBattle,
							DefaultStance: proto.WarriorStance_WarriorStanceBerserker,
							HasBsT2:       hasT2,
						},
					},
				},
			},
		}
	}

	build := func(party *proto.PartyBuffs, players ...*proto.Player) *core.Environment {
		t.Helper()
		raid := &proto.Raid{
			Parties:          []*proto.Party{{Players: players, Buffs: party}},
			Buffs:            &proto.RaidBuffs{},
			Debuffs:          &proto.Debuffs{},
			NumActiveParties: 1,
		}
		env, _, _ := core.NewEnvironment(raid, core.MakeSingleTargetEncounter(0), false, true)
		return env
	}

	withSet := core.BattleShoutValue(0) + core.BattleShoutT2Bonus

	t.Run("an external shout that snapshots the set outbids a warrior without it", func(t *testing.T) {
		sim := &core.Simulation{}
		env := build(&proto.PartyBuffs{BattleShout: true, SnapshotBsT2: true}, warriorProto("Bare", false))
		war := env.Raid.Parties[0].Players[0].(*DpsWarrior)

		external := war.GetAura("Battle Shout (External)")
		if external == nil {
			t.Fatal("the party's copy is not registered")
		}
		external.Activate(sim)
		if !external.IsActive() || external.ExclusiveEffects[0].Priority != withSet {
			t.Fatalf("the party's copy is active %v at %v, want it up at %v",
				external.IsActive(), external.ExclusiveEffects[0].Priority, withSet)
		}

		if war.BattleShout.ExtraCastCondition(sim, &war.Unit) {
			t.Error("the warrior would shout for 139 while a 169 copy holds the category, wasting the global")
		}
	})

	t.Run("a warrior keeps refreshing the copy that holds the category", func(t *testing.T) {
		sim := &core.Simulation{}
		env := build(&proto.PartyBuffs{}, warriorProto("Set", true), warriorProto("Bare", false))
		bare := env.Raid.Parties[0].Players[1].(*DpsWarrior)

		own := bare.GetAura("Battle Shout (Player)")
		if own == nil {
			t.Fatal("the warrior carries no aura labelled \"Battle Shout (Player)\"")
		}
		own.Activate(sim)
		if !own.IsActive() || own.ExclusiveEffects[0].Priority != withSet {
			t.Fatalf("the shared copy is active %v at %v, want it up at the set warrior's %v",
				own.IsActive(), own.ExclusiveEffects[0].Priority, withSet)
		}

		if !bare.BattleShout.ExtraCastCondition(sim, &bare.Unit) {
			t.Error("the warrior stopped refreshing the copy it is keeping up, because the set made it bid more")
		}
	})

	t.Run("nothing holding the category is always worth a shout", func(t *testing.T) {
		sim := &core.Simulation{}
		env := build(&proto.PartyBuffs{}, warriorProto("Bare", false))
		war := env.Raid.Parties[0].Players[0].(*DpsWarrior)

		if !war.BattleShout.ExtraCastCondition(sim, &war.Unit) {
			t.Error("the warrior would not shout with nothing holding the category")
		}
	})
}
