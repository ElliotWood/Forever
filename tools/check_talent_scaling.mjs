// Compares the per-point scaling in the Go against the rank values in the talent trees.
//
//   node tools/check_talent_scaling.mjs
//
// The trees decide what a player is shown; the Go decides what the sim computes. Nothing
// tied the two together, and they have drifted twice: Master of Defense paid 10 Rage at
// rank 2 where its tree read 5, and Shadow and Flame was held flat where its tree scaled.
// Both looked right on the site while the sim did something else.
//
// This reads every `<coefficient> * float64(x.Talents.<Name>)` in sim/ and checks that the
// series it produces appears among the tree's own rank values. It is deliberately loose:
// a talent's Go value often is not a number the tooltip prints (a multiplier against a
// base, a value in different units), so a mismatch is a question rather than a verdict.
// Findings are printed for a human to judge; nothing fails the build.
import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const root = path.join(path.dirname(fileURLToPath(import.meta.url)), '..');
const treeDir = path.join(root, 'ui', 'core', 'talents', 'trees');

// fieldName -> the numbers each rank of that talent prints, per class.
const treeRanks = {};
for (const file of fs.readdirSync(treeDir).filter(f => f.endsWith('.json'))) {
  const className = file.replace('.json', '');
  for (const tree of JSON.parse(fs.readFileSync(path.join(treeDir, file), 'utf8'))) {
    for (const talent of tree.talents) {
      if (!talent.fieldName || !talent.ranks) continue;
      treeRanks[className + '|' + talent.fieldName.toLowerCase()] =
        talent.ranks.map(values => values.filter(v => typeof v === 'number'));
    }
  }
}

// Go writes the proto field in PascalCase; the trees use camelCase.
const toFieldName = goName => goName.charAt(0).toLowerCase() + goName.slice(1);

const scalingRegex = /([\d.]+)\s*\*\s*float64\(\w+\.Talents\.([A-Za-z]+)\)/g;
const findings = [];
let checked = 0;

function walk(dir) {
  for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
    const full = path.join(dir, entry.name);
    if (entry.isDirectory()) { walk(full); continue; }
    if (!entry.name.endsWith('.go') || entry.name.endsWith('_test.go')) continue;

    // sim/<class>/... tells us which tree to compare against.
    const rel = path.relative(root, full).replace(/\\/g, '/');
    const className = rel.split('/')[1];
    if (!className || className === 'core' || className === 'common' || className === 'encounters') continue;

    const source = fs.readFileSync(full, 'utf8');
    source.split(/\r?\n/).forEach((line, i) => {
      for (const match of line.matchAll(scalingRegex)) {
        const coefficient = Number(match[1]);
        const key = className + '|' + toFieldName(match[2]).toLowerCase();
        const ranks = treeRanks[key];
        if (!ranks || !ranks.length) continue;
        checked++;

        // What the Go pays at each rank, against every number that rank prints.
        const mismatches = [];
        for (let rank = 0; rank < ranks.length; rank++) {
          const goValue = coefficient * (rank + 1);
          const printed = ranks[rank];
          if (!printed.length) continue;
          // Accept the value itself, or it scaled by 100 either way - a tooltip says 2%
          // where the Go carries 0.02.
          const found = printed.some(v => Math.abs(v - goValue) < 0.005 ||
            Math.abs(v - goValue * 100) < 0.005 || Math.abs(v * 100 - goValue) < 0.005);
          if (!found) mismatches.push(`rank ${rank + 1}: go ${goValue}, tree ${printed.join('/')}`);
        }
        // Every rank disagreeing usually means the Go value is simply a different quantity
        // from anything the tooltip prints. A partial disagreement is the interesting case:
        // the two agree on some ranks and not others, which is what drift looks like.
        if (mismatches.length && mismatches.length < ranks.length) {
          findings.push({ file: rel, line: i + 1, talent: match[2], coefficient, mismatches });
        }
      }
    });
  }
}

walk(path.join(root, 'sim'));

console.log(`checked ${checked} per-point scaling sites against the trees`);
if (!findings.length) {
  console.log('none disagree with their tree on some ranks but not others');
} else {
  console.log(`\n${findings.length} scale differently from their tree on some ranks:\n`);
  for (const f of findings) {
    console.log(`  ${f.file}:${f.line}  ${f.talent} (${f.coefficient} per point)`);
    for (const m of f.mismatches) console.log(`      ${m}`);
  }
  console.log('\nEach is a question, not a verdict: a Go value is often not a number the');
  console.log('tooltip prints. Check the tooltip before changing anything.');
}
