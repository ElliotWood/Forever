// Runs every build a spec can be given from what is already in the repository, so the
// leaderboard is a result rather than a guess about which build to show.
//
// The rankings page picks one preset per spec and runs the 26 of them together in the
// browser, which takes about forty seconds. That is fine for one build each and impossible
// for all of them: the gear sets and rotations already sitting in ui/<spec>/ multiply out to
// a few hundred runs, and nobody holds a tab open for that. So this runs headless, off the
// same files, and the site reads what it produced.
//
// It lives as a test rather than a command for one reason: the talents, spec options and
// consumes that make a build are defined in each spec's own _test.go, and a normal package
// cannot see them. Duplicating twenty struct literals into a cmd/ tool would mean the arena
// could silently drift from what the regression tests actually run, which is the one thing
// that would make its numbers worthless. Each spec contributes ten lines next to the
// definitions themselves instead.
//
// Every build meets the same fixed environment - one buff set, one consumable set, one
// encounter - because that, not a shared raid, is what makes two numbers comparable. The
// rankings page achieves the same thing by putting everyone in one raid; at this count that
// is not an option, and a fixed environment is if anything fairer, since no build gets a
// better draw of party members than another.
package arenalib

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
)

// Where each spec drops its results. The merge step in tools/arena reads the directory;
// separate files per spec because `go test` runs packages concurrently and a shared file
// would be a race.
const outDirEnv = "ARENA_OUT"

// Searching the talent trees costs roughly a sim run per talent per step, so it is off by
// default and belongs to the scheduled job rather than the one that runs on every push.
const optimiseEnv = "ARENA_OPTIMISE"

// Sim runs the talent search may spend per spec. The wall clock, not the answer, is what
// this protects: the climb stops on its own when no single point move helps.
const optimiseBudget = 900

// Enough that two builds a few DPS apart are distinguishable, and the whole arena still
// finishes inside a CI job. A single-player run is far cheaper than the rankings page's
// 26-player one, so this buys more precision than that page does at a fraction of the cost.
const iterations = int32(5000)

// One row of the leaderboard: a build, what it did, and what its damage was made of.
type Result struct {
	Spec     string  `json:"spec"`
	Talents  string  `json:"talents"`
	Build    string  `json:"build"`
	Gear     string  `json:"gear"`
	Rotation string  `json:"rotation"`
	Dps      float64 `json:"dps"`
	// Damage by spell id, for the confidence column. Weighting happens in the merge step,
	// where the manifest is read once rather than once per spec.
	Damage map[string]float64 `json:"damage"`
	// Damage from white swings, which have no spell id and no manifest entry to have.
	WeaponDamage float64 `json:"weaponDamage"`
	// Found by searching the talent trees rather than written down by a person.
	Optimised bool `json:"optimised,omitempty"`
}

// What a spec needs to contribute. Everything here already exists in the spec's test file.
type Spec struct {
	// The ui/ directory, which is also the key the site looks results up by.
	Dir   string
	Class proto.Class
	Race  proto.Race

	SpecOptions interface{}
	Consumes    core.ConsumesCombo
	Cooldowns   *proto.Cooldowns
	Buffs       core.BuffsCombo

	IsTank             bool
	IsHealer           bool
	InFrontOfTarget    bool
	DistanceFromTarget float64

	// Extra talent builds beyond the community ones read from the UI presets. Rarely needed:
	// the presets are what the site shows, so they are what the arena should rank.
	ExtraTalents []TalentBuild
}

type TalentBuild struct {
	Name    string
	Talents string
}

// Run enumerates the spec's builds and writes them out. A no-op unless ARENA_OUT is set, so
// an ordinary `go test ./...` never pays for it.
func Run(t *testing.T, spec Spec) {
	outDir := os.Getenv(outDirEnv)
	if outDir == "" {
		t.Skipf("%s is not set, so the arena is not being generated", outDirEnv)
	}

	uiDir := filepath.Join("..", "..", "..", "ui", spec.Dir)
	if _, err := os.Stat(uiDir); err != nil {
		// Spec directories are two or three levels deep depending on the class, and getting it
		// wrong should say so rather than produce an empty spec.
		uiDir = filepath.Join("..", "..", "ui", spec.Dir)
	}

	talents := append(communityTalents(t, uiDir), spec.ExtraTalents...)
	// A build the search found on some previous run is carried forward and re-simulated
	// rather than inherited. Without this a push-triggered rebuild - which does not search -
	// would quietly drop every optimised build the weekly one found; with it, the number is
	// always produced by the sim as it stands today even when the search has not run since.
	carried := previouslyOptimised(uiDir, spec.Dir)
	talents = append(talents, carried...)
	wasOptimised := map[string]bool{}
	for _, build := range carried {
		wasOptimised[build.Talents] = true
	}
	gearSets := namesIn(t, filepath.Join(uiDir, "gear_sets"), ".gear.json")
	rotations := namesIn(t, filepath.Join(uiDir, "apls"), ".apl.json")
	if len(talents) == 0 || len(gearSets) == 0 {
		t.Fatalf("%s: %d talent builds and %d gear sets, so there is nothing to rank", spec.Dir, len(talents), len(gearSets))
	}
	// A spec with no APL on disk still runs; the sim falls back to its own default rotation.
	if len(rotations) == 0 {
		rotations = []string{""}
	}

	results := []Result{}
	for _, talent := range talents {
		for _, gear := range gearSets {
			for _, rotation := range rotations {
				results = append(results, run(spec, uiDir, talent, gear, rotation))
			}
		}
	}

	for i := range results {
		results[i].Optimised = wasOptimised[results[i].Talents]
	}
	sort.Slice(results, func(i, j int) bool { return results[i].Dps > results[j].Dps })

	// The search starts from the best build found above rather than from an arbitrary one:
	// a climb goes to the nearest peak, so where it starts is most of what it finds.
	//
	// On launch gear specifically, even when the spec has a better set on file. The
	// leaderboard compares specs on launch gear because that is the one tier all of them
	// have, and a build optimised against phase 2 gear would never appear in it - the one
	// view most people read would be the one view missing the searched builds.
	if os.Getenv(optimiseEnv) != "" {
		top := results[0]
		for _, result := range results {
			if strings.Contains(result.Gear, "launch") {
				top = result
				break
			}
		}
		found, runs, err := optimise(spec, uiDir, TalentBuild{Name: top.Build, Talents: top.Talents}, top.Gear, top.Rotation, optimiseBudget)
		if err != nil {
			t.Fatalf("%s: %s", spec.Dir, err)
		}
		if found.Talents != top.Talents && !hasTalents(results, found.Talents) {
			row := run(spec, uiDir, found, top.Gear, top.Rotation)
			row.Optimised = true
			results = append(results, row)
			sort.Slice(results, func(i, j int) bool { return results[i].Dps > results[j].Dps })
			t.Logf("%s: %d runs of search moved %s from %.1f to %.1f dps (%s)",
				spec.Dir, runs, top.Build, top.Dps, row.Dps, found.Talents)
		} else {
			t.Logf("%s: %d runs of search found nothing better than %s", spec.Dir, runs, top.Build)
		}
	}

	write(t, filepath.Join(outDir, spec.Dir+".json"), results)
	t.Logf("%s: %d builds, best %.1f dps", spec.Dir, len(results), results[0].Dps)
}

func run(spec Spec, uiDir string, talent TalentBuild, gear string, rotation string) Result {
	return runAt(spec, uiDir, talent, gear, rotation, iterations)
}

func runAt(spec Spec, uiDir string, talent TalentBuild, gear string, rotation string, iterations int32) Result {
	gearCombo := core.GetGearSet(filepath.Join(uiDir, "gear_sets"), gear)
	rotationProto := &proto.APLRotation{}
	if rotation != "" {
		rotationProto = core.GetAplRotation(filepath.Join(uiDir, "apls"), rotation).Rotation
	}

	distance := spec.DistanceFromTarget
	if distance == 0 {
		distance = 5
	}

	player := core.WithSpec(&proto.Player{
		Class:              spec.Class,
		Race:               spec.Race,
		Equipment:          gearCombo.GearSet,
		Consumes:           spec.Consumes.Consumes,
		Buffs:              spec.Buffs.Player,
		TalentsString:      talent.Talents,
		Profession1:        proto.Profession_Engineering,
		Rotation:           rotationProto,
		Cooldowns:          spec.Cooldowns,
		InFrontOfTarget:    spec.InFrontOfTarget,
		DistanceFromTarget: distance,
		ReactionTimeMs:     150,
		ChannelClipDelayMs: 50,
	}, spec.SpecOptions)

	raid := core.SinglePlayerRaidProto(player, spec.Buffs.Party, spec.Buffs.Raid, spec.Buffs.Debuffs)
	if spec.IsTank {
		raid.Tanks = append(raid.Tanks, &proto.UnitReference{Type: proto.UnitReference_Player, Index: 0})
	}
	if spec.IsHealer {
		raid.TargetDummies = 1
	}

	result := core.RunRaidSim(&proto.RaidSimRequest{
		Raid:      raid,
		Encounter: core.MakeSingleTargetEncounter(0),
		SimOptions: &proto.SimOptions{
			Iterations: iterations,
			IsTest:     true,
			Debug:      false,
			RandomSeed: 101,
			Ruleset:    proto.Ruleset_RulesetForever,
		},
	})

	row := Result{
		Spec:     spec.Dir,
		Talents:  talent.Talents,
		Build:    talent.Name,
		Gear:     gear,
		Rotation: rotation,
		Damage:   map[string]float64{},
	}
	if result.Error != nil || result.RaidMetrics == nil || len(result.RaidMetrics.Parties) == 0 {
		return row
	}

	metrics := result.RaidMetrics.Parties[0].Players[0]
	row.Dps = metrics.Dps.Avg
	collect(&row, metrics)
	for _, pet := range metrics.Pets {
		collect(&row, pet)
	}
	return row
}

// Damage per spell, summed over targets. A pet's damage is the build's damage, and the
// manifest classifies pet abilities the same as any other.
func collect(row *Result, unit *proto.UnitMetrics) {
	for _, action := range unit.Actions {
		damage := 0.0
		for _, target := range action.Targets {
			damage += target.Damage
		}
		if damage <= 0 {
			continue
		}
		if spellId := action.Id.GetSpellId(); spellId != 0 {
			row.Damage[fmt.Sprint(spellId)] += damage
		} else if action.Id.GetOtherId() != proto.OtherAction_OtherActionNone {
			row.WeaponDamage += damage
		} else {
			// An item with no spell behind it. Counted as unclassified rather than as weapon
			// damage, because it is an effect somebody has to have got right.
			row.Damage["0"] += damage
		}
	}
}

// The community builds the site already shows, read straight out of the UI preset that
// defines them, so the arena and the talent picker can never disagree about what a build is.
//
// Read with a regex rather than by running the TypeScript: these are single-line literals of
// a fixed shape, and the alternative is a Node step in the middle of a Go test. A spec that
// yields nothing fails the test rather than quietly ranking one build.
var talentPreset = regexp.MustCompile(`makePresetTalents\(\s*'([^']+)'\s*,\s*SavedTalents\.create\(\{\s*talentsString:\s*'([^']+)'`)

func communityTalents(t *testing.T, uiDir string) []TalentBuild {
	source, err := os.ReadFile(filepath.Join(uiDir, "presets.ts"))
	if err != nil {
		t.Fatalf("no presets to read builds from: %s", err)
	}

	// The same rule the raid page uses to tell a community build from a leftover preset.
	named := regexp.MustCompile(`\d+/\d+/\d+$`)
	builds := []TalentBuild{}
	for _, match := range talentPreset.FindAllStringSubmatch(string(source), -1) {
		if named.MatchString(match[1]) {
			builds = append(builds, TalentBuild{Name: match[1], Talents: match[2]})
		}
	}
	if len(builds) == 0 {
		// Same fallback as the raid page: a spec whose presets are not named in the x/y/z
		// style still gets its first build ranked rather than disappearing from the table.
		for _, match := range talentPreset.FindAllStringSubmatch(string(source), -1) {
			builds = append(builds, TalentBuild{Name: match[1], Talents: match[2]})
			break
		}
	}
	return builds
}

func namesIn(t *testing.T, dir string, suffix string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	names := []string{}
	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), suffix) {
			continue
		}
		name := strings.TrimSuffix(entry.Name(), suffix)
		// Several specs ship an empty gear set for the picker to start from. Simulating a
		// naked character produces a number, which is the problem: it would sit at the bottom
		// of a leaderboard looking like a finding about the build.
		if name == "blank" {
			continue
		}
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func write(t *testing.T, path string, results []Result) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	encoded, err := json.MarshalIndent(results, "", "\t")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, encoded, 0o644); err != nil {
		t.Fatal(err)
	}
}

// Marks the carried-forward builds so the page can still label them, and keeps the search
// from re-adding a build that is already in the table.
func hasTalents(results []Result, talents string) bool {
	for _, result := range results {
		if result.Talents == talents {
			return true
		}
	}
	return false
}

// The optimised builds a previous run committed, read back out of the published file.
//
// Deliberately talents only. Carrying the DPS forward would leave a number describing a sim
// that no longer exists sitting in a table of numbers that do, and there is no way to tell
// them apart by looking.
func previouslyOptimised(uiDir string, spec string) []TalentBuild {
	path := filepath.Join(filepath.Dir(uiDir), "arena", "results.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var published struct {
		Builds []struct {
			Spec      string `json:"spec"`
			Build     string `json:"build"`
			Talents   string `json:"talents"`
			Optimised bool   `json:"optimised"`
		} `json:"builds"`
	}
	if json.Unmarshal(data, &published) != nil {
		return nil
	}

	seen := map[string]bool{}
	builds := []TalentBuild{}
	for _, build := range published.Builds {
		if build.Spec == spec && build.Optimised && !seen[build.Talents] {
			seen[build.Talents] = true
			builds = append(builds, TalentBuild{Name: build.Build, Talents: build.Talents})
		}
	}
	return builds
}
