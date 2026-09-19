#!/usr/bin/env python3
"""Render one class's page from template.html + data/<class>.json.

    python3 tools/spellbook/build_page.py <Class> tools/spellbook/palette-<class>.json

Writes tools/spellbook/out/<class>-spellbook.html: one standalone file with the
data, the icons and all styling inlined, ready to publish or open locally.
"""
import json, os, sys

SP = os.path.dirname(os.path.abspath(__file__))
HEAD = """<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1, viewport-fit=cover">
<style>
  :root{color-scheme:dark light;padding-top:env(safe-area-inset-top,0px);padding-bottom:env(safe-area-inset-bottom,0px)}
  body{margin:0;font-family:system-ui,-apple-system,sans-serif;font-size:14px;background:#0e1216}
  img{max-width:100%}
  [hidden]{display:none!important}
</style>
"""


def main(cls, palette_path):
    data = json.load(open(os.path.join(SP, "data", cls.lower() + ".json")))
    pal = json.load(open(palette_path))
    # Tree colours are keyed by skill-line name, not position: which lines a class
    # has varies once the SoD id band is excluded (Engraving disappears entirely).
    trees = data["cfg"]["trees"]
    for key in ("treesLight", "treesDark"):
        missing = [t for t in trees if t not in pal[key]]
        if missing:
            sys.exit(f"{key} has no colour for {cls} lines: {missing}")

    data["cfg"]["notes"] = pal["notes"]
    icons_path = os.path.join(SP, "data", cls.lower() + "-icons.json")
    if not os.path.exists(icons_path):
        sys.exit(f"missing {icons_path} - run fetch_icons.py {cls} first")
    data["icons"] = json.load(open(icons_path))

    tpl = open(os.path.join(SP, "template.html"), encoding="utf-8").read()
    repl = {
        "__TITLE__": f"{cls} Spellbook 1.60",
        "__CLASS__": cls,
        "__ACCENT_L__": pal["accentLight"],
        "__ACCENT_SOFT_L__": pal["accentSoftLight"],
        "__ACCENT_D__": pal["accentDark"],
        "__ACCENT_D2__": pal["accentDarkInk"],
        "__ACCENT_SOFT_D__": pal["accentSoftDark"],
        "__FORM_LABEL__": pal["formLabel"],
        "__CHART_SUB__": pal["chartSub"],
        "__TREES_L__": " ".join(f"--t{i}:{pal['treesLight'][t]};" for i, t in enumerate(trees)),
        "__TREES_D__": " ".join(f"--t{i}:{pal['treesDark'][t]};" for i, t in enumerate(trees)),
    }
    for k, v in repl.items():
        assert k in tpl, "missing placeholder " + k
        tpl = tpl.replace(k, v)
    tpl = tpl.replace("__DATA__", json.dumps(data, separators=(",", ":")).replace("</", "<\\/"))

    cut = tpl.index('<header class="top">')
    out = HEAD + tpl[:cut] + "</head>\n<body>\n" + tpl[cut:] + "\n</body>\n</html>\n"
    path = os.path.join(SP, "out", cls.lower() + "-spellbook.html")
    open(path, "w", encoding="utf-8").write(out)
    m = data["meta"]
    print(f"{cls}: {m['abilities']} spells / {m['families']} abilities / {m['talents']} talents "
          f"-> {path} ({len(out.encode()) // 1024} KB)")


if __name__ == "__main__":
    main(sys.argv[1], sys.argv[2])
