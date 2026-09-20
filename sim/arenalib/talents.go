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

// Which talents can move this spec's damage at all.
//
// Most of a tree cannot. Defensives, PvP talents, the ones the sim does not implement -
// spending a point there changes nothing, and every pass that prices them spends a run to
// rediscover that. Measuring once and skipping them afterwards roughly halves the cost of a
// step, which is what pays for searching several builds instead of one.
//
// The probe deliberately leaves the legal space: it adds a point on top of the build without
// taking one away, so the character briefly holds 52. The sim does not care - it applies
// whatever the string says - and nothing probed this way is ever published. Only the answer
// "did the number move" comes back out.
//
// Conservative in the direction that matters: a talent is dropped only when the number does
// not move at all, so anything with the faintest effect stays in the search.
func relevantTalents(trees []tree, base allocation, measure func(allocation) float64, baseDps float64) map[[2]int]bool {
	relevant := map[[2]int]bool{}
	for i, tree := range trees {
		for j, t := range tree.Talents {
			if t.NotSimulated {
				continue
			}
			probe := base.clone()
			if probe[i][j] < t.MaxPoints {
				probe[i][j] = t.MaxPoints
			} else {
				probe[i][j] = 0
			}
			if measure(probe) != baseDps {
				relevant[[2]int{i, j}] = true
			}
		}
	}
	return relevant
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
// Holds a tree at a fixed number of points while everything around it is optimised.
//
// The shapes players actually compare: 31 points for a tree's capstone, 21 for the talent
// two rows above it. An unconstrained climb finds one build and says nothing about what the
// spec can do if you commit to a different tree, which is the question anyone choosing a
// build is actually asking.
type anchor struct {
	tree   int
	points int
	label  string
}

func optimise(spec Spec, uiDir string, start TalentBuild, gear string, rotation string, budget int) (TalentBuild, int, error) {
	return optimiseAnchored(spec, uiDir, start, gear, rotation, budget, nil)
}

func optimiseAnchored(spec Spec, uiDir string, start TalentBuild, gear string, rotation string, budget int, hold *anchor) (TalentBuild, int, error) {
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

	// The starting build has to satisfy the anchor before the climb can preserve it, so points
	// are moved into or out of the held tree first, cheapest first, until the total is right.
	if hold != nil {
		var err error
		current, err = reshape(trees, current, *hold)
		if err != nil {
			return start, 0, nil
		}
	}

	// Anything the anchor forbids is simply not a move.
	allowed := func(from allocation, to allocation) bool {
		if hold == nil {
			return true
		}
		return sum(to[hold.tree]) == hold.points
	}

	best := candidate{points: current, dps: measure(current)}
	relevant := relevantTalents(trees, current, measure, best.dps)

	// Spend anything the starting build left on the table before trading points around.
	// Not hypothetical: the mage Frost community build spends 49 of its 51 points, and
	// every run of it on this site has been two points short of a character. A trade-only
	// search would have carried that forward forever, because moving a point never notices
	// that there is a point nobody moved.
	for best.points.total() < talentBudget && runs < budget {
		added := false
		for _, add := range rank(trees, best.points, measure, best.dps, true, relevant) {
			trial := best.points.clone()
			trial[add.tree][add.talent]++
			if !trial.valid(trees) || !allowed(best.points, trial) {
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
		removals := rank(trees, best.points, measure, best.dps, false, relevant)
		additions := rank(trees, best.points, measure, best.dps, true, relevant)

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
				if !trial.valid(trees) || trial.total() != best.points.total() || !allowed(best.points, trial) {
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
func rank(trees []tree, points allocation, measure func(allocation) float64, base float64, adding bool, relevant map[[2]int]bool) []scoredMove {
	moves := []scoredMove{}
	for i, tree := range trees {
		for j := range tree.Talents {
			if !relevant[[2]int{i, j}] {
				continue
			}
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

func sum(points []int) int {
	total := 0
	for _, p := range points {
		total += p
	}
	return total
}

// Moves points into or out of the held tree until it holds exactly what the anchor asks.
//
// Blind on purpose - it takes from wherever is legal rather than measuring, because the
// climb that follows is what decides where points belong. This only has to produce a legal
// build of the right shape for it to start from.
func reshape(trees []tree, points allocation, hold anchor) (allocation, error) {
	out := points.clone()

	for sum(out[hold.tree]) < hold.points {
		placed := false
		for j, t := range trees[hold.tree].Talents {
			if out[hold.tree][j] >= t.MaxPoints {
				continue
			}
			trial := out.clone()
			trial[hold.tree][j]++
			if trial.valid(trees) {
				out, placed = trial, true
				break
			}
		}
		if !placed {
			return nil, fmt.Errorf("cannot reach %d points in %s", hold.points, trees[hold.tree].Name)
		}
	}
	for sum(out[hold.tree]) > hold.points {
		removed := false
		// Deepest first: a point in the last row is never holding another one up.
		for j := len(trees[hold.tree].Talents) - 1; j >= 0; j-- {
			if out[hold.tree][j] == 0 {
				continue
			}
			trial := out.clone()
			trial[hold.tree][j]--
			if trial.valid(trees) {
				out, removed = trial, true
				break
			}
		}
		if !removed {
			return nil, fmt.Errorf("cannot come down to %d points in %s", hold.points, trees[hold.tree].Name)
		}
	}

	// The points taken out of the held tree have to land somewhere, and the ones added to it
	// had to come from somewhere. Balanced blindly for the same reason as above.
	for out.total() > talentBudget {
		removed := false
		for i := range trees {
			if hold.tree == i {
				continue
			}
			for j := len(trees[i].Talents) - 1; j >= 0; j-- {
				if out[i][j] == 0 {
					continue
				}
				trial := out.clone()
				trial[i][j]--
				if trial.valid(trees) {
					out, removed = trial, true
					break
				}
			}
			if removed {
				break
			}
		}
		if !removed {
			return nil, fmt.Errorf("cannot come down to %d points", talentBudget)
		}
	}
	for out.total() < talentBudget {
		placed := false
		for i := range trees {
			if hold.tree == i {
				continue
			}
			for j, t := range trees[i].Talents {
				if out[i][j] >= t.MaxPoints {
					continue
				}
				trial := out.clone()
				trial[i][j]++
				if trial.valid(trees) {
					out, placed = trial, true
					break
				}
			}
			if placed {
				break
			}
		}
		if !placed {
			return nil, fmt.Errorf("cannot reach %d points", talentBudget)
		}
	}

	if !out.valid(trees) || out.total() != talentBudget || sum(out[hold.tree]) != hold.points {
		return nil, fmt.Errorf("reshaping produced a build that does not hold")
	}
	return out, nil
}
