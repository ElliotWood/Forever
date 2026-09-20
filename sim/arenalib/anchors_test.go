package arenalib

// Reshaping is the one piece of the search that moves points without measuring anything, so
// nothing downstream would notice it going wrong - the climb would simply start from a build
// of the wrong shape and report it confidently. Every class, every tree, every anchor.

import (
	"testing"

	"github.com/wowsims/classic/sim/core/proto"
)

func TestReshapeHoldsItsShape(t *testing.T) {
	classes := []proto.Class{
		proto.Class_ClassDruid, proto.Class_ClassHunter, proto.Class_ClassMage,
		proto.Class_ClassPaladin, proto.Class_ClassPriest, proto.Class_ClassRogue,
		proto.Class_ClassShaman, proto.Class_ClassWarlock, proto.Class_ClassWarrior,
	}

	for _, class := range classes {
		trees, err := loadTrees(class)
		if err != nil {
			t.Fatalf("%s: %s", class, err)
		}

		// From an empty build, which is the hardest start: every point has to be placed, and
		// placed somewhere the row gates and prerequisites allow.
		empty := make(allocation, len(trees))
		for i := range trees {
			empty[i] = make([]int, len(trees[i].Talents))
		}

		anchors, err := anchorsFor(class)
		if err != nil {
			t.Fatalf("%s: %s", class, err)
		}
		for _, hold := range anchors {
			shaped, err := reshape(trees, empty, hold)
			if err != nil {
				t.Errorf("%s, %s: %s", class, hold.label, err)
				continue
			}
			if got := sum(shaped[hold.tree]); got != hold.points {
				t.Errorf("%s, %s: held tree has %d points, not %d", class, hold.label, got, hold.points)
			}
			if got := shaped.total(); got != talentBudget {
				t.Errorf("%s, %s: build spends %d points, not %d", class, hold.label, got, talentBudget)
			}
			if !shaped.valid(trees) {
				t.Errorf("%s, %s: produced a build the game would not allow (%s)", class, hold.label, shaped)
			}
		}
	}
}

// Reshaping a build that is already the right shape must leave it alone rather than churn it
// into a different one, or an anchored climb would start somewhere arbitrary every time.
func TestReshapeLeavesAGoodShapeAlone(t *testing.T) {
	trees, err := loadTrees(proto.Class_ClassWarrior)
	if err != nil {
		t.Fatal(err)
	}

	// Fury 17/34/0: 34 points in the Fury tree, 51 in total.
	build := parseTalents(trees, "30305213-550501015050010051")
	shaped, err := reshape(trees, build, anchor{tree: 1, points: 34, label: "34 Fury"})
	if err != nil {
		t.Fatal(err)
	}
	if shaped.String() != build.String() {
		t.Errorf("reshaped a build that already held its shape: %s -> %s", build, shaped)
	}
}
