package database

// Re-renders the buff files from the client rows committed in assets/db_inputs/spell_store_inputs.json
// and asserts the committed ones agree, plus the resolver invariants the rendered code cannot show.
// No client database: the buffs render from the same capture the store does, so this runs in CI.

import (
	"bytes"
	"fmt"
	"maps"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/wowsims/forever/tools/database/buffmanifest"
)

func resolveCommittedBuffs(t *testing.T) []ResolvedBuff {
	t.Helper()
	inRepositoryRoot(t)

	inputs, err := readStoreInputs(spellStoreInputsPath)
	if err != nil {
		t.Fatalf("%v", err)
	}
	rows, err := resolveBuffManifest(inputs)
	if err != nil {
		t.Fatalf("resolving the manifest: %v", err)
	}
	return rows
}

func TestBuffFilesRegenerateFromTheCommittedInputs(t *testing.T) {
	inRepositoryRoot(t)

	inputs, err := readStoreInputs(spellStoreInputsPath)
	if err != nil {
		t.Fatalf("%v", err)
	}
	files, err := renderBuffOutputs(inputs)
	if err != nil {
		t.Fatalf("rendering the buff files: %v", err)
	}

	for name, rendered := range files {
		committed, err := os.ReadFile(name)
		if err != nil {
			t.Errorf("reading %s: %v", name, err)
			continue
		}
		if !bytes.Equal(committed, rendered) {
			t.Errorf("%s is not what the committed inputs render, regenerate with `go run ./tools/database/gen_spelldata`:\n%s",
				name, unifiedBuffDiff(string(committed), string(rendered)))
		}
	}
}

func TestResolvedBuffInvariants(t *testing.T) {
	rows := resolveCommittedBuffs(t)
	if len(rows) != len(buffmanifest.Manifest) {
		t.Fatalf("resolved %d rows for a manifest of %d", len(rows), len(buffmanifest.Manifest))
	}

	for _, row := range rows {
		if row.Kind == buffmanifest.KindAbsent || row.Kind == buffmanifest.KindFlag {
			if row.SpellID != 0 {
				t.Errorf("%s: %s resolved to spell %d", row.Field, row.Kind, row.SpellID)
			}
			continue
		}
		if row.Proto == buffmanifest.ProtoTristate && row.TalentRanks == 0 && row.ImpAction == nil {
			t.Errorf("%s: declared ProtoTristate but neither a talent nor an ImpAction prices it", row.Field)
		}
		if want, pinned := pinnedTalentRanks[row.Field]; pinned && row.TalentRanks != want {
			t.Errorf("%s: the talent takes %d points, want %d", row.Field, row.TalentRanks, want)
		}
		if want, pinned := pinnedCategories[row.Field]; pinned && row.Category != want {
			t.Errorf("%s: the aura competes under %q, want %q", row.Field, row.Category, want)
		}
		if want, pinned := pinnedAmounts[row.Field]; pinned {
			got := map[string]float64{}
			for _, applied := range row.Applied {
				got[applied.Kind] = applied.Value
			}
			if !maps.Equal(got, want) {
				t.Errorf("%s: grants %v, want %v", row.Field, got, want)
			}
		}
	}
}

// Every raid and individual row names the group it reaches, so each is checked two ways: that the
// client agrees with the scope, and that it states a group at all. A party or debuff row the client
// states nothing for is exempt, because its scope is then the sim's grouping rather than a game
// fact: the buff Windfury's proc lands on its wielder, a debuff the client reaches by the area around
// its caster.
func TestScopeMatchesTheClientTargeting(t *testing.T) {
	for _, row := range resolveCommittedBuffs(t) {
		if row.SpellID == 0 {
			continue
		}
		complaint := scopeComplaint(row)
		if reason, exempt := scopeExemptions[row.Field]; exempt {
			if complaint == "" {
				t.Errorf("%s: spell %d now states %s, so the exemption %q is stale",
					row.Field, row.SpellID, row.Scope, reason)
			}
			continue
		}
		if complaint != "" {
			t.Errorf("%s: %s", row.Field, complaint)
		}
	}
}

// The rows whose spell the client targets somewhere else than the scope reaches,
// each with the reason the scope is still what the sim wants.
var scopeExemptions = map[string]string{
	"thorns": "the druid trains no raid-wide spell of the family, so the raid scope is the sim spreading the one cast 9910 states over every player",
}

// What the client says against the scope the manifest files the row under, and
// "" when the two agree or the row is one the client may state nothing for.
func scopeComplaint(row ResolvedBuff) string {
	scope, stated := clientScope(row)
	switch {
	case stated && scope == row.Scope:
		return ""
	case stated:
		return fmt.Sprintf("the manifest says %s, spell %d states %s", row.Scope, row.SpellID, scope)
	case row.Scope == buffmanifest.ScopeRaid || row.Scope == buffmanifest.ScopeIndividual:
		return fmt.Sprintf("the manifest says %s, spell %d states no group it reaches", row.Scope, row.SpellID)
	}
	return ""
}

// What the parse attaches for a row, by the kind it attaches as. The two auras
// the client states as a bare A_MOD_CRIT_PCT carry no school: what they are
// worth is the client's 3, and it lands on every kind of crit.
//
// Expose Armor states nothing on the effect itself: its -450 armor is per combo
// point, and the raid config's debuff is the five-point finisher.
//
// The paladin auras state a healing-taken row of 0 beside the aura, which the
// manifest's SkipAuras leaves out.
var pinnedAmounts = map[string]map[string]float64{
	"leader_of_the_pack": {"stat PhysicalCritPercent+SpellCritPercent": 3},
	"moonkin_aura":       {"stat PhysicalCritPercent+SpellCritPercent": 3},
	"expose_armor":       {"stat Armor": -2250},
	"devotion_aura":      {"stat Armor": 735},
	"concentration_aura": {"pushback": -0.35},
	"retribution_aura":   {},
}

// Mana Spring is the only buff an improving talent prices: Restorative Totems, five points.
var pinnedTalentRanks = map[string]int32{
	"mana_spring_totem": 5,
}

// What a resistance row competes under, as (stats, own aura). A source that has
// no exclusivity beyond the school itself keeps no category of its own, which is
// how every totem, Aspect of the Wild and Shadow Protection read; a paladin aura
// also holds its own slot. Armor is not a school, so Devotion Aura has only the
// slot.
var pinnedCategories = map[string]string{
	"frost_resistance_totem":      "",
	"aspect_of_the_wild":          "",
	"prayer_of_shadow_protection": "",
	"frost_resistance_aura":       "FrostResistanceAura",
	"shadow_resistance_aura":      "ShadowResistanceAura",
	"devotion_aura":               "DevotionAura",
}

// The first differing line of each file with a little context, which is all a
// reader needs to see whether the generated file or the manifest moved.
func unifiedBuffDiff(committed string, rendered string) string {
	committedLines := strings.Split(committed, "\n")
	renderedLines := strings.Split(rendered, "\n")

	var b strings.Builder
	for i := 0; i < max(len(committedLines), len(renderedLines)); i++ {
		var left, right string
		if i < len(committedLines) {
			left = committedLines[i]
		}
		if i < len(renderedLines) {
			right = renderedLines[i]
		}
		if left == right {
			continue
		}
		b.WriteString("@@ line " + strconv.Itoa(i+1) + " @@\n")
		b.WriteString("-" + left + "\n")
		b.WriteString("+" + right + "\n")
		if b.Len() > 4000 {
			b.WriteString("... truncated\n")
			break
		}
	}
	return b.String()
}
