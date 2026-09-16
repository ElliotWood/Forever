// Reports what the upstream talent sources say that this checkout does not.
//
//   node tools/check_talent_sources.mjs
//
// Two sources publish machine-readable Forever talent data:
//
//   talentsforever.com/data.json   marks a talent `complete` once every rank has been read
//                                  off the game rather than extrapolated from rank 1
//   nether.wowhead.com             the calculator's own data file
//
// A confirmed rank is the one kind of number this fork cannot argue with, so this prints
// any that disagree with ui/core/talents/trees/*.json, and says how many talents each
// source has confirmed so far. Nothing is written; the fix is a deliberate edit.
import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const root = path.join(path.dirname(fileURLToPath(import.meta.url)), '..');
const treeDir = path.join(root, 'ui', 'core', 'talents', 'trees');
const confirmedPath = path.join(root, 'assets', 'confirmed_talents.json');

const normalize = name => name.toLowerCase().replace(/[^a-z0-9]+/g, ' ').trim();
const numbersIn = text => (String(text).match(/\d+(\.\d+)?/g) ?? []).map(Number);

// The sim's own view: each talent's rank text with its placeholders filled in.
const sim = {};
for (const file of fs.readdirSync(treeDir).filter(f => f.endsWith('.json'))) {
  const className = file.replace('.json', '');
  for (const tree of JSON.parse(fs.readFileSync(path.join(treeDir, file), 'utf8'))) {
    for (const talent of tree.talents) {
      if (!talent.description || !talent.ranks) continue;
      sim[className + '|' + normalize(talent.name)] = talent.ranks.map(values => {
        let text = talent.description;
        values.forEach((v, i) => { text = text.split(`{${i}}`).join(String(v)); });
        return text;
      });
    }
  }
}

let response;
try {
  response = await fetch('https://talentsforever.com/data.json');
} catch (error) {
  console.error(`could not reach talentsforever.com: ${error.message}`);
  process.exitCode = 1;
  throw error;
}
if (!response.ok) {
  console.error(`talentsforever.com returned ${response.status}`);
  process.exitCode = 1;
  throw new Error('upstream unavailable');
}
const upstream = await response.json();

let confirmedTalents = 0, confirmedRanks = 0, comparedRanks = 0;
const disagreements = [];

for (const [className, classData] of Object.entries(upstream.talents ?? {})) {
  for (const tree of classData.trees ?? []) {
    for (const talent of tree.talents ?? []) {
      if (talent.complete === false) continue;
      confirmedTalents++;

      const descriptions = Array.isArray(talent.desc) ? talent.desc : Object.values(talent.desc ?? {});
      const simRanks = sim[className.toLowerCase() + '|' + normalize(talent.name)];
      for (let rank = 0; rank < descriptions.length; rank++) {
        confirmedRanks++;
        if (!simRanks || !simRanks[rank]) continue;
        comparedRanks++;

        const ours = numbersIn(simRanks[rank]);
        const theirs = numbersIn(descriptions[rank]);
        // The two write the same duration in different units, so a value that matches
        // once multiplied by 60 is the same number spelled differently.
        const same = ours.length === theirs.length &&
          [...ours].sort().every((v, i) => {
            const t = [...theirs].sort()[i];
            return v === t || v === t * 60 || t === v * 60;
          });
        if (!same) {
          disagreements.push({ className, name: talent.name, rank: rank + 1, ours: simRanks[rank], theirs: descriptions[rank] });
        }
      }
    }
  }
}

const stored = JSON.parse(fs.readFileSync(confirmedPath, 'utf8'));
const storedCount = Object.values(stored.talents ?? {}).reduce((n, t) => n + Object.keys(t).length, 0);

console.log(`talentsforever.com generated ${upstream.generated}`);
console.log(`  confirmed: ${confirmedTalents} talents, ${confirmedRanks} ranks`);
console.log(`  compared against the trees: ${comparedRanks} ranks`);
console.log(`  assets/confirmed_talents.json holds ${storedCount}` +
  (storedCount === confirmedTalents ? ' (current)' : ' — run tools/refresh_confirmed_talents.mjs'));

if (!disagreements.length) {
  console.log('\nevery confirmed rank agrees with this checkout');
} else {
  console.log(`\n${disagreements.length} confirmed rank(s) disagree with this checkout:\n`);
  for (const d of disagreements) {
    console.log(`  ${d.className} ${d.name} rank ${d.rank}`);
    console.log(`    sim:      ${d.ours}`);
    console.log(`    upstream: ${d.theirs}`);
  }
  // Set rather than exit, so the fetch keep-alive socket closes on its own. Calling
  // process.exit here trips an assertion in libuv on Windows.
  process.exitCode = 1;
}
