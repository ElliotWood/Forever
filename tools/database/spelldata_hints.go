package database

import (
	"regexp"
	"strings"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/tools/database/dbc"
)

// What the tooltip says about a proc that its columns cannot, baked onto the row at generation.
// Spell.Description_lang is generator input only - the store carries the reading, never the text -
// so the sim reads a decided shape instead of re-parsing English at runtime.
//
// The four shapes are the ones docs/spell_data.md names under "Proc chances": the column is the
// roll, an effect's value is the roll, the column is a sentinel for "fires on its own condition",
// or nothing states a rate at all and an override owes one.

// The tooltip's own name for SpellAuraOptions.ProcChance. A chance stated as "$<id>h" is another
// spell's column and cannot match, since the id sits between the $ and the h. The chance is
// sometimes scaled - "${$h/2}% chance" - which still says the column is where it comes from.
var tooltipOwnChance = regexp.MustCompile(`\$h`)

// A chance stated as an effect's value: "$m2% chance", "$s1%". Only the first three effects can be
// named this way, and the digit is the client's EffectIndex plus one.
var tooltipEffectChance = regexp.MustCompile(`\$[ms]([123])%`)

// Reads the proc shape off the tooltip and the aura columns and writes it onto the row.
func applyTooltipHints(t *spellTables, s *storeSpell) {
	description := t.descriptions[s.ID]

	s.ProcHint = procTooltipHints(description)
	s.ProcChanceSource, s.ProcChanceEffect = procChanceSource(description, s)
}

func procChanceSource(description string, s *storeSpell) (storeProcChanceSource, int8) {
	if tooltipOwnChance.MatchString(description) {
		return procChanceColumn, 0
	}

	// A tooltip names several effects - Shield Specialization states "$s1% block" and "$m2% chance
	// to generate rage" - so the one the chance is on is the one carrying the proc aura, not the
	// first one the text mentions.
	for _, m := range tooltipEffectChance.FindAllStringSubmatch(description, -1) {
		index := int8(m[1][0] - '0')
		if isProcEffect(s.effect(uint8(index - 1))) {
			return procChanceEffectN, index
		}
	}

	// 100 and 101 are both the client's "no roll here": the aura fires whenever its own condition
	// is met, and on a spell that is not a proc at all they mean nothing.
	if s.ProcChance == 100 || s.ProcChance == 101 {
		return procChanceAlways, 0
	}

	// A proc aura with no chance anywhere is the shape whose rate the client does not carry, so it
	// has to come from an override into RPPM.
	if s.ProcChance == 0 && s.triggersAProc() {
		return procChancePPM, 0
	}

	return procChanceColumn, 0
}

// An effect a tooltip can state a chance for: the two trigger auras, and the dummy the server hangs
// a hand-written proc off - Enrage's chance sits on one.
func isProcEffect(e *storeEffect) bool {
	if e == nil {
		return false
	}
	aura := dbc.EffectAuraType(e.Aura)
	return aura == dbc.A_PROC_TRIGGER_SPELL || aura == dbc.A_PROC_TRIGGER_SPELL_WITH_VALUE ||
		aura == dbc.A_DUMMY
}

// Whether the spell fires something through the client's own proc machinery, which is what makes a
// missing rate a rate somebody owes. The dummy is deliberately out: Rip, Rupture and Arcane
// Missiles all carry one, and none of them is a proc.
func (s *storeSpell) triggersAProc() bool {
	for i := range s.Effects {
		aura := dbc.EffectAuraType(s.Effects[i].Aura)
		if aura == dbc.A_PROC_TRIGGER_SPELL || aura == dbc.A_PROC_TRIGGER_SPELL_WITH_VALUE {
			return true
		}
	}
	return false
}

// The effect the client files under an index, which has gaps: EffectIndex is the tooltip's
// numbering, while the slice is packed.
func (s *storeSpell) effect(index uint8) *storeEffect {
	for i := range s.Effects {
		if s.Effects[i].Index == index {
			return &s.Effects[i]
		}
	}
	return nil
}

// The store's ProcChanceSource, mirrored here the way the row structs are. The order is the store's
// own: the emitted name comes from this table, not from the number.
type storeProcChanceSource uint8

const (
	procChanceColumn storeProcChanceSource = iota
	procChanceEffectN
	procChanceAlways
	procChancePPM
)

func (s storeProcChanceSource) String() string {
	switch s {
	case procChanceEffectN:
		return "ProcChanceEffectN"
	case procChanceAlways:
		return "ProcChanceAlways"
	case procChancePPM:
		return "ProcChancePPM"
	default:
		return "ProcChanceColumn"
	}
}

// The hint bits as the generated file names them, so a row states which words the reading came
// from rather than a number.
func formatProcHint(hint core.ProcHint) string {
	names := []struct {
		bit  core.ProcHint
		name string
	}{
		{core.ProcHintCastTrigger, "core.ProcHintCastTrigger"},
		{core.ProcHintCrit, "core.ProcHintCrit"},
		{core.ProcHintHeals, "core.ProcHintHeals"},
		{core.ProcHintPureHeal, "core.ProcHintPureHeal"},
		{core.ProcHintNamedAbility, "core.ProcHintNamedAbility"},
		{core.ProcHintOutcomeTaken, "core.ProcHintOutcomeTaken"},
	}

	var set []string
	for _, n := range names {
		if hint.Matches(n.bit) {
			set = append(set, n.name)
		}
	}
	return strings.Join(set, " | ")
}
