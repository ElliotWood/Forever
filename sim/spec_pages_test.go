package sim

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var unlaunchedSpecRegex = regexp.MustCompile(`\[Spec\.Spec(\w+)\]:\s*\{[^}]*status:\s*LaunchStatus\.Unlaunched`)
var specConfigRegex = regexp.MustCompile(`registerSpecConfig\(\s*Spec\.Spec(\w+)`)

// The site builds a page for every ui/ directory that has an entry point and styles.
// Nothing used to check that the sim could run that spec, so the four unlaunched healer
// specs still shipped pages: they loaded fine and then failed with "No agent factory" the
// moment anyone pressed Simulate. The raid sim already filtered its picker on the same
// launch statuses; only the individual pages were missing the check.
//
// discoverSpecPages now skips an unlaunched spec's page. This checks that it still does,
// because the failure is silent - a page that should not exist builds and deploys without
// complaint.
func TestPagesAreNotBuiltForUnlaunchedSpecs(t *testing.T) {
	uiRoot := filepath.Join("..", "ui")

	statuses, err := os.ReadFile(filepath.Join(uiRoot, "core", "launched_sims.ts"))
	if err != nil {
		t.Fatal(err)
	}
	unlaunched := map[string]bool{}
	for _, match := range unlaunchedSpecRegex.FindAllStringSubmatch(string(statuses), -1) {
		unlaunched[match[1]] = true
	}
	if len(unlaunched) == 0 {
		t.Fatal("no unlaunched specs found in ui/core/launched_sims.ts; has the shape changed?")
	}

	plugin, err := os.ReadFile(filepath.Join("..", "tools", "vite", "spec_pages.mts"))
	if err != nil {
		t.Fatal(err)
	}
	// The plugin names the status inside a regex, where the dot is escaped.
	if !strings.Contains(string(plugin), `LaunchStatus\.Unlaunched`) {
		t.Fatal("tools/vite/spec_pages.mts no longer filters on launch status, so a spec with no sim can ship a page again")
	}

	// The page directories the plugin is relying on that filter to skip. Naming them keeps
	// the set deliberate: launching one of these specs without removing it here, or adding
	// a new unlaunched spec, is reported rather than silently changing what deploys.
	expected := map[string]bool{
		"healing_priest":     true,
		"holy_paladin":       true,
		"restoration_druid":  true,
		"restoration_shaman": true,
	}

	entries, err := os.ReadDir(uiRoot)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if _, err := os.Stat(filepath.Join(uiRoot, entry.Name(), "index.ts")); err != nil {
			continue
		}
		if _, err := os.Stat(filepath.Join(uiRoot, "scss", "sims", entry.Name(), "index.scss")); err != nil {
			continue
		}
		sim, err := os.ReadFile(filepath.Join(uiRoot, entry.Name(), "sim.ts"))
		if err != nil {
			continue
		}
		match := specConfigRegex.FindStringSubmatch(string(sim))
		if match == nil || !unlaunched[match[1]] {
			continue
		}

		if !expected[entry.Name()] {
			t.Errorf("ui/%s registers the unlaunched spec %s, so its page is no longer built; add it to this list once that is deliberate", entry.Name(), match[1])
		}
		delete(expected, entry.Name())
	}

	for name := range expected {
		t.Errorf("ui/%s is listed here as unlaunched, but it no longer is: remove it from this list so its page ships again", name)
	}
}
