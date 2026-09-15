# Database Updates

This doc outlines the various methods for updating database data in WoWSims Classic.

## Introduction

WoWSims uses a number of JSON database files to store our data.
These can be found in [assets/database/](https://github.com/wowsims/classic/tree/master/assets/database):

-   [db.json](https://github.com/wowsims/classic/blob/master/assets/database/db.json) - Our primary database file. This is where we store the majority of our item, spell, and other data used in the sims.
-   [leftover_db.json](https://github.com/wowsims/classic/blob/master/assets/database/leftover_db.json) - This is a leftover database file where items filtered out for having invalid data (determined in `main.go`) are placed.

We also store scraped data from external sources like Wowhead, Atlasloot, etc. in order to build database data in the first place.
These can be found in [assets/db_inputs/](https://github.com/wowsims/classic/tree/master/assets/db_inputs):

-   [atlasloot_db.json](https://github.com/wowsims/classic/blob/master/assets/db_inputs/atlasloot_db.json) - Scraped data from the AtlasLoot addon's GitHub repo. This is primarily used for item source data and faction restrictions.
-   [wago_db2_items.csv](https://github.com/wowsims/classic/blob/master/assets/db_inputs/wago_db2_items.csv) - Scraped item data from Wago DB. This is currently used for additional item set and faction restriction information.
-   [wowhead_gearplannerdb.txt](https://github.com/wowsims/classic/blob/master/assets/db_inputs/wowhead_gearplannerdb.txt) - Scraped item data from Wowhead Gear Planner. This is one of our primary sources of item data along with Wowhead item tooltips.
-   [wowhead_item_tooltips.csv](https://github.com/wowsims/classic/blob/master/assets/db_inputs/wowhead_item_tooltips.csv) - Scraped item tooltips from Wowhead. We store the full tooltip along with the item IDs.
-   [wowhead_rune_tooltips.csv](https://github.com/wowsims/classic/blob/master/assets/db_inputs/wowhead_rune_tooltips.csv) - Scraped rune tooltips from Wowhead. We store the full tooltip along with the rune engraving spell IDs.
-   [wowhead_spell_tooltips.csv](https://github.com/wowsims/classic/blob/master/assets/db_inputs/wowhead_spell_tooltips.csv) - Scraped spell tooltips from Wowhead. We store the full tooltip along with the spell IDs.

The entry point for running database scripts is [tools/database/gen_db/main.go](https://github.com/wowsims/classic/blob/master/tools/database/gen_db/main.go).
This file is executed by running the `make items` command, or by running one of a set of commands listed in the comments of `gen_db/main.go`.

-   `go run ./tools/database/gen_db -outDir=assets -gen=atlasloot`
    -   Run the Atlasloot scraper. You do not need to delete existing data beforehand.
-   `go run ./tools/database/gen_db -outDir=assets -gen=wago-db2-items`
    -   Run the Wago DB2 scraper. You do not need to delete existing data beforehand.
-   `go run ./tools/database/gen_db -outDir=assets -gen=wowhead-gearplannerdb`
    -   Run the Wowhead Gear Planner scraper. You do not need to delete existing data beforehand.
-   `go run ./tools/database/gen_db -outDir=assets -gen=wowhead-items`
    -   Run the Wowhead Item Tooltips scraper. This only pulls in new entries, so if you want to regenerate existing data you will need to delete the existing entries. Often times you will want to pass IDs to restrict the range of items pulled in at once. This can be done using the `-id=`, `-minid=`, and `maxid=` commands.
    -   Example: `go run ./tools/database/gen_db -outDir=assets -gen=wowhead-items -minid=200000 -maxid=235000`
-   `go run ./tools/database/gen_db -outDir=assets -gen=wowhead-spells`
    -   Run the Wowhead Spell Tooltips scraper. This only pulls in new entries, so if you want to regenerate existing data you will need to delete the existing entries. Often times you will want to pass IDs to restrict the range of items pulled in at once. This can be done using the `-id=`, `-minid=`, and `maxid=` commands.
    -   Example: `go run ./tools/database/gen_db -outDir=assets -gen=wowhead-spells -minid=200000 -maxid=235000`
-   `python3 tools/scrape_runes.py assets/db_inputs/wowhead_rune_tooltips.csv`
    -   Run the Rune Tooltips scraper. This command is slightly different in that it actually visits Wowhead and parses data from the page using Selenium. You do not need to delete existing data beforehand.

After performing any of these commands, you should then use `make items` in order to rebuild the database with the updated input data.

## Overrides

In addition to our db inputs, we can also define overrides for both adding and removing data.
We have three different overrides files:

-   [tools/database/overrides.go]https://github.com/wowsims/classic/blob/master/tools/database/overrides.go has several different types of overrides for items, including allowlists and denylists.
-   [tools/database/enchant_overrides.go](https://github.com/wowsims/classic/blob/master/tools/database/enchant_overrides.go) has all of our enchant data and is where new enchant entries should be added
-   [tools/database/rune_overrides.go](https://github.com/wowsims/classic/blob/master/tools/database/rune_overrides.go) has overrides related to runes, including overrides that fix some rune information scraped from wowhead, and a way to block runes from the sim until they can be implemented.

These overrides should be used sparingly when possible, but are often a necessary part of filtering data in the sim's database.

## Updating the Item Database

This involves performing several commands to update the multiple sources of item data.

1. Run the Wowhead Gear Planner scraper
2. Run the Wowhead Item Tooltip scraper with the range of IDs you want to add/update, deleting any existing entries in that range if there are any.
3. Optionally run the Atlasloot scraper to grab any new source data
4. Finally run `make items` to regenerate the database.

You should now also run tests because they could have changed if existing items were changed, then you're ready to commit everything.

## Refreshing for Forever

Every scraper reads Wowhead's Classic Era database by default. Wowhead also runs a Forever
database at wowhead.com/forever, which will carry the beta client's items and spells once it
is datamined. Pass `-wowhead=forever` to point the scrapers at it:

1. `go run ./tools/database/gen_db -outDir=assets -wowhead=forever -gen=wowhead-gearplannerdb`
2. `go run ./tools/database/gen_db -outDir=assets -wowhead=forever -gen=wowhead-items -minid=<first new id> -maxid=<last new id>`
   (delete existing rows in that range from `wowhead_item_tooltips.csv` first if they should be re-read)
3. `go run ./tools/database/gen_db -outDir=assets -wowhead=forever -gen=wowhead-spells -minid=... -maxid=...`
4. `go run ./tools/database/gen_db -outDir=assets -wowhead=forever -gen=atlasloot`, if AtlasLoot has Forever data by then; otherwise the Classic source data stands and the new dungeons' items arrive without a source until it does.
5. `make items`, then review `tools/database/launch_content.go`: it filters the Classic database to what Forever has at launch by phase, zone and item level, and new Forever items (the nine dungeons, the four zones) will carry no Classic phase, so check they come through rather than being filtered out.

The Wago DB2 export is pinned to a Classic build id and is not covered by the flag.
