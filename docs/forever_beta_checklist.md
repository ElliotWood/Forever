# Forever beta re-verification checklist

Every number in the Forever ruleset that was read off a BlizzCon 2026 demo tooltip rather than game data, in one place, so the pass against the beta client is a checklist rather than an archaeology exercise.

The demo mostly showed rank 1 of each talent. Where a talent has more ranks than the demo displayed, the implementation assumes a scaling (linear unless the tooltip said otherwise) and says so in a `TODO` beside the number. This file lists those `TODO`s, gathered from `grep -rn TODO sim/` and grouped by class. When a value is confirmed, fix the number if it moved, delete the comment, and regenerate the affected `.results` files with `make test && make update-tests`.

## How to run the pass

1. Re-export the talent trees from the beta client (the community calculators rebuild from it) and diff against `ui/core/talents/trees/*.json`: talent set, grid positions, rank counts and prerequisite arrows. `go test ./sim/ -run TestTalentTreesMatchTheirProtos` then pins the trees, protos and `TalentTreeSizes` together, and `-run TestPresetBuildsAreLegal` checks every shipped build still fits.
2. Work through the class sections below against the beta tooltips.
3. Re-check the racials against the beta spellbook, in particular the two figures still missing: the cooldown and cost reduction of the gnome's Eureka!, and whether any racial cooldown differs from the three minutes assumed where none was published.
4. World buffs do not work inside Forever raids (reported 13 September from the demo; the sim ignores them under the Forever ruleset and hides the picker). Confirm on the beta client, and confirm whether the campsite buffs that replace them have combat numbers.
5. Re-run the DPS sweep across every spec and compare with the numbers recorded in the pull request history; anything that moves more than its change explains is worth a second look.

## Druid (16)

- `sim/druid/berserk.go:16` — The tooltip didn't show a cooldown, the 3 minutes are taken from the Classic Berserk.
- `sim/druid/demoralizing_roar.go:24` — assumed baseline, beta will confirm - Feral Aggression is gone from the tree, so the attack power reduction is taken at full strength like the warrior's Demoralizing Shout.
- `sim/druid/druid.go:122` — Improved Mark of the Wild is gone from the tree and is assumed to be baseline, beta will confirm.
- `sim/druid/faerie_fire.go:26` — The feral version's talent is gone from the tree and is assumed to be baseline, beta will confirm.
- `sim/druid/forms.go:223` — Beta will show whether the carryover share and the out of form regen both scale per rank, linear scaling is assumed here.
- `sim/druid/lacerate.go:16` — assumed from Season of Discovery, beta will confirm the cost, the damage and the threat. Shredding Attacks names Lacerate, so the bear has it, but no tooltip for it has been seen.
- `sim/druid/mangle.go:17` — Only the tooltip was seen, the Energy cost is taken from the Classic Mangle (Cat).
- `sim/druid/mangle.go:71` — Only the tooltip was seen, the Rage cost, the 6 sec cooldown and the threat are taken from the Classic Mangle (Bear).
- `sim/druid/talents.go:175` — Only rank 1 was seen, the damage bonus is assumed to scale linearly. Beta will confirm.
- `sim/druid/talents.go:301` — Only rank 1 was seen, the cast time reduction is assumed to scale linearly. Beta will confirm.
- `sim/druid/talents.go:388` — Only rank 1 was seen, both amounts are assumed to scale linearly. Beta will confirm.
- `sim/druid/talents.go:449` — Only rank 1 was seen, the dodge chance is assumed to scale linearly. Beta will confirm.
- `sim/druid/talents.go:452` — Every rank reads the same 20% chance for 5 Rage, so the proc doesn't grow past rank 1. Beta will confirm.
- `sim/druid/talents.go:515` — Only rank 1 was seen, the damage bonus is assumed to scale linearly. Beta will confirm.
- `sim/druid/tigers_fury.go:45` — Tiger's Fury is given Wrath's shape (no Energy cost, 30 sec cooldown) because King of the Jungle is Wrath's talent word for word and Classic's costless-to-spam Tiger's Fury would let it mint Energy; the +40 damage is Classic's rank 4. Beta will confirm the cooldown, the cost and the damage bonus.
- `sim/druid/wrath.go:53` — Only rank 1 was seen, the mana cost reduction is assumed to scale linearly. Beta will confirm.

## Hunter (6)

- `sim/hunter/aimed_shot.go:74` — assumed baseline, beta will confirm
- `sim/hunter/aspects.go:52` — only rank 1 was observed, the proc chance is assumed to scale per rank.
- `sim/hunter/rapid_fire.go:15` — only rank 1 of Rapid Killing was observed, the reduction is assumed to scale per rank.
- `sim/hunter/serpent_sting.go:43` — only rank 1 of Improved Stings was observed, the damage bonus is assumed to scale per rank.
- `sim/hunter/talents.go:299` — only rank 1 was observed, the cost reduction and the proc chance are assumed to scale per rank.
- `sim/hunter/talents.go:352` — only rank 1 was observed, the regeneration is assumed to scale per rank.

## Mage (2)

- `sim/mage/fire_blast.go:39` — only rank 1 of Wake of Fire was shown, so mage.json copies its 1 sec into rank 2.
- `sim/mage/talents.go:655` — both ranks read 15% on the demo tooltip, beta will confirm whether rank 2 is higher.

## Paladin (18)

- `sim/paladin/consecration.go:10` — assumed baseline, beta will confirm - Consecration is no longer a talent and the Forever tree builds on top of it through Consecrated Ground and Holy Conduit.
- `sim/paladin/hammer_of_wrath.go:29` — Both ranks of Instrument of Law read 0.5 sec in the tooltip data.
- `sim/paladin/holy_strike.go:15` — assumed baseline, beta will confirm - the 110% weapon damage and the 6 second cooldown are both guesses.
- `sim/paladin/holy_strike.go:20` — Only rank 1 of Improved Holy Strike was seen, the second second of cooldown is assumed to scale linearly.
- `sim/paladin/holy_strike.go:32` — Every rank of Iron Creed reads the same 5% threat, the rest are assumed to scale linearly.
- `sim/paladin/holy_strike.go:82` — Every rank of Iron Creed reads the same 2% for 6 sec, the rest are assumed to scale linearly.
- `sim/paladin/sotc.go:34` — assumed baseline, beta will confirm - Improved Seal of the Crusader is gone from the tree and the raid reads the improved Judgement of the Crusader through Debuffs either way.
- `sim/paladin/swift_judgement.go:11` — assumed baseline, beta will confirm - the tooltip carries no cooldown, so it is given a minute, long enough that it buys one extra Judgement rather than a second rotation.
- `sim/paladin/talents.go:18` — Only rank 1 of Divine Precision was seen, ranks 2 and 3 are extrapolated from it.
- `sim/paladin/talents.go:45` — Only rank 1 of Champion of the Light was seen, and the extrapolated ranks 2 and 3 are a large chunk of a Forever paladin's spell power.
- `sim/paladin/talents.go:87` — Every rank of Redoubt reads the same 10% chance for 6% block, so ranks 2-5 do nothing.
- `sim/paladin/talents.go:168` — The mana return reads 33% for 6% of maximum mana at every rank.
- `sim/paladin/talents.go:228` — Every rank reads 1% per stack up to 5 stacks, so ranks 2 and 3 do nothing.
- `sim/paladin/talents.go:261` — The self buff reads 1% at every rank. The 42 attack power the target loses is not modelled, nothing in the sim reads an enemy's attack power.
- `sim/paladin/talents.go:298` — The tooltip caps the bonus at the first 4 or 8 enemies to enter the Consecration, which is not modelled here - everything standing in it gets the bonus.
- `sim/paladin/talents.go:334` — Both ranks read 10% in the tooltip data.
- `sim/paladin/templars_bulwark.go:11` — assumed baseline, beta will confirm - the tooltip carries no cooldown, so it shares the 5 minutes of the two Forbearance abilities Sacred Duty shortens alongside it.
- `sim/paladin/templars_bulwark.go:36` — Both ranks of Sacred Duty read 30 sec, the second is assumed to scale linearly.

## Priest (8)

- `sim/priest/devouring_plague.go:53` — only rank 1 of Devouring Contagion was shown, beta will confirm the rank 2 value
- `sim/priest/holy_nova.go:13` — beta will confirm the higher ranks and the mana cost.
- `sim/priest/mind_flay.go:78` — only rank 1 of Improved Mind Flay was shown, beta will confirm the rank 2 value
- `sim/priest/penance.go:18` — beta will confirm the cost, the cooldown and the level 60 damage.
- `sim/priest/priest.go:81` — beta will confirm whether they were made baseline or removed outright.
- `sim/priest/talents.go:40` — only rank 1 was shown, beta will confirm that the two halves scale at 5% and 1% per point
- `sim/priest/talents.go:57` — only rank 1 was shown, beta will confirm the 2% per point
- `sim/priest/talents.go:398` — beta will confirm the 50%.

## Rogue (9)

- `sim/rogue/backstab.go:28` — Only rank 1 of Puncturing Wounds was seen, the extra combo point chance is assumed to scale linearly. Beta will confirm.
- `sim/rogue/expose_armor.go:11` — assumed baseline, beta will confirm. The raid reads this debuff, so the Classic 2/2 armor value is treated as baseline rather than deleted.
- `sim/rogue/expose_armor.go:36` — Only rank 1 was seen, the Energy discount is assumed to scale linearly while the refund and the 5 combo point trigger stay put. Beta will confirm.
- `sim/rogue/hack_and_slash.go:14` — Only rank 1 was seen, all three effects are assumed to scale linearly. Beta will confirm.
- `sim/rogue/mutilate.go:43` — The tooltip showed no Energy cost, the 60 is taken from the Classic Mutilate.
- `sim/rogue/mutilate.go:55` — The tooltip didn't repeat the Classic dagger requirement, it's assumed to still apply.
- `sim/rogue/talents.go:273` — Only rank 1 was seen, the proc chance is assumed to scale linearly. Beta will confirm.
- `sim/rogue/talents.go:350` — Only rank 1 was seen and the damage bonus is assumed to scale linearly. The tooltip data extrapolates the health threshold along with it, which it cannot be, so rank 1's 35% is used for every rank. Beta will confirm.
- `sim/rogue/venom.go:47` — The tooltip showed no Energy cost, the 25 matches the other Rogue finishers.

## Shaman (13)

- `sim/shaman/air_totems.go:50` — The sim won't respect the value of a totem dropped via the APL. It uses hard-coded values from buffs.go bonusDamage := WindfuryTotemBonusDamage[rank]
- `sim/shaman/lava_burst.go:11` — Only the damage range and the Flame Shock bonus were on the tooltip. The cast time, cooldown, mana cost and coefficient are taken from the spell of the same name, beta will confirm them.
- `sim/shaman/lightning_overload.go:25` — Only rank 1 was seen, beta will confirm that the ranks go up in steps of 3%.
- `sim/shaman/talents.go:152` — Only rank 1 was seen, beta will confirm that the ranks stack to 0.51 sec.
- `sim/shaman/talents.go:157` — Only rank 1 was seen, beta will confirm that rank 2 doubles both halves.
- `sim/shaman/talents.go:263` — Only rank 1 was seen, beta will confirm that the ranks go up in steps of 20%.
- `sim/shaman/talents.go:427` — Only rank 1 was seen, beta will confirm whether both chances really double at rank 2.
- `sim/shaman/talents.go:465` — The tooltip never showed a proc rate and only rank 1 was seen, beta will confirm both.
- `sim/shaman/talents.go:513` — The tooltip showed no cooldown, beta will confirm it. 3 minutes matches the other class cooldowns of this size.
- `sim/shaman/totems.go:9` — Assumed baseline rather than deleted, beta will confirm it.
- `sim/shaman/water_shield.go:18` — "Only one globe will activate every few seconds", the tooltip never said how long.
- `sim/shaman/water_totems.go:115` — The sim won't respect the value of a totem dropped via the APL. It uses hard-coded values from buffs.go manaRestoreBase := ManaSpringTotemManaRestore[rank]
- `sim/shaman/windfury_weapon.go:78` — Classic lets both weapons carry the imbue and gives the extra attacks to the hand that procced, beta will confirm that Forever kept both halves of that.

## Warlock (1)

- `sim/warlock/talents.go:451` — Beta will show whether the 33% of level is per rank or the full value

## Warrior (9)

- `sim/warrior/demoralizing_shout.go:15` — assumed baseline, beta will confirm
- `sim/warrior/shield_wall.go:29` — only rank 1 was shown, beta will confirm that rank 2 is another 5.5 minutes.
- `sim/warrior/shouts.go:55` — assumed baseline, beta will confirm
- `sim/warrior/talents.go:76` — only rank 1 was shown, beta will confirm the 1% crit / 3% armor / 1% extra attack per point.
- `sim/warrior/talents.go:175` — only rank 1 was shown, beta will confirm the 12% per point.
- `sim/warrior/talents.go:366` — only rank 1 was shown, beta will confirm that rank 2 doubles both the chance and the rage.
- `sim/warrior/talents.go:392` — only rank 1 was shown, beta will confirm the 2% per point.
- `sim/warrior/talents.go:403` — only rank 1 was shown, beta will confirm the 2% per point.
- `sim/warrior/talents.go:412` — only rank 1 was shown, beta will confirm the 1 Rage per point.

## Baseline ability changes

Forever also changes abilities that are not talents. None of their tooltips were shown at BlizzCon; what follows comes from the panel and the coverage of it, so every line is a claim to check against the beta spellbook rather than a number to confirm.

Warrior, the only class with concrete changes reported:

- Slam no longer resets the swing timer — modelled in `sim/warrior/slam.go`.
- Thunder Clap can be used in Defensive Stance — modelled in `sim/warrior/thunder_clap.go`. The shipped protection rotation does not cast it, so no numbers move until an APL does.
- Improved Shield Wall shortens the cooldown instead of lengthening the duration — modelled in `sim/warrior/shield_wall.go`.
- Tactical Mastery is baseline, with Improved Tactical Mastery on top — modelled in `sim/warrior/stances.go`.
- Victory Rush is baseline — not modelled. It needs a killing blow, which a boss encounter never gives before the fight ends.

Other classes: the panel spoke of baseline changes across every class without listing them, and nothing more specific has been published. When the beta client is datamined, diff each class spellbook against Classic Era and add every changed ability here with the file that models it, or the reason it is left out.

## Talents the sim does not read (23)

These are in the trees and the picker marks them as not simulated; spending points in them changes nothing. Most are utility or PvP talents the Classic sim never modelled either. They are listed so the beta pass can confirm none of them turned into something a raid rotation cares about.

- Druid / Balance: Overgrowth
- Druid / Restoration: Gift of the Earthmother
- Druid / Restoration: Wild Growth
- Hunter / Survival: Strider Kick
- Paladin / Holy: Voice of Truth
- Paladin / Holy: Light's Vigil
- Paladin / Holy: Infusion of Light
- Paladin / Protection: Improved Seal of Fury
- Priest / Discipline: Wand Specialization (the sim has no wand attacks at all, so a shadow build's two points here are idle)
- Priest / Discipline: Soul Warding
- Priest / Discipline: Renewed Hope
- Priest / Discipline: Divine Aegis
- Priest / Holy: Binding Heal
- Priest / Holy: Litany of Light
- Priest / Shadow Magic: Early Demise
- Rogue / Subtlety: Improved Distract
- Shaman / Restoration: Riptide
- Warlock / Demonology: Demonic Aegis
- Warlock / Demonology: Improved Felhunter
- Warlock / Destruction: Molten Skin
- Warrior / Arms: Spearing Strike
- Warrior / Fury: Blood Craze
- Warrior / Protection: Vanguard
