package database

// Resolves tools/database/buffmanifest against the spell store's captured client
// rows and renders sim/core/buffs/buffs_auto_gen.go and
// sim/core/buffs/debuffs_auto_gen.go from the result, in the same pass that
// renders the store.
//
// The manifest names the spell each proto field reads; which of its effects are
// which stats, its stacks, its timing and the talent that improves it are read
// here, and the generated constructors read the values themselves off the store
// at runtime. A row the generator cannot express in the support API renders as a
// commented shell carrying the reason.

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"text/template"

	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/spelldata"
	"github.com/wowsims/forever/sim/core/stats"
	"github.com/wowsims/forever/tools/database/buffmanifest"
	"github.com/wowsims/forever/tools/database/dbc"
)

const buffsGenFile = "sim/core/buffs/buffs_auto_gen.go"
const debuffsGenFile = "sim/core/buffs/debuffs_auto_gen.go"

// The combo points a finisher the raid config simply has on the target is cast
// with.
const maxComboPoints = 5

// StatAmount is one stat the aura grants: its amount at level 60, and Expr, the
// Go expression the generated file reads that amount from the store with.
type StatAmount struct {
	Stat           stats.Stat
	Amount         float64
	Expr           string
	Multiplicative bool
}

// PseudoMod is one PseudoStats field the aura modifies.
type PseudoMod struct {
	Kind           string // PseudoStats field name, e.g. "ThreatMultiplier"
	Amount         float64
	Expr           string
	Multiplicative bool
	SchoolMask     int32 // client school bits, 0 = every school
}

// ResolvedEffect is one SpellEffect row of the anchor spell, with its value
// already derived for level 60 and truncated toward zero. Ref is the Go
// expression that reaches the effect in the store, and Amount the one that reads
// its value in the client's units.
type ResolvedEffect struct {
	Index          int32
	Effect         dbc.SpellEffectType
	Aura           dbc.EffectAuraType
	Misc           int32
	Value          float64
	PerResource    float64
	PeriodMs       int32
	ImplicitTarget dbc.ImplicitTarget
	Ref            string
	Amount         string
}

// ResolvedBuff is one manifest row plus everything the database states about it.
type ResolvedBuff struct {
	buffmanifest.BuffSpec

	SpellID     int32 // the spell the aura's numbers are read from
	CastSpellID int32 // the cast, where the manifest pins one; SpellID otherwise
	DurationMs  int32 // -1 or 0 never expires
	CooldownMs  int32 // 0 when the client states none; read from the cast
	MaxStacks   int32
	SchoolMask  int32
	Effects     []ResolvedEffect

	Stats  []StatAmount
	Pseudo []PseudoMod

	// OverrideStats is the manifest's StatOverride resolved to sim stats.
	OverrideStats []stats.Stat

	// TalentRanks is how many points the improving talent takes, 0 for a row
	// no talent prices.
	TalentRanks   int32
	TalentApplies buffmanifest.TalentApplies

	// TalentOnPseudo says the talent prices Pseudo[0] rather than Stats[0].
	TalentOnPseudo bool
	// TalentSpellID is the spell of the trait node that prices the improvement,
	// which is the icon the UI shows for the improved state. TalentPosition is
	// the effect of that spell the improvement is read from, counted the way
	// EffectN counts.
	TalentSpellID  int32
	TalentPosition int32

	// DurationFromCast says the aura states no duration and the cast times it.
	DurationFromCast bool

	ScopeFromClient buffmanifest.BuffScope

	DBName string

	// Note is what the generated file says about the row above its constructor,
	// for a value the client states in a way the row had to be read through.
	Note string

	Supported bool
	Reason    string
	Warnings  []string
}

func (r *ResolvedBuff) warn(format string, args ...any) {
	r.Warnings = append(r.Warnings, fmt.Sprintf(format, args...))
}

func (r *ResolvedBuff) unsupported(format string, args ...any) {
	r.Supported = false
	r.Reason = fmt.Sprintf(format, args...)
}

// Every manifest row, read out of the client rows the store is built from.
func resolveBuffManifest(in *storeInputs) ([]ResolvedBuff, error) {
	t := in.tables()
	rows := make([]ResolvedBuff, 0, len(buffmanifest.Manifest))
	for _, spec := range buffmanifest.Manifest {
		row, err := resolveBuff(t, in.TraitNodes, spec)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", spec.Field, err)
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func resolveBuff(t *spellTables, nodes []traitNode, spec buffmanifest.BuffSpec) (ResolvedBuff, error) {
	row := ResolvedBuff{BuffSpec: spec, Supported: true, ScopeFromClient: spec.Scope}

	compiled, err := compiledProtoType(spec)
	if err != nil {
		return row, err
	}

	// A pet inherits the buff by finding the aura on its owner, so the row has
	// to name it whether or not the client describes the buff at all.
	if spec.Pet == buffmanifest.PetInheritOwnerAura && spec.Label == "" && spec.Name == "" && spec.Category == "" {
		return row, fmt.Errorf("states %s but names no aura for the pet to find on its owner", spec.Pet)
	}
	// A pet policy is about what its owner's buffs reach it, which a debuff is
	// not, and applyGeneratedPetBuffs never sees the debuff message.
	if spec.Scope == buffmanifest.ScopeDebuff && spec.Pet != buffmanifest.PetNormal {
		return row, fmt.Errorf("states %s, which only a buff row can", spec.Pet)
	}

	if spec.Kind == buffmanifest.KindAbsent || spec.Kind == buffmanifest.KindFlag {
		row.unsupported("%s", spec.Notes)
		return row, nil
	}

	if err := loadBuffSpell(t, &row); err != nil {
		return row, err
	}
	if err := resolveStatOverride(&row); err != nil {
		return row, err
	}
	mapEffects(&row)
	if err := resolveTalent(t, nodes, &row); err != nil {
		return row, err
	}
	validateScope(&row)

	if compiled != buffProtoTypeNames[spec.Proto] && row.Supported {
		row.unsupported("proto field not yet retyped: the sim compiled against %s, the manifest declares %s",
			compiled, buffProtoTypeNames[spec.Proto])
	}
	return row, nil
}

// The spell's name, school, duration, stacks and effects, and the cooldown of a
// buff other players cast on their own: the shared timer the sim hands the next
// source is the cast's cooldown, and a totem's aura states neither how long the
// totem stands nor when the next one may be dropped. Mana Tide is that row -
// 17360 carries the mana, the cast 17359 the 13 seconds and the 5 minutes.
func loadBuffSpell(t *spellTables, row *ResolvedBuff) error {
	if row.SpellID = row.BuffSpec.SpellID; row.SpellID == 0 {
		return fmt.Errorf("names no spell")
	}
	row.CastSpellID = row.SpellID
	if row.CastID != 0 {
		row.CastSpellID = row.CastID
	}
	for _, id := range []int32{row.SpellID, row.CastSpellID} {
		if _, named := t.Names[id]; !named {
			return fmt.Errorf("spell %d is no spell in this client", id)
		}
	}

	spell := t.row(row.SpellID)
	row.DBName = spell.Name
	row.SchoolMask = int32(spell.School)
	row.DurationMs = spell.DurationMs
	row.MaxStacks = int32(spell.MaxStack)

	for _, e := range spell.Effects {
		value := (&spelldata.Effect{BasePoints: e.BasePoints, PPL: e.PPL,
			SpellLevel: e.SpellLevel, MaxLevel: e.MaxLevel}).Average(RankLevel)
		row.Effects = append(row.Effects, ResolvedEffect{
			Index:          int32(e.Index),
			Effect:         e.Type,
			Aura:           e.Aura,
			Misc:           e.Misc,
			Value:          value,
			PerResource:    e.PointsPerResource,
			PeriodMs:       e.PeriodMs,
			ImplicitTarget: dbc.ImplicitTarget(e.Target[0]),
		})
	}
	setEffectRefs(row)

	if row.Kind != buffmanifest.KindExternalCD {
		return nil
	}
	cast := t.row(row.CastSpellID)
	row.CooldownMs = max(cast.CooldownMs, cast.CategoryCooldownMs)
	if row.CooldownMs == 0 {
		row.unsupported("spell %d states no cooldown, which an external cooldown is scheduled by", row.CastSpellID)
		return nil
	}
	if row.DurationMs <= 0 && row.CastSpellID != row.SpellID {
		row.DurationFromCast = true
		row.DurationMs = cast.DurationMs
	}
	if row.DurationMs <= 0 {
		row.unsupported("spell %d states no duration, which the external cooldown's aura needs", row.SpellID)
	}
	return nil
}

// The improving talent the manifest pins, which the store keeps as a ladder of
// its ranks. Its effect has to be a modifier whose class mask reaches the buff's
// spell; one that modifies misc 1 scales the duration, anything else the value.
func resolveTalent(t *spellTables, nodes []traitNode, row *ResolvedBuff) error {
	if row.Talent == nil {
		if row.Proto == buffmanifest.ProtoTristate && row.ImpAction == nil {
			return fmt.Errorf("declared ProtoTristate but the manifest names neither a talent nor an ImpAction")
		}
		return nil
	}
	if row.Proto != buffmanifest.ProtoTristate {
		return fmt.Errorf("names talent %q but is not ProtoTristate", row.Talent.Name)
	}

	id := row.Talent.SpellID
	node := slices.IndexFunc(nodes, func(n traitNode) bool { return n.SpellID == id })
	if node < 0 {
		return fmt.Errorf("talent %d is no node of a class tree", id)
	}
	talent := t.row(id)
	if talent.Name != row.Talent.Name {
		return fmt.Errorf("talent %d is %q, the manifest says %q", id, talent.Name, row.Talent.Name)
	}
	position := slices.IndexFunc(talent.Effects, func(e storeEffect) bool { return int32(e.Index) == row.Talent.Effect })
	if position < 0 {
		return fmt.Errorf("talent %d has no effect %d", id, row.Talent.Effect)
	}
	mod := talent.Effects[position]
	if mod.Aura != dbcenums.A_ADD_FLAT_MODIFIER && mod.Aura != dbcenums.A_ADD_PCT_MODIFIER {
		return fmt.Errorf("talent %d effect %d is aura %d, not a spell modifier", id, row.Talent.Effect, mod.Aura)
	}
	if !mod.ClassFlags.Matches(t.row(row.SpellID).ClassFlags) {
		return fmt.Errorf("talent %d effect %d does not reach spell %d", id, row.Talent.Effect, row.SpellID)
	}

	row.TalentSpellID = id
	row.TalentPosition = int32(position + 1)
	row.TalentApplies = row.Talent.Applies
	if mod.Misc == int32(dbcenums.SPELLMOD_DURATION) {
		row.TalentApplies = buffmanifest.TalentScalesDuration
	}
	if row.TalentApplies != buffmanifest.TalentScalesDuration {
		_, onPseudo, ok := row.talentTarget()
		if !ok {
			row.warn("talent %q has nothing to scale: spell %d states no amount this generator maps", row.Talent.Name, row.SpellID)
			return nil
		}
		row.TalentOnPseudo = onPseudo
	}
	row.TalentRanks = nodes[node].MaxRanks
	return nil
}

// Who the client says the aura reaches. The manifest wins - the sim's scopes are
// a UI grouping as much as a game fact - but a disagreement is worth printing.
func validateScope(row *ResolvedBuff) {
	scope, known := clientScope(*row)
	if !known {
		return
	}
	row.ScopeFromClient = scope
	if scope != row.Scope {
		row.warn("scope mismatch: the manifest says %s, spell %d reads as %s", row.Scope, row.SpellID, scope)
	}
}

// The group the client states a buff reaches: an area aura names it in the
// effect itself, and an aura the client applies over an area or on one ally
// names it in the effect's target.
func clientScope(row ResolvedBuff) (buffmanifest.BuffScope, bool) {
	for _, effect := range row.Effects {
		switch effect.Effect {
		case dbcenums.E_APPLY_AREA_AURA_RAID:
			return buffmanifest.ScopeRaid, true
		case dbcenums.E_APPLY_AREA_AURA_PARTY:
			return buffmanifest.ScopeParty, true
		case dbcenums.E_APPLY_AURA:
			switch effect.ImplicitTarget {
			case dbc.TARGET_UNIT_CASTER_AREA_RAID:
				return buffmanifest.ScopeRaid, true
			case dbc.TARGET_UNIT_CASTER_AREA_PARTY:
				return buffmanifest.ScopeParty, true
			case dbc.TARGET_UNIT_TARGET_ALLY, dbc.TARGET_UNIT_TARGET_ALLY_OR_RAID:
				return buffmanifest.ScopeIndividual, true
			}
		}
	}
	return buffmanifest.ScopeIndividual, false
}

// The repository root, found by walking up from the working directory until
// go.mod appears. The generator runs from the root and the regeneration test
// from tools/database, and both have to reach sim/core/buffs.
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
	buffmanifest.ProtoBool:     "bool",
	buffmanifest.ProtoTristate: "proto.TristateEffect",
	buffmanifest.ProtoInt32:    "int32",
	buffmanifest.ProtoDouble:   "float64",
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

func isAuraApplication(effect dbc.SpellEffectType) bool {
	return effect == dbcenums.E_APPLY_AURA || effect == dbcenums.E_APPLY_AREA_AURA_PARTY ||
		effect == dbcenums.E_APPLY_AREA_AURA_RAID
}

// How the generated file reaches each effect: by its aura and misc value where
// no other effect of the spell shares them, and by position where one does.
func setEffectRefs(row *ResolvedBuff) {
	for i := range row.Effects {
		e := &row.Effects[i]
		shared := 0
		for _, other := range row.Effects {
			if other.Aura == e.Aura && other.Misc == e.Misc {
				shared++
			}
		}
		name, named := dbcenums.Named(e.Aura)
		if shared == 1 && e.Aura != 0 && named {
			e.Ref = fmt.Sprintf("%s.Effect(dbcenums.%s, %d)", spellVar(row.Go), name, e.Misc)
		} else {
			e.Ref = fmt.Sprintf("%s.EffectN(%d)", spellVar(row.Go), i+1)
		}
		e.Amount = "amount(" + e.Ref + ")"
	}
}

func spellVar(stem string) string {
	return lowerFirst(stem) + "Spell"
}

func castVar(stem string) string {
	return lowerFirst(stem) + "Cast"
}

func talentVar(stem string) string {
	return lowerFirst(stem) + "Talent"
}

func lowerFirst(s string) string {
	return strings.ToLower(s[:1]) + s[1:]
}

// The stats or pseudo-stats the row's kind says its effects are. An aura this
// cannot express is not an error: the row becomes a shell naming the aura, which
// is the signal that either the support API or the manifest has to grow.
func mapEffects(row *ResolvedBuff) {
	var unmapped []string
	damageShield := false
	for i := range row.Effects {
		e := row.Effects[i]
		if !isAuraApplication(e.Effect) {
			continue
		}
		if e.Aura == dbcenums.A_DAMAGE_SHIELD {
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
			row.Effects[i].Amount = "fullComboPoints(" + e.Ref + ")"
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

	// The sim puts every resistance stat into its school's category by itself, so
	// a row whose manifest category is that school says the same thing twice and
	// has no exclusivity of its own beyond it.
	if isSchoolResistanceCategory(row.Category) {
		row.Category = ""
	}

	switch row.Kind {
	case buffmanifest.KindDamageShield:
		if !damageShield {
			row.unsupported("no A_DAMAGE_SHIELD effect on spell %d", row.SpellID)
			return
		}
	case buffmanifest.KindItemCount:
		// The count multiplies every amount, which a multiplier cannot be read
		// through: two of the same staff would raise a stat by twice its factor.
		for _, stat := range row.Stats {
			if stat.Multiplicative {
				row.unsupported("%s is a multiplier, which a count of items cannot scale", stat.Stat.StatName())
				return
			}
		}
		for _, mod := range row.Pseudo {
			if mod.Multiplicative {
				row.unsupported("%s is a multiplier, which a count of items cannot scale", mod.Kind)
				return
			}
		}
	case buffmanifest.KindExternalCD, buffmanifest.KindProc, buffmanifest.KindManual,
		buffmanifest.KindDebuffUptime:
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
		out = append(out, StatAmount{Stat: stat, Amount: e.Value, Expr: e.Amount})
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
		return []StatAmount{{Stat: stat, Amount: amount, Expr: e.Amount}}, true
	}
	each := func(sts []stats.Stat, amount float64, expr string, multiplicative bool) ([]StatAmount, bool) {
		var out []StatAmount
		for _, stat := range sts {
			out = append(out, StatAmount{Stat: stat, Amount: amount, Expr: expr, Multiplicative: multiplicative})
		}
		return out, true
	}
	mainStats := []stats.Stat{stats.Strength, stats.Agility, stats.Stamina, stats.Intellect, stats.Spirit}
	multiplier, multiplierExpr := 1+e.Value/100, "1 + "+e.Amount+"/100"

	switch e.Aura {
	case dbcenums.A_MOD_STAT:
		if e.Misc == -1 {
			return each(mainStats, e.Value, e.Amount, false)
		}
		stat, ok := dbc.MapMainStatToStat(int(e.Misc))
		if !ok {
			return nil, false
		}
		return flat(stats.Stat(stat), e.Value)
	case dbcenums.A_MOD_ATTACK_POWER:
		return flat(stats.AttackPower, e.Value)
	case dbcenums.A_MOD_RANGED_ATTACK_POWER:
		return flat(stats.RangedAttackPower, e.Value)
	case dbcenums.A_MOD_INCREASE_HEALTH:
		return flat(stats.Health, e.Value)
	case dbcenums.A_MOD_HEALING_DONE:
		return flat(stats.HealingPower, e.Value)
	case dbcenums.A_MOD_POWER_REGEN:
		// The client states mana per 5 seconds directly on this aura.
		return flat(stats.MP5, e.Value)
	case dbcenums.A_PERIODIC_ENERGIZE:
		// The client's amount is already whole; what the conversion to mana per
		// five seconds produces is not, and truncating it would lose part of a
		// tick the aura really restores.
		if e.PeriodMs <= 0 {
			return nil, false
		}
		return []StatAmount{{Stat: stats.MP5, Amount: e.Value * 5000 / float64(e.PeriodMs),
			Expr: fmt.Sprintf("manaPerFive(%s, %s)", e.Ref, e.Amount)}}, true
	case dbcenums.A_MOD_SPELL_CRIT_CHANCE:
		return flat(stats.SpellCritPercent, e.Value)
	case dbcenums.A_MOD_HIT_CHANCE:
		return flat(stats.PhysicalHitPercent, e.Value)
	case dbcenums.A_MOD_EXPERTISE:
		return flat(stats.ExpertisePercent, e.Value)
	case dbcenums.A_MOD_RESISTANCE:
		var sts []stats.Stat
		for bit, stat := range resistanceBits {
			if e.Misc&bit != 0 {
				sts = append(sts, stat)
			}
		}
		slices.Sort(sts)
		// Bit 2 is Holy, which has no resistance stat in the sim; a mask of only
		// that bit resolves to nothing rather than to a wrong stat.
		if len(sts) == 0 {
			return nil, false
		}
		return each(sts, e.Value, e.Amount, false)
	case dbcenums.A_MOD_DAMAGE_DONE, dbcenums.A_MOD_RATING:
		stat := dbc.ConvertEffectAuraToStatIndex(e.Aura, int(e.Misc))
		if stat < 0 {
			return nil, false
		}
		return flat(stats.Stat(stat), e.Value)
	case dbcenums.A_MOD_TOTAL_STAT_PERCENTAGE:
		if e.Misc == -1 {
			return each(mainStats, multiplier, multiplierExpr, true)
		}
		stat, ok := dbc.MapMainStatToStat(int(e.Misc))
		if !ok {
			return nil, false
		}
		return each([]stats.Stat{stats.Stat(stat)}, multiplier, multiplierExpr, true)
	case dbcenums.A_MOD_ATTACK_POWER_PCT:
		return each([]stats.Stat{stats.AttackPower}, multiplier, multiplierExpr, true)
	}
	return nil, false
}

// Whether a manifest category names a resistance school, which is the category
// registerGeneratedSchoolResistances puts that school's stat into anyway. The
// names are sim/core's ResistanceCategory* constants.
func isSchoolResistanceCategory(category string) bool {
	switch category {
	case "ResistanceArcane", "ResistanceFire", "ResistanceFrost", "ResistanceNature", "ResistanceShadow":
		return true
	}
	return false
}

// Holy, fire, nature, frost, shadow and arcane together, which is every school
// the sim counts as spell damage.
const everySpellSchoolMask int32 = 126

// The spell schools plus physical, which is everything a unit can deal.
const everySchoolMask int32 = 127

var resistanceBits = map[int32]stats.Stat{
	1:  stats.Armor,
	4:  stats.FireResistance,
	8:  stats.NatureResistance,
	16: stats.FrostResistance,
	32: stats.ShadowResistance,
	64: stats.ArcaneResistance,
}

func pseudoModsOf(e ResolvedEffect) ([]PseudoMod, bool) {
	multiplier, multiplierExpr := 1+e.Value/100, "1 + "+e.Amount+"/100"
	multiply := func(kind string, schoolMask int32) ([]PseudoMod, bool) {
		return []PseudoMod{{Kind: kind, Amount: multiplier, Expr: multiplierExpr, Multiplicative: true, SchoolMask: schoolMask}}, true
	}
	add := func(kind string) ([]PseudoMod, bool) {
		return []PseudoMod{{Kind: kind, Amount: e.Value, Expr: e.Amount}}, true
	}

	switch e.Aura {
	case dbcenums.A_MOD_THREAT:
		return multiply("ThreatMultiplier", 0)
	case dbcenums.A_MOD_DAMAGE_PERCENT_DONE:
		// A mask of every school, physical included, raises everything the unit
		// deals; anything narrower is per school, so a buff the client states
		// for the magic schools does not raise a melee swing.
		if e.Misc == everySchoolMask {
			return multiply("DamageDealtMultiplier", 0)
		}
		// A mask of 0 names no school, so the effect raises nothing: emitting it
		// would be a no-op the reader of the generated file has to work out.
		if e.Misc == 0 {
			return nil, false
		}
		return multiply("SchoolDamageDealtMultiplier", e.Misc)
	case dbcenums.A_MOD_HEALING_DONE_PERCENT:
		return multiply("HealingDealtMultiplier", 0)
	case dbcenums.A_REDUCE_PUSHBACK:
		// PseudoStats.PushbackChance is the chance of being pushed back and
		// starts at 1, so the client's "35% less pushback" is -0.35 there.
		return []PseudoMod{{Kind: "PushbackChance", Amount: -e.Value / 100, Expr: "-" + e.Amount + "/100"}}, true
	case dbcenums.A_MOD_MELEE_HASTE_3, dbcenums.A_MOD_ATTACKSPEED:
		return multiply("MeleeSpeedMultiplier", 0)
	case dbcenums.A_MOD_DAMAGE_PERCENT_TAKEN:
		return multiply("SchoolDamageTakenMultiplier", e.Misc)
	case dbcenums.A_RANGED_ATTACK_POWER_ATTACKER_BONUS:
		return add("BonusRangedAttackPower")
	case dbcenums.A_MOD_DAMAGE_TAKEN:
		// The sim splits flat damage taken into a physical and a spell field,
		// so the school mask picks which one the effect is. A mask that names
		// some spell schools and not others has neither: the spell field would
		// raise what every school does to the target.
		if e.Misc&1 != 0 {
			return add("BonusPhysicalDamageTaken")
		}
		if e.Misc&everySpellSchoolMask == everySpellSchoolMask {
			return add("BonusSpellDamageTaken")
		}
		return nil, false
	}
	return nil, false
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
		if e.Aura == dbcenums.A_DAMAGE_SHIELD {
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
	HasSpell      bool
	SpellVar      string
	CastVar       string
	CastID        int32
	TalentVar     string
	TalentID      int32
	TalentRanks   int32
	Category      string
	CategoryVar   string
	ValueExpr     string
	HasValue      bool
	ValueOnPseudo bool
	Duration      string
	Cooldown      string
	HasCooldown   bool
	ExtraParams   string
	OwnerAura     string
	OwnerAuraVar  string
	Constructor   string
	ApplyIf       string
	ApplyBody     string
	HasApply      bool
}

// Every file the manifest renders, by the path it is written to: the two Go
// files, the settings inputs and the proto messages.
func renderBuffOutputs(in *storeInputs) (map[string][]byte, error) {
	rows, err := resolveBuffManifest(in)
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		for _, warning := range row.Warnings {
			fmt.Fprintf(progress, "buffs: %s: %s\n", row.Field, warning)
		}
	}

	files, err := renderBuffFiles(rows)
	if err != nil {
		return nil, err
	}
	if files[buffsDebuffsTSFile], err = RenderBuffsDebuffsTS(rows); err != nil {
		return nil, err
	}
	files[buffsProtoFile] = buffmanifest.RenderProto(buffmanifest.Manifest)
	return files, nil
}

const buffsProtoFile = "proto/buffs.proto"

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
	var shared []sharedCategoryRow
	needsTime, needsStats, needsEnums := false, false, false
	for _, row := range resolved {
		if (row.Scope == buffmanifest.ScopeDebuff) != debuffs {
			continue
		}
		rendered := renderRow(row)
		if rendered.Supported {
			needsTime = true
			needsStats = needsStats || len(row.Stats) > 0
			needsEnums = needsEnums || strings.Contains(rendered.ValueExpr+rendered.Constructor, "dbcenums.")
			if row.SharedCategory != "" && !slices.ContainsFunc(shared, func(c sharedCategoryRow) bool { return c.Name == row.SharedCategory }) {
				shared = append(shared, sharedCategoryRow{Var: sharedCategoryVar(row.SharedCategory), Name: row.SharedCategory})
			}
		}
		rows = append(rows, rendered)
	}
	slices.SortFunc(shared, func(a, b sharedCategoryRow) int { return strings.Compare(a.Name, b.Name) })

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
		"Rows": rows, "NeedsTime": needsTime, "NeedsStats": needsStats, "NeedsEnums": needsEnums,
		"SharedCategories": shared, "PetRows": petBuffRows(resolved),
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

// petBuffRow is one line of applyGeneratedPetBuffs: the field the policy reads,
// what the field is set to when the pet does not get the buff, and the aura the
// pet inherits by standing next to its owner.
type petBuffRow struct {
	Access    string
	Zero      string
	OwnerAura string
	Strip     bool
	Inherit   bool
	StripLate bool
}

// Every row whose pet policy says something, in manifest order. A row the
// generator could not express is included too: the policy is about the proto
// field, which the hand-written apply block reads just the same.
func petBuffRows(resolved []ResolvedBuff) []petBuffRow {
	var rows []petBuffRow
	for _, row := range resolved {
		if row.Scope == buffmanifest.ScopeDebuff || row.Pet == buffmanifest.PetNormal {
			continue
		}
		out := petBuffRow{
			Access:    buffScopeField(row.Scope) + "." + row.GoField(),
			Zero:      buffProtoZero(row.Proto),
			Strip:     row.Pet == buffmanifest.PetStrip,
			Inherit:   row.Pet == buffmanifest.PetInheritOwnerAura,
			StripLate: row.Pet == buffmanifest.PetStripWhenSummonedLate,
		}
		if out.Inherit {
			out.OwnerAura = petOwnerAuraVar(row)
		}
		rows = append(rows, out)
	}
	return rows
}

// The label of the aura the owner carries, which for a row the client does not
// describe is the name its hand-written constructor registers, and the manifest
// keeps that as the row's exclusive category.
func petOwnerAura(row ResolvedBuff) string {
	if label := buffLabel(row); label != "" {
		return label
	}
	return row.Category
}

// sharedCategoryRow is one `var <Name>Category = "<Name>"` line: several rows
// name the same shared category, so the generated file declares it once.
type sharedCategoryRow struct {
	Var  string
	Name string
}

// The identifier the generated file gives a shared category, which every row
// that joins it and every caller in sim/core/buffs name.
func sharedCategoryVar(category string) string {
	return category + "Category"
}

// The identifier the generated file gives that label, which the buff's own
// hand-written constructor names as well.
func petOwnerAuraVar(row ResolvedBuff) string {
	return row.Go + "AuraLabel"
}

func buffProtoZero(protoType buffmanifest.BuffProtoType) string {
	switch protoType {
	case buffmanifest.ProtoBool:
		return "false"
	case buffmanifest.ProtoTristate:
		return "proto.TristateEffect_TristateEffectMissing"
	}
	return "0"
}

func buffScopeField(scope buffmanifest.BuffScope) string {
	switch scope {
	case buffmanifest.ScopeRaid:
		return "raid"
	case buffmanifest.ScopeParty:
		return "party"
	case buffmanifest.ScopeDebuff:
		return "debuffs"
	}
	return "individual"
}

func renderRow(row ResolvedBuff) buffRow {
	out := buffRow{
		Go: row.Go, Field: row.Field, Label: buffLabel(row),
		SpellID: row.SpellID, Kind: row.Kind.String(), Reason: row.Reason, Note: row.Note,
		Supported: row.Supported, HasSpell: row.SpellID != 0,
	}
	if row.Pet == buffmanifest.PetInheritOwnerAura {
		out.OwnerAura, out.OwnerAuraVar = strconv.Quote(petOwnerAura(row)), petOwnerAuraVar(row)
	}

	// A row whose amounts are worth one item each takes the number of them, so
	// that a party with three of the same staff gets three times the aura.
	if row.Kind == buffmanifest.KindItemCount {
		out.ExtraParams = ", count float64"
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
	out.SpellVar = spellVar(row.Go)
	timing := out.SpellVar
	if row.CastSpellID != row.SpellID && (row.DurationFromCast || row.CooldownMs > 0) {
		out.CastVar, out.CastID = castVar(row.Go), row.CastSpellID
		timing = out.CastVar
	}
	if row.TalentRanks > 0 {
		out.TalentVar, out.TalentID, out.TalentRanks = talentVar(row.Go), row.TalentSpellID, row.TalentRanks
	}
	out.ValueExpr, out.ValueOnPseudo, out.HasValue = buffValueExpr(row, out)
	out.Duration = buffDurationExpr(row, out)
	if row.DurationFromCast {
		out.Duration = "auraDuration(" + out.CastVar + ")"
	}
	if row.CooldownMs > 0 {
		out.Cooldown, out.HasCooldown = "cooldown("+timing+")", true
	}
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

// The first stat amount, talent-scaled when the tree prices the talent: the
// talent's rank reads its modifier off the ladder the store keeps for it.
func buffValueExpr(row ResolvedBuff, rendered buffRow) (string, bool, bool) {
	if rendered.TalentVar != "" && row.TalentApplies != buffmanifest.TalentScalesDuration {
		if target, onPseudo, ok := row.talentTarget(); ok {
			target.Amount = fmt.Sprintf("talentScaled(%s, %s)", target.Amount, talentModifier(row, rendered))
			return row.convertedExpr(target, onPseudo), onPseudo, true
		}
	}
	if len(row.Stats) > 0 {
		return row.Stats[0].Expr, false, true
	}
	if len(row.Pseudo) > 0 {
		return row.Pseudo[0].Expr, true, true
	}
	for _, e := range row.Effects {
		if e.Aura == dbcenums.A_DAMAGE_SHIELD {
			return e.Amount, false, true
		}
	}
	return "", false, false
}

func talentModifier(row ResolvedBuff, rendered buffRow) string {
	return fmt.Sprintf("%s.Rank(talentPoints).EffectN(%d)", rendered.TalentVar, row.TalentPosition)
}

// What the effect is worth once the kind mapping has converted it, as the
// expression the generated file reads it with.
func (row ResolvedBuff) convertedExpr(e ResolvedEffect, onPseudo bool) string {
	if onPseudo {
		if mods, ok := pseudoModsOf(e); ok {
			return mods[0].Expr
		}
		return e.Amount
	}
	if amounts, ok := row.statAmounts(e); ok {
		return amounts[0].Expr
	}
	return e.Amount
}

func buffDurationExpr(row ResolvedBuff, rendered buffRow) string {
	if rendered.TalentVar != "" && row.TalentApplies == buffmanifest.TalentScalesDuration {
		return fmt.Sprintf("talentScaledDuration(%s, %s)", rendered.SpellVar, talentModifier(row, rendered))
	}
	return "auraDuration(" + rendered.SpellVar + ")"
}

// The call the constructor makes into the hand-written support API.
func buffConstructor(row ResolvedBuff, rendered buffRow) string {
	var b strings.Builder
	config := buffConfigLiteral(row, rendered)

	switch row.Kind {
	case buffmanifest.KindDamageShield:
		school := buffSchoolName(row.SchoolMask)
		fmt.Fprintf(&b, "return core.NewGeneratedDamageShield(unit, %s, %s, %s(talentPoints))",
			config, school, row.Go+"Value")
	default:
		if row.Scope == buffmanifest.ScopeDebuff {
			fmt.Fprintf(&b, "return core.NewGeneratedDebuff(unit, %s)", config)
		} else {
			fmt.Fprintf(&b, "return core.NewGeneratedStatAura(unit, %s)", config)
		}
	}
	return b.String()
}

func buffConfigLiteral(row ResolvedBuff, rendered buffRow) string {
	var b strings.Builder
	b.WriteString("core.GeneratedBuff{\n")
	fmt.Fprintf(&b, "Label: %q + core.Ternary(isPlayer, \"Player\", \"External\") + \")\",\n", rendered.Label+" (")
	fmt.Fprintf(&b, "ActionID: core.ActionID{SpellID: %s.ID}.WithTag(core.TernaryInt32(isPlayer, 0, -1)),\n", rendered.SpellVar)
	fmt.Fprintf(&b, "Duration: %sDuration(talentPoints),\n", row.Go)
	if row.MaxStacks > 0 {
		fmt.Fprintf(&b, "MaxStacks: %d,\n", row.MaxStacks)
	}
	if rendered.CategoryVar != "" {
		fmt.Fprintf(&b, "Category: %s,\n", rendered.CategoryVar)
	}
	if row.SharedCategory != "" {
		fmt.Fprintf(&b, "SharedCategory: %s,\n", sharedCategoryVar(row.SharedCategory))
	}
	if row.SingleAura {
		b.WriteString("SingleAura: true,\n")
	}
	b.WriteString("IsPlayer: isPlayer,\n")

	// The talent prices one amount, so the call to <Go>Value goes where that
	// amount sits and every other amount is read off its own effect.
	value := row.Go + "Value(talentPoints)"
	// Every amount of an item-count row is per item.
	scale := ""
	if rendered.ExtraParams != "" {
		scale = " * count"
	}

	if len(row.Stats) > 0 {
		b.WriteString("Stats: []core.StatConfig{\n")
		for i, stat := range row.Stats {
			amount := stat.Expr
			if i == 0 && rendered.HasValue && !rendered.ValueOnPseudo {
				amount = value
			}
			fmt.Fprintf(&b, "{Stat: stats.%s, Amount: %s%s, IsMultiplicative: %t},\n", stat.Stat.StatName(), amount, scale, stat.Multiplicative)
		}
		b.WriteString("},\n")
	}
	if len(row.Pseudo) > 0 {
		b.WriteString("Pseudo: []core.PseudoConfig{\n")
		for i, mod := range row.Pseudo {
			amount := mod.Expr
			if i == 0 && rendered.HasValue && rendered.ValueOnPseudo {
				amount = value
			}
			fmt.Fprintf(&b, "{Kind: core.PseudoStat%s, Amount: %s%s, IsMultiplicative: %t, SchoolMask: %d},\n",
				mod.Kind, amount, scale, mod.Multiplicative, mod.SchoolMask)
		}
		b.WriteString("},\n")
	}
	b.WriteString("}")
	return b.String()
}

// The apply block: the condition the proto field is read by, and the call that
// puts the buff on the unit. A kind whose behaviour is a cooldown, a proc or an
// uptime calls a driver of a fixed name that sim/core/buffs/drivers.go declares,
// and so does a row the manifest marks as driven. A driver is handed the whole
// scope message rather than its own field, because a driven buff often reads a
// second one: Grace of Air is 9 seconds long while the party is twisting totems.
func buffApply(row ResolvedBuff) (string, string, bool) {
	unit, field := "char", buffScopeField(row.Scope)
	if row.Scope == buffmanifest.ScopeDebuff {
		unit = "target"
	}
	access := field + "." + row.GoField()

	var cond, points string
	switch row.Proto {
	case buffmanifest.ProtoBool:
		cond, points = access, "0"
	case buffmanifest.ProtoTristate:
		cond = access + " != proto.TristateEffect_TristateEffectMissing"
		points = fmt.Sprintf("core.GetTristateValueInt32(%s, 0, %d)", access, row.TalentRanks)
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
		body = fmt.Sprintf("core.MakePermanent(%sAura(%s, false, %s))", row.Go, target, points)
	}
	return cond, body, true
}

func buffSchoolName(mask int32) string {
	name := schoolName(mask)
	if name == "" {
		return "core.SpellSchoolNone"
	}
	return name
}

func formatFloat(v float64) string {
	if v == math.Trunc(v) {
		return strconv.FormatFloat(v, 'f', 1, 64)
	}
	return strconv.FormatFloat(v, 'f', -1, 64)
}

