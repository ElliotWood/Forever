package spelldata

import (
	"math"
	"time"

	"github.com/wowsims/forever/sim/core"
)

// The client's own number, before any unit the aura gives it.
func (e *Effect) BaseValue() float64 {
	return e.BasePoints
}

// The client states a percentage as an integer: Improved Righteous Fury reads 16, not 0.16.
func (e *Effect) Percent() float64 {
	return e.BasePoints / 100
}

// The client states rage and energy on a 0-1000 bar: a -30 cost modifier is 3 rage, Charge's energize
// of 150 is 15.
func (e *Effect) Tenths() float64 {
	return e.BasePoints / 10
}

// The value read as a time, which is the unit a duration or delay modifier states it in.
func (e *Effect) TimeValue() time.Duration {
	return time.Duration(e.BasePoints) * time.Millisecond
}

// EffectAuraPeriod: how long one tick of the aura lasts.
func (e *Effect) Period() time.Duration {
	return millis(e.PeriodMs)
}

func (e *Effect) Coeff() float64 {
	return e.SPCoef
}

func (e *Effect) APCoeff() float64 {
	return e.APCoef
}

// The spell this effect fires, or Nil where it fires nothing.
func (e *Effect) Trigger() *Spell {
	return Find(e.TriggerID)
}

// The amount at a caster level: the base points plus the per-level gain over the spell's own level,
// stopping at the level the spell stops scaling at.
//
// float32 is load-bearing: EffectRealPointsPerLevel is a float32 widened into the DB
// (3.79999995231628), and multiplying in float64 moves the result off the tooltip on six rows. The
// client resolves the amount to a whole number, which is the floor below. Same arithmetic as
// DeriveRankAmount in tools/database/spelldata.go, which the generated rank tables are built with.
func (e *Effect) Average(level int32) float64 {
	owner := Find(e.SpellID)

	// MaxLevel 0 is the client's "no cap", which a spell with no SpellLevels row carries.
	lvl := level
	if owner.MaxLevel > 0 && int32(owner.MaxLevel) < lvl {
		lvl = int32(owner.MaxLevel)
	}

	delta := lvl - int32(owner.SpellLevel)
	if delta < 0 {
		delta = 0
	}

	base := float32(e.BasePoints) + float32(float32(delta)*float32(e.PPL))
	return math.Floor(float64(base))
}

// The low end of the roll. EffectVariance is the spread around the average, so an effect that does not
// roll answers the average at both ends.
func (e *Effect) Min(level int32) float64 {
	return e.Average(level) * (1 - e.Variance/2)
}

func (e *Effect) Max(level int32) float64 {
	return e.Average(level) * (1 + e.Variance/2)
}

// The amount for one cast: rolled where the client states a spread, the average where it does not.
func (e *Effect) Roll(sim *core.Simulation, level int32) float64 {
	if e.Variance == 0 {
		return e.Average(level)
	}
	return sim.Roll(e.Min(level), e.Max(level))
}
