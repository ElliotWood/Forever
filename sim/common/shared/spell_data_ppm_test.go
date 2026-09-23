package shared

import "testing"

func ppmTestTable() SpellDataTable {
	return SpellDataTable{
		{Rank: 1, SpellID: 100, ProcChance: 101},
		{Rank: 2, SpellID: 200, ProcChance: 101},
	}
}

func TestPPMPerRank(t *testing.T) {
	src := ppmTestTable()
	out := WithSpellDataPPMs(src, map[int32]float64{1: 1.5, 2: 3})

	if got := out.PPMAt(1); got != 1.5 {
		t.Errorf("rank 1 PPM = %v, want 1.5", got)
	}
	if got := out.PPMAt(2); got != 3 {
		t.Errorf("rank 2 PPM = %v, want 3", got)
	}
	if got := out.PPMAt(0); got != 0 {
		t.Errorf("untaken rank PPM = %v, want 0", got)
	}
	if got := src.PPMAt(1); got != 0 {
		t.Errorf("source table was mutated: %v", got)
	}
	if got := out.ProcChanceAt(2); got != 1.01 {
		t.Errorf("proc chance column lost: %v", got)
	}
}

func TestPPMOneValueForEveryRank(t *testing.T) {
	out := WithSpellDataPPM(ppmTestTable(), 2)
	if out.PPMAt(1) != 2 || out.PPMAt(2) != 2 {
		t.Errorf("PPM = %v / %v, want 2 on every rank", out.PPMAt(1), out.PPMAt(2))
	}
}

func TestPPMMissingRankPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("expected a panic for a rank with no PPM")
		}
	}()
	WithSpellDataPPMs(ppmTestTable(), map[int32]float64{1: 1.5})
}

func TestPPMAlreadySetPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("expected a panic for a table that already carries a PPM")
		}
	}()
	WithSpellDataPPM(WithSpellDataPPM(ppmTestTable(), 2), 3)
}
