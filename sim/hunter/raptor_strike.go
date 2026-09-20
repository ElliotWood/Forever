package hunter

import (
	"github.com/wowsims/forever/sim/core"
)

// TODO: To be implemented. The Forever client DOES ship a rank ladder for this, and spellData
// now carries it -- the family was previously dropped as ambiguous because Forever re-issues
// the ability as a second spell per rank. The body below is the TBC implementation, awaiting
// a port onto the recovered ladder.
func (hunter *Hunter) registerRaptorStrikeSpell() {
	// Registered unconditionally, so this returns instead of panicking -- a panic
	// here would stop the sim from starting at all rather than flagging one ability.
}

// Returns true if the regular melee swing should be used, false otherwise.
func (hunter *Hunter) TryRaptorStrike(sim *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
	if mhSwingSpell.ActionID.Tag != 1 || !hunter.RaptorStrike.CanCast(sim, hunter.CurrentTarget) {
		return mhSwingSpell
	}

	return hunter.RaptorStrike
}
