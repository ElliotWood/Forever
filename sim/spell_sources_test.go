package sim

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// The Go sim carries an integer for every ability and nothing else. Every name, icon and
// hover tooltip the site shows is resolved in the browser from Wowhead's Classic data for
// that integer, so an ability Forever changed is described by its Classic ancestor unless
// something says otherwise. That is the same trap the talent trees fell into, where 42
// talents wore another talent's tooltip and Improved Revenge kept Classic's stun long
// after Forever turned it into damage.
//
// ui/core/spells/<class>.json is that something: one entry per spell id the sim registers,
// saying where the numbers came from and, when Forever changed them, what the ability
// actually does. The UI renders those entries instead of the Wowhead tooltip. This test
// keeps the manifest and the sim in step - an ability cannot be registered without saying
// where its numbers came from, and the number still waiting on an answer can only fall.

// Raised by hand when an id is deliberately left undeclared, never by a tool. Every entry
// needs a reason, because an undeclared id is an ability the site describes wrongly.
const unreviewedSpellBudget = 313

type spellSource struct {
	Ability     string   `json:"ability"`
	File        string   `json:"file"`
	Source      string   `json:"source"`
	ForeverID   int      `json:"foreverId,omitempty"`
	Tooltip     string   `json:"tooltip,omitempty"`
	Note        string   `json:"note,omitempty"`
	Assumptions []string `json:"assumptions,omitempty"`
}

var validSpellSources = map[string]bool{
	// Forever did not change the ability, so Wowhead's Classic tooltip describes it.
	"classic": true,
	// Forever changed it and a published Forever source gave the numbers.
	"forever": true,
	// Forever changed it or it is new, and at least one number is a guess.
	"assumed": true,
	// Not classified yet. Held to unreviewedSpellBudget.
	"unreviewed": true,
}

// Every SpellID literal the sim registers, by id, with the files that use it. Reading the
// syntax tree rather than grepping keeps ids in comments and in strings out of the count.
func registeredSpellIDs(t *testing.T) map[int][]string {
	t.Helper()

	ids := map[int]map[string]bool{}
	fset := token.NewFileSet()

	err := filepath.WalkDir("..", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			// sim/ is the whole of the simulator; nothing outside it registers a spell.
			if !strings.HasPrefix(filepath.ToSlash(path), "../sim") && path != ".." {
				return fs.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		file, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		rel := strings.TrimPrefix(filepath.ToSlash(path), "../")

		ast.Inspect(file, func(node ast.Node) bool {
			lit, ok := node.(*ast.CompositeLit)
			if !ok || !isActionID(lit.Type) {
				return true
			}
			for _, elt := range lit.Elts {
				kv, ok := elt.(*ast.KeyValueExpr)
				if !ok {
					continue
				}
				key, ok := kv.Key.(*ast.Ident)
				if !ok || key.Name != "SpellID" {
					continue
				}
				value, ok := kv.Value.(*ast.BasicLit)
				if !ok || value.Kind != token.INT {
					continue
				}
				id, err := strconv.Atoi(value.Value)
				// A zero id is the empty action, which names nothing and needs no entry.
				if err != nil || id == 0 {
					continue
				}
				if ids[id] == nil {
					ids[id] = map[string]bool{}
				}
				ids[id][rel] = true
			}
			return true
		})
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	out := map[int][]string{}
	for id, files := range ids {
		for file := range files {
			out[id] = append(out[id], file)
		}
		sort.Strings(out[id])
	}
	return out
}

func isActionID(expr ast.Expr) bool {
	switch typ := expr.(type) {
	case *ast.Ident:
		return typ.Name == "ActionID"
	case *ast.SelectorExpr:
		return typ.Sel.Name == "ActionID"
	}
	return false
}

func loadSpellSources(t *testing.T) map[int]spellSource {
	t.Helper()

	dir := filepath.Join("..", "ui", "core", "spells")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}

	sources := map[int]spellSource{}
	owner := map[int]string{}
	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			t.Fatal(err)
		}
		var file map[string]spellSource
		if err := json.Unmarshal(data, &file); err != nil {
			t.Fatalf("%s: %v", entry.Name(), err)
		}
		for key, source := range file {
			id, err := strconv.Atoi(key)
			if err != nil {
				t.Errorf("%s: %q is not a spell id", entry.Name(), key)
				continue
			}
			if previous, ok := owner[id]; ok {
				t.Errorf("spell %d is declared in both %s and %s", id, previous, entry.Name())
				continue
			}
			owner[id] = entry.Name()
			sources[id] = source
		}
	}
	return sources
}

func TestEveryRegisteredSpellSaysWhereItsNumbersCameFrom(t *testing.T) {
	registered := registeredSpellIDs(t)
	sources := loadSpellSources(t)

	if len(registered) == 0 {
		t.Fatal("found no spell ids in sim/, the walk is broken rather than the manifest")
	}

	var undeclared []int
	for id := range registered {
		if _, ok := sources[id]; !ok {
			undeclared = append(undeclared, id)
		}
	}
	sort.Ints(undeclared)
	for _, id := range undeclared {
		t.Errorf("spell %d (%s) is not in ui/core/spells, so the site describes it with Wowhead's Classic entry",
			id, strings.Join(registered[id], ", "))
	}

	var orphaned []int
	for id := range sources {
		if _, ok := registered[id]; !ok {
			orphaned = append(orphaned, id)
		}
	}
	sort.Ints(orphaned)
	for _, id := range orphaned {
		t.Errorf("spell %d is declared in ui/core/spells but nothing in sim/ registers it", id)
	}
}

func TestSpellSourcesAreWellFormed(t *testing.T) {
	for id, source := range loadSpellSources(t) {
		if source.Ability == "" {
			t.Errorf("spell %d: no ability name", id)
		}
		if source.File == "" {
			t.Errorf("spell %d (%s): no file", id, source.Ability)
		}
		if !validSpellSources[source.Source] {
			t.Errorf("spell %d (%s): %q is not a source", id, source.Ability, source.Source)
		}

		// An ability Forever changed has to carry its own words, or the site falls back to
		// the Classic tooltip and the entry has bought nothing.
		if source.Source == "forever" || source.Source == "assumed" {
			if source.Tooltip == "" {
				t.Errorf("spell %d (%s): %s with no tooltip of its own", id, source.Ability, source.Source)
			}
		} else if source.Tooltip != "" {
			t.Errorf("spell %d (%s): %s does not need a tooltip, Wowhead's is right", id, source.Ability, source.Source)
		}

		// A guess that does not say what was guessed cannot be checked against the beta.
		if source.Source == "assumed" && len(source.Assumptions) == 0 {
			t.Errorf("spell %d (%s): assumed but lists nothing to confirm", id, source.Ability)
		}
		if source.Source != "assumed" && len(source.Assumptions) > 0 {
			t.Errorf("spell %d (%s): %s does not assume anything", id, source.Ability, source.Source)
		}
	}
}

func TestUnreviewedSpellsOnlyShrink(t *testing.T) {
	var unreviewed []int
	for id, source := range loadSpellSources(t) {
		if source.Source == "unreviewed" {
			unreviewed = append(unreviewed, id)
		}
	}
	sort.Ints(unreviewed)

	if len(unreviewed) > unreviewedSpellBudget {
		t.Errorf("%d spells are unreviewed and the budget is %d: %v", len(unreviewed), unreviewedSpellBudget, unreviewed)
	}
	if len(unreviewed) < unreviewedSpellBudget {
		t.Errorf("only %d spells are unreviewed, lower unreviewedSpellBudget to %d", len(unreviewed), len(unreviewed))
	}
}
