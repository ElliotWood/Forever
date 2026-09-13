# WoW Forever Sim — Claude Code handoff spec

**Goal:** fork `wowsims/classic` into a WoW: Forever simulator. Phase 1 = Warlock end-to-end (engine rule changes + Forever talent trees + tests + UI), because it has the contested answer. Phases 2–4 = Druid, Shaman, Priest.

**Repo:** fork `https://github.com/wowsims/classic` to a **private** repo under `github.com/elliotwood` (personal project → personal GitHub), branch `forever`. Keep upstream as a remote; Classic Era is the baseline Forever launches on, so merging upstream fixes stays valuable.

**Status date:** 14 Sep 2026. Forever beta client opens **Thu 17 Sep** — all talent data below is from BlizzCon demo tooltip captures (rank 1 only) and must be re-verified against the beta client / Wowhead Forever DB (`wowhead.com/forever`) as soon as it's datamined. Design the data layer so that swap is a data change, not a code change.

---

## 0. Toolchain (native machine — none of the sandbox workarounds apply)

- Go ≥ 1.23 (`go.mod` says 1.23.0, toolchain 1.23.4). Standard module proxy works on your network; do **not** carry over the `replace` hacks I used.
- `protoc` + `protoc-gen-go` (`go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`), Node/npm (`npm install` pulls `protoc-gen-ts`).
- `make proto` regenerates `sim/core/proto/*.pb.go` and `ui/core/proto/*.ts` from `proto/*.proto`.
- `make test` = `go test --tags=with_db ./sim/...` — the per-spec tests compare against committed `*.results` files (e.g. `sim/warlock/dps/TestWarlockDSRuin.results`). Any engine change that alters numbers requires regenerating them (`go test ... -update` — check the test harness flag name in `sim/core/test_suite.go`). Regenerate deliberately, per change, and review the diff.
- Windows: WSL2 is the smooth path (the makefile assumes a POSIX shell, `ulimit`, `find`). Windows-native works for `go build/test` but not `make`.

## 1. Repo map (what matters)

| Concern | Location |
|---|---|
| Engine core | `sim/core/` — `dot.go`, `spell_outcome.go`, `spell_result.go`, `aura.go`, `debuffs.go`, `buffs.go`, `apl_*.go` |
| Class impl | `sim/warlock/*.go` (one file per spell; `talents.go` applies talents; `pet.go`, `imp.go`, `succubus.go`, …) |
| Talent field defs | `proto/warlock.proto` → `message WarlockTalents` (field numbers are the wire format; the UI talent JSON references `fieldName` in camelCase) |
| UI talent tree layout | `ui/core/talents/trees/warlock.json` — list of 3 trees, each `{name, backgroundUrl, talents:[{fieldName, location:{rowIdx,colIdx}, spellIds:[per rank], maxPoints}]}`. Prereq arrows are also in here (check for a `prereqLocation` key in other trees). |
| Spec tests | `sim/warlock/dps/dps_warlock_test.go` + `.results` |
| Presets / APLs | `ui/warlock/presets.ts`, `ui/warlock/apls/*.json`, `sim/warlock/dps/dps_warlock.go` |
| Talent scraper (for reference) | `tools/scrape_talents_proto.py` scrapes Wowhead classic talent-calc via Selenium — adapt to `wowhead.com/forever/talent-calc/<class>` once the beta data is live |

## 2. Engine rule changes for Forever

### 2.1 Periodic damage can crit (confirmed by tooltip language)
Evidence: Nature's Grace "non-periodic spell criticals"; Primal Fury "non-periodic critical strikes from Cat Form abilities"; Pandemic (crit *damage* bonus for DoTs) exists at all. Applies to spell DoTs **and** bleeds.

The engine already supports it: `dot.Snapshot()` sets `dot.SnapshotCritChance` for every dot (`spell_result.go:392/394/547`), and `OutcomeSnapshotCrit` / `OutcomeMagicHitAndSnapshotCrit` / `OutcomeTickPhysicalCrit` exist in `spell_outcome.go`. Classic dots just call `dot.OutcomeTick` (no crit roll).

**Implementation:** add a ruleset switch rather than editing 100 spell files:
```go
// sim/core/ruleset.go (new)
type Ruleset int
const ( RulesetClassic Ruleset = iota; RulesetForever )
var ActiveRuleset = RulesetForever  // wire to a proto option on SimOptions later
```
Then in `spell_outcome.go`, make `(dot *Dot) OutcomeTick` delegate to `OutcomeSnapshotCrit` when `ActiveRuleset == RulesetForever` (and `OutcomeTickMagicHit` → `OutcomeMagicHitAndSnapshotCrit`). Add a `SpellFlagNoDotCrit` escape hatch for anything that shouldn't crit. Check `CritMultiplier` for dots: Classic spell crit bonus is 50%; Forever crit-damage talents (Pandemic, Ruin, Shadowform, Elemental Fury, Vengeance, Predatory Instincts) must apply to the right subset — see per-class notes.

Also check `spell_outcome.go:972` (expected-damage path used by APL values) already accounts for snapshot crit — it does.

### 2.2 Curse / Bane split (Warlock)
Curse of Agony → **Bane of Agony**; Curse of Doom → Bane of Doom; new **Bane of Havoc** (Destro talent). Tooltip: "only one Bane per Warlock can be active on any one target" — Banes and Curses are separate exclusivity groups, so a lock runs Bane of Agony + Curse of Elements/Shadow/Recklessness simultaneously. In `sim/warlock/curses.go` split the exclusive-effect category into two. Amplify Curse now boosts "Curse of Weakness or Bane of Agony by 50%, Curse of Exhaustion 20%".

### 2.3 Improved Shadow Bolt: no charge limit
Forever text: "Shadow Bolt critical strikes increase Shadow damage taken by the target from your attacks by 4% for 12 sec" (per rank). Personal debuff (not raid-shared), no 4-charge consumption. Rewrite `applyImprovedShadowBolt` accordingly.

### 2.4 Other engine-touching changes seen in tooltips
- **Nightfall** procs from Corruption, Drain Soul **and Drain Life** (2%/rank).
- **Demonic Sacrifice** roles swapped: Imp → +15% **Shadow**, Sayaad → +15% **Fire**, Voidwalker → mana, Felhunter → health. 2-hr duration. Row 3.
- **Demonic Pact** (Demo capstone): sacrifice effect not cancelled by summoning a *different* demon; resummoning the sacrificed one cancels it. Whether a *second* sacrifice replaces or stacks is unknown — tooltip is singular ("an effect") → implement as **replace**, flag for beta test.
- **Furor (Druid)**: cat shift regains 20% of last energy + 2/sec out of form (cap 20). Powershifting rework — energy model change.
- **Water Shield (Shaman, Resto row 3)**: 2% max mana per globe on hit or healing crit, 3 globes.
- **Stormstrike**: +20% Nature damage *dealt to the target* for 12 s (not the Classic charge-based debuff).
- **Nature's Grace (Druid)**: 10% spell haste + GCD reduction for 3 s after non-periodic crit (was −0.5 s next cast).
- **Devouring Plague / Desperate Prayer etc.**: priest racial spells appear class-wide; **Shadow Word: Death** exists (Early Demise talent references it).
- **Racials**: revamped to 2 active + 2 passive per race; not yet captured in full. New combos: Gnome Priest, Human Hunter, Dwarf Shaman, Orc Mage, Troll Warlock, Undead Paladin, Skyborne (Horde: Shaman; Alliance: Mage; both: Warrior/Hunter/Rogue/Druid).

## 3. Warlock Forever talent trees (Phase 1 data)

Source: BlizzCon demo footage transcribed at classicwowforever.com/talents/warlock (rank-1 text). Row gating = 5 points per row, 51 total. **Prerequisite arrows were not captured** — the tree screenshot is at `classicwowforever.com/talents/warlock.png`; read arrows off it or wait for beta. Rank scaling is assumed linear from rank 1 (flag every place this matters).

Proto: add new fields to `WarlockTalents` **appending new field numbers** (don't renumber existing ones); mark removed Classic talents `reserved`. Removed from tree: Dark Pact, Improved Curse of Weakness, Improved Drain Soul, Improved Drain Life, Grim Reach, Improved Drain Mana, Improved Curse of Exhaustion, Improved Healthstone, Fel Intellect, Fel Stamina, Improved Subjugate Demon, Improved Firestone, Improved Spellstone, Improved Firebolt, Improved Lash of Pain, Devastation, Improved Searing Pain, Improved Immolate, Emberstorm.

**Affliction** (row · ranks · rank-1 effect)
1 · Improved Life Tap 2 · Life Tap mana +10%
1 · Suppression 5 · +1% hit all spells/attacks, −4% threat
1 · Improved Corruption 5 · cast −0.4 s, damage +2%
2 · Malediction 5 · all periodic damage +1%
2 · Soul Harvesting 2 · Drain Soul kill → 10 s of 50% regen while casting +50% regen
2 · Improved Drains 3 · Drain Life/Soul +2% per other Affliction effect (max 6%), ×3 below 20% hp, +3 yd Drain Life
3 · Improved Bane of Agony 2 · +5%
3 · Fel Concentration 3 · 23% pushback resist on drains
3 · Amplify Curse 1 · 3-min CD, next CoW/BoA +50% or CoEx +20%
3 · Pandemic 3 · +33% crit damage bonus on Corruption, BoA, BoD, Drain Soul, Drain Life, Siphon Life, Drain Hope
4 · Malevolence 5 · +1% Shadow crit
4 · Nightfall 2 · 2% per Corruption/Drain Soul/Drain Life damage → instant Shadow Bolt
4 · Curse of Exhaustion 1
5 · Siphon Life 1 · (demo: 15 hp/3 s; use Classic r4 values at 60)
5 · Soul Siphon 3 · drain tick rate +17%, Drain Life healing −10%
6 · Shadow Mastery 5 · +1% Shadow damage
7 · Drain Hope 1 · channel 6 s, 52/s at demo level; other Shadow DoTs on target +10% while active

**Demonology**
1 · Improved Health Funnel 2 · +20% transfer, −15% cost, −50% threat, usable at any pet hp
1 · Improved Imp 3 · Firebolt +10%, Fire Shield +10%
1 · Demonic Embrace 5 · +3% Stamina
1 · Unholy Power 5 · pet damage +2%
2 · Demonic Aegis 2 · Demon Skin/Armor +15%
2 · Improved Voidwalker 3 · +10%
2 · Fel Vitality 3 · pet hp/mana +5%, **your max mana +5%**
2 · Demonic Energies 2 · heal pet 8% of spell damage; Life Tap gives pet 50% of mana gained
3 · Improved Sayaad 3 · Succubus **and Incubus** effects +10%
3 · Demonic Sacrifice 1 · see §2.4
3 · Master Summoner 2 · summon −2 s, −20% mana
4 · Decimation 2 · Soul Fire CD −45%; SB/Searing Pain on target <35% → +3% damage and next 10 s Soul Fire −20% cast, no shard
4 · Fel Domination 1
4 · Demonic Brand 3 · Searing Pain −17% threat, brands target 10 s; pet's next 2 attacks high threat + 39–42 dmg
5 · Improved Felhunter 3
5 · Soul Link 1 · 30% damage to pet, **+3% damage both**
5 · Demonic Knowledge 3 · spell damage/healing +33% of your level per rank (**per-rank vs flat is the single biggest unknown — it decides Pact vs DS/Ruin**)
6 · Master Demonologist 5 · with active demon: Imp +2% Fire, VW −2% Physical taken, Sayaad **+2% Shadow**, Felhunter −2% Magic taken
7 · Demonic Pact 1 · see §2.4

**Destruction**
1 · Destructive Reach 2 · +10% range
1 · Improved Shadow Bolt 5 · see §2.3
1 · Bane 5 · SB/Immolate/**Incinerate** −0.1 s, Soul Fire −0.4 s
2 · Molten Skin 5 · −2% damage taken
2 · Cataclysm 3 · Destro mana −3%
2 · Aftermath 5 · Immolate initial +10%, Conflagrate 20% daze
3 · **Ruin 5** · +20% crit damage bonus for Destruction spells (= Classic +100% at 5/5, but 15-point reach)
3 · Shadowburn 1 · (demo 102–111; use Classic rank at 60)
4 · Intensity 3
4 · Agonizing Flames 3 · Searing Pain crit +3%, **all Destruction damage +3%**
4 · Conflagrate 1
5 · Pyroclasm 2 · 13% stun on Soul Fire / RoF / Hellfire
5 · Bane of Havoc 1 · 5 min, 15% of damage to others copied to banned target
5 · Fire and Brimstone 3 · Conflagrate crit +8%
6 · Shadow and Flame 5 · Conflag hit → Shadow +2% 20 s; Shadowburn hit → Fire +2% 20 s; Conflag 20% not to consume Immolate; Shadowburn 20% shard refund
7 · Incinerate 1 · 2.5 s cast, 125–140 at demo, +25% if Immolate up

## 4. Validation targets (from my Python expected-value sim, 300 s, 450 SP / 8% crit / 6.2k mana)

Use these as sanity anchors, not truth. Expect the Go sim to differ but the **ordering and the mana finding should reproduce**:

| Build | DPS | Life Taps | Notes |
|---|---|---|---|
| 2/31/18 Demonic Pact (Imp sac + Succubus out, ISB/Bane/Ruin/AgFlames) | 542 | 46 | wins every sensitivity (DK flat: 526; pet 30 dps: 514; 90 s: 585) |
| 32/0/19 deep Affliction, pet out | 528 | 44 | local search dropped Drain Hope & Nightfall |
| 21/11/19 DS/Ruin + Pandemic 3 | 508 | 43 | |

Key finding to reproduce: ~22% of fight time is Life Tap GCDs → Improved Life Tap 2/2 and Fel Vitality 3/3 are in every optimal build; Improved Corruption beyond 2/5 is not worth the points. If the Go sim disagrees, investigate mana model first (Spirit regen, Demonic Energies, Soul Harvesting are unmodelled in mine).

## 5. Phased plan with acceptance criteria

**Phase 0 — scaffold (½ day)**
- Fork, branch `forever`, `make proto && make test` green on unmodified tree. Commit the baseline `.results`.
- Add `sim/core/ruleset.go` + `SimOptions` proto field `ruleset` (default Forever in this fork). CI: run tests under both rulesets.

**Phase 1 — engine rules (1 day)**
- §2.1 dot crits behind ruleset. Acceptance: Corruption `CritTicks` metric > 0 in a Forever run; Classic ruleset `.results` unchanged byte-for-byte.
- §2.2 Curse/Bane split; §2.3 ISB rewrite. Acceptance: BoA + CoE both active in a sim log.

**Phase 2 — Warlock talents (2–3 days)**
- Proto fields, `talents.go` implementations, `trees/warlock.json` with new layout (rowIdx/colIdx from the tree screenshot; spellIds are unknown pre-beta — use placeholder negative ids and a TODO; the UI only needs them for tooltips/icons).
- New spells: Incinerate, Drain Hope, Bane of Havoc (Havoc can be stubbed as no-op single-target), Decimation aura, Demonic Pact logic in `applyDemonicSacrifice`.
- Tests: `TestWarlockForeverPact`, `TestWarlockForeverDSRuin`, `TestWarlockForeverAffliction` with the three builds in §4; APLs in `ui/warlock/apls/forever_*.json`.
- Acceptance: DPS ordering Pact > Aff-pet > DS/Ruin reproduces (or a documented reason it doesn't).

**Phase 3 — UI (1 day)**: talent picker renders new trees; a "Ruleset: Forever" toggle; presets. Deploy via the existing `make dist` → GitHub Pages on the private repo (or Vercel, which you already have connected).

**Phase 4 — Druid / Shaman / Priest** — trees are transcribed at `classicwowforever.com/talents/{druid,shaman,priest}`; same recipe. Feral first (energy/Furor rework + Berserk/Mangle/Rend and Tear), it's where rotation policy matters most.

**Beta day (17 Sep)**: adapt `tools/scrape_talents_proto.py` to `wowhead.com/forever/talent-calc/<class>`, diff against the transcribed data, fix ranks/prereqs. Everything in §3 marked "demo" gets replaced.

## 6. Open questions to resolve in beta (ranked by DPS impact)
1. Demonic Knowledge per-rank or flat (decides the top warlock build).
2. Prerequisite arrows for all four classes.
3. Rank scaling of every 3–5 rank talent (linear assumed).
4. Does a second Demonic Sacrifice replace or stack under Pact?
5. Holy Precision 6% at rank 1 — real (6/12/18) or a rank-3 capture?
6. Do bleeds use melee crit (Predatory Instincts "melee abilities") — affects Feral crit-damage attribution.
7. Soul Shard economy with Shadowburn's 20% refund (Shadow and Flame).

## 7. Kickoff prompt for Claude Code

> Read FOREVER_SIM_HANDOFF.md in the repo root. We are forking wowsims/classic into a WoW: Forever simulator on branch `forever`. Start with Phase 0: confirm `make proto` and `make test` pass on the unmodified tree, then add the `Ruleset` switch (sim/core/ruleset.go + SimOptions proto field) and implement §2.1 (periodic crits behind the Forever ruleset) with the escape-hatch flag. Show me the Corruption CritTicks metric from a Forever-ruleset warlock test before moving on. Do not regenerate `.results` files without telling me which tests changed and why. Use private GitHub only.

---
*Companion artefacts from the analysis session: `forever_talent_optimiser.py` (exact DP), `forever_optimiser_dotcrit.py` (direct/periodic split model), `forever_warlock_sim_bnb.py` (rotation sim + branch-and-bound), and their results files.*
