package database

// Rebuilds sim/core/dbcenums/forms_auto_gen.go from the SpellShapeshiftForm rows committed in
// assets/db_inputs/spell_store_inputs.json and asserts the committed file is what comes out. No
// client database and no build tag, the same as the store's own regeneration check.

import (
	"bytes"
	"os"
	"testing"
)

func TestFormsRegenerateFromTheCommittedInputs(t *testing.T) {
	inRepositoryRoot(t)

	inputs, err := readStoreInputs(spellStoreInputsPath)
	if err != nil {
		t.Fatalf("%v", err)
	}

	rendered, err := renderFormsFile(inputs.Forms)
	if err != nil {
		t.Fatalf("rendering the forms file: %v", err)
	}

	const formsPath = "sim/core/dbcenums/forms_auto_gen.go"
	committed, err := os.ReadFile(formsPath)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if bytes.Equal(committed, rendered) {
		return
	}

	line, want, got := firstDifference(committed, rendered)
	t.Errorf("%s is not what the committed inputs render, from line %d:\n  committed: %s\n  rendered:  %s\n"+
		"regenerate both with `go run ./tools/database/gen_spelldata`",
		formsPath, line, want, got)
}
