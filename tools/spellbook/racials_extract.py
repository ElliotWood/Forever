#!/usr/bin/env python3
"""Extract every racial ability from wowsims-forever.db into the page payload.

    python3 tools/spellbook/racials_extract.py      (run from the repo root)

Same payload shape as class_extract.py, with each race standing in for a skill
line, so the page template renders it unchanged. There are no talents and no
learn levels to chart, and the template hides both when they are empty.
"""
import json, os, re, sqlite3

DB = "tools/database/wowsims-forever.db"
SP = os.path.dirname(os.path.abspath(__file__))
A = "Field_1_60_1_69876_"
# SchoolMask is a mask, not an index, and the client names the combinations itself -
# these are its own STRING_SCHOOL_* strings out of GlobalStrings, not invented labels.
# Two combinations are live on these pages: 20 Frostfire (Frostfire Bolt) and
# 28 Elemental (Call of the Elements / Ancestors / Spirits). An earlier single-bit
# lookup rendered those two as the bare numbers "20" and "28".
SCHOOL_BITS = {1: "Physical", 2: "Holy", 4: "Fire", 8: "Nature", 16: "Frost",
               32: "Shadow", 64: "Arcane"}
SCHOOL = {
    1: "Physical", 2: "Holy", 4: "Fire", 8: "Nature", 16: "Frost", 32: "Shadow", 64: "Arcane",
    3: "Holystrike", 5: "Flamestrike", 9: "Stormstrike", 17: "Froststrike",
    33: "Shadowstrike", 65: "Spellstrike",
    6: "Radiant", 10: "Holystorm", 18: "Holyfrost", 34: "Twilight", 66: "Divine",
    12: "Volcanic", 20: "Frostfire", 36: "Shadowflame", 68: "Spellfire",
    24: "Froststorm", 40: "Plague", 72: "Astral",
    48: "Shadowfrost", 80: "Spellfrost", 96: "Spellshadow",
    28: "Elemental", 104: "Unknown", 124: "Chromatic", 126: "Magic", 127: "Chaos",
}


# SkillLineAbility.RaceMasks, a pair of 32-bit words. The column has no community
# name for build 69876 so it is addressed positionally as field 017; that is the
# same column the named 69893 extract calls RaceMasks, agreeing on all 7,824 rows.
# Word-0 bits 0-7 are named from the racial skill lines themselves - each line's
# abilities carry exactly one bit (Racial - Undead -> 16, Orc Racial -> 2, ...),
# so the mapping is read off the data rather than assumed.
RACE_BITS = {1: "Human", 2: "Orc", 4: "Dwarf", 8: "Night Elf",
             16: "Undead", 32: "Tauren", 64: "Gnome", 128: "Troll"}
# Word 1 bits 0 and 1 are the two Skyborne races, one per faction. Derived from
# the Riding lines: every Alliance mount row (Horse/Ram/Tiger/Mechanostrider)
# carries word-1 bit 0 on top of its three Alliance races, every Horde one
# (Wolf/Kodo/Raptor/Undead Horsemanship) carries bit 1, and Galestrider Riding -
# the Skyborne mount - is the eight classic races with neither bit.
RACE_BITS_HI = {1: "Skyborne (Alliance)", 2: "Skyborne (Horde)"}
# Two exact masks stand for whole factions. Not bit-decoded: they set bits for
# races this build has no other trace of. They are identified by containment -
# the first holds Human/Dwarf/Night Elf/Gnome, the second Orc/Undead/Tauren/Troll,
# and the two share no bit at all.
FACTION_MASK = {(-1321907123, 1427461461): "Alliance",
                (1309324210, -1440044374): "Horde"}


def race_names(raw):
    """Decode one RaceMasks value. None means the row places no race restriction."""
    lo, hi = json.loads(raw) if isinstance(raw, str) else (raw or [0, 0])
    if (lo, hi) in ((0, 0), (-1, -1)):
        return None
    if (lo, hi) in FACTION_MASK:
        return [FACTION_MASK[(lo, hi)]]
    out = []
    for i in range(32):
        if lo >> i & 1:
            out.append(RACE_BITS.get(1 << i, f"race bit {i}"))
    for i in range(32):
        if hi >> i & 1:
            out.append(RACE_BITS_HI.get(1 << i, f"race bit {i + 32}"))
    return out or None


def merge_races(masks):
    """Several SkillLineAbility rows can grant one spell. The spell is restricted
    only if every row is, and being learnable on both faction rows is no
    restriction either (Sense Undead, the city Teleports' pair)."""
    seen = []
    for m in masks:
        n = race_names(m)
        if n is None:
            return None
        for x in n:
            if x not in seen:
                seen.append(x)
    if "Alliance" in seen and "Horde" in seen:
        return None
    return seen or None


def school_name(mask):
    if mask in SCHOOL:
        return SCHOOL[mask]
    parts = [n for b, n in SCHOOL_BITS.items() if mask & b]
    return " + ".join(parts) if parts else f"mask {mask}"
# The client names these lines inconsistently: some "X Racial", some "Racial - X".
RACE = {"Dwarven Racial": "Dwarf", "Night Elf Racial": "Night Elf", "Orc Racial": "Orc",
        "Racial - Gnome": "Gnome", "Racial - Human": "Human", "Racial - Troll": "Troll",
        "Racial - Undead": "Undead", "Skyborne Racial": "Skyborne", "Tauren Racial": "Tauren"}
# Each is (label, Attributes word index, bit). Four of these are the flags the sim
# itself names in tools/database/dbc/enums.go; the three added below are not in
# that file and are identified from the client data, with the evidence recorded.
# IS_CHANNELLED / IS_SELF_CHANNELLED: verified by census rather than taken on
# trust - Attributes[1] 0x4 lands on Mind Flay, Drain Life, Arcane Missiles and
# Penance, 0x40 on Blizzard, Hurricane, Volley, Tranquility and Bladestorm, i.e.
# the target channel and the ground/self channel. A channel stores 0 in
# SpellCastTimes, so without this flag every channel reads as "Instant".
# PERIODIC_CAN_CRIT is Attributes[8] 0x200, named on wowdev.wiki
# (Spell.dbc/Attributes) as SPELL_ATTR8_PERIODIC_CAN_CRIT, "allows DoTs and HoTs
# to be critical". 225 rows here, all 49 names damage-over-time.
# Not to be confused with Attributes[10] 0x4000 SPELL_ATTR10_ROLLING_PERIODIC
# (the Ignite/Blackout Kick roll-the-remainder flag), which an earlier pass
# mislabelled this bit as. That one is checked and carried by ZERO spells in
# build 69876, so it is deliberately not extracted or displayed.
ATTR_FLAGS = [
    ("IS_CHANNELLED", 1, 0x4,
     "Cast is a channel on the target; SpellCastTimes holds 0, the length is the aura duration."),
    ("IS_SELF_CHANNELLED", 1, 0x40,
     "Channelled on the caster or the ground rather than a target (Blizzard, Hurricane, Volley)."),
    ("CANT_CRIT", 2, 0x20000000,
     "The direct hit can never critically strike."),
    ("CAN_PROC_FROM_PROCS", 3, 0x4000000,
     "May be triggered by an event that was itself a proc, instead of only by a real cast or swing."),
    ("PERIODIC_CAN_CRIT", 8, 0x200,
     "Individual ticks of the damage or heal over time can critically strike."),
    ("SCALES_WITH_ITEM_LEVEL", 11, 0x4,
     "Effect values scale from the item level of the item that granted the spell."),
    ("ONLY_PROC_FROM_CLASS_ABILITIES", 12, 0x80000000,
     "Only class abilities can trigger it - auto attacks and item effects cannot."),
]


def attr_table(flags):
    """Bit and meaning for each flag name, so the page can table them."""
    return {name: {"bit": f"Attributes[{word}] 0x{bit:X}", "note": note}
            for name, word, bit, note in flags}


def attr_flags(q, spell_id):
    row = q("select " + ",".join(f"coalesce(Attributes_{i},0)" for i in range(16)) +
            " from SpellMisc where SpellID=?", (spell_id,)).fetchone()
    if not row:
        return []
    return [name for name, word, bit, _ in ATTR_FLAGS if row[word] & bit]

# The client never links a proc listener to the spell it actually casts: there is no
# EffectTriggerSpell on any of them. Three indirect routes do work, in order of trust -
# an id written into this spell's own description (Cannibalize names $20578s1), a spell
# whose description names this one (Touch of the Grave's drain cites $1260189s1), and a
# same-named spell close by in id space (which is how rank 2 of Touch of the Grave, 10%,
# reaches the same drain). Candidates that carry a proc mask of their own are dropped -
# those are sibling listeners, not the effect - as are other abilities already on the page.
# Restricted to abilities that actually proc: applied more widely the same-name route
# collects unrelated NPC spells, and a bare LIKE on "$122" also matches "$12257".
def proc_spells(q, one, spell_id, name, desc, effects, own_ids):
    cand = {e["trig"] for e in effects if e.get("trig")}
    cand |= {int(m) for m in re.findall(r"\$(\d{2,})(?!\d)", desc or "")}
    for r in q("select ID,Description_lang from Spell where Description_lang like ?",
               (f"%${spell_id}%",)):
        if re.search(rf"\${spell_id}(?!\d)", r["Description_lang"] or ""):
            cand.add(r["ID"])
    for r in q("select ID from SpellName where Name_lang=? and ID!=?", (name, spell_id)):
        if abs(r["ID"] - spell_id) < 10000:
            cand.add(r["ID"])

    out = []
    for c in sorted(cand):
        if c == spell_id or c in own_ids:
            continue
        ao = one("select ProcTypeMask_0 p from SpellAuraOptions where SpellID=?", c)
        if ao.get("p") or 0:
            continue
        nm = one("select Name_lang n from SpellName where ID=?", c).get("n")
        if not nm:
            continue
        misc = one("select SchoolMask, Attributes_2 a2, Attributes_8 a8 from SpellMisc where SpellID=?", c)
        dur = one("select d.Duration ms from SpellDuration d join SpellMisc m on m.DurationIndex=d.ID where m.SpellID=?", c)
        eff = [dict(e) for e in q("""select EffectIndex i, Effect e, EffectAura aura, EffectBasePointsF bp,
                                     EffectAuraPeriod per, EffectBonusCoefficient coef
                                     from SpellEffect where SpellID=? and Effect!=0 order by EffectIndex""", (c,))]
        if not eff:
            continue
        sm = misc.get("SchoolMask") or 1
        out.append({
            "id": c, "name": nm,
            "school": school_name(sm),
            "schoolMask": sm,
            "schoolBits": [n for b, n in SCHOOL_BITS.items() if sm & b],
            "duration": dur.get("ms") or 0,
            "cantCrit": bool((misc.get("a2") or 0) & 0x20000000),
            "periodicCanCrit": bool((misc.get("a8") or 0) & 0x200),
            "effects": eff,
        })
    return out



CLASS_BIT = {1: "Warrior", 2: "Paladin", 4: "Hunter", 8: "Rogue", 16: "Priest",
             64: "Shaman", 128: "Mage", 256: "Warlock", 1024: "Druid"}


def classes_of(mask):
    if mask in (0, -1):
        return []
    return [n for b, n in sorted(CLASS_BIT.items()) if mask & b]


def tidy(text):
    text = re.sub(r"\bby -(\d)", r"by \1", text)
    return re.sub(r"(\d)\.0(?!\d)", r"\1", text)


def main():
    db = sqlite3.connect(DB)
    db.row_factory = sqlite3.Row
    q = db.execute

    def one(sql, *a):
        r = q(sql, a).fetchone()
        return dict(r) if r else {}

    rows = q(f"""select a.{A}004 spell, sl.DisplayName_lang line, a.{A}006 classmask,
                        a.{A}017 racemask
                 from SkillLineAbility a join SkillLine sl on sl.ID = a.{A}003
                 where sl.DisplayName_lang like '%Racial%'""").fetchall()

    out, seen = [], set()
    for r in rows:
        sid = r["spell"]
        race = RACE.get(r["line"], r["line"])
        if (sid, race) in seen:
            continue
        seen.add((sid, race))
        nm = one("select Name_lang n from SpellName where ID=?", sid).get("n")
        if not nm:
            continue
        sp = one("select NameSubtext_lang sub, Description_lang d, AuraDescription_lang ad from Spell where ID=?", sid)
        misc = one("select SchoolMask, SpellIconFileDataID icon from SpellMisc where SpellID=?", sid)
        cd = one("select RecoveryTime, CategoryRecoveryTime, StartRecoveryTime from SpellCooldowns where SpellID=?", sid)
        cat = one("select Category from SpellCategories where SpellID=?", sid)
        pw = one("select PowerType, ManaCost, PowerCostPct from SpellPower where SpellID=? order by OrderIndex", sid)
        ct = one("""select ct.Base b from SpellMisc m join SpellCastTimes ct on ct.ID=m.CastingTimeIndex
                    where m.SpellID=?""", sid)
        dur = one("""select d.Duration ms from SpellDuration d join SpellMisc m on m.DurationIndex=d.ID
                     where m.SpellID=?""", sid)
        rng = one("""select r.RangeMax_1 mx from SpellRange r join SpellMisc m on m.RangeIndex=r.ID
                     where m.SpellID=?""", sid)
        eq = one("select EquippedItemClass c, EquippedItemSubclass sc from SpellEquippedItems where SpellID=?", sid)
        ao = one("""select ProcTypeMask_0 p0, ProcTypeMask_1 p1, ProcChance pc,
                         ProcCategoryRecovery icd, ProcCharges chg, CumulativeAura stacks
                  from SpellAuraOptions where SpellID=?""", sid)
        effs = [dict(e) for e in q("""select EffectIndex i, Effect e, EffectAura aura, EffectBasePointsF bp,
                                      EffectBonusCoefficient coef, BonusCoefficientFromAP ap,
                                      EffectAuraPeriod per, EffectChainTargets ct, EffectTriggerSpell trig
                                      from SpellEffect where SpellID=? order by EffectIndex""", (sid,))]
        ptype, raw = pw.get("PowerType"), pw.get("ManaCost") or 0
        power = raw / 10.0 if ptype == 1 else raw
        if power == int(power):
            power = int(power)
        out.append({
            "id": sid, "name": nm, "rank": sp.get("sub") or "",
            "desc": sp.get("d") or "", "auraDesc": sp.get("ad") or "",
            "tree": race, "supercedes": 0, "spellLevel": 0,
            # The Skyborne line covers two races, one per faction, so the mask is
            # the only thing that separates Read Ley Line from Skysight.
            "races": race_names(r["racemask"]),
            "school": school_name(misc.get("SchoolMask") or 1),
            "schoolMask": misc.get("SchoolMask") or 1,
            "schoolBits": [n for b, n in SCHOOL_BITS.items() if (misc.get("SchoolMask") or 1) & b],
            "power": power, "powerPct": round(pw.get("PowerCostPct") or 0, 1), "powerType": ptype,
            "cd": max(cd.get("RecoveryTime") or 0, cd.get("CategoryRecoveryTime") or 0),
            "cdOwn": cd.get("RecoveryTime") or 0, "cdCat": cd.get("CategoryRecoveryTime") or 0,
            "catId": cat.get("Category") or 0,
            "castTime": max(ct.get("b") or 0, 0), "castWeapon": (ct.get("b") or 0) < 0,
            "gcd": cd.get("StartRecoveryTime") or 0, "stances": [],
            "range": rng.get("mx") or 0, "duration": dur.get("ms") or 0,
            "effects": effs, "iconFdid": misc.get("icon") or 0,
            "rankValues": {}, "sharesWith": [],
            # Extras specific to this page.
            "classes": classes_of(r["classmask"]),
            "weaponClass": eq.get("c", -1), "weaponSub": eq.get("sc", 0),
            "attrFlags": attr_flags(q, sid),
            "procMask": (ao.get("p0") or 0) & 0xFFFFFFFF,
            "procMask2": (ao.get("p1") or 0) & 0xFFFFFFFF,
            "procChance": ao.get("pc") or 0,
            # ProcCategoryRecovery is the internal cooldown, in ms.
            "procIcd": ao.get("icd") or 0,
            "procCharges": ao.get("chg") or 0,
            "procStacks": ao.get("stacks") or 0,
        })

    # icons, through the community listfile
    want = {a["iconFdid"] for a in out if a["iconFdid"]}
    icon_by_fdid = {}
    with open("tools/db2tool/listfile.csv", encoding="utf-8", errors="replace") as fh:
        for line in fh:
            sc = line.find(";")
            if sc < 0:
                continue
            try:
                fd = int(line[:sc])
            except ValueError:
                continue
            if fd in want and "interface/icons/" in line:
                icon_by_fdid[fd] = line[sc + 1:].strip().rsplit("/", 1)[-1][:-4]
    for a in out:
        a["icon"] = icon_by_fdid.get(a.pop("iconFdid"), "")

    resolved = json.load(open(os.path.join(SP, "data", "descriptions.json")))
    byid = {a["id"]: a for a in out}
    extern = {}
    for a in out:
        r = resolved.get(str(a["id"]))
        if r:
            a["descRaw"], a["desc"] = a["desc"], tidy(r)
        trig = []
        for e in a["effects"]:
            ts = e.get("trig") or 0
            if ts and ts not in trig:
                trig.append(ts)
        a["triggers"] = trig
        for ts in trig:
            if ts not in byid and ts not in extern:
                extern[ts] = one("select Name_lang n from SpellName where ID=?", ts).get("n") or ""
        a["siblings"] = []
        a["talent"] = None

    fams = [{"root": a["id"], "name": a["name"], "tree": a["tree"], "ranks": [a["id"]],
             "variant": 1, "variants": 1} for a in out]

    enums = json.load(open(os.path.join(SP, "data", "enums.json")))
    # Resolve what each proc actually casts, so the tooltip can show the real numbers
    # rather than just "this procs". Only abilities with a proc mask qualify.
    own_ids = {a["id"] for a in out}
    for a in out:
        a["procSpells"] = (proc_spells(q, one, a["id"], a["name"], a["desc"], a["effects"], own_ids)
                           if a.get("procMask") else [])

    used_p = sorted({str(b) for a in out for b in range(32) if (a.get("procMask") or 0) >> b & 1})
    # Proc spells contribute effect and aura ids of their own, which the page needs
    # named too - without this the tooltip prints bare numbers for them.
    used_e = sorted({str(e["e"]) for a in out for e in a["effects"]} |
                    {str(e["e"]) for a in out for sp in a.get("procSpells", []) for e in sp["effects"]})
    used_a = sorted({str(e["aura"]) for a in out for e in a["effects"]} |
                    {str(e["aura"]) for a in out for sp in a.get("procSpells", []) for e in sp["effects"]})
    trees = sorted({a["tree"] for a in out})

    payload = {
        "meta": {"build": "1.60.1.69876", "product": "wow_classic_beta", "cls": "Racial",
                 "abilities": len(out), "families": len(fams), "talents": 0},
        "cfg": {"className": "Racial", "trees": trees, "chartTrees": [], "forms": [],
                "formWord": " form", "talentTabs": [], "notes": []},
        "abilities": out, "families": fams, "talents": [], "extern": extern,
        "effectNames": {k: enums["effect"].get(k, "EFFECT_" + k) for k in used_e},
        "auraNames": {k: enums["aura"].get(k, "AURA_" + k) for k in used_a},
        "procNames": {k: enums.get("proc", {}).get(k, {"name": "bit " + k, "note": ""}) for k in used_p},
        "attrNames": attr_table(ATTR_FLAGS),
    }
    json.dump(payload, open(os.path.join(SP, "data", "racial.json"), "w"), separators=(",", ":"))
    print(f"Racial: {len(out)} abilities across {len(trees)} races -> {trees}")


if __name__ == "__main__":
    main()
