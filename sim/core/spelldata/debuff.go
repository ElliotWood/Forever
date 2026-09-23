package spelldata

import (
	"slices"

	"github.com/wowsims/forever/sim/core/dbcenums"
)

// The speed auras that slow a unit's attacks or casts.
var speedAuras = []dbcenums.EffectAuraType{
	dbcenums.A_MOD_ATTACKSPEED, dbcenums.A_MOD_CASTING_SPEED_NOT_STACK, dbcenums.A_MOD_MELEE_HASTE_3,
}

// The positions, counted from 1 the way EffectN counts, of the effects that slow the enemy the spell
// lands on: a speed aura of a negative value on an enemy target, as Frostguard's Chilled 16927 states.
// A row that stacks answers none, since an exclusive slow cannot follow stacks.
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

// The positions of the effects that change a stat of the enemy the spell lands on, in the stats
// ParseEffects applies to a unit: armor and resistances, attack power, flat damage done. Annihilator's
// Armor Shatter 16928 takes 165 armor per stack.
func (s *Spell) StatDebuffEffects() []int32 {
	var debuffs []int32
	for i := range s.Effects {
		e := &s.Effects[i]
		if appliesAura(e.Type) && e.HitsAnEnemy() && debuffsAStat(e) {
			debuffs = append(debuffs, int32(i+1))
		}
	}
	return debuffs
}

func debuffsAStat(e *Effect) bool {
	switch e.Aura {
	case dbcenums.A_MOD_RESISTANCE:
		return len(resistanceStats(e.Misc)) > 0
	case dbcenums.A_MOD_DAMAGE_DONE:
		return len(damageDoneStats(e.Misc)) > 0
	case dbcenums.A_MOD_ATTACK_POWER, dbcenums.A_MOD_RANGED_ATTACK_POWER:
		return true
	}
	return false
}

// Whether the row puts a debuff the sim models on the enemy it lands on.
func (s *Spell) DebuffsTheTarget() bool {
	return len(s.SlowEffects()) > 0 || len(s.StatDebuffEffects()) > 0
}

// Whether any aura the row applies lands on an enemy. Such a row is never a buff on the wearer.
func (s *Spell) AppliesAnAuraToAnEnemy() bool {
	return slices.ContainsFunc(s.Effects, func(e Effect) bool { return appliesAura(e.Type) && e.HitsAnEnemy() })
}
