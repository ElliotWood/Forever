package buffs

import (
	"math"
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/spelldata"
)

// The readers the generated constructors turn a store row into a buff with. An amount is the
// caster's at core.CharacterLevel, in the client's units until the constructor converts it.

func amount(e *spelldata.Effect) float64 {
	return e.Average(core.CharacterLevel)
}

func manaPerFive(tick *spelldata.Effect, amount float64) float64 {
	return amount * 5000 / float64(tick.PeriodMs)
}

// The raid config puts a finisher on the target at full combo points; the client states the amount
// per point spent.
const maxComboPoints = 5

func fullComboPoints(e *spelldata.Effect) float64 {
	return float64(e.PointsPerResource) * maxComboPoints
}

// The client states a permanent aura as -1 and an aura with no duration of its own, a totem's, as 0.
func auraDuration(s *spelldata.Spell) time.Duration {
	if s.DurationMs <= 0 {
		return core.NeverExpires
	}
	return s.Duration()
}

func cooldown(s *spelldata.Spell) time.Duration {
	return max(s.Cooldown(), s.CategoryCooldown())
}

// An improving talent's modifier on the buff's amount, truncated the way the client resolves one: a
// percentage of it, or a flat addition. An untaken talent is NilEffect, which adds nothing.
func talentScaled(amount float64, mod *spelldata.Effect) float64 {
	if mod.Aura == dbcenums.A_ADD_PCT_MODIFIER {
		return math.Trunc(amount * (1 + mod.BaseValue()/100))
	}
	return math.Trunc(amount + mod.BaseValue())
}

func talentScaledDuration(s *spelldata.Spell, mod *spelldata.Effect) time.Duration {
	if s.DurationMs <= 0 {
		return core.NeverExpires
	}
	return time.Duration(talentScaled(float64(s.DurationMs), mod)) * time.Millisecond
}
