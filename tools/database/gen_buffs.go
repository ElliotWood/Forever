package database

// Resolves tools/database/buffmanifest against the client database and renders
// sim/core/buffs_auto_gen.go and sim/core/debuffs_auto_gen.go from the result.
//
// The manifest names which spell each proto field is; everything else - anchor
// rank, values at level 60, duration, stacks, talent scaling, owner class - is
// read here. A row the generator cannot express in the support API renders as a
// commented shell carrying the reason, which is also how a row whose hand-written
// constructor is still in sim/core renders: deleting that constructor and its
// apply block is what switches the row over to generated code.

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
	"github.com/wowsims/forever/tools/database/buffmanifest"
	"github.com/wowsims/forever/tools/database/dbc"
)

const buffsGenFile = "sim/core/buffs_auto_gen.go"
const debuffsGenFile = "sim/core/debuffs_auto_gen.go"

// The two skill lines that grant runes rather than class abilities. Both are
// CategoryID 7, so the anchor query has to name them.
const skillLineEngraving = 2851
const skillLineRunes = 2853

// The combo points a finisher the raid config simply has on the target is cast
// with.
const maxComboPoints = 5

// StatAmount is one stat the aura grants, at level 60.
type StatAmount struct {
	Stat           stats.Stat
	Amount         float64
	Multiplicative bool
}

// PseudoMod is one PseudoStats field the aura modifies.
type PseudoMod struct {
	Kind           string // PseudoStats field name, e.g. "ThreatMultiplier"
	Amount         float64
	Multiplicative bool
	SchoolMask     int32 // client school bits, 0 = every school
}

// ResolvedEffect is one SpellEffect row of the anchor spell, with its value
// already derived for level 60 and truncated toward zero.
type ResolvedEffect struct {
	Index          int32
	Effect         dbc.SpellEffectType
	Aura           dbc.EffectAuraType
	Misc           int32
	Value          float64
	PerResource    float64
	PeriodMs       int32
	ImplicitTarget dbc.ImplicitTarget
}

// ResolvedBuff is one manifest row plus everything the database states about it.
type ResolvedBuff struct {
	buffmanifest.BuffSpec

	SpellID    int32 // the anchor rank, or the aura family member when AuraName is set
	Rank       int32 // 0 when the family has no rank subtext
	DurationMs int32 // -1 or 0 never expires
	MaxStacks  int32
	SchoolMask int32
	Effects    []ResolvedEffect

	Stats  []StatAmount
	Pseudo []PseudoMod

	// OverrideStats is the manifest's StatOverride resolved to sim stats.
	OverrideStats []stats.Stat

	// TalentCurve holds the finished value per talent point: index 0 is the
	// untalented value, index n the value with n points spent. Nil when the
	// improving talent has no node in the owner's live tree.
	TalentCurve   []float64
	TalentApplies buffmanifest.TalentApplies

	// TalentOnPseudo says the curve prices Pseudo[0] rather than Stats[0].
	TalentOnPseudo bool
	// TalentSpellID is the spell of the trait node that prices the improvement,
	// which is the icon the UI shows for the improved state.
	TalentSpellID int32

	OwnerClassMask int32
	ScopeFromDB    buffmanifest.BuffScope

	// StatCategory is what the row's stats compete under, which for a resistance
	// buff is the school every other source of that resistance also bids in. The
	// manifest's Category stays the aura's own.
	StatCategory string

	DBName string

	// Note is what the generated file says about the row above its constructor,
	// for a value the client states in a way the row had to be read through.
	Note string

	Supported   bool
	Reason      string
	HandWritten bool
	ProtoStale  bool
	Warnings    []string
}

// MaxTalentPoints is the rank cap the apply block passes to GetTristateValueInt32.
func (r ResolvedBuff) MaxTalentPoints() int32 {
	if len(r.TalentCurve) == 0 {
		return 0
	}
	return int32(len(r.TalentCurve) - 1)
}

func (r *ResolvedBuff) warn(format string, args ...any) {
	r.Warnings = append(r.Warnings, fmt.Sprintf(format, args...))
}

func (r *ResolvedBuff) unsupported(format string, args ...any) {
	r.Supported = false
	r.Reason = fmt.Sprintf(format, args...)
}

type buffResolver struct {
	db          *sql.DB
	runeGranted map[int32]bool
	trees       map[int]int     // class mask -> live trait tree
	traits      map[int][]trait // trait tree -> its definitions, loaded on demand
	handWritten map[string]bool // Go stem of every hand-written <Go>Aura in sim/core
	handApplied map[string]bool // proto field the hand-written apply functions still read
}

// One node of a class tree, as far as a buff cares: which spell it modifies and
// by how much.
type trait struct {
	DefID    int32
	SpellID  int32
	MaxRanks int32
	Name     string
}

var buffRankSubtext = regexp.MustCompile(`^Rank (\d+)$`)

// ResolveBuffManifest reads every manifest row out of the client database.
func ResolveBuffManifest(helper *DBHelper) ([]ResolvedBuff, error) {
	res, err := newBuffResolver(helper)
	if err != nil {
		return nil, err
	}

	rows := make([]ResolvedBuff, 0, len(buffmanifest.Manifest))
	for _, spec := range buffmanifest.Manifest {
		row, err := res.resolve(spec)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", spec.Field, err)
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func newBuffResolver(helper *DBHelper) (*buffResolver, error) {
	runes, err := loadRuneGrantedSpells(helper.db)
	if err != nil {
		return nil, err
	}
	trees, err := selectTraitTrees(helper)
	if err != nil {
		return nil, err
	}
	root, err := repoRoot()
	if err != nil {
		return nil, err
	}
	constructors, applied, err := scanHandWrittenBuffs(filepath.Join(root, "sim", "core"))
	if err != nil {
		return nil, err
	}
	return &buffResolver{
		db:          helper.db,
		runeGranted: runes,
		trees:       trees,
		traits:      map[int][]trait{},
		handWritten: constructors,
		handApplied: applied,
	}, nil
}

// The spells a rune grants. An `Engrave <slot> - <name>` spell enchants the item
// with a SpellItemEnchantment whose EffectArg columns hold the spells the rune
// teaches, and those are Season of Discovery content rather than abilities a
// class trains.
func loadRuneGrantedSpells(db *sql.DB) (map[int32]bool, error) {
	rows, err := db.Query(`
		SELECT COALESCE(sie.EffectArg_0, 0), COALESCE(sie.EffectArg_1, 0), COALESCE(sie.EffectArg_2, 0)
		FROM SpellEffect se
		JOIN SpellName n ON n.ID = se.SpellID
		JOIN SpellItemEnchantment sie ON sie.ID = se.EffectMiscValue_0
		WHERE se.Effect = 54 AND n.Name_lang LIKE 'Engrave % - %'`)
	if err != nil {
		return nil, fmt.Errorf("rune-granted spells: %w", err)
	}
	defer rows.Close()

	granted := map[int32]bool{}
	for rows.Next() {
		var a, b, c int32
		if err := rows.Scan(&a, &b, &c); err != nil {
			return nil, err
		}
		for _, id := range []int32{a, b, c} {
			if id != 0 {
				granted[id] = true
			}
		}
	}
	return granted, rows.Err()
}

// The repository root, found by walking up from the working directory until
// go.mod appears. The generator runs from the root and the regeneration test
// from tools/database, and both have to reach sim/core.
func repoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("no go.mod above the working directory, so the repository root is unknown")
		}
		dir = parent
	}
}

// Which buffs sim/core still implements by hand: the Go stems that already have a
// <Go>Aura constructor, and the proto fields applyBuffEffects and
// applyDebuffEffects still read. Both pin a row to a shell - generating a second
// aura for a buff the hand-written code still applies would apply it twice - and
// deleting both is what hands the row over to the generator.
func scanHandWrittenBuffs(coreDir string) (map[string]bool, map[string]bool, error) {
	entries, err := os.ReadDir(coreDir)
	if err != nil {
		return nil, nil, fmt.Errorf("reading %s: %w", coreDir, err)
	}

	constructors := map[string]bool{}
	applied := map[string]bool{}
	fset := token.NewFileSet()

	var names []string
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_auto_gen.go") {
			continue
		}
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		file, err := parser.ParseFile(fset, filepath.Join(coreDir, name), nil, 0)
		if err != nil {
			return nil, nil, fmt.Errorf("parsing %s: %w", name, err)
		}
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv != nil {
				continue
			}
			if stem, found := strings.CutSuffix(fn.Name.Name, "Aura"); found && stem != "" {
				constructors[stem] = true
			}
			if fn.Name.Name != "applyBuffEffects" && fn.Name.Name != "applyDebuffEffects" {
				continue
			}
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				if sel, ok := n.(*ast.SelectorExpr); ok {
					applied[sel.Sel.Name] = true
				}
				return true
			})
		}
	}
	return constructors, applied, nil
}

// The proto message each scope's fields live on. The generated apply blocks read
// the field through the type the sim compiled against, so the manifest's declared
// type is checked against this rather than against proto/buffs.proto: a row
// whose proto field has not been retyped yet would otherwise generate code that
// cannot compile.
var buffScopeMessages = map[buffmanifest.BuffScope]reflect.Type{
	buffmanifest.ScopeRaid:       reflect.TypeOf(proto.RaidBuffs{}),
	buffmanifest.ScopeParty:      reflect.TypeOf(proto.PartyBuffs{}),
	buffmanifest.ScopeIndividual: reflect.TypeOf(proto.IndividualBuffs{}),
	buffmanifest.ScopeDebuff:     reflect.TypeOf(proto.Debuffs{}),
}

var buffProtoTypeNames = map[buffmanifest.BuffProtoType]string{
	buffmanifest.ProtoBool:      "bool",
	buffmanifest.ProtoTristate:  "proto.TristateEffect",
	buffmanifest.ProtoInt32:     "int32",
	buffmanifest.ProtoDouble:    "float64",
	buffmanifest.ProtoEnumDrums: "proto.Drums",
}

// The Go type the compiled proto states for a field, spelled the way the manifest
// spells it.
func compiledProtoType(spec buffmanifest.BuffSpec) (string, error) {
	message, ok := buffScopeMessages[spec.Scope]
	if !ok {
		return "", fmt.Errorf("unknown scope %s", spec.Scope)
	}
	field, ok := message.FieldByName(spec.GoField())
	if !ok {
		return "", fmt.Errorf("%s has no field %s", message.Name(), spec.GoField())
	}
	switch field.Type.Kind() {
	case reflect.Bool:
		return "bool", nil
	case reflect.Float64:
		return "float64", nil
	case reflect.Int32:
		if name := field.Type.Name(); name != "int32" {
			return "proto." + name, nil
		}
		return "int32", nil
	}
	return field.Type.String(), nil
}

func (res *buffResolver) resolve(spec buffmanifest.BuffSpec) (ResolvedBuff, error) {
	row := ResolvedBuff{BuffSpec: spec, Supported: true, ScopeFromDB: spec.Scope}

	compiled, err := compiledProtoType(spec)
	if err != nil {
		return row, err
	}

	switch spec.Kind {
	case buffmanifest.KindAbsent:
		row.unsupported("%s", spec.Notes)
		return row, nil
	case buffmanifest.KindFlag, buffmanifest.KindEnum:
		row.unsupported("%s", spec.Notes)
		return row, nil
	}

	if err := res.resolveSpell(&row); err != nil {
		return row, err
	}
	if !row.Supported {
		return row, nil
	}

	if err := resolveStatOverride(&row); err != nil {
		return row, err
	}

	res.mapEffects(&row)
	if err := res.resolveTalent(&row); err != nil {
		return row, err
	}

	res.validateScope(&row)
	if err := res.validateOwner(&row); err != nil {
		return row, err
	}
	if err := res.validateImpAction(&row); err != nil {
		return row, err
	}

	// A row that is already a shell keeps the reason the database gave it; the
	// two checks below are about what sim/core and the proto still look like, and
	// both are true of plenty of rows the client could not describe either.
	switch {
	case res.handWritten[spec.Go]:
		row.HandWritten = true
		if row.Supported {
			row.unsupported("hand-written constructor still present")
		}
	case res.handApplied[spec.GoField()]:
		row.HandWritten = true
		if row.Supported {
			row.unsupported("hand-written apply block still present")
		}
	}

	if compiled != buffProtoTypeNames[spec.Proto] {
		row.ProtoStale = true
		if row.Supported {
			row.unsupported("proto field not yet retyped: the sim compiled against %s, the manifest declares %s",
				compiled, buffProtoTypeNames[spec.Proto])
		}
	}
	return row, nil
}

type buffCandidate struct {
	SpellID       int32
	Subtext       string
	Rank          int32
	ClassMask     int32
	SkillLine     int32
	AcquireMethod int32
	Supercedes    int32
}

// The spell the row's values are read from: the explicit anchor, or the top rank
// of the named family, and then the aura family member when the cast is a summon
// or a dummy.
func (res *buffResolver) resolveSpell(row *ResolvedBuff) error {
	var subtext string

	switch {
	case row.Anchor != 0:
		exists, err := res.spellExists(row.Anchor)
		if err != nil {
			return err
		}
		if !exists {
			row.unsupported("anchor spell %d has no SpellName row", row.Anchor)
			return nil
		}
		row.SpellID = row.Anchor
		if err := res.db.QueryRow(`SELECT COALESCE(NameSubtext_lang, '') FROM Spell WHERE ID = ?`,
			row.Anchor).Scan(&subtext); err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}
	case row.Name != "":
		cands, err := res.anchorCandidates(row.Name, ownerClassMask(row.Owner))
		if err != nil {
			return err
		}
		if len(cands) == 0 {
			row.unsupported("no SkillLineAbility row grants %q to %s", row.Name, row.Owner)
			return nil
		}
		chosen, err := res.pickTopRank(row, cands)
		if err != nil {
			return err
		}
		row.SpellID = chosen.SpellID
		row.Rank = chosen.Rank
		subtext = chosen.Subtext
	default:
		row.unsupported("the manifest row states neither Name nor Anchor")
		return nil
	}

	if row.AuraName != "" {
		id, warning, err := res.resolveAuraFamily(row.AuraName, subtext)
		if err != nil {
			return err
		}
		if id == 0 {
			row.unsupported("no spell named %q carries an aura effect", row.AuraName)
			return nil
		}
		if warning != "" {
			row.warn("%s", warning)
		}
		row.SpellID = id
	}

	if res.runeGranted[row.SpellID] {
		row.unsupported("spell %d is granted by a rune", row.SpellID)
		return nil
	}

	return res.loadSpellData(row)
}

func (res *buffResolver) spellExists(spellID int32) (bool, error) {
	var found int32
	err := res.db.QueryRow(`SELECT ID FROM SpellName WHERE ID = ?`, spellID).Scan(&found)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return err == nil, err
}

// Every spell of the named family the owner can learn. CategoryID 7 is the class
// and profession skill lines; Engraving and Runes sit in it too and grant spells
// no class trains, so they are excluded by id.
func (res *buffResolver) anchorCandidates(name string, classMask int32) ([]buffCandidate, error) {
	rows, err := res.db.Query(`
		SELECT sla.Spell, COALESCE(s.NameSubtext_lang, ''), sla.ClassMask, sla.SkillLine,
		       sla.AcquireMethod, COALESCE(sla.SupercedesSpell, 0)
		FROM SkillLineAbility sla
		JOIN SpellName n ON n.ID = sla.Spell
		JOIN Spell s ON s.ID = sla.Spell
		JOIN SkillLine sl ON sl.ID = sla.SkillLine AND sl.CategoryID = 7
		  AND sl.ID NOT IN (?, ?)
		WHERE n.Name_lang = ? AND ((sla.ClassMask & ?) != 0 OR sla.ClassMask = 0)
		ORDER BY sla.Spell`, skillLineEngraving, skillLineRunes, name, classMask)
	if err != nil {
		return nil, fmt.Errorf("anchor candidates for %q: %w", name, err)
	}
	defer rows.Close()

	var out []buffCandidate
	for rows.Next() {
		var c buffCandidate
		if err := rows.Scan(&c.SpellID, &c.Subtext, &c.ClassMask, &c.SkillLine,
			&c.AcquireMethod, &c.Supercedes); err != nil {
			return nil, err
		}
		if m := buffRankSubtext.FindStringSubmatch(c.Subtext); m != nil {
			rank, _ := strconv.Atoi(m[1])
			c.Rank = int32(rank)
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// The rank the player ends up with. A family with rank subtexts is ordered by
// them; a family without one - Blessing of Kings, Innervate, Commanding Shout -
// is ordered by SupercedesSpell, where the spell no other row supersedes is the
// last one learnt.
func (res *buffResolver) pickTopRank(row *ResolvedBuff, cands []buffCandidate) (buffCandidate, error) {
	classMask := ownerClassMask(row.Owner)
	var owned []buffCandidate
	for _, c := range cands {
		if classMask != 0 && c.ClassMask&classMask != 0 {
			owned = append(owned, c)
		}
	}
	if len(owned) == 0 {
		owned = cands
	}

	byRank := map[int32][]buffCandidate{}
	ranked := false
	for _, c := range owned {
		byRank[c.Rank] = append(byRank[c.Rank], c)
		if c.Rank > 0 {
			ranked = true
		}
	}

	var chosen buffCandidate
	if ranked {
		var ranks []int32
		for rank := range byRank {
			if rank > 0 {
				ranks = append(ranks, rank)
			}
		}
		sort.Slice(ranks, func(i, j int) bool { return ranks[i] < ranks[j] })
		// Narrowed per rank rather than over the whole family: Forever re-issues
		// single ranks as granted copies of a trained ability, so a family-wide
		// narrowing would drop every rank the client only grants.
		top := lowestBuffAcquireMethod(byRank[ranks[len(ranks)-1]])
		chosen = top[0]
		if len(top) > 1 {
			row.warn("rank %d of %q is ambiguous between %s, taking %d",
				chosen.Rank, row.Name, spellIDList(top), chosen.SpellID)
		}
		if err := res.checkLadder(row, byRank, ranks); err != nil {
			return chosen, err
		}
	} else {
		superseded := map[int32]bool{}
		for _, c := range owned {
			if c.Supercedes != 0 {
				superseded[c.Supercedes] = true
			}
		}
		var tops []buffCandidate
		for _, c := range owned {
			if !superseded[c.SpellID] {
				tops = append(tops, c)
			}
		}
		if len(tops) == 0 {
			tops = owned
		}
		chosen = tops[0]
		if len(tops) > 1 {
			row.warn("%q has no rank subtext and %s all end the supersede chain, taking %d",
				row.Name, spellIDList(tops), chosen.SpellID)
		}
	}
	return chosen, nil
}

// A family whose later rank is worth less than its earlier one, which is a
// client-data fact the sim has to be told about rather than a resolution error.
func (res *buffResolver) checkLadder(row *ResolvedBuff, byRank map[int32][]buffCandidate, ranks []int32) error {
	prevRank, prevValue := int32(0), math.Inf(-1)
	var prevID int32
	for _, rank := range ranks {
		spellID := byRank[rank][0].SpellID
		value, ok, err := res.primaryValue(spellID)
		if err != nil {
			return err
		}
		if !ok {
			continue
		}
		if prevValue > math.Inf(-1) && math.Abs(value) < math.Abs(prevValue) {
			row.warn("non-monotonic ladder: rank %d (%d) is worth %v where rank %d (%d) is worth %v",
				rank, spellID, value, prevRank, prevID, prevValue)
		}
		prevRank, prevValue, prevID = rank, value, spellID
	}
	return nil
}

// What an effect is worth at level 60. DeriveRankAmount floors the scaled
// amount, which rounds a negative one away from zero; the client's tooltip
// truncates toward zero instead, so Demoralizing Shout reads -204 and not -205.
// The scaling below is DeriveRankAmount's own, re-run in float32 for the
// negative case only.
func buffValueAt60(e RankEffect, spellLevel int32, maxLevel int32) float64 {
	floored, _ := DeriveRankAmount(e, spellLevel, maxLevel)
	if floored >= 0 {
		return floored
	}

	maxScalingLevel := maxLevel
	if maxScalingLevel <= 0 {
		maxScalingLevel = RankLevel
	}
	level := int32(RankLevel)
	if maxScalingLevel < level {
		level = maxScalingLevel
	}
	delta := level - spellLevel
	if delta < 0 {
		delta = 0
	}
	base := float32(e.BasePoints) + float32(float32(delta)*float32(e.PointsPerLvl))
	return math.Trunc(float64(base))
}

// The number a rank is judged by: the first aura effect's value at level 60.
func (res *buffResolver) primaryValue(spellID int32) (float64, bool, error) {
	spell, err := LoadRankSpell(res.db, spellID)
	if err != nil {
		return 0, false, err
	}
	for _, e := range spell.Effects {
		if !isAuraApplication(e.Effect) {
			continue
		}
		return buffValueAt60(e, spell.SpellLevel, spell.MaxLevel), true, nil
	}
	return 0, false, nil
}

func isAuraApplication(effect dbc.SpellEffectType) bool {
	return effect == dbc.E_APPLY_AURA || effect == dbc.E_APPLY_AREA_AURA_PARTY ||
		effect == dbc.E_APPLY_AREA_AURA_RAID
}

// The aura of a totem or of a dummy passive, which the client only ties to the
// cast by name. The rank subtext of the cast picks the matching rank; where
// neither carries one - Leader of the Pack is 17007 and 24932, both rankless -
// the spell that applies a party or raid area aura is the one other players see,
// and the self-only passive is the caster's own.
func (res *buffResolver) resolveAuraFamily(name string, subtext string) (int32, string, error) {
	rows, err := res.db.Query(`
		SELECT n.ID, COALESCE(s.NameSubtext_lang, ''), e.Effect, COALESCE(e.ImplicitTarget_0, 0)
		FROM SpellName n
		JOIN Spell s ON s.ID = n.ID
		JOIN SpellEffect e ON e.SpellID = n.ID
		WHERE n.Name_lang = ? AND e.Effect IN (?, ?, ?)
		ORDER BY n.ID, e.EffectIndex`,
		name, dbc.E_APPLY_AURA, dbc.E_APPLY_AREA_AURA_PARTY, dbc.E_APPLY_AREA_AURA_RAID)
	if err != nil {
		return 0, "", fmt.Errorf("aura family %q: %w", name, err)
	}
	defer rows.Close()

	type auraCandidate struct {
		SpellID int32
		Subtext string
		Shared  bool
	}
	seen := map[int32]int{}
	var cands []auraCandidate
	for rows.Next() {
		var id, target int32
		var sub string
		var effect dbc.SpellEffectType
		if err := rows.Scan(&id, &sub, &effect, &target); err != nil {
			return 0, "", err
		}
		shared := effect == dbc.E_APPLY_AREA_AURA_PARTY || effect == dbc.E_APPLY_AREA_AURA_RAID ||
			isSharedTarget(dbc.ImplicitTarget(target))
		if index, ok := seen[id]; ok {
			cands[index].Shared = cands[index].Shared || shared
			continue
		}
		seen[id] = len(cands)
		cands = append(cands, auraCandidate{SpellID: id, Subtext: sub, Shared: shared})
	}
	if err := rows.Err(); err != nil {
		return 0, "", err
	}
	if len(cands) == 0 {
		return 0, "", nil
	}

	var matching []auraCandidate
	for _, c := range cands {
		if c.Subtext == subtext {
			matching = append(matching, c)
		}
	}
	if len(matching) == 0 {
		matching = cands
	}
	if len(matching) == 1 {
		return matching[0].SpellID, "", nil
	}

	var shared []auraCandidate
	for _, c := range matching {
		if c.Shared {
			shared = append(shared, c)
		}
	}
	if len(shared) == 1 {
		return shared[0].SpellID, fmt.Sprintf(
			"aura name %q is ambiguous; taking %d, the one that applies a party or raid aura",
			name, shared[0].SpellID), nil
	}
	if len(shared) > 1 {
		matching = shared
	}
	var ids []string
	for _, c := range matching {
		ids = append(ids, strconv.Itoa(int(c.SpellID)))
	}
	return matching[0].SpellID, fmt.Sprintf("aura name %q is ambiguous between %s, taking %d",
		name, strings.Join(ids, ", "), matching[0].SpellID), nil
}

func isSharedTarget(target dbc.ImplicitTarget) bool {
	switch target {
	case dbc.TARGET_UNIT_CASTER_AREA_PARTY, dbc.TARGET_UNIT_CASTER_AREA_RAID,
		dbc.TARGET_UNIT_TARGET_ALLY, dbc.TARGET_UNIT_TARGET_RAID, dbc.TARGET_UNIT_TARGET_ALLY_OR_RAID:
		return true
	}
	return false
}

// Duration, stacks, school and every effect of the resolved spell, at level 60.
func (res *buffResolver) loadSpellData(row *ResolvedBuff) error {
	spell, err := LoadRankSpell(res.db, row.SpellID)
	if err != nil {
		return err
	}
	row.SchoolMask = spell.SchoolMask

	if err := res.db.QueryRow(`SELECT Name_lang FROM SpellName WHERE ID = ?`, row.SpellID).
		Scan(&row.DBName); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	// LoadRankSpell reads the duration through COALESCE, which cannot tell a
	// spell with no duration row from one the client states as 0. Both mean the
	// aura never expires, so the distinction does not matter here.
	if err := scanOptional(res.db, `
		SELECT COALESCE(d.Duration, 0)
		FROM SpellMisc m LEFT JOIN SpellDuration d ON d.ID = m.DurationIndex
		WHERE m.SpellID = ?`, row.SpellID, &row.DurationMs); err != nil {
		return fmt.Errorf("duration of spell %d: %w", row.SpellID, err)
	}
	if err := scanOptional(res.db,
		`SELECT COALESCE(CumulativeAura, 0) FROM SpellAuraOptions WHERE SpellID = ?`,
		row.SpellID, &row.MaxStacks); err != nil {
		return fmt.Errorf("stacks of spell %d: %w", row.SpellID, err)
	}

	extra, err := res.effectTargets(row.SpellID)
	if err != nil {
		return err
	}

	for _, e := range spell.Effects {
		resolved := ResolvedEffect{
			Index:    e.Index,
			Effect:   e.Effect,
			Aura:     e.Aura,
			Misc:     e.MiscValue,
			Value:    buffValueAt60(e, spell.SpellLevel, spell.MaxLevel),
			PeriodMs: e.AuraPeriod,
		}
		if t, ok := extra[e.Index]; ok {
			resolved.ImplicitTarget = t.target
			resolved.PerResource = t.perResource
		}
		row.Effects = append(row.Effects, resolved)
	}
	return nil
}

type effectTarget struct {
	target      dbc.ImplicitTarget
	perResource float64
}

func (res *buffResolver) effectTargets(spellID int32) (map[int32]effectTarget, error) {
	rows, err := res.db.Query(`
		SELECT EffectIndex, COALESCE(ImplicitTarget_0, 0), COALESCE(EffectPointsPerResource, 0)
		FROM SpellEffect WHERE SpellID = ?`, spellID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := map[int32]effectTarget{}
	for rows.Next() {
		var index, target int32
		var perResource float64
		if err := rows.Scan(&index, &target, &perResource); err != nil {
			return nil, err
		}
		out[index] = effectTarget{target: dbc.ImplicitTarget(target), perResource: perResource}
	}
	return out, rows.Err()
}

// The stats or pseudo-stats the row's kind says its effects are. An aura this
// cannot express is not an error: the row becomes a shell naming the aura, which
// is the signal that either the support API or the manifest has to grow.
func (res *buffResolver) mapEffects(row *ResolvedBuff) {
	var unmapped []string
	damageShield := false
	for i := range row.Effects {
		e := row.Effects[i]
		if !isAuraApplication(e.Effect) {
			continue
		}
		if e.Aura == dbc.A_DAMAGE_SHIELD {
			damageShield = true
			continue
		}
		// The value a spell states per combo point is 0 on the spell itself.
		// The raid config is one debuff that is simply on the target, which is
		// the finisher at full combo points; a caster spending fewer of them
		// takes its own value through a driver.
		if e.PerResource != 0 && e.Value == 0 {
			if row.Kind != buffmanifest.KindDebuffStat {
				row.unsupported("effect %d is worth %v per combo point, which the support API cannot express",
					e.Index, e.PerResource)
				return
			}
			row.Effects[i].Value = e.PerResource * maxComboPoints
			e = row.Effects[i]
			row.Note = fmt.Sprintf("Effect %d is worth %s per combo point; this is the %d-point finisher.",
				e.Index, formatFloat(e.PerResource), maxComboPoints)
		}
		if amounts, ok := row.statAmounts(e); ok {
			row.Stats = append(row.Stats, amounts...)
			continue
		}
		if mods, ok := pseudoModsOf(e); ok {
			row.Pseudo = append(row.Pseudo, mods...)
			continue
		}
		if !slices.Contains(unmapped, strconv.Itoa(int(e.Aura))) {
			unmapped = append(unmapped, strconv.Itoa(int(e.Aura)))
		}
	}

	if row.Kind == buffmanifest.KindResistance && len(row.Stats) > 0 {
		row.StatCategory = resistanceCategoryOf(row.Stats[0].Stat)
		// A row whose manifest category is the school itself states the same
		// thing twice: the school is the competition, and the aura has no
		// exclusivity of its own beyond it.
		if row.StatCategory == row.Category {
			row.Category = ""
		}
	}

	switch row.Kind {
	case buffmanifest.KindDamageShield:
		if !damageShield {
			row.unsupported("no A_DAMAGE_SHIELD effect on spell %d", row.SpellID)
			return
		}
	case buffmanifest.KindExternalCD, buffmanifest.KindProc, buffmanifest.KindManual,
		buffmanifest.KindDebuffUptime, buffmanifest.KindItemCount:
		// Driver kinds: the hand-written driver decides what the numbers mean.
	default:
		if len(row.Stats) == 0 && len(row.Pseudo) == 0 {
			row.unsupported("spell %d states no aura effect this generator maps (auras %s)",
				row.SpellID, strings.Join(unmapped, ", "))
			return
		}
	}

	for _, aura := range unmapped {
		row.warn("aura %s of spell %d is left out: the generator maps no stat or pseudo-stat to it",
			aura, row.SpellID)
	}
}

// The stats one effect grants. A row that names stats itself puts the effect's
// value on exactly those, which is how an aura the client leaves unqualified -
// A_MOD_CRIT_PCT says "critical strike chance" and no more - reaches the melee
// or the spell crit stat.
func (row ResolvedBuff) statAmounts(e ResolvedEffect) ([]StatAmount, bool) {
	if len(row.OverrideStats) == 0 {
		return statAmountsOf(e)
	}
	out := make([]StatAmount, 0, len(row.OverrideStats))
	for _, stat := range row.OverrideStats {
		out = append(out, StatAmount{Stat: stat, Amount: e.Value})
	}
	return out, true
}

// The manifest's StatOverride, read as sim stats. A row that states one may
// only have a single aura effect: every effect would otherwise be mapped onto
// the same stats and the amounts would add up.
func resolveStatOverride(row *ResolvedBuff) error {
	if len(row.StatOverride) == 0 {
		return nil
	}

	auraEffects := 0
	for _, e := range row.Effects {
		if isAuraApplication(e.Effect) {
			auraEffects++
		}
	}
	if auraEffects != 1 {
		return fmt.Errorf("states a StatOverride, but spell %d has %d aura effects",
			row.SpellID, auraEffects)
	}

	for _, name := range row.StatOverride {
		stat, ok := statByName(name)
		if !ok {
			return fmt.Errorf("StatOverride names %q, which is not a sim stat", name)
		}
		row.OverrideStats = append(row.OverrideStats, stat)
	}
	return nil
}

func statByName(name string) (stats.Stat, bool) {
	for stat := stats.Stat(0); stat < stats.SimStatsLen; stat++ {
		if stat.StatName() == name {
			return stat, true
		}
	}
	return 0, false
}

func statAmountsOf(e ResolvedEffect) ([]StatAmount, bool) {
	flat := func(stat stats.Stat, amount float64) ([]StatAmount, bool) {
		return []StatAmount{{Stat: stat, Amount: amount}}, true
	}

	switch e.Aura {
	case dbc.A_MOD_STAT:
		if e.Misc == -1 {
			var out []StatAmount
			for _, stat := range []stats.Stat{stats.Strength, stats.Agility, stats.Stamina, stats.Intellect, stats.Spirit} {
				out = append(out, StatAmount{Stat: stat, Amount: e.Value})
			}
			return out, true
		}
		stat, ok := dbc.MapMainStatToStat(int(e.Misc))
		if !ok {
			return nil, false
		}
		return flat(stats.Stat(stat), e.Value)
	case dbc.A_MOD_ATTACK_POWER:
		return flat(stats.AttackPower, e.Value)
	case dbc.A_MOD_RANGED_ATTACK_POWER:
		return flat(stats.RangedAttackPower, e.Value)
	case dbc.A_MOD_INCREASE_HEALTH:
		return flat(stats.Health, e.Value)
	case dbc.A_MOD_HEALING_DONE:
		return flat(stats.HealingPower, e.Value)
	case dbc.A_MOD_POWER_REGEN:
		// The client states mana per 5 seconds directly on this aura.
		return flat(stats.MP5, e.Value)
	case dbc.A_PERIODIC_ENERGIZE:
		if e.PeriodMs <= 0 {
			return nil, false
		}
		return flat(stats.MP5, math.Trunc(e.Value*5000/float64(e.PeriodMs)))
	case dbc.A_MOD_SPELL_CRIT_CHANCE:
		return flat(stats.SpellCritPercent, e.Value)
	case dbc.A_MOD_HIT_CHANCE:
		return flat(stats.PhysicalHitPercent, e.Value)
	case dbc.A_MOD_RESISTANCE:
		var out []StatAmount
		for bit, stat := range resistanceBits {
			if e.Misc&bit != 0 {
				out = append(out, StatAmount{Stat: stat, Amount: e.Value})
			}
		}
		sort.Slice(out, func(i, j int) bool { return out[i].Stat < out[j].Stat })
		// Bit 2 is Holy, which has no resistance stat in the sim; a mask of only
		// that bit resolves to nothing rather than to a wrong stat.
		return out, len(out) > 0
	case dbc.A_MOD_DAMAGE_DONE:
		stat := dbc.ConvertEffectAuraToStatIndex(e.Aura, int(e.Misc))
		if stat < 0 {
			return nil, false
		}
		return flat(stats.Stat(stat), e.Value)
	case dbc.A_MOD_RATING:
		stat := dbc.ConvertEffectAuraToStatIndex(e.Aura, int(e.Misc))
		if stat < 0 {
			return nil, false
		}
		return flat(stats.Stat(stat), e.Value)
	case dbc.A_MOD_TOTAL_STAT_PERCENTAGE:
		multiplier := 1 + e.Value/100
		if e.Misc == -1 {
			var out []StatAmount
			for _, stat := range []stats.Stat{stats.Strength, stats.Agility, stats.Stamina, stats.Intellect, stats.Spirit} {
				out = append(out, StatAmount{Stat: stat, Amount: multiplier, Multiplicative: true})
			}
			return out, true
		}
		stat, ok := dbc.MapMainStatToStat(int(e.Misc))
		if !ok {
			return nil, false
		}
		return []StatAmount{{Stat: stats.Stat(stat), Amount: multiplier, Multiplicative: true}}, true
	case dbc.A_MOD_ATTACK_POWER_PCT:
		return []StatAmount{{Stat: stats.AttackPower, Amount: 1 + e.Value/100, Multiplicative: true}}, true
	}
	return nil, false
}

// The exclusive category a resistance stat competes under, spelled the way
// sim/core spells it. Armor is not a resistance school and has none.
func resistanceCategoryOf(stat stats.Stat) string {
	switch stat {
	case stats.ArcaneResistance:
		return "ResistanceArcane"
	case stats.FireResistance:
		return "ResistanceFire"
	case stats.FrostResistance:
		return "ResistanceFrost"
	case stats.NatureResistance:
		return "ResistanceNature"
	case stats.ShadowResistance:
		return "ResistanceShadow"
	}
	return ""
}

// Holy, fire, nature, frost, shadow and arcane together, which is every school
// the sim counts as spell damage.
const everySpellSchoolMask int32 = 126

var resistanceBits = map[int32]stats.Stat{
	1:  stats.Armor,
	4:  stats.FireResistance,
	8:  stats.NatureResistance,
	16: stats.FrostResistance,
	32: stats.ShadowResistance,
	64: stats.ArcaneResistance,
}

func pseudoModsOf(e ResolvedEffect) ([]PseudoMod, bool) {
	switch e.Aura {
	case dbc.A_MOD_THREAT:
		return []PseudoMod{{Kind: "ThreatMultiplier", Amount: 1 + e.Value/100, Multiplicative: true}}, true
	case dbc.A_MOD_DAMAGE_PERCENT_DONE:
		return []PseudoMod{{Kind: "DamageDealtMultiplier", Amount: 1 + e.Value/100, Multiplicative: true}}, true
	case dbc.A_REDUCE_PUSHBACK:
		// PseudoStats.PushbackChance is the chance of being pushed back and
		// starts at 1, so the client's "35% less pushback" is -0.35 there.
		return []PseudoMod{{Kind: "PushbackChance", Amount: -e.Value / 100}}, true
	case dbc.A_MOD_MELEE_HASTE_3, dbc.A_MOD_ATTACKSPEED:
		return []PseudoMod{{Kind: "MeleeSpeedMultiplier", Amount: 1 + e.Value/100, Multiplicative: true}}, true
	case dbc.A_MOD_DAMAGE_PERCENT_TAKEN:
		return []PseudoMod{{
			Kind: "SchoolDamageTakenMultiplier", Amount: 1 + e.Value/100,
			Multiplicative: true, SchoolMask: e.Misc,
		}}, true
	case dbc.A_RANGED_ATTACK_POWER_ATTACKER_BONUS:
		return []PseudoMod{{Kind: "BonusRangedAttackPower", Amount: e.Value}}, true
	case dbc.A_MOD_DAMAGE_TAKEN:
		// The sim splits flat damage taken into a physical and a spell field,
		// so the school mask picks which one the effect is. A mask that names
		// some spell schools and not others has neither: the spell field would
		// raise what every school does to the target.
		if e.Misc&1 != 0 {
			return []PseudoMod{{Kind: "BonusPhysicalDamageTaken", Amount: e.Value}}, true
		}
		if e.Misc&everySpellSchoolMask == everySpellSchoolMask {
			return []PseudoMod{{Kind: "BonusSpellDamageTaken", Amount: e.Value}}, true
		}
		return nil, false
	}
	return nil, false
}

// The improving talent, read only from the owner's live tree. A talent that
// modifies the buff's spell family by a percentage or a flat amount scales its
// value; one that modifies misc 1 scales its duration; every other modifier
// (radius, cost, range) leaves the buff alone.
func (res *buffResolver) resolveTalent(row *ResolvedBuff) error {
	classMask := ownerClassMask(row.Owner)
	if classMask == 0 {
		if row.Proto == buffmanifest.ProtoTristate {
			return fmt.Errorf("declared ProtoTristate but the row has no owner class, so no tree can be searched")
		}
		return nil
	}
	treeID, ok := res.trees[int(classMask)]
	if !ok {
		return fmt.Errorf("the client database holds no talent tree for %s", row.Owner)
	}

	family, err := res.spellClassFamily(row.SpellID)
	if err != nil {
		return err
	}

	matches, ignored, err := res.talentMatches(treeID, family)
	if err != nil {
		return err
	}
	for _, m := range ignored {
		row.warn("%s (trait definition %d, effect %d) modifies spell %d by misc %d, which no buff reads",
			m.Name, m.DefID, m.EffectIndex, row.SpellID, m.Misc)
	}

	if row.Talent == nil {
		if row.Proto == buffmanifest.ProtoTristate {
			return fmt.Errorf("declared ProtoTristate but the manifest names no talent")
		}
		for _, m := range matches {
			row.warn("declared %s, but %s (trait definition %d, effect %d, misc %d) modifies spell %d in tree %d",
				row.Proto, m.Name, m.DefID, m.EffectIndex, m.Misc, row.SpellID, treeID)
		}
		return nil
	}

	var chosen *talentMatch
	for i := range matches {
		if matches[i].Name == row.Talent.Name && matches[i].EffectIndex == row.Talent.Effect {
			chosen = &matches[i]
			break
		}
	}
	if chosen == nil {
		if row.Proto == buffmanifest.ProtoTristate {
			return fmt.Errorf("declared ProtoTristate with talent %q effect %d, "+
				"but no node of tree %d modifies spell %d that way",
				row.Talent.Name, row.Talent.Effect, treeID, row.SpellID)
		}
		row.warn("talent %q effect %d has no node in tree %d that modifies spell %d",
			row.Talent.Name, row.Talent.Effect, treeID, row.SpellID)
		return nil
	}

	curves, err := traitCurves(res.db, chosen.DefID)
	if err != nil {
		return err
	}
	values, ok := curveRanks(curves[chosen.EffectIndex], chosen.MaxRanks)
	if !ok {
		row.warn("talent %q states no rank curve for effect %d, so the improved state is dropped",
			chosen.Name, chosen.EffectIndex)
		return nil
	}

	row.TalentSpellID = chosen.SpellID
	row.TalentApplies = row.Talent.Applies
	if chosen.Misc == 1 {
		row.TalentApplies = buffmanifest.TalentScalesDuration
	}

	row.TalentCurve = make([]float64, chosen.MaxRanks+1)
	if row.TalentApplies == buffmanifest.TalentScalesDuration {
		for rank := int32(0); rank <= chosen.MaxRanks; rank++ {
			duration := float64(row.DurationMs)
			if rank > 0 {
				duration = math.Trunc(applyTalentPoints(duration, values[rank], chosen.Aura))
			}
			row.TalentCurve[rank] = duration
		}
		return nil
	}

	target, onPseudo, ok := row.talentTarget()
	if !ok {
		row.TalentCurve = nil
		row.warn("talent %q has nothing to scale: spell %d states no amount this generator maps",
			chosen.Name, row.SpellID)
		return nil
	}
	row.TalentOnPseudo = onPseudo

	// The talent scales the number the client states, not the number the sim
	// stores: a percentage aura reaches the sim as 1 + value/100, and adding
	// talent points to that would be arithmetic on the wrong quantity.
	for rank := int32(0); rank <= chosen.MaxRanks; rank++ {
		scaled := target
		if rank > 0 {
			scaled.Value = math.Trunc(applyTalentPoints(target.Value, values[rank], chosen.Aura))
		}
		row.TalentCurve[rank] = row.convertedAmount(scaled, onPseudo)
	}
	return nil
}

// The effect the talent's curve is read through, and whether it lands on a
// pseudo-stat. It is the first effect the kind mapping took an amount from, so
// the curve and Stats[0] or Pseudo[0] describe the same number.
func (row ResolvedBuff) talentTarget() (ResolvedEffect, bool, bool) {
	// The damage shield is stepped over here exactly as the kind mapping steps
	// over it, so a spell that shields and buffs a stat prices the stat, and it
	// is only fallen back on when the spell states nothing else.
	var shield ResolvedEffect
	var shielded bool
	for _, e := range row.Effects {
		if !isAuraApplication(e.Effect) {
			continue
		}
		if e.Aura == dbc.A_DAMAGE_SHIELD {
			if !shielded {
				shield, shielded = e, true
			}
			continue
		}
		if _, ok := row.statAmounts(e); ok {
			return e, false, true
		}
		if _, ok := pseudoModsOf(e); ok {
			return e, true, true
		}
	}
	return shield, false, shielded
}

// What the effect is worth once the kind mapping has converted it.
func (row ResolvedBuff) convertedAmount(e ResolvedEffect, onPseudo bool) float64 {
	if onPseudo {
		if mods, ok := pseudoModsOf(e); ok {
			return mods[0].Amount
		}
		return e.Value
	}
	if amounts, ok := row.statAmounts(e); ok {
		return amounts[0].Amount
	}
	return e.Value
}

// A_ADD_PCT_MODIFIER states a percentage of the spell's own number;
// A_ADD_FLAT_MODIFIER states an amount to add to it.
func applyTalentPoints(base float64, points float64, aura dbc.EffectAuraType) float64 {
	if aura == dbc.A_ADD_PCT_MODIFIER {
		return base * (1 + points/100)
	}
	return base + points
}

type talentMatch struct {
	DefID       int32
	SpellID     int32
	Name        string
	MaxRanks    int32
	EffectIndex int32
	Aura        dbc.EffectAuraType
	Misc        int32
}

type spellFamily struct {
	ClassSet int32
	Mask     [4]int32
	Known    bool
}

func (res *buffResolver) spellClassFamily(spellID int32) (spellFamily, error) {
	var f spellFamily
	err := res.db.QueryRow(`
		SELECT SpellClassSet, COALESCE(SpellClassMask_0, 0), COALESCE(SpellClassMask_1, 0),
		       COALESCE(SpellClassMask_2, 0), COALESCE(SpellClassMask_3, 0)
		FROM SpellClassOptions WHERE SpellID = ?`, spellID).
		Scan(&f.ClassSet, &f.Mask[0], &f.Mask[1], &f.Mask[2], &f.Mask[3])
	if errors.Is(err, sql.ErrNoRows) {
		return f, nil
	}
	if err != nil {
		return f, fmt.Errorf("spell class options of %d: %w", spellID, err)
	}
	f.Known = true
	return f, nil
}

// Every node of the tree whose spell modifies the family, filtered down to the
// modifiers that reach a buff: 0, 3 and 8 change the value, 1 changes the
// duration.
func (res *buffResolver) talentMatches(treeID int, family spellFamily) ([]talentMatch, []talentMatch, error) {
	if !family.Known {
		return nil, nil, nil
	}

	traits, err := res.treeTraits(treeID)
	if err != nil {
		return nil, nil, err
	}

	var out, ignored []talentMatch
	for _, t := range traits {
		rows, err := res.db.Query(`
			SELECT EffectIndex, EffectAura, COALESCE(EffectMiscValue_0, 0),
			       COALESCE(EffectSpellClassMask_0, 0), COALESCE(EffectSpellClassMask_1, 0),
			       COALESCE(EffectSpellClassMask_2, 0), COALESCE(EffectSpellClassMask_3, 0)
			FROM SpellEffect
			WHERE SpellID = ? AND EffectAura IN (?, ?)
			ORDER BY EffectIndex`,
			t.SpellID, dbc.A_ADD_FLAT_MODIFIER, dbc.A_ADD_PCT_MODIFIER)
		if err != nil {
			return nil, nil, err
		}
		for rows.Next() {
			var m talentMatch
			var mask [4]int32
			if err := rows.Scan(&m.EffectIndex, &m.Aura, &m.Misc,
				&mask[0], &mask[1], &mask[2], &mask[3]); err != nil {
				rows.Close()
				return nil, nil, err
			}
			overlap := false
			for i := range mask {
				if mask[i]&family.Mask[i] != 0 {
					overlap = true
				}
			}
			if !overlap {
				continue
			}
			m.DefID, m.SpellID, m.Name, m.MaxRanks = t.DefID, t.SpellID, t.Name, t.MaxRanks
			switch m.Misc {
			case 0, 1, 3, 8:
				out = append(out, m)
			default:
				ignored = append(ignored, m)
			}
		}
		if err := rows.Err(); err != nil {
			rows.Close()
			return nil, nil, err
		}
		rows.Close()
	}
	return out, ignored, nil
}

// The tree's nodes, keyed by tree so a run that resolves nine classes reads each
// tree once. SpellClassSet is not checked here because the tree is the class's
// own; the class mask overlap is what ties the node to the buff.
func (res *buffResolver) treeTraits(treeID int) ([]trait, error) {
	if cached, ok := res.traits[treeID]; ok {
		return cached, nil
	}

	rows, err := res.db.Query(`
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
		return nil, fmt.Errorf("traits of tree %d: %w", treeID, err)
	}
	defer rows.Close()

	var out []trait
	for rows.Next() {
		var t trait
		if err := rows.Scan(&t.DefID, &t.SpellID, &t.MaxRanks, &t.Name); err != nil {
			return nil, err
		}
		if t.SpellID == 0 || t.Name == "" {
			continue
		}
		out = append(out, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	res.traits[treeID] = out
	return out, nil
}

// Who the client says the aura reaches. The manifest wins - the sim's scopes are
// a UI grouping as much as a game fact - but a disagreement is worth printing.
func (res *buffResolver) validateScope(row *ResolvedBuff) {
	scope, known := buffmanifest.ScopeIndividual, false
	for _, e := range row.Effects {
		switch {
		case e.Effect == dbc.E_APPLY_AREA_AURA_RAID:
			scope, known = buffmanifest.ScopeRaid, true
		case e.Effect == dbc.E_APPLY_AREA_AURA_PARTY:
			scope, known = buffmanifest.ScopeParty, true
		case !isAuraApplication(e.Effect):
			continue
		default:
			switch e.ImplicitTarget {
			case dbc.TARGET_UNIT_TARGET_ENEMY, dbc.TARGET_UNIT_SRC_AREA_ENEMY, dbc.TARGET_SRC_CASTER:
				scope, known = buffmanifest.ScopeDebuff, true
			case dbc.TARGET_UNIT_CASTER_AREA_RAID, dbc.TARGET_UNIT_TARGET_RAID, dbc.TARGET_UNIT_TARGET_ALLY_OR_RAID:
				scope, known = buffmanifest.ScopeRaid, true
			case dbc.TARGET_UNIT_CASTER_AREA_PARTY:
				scope, known = buffmanifest.ScopeParty, true
			case dbc.TARGET_UNIT_TARGET_ALLY, dbc.TARGET_UNIT_TARGET_ANY, dbc.TARGET_UNIT_CASTER:
				scope, known = buffmanifest.ScopeIndividual, true
			}
		}
		if known {
			break
		}
	}
	if !known {
		return
	}
	row.ScopeFromDB = scope
	if scope != row.Scope {
		row.warn("scope mismatch: the manifest says %s, spell %d reads as %s", row.Scope, row.SpellID, scope)
	}
}

func (res *buffResolver) validateOwner(row *ResolvedBuff) error {
	// Read as rows and combined here: a spell granted under two masks wants their
	// bitwise or, and SUM() would turn 1 and 3 into 4.
	rows, err := res.db.Query(
		`SELECT DISTINCT ClassMask FROM SkillLineAbility WHERE Spell = ?`, row.SpellID)
	if err != nil {
		return err
	}
	for rows.Next() {
		var granted int32
		if err := rows.Scan(&granted); err != nil {
			rows.Close()
			return err
		}
		row.OwnerClassMask |= granted
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return err
	}
	rows.Close()

	mask := ownerClassMask(row.Owner)
	// A ClassMask of 0 is the client saying "no class in particular", which is
	// how Trueshot Aura and every item-granted aura read, so it validates nothing.
	if mask == 0 || row.OwnerClassMask == 0 {
		return nil
	}
	if row.OwnerClassMask&mask == 0 {
		row.warn("owner mismatch: the manifest says %s, spell %d is granted to class mask %d",
			row.Owner, row.SpellID, row.OwnerClassMask)
	}
	return nil
}

// An improved icon the manifest points at an item only exists if the item does.
func (res *buffResolver) validateImpAction(row *ResolvedBuff) error {
	if row.ImpAction == nil || row.ImpAction.ItemID == 0 {
		return nil
	}
	var id int32
	err := res.db.QueryRow(`SELECT ID FROM Item WHERE ID = ?`, row.ImpAction.ItemID).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		row.warn("improved icon item %d has no Item row", row.ImpAction.ItemID)
		row.ImpAction = nil
		return nil
	}
	return err
}

// SkillLineAbility.ClassMask is a bitmask over the client's class ids, which is
// not the order proto.Class uses.
func ownerClassMask(class proto.Class) int32 {
	for _, c := range dbc.Classes {
		if c.ProtoClass == class {
			return int32(1) << (c.ID - 1)
		}
	}
	return 0
}

func lowestBuffAcquireMethod(cands []buffCandidate) []buffCandidate {
	best := int32(math.MaxInt32)
	for _, c := range cands {
		if c.AcquireMethod < best {
			best = c.AcquireMethod
		}
	}
	var kept []buffCandidate
	for _, c := range cands {
		if c.AcquireMethod == best {
			kept = append(kept, c)
		}
	}
	return kept
}

func spellIDList(cands []buffCandidate) string {
	var ids []string
	for _, c := range cands {
		ids = append(ids, strconv.Itoa(int(c.SpellID)))
	}
	return strings.Join(ids, ", ")
}

///////////////////////////////////////////////////////////////////////////
//							Rendering
///////////////////////////////////////////////////////////////////////////

// buffRow is what the templates see: every expression the generated file needs,
// already spelled as Go source.
type buffRow struct {
	Go            string
	Field         string
	Label         string
	SpellID       int32
	Kind          string
	Reason        string
	Note          string
	Supported     bool
	HasWowhead    bool
	Category      string
	CategoryVar   string
	ValueExpr     string
	HasValue      bool
	ValueOnPseudo bool
	Duration      string
	Constructor   string
	ApplyIf       string
	ApplyBody     string
	HasApply      bool
}

// RenderBuffFiles renders both generated files without writing them.
func RenderBuffFiles(helper *DBHelper) (map[string][]byte, error) {
	rows, err := ResolveBuffManifest(helper)
	if err != nil {
		return nil, err
	}
	return renderBuffFiles(rows)
}

func renderBuffFiles(rows []ResolvedBuff) (map[string][]byte, error) {
	buffs, err := renderBuffFile(rows, false)
	if err != nil {
		return nil, err
	}
	debuffs, err := renderBuffFile(rows, true)
	if err != nil {
		return nil, err
	}
	return map[string][]byte{buffsGenFile: buffs, debuffsGenFile: debuffs}, nil
}

func renderBuffFile(resolved []ResolvedBuff, debuffs bool) ([]byte, error) {
	var rows []buffRow
	needsTime, needsStats := false, false
	for _, row := range resolved {
		if (row.Scope == buffmanifest.ScopeDebuff) != debuffs {
			continue
		}
		rendered := renderRow(row)
		if rendered.Supported {
			needsTime = true
			needsStats = needsStats || len(row.Stats) > 0
		}
		rows = append(rows, rendered)
	}

	name := "buffs"
	tmplStr := TmplStrBuffs
	if debuffs {
		name, tmplStr = "debuffs", TmplStrDebuffs
	}

	tmpl, err := template.New(name).Parse(tmplStr)
	if err != nil {
		return nil, err
	}

	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, map[string]any{
		"Rows": rows, "NeedsTime": needsTime, "NeedsStats": needsStats,
	}); err != nil {
		return nil, fmt.Errorf("rendering %s: %w", name, err)
	}

	// The template cannot indent a commented-out shell the way gofmt wants, and
	// writing unformatted Go would make every later run diff on whitespace.
	out, err := format.Source(rendered.Bytes())
	if err != nil {
		return nil, fmt.Errorf("generated %s is not valid Go: %w", name, err)
	}
	return out, nil
}

func renderRow(row ResolvedBuff) buffRow {
	out := buffRow{
		Go: row.Go, Field: row.Field, Label: buffLabel(row),
		SpellID: row.SpellID, Kind: row.Kind.String(), Reason: row.Reason, Note: row.Note,
		Supported: row.Supported, HasWowhead: row.SpellID != 0,
	}
	if !row.Supported {
		if out.Reason == "" {
			out.Reason = "no reason given"
		}
		return out
	}

	if row.Category != "" {
		out.Category = row.Category
		out.CategoryVar = row.Go + "Category"
	}
	out.ValueExpr, out.ValueOnPseudo, out.HasValue = buffValueExpr(row)
	out.Duration = buffDurationExpr(row)
	out.Constructor = buffConstructor(row, out)
	out.ApplyIf, out.ApplyBody, out.HasApply = buffApply(row)
	return out
}

// The UI label, which is the client's name for the spell unless the manifest
// overrides it.
func buffLabel(row ResolvedBuff) string {
	if row.Label != "" {
		return row.Label
	}
	if row.Name != "" {
		return row.Name
	}
	return row.DBName
}

// The first stat amount, talent-scaled when the tree prices the talent.
func buffValueExpr(row ResolvedBuff) (string, bool, bool) {
	if len(row.TalentCurve) > 0 && row.TalentApplies != buffmanifest.TalentScalesDuration {
		return floatSliceLiteral(row.TalentCurve) + "[talentPoints]", row.TalentOnPseudo, true
	}
	if len(row.Stats) > 0 {
		return formatFloat(row.Stats[0].Amount), false, true
	}
	if len(row.Pseudo) > 0 {
		return formatFloat(row.Pseudo[0].Amount), true, true
	}
	for _, e := range row.Effects {
		if e.Aura == dbc.A_DAMAGE_SHIELD {
			return formatFloat(e.Value), false, true
		}
	}
	return "", false, false
}

func buffDurationExpr(row ResolvedBuff) string {
	if len(row.TalentCurve) > 0 && row.TalentApplies == buffmanifest.TalentScalesDuration {
		return durationSliceLiteral(row.TalentCurve) + "[talentPoints]"
	}
	if row.DurationMs <= 0 {
		return "NeverExpires"
	}
	return fmt.Sprintf("%d * time.Millisecond", row.DurationMs)
}

// The call the constructor makes into the hand-written support API.
func buffConstructor(row ResolvedBuff, rendered buffRow) string {
	var b strings.Builder
	config := buffConfigLiteral(row, rendered)

	switch row.Kind {
	case buffmanifest.KindDamageShield:
		school := buffSchoolName(row.SchoolMask)
		fmt.Fprintf(&b, "return newGeneratedDamageShield(unit, %s, %s, %s(talentPoints))",
			config, school, row.Go+"Value")
	default:
		if row.Scope == buffmanifest.ScopeDebuff {
			fmt.Fprintf(&b, "return newGeneratedDebuff(unit, %s)", config)
		} else {
			fmt.Fprintf(&b, "return newGeneratedStatAura(unit, %s)", config)
		}
	}
	return b.String()
}

func buffConfigLiteral(row ResolvedBuff, rendered buffRow) string {
	var b strings.Builder
	b.WriteString("GeneratedBuff{\n")
	fmt.Fprintf(&b, "Label: %q + Ternary(isPlayer, \"Player\", \"External\") + \")\",\n", rendered.Label+" (")
	fmt.Fprintf(&b, "ActionID: ActionID{SpellID: %d}.WithTag(TernaryInt32(isPlayer, 0, -1)),\n", row.SpellID)
	fmt.Fprintf(&b, "Duration: %sDuration(talentPoints),\n", row.Go)
	if row.MaxStacks > 0 {
		fmt.Fprintf(&b, "MaxStacks: %d,\n", row.MaxStacks)
	}
	if row.StatCategory != "" {
		fmt.Fprintf(&b, "StatCategory: %q,\n", row.StatCategory)
	}
	if rendered.CategoryVar != "" {
		fmt.Fprintf(&b, "Category: %s,\n", rendered.CategoryVar)
	}
	if row.SharedCategory != "" {
		fmt.Fprintf(&b, "SharedCategory: %q,\n", row.SharedCategory)
	}
	if row.SingleAura {
		b.WriteString("SingleAura: true,\n")
	}
	b.WriteString("IsPlayer: isPlayer,\n")

	// The talent curve prices one amount, so the call to <Go>Value goes where
	// that amount sits and every other amount is a literal.
	value := row.Go + "Value(talentPoints)"

	if len(row.Stats) > 0 {
		b.WriteString("Stats: []StatConfig{\n")
		for i, stat := range row.Stats {
			amount := formatFloat(stat.Amount)
			if i == 0 && rendered.HasValue && !rendered.ValueOnPseudo {
				amount = value
			}
			fmt.Fprintf(&b, "{stats.%s, %s, %t},\n", stat.Stat.StatName(), amount, stat.Multiplicative)
		}
		b.WriteString("},\n")
	}
	if len(row.Pseudo) > 0 {
		b.WriteString("Pseudo: []PseudoConfig{\n")
		for i, mod := range row.Pseudo {
			amount := formatFloat(mod.Amount)
			if i == 0 && rendered.HasValue && rendered.ValueOnPseudo {
				amount = value
			}
			fmt.Fprintf(&b, "{PseudoStat%s, %s, %t, %d},\n",
				mod.Kind, amount, mod.Multiplicative, mod.SchoolMask)
		}
		b.WriteString("},\n")
	}
	b.WriteString("}")
	return b.String()
}

// The apply block: the condition the proto field is read by, and the call that
// puts the buff on the unit. A kind whose behaviour is a cooldown, a proc or an
// uptime calls a driver of a fixed name that sim/core/buffs_manual.go declares,
// and so does a row the manifest marks as driven. A driver is handed the whole
// scope message rather than its own field, because a driven buff often reads a
// second one: Grace of Air is 9 seconds long while the party is twisting totems.
func buffApply(row ResolvedBuff) (string, string, bool) {
	unit, field := "char", "individual"
	switch row.Scope {
	case buffmanifest.ScopeRaid:
		field = "raid"
	case buffmanifest.ScopeParty:
		field = "party"
	case buffmanifest.ScopeDebuff:
		field, unit = "debuffs", "target"
	}
	access := field + "." + row.GoField()

	var cond, points string
	switch row.Proto {
	case buffmanifest.ProtoBool:
		cond, points = access, "0"
	case buffmanifest.ProtoTristate:
		cond = access + " != proto.TristateEffect_TristateEffectMissing"
		points = fmt.Sprintf("GetTristateValueInt32(%s, 0, %d)", access, row.MaxTalentPoints())
	case buffmanifest.ProtoInt32, buffmanifest.ProtoDouble:
		cond, points = access+" > 0", "0"
	default:
		return "", "", false
	}

	var body string
	switch {
	case row.Driver,
		row.Kind == buffmanifest.KindExternalCD, row.Kind == buffmanifest.KindProc,
		row.Kind == buffmanifest.KindManual, row.Kind == buffmanifest.KindDebuffUptime,
		row.Kind == buffmanifest.KindItemCount:
		scope := field
		if row.Scope == buffmanifest.ScopeDebuff {
			scope = "debuffs, raid"
		}
		body = fmt.Sprintf("drive%s(%s, %s)", row.Go, unit, scope)
	default:
		target := "&char.Unit"
		if row.Scope == buffmanifest.ScopeDebuff {
			target = "target"
		}
		body = fmt.Sprintf("MakePermanent(%sAura(%s, false, %s))", row.Go, target, points)
	}
	return cond, body, true
}

func buffSchoolName(mask int32) string {
	name := schoolName(mask)
	if name == "" {
		return "SpellSchoolNone"
	}
	return strings.ReplaceAll(name, "core.", "")
}

func floatSliceLiteral(values []float64) string {
	parts := make([]string, len(values))
	for i, v := range values {
		parts[i] = formatFloat(v)
	}
	return "[]float64{" + strings.Join(parts, ", ") + "}"
}

func durationSliceLiteral(values []float64) string {
	parts := make([]string, len(values))
	for i, v := range values {
		parts[i] = fmt.Sprintf("%d * time.Millisecond", int64(v))
	}
	return "[]time.Duration{" + strings.Join(parts, ", ") + "}"
}

func formatFloat(v float64) string {
	if v == math.Trunc(v) {
		return strconv.FormatFloat(v, 'f', 1, 64)
	}
	return strconv.FormatFloat(v, 'f', -1, 64)
}

// GenerateBuffFiles writes both generated files and prints what the run made of
// the manifest.
func GenerateBuffFiles(helper *DBHelper) error {
	rows, err := ResolveBuffManifest(helper)
	if err != nil {
		return err
	}
	files, err := renderBuffFiles(rows)
	if err != nil {
		return err
	}

	root, err := repoRoot()
	if err != nil {
		return err
	}
	for name, out := range files {
		if err := os.WriteFile(filepath.Join(root, name), out, 0644); err != nil {
			return fmt.Errorf("writing %s: %w", name, err)
		}
	}

	printBuffSummary(rows)
	return nil
}

func printBuffSummary(rows []ResolvedBuff) {
	var generated, handWritten, manual, ghost, absent, unsupported int
	for _, row := range rows {
		switch {
		case row.HandWritten:
			handWritten++
		case row.Kind == buffmanifest.KindAbsent:
			absent++
		case row.Kind == buffmanifest.KindManual:
			manual++
		case row.Supported:
			generated++
		default:
			// The client describes the row, this generator cannot express it,
			// which is not the same thing as the client not having it.
			unsupported++
		}
		if row.Talent != nil && len(row.TalentCurve) == 0 {
			ghost++
		}
	}

	fmt.Printf("buffs: resolved %d rows: %d generated, %d hand-written, %d manual, %d ghost talents, %d absent, %d unsupported\n",
		len(rows), generated, handWritten, manual, ghost, absent, unsupported)

	var stale int
	for _, row := range rows {
		if row.ProtoStale {
			stale++
		}
		for _, warning := range row.Warnings {
			fmt.Printf("  %s: %s\n", row.Field, warning)
		}
		if row.DBName != "" && buffLabel(row) != row.DBName {
			fmt.Printf("  %s: label drift: the row is labelled %q, the client names spell %d %q\n",
				row.Field, buffLabel(row), row.SpellID, row.DBName)
		}
	}
	fmt.Printf("buffs: %d rows wait on the proto field being retyped\n", stale)
}
