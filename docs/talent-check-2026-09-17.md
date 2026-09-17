# Forever talent trees: vendored sim data vs beta client 1.60.1.69893

> **Checked 2026-09-17: none of the seven structural differences below were applied.** An independent
> source contradicts every one of them, and two would break shipped presets or collide two talents in one
> slot. See "Verdict on the seven structural differences" at the end before acting on anything here.

Vendored = `ui/core/talents/trees/*.json` in ElliotWood/Forever (the trees the sim runs on). Client = wago.tools DB2 export of build 1.60.1.69893: TraitTree/TraitNode/TraitNodeEntry/TraitDefinition/TraitEdge/TraitNodeGroup*/TraitCond, names from TraitDefinition.OverrideName_lang or SpellName.

## Key finding on the source tables

- The beta `Talent` and `TalentTab` tables are identical to Classic Era 1.15.9.69722 (432 rows, only ClassID zeroed): they are leftover Classic placeholders and are NOT the Forever trees.
- Forever ships its trees in the retail-style Trait system: one TraitTree per class (1082 Shaman, 1089 Druid, 1091 Hunter, 1100 Paladin, 1111 Rogue, 1112 Mage, 1114 Priest, 1116 Warlock, 1117 Warrior), currency 3820 (max 51 points), three tabs laid side by side (PosX bases 1020/5020/9080, 600 per column), rows at PosY 2130 + 600*row, row gates via TraitNodeGroup conditions at 5/10/.../30 points. Trees 1081 (16 Elemental nodes) and 1083 (36 Feral nodes) look like partial prototypes and are ignored; 1187-1189 are a 16-point profession/"Legacy" system.
- Tab names are not in the Trait tables, so client tabs are identified by talent-name overlap with the vendored trees.

## Conventions normalised

- Vendored rowIdx/colIdx are 0-based; client row/col are derived 0-based from PosY/PosX, so no offset.
- Names compared case- and punctuation-insensitively; unmatched talents at the same tab/row/col are paired as renames.
- Max rank = TraitNodeEntry.MaxRanks (client) vs maxPoints (vendored). Prereq = TraitEdge (types 2/3, visual-only type 0 ignored, edges pointing upward ignored) vs vendored prereqLocation. Spell id = TraitDefinition.SpellID vs vendored spellIds[0] (does not count against exact match).
- Client quirks handled: nodes with a stray extra zero in PosX/PosY scaled back; where two nodes share a definition the higher (newer) node id is kept and the other reported as stale.

## Headline

Vendored talents 470, client talents 469, exact layout matches 461 (98.1% of vendored), matched with differences 7, only vendored 2, only client 1. Among differences: max rank 0, position/tree 4, prerequisite 1, rename 2; spell id differs on 133 matched talents.

## Summary

| Class | Tree | Vendored | Client | Exact matches | Position/rank/prereq mismatches | Missing (only vendored) | Extra (only client) |
|---|---|---|---|---|---|---|---|
| Druid | Balance | 17 | 16 | 16 | 0 | 1 | 0 |
| Druid | Feral Combat | 19 | 19 | 19 | 0 | 0 | 0 |
| Druid | Restoration | 16 | 16 | 16 | 0 | 0 | 0 |
| Hunter | Beast Mastery | 16 | 16 | 16 | 0 | 0 | 0 |
| Hunter | Marksmanship | 16 | 17 | 16 | 0 | 0 | 1 |
| Hunter | Survival | 18 | 18 | 18 | 0 | 0 | 0 |
| Mage | Arcane | 18 | 18 | 18 | 0 | 0 | 0 |
| Mage | Fire | 17 | 17 | 17 | 0 | 0 | 0 |
| Mage | Frost | 19 | 19 | 19 | 0 | 0 | 0 |
| Paladin | Holy | 18 | 18 | 18 | 0 | 0 | 0 |
| Paladin | Protection | 16 | 16 | 16 | 0 | 0 | 0 |
| Paladin | Retribution | 18 | 18 | 18 | 0 | 0 | 0 |
| Priest | Discipline | 18 | 18 | 18 | 0 | 0 | 0 |
| Priest | Holy | 17 | 17 | 17 | 0 | 0 | 0 |
| Priest | Shadow Magic | 18 | 18 | 18 | 0 | 0 | 0 |
| Rogue | Assassination | 17 | 17 | 17 | 0 | 0 | 0 |
| Rogue | Combat | 17 | 17 | 15 | 2 | 0 | 0 |
| Rogue | Subtlety | 19 | 19 | 19 | 0 | 0 | 0 |
| Shaman | Elemental Combat | 16 | 16 | 16 | 0 | 0 | 0 |
| Shaman | Enhancement | 18 | 18 | 18 | 0 | 0 | 0 |
| Shaman | Restoration | 16 | 16 | 14 | 2 | 0 | 0 |
| Warlock | Affliction | 17 | 17 | 16 | 1 | 0 | 0 |
| Warlock | Demonology | 19 | 19 | 19 | 0 | 0 | 0 |
| Warlock | Destruction | 16 | 16 | 16 | 0 | 0 | 0 |
| Warrior | Arms | 17 | 17 | 17 | 0 | 0 | 0 |
| Warrior | Fury | 18 | 18 | 18 | 0 | 0 | 0 |
| Warrior | Protection | 19 | 18 | 16 | 2 | 1 | 0 |

# Detail

## Druid (client TraitTree 1089)

Vendored tree order: Balance, Feral Combat, Restoration. Client tabs left to right map to: Balance, Feral Combat, Restoration. Client has no tab names in Trait tables (TalentTab is the unchanged Classic table).


### Balance

- **Genesis**: spell id: vendored 11124, client 1223081
- **Nature's Majesty**: spell id: vendored 18459, client 1223082
- **Nature's Splendor**: spell id: vendored 0, client 1223083
- **Overgrowth**: spell id: vendored 0, client 17245
- **Eclipse**: spell id: vendored 0, client 408248
- ONLY VENDORED **Balance of Nature** (row 2, col 3, 5 ranks)

### Feral Combat

- **Feral Charge**: spell id: vendored 16979, client 1238122
- **Shredding Attacks**: spell id: vendored 0, client 16966
- **Mangle**: spell id: vendored 0, client 407995
- **Predatory Instincts**: spell id: vendored 16493, client 1223242
- **King of the Jungle**: spell id: vendored 0, client 417046
- **Natural Reaction**: spell id: vendored 0, client 417051
- **Rend and Tear**: spell id: vendored 19616, client 1223246
- **Berserk**: spell id: vendored 0, client 417141

### Restoration

- **Naturalist**: spell id: vendored 0, client 17069
- **Gift of the Earthmother**: spell id: vendored 0, client 414673
- **Living Spirit**: spell id: vendored 16836, client 1309631
- **Wild Growth**: spell id: vendored 0, client 408120

## Hunter (client TraitTree 1091)

Vendored tree order: Beast Mastery, Marksmanship, Survival. Client tabs left to right map to: Beast Mastery, Marksmanship, Survival. Client has no tab names in Trait tables (TalentTab is the unchanged Classic table).

Stale duplicate client nodes ignored: Lightning Reflexes (node 104982, raw tab2 r6 c2)


### Beast Mastery

- **Focused Fire**: spell id: vendored 0, client 1223755
- **Summon Hawk**: spell id: vendored 0, client 1293241

### Marksmanship

- **Improved Stings**: spell id: vendored 19464, client 1310661
- **Careful Aim**: spell id: vendored 30902, client 1223984
- **Rapid Killing**: spell id: vendored 0, client 415405
- **Lone Wolf**: spell id: vendored 0, client 415370
- **Trueshot Aura**: spell id: vendored 19506, client 1299346
- **Rapid Recuperation**: spell id: vendored 0, client 1223987
- **Sniper Shot**: spell id: vendored 19434, client 1310687
- ONLY CLIENT **Improved Serpent Sting**: row 3, col 3 (0-based), max rank 5, spell 19464, prereq none
  - *client data note*: raw PosX/PosY 6820,39300 off-grid (extra zero), scaled
  - vendored ['Improved Stings'] carries this spell id; client also has ['Improved Stings'] as a separate newer node, so this is probably a stale node the client never cleaned up (inference)

### Survival

- **Predator's Edge**: spell id: vendored 0, client 1310627
- **Resourcefulness**: spell id: vendored 0, client 440529
- **Expose Prey**: spell id: vendored 0, client 1310532
- **Survivalist's Discipline**: spell id: vendored 0, client 1310496
- **Strider Kick**: spell id: vendored 0, client 1317257
- *client data note* Lightning Reflexes: duplicate stale node 104982 at tab2 r6 c2 ignored
- **Lacerating Strikes**: spell id: vendored 0, client 1310533

## Mage (client TraitTree 1112)

Vendored tree order: Arcane, Fire, Frost. Client tabs left to right map to: Arcane, Fire, Frost. Client has no tab names in Trait tables (TalentTab is the unchanged Classic table).


### Arcane

- **Arcane Geometry**: spell id: vendored 11100, client 11247
- **Arcane Blast**: spell id: vendored 30451, client 400574
- **Missile Barrage**: spell id: vendored 44404, client 400588

### Fire

- **Flame Throwing**: spell id: vendored 11078, client 11100
- **Hot Streak**: spell id: vendored 44445, client 400624

### Frost

- **Ice Lance**: spell id: vendored 30455, client 1312002
- **Fingers of Frost**: spell id: vendored 44543, client 400647

## Paladin (client TraitTree 1100)

Vendored tree order: Holy, Protection, Retribution. Client tabs left to right map to: Holy, Protection, Retribution. Client has no tab names in Trait tables (TalentTab is the unchanged Classic table).


### Holy

- **Improved Holy Strike**: spell id: vendored 0, client 1310902
- **Voice of Truth**: spell id: vendored 0, client 1310897
- **Reverence**: spell id: vendored 0, client 1310899
- **Purifying Power**: spell id: vendored 31825, client 429144
- **Infusion of Light**: spell id: vendored 53569, client 426065
- **Divine Precision**: spell id: vendored 0, client 1310904
- **Holy Shock**: spell id: vendored 20473, client 1311606
- **Consecrated Ground**: spell id: vendored 0, client 1310905
- **Light's Vigil**: spell id: vendored 0, client 1310911

### Protection

- **Improved Seal of Fury**: spell id: vendored 0, client 1314103
- **Shield Specialization**: spell id: vendored 20148, client 20150
- **Sacred Duty**: spell id: vendored 31848, client 1224697
- **Swift Judgement**: spell id: vendored 0, client 1310994
- **Templar's Bulwark**: spell id: vendored 0, client 1311015
- **Iron Creed**: spell id: vendored 0, client 1311034

### Retribution

- **Holy Conduit**: spell id: vendored 0, client 1237268
- **Sanctified Judgement**: spell id: vendored 31876, client 1311074
- **Sacred Arbiter**: spell id: vendored 0, client 1311087
- **Crusade**: spell id: vendored 31866, client 1311083
- **Champion of the Light**: spell id: vendored 31837, client 1311084
- **Instrument of Law**: spell id: vendored 0, client 1311085
- **Twist of Light**: spell id: vendored 0, client 1310735

## Priest (client TraitTree 1114)

Vendored tree order: Discipline, Holy, Shadow Magic. Client tabs left to right map to: Discipline, Holy, Shadow Magic. Client has no tab names in Trait tables (TalentTab is the unchanged Classic table).

Stale duplicate client nodes ignored: Holy Specialization (node 105865, raw tab2 r0 c0)


### Discipline

- **Power in Light**: spell id: vendored 0, client 1309969
- **Twin Disciplines**: spell id: vendored 0, client 1225132
- **Holy Precision**: spell id: vendored 0, client 1309957
- **Soul Warding**: spell id: vendored 0, client 402000
- **Penance**: spell id: vendored 0, client 402174
- **Renewed Hope**: spell id: vendored 0, client 425280
- **Divine Aegis**: spell id: vendored 0, client 431622

### Holy

- **Twilight Focus**: spell id: vendored 0, client 14913
- *client data note* Holy Specialization: duplicate stale node 105865 at tab2 r0 c0 ignored
- **Binding Heal**: spell id: vendored 0, client 401937
- **Litany of Light**: spell id: vendored 0, client 1317006
- **Prayer of Mending**: spell id: vendored 0, client 401859

### Shadow Magic

- **Blackout**: spell id: vendored 15268, client 15326
- **Improved Mind Flay**: spell id: vendored 0, client 1225139
- **Devouring Contagion**: spell id: vendored 0, client 1309950
- **Early Demise**: spell id: vendored 0, client 1310076

## Rogue (client TraitTree 1111)

Vendored tree order: Assassination, Combat, Subtlety. Client tabs left to right map to: Assassination, Combat, Subtlety. Client has no tab names in Trait tables (TalentTab is the unchanged Classic table).


### Assassination

- **Mutilate**: spell id: vendored 0, client 1310707
- **Venom**: spell id: vendored 0, client 1310703

### Combat

- **Puncturing Wounds**: spell id: vendored 0, client 1224716
- **Hack and Slash**: spell id: vendored 0, client 13960
- **Aggression**: prereq: vendored ['hackandslash'], client none
- **Flawless Execution**: name: vendored "Restless Blades", client "Flawless Execution" (matched by position); spell id: vendored 0, client 1310711

### Subtlety

- **Dirty Tricks**: spell id: vendored 0, client 1224782
- **Improved Distract**: spell id: vendored 0, client 14084
- **Quietus**: spell id: vendored 0, client 1310728
- **Cutthroat**: spell id: vendored 0, client 462708
- **Thousand Cuts**: spell id: vendored 0, client 1310721

## Shaman (client TraitTree 1082)

Vendored tree order: Elemental Combat, Enhancement, Restoration. Client tabs left to right map to: Elemental Combat, Enhancement, Restoration. Client has no tab names in Trait tables (TalentTab is the unchanged Classic table).


### Elemental Combat

- **Call of Thunder**: spell id: vendored 16041, client 16120
- **Lightning Overload**: spell id: vendored 30675, client 408438
- **Earthbound**: spell id: vendored 2484, client 1222988
- **Lava Burst**: spell id: vendored 51505, client 408490

### Enhancement

- **Mental Dexterity**: spell id: vendored 51883, client 415140
- **Shamanistic Focus**: spell id: vendored 43338, client 1223030
- **Improved Stormstrike**: spell id: vendored 51521, client 1223031
- **Maelstrom Weapon**: spell id: vendored 51528, client 408498
- **Rage of the Farseer**: spell id: vendored 2825, client 425336

### Restoration

- **Tidal Mastery**: position (row,col 0-based): vendored (0,2), client (3,0)
- **Mindfulness**: spell id: vendored 17106, client 1223033
- **Water Shield**: spell id: vendored 52127, client 408510
- **Totemic Focus**: position (row,col 0-based): vendored (3,0), client (0,2)
- **Riptide**: spell id: vendored 61295, client 408521

## Warlock (client TraitTree 1116)

Vendored tree order: Affliction, Demonology, Destruction. Client tabs left to right map to: Affliction, Demonology, Destruction. Client has no tab names in Trait tables (TalentTab is the unchanged Classic table).


### Affliction

- **Malediction**: spell id: vendored 32477, client 1225177
- **Soul Harvesting**: spell id: vendored 18213, client 437032
- **Improved Drains**: spell id: vendored 17804, client 403511
- **Pandemic**: spell id: vendored 58435, client 427712
- **Malevolence**: spell id: vendored 0, client 1310949
- **Soul Siphon**: spell id: vendored 17864, client 17804
- **Wrack**: name: vendored "Drain Hope", client "Wrack" (matched by position); spell id: vendored 18220, client 1316697

### Demonology

- **Demonic Aegis**: spell id: vendored 30143, client 1235316
- **Demonic Energies**: spell id: vendored 18748, client 1225214
- **Decimation**: spell id: vendored 63156, client 440870
- **Demonic Brand**: spell id: vendored 18821, client 1293695
- **Improved Felhunter**: spell id: vendored 54037, client 1225217
- **Demonic Knowledge**: spell id: vendored 35691, client 412732
- **Demonic Pact**: spell id: vendored 47236, client 425464

### Destruction

- **Molten Skin**: spell id: vendored 63349, client 1225220
- **Conflagrate**: spell id: vendored 17962, client 1293817
- **Pyroclasm**: spell id: vendored 18096, client 18073
- **Bane of Havoc**: spell id: vendored 80240, client 1225228
- **Fire and Brimstone**: spell id: vendored 47266, client 412751
- **Shadow and Flame**: spell id: vendored 30288, client 426316
- **Incinerate**: spell id: vendored 29722, client 412758

## Warrior (client TraitTree 1117)

Vendored tree order: Arms, Fury, Protection. Client tabs left to right map to: Arms, Fury, Protection. Client has no tab names in Trait tables (TalentTab is the unchanged Classic table).


### Arms

- **Improved Tactical Mastery**: spell id: vendored 0, client 12295
- **Spearing Strike**: spell id: vendored 0, client 1310222
- **Bloodthrill**: spell id: vendored 0, client 1289682
- **Weaponmaster**: spell id: vendored 0, client 1290261

### Fury

- **Iron Will**: spell id: vendored 12300, client 12962
- **Blood Craze**: spell id: vendored 0, client 16487
- **Boundless Rage**: spell id: vendored 0, client 1310236
- **Raging Blows**: spell id: vendored 0, client 1310315
- **Precision**: spell id: vendored 20189, client 1225295

### Protection

- **Master of Defense**: spell id: vendored 0, client 1310316
- **Defiance**: spell id: vendored 12303, client 12792
- **Vanguard**: spell id: vendored 0, client 1310317
- **Focused Rage**: position (row,col 0-based): vendored (5,0), client (5,2); spell id: vendored 0, client 29787
- **Bastion**: position (row,col 0-based): vendored (5,2), client (4,3); spell id: vendored 0, client 16538
- ONLY VENDORED **Vitality** (row 4, col 3, 5 ranks)

## Verdict on the seven structural differences (checked 2026-09-17)

All seven were checked against the sim trees and against the independent talentsforever crawl
(`C:/repo/forver/crawl/talentsforever/data.json`, 470 talents — the same count the sim vendors).
**None were applied.** Each structural claim is contradicted by the crawl:

| Claim in this report | Crawl says | Action |
|---|---|---|
| Shaman Tidal Mastery (0,2)->(3,0), Totemic Focus (3,0)->(0,2) | Tidal Mastery r0c2, Totemic Focus r3c0 — the vendored layout | not applied |
| Warrior Focused Rage (5,0)->(5,2), Bastion (5,2)->(4,3) | Bastion at (4,3) collides with Vitality, which the crawl confirms at (4,3) | not applied |
| Rogue Aggression prereq: client none | crawl records `req="Hack and Slash"` | not applied |
| Druid Balance of Nature missing from client | present at r2c3 | not applied |
| Warrior Vitality missing from client | present at r4c3 | not applied |
| Rogue Restless Blades renamed to Flawless Execution | Restless Blades at r3c1 | not applied |
| Warlock Drain Hope renamed to Wrack | Drain Hope at r6c1, req Siphon Life | not applied |
| Hunter Improved Serpent Sting only in client | absent from the crawl | not applied (this report already called it a stale node) |

Two further signals that the Trait-table reading is off rather than the trees:

- Applying the shaman swap makes a shipped Elemental preset **illegal**: it spends 14 points in
  Restoration, and Tidal Mastery at (3,0) needs 15 spent above it. `TestPresetBuildsAreLegal` catches this.
- Applying the warrior swap puts **two talents in one slot** (Bastion and Vitality at (4,3)).

The likely cause is the one this report already flags: several client nodes carry off-grid PosX/PosY
that had to be rescaled, and where two nodes share a definition the higher id was kept. Those two
heuristics can plausibly swap a pair of positions. The spell-id column is a different matter and is
not disputed here — it is unverified either way, and the trees do not rely on it.

**Before acting on a future run of this tool**, confirm any position or prereq change against a second
source. The guards that caught these: `TestTalentTreesMatchTheirProtos` and `TestPresetBuildsAreLegal`.

## Verdict reversed (17 September, later)

The rejection above rested on the talentsforever crawl being a reading of the beta. Its own data export says it
was read off BlizzCon demo footage and Blizzard's slides, so it predates build 1.60.1.69893. The client has no
trait node at all for Restless Blades, Drain Hope, Balance of Nature or Vitality, and Wowhead Forever lists
Flawless Execution (1310711) and Wrack (1316697). Applied: the two renames, both removals, the shaman and warrior
position changes and Aggression's dropped prerequisite. Presets that spent points in a removed talent moved them
(Moonkin into Genesis, Protection warrior into Improved Disarm and Improved Shield Bash), and the Elemental preset
spends its Restoration points in Totemic Focus, which now sits where Tidal Mastery was.

Not applied: the hunter's Improved Serpent Sting. The node carries Classic's spell id 19464, a malformed PosY and
the lowest node id in the tab, next to the newer Improved Stings that took over its effect; it reads as a node the
client never deleted.
