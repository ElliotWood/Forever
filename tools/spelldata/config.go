package main

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/parser"
	"go/scanner"
	"go/token"
	"math/bits"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/spelldata"
)

var spellConfigPattern = regexp.MustCompile(`\bspelldata\.SpellConfig\b`)

func exactName(typeName string, value uint64) string {
	for _, c := range coreConstants(typeName) {
		if c.value == value && !strings.HasSuffix(c.name, "Len") {
			return c.name
		}
	}
	return strconv.FormatUint(value, 10)
}

func bitNames(typeName string, value uint64) []string {
	var out []string
	for value != 0 {
		bit := value & -value
		value &^= bit
		name := fmt.Sprintf("bit %d", bits.TrailingZeros64(bit))
		for _, c := range coreConstants(typeName) {
			if c.value == bit {
				name = c.name
				break
			}
		}
		out = append(out, name)
	}
	return out
}

type configOption struct {
	label  string
	source string
	opt    spelldata.SpellOpt
	err    error
}

var optionBuilders = map[string]struct {
	args  int
	build func(args []uint64) spelldata.SpellOpt
}{
	"Melee": {1, func(a []uint64) spelldata.SpellOpt { return spelldata.Melee(core.ProcMask(a[0])) }},
	"Magic": {1, func(a []uint64) spelldata.SpellOpt { return spelldata.Magic(core.ProcMask(a[0])) }},
	"Proc":  {0, func([]uint64) spelldata.SpellOpt { return spelldata.Proc() }},
	"Flags": {1, func(a []uint64) spelldata.SpellOpt { return spelldata.Flags(core.SpellFlag(a[0])) }},
	"Tag":   {1, func(a []uint64) spelldata.SpellOpt { return spelldata.Tag(int32(a[0])) }},
}

func readOption(expr ast.Expr) configOption {
	source := nodeText(expr)
	option := configOption{label: source, source: source}

	call, ok := expr.(*ast.CallExpr)
	var sel *ast.SelectorExpr
	if ok {
		sel, ok = call.Fun.(*ast.SelectorExpr)
	}
	if ok {
		pkg, isIdent := sel.X.(*ast.Ident)
		ok = isIdent && pkg.Name == "spelldata"
	}
	if !ok {
		option.err = fmt.Errorf("not a spelldata option call")
		return option
	}

	args := make([]string, 0, len(call.Args))
	for _, arg := range call.Args {
		args = append(args, strings.ReplaceAll(nodeText(arg), "core.", ""))
	}
	option.label = sel.Sel.Name + "(" + strings.Join(args, " | ") + ")"
	if len(option.label) > 48 {
		option.label = sel.Sel.Name + "(…)"
	}

	builder, known := optionBuilders[sel.Sel.Name]
	if !known {
		option.err = fmt.Errorf("spelldata.%s is not an option this reads", sel.Sel.Name)
		return option
	}
	if len(call.Args) != builder.args {
		option.err = fmt.Errorf("spelldata.%s takes %s", sel.Sel.Name, arguments(builder.args))
		return option
	}

	values := make([]uint64, 0, len(call.Args))
	for _, arg := range call.Args {
		value, err := evalConst(arg)
		if err != nil {
			option.err = err
			return option
		}
		v, exact := constant.Uint64Val(constant.ToInt(value))
		if !exact {
			i, ok := constant.Int64Val(constant.ToInt(value))
			if !ok {
				option.err = fmt.Errorf("%s is not an integer", nodeText(arg))
				return option
			}
			v = uint64(i)
		}
		values = append(values, v)
	}
	option.opt = builder.build(values)
	return option
}

func callText(text string, start int) (string, error) {
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(text)-start)
	s.Init(file, []byte(text[start:]), nil, scanner.ScanComments)

	depth := 0
	for {
		pos, tok, _ := s.Scan()
		switch tok {
		case token.EOF:
			return "", fmt.Errorf("the SpellConfig call is not closed")
		case token.LPAREN:
			depth++
		case token.RPAREN:
			depth--
			if depth == 0 {
				return text[start : start+file.Offset(pos)+1], nil
			}
		}
	}
}

type configRow struct {
	field string
	value string
	from  string
}

type configResult struct {
	spell   *spelldata.Spell
	pick    string
	rows    []configRow
	skipped []configOption
}

func evalSpellConfig(call string, declarations map[string]declaration, pkg string, trace *tracer) (*configResult, error) {
	node, err := parser.ParseExpr(call)
	if err != nil {
		return nil, fmt.Errorf("%q is not a SpellConfig call", call)
	}
	expr, ok := node.(*ast.CallExpr)
	if !ok || len(expr.Args) < 2 {
		return nil, fmt.Errorf("%q is not a SpellConfig call with a row", call)
	}

	s, pick, err := configPick(expr.Args[1], declarations, pkg, trace)
	if err != nil {
		return nil, err
	}

	result := &configResult{spell: s, pick: pick}
	labels := []string{"row"}
	var opts []spelldata.SpellOpt
	for _, arg := range expr.Args[2:] {
		option := readOption(arg)
		if option.err != nil {
			trace.add("  option %s unevaluated: %v", option.source, option.err)
			result.skipped = append(result.skipped, option)
			continue
		}
		trace.add("  option %s", option.label)
		opts = append(opts, option.opt)
		labels = append(labels, option.label)
	}

	stages := make([]core.SpellConfig, 0, len(opts)+1)
	for i := 0; i <= len(opts); i++ {
		stages = append(stages, spelldata.SpellConfig(&core.Unit{}, s, opts[:i]...))
	}
	result.rows = attribute(stages, labels)
	return result, nil
}

func configPick(arg ast.Expr, declarations map[string]declaration, pkg string, trace *tracer) (*spelldata.Spell, string, error) {
	if id, ok := findCall(arg); ok {
		s := spelldata.Find(id)
		if s == spelldata.Nil {
			return nil, "", fmt.Errorf("spell %d is not in the store", id)
		}
		return s, nodeText(arg), nil
	}

	trace.add("  row %s", nodeText(arg))
	c, err := walkChain(arg, func(token.Pos) int { return 0 })
	if err != nil {
		return nil, "", err
	}
	c, err = resolveChain(c, declarations, trace)
	if err != nil {
		return nil, "", err
	}
	result, err := evalExpr(ladderFamilies(), c, pkg)
	if err != nil {
		return nil, "", err
	}
	if result.kind != kindSpell {
		return nil, "", fmt.Errorf("%s reads a %s, not a row", nodeText(arg), result.kind)
	}
	return result.spell, c.text(false), nil
}

// The id of a `spelldata.MustFind(id)` or `spelldata.Find(id)`.
func findCall(arg ast.Expr) (int32, bool) {
	call, ok := arg.(*ast.CallExpr)
	if !ok || len(call.Args) != 1 {
		return 0, false
	}
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "MustFind" && sel.Sel.Name != "Find" {
		return 0, false
	}
	if pkg, ok := sel.X.(*ast.Ident); !ok || pkg.Name != "spelldata" {
		return 0, false
	}
	id, err := evalInt(call.Args[0])
	return int32(id), err == nil
}

type configField struct {
	name  string
	read  func(*core.SpellConfig) string
	flags string
	bits  func(*core.SpellConfig) uint64
}

func nonZero[T comparable](v T, format func(T) string) string {
	var zero T
	if v == zero {
		return ""
	}
	return format(v)
}

func floatText(v float64) string          { return strconv.FormatFloat(v, 'f', -1, 64) }
func durationText(v time.Duration) string { return v.String() }
func intText[T int32 | int](v T) string   { return strconv.Itoa(int(v)) }
func boolText(bool) string                { return "true" }

var configFields = []configField{
	{name: "ActionID", read: func(c *core.SpellConfig) string {
		if c.ActionID.Tag != 0 {
			return fmt.Sprintf("SpellID %d, Tag %d", c.ActionID.SpellID, c.ActionID.Tag)
		}
		return nonZero(c.ActionID.SpellID, func(id int32) string { return fmt.Sprintf("SpellID %d", id) })
	}},
	{name: "Rank", read: func(c *core.SpellConfig) string { return nonZero(c.Rank, intText[int32]) }},
	{name: "SpellSchool", read: func(c *core.SpellConfig) string {
		return nonZero(c.SpellSchool, func(v core.SpellSchool) string { return exactName("SpellSchool", uint64(v)) })
	}},
	{name: "DefenseType", read: func(c *core.SpellConfig) string {
		return nonZero(c.DefenseType, func(v core.DefenseType) string { return exactName("DefenseType", uint64(v)) })
	}},
	{name: "Flags", flags: "SpellFlag", bits: func(c *core.SpellConfig) uint64 { return uint64(c.Flags) }},
	{name: "ProcMask", flags: "ProcMask", bits: func(c *core.SpellConfig) uint64 { return uint64(c.ProcMask) }},
	{name: "Cast.DefaultCast.CastTime", read: func(c *core.SpellConfig) string { return nonZero(c.Cast.DefaultCast.CastTime, durationText) }},
	{name: "Cast.DefaultCast.GCD", read: func(c *core.SpellConfig) string { return nonZero(c.Cast.DefaultCast.GCD, durationText) }},
	{name: "Cast.DefaultCast.NonEmpty", read: func(c *core.SpellConfig) string { return nonZero(c.Cast.DefaultCast.NonEmpty, boolText) }},
	{name: "Cast.IgnoreHaste", read: func(c *core.SpellConfig) string { return nonZero(c.Cast.IgnoreHaste, boolText) }},
	{name: "Cast.CD", read: func(c *core.SpellConfig) string { return nonZero(c.Cast.CD.Duration, durationText) }},
	{name: "Cast.SharedCD", read: func(c *core.SpellConfig) string { return nonZero(c.Cast.SharedCD.Duration, durationText) }},
	{name: "ManaCost.FlatCost", read: func(c *core.SpellConfig) string { return nonZero(c.ManaCost.FlatCost, intText[int32]) }},
	{name: "ManaCost.BaseCostPercent", read: func(c *core.SpellConfig) string { return nonZero(c.ManaCost.BaseCostPercent, floatText) }},
	{name: "RageCost.Cost", read: func(c *core.SpellConfig) string { return nonZero(c.RageCost.Cost, intText[int32]) }},
	{name: "RageCost.Refund", read: func(c *core.SpellConfig) string { return nonZero(c.RageCost.Refund, floatText) }},
	{name: "EnergyCost.Cost", read: func(c *core.SpellConfig) string { return nonZero(c.EnergyCost.Cost, intText[int32]) }},
	{name: "EnergyCost.Refund", read: func(c *core.SpellConfig) string { return nonZero(c.EnergyCost.Refund, floatText) }},
	{name: "FocusCost.Cost", read: func(c *core.SpellConfig) string { return nonZero(c.FocusCost.Cost, intText[int32]) }},
	{name: "FocusCost.Refund", read: func(c *core.SpellConfig) string { return nonZero(c.FocusCost.Refund, floatText) }},
	{name: "DamageMultiplier", read: func(c *core.SpellConfig) string { return nonZero(c.DamageMultiplier, floatText) }},
	{name: "ThreatMultiplier", read: func(c *core.SpellConfig) string { return nonZero(c.ThreatMultiplier, floatText) }},
	{name: "BonusCoefficient", read: func(c *core.SpellConfig) string { return nonZero(c.BonusCoefficient, floatText) }},
	{name: "MinRange", read: func(c *core.SpellConfig) string { return nonZero(c.MinRange, floatText) }},
	{name: "MaxRange", read: func(c *core.SpellConfig) string { return nonZero(c.MaxRange, floatText) }},
	{name: "MissileSpeed", read: func(c *core.SpellConfig) string { return nonZero(c.MissileSpeed, floatText) }},
	{name: "ClassFlags", read: func(c *core.SpellConfig) string { return classFlagsPhrase(c.ClassFlags) }},
}

func attribute(stages []core.SpellConfig, labels []string) []configRow {
	final := &stages[len(stages)-1]
	var rows []configRow

	for _, field := range configFields {
		if field.bits != nil {
			value := field.bits(final)
			if value == 0 {
				continue
			}
			var from []string
			for i := range stages {
				added := field.bits(&stages[i]) &^ previousBits(stages, i, field.bits) & value
				if added != 0 && !contains(from, labels[i]) {
					from = append(from, labels[i])
				}
			}
			rows = append(rows, configRow{field.name, strings.Join(bitNames(field.flags, value), " | "), strings.Join(from, ", ")})
			continue
		}

		value := field.read(final)
		if value == "" {
			continue
		}
		var from []string
		previous := ""
		for i := range stages {
			current := field.read(&stages[i])
			if current != previous && current != "" {
				from = append(from, labels[i])
			}
			previous = current
		}
		rows = append(rows, configRow{field.name, value, strings.Join(from, ", ")})
	}
	return rows
}

func previousBits(stages []core.SpellConfig, i int, read func(*core.SpellConfig) uint64) uint64 {
	if i == 0 {
		return 0
	}
	return read(&stages[i-1])
}

func contains(list []string, s string) bool {
	for _, item := range list {
		if item == s {
			return true
		}
	}
	return false
}

const configFootnote = "Assignments to the config after the call are not folded in."

func configMarkdown(result *configResult) string {
	var md strings.Builder
	fmt.Fprintf(&md, "`SpellConfig` of %s\n\n", title(result.spell))
	if result.pick != "" {
		fmt.Fprintf(&md, "`%s`\n\n", result.pick)
	}
	md.WriteString("---\n| field | value | from |\n|---|---|---|\n")
	for _, row := range result.rows {
		fmt.Fprintf(&md, "| **%s** | %s | %s |\n", row.field, codeCell(row.value), cell(row.from))
	}
	for _, option := range result.skipped {
		fmt.Fprintf(&md, "\nunevaluated: `%s` (%s)\n", option.source, option.err)
	}
	fmt.Fprintf(&md, "\n%s\n", configFootnote)
	return md.String()
}

func writeConfigText(out *strings.Builder, result *configResult) {
	fmt.Fprintf(out, "SpellConfig of %s\n", title(result.spell))
	if result.pick != "" {
		fmt.Fprintf(out, "%s\n", result.pick)
	}
	out.WriteString("\n")
	fieldWidth, valueWidth := 0, 0
	for _, row := range result.rows {
		fieldWidth = max(fieldWidth, len(row.field))
		valueWidth = max(valueWidth, len(row.value))
	}
	for _, row := range result.rows {
		fmt.Fprintf(out, "%-*s  %-*s  %s\n", fieldWidth, row.field, valueWidth, row.value, row.from)
	}
	for _, option := range result.skipped {
		fmt.Fprintf(out, "\nunevaluated: %s (%s)\n", option.source, option.err)
	}
	fmt.Fprintf(out, "\n%s\n", configFootnote)
}
