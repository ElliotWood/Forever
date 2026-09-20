package paladin

import (
	"github.com/wowsims/forever/sim/common/shared"
)

// Exorcism
// https://www.wowhead.com/forever/spell=10314
//
// Causes X to Y Holy damage to an Undead or Demon target.
//
// TODO: To be implemented. The Forever client ships no rank ladder the generator can
// read for this ability -- it survives as a single spell with no "Rank N" subtext and
// no ranked SkillLineAbility row -- so there is no data to build the spell from.
func (paladin *Paladin) registerExorcism(rankConfig shared.SpellData) {
	// Registered unconditionally, so this returns instead of panicking -- a panic
	// here would stop the sim from starting at all rather than flagging one ability.
}
