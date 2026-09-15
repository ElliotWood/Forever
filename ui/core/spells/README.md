# Where every ability's numbers came from

The Go sim carries an integer for each ability and nothing else. Every name, icon and hover
tooltip the site shows is resolved in the browser from Wowhead's Classic data for that
integer, so an ability Forever changed is described by its Classic ancestor unless something
says otherwise.

That is the same trap the talent trees fell into. Forty two talents wore another talent's
tooltip, and Improved Revenge kept Classic's stun long after Forever turned it into damage,
because the two share a name. The fix there was to stop letting the Classic database
describe a Forever talent at all. These files do the same for abilities.

One file per class, plus `core.json` for the raid buffs, debuffs and racials, `common.json`
for item and enchant effects and `encounters.json` for boss mechanics. Each is an object
keyed by the spell id the sim registers.

```json
{
	"13953": {
		"ability": "Holy Strike",
		"file": "sim/paladin/holy_strike.go",
		"source": "forever",
		"foreverId": 17143,
		"tooltip": "20 Mana\nMelee Range\nInstant\n12 sec cooldown\nDeals 40% weapon damage plus 36 to 46 Holy damage.",
		"note": "Forever's own id is not in the item database, so Classic's unused 13953 lends the name and the icon."
	}
}
```

## Fields

- `ability` — what Forever calls it.
- `file` — the Go file that registers it.
- `source` — one of:
  - `classic`, Forever did not change the ability, so Wowhead's tooltip describes it and the
    site keeps showing it.
  - `forever`, Forever changed it and a published Forever source gave the numbers. The site
    shows `tooltip` instead.
  - `assumed`, Forever changed it or it is new, and at least one number is a guess. The site
    shows `tooltip` and says the number is unconfirmed.
  - `unreviewed`, not classified yet.
- `foreverId` — Forever's own spell id, when it is known and differs from the one the sim
  carries.
- `tooltip` — required for `forever` and `assumed`, and rejected for `classic`. Written from
  the implementation's own numbers so the two cannot drift.
- `note` — one sentence on why the id or the number is what it is.
- `assumptions` — required for `assumed`, one short line per guessed number, phrased so a
  beta tester can check it. These feed the beta checklist.

`classic` is a claim that the ability is untouched, not a shrug. An ability nobody has
checked is `assumed` or `unreviewed`.

`sim/spell_sources_test.go` holds four rules: every id the sim registers appears here, every
entry is well formed, the count of `unreviewed` entries only ever falls, and every
registration site resolves to an id.

That last rule is the one that decides whether the rest means anything. Only a minority of
the sim's abilities write their id as a literal in the `ActionID`; the ranked ones build it
from a table, either a local one indexed by rank, a package level one indexed by level, or
a rank struct passed in. A walk that reads literals alone sees 314 ids and misses every
rank of Shred, Rip, Moonfire and Tiger's Fury, which is most of what Forever changed about
the druid. Following the tables finds 991. Nineteen sites still build their id somewhere
the walk cannot follow; they are listed by the test and that count only falls too.
