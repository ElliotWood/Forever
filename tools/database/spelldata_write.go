package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

// Where the generator says what it did. -check silences it, so that a clean check prints nothing at
// all and a stale one prints the paths alone.
var progress io.Writer = os.Stderr

// Writes the rendered files, but only once the sim compiles against them.
//
// The generated files are the sim's data: one that does not compile takes the whole repository down
// with it, gen_db included, and gen_db is what rebuilds the database the generator reads. So the
// files go to a staging directory first and are type-checked there through go build's overlay,
// which lets the compiler read the staged bytes in the place of the committed ones without any of
// them being in the tree. A failure leaves the tree exactly as it was and prints what the compiler
// said.
func writeSpellDataFiles(files map[string][]byte) error {
	staging, err := os.MkdirTemp("", "spelldata-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(staging)

	overlay := map[string]string{}
	for path, out := range files {
		absolute, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		staged := filepath.Join(staging, strings.ReplaceAll(path, string(filepath.Separator), "_"))
		if err := os.WriteFile(staged, out, 0644); err != nil {
			return err
		}
		overlay[absolute] = staged
	}

	if err := buildStaged(staging, overlay, packagesOf(files)); err != nil {
		return err
	}

	for path, out := range files {
		if err := os.WriteFile(path, out, 0644); err != nil {
			return err
		}
	}
	return nil
}

// The committed files that are not what the generator writes today. Named rather than rewritten:
// the caller is a check, and a check that edits the tree is not one.
func checkSpellDataFiles(files map[string][]byte) ([]string, error) {
	var stale []string
	for path, out := range files {
		committed, err := os.ReadFile(path)
		if os.IsNotExist(err) {
			stale = append(stale, path)
			continue
		}
		if err != nil {
			return nil, err
		}
		if !bytes.Equal(committed, out) {
			stale = append(stale, path)
		}
	}
	sort.Strings(stale)
	return stale, nil
}

func buildStaged(staging string, overlay map[string]string, packages []string) error {
	document, err := json.Marshal(struct{ Replace map[string]string }{overlay})
	if err != nil {
		return err
	}
	path := filepath.Join(staging, "overlay.json")
	if err := os.WriteFile(path, document, 0644); err != nil {
		return err
	}

	root, err := os.Getwd()
	if err != nil {
		return err
	}

	cmd := exec.Command(goTool(), append([]string{"build", "-overlay", path}, packages...)...)
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("the sim does not compile against the generated files, so none was written:\n%s", out)
	}
	return nil
}

// The packages to type-check: the ones the rendered files belong to, which is the store, the shared
// enums and one per class.
func packagesOf(files map[string][]byte) []string {
	seen := map[string]bool{}
	for path := range files {
		seen["./"+filepath.ToSlash(filepath.Dir(path))+"/"] = true
	}

	packages := make([]string, 0, len(seen))
	for pkg := range seen {
		packages = append(packages, pkg)
	}
	sort.Strings(packages)
	return packages
}

// The go binary of the toolchain this generator was built with, so the check reads the same
// compiler the caller does.
func goTool() string {
	if root := runtime.GOROOT(); root != "" {
		path := filepath.Join(root, "bin", "go")
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	if path, err := exec.LookPath("go"); err == nil {
		return path
	}
	return "go"
}

// Regenerates every spell data file, once the sim compiles against all of them, and writes the
// client rows the store was built from beside them: they are what lets the store be regenerated and
// checked without the client database, so they are written from the same pass that wrote it and
// only once that pass has landed.
func GenerateSpellDataFiles(helper *DBHelper) error {
	files, inputs, err := renderSpellDataFiles(helper)
	if err != nil {
		return err
	}
	if err := writeSpellDataFiles(files); err != nil {
		return err
	}
	return writeStoreInputs(inputs)
}

// The committed files that no longer match what the generator writes, for a caller that wants to
// know rather than to regenerate. Silent on the way there: its output is the list it returns.
func CheckSpellDataFiles(helper *DBHelper) ([]string, error) {
	progress = io.Discard
	defer func() { progress = os.Stderr }()

	files, _, err := renderSpellDataFiles(helper)
	if err != nil {
		return nil, err
	}
	return checkSpellDataFiles(files)
}
