package sim

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var talentFieldRegex = regexp.MustCompile(`\.Talents\.([A-Za-z]+)`)

// A talent the sim does not implement shows "Not simulated yet - points spent here do not
// affect results." in the picker, but only if its tree entry says notSimulated. Three
// talents were missing the flag, so a player spending points in them was told nothing and
// got nothing.
//
// Nothing reads a talent it has no field for, so the useful direction is the other one:
// every talent in the trees should either be read somewhere in its class package or be
// marked. This finds the ones that are neither.
func TestUnsimulatedTalentsAreMarked(t *testing.T) {
	for _, class := range talentClasses {
		read := talentFieldsReadBy(t, filepath.Join("..", "sim", class.name))

		for _, tree := range loadTrees(t, class.name) {
			for i := range tree.Talents {
				talent := &tree.Talents[i]
				if talent.FieldName == "" || talent.NotSimulated {
					continue
				}
				// The Go writes the proto field in PascalCase.
				if read[strings.ToUpper(talent.FieldName[:1])+talent.FieldName[1:]] {
					continue
				}
				t.Errorf("%s/%s: %q is not read anywhere in sim/%s and is not marked notSimulated, so the picker tells a player nothing while their points do nothing",
					class.name, tree.Name, talent.Name, class.name)
			}
		}
	}
}

// Every Talents.<Field> named anywhere under dir.
func talentFieldsReadBy(t *testing.T, dir string) map[string]bool {
	t.Helper()

	read := map[string]bool{}
	err := filepath.WalkDir(dir, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		source, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		for _, match := range talentFieldRegex.FindAllStringSubmatch(string(source), -1) {
			read[match[1]] = true
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return read
}
