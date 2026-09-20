package paladin

import (
	"github.com/wowsims/forever/sim/common/shared"
)

// Exorcism
// https://www.wowhead.com/forever/spell=10314
//
// Causes X to Y Holy damage to an Undead or Demon target.
//
// TODO: To be implemented. The Forever client DOES ship a rank ladder for this, and spellData
// now carries it -- the family was previously dropped as ambiguous because Forever re-issues
// the ability as a second spell per rank. The body below is the TBC implementation, awaiting
// a port onto the recovered ladder.
func (paladin *Paladin) registerExorcism(rankConfig shared.SpellData) {
	// Registered unconditionally, so this returns instead of panicking -- a panic
	// here would stop the sim from starting at all rather than flagging one ability.
}
