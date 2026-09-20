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
// TODO: To be implemented. The Forever client ships a rank ladder for this, but each rank has
// two candidate spells that the resolver cannot separate (rank 2 is 25912 and 25914), so
// spellData carries no table for it yet. See the "Not generated" list in the generated file.
func (paladin *Paladin) registerHolyShock(rankConfig shared.SpellData) {
	// Registered unconditionally, so this returns instead of panicking -- a panic
	// here would stop the sim from starting at all rather than flagging one ability.
}
