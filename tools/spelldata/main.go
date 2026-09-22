// Command spelldata prints one client spell row the way the store carries it, so a bare id in
// hand-written code can be read without grepping the generated table.
//
// Run from the repository root:
//
//	go run ./tools/spelldata 11574           # the row, humanised and with the client's own literal
//	go run ./tools/spelldata Whirlwind       # every row with that name, then the highest rank
//	go run ./tools/spelldata Rend -all       # every match in full
//	go run ./tools/spelldata 11574 -json     # the same as JSON, which tools/vscode-spelldata reads
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/wowsims/forever/sim/core/spelldata"
)

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "spelldata:", err)
		os.Exit(1)
	}
}

type options struct {
	query string
	json  bool
	all   bool
}

// The flags are taken in any position: `11574 -json` is how a caller writes it, and the stdlib flag
// package stops reading flags at the first positional argument.
func parseArgs(args []string) (options, error) {
	var opts options
	for _, arg := range args {
		switch arg {
		case "-json", "--json":
			opts.json = true
		case "-all", "--all":
			opts.all = true
		default:
			if strings.HasPrefix(arg, "-") {
				return opts, fmt.Errorf("unknown argument %q", arg)
			}
			if opts.query != "" {
				return opts, fmt.Errorf("name one spell, not both %q and %q", opts.query, arg)
			}
			opts.query = arg
		}
	}
	if opts.query == "" {
		return opts, fmt.Errorf("usage: go run ./tools/spelldata <id | name> [-all] [-json]")
	}
	return opts, nil
}

func run(args []string, out io.Writer) error {
	opts, err := parseArgs(args)
	if err != nil {
		return err
	}

	if id, err := strconv.ParseInt(opts.query, 10, 32); err == nil {
		s := spelldata.Find(int32(id))
		if s == spelldata.Nil {
			return fmt.Errorf("spell %d is not in the store", id)
		}
		if opts.json {
			return writeJSON(out, asJSON(s))
		}
		writeText(out, s)
		return nil
	}

	matches := spelldata.ByName(opts.query)
	if len(matches) == 0 {
		return fmt.Errorf("no spell is named %q", opts.query)
	}

	picked := matches
	if !opts.all {
		picked = []*spelldata.Spell{highestRank(matches)}
	}

	if opts.json {
		rows := make([]spellJSON, 0, len(picked))
		for _, s := range picked {
			rows = append(rows, asJSON(s))
		}
		return writeJSON(out, rows)
	}

	if len(matches) > 1 {
		for _, s := range matches {
			fmt.Fprintf(out, "%-8d %s  %s\n", s.ID, s.Name, s.Rank)
		}
		fmt.Fprintln(out)
	}
	for i, s := range picked {
		if i > 0 {
			fmt.Fprintln(out)
		}
		writeText(out, s)
	}
	return nil
}

// The rank a caller means out of several rows with one name: the highest one the client states, and
// the lowest id among rows that state no rank, which is the ability itself ahead of the item and set
// rows that share its name.
func highestRank(matches []*spelldata.Spell) *spelldata.Spell {
	best := matches[0]
	for _, s := range matches[1:] {
		if s.RankNumber() > best.RankNumber() {
			best = s
		}
	}
	return best
}

func writeText(out io.Writer, s *spelldata.Spell) {
	fmt.Fprintln(out, join(title(s), strings.Join(ladderRefs(s.ID), "  ")))
	for _, line := range header(s) {
		fmt.Fprintln(out, line)
	}
	if proc := procSummary(s); proc != "" {
		fmt.Fprintf(out, "%-9s %s\n", "proc", proc)
	}
	if refs := refList(s); len(refs) > 0 {
		fmt.Fprintf(out, "%-9s %s\n", "refs", strings.Join(refs, ", "))
	}

	effects := effectLines(s)
	if len(effects) > 0 {
		fmt.Fprintln(out)
	}
	for i, line := range effects {
		fmt.Fprintf(out, "effect %-2d %s\n", i+1, line.Human)
		fmt.Fprintf(out, "%9s %s\n", "", line.Literal)
	}

	fmt.Fprintln(out)
	fmt.Fprintln(out, wowheadURL(s.ID))
}

type spellJSON struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Rank string `json:"rank"`

	// The heading the text form prints, and the ladder calls a class file reaches this id through.
	Title  string   `json:"title"`
	Ladder []string `json:"ladder"`

	Header  []string `json:"header"`
	Effects []Line   `json:"effects"`

	// Empty on a row that is not a proc and on one the tooltip names no spell from.
	Proc string   `json:"proc"`
	Refs []string `json:"refs"`

	Wowhead string `json:"wowhead"`
}

func asJSON(s *spelldata.Spell) spellJSON {
	effects := effectLines(s)
	if effects == nil {
		effects = []Line{}
	}
	ladder := ladderRefs(s.ID)
	if ladder == nil {
		ladder = []string{}
	}
	return spellJSON{
		ID:      s.ID,
		Name:    s.Name,
		Rank:    s.Rank,
		Title:   title(s),
		Ladder:  ladder,
		Header:  header(s),
		Effects: effects,
		Proc:    procSummary(s),
		Refs:    refList(s),
		Wowhead: wowheadURL(s.ID),
	}
}

func writeJSON(out io.Writer, v any) error {
	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
