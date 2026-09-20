package arenalib

// Searching the talent trees, because enumerating them is not a thing anyone can do.
//
// There are 89,776,730,783,606,094 builds a warrior could actually spend - counted by
// tools/talents/count_builds.py from the trees themselves, enforcing rank caps, row gates
// and the prerequisite arrows. Legality matters and does not rescue this: it cuts the count
// about fourfold from the same walk without the arrows, leaving two billion years at a
// second a build.
//
// Nor does restricting it to builds anyone would run. Count only the all-or-nothing ones,
// every talent maxed or untouched, which is roughly what an optimal build looks like, and a
// warrior still has 57,341,667 of them and a mage 1,261,940,421. Eighteen months and forty
// years respectively. "Try every configuration" is not a large job, it is an impossible one,
// and the honest response is not to run a subset and call it the answer.
//
// So: climb instead. Start from a build somebody already believed in, and repeatedly ask
// what one point is worth. Take the point whose removal costs least, give it to the talent
// whose addition gains most, and keep going while that trade is positive. It finds a local
// maximum rather than the global one, and the difference matters least exactly where people
// care most - near a build the community has already converged on.
//
// Two things keep it honest. The marginal pass is an approximation, because talents interact
// and measuring them one at a time cannot see that; so the combined move is measured for
// real before it is accepted, and rejected if the pair is worth less than the parts
// suggested. And every run shares a random seed, so two builds are compared over the same
// stream of rolls and the difference between them is signal rather than noise.

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/wowsims/classic/sim/core/proto"
)

// Fewer iterations while searching than when reporting: a thousand comparisons at survey
// quality beats twenty at publication quality, and the winner is re-run properly at the end.
const searchIterations = int32(2000)

// A move has to be worth more than this to be taken. Below it the trade is inside the noise
// two 2000 iteration runs can produce even on a shared seed, and chasing it would walk the
// build sideways for hours.
const minimumGain = 0.001

// A character at 60 has 51 talent points. A build holding fewer is not a choice, it is a
// mistake, and the search fixes it rather than preserving it.
const talentBudget = 51

// How many steps the climb may take. Each step costs roughly one sim run per talent, so this
// is the knob that decides whether the job finishes inside its CI timeout.
const maxSteps = 20

type talent struct {
	Name      string `json:"name"`
	MaxPoints int    `json:"maxPoints"`
	Location  struct {
		RowIdx int `json:"rowIdx"`
		ColIdx int `json:"colIdx"`
	} `json:"location"`
	PrereqLocation *struct {
		RowIdx int `json:"rowIdx"`
		ColIdx int `json:"colIdx"`
	} `json:"prereqLocation"`
	// Talents the sim does not implement. Points here are legal in game and worth nothing
	// here, so the search will always empty them - which is a fact about the sim, not advice.
	NotSimulated bool `json:"notSimulated"`
}

type tree struct {
	Name    string   `json:"name"`
	Talents []talent `json:"talents"`
}

// A build as points per talent, in the same row-major order the talents string uses.
type allocation [][]int

func loadTrees(class proto.Class) ([]tree, error) {
	// Class_ClassWarrior -> warrior.
	name := strings.ToLower(strings.TrimPrefix(class.String(), "Class_Class"))
	if idx := strings.LastIndex(class.String(), "Class"); idx >= 0 {
		name = strings.ToLower(class.String()[idx+len("Class"):])
	}

	for _, dir := range []string{
		filepath.Join("..", "..", "..", "ui", "core", "talents", "trees"),
		filepath.Join("..", "..", "ui", "core", "talents", "trees"),
	} {
		data, err := os.ReadFile(filepath.Join(dir, name+".json"))
		if err != nil {
			continue
		}
		var trees []tree
		if err := json.Unmarshal(data, &trees); err != nil {
			return nil, err
		}
		for i := range trees {
			sort.SliceStable(trees[i].Talents, func(a, b int) bool {
				left, right := trees[i].Talents[a].Location, trees[i].Talents[b].Location
				if left.RowIdx != right.RowIdx {
					return left.RowIdx < right.RowIdx
				}
				return left.ColIdx < right.ColIdx
			})
		}
		return trees, nil
	}
	return nil, fmt.Errorf("no talent tree for %s", name)
}

func parseTalents(trees []tree, talents string) allocation {
	points := make(allocation, len(trees))
	parts := strings.Split(talents, "-")
	for i := range trees {
		points[i] = make([]int, len(trees[i].Talents))
		if i >= len(parts) {
			continue
		}
		for j, char := range parts[i] {
			if j >= len(points[i]) {
				break
			}
			if char >= '0' && char <= '9' {
				points[i][j] = int(char - '0')
			}
		}
	}
	return points
}

// The inverse, in the form the sim reads. Trailing zeros are trimmed the way every talent
// calculator writes them, so an optimised build can be pasted straight into the site.
func (points allocation) String() string {
	parts := make([]string, len(points))
	for i, tree := range points {
		var sb strings.Builder
		for _, p := range tree {
			sb.WriteByte(byte('0' + p))
		}
		parts[i] = strings.TrimRight(sb.String(), "0")
	}
	// A talents string keeps its separators even when the later trees are empty, otherwise
	// the second tree's points would be read as the first tree's.
	for len(parts) > 0 && parts[len(parts)-1] == "" {
		parts = parts[:len(parts)-1]
	}
	return strings.Join(parts, "-")
}

func (points allocation) total() int {
	sum := 0
	for _, tree := range points {
		for _, p := range tree {
			sum += p
		}
	}
	return sum
}

func (points allocation) clone() allocation {
	out := make(allocation, len(points))
	for i, tree := range points {
		out[i] = append([]int(nil), tree...)
	}
	return out
}

// Whether the game would let somebody spend points this way: every talent inside its rank
// cap, every row reached by five points per row already spent in that tree, and every talent
// with an arrow into it fed by a maxed prerequisite.
func (points allocation) valid(trees []tree) bool {
	for i, tree := range trees {
		byLocation := map[[2]int]int{}
		for j, t := range tree.Talents {
			byLocation[[2]int{t.Location.RowIdx, t.Location.ColIdx}] = j
		}

		// Points in a row are only legal once enough points sit in the rows above it, and
		// "above" means anywhere in the tree, not on the path.
		above := map[int]int{}
		for j, t := range tree.Talents {
			above[t.Location.RowIdx] += points[i][j]
		}

		for j, t := range tree.Talents {
			if points[i][j] < 0 || points[i][j] > t.MaxPoints {
				return false
			}
			if points[i][j] == 0 {
				continue
			}

			earlier := 0
			for row, p := range above {
				if row < t.Location.RowIdx {
					earlier += p
				}
			}
			if earlier < 5*t.Location.RowIdx {
				return false
			}

			if t.PrereqLocation != nil {
				prereq, ok := byLocation[[2]int{t.PrereqLocation.RowIdx, t.PrereqLocation.ColIdx}]
				if !ok || points[i][prereq] < tree.Talents[prereq].MaxPoints {
					return false
				}
			}
		}
	}
	return true
}

// One evaluated build.
type candidate struct {
	points allocation
	dps    float64
}

// Optimise climbs from a starting build, returning the best it reached and how many sim runs
// it took to get there.
//
// The shape of each step: price every point that could be taken out, price every point that
// could be put in, then actually run the best-looking trade rather than trusting the two
// halves to add up. Talents interact - a point in Flurry is worth more next to Enrage - and
// the separable estimate is only a way of deciding what to measure properly.
func optimise(spec Spec, uiDir string, start TalentBuild, gear string, rotation string, budget int) (TalentBuild, int, error) {
	trees, err := loadTrees(spec.Class)
	if err != nil {
		return start, 0, err
	}

	current := parseTalents(trees, start.Talents)
	if !current.valid(trees) {
		return start, 0, fmt.Errorf("%s is not a legal build to start from", start.Name)
	}

	runs := 0
	measure := func(points allocation) float64 {
		runs++
		return runAt(spec, uiDir, TalentBuild{Talents: points.String()}, gear, rotation, searchIterations).Dps
	}

	best := candidate{points: current, dps: measure(current)}

	// Spend anything the starting build left on the table before trading points around.
	// Not hypothetical: the mage Frost community build spends 49 of its 51 points, and
	// every run of it on this site has been two points short of a character. A trade-only
	// search would have carried that forward forever, because moving a point never notices
	// that there is a point nobody moved.
	for best.points.total() < talentBudget && runs < budget {
		added := false
		for _, add := range rank(trees, best.points, measure, best.dps, true) {
			trial := best.points.clone()
			trial[add.tree][add.talent]++
			if !trial.valid(trees) {
				continue
			}
			best = candidate{points: trial, dps: measure(trial)}
			added = true
			break
		}
		if !added {
			break
		}
	}

	for step := 0; step < maxSteps && runs < budget; step++ {
		removals := rank(trees, best.points, measure, best.dps, false)
		additions := rank(trees, best.points, measure, best.dps, true)

		// Only the few most promising pairs are run for real. The estimate is good enough to
		// rank candidates and not good enough to accept one.
		improved := false
		for _, add := range topN(additions, 3) {
			for _, remove := range topN(removals, 3) {
				if runs >= budget {
					break
				}
				if add.tree == remove.tree && add.talent == remove.talent {
					continue
				}
				if add.delta+remove.delta <= 0 {
					continue
				}
				trial := best.points.clone()
				trial[remove.tree][remove.talent]--
				trial[add.tree][add.talent]++
				if !trial.valid(trees) || trial.total() != best.points.total() {
					continue
				}
				dps := measure(trial)
				if dps > best.dps*(1+minimumGain) {
					best = candidate{points: trial, dps: dps}
					improved = true
					break
				}
			}
			if improved {
				break
			}
		}
		if !improved {
			break
		}
	}

	name := start.Name
	if best.points.String() != start.Talents {
		name = fmt.Sprintf("%s, optimised", start.Name)
	}
	return TalentBuild{Name: name, Talents: best.points.String()}, runs, nil
}

type scoredMove struct {
	tree, talent int
	delta        float64
}

// Prices every legal single-point change of one direction, best first.
func rank(trees []tree, points allocation, measure func(allocation) float64, base float64, adding bool) []scoredMove {
	moves := []scoredMove{}
	for i, tree := range trees {
		for j := range tree.Talents {
			trial := points.clone()
			if adding {
				if points[i][j] >= tree.Talents[j].MaxPoints {
					continue
				}
				trial[i][j]++
			} else {
				if points[i][j] == 0 {
					continue
				}
				trial[i][j]--
			}
			if !trial.valid(trees) {
				continue
			}
			moves = append(moves, scoredMove{i, j, measure(trial) - base})
		}
	}
	sort.Slice(moves, func(a, b int) bool { return moves[a].delta > moves[b].delta })
	return moves
}

func topN[T any](items []T, n int) []T {
	if len(items) < n {
		return items
	}
	return items[:n]
}
