// foreverdesc resolves the client's tooltip format tokens ($s1, $m1, $o1, $d, ...)
// in every spell description in tools/database/wowsims-forever.db, using the repo's
// own tools/tooltip parser.
//
// The parser needs a TooltipDataProvider. The one that ships with it reads the TBC
// JSON dumps under assets/db_inputs/dbc, which do not exist for this database, so
// this provider is backed by wowsims-forever.db directly.
//
//	go run ./tools/foreverdesc -db tools/database/wowsims-forever.db -o descriptions.json
package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/tools/database/dbc"
	"github.com/wowsims/forever/tools/tooltip"
	_ "modernc.org/sqlite"
)

// The client is a level-60 one. Character-dependent terms resolve against a bare
// character (no attack power, no spell power) so what comes out is the spell's own
// base value, the way a spell database shows it rather than a geared character sheet.
const (
	playerLevel = 60.0
	attackPower = 0.0
	spellDamage = 0.0
)

type effect struct {
	idx                                     int64
	effectType, aura                        int64
	basePoints, realPointsPerLevel          float64
	coefficient, bonusCoefficient, apCoeff  float64
	amplitude, chainAmplitude, pointsPerRes float64
	auraPeriod, chainTargets                int64
	radiusMax                               float64
	variance                                float64
}

type spell struct {
	name, desc, variables string
	classSet              int64
	duration              int64 // ms
	maxRange              float64
	procChance            float64
	procCharges, stacks   int64
	procCategoryRecovery  int64
	maxTargets            int64
	effects               []effect
}

// rankCurves[spellID][effectIndex] = that effect's value at each rank, from
// TraitDefinition -> TraitDefinitionEffectPoints -> CurvePoint.
type provider struct {
	spells     map[int64]*spell
	rankCurves map[int64]map[int64][]float64
	// When ovrSpell is set, base points for that spell come from ovrPoints
	// instead of SpellEffect, so the parser renders the description at one rank.
	ovrSpell  int64
	ovrPoints map[int64]float64
}

func (p *provider) override(spellID int64, rank int) {
	p.ovrSpell, p.ovrPoints = 0, nil
	curves, ok := p.rankCurves[spellID]
	if !ok {
		return
	}
	pts := map[int64]float64{}
	for ei, vals := range curves {
		if rank < len(vals) {
			pts[ei] = vals[rank]
		}
	}
	if len(pts) > 0 {
		p.ovrSpell, p.ovrPoints = spellID, pts
	}
}

func (p *provider) basePoints(id, idx int64, fallback float64) float64 {
	if id == p.ovrSpell && p.ovrPoints != nil {
		if v, ok := p.ovrPoints[idx]; ok {
			return v
		}
	}
	return fallback
}

func (p *provider) get(id int64) *spell { return p.spells[id] }

func (p *provider) eff(id, idx int64) *effect {
	s := p.get(id)
	if s == nil {
		return nil
	}
	for i := range s.effects {
		if s.effects[i].idx == idx {
			return &s.effects[i]
		}
	}
	return nil
}

func (p *provider) GetAttackPower() float64 { return attackPower }
func (p *provider) GetSpellDamage() float64 { return spellDamage }
func (p *provider) GetPlayerLevel() float64 { return playerLevel }
func (p *provider) GetSpecNum() int64       { return 0 }
func (p *provider) IsMaleGender() bool      { return true }
func (p *provider) HasAura(int64) bool      { return true }
func (p *provider) HasPassive(int64) bool   { return true }
func (p *provider) KnowsSpell(int64) bool   { return true }

func (p *provider) GetMainHandWeapon() *core.Weapon {
	return &core.Weapon{BaseDamageMin: 100, BaseDamageMax: 200, SwingSpeed: 2.6}
}
func (p *provider) GetOffHandWeapon() *core.Weapon {
	return &core.Weapon{BaseDamageMin: 100, BaseDamageMax: 200, SwingSpeed: 2.6}
}

func (p *provider) GetSpellName(id int64) string {
	if s := p.get(id); s != nil {
		return s.name
	}
	return ""
}
func (p *provider) GetSpellDescription(id int64) string {
	if s := p.get(id); s != nil {
		return s.desc
	}
	return ""
}
func (p *provider) GetDescriptionVariableString(id int64) string {
	if s := p.get(id); s != nil {
		return s.variables
	}
	return ""
}
func (p *provider) GetSpellIcon(int64) string { return "" }

func (p *provider) GetSpellDuration(id int64) time.Duration {
	s := p.get(id)
	if s == nil || s.duration < 0 {
		return 0
	}
	return time.Duration(s.duration) * time.Millisecond
}
func (p *provider) GetSpellRange(id int64) float64 {
	if s := p.get(id); s != nil {
		return s.maxRange
	}
	return 0
}
func (p *provider) GetSpellMaxTargets(id int64) int64 {
	if s := p.get(id); s != nil && s.maxTargets > 0 {
		return s.maxTargets
	}
	return 1
}
func (p *provider) GetSpellStacks(id int64) int64 {
	s := p.get(id)
	if s == nil {
		return 0
	}
	if s.procCharges > 0 {
		return s.procCharges
	}
	return s.stacks
}
func (p *provider) GetSpellProcChance(id int64) float64 {
	if s := p.get(id); s != nil {
		return s.procChance
	}
	return 0
}
func (p *provider) GetSpellProcCooldown(id int64) time.Duration {
	if s := p.get(id); s != nil {
		return time.Duration(s.procCategoryRecovery) * time.Millisecond
	}
	return 1
}
func (p *provider) GetSpellPPM(int64) float64 { return 0 }

func (p *provider) GetEffectAmplitude(id, idx int64) float64 {
	if e := p.eff(id, idx); e != nil {
		return e.amplitude
	}
	return 0
}
func (p *provider) GetEffectChainAmplitude(id, idx int64) float64 {
	if e := p.eff(id, idx); e != nil {
		return e.chainAmplitude
	}
	return 0
}
func (p *provider) GetEffectPointsPerResource(id, idx int64) float64 {
	if e := p.eff(id, idx); e != nil {
		return e.pointsPerRes
	}
	return 0
}
func (p *provider) GetEffectMaxTargets(id, idx int64) int64 {
	if e := p.eff(id, idx); e != nil {
		return e.chainTargets
	}
	return 0
}
func (p *provider) GetEffectPeriod(id, idx int64) time.Duration {
	if e := p.eff(id, idx); e != nil {
		return time.Duration(e.auraPeriod) * time.Millisecond
	}
	return 0
}
func (p *provider) GetEffectRadius(id, idx int64) float64 {
	if e := p.eff(id, idx); e != nil {
		return e.radiusMax
	}
	return 0
}

// A damage or healing range is stored as a base value plus a fractional Variance
// rather than the old die-sides pair, so $m and $M are the two ends of that
// spread: Frostbolt rank 1 is base 19 variance 0.105, which is 18 to 20.
func (p *provider) spread(id, idx int64, high bool) float64 {
	e := p.eff(id, idx)
	if e == nil {
		return 0
	}
	base := p.basePoints(id, idx, e.basePoints)
	if e.variance == 0 {
		return base
	}
	half := e.variance / 2
	if high {
		return math.Round(base * (1 + half))
	}
	return math.Round(base * (1 - half))
}

// GetEffectBaseValue backs $m, the bottom of the range.
func (p *provider) GetEffectBaseValue(id, idx int64) float64 {
	return p.spread(id, idx, false)
}

// GetEffectMaxValue backs $M, the top of the range.
func (p *provider) GetEffectMaxValue(id, idx int64) float64 {
	return p.spread(id, idx, true)
}

func classOf(classSet int64) proto.Class {
	switch classSet {
	case 3:
		return proto.Class_ClassMage
	case 4:
		return proto.Class_ClassWarrior
	case 5:
		return proto.Class_ClassWarlock
	case 6:
		return proto.Class_ClassPriest
	case 7:
		return proto.Class_ClassDruid
	case 8:
		return proto.Class_ClassRogue
	case 9:
		return proto.Class_ClassHunter
	case 10:
		return proto.Class_ClassPaladin
	case 11:
		return proto.Class_ClassShaman
	}
	return proto.Class_ClassUnknown
}

// GetEffectScaledValue follows DBCTooltipDataProvider.GetEffectScaledValue.
func (p *provider) GetEffectScaledValue(id, idx int64) float64 {
	s := p.get(id)
	e := p.eff(id, idx)
	if s == nil || e == nil {
		return 1
	}
	base := 0.0
	if e.coefficient > 0 && s.classSet > 0 {
		cls := classOf(s.classSet)
		scale := 1710.0
		if cls != proto.Class_ClassUnknown {
			scale = core.ClassBaseScaling[cls]
		}
		base += scale * e.coefficient
	} else {
		base += p.basePoints(id, idx, e.basePoints)
		base += e.realPointsPerLevel * playerLevel
	}

	shouldScale := false
	switch dbc.SpellEffectType(e.effectType) {
	case dbc.E_SCHOOL_DAMAGE:
		shouldScale = true
	case dbc.E_APPLY_AURA, dbc.E_APPLY_AREA_AURA_ENEMY, dbc.E_APPLY_AREA_AURA_FRIEND,
		dbc.E_APPLY_AREA_AURA_PARTY, dbc.E_APPLY_AREA_AURA_OWNER, dbc.E_APPLY_AREA_AURA_RAID,
		dbc.E_APPLY_AREA_AURA_PARTY_NONRANDOM, dbc.E_APPLY_AREA_AURA_PET, dbc.E_APPLY_AURA_ON_PET:
		switch dbc.EffectAuraType(e.aura) {
		case dbc.A_PERIODIC_DAMAGE, dbc.A_PERIODIC_HEAL:
			shouldScale = true
		}
	}
	if !shouldScale {
		return base
	}
	if e.apCoeff > 0 {
		base += p.GetAttackPower() * e.apCoeff
	}
	if e.bonusCoefficient > 0 {
		base += p.GetSpellDamage() * e.bonusCoefficient
	}
	return base
}

func (p *provider) GetEffectEnchantValue(int64, int64) float64 { return 0 }

func load(db *sql.DB) (map[int64]*spell, error) {
	spells := map[int64]*spell{}
	at := func(id int64) *spell {
		s, ok := spells[id]
		if !ok {
			s = &spell{}
			spells[id] = s
		}
		return s
	}

	rows, err := db.Query(`select s.ID, coalesce(n.Name_lang,''), coalesce(s.Description_lang,'') from Spell s
	                       left join SpellName n on n.ID = s.ID`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var id int64
		var name, desc string
		if err := rows.Scan(&id, &name, &desc); err != nil {
			return nil, err
		}
		s := at(id)
		s.name, s.desc = name, desc
	}
	rows.Close()

	type q struct {
		sql   string
		apply func(*spell, []float64)
	}
	scalars := []q{
		{`select m.SpellID, coalesce(m.SchoolMask,0), coalesce(d.Duration,0), coalesce(r.RangeMax_1,0), coalesce(c.SpellClassSet,0)
		  from SpellMisc m
		  left join SpellDuration d on d.ID = m.DurationIndex
		  left join SpellRange r on r.ID = m.RangeIndex
		  left join SpellClassOptions c on c.SpellID = m.SpellID
		  where m.SpellID is not null`,
			func(s *spell, v []float64) {
				s.duration = int64(v[1])
				s.maxRange = v[2]
				s.classSet = int64(v[3])
			}},
		{`select SpellID, coalesce(ProcChance,0), coalesce(ProcCharges,0), coalesce(CumulativeAura,0), coalesce(ProcCategoryRecovery,0)
		  from SpellAuraOptions where SpellID is not null`,
			func(s *spell, v []float64) {
				s.procChance = v[0]
				s.procCharges = int64(v[1])
				s.stacks = int64(v[2])
				s.procCategoryRecovery = int64(v[3])
			}},
		{`select SpellID, coalesce(MaxTargets,0), 0, 0, 0 from SpellTargetRestrictions where SpellID is not null`,
			func(s *spell, v []float64) { s.maxTargets = int64(v[0]) }},
		{`select x.SpellID, 0, 0, 0, 0 from SpellXDescriptionVariables x where x.SpellID is not null`, nil},
	}
	for _, item := range scalars {
		if item.apply == nil {
			continue
		}
		rs, err := db.Query(item.sql)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", item.sql[:40], err)
		}
		for rs.Next() {
			var id int64
			v := make([]float64, 4)
			if err := rs.Scan(&id, &v[0], &v[1], &v[2], &v[3]); err != nil {
				rs.Close()
				return nil, err
			}
			item.apply(at(id), v)
		}
		rs.Close()
	}

	rs, err := db.Query(`select x.SpellID, coalesce(v.Variables,'') from SpellXDescriptionVariables x
	                     join SpellDescriptionVariables v on v.ID = x.SpellDescriptionVariablesID
	                     where x.SpellID is not null`)
	if err != nil {
		return nil, err
	}
	for rs.Next() {
		var id int64
		var vars string
		if err := rs.Scan(&id, &vars); err != nil {
			rs.Close()
			return nil, err
		}
		at(id).variables = vars
	}
	rs.Close()

	rs, err = db.Query(`select e.SpellID, e.EffectIndex, coalesce(e.Effect,0), coalesce(e.EffectAura,0),
	                      coalesce(e.EffectBasePointsF,0), coalesce(e.EffectRealPointsPerLevel,0),
	                      coalesce(e.Coefficient,0), coalesce(e.EffectBonusCoefficient,0),
	                      coalesce(e.BonusCoefficientFromAP,0), coalesce(e.EffectAmplitude,0),
	                      coalesce(e.EffectChainAmplitude,0), coalesce(e.EffectPointsPerResource,0),
	                      coalesce(e.EffectAuraPeriod,0), coalesce(e.EffectChainTargets,0),
	                      coalesce(r.RadiusMax,0), coalesce(e.Variance,0)
	                    from SpellEffect e
	                    left join SpellRadius r on r.ID = e.EffectRadiusIndex_0
	                    where e.SpellID is not null order by e.SpellID, e.EffectIndex`)
	if err != nil {
		return nil, err
	}
	for rs.Next() {
		var id int64
		var e effect
		if err := rs.Scan(&id, &e.idx, &e.effectType, &e.aura, &e.basePoints, &e.realPointsPerLevel,
			&e.coefficient, &e.bonusCoefficient, &e.apCoeff, &e.amplitude, &e.chainAmplitude,
			&e.pointsPerRes, &e.auraPeriod, &e.chainTargets, &e.radiusMax, &e.variance); err != nil {
			rs.Close()
			return nil, err
		}
		s := at(id)
		s.effects = append(s.effects, e)
	}
	rs.Close()
	return spells, nil
}

// loadCurves reads the modern trait system's per-rank effect values.
func loadCurves(db *sql.DB) (map[int64]map[int64][]float64, error) {
	out := map[int64]map[int64][]float64{}
	rows, err := db.Query(`select d.SpellID, ep.EffectIndex, ep.CurveID
	                       from TraitDefinition d
	                       join TraitDefinitionEffectPoints ep on ep.TraitDefinitionID = d.ID
	                       where d.SpellID is not null and ep.CurveID is not null`)
	if err != nil {
		return nil, err
	}
	type key struct{ spell, idx int64 }
	seen := map[key]bool{}
	var want []struct {
		spell, idx, curve int64
	}
	for rows.Next() {
		var sp, idx, cv int64
		if err := rows.Scan(&sp, &idx, &cv); err != nil {
			rows.Close()
			return nil, err
		}
		if seen[key{sp, idx}] {
			continue
		}
		seen[key{sp, idx}] = true
		want = append(want, struct{ spell, idx, curve int64 }{sp, idx, cv})
	}
	rows.Close()

	for _, w := range want {
		pr, err := db.Query(`select Pos_1 from CurvePoint where CurveID = ? order by OrderIndex`, w.curve)
		if err != nil {
			return nil, err
		}
		var vals []float64
		for pr.Next() {
			var v float64
			if err := pr.Scan(&v); err != nil {
				pr.Close()
				return nil, err
			}
			vals = append(vals, v)
		}
		pr.Close()
		if len(vals) == 0 {
			continue
		}
		if out[w.spell] == nil {
			out[w.spell] = map[int64][]float64{}
		}
		out[w.spell][w.idx] = vals
	}
	return out, nil
}

func main() {
	dbPath := flag.String("db", "tools/database/wowsims-forever.db", "sqlite database")
	out := flag.String("o", "descriptions.json", "output json")
	flag.Parse()

	db, err := sql.Open("sqlite", *dbPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer db.Close()

	spells, err := load(db)
	if err != nil {
		fmt.Fprintln(os.Stderr, "loading:", err)
		os.Exit(1)
	}
	curves, err := loadCurves(db)
	if err != nil {
		fmt.Fprintln(os.Stderr, "loading curves:", err)
		os.Exit(1)
	}
	p := &provider{spells: spells, rankCurves: curves}

	resolved := map[string]string{}
	var ok, failed int
	for id, s := range spells {
		if s.desc == "" {
			continue
		}
		parsed, err := tooltip.ParseTooltip(s.desc, p, id)
		if err != nil || parsed == nil {
			failed++
			continue
		}
		text := parsed.String()
		if text != "" && text != s.desc {
			resolved[fmt.Sprint(id)] = text
			ok++
		}
	}
	buf, _ := json.Marshal(resolved)
	if err := os.WriteFile(*out, buf, 0o644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Per-rank descriptions for curve-backed spells, so a talent reads correctly
	// at whatever rank the reader has selected.
	ranked := map[string][]string{}
	for id, curves := range p.rankCurves {
		s := p.get(id)
		if s == nil || s.desc == "" {
			continue
		}
		n := 0
		for _, vals := range curves {
			if len(vals) > n {
				n = len(vals)
			}
		}
		var texts []string
		for r := 0; r < n; r++ {
			p.override(id, r)
			parsed, err := tooltip.ParseTooltip(s.desc, p, id)
			if err != nil || parsed == nil {
				texts = append(texts, "")
				continue
			}
			texts = append(texts, parsed.String())
		}
		p.ovrSpell, p.ovrPoints = 0, nil
		if len(texts) > 0 {
			ranked[fmt.Sprint(id)] = texts
		}
	}
	rbuf, _ := json.Marshal(ranked)
	rpath := strings.TrimSuffix(*out, ".json") + "-ranked.json"
	if err := os.WriteFile(rpath, rbuf, 0o644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("resolved %d descriptions (%d failed) -> %s\n", ok, failed, *out)
	fmt.Printf("per-rank descriptions for %d curve-backed spells -> %s\n", len(ranked), rpath)
}
