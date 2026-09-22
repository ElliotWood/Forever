package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestFamilyIndex(t *testing.T) {
	var out bytes.Buffer
	if err := run([]string{"-family", "warrior/Execute"}, &out); err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(out.String(), "\n")
	want := []string{
		"warrior spellData.Execute",
		"5308     Execute  Rank 1   Rank(1)",
		"20658    Execute  Rank 2   Rank(2)",
		"20660    Execute  Rank 3   Rank(3)",
		"20661    Execute  Rank 4   Rank(4)",
		"20662    Execute  Rank 5   Highest()",
		"",
		"20662 Execute (Rank 5) warrior spellData.Execute.Highest()",
	}
	for i, line := range want {
		if i >= len(lines) || lines[i] != line {
			t.Fatalf("line %d is %q, want %q", i+1, lines[i], line)
		}
	}
}

// A talent is one spell whose ranks are a curve, so the index is one line naming the rank count.
func TestFamilyIndexTalent(t *testing.T) {
	var out bytes.Buffer
	if err := run([]string{"-family", "warrior/Cruelty"}, &out); err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(out.String(), "\n")
	if lines[1] != "12320    Cruelty           Rank(n), n up to 5" {
		t.Fatalf("the rank index is %q", lines[1])
	}
	if !strings.HasPrefix(lines[3], "12320 Cruelty ") {
		t.Fatalf("the card is %q", lines[3])
	}
}

func TestFamilyWithoutAPackage(t *testing.T) {
	var out bytes.Buffer
	if err := run([]string{"-family", "Execute"}, &out); err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(out.String(), "warrior spellData.Execute\n") {
		t.Fatalf("the heading is %q", strings.SplitN(out.String(), "\n", 2)[0])
	}
}

func TestFamilyJSON(t *testing.T) {
	var out bytes.Buffer
	if err := run([]string{"-family", "warrior/Execute", "-json"}, &out); err != nil {
		t.Fatal(err)
	}

	var got familyJSON
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.Family != "warrior/Execute" {
		t.Errorf("family %q", got.Family)
	}
	if len(got.Ranks) != 5 {
		t.Fatalf("%d ranks", len(got.Ranks))
	}
	if got.Ranks[0] != (familyRankJSON{ID: 5308, Name: "Execute", Rank: "Rank 1", Accessor: "Rank(1)"}) {
		t.Errorf("rank 1 is %+v", got.Ranks[0])
	}
	if got.Ranks[4].Accessor != "Highest()" {
		t.Errorf("the top rank is reached by %q", got.Ranks[4].Accessor)
	}
	if got.Highest.ID != 20662 || len(got.Highest.Effects) == 0 {
		t.Errorf("the highest rank is %+v", got.Highest)
	}
}

func TestExpr(t *testing.T) {
	cases := []struct {
		expr string
		id   int32
	}{
		{"spellData.Execute.Highest()", 20662},
		{"spellData.Execute.Rank(3)", 20660},
		{"spellData.Execute.ByID(20658)", 20658},
		{"Execute.Rank(1)", 5308},
		// Every rank of a talent is the one spell the curve belongs to.
		{"spellData.Cruelty.Rank(2)", 12320},
		{"spellData.Cruelty.Highest()", 12320},
	}

	for _, c := range cases {
		var out bytes.Buffer
		if err := run([]string{"-expr", c.expr, "-package", "warrior", "-json"}, &out); err != nil {
			t.Errorf("%s: %v", c.expr, err)
			continue
		}

		var got exprJSON
		if err := json.Unmarshal(out.Bytes(), &got); err != nil {
			t.Errorf("%s: %v", c.expr, err)
			continue
		}
		if got.Resolved != c.id || got.ID != c.id {
			t.Errorf("%s resolves to %d and prints %d, want %d", c.expr, got.Resolved, got.ID, c.id)
		}
		if got.Expr != c.expr {
			t.Errorf("%s echoes %q", c.expr, got.Expr)
		}
	}
}

func TestExprText(t *testing.T) {
	var out bytes.Buffer
	if err := run([]string{"-expr", "spellData.Execute.Highest()", "-package", "warrior"}, &out); err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(out.String(), "\n")
	if lines[0] != "spellData.Execute.Highest() = 20662" {
		t.Fatalf("the first line is %q", lines[0])
	}
	if !strings.HasPrefix(lines[2], "20662 Execute (Rank 5)") {
		t.Fatalf("the card is %q", lines[2])
	}
}

func TestExprRefused(t *testing.T) {
	cases := []struct {
		args []string
		want string
	}{
		{[]string{"-expr", "spellData.Execute.Rank(9)", "-package", "warrior"}, "has 5 ranks, not rank 9"},
		{[]string{"-expr", "spellData.Execute.ByID(25236)", "-package", "warrior"}, "no rank with id 25236"},
		{[]string{"-expr", "spellData.Nope.Highest()", "-package", "warrior"}, `warrior states no ladder named "Nope"`},
		{[]string{"-expr", "executeRank.EffectN(1)", "-package", "warrior"}, "is not a ladder call"},
		{[]string{"-family", "warrior/Nope"}, `warrior states no ladder named "Nope"`},
		{[]string{"-family", "Execute", "-expr", "spellData.Execute.Highest()"}, "ask for one thing"},
		{[]string{"-family"}, "-family takes a value"},
	}

	for _, c := range cases {
		var out bytes.Buffer
		err := run(c.args, &out)
		if err == nil {
			t.Errorf("%v answered %q", c.args, out.String())
			continue
		}
		if !strings.Contains(err.Error(), c.want) {
			t.Errorf("%v failed with %q, want %q in it", c.args, err, c.want)
		}
	}
}

// The class files a session is editing while this runs are not a fixed index, so the tiebreak between
// classes is tested against one that is.
func TestFamilyAcrossClasses(t *testing.T) {
	index := map[string]*ladderFamily{
		"warrior/Execute": {pkg: "warrior", field: "Execute", ranks: []familyRank{{id: 1, accessor: "Highest()"}}},
		"paladin/Execute": {pkg: "paladin", field: "Execute", ranks: []familyRank{{id: 2, accessor: "Highest()"}}},
		"rogue/Rupture":   {pkg: "rogue", field: "Rupture", ranks: []familyRank{{id: 3, accessor: "Highest()"}}},
	}

	if _, err := findFamily(index, "Execute", ""); err == nil ||
		!strings.Contains(err.Error(), "Execute is a ladder in paladin, warrior") {
		t.Errorf("a name two classes state answered %v", err)
	}
	if family, err := findFamily(index, "Execute", "paladin"); err != nil || family.pkg != "paladin" {
		t.Errorf("the package did not choose: %v", err)
	}
	if family, err := findFamily(index, "warrior/Execute", "paladin"); err != nil || family.pkg != "warrior" {
		t.Errorf("the spec's own class did not win: %v", err)
	}
	if family, err := findFamily(index, "Rupture", ""); err != nil || family.pkg != "rogue" {
		t.Errorf("a name one class states did not answer: %v", err)
	}
	if _, err := findFamily(index, "Whirlwind", ""); err == nil ||
		!strings.Contains(err.Error(), `no class file states a ladder named "Whirlwind"`) {
		t.Errorf("an unknown name answered %v", err)
	}
}
