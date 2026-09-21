package main

import (
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/wowsims/forever/sim/core/proto"
)

const oldResults = `dps_results: {
 key: "TestWarrior-Settings-Fury"
 value: {
  dps: 1000
  tps: 800
  dtps: 10
 }
}
dps_results: {
 key: "TestWarrior-Settings-Prot"
 value: {
  dps: 500
  tps: 900
  dtps: 200
  hps: 5
 }
}
`

const newResults = `dps_results: {
 key: "TestWarrior-Settings-Fury"
 value: {
  dps: 1000
  tps: 800
  dtps: 10
 }
}
dps_results: {
 key: "TestWarrior-Settings-Prot"
 value: {
  dps: 501
  tps: 900
  dtps: 200
  hps: 5
 }
}
`

func write(t *testing.T, name, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(body), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestReadsAPrototextResultsFile(t *testing.T) {
	results, err := read(write(t, "old.results", oldResults))
	if err != nil {
		t.Fatal(err)
	}
	if got := len(results.DpsResults); got != 2 {
		t.Fatalf("read %d tests, want 2", got)
	}
	if got := results.DpsResults["TestWarrior-Settings-Prot"].Dps; got != 500 {
		t.Errorf("Prot DPS %v, want 500", got)
	}
}

func TestOneMovedMetricIsFlaggedAndTheRestAreNot(t *testing.T) {
	old, err := read(write(t, "old.results", oldResults))
	if err != nil {
		t.Fatal(err)
	}
	new, err := read(write(t, "new.results", newResults))
	if err != nil {
		t.Fatal(err)
	}

	deltas := compare(old.DpsResults, new.DpsResults)
	if len(deltas) != 8 {
		t.Fatalf("compared %d metrics, want 8", len(deltas))
	}

	top := deltas[0]
	if top.Key != "TestWarrior-Settings-Prot" || top.Metric != "DPS" {
		t.Fatalf("largest move is %s %s, want the Prot DPS", top.Key, top.Metric)
	}
	if !top.Moved() {
		t.Errorf("a move of %v%% is not flagged", top.Relative()*100)
	}
	for _, d := range deltas[1:] {
		if d.Absolute() != 0 {
			t.Errorf("%s %s moved by %v, and only the Prot DPS should have", d.Key, d.Metric, d.Absolute())
		}
		if d.Moved() {
			t.Errorf("%s %s is flagged and did not move", d.Key, d.Metric)
		}
	}
}

func TestAMoveUnderTheThresholdIsReportedAndNotFlagged(t *testing.T) {
	old := map[string]*proto.DpsTestResult{"A": {Dps: 1000}}
	new := map[string]*proto.DpsTestResult{"A": {Dps: 1000.001}}

	deltas := compare(old, new)
	if deltas[0].Metric != "DPS" || deltas[0].Absolute() == 0 {
		t.Fatalf("the DPS move is not the first delta: %+v", deltas[0])
	}
	if deltas[0].Moved() {
		t.Errorf("a move of %v%% is flagged, and the threshold is %v%%", deltas[0].Relative()*100, reviewThreshold*100)
	}
}

func TestAMetricThatWasZeroSortsFirst(t *testing.T) {
	old := map[string]*proto.DpsTestResult{"A": {Dps: 1000}}
	new := map[string]*proto.DpsTestResult{"A": {Dps: 1000, Dtps: 7}}

	deltas := compare(old, new)
	if deltas[0].Metric != "DTPS" || !math.IsInf(deltas[0].Relative(), 1) {
		t.Fatalf("a metric that gained a value is not the largest move: %+v", deltas[0])
	}
	if !deltas[0].Moved() {
		t.Error("a metric that gained a value is not flagged")
	}
}

func TestAKeyInOnlyOneFileIsFlagged(t *testing.T) {
	old := map[string]*proto.DpsTestResult{"A": {Dps: 1}}
	new := map[string]*proto.DpsTestResult{"B": {Dps: 1}}

	deltas := compare(old, new)
	if len(deltas) != 2 {
		t.Fatalf("compared %d entries, want the two unmatched keys", len(deltas))
	}
	for _, d := range deltas {
		if d.Missing == "" || !d.Moved() {
			t.Errorf("%s is not flagged as unmatched: %+v", d.Key, d)
		}
	}
}
