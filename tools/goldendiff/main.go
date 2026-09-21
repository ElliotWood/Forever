// Compares two .results files and says which goldens moved and by how much.
//
//	go run ./tools/goldendiff sim/warrior/dps/TestDpsWarrior.results sim/warrior/dps/TestDpsWarrior.results.tmp
//
// The four numbers a DPS test records are compared per key and sorted by the largest relative move,
// so a change that was meant to touch one ability shows whether it touched anything else. A move
// above the threshold is marked REVIEW and makes the command exit 1, which is the signal that a
// golden needs a reason before it is committed.
package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"sort"

	"github.com/wowsims/forever/sim/core/proto"
	"google.golang.org/protobuf/encoding/prototext"
)

// The relative move a golden is allowed before it is called out. Iteration noise is seeded away in
// the test suite, so anything this side of a rounding difference is a real change.
const reviewThreshold = 0.0005

func main() {
	flag.Parse()
	if flag.NArg() != 2 {
		fmt.Fprintln(os.Stderr, "usage: goldendiff <old.results> <new.results>")
		os.Exit(2)
	}

	old, err := read(flag.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	new, err := read(flag.Arg(1))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	deltas := compare(old.DpsResults, new.DpsResults)
	if report(os.Stdout, deltas) {
		os.Exit(1)
	}
}

func read(path string) (*proto.TestSuiteResult, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	results := &proto.TestSuiteResult{}
	if err := prototext.Unmarshal(data, results); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return results, nil
}

// One metric of one test, as the two files state it.
type delta struct {
	Key    string
	Metric string
	Old    float64
	New    float64

	// Empty where both files carry the key.
	Missing string
}

func (d delta) Absolute() float64 {
	return d.New - d.Old
}

// The move as a share of where it started. A metric that was zero and is not has no share to state,
// so it sorts to the top as an infinite one rather than as no move at all.
func (d delta) Relative() float64 {
	if d.Old == 0 {
		if d.New == 0 {
			return 0
		}
		return math.Inf(1)
	}
	return d.Absolute() / math.Abs(d.Old)
}

func (d delta) Moved() bool {
	return d.Missing != "" || math.Abs(d.Relative()) > reviewThreshold
}

func compare(old, new map[string]*proto.DpsTestResult) []delta {
	keys := map[string]bool{}
	for key := range old {
		keys[key] = true
	}
	for key := range new {
		keys[key] = true
	}

	var out []delta
	for key := range keys {
		before, inOld := old[key]
		after, inNew := new[key]
		switch {
		case !inOld:
			out = append(out, delta{Key: key, Metric: "the test itself", Missing: "only in the new file"})
			continue
		case !inNew:
			out = append(out, delta{Key: key, Metric: "the test itself", Missing: "only in the old file"})
			continue
		}
		for _, m := range []struct {
			name     string
			old, new float64
		}{
			{"DPS", before.Dps, after.Dps},
			{"TPS", before.Tps, after.Tps},
			{"DTPS", before.Dtps, after.Dtps},
			{"HPS", before.Hps, after.Hps},
		} {
			out = append(out, delta{Key: key, Metric: m.name, Old: m.old, New: m.new})
		}
	}

	// Largest move first, and by key and metric under that so two runs of the same pair of files
	// print the same order.
	sort.Slice(out, func(i, j int) bool {
		a, b := math.Abs(out[i].Relative()), math.Abs(out[j].Relative())
		if a != b {
			return a > b
		}
		if out[i].Key != out[j].Key {
			return out[i].Key < out[j].Key
		}
		return out[i].Metric < out[j].Metric
	})
	return out
}

// Prints every metric that differs at all, marking the ones past the threshold, and answers whether
// any was marked. A run where nothing differs prints the totals alone.
func report(w *os.File, deltas []delta) bool {
	flagged := 0
	moved := 0
	for _, d := range deltas {
		if d.Missing != "" {
			flagged++
			fmt.Fprintf(w, "REVIEW %-60s %-16s %s\n", d.Key, d.Metric, d.Missing)
			continue
		}
		if d.Absolute() == 0 {
			continue
		}
		moved++
		mark := "      "
		if d.Moved() {
			flagged++
			mark = "REVIEW"
		}
		fmt.Fprintf(w, "%s %-60s %-16s %12.5f -> %12.5f  %+.5f  %+.4f%%\n",
			mark, d.Key, d.Metric, d.Old, d.New, d.Absolute(), d.Relative()*100)
	}

	fmt.Fprintf(w, "%d metrics compared, %d moved, %d past %.4f%%\n",
		len(deltas), moved, flagged, reviewThreshold*100)
	return flagged > 0
}
