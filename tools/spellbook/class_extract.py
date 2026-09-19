#!/usr/bin/env python3
"""Extract one class's spellbook and talent tree from wowsims-forever.db.

    python3 tools/spellbook/class_extract.py <Class>      (run from the repo root)

Writes tools/spellbook/data/<class>.json, the payload the page template consumes.

Two things in this database are not where a TBC-shaped pipeline would look:

* SkillLineAbility has no community column names for build 69876, so its columns
  are addressed positionally (see A below).
* Talents are NOT the legacy Talent/TalentTab tables. Those still hold the vanilla
  layout and are vestigial: for Warrior they miss 9 talents the client actually
  uses and list 10 it does not. The live tree is the modern trait system --
  TraitNode (grid position) / TraitNodeEntry (max ranks) / TraitDefinition (spell)
  / TraitEdge (prerequisites) / TraitDefinitionEffectPoints -> Curve (rank values).
  TalentTab is still used, but only to name the three specs and order them.
"""
import json, os, re, sqlite3, sys

DB = "tools/database/wowsims-forever.db"
SP = os.path.dirname(os.path.abspath(__file__))
A = "Field_1_60_1_69876_"          # SkillLineAbility's unnamed columns
SL_LINE, SL_SPELL, SL_RANKREQ, SL_MASK, SL_SUPER = "003", "004", "005", "006", "007"
SL_RACE = "017"                    # RaceMasks; see race_names() for how it is pinned

CLASS_MASK = {"Warrior": 1, "Paladin": 2, "Hunter": 4, "Rogue": 8, "Priest": 16,
              "Shaman": 64, "Mage": 128, "Warlock": 256, "Druid": 1024}
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


SKIP_LINES = {183, 777}            # GENERIC (DND) and Mounts are not abilities
# Season of Discovery's spell space. Forever imported SoD wholesale at the same
# ids: of the class abilities in this band, 85% also exist in the wow_classic_era
# client, and all 458 Engraving/Runes entries match by id. Excluded here so the
# pages show the classic core plus Forever's own 1.0M+ additions.
# Talents are deliberately NOT filtered - 40 of them sit in this band (Penance,
# Riptide, Incinerate, Hot Streak...) and removing them would leave holes in the
# tree grid and break prerequisite chains.
# The bound is not a round number picked by eye: sorting the class-ability ids in
# the era build shows a 367,178-wide gap between the classic block (ends at 31018)
# and the SoD block (starts at 398196, ends at 469145), and Forever's own content
# starts far above at 1214168. Nothing at all sits between 31018 and 398196, so
# any cut inside that gap is equivalent; 100000 is used for clarity. A 400000 cut
# would have left 33 abilities behind, including every class's Chest/Hands/Legs
# Rune Ability at 399954/399966/399967 and Rogue's Mutilate and Shadowstrike.
SOD_ID_MIN, SOD_ID_MAX = 100000, 1000000


def is_sod(spell_id):
    return SOD_ID_MIN <= spell_id < SOD_ID_MAX
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
FORM = {1: "Cat", 3: "Travel", 4: "Aquatic", 5: "Bear", 8: "Dire Bear", 15: "Cat",
        17: "Battle", 18: "Defensive", 19: "Berserker", 28: "Shadowform",
        30: "Stealth", 31: "Stealth"}


def tidy(text):
    """Cosmetic fixes on a resolved tooltip string; descRaw keeps the original.

    This client stores reductions as negative base points (the TBC database has 0
    for the same spells), and the sentence already says "Reduces", so the minus
    would read twice. A trailing .0 carries no information.
    """
    text = re.sub(r"\bby -(\d)", r"by \1", text)
    return re.sub(r"(\d)\.0(?!\d)", r"\1", text)


def curve_values(q, curve_id):
    return [r[0] for r in q("select Pos_1 from CurvePoint where CurveID=? order by OrderIndex", (curve_id,))]


def trait_points(q, spell_id):
    """{effectIndex: [value per rank]} for a spell, via the trait system.

    The spell's own EffectBasePointsF is only a snapshot of one rank and is not
    reliably the first or the last: Toughness stores its rank-5 value, Deflection
    its rank-1, and Deep Wounds stores 0 while its curve runs 20/40/60.
    """
    out = {}
    for ei, cid in q("""select ep.EffectIndex, ep.CurveID from TraitDefinition d
                        join TraitDefinitionEffectPoints ep on ep.TraitDefinitionID=d.ID
                        where d.SpellID=? and ep.CurveID is not null""", (spell_id,)):
        if ei not in out:
            vals = curve_values(q, cid)
            if vals:
                out[ei] = vals
    return out


# The trees are laid out on a fixed grid: four columns per spec 600 apart, specs
# about 2200 apart, seven rows 600 apart. A handful of nodes carry small jitter
# (10-200 units) and three across all nine classes are parked far off-grid
# entirely (Hunter's Lightning Reflexes and Improved Serpent Sting, Priest's Holy
# Specialization). Jitter is snapped; the far outliers keep no position and are
# reported separately rather than given an invented one.
GRID_X, GRID_Y = 1020, 2130      # first column, first row
STEP = 600
SPEC_GAP = 1200
SNAP_X, SNAP_Y = 300, 100


def spec_layout(q, tree):
    """Return (snapX, snapY) helpers mapping a node position onto the grid.

    Each returns None when the coordinate is too far off-grid to place.
    """
    xs = sorted({r[0] for r in q("select distinct PosX from TraitNode where TraitTreeID=?", (tree,))})
    ys = sorted({r[0] for r in q("select distinct PosY from TraitNode where TraitTreeID=?", (tree,))})

    def cluster(vals, tol):
        out, cur = [], [vals[0]]
        for a, b in zip(vals, vals[1:]):
            if b - a <= tol:
                cur.append(b)
            else:
                out.append(cur)
                cur = [b]
        out.append(cur)
        return [sum(c) // len(c) for c in out]

    # canonical rows: clusters that sit on the regular 600 spacing
    rows = [c for c in cluster(ys, SNAP_Y) if (c - GRID_Y) % STEP < SNAP_Y or (GRID_Y - c) % STEP < SNAP_Y]
    rows = [c for c in rows if abs(c - GRID_Y) <= STEP * 20]
    cols = [c for c in cluster(xs, SNAP_X) if abs(c - GRID_X) <= STEP * 20]
    specs, cur = [], [cols[0]]
    for a, b in zip(cols, cols[1:]):
        if b - a < SPEC_GAP:
            cur.append(b)
        else:
            specs.append(cur)
            cur = [b]
    specs.append(cur)

    def snap_y(y):
        for i, r in enumerate(sorted(rows)):
            if abs(y - r) <= SNAP_Y:
                return i
        return None

    def snap_x(x):
        for gi, g in enumerate(specs):
            for ci, c in enumerate(g):
                if abs(x - c) <= SNAP_X:
                    return gi, ci
        return None

    return snap_x, snap_y


def find_tree(q, tabs):
    """The class's trait tree, found by seeing where its legacy talent spells live."""
    votes = {}
    for (sr,) in q(f"select SpellRank from Talent where TabID in ({','.join(map(str, tabs))})"):
        ranks = [x for x in json.loads(sr) if x]
        if not ranks:
            continue
        for (tid,) in q("""select n.TraitTreeID from TraitDefinition d
                           join TraitNodeEntry ne on ne.TraitDefinitionID=d.ID
                           join TraitNodeXTraitNodeEntry x on x.TraitNodeEntryID=ne.ID
                           join TraitNode n on n.ID=x.TraitNodeID where d.SpellID=?""", (ranks[0],)):
            votes[tid] = votes.get(tid, 0) + 1
    return max(votes, key=votes.get) if votes else None


def talents_from_tree(q, tree, spec_names, one):
    snap_x, snap_y = spec_layout(q, tree)
    rows = q("""select n.ID, n.PosX, n.PosY, ne.MaxRanks, d.ID, d.SpellID,
                       coalesce(sn.Name_lang, d.OverrideName_lang, '') nm,
                       coalesce(nullif(d.OverrideDescription_lang,''), sp.Description_lang, '') descr,
                       d.OverrideIcon
                from TraitNode n
                join TraitNodeXTraitNodeEntry x on x.TraitNodeID=n.ID
                join TraitNodeEntry ne on ne.ID=x.TraitNodeEntryID
                left join TraitDefinition d on d.ID=ne.TraitDefinitionID
                left join SpellName sn on sn.ID=d.SpellID
                left join Spell sp on sp.ID=d.SpellID
                where n.TraitTreeID=?""", (tree,)).fetchall()
    talents, by_node = [], {}
    for nid, px, py, maxr, defid, spell, nm, descr, ovicon in rows:
        xy, tier = snap_x(px), snap_y(py)
        off = xy is None or tier is None
        gi, ci = xy if xy else (0, 0)
        icon_fdid = ovicon or (one("select SpellIconFileDataID i from SpellMisc where SpellID=?", spell).get("i") or 0)
        t = {"id": nid, "tab": spec_names[gi] if gi < len(spec_names) else f"Spec {gi + 1}",
             "tier": tier if tier is not None else 0, "col": ci, "offGrid": off,
             "name": nm or f"Spell {spell}",
             "desc": descr or "", "ranks": [spell] if spell else [],
             "maxRank": maxr or 1, "nameFrom": spell or 0,
             "prereq": [], "prereqRank": [], "iconFdid": icon_fdid,
             "effects": talent_effects(q, spell) if spell else []}
        talents.append(t)
        by_node[nid] = t
    # TraitEdge.Type follows the modern trait system: 0 draws a line but gates
    # nothing, 2 is sufficient-for-availability and 3 is required-for-availability.
    # Treating every row as a prerequisite produces an upward arrow on Druid's
    # Nature's Majesty / Nature's Splendor pair (a Type 0 edge pointing back up)
    # and would deadlock both talents, since each would require the other.
    #
    # A pair can also appear twice, once in each direction (Hunter's Intimidation
    # and Bestial Wrath), so direction comes from the grid rather than from which
    # column the row happens to put first: the node higher up the tree is the
    # prerequisite. Two genuine sideways links sit on one row (Paladin's Holy
    # Shock into Divine Precision, Priest's Mind Flay into Improved Mind Flay);
    # those keep the stored Left -> Right order.
    seen_pairs = set()
    for lid, rid, etype in q("""select e.LeftTraitNodeID, e.RightTraitNodeID, e.Type
                                from TraitEdge e
                                join TraitNode n on n.ID=e.LeftTraitNodeID
                                where n.TraitTreeID=? and e.Type in (2, 3)""", (tree,)):
        if lid not in by_node or rid not in by_node:
            continue
        pair = frozenset((lid, rid))
        if pair in seen_pairs:
            continue
        seen_pairs.add(pair)
        src, dst = by_node[lid], by_node[rid]
        if src["tier"] > dst["tier"]:
            src, dst = dst, src
        if src["id"] in dst["prereq"]:
            continue
        dst["prereq"].append(src["id"])
        dst["prereqRank"].append(src["maxRank"])
    return talents


def talent_effects(q, spell):
    """One row per effect index, carrying that effect's value at each rank."""
    traits = trait_points(q, spell)
    rows = q("""select EffectIndex, Effect, EffectAura, EffectBasePointsF, EffectBonusCoefficient,
                       BonusCoefficientFromAP, EffectAuraPeriod, EffectMiscValue_0, EffectTriggerSpell
                from SpellEffect where SpellID=? order by EffectIndex""", (spell,)).fetchall()
    out = []
    for i, e, aura, bp, coef, ap, per, mv, trig in rows:
        if not e:
            continue
        curve = traits.get(i)
        out.append({"i": i, "e": e, "aura": aura, "per": per, "mv": mv, "coef": coef,
                    "ap": ap, "trig": trig or 0,
                    "bps": curve if curve else [bp], "src": "curve" if curve else "spell"})
    return out


def main(cls):
    mask = CLASS_MASK[cls]
    db = sqlite3.connect(DB)
    db.row_factory = sqlite3.Row
    q = db.execute

    def one(sql, *a):
        r = q(sql, a).fetchone()
        return dict(r) if r else {}

    lines = {r["ID"]: r["DisplayName_lang"] for r in q(f"""
        select sl.ID, sl.DisplayName_lang, count(*) n
        from SkillLineAbility a join SkillLine sl on sl.ID = a.{A}{SL_LINE}
        where (a.{A}{SL_MASK} & ?) != 0 and sl.CategoryID = 7
        group by sl.ID having n >= 3""", (mask,)) if r["ID"] not in SKIP_LINES}

    # Rows with ClassMask 0 are not class-restricted, and the client does use some
    # of them: Mage Evocation, Hunter Aimed Shot, Paladin Judgement of Command,
    # Priest Shadow Word: Death and Druid Innervate all sit on their class's own
    # line with mask 0 and would otherwise be invisible. Most mask-0 rows are not
    # abilities though - talent rank spells, rune plumbing, pet passives - so they
    # are collected here and filtered once the talent list exists.
    # The line list itself is built only from masked rows, which keeps a class from
    # picking up another's spells where two share a line name (Protection is both
    # Warrior and Paladin; Holy is Paladin and Priest).
    rows = q(f"""select {A}{SL_LINE} sl, {A}{SL_SPELL} spell, {A}{SL_SUPER} supercedes,
                        {A}{SL_MASK} classmask, {A}{SL_RACE} racemask
                 from SkillLineAbility
                 where {A}{SL_LINE} in ({','.join(map(str, lines))})
                   and (({A}{SL_MASK} & ?) != 0 or {A}{SL_MASK} = 0)""",
             (mask,)).fetchall()

    # One spell can sit on several rows - a race-locked row beside an open one, or
    # the two faction copies of a city portal - so the restriction is the merge of
    # every row that grants it, not whichever row is read last.
    race_rows = {}
    for r in rows:
        race_rows.setdefault(r["spell"], []).append(r["racemask"])

    out = []
    for r in rows:
        sid = r["spell"]
        if is_sod(sid):
            continue
        sp = one("select NameSubtext_lang sub, Description_lang d, AuraDescription_lang ad from Spell where ID=?", sid)
        misc = one("select SchoolMask, SpellIconFileDataID icon from SpellMisc where SpellID=?", sid)
        lvl = one("select SpellLevel from SpellLevels where SpellID=?", sid)
        cd = one("select RecoveryTime, CategoryRecoveryTime, StartRecoveryTime from SpellCooldowns where SpellID=?", sid)
        cat = one("select Category from SpellCategories where SpellID=?", sid)
        # SpellCastTimes.Base is milliseconds; a negative base means the cast is the
        # ranged weapon's swing time (Arcane Shot, the Stings, Quick Shot).
        ct = one("""select ct.Base b from SpellMisc m
                    join SpellCastTimes ct on ct.ID = m.CastingTimeIndex
                    where m.SpellID=?""", sid)
        pw = one("select PowerType, ManaCost, PowerCostPct from SpellPower where SpellID=? order by OrderIndex", sid)
        ss = one("select ShapeshiftMask_0 m from SpellShapeshift where SpellID=?", sid)
        eq = one("select EquippedItemClass c, EquippedItemSubclass sc from SpellEquippedItems where SpellID=?", sid)
        ao = one("""select ProcTypeMask_0 p0, ProcTypeMask_1 p1, ProcChance pc,
                         ProcCategoryRecovery icd, ProcCharges chg, CumulativeAura stacks
                  from SpellAuraOptions where SpellID=?""", sid)
        rng = one("select r.RangeMax_1 mx from SpellRange r join SpellMisc m on m.RangeIndex=r.ID where m.SpellID=?", sid)
        dur = one("select d.Duration ms from SpellDuration d join SpellMisc m on m.DurationIndex=d.ID where m.SpellID=?", sid)
        effs = [dict(e) for e in q("""select EffectIndex i, Effect e, EffectAura aura, EffectBasePointsF bp,
                                      EffectBonusCoefficient coef, BonusCoefficientFromAP ap,
                                      EffectAuraPeriod per, EffectChainTargets ct, EffectTriggerSpell trig
                                      from SpellEffect where SpellID=? order by EffectIndex""", (sid,))]
        m = ss.get("m") or 0
        forms, seen = [], set()
        for b in range(32):
            if m >> b & 1:
                nm = FORM.get(b + 1, "form " + str(b + 1))
                if nm not in seen:
                    seen.add(nm)
                    forms.append(nm)
        # Rage is stored x10; mana, energy and focus are raw. A percentage of base
        # mana lives in PowerCostPct with ManaCost 0 (Judgement 6%, Multi-Shot 13.9%).
        ptype, raw = pw.get("PowerType"), pw.get("ManaCost") or 0
        power = raw / 10.0 if ptype == 1 else raw
        if power == int(power):
            power = int(power)
        out.append({
            "id": sid, "name": one("select Name_lang n from SpellName where ID=?", sid).get("n") or f"Spell {sid}",
            "rank": sp.get("sub") or "", "desc": sp.get("d") or "", "auraDesc": sp.get("ad") or "",
            "tree": lines[r["sl"]], "supercedes": r["supercedes"] or 0,
            "spellLevel": lvl.get("SpellLevel") or 0,
            "races": merge_races(race_rows.get(sid, [])),
            "school": school_name(misc.get("SchoolMask") or 1),
            "schoolMask": misc.get("SchoolMask") or 1,
            "schoolBits": [n for b, n in SCHOOL_BITS.items() if (misc.get("SchoolMask") or 1) & b],
            "power": power, "powerPct": round(pw.get("PowerCostPct") or 0, 1), "powerType": ptype,
            # Two different cooldowns: the spell's own, and the category one it shares
            # with its siblings. 715 class abilities use a category cooldown against
            # 227 with their own, so collapsing them into one number hides the
            # mechanic that actually governs most abilities.
            "cd": max(cd.get("RecoveryTime") or 0, cd.get("CategoryRecoveryTime") or 0),
            "cdOwn": cd.get("RecoveryTime") or 0,
            "cdCat": cd.get("CategoryRecoveryTime") or 0,
            "catId": cat.get("Category") or 0,
            "castTime": max(ct.get("b") or 0, 0),
            "castWeapon": (ct.get("b") or 0) < 0,
            "gcd": cd.get("StartRecoveryTime") or 0,
            "stances": forms, "range": rng.get("mx") or 0, "duration": dur.get("ms") or 0,
            "effects": effs, "iconFdid": misc.get("icon") or 0,
            # Weapon gate and proc mask: the difference between a flat passive and
            # something that triggers on an attack is not visible from the text.
            "weaponClass": eq.get("c", -1), "weaponSub": eq.get("sc") or 0,
            "attrFlags": attr_flags(q, sid),
            "procMask": (ao.get("p0") or 0) & 0xFFFFFFFF,
            "procMask2": (ao.get("p1") or 0) & 0xFFFFFFFF,
            "procChance": ao.get("pc") or 0,
            # ProcCategoryRecovery is the internal cooldown, in ms.
            "procIcd": ao.get("icd") or 0,
            "procCharges": ao.get("chg") or 0,
            "procStacks": ao.get("stacks") or 0,
            "maskZero": r["classmask"] == 0,
            "rankValues": {str(ei): v for ei, v in trait_points(q, sid).items()},
        })

    # --- talents, from the trait tree ---
    tabs = {r["ID"]: (r["Name_lang"], r["OrderIndex"]) for r in
            q("select ID,Name_lang,OrderIndex from TalentTab where ClassMask=?", (mask,))}
    spec_names = [n for n, _ in sorted(tabs.values(), key=lambda t: t[1])]
    tree = find_tree(q, list(tabs))
    talents = talents_from_tree(q, tree, spec_names, one) if tree else []

    # Keep only the mask-0 rows that are real abilities: named, with a genuine
    # learn level, and not a talent rank under another guise.
    talent_spell_ids = {s for t in talents for s in t["ranks"]}
    talent_names = {t["name"] for t in talents}
    out = [a for a in out if not a["maskZero"] or (
        a["name"] and not a["name"].startswith("Spell ")
        and a["spellLevel"] > 1
        and a["id"] not in talent_spell_ids
        and a["name"] not in talent_names)]
    for a in out:
        del a["maskZero"]

    # --- rank ladders: group by (line, name), then split parallel ladders ---
    def rank_no(a):
        m = re.search(r"(\d+)", a["rank"] or "")
        return int(m.group(1)) if m else 0

    groups = {}
    for a in out:
        groups.setdefault((a["tree"], a["name"]), []).append(a)
    fams = []
    for (tree_name, name), members in groups.items():
        numbered = [a for a in members if rank_no(a) > 0]
        unnumbered = [a for a in members if rank_no(a) == 0]
        ladders = []
        if numbered:
            buckets = {}
            for a in numbered:
                buckets.setdefault(rank_no(a), []).append(a)
            for b in buckets.values():
                b.sort(key=lambda a: a["id"])
            for j in range(max(len(b) for b in buckets.values())):
                picked = [b[j] for _, b in sorted(buckets.items()) if j < len(b)]
                if picked:
                    picked.sort(key=lambda a: (rank_no(a), a["spellLevel"], a["id"]))
                    ladders.append(picked)
        if unnumbered:
            unnumbered.sort(key=lambda a: (a["spellLevel"], a["id"]))
            ladders.append(unnumbered)
        for j, picked in enumerate(ladders):
            fams.append({"root": picked[0]["id"], "name": name, "tree": tree_name,
                         "ranks": [a["id"] for a in picked],
                         "variant": j + 1, "variants": len(ladders)})

    # --- icons, via the community listfile: FileDataID -> interface/icons/<name>.blp ---
    want = {a["iconFdid"] for a in out if a["iconFdid"]} | {t["iconFdid"] for t in talents if t["iconFdid"]}
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
    for item in list(out) + list(talents):
        item["icon"] = icon_by_fdid.get(item.pop("iconFdid"), "")

    # --- cooldown categories: who each ability shares its cooldown with ---
    by_cat = {}
    for a in out:
        if a["cdCat"] and a["catId"]:
            by_cat.setdefault(a["catId"], set()).add(a["name"])
    for a in out:
        peers = by_cat.get(a["catId"], set()) - {a["name"]} if a["cdCat"] else set()
        a["sharesWith"] = sorted(peers)

    # --- cross-links ---
    byid = {a["id"]: a for a in out}
    fam_of = {i: f for f in fams for i in f["ranks"]}
    talent_of = {}
    for t in talents:
        for sid in t["ranks"]:
            talent_of[sid] = {"name": t["name"], "tab": t["tab"], "tier": t["tier"],
                              "rank": 1, "maxRank": t["maxRank"], "id": t["id"]}
    extern = {}
    for a in out:
        trig = []
        for e in a["effects"]:
            ts = e.get("trig") or 0
            if ts and ts not in trig:
                trig.append(ts)
        a["triggers"] = trig
        for ts in trig:
            if ts not in byid and ts not in extern:
                extern[ts] = one("select Name_lang n from SpellName where ID=?", ts).get("n") or ""
        fam = fam_of.get(a["id"])
        a["siblings"] = [i for i in fam["ranks"] if i != a["id"]] if fam else []
        a["talent"] = talent_of.get(a["id"])
    for t in talents:
        for e in t["effects"]:
            ts = e.get("trig") or 0
            if ts and ts not in byid and ts not in extern:
                extern[ts] = one("select Name_lang n from SpellName where ID=?", ts).get("n") or ""

    # --- descriptions resolved by tools/foreverdesc (flat, and one per rank) ---
    resolved = _load(os.path.join(SP, "data", "descriptions.json"))
    ranked = _load(os.path.join(SP, "data", "descriptions-ranked.json"))
    for a in out:
        r = resolved.get(str(a["id"]))
        if r:
            a["descRaw"], a["desc"] = a["desc"], tidy(r)
        br = ranked.get(str(a["id"]))
        if br:
            a["descByRank"] = [tidy(x) if x else "" for x in br]
    for t in talents:
        key = str(t["nameFrom"])
        r = resolved.get(key)
        if r:
            t["descRaw"], t["desc"] = t["desc"], tidy(r)
        br = ranked.get(key)
        if br:
            t["descByRank"] = [tidy(x) if x else "" for x in br]
            if t["descByRank"][0]:
                t["desc"] = t["descByRank"][0]

    enums = json.load(open(os.path.join(SP, "data", "enums.json")))
    used_p = sorted({str(b) for a in out for b in range(32) if (a.get("procMask") or 0) >> b & 1})
    used_e = sorted({str(e["e"]) for a in out for e in a["effects"]} |
                    {str(e["e"]) for t in talents for e in t["effects"]})
    used_a = sorted({str(e["aura"]) for a in out for e in a["effects"]} |
                    {str(e["aura"]) for t in talents for e in t["effects"]})

    live_lines = {a["tree"] for a in out}
    chart = [n for n in sorted(live_lines) if n not in ("Engraving", "Runes")]
    chart = [n for n in chart if any(a["tree"] == n and a["spellLevel"] for a in out)]
    tree_order = sorted(live_lines, key=lambda n: (n in ("Engraving", "Runes"), n))
    forms_present = []
    for a in out:
        for f in a["stances"]:
            if f not in forms_present:
                forms_present.append(f)

    payload = {
        "meta": {"build": "1.60.1.69876", "product": "wow_classic_beta", "cls": cls,
                 "abilities": len(out), "families": len(fams), "talents": len(talents),
                 "traitTree": tree},
        "cfg": {"className": cls, "trees": tree_order, "chartTrees": chart,
                "forms": forms_present, "formWord": " stance" if cls == "Warrior" else " form",
                "talentTabs": spec_names, "notes": []},
        "abilities": out, "families": fams, "talents": talents, "extern": extern,
        "effectNames": {k: enums["effect"].get(k, "EFFECT_" + k) for k in used_e},
        "auraNames": {k: enums["aura"].get(k, "AURA_" + k) for k in used_a},
        "procNames": {k: enums.get("proc", {}).get(k, {"name": "bit " + k, "note": ""}) for k in used_p},
        "attrNames": attr_table(ATTR_FLAGS),
    }
    path = os.path.join(SP, "data", cls.lower() + ".json")
    json.dump(payload, open(path, "w"), separators=(",", ":"))
    print(f"{cls}: {len(out)} spells, {len(fams)} abilities, {len(talents)} talents "
          f"(trait tree {tree}), lines={tree_order}")


def _load(path):
    return json.load(open(path)) if os.path.exists(path) else {}


if __name__ == "__main__":
    main(sys.argv[1])
