package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Every talent at max, so each talent-gated spell registers and its icon makes it into the
// database. Not a legal build - it spends far more than 51 points - which is fine because
// nothing simulates it.
//
// Read from the tree rather than written down. The string used to be a literal with a comment
// asking whoever changed a tree to regenerate it by hand; nobody did, Shatter went from five
// ranks to three, and the generator panicked on `[]float64{0, 17, 33, 50}[5]` - so `make items`
// could not rebuild the item database at all until someone read a stack trace.
func maxTalents(class string) string {
	type talent struct {
		Location struct {
			RowIdx int `json:"rowIdx"`
			ColIdx int `json:"colIdx"`
		} `json:"location"`
		MaxPoints int `json:"maxPoints"`
	}
	var trees []struct {
		Talents []talent `json:"talents"`
	}

	path := filepath.Join("ui", "core", "talents", "trees", class+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("reading %s: %s", path, err))
	}
	if err := json.Unmarshal(data, &trees); err != nil {
		panic(fmt.Sprintf("parsing %s: %s", path, err))
	}

	parts := make([]string, 0, len(trees))
	for _, tree := range trees {
		talents := tree.Talents
		// The talent string is row-major, which is the order the tree file happens to use but
		// does not promise.
		sort.SliceStable(talents, func(i, j int) bool {
			if talents[i].Location.RowIdx != talents[j].Location.RowIdx {
				return talents[i].Location.RowIdx < talents[j].Location.RowIdx
			}
			return talents[i].Location.ColIdx < talents[j].Location.ColIdx
		})
		var b strings.Builder
		for _, t := range talents {
			fmt.Fprintf(&b, "%d", t.MaxPoints)
		}
		parts = append(parts, b.String())
	}
	return strings.Join(parts, "-")
}
