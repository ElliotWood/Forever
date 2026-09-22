package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"sort"
	"strconv"
	"sync"
)

// Where a spell id sits in a class file's ladders: the package, the field name the generator gave the
// family, and the call that reaches this rank of it.
type ladderRef struct {
	pkg   string
	field string
	call  string
}

func (r ladderRef) String() string {
	return fmt.Sprintf("%s spellData.%s.%s", r.pkg, r.field, r.call)
}

var (
	loadLadders sync.Once
	ladderIndex = map[int32][]ladderRef{}
)

// The ladder calls that reach this id, in package and field order. A class whose generated file is
// still a shared.SpellDataTable states no ladder and answers nothing, and so does a file that does
// not parse - another session editing a class package must not stop the printer.
func ladderRefs(id int32) []string {
	loadLadders.Do(scanLadders)

	refs := ladderIndex[id]
	out := make([]string, 0, len(refs))
	for _, ref := range refs {
		out = append(out, ref.String())
	}
	return out
}

// Every `spelldata.Ranked(...)` and `spelldata.Talent(...)` in the generated class files, read as
// source rather than through an import: tools/spelldata must not depend on a class package.
func scanLadders() {
	root, err := moduleRoot()
	if err != nil {
		return
	}

	files, err := filepath.Glob(filepath.Join(root, "sim", "*", "spell_data_auto_gen.go"))
	if err != nil {
		return
	}
	sort.Strings(files)

	for _, path := range files {
		file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
		if err != nil {
			continue
		}
		collectLadders(file)
	}
}

func collectLadders(file *ast.File) {
	pkg := file.Name.Name

	ast.Inspect(file, func(node ast.Node) bool {
		kv, ok := node.(*ast.KeyValueExpr)
		if !ok {
			return true
		}
		field, ok := kv.Key.(*ast.Ident)
		if !ok {
			return true
		}
		ctor, args := ladderCall(kv.Value)
		if ctor == "" {
			return true
		}

		for _, ref := range ladderCalls(pkg, field.Name, ctor, args) {
			ladderIndex[ref.id] = append(ladderIndex[ref.id], ref.ref)
		}
		return true
	})
}

// The constructor name and its integer arguments, for a value that is a call on the spelldata
// package. Anything else answers an empty name.
func ladderCall(value ast.Expr) (string, []int32) {
	call, ok := value.(*ast.CallExpr)
	if !ok {
		return "", nil
	}
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return "", nil
	}
	if pkg, ok := sel.X.(*ast.Ident); !ok || pkg.Name != "spelldata" {
		return "", nil
	}

	args := make([]int32, 0, len(call.Args))
	for _, arg := range call.Args {
		lit, ok := arg.(*ast.BasicLit)
		if !ok || lit.Kind != token.INT {
			return "", nil
		}
		value, err := strconv.ParseInt(lit.Value, 0, 32)
		if err != nil {
			return "", nil
		}
		args = append(args, int32(value))
	}
	return sel.Sel.Name, args
}

type ladderEntry struct {
	id  int32
	ref ladderRef
}

// How a class file reaches one rank, by the constructor the generator wrote. Ranked states one spell
// per rank, so the last id is what Highest() answers and an earlier one is what ByID() names; Talent
// states one spell whose ranks are a curve, so every rank is that id and Rank(n) is the only way in.
func ladderCalls(pkg, field, ctor string, args []int32) []ladderEntry {
	switch ctor {
	case "Talent":
		if len(args) != 2 {
			return nil
		}
		return []ladderEntry{{id: args[0], ref: ladderRef{pkg: pkg, field: field,
			call: fmt.Sprintf("Rank(n), n up to %d", args[1])}}}

	case "Ranked":
		var out []ladderEntry
		for i, id := range args {
			call := fmt.Sprintf("ByID(%d)", id)
			if i == len(args)-1 {
				call = "Highest()"
			}
			out = append(out, ladderEntry{id: id, ref: ladderRef{pkg: pkg, field: field, call: call}})
		}
		return out
	}
	return nil
}
