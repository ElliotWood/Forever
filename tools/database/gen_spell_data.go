package database

import (
	"database/sql"
	"errors"
	"fmt"
	"go/format"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/wowsims/forever/tools/database/dbc"
)

var rankSubtext = regexp.MustCompile(`^Rank (\d+)$`)

type generatedRow struct {
	Rank            int32
	SpellID         int32
	Cost            int32
	CastTimeMs      int32
	GCDMs           int32
	CooldownMs      int32
	MinRange        float64
	MaxRange        float64
	MissileSpeed    float64
	ProcChance      int32
	FlatThreatBonus float64
	SchoolMask      int32
	DefenseType     int32
	Effects         []generatedEffect
	Direct          *generatedAmount
	Heal            *generatedAmount
	Periodic        *generatedAmount
	Energize        *generatedAmount
}

type generatedEffect struct {
	Index    int32
	Effect   dbc.SpellEffectType
	Aura     dbc.EffectAuraType
	Misc     int32
	Value    float64
	ValueMax float64
}

type generatedAmount struct {
	Min    float64
	Max    float64
	Coef   float64
	APCoef float64

	PeriodMs int32
	Ticks    int32
}

type rankCandidate struct {
	SpellID   int32
	Rank      int32
	ClassMask int
	SkillLine int32
}

type rankLadder struct {
	Name  string
	Field string
	Ranks map[int32]int32

	// What each rank's effects are worth, keyed by rank and then by the client's EffectIndex. Set
	// only on a ladder the talent tree supplied, where every rank is the same spell and the numbers
	// sit on the tree's curves rather than on a spell per rank. Nil means the rank's own spell states
	// its numbers, which is what every "Rank N" family does.
	Points map[int32]map[int32]float64
}

// One talent as the tree states it: a single spell, the rank cap the game enforces, and the value
// each rank gives. The rank spells the legacy Talent table still lists are gone from this client -
// Anticipation's 12750-12753 have no row in SpellName, SpellEffect or SpellMisc - so the tree is the
// only place a talent's ranks are described.
type traitLadder struct {
	SpellID  int32
	MaxRanks int32
	Points   map[int32]map[int32]float64
}

// SkillLineAbility.ClassMask is a bitmask over the class index dbc.Classes already carries, so the bit
// is that index shifted rather than a second table to keep in step with it.
func classMaskOf(class dbc.DbcClass) int {
	return 1 << (class.ID - 1)
}

// The Go identifier a family is reached by: "Shadow Word: Pain" -> ShadowWordPain.
func fieldNameOf(spellName string) string {
	var b strings.Builder
	upper := true
	for _, r := range spellName {
		switch {
		case r == '\'' || r == '\u2019':
			// Dropped without breaking the word, so "Avenger's Shield" is AvengersShield rather than
			// AvengerSShield.
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			if upper {
				b.WriteRune(unicode.ToUpper(r))
				upper = false
			} else {
				b.WriteRune(r)
			}
		default:
			upper = true
		}
	}

	name := b.String()
	if name == "" || unicode.IsDigit(rune(name[0])) {
		return ""
	}
	return name
}

// Every multi-rank family a class can learn, from two sources: every spell whose subtext reads
// "Rank N" in one of the class's skill lines, and every talent in the class's tree. Nothing
// hand-maintained.
func discoverLadders(db *sql.DB, class dbc.DbcClass, treeID int) ([]rankLadder, []string, []string, error) {
	mask := classMaskOf(class)

	exclusive, err := exclusiveSkillLines(db, mask)
	if err != nil {
		return nil, nil, nil, err
	}

	rows, err := db.Query(`
		SELECT n.Name_lang, sla.Spell, s.NameSubtext_lang, sla.ClassMask, sla.SkillLine
		FROM SkillLineAbility sla
		JOIN SpellName n ON n.ID = sla.Spell
		JOIN Spell s ON s.ID = sla.Spell
		WHERE sla.SkillLine IN (
			SELECT DISTINCT sla2.SkillLine
			FROM SkillLineAbility sla2
			JOIN SkillLine sl2 ON sl2.ID = sla2.SkillLine AND sl2.CategoryID = 7
			WHERE (sla2.ClassMask & ?) != 0
		)
		AND s.NameSubtext_lang LIKE 'Rank %'
		ORDER BY n.Name_lang, sla.Spell`, mask)
	if err != nil {
		return nil, nil, nil, err
	}
	defer rows.Close()

	byName := map[string]map[int32][]rankCandidate{}
	for rows.Next() {
		var name, subtext string
		var c rankCandidate
		if err := rows.Scan(&name, &c.SpellID, &subtext, &c.ClassMask, &c.SkillLine); err != nil {
			return nil, nil, nil, err
		}
		m := rankSubtext.FindStringSubmatch(subtext)
		if m == nil {
			continue
		}
		rank, _ := strconv.Atoi(m[1])
		c.Rank = int32(rank)
		if byName[name] == nil {
			byName[name] = map[int32][]rankCandidate{}
		}
		byName[name][c.Rank] = append(byName[name][c.Rank], c)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, nil, err
	}

	traits, traitSkipped, traitPartial, err := discoverTraitLadders(db, treeID)
	if err != nil {
		return nil, nil, nil, err
	}

	seen := map[string]bool{}
	var names []string
	for _, source := range []map[string]bool{keysOf(byName), keysOf(traits), keysOf(traitSkipped)} {
		for name := range source {
			if !seen[name] {
				seen[name] = true
				names = append(names, name)
			}
		}
	}
	sort.Strings(names)

	var ladders []rankLadder
	var skipped, partial []string
	taken := map[string]string{}

	for _, name := range names {
		field := fieldNameOf(name)
		if field == "" {
			skipped = append(skipped, fmt.Sprintf("%s: no usable Go identifier", name))
			continue
		}
		if owner, clash := taken[field]; clash {
			skipped = append(skipped, fmt.Sprintf("%s: identifier %s already taken by %s", name, field, owner))
			continue
		}

		var ladder map[int32]int32
		var points map[int32]map[int32]float64
		var resolveErr error
		if claimedByClass(byName[name], mask, exclusive) {
			ladder, resolveErr = resolveLadder(db, name, byName[name], mask)
		}

		// The tree wins where it states more ranks than the "Rank N" spells resolved to, because the
		// rank cap it carries is the one the talent proto and the UI trees are built from: Precision
		// kept its subtext on rank 1 alone, and leaving it as that one rank would let a 3-point
		// talent be read past its own ladder. Below that the "Rank N" spells win - they state a cost
		// and a cast time per rank, which the tree does not - and a family the resolver refused
		// stays refused, since a tree node that grants an ability does not describe its ranks.
		if t, ok := traits[name]; ok && resolveErr == nil && int(t.MaxRanks) > len(ladder) {
			ladder, points = map[int32]int32{}, t.Points
			for rank := int32(1); rank <= t.MaxRanks; rank++ {
				ladder[rank] = t.SpellID
			}
			if note := traitPartial[name]; note != "" {
				partial = append(partial, fmt.Sprintf("%s: %s", name, note))
			}
		}

		if resolveErr != nil {
			skipped = append(skipped, fmt.Sprintf("%s: %s", name, resolveErr))
			continue
		}
		if len(ladder) == 0 {
			if reason := traitSkipped[name]; reason != "" {
				skipped = append(skipped, fmt.Sprintf("%s: %s", name, reason))
			}
			continue
		}

		taken[field] = name
		ladders = append(ladders, rankLadder{Name: name, Field: field, Ranks: ladder, Points: points})
	}

	return ladders, skipped, partial, nil
}

func keysOf[V any](m map[string]V) map[string]bool {
	keys := make(map[string]bool, len(m))
	for k := range m {
		keys[k] = true
	}
	return keys
}

// Every talent in one class's tree, by name, plus the ones that had to be left out and the ones that
// came out short an effect. The tree states one spell per talent and a curve per effect, where the
// curve reads rank -> value, so a rank's numbers are the curve's and never the spell's own: those
// hold the top rank on some talents and a stale number on others, and 240 of the 640 curves in this
// build disagree with them.
func discoverTraitLadders(db *sql.DB, treeID int) (map[string]traitLadder, map[string]string, map[string]string, error) {
	rows, err := db.Query(`
		SELECT DISTINCT d.ID, COALESCE(d.SpellID, 0), e.MaxRanks,
		       COALESCE(NULLIF(sn.Name_lang, ''), d.OverrideName_lang, '')
		FROM TraitNode tn
		JOIN TraitNodeXTraitNodeEntry x ON x.TraitNodeID = tn.ID
		JOIN TraitNodeEntry e ON e.ID = x.TraitNodeEntryID
		JOIN TraitDefinition d ON d.ID = e.TraitDefinitionID
		LEFT JOIN SpellName sn ON sn.ID = d.SpellID
		WHERE tn.TraitTreeID = ?
		ORDER BY d.ID`, treeID)
	if err != nil {
		return nil, nil, nil, err
	}
	defer rows.Close()

	type traitDef struct {
		ID       int32
		SpellID  int32
		MaxRanks int32
		Name     string
	}
	var defs []traitDef
	for rows.Next() {
		var d traitDef
		if err := rows.Scan(&d.ID, &d.SpellID, &d.MaxRanks, &d.Name); err != nil {
			return nil, nil, nil, err
		}
		defs = append(defs, d)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, nil, err
	}

	ladders := map[string]traitLadder{}
	skipped := map[string]string{}
	partial := map[string]string{}
	spellOf := map[string]int32{}

	for _, d := range defs {
		if d.Name == "" || d.SpellID == 0 {
			continue
		}
		// A node appears once per entry, and the client keeps retired twins around, so a second
		// definition on the same spell is the same talent and a second one on another spell is two
		// talents the generated file could only reach by one name.
		if first, ok := spellOf[d.Name]; ok {
			if first != d.SpellID {
				delete(ladders, d.Name)
				skipped[d.Name] = fmt.Sprintf(
					"two talent tree nodes carry the name, on spells %d and %d", first, d.SpellID)
			}
			continue
		}
		spellOf[d.Name] = d.SpellID

		// A one-rank node grants an ability rather than describing a ladder - Hemorrhage and Water
		// Shield are nodes on a spell the game teaches - and the one row it would yield names the
		// node's spell, which is not the spell the ability's own ranks are keyed on.
		if d.MaxRanks <= 1 {
			continue
		}

		effects, err := effectIndicesOf(db, d.SpellID)
		if err != nil {
			return nil, nil, nil, err
		}
		if len(effects) == 0 {
			skipped[d.Name] = fmt.Sprintf("the client states no effects for spell %d", d.SpellID)
			continue
		}

		curves, err := traitCurves(db, d.ID)
		if err != nil {
			return nil, nil, nil, err
		}

		points := map[int32]map[int32]float64{}
		var uncovered []string
		for _, index := range effects {
			values, ok := curveRanks(curves[index], d.MaxRanks)
			if !ok {
				uncovered = append(uncovered, strconv.Itoa(int(index)))
				continue
			}
			for rank, value := range values {
				if points[rank] == nil {
					points[rank] = map[int32]float64{}
				}
				points[rank][index] = value
			}
		}

		if len(points) == 0 {
			skipped[d.Name] = fmt.Sprintf("the talent tree states no per-rank value for spell %d", d.SpellID)
			continue
		}
		if len(uncovered) > 0 {
			partial[d.Name] = fmt.Sprintf("effect %s of spell %d has no rank curve and is left out",
				strings.Join(uncovered, ", "), d.SpellID)
		}
		ladders[d.Name] = traitLadder{SpellID: d.SpellID, MaxRanks: d.MaxRanks, Points: points}
	}

	return ladders, skipped, partial, nil
}

func effectIndicesOf(db *sql.DB, spellID int32) ([]int32, error) {
	rows, err := db.Query(`SELECT EffectIndex FROM SpellEffect WHERE SpellID = ? ORDER BY EffectIndex`, spellID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var indices []int32
	for rows.Next() {
		var index int32
		if err := rows.Scan(&index); err != nil {
			return nil, err
		}
		indices = append(indices, index)
	}
	return indices, rows.Err()
}

// Each curve as rank -> value. OperationType is 0 on all 640 rows in this build, meaning the point is
// the value the rank is set to; anything else would modify the spell's own number instead, and is
// dropped rather than guessed at - the effect then reads as having no curve.
func traitCurves(db *sql.DB, defID int32) (map[int32]map[int32]float64, error) {
	rows, err := db.Query(`
		SELECT p.EffectIndex, p.OperationType, cp.Pos_0, cp.Pos_1
		FROM TraitDefinitionEffectPoints p
		JOIN CurvePoint cp ON cp.CurveID = p.CurveID
		WHERE p.TraitDefinitionID = ?
		ORDER BY p.EffectIndex, cp.Pos_0`, defID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	curves := map[int32]map[int32]float64{}
	for rows.Next() {
		var index, operation int32
		var at, value float64
		if err := rows.Scan(&index, &operation, &at, &value); err != nil {
			return nil, err
		}
		rank := int32(at)
		if operation != 0 || float64(rank) != at {
			continue
		}
		if curves[index] == nil {
			curves[index] = map[int32]float64{}
		}
		curves[index][rank] = value
	}
	return curves, rows.Err()
}

func curveRanks(curve map[int32]float64, maxRanks int32) (map[int32]float64, bool) {
	values := make(map[int32]float64, maxRanks)
	for rank := int32(1); rank <= maxRanks; rank++ {
		value, ok := curve[rank]
		if !ok {
			return nil, false
		}
		values[rank] = value
	}
	return values, true
}

// The skill lines only this class appears in. 11 are shared: "Holy" carries both the paladin and the
// priest bit, which is how paladin Holy Shock reached the priest file. Inside an exclusive line a
// ClassMask-0 spell must be this class's, which is what makes Ignite attributable.
func exclusiveSkillLines(db *sql.DB, mask int) (map[int32]bool, error) {
	rows, err := db.Query(`
		SELECT SkillLine, group_concat(DISTINCT ClassMask)
		FROM SkillLineAbility WHERE ClassMask != 0 GROUP BY SkillLine`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lines := map[int32]bool{}
	for rows.Next() {
		var line int32
		var masks string
		if err := rows.Scan(&line, &masks); err != nil {
			return nil, err
		}
		union := 0
		for _, part := range strings.Split(masks, ",") {
			m, err := strconv.Atoi(part)
			if err != nil {
				continue
			}
			union |= m
		}
		if union == mask {
			lines[line] = true
		}
	}
	return lines, rows.Err()
}

// A family is this class's if one of its ranks carries the class bit, or if its ranks sit in a skill
// line no other class appears in.
func claimedByClass(byRank map[int32][]rankCandidate, mask int, exclusive map[int32]bool) bool {
	for _, cands := range byRank {
		for _, c := range cands {
			if c.ClassMask&mask != 0 || exclusive[c.SkillLine] {
				return true
			}
		}
	}
	return false
}

// Picks one spell per rank. Lightning Bolt's Elemental Overload twins (45284-45293) share the name,
// skill line and class set, differing only by ClassMask 0. So the class bit wins where it exists, and
// a ClassMask-0 candidate is taken only when no real one was found - that is how Holy Shield 1-3
// resolve.
// Where no rule can separate two candidates and none should try: Seal of Righteousness rank 1 is
// 20154 and 21084, same name, rank, effect shape and 20 mana. The sim has always used 21084.
var ladderPins = map[string]map[int32]int32{
	"Seal of Righteousness": {1: 21084},
}

func resolveLadder(db *sql.DB, name string, byRank map[int32][]rankCandidate, mask int) (map[int32]int32, error) {
	ladder := map[int32]int32{}

	// Sorted, so that when several ranks are ambiguous the one named in the skipped list is always the
	// lowest rather than whichever the map handed over first.
	rankNums := make([]int32, 0, len(byRank))
	for rank := range byRank {
		rankNums = append(rankNums, rank)
	}
	sort.Slice(rankNums, func(i, j int) bool { return rankNums[i] < rankNums[j] })

	var ambiguous []string
	for _, rank := range rankNums {
		cands := byRank[rank]
		var chosen []rankCandidate
		for _, c := range cands {
			if c.ClassMask&mask != 0 {
				chosen = append(chosen, c)
			}
		}
		if len(chosen) == 0 {
			chosen = cands
		}

		ids := map[int32]bool{}
		for _, c := range chosen {
			ids[c.SpellID] = true
		}

		if len(ids) == 1 {
			ladder[rank] = chosen[0].SpellID
			continue
		}

		// Holy Shock is one name over three spells per rank: a dummy the player casts, plus a damage
		// and a heal spell the client never exposes. Only the castable one carries a SpellPower row,
		// and it is the one the sim registers, so that is the tie-break.
		castable, err := castableOf(db, ids)
		if err != nil {
			return nil, err
		}
		if castable == 0 {
			// Judgement of Command is the same dispatcher shape with no mana on either half - the
			// parent Judgement pays - so SpellPower cannot separate them. One E_DUMMY against one
			// E_SCHOOL_DAMAGE can: the dummy is what the seal triggers and what the sim registers,
			// and the damage spell rides along as a sibling, exactly as Holy Shock does.
			castable, err = dispatcherOf(db, ids)
			if err != nil {
				return nil, err
			}
		}
		if pinned, ok := ladderPins[name][rank]; ok && ids[pinned] {
			ladder[rank] = pinned
			continue
		}
		if castable == 0 {
			// Every ambiguous rank is collected rather than returning on the first, so the skipped
			// list in the class file names all of them. Judgement of Command is ambiguous on all six.
			var list []string
			for id := range ids {
				list = append(list, strconv.Itoa(int(id)))
			}
			sort.Strings(list)
			ambiguous = append(ambiguous, fmt.Sprintf("rank %d between spells %s", rank, strings.Join(list, ", ")))
			continue
		}
		ladder[rank] = castable
	}
	if len(ambiguous) > 0 {
		return nil, fmt.Errorf("ambiguous %s", strings.Join(ambiguous, "; "))
	}

	maxRank := int32(0)
	for rank := range ladder {
		if rank > maxRank {
			maxRank = rank
		}
	}
	for rank := int32(1); rank <= maxRank; rank++ {
		if _, ok := ladder[rank]; !ok {
			return nil, fmt.Errorf("missing rank %d of %d", rank, maxRank)
		}
	}
	return ladder, nil
}

// The one spell among these whose first effect is a dummy, when exactly one other is a direct
// damage effect. That pairing is the client's dispatcher shape: the dummy is the spell the game
// grants and triggers, the damage spell is never exposed. Anything else returns 0 rather than
// guessing - two talent auras that merely differ somewhere are not a dispatcher.
func dispatcherOf(db *sql.DB, ids map[int32]bool) (int32, error) {
	var dummy, damage int32
	for id := range ids {
		var effect dbc.SpellEffectType
		err := db.QueryRow(
			`SELECT Effect FROM SpellEffect WHERE SpellID = ? ORDER BY EffectIndex LIMIT 1`, id).Scan(&effect)
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		if err != nil {
			return 0, err
		}
		switch effect {
		case dbc.E_DUMMY:
			if dummy != 0 {
				return 0, nil
			}
			dummy = id
		case dbc.E_SCHOOL_DAMAGE:
			if damage != 0 {
				return 0, nil
			}
			damage = id
		default:
			return 0, nil
		}
	}
	if dummy == 0 || damage == 0 {
		return 0, nil
	}
	return dummy, nil
}

// The one spell among these that has a mana cost, or 0 when that does not single one out.
func castableOf(db *sql.DB, ids map[int32]bool) (int32, error) {
	var found int32
	for id := range ids {
		var n int
		if err := db.QueryRow(`SELECT count(*) FROM SpellPower WHERE SpellID = ?`, id).Scan(&n); err != nil {
			return 0, err
		}
		if n == 0 {
			continue
		}
		if found != 0 {
			return 0, nil
		}
		found = id
	}
	return found, nil
}

// Which effect supplies which field. Derived from the effect types rather than declared per family,
// because the client data already says it: a SCHOOL_DAMAGE effect is direct damage, a HEAL effect is a
// heal, an ENERGIZE effect is Lay on Hands' mana restore, a periodic aura is a tick.
// Where points is set the rank's numbers come from the talent tree's curves, and only the effects it
// prices are kept: the rest would have to be read off a spell that states one rank for all of them.
// The same-name sibling search is not wanted there either - it exists for the ability dispatchers,
// and a talent priced per rank has nothing to borrow.
func buildRow(db *sql.DB, rank int32, spellID int32, mask int, points map[int32]float64) (generatedRow, error) {
	var spell RankSpell
	var candidates []RankEffect
	var err error
	if points != nil {
		spell, err = LoadRankSpell(db, spellID)
		candidates = spell.Effects
	} else {
		spell, candidates, err = RankCandidates(db, spellID, mask)
	}
	if err != nil {
		return generatedRow{}, err
	}
	if points != nil {
		spell.Effects = pricedEffects(spell.Effects, points)
		candidates = spell.Effects
	}

	derive := func(e RankEffect) (float64, float64) {
		if value, ok := points[e.Index]; ok {
			return value, value
		}
		return DeriveRankAmount(e, spell.SpellLevel, spell.MaxLevel)
	}

	row := generatedRow{
		Rank: rank, SpellID: spellID,
		CastTimeMs: spell.CastTimeMs, GCDMs: spell.GCDMs, CooldownMs: spell.CooldownMs,
		MinRange: spell.MinRange, MaxRange: spell.MaxRange, MissileSpeed: spell.MissileSpeed,
		ProcChance: spell.ProcChance, SchoolMask: spell.SchoolMask, DefenseType: spell.DefenseType,
	}
	if spell.ManaCost.Valid {
		row.Cost = NormalizePowerCost(int32(spell.ManaCost.Int64), spell.PowerType)
	}

	amountOf := func(e RankEffect) *generatedAmount {
		min, max := derive(e)
		a := &generatedAmount{Min: min, Max: max, Coef: e.Coefficient, APCoef: e.APCoef}

		// A tick count is the duration over the period, which is how every hand-written
		// NumberOfTicks in the sim was arrived at.
		if e.AuraPeriod > 0 {
			a.PeriodMs = e.AuraPeriod
			if spell.DurationMs > 0 {
				a.Ticks = spell.DurationMs / e.AuraPeriod
			}
		}
		return a
	}

	for _, e := range spell.Effects {
		min, max := derive(e)
		if max == min {
			max = 0
		}
		row.Effects = append(row.Effects, generatedEffect{
			Index: e.Index, Effect: e.Effect, Aura: e.Aura, Misc: e.MiscValue, Value: min, ValueMax: max,
		})
	}

	for _, e := range candidates {
		switch {
		case (e.Effect == dbc.E_SCHOOL_DAMAGE || IsWeaponDamageEffect(e.Effect)) && row.Direct == nil:
			row.Direct = amountOf(e)
		case e.Effect == dbc.E_HEAL && row.Heal == nil:
			row.Heal = amountOf(e)
		case e.Effect == dbc.E_ENERGIZE && row.Energize == nil:
			row.Energize = amountOf(e)
		case IsThreatEffect(e.Effect) && row.FlatThreatBonus == 0:
			min, _ := derive(e)
			row.FlatThreatBonus = min
		case IsPeriodicAura(e.Aura) && row.Periodic == nil:
			row.Periodic = amountOf(e)
		}
	}

	// An aura effect reached by the fallbacks below still lands in the role its shape says it has:
	// Frenzied Regeneration's aura ticks every second, and filing that under Direct would hand the
	// call site a periodic value through a field that promises a direct one.
	fallback := func(e RankEffect) {
		if e.AuraPeriod > 0 {
			row.Periodic = amountOf(e)
		} else {
			row.Direct = amountOf(e)
		}
	}

	// Holy Shield keeps its per-block damage on an aura effect that is none of the roles above, and it
	// is not the only aura effect on the spell: one index holds the block value and another the damage.
	// The damage is the one that scales with spell power, so a nonzero coefficient is what picks it.
	if !row.hasValue() {
		for _, e := range candidates {
			if e.Aura != 0 && e.BasePoints > 0 && e.Coefficient > 0 {
				fallback(e)
				break
			}
		}
	}
	if !row.hasValue() {
		for _, e := range candidates {
			if e.Aura != 0 && e.BasePoints > 0 {
				fallback(e)
				break
			}
		}
	}

	return row, nil
}

func pricedEffects(effects []RankEffect, points map[int32]float64) []RankEffect {
	var kept []RankEffect
	for _, e := range effects {
		if _, ok := points[e.Index]; ok {
			kept = append(kept, e)
		}
	}
	return kept
}

// A low rank can legitimately carry no numbers at all - Lay on Hands rank 1 heals a share of max health
// and restores no mana, so it has no ENERGIZE effect where ranks 2-4 do.
func (row generatedRow) hasValue() bool {
	return row.Direct != nil || row.Heal != nil || row.Periodic != nil || row.Energize != nil
}

func GenerateSpellDataFiles(helper *DBHelper) error {
	if err := RequireSpellCastTimes(helper.db); err != nil {
		return err
	}

	// Rendered in full before anything is written, so a class that fails validation cannot leave half
	// the packages regenerated and half stale.
	namer, err := newRankEnumNamer()
	if err != nil {
		return err
	}

	// One tree per class, picked the same way the talent protos pick it, so the rank caps here and the
	// ones the sim's Talents message carries are the same numbers.
	trees, err := selectTraitTrees(helper)
	if err != nil {
		return err
	}

	rendered := map[string][]byte{}
	for _, class := range dbc.Classes {
		pkg := strings.ToLower(dbc.ClassNameFromDBC(class))
		tree, ok := trees[classMaskOf(class)]
		if !ok {
			return fmt.Errorf("%s: the client database holds no talent tree for this class", pkg)
		}
		out, err := renderClassFile(helper.db, pkg, class, namer, tree)
		if err != nil {
			return fmt.Errorf("%s: %w", pkg, err)
		}
		rendered[pkg] = out
	}

	// Rendered before any write too: it holds exactly the names the class files above turned out to
	// reference, and a class file naming a constant this file does not declare breaks the sim - and
	// with it gen_db, which imports the sim.
	enums, err := namer.render()
	if err != nil {
		return err
	}

	for pkg, out := range rendered {
		if err := os.WriteFile(fmt.Sprintf("sim/%s/spell_data_auto_gen.go", pkg), out, 0644); err != nil {
			return err
		}
	}
	return os.WriteFile("sim/common/shared/spell_data_enums_auto_gen.go", enums, 0644)
}

func renderClassFile(db *sql.DB, pkg string, class dbc.DbcClass, namer *rankEnumNamer, treeID int) ([]byte, error) {
	ladders, skipped, partial, err := discoverLadders(db, class, treeID)
	if err != nil {
		return nil, err
	}

	// Named rather than dropped silently, so a family the resolver could not make sense of is visible
	// here instead of merely absent. Kept out of the body below, whose text decides which imports the
	// file needs - a family name containing "time." would otherwise add an unused one.
	var notGenerated strings.Builder
	for _, block := range []struct {
		header string
		lines  []string
	}{
		{"// Not generated:\n", skipped},
		{"// Generated without an effect the talent tree states no rank value for:\n", partial},
	} {
		if len(block.lines) == 0 {
			continue
		}
		notGenerated.WriteString(block.header)
		for _, s := range block.lines {
			fmt.Fprintf(&notGenerated, "//   %s\n", s)
		}
		notGenerated.WriteString("\n")
	}

	var b strings.Builder
	b.WriteString("type generatedSpellData struct {\n")
	for _, l := range ladders {
		fmt.Fprintf(&b, "\t%s shared.SpellDataTable\n", l.Field)
	}
	b.WriteString("}\n\nvar spellData = generatedSpellData{\n")

	mask := classMaskOf(class)
	for _, l := range ladders {
		ranks := make([]int32, 0, len(l.Ranks))
		for rank := range l.Ranks {
			ranks = append(ranks, rank)
		}
		sort.Slice(ranks, func(i, j int) bool { return ranks[i] < ranks[j] })

		fmt.Fprintf(&b, "\t%s: shared.SpellDataTable{\n", l.Field)
		for _, rank := range ranks {
			row, err := buildRow(db, rank, l.Ranks[rank], mask, l.Points[rank])
			if err != nil {
				return nil, fmt.Errorf("%s rank %d: %w", l.Name, rank, err)
			}
			fmt.Fprintf(&b, "\t\t%s\n", formatRow(row, namer))
		}
		b.WriteString("\t},\n")
	}
	b.WriteString("}\n")

	body := b.String()

	// Imports are gated on the body actually using them: an unused import does not compile, and this
	// file is one the generator itself needs in order to run again.
	var head strings.Builder
	fmt.Fprintf(&head, "// Code generated by tools/database/gen_spelldata. DO NOT EDIT.\n\n")
	fmt.Fprintf(&head, "package %s\n\n", pkg)
	var std, mod []string
	if strings.Contains(body, "time.") {
		std = append(std, `"time"`)
	}
	mod = append(mod, `"github.com/wowsims/forever/sim/common/shared"`)
	if strings.Contains(body, "core.") {
		mod = append(mod, `"github.com/wowsims/forever/sim/core"`)
	}
	head.WriteString("import (\n")
	for _, i := range std {
		fmt.Fprintf(&head, "\t%s\n", i)
	}
	if len(std) > 0 {
		head.WriteString("\n")
	}
	for _, i := range mod {
		fmt.Fprintf(&head, "\t%s\n", i)
	}
	head.WriteString(")\n\n")

	out, err := format.Source([]byte(head.String() + notGenerated.String() + body))
	if err != nil {
		return nil, fmt.Errorf("generated %s file does not parse, refusing to write it: %w", pkg, err)
	}
	return out, nil
}

// core's SpellSchool bits are the client's, so this is a rendering and not a translation: the name
// emitted for a bit holds that same bit. TestSpellSchoolBitsMatchTheClient asserts the pairing.
var schoolBitNames = map[int32]string{
	1:  "core.SpellSchoolPhysical",
	2:  "core.SpellSchoolHoly",
	4:  "core.SpellSchoolFire",
	8:  "core.SpellSchoolNature",
	16: "core.SpellSchoolFrost",
	32: "core.SpellSchoolShadow",
	64: "core.SpellSchoolArcane",
}

// Nine spells in the build carry two schools - Frostfire and the Nature/Shadow plagues - and the OR of
// the two names is the value core's own named combination holds, so they are spelled out rather than
// matched against a list that would need keeping in step.
func schoolName(mask int32) string {
	var names []string
	for bit := int32(1); bit <= 64; bit <<= 1 {
		if mask&bit != 0 {
			name, ok := schoolBitNames[bit]
			if !ok {
				return ""
			}
			names = append(names, name)
		}
	}
	if len(names) == 0 {
		return ""
	}
	return strings.Join(names, " | ")
}

func defenseTypeName(t int32) string {
	switch t {
	case 1:
		return "core.DefenseTypeMagic"
	case 2:
		return "core.DefenseTypeMelee"
	case 3:
		return "core.DefenseTypeRanged"
	}
	return ""
}

func formatRow(row generatedRow, namer *rankEnumNamer) string {
	parts := []string{fmt.Sprintf("Rank: %d", row.Rank), fmt.Sprintf("SpellID: %d", row.SpellID)}
	if row.Cost > 0 {
		parts = append(parts, fmt.Sprintf("Cost: %d", row.Cost))
	}
	if row.CastTimeMs > 0 {
		parts = append(parts, fmt.Sprintf("CastTime: %s", millis(row.CastTimeMs)))
	}
	if row.GCDMs > 0 {
		parts = append(parts, fmt.Sprintf("GCD: %s", millis(row.GCDMs)))
	}
	if row.CooldownMs > 0 {
		parts = append(parts, fmt.Sprintf("Cooldown: %s", millis(row.CooldownMs)))
	}
	if row.MinRange > 0 {
		parts = append(parts, fmt.Sprintf("MinRange: %s", num(row.MinRange)))
	}
	if row.MaxRange > 0 {
		parts = append(parts, fmt.Sprintf("MaxRange: %s", num(row.MaxRange)))
	}
	if row.MissileSpeed > 0 {
		parts = append(parts, fmt.Sprintf("MissileSpeed: %s", num(row.MissileSpeed)))
	}
	if row.ProcChance > 0 {
		parts = append(parts, fmt.Sprintf("ProcChance: %d", row.ProcChance))
	}
	if row.FlatThreatBonus != 0 {
		parts = append(parts, fmt.Sprintf("FlatThreatBonus: %s", num(row.FlatThreatBonus)))
	}
	if name := schoolName(row.SchoolMask); name != "" {
		parts = append(parts, "SpellSchool: "+name)
	}
	if name := defenseTypeName(row.DefenseType); name != "" {
		parts = append(parts, "DefenseType: "+name)
	}
	if len(row.Effects) > 0 {
		var es []string
		for _, e := range row.Effects {
			f := fmt.Sprintf("{Index: %d, Effect: %s, Aura: %s, Misc: %d, Value: %s",
				e.Index, namer.Effect(e.Effect), namer.Aura(e.Aura), e.Misc, num(e.Value))
			if e.ValueMax > 0 {
				f += ", ValueMax: " + num(e.ValueMax)
			}
			es = append(es, f+"}")
		}
		parts = append(parts, "Effects: []shared.SpellDataEffect{"+strings.Join(es, ", ")+"}")
	}
	for _, role := range []struct {
		name  string
		value *generatedAmount
	}{
		{"Direct", row.Direct},
		{"Heal", row.Heal},
		{"Periodic", row.Periodic},
		{"Energize", row.Energize},
	} {
		if role.value != nil {
			parts = append(parts, role.name+": "+formatValue(*role.value))
		}
	}
	return "{" + strings.Join(parts, ", ") + "},"
}

// Picks the variant from the shape of the data: a periodic value if it ticks, a range if the client
// rolls it, a flat number otherwise.
func formatValue(a generatedAmount) string {
	tail := fmt.Sprintf("Coef: %s", num(a.Coef))
	if a.APCoef > 0 {
		tail += fmt.Sprintf(", APCoef: %s", num(a.APCoef))
	}

	switch {
	case a.PeriodMs > 0:
		tick := num(a.Min)
		if a.Max > a.Min {
			tick += fmt.Sprintf(", TickMax: %s", num(a.Max))
		}
		out := fmt.Sprintf("shared.SpellDataPeriodic{Tick: %s, %s, TickLength: %s", tick, tail, millis(a.PeriodMs))
		if a.Ticks > 0 {
			out += fmt.Sprintf(", NumberOfTicks: %d", a.Ticks)
		}
		return out + "}"
	case a.Max > a.Min:
		return fmt.Sprintf("shared.SpellDataRange{Min: %s, Max: %s, %s}", num(a.Min), num(a.Max), tail)
	default:
		return fmt.Sprintf("shared.SpellDataFlat{Value: %s, %s}", num(a.Min), tail)
	}
}

// Written as a time.Duration expression rather than a bare number, so the generated file reads the way
// a hand-written cast time or tick length does.
func millis(ms int32) string {
	return fmt.Sprintf("%d * time.Millisecond", ms)
}

func num(f float64) string {
	return strconv.FormatFloat(f, 'g', -1, 64)
}
