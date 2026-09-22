package main

import (
	"testing"

	"github.com/wowsims/forever/sim/core/spelldata"
)

func TestEffectLines(t *testing.T) {
	cases := []struct {
		name    string
		id      int32
		effect  int
		human   string
		literal string
	}{
		{
			name:    "a tick",
			id:      11574,
			effect:  1,
			human:   "21 physical damage every 3 s to the enemy (7 ticks)",
			literal: "E_APPLY_AURA A_PERIODIC_DAMAGE base=21 period=3000ms target=[6,0]",
		},
		{
			name:    "school damage",
			id:      116,
			effect:  2,
			human:   "19 frost damage to the enemy",
			literal: "E_SCHOOL_DAMAGE base=19 ppl=0.5 variance=0.10526316 sp=0.407 target=[6,0]",
		},
		{
			name:    "a shape no branch words",
			id:      1680,
			effect:  1,
			human:   "",
			literal: "E_NORMALIZED_WEAPON_DMG base=0 sp=1 target=[22,15]",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			lines := effectLines(spelldata.MustFind(c.id))
			if len(lines) < c.effect {
				t.Fatalf("spell %d has %d effects, want at least %d", c.id, len(lines), c.effect)
			}
			line := lines[c.effect-1]
			if line.Human != c.human {
				t.Errorf("human = %q, want %q", line.Human, c.human)
			}
			if line.Literal != c.literal {
				t.Errorf("literal = %q, want %q", line.Literal, c.literal)
			}
		})
	}
}

func TestHeader(t *testing.T) {
	want := []string{
		"school    physical",
		"defense   melee",
		"mechanic  bleed",
		"duration  21 s",
		"gcd       1.5 s",
		"cost      10 rage",
	}
	got := header(spelldata.MustFind(11574))
	if len(got) != len(want) {
		t.Fatalf("header = %q, want %q", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("header line %d = %q, want %q", i, got[i], want[i])
		}
	}
	if title := title(spelldata.MustFind(11574)); title != "11574 Rend (Rank 7)" {
		t.Errorf("title = %q", title)
	}
}

func TestHighestRank(t *testing.T) {
	if s := highestRank(spelldata.ByName("Rend")); s.ID != 11574 {
		t.Errorf("Rend picks %d %s, want 11574 Rank 7", s.ID, s.Rank)
	}
	if s := highestRank(spelldata.ByName("Whirlwind")); s.ID != 1680 {
		t.Errorf("Whirlwind picks %d, want the warrior ability 1680", s.ID)
	}
}

func TestParseArgs(t *testing.T) {
	for _, args := range [][]string{{"11574", "-json"}, {"-json", "11574"}} {
		opts, err := parseArgs(args)
		if err != nil {
			t.Fatalf("parseArgs(%q): %v", args, err)
		}
		if opts.query != "11574" || !opts.json {
			t.Errorf("parseArgs(%q) = %+v", args, opts)
		}
	}
	if _, err := parseArgs(nil); err == nil {
		t.Error("parseArgs with no arguments wants a usage error")
	}
}
