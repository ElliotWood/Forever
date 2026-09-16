// Refreshes assets/confirmed_talents.json from talentsforever.com.
//
// That site marks a talent `complete` once every rank has been read off the game rather
// than extrapolated from rank 1. Those are the only ranks that can contradict the sim, so
// they are kept here and checked by TestConfirmedTalentRanksMatchTheSim.
//
//   node tools/refresh_confirmed_talents.mjs
//
// Numbers only: the check compares the values a rank states, not its wording, because the
// two sources phrase the same effect differently ("90 sec" against "1.5 min") often enough
// that comparing prose would be all noise.
import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const SOURCE = 'https://talentsforever.com/data.json';
const OUT = path.join(path.dirname(fileURLToPath(import.meta.url)), '..', 'assets', 'confirmed_talents.json');

const response = await fetch(SOURCE);
if (!response.ok) {
  console.error(`${SOURCE} returned ${response.status}`);
  process.exit(1);
}
const data = await response.json();

const talents = {};
let count = 0;
for (const [className, classData] of Object.entries(data.talents ?? {})) {
  for (const tree of classData.trees ?? []) {
    for (const talent of tree.talents ?? []) {
      if (talent.complete === false) continue;
      const descriptions = Array.isArray(talent.desc) ? talent.desc : Object.values(talent.desc ?? {});
      if (!descriptions.length) continue;
      (talents[className] ??= {})[talent.name] =
        descriptions.map(text => (text.match(/\d+(\.\d+)?/g) ?? []).map(Number));
      count++;
    }
  }
}

fs.writeFileSync(OUT, JSON.stringify({
  _readme: 'Confirmed talent rank values from talentsforever.com, kept so a test can check the sim has not drifted from an observation. Numbers only, in the order they appear in each rank tooltip. Regenerate with tools/refresh_confirmed_talents.mjs.',
  source: 'talentsforever.com (CC BY 4.0)',
  generated: data.generated,
  talents,
}, null, 1) + '\n');

console.log(`wrote ${count} confirmed talents, generated ${data.generated}`);
