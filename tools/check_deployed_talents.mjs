// Checks that the deployed site actually serves the talent data in this checkout.
//
//   node tools/check_deployed_talents.mjs [siteUrl]
//
// A green deploy workflow says the build ran, not that the numbers a player sees are the
// ones in the trees. This fetches a page, follows its bundle to the chunk holding the
// talent configs, and compares every rank against ui/core/talents/trees/*.json.
//
// Run it after a deploy, or when a fix is supposed to be live and something looks stale.
import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const site = (process.argv[2] ?? 'https://elliotwood.github.io/Forever/classic/').replace(/\/?$/, '/');
const treeDir = path.join(path.dirname(fileURLToPath(import.meta.url)), '..', 'ui', 'core', 'talents', 'trees');

async function text(url) {
  const response = await fetch(url);
  if (!response.ok) throw new Error(`${url} returned ${response.status}`);
  return response.text();
}

// The talent configs are bundled into a chunk the spec pages share, so any page reaches it.
const page = await text(site + 'enhancement_shaman/');
const entry = page.match(/src="([^"]*bundle[^"]*\.entry\.js)"/)?.[1];
if (!entry) {
  console.error('no entry bundle in the page; has the build output changed shape?');
  process.exit(1);
}

const entrySource = await text(new URL(entry, site).href);
const chunks = [...new Set([...entrySource.matchAll(/from"\.\.\/\.\.\/([^"]+\.chunk\.js)"/g)].map(m => m[1]))];

// More than one chunk mentions fieldName, so the talent chunk is the one that actually
// carries talent rank tables rather than the first that looks plausible.
let bundle = null;
for (const chunk of chunks) {
  const source = await text(site + 'bundle/' + chunk);
  if (/maxPoints:\d+,name:"/.test(source)) { bundle = source; break; }
}
if (!bundle) {
  console.error('no chunk carried the talent configs; has the build output changed shape?');
  process.exit(1);
}

// Talent names repeat across classes - Deflection belongs to both the rogue and the
// warrior - so the bundle is cut into its 27 trees first and each talent is filed under
// its own tree. The trees are delimited by backgroundUrl, which only a tree carries.
const treeStarts = [...bundle.matchAll(/backgroundUrl:/g)].map(m => m.index);
if (treeStarts.length === 0) {
  console.error('found no trees in the bundle; has the build output changed shape?');
  process.exit(1);
}

// Each tree's talents, keyed by the fieldName set so a tree can be recognised without
// relying on the order the bundler emitted them in.
const deployedTrees = treeStarts.map((start, i) => {
  const segment = bundle.slice(start, treeStarts[i + 1] ?? bundle.length);
  const talents = {};
  for (const match of segment.matchAll(/fieldName:"([^"]+)"((?:(?!fieldName:")[\s\S])*?)ranks:(\[\[[^\]]*\](?:,\[[^\]]*\])*\])/g)) {
    talents[match[1]] = match[3].replace(/\s+/g, '');
  }
  return talents;
});
if (!deployedTrees.some(t => Object.keys(t).length)) {
  console.error('parsed no talents out of the bundle; the minified shape has probably changed');
  process.exit(1);
}

// Rank tables hold numbers plus the occasional pluralisation string, so both sides are
// parsed and re-serialised rather than compared as text.
function canonical(ranks) {
  return JSON.stringify(JSON.parse(ranks.replace(/([[,])\./g, '$10.')));
}

let checked = 0;
const stale = [];
const unmatched = [];
for (const file of fs.readdirSync(treeDir).filter(f => f.endsWith('.json'))) {
  const className = file.replace('.json', '');
  for (const tree of JSON.parse(fs.readFileSync(path.join(treeDir, file), 'utf8'))) {
    // The deployed tree is the one holding this tree's talents, found by fieldName so
    // identically named talents in other classes cannot be mistaken for it.
    const fieldNames = tree.talents.map(t => t.fieldName).filter(Boolean);
    const deployed = deployedTrees.find(t => fieldNames.length && fieldNames.every(n => n in t));
    if (!deployed) {
      unmatched.push(`${className} ${tree.name}`);
      continue;
    }

    for (const talent of tree.talents) {
      const live = deployed[talent.fieldName];
      if (!live || !talent.ranks) continue;
      checked++;
      // The bundle is minified, so 0.4 is written .4 and 1.0 as 1. Comparing the parsed
      // numbers keeps the check on the values rather than how they are spelled.
      const local = canonical(JSON.stringify(talent.ranks));
      if (local !== canonical(live)) {
        stale.push(`${className} ${talent.name}\n    deployed: ${live}\n    local:    ${local}`);
      }
    }
  }
}

if (unmatched.length) {
  console.log(`could not find ${unmatched.length} tree(s) in the bundle: ${unmatched.join(', ')}`);
}

console.log(`compared ${checked} talents against ${site}`);
if (stale.length) {
  console.log(`\n${stale.length} differ from this checkout:\n`);
  for (const entry of stale) console.log('  ' + entry);
  console.log('\nEither the deploy has not run yet, or it did not pick these up.');
  process.exit(1);
}
console.log('every deployed talent matches this checkout');
