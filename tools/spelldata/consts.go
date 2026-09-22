package main

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/parser"
	"go/token"
	"go/types"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// Where the constants a class file passes to an accessor or an option are declared.
var constantSources = []struct {
	path  string
	files []string
}{
	{"github.com/wowsims/forever/sim/core", []string{"sim/core/flags.go", "sim/core/constants.go"}},
	{"github.com/wowsims/forever/sim/core/dbcenums", []string{"sim/core/dbcenums/*.go"}},
	{"github.com/wowsims/forever/sim/core/spelldata", []string{"sim/core/spelldata/parse_effects_table.go"}},
}

var (
	loadConstants sync.Once
	constFset     = token.NewFileSet()
	constPackages = map[string]*types.Package{}
	evalPackage   *types.Package
	evalPos       token.Pos
	coreNames     = map[string][]namedConstant{}
)

type namedConstant struct {
	name  string
	value uint64
}

// Only the packages above are imported, and only by each other: type-checking proto and stats from
// source takes seconds, and the constants do not use them. The errors that leaves behind are ignored;
// the constants still resolve.
type checkedImports map[string]*types.Package

func (c checkedImports) Import(path string) (*types.Package, error) {
	if pkg, ok := c[path]; ok {
		return pkg, nil
	}
	return nil, fmt.Errorf("%s is not imported", path)
}

func scanConstants() {
	root, err := moduleRoot()
	if err != nil {
		return
	}
	conf := types.Config{Importer: checkedImports(constPackages), Error: func(error) {}}

	for _, source := range constantSources {
		var files []*ast.File
		for _, pattern := range source.files {
			paths, _ := filepath.Glob(filepath.Join(root, pattern))
			for _, path := range paths {
				if strings.HasSuffix(path, "_test.go") {
					continue
				}
				file, err := parser.ParseFile(constFset, path, nil, parser.SkipObjectResolution)
				if err != nil {
					return
				}
				files = append(files, file)
			}
		}
		constPackages[source.path], _ = conf.Check(source.path, constFset, files, nil)
	}

	var imports strings.Builder
	for _, source := range constantSources {
		fmt.Fprintf(&imports, "import %q\n", source.path)
	}
	file, err := parser.ParseFile(constFset, "eval.go", "package eval\n"+imports.String(), 0)
	if err != nil {
		return
	}
	evalPackage, _ = conf.Check("eval", constFset, []*ast.File{file}, nil)
	evalPos = file.Name.Pos()

	scope := constPackages["github.com/wowsims/forever/sim/core"].Scope()
	for _, name := range scope.Names() {
		c, ok := scope.Lookup(name).(*types.Const)
		if !ok {
			continue
		}
		named, ok := c.Type().(*types.Named)
		if !ok {
			continue
		}
		value, exact := constant.Uint64Val(constant.ToInt(c.Val()))
		if !exact {
			continue
		}
		typeName := named.Obj().Name()
		coreNames[typeName] = append(coreNames[typeName], namedConstant{name: name, value: value})
	}
	for _, names := range coreNames {
		sort.SliceStable(names, func(i, j int) bool { return names[i].value < names[j].value })
	}
}

func coreConstants(typeName string) []namedConstant {
	loadConstants.Do(scanConstants)
	return coreNames[typeName]
}

// A variable or a call is refused: this reads, it does not run the package around it.
func evalConst(expr ast.Expr) (constant.Value, error) {
	loadConstants.Do(scanConstants)
	if evalPackage == nil {
		return nil, fmt.Errorf("sim/core's constants did not load")
	}
	tv, err := types.Eval(constFset, evalPackage, evalPos, nodeText(expr))
	if err != nil || tv.Value == nil {
		return nil, fmt.Errorf("%s is not a literal or a constant of sim/core, dbcenums or spelldata", nodeText(expr))
	}
	return tv.Value, nil
}

func evalInt(expr ast.Expr) (int64, error) {
	value, err := evalConst(expr)
	if err != nil {
		return 0, err
	}
	n, exact := constant.Int64Val(constant.ToInt(value))
	if !exact {
		return 0, fmt.Errorf("%s is not an integer", nodeText(expr))
	}
	return n, nil
}
