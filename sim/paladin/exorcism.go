package paladin

import (
	"github.com/wowsims/forever/sim/common/shared"
)

// Exorcism
// https://www.wowhead.com/forever/spell=10314
//
// Causes X to Y Holy damage to an Undead or Demon target.
//
// TODO: To be implemented. spellData.Exorcism holds the six trainer ranks, 879 to 10314.
// The client carries a second ladder, 415068 to 415073, that the Season of Discovery
// passive Exorcist (415076) swaps onto the action bar so the spell can hit any target.
// The generator drops those stand-ins, and nothing in Forever teaches Exorcist, so the
// trainer ranks are the ones to build from.
func (paladin *Paladin) registerExorcism(rankConfig shared.SpellData) {
	// Registered unconditionally, so this returns instead of panicking -- a panic
	// here would stop the sim from starting at all rather than flagging one ability.
}
