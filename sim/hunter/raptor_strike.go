package hunter

import (
	"github.com/wowsims/forever/sim/core"
)

// TODO: To be implemented. spellData.RaptorStrike holds the eight trainer ranks, 2973 to 14266.
// The client also carries two Season of Discovery ladders under the same name: 415335 to
// 415343, which the rune passive Melee Specialist (415352) swaps onto the action bar, and
// 409691 to 409755, a mana-less copy nothing references. The generator drops the first by
// its override link and the second because only the trainer rank carries a SpellPower row.
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
