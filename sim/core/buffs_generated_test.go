package core

import (
	"testing"

	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

// applyBuffEffects takes an Agent, and the part of one a buff reaches is the
// Character it wraps.
type generatedBuffTestAgent struct {
	*Character
}

func (agent generatedBuffTestAgent) GetCharacter() *Character       { return agent.Character }
func (agent generatedBuffTestAgent) Initialize()                    {}
func (agent generatedBuffTestAgent) ApplyTalents()                  {}
func (agent generatedBuffTestAgent) Reset(_ *Simulation)            {}
func (agent generatedBuffTestAgent) OnEncounterStart(_ *Simulation) {}

func newGeneratedBuffTestCharacter() *Character {
	char := &Character{
		Unit: Unit{
			Type:        PlayerUnit,
			Level:       60,
			auraTracker: newAuraTracker(),
			Env:         &Environment{},
		},
	}
	char.PseudoStats = stats.NewPseudoStats()
	return char
}

func TestPartyBattleShoutAppliesTheGeneratedAura(t *testing.T) {
	char := newGeneratedBuffTestCharacter()

	applyBuffEffects(generatedBuffTestAgent{char},
		&proto.RaidBuffs{}, &proto.PartyBuffs{BattleShout: true}, &proto.IndividualBuffs{})
	char.applyBuildPhaseAuras(CharacterBuildPhaseBuffs)

	if got := char.stats[stats.AttackPower]; got != 139 {
		t.Errorf("the party's Battle Shout applied %v attack power, want the client's 139", got)
	}

	aura := char.GetAura("Battle Shout (External)")
	if aura == nil {
		t.Fatalf("no aura is labelled %q; the unit has %v", "Battle Shout (External)", auraLabels(char))
	}
	if want := (ActionID{SpellID: 25289, Tag: -1}); aura.ActionID != want {
		t.Errorf("the external copy is %v, want %v", aura.ActionID, want)
	}
	if aura.Duration != BattleShoutDuration(0) {
		t.Errorf("the external copy lasts %v, want %v", aura.Duration, BattleShoutDuration(0))
	}

	category := char.ExclusiveEffectManager.GetExclusiveEffectCategory(BattleShoutCategory)
	if !category.SingleAura {
		t.Error("the Battle Shout category is not single-aura, so a second copy could sit next to it")
	}
	if len(category.effects) != 1 || category.effects[0].Priority != 139 {
		t.Errorf("the category holds %d effects, first bid %v; want one bidding 139",
			len(category.effects), category.effects[0].Priority)
	}

	// The chain the driver installs finds the player's own shout by this tag.
	if tagged := char.GetAurasWithTag(BattleShoutCategory); len(tagged) != 1 || tagged[0] != aura {
		t.Errorf("%d auras carry the Battle Shout tag, want only the external copy", len(tagged))
	}
}

// A warrior who casts Battle Shout and has the external one ticked: both copies
// are 139 attack power, and both are equally long, so the tie goes to whichever
// the build phase activates last. ExclusiveEffect.Activate only turns a newcomer
// away when the incumbent's remaining duration is longer than the newcomer's, so
// the player's copy takes the category and the external one is deactivated - the
// character sheet shows 139 either way.
func TestPlayerBattleShoutTakesTheCategoryOnATie(t *testing.T) {
	char := newGeneratedBuffTestCharacter()

	applyBuffEffects(generatedBuffTestAgent{char},
		&proto.RaidBuffs{}, &proto.PartyBuffs{BattleShout: true}, &proto.IndividualBuffs{})

	player := BattleShoutAura(&char.Unit, true, 0)
	if want := (ActionID{SpellID: 25289, Tag: 0}); player.ActionID != want {
		t.Errorf("the player's copy is %v, want %v", player.ActionID, want)
	}
	if player.Label != "Battle Shout (Player)" {
		t.Errorf("the player's copy is labelled %q, want %q", player.Label, "Battle Shout (Player)")
	}
	// What a warrior whose DefaultShout is Battle does with its own aura.
	player.BuildPhase = CharacterBuildPhaseBuffs

	char.applyBuildPhaseAuras(CharacterBuildPhaseBuffs)

	external := char.GetAura("Battle Shout (External)")
	if !player.IsActive() {
		t.Error("the player's own shout did not activate next to the external one")
	}
	if external.IsActive() {
		t.Error("the external copy stayed active, so both copies of the shout are on the character")
	}
	if got := char.stats[stats.AttackPower]; got != 139 {
		t.Errorf("both copies together applied %v attack power, want 139", got)
	}
	if tagged := char.GetAurasWithTag(BattleShoutCategory); len(tagged) != 2 {
		t.Errorf("%d auras carry the Battle Shout tag, want both copies", len(tagged))
	}
}

func auraLabels(char *Character) []string {
	labels := make([]string, 0, len(char.auras))
	for _, aura := range char.auras {
		labels = append(labels, aura.Label)
	}
	return labels
}
