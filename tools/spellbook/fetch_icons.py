#!/usr/bin/env python3
"""Download the Wowhead icons one class's page needs and inline them.

    python3 tools/spellbook/fetch_icons.py <Class>

The Artifact sandbox blocks external images, so a remote <img src> would silently
fail there; the icons ship inside the page as data URIs instead. They are 36x36
JPEGs of about 1 KB, and the icons/ cache is shared between classes.
"""
import base64, json, os, sys, time, urllib.error, urllib.request

SP = os.path.dirname(os.path.abspath(__file__))
CACHE = os.path.join(SP, "icons")
URL = "https://wow.zamimg.com/images/wow/icons/medium/{}.jpg"


def fetch(name):
    path = os.path.join(CACHE, name + ".jpg")
    if os.path.exists(path) and os.path.getsize(path) > 0:
        return open(path, "rb").read()
    req = urllib.request.Request(URL.format(name), headers={"User-Agent": "wowsims-spellbook/1.0"})
    for attempt in range(3):
        try:
            with urllib.request.urlopen(req, timeout=20) as r:
                data = r.read()
            open(path, "wb").write(data)
            return data
        except urllib.error.HTTPError as e:
            if e.code == 404:
                return None
            time.sleep(1 + attempt)
        except Exception:
            time.sleep(1 + attempt)
    return None


def main(cls):
    os.makedirs(CACHE, exist_ok=True)
    d = json.load(open(os.path.join(SP, "data", cls.lower() + ".json")))
    path = os.path.join(SP, "data", cls.lower() + "-icons.json")
    have = json.load(open(path)) if os.path.exists(path) else {}
    names = sorted({a["icon"] for a in d["abilities"] if a.get("icon")} |
                   {t["icon"] for t in d["talents"] if t.get("icon")})
    out, missing, fetched = {}, [], 0
    for n in names:
        if n in have:
            out[n] = have[n]
            continue
        raw = fetch(n)
        if raw is None:
            missing.append(n)
            continue
        out[n] = "data:image/jpeg;base64," + base64.b64encode(raw).decode()
        fetched += 1
    json.dump(out, open(path, "w"), separators=(",", ":"))
    print(f"{cls}: {len(out)}/{len(names)} icons ({fetched} newly downloaded)"
          + (f", missing: {missing}" if missing else ""))


if __name__ == "__main__":
    main(sys.argv[1])
