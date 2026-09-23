package spelldata

import (
	"fmt"
	"math"

	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/stats"
)

// The rows a raid buff reads and nothing else does, which BuffAuras turns on. A class's own spell
// states these auras for a resource or a value the class wires itself: Bloodrage's rage over time is
// an A_PERIODIC_ENERGIZE, and Berserker Stance states an attack power percentage of 0.
var buffAuraTable = map[dbcenums.EffectAuraType]row{
	// Mana over time, which the sim keeps as mana per five seconds. The client's tick is whole; what
	// the conversion makes of it is not, and it is not truncated.
	dbcenums.A_PERIODIC_ENERGIZE: func(p *parser, e *Effect, v float64) *attachment {
		if e.Misc != int32(dbcenums.POWER_MANA) || e.PeriodMs <= 0 {
			return nil
		}
		return p.statBuff(stats.MP5, v*5000/float64(e.PeriodMs))
	},

	dbcenums.A_MOD_ATTACK_POWER_PCT: func(p *parser, e *Effect, v float64) *attachment {
		return p.statMultiplier([]stats.Stat{stats.AttackPower}, percentMultiplier(v))
	},
}

// What a buff built from the row with these options leaves out, one note per aura effect the parse
// skips. An empty answer means every aura effect the row states is attached. The notes come from
// DryRun, so a buff the generator writes is one the sim builds the same way.
func BuffUnsupported(s *Spell, onCharacter bool, opts ...ParseOpt) []string {
	if s == nil || s == Nil || s.ID == 0 {
		return []string{"no row in the store"}
	}

	var notes []string
	for _, e := range DryRun(s, onCharacter, opts...).Skipped {
		pos := 0
		for i := range s.Effects {
			if &s.Effects[i] == e {
				pos = i + 1
			}
		}

		if e.PointsPerResource != 0 && e.BasePoints == 0 {
			notes = append(notes, fmt.Sprintf("effect %d is worth %v per combo point", pos, e.PointsPerResource))
			continue
		}
		notes = append(notes, fmt.Sprintf("effect %d %s(%d) misc %d", pos, auraName(e.Aura), e.Aura, e.Misc))
	}
	return notes
}

// An improving talent's modifier on an amount, truncated the way the client resolves one: a
// percentage of it, or a flat addition. An untaken talent is NilEffect, which adds nothing.
func Scaled(amount float64, mod *Effect) float64 {
	if mod.Aura == dbcenums.A_ADD_PCT_MODIFIER {
		return math.Trunc(amount * (1 + mod.BaseValue()/100))
	}
	return math.Trunc(amount + mod.BaseValue())
}
