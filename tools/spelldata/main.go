// Command spelldata prints one client spell row the way the store carries it, so a bare id in
// hand-written code can be read without grepping the generated table.
//
// Run from the repository root:
//
//	go run ./tools/spelldata 11574           # the row, humanised and with the client's own literal
//	go run ./tools/spelldata Whirlwind       # every row with that name, then the highest rank
//	go run ./tools/spelldata Rend -all       # every match in full
//	go run ./tools/spelldata 11574 -json     # the same as JSON, which tools/vscode-spelldata reads
//	go run ./tools/spelldata -family warrior/Execute   # the ladder's ranks, then its highest in full
//	go run ./tools/spelldata -expr 'spellData.Execute.Rank(3)' -package warrior   # the row that reaches
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
	query  string
	family string
	expr   string
	pkg    string
	json   bool
	all    bool
}

// The flags are taken in any position: `11574 -json` is how a caller writes it, and the stdlib flag
// package stops reading flags at the first positional argument. The ones that take a value read it as
// the next argument or after an `=`.
func parseArgs(args []string) (options, error) {
	var opts options
	valued := map[string]*string{"family": &opts.family, "expr": &opts.expr, "package": &opts.pkg}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-json", "--json":
			opts.json = true
			continue
		case "-all", "--all":
			opts.all = true
			continue
		}

		if name, value, ok := strings.Cut(strings.TrimLeft(arg, "-"), "="); ok && strings.HasPrefix(arg, "-") {
			into, known := valued[name]
			if !known {
				return opts, fmt.Errorf("unknown argument %q", arg)
			}
			*into = value
			continue
		}
		if into, known := valued[strings.TrimLeft(arg, "-")]; known && strings.HasPrefix(arg, "-") {
			if i+1 >= len(args) {
				return opts, fmt.Errorf("%s takes a value", arg)
			}
			i++
			*into = args[i]
			continue
		}

		if strings.HasPrefix(arg, "-") {
			return opts, fmt.Errorf("unknown argument %q", arg)
		}
		if opts.query != "" {
			return opts, fmt.Errorf("name one spell, not both %q and %q", opts.query, arg)
		}
		opts.query = arg
	}

	if opts.query != "" && (opts.family != "" || opts.expr != "") || opts.family != "" && opts.expr != "" {
		return opts, fmt.Errorf("ask for one thing: a spell, -family or -expr")
	}
	if opts.query == "" && opts.family == "" && opts.expr == "" {
		return opts, fmt.Errorf("usage: go run ./tools/spelldata <id | name> [-all] [-json] | -family <class>/<Family> | -expr <call> [-package <class>]")
	}
	return opts, nil
}

func run(args []string, out io.Writer) error {
	opts, err := parseArgs(args)
	if err != nil {
		return err
	}

	if opts.family != "" {
		return runFamily(out, opts)
	}
	if opts.expr != "" {
		return runExpr(out, opts)
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

// A ladder as a whole: its rank index, then the highest rank in full, which is the rank most callers
// are after and the one a bare Highest() reaches.
func runFamily(out io.Writer, opts options) error {
	family, err := findFamily(ladderFamilies(), opts.family, opts.pkg)
	if err != nil {
		return err
	}

	s := spelldata.Find(family.highest())
	if s == spelldata.Nil {
		return fmt.Errorf("spell %d is not in the store", family.highest())
	}

	if opts.json {
		return writeJSON(out, familyJSON{
			Family:  family.key(),
			Ranks:   familyRows(family),
			Highest: asJSON(s),
		})
	}

	writeFamilyText(out, family)
	writeText(out, s)
	return nil
}

func runExpr(out io.Writer, opts options) error {
	_, id, err := resolveExpr(ladderFamilies(), opts.expr, opts.pkg)
	if err != nil {
		return err
	}

	s := spelldata.Find(id)
	if s == spelldata.Nil {
		return fmt.Errorf("spell %d is not in the store", id)
	}

	if opts.json {
		return writeJSON(out, exprJSON{spellJSON: asJSON(s), Expr: opts.expr, Resolved: id})
	}

	fmt.Fprintf(out, "%s = %d\n\n", opts.expr, id)
	writeText(out, s)
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

type exprJSON struct {
	spellJSON
	Expr     string `json:"expr"`
	Resolved int32  `json:"resolved"`
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
