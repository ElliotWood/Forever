// Regenerates sim/<class>/spell_data_auto_gen.go from the client database.
//
// Its own binary, not a mode of gen_db: gen_db imports the sim, so a stale generated file would stop
// the generator that fixes it from compiling.
//
//	go run ./tools/database/gen_spelldata
//	go run ./tools/database/gen_spelldata -check
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/wowsims/forever/tools/database"
)

var dbPath = flag.String("dbPath", "./tools/database/wowsims.db", "Location of the wowsims.db file produced by tools/db2tool")
var check = flag.Bool("check", false, "Name the generated files that are not what this generator writes, and write nothing")

func main() {
	flag.Parse()
	database.DatabasePath = *dbPath

	helper, err := database.NewDBHelper()
	if err != nil {
		log.Fatalf("failed to open %s: %v", *dbPath, err)
	}
	defer helper.Close()

	if *check {
		stale, err := database.CheckSpellDataFiles(helper)
		if err != nil {
			log.Fatalf("failed to generate spell data tables: %v", err)
		}
		for _, path := range stale {
			fmt.Fprintln(os.Stderr, path)
		}
		if len(stale) > 0 {
			os.Exit(1)
		}
		return
	}

	if err := database.GenerateSpellDataFiles(helper); err != nil {
		log.Fatalf("failed to generate spell data tables: %v", err)
	}
}
