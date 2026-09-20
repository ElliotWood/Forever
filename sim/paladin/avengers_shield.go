package paladin

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

func (paladin *Paladin) getAvengersShieldTimer() *core.Timer {
	if paladin.avengersShieldTimer == nil {
		paladin.avengersShieldTimer = paladin.NewTimer()
	}
	return paladin.avengersShieldTimer
}

// Avenger's Shield (Talent)
// https://www.wowhead.com/forever/spell=31935
//
// Hurls a holy shield at the enemy, dealing Holy damage, dazing them and
// then jumping to additional nearby enemies. Affects 3 total targets.
//
// TODO: uncalled -- Forever drops the Avenger's Shield talent; re-gate before wiring
// back into registerTalentSpells.
// TODO: To be implemented. The Forever client ships no rank ladder the generator can
// read for this ability -- it survives as a single spell with no "Rank N" subtext and
// no ranked SkillLineAbility row -- so there is no data to build the spell from.
func (paladin *Paladin) registerAvengersShield(rankConfig shared.SpellData) {
	// Registered unconditionally, so this returns instead of panicking -- a panic
	// here would stop the sim from starting at all rather than flagging one ability.
}
