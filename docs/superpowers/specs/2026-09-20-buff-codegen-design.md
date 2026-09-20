# Generate buffs/debuffs from the Forever client DB

Visual version of this plan (pipeline, tag contract, deltas, phases): https://claude.ai/artifact/3zAUVTxXN3xNWEKEU7upHz

## Context

`sim/core/buffs.go` (2102 lines) and `sim/core/debuffs.go` (904 lines) are hand-written TBC code. Every value and spell ID is hardcoded (Battle Shout = 306 AP / spell 2048; spell 2048 does not exist in the Forever DB). Nothing in them reads the generated `sim/<class>/spell_data_auto_gen.go` tables. Adding one buff today touches 4-6 places by hand: proto field, TS input def, TS registry array, Go apply block, per-spec presets, and a per-name `case` in `ui/sim/proto/action_id/index.ts` for metrics naming. Only the warrior's Battle Shout has a working class→raid bridge; every other class bridge is a `panic("To be implemented")` stub.

Goal: one checked-in **manifest** (one row per buff/debuff) + the existing `tools/database` parsers derive **proto + Go + TS** from `tools/database/wowsims.db`. Forever's actual values, ranks, durations, icons and talent availability replace TBC numbers. Class-owned buffs get one integration contract so the player-cast and external versions separate cleanly in UI and metrics (Battle Shout is the pilot).

**Why a manifest.** The DB has no "raid buff" flag. `ImplicitTarget` + skill-line gating narrows candidates but still leaves ungranted twins (BoK 1213408, MotW 1291335/1310503) and NPC copies. The manifest names _which_ spells are the sim's buffs; the DB supplies everything else.

## Decisions (user-confirmed)

| Decision          | Choice                                                                                                                                                                                                                                                                   |
| ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| Proto messages    | Generated from manifest into `proto/buffs.proto` (checked in; protoc needs it to bootstrap). Manifest owns field numbers; **existing numbers kept**.                                                                                                                     |
| Values            | **DB is truth.** Expected deltas listed below for review.                                                                                                                                                                                                                |
| Scope             | **Every kind** in manifest. Generator emits proto field + TS input + Go shell for all; hand-written drivers plug in by fixed name. Manual additions stay possible via `KindManual`.                                                                                      |
| Class integration | One tag convention: `ActionID.Tag 0` = player-cast, `-1` = external. Generic "(External)" rule in TS replaces per-name cases. Owner class emitted so the UI relabels the external picker. Pilot = Battle Shout, re-enable `warrior.registerShouts()`.                    |
| Tooling           | Extend existing parsers: `tools/database/gen_spelldata` binary, `dbc` package, `gen_effects.go` template/writer pattern, `gen_db -gen=go-to-ts` for TS.                                                                                                                  |
| Work location     | New worktree + branch `forever-buff-codegen`. **Not merged** by this work; PR left open for review. GitHub issue `[Core] Implement all known buffs/debuffs` created as step 0 (body below).                                                                              |
| Execution         | Opus workers (`opus-worker` / `opus-high`) on disjoint files per phase. Another agent is active on this machine: every test/build run capped at **50% of 18 cores = 9**: `GOMAXPROCS=9 /usr/local/go/bin/go test -p 9 ...`, `npx vitest run --maxWorkers 9`, `make -j9`. |

**Calls made without an explicit user answer (approve or override):**

- Ghost improved talents (not in the Forever Trait tree) change field type `TristateEffect → bool` and bump the proto to version 17. User approved DB-is-truth for _values_; this extends it to schema type. Alternative: keep `TristateEffect` and ignore the Improved state (no version bump, dead UI state).
- Owner-class picker is **relabelled "(External)", not hidden** (hiding would leave an enabled invisible buff; protection warrior defaults it on).
- Issue/commit prefix spelled `[Core]` to match the repo's commit convention (user wrote `[CORE]`).
- Correction to an earlier question: Sunder Armor in the DB is −450 per stack ×5 (not −90/stack as written in the question).

## Architecture

```
                 tools/database/wowsims.db (SQLite, gitignored, `make db`)
                                   │
  tools/database/buffmanifest/    │   Spell/SpellEffect/SpellMisc/SpellDuration/SpellAuraOptions
   (one BuffSpec per proto field)  │   SkillLineAbility/SkillLine/Trait*/SpellItemEnchantment
                 │                 │
                 ▼                 ▼
        tools/database/gen_buffs.go  ── ResolveBuffManifest(db) → []ResolvedBuff
                 │        (anchor → rank ladder → effects → duration/stacks → talent curve → owner)
                 │
     ┌───────────┼──────────────────────────────┐
     ▼           ▼                              ▼
proto/buffs.proto   sim/core/buffs_auto_gen.go    ui/features/settings/model/
(checked in)        sim/core/debuffs_auto_gen.go  buffs_debuffs_auto_gen.ts (checked in)
     │                     │                              │
  protoc              buffs_gen_support.go            buffs_debuffs.ts re-exports +
     │                buffs_manual.go (drivers)        manual rows + composed registries
     ▼                     │                              │
*.pb.go / *.ts     applyGeneratedBuffs/Debuffs      SettingsTabBody → relevantStatOptions
                           │                          → applyOwnerClassLabels
                     class code calls XAura(unit, isPlayer=true, talentPts)
                           │
                     metrics: tag 0 "Battle Shout", tag -1 "Battle Shout (External)"
```

Regen guard: `tools/database/buffs_regen_test.go` renders all three outputs and byte-diffs against committed files (skips when `wowsims.db` absent, like `spelldata_regen_test.go:84-87`).

## Manifest (`tools/database/buffmanifest/{manifest,census}.go`)

Go literal, same precedent as `overrides.go`, but in its **own leaf package `tools/database/buffmanifest`** (plain strings/ints, at most `sim/core/proto` enums). Package `database` imports `sim/core` (atlasloot.go, gen_effects.go), so the manifest cannot live there: the proto emitter must compile while `*.pb.go` is stale, and the generator must compile while `buffs_auto_gen.go` is stale.

```go
type BuffSpec struct {
    Field      string        // proto field name, owns the number: "battle_shout"
    Number     int32         // 28 — existing numbers preserved
    Scope      BuffScope     // ScopeRaid | ScopeParty | ScopeIndividual | ScopeDebuff
    Proto      BuffProtoType // ProtoBool | ProtoTristate | ProtoInt32 | ProtoDouble | ProtoEnumDrums
    Kind       BuffKind      // see taxonomy
    Go         string        // identifier stem → BattleShoutAura / BattleShoutValue / BattleShoutCategory
    Name       string        // SpellName.Name_lang of the castable family
    AuraName   string        // aura family when cast is E_SUMMON (totems: "Strength of Earth") or dummy (LotP)
    Anchor     int32         // explicit spell id; 0 = resolve Name via SkillLineAbility
    Owner      proto.Class   // validated against SkillLineAbility.ClassMask, not derived
    Talent     *TalentMod    // {Name "Booming Voice", Effect idx, Applies ScalesValue|ScalesDuration|AddsStat}
    Category   string        // exclusive category ("" = none)
    SharedCategory string    // second category the aura joins without an effect of its own ("PaladinAura"), player copy only
    SingleAura bool
    Driver     bool          // apply block calls drive<Go> instead of activating the aura outright
    Pet        PetPolicy     // PetNormal | PetStrip | PetInheritOwnerAura | PetCapAtRegular | PetStripWhenSummonedLate
    StatOverride []string    // sim stats an untyped aura lands on (A_MOD_CRIT_PCT: LotP melee, Moonkin spell); one aura effect only
    Stats      []proto.Stat  // UI relevance tags (heuristic, manifest-authoritative)
    ImpAction  *ActionRef    // override for improved icon when it is an item (30446, 32387)
    Label      string        // override only when DB name is wrong for UI
    Notes      string        // emitted as comment on manual/absent shells
}
```

Shipped as built, plus `var Retired = map[BuffScope][]int32{...}`: the field numbers api version 17 gave up, which the emitter turns into `reserved` lines and no later row may take.

Seed rows in current registry order (`PARTY_BUFFS_CONFIG` :325-368, `BUFFS_CONFIG` :370-385, `DEBUFFS_CONFIG` :466-490) so UI order is unchanged.

## Kind taxonomy

| Kind                                           | DB source                                                                                                                                                                                                                                                                                                                                                           | Emitted Go                                                                                                                                 |
| ---------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------ |
| `KindStatFlat`                                 | `E_APPLY_AURA`(6)/`AREA_AURA_PARTY`(35)/`AREA_AURA_RAID`(65) with `A_MOD_STAT` (misc -1 = all five), `A_MOD_ATTACK_POWER`, `A_MOD_RANGED_ATTACK_POWER`, `A_MOD_DAMAGE_DONE`, `A_MOD_SPELL_CRIT_CHANCE`, `A_MOD_CRIT_PCT`, `A_MOD_HIT_CHANCE`, `A_MOD_POWER_REGEN` (MP5 = v×5000/period). Reuse `dbc.ParseStatEffect` (spell_effect.go:288-350), `MapMainStatToStat` | `newGeneratedStatAura(unit, GeneratedBuff{...})`                                                                                           |
| `KindStatPct`                                  | `A_MOD_TOTAL_STAT_PERCENTAGE`, `A_MOD_ATTACK_POWER_PCT`                                                                                                                                                                                                                                                                                                             | same, `IsMultiplicative`, `1+v/100`                                                                                                        |
| `KindResistance`                               | `A_MOD_RESISTANCE` school bitmask (1 armor, 4 fire, 8 nature, 16 frost, 32 shadow, 64 arcane; bit 2 holy → comment)                                                                                                                                                                                                                                                 | `makeExclusiveFlatStatBuff(..., ResistanceCategoryX)`                                                                                      |
| `KindPseudoMult`                               | `A_MOD_THREAT`, `A_MOD_DAMAGE_PERCENT_DONE`, `A_REDUCE_PUSHBACK`, `A_MOD_MELEE_HASTE_3`                                                                                                                                                                                                                                                                             | `AttachMultiplicativePseudoStatBuff` / additive                                                                                            |
| `KindDamageShield`                             | `A_DAMAGE_SHIELD` + `SpellMisc.SchoolMask`                                                                                                                                                                                                                                                                                                                          | `newGeneratedDamageShield` (factored from `ThornsAura` buffs.go:377)                                                                       |
| `KindDebuffStat`                               | stat mapping on target; `EffectPointsPerResource` → per-combo-point                                                                                                                                                                                                                                                                                                 | `newGeneratedDebuff`, `ExclusiveEffect` priority = value                                                                                   |
| `KindDebuffStacking`                           | `SpellAuraOptions.CumulativeAura` → MaxStacks                                                                                                                                                                                                                                                                                                                       | priority = value×stacks (factored from `SunderArmorAura` debuffs.go:759)                                                                   |
| `KindDebuffDamageTaken`                        | `A_MOD_DAMAGE_PERCENT_TAKEN` mask, `A_MOD_DAMAGE_TAKEN` flat                                                                                                                                                                                                                                                                                                        | `damageTakenDebuff` / additive pseudo                                                                                                      |
| `KindDebuffAtkSpeed`                           | `A_MOD_MELEE_HASTE_3` negative                                                                                                                                                                                                                                                                                                                                      | `AtkSpeedReductionEffect`                                                                                                                  |
| `KindDebuffUptime`                             | proto `double` + `EffectTriggerSpell` one hop                                                                                                                                                                                                                                                                                                                       | shell + driver calling `ApplyFixedUptimeAura`                                                                                              |
| `KindExternalCD`                               | duration, `SpellCooldowns.RecoveryTime`                                                                                                                                                                                                                                                                                                                             | generated aura + `registerXxxCD` via `registerExternalConsecutiveCDApproximation` (buffs.go:1859); `ShouldActivate` from `buffs_manual.go` |
| `KindProc`                                     | anchor only                                                                                                                                                                                                                                                                                                                                                         | shell; driver (Windfury, JoL/JoW)                                                                                                          |
| `KindItemCount`                                | anchor effects × count                                                                                                                                                                                                                                                                                                                                              | stat aura × `float64(count)`                                                                                                               |
| `KindManual`                                   | anchor only (**escape hatch**)                                                                                                                                                                                                                                                                                                                                      | shell + `// manual driver: <Notes>` + call `manualXxx(...)` that `buffs_manual.go` must define (missing = compile error)                   |
| `KindFlag` / `KindEnum`                        | none                                                                                                                                                                                                                                                                                                                                                                | proto field only                                                                                                                           |
| `KindAbsent` (no row left; retired in phase 8) | no `SpellName` row / aura 270,271 caster-only                                                                                                                                                                                                                                                                                                                       | commented-out shell with reason (the `{{- else if not .Supported}}` idiom from `gen_effects_templates.go`)                                 |

## Generator (`tools/database/gen_buffs.go`, `gen_buffs_templates.go`)

Hooked into `tools/database/gen_spelldata/main.go` after `GenerateSpellDataFiles`. New `make spelldata` target. `RenderBuffFiles(helper) map[string][]byte` renders everything, `format.Source` gate before any write (as `renderClassFile` does).

**Anchor resolution** (no `Rank %` filter — BoK, Greater BoK, Innervate, PI have no subtext):

```sql
SELECT sla.Spell, s.NameSubtext_lang, sla.ClassMask, sla.SkillLine, sla.AcquireMethod, sla.SupercedesSpell
FROM SkillLineAbility sla JOIN SpellName n ON n.ID = sla.Spell JOIN Spell s ON s.ID = sla.Spell
JOIN SkillLine sl ON sl.ID = sla.SkillLine AND sl.CategoryID = 7 AND sl.ID NOT IN (2851, 2853)
WHERE n.Name_lang = ? AND (sla.ClassMask & ? != 0 OR sla.ClassMask = 0)
```

Rank = subtext if present, else the spell no other row `SupercedesSpell`. Pick highest; warn on non-monotonic ladder (Trueshot r5=50 < r4=75). `AuraName` set → resolve aura family by name + same rank subtext among `Effect IN (6,35,65)` (no creature table to follow `E_SUMMON`). `Anchor != 0` → use directly.

**Traps the generator must handle (all verified in DB this session):**

- `SkillLine` 2851 Engraving / 2853 Runes are `CategoryID = 7`. Exclude explicitly. Rune-granted check: `Engrave % - %` spells have `Effect=54`, `EffectMiscValue_0 = SpellItemEnchantment.ID`, `EffectArg_0..2` = rune spell. SoD-band ids 403215 (Commanding Shout), 402004, 407995 (Mangle) are NOT rune-granted → Forever content.
- Improved talents: look up **only in the owner's live `Trait*` tree** (trees: warrior 1117, paladin 1100, hunter 1091, rogue 1111, priest 1114, shaman 1082, mage 1112, warlock 1116, druid 1089). Legacy `Talent` table has ClassID/SpellID zeroed. Detection rule: trait spell with `A_ADD_PCT_MODIFIER`(108)/`A_ADD_FLAT_MODIFIER`(107) whose `EffectSpellClassMask` overlaps the buff's `SpellClassOptions.SpellClassMask`, same `SpellClassSet`. Misc → `Applies`: 1 → ScalesDuration; 0/3/8 → ScalesValue; anything else (6 radius, 14 cost, 7 crit…) ignored with a stdout note. No matching row → `TalentGhost`; field emitted as **bool**. Verified: Forever Booming Voice (12321, 5 ranks) is misc 6 radius only, so Battle Shout has **no** talent hook and a flat 3 min duration.
- Duration via `SpellMisc.DurationIndex → SpellDuration.Duration` (`-1` = NeverExpires). `shared.SpellData` lacks it; query directly.
- Level scaling: `DeriveRankAmount(e, SpellLevel, MaxLevel)` at 60, then **truncate toward zero** (-204.2 → -204).
- Stacks: `SpellAuraOptions.CumulativeAura`.
- Scope validation only (manifest wins): `ImplicitTarget_0` 1+Effect 35 party-area, 20 caster-area-party, 21 single friendly, 22/15 area enemy, 56/153 raid, 6 enemy, 57 raid-target. Add `type ImplicitTarget int` consts to `dbc/enums.go`.
- Parallel ungranted ladders excluded by the SLA gate (MotW 1291335/1310503, BoK 1213408, AI 364161, FF 1288561/1289452, CoE 11723/11724).

Stdout summary: `resolved N rows: a generated, b manual, c ghost talents, d absent` + scope mismatches + label drift vs DB name.

## Generated outputs

### `proto/buffs.proto` (generated in place, committed)

`syntax proto3; package proto; import "common.proto";` — exactly `RaidBuffs`, `PartyBuffs`, `IndividualBuffs`, `Debuffs`. `TristateEffect` stays in `common.proto`, and so does `Drums`, which `ConsumesSpec.drums_id` still names; `StrengthOfEarthType` was unreferenced and is gone. **As shipped:** api version 17 also retires the 33 fields the Forever client describes no spell for, so each message opens with a `reserved` line and the `// Next index` comment counts past the reserved numbers as well as the live ones. Emitter = tiny `tools/gen_buffs_proto/main.go` importing only the manifest package (bootstrap: protoc must run before `gen_db` can compile). Edits: delete `common.proto:453-580`; add `import "buffs.proto";` in `api.proto` and `ui.proto`. Makefile: `proto/buffs.proto` rule + prerequisite on `sim/core/proto/api.pb.go` (:192) and `ui/generated/proto/api.ts` (:53). Proto messages will not be `DO NOT EDIT`-only; the manifest is edited instead.

**Versioning (mandatory):** ghost-talent fields change `TristateEffect → bool` (~15 fields). Settings persist as JSON and `migrateOldProto` runs on the **already-parsed** proto (`ui/sim/state/serialization.ts:39`, `ui/app/proto_version.ts:60`), so protobuf-ts `fromJson` throws on `"TristateEffectImproved"` in a bool field before any converter runs. Therefore: bump `current_version_number` 16→17 in `common.proto`; add a **raw-JSON pre-pass** at every `IndividualSimSettings.fromJson` / `Player.fromJson` entry (`ui/app/preset_utils.ts:279`, storage load in `ui/sim/state/serialization.ts`, share-link decode) that, when `apiVersion < 17`, rewrites the listed field names from tristate string → boolean; then the normal v17 stamp in `proto_version.ts`; update `ui/app/storage_keys.test.ts:58`; `make update-tests` for `TestProtoVersioning.results`. Binary/share-link payloads are fine either way (varint 1 or 2 decodes to `true`). Add `proto_version.test.ts` case: v16 JSON with `battleShout: "TristateEffectImproved"` loads as `true`. Also switch `buf.yaml` breaking to `PACKAGE` so the same-package file move is not itself flagged (probe with `npx buf breaking` first; fall back to FILE + bump, which we need anyway).

### Go: `sim/core/buffs_auto_gen.go`, `sim/core/debuffs_auto_gen.go`

Per row (data literals + calls into a thin hand-written support API so generated code is nearly type-error-proof):

```go
// Battle Shout - https://www.wowhead.com/forever/spell=25289
var BattleShoutCategory = "BattleShout"
func BattleShoutValue(talentPoints int32) float64 { return 139 }
func BattleShoutDuration(talentPoints int32) time.Duration { return 180 * time.Second }
func BattleShoutAura(unit *Unit, isPlayer bool, talentPoints int32) *Aura {
    return newGeneratedStatAura(unit, GeneratedBuff{
        Label:    "Battle Shout (" + Ternary(isPlayer, "Player", "External") + ")",
        ActionID: ActionID{SpellID: 25289}.WithTag(TernaryInt32(isPlayer, 0, -1)),
        Duration: BattleShoutDuration(talentPoints), IsPlayer: isPlayer,
        Category: BattleShoutCategory, SingleAura: true,
        Stats:    []StatConfig{{stats.AttackPower, BattleShoutValue(talentPoints), false}},
    })
}
```

Plus `applyGeneratedBuffs(char, raid, party, individual)` / `applyGeneratedDebuffs(target, debuffs, raid)` — one `if` block per row, calling `MakePermanent(...)` or the named driver (`driveBattleShout`, `windfuryTotemDriver`, `innervateShouldActivate`, ...). `applyBuffEffects` (called character.go:304) and `applyDebuffEffects` (environment.go:93) become the generated call + shrinking hand-written remainder.

New hand-written files: `sim/core/buffs_gen_support.go` (`GeneratedBuff`, `newGeneratedStatAura`, `newGeneratedDebuff`, `newGeneratedDamageShield`, `newGeneratedExternalCD`), `sim/core/buffs_manual.go` (drivers).

**As shipped:** `makeStatBuff`'s tag rewrite `0→-1` is gone and `BuildPhase = Ternary(Tag == -1, Buffs, None)` reads the explicit tag; its one remaining caller is the draenei racial, which states `-1`. Every other hand-written buff body whose row is generated, retired or unsupported is deleted, leaving `BuffConfig`/`StatConfig` and the `registerExlusiveEffects`/`registerStatEffect` family the support API calls, `ApplyFixedShoutAura`, `registerExternalConsecutiveCDApproximation`, `registerBloodlustCD`/`BloodlustAura`, `ScheduledAura`, and the two rows the client has a named gap for (`MangleAura`, `ImprovedSealOfTheCrusaderAura`). Paladin auras `1/-1` (buffs.go:776-778) and shouts `0/1` collapse to `0/-1`. Class stubs (`sim/paladin/auras.go`, `sim/shaman/totems.go`, …) get their commented bodies updated to the new signatures, still stubbed (repo rule: stub, don't delete).

Import cycle: generated values are literals in `sim/core`; never read `sim/<class>` tables. `tools/database` imports `sim/core`, so a stale `*_auto_gen.go` (e.g. after a proto field retype) breaks the generator itself. Recovery, documented in `docs/spell_data.md`: cut the lines of the generated file that name the field or symbol that is going away by hand, `go build ./...`, then run `gen_spelldata`, which writes the whole file back. No `make buffs-regen` target was added - the cut is a handful of lines and differs per change, so a target would only hide it. Plain `git checkout` is not a recovery when the committed file is the stale one.

### TS: `ui/features/settings/model/buffs_debuffs_auto_gen.ts` (**committed**, not gitignored)

CI runs `make go-to-ts` before `type-check` with **no `wowsims.db`**, and this file needs the DB (icons, labels, spell IDs). So it is emitted by `gen_spelldata` in the same DB run as the Go files, like `spell_data_auto_gen.go`. Emitter `tools/database/gen_buffs_debuffs_ts.go` → `GenerateBuffsDebuffsTSFile()`. Add a negated `.gitignore` entry (`!ui/features/settings/model/buffs_debuffs_auto_gen.ts`) and explicit ignore entries in the oxlint/oxfmt configs (today they skip `*_auto_gen.ts` only because it is gitignored). Not added to `AUTO_GEN_FILES_TS`. **Factory is chosen from the emitted proto type**, never from the legacy kind: bool → `makeBoolean*Input`, tristate (talent present in the live tree) → `makeTristate*Input` with `impId` = the talent's top-rank spell, int32 → `makeMultistate*Input`, quadstate only when a second bool field survives. **Item anchors are checked against the `Item` table**; missing ones (30446 Solarian's Sapphire, 32387, the four pendants; T2 bonus 23563 as an aura source) leave the row with nothing to resolve. **As shipped:** such a row is retired rather than kept inert - it leaves the manifest and the proto, its number is `reserved`, and the v17 pre-pass drops its key from an older payload. The manifest has no `KindAbsent` row left. Per row: `export const BattleShout = makeBooleanPartyBuffInput({ actionId: ActionId.fromSpellId(25289), fieldName: 'battleShout', label: 'Battle Shout' })` + registries `GENERATED_{RAID,PARTY,INDIVIDUAL,DEBUFFS}_CONFIG: RenderableStatOptions[]` with `{ config, stats, ownerClass }`. `StatOption` carries `ownerClass`, so the registries use `RenderableStatOptions` directly.

`buffs_debuffs.ts` shrinks to: `export * from './buffs_debuffs_auto_gen'` + the three rows the manifest cannot produce (`Bloodlust`, `BlessingOfSalvation` role-gated, `ShadowPriestDPS`) + the `Innervate`/`PowerInfusion`/`ManaTideTotem` aliases the specs name + composed registries. The rows for fields api version 17 retired - the drums swatch, the draenei racials, the necks, the absent debuffs - are gone with their fields. `ui/specs/**` imports of `BuffDebuffInputs.X` keep working. Windfury `StatParryRating` sentinel removed; feral specs exclude `[BuffDebuffInputs.WindfuryTotem]` instead.

24 files change `RaidBuffs|PartyBuffs|IndividualBuffs|Debuffs` import from `@generated/proto/common` to `@generated/proto/buffs` (list in section 5 of the proto/TS design: specs presets, icon_inputs, input_helpers, saved_settings, serialization, raid/party, …). No field renames → zero preset value edits.

## Class integration + metrics

- Contract: every generated constructor `XAura(unit *Unit, isPlayer bool, talentPoints int32)`. Class code passes `isPlayer=true` + its real talent points; `applyGeneratedBuffs` passes `false, 0`.
- UI: add `ownerClass?: Class` to `StatOption` and `applyOwnerClassLabels(options, player)` in `stat_options.ts`, called in `SettingsTabBody.tsx` **after** `relevantStatOptions` (include/exclude matches config by reference). **Relabel "(External)", don't hide**: `ApplyFixedShoutAura` chains the external shout after the player's expires, and protection warrior defaults it on; hiding would leave an invisible enabled buff. No `showWhen` plumbing on tristate/quadstate factories is needed for this (relabel happens in the registry layer).
- Metrics (`ui/sim/proto/action_id/index.ts`): the `default:` branch already maps tag -1 → " (External)". Delete `case 'Battle Shout'/'Commanding Shout'` (:481-488; tag 1 gone, tag 3 never emitted by Go) and the six-name `(Self)/(External)` block (:316-325). Keep the raid-index block (:441-459 Bloodlust/Innervate/Mana Tide/PI/Ferocious Inspiration — there -1 means "(raid)") and Drums. `AuraMetricsTable` grouping unchanged: parent "Battle Shout" with sub-rows "Battle Shout" + "Battle Shout (External)".
- `ui/specs/warrior/shared/inputs.ts:15` ShoutPicker icon 2048 → 25289 (2048 has no DB row; tooltip lookup fails).

## Battle Shout pilot (Phase 2)

1. Manifest row: `{Field:"battle_shout", Number:28, Scope:ScopeParty, Proto:ProtoBool (Imp BS not in tree), Kind:KindStatFlat, Go:"BattleShout", Name:"Battle Shout", Owner:Warrior, Talent:nil (Booming Voice is radius-only in Forever), Category:"BattleShoutCategory", SingleAura:true, Stats:[AP,RAP]}`. Rows `bs_solarian_sapphire` (item 30446 absent from `Item`) and `snapshot_bs_*` have no UI input. **As shipped:** all four are retired in phase 8 and their numbers are `reserved`.
2. Generate → `BattleShoutAura/Value/Duration/Category` + apply block calling `driveBattleShout`.
3. `buffs_manual.go`: `driveBattleShout` wraps `ApplyFixedShoutAura` (it already reads `Tag == 0` as player, buffs.go:676).
4. Delete buffs.go:624-668 and :181-192. (`makeStatBuff`'s tag rewrite stays until Phase 8: generated code uses `newGeneratedStatAura`, and still-hand-written stat buffs depend on the rewrite for `BuildPhase Buffs`.)
5. `sim/warrior/shouts.go`: `core.BattleShoutAura(&warrior.Unit, DefaultShout == Battle, Talents.BoomingVoice)`; `ExtraCastCondition` compares to `core.BattleShoutValue(...)`; drop `commandingPresenceMultiplier`, sapphire/T2 args from this path (struct fields stay until proto flags are retired); `battleShoutRank = spellData.BattleShout.BySpellID(25289)`.
6. `sim/warrior/warrior.go:175` uncomment `registerShouts()`. `demoralizing_shout.go:17` keeps hand-written `DemoralizingShoutAura` until Phase 4.
7. Expect: aura 25289 tag 0 "Battle Shout (Player)" +139 AP, flat 3 min; external tag -1. `warrior.Talents.BoomingVoice` is still passed as `talentPoints` for signature uniformity and ignored. `TestDpsWarrior.results` / `TestProtectionWarrior.results` move (AP 306→139) → `make update-tests`, DPS delta in PR body.

## Expected value deltas (DB @60 vs hand-written) — review list

| Buff                                    | Now                       | DB                                         | Note                                                               |
| --------------------------------------- | ------------------------- | ------------------------------------------ | ------------------------------------------------------------------ |
| Battle Shout                            | 306 AP, 2 min, spell 2048 | 139 AP, 3 min, 25289                       | Imp BS ghost; Booming Voice = radius only, no sim effect           |
| Commanding Shout                        | 1080 HP                   | 42 Stam, 5 min, 403215                     |                                                                    |
| Arcane Brilliance                       | 40 Int                    | 31                                         |                                                                    |
| PW:Fortitude / Divine Spirit            | 79 Stam / 50 Spi          | 70 / 40                                    | Imp ghosts; no DS 10% conversion                                   |
| Gift of the Wild                        | 340 armor/14 stats/25 res | 385/16/27                                  | includes holy (no stat)                                            |
| Thorns / Retribution Aura               | 25 / 26                   | 22 / 30                                    |                                                                    |
| Devotion Aura                           | 861 armor                 | 735                                        |                                                                    |
| Trueshot Aura                           | 125 AP+RAP                | 50 RAP (r4 = 75!)                          | non-monotonic, flag                                                |
| Blood Pact                              | 70 Stam                   | 54                                         |                                                                    |
| Moonkin / LotP                          | 5% crit (+rating)         | 3%                                         | `A_MOD_CRIT_PCT` can't split melee/spell → manifest `StatOverride` |
| Blessing of Might                       | 220 AP+RAP                | 133 melee AP only                          |                                                                    |
| Blessing of Wisdom                      | 41 MP5                    | 40                                         |                                                                    |
| SoE / GoA / Mana Spring                 | 86 Str / 77 Agi / 50 MP5  | 53 / 89 / 25                               |                                                                    |
| Windfury Totem                          | 445 AP                    | 246                                        |                                                                    |
| Resist auras/totems, Aspect of the Wild | 70                        | 60                                         |                                                                    |
| Power Infusion                          | cast speed/cost           | +20% dmg & healing done, 15 s              | **semantics change**                                               |
| Sunder Armor                            | -520 ×5                   | -450 ×5, threat 1013                       |                                                                    |
| Expose Armor                            | -2050 @5cp                | -450/cp = -2250                            | Imp EA in tree                                                     |
| Faerie Fire                             | -610                      | -505, 40 s                                 |                                                                    |
| Curse of Recklessness                   | -800 armor / +135 AP      | -505 / dummy 90                            |                                                                    |
| Curse of Elements                       | 4 schools 10%, -88        | mask 126 incl. holy+nature 10%, -75, 5 min | Malediction in tree                                                |
| Hunter's Mark                           | 110 +11/stack             | 71 flat                                    |                                                                    |
| Demo Roar / Shout                       | -248 / -300, 30 s         | -204 / -204, 45 s                          |                                                                    |
| Thunder Clap                            | -10%                      | -20%                                       | Imp TC in tree                                                     |
| Scorpid Sting                           | -5 hit                    | -2                                         |                                                                    |
| JotC                                    | 219 holy taken            | 161, 40 s                                  | ISotC ghost                                                        |

**Ghost improved talents** (not in live tree → field becomes bool): Imp Battle Shout, BoM, BoW, Devotion, Retribution, Concentration, Sanctity, PW:F, Divine Spirit, MotW, Blood Pact, LotP, Moonkin, Enhancing Totems, Imp WF, Imp Hunter's Mark, Imp Demo Shout, Feral Aggression, Brambles, Imp FF, ISotC.

**Absent in client**: Sanctuary, Sanctity (no SLA), Tranquil Air cast, Wrath of Air, Totem of Wrath, Drums, four neck items, Draenei presence, Misery, Ferocious Inspiration, Shadow Embrace, Screech, Unleashed Rage, Blood Frenzy, Expose Weakness, ISB debuff 17800. **As shipped:** these are retired, not kept as inert fields - the row leaves the manifest, the field leaves the proto, its number joins `buffmanifest.Retired` and is `reserved`, and the v17 pre-pass drops the key from an older payload.

**Caster-only auras** (270/271: Shadow Weaving, Improved Scorch, Hemorrhage) grant the simmed player nothing as external debuffs, so they are retired with the rest.

## Policy defaults (assumptions, flagged)

- Bloodlust: no SLA row for any candidate (1245940 / 1222564 / 468408). Keep `KindManual` with today's `registerBloodlustCD`, `Anchor: 0`, open item in issue.
- Ghost tristate → bool + v17 converter (DB truth). Reversible per row by setting `Proto: ProtoTristate` + explicit `Talent`.
- Absent fields keep their proto number and UI input is not emitted (no icon to show); issue tracks them.

## Phases

0. **Setup**: worktree `../wowsims-forever-buff-codegen`, branch `forever-buff-codegen`; copy this plan into the worktree as `docs/superpowers/specs/2026-09-20-buff-codegen-design.md` and commit; create GitHub issue (body below); `npx buf breaking` PACKAGE probe.
1. **Generator scaffolding** — `tools/database/buffmanifest/{manifest,census}.go` (full census, most rows Manual/Absent), `gen_buffs.go`, `gen_buffs_templates.go`, `gen_buffs_debuffs_ts.go`, `tools/gen_buffs_proto/main.go`, `dbc/enums.go` (ImplicitTarget), `gen_spelldata/main.go` hook, `gen_db/main.go` go-to-ts hook, makefile (`spelldata`, `proto/buffs.proto`, `AUTO_GEN_FILES_TS`), `buffs_regen_test.go` + manifest-invariants test, `sim/core/buffs_gen_support.go`, empty `sim/core/buffs_manual.go`. Gate: all outputs compile with every row a shell; `proto/buffs.proto` byte-identical to today's messages.
2. **Proto split (types unchanged)** — `proto/common.proto`, `api.proto`, `ui.proto`, `buf.yaml`, 24-file TS import sweep. Gate: `make sim/core/proto/api.pb.go ui/generated/proto/api.ts && npm run type-check && go test -tags with_db ./sim/core -run TestProtoVersioning`. **Tristate→bool retypes happen per buff in its migration phase** (4, 5, 6); the version 17 bump, the **raw-JSON pre-pass** (see Versioning: `migrateOldProto` runs on the parsed proto, so `fromJson` would throw first), `storage_keys.test.ts` and `make update-tests` land with the first retype in Phase 4 and the pre-pass field list grows per phase. **Every retype needs a type-change sweep**: TS presets/spec defaults setting `TristateEffect.*` on that field (14 files under `ui/specs/**`, `grep -rl TristateEffect ui/specs`) and Go test configs setting `proto.TristateEffect_*` (59 usages, `grep -rn TristateEffect_ sim --include=*.go | grep -v pb.go`) flip to `true`; the "zero preset edits" claim holds only for untouched fields.
3. **UI registry plumbing** — `stat_options.ts` ownerClass + `applyOwnerClassLabels` + test, `SettingsTabBody.tsx`, feral specs Windfury exclude, `buffs_debuffs.ts` re-export/compose. (No `showWhen` plumbing: relabel lives in the registry layer.)
4. **Battle Shout pilot** — as above + `action_id/index.ts` case deletions + `warrior/shared/inputs.ts` + `battle_shout` retype sweep + v17 pre-pass. Gate: warrior suites, aura table shows base vs "(External)", v16 JSON with `battleShout: "TristateEffectImproved"` loads as `true`.
5. **Stat/resistance/pct buffs** — flip AB, PWF, DS, GotW, SP, BoK, BoM, BoW, BoS, Blood Pact, Moonkin, LotP, Devotion, Concentration, resist auras/totems, Aspect of the Wild, SoE, GoA, Mana Spring, Atiesh, Commanding Shout; delete from buffs.go; update paladin/shaman stub comments to new signatures.
6. **Debuffs** — Sunder, EA, FF, CoR, CoE, Demo Roar/Shout, TC, HM, Scorpid, IS, Gift of Arthas, Winter's Chill, JotC; update `sim/warrior/{sunder_armor,demoralizing_shout,thunder_clap}.go` call sites.
7. **Drivers** — external CDs (Innervate, PI, Mana Tide, Pain Suppression), damage shields, Windfury/JoL/JoW/ISB, `applyPetBuffEffects` from `Pet` policy.
8. **Cleanup + docs** — remove `makeStatBuff` tag rewrite (buffs.go:113-115); buffs.go/debuffs.go down to helpers; `docs/spell_data.md` "Buffs" section (manifest, kinds, recovery path). `.results` regeneration as its own reviewed commit.

Each phase = own commit(s) on `forever-buff-codegen` in worktree `../wowsims-forever-buff-codegen`; Opus workers (`opus-worker`) on disjoint files per phase, `opus-high` for review passes; `sonnet-verifier` only for gate runs. Nothing is merged; final state is an open PR. All Go/TS test invocations use the 9-core cap above because another agent shares the machine.

## Verification

- `/usr/local/go/bin/go run ./tools/database/gen_spelldata` then `git diff --stat` shows only intended generated files.
- `/usr/local/go/bin/go build ./... && /usr/local/go/bin/go test -tags with_db ./sim/core/... ./sim/warrior/... ./tools/database/...` (never via rtk; it masks failures).
- `make go-to-ts && npm run type-check && npm run lint:js && npm run fmt && npm run test:unit`.
- Regen guard: `go test ./tools/database/ -run TestGeneratedBuffFiles` (skips without DB; add to CI workflow since `run_tests.yml:129` only runs `./sim/...`).
- Manual: `make host`, open warrior DPS → settings shows "Battle Shout (External)" with DB icon next to ShoutPicker; hunter shows "Battle Shout"; feral cat no Windfury; run sim with DefaultShout=Battle + external set → aura table parent "Battle Shout" with two sub-rows. Load a v16 saved settings JSON with `battleShout: "TristateEffectImproved"` → converter yields `true`, toast shown.
- Before/after: keep HEAD `.results`, list per-test DPS delta in PR body alongside the value table.

## GitHub issue (step 0) — `[Core] Implement all known buffs/debuffs from the Forever client`

Label: `enhancement`. Body = checklist grouped by scope, one line per manifest row: field, Forever anchor spell, kind, status (`generated` / `manual driver` / `class bridge stubbed` / `ghost talent → bool` / `absent in client` / `open: anchor unknown`), owner class. Sections: "Generated from DB" (list), "Needs hand-written driver" (Bloodlust, Innervate, PI, Mana Tide, Pain Suppression, Windfury, JoL/JoW, ISB, Thorns, Retribution, Shadow Priest DPS, Drums), "Class bridges to port" (paladin auras, shaman totems, druid FF/roar/LotP/Moonkin, rogue EA, warlock CoE/CoR, mage Scorch/WC, priest PWF/DS/SP/PI, hunter HM/Trueshot/Scorpid), "Ghost talents" (list above), "Absent in client" (list above), "Open questions" (Bloodlust anchor, PI semantics, Trueshot ladder, caster-only auras). Link this plan's branch.
