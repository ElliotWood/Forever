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
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/wowsims/forever/tools/database/buffmanifest"
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
		if row.Kind == buffmanifest.KindAbsent || row.Kind == buffmanifest.KindFlag ||
			row.Kind == buffmanifest.KindEnum {
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
		if row.Proto == buffmanifest.ProtoTristate && len(row.TalentCurve) == 0 {
			t.Errorf("%s: declared ProtoTristate but no talent in the owner's tree prices it", row.Field)
		}
		if want, pinned := pinnedTalentCurves[row.Field]; pinned && !slices.Equal(row.TalentCurve, want) {
			t.Errorf("%s: talent curve is %v, want %v", row.Field, row.TalentCurve, want)
		}
	}
}

// Mana Spring is the only buff an improving talent still prices, so it is the
// only place the curve can be checked against the client until the rest of the
// manifest stops rendering as shells. Restorative Totems modifies the aura's own
// number - 10 mana per 2 seconds - by 5% a point, and the client states that
// number as a whole one, so ranks 1 and 2 both come out at 10 per tick.
var pinnedTalentCurves = map[string][]float64{
	"mana_spring_totem": {25, 25, 27, 27, 30, 30},
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
