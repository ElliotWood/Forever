package main

import (
	"fmt"
	"io"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/wowsims/forever/sim/core/spelldata"
)

// The ladder a caller names: `warrior/Execute`, or `Execute` alone when one class states it, with the
// package the caller's file sits in as the tiebreak. A name several classes state and no package to
// choose between them is an error rather than a guess.
func findFamily(index map[string]*ladderFamily, spec, pkg string) (*ladderFamily, error) {
	field := spec
	if class, name, ok := strings.Cut(spec, "/"); ok {
		pkg, field = class, name
	}
	if field == "" {
		return nil, fmt.Errorf("name a ladder, as <class>/<Family> or <Family>")
	}

	if pkg != "" {
		family, ok := index[pkg+"/"+field]
		if !ok {
			return nil, fmt.Errorf("%s states no ladder named %q", pkg, field)
		}
		return family, nil
	}

	var found []*ladderFamily
	for _, family := range index {
		if family.field == field {
			found = append(found, family)
		}
	}
	switch len(found) {
	case 0:
		return nil, fmt.Errorf("no class file states a ladder named %q", field)
	case 1:
		return found[0], nil
	}

	classes := make([]string, 0, len(found))
	for _, family := range found {
		classes = append(classes, family.pkg)
	}
	sort.Strings(classes)
	return nil, fmt.Errorf("%s is a ladder in %s - name one as <class>/%s",
		field, strings.Join(classes, ", "), field)
}

// The ladder calls a class file reaches one rank through, which is what -expr takes: the field, with or
// without the `spellData.` the class files write, and one of the three accessors that names a rank.
var exprPattern = regexp.MustCompile(`^(?:spellData\.)?(\w+)\.(?:(Highest)\(\)|(Rank)\((\d+)\)|(ByID)\((\d+)\))$`)

// The id a ladder call names.
func resolveExpr(index map[string]*ladderFamily, expr, pkg string) (*ladderFamily, int32, error) {
	match := exprPattern.FindStringSubmatch(strings.TrimSpace(expr))
	if match == nil {
		return nil, 0, fmt.Errorf("%q is not a ladder call: write spellData.<Family>.Highest(), .Rank(n) or .ByID(id)", expr)
	}

	family, err := findFamily(index, match[1], pkg)
	if err != nil {
		return nil, 0, err
	}

	switch {
	case match[2] != "":
		return family, family.ranks[len(family.ranks)-1].id, nil

	case match[3] != "":
		n, err := strconv.ParseInt(match[4], 10, 32)
		if err != nil {
			return nil, 0, err
		}
		id, ok := family.rankID(int32(n))
		if !ok {
			return nil, 0, fmt.Errorf("%s spellData.%s has %d ranks, not rank %d",
				family.pkg, family.field, family.rankCount(), n)
		}
		return family, id, nil

	default:
		id, err := strconv.ParseInt(match[6], 10, 32)
		if err != nil {
			return nil, 0, err
		}
		if !family.carries(int32(id)) {
			return nil, 0, fmt.Errorf("%s spellData.%s has no rank with id %d",
				family.pkg, family.field, id)
		}
		return family, int32(id), nil
	}
}

func (f *ladderFamily) rankCount() int32 {
	if f.talentRanks > 0 {
		return f.talentRanks
	}
	return int32(len(f.ranks))
}

// A talent's ranks are all the one spell, so every rank the curve states answers that id.
func (f *ladderFamily) rankID(n int32) (int32, bool) {
	if n <= 0 || n > f.rankCount() {
		return 0, false
	}
	if f.talentRanks > 0 {
		return f.ranks[0].id, true
	}
	return f.ranks[n-1].id, true
}

func (f *ladderFamily) carries(id int32) bool {
	for _, rank := range f.ranks {
		if rank.id == id {
			return true
		}
	}
	return false
}

func (f *ladderFamily) highest() int32 {
	return f.ranks[len(f.ranks)-1].id
}

type familyRankJSON struct {
	ID       int32  `json:"id"`
	Name     string `json:"name"`
	Rank     string `json:"rank"`
	Accessor string `json:"accessor"`
}

type familyJSON struct {
	Family  string           `json:"family"`
	Ranks   []familyRankJSON `json:"ranks"`
	Highest spellJSON        `json:"highest"`
}

// The rank index, with each id's name and rank read out of the store. An id the store does not carry
// says so in place of a name rather than stopping the listing, since the ladder is still worth reading.
func familyRows(f *ladderFamily) []familyRankJSON {
	rows := make([]familyRankJSON, 0, len(f.ranks))
	for _, rank := range f.ranks {
		row := familyRankJSON{ID: rank.id, Accessor: rank.accessor, Name: "not in the store"}
		if s := spelldata.Find(rank.id); s != spelldata.Nil {
			row.Name, row.Rank = s.Name, s.Rank
		}
		rows = append(rows, row)
	}
	return rows
}

func writeFamilyText(out io.Writer, f *ladderFamily) {
	rows := familyRows(f)

	width := 0
	for _, row := range rows {
		width = max(width, len(row.Name))
	}

	fmt.Fprintf(out, "%s spellData.%s\n", f.pkg, f.field)
	for _, row := range rows {
		fmt.Fprintf(out, "%-8d %-*s  %-8s %s\n", row.ID, width, row.Name, row.Rank, row.Accessor)
	}
	fmt.Fprintln(out)
}
