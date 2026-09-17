# Spell data from the Forever beta client (17 Sep 2026)

Beta build `1.60.1.69893` against Classic Era `1.15.9.69722`, both from wago.tools DB2 CSVs. **Nothing here is applied
to the sim yet.**

## Coefficients are in the client

`SpellEffect.EffectBonusCoefficient` is the spell power coefficient and it is populated (11,398 effects). The column
that is zero everywhere is `Coefficient`, which is retail's scaling-class coefficient and unused here; the earlier
checklist note that "the client cannot confirm coefficients" looked at that column. The Era client's
`EffectBonusCoefficient` matches the values this sim already uses, which is what makes the Forever column trustworthy.
`BonusCoefficientFromAP` is set on only 3 effects.

How the sim side was read: every `RegisterSpell` call during `go test ./sim/<class>/...` was logged (spell id,
`BonusCoefficient`, dot `BonusCoefficient`) with a temporary hook, giving 392 spell ids with a coefficient, then
compared per spell id with both clients. A value of 1 in the client means "not set" and is ignored.

Pattern in what Forever changed:
- **No downranking penalty.** Low ranks now carry the full cast-time coefficient (Fireball rank 1: 0.123 -> 0.429).
- **Max-rank coefficients moved both ways.** Lightning Bolt 0.857 -> 0.714, Chain Lightning 0.714 -> 0.571,
  Shadow Word: Pain and Corruption 0.167 -> 0.2 per tick, Mind Flay 0.15 -> 0.167, Devouring Plague 0.063 -> 0.1,
  Blizzard 0.042 -> 0.03, Rain of Fire 0.083 -> 0.03, Searing Totem 0.083 -> 0.017, Moonfire 0.13 / 0.15.
- **Some coefficients are gone** from the effect rows: Arcane Shot, Serpent Sting, Consecration. Those likely moved to
  a different spell or effect and need a look before reading them as zero.

## Base damage was retuned too

`EffectDieSides` no longer exists; a range is `EffectBasePointsF` with `Variance` (low = base x (1 - v/2), high =
base x (1 + v/2)). Every rank moved, low ranks included, on a smooth curve, so this reads as a rebalance rather than
placeholders:

| Spell (max rank) | Era | Forever | Cast / cost |
|---|---|---|---|
| Fireball 25306 | 596-760, coef 1.0 | 425-541, coef 1.0 | unchanged |
| Lightning Bolt 15208 | 419-467, 0.857 | 185-207, 0.714 | 3.0 -> 2.5 sec, 265 -> 220 mana |
| Chain Lightning 10605 | 493-551, 0.714 | 123 +-11%, 0.571 | 2.5 -> 2.0 sec, 605 -> 485 mana |
| Shadow Bolt 25307 | 482-538, 0.857 | 253-283, 0.857 | unchanged |
| Corruption 25311 (per tick) | 136, 0.167 | 73, 0.2 | unchanged |

## What else the client gives

- `SpellCastTimes` via `SpellMisc.CastingTimeIndex`, `SpellPower` (mana cost), `SpellCooldowns`, `SpellDuration`,
  `SpellEffect.EffectAuraPeriod` (tick length): enough to rebuild each spell's cast, cost, cooldown and dot shape.
- `SpellEffect.EffectRealPointsPerLevel`, `SpellLevels`: level scaling for the few spells that use it.
- `SkillLineAbility`: which spells each class learns, to catch new Forever abilities.
- `ItemSparse`, `ItemSet`, `SpellItemEnchantment`: already watched by `tools/data_watch/wago_db2_diff.py`.

## Sim spells whose coefficient differs from the Forever client

| result | spells |
|---|---|
| FOREVER CHANGED IT, sim still has Era | 124 |
| sim matches neither | 16 |
| sim matches Forever | 168 |
| no coefficient in either client | 84 |

### Forever changed it, the sim still has Era's

| spell | name | sim | Era client | Forever client |
|---|---|---|---|---|
| 10 | Blizzard | 0.042 | 0.042 | 0.03 |
| 116 | Frostbolt | 0.163 | 0.163 | 0.407 |
| 133 | Fireball | 0.123 | 0.123 | 0.429 |
| 143 | Fireball | 0.271 | 0.271 | 0.571 |
| 145 | Fireball | 0.5 | 0.5 | 0.714 |
| 172 | Corruption | 0.08 | 0.08 | 0.2 |
| 205 | Frostbolt | 0.269 | 0.269 | 0.489 |
| 348 | Immolate | 0.037, 0.058 | 0.037, 0.058 | 0.13, 0.2 |
| 403 | Lightning Bolt | 0.123 | 0.123 | 0.429 |
| 421 | Chain Lightning | 0.714 | 0.714 | 0.571 |
| 529 | Lightning Bolt | 0.314 | 0.314 | 0.571 |
| 548 | Lightning Bolt | 0.554 | 0.554 | 0.714 |
| 585 | Smite | 0.123 | 0.123 | 0.429 |
| 589 | Shadow Word: Pain | 0.067 | 0.067 | 0.2 |
| 591 | Smite | 0.271 | 0.271 | 0.571 |
| 594 | Shadow Word: Pain | 0.104 | 0.104 | 0.2 |
| 598 | Smite | 0.554 | 0.554 | 0.714 |
| 686 | Shadow Bolt | 0.14 | 0.14 | 0.486 |
| 689 | Drain Life | 0.078 | 0.078 | 0.1 |
| 695 | Shadow Bolt | 0.299 | 0.299 | 0.629 |
| 705 | Shadow Bolt | 0.56 | 0.56 | 0.8 |
| 707 | Immolate | 0.081, 0.125 | 0.081, 0.125 | 0.13, 0.2 |
| 837 | Frostbolt | 0.463 | 0.463 | 0.597 |
| 915 | Lightning Bolt | 0.857 | 0.857 | 0.714 |
| 930 | Chain Lightning | 0.714 | 0.714 | 0.571 |
| 943 | Lightning Bolt | 0.857 | 0.857 | 0.714 |
| 970 | Shadow Word: Pain | 0.154 | 0.154 | 0.2 |
| 980 | Bane of Agony | 0.046 | 0.046 | 0.133 |
| 992 | Shadow Word: Pain | 0.167 | 0.167 | 0.2 |
| 1014 | Bane of Agony | 0.077 | 0.077 | 0.133 |
| 1120 | Drain Soul | 0.063 | 0.063 | 0.1 |
| 1449 | Arcane Explosion | 0.111 | 0.111 | 0.143 |
| 2120 | Flamestrike | 0.017, 0.134 | 0.134, 0.017 | 0.157 |
| 2121 | Flamestrike | 0.02, 0.157 | 0.157, 0.02 | 0.157 |
| 2136 | Fire Blast | 0.204 | 0.204 | 0.429 |
| 2137 | Fire Blast | 0.332 | 0.332 | 0.429 |
| 2767 | Shadow Word: Pain | 0.167 | 0.167 | 0.2 |
| 2860 | Chain Lightning | 0.714 | 0.714 | 0.517 |
| 2944 | Devouring Plague | 0.063 | 0.063 | 0.1 |
| 3044 | Arcane Shot | 0.204 | 0.204 | - |
| 3140 | Fireball | 0.793 | 0.793 | 0.857 |
| 3606 | Attack | 0.052 | 0.052 | 0.017 |
| 5176 | Wrath | 0.123 | 0.123 | 0.429 |
| 5177 | Wrath | 0.231 | 0.231 | 0.486 |
| 5178 | Wrath | 0.443 | 0.443 | 0.571 |
| 5676 | Searing Pain | 0.396 | 0.396 | 0.429 |
| 5740 | Rain of Fire | 0.083 | 0.083 | 0.03 |
| 6041 | Lightning Bolt | 0.857 | 0.857 | 0.714 |
| 6141 | Blizzard | 0.042 | 0.042 | 0.03 |
| 6217 | Bane of Agony | 0.083 | 0.083 | 0.133 |
| 6219 | Rain of Fire | 0.083 | 0.083 | 0.03 |
| 6222 | Corruption | 0.155 | 0.155 | 0.2 |
| 6223 | Corruption | 0.167 | 0.167 | 0.2 |
| 6350 | Attack | 0.083 | 0.083 | 0.017 |
| 6351 | Attack | 0.083 | 0.083 | 0.017 |
| 6352 | Attack | 0.083 | 0.083 | 0.017 |
| 7648 | Corruption | 0.167 | 0.167 | 0.2 |
| 8042 | Earth Shock | 0.154 | 0.154 | 0.386 |
| 8044 | Earth Shock | 0.212 | 0.212 | 0.386 |
| 8045 | Earth Shock | 0.299 | 0.299 | 0.386 |
| 8050 | Flame Shock | 0.063, 0.134 | 0.134, 0.063 | 0.214, 0.1 |
| 8052 | Flame Shock | 0.093, 0.198 | 0.198, 0.093 | 0.214, 0.1 |
| 8092 | Mind Blast | 0.268 | 0.268 | 0.429 |
| 8102 | Mind Blast | 0.364 | 0.364 | 0.429 |
| 8422 | Flamestrike | 0.02, 0.157 | 0.157, 0.02 | 0.157 |
| 8423 | Flamestrike | 0.02, 0.157 | 0.157, 0.02 | 0.157 |
| 8427 | Blizzard | 0.042 | 0.042 | 0.03 |
| 8921 | Moonfire | 0.052, 0.06 | 0.052, 0.06 | 0.13, 0.15 |
| 8924 | Moonfire | 0.081, 0.094 | 0.081, 0.094 | 0.13, 0.15 |
| 8925 | Moonfire | 0.111, 0.128 | 0.111, 0.128 | 0.13, 0.15 |
| 10185 | Blizzard | 0.042 | 0.042 | 0.03 |
| 10186 | Blizzard | 0.042 | 0.042 | 0.03 |
| 10187 | Blizzard | 0.042 | 0.042 | 0.03 |
| 10215 | Flamestrike | 0.02, 0.157 | 0.157, 0.02 | 0.157 |
| 10216 | Flamestrike | 0.02, 0.157 | 0.157, 0.02 | 0.157 |
| 10391 | Lightning Bolt | 0.857 | 0.857 | 0.714 |
| 10392 | Lightning Bolt | 0.857 | 0.857 | 0.714 |
| 10435 | Attack | 0.083 | 0.083 | 0.017 |
| 10436 | Attack | 0.083 | 0.083 | 0.017 |
| 10605 | Chain Lightning | 0.714 | 0.714 | 0.571 |
| 10892 | Shadow Word: Pain | 0.167 | 0.167 | 0.2 |
| 10893 | Shadow Word: Pain | 0.167 | 0.167 | 0.2 |
| 10894 | Shadow Word: Pain | 0.167 | 0.167 | 0.2 |
| 11671 | Corruption | 0.167 | 0.167 | 0.2 |
| 11672 | Corruption | 0.167 | 0.167 | 0.2 |
| 11677 | Rain of Fire | 0.083 | 0.083 | 0.03 |
| 11678 | Rain of Fire | 0.083 | 0.083 | 0.03 |
| 11711 | Bane of Agony | 0.083 | 0.083 | 0.133 |
| 11712 | Bane of Agony | 0.083 | 0.083 | 0.133 |
| 11713 | Bane of Agony | 0.083 | 0.083 | 0.133 |
| 14281 | Arcane Shot | 0.3 | 0.3 | - |
| 14282 | Arcane Shot | 0.429 | 0.429 | - |
| 14283 | Arcane Shot | 0.429 | 0.429 | - |
| 14284 | Arcane Shot | 0.429 | 0.429 | - |
| 14285 | Arcane Shot | 0.429 | 0.429 | - |
| 14286 | Arcane Shot | 0.429 | 0.429 | - |
| 14287 | Arcane Shot | 0.429 | 0.429 | - |
| 14295 | Volley | 0.056 | 0.056 | 0.03 |
| 15207 | Lightning Bolt | 0.857 | 0.857 | 0.714 |
| 15208 | Lightning Bolt | 0.857 | 0.857 | 0.714 |
| 15407 | Mind Flay | 0.15 | 0.15 | 0.167 |
| 17311 | Mind Flay | 0.15 | 0.15 | 0.167 |
| 17312 | Mind Flay | 0.15 | 0.15 | 0.167 |
| 17313 | Mind Flay | 0.15 | 0.15 | 0.167 |
| 17314 | Mind Flay | 0.15 | 0.15 | 0.167 |
| 18807 | Mind Flay | 0.15 | 0.15 | 0.167 |
| 19276 | Devouring Plague | 0.063 | 0.063 | 0.1 |
| 19277 | Devouring Plague | 0.063 | 0.063 | 0.1 |
| 19278 | Devouring Plague | 0.063 | 0.063 | 0.1 |
| 19279 | Devouring Plague | 0.063 | 0.063 | 0.1 |
| 19280 | Devouring Plague | 0.063 | 0.063 | 0.1 |
| 20116 | Consecration | 0.042 | 0.042 | - |
| 20187 | Judgement of Righteousness | 0.144 | 0.144, 0.058 | 0.5, 0.058 |
| 20280 | Judgement of Righteousness | 0.312 | 0.312, 0.125 | 0.5, 0.125 |
| 20281 | Judgement of Righteousness | 0.462 | 0.462, 0.185 | 0.5, 0.185 |
| 20922 | Consecration | 0.042 | 0.042 | - |
| 20923 | Consecration | 0.042 | 0.042 | - |
| 20924 | Consecration | 0.042 | 0.042 | - |
| 25295 | Serpent Sting | 0.2 | 0.2 | - |
| 25311 | Corruption | 0.167 | 0.167 | 0.2 |
| 25742 | Seal of Righteousness | 0.029, 0.032 | 0.029 | 0.1 |
| 26364 | Lightning Shield | 0.147 | 0.147 | 0.267 |
| 26365 | Lightning Shield | 0.227 | 0.227 | 0.267 |
| 26573 | Consecration | 0.042 | 0.042 | - |

### Matches neither client (sim-specific splits such as Pyroblast's dot and Seal of Righteousness' weapon-speed scaling)

| spell | name | sim | Era client | Forever client |
|---|---|---|---|---|
| 603 | Bane of Doom | 1.0 | - | 4.0 |
| 11366 | Pyroblast | 0.15, 1.0 | 0.15 | 0.15 |
| 12505 | Pyroblast | 0.15, 1.0 | 0.15 | 0.15 |
| 12522 | Pyroblast | 0.15, 1.0 | 0.15 | 0.15 |
| 12523 | Pyroblast | 0.15, 1.0 | 0.15 | 0.15 |
| 12524 | Pyroblast | 0.15, 1.0 | 0.15 | 0.15 |
| 12525 | Pyroblast | 0.15, 1.0 | 0.15 | 0.15 |
| 12526 | Pyroblast | 0.15, 1.0 | 0.15 | 0.15 |
| 18809 | Pyroblast | 0.15, 1.0 | 0.15 | 0.15 |
| 25713 | Seal of Righteousness | 0.1, 0.11 | 0.1 | 0.1 |
| 25735 | Seal of Righteousness | 0.1, 0.11 | 0.1 | 0.1 |
| 25736 | Seal of Righteousness | 0.1, 0.11 | 0.1 | 0.1 |
| 25737 | Seal of Righteousness | 0.1, 0.11 | 0.1 | 0.1 |
| 25738 | Seal of Righteousness | 0.1, 0.11 | 0.1 | 0.1 |
| 25739 | Seal of Righteousness | 0.093, 0.102 | 0.093 | 0.1 |
| 25740 | Seal of Righteousness | 0.063, 0.069 | 0.063 | 0.1 |
