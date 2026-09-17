package sim

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// The trees decide what a player is shown; the Go decides what the sim computes. Nothing
// tied the two together and they drifted four times - Master of Defense paid 10 Rage where
// its tree read 5, Improved Rend 36% where its tree read 35, Shadow and Flame was flat
// where its tree scaled, Fire and Brimstone 24% against 25. Each looked right on the site
// while the sim did something else, which is worse than a plain wrong number because
// nothing looks wrong.
//
// This reads the two ways the Go pays a talent out and checks the series each produces
// against that talent's own rank values:
//
//	0.05 * float64(paladin.Talents.IronCreed)
//	[]float64{0, 8, 17, 25}[warlock.Talents.FireAndBrimstone]
//
// A Go value often is not a number the tooltip prints - a multiplier against a base, a
// different unit - so only a talent that agrees on some ranks and disagrees on others is
// reported. That is what drift looks like; a units mismatch disagrees everywhere.
var coefficientRegex = regexp.MustCompile(`([\d.]+)\s*\*\s*float64\(\w+\.Talents\.([A-Za-z]+)\)`)
var rankTableRegex = regexp.MustCompile(`\[\]float64\{([0-9.,\s]+)\}\[\w+\.Talents\.([A-Za-z]+)\]`)

// Talents whose Go value is knowingly not the tree's, each waiting on an observation. The
// tree figure is the community site's rounding rather than something anyone has seen, so
// matching it would trade one guess for another.
var scalingExceptions = map[string]string{
	"mage|arcaneSubtlety":      "spell penetration: tree says 15 at rank 2, the Go doubles rank 1's 8. Only rank 1 observed.",
	"priest|spiritualGuidance": "damage half: tree says 1/3/5/6/8, the Go takes 1/2/3/4/5. Only rank 1 observed.",
}

func TestGoTalentScalingMatchesTheTrees(t *testing.T) {
	trees := map[string][][]float64{}

	for _, class := range talentClasses {
		for _, tree := range loadTrees(t, class.name) {
			for i := range tree.Talents {
				talent := &tree.Talents[i]
				if talent.FieldName == "" || len(talent.Ranks) == 0 {
					continue
				}
				var ranks [][]float64
				for _, values := range talent.Ranks {
					var numbers []float64
					for _, raw := range values {
						var number float64
						if err := json.Unmarshal(raw, &number); err == nil {
							numbers = append(numbers, number)
						}
					}
					ranks = append(ranks, numbers)
				}
				trees[class.name+"|"+talent.FieldName] = ranks
			}
		}
	}

	checked := 0
	for _, class := range talentClasses {
		root := filepath.Join("..", "sim", class.name)
		err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
				return nil
			}
			source, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			for lineNumber, line := range strings.Split(string(source), "\n") {
				for _, match := range coefficientRegex.FindAllStringSubmatch(line, -1) {
					coefficient, convErr := strconv.ParseFloat(match[1], 64)
					if convErr != nil {
						continue
					}
					checked++
					compareScaling(t, trees, class.name, match[2], path, lineNumber+1,
						fmt.Sprintf("%g per point", coefficient),
						func(rank int) (float64, bool) { return coefficient * float64(rank), true })
				}

				for _, match := range rankTableRegex.FindAllStringSubmatch(line, -1) {
					var table []float64
					valid := true
					for _, field := range strings.Split(match[1], ",") {
						value, convErr := strconv.ParseFloat(strings.TrimSpace(field), 64)
						if convErr != nil {
							valid = false
							break
						}
						table = append(table, value)
					}
					if !valid {
						continue
					}
					checked++
					compareScaling(t, trees, class.name, match[2], path, lineNumber+1,
						"table "+match[1],
						func(rank int) (float64, bool) {
							if rank >= len(table) {
								return 0, false
							}
							return table[rank], true
						})
				}
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	if checked == 0 {
		t.Fatal("found no talent scaling in sim/; the shapes this looks for have probably changed")
	}
	t.Logf("checked %d talent scaling sites against the trees", checked)
}

// valueAt(rank) is what the Go pays for that many points, 1-based, and false when the Go
// has nothing for that rank.
func compareScaling(t *testing.T, trees map[string][][]float64, className, goName, path string, line int, describedAs string, valueAt func(int) (float64, bool)) {
	t.Helper()

	fieldName := strings.ToLower(goName[:1]) + goName[1:]
	ranks, ok := trees[className+"|"+fieldName]
	if !ok || len(ranks) == 0 {
		return
	}
	if reason, skip := scalingExceptions[className+"|"+fieldName]; skip {
		_ = reason
		return
	}

	var mismatches []string
	for rank, printed := range ranks {
		goValue, has := valueAt(rank + 1)
		if !has || len(printed) == 0 {
			continue
		}
		// Accept the value itself, or scaled by 100 either way: a tooltip says 2% where
		// the Go carries 0.02.
		matched := false
		for _, value := range printed {
			if withinTolerance(value, goValue) || withinTolerance(value, goValue*100) || withinTolerance(value*100, goValue) {
				matched = true
				break
			}
		}
		if !matched {
			mismatches = append(mismatches, fmt.Sprintf("rank %d: go %g, tree %v", rank+1, goValue, printed))
		}
	}

	// Disagreeing at every rank means the Go value is a different quantity from anything
	// the tooltip prints, which is common and fine.
	if len(mismatches) > 0 && len(mismatches) < len(ranks) {
		t.Errorf("%s:%d %s (%s) scales differently from its tree on some ranks but not others:\n      %s",
			filepath.ToSlash(path), line, goName, describedAs, strings.Join(mismatches, "\n      "))
	}
}

func withinTolerance(a, b float64) bool {
	diff := a - b
	return diff < 0.005 && diff > -0.005
}
