package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/spelldata"
)

// One effect as both readings: what it does in words where a branch below knows the shape, and the
// client's own columns, which is what a developer checks the words against. Human is empty on an
// effect no branch covers, and the literal is then all there is to print.
type Line struct {
	Human   string `json:"human"`
	Literal string `json:"literal"`
}

func title(s *spelldata.Spell) string {
	if s.Rank == "" {
		return fmt.Sprintf("%d %s", s.ID, s.Name)
	}
	return fmt.Sprintf("%d %s (%s)", s.ID, s.Name, s.Rank)
}

// The spell's own columns, one label per line, leaving out what the row does not state.
func header(s *spelldata.Spell) []string {
	var out []string
	add := func(label, value string) {
		if value != "" {
			out = append(out, fmt.Sprintf("%-9s %s", label, value))
		}
	}

	add("school", schoolName(s.School))
	add("defense", defenseName(s.DefenseType))
	if dbcenums.Mechanic(s.Mechanic) == dbcenums.MECHANIC_BLEED {
		add("mechanic", "bleed")
	} else if s.Mechanic != 0 {
		add("mechanic", fmt.Sprintf("mechanic %d", s.Mechanic))
	}
	if s.DurationMs == -1 {
		add("duration", "permanent")
	} else {
		add("duration", seconds(s.DurationMs))
	}
	add("cast", seconds(s.CastTimeMs))
	add("cooldown", seconds(s.CooldownMs))
	if s.CategoryCooldownMs != 0 {
		add("cooldown", seconds(s.CategoryCooldownMs)+" (category)")
	}
	add("gcd", seconds(s.GCDMs))
	for i := range s.Powers {
		add("cost", cost(s, &s.Powers[i]))
	}
	return out
}

// SpellPower states either a flat amount out of the bar or a share of the caster's maximum.
func cost(s *spelldata.Spell, p *spelldata.Power) string {
	if p.CostPct != 0 {
		return fmt.Sprintf("%s%% of maximum %s", number(float64(p.CostPct)), powerName(p.Type))
	}
	if p.Cost == 0 {
		return ""
	}
	return join(number(s.PowerCost(p.Type)), powerName(p.Type))
}

func effectLines(s *spelldata.Spell) []Line {
	var out []Line
	for i := range s.Effects {
		e := &s.Effects[i]
		out = append(out, Line{Human: humanise(s, e), Literal: literal(e)})
	}
	return out
}

func wowheadURL(id int32) string {
	return fmt.Sprintf("https://www.wowhead.com/forever/spell=%d", id)
}

// What the effect does, for the shapes the client states plainly enough to word without reading the
// tooltip: a tick, a direct hit, a weapon multiplier, a power gain and a stat aura. Anything else
// answers empty and prints as its columns alone.
func humanise(s *spelldata.Spell, e *spelldata.Effect) string {
	switch e.Type {
	case dbcenums.E_SCHOOL_DAMAGE:
		return join(number(e.BasePoints), schoolName(s.School), "damage", targetPhrase(e))
	case dbcenums.E_HEAL:
		return join(number(e.BasePoints), "healing", targetPhrase(e))
	case dbcenums.E_WEAPON_PERCENT_DAMAGE:
		return join(number(e.BasePoints)+"% weapon damage", targetPhrase(e))
	case dbcenums.E_ENERGIZE:
		return join("restores", powerAmount(e), targetPhrase(e))
	case dbcenums.E_APPLY_AURA:
		return humaniseAura(s, e)
	}
	return ""
}

func humaniseAura(s *spelldata.Spell, e *spelldata.Effect) string {
	switch e.Aura {
	case dbcenums.A_PERIODIC_DAMAGE:
		return join(number(e.BasePoints), schoolName(s.School), "damage", every(e), targetPhrase(e), ticks(s, e))
	case dbcenums.A_PERIODIC_HEAL:
		return join(number(e.BasePoints), "healing", every(e), targetPhrase(e), ticks(s, e))
	case dbcenums.A_PERIODIC_ENERGIZE:
		return join("restores", powerAmount(e), every(e), targetPhrase(e), ticks(s, e))
	case dbcenums.A_MOD_STAT:
		return join(signed(e.BasePoints), statName(e.Misc))
	case dbcenums.A_MOD_TOTAL_STAT_PERCENTAGE:
		return join(signed(e.BasePoints)+"%", statName(e.Misc))
	}
	return ""
}

// The effect as the client states it: the columns it fills, in the units the row keeps them in, so the
// words above can be checked against them.
func literal(e *spelldata.Effect) string {
	parts := []string{effectTypeName(e.Type)}
	if e.Aura != 0 {
		parts = append(parts, auraName(e.Aura))
	}
	parts = append(parts, "base="+number(e.BasePoints))

	add := func(format string, args ...any) {
		parts = append(parts, fmt.Sprintf(format, args...))
	}
	if e.PPL != 0 {
		add("ppl=%s", number(e.PPL))
	}
	if e.Variance != 0 {
		add("variance=%s", number(e.Variance))
	}
	if e.SPCoef != 0 {
		add("sp=%s", number(e.SPCoef))
	}
	if e.APCoef != 0 {
		add("ap=%s", number(e.APCoef))
	}
	if e.PeriodMs != 0 {
		add("period=%dms", e.PeriodMs)
	}
	if e.Misc != 0 {
		add("misc=%d", e.Misc)
	}
	if e.Misc2 != 0 {
		add("misc2=%d", e.Misc2)
	}
	if e.TriggerID != 0 {
		add("trigger=%d", e.TriggerID)
	}
	if e.ChainTargets != 0 {
		add("chain=%d", e.ChainTargets)
	}
	if e.RadiusMax != 0 {
		add("radius=%s", number(float64(e.RadiusMax)))
	}
	if e.Mechanic != 0 {
		add("mechanic=%d", e.Mechanic)
	}
	add("target=[%d,%d]", e.Target[0], e.Target[1])
	return strings.Join(parts, " ")
}

func every(e *spelldata.Effect) string {
	if e.PeriodMs == 0 {
		return ""
	}
	return "every " + seconds(e.PeriodMs)
}

// How many times the aura ticks over its spell's duration, where both are stated.
func ticks(s *spelldata.Spell, e *spelldata.Effect) string {
	if s.DurationMs <= 0 || e.PeriodMs <= 0 {
		return ""
	}
	return fmt.Sprintf("(%d ticks)", s.DurationMs/e.PeriodMs)
}

// The two implicit targets the store's rows state often enough to word. The rest are left to the
// literal's target=[a,b].
func targetPhrase(e *spelldata.Effect) string {
	switch e.Target[0] {
	case 1:
		return "to the caster"
	case 6:
		return "to the enemy"
	}
	return ""
}

// EffectMiscValue_0 is the power type on an energize effect, and rage is on the client's 0-1000 bar
// the way a cost is.
func powerAmount(e *spelldata.Effect) string {
	amount := e.BasePoints
	if int8(e.Misc) == powerTypeRage {
		amount /= 10
	}
	bar := powerName(int8(e.Misc))
	if amount == 1 {
		bar = strings.TrimSuffix(bar, "s")
	}
	return join(number(amount), bar)
}

const powerTypeRage int8 = 1

// SpellPower.PowerType, of which the store's rows carry six: the client's health is -2.
func powerName(t int8) string {
	switch t {
	case -2:
		return "health"
	case 0:
		return "mana"
	case powerTypeRage:
		return "rage"
	case 2:
		return "focus"
	case 3:
		return "energy"
	case 4:
		return "combo points"
	}
	return fmt.Sprintf("power %d", t)
}

// A_MOD_STAT states the stat in its misc value, and -1 is the client's "every stat".
func statName(misc int32) string {
	switch misc {
	case -1:
		return "to all stats"
	case 0:
		return "strength"
	case 1:
		return "agility"
	case 2:
		return "stamina"
	case 3:
		return "intellect"
	case 4:
		return "spirit"
	}
	return fmt.Sprintf("stat %d", misc)
}

// The school mask as the schools it holds, since a row can carry more than one bit.
func schoolName(mask uint8) string {
	if mask == 0 {
		return ""
	}
	var names []string
	for _, s := range []struct {
		bit  core.SpellSchool
		name string
	}{
		{core.SpellSchoolPhysical, "physical"},
		{core.SpellSchoolHoly, "holy"},
		{core.SpellSchoolFire, "fire"},
		{core.SpellSchoolNature, "nature"},
		{core.SpellSchoolFrost, "frost"},
		{core.SpellSchoolShadow, "shadow"},
		{core.SpellSchoolArcane, "arcane"},
	} {
		if core.SpellSchool(mask)&s.bit != 0 {
			names = append(names, s.name)
		}
	}
	return strings.Join(names, "+")
}

func defenseName(t uint8) string {
	switch core.DefenseType(t) {
	case core.DefenseTypeMagic:
		return "magic"
	case core.DefenseTypeMelee:
		return "melee"
	case core.DefenseTypeRanged:
		return "ranged"
	}
	return ""
}

func seconds(ms int32) string {
	if ms == 0 {
		return ""
	}
	return number(float64(ms)/1000) + " s"
}

// Read back as the float32 the client's columns are, so a base point of 5.449999809265137 - the
// widening of the client's 5.45 - prints as the client states it.
func number(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 32)
}

func signed(v float64) string {
	if v >= 0 {
		return "+" + number(v)
	}
	return number(v)
}

func join(parts ...string) string {
	var kept []string
	for _, p := range parts {
		if p != "" {
			kept = append(kept, p)
		}
	}
	return strings.Join(kept, " ")
}
