package spelldata

import (
	"slices"

	"github.com/wowsims/forever/sim/core/dbcenums"
)

// The speed auras ParseEffects multiplies a unit's attacks or casts by.
var speedAuras = []dbcenums.EffectAuraType{
	dbcenums.A_MOD_ATTACKSPEED, dbcenums.A_MOD_CASTING_SPEED_NOT_STACK, dbcenums.A_MOD_MELEE_HASTE_3,
}

// The positions, counted from 1 the way EffectN counts, of the effects that slow the enemy the spell
// lands on: a speed aura of a negative value on an enemy target, as Frostguard's Chilled 16927 states.
// A row that stacks answers none, since a speed multiplier cannot follow stacks.
func (s *Spell) SlowEffects() []int32 {
	if s.MaxStack > 0 {
		return nil
	}

	var slows []int32
	for i := range s.Effects {
		e := &s.Effects[i]
		if appliesAura(e.Type) && slices.Contains(speedAuras, e.Aura) && e.BasePoints < 0 && e.HitsAnEnemy() {
			slows = append(slows, int32(i+1))
		}
	}
	return slows
}
