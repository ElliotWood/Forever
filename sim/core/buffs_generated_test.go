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

// A flat stat row: no category at all, so the two sources of stamina add up.
func TestGeneratedFlatStatBuffsAddUp(t *testing.T) {
	char := newGeneratedBuffTestCharacter()

	applyBuffEffects(generatedBuffTestAgent{char},
		&proto.RaidBuffs{PowerWordFortitude: true}, &proto.PartyBuffs{BloodPact: true}, &proto.IndividualBuffs{})
	char.applyBuildPhaseAuras(CharacterBuildPhaseBuffs)

	if got := char.stats[stats.Stamina]; got != 124 {
		t.Errorf("Power Word: Fortitude and Blood Pact applied %v stamina, want the client's 70 + 54", got)
	}

	fortitude := char.GetAura("Power Word: Fortitude (External)")
	if fortitude == nil {
		t.Fatalf("no aura is labelled %q; the unit has %v", "Power Word: Fortitude (External)", auraLabels(char))
	}
	if want := (ActionID{SpellID: 10938, Tag: -1}); fortitude.ActionID != want {
		t.Errorf("the external copy is %v, want %v", fortitude.ActionID, want)
	}
}

// A percentage row: every stat the client's A_MOD_TOTAL_STAT_PERCENTAGE names
// goes through a multiplying dependency rather than a flat amount.
func TestGeneratedBlessingOfKingsMultipliesEveryStat(t *testing.T) {
	char := newGeneratedBuffTestCharacter()
	char.stats = stats.Stats{
		stats.Strength: 100, stats.Agility: 100, stats.Stamina: 100,
		stats.Intellect: 100, stats.Spirit: 100, stats.Armor: 100,
	}

	applyBuffEffects(generatedBuffTestAgent{char},
		&proto.RaidBuffs{}, &proto.PartyBuffs{}, &proto.IndividualBuffs{BlessingOfKings: true})
	measureGeneratedBuffStats(char)

	for _, stat := range []stats.Stat{stats.Strength, stats.Agility, stats.Stamina, stats.Intellect, stats.Spirit} {
		if got := char.stats[stat]; got != 110 {
			t.Errorf("%s is %v, want 100 multiplied by the client's 1.1", stat.StatName(), got)
		}
	}
	if got := char.stats[stats.Armor]; got != 100 {
		t.Errorf("armor is %v, want the blessing to have left it alone", got)
	}
}

// Two sources of shadow resistance compete for the school, while everything
// else Gift of the Wild grants is applied outright.
func TestGeneratedResistancesCompeteAcrossBuffs(t *testing.T) {
	char := newGeneratedBuffTestCharacter()

	applyBuffEffects(generatedBuffTestAgent{char},
		&proto.RaidBuffs{GiftOfTheWild: true, ShadowProtection: true},
		&proto.PartyBuffs{}, &proto.IndividualBuffs{})
	char.applyBuildPhaseAuras(CharacterBuildPhaseBuffs)

	if got := char.stats[stats.ShadowResistance]; got != 60 {
		t.Errorf("shadow resistance is %v, want only Shadow Protection's 60", got)
	}
	if got := char.stats[stats.FrostResistance]; got != 27 {
		t.Errorf("frost resistance is %v, want Gift of the Wild's 27", got)
	}
	if got := char.stats[stats.Armor]; got != 385 {
		t.Errorf("armor is %v, want Gift of the Wild's 385", got)
	}
	if got := char.stats[stats.Stamina]; got != 16 {
		t.Errorf("stamina is %v, want Gift of the Wild's 16", got)
	}

	school := char.ExclusiveEffectManager.GetExclusiveEffectCategory(
		ResistanceCategoryShadow + stats.ShadowResistance.StatName() + "Add")
	if len(school.effects) != 2 {
		t.Errorf("the shadow school holds %d effects, want both buffs", len(school.effects))
	}
}

// A paladin aura holds its own slot and joins the shared one that stops a
// paladin from having two of them up. Only the player's own copy joins it: the
// external copy has to be able to sit next to the one the paladin casts.
func TestGeneratedPaladinAuraJoinsTheSharedCategoryOnThePlayerCopyOnly(t *testing.T) {
	char := newGeneratedBuffTestCharacter()

	applyBuffEffects(generatedBuffTestAgent{char},
		&proto.RaidBuffs{}, &proto.PartyBuffs{DevotionAura: true}, &proto.IndividualBuffs{})
	char.applyBuildPhaseAuras(CharacterBuildPhaseBuffs)

	if got := char.stats[stats.Armor]; got != 735 {
		t.Errorf("Devotion Aura applied %v armor, want the client's 735", got)
	}

	shared := char.ExclusiveEffectManager.GetExclusiveEffectCategory("PaladinAura")
	if len(shared.effects) != 0 {
		t.Errorf("the external copy joined the shared category with %d effects, want none", len(shared.effects))
	}

	own := char.ExclusiveEffectManager.GetExclusiveEffectCategory(DevotionAuraCategory)
	if !own.SingleAura {
		t.Error("the aura's own category is not single-aura, so a second copy could sit next to it")
	}

	DevotionAuraAura(&char.Unit, true, 0)
	if len(shared.effects) != 1 {
		t.Errorf("the player's copy put %d effects in the shared category, want one", len(shared.effects))
	}
}

// A totem the client only ties to its cast by name, improved by the talent the
// live tree still prices.
func TestGeneratedManaSpringTotemTakesTheTalentedValue(t *testing.T) {
	char := newGeneratedBuffTestCharacter()

	applyBuffEffects(generatedBuffTestAgent{char}, &proto.RaidBuffs{},
		&proto.PartyBuffs{ManaSpringTotem: proto.TristateEffect_TristateEffectImproved},
		&proto.IndividualBuffs{})
	char.applyBuildPhaseAuras(CharacterBuildPhaseBuffs)

	if got := char.stats[stats.MP5]; got != 30 {
		t.Errorf("the improved totem applied %v MP5, want the curve's top value 30", got)
	}

	aura := char.GetAura("Mana Spring Totem (External)")
	if aura == nil {
		t.Fatalf("no aura is labelled %q; the unit has %v", "Mana Spring Totem (External)", auraLabels(char))
	}
	if want := (ActionID{SpellID: 10494, Tag: -1}); aura.ActionID != want {
		t.Errorf("the totem's aura is %v, want the aura family member %v", aura.ActionID, want)
	}

	category := char.ExclusiveEffectManager.GetExclusiveEffectCategory(
		ManaSpringTotemCategory + stats.MP5.StatName() + "Add")
	if len(category.effects) != 1 || category.effects[0].Priority != 30 {
		t.Errorf("the totem's category holds %d effects, first bid %v; want one bidding 30",
			len(category.effects), category.effects[0].Priority)
	}
}

// Commanding Shout is driven like its sibling: the external copy chains behind
// the player's own cast, which the driver finds by the aura's tag.
func TestPartyCommandingShoutAppliesTheGeneratedAura(t *testing.T) {
	char := newGeneratedBuffTestCharacter()

	applyBuffEffects(generatedBuffTestAgent{char},
		&proto.RaidBuffs{}, &proto.PartyBuffs{CommandingShout: true}, &proto.IndividualBuffs{})
	char.applyBuildPhaseAuras(CharacterBuildPhaseBuffs)

	if got := char.stats[stats.Stamina]; got != 42 {
		t.Errorf("the party's Commanding Shout applied %v stamina, want the client's 42", got)
	}

	aura := char.GetAura("Commanding Shout (External)")
	if aura == nil {
		t.Fatalf("no aura is labelled %q; the unit has %v", "Commanding Shout (External)", auraLabels(char))
	}
	if want := (ActionID{SpellID: 403215, Tag: -1}); aura.ActionID != want {
		t.Errorf("the external copy is %v, want %v", aura.ActionID, want)
	}
	if aura.Duration != CommandingShoutDuration(0) {
		t.Errorf("the external copy lasts %v, want %v", aura.Duration, CommandingShoutDuration(0))
	}

	category := char.ExclusiveEffectManager.GetExclusiveEffectCategory(CommandingShoutCategory)
	if !category.SingleAura {
		t.Error("the Commanding Shout category is not single-aura, so a second copy could sit next to it")
	}
	if tagged := char.GetAurasWithTag(CommandingShoutCategory); len(tagged) != 1 || tagged[0] != aura {
		t.Errorf("%d auras carry the Commanding Shout tag, want only the external copy", len(tagged))
	}
}

// The two rows whose stat the manifest names, and the one whose only amount is
// a pseudo-stat the client states as a reduction.
func TestGeneratedPartyAurasApplyTheStatsTheManifestNames(t *testing.T) {
	char := newGeneratedBuffTestCharacter()

	applyBuffEffects(generatedBuffTestAgent{char}, &proto.RaidBuffs{},
		&proto.PartyBuffs{LeaderOfThePack: true, MoonkinAura: true, ConcentrationAura: true},
		&proto.IndividualBuffs{})
	char.applyBuildPhaseAuras(CharacterBuildPhaseBuffs)

	if got := char.stats[stats.PhysicalCritPercent]; got != 3 {
		t.Errorf("Leader of the Pack applied %v physical crit, want the client's 3", got)
	}
	if got := char.stats[stats.SpellCritPercent]; got != 3 {
		t.Errorf("Moonkin Aura applied %v spell crit, want the client's 3", got)
	}
	if got := char.PseudoStats.PushbackChance; got != 0.65 {
		t.Errorf("Concentration Aura left the pushback chance at %v, want 1 reduced by the client's 35%%", got)
	}
}

// The build phase enables a multiplying dependency without recomputing the
// unit's stats; applyAllEffects measures them afterwards, and so does this.
func measureGeneratedBuffStats(char *Character) {
	char.applyBuildPhaseAuras(CharacterBuildPhaseBuffs)
	char.stats = char.SortAndApplyStatDependencies(char.stats).FloorGameStats()
}

func auraLabels(char *Character) []string {
	labels := make([]string, 0, len(char.auras))
	for _, aura := range char.auras {
		labels = append(labels, aura.Label)
	}
	return labels
}
