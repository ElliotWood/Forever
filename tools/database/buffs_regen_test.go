package database

// Re-renders the two generated buff files from the client database and asserts
// the committed ones agree, plus the resolver invariants that cannot be checked
// without a database.
//
// Skips when tools/database/wowsims.db is absent, which is why CI is unaffected.
//
// RenderBuffFiles reads the live sim/core tree to decide which rows sim/core
// still implements by hand, so declaring any func named <Something>Aura there,
// or reading a buff's proto field in applyBuffEffects or applyDebuffEffects,
// changes what this test expects. That is the migration switch working, not a
// broken test: regenerate.

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/wowsims/forever/tools/database/buffmanifest"
	"github.com/wowsims/forever/tools/database/dbc"
)

func openBuffTestDB(t *testing.T) *DBHelper {
	t.Helper()

	DatabasePath = "wowsims.db"
	if _, err := os.Stat(DatabasePath); err != nil {
		t.Skipf("no client database at %s - run `make db` from a local WoW install to enable this gate", DatabasePath)
	}

	helper, err := NewDBHelper()
	if err != nil {
		t.Fatalf("opening %s: %v", DatabasePath, err)
	}
	t.Cleanup(func() { helper.Close() })
	return helper
}

func TestGeneratedBuffFiles(t *testing.T) {
	helper := openBuffTestDB(t)

	files, err := RenderBuffFiles(helper)
	if err != nil {
		t.Fatalf("rendering the buff files: %v", err)
	}

	root, err := repoRoot()
	if err != nil {
		t.Fatalf("finding the repository root: %v", err)
	}

	for name, rendered := range files {
		committed, err := os.ReadFile(filepath.Join(root, name))
		if err != nil {
			t.Errorf("reading %s: %v", name, err)
			continue
		}
		if bytes.Equal(committed, rendered) {
			continue
		}
		t.Errorf("%s does not match the client database, regenerate with `make spelldata`:\n%s",
			name, unifiedBuffDiff(string(committed), string(rendered)))
	}
}

func TestResolvedBuffInvariants(t *testing.T) {
	helper := openBuffTestDB(t)

	rows, err := ResolveBuffManifest(helper)
	if err != nil {
		t.Fatalf("resolving the manifest: %v", err)
	}
	if len(rows) != len(buffmanifest.Manifest) {
		t.Fatalf("resolved %d rows for a manifest of %d", len(rows), len(buffmanifest.Manifest))
	}

	runes, err := loadRuneGrantedSpells(helper.db)
	if err != nil {
		t.Fatalf("loading the rune-granted spells: %v", err)
	}

	for _, row := range rows {
		if row.Kind == buffmanifest.KindAbsent || row.Kind == buffmanifest.KindFlag {
			if row.SpellID != 0 {
				t.Errorf("%s: %s resolved to spell %d", row.Field, row.Kind, row.SpellID)
			}
			continue
		}

		if row.SpellID == 0 {
			// A row the client has no spell for states why in its reason; that
			// is the shell the generator emits rather than an invariant break.
			if row.Supported {
				t.Errorf("%s: generated without a spell", row.Field)
			}
			continue
		}

		var name string
		err := helper.db.QueryRow(`SELECT Name_lang FROM SpellName WHERE ID = ?`, row.SpellID).Scan(&name)
		if errors.Is(err, sql.ErrNoRows) {
			t.Errorf("%s: resolved to spell %d, which has no SpellName row", row.Field, row.SpellID)
			continue
		}
		if err != nil {
			t.Fatalf("%s: reading spell %d: %v", row.Field, row.SpellID, err)
		}
		if runes[row.SpellID] {
			t.Errorf("%s: resolved to spell %d, which a rune grants", row.Field, row.SpellID)
		}
		if row.Proto == buffmanifest.ProtoTristate && len(row.TalentCurve) == 0 && row.ImpAction == nil {
			t.Errorf("%s: declared ProtoTristate but neither a talent in the owner's tree nor an ImpAction prices it", row.Field)
		}
		if want, pinned := pinnedTalentCurves[row.Field]; pinned && !slices.Equal(row.TalentCurve, want) {
			t.Errorf("%s: talent curve is %v, want %v", row.Field, row.TalentCurve, want)
		}
		if want, pinned := pinnedCategories[row.Field]; pinned && row.Category != want {
			t.Errorf("%s: the aura competes under %q, want %q", row.Field, row.Category, want)
		}
		if want, pinned := pinnedStatAmounts[row.Field]; pinned {
			got := map[string]float64{}
			for _, amount := range row.Stats {
				got[amount.Stat.StatName()] = amount.Amount
			}
			if !maps.Equal(got, want) {
				t.Errorf("%s: grants %v, want %v", row.Field, got, want)
			}
		}
	}
}

// The group the client states a buff reaches: an area aura names it in the
// effect itself, and an aura the client applies over an area or on one ally
// names it in the effect's target. Every raid and individual row names one, so
// each is checked two ways: that the client agrees with the scope, and that it
// states a group at all. A party or debuff row the client states nothing for is
// exempt, because its scope is then the sim's grouping rather than a game fact:
// the buff Windfury's proc lands on its wielder, a debuff the client reaches by
// the area around its caster.
func TestScopeMatchesTheClientTargeting(t *testing.T) {
	helper := openBuffTestDB(t)

	rows, err := ResolveBuffManifest(helper)
	if err != nil {
		t.Fatalf("resolving the manifest: %v", err)
	}

	for _, row := range rows {
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

func clientScope(row ResolvedBuff) (buffmanifest.BuffScope, bool) {
	for _, effect := range row.Effects {
		switch effect.Effect {
		case dbc.E_APPLY_AREA_AURA_RAID:
			return buffmanifest.ScopeRaid, true
		case dbc.E_APPLY_AREA_AURA_PARTY:
			return buffmanifest.ScopeParty, true
		case dbc.E_APPLY_AURA:
			switch effect.ImplicitTarget {
			case dbc.TARGET_UNIT_CASTER_AREA_RAID:
				return buffmanifest.ScopeRaid, true
			case dbc.TARGET_UNIT_CASTER_AREA_PARTY:
				return buffmanifest.ScopeParty, true
			case dbc.TARGET_UNIT_TARGET_ALLY, dbc.TARGET_UNIT_TARGET_ALLY_OR_RAID:
				return buffmanifest.ScopeIndividual, true
			}
		}
	}
	return buffmanifest.ScopeIndividual, false
}

// The two auras the client states as a bare A_MOD_CRIT_PCT, which carries no
// school: what they are worth is the client's 3, and it lands on every kind of
// crit, which is what the manifest's StatOverride says. Nothing else can see
// that mapping while both rows render as shells.
//
// Expose Armor states nothing on the effect itself: its -450 armor is per combo
// point, and the raid config's debuff is the five-point finisher.
var pinnedStatAmounts = map[string]map[string]float64{
	"leader_of_the_pack": {"PhysicalCritPercent": 3, "SpellCritPercent": 3},
	"moonkin_aura":       {"PhysicalCritPercent": 3, "SpellCritPercent": 3},
	"expose_armor":       {"Armor": -2250},
}

// Mana Spring is the only buff an improving talent still prices, so it is the
// only place the curve can be checked against the client. Restorative Totems
// modifies the aura's own number - 10 mana per 2 seconds - by 5% a point, and
// the client states that number as a whole one, so ranks 1 and 2 both come out
// at 10 per tick; the conversion to mana per 5 seconds keeps the half the 11 of
// ranks 3 and 4 is worth.
var pinnedTalentCurves = map[string][]float64{
	"mana_spring_totem": {25, 25, 27.5, 27.5, 30, 30},
}

// What a resistance row competes under, as (stats, own aura). A source that has
// no exclusivity beyond the school itself keeps no category of its own, which is
// how every totem, Aspect of the Wild and Shadow Protection read; a paladin aura
// also holds its own slot. Armor is not a school, so Devotion Aura has only the
// slot. Nothing else can see this while every row renders as a shell.
// A row whose manifest category is the resistance school itself keeps none of its
// own: the school category the sim puts the stat into is the whole competition.
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
