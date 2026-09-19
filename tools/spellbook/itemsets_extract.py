#!/usr/bin/env python3
"""Extract every item set from the Classic Forever client into data/itemsets.json.

    python3 tools/spellbook/itemsets_extract.py      (run from the repo root)

Reads tools/database/forever-items.db, which is the same 80 tables as
wowsims-forever.db but re-decoded after WoWDBDefs gained real column names for
build 1.60.1.69876. That matters here: Item and ItemSparse were both
placeholder-named in the original extraction, and ItemSparse has 69 columns
against the 75 of the 1.15.x builds, so a positional mapping would not have
worked. The refresh also, incidentally, confirms the SkillLineAbility column
guesses the spellbook has been relying on.

One gap the data itself imposes: ItemSparse holds 19,171 of the client's 31,675
items, and only 824 of the 3,049 item-set members are among them. Everything a
set needs to be *grouped* - armour type, slot, icon - lives in Item and is
present for all of them, but name, item level and stats live in ItemSparse. The
missing rows are hotfix rows: this machine's _classic_beta_ DBCache.bin is for
build 61916, so db2tool's overlay had nothing to apply for 69876. Launching the
client once and re-extracting fills them in. Items without a row are emitted
with "sparse": false so the page can say so rather than imply the slot is empty.
"""
import json, os, re, sqlite3, sys

DB = "tools/database/forever-items.db"
SOD = "tools/database/wowsims-sod.db"
SP = os.path.dirname(os.path.abspath(__file__))

CLASS_MASK = {1: "Warrior", 2: "Paladin", 4: "Hunter", 8: "Rogue", 16: "Priest",
              64: "Shaman", 128: "Mage", 256: "Warlock", 1024: "Druid"}

QUALITY = {0: "Poor", 1: "Common", 2: "Uncommon", 3: "Rare", 4: "Epic",
           5: "Legendary", 6: "Artifact", 7: "Heirloom"}

SLOT = {1: "Head", 2: "Neck", 3: "Shoulder", 4: "Shirt", 5: "Chest", 6: "Waist",
        7: "Legs", 8: "Feet", 9: "Wrist", 10: "Hands", 11: "Finger", 12: "Trinket",
        13: "One-Hand", 14: "Off Hand", 15: "Ranged", 16: "Back", 17: "Two-Hand",
        18: "Bag", 19: "Tabard", 20: "Chest", 21: "Main Hand", 22: "Off Hand",
        23: "Held in Off-hand", 24: "Ammo", 25: "Thrown", 26: "Ranged", 28: "Relic"}

# ItemSparse.Resistances is the classic six-wide array and index 0 is armour,
# not a resistance - tools/database/dbc/item.go reads it the same way.
RESIST = {1: "Holy", 2: "Fire", 3: "Nature", 4: "Frost", 5: "Shadow", 6: "Arcane"}

# ITEM_MOD_* from tools/database/dbc/enums.go, given the label the game uses.
STAT = {0: "Mana", 1: "Health", 3: "Agility", 4: "Strength", 5: "Intellect",
        6: "Spirit", 7: "Stamina", 12: "Defense Rating", 13: "Dodge Rating",
        14: "Parry Rating", 15: "Block Rating", 16: "Melee Hit Rating",
        17: "Ranged Hit Rating", 18: "Spell Hit Rating", 19: "Melee Crit Rating",
        20: "Ranged Crit Rating", 21: "Spell Crit Rating", 28: "Melee Haste Rating",
        29: "Ranged Haste Rating", 30: "Spell Haste Rating", 31: "Hit Rating",
        32: "Crit Rating", 35: "Resilience Rating", 36: "Haste Rating",
        37: "Expertise Rating", 38: "Attack Power", 39: "Ranged Attack Power",
        40: "Feral Attack Power", 41: "Healing Power", 42: "Spell Damage",
        43: "Mana per 5 sec", 44: "Armor Penetration Rating", 45: "Spell Power",
        46: "Health per 5 sec", 47: "Spell Penetration", 48: "Block Value"}


def arr(v):
    return json.loads(v) if isinstance(v, str) else (v or [])


def classes_from_mask(mask):
    if not mask or mask == -1 or mask == 0xFFFFFFFF:
        return []
    out = [n for bit, n in CLASS_MASK.items() if mask & bit]
    return [] if len(out) == len(CLASS_MASK) else out


def main():
    if not os.path.exists(DB):
        sys.exit(f"missing {DB} - see the module docstring for how it is produced")
    c = sqlite3.connect(DB)
    c.row_factory = sqlite3.Row
    q = lambda sql, a=(): c.execute(sql, a).fetchall()

    # Forever-only is decided by comparing against the Season of Discovery client,
    # the same discriminator the spellbook uses for SoD-space spells. A set that
    # exists in both but had its member list rebuilt is caught by the second test
    # below, on item ids.
    sod_sets, sod_items = set(), set()
    if os.path.exists(SOD):
        s = sqlite3.connect(SOD)
        sod_sets = {r[0] for r in s.execute("select ID from ItemSet")}
        sod_items = {r[0] for r in s.execute("select ID from Item")}

    subclass = {(r["ClassID"], r["SubClassID"]): r["DisplayName_lang"]
                for r in q("select ClassID,SubClassID,DisplayName_lang from ItemSubClass")}
    item_rows = {r["ID"]: r for r in q(
        "select ID,ClassID,SubclassID,InventoryType,IconFileDataID from Item")}
    sparse_rows = {r["ID"]: r for r in q("""select ID,Display_lang,Description_lang,ItemLevel,
        RequiredLevel,InventoryType,OverallQualityID,ItemSet,AllowableClass,Bonding,
        StatModifier_bonusStat,StatModifier_bonusAmount,Resistances,MinDamage,MaxDamage,
        ItemDelay,DmgVariance from ItemSparse""")}

    spell_names = {r["ID"]: r["Name_lang"] for r in q("select ID,Name_lang from SpellName")}
    descs = _load(os.path.join(SP, "data", "descriptions.json"))

    def armor_type(it):
        # Only class 4 (armour) subclasses name a material; a set of rings or
        # trinkets is class 4 subclass 0 "Miscellaneous", and weapons are class 2.
        if not it:
            return "Unknown"
        if it["ClassID"] == 4 and it["SubclassID"] in (1, 2, 3, 4):
            return subclass.get((4, it["SubclassID"]), "Unknown")
        if it["ClassID"] == 4 and it["SubclassID"] == 6:
            return "Shield"
        if it["ClassID"] == 2:
            return "Weapon"
        return "Other"

    sets = []
    for r in q("select ID,Name_lang,SetFlags,RequiredSkill,RequiredSkillRank,ItemID from ItemSet"):
        member_ids = [i for i in arr(r["ItemID"]) if i]
        if not member_ids:
            continue

        members, types, slots, classes, ilvls, missing = [], [], [], set(), [], 0
        for iid in member_ids:
            it, sp = item_rows.get(iid), sparse_rows.get(iid)
            entry = {"id": iid, "sparse": sp is not None,
                     "icon": it["IconFileDataID"] if it else 0,
                     "slot": SLOT.get((sp or it or {})["InventoryType"] if (sp or it) else 0, ""),
                     "armor": armor_type(it)}
            if it:
                types.append(entry["armor"])
                if entry["slot"]:
                    slots.append(entry["slot"])
            if sp is None:
                missing += 1
            else:
                stats = []
                for st, amt in zip(arr(sp["StatModifier_bonusStat"]),
                                   arr(sp["StatModifier_bonusAmount"])):
                    if st < 0 or not amt:
                        continue
                    stats.append({"name": STAT.get(st, f"stat {st}"), "value": round(amt)})
                res = arr(sp["Resistances"])
                entry.update({
                    "name": sp["Display_lang"] or "", "desc": sp["Description_lang"] or "",
                    "ilvl": sp["ItemLevel"] or 0, "reqLevel": sp["RequiredLevel"] or 0,
                    "quality": QUALITY.get(sp["OverallQualityID"], str(sp["OverallQualityID"])),
                    "qualityId": sp["OverallQualityID"] or 0,
                    "bonding": sp["Bonding"] or 0,
                    "stats": stats,
                    "armorValue": res[0] if res else 0,
                    "resist": [{"name": RESIST[i], "value": res[i]}
                               for i in range(1, min(len(res), 7)) if res[i]],
                    "minDmg": round(sp["MinDamage"] or 0, 1),
                    "maxDmg": round(sp["MaxDamage"] or 0, 1),
                    "speed": round((sp["ItemDelay"] or 0) / 1000.0, 2),
                })
                if sp["ItemLevel"]:
                    ilvls.append(sp["ItemLevel"])
                for n in classes_from_mask(sp["AllowableClass"]):
                    classes.add(n)
            members.append(entry)

        bonuses = []
        for b in q("select SpellID,Threshold from ItemSetSpell where ItemSetID=? order by Threshold",
                   (r["ID"],)):
            sid = b["SpellID"]
            bonuses.append({"spellId": sid, "pieces": b["Threshold"],
                            "name": spell_names.get(sid, f"Spell {sid}"),
                            "text": descs.get(str(sid), "")})

        # Tier lives in the bonus spell's internal name ("1.60.0 - Item - Tier 1 -
        # Paladin - Protection 5P Bonus"), which is the only place it is written
        # down. Classic sets predate the convention and simply have no tier.
        tier, spec = 0, ""
        for b in bonuses:
            m = re.search(r"Tier (\d+)", b["name"])
            if m:
                tier = int(m.group(1))
            m = re.search(r"- (\w+) - (\w+) \d+P", b["name"])
            if m and not spec:
                spec = m.group(2)
        # Class from the members' own AllowableClass, never from the set name:
        # a classic set will not match the modern naming convention at all.
        if not classes:
            for b in bonuses:
                m = re.search(r"- (Warrior|Paladin|Hunter|Rogue|Priest|Shaman|Mage|Warlock|Druid) -",
                              b["name"])
                if m:
                    classes.add(m.group(1))

        kinds = [t for t in types if t not in ("Other", "Unknown")]
        primary = max(set(kinds), key=kinds.count) if kinds else "Other"
        new_by_id = r["ID"] not in sod_sets
        new_by_items = bool(sod_items) and all(i not in sod_items for i in member_ids)
        sets.append({
            "id": r["ID"], "name": r["Name_lang"] or f"Set {r['ID']}",
            "armorType": primary,
            "armorTypes": sorted(set(types)),
            "classes": sorted(classes), "tier": tier, "spec": spec,
            "pieces": len(member_ids), "members": members, "bonuses": bonuses,
            "ilvl": max(ilvls) if ilvls else 0,
            "missingStats": missing,
            "isNew": new_by_id or new_by_items,
            "newReason": "set id not in the SoD client" if new_by_id else
                         ("every item id is new" if new_by_items else ""),
            "reqSkill": r["RequiredSkill"] or 0, "reqSkillRank": r["RequiredSkillRank"] or 0,
        })

    sets.sort(key=lambda s: (-s["isNew"], -s["tier"], s["name"]))
    payload = {
        "meta": {"build": "1.60.1.69876", "product": "wow_classic_beta",
                 "sets": len(sets), "new": sum(1 for s in sets if s["isNew"]),
                 "items": sum(s["pieces"] for s in sets),
                 "withStats": sum(sum(1 for m in s["members"] if m["sparse"]) for s in sets)},
        "sets": sets,
    }
    path = os.path.join(SP, "data", "itemsets.json")
    json.dump(payload, open(path, "w"), separators=(",", ":"))
    m = payload["meta"]
    print(f"{m['sets']} sets ({m['new']} new), {m['items']} items, "
          f"{m['withStats']} with ItemSparse stats -> {path}")


def _load(path):
    return json.load(open(path)) if os.path.exists(path) else {}


if __name__ == "__main__":
    main()
