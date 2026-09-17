# Forever talents: per-rank values and tooltips vs beta client (17 Sep 2026)

Compared `tools/forever_talents` trees (the `description` + `ranks` each talent tooltip is built from) with the Forever
beta client `1.60.1.69893` (wago.tools Trait tables) and Wowhead Forever tooltips.
**Scope: tree data and tooltips only. The Go sim's own talent code was not checked.**

## How the client stores rank values (verified)

Each live talent's trait definition has a `TraitDefinitionEffectPoints` curve per effect; `CurvePoint(rank)` is the value at
that rank. Checked on Ignite (8/16/24/32/40), Ruin (20..100), Improved Life Tap (10/20), Deflection (2..10).
`SpellEffect.EffectBasePointsF` is only the max-rank value, sometimes stale. Durations are stored in milliseconds and
rage in tenths.

## Summary (464 talents matched by name)

| Result | Talents |
|---|---|
| Values match at every rank | 206 |
| Off by at most 1 (rounding: sim split a Classic total evenly, client rounds per rank) | 46 |
| Same value, different units (ms, rage ×10) | 14 |
| **Real value errors** | **~30** |
| Different number of values in the tooltip (not auto-comparable, needs a human) | 112 |
| No numbers in the client tooltip | 54 |
| Tooltip wording differs | 4 |

## 1. Tree stores the per-point value at every rank (tooltip shows rank 1 even at max rank)

All `ranksSource: manual`. Client value at max rank = sim value × max rank, so the `ranks` arrays are not cumulative.
Check whether the sim multiplies by points spent; the tooltip is wrong either way.

| Talent | Tree says (every rank) | Client at max rank |
|---|---|---|
| Rogue / Combat / Hack and Slash | 1, 1, 3 | 5, 5, 15 |
| Rogue / Combat / Puncturing Wounds | 10, 5, 15 | 30, 15, 45 |
| Rogue / Assassination / Improved Expose Armor | 5, 1 | 10, 2 |
| Warrior / Arms / Weaponmaster | 1, 3, 1 | 5, 15, 5 |
| Warrior / Fury / Dual Wield Specialization | 5, 20, 2 | 25, 100, 10 |
| Warrior / Protection / Improved Shield Wall | 5.5 min | 11 min |
| Warlock / Demonology / Improved Health Funnel | 20, 15, 50 | 40, 30, 100 |
| Warlock / Demonology / Improved Felhunter | 10%, 2 sec | 30%, 6 sec |
| Mage / Frost / Permafrost | 11, 3 | 33, 10 |
| Paladin / Protection / Sacred Duty | 2%, 30 sec | 4%, 60 sec |
| Paladin / Retribution / Instrument of Law | 0.5 sec, 10% | 1.0 sec, 20% |
| Priest / Discipline / Renewed Hope | 2%, 1 sec | 10%, 5 sec |
| Shaman / Elemental / Improved Fire Nova | 10%, 2 sec | 20%, 4 sec |
| Shaman / Restoration / Improved Reincarnation | 10 min, 2%, 10% | 20 min, 4%, 20% |
| Druid / Balance / Moonfury | 5% | 10% |

## 2. Different numbers

| Talent | Tree (max rank) | Client (max rank) | Tree source |
|---|---|---|---|
| Druid / Balance / Moonglow | 9% | 25% | classic-prior |
| Rogue / Assassination / Lethality | 30% | 20% | classic-prior |
| Rogue / Subtlety / Quietus | 10%, 175 | 10%, 35 (client: "targets below 35% health") | extrapolated |
| Warrior / Fury / Enrage | 150% | 10% | classic-prior |
| Warlock / Affliction / Soul Siphon | 51, 30 | 12, 12 | extrapolated |
| Priest / Shadow / Early Demise | 40% | 20% | extrapolated |
| Priest / Holy / Spiritual Guidance | 25%, 5% | 25%, 8% | extrapolated |
| Paladin / Holy / Illumination | 100%, 250% | 100%, 50% | extrapolated |
| Paladin / Protection / Reckoning | 40%, 100% | 40%, 40% | extrapolated |
| Mage / Frost / Improved Blizzard | 45% | 40% | classic-prior |
| Mage / Arcane / Improved Counterspell | 2 sec | 4 sec | manual |
| Hunter / Survival / Entrapment | 1 sec | 5 sec | manual |
| Hunter / Beast Mastery / Summon Hawk | 53, 18 | 32, 5 | observed |
| Hunter / Marksmanship / Improved Stings | 6%, 2 sec, 15 sec | 20%, 6 sec, 45 sec | manual |
| Priest / Discipline / Improved Mana Burn | 0.5 sec | 1.0 sec | manual |
| Shaman / Enhancement / Improved Ghost Wolf | 2 sec | 3 sec | classic-prior |
| Shaman / Elemental / Elemental Alacrity | 0.17 sec | 0.5 sec | manual |

Pattern: most `extrapolated` and several `classic-prior` values are wrong. `classic-prior` rounding (46 talents, e.g.
9 vs 10, 51 vs 50) is harmless in play but should take the client's per-rank numbers.

## 3. Tooltip wording

- Wording differs from the client: Druid Feral Charge; Precision (Paladin, Rogue, Warrior).
- 112 talents have a different number of values in the tooltip than the client (e.g. the client puts a duration in
  `$d`, the tree puts it in `{1}`). Not auto-comparable; full list in the working report.

## 4. Wowhead Forever tooltips are not reliable yet

- Wowhead renders the spell's base points, not rank values. Some tooltips show the max-rank number, others rank 1.
- Talents whose values exist only in the rank curve show **0%** on Wowhead: Illumination, Reckoning, Enrage.
- Use the client tables, not Wowhead tooltips, as the source of truth for talent numbers until Wowhead fixes this.
