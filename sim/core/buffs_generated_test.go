package core

import (
	"testing"
	"time"

	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
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
	if BattleShoutDuration(0) != 3*time.Minute {
		t.Errorf("Battle Shout lasts %v, want the client's 3 minutes", BattleShoutDuration(0))
	}
	if aura.Duration != BattleShoutDuration(0) {
		t.Errorf("the external copy lasts %v, want %v", aura.Duration, BattleShoutDuration(0))
	}

	category := char.ExclusiveEffectManager.GetExclusiveEffectCategory(BattleShoutCategory)
	if !category.SingleAura {
		t.Error("the Battle Shout category is not single-aura, so a second copy could sit next to it")
	}
	if len(category.effects) != 1 {
		t.Errorf("the category holds %d effects, want one", len(category.effects))
	} else if category.effects[0].Priority != 139 {
		t.Errorf("the only effect in the category bids %v, want 139", category.effects[0].Priority)
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

// The party's snapshot flag says the warrior who shouted for the party wears
// three pieces of Battlegear of Wrath, so the external copy is worth the
// client's 139 plus the set's 30.
func TestPartyBattleShoutSnapshotsTheTierTwoBonus(t *testing.T) {
	for _, row := range []struct {
		name  string
		party *proto.PartyBuffs
		want  float64
	}{
		{"with the set", &proto.PartyBuffs{BattleShout: true, SnapshotBsT2: true}, 169},
		{"without it", &proto.PartyBuffs{BattleShout: true}, 139},
	} {
		t.Run(row.name, func(t *testing.T) {
			char := newGeneratedBuffTestCharacter()

			applyBuffEffects(generatedBuffTestAgent{char}, &proto.RaidBuffs{}, row.party, &proto.IndividualBuffs{})
			char.applyBuildPhaseAuras(CharacterBuildPhaseBuffs)

			if got := char.stats[stats.AttackPower]; got != row.want {
				t.Errorf("the party's Battle Shout applied %v attack power, want %v", got, row.want)
			}
			if got := char.GetAura("Battle Shout (External)").ExclusiveEffects[0].Priority; got != row.want {
				t.Errorf("the external copy bids %v, want %v", got, row.want)
			}
		})
	}
}

// Two copies of the shout that are no longer worth the same: the stronger one
// takes the category whichever side it is on, and the character sheet reads it
// once rather than both added up.
func TestTheStrongerBattleShoutTakesTheCategory(t *testing.T) {
	for _, row := range []struct {
		name         string
		party        *proto.PartyBuffs
		playerBonus  float64
		wantPlayerUp bool
	}{
		{"the player wears the set", &proto.PartyBuffs{BattleShout: true}, BattleShoutT2Bonus, true},
		{"the external warrior does", &proto.PartyBuffs{BattleShout: true, SnapshotBsT2: true}, 0, false},
	} {
		t.Run(row.name, func(t *testing.T) {
			char := newGeneratedBuffTestCharacter()

			applyBuffEffects(generatedBuffTestAgent{char}, &proto.RaidBuffs{}, row.party, &proto.IndividualBuffs{})

			player := BattleShoutAura(&char.Unit, true, 0)
			if row.playerBonus != 0 {
				AddGeneratedFlatBonus(player, stats.AttackPower, row.playerBonus)
			}
			player.BuildPhase = CharacterBuildPhaseBuffs

			char.applyBuildPhaseAuras(CharacterBuildPhaseBuffs)

			external := char.GetAura("Battle Shout (External)")
			if player.IsActive() != row.wantPlayerUp {
				t.Errorf("the player's copy is active: %v, want %v", player.IsActive(), row.wantPlayerUp)
			}
			if external.IsActive() == row.wantPlayerUp {
				t.Errorf("the external copy is active: %v, want %v", external.IsActive(), !row.wantPlayerUp)
			}
			if got := char.stats[stats.AttackPower]; got != 169 {
				t.Errorf("both copies together applied %v attack power, want the stronger one's 169", got)
			}
		})
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

	shared := char.ExclusiveEffectManager.GetExclusiveEffectCategory("PaladinAura")
	if len(shared.effects) != 0 {
		t.Errorf("the external copy joined the shared category with %d effects, want none", len(shared.effects))
	}

	player := DevotionAuraAura(&char.Unit, true, 0)
	if len(shared.effects) != 1 {
		t.Errorf("the player's copy put %d effects in the shared category, want one", len(shared.effects))
	}
	player.BuildPhase = CharacterBuildPhaseBuffs

	measureGeneratedBuffStats(char)

	own := char.ExclusiveEffectManager.GetExclusiveEffectCategory(DevotionAuraCategory)
	if !own.SingleAura {
		t.Error("the aura's own category is not single-aura, so a second copy could sit next to it")
	}

	// Both copies bid the same 735 and both are permanent, so the paladin's own
	// cast takes the slot by activating after the external one, and the armor on
	// the character sheet does not move.
	external := char.GetAura("Devotion Aura (External)")
	if !player.IsActive() {
		t.Error("the paladin's own aura did not activate next to the external one")
	}
	if external.IsActive() {
		t.Error("the external copy stayed active, so both copies of the aura are on the character")
	}
	if got := char.stats[stats.Armor]; got != 735 {
		t.Errorf("both copies together applied %v armor, want the client's 735", got)
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

// The two rows whose stats the manifest names. Each is worth the client's 3 on
// every kind of crit, and neither touches RangedCritPercent, which the sim adds
// on top of the physical one, so a ranged attack gains 3 and not 6.
func TestGeneratedCritAurasApplyTheClientsThreeToEveryKindOfCrit(t *testing.T) {
	for _, row := range []struct {
		name  string
		party *proto.PartyBuffs
	}{
		{"Leader of the Pack", &proto.PartyBuffs{LeaderOfThePack: true}},
		{"Moonkin Aura", &proto.PartyBuffs{MoonkinAura: true}},
	} {
		t.Run(row.name, func(t *testing.T) {
			char := newGeneratedBuffTestCharacter()

			applyBuffEffects(generatedBuffTestAgent{char}, &proto.RaidBuffs{}, row.party, &proto.IndividualBuffs{})
			char.applyBuildPhaseAuras(CharacterBuildPhaseBuffs)

			if got := char.stats[stats.PhysicalCritPercent]; got != 3 {
				t.Errorf("%s applied %v physical crit, want the client's 3", row.name, got)
			}
			if got := char.stats[stats.SpellCritPercent]; got != 3 {
				t.Errorf("%s applied %v spell crit, want the client's 3", row.name, got)
			}
			if got := char.stats[stats.RangedCritPercent]; got != 0 {
				t.Errorf("%s applied %v on top of the physical crit a ranged attack already reads", row.name, got)
			}
		})
	}
}

// The row whose only amount is a pseudo-stat the client states as a reduction.
func TestGeneratedPartyConcentrationAuraReducesPushback(t *testing.T) {
	char := newGeneratedBuffTestCharacter()

	applyBuffEffects(generatedBuffTestAgent{char}, &proto.RaidBuffs{},
		&proto.PartyBuffs{ConcentrationAura: true}, &proto.IndividualBuffs{})
	char.applyBuildPhaseAuras(CharacterBuildPhaseBuffs)

	if got := char.PseudoStats.PushbackChance; got != 0.65 {
		t.Errorf("Concentration Aura left the pushback chance at %v, want 1 reduced by the client's 35%%", got)
	}
}

// Rank 5 of Trueshot Aura, which the resolver takes, states one
// A_MOD_RANGED_ATTACK_POWER effect worth 50: melee attack power is untouched and
// the 75 of rank 4 is not what the party gets.
func TestGeneratedTrueshotAuraGivesTheTopRanksRangedAttackPower(t *testing.T) {
	char := newGeneratedBuffTestCharacter()

	applyBuffEffects(generatedBuffTestAgent{char},
		&proto.RaidBuffs{}, &proto.PartyBuffs{TrueshotAura: true}, &proto.IndividualBuffs{})
	char.applyBuildPhaseAuras(CharacterBuildPhaseBuffs)

	if got := char.stats[stats.RangedAttackPower]; got != 50 {
		t.Errorf("Trueshot Aura applied %v ranged attack power, want the client's 50", got)
	}
	if got := char.stats[stats.AttackPower]; got != 0 {
		t.Errorf("Trueshot Aura applied %v melee attack power, want none", got)
	}

	aura := char.GetAura("Trueshot Aura (External)")
	if aura == nil {
		t.Fatalf("no aura is labelled %q; the unit has %v", "Trueshot Aura (External)", auraLabels(char))
	}
	if want := (ActionID{SpellID: 20906, Tag: -1}); aura.ActionID != want {
		t.Errorf("the external copy is %v, want %v", aura.ActionID, want)
	}
	if aura.Duration != NeverExpires {
		t.Errorf("the party's Trueshot Aura lasts %v, want it permanent", aura.Duration)
	}
}

// The build phase enables a multiplying dependency without recomputing the
// unit's stats; applyAllEffects measures them afterwards, and so does this.
func measureGeneratedBuffStats(char *Character) {
	char.applyBuildPhaseAuras(CharacterBuildPhaseBuffs)
	char.stats = char.SortAndApplyStatDependencies(char.stats).FloorGameStats()
}

// Every buff a pet is given or denied, field by field, for a pet that is out
// when the fight starts and one that is summoned during it.
func TestGeneratedPetBuffsStripExactlyTheRowsThePolicyNames(t *testing.T) {
	ticked := func() (*proto.RaidBuffs, *proto.PartyBuffs, *proto.IndividualBuffs) {
		return &proto.RaidBuffs{
				Thorns: true, ArcaneBrilliance: true, DivineSpirit: true,
				GiftOfTheWild: true, PowerWordFortitude: true, ShadowProtection: true,
			},
			&proto.PartyBuffs{
				WindfuryTotem: true, BattleShout: true, GraceOfAirTotem: true,
				MoonkinAura: true, LeaderOfThePack: true,
				AtieshMage: 2,
			},
			&proto.IndividualBuffs{
				Innervates: 1, PowerInfusions: 1, ShadowPriestDps: 500,
				BlessingOfKings: true, BlessingOfMight: true, BlessingOfWisdom: true,
				BlessingOfSalvation: true,
			}
	}

	owner := newGeneratedBuffTestCharacter()
	pet := &Pet{Character: *newGeneratedBuffTestCharacter(), Owner: owner, enabledOnStart: true}

	raid, party, individual := ticked()
	applyGeneratedPetBuffs(pet, raid, party, individual)

	// Stripped whenever the pet is out: Thorns is not given to it, it cannot
	// gain an extra attack, and nobody spends a cooldown on a pet.
	if raid.Thorns || party.WindfuryTotem ||
		individual.Innervates != 0 || individual.PowerInfusions != 0 {
		t.Errorf("a pet out from the start kept %v, %v, %v", raid, party, individual)
	}

	// Everything else the party grants is still the pet's.
	if !party.BattleShout || !party.GraceOfAirTotem || !party.MoonkinAura ||
		!party.LeaderOfThePack || party.AtieshMage != 2 {
		t.Errorf("a pet out from the start lost a buff no policy strips: %v, %v", party, individual)
	}
	// A pet that was there from the start keeps every targeted buff.
	if !raid.ArcaneBrilliance || !raid.DivineSpirit || !raid.GiftOfTheWild ||
		!raid.PowerWordFortitude || !raid.ShadowProtection || !individual.BlessingOfKings ||
		!individual.BlessingOfMight || !individual.BlessingOfWisdom ||
		!individual.BlessingOfSalvation || individual.ShadowPriestDps != 500 {
		t.Errorf("a pet out from the start lost a targeted buff: %v, %v", raid, individual)
	}

	late := &Pet{Character: *newGeneratedBuffTestCharacter(), Owner: owner}
	raid, party, individual = ticked()
	applyGeneratedPetBuffs(late, raid, party, individual)

	if raid.ArcaneBrilliance || raid.DivineSpirit || raid.GiftOfTheWild ||
		raid.PowerWordFortitude || raid.ShadowProtection {
		t.Errorf("a pet summoned late kept a raid buff cast at the pull: %v", raid)
	}
	if individual.BlessingOfKings || individual.BlessingOfMight || individual.BlessingOfWisdom ||
		individual.BlessingOfSalvation || individual.ShadowPriestDps != 0 {
		t.Errorf("a pet summoned late kept a targeted individual buff: %v", individual)
	}
	if !party.BattleShout || !party.GraceOfAirTotem {
		t.Error("a pet summoned late lost a party aura, which no policy strips")
	}
}

// Grace of Air is the row whose uptime another field decides: with totem
// twisting on it is 9 seconds long and re-cast every 10, and without it the
// totem simply stands.
func TestGeneratedGraceOfAirFollowsTotemTwisting(t *testing.T) {
	standing := newGeneratedBuffTestCharacter()
	applyBuffEffects(generatedBuffTestAgent{standing},
		&proto.RaidBuffs{}, &proto.PartyBuffs{GraceOfAirTotem: true}, &proto.IndividualBuffs{})
	standing.applyBuildPhaseAuras(CharacterBuildPhaseBuffs)

	if got := standing.stats[stats.Agility]; got != 89 {
		t.Errorf("the totem applied %v agility, want the client's 89", got)
	}
	aura := standing.GetAura("Grace of Air Totem (External)")
	if aura == nil {
		t.Fatalf("no aura is labelled %q; the unit has %v", "Grace of Air Totem (External)", auraLabels(standing))
	}
	if aura.Duration != NeverExpires {
		t.Errorf("the totem's aura lasts %v, want it to stand for the fight", aura.Duration)
	}

	twisting := newGeneratedBuffTestCharacter()
	applyBuffEffects(generatedBuffTestAgent{twisting}, &proto.RaidBuffs{},
		&proto.PartyBuffs{GraceOfAirTotem: true, TotemTwisting: true}, &proto.IndividualBuffs{})
	twisting.applyBuildPhaseAuras(CharacterBuildPhaseBuffs)

	if got := twisting.stats[stats.Agility]; got != 89 {
		t.Errorf("the twisted totem applied %v agility, want the client's 89", got)
	}
	twisted := twisting.GetAura("Grace of Air Totem (External)")
	if twisted.Duration != time.Second*9 {
		t.Errorf("the twisted totem's aura lasts %v, want 9 seconds of every 10", twisted.Duration)
	}

	// The schedule itself: the first cast lands one totem cycle in, not at the
	// pull, so the aura is down for the opening 10 seconds.
	sim := setupFakeSimWithBuffs(&proto.RaidBuffs{},
		&proto.PartyBuffs{GraceOfAirTotem: true, TotemTwisting: true}, &proto.IndividualBuffs{})
	scheduled := sim.Raid.Parties[0].Players[0].GetCharacter().GetAura("Grace of Air Totem (External)")
	if scheduled == nil {
		t.Fatal("the twisting shaman's totem registered no aura")
	}
	if scheduled.IsActive() {
		t.Error("the twisted totem is up at the pull, want the first cast one cycle in")
	}

	for sim.CurrentTime < time.Second*10 && !scheduled.IsActive() {
		sim.Step()
	}
	if !scheduled.IsActive() {
		t.Errorf("the twisted totem is still down at %v, want it up after the first 10 second tick",
			sim.CurrentTime)
	}
	if sim.CurrentTime != time.Second*10 {
		t.Errorf("the first cast landed at %v, want 10 seconds in", sim.CurrentTime)
	}

	standingSim := setupFakeSimWithBuffs(&proto.RaidBuffs{},
		&proto.PartyBuffs{GraceOfAirTotem: true}, &proto.IndividualBuffs{})
	permanent := standingSim.Raid.Parties[0].Players[0].GetCharacter().GetAura("Grace of Air Totem (External)")
	if !permanent.IsActive() || permanent.Duration != NeverExpires {
		t.Errorf("without twisting the totem is active %v for %v, want it up from the pull for the fight",
			permanent.IsActive(), permanent.Duration)
	}
}

// Each staff in the party is worth its own copy of the aura's amounts: 11 MP5
// per druid staff, 2 spell crit per mage one, 62 healing per priest one and
// 33 spell damage plus 33 healing per warlock one.
func TestGeneratedAtieshStavesCountTheStaves(t *testing.T) {
	char := newGeneratedBuffTestCharacter()

	applyBuffEffects(generatedBuffTestAgent{char}, &proto.RaidBuffs{},
		&proto.PartyBuffs{AtieshDruid: 1, AtieshMage: 3, AtieshPriest: 2, AtieshWarlock: 1},
		&proto.IndividualBuffs{})
	char.applyBuildPhaseAuras(CharacterBuildPhaseBuffs)

	for stat, want := range map[stats.Stat]float64{
		stats.MP5:              11,
		stats.SpellCritPercent: 6,
		stats.HealingPower:     62*2 + 33,
		stats.SpellDamage:      33,
	} {
		if got := char.stats[stat]; got != want {
			t.Errorf("the party's staves applied %v %s, want %v", got, stat.StatName(), want)
		}
	}

	if char.GetAura("Atiesh - Mage (External)") == nil {
		t.Fatalf("the unit has %v, want an aura for the mage's staff", auraLabels(char))
	}
}

// The windfury proc grants the attack power the client states for the second
// the aura lasts; the totem aura around it holds the category and the trigger.
func TestGeneratedWindfuryTotemProcAppliesTheClientsAttackPower(t *testing.T) {
	sim := setupFakeSimWithBuffs(&proto.RaidBuffs{}, &proto.PartyBuffs{WindfuryTotem: true}, &proto.IndividualBuffs{})
	char := sim.Raid.Parties[0].Players[0].GetCharacter()

	proc := char.GetAura("Windfury Totem (External)")
	if proc == nil {
		t.Fatalf("no aura is labelled %q; the unit has %v", "Windfury Totem (External)", auraLabels(char))
	}
	if proc.Duration != time.Second {
		t.Errorf("the proc buff lasts %v, want the client's 1 second", proc.Duration)
	}

	before := char.stats[stats.AttackPower]
	proc.Activate(sim)
	if got := char.stats[stats.AttackPower] - before; got != 246 {
		t.Errorf("the proc applied %v attack power, want the client's 246", got)
	}
	proc.Deactivate(sim)
	if got := char.stats[stats.AttackPower]; got != before {
		t.Errorf("the proc left %v attack power behind when it expired", got-before)
	}

	totem := char.GetAura("Windfury Totem")
	if totem == nil {
		t.Fatalf("no aura is labelled %q; the unit has %v", "Windfury Totem", auraLabels(char))
	}
	category := char.ExclusiveEffectManager.GetExclusiveEffectCategory(WindfuryTotemCategory)
	if len(category.effects) != 1 || category.effects[0].Priority != WindfuryTotemValue(0) {
		t.Errorf("the category holds %d effects, first bid %v; want one bidding 246",
			len(category.effects), category.effects[0].Priority)
	}
	if char.GetAura("Windfury Totem Trigger") == nil {
		t.Error("the driver registered no proc trigger for the totem")
	}
}

// A generated damage shield deals the client's damage back to whoever lands a
// melee hit, and nothing to a spell.
func TestGeneratedThornsStrikesBackAtAMeleeHit(t *testing.T) {
	sim := setupFakeSimWithBuffs(&proto.RaidBuffs{Thorns: true}, &proto.PartyBuffs{}, &proto.IndividualBuffs{})
	char := sim.Raid.Parties[0].Players[0].GetCharacter()
	attacker := sim.Encounter.AllTargetUnits[0]

	if ThornsValue(0) != 22 {
		t.Errorf("Thorns is worth %v, want the client's 22", ThornsValue(0))
	}

	aura := char.GetAura("Thorns (External)")
	if aura == nil {
		t.Fatalf("no aura is labelled %q; the unit has %v", "Thorns (External)", auraLabels(char))
	}
	if !aura.IsActive() {
		t.Fatal("the raid's Thorns is not up at the start of the fight")
	}

	shield := char.GetSpell(ActionID{SpellID: 9910, Tag: 1})
	if shield == nil {
		t.Fatal("the shield registered no spell to deal its damage with")
	}

	// The shield's own swings would land on the character too, so they are
	// cancelled and every hit the test asks about is delivered by hand.
	attacker.AutoAttacks.CancelAutoSwing(sim)

	landed := &SpellResult{Target: &char.Unit, Outcome: OutcomeHit}
	aura.OnSpellHitTaken(aura, sim, attacker.AutoAttacks.MHAuto(), landed)
	sim.Step()
	if got := shield.SpellMetrics[attacker.UnitIndex].Casts; got != 1 {
		t.Errorf("a melee hit taken cast the shield %d times, want once", got)
	}
	if got := shield.SpellMetrics[attacker.UnitIndex].TotalDamage; got != 22 {
		t.Errorf("the shield dealt %v damage, want the client's 22", got)
	}

	caster := sim.Raid.Parties[0].Players[0].(*FakeAgent)
	aura.OnSpellHitTaken(aura, sim, caster.Spell, landed)
	sim.Step()
	if got := shield.SpellMetrics[attacker.UnitIndex].Casts; got != 1 {
		t.Errorf("a shadow spell taken cast the shield %d times, want the shield to ignore it", got)
	}
}

// Retribution Aura is the paladin's slot: the external copy keeps its own
// category and stays out of the shared one only the paladin's own cast joins.
func TestGeneratedRetributionAuraHoldsThePaladinSlot(t *testing.T) {
	char := newGeneratedBuffTestCharacter()

	applyBuffEffects(generatedBuffTestAgent{char},
		&proto.RaidBuffs{}, &proto.PartyBuffs{RetributionAura: true}, &proto.IndividualBuffs{})

	external := char.GetAura("Retribution Aura (External)")
	if external == nil {
		t.Fatalf("no aura is labelled %q; the unit has %v", "Retribution Aura (External)", auraLabels(char))
	}
	if RetributionAuraValue(0) != 30 {
		t.Errorf("the aura is worth %v holy damage, want the client's 30", RetributionAuraValue(0))
	}

	category := char.ExclusiveEffectManager.GetExclusiveEffectCategory(RetributionAuraCategory)
	if !category.SingleAura || len(category.effects) != 1 || category.effects[0].Priority != 30 {
		t.Errorf("the category is single-aura %v with %d effects, first bid %v; want one bidding 30",
			category.SingleAura, len(category.effects), category.effects[0].Priority)
	}
	if shared := char.ExclusiveEffectManager.GetExclusiveEffectCategory(PaladinAuraCategory); len(shared.effects) != 0 {
		t.Errorf("the external copy put %d effects in the paladin's shared category, want none",
			len(shared.effects))
	}

	player := RetributionAuraAura(&char.Unit, true, 0)
	if shared := char.ExclusiveEffectManager.GetExclusiveEffectCategory(PaladinAuraCategory); len(shared.effects) != 1 {
		t.Errorf("the paladin's own copy put %d effects in the shared category, want one", len(shared.effects))
	}
	if player == external {
		t.Error("the paladin's own copy and the external one are the same aura")
	}
}

// A whole environment, because an external cooldown registers a spell, a timer
// per source and a major cooldown, none of which a bare Character has.
func setupFakeSimWithBuffs(raid *proto.RaidBuffs, party *proto.PartyBuffs, individual *proto.IndividualBuffs) *Simulation {
	sim := NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 100},
		Raid: &proto.Raid{
			Parties: []*proto.Party{
				{
					Players: []*proto.Player{
						{
							Name:      "Caster",
							Class:     proto.Class_ClassShaman,
							Buffs:     individual,
							Spec:      &proto.Player_ElementalShaman{},
							Equipment: &proto.EquipmentSpec{},
						},
					},
					Buffs: party,
				},
			},
			Buffs: raid,
		},
		Encounter: &proto.Encounter{
			Targets: []*proto.Target{{
				Name: "target", Level: 60, MobType: proto.MobType_MobTypeDemon,
				SwingSpeed: 2, MinBaseDamage: 100,
			}},
			Duration: 180,
		},
	}, simsignals.CreateSignals())
	sim.Reset()

	return sim
}

// Two druids innervating one character take turns: the second waits for the
// first one's aura to fall off, and a third cast waits for the six minutes the
// client states.
func TestGeneratedInnervatesTakeTurnsBetweenTheirSources(t *testing.T) {
	sim := setupFakeSimWithBuffs(&proto.RaidBuffs{}, &proto.PartyBuffs{}, &proto.IndividualBuffs{Innervates: 2})
	char := sim.Raid.Parties[0].Players[0].GetCharacter()

	aura := char.GetAura("Innervates (External)")
	if aura == nil {
		t.Fatalf("no aura is labelled %q; the unit has %v", "Innervates (External)", auraLabels(char))
	}
	if aura.Duration != InnervatesDuration(0) || aura.Duration != time.Second*20 {
		t.Errorf("the innervate lasts %v, want the client's 20 seconds", aura.Duration)
	}
	if InnervatesCooldown() != time.Minute*6 {
		t.Errorf("the innervate's cooldown is %v, want the client's 6 minutes", InnervatesCooldown())
	}

	spell := char.GetSpell(ActionID{SpellID: 29166, Tag: -1})
	if spell == nil {
		t.Fatal("no spell stands for the external innervates")
	}

	// What the sim asks before it spends one of the sources: its own timer and
	// the condition the external approximation installs.
	ready := func() bool {
		return spell.CD.Timer.IsReady(sim) && spell.ExtraCastCondition(sim, &char.Unit)
	}

	if !ready() {
		t.Fatal("the first of two innervates cannot be cast at the start of the fight")
	}
	spell.SkipCastAndApplyEffects(sim, &char.Unit)
	if !aura.IsActive() {
		t.Fatal("casting the first innervate left the character without the aura")
	}
	if ready() {
		t.Error("a second innervate lands while the first one is still up")
	}

	sim.CurrentTime = aura.Duration + time.Second
	aura.Deactivate(sim)
	if !ready() {
		t.Fatal("the second druid cannot innervate once the first aura has fallen off")
	}
	spell.SkipCastAndApplyEffects(sim, &char.Unit)

	sim.CurrentTime += aura.Duration + time.Second
	aura.Deactivate(sim)
	if ready() {
		t.Error("a third innervate lands before either druid's six minutes are up")
	}

	sim.CurrentTime = InnervatesCooldown() + time.Minute
	if !ready() {
		t.Error("the first druid cannot innervate again six minutes later")
	}
}

// Power Infusion is +20% damage and healing done while it is up, and nothing
// once it has expired. The damage half is the six schools the client's mask
// names, which does not include physical.
func TestGeneratedPowerInfusionRaisesDamageAndHealingDone(t *testing.T) {
	sim := setupFakeSimWithBuffs(&proto.RaidBuffs{}, &proto.PartyBuffs{}, &proto.IndividualBuffs{PowerInfusions: 1})
	char := sim.Raid.Parties[0].Players[0].GetCharacter()

	aura := char.GetAura("Power Infusions (External)")
	if aura == nil {
		t.Fatalf("no aura is labelled %q; the unit has %v", "Power Infusions (External)", auraLabels(char))
	}
	if aura.Duration != time.Second*15 || PowerInfusionsCooldown() != time.Minute*3 {
		t.Errorf("the infusion lasts %v on a %v cooldown, want the client's 15 seconds and 3 minutes",
			aura.Duration, PowerInfusionsCooldown())
	}

	healing := char.PseudoStats.HealingDealtMultiplier
	physical := generatedSchoolIndexes(1)[0]
	before := char.PseudoStats.SchoolDamageDealtMultiplier[physical]

	aura.Activate(sim)
	// Spell 10060 states mask 126, which is every school but physical, so a
	// melee swing is not part of what the infusion raises.
	for _, school := range generatedSchoolIndexes(126) {
		if got := char.PseudoStats.SchoolDamageDealtMultiplier[school]; got != 1.2 {
			t.Errorf("school %d deals %v times the damage, want the client's 1.2", school, got)
		}
	}
	if got := char.PseudoStats.SchoolDamageDealtMultiplier[physical]; got != before {
		t.Errorf("physical damage dealt is %v, want the mask to have left it at %v", got, before)
	}
	if got := char.PseudoStats.HealingDealtMultiplier; got != healing*1.2 {
		t.Errorf("the infusion multiplies healing dealt by %v, want the client's 1.2", got/healing)
	}

	aura.Deactivate(sim)
	for _, school := range generatedSchoolIndexes(126) {
		if got := char.PseudoStats.SchoolDamageDealtMultiplier[school]; got != 1 {
			t.Errorf("school %d still deals %v times the damage once the infusion expired", school, got)
		}
	}
	if char.PseudoStats.HealingDealtMultiplier != healing {
		t.Error("the infusion left part of itself behind when it expired")
	}
}

// Mana Tide's aura states 290 mana every 3 seconds, which reaches the sim as
// the mana per 5 seconds it is worth; the totem cast states how long it stands.
func TestGeneratedManaTideTotemRestoresTheClientsMana(t *testing.T) {
	sim := setupFakeSimWithBuffs(&proto.RaidBuffs{}, &proto.PartyBuffs{ManaTideTotems: 1}, &proto.IndividualBuffs{})
	char := sim.Raid.Parties[0].Players[0].GetCharacter()

	aura := char.GetAura("Mana Tide Totem (External)")
	if aura == nil {
		t.Fatalf("no aura is labelled %q; the unit has %v", "Mana Tide Totem (External)", auraLabels(char))
	}
	if aura.Duration != time.Second*13 || ManaTideTotemsCooldown() != time.Minute*5 {
		t.Errorf("the totem stands for %v on a %v cooldown, want the client's 13 seconds and 5 minutes",
			aura.Duration, ManaTideTotemsCooldown())
	}
	if want := 290 * 5000.0 / 3000.0; ManaTideTotemsValue(0) != want {
		t.Errorf("the totem is worth %v MP5, want the 290 per 3 seconds the client states as %v",
			ManaTideTotemsValue(0), want)
	}

	before := char.stats[stats.MP5]
	aura.Activate(sim)
	if got := char.stats[stats.MP5] - before; got != ManaTideTotemsValue(0) {
		t.Errorf("the totem applied %v MP5, want %v", got, ManaTideTotemsValue(0))
	}
}

// A buff that is not simply up is not part of the stats the character sheet is
// measured with: the cooldowns, the windfury proc and the shadow priest's mana,
// the last of which has a fight to run before it can give any away.
func TestGeneratedDrivenBuffsAreNotBuildPhaseAuras(t *testing.T) {
	sim := setupFakeSimWithBuffs(&proto.RaidBuffs{},
		&proto.PartyBuffs{ManaTideTotems: 1, WindfuryTotem: true},
		&proto.IndividualBuffs{Innervates: 1, PowerInfusions: 1, ShadowPriestDps: 500})
	char := sim.Raid.Parties[0].Players[0].GetCharacter()

	for _, label := range []string{
		"Innervates (External)", "Power Infusions (External)", "Mana Tide Totem (External)",
		"Windfury Totem (External)", "Vampiric Touch (External)",
	} {
		aura := char.GetAura(label)
		if aura == nil {
			t.Fatalf("no aura is labelled %q; the unit has %v", label, auraLabels(char))
		}
		if aura.BuildPhase != CharacterBuildPhaseNone {
			t.Errorf("%s is measured in build phase %v, want none of them", label, aura.BuildPhase)
		}
	}

	// The build phase runs on a simulation that has no fight in it, which the
	// mana the shadow priest gives away would otherwise be scheduled against.
	bare := newGeneratedBuffTestCharacter()
	applyBuffEffects(generatedBuffTestAgent{bare}, &proto.RaidBuffs{}, &proto.PartyBuffs{},
		&proto.IndividualBuffs{ShadowPriestDps: 500})
	bare.applyBuildPhaseAuras(CharacterBuildPhaseBuffs)
}

func auraLabels(char *Character) []string {
	labels := make([]string, 0, len(char.auras))
	for _, aura := range char.auras {
		labels = append(labels, aura.Label)
	}
	return labels
}
