package paladin

import (
	"github.com/wowsims/forever/sim/common/shared"
)

// Holy Shock
// https://www.wowhead.com/forever/spell=20473
//
// Blasts the target with Holy energy, causing X to Y Holy damage to an enemy,
// or X*1.267 to Y*1.267 healing to an ally.
//
// TODO: To be implemented. The client ships TWO complete parallel ladders under this one
// name, not one ambiguous ladder: a damage chain (25912, 25911, 25902 -- SpellEffect 2) and
// a heal chain (25914, 25913, 25903 -- SpellEffect 10). discoverLadders keys ladders by spell
// name alone, so the two collide at every rank and the resolver refuses the family outright.
// Splitting parallel ladders by effect type, rather than picking one, is what this needs;
// Renew and Mutilate have the same shape. Until then spellData carries no HolyShock table.
func (paladin *Paladin) registerHolyShock(rankConfig shared.SpellData) {
	// Registered unconditionally, so this returns instead of panicking -- a panic
	// here would stop the sim from starting at all rather than flagging one ability.
}
