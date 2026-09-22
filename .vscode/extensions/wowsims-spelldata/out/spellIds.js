"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.SPELL_ID_PATTERNS = void 0;
exports.spellIdAt = spellIdAt;
// Every way hand-written code names a spell id: the store's accessors and a ladder's ByID in Go,
// core.ActionID's SpellID, an APL file's "spellId" and the TS ActionId.fromSpellId. Each pattern
// captures the id in group 1.
exports.SPELL_ID_PATTERNS = [
    /\bMustFind\(\s*(\d+)\s*\)/g,
    /\bFind\(\s*(\d+)\s*\)/g,
    /\.ByID\(\s*(\d+)\s*\)/g,
    /\bSpellID:\s*(\d+)/g,
    /"spellId":\s*(\d+)/g,
    /\bfromSpellId\(\s*(\d+)\s*\)/g,
    /\bspellId:\s*(\d+)/g,
];
// The id the cursor sits in, anywhere within a match - on the accessor's name as readily as on the
// digits. The first pattern that covers the column wins.
function spellIdAt(lineText, column) {
    for (const pattern of exports.SPELL_ID_PATTERNS) {
        const scan = new RegExp(pattern.source, 'g');
        let match;
        while ((match = scan.exec(lineText)) !== null) {
            const start = match.index;
            const end = start + match[0].length;
            if (column >= start && column <= end) {
                return Number(match[1]);
            }
        }
    }
    return undefined;
}
//# sourceMappingURL=spellIds.js.map