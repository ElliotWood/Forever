package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/wowsims/forever/sim/core/dbcenums"
)

// The client's names for an effect type and an aura are Go constants, and a Go program cannot ask a
// constant for its name, so they are parsed out of the source the way tools/database's
// gen_spell_data_enums.go does. A name the parse cannot reach reads as its number, which is what the
// tool prints when it is run from outside the repository.
const dbcEnumsDir = "sim/core/dbcenums"

var (
	loadNames  sync.Once
	effectName = map[int32]string{}
	auraNames  = map[int32]string{}
)

func effectTypeName(t dbcenums.SpellEffectType) string {
	loadNames.Do(parseEnumNames)
	if name, ok := effectName[int32(t)]; ok {
		return name
	}
	return fmt.Sprintf("E_%d", t)
}

func auraName(a dbcenums.EffectAuraType) string {
	loadNames.Do(parseEnumNames)
	if name, ok := auraNames[int32(a)]; ok {
		return name
	}
	return fmt.Sprintf("A_%d", a)
}

func parseEnumNames() {
	root, err := moduleRoot()
	if err != nil {
		return
	}
	dir := filepath.Join(root, dbcEnumsDir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	wanted := map[string]map[int32]string{
		"SpellEffectType": effectName,
		"EffectAuraType":  auraNames,
	}
	fset := token.NewFileSet()
	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".go") || strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		file, err := parser.ParseFile(fset, filepath.Join(dir, entry.Name()), nil, 0)
		if err != nil {
			continue
		}
		collectEnumNames(file, wanted)
	}
}

func collectEnumNames(file *ast.File, wanted map[string]map[int32]string) {
	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.CONST {
			continue
		}
		for _, spec := range gen.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok || vs.Type == nil || len(vs.Names) != 1 || len(vs.Values) != 1 {
				continue
			}
			ident, ok := vs.Type.(*ast.Ident)
			if !ok {
				continue
			}
			into, ok := wanted[ident.Name]
			if !ok {
				continue
			}
			lit, ok := vs.Values[0].(*ast.BasicLit)
			if !ok || lit.Kind != token.INT {
				continue
			}
			value, err := strconv.ParseInt(lit.Value, 0, 32)
			if err != nil {
				continue
			}
			into[int32(value)] = vs.Names[0].Name
		}
	}
}

// The repository the tool is reading, found by walking up from the working directory. The extension
// spawns the tool with the folder holding go.mod as its working directory, and `go run` is normally
// typed there too.
func moduleRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("no go.mod above %s", dir)
		}
		dir = parent
	}
}
