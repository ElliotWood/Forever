// Regenerates sim/<class>/spell_data_auto_gen.go from the client database.
//
// Its own binary, not a mode of gen_db: gen_db imports the sim, so a stale generated file would stop
// the generator that fixes it from compiling.
//
//	go run ./tools/database/gen_spelldata
package main

import (
	"flag"
	"log"

	"github.com/wowsims/forever/tools/database"
)

var dbPath = flag.String("dbPath", "./tools/database/wowsims.db", "Location of the wowsims.db file produced by tools/db2tool")

func main() {
	flag.Parse()
	database.DatabasePath = *dbPath

	helper, err := database.NewDBHelper()
	if err != nil {
		log.Fatalf("failed to open %s: %v", *dbPath, err)
	}
	defer helper.Close()

	if err := database.GenerateSpellDataFiles(helper); err != nil {
		log.Fatalf("failed to generate spell data tables: %v", err)
	}

	if err := database.GenerateBuffFiles(helper); err != nil {
		log.Fatalf("failed to generate buff files: %v", err)
	}

	if err := database.GenerateBuffsDebuffsTSFile(helper); err != nil {
		log.Fatalf("failed to generate the settings buff inputs: %v", err)
	}
}
