package database

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"slices"
	"sync"

	"github.com/wowsims/forever/tools/database/dbc"
)

const RankLevel = 60

// Rage is stored in tenths: Heroic Strike costs 150, not 15. Mana, energy and focus are not.
const powerTypeRage = 1

func NormalizePowerCost(cost int32, powerType int32) int32 {
	if powerType == powerTypeRage {
		return cost / 10
	}
	return cost
}

func IsWeaponDamageEffect(effect dbc.SpellEffectType) bool {
	return effect == dbc.E_WEAPON_DAMAGE_NOSCHOOL || effect == dbc.E_WEAPON_DAMAGE ||
		effect == dbc.E_NORMALIZED_WEAPON_DMG
}

// Devouring Plague ticks as a leech rather than as plain periodic damage, so "is this a DoT" cannot be
// a single aura check.
// The client states a flat threat amount on the abilities whose point is threat: Feint and Cower
// shed it, Distracting Shot adds it. 22 ranked spells across four families carry one.
func IsThreatEffect(effect dbc.SpellEffectType) bool {
	return effect == dbc.E_THREAT || effect == dbc.E_THREAT_ALL
}

func IsPeriodicAura(aura dbc.EffectAuraType) bool {
	return aura == dbc.A_PERIODIC_DAMAGE || aura == dbc.A_PERIODIC_LEECH
}

type RankEffect struct {
	Index        int32
	Effect       dbc.SpellEffectType
	Aura         dbc.EffectAuraType
	BasePoints   int32
	DieSides     int32
	PointsPerLvl float64
	Coefficient  float64
	APCoef       float64
	MiscValue    int32
	AuraPeriod   int32
	OwnerSpellID int32
}

type RankSpell struct {
	SpellID    int32
	SpellLevel int32
	MaxLevel   int32
	ManaCost   sql.NullInt64
	PowerType  int32
	DurationMs int32
	CastTimeMs int32
	GCDMs      int32
	CooldownMs int32
	MinRange   float64
	MaxRange   float64

	// How fast the projectile flies, in yards per second, which core turns into the delay between
	// the cast landing and the damage arriving. Zero for a spell that hits the instant it is cast.
	MissileSpeed float64

	// SpellAuraOptions.ProcChance, as a percentage. 100 means the aura fires on its own condition
	// rather than on a roll, which is how Flurry and Enrage read.
	ProcChance int32

	// SpellMisc.SchoolMask, in the client's bit order, which is not the sim's - see schoolName.
	SchoolMask int32

	// SpellCategories.DefenseType: 1 magic, 2 melee, 3 ranged, 0 none. Same order core uses.
	DefenseType int32

	Effects []RankEffect
}

// Absent from a database extracted before SpellCastTimes went into generator-settings.json. Keyed by
// handle, not once per process: gen_db opens two databases and one must not answer for the other.
var castTimesByDB sync.Map

func castTimesAvailable(db *sql.DB) bool {
	if cached, ok := castTimesByDB.Load(db); ok {
		return cached.(bool)
	}

	var n int
	err := db.QueryRow(
		`SELECT count(*) FROM sqlite_master WHERE type = 'table' AND name = 'SpellCastTimes'`).Scan(&n)
	available := err == nil && n > 0
	castTimesByDB.Store(db, available)
	return available
}

func RequireSpellCastTimes(db *sql.DB) error {
	if castTimesAvailable(db) {
		return nil
	}
	return errors.New("the client database has no SpellCastTimes table, so every generated cast time would " +
		"be zero - add \"SpellCastTimes\" to tools/database/generator-settings.json and re-run `make db`")
}

// The calibrated rule, shared with the regeneration check. Scored 86/87 against the hand tables.
//
// float32 is load-bearing: EffectRealPointsPerLevel is a float32 widened into the DB
// (3.79999995231628), and multiplying in float64 breaks 6 rows. Reproduces the tooltip, which is not
// the same as reproducing the server roll.
func DeriveRankAmount(e RankEffect, spellLevel, maxLevel int32) (min float64, max float64) {
	cap := maxLevel
	if cap <= 0 {
		cap = RankLevel
	}
	lvl := int32(RankLevel)
	if cap < lvl {
		lvl = cap
	}
	delta := lvl - spellLevel
	if delta < 0 {
		delta = 0
	}

	base := float32(e.BasePoints) + float32(float32(delta)*float32(e.PointsPerLvl))
	// The old EffectBasePoints column stored the value MINUS ONE, so the roll was
	// basePoints+1 .. basePoints+dieSides and the +1 below was the true minimum. This
	// client stores the real value in EffectBasePointsF, so the +1 is gone: Improved
	// Battle Shout reads 5/10/15/20/25 and must generate as 5/10/15/20/25, not 6..26.
	min = math.Floor(float64(base))
	max = math.Ceil(float64(base)) + float64(e.DieSides)
	if e.DieSides <= 0 {
		max = min
	}
	return min, max
}

func LoadRankSpell(db *sql.DB, spellID int32) (RankSpell, error) {
	s := RankSpell{SpellID: spellID}
	err := db.QueryRow(`
		SELECT l.SpellLevel, l.MaxLevel,
		       (SELECT ManaCost FROM SpellPower WHERE SpellID = l.SpellID ORDER BY OrderIndex LIMIT 1),
		       COALESCE((SELECT PowerType FROM SpellPower WHERE SpellID = l.SpellID ORDER BY OrderIndex LIMIT 1), 0)
		FROM SpellLevels l WHERE l.SpellID = ?`, spellID).Scan(&s.SpellLevel, &s.MaxLevel, &s.ManaCost, &s.PowerType)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		// Passive talents such as the warrior's Blood Craze have no SpellLevels row at all. No level
		// data means no level scaling, which is what setting the spell's own level to the cap gives.
		s.SpellLevel = RankLevel
	case err != nil:
		return s, fmt.Errorf("levels for spell %d: %w", spellID, err)
	}

	// Duration and its period are what a DoT's NumberOfTicks and TickLength are derived from.
	if err := scanOptional(db, `
		SELECT COALESCE(d.Duration, 0)
		FROM SpellMisc m LEFT JOIN SpellDuration d ON d.ID = m.DurationIndex
		WHERE m.SpellID = ?`, spellID, &s.DurationMs); err != nil {
		return s, fmt.Errorf("duration for spell %d: %w", spellID, err)
	}

	// Cooldown takes the longer of the two: Fire Blast and Cone of Cold use the shared
	// CategoryRecoveryTime, everything else its own RecoveryTime.
	if err := scanOptional(db, `
		SELECT COALESCE(max(RecoveryTime, CategoryRecoveryTime), 0), COALESCE(StartRecoveryTime, 0)
		FROM SpellCooldowns WHERE SpellID = ?`, spellID, &s.CooldownMs, &s.GCDMs); err != nil {
		return s, fmt.Errorf("cooldown for spell %d: %w", spellID, err)
	}

	// RangeMin is nonzero on only 212 spells in this build - the dead zone on a charge, and a handful
	// of ranged abilities - but where it exists core gates the cast on it exactly as it does MaxRange.
	if err := scanOptional(db, `
		SELECT COALESCE(r.RangeMin_0, 0), COALESCE(r.RangeMax_0, 0)
		FROM SpellMisc m JOIN SpellRange r ON r.ID = m.RangeIndex
		WHERE m.SpellID = ?`, spellID, &s.MinRange, &s.MaxRange); err != nil {
		return s, fmt.Errorf("range for spell %d: %w", spellID, err)
	}

	if err := scanOptional(db,
		`SELECT COALESCE(Speed, 0) FROM SpellMisc WHERE SpellID = ?`, spellID, &s.MissileSpeed); err != nil {
		return s, fmt.Errorf("missile speed for spell %d: %w", spellID, err)
	}

	if err := scanOptional(db,
		`SELECT COALESCE(ProcChance, 0) FROM SpellAuraOptions WHERE SpellID = ?`, spellID, &s.ProcChance); err != nil {
		return s, fmt.Errorf("proc chance for spell %d: %w", spellID, err)
	}

	if err := scanOptional(db,
		`SELECT COALESCE(SchoolMask, 0) FROM SpellMisc WHERE SpellID = ?`, spellID, &s.SchoolMask); err != nil {
		return s, fmt.Errorf("school for spell %d: %w", spellID, err)
	}

	if err := scanOptional(db,
		`SELECT COALESCE(DefenseType, 0) FROM SpellCategories WHERE SpellID = ?`, spellID, &s.DefenseType); err != nil {
		return s, fmt.Errorf("defense type for spell %d: %w", spellID, err)
	}

	if castTimesAvailable(db) {
		if err := scanOptional(db, `
			SELECT COALESCE(ct.Base, 0)
			FROM SpellMisc m JOIN SpellCastTimes ct ON ct.ID = m.CastingTimeIndex
			WHERE m.SpellID = ?`, spellID, &s.CastTimeMs); err != nil {
			return s, fmt.Errorf("cast time for spell %d: %w", spellID, err)
		}
	}

	s.Effects, err = RankEffectsOf(db, spellID)
	return s, err
}

// A spell with no row in one of the optional tables is ordinary - most spells have no cooldown - and
// leaves the destination at its zero, which is what core reads as "ungated". Any other error means the
// schema moved, and silently zeroing a cast time or a range on that is the failure this loader exists
// to avoid.
func scanOptional(db *sql.DB, query string, spellID int32, dest ...any) error {
	err := db.QueryRow(query, spellID).Scan(dest...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	return err
}

func RankEffectsOf(db *sql.DB, spellID int32) ([]RankEffect, error) {
	rows, err := db.Query(`
		-- EffectBasePoints became the REAL EffectBasePointsF in this client's layout, and
		-- EffectDieSides is gone entirely (Variance carries the spread now). BasePoints stays
		-- integral here because the dummy-target heuristic below reads base+dieSides as a
		-- spell id.
		--
		-- TODO: ~1.8% of SpellEffect rows have a fractional EffectBasePointsF and lose it to
		-- this cast. TODO: with DieSides pinned to 0, DeriveRankAmount collapses min and max
		-- onto the same value, so generated rank tables no longer carry a damage range.
		SELECT EffectIndex, Effect, EffectAura, CAST(EffectBasePointsF AS INTEGER), 0,
		       EffectRealPointsPerLevel, EffectBonusCoefficient, BonusCoefficientFromAP, EffectAuraPeriod,
		       COALESCE(EffectMiscValue_0, 0)
		FROM SpellEffect WHERE SpellID = ? ORDER BY EffectIndex`, spellID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []RankEffect
	for rows.Next() {
		e := RankEffect{OwnerSpellID: spellID}
		if err := rows.Scan(&e.Index, &e.Effect, &e.Aura, &e.BasePoints, &e.DieSides, &e.PointsPerLvl, &e.Coefficient, &e.APCoef, &e.AuraPeriod, &e.MiscValue); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

// Judgement of Command's rows sit on the Retribution skill line with ClassMask 0, so the sibling
// search by name finds nothing. The dummy names its target itself: base points plus the one die side
// derive to the damage spell's ID - 20425 states 20466+1 = 20467 - and that spell shares the name and
// the rank subtext, which is what is checked before its effects are taken. Following the pointer
// rather than widening the name search keeps Blizzard's tick spell out of its parent's Direct.
func dummyTargetEffects(db *sql.DB, spell RankSpell) ([]RankEffect, error) {
	if len(spell.Effects) != 1 || spell.Effects[0].Effect != dbc.E_DUMMY {
		return nil, nil
	}
	target, _ := DeriveRankAmount(spell.Effects[0], spell.SpellLevel, spell.MaxLevel)
	id := int32(target)
	if float64(id) != target || id <= 0 {
		return nil, nil
	}

	var same int
	if err := db.QueryRow(`
		SELECT count(*)
		FROM Spell s
		JOIN SpellName n ON n.ID = s.ID
		WHERE s.ID = ?
		  AND n.Name_lang = (SELECT Name_lang FROM SpellName WHERE ID = ?)
		  AND s.NameSubtext_lang = (SELECT NameSubtext_lang FROM Spell WHERE ID = ?)`,
		id, spell.SpellID, spell.SpellID).Scan(&same); err != nil {
		return nil, err
	}
	if same == 0 {
		return nil, nil
	}
	return RankEffectsOf(db, id)
}

// Holy Shock's registered spells carry only Effect=3 (dummy) and have no EffectTriggerSpell edge to the
// damage and heal spells that share their name and rank - the association exists nowhere but the name.
// Restricted to the family's class so an NPC copy of the name cannot be picked up.
func SiblingRankEffects(db *sql.DB, spellID int32, classBit int) ([]RankEffect, error) {
	rows, err := db.Query(`
		SELECT DISTINCT sla.Spell
		FROM SkillLineAbility sla
		JOIN SpellName n ON n.ID = sla.Spell
		JOIN Spell s ON s.ID = sla.Spell
		WHERE n.Name_lang = (SELECT Name_lang FROM SpellName WHERE ID = ?)
		  AND s.NameSubtext_lang = (SELECT NameSubtext_lang FROM Spell WHERE ID = ?)
		  AND (sla.ClassMask & ?) != 0
		  AND sla.Spell != ?
		ORDER BY sla.Spell`, spellID, spellID, classBit, spellID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int32
	for rows.Next() {
		var id int32
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	var out []RankEffect
	for _, id := range ids {
		effs, err := RankEffectsOf(db, id)
		if err != nil {
			return nil, err
		}
		out = append(out, effs...)
	}
	return out, nil
}

func HasValueEffect(effects []RankEffect) bool {
	for _, e := range effects {
		if e.Effect == dbc.E_SCHOOL_DAMAGE || e.Effect == dbc.E_HEAL || e.Effect == dbc.E_ENERGIZE ||
			IsWeaponDamageEffect(e.Effect) || IsPeriodicAura(e.Aura) {
			return true
		}
	}
	return false
}

// Every effect that could carry the numbers a rank row wants, including the ones reached only through
// a same-name sibling.
func RankCandidates(db *sql.DB, spellID int32, classBit int) (RankSpell, []RankEffect, error) {
	spell, err := LoadRankSpell(db, spellID)
	if err != nil {
		return spell, nil, err
	}

	candidates := spell.Effects
	if !HasValueEffect(candidates) {
		sibs, err := SiblingRankEffects(db, spellID, classBit)
		if err != nil {
			return spell, nil, err
		}
		if len(sibs) == 0 {
			sibs, err = dummyTargetEffects(db, spell)
			if err != nil {
				return spell, nil, err
			}
		}
		candidates = slices.Concat(spell.Effects, sibs)
	}
	return spell, candidates, nil
}
