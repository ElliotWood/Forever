package arenalib

// The enumeration produces builds nobody wrote and nobody reviews, so the only thing
// standing between it and publishing an illegal build is this.

import (
	"testing"

	"github.com/wowsims/classic/sim/core/proto"
)

func TestEnumerationOnlyProducesLegalBuilds(t *testing.T) {
	// Warrior has the arrows, the row gates and enough talents to be interesting, and is
	// small enough to walk exhaustively in a unit test.
	trees, err := loadTrees(proto.Class_ClassWarrior)
	if err != nil {
		t.Fatal(err)
	}

	// Stand in for a real relevance probe: the first three talents of each row, which gives
	// a spread of shallow and deep, some with prerequisites and some without.
	relevant := map[[2]int]bool{}
	value := map[[2]int]float64{}
	for i, tree := range trees {
		for j := range tree.Talents {
			if j%3 == 0 {
				relevant[[2]int{i, j}] = true
				value[[2]int{i, j}] = float64(j)
			}
		}
	}

	seen, completed, illegal := 0, 0, 0
	enumerate(trees, relevant, value, func(candidate scored) {
		seen++
		built, ok := complete(trees, candidate.points)
		if !ok {
			return
		}
		completed++
		if built.total() != talentBudget {
			illegal++
			if illegal < 4 {
				t.Errorf("completed to a build spending %d points: %s", built.total(), built)
			}
			return
		}
		if !built.valid(trees) {
			illegal++
			if illegal < 4 {
				t.Errorf("completed to a build the game would not allow: %s", built)
			}
		}
	})

	if seen == 0 {
		t.Fatal("enumerated nothing at all, so the test proves nothing")
	}
	if completed == 0 {
		t.Fatal("nothing completed, so the test proves nothing")
	}
	t.Logf("%d assignments, %d completed to builds, %d illegal", seen, completed, illegal)
}

// Completion has to refuse rather than improvise: an assignment that cannot become a legal
// 51 point build must come back as impossible, not as something close.
func TestCompletionRefusesTheImpossible(t *testing.T) {
	trees, err := loadTrees(proto.Class_ClassWarrior)
	if err != nil {
		t.Fatal(err)
	}

	empty := make(allocation, len(trees))
	for i := range trees {
		empty[i] = make([]int, len(trees[i].Talents))
	}

	t.Run("fills an empty build to 51", func(t *testing.T) {
		built, ok := complete(trees, empty)
		if !ok {
			t.Fatal("refused to complete an empty build, which is always possible")
		}
		if built.total() != talentBudget || !built.valid(trees) {
			t.Errorf("completed to something illegal: %s (%d points)", built, built.total())
		}
	})

	t.Run("feeds a prerequisite it was not given", func(t *testing.T) {
		// Find a talent with an arrow into it, take that talent alone, and check the
		// completion maxes its prerequisite rather than publishing an unfed arrow.
		for i, tree := range trees {
			for j, talent := range tree.Talents {
				if talent.PrereqLocation == nil {
					continue
				}
				bare := empty.clone()
				bare[i][j] = talent.MaxPoints
				built, ok := complete(trees, bare)
				if !ok {
					continue
				}
				if !built.valid(trees) {
					t.Errorf("completed %s without feeding its prerequisite: %s", talent.Name, built)
				}
				return
			}
		}
		t.Skip("this class has no prerequisite arrows")
	})
}
