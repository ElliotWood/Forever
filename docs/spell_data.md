# Spell Data

Every ranked spell in the game has a table generated from the client database, checked in at
`sim/<class>/spell_data_auto_gen.go`. A spell reads its numbers from that table instead of carrying
hand-transcribed literals.

- [Using a rank](#using-a-rank)
- [The value shapes](#the-value-shapes)
- [Reaching a single effect](#reaching-a-single-effect)
- [Talents](#talents)
- [Worked examples](#worked-examples)
- [Attack power](#attack-power)
- [Regenerating](#regenerating)
- [Traps](#traps)
- [Buffs and debuffs](#buffs-and-debuffs)

## Using a rank

Each class package has exactly one generated global, `spellData`, with a field per spell family:

```go
spellData.Exorcism          // the whole ladder, ranks 1-7
spellData.Fireball          // ranks 1-14
```

Pick the rank your spell registers **by spell ID**:

```go
var exorcismRanks = spellData.Exorcism.BySpellID(27138)
```

That is the identity the sim already uses everywhere - `ActionID`, saved APLs and the icon database all
key on the spell ID - so it cannot drift onto a different rank, and a regeneration that drops the ID
fails loudly instead of quietly substituting another.

The other accessors:

|                    |                                                                                               |
| ------------------ | --------------------------------------------------------------------------------------------- |
| `BySpellID(27138)` | the rank registered under that spell ID. Prefer this.                                         |
| `ByRank(6)`        | the rank numbered 6                                                                           |
| `Ranks(6, 8)`      | a subset, **in the order given**, which is registration order                                 |
| `HighestRank()`    | the highest rank _in the data_, which is not always one the game grants - see [Traps](#traps) |
| `RegisterAll(f)`   | calls `f` once per rank, in declaration order                                                 |

## The value shapes

A rank's value is discriminated by shape, so a variant only carries fields that mean something for it:

```go
shared.SpellDataFlat     {Value, Coef, APCoef}                          // a mana restore, a talent's number
shared.SpellDataRange    {Min, Max, Coef, APCoef}                       // damage or healing the client rolls
shared.SpellDataPeriodic {Tick, TickMax, TickLength, NumberOfTicks, Coef, APCoef} // a tick and its schedule
```

They sit on the four roles a rank can carry, any of which may be nil:

```go
rank.Direct     // Effect = SCHOOL_DAMAGE
rank.Heal       // Effect = HEAL
rank.Periodic   // a periodic aura
rank.Energize   // Effect = ENERGIZE, e.g. Lay on Hands' mana restore
```

Asking what a value is worth on this cast is a single call, because the question means something for
all three shapes - a range rolls between its ends, a flat value and a tick are already the answer. The
method is named for what it produces rather than for how, so a static ability does not read as if it
rolled:

```go
baseDamage := rank.Direct.Damage(sim)     // instead of CalcAndRollDamageRange(sim, min, max)
tickDamage := rank.Periodic.Damage(sim)   // a tick is the answer unless the client rolls it
```

The coefficients are methods, named for the `core.SpellConfig` fields they feed:

```go
BonusCoefficient: rank.Direct.BonusCoefficient(),     // spell power
                  rank.Periodic.BonusCoefficient(),   // same on a tick
                  rank.Direct.APBonusCoefficient(),   // attack power
```

`Range()` gives both ends of a value at once:

```go
low, high := rank.Direct.Range()   // equal for a flat value or a tick
```

The methods assume the role is there. Where it may not be - `Energize` is nil on Lay on Hands rank 1 -
use the package helpers instead, which read a nil value as zero:

```go
shared.SpellDataMin(rank.Energize)     // 0 rather than a panic
shared.SpellDataMax(rank.Direct)
shared.SpellDataCoef(rank.Periodic)
shared.SpellDataAPCoef(rank.Direct)
```

Tick length and count live only on the periodic shape, so assert for them:

```go
p := rank.Periodic.(shared.SpellDataPeriodic)
p.TickLength     // time.Duration, feeds core.DotConfig.TickLength
p.NumberOfTicks  // duration over the tick length, feeds core.DotConfig.NumberOfTicks
```

A rank also carries what the client knows about casting it:

|                               |                                                                                                                                            |
| ----------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------ |
| `Cost`                        | in the units the sim uses - see the rage trap below                                                                                        |
| `CastTime`, `GCD`, `Cooldown` | zero for a channel, whose duration carries it                                                                                              |
| `MinRange`, `MaxRange`        | `core` gates the cast on both; zero means ungated. `MinRange` is the dead zone on a charge, and is nonzero on only 212 spells in the build |
| `MissileSpeed`                | yards per second, which `core` turns into the delay before the damage lands. Zero is an instant hit                                        |
| `SpellSchool`, `DefenseType`  | as `core` names them. The client's school bits are `core`'s bits, so this is the number the DBC states rather than a translation of it     |

`MissileSpeed` is the one to be careful with: giving a spell a speed it did not have delays its damage
and moves goldens, so check the sim is not already modelling it elsewhere. Arcane Missiles is the case
to know - the channel carries no speed because the missile spell does, and that one is not a ranked row.

## Reaching a single effect

The role fields describe one effect each, which is all a castable spell needs. A talent routinely
carries two or three that the sim reads separately, and only one of them can be `Direct`:

```go
irf := spellData.ImprovedRighteousFury.ByRank(3)

irf.Effect(shared.A_ADD_PCT_MODIFIER, 8).Value    //  50  threat bonus
irf.Effect(shared.A_ADD_FLAT_MODIFIER, 12).Value  //  -6  damage taken
irf.Effects[1].Value                              //  -6  the same effect, by index
```

The aura and effect names are generated into `sim/common/shared/spell_data_enums_auto_gen.go`, mirrored
from `tools/database/dbc/enums.go` and holding only the values the tables use, so the two cannot drift.
`Misc` stays a plain int, because what it selects depends on the aura - see
[The Misc value](#the-misc-value).

Name the effect by aura rather than reading `Direct` whenever a spell has more than one. Which effect
lands in `Direct` is the generator's choice, not a promise, so a caller that depends on it breaks
silently the day the ordering changes.

`Effect` panics when nothing matches, and also when **two** effects match: 186 ranked spells carry a
duplicate aura/misc pair, and returning the first is how a caller ends up reading the wrong half of a
talent. Index into `Effects` where the pair cannot tell them apart.

**The value is in the client's units.** A percentage is an integer here - Improved Righteous Fury's
threat bonus reads `16`, not `0.16` - so the `/100` stays at the call site. It is deliberately not
folded into the generator the way the rage `/10` is: whether a value is a percentage depends on the
aura, so a blanket rule would be wrong for some rows and invisible when it was.

## Talents

A talent is read by the points spent in it, not registered at a rank it has, so the ladder has its own
three readers. All of them answer the identity at rank 0 - an untaken talent - where `ByRank` would
panic:

```go
spellData.Moonfury.FractionAt(rank)        // 0.10 at 5/5 - the client's 10, over 100
spellData.NaturesReach.ValueAt(rank)       // 20 at 2/2  - the client's number as it stands
spellData.LivingSpirit.MultiplierAt(rank)  // 1.15 at 5/5 - 1 + the fraction
```

That replaces the `<literal> * float64(x.Talents.Y)` idiom, and with it the `if rank > 0` guard the
caller would otherwise need.

**`MultiplierAt` takes its sign from the data.** Improved Righteous Fury states its damage reduction as
-2 / -4 / -6, so rank 3 gives 0.94 and nobody writes the minus. Where the sim's parameter runs the other
way - `AddReducedCritTakenPercent` wants a positive amount for a reduction the client states negative -
negate at the call site so the disagreement is visible.

**Ladders are not always the per-point literal times the rank.** Most are: of 125 percent talents, 122
scale linearly, so `0.02 * rank` was already right and the table only adds provenance. The ones that do
not are the reason to read it - Improved Righteous Fury is 16 / 33 / 50, not 16 / 32 / 48, and shaman
Elemental Weapons is 7 / 14 / 20, not 7 / 14 / 21.

### Picking the effect

A talent with one effect per rank needs nothing further. One with several does, and `ValueAt` panics
rather than guess:

```go
spellData.ImprovedRighteousFury.
    Effect(shared.A_ADD_PCT_MODIFIER, shared.SPELLMOD_ALL_EFFECTS).MultiplierAt(rank)   // 1.50 threat
spellData.ImprovedRighteousFury.
    Effect(shared.A_ADD_FLAT_MODIFIER, shared.SPELLMOD_EFFECT2).MultiplierAt(rank)      // 0.94 taken
```

**Do not pick the effect by which one matches the number.** Survival of the Fittest states +1/2/3% to
all stats and -1/-2/-3% crit taken; both ladders fit, and an automatic pass attached the stat effect to
the crit-taken call site. What the call site does decides it, and the mod's `Kind` usually says:
Improved Moonfire's two mods are `SpellMod_DamageDone_Flat` and `SpellMod_BonusCrit_Percent`, so one
takes `SPELLMOD_DAMAGE` and the other `SPELLMOD_CRITICAL_CHANCE`.

Where a talent modifies damage and its DoT with the same ladder, the sim has one mod against the
client's two. Either aura reads the same number; `SPELLMOD_DAMAGE` is the convention here.

### Proc chances

`SpellAuraOptions.ProcChance` is a separate source from the effects, and `ProcChanceAt` reads it as the
fraction a `ProcTrigger` takes:

```go
ProcChance: spellData.SealFate.ProcChanceAt(rogue.Talents.SealFate)   // 0.20 at 1/5, 1.00 at 5/5
```

**A 100 does not always mean a 100% roll.** Flurry and Enrage read 100 because they fire on their own
condition - a crit - rather than on a chance, and the number the sim wants for those is somewhere else
entirely. Check what the talent actually does before wiring it.

### The high end of an effect

`SpellDataEffect.Value` is the low end. An aura with one die side and a fractional base has two ends a
whole number apart, and the game shows the higher: Seal of the Crusader rank 1 states 39.2 attack power
and buffs for 41. `ValueMax` holds it, and is zero on the 82% of effects where the two agree, so read
`High()` rather than `ValueMax` - rank 4's base is whole, so it has no `ValueMax` and its answer is
`Value`.

```go
spellData.SealOfTheCrusader.ByRank(rank).Effects[0].High()   // 41 at rank 1, 183 at rank 4
```

### The Misc value

`Misc` says what an effect applies to, and what it means depends on the aura: a modified spell property
under `A_ADD_PCT_MODIFIER` and `A_ADD_FLAT_MODIFIER`, a stat under `A_MOD_TOTAL_STAT_PERCENTAGE`, a
school mask under `A_MOD_DAMAGE_DONE`. There is no single enum for it, so it stays an int.

For the two modifier auras the `SPELLMOD_*` constants name it. Those are hand-written in
`sim/common/shared/spell_data_talents.go`, because the client ships no name list - each carries the
talents it was read off. All 23 were then checked against [TrinityCore's `SpellModOp`][tc] (3.3.5) and
[cmangos-tbc's][cm] (2.4.3), which agree with every value, and with every name except 24 and 27 where
cmangos says `SPELL_BONUS_DAMAGE` and `MULTIPLE_VALUE`. No modifier effect in the tables uses a value
outside those 23; the ones the cores name and TBC does not use are 13, 17, 20, 21 and 26.

[tc]: https://github.com/TrinityCore/TrinityCore/blob/3.3.5/src/server/game/Spells/SpellDefines.h
[cm]: https://github.com/cmangos/mangos-tbc/blob/master/src/game/Spells/SpellDefines.h

## Worked examples

### Direct damage

```go
var exorcismRanks = spellData.Exorcism.BySpellID(27138)

func (paladin *Paladin) registerExorcism() {
	paladin.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: exorcismRanks.SpellID},
		Rank:             exorcismRanks.Rank,
		ManaCost:         core.ManaCostOptions{FlatCost: exorcismRanks.Cost},
		BonusCoefficient: exorcismRanks.Direct.BonusCoefficient(),
		MaxRange:         exorcismRanks.MaxRange,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{GCD: exorcismRanks.GCD, CastTime: exorcismRanks.CastTime},
			CD:          core.Cooldown{Timer: paladin.getExorcismTimer(), Duration: exorcismRanks.Cooldown},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, exorcismRanks.Direct.Damage(sim), spell.OutcomeMagicHitAndCrit)
		},
	})
}
```

### A damage-over-time effect

The tick and its schedule both come from the table, so `NumberOfTicks` and `TickLength` stop being
hand-written:

```go
var swpRanks = spellData.ShadowWordPain.BySpellID(25368)

func (priest *Priest) registerShadowWordPain() {
	tick := swpRanks.Periodic.(shared.SpellDataPeriodic)

	priest.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: swpRanks.SpellID},
		ManaCost: core.ManaCostOptions{FlatCost: swpRanks.Cost},

		Dot: core.DotConfig{
			Aura:             core.Aura{Label: "ShadowWordPain-" + swpRanks.GetRankLabel()},
			NumberOfTicks:    tick.NumberOfTicks,
			TickLength:       tick.TickLength,
			BonusCoefficient: tick.Coef,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.Snapshot(target, tick.Tick)
			},
		},
	})
}
```

### A heal, and a mana restore

```go
holyLight := spellData.HolyLight.BySpellID(27136)
low, high := holyLight.Heal.Range()          // both ends in one call

layOnHands := spellData.LayOnHands.BySpellID(27154)
mana := shared.SpellDataMin(layOnHands.Energize)   // 900; rank 1 has no Energize at all, and reads 0
```

### Registering several ranks

Downranking registers more than one, and the spec chooses which:

```go
// Starfire's ladder runs 1-8; the sim registers only these two.
spellData.Starfire.Ranks(6, 8).RegisterAll(druid.registerStarfireSpell)

func (druid *Druid) registerStarfireSpell(rank shared.SpellData) {
	druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
		ActionID: core.ActionID{SpellID: rank.SpellID},
		Rank:     rank.Rank,
		ManaCost: core.ManaCostOptions{FlatCost: rank.Cost},
		Cast:     core.CastConfig{DefaultCast: core.Cast{GCD: rank.GCD, CastTime: rank.CastTime}},
		// ...
	})
}
```

Each rank is its own registered spell with its own ActionID, which is what lets an APL name a downrank.

## Attack power

Attack power scaling is **not** in the client data - one effect in 38357 carries a nonzero
`BonusCoefficientFromAP` - so melee coefficients live in server script and have to be supplied by hand:

One coefficient for the whole ladder:

```go
// Rupture's ranks are all periodic, so the coefficient goes on the tick.
var ruptureRanks = shared.WithSpellDataPeriodicAPCoef(spellData.Rupture, 0.18)
```

Or one per rank, the way the spell power coefficient already varies because each row carries its own:

```go
var ruptureRanks = shared.WithSpellDataPeriodicAPCoefs(spellData.Rupture, map[int32]float64{
	1: 0.04, 2: 0.06, 3: 0.08, 4: 0.10, 5: 0.12, 6: 0.15, 7: 0.18,
})
```

**Every rank in the table has to be named.** Leave one out and it panics rather than scaling that rank
off nothing, and naming a rank the ladder does not have panics too - so a ladder that gains a rank in a
later client build fails loudly instead of quietly mis-scaling.

|                                        |                                |
| -------------------------------------- | ------------------------------ |
| `WithSpellDataAPCoef(t, c)`            | one coefficient, on `Direct`   |
| `WithSpellDataPeriodicAPCoef(t, c)`    | one coefficient, on `Periodic` |
| `WithSpellDataAPCoefs(t, map)`         | per rank, on `Direct`          |
| `WithSpellDataPeriodicAPCoefs(t, map)` | per rank, on `Periodic`        |

All four return a copy, so the generated table keeps what the database said. All four panic if the role
is nil on any rank - check the generated table first, `spellData.Mangle` is the _learn-spell_ entry
(`Effect = 36`) and carries no value at all - and if the table already carries a coefficient, because a
value that appears upstream should be noticed, not silently shadowed.

## A row that is more than one spell

A seal is three spells for one rank - the aura, the proc it triggers and the judgement - so it keeps
its own row type rather than becoming a `SpellData`, and `SpellDataTableOf` is generic for exactly
that. What the client states is read from the tables; what it does not is passed in, and the shorter
name goes to the common case so a family that diverges reads differently from one that does not.

```go
// everything the client states, judgement damage included
sealOf(spellData.SealOfCommand, spellData.JudgementOfCommand, rank, proc{...})

// for the families whose judgement damage has to be supplied by hand
sealWithJudgement(spellData.SealOfLight, spellData.JudgementOfLight, rank, proc{...}, judge{...})
```

This is the pattern for any composite that follows - totems, poisons. Two things it taught: put the
reason for each literal on its own line rather than in a block at the top, and do not trust a golden
to verify it. The paladin goldens carry no seal spell ID at all, so the port was checked by dumping
every constructed row against the literals it replaced.

## Regenerating

```
go run ./tools/database/gen_spelldata
```

Reads `tools/database/wowsims.db` and rewrites every `sim/<class>/spell_data_auto_gen.go`. It is its
own binary rather than a mode of `gen_db` on purpose: `gen_db` imports the sim, and the sim reads these
tables, so a stale generated file would stop the generator that fixes it from compiling.

Nothing lists which spells to generate. The generator walks `dbc.Classes`, takes each class's own skill
lines and emits a table for every spell in them whose subtext reads `Rank N`. A new family appears on
its own; if one you expect is missing, look at the `// Not generated:` comment at the head of the class
file, which names every family that could not be resolved and why.

`go test ./tools/database/ -run GeneratedRankTables` re-derives amounts and coefficients from the
database and compares them to the committed tables. It covers the 23 families listed in
`spelldata_regen_test.go` - 514 of the 3327 rows - so it is not a substitute for regenerating and
checking the diff is empty, which is the only check that covers every row. It skips when `wowsims.db`
is absent.

## Traps

**The data contains ranks the game never grants.** Fireball 38692 and Frostbolt 38697 are rank 14
entries at level 70 whose `SkillLineAbility` rows and spell attributes are byte-for-byte
indistinguishable from the real rank 13s, so the generator cannot filter them. Reaching for
`HighestRank()` on those families silently casts a spell that does not exist - it cost 1.2% DPS when it
happened during the mage port. Use `BySpellID`.

**`HighestRank()` is not "the rank my spec casts".** It is the largest rank number present, which is
also not the last element: Flamestrike is declared rank 7 then rank 6.

**A rank can carry nothing in a role.** Lay on Hands rank 1 restores no mana where ranks 2-4 do, so
`rank.Energize` is nil there. The `SpellData*` helpers read nil as zero; a direct field access does not.

**Never delete a generated file before regenerating it.** The class package stops compiling, and
`gen_db` - which imports the sim - then cannot build either. Regenerate over the top, or
`git checkout HEAD -- <path>` to get back.

**A melee ability states its bonus as weapon damage, not school damage.** Sinister Strike's +98 is
`Effect = 121` (normalised weapon damage); the generator reads 17, 58 and 121 alongside school damage,
so those land in `Direct`. `Effect = 31` (weapon percent damage) is deliberately excluded - it is a
multiplier on the swing, not an amount a rank can carry.

**Rage costs are divided by ten on the way in.** The client stores rage on a 0-1000 bar, so Heroic
Strike's cost reads 150 where the player sees 15. `Cost` is always in the units the sim uses; mana,
energy and focus need no conversion, and only rage does. An `Energize` effect that restores rage would
still be in tenths - nothing generated today does.

**A new client table needs a settings line.** `SpellCastTimes` was missing from
`generator-settings.json` and cast times read zero until it was added and `make db` re-run. Adding a
table is one line; the extractor needs no code.

## Buffs and debuffs

Every raid, party, individual and enemy-debuff proto field has one row in
`tools/database/buffmanifest`, and `tools/database/gen_spelldata` turns those rows into
`sim/core/buffs_auto_gen.go`, `sim/core/debuffs_auto_gen.go` and
`ui/features/settings/model/buffs_debuffs_auto_gen.ts`. `tools/gen_buffs_proto` renders
`proto/buffs.proto` from the same rows. The manifest owns the field numbers; the client database owns
the values.

### The manifest row

`BuffSpec` in `tools/database/buffmanifest/manifest.go`:

| Field                      | What it is                                                                                                                                                                                                                                                        |
| -------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `Field`, `Number`, `Scope` | the proto field, its number and the message it lives on. Nothing else may take a number, and a retired one goes into `Retired`                                                                                                                                    |
| `Proto`                    | `ProtoBool`, `ProtoTristate`, `ProtoInt32`, `ProtoDouble` or `ProtoEnumDrums`. Declared, not derived, so the emitter runs while the compiled protos are stale; the resolver checks it against the live trait tree                                                 |
| `Kind`                     | what the generator emits, below                                                                                                                                                                                                                                   |
| `Go`                       | the identifier stem: `BattleShout` gives `BattleShoutAura`, `BattleShoutValue`, `BattleShoutDuration`, `BattleShoutCategory`                                                                                                                                      |
| `Name`                     | the castable family's `SpellName.Name_lang`, resolved through `SkillLineAbility` for the owning class                                                                                                                                                             |
| `AuraName`                 | the aura family, when the cast is a summon or a dummy: the totem cast names `Strength of Earth Totem` and the aura `Strength of Earth`                                                                                                                            |
| `Anchor`                   | an explicit spell id, for a row with no castable name at all (the Atiesh staves)                                                                                                                                                                                  |
| `Owner`                    | the class that casts it, which narrows the `SkillLineAbility` lookup and marks the row "(External)" on that class's settings tab                                                                                                                                  |
| `Talent`                   | the improving talent family, its effect index, and whether it scales the value, the duration or adds a stat. Only a `ProtoTristate` row may state one                                                                                                             |
| `Category`                 | the exclusive-effect category the aura bids in, `""` for none                                                                                                                                                                                                     |
| `SharedCategory`           | a second category the aura joins without an effect of its own, which is how the paladin auras exclude each other across schools. Applied to the player's copy only, and declared once in the generated file as `<Name>Category`                                   |
| `SingleAura`               | the category holds one aura at a time, so the loser is deactivated rather than outbid                                                                                                                                                                             |
| `Driver`                   | the apply block hands the row to `drive<Go>` instead of activating the aura outright                                                                                                                                                                              |
| `Pet`                      | `PetNormal`, `PetStrip`, `PetInheritOwnerAura`, `PetCapAtRegular` or `PetStripWhenSummonedLate`                                                                                                                                                                   |
| `StatOverride`             | the sim stats the value lands on, for an aura the client states without naming one: `A_MOD_CRIT_PCT` carries no school, so Leader of the Pack says `PhysicalCritPercent` and Moonkin Aura `SpellCritPercent`. A row that states one may have only one aura effect |
| `Stats`                    | the UI relevance tags a spec's `epStats` and `displayStats` are matched against                                                                                                                                                                                   |
| `ImpAction`                | the icon the improved state of a tristate input shows, when it is an item rather than the talent                                                                                                                                                                  |
| `Label`                    | a UI label override; the client's name is the default                                                                                                                                                                                                             |
| `Notes`                    | why a `KindManual`, `KindAbsent` or `KindFlag` row is one. Required for those three                                                                                                                                                                               |

### The kinds

`KindStatFlat`, `KindStatPct` and `KindResistance` are stat buffs; `KindPseudoMult` moves a
pseudo-stat; `KindDamageShield` is a retaliation proc; `KindProc` and `KindExternalCD` need a driver
for the trigger or the cooldown; `KindItemCount` takes a count and applies its amounts per item;
`KindDebuffStat`, `KindDebuffStacking`, `KindDebuffDamageTaken`, `KindDebuffAtkSpeed` and
`KindDebuffUptime` are the debuff shapes. `KindManual` is a row the sim models by hand,
`KindFlag`/`KindEnum` are sim toggles rather than buffs, and `KindAbsent` is a field the Forever
client describes no spell for. The last three resolve to a commented shell naming the reason.

### Resolving a row

An `Anchor` is taken as it stands. Otherwise `Name` is looked up in `SkillLineAbility` for the owner's
class, the highest rank is taken, and a non-monotonic ladder is warned about and recorded above the
constructor. `AuraName` then hops from the cast to the aura family member with the same rank subtext.
A spell an engraving rune grants is refused: those are class runes, not raid buffs. A `Talent` is
matched in the owner's live trait tree and priced per point into a curve; declaring one for a talent
that is not in the tree is an error, and not declaring one where a value-scaling talent exists is a
warning.

### What stays hand-written

The generator emits a real constructor only when `sim/core` declares none: it parses the non-generated
files and treats a row as hand-written if `<Go>Aura` is declared there **or** if an apply block still
reads the row's proto field. Either way the row renders as a commented shell, so migrating one means
deleting both the constructor and the apply-block branch and regenerating.

### Drivers

A row the generator cannot express outright states `Driver: true`, or is a kind that always needs one,
and the apply block calls `drive<Go>` instead. The contract is in `sim/core/buffs_manual.go`: a buff
row's driver takes the `*Character` and the whole scope message, a debuff row's takes the `*Unit`, the
debuffs and the raid. Handing over the message rather than the one field is what lets Grace of Air read
`party.TotemTwisting`. The driver builds the generated aura with `<Go>Aura(...)` and adds what the
client does not state: a proc trigger, a cooldown, a delay, a regen. Declaring the function is what
makes the row compile.

### Regenerating

```
go run ./tools/database/gen_spelldata
```

writes all three generated files. `tools/database` imports `sim/core`, so the generator cannot run
while the tree it generates into does not compile. A change that breaks it - retiring a proto field,
deleting a hand-written constructor the generated file still references - needs the generated file cut
by hand first: delete the lines that name the field or symbol that is going away, `go build ./...`,
then run the generator, which writes the whole file back. Running it twice and seeing an empty
`git status` is the check that it converged.

### The guard tests

`go test ./tools/database/... ./tools/gen_buffs_proto/...` needs no client database and runs in CI.
Without one, `TestGeneratedBuffFiles`, `TestResolvedBuffInvariants`, `TestGeneratedBuffsDebuffsTS`
and `TestGeneratedRankTablesMatchTheDatabase` skip; the other 27 run.

| Test                                                                                                              | What it holds                                                                                                                   |
| ----------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------- |
| `TestUniqueScopeField`, `TestUniqueScopeNumber`, `TestUniqueGoStem`                                               | no two rows collide                                                                                                             |
| `TestNoLiveRowTakesARetiredNumberOrName`                                                                          | neither a field number nor a field name api version 17 gave up is ever handed out again                                         |
| `TestProtoTypeMatchesKind`, `TestTalentImpliesTristate`, `TestShellRowsHaveNotes`, `TestResolvableRowsHaveAnchor` | the schema rules above                                                                                                          |
| `TestFieldNaming`, `TestFieldNamesRoundTrip`                                                                      | `GoField()` and `TSField()` reproduce protoc's and protobuf-ts's camel case                                                     |
| `TestRenderMatchesCommittedFile`                                                                                  | `proto/buffs.proto` is what the manifest renders                                                                                |
| `TestRenderReservesEveryRetiredNumber`, `TestRenderNextIndex`                                                     | both `reserved` lines per message — numbers and names — and the next free number                                                |
| `TestRetypedFieldsAreBool`, `TestRetypedFieldsMatchTheMigration`                                                  | the 25 fields api version 17 retyped are bool, and `ui/sim/proto/buff_field_migration.ts` names the same 25                     |
| `TestRetiredFieldsMatchTheMigration`                                                                              | `buff_field_migration.ts` drops the same 33 retired names the manifest retires, scope by scope                                  |
| `TestRenderedBuffFilesMatchTheFixtures`, `TestRenderedBuffFilesCompile`                                           | synthetic rows render to the committed fixtures, and those fixtures compile against the real `sim/core` through a build overlay |
| `TestRenderBuffsDebuffsTS*`                                                                                       | the settings inputs each proto type and kind renders                                                                            |
| `TestGeneratedBuffFiles`, `TestGeneratedBuffsDebuffsTS`                                                           | with a database, the committed files are byte-for-byte what the generator emits                                                 |
| `TestResolvedBuffInvariants`                                                                                      | with a database, the pinned talent curves, categories and stat amounts                                                          |

Rewrite the fixtures with `UPDATE_BUFF_FIXTURES=1 go test ./tools/database/`.

### Traps

**A retired field number is gone for good.** `buffmanifest.Retired` lists the 33 api version 17 gave
up; a share link or a saved setting from before the bump still carries them on the wire, and a new row
reusing one would read the old value. The emitter turns them into `reserved` lines and counts the next
free index past them.

**A proto change here is a UI migration too.** `ui/sim/proto/buff_field_migration.ts` runs on the raw
JSON before `fromJson`, because the parser throws on an enum name in a bool field and on an unknown key
long before the version converters run. A retyped field goes in its list, a retired one in the other,
and every `fromJson` of a settings envelope passes `ignoreUnknownFields`.

**A `PetInheritOwnerAura` row must name the aura the pet looks for.** The pet finds it on its owner by
label, so the row needs a `Label`, a `Name` or a `Category` even when the client describes the buff not
at all.

**A talent curve only scales the row's first stat.** A row whose talent improves a second amount would
need the generator extended; nothing in the manifest does today.

**A party or raid flag means an external caster provides the buff.** The generated apply block builds
the row's `isPlayer=false` copy whenever the proto field is set, so a class port that registers its own
`isPlayer=true` copy has two copies on the character. It either stops setting the flag in
`AddPartyBuffs`/`AddRaidBuffs`, or the row states a `Category` with `SingleAura` and the two copies bid
against each other, so the character sheet shows the buff once. The higher bid deactivates the other
copy; on a tie the incumbent keeps the category when its remaining duration is the longer one, which
is why a druid casting its own Thorns is turned away while the raid's permanent copy is up - both deal
the same 22, so the character strikes back for the same either way. Battle Shout is the worked example:
both copies are worth the same and neither is permanent, and
`TestPlayerBattleShoutTakesTheCategoryOnATie` holds the player's own to the tie. The rows this decides
are `thorns`, `leader_of_the_pack`, `moonkin_aura` and `trueshot_aura`: `sim/druid/druid.go`,
`sim/druid/feralcat` and `sim/druid/feralbear` raise the party's Leader of the Pack or Moonkin Aura
from a talent, `sim/hunter/hunter.go` raises Trueshot Aura, and the druid's own Thorns waits on the
druid port. `thorns` and `battle_shout` carry the category; the other three do not, so their ports
have to choose.
