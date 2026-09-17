package sim

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"math"
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
// different unit - so a talent that agrees on some ranks and disagrees on others is
// reported. That is what drift looks like; a units mismatch disagrees everywhere.
//
// Disagreeing everywhere is also reported when the Go pays the tree's values divided by
// the SAME factor at every rank, which is one quantity scaled wrong rather than two
// quantities. Moonfury and Lethality both shipped that way. Two shapes are excluded there
// because they are ordinary rather than drift: a talent with another Go value already
// matching its tree exactly, and a Go value that is one of the fixed numbers the tooltip
// states outright (warrior Enrage pays 0.3 for the "30% chance" its own text names).
//
// Known gap: a Go rate against a tree whose values are NOT proportional to it slips
// through both rules, since it neither agrees anywhere nor holds a constant ratio.
// Druid Moonglow did exactly that - a 3-per-point Go against a tree reading 8/17/25.
var coefficientRegex = regexp.MustCompile(`([\d.]+)\s*\*\s*float64\(\w+\.Talents\.([A-Za-z]+)\)`)
var rankTableRegex = regexp.MustCompile(`\[\]float64\{([0-9.,\s]+)\}\[\w+\.Talents\.([A-Za-z]+)\]`)

// Numbers written into a tooltip's text rather than left as a {0} placeholder.
var descriptionNumberRegex = regexp.MustCompile(`\d+(?:\.\d+)?`)

// Talents whose Go value is knowingly not the tree's, each waiting on an observation. The
// tree figure is the community site's rounding rather than something anyone has seen, so
// matching it would trade one guess for another.
var scalingExceptions = map[string]string{
	"mage|arcaneSubtlety":      "spell penetration: tree says 15 at rank 2, the Go doubles rank 1's 8. Only rank 1 observed.",
	"priest|spiritualGuidance": "damage half: tree says 1/3/5/6/8, the Go takes 1/2/3/4/5. Only rank 1 observed.",
}

// A talent whose Go pays its tree's values divided by the same factor at every rank.
// Held until the whole sweep is done: a talent with several Go values (Maelstrom Weapon
// has a proc rate, a cast time and a cost) can have one of them land at a tidy ratio to
// the tooltip by coincidence, and the giveaway is that another site for the same talent
// matches the tree exactly.
type ratioFinding struct {
	key         string
	path        string
	line        int
	goName      string
	describedAs string
	factor      float64
	goRank1     float64
	mismatches  []string
}

func TestGoTalentScalingMatchesTheTrees(t *testing.T) {
	trees := map[string][][]float64{}
	descriptionLiterals := map[string][]float64{}
	var ratioFindings []ratioFinding
	agrees := map[string]bool{}

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

				// Numbers written into the tooltip text rather than left as a {0}
				// placeholder are fixed parts of the talent, not per-rank values. A Go
				// constant matching one of them is that number, not a drifted rank value:
				// warrior Enrage pays 0.3 per point for the "30% chance" its text states.
				for _, literal := range descriptionNumberRegex.FindAllString(talent.Description, -1) {
					if number, err := strconv.ParseFloat(literal, 64); err == nil {
						descriptionLiterals[class.name+"|"+talent.FieldName] =
							append(descriptionLiterals[class.name+"|"+talent.FieldName], number)
					}
				}
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
						func(rank int) (float64, bool) { return coefficient * float64(rank), true }, &ratioFindings, agrees)
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
						}, &ratioFindings, agrees)
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

	// Report the constant-ratio findings now that every site has been seen, skipping any
	// talent that also has a Go value agreeing with its tree exactly — that one is the
	// tooltip's quantity, and the odd ratio elsewhere is a coincidence between two
	// different numbers.
	for _, finding := range ratioFindings {
		if agrees[finding.key] {
			continue
		}
		if matchesDescriptionLiteral(finding.goRank1, descriptionLiterals[finding.key]) {
			continue
		}
		t.Errorf("%s:%d %s (%s) pays its tree's values divided by %g at every rank, so the two disagree about one quantity rather than describing different ones:\n      %s",
			filepath.ToSlash(finding.path), finding.line, finding.goName, finding.describedAs,
			finding.factor, strings.Join(finding.mismatches, "\n      "))
	}

	t.Logf("checked %d talent scaling sites against the trees", checked)
}

// valueAt(rank) is what the Go pays for that many points, 1-based, and false when the Go
// has nothing for that rank.
func compareScaling(t *testing.T, trees map[string][][]float64, className, goName, path string, line int, describedAs string, valueAt func(int) (float64, bool), findings *[]ratioFinding, agrees map[string]bool) {
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

	// Disagreeing at every rank usually means the Go value is a different quantity from
	// anything the tooltip prints — seconds against a percentage, say — which is common
	// and fine. The exception is a disagreement that is the SAME factor at every rank:
	// that is the one quantity scaled wrong, not two quantities. Moonglow, Moonfury and
	// Lethality all shipped that way when the trees took new per-rank values and the Go
	// kept its old per-point rate.
	if len(mismatches) == 0 {
		agrees[className+"|"+fieldName] = true
		return
	}

	if len(mismatches) == len(ranks) && len(ranks) > 1 {
		rank1, _ := valueAt(1)
		if factor, constant := constantRatio(ranks, valueAt); constant {
			*findings = append(*findings, ratioFinding{
				key: className + "|" + fieldName, path: path, line: line,
				goName: goName, describedAs: describedAs, factor: factor, mismatches: mismatches,
				goRank1: rank1,
			})
		}
		return
	}

	if len(mismatches) > 0 && len(mismatches) < len(ranks) {
		t.Errorf("%s:%d %s (%s) scales differently from its tree on some ranks but not others:\n      %s",
			filepath.ToSlash(path), line, goName, describedAs, strings.Join(mismatches, "\n      "))
	}
}

// constantRatio reports whether the Go pays the tree's first printed value divided by the
// same factor at every rank. Only a near-whole factor counts: a tidy 2x or 3x is a value
// that was meant to track the tree and drifted, while an arbitrary ratio is the ordinary
// case of the Go carrying a different quantity from the one the tooltip prints.
func constantRatio(ranks [][]float64, valueAt func(int) (float64, bool)) (float64, bool) {
	// A tooltip that prints several numbers per rank describes several quantities, and the
	// Go value may legitimately track one of the later ones or stack with a second site to
	// reach the first. Paladin Holy Power reads [3 1] — its 2 per point on Holy Shock adds
	// to 1 per point on everything. Too ambiguous to call drift.
	for _, printed := range ranks {
		if len(printed) > 1 {
			return 0, false
		}
	}

	var factor float64
	for rank, printed := range ranks {
		goValue, has := valueAt(rank + 1)
		if !has || len(printed) == 0 || goValue == 0 {
			return 0, false
		}
		ratio := printed[0] / goValue
		if ratio <= 0 {
			return 0, false
		}
		if rank == 0 {
			factor = ratio
			// A factor of 1 would have matched already, and anything under it means the Go
			// pays more than the tooltip prints, which is the ordinary "different quantity"
			// case rather than drift.
			if factor <= 1.01 {
				return 0, false
			}
			continue
		}
		// Proportional, so the same series scaled — not two unrelated numbers that happen
		// to differ by a constant amount.
		if math.Abs(ratio-factor) > factor*0.005 {
			return 0, false
		}
	}
	return factor, factor != 0
}

// matchesDescriptionLiteral reports whether the Go's per-point value is one of the fixed
// numbers the tooltip states outright, in which case it is that number rather than a rank
// value that drifted. Warrior Enrage pays 0.3 per point against a tree reading 2/4/6/8/10,
// which looks like a constant ratio but is the "30% chance" its own text names.
func matchesDescriptionLiteral(goValue float64, literals []float64) bool {
	for _, literal := range literals {
		if withinTolerance(literal, goValue) || withinTolerance(literal, goValue*100) {
			return true
		}
	}
	return false
}

func withinTolerance(a, b float64) bool {
	diff := a - b
	return diff < 0.005 && diff > -0.005
}
