#!/usr/bin/env python3
"""Compare each class's spell ids between the SoD (wow_classic_era) and Forever
(wow_classic_beta) databases.

    python3 tools/spellbook/compare_sod.py [--json out.json]

Run from the repo root. Spell id is the identity here: matching on rank instead is
misleading, because Forever inserts new low ranks (Warrior Slam gains a rank 1, so
every older Slam id shifts down a rank while keeping its id).

The two builds address SkillLineAbility differently - the era build ships named
columns, the beta build has none for build 69876 and is addressed positionally.
The era column names confirm that mapping: 003 SkillLine, 004 Spell, 006 ClassMask.
"""
import json, sqlite3, sys
from collections import defaultdict

A = "Field_1_60_1_69876_"
CLASS_MASK = {"Warrior": 1, "Paladin": 2, "Hunter": 4, "Rogue": 8, "Priest": 16,
              "Shaman": 64, "Mage": 128, "Warlock": 256, "Druid": 1024}
SKIP_LINES = {183, 777}


def abilities(db, mask, named):
    """{spellId: (line, name, rank)} for one class."""
    line_col = "a.SkillLine" if named else f"a.{A}003"
    spell_col = "a.Spell" if named else f"a.{A}004"
    mask_col = "a.ClassMask" if named else f"a.{A}006"
    rows = db.execute(f"""
        select {spell_col}, sl.DisplayName_lang, coalesce(n.Name_lang,''), coalesce(s.NameSubtext_lang,'')
        from SkillLineAbility a
        join SkillLine sl on sl.ID = {line_col}
        left join SpellName n on n.ID = {spell_col}
        left join Spell s on s.ID = {spell_col}
        where ({mask_col} & ?) != 0 and sl.CategoryID = 7 and sl.ID not in ({','.join(map(str, SKIP_LINES))})
        """, (mask,)).fetchall()
    return {sid: (line, name, rank) for sid, line, name, rank in rows if name}


def main(json_path=None):
    sod = sqlite3.connect("tools/database/wowsims-sod.db")
    fvr = sqlite3.connect("tools/database/wowsims-forever.db")
    report = {}
    print(f"{'class':9}{'shared ids':>11}{'sod only':>10}{'forever only':>14}{'rank relabelled':>17}")
    for cls, mask in CLASS_MASK.items():
        a, b = abilities(sod, mask, True), abilities(fvr, mask, False)
        shared = set(a) & set(b)
        relabelled = [i for i in shared if a[i][2] != b[i][2]]
        only_sod = sorted(set(a) - set(b))
        only_fvr = sorted(set(b) - set(a))
        report[cls] = {
            "sharedIds": len(shared),
            "rankRelabelled": [{"spell": i, "name": a[i][1], "sodRank": a[i][2], "foreverRank": b[i][2]}
                               for i in sorted(relabelled)],
            "sodOnly": [{"spell": i, "line": a[i][0], "name": a[i][1], "rank": a[i][2]} for i in only_sod],
            "foreverOnly": [{"spell": i, "line": b[i][0], "name": b[i][1], "rank": b[i][2]} for i in only_fvr],
        }
        print(f"  {cls:8}{len(shared):>11}{len(only_sod):>10}{len(only_fvr):>14}{len(relabelled):>17}")

    tot_s = sum(len(v["sodOnly"]) for v in report.values())
    tot_f = sum(len(v["foreverOnly"]) for v in report.values())
    tot_r = sum(len(v["rankRelabelled"]) for v in report.values())
    print(f"\nshared ids: {sum(v['sharedIds'] for v in report.values())}"
          f"   sod-only: {tot_s}   forever-only: {tot_f}   rank relabelled: {tot_r}")
    if json_path:
        json.dump(report, open(json_path, "w"), indent=1)
        print(f"written -> {json_path}")
    return report


if __name__ == "__main__":
    main(sys.argv[2] if len(sys.argv) > 2 and sys.argv[1] == "--json" else None)
