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
// TODO: To be implemented. The Forever client ships no rank ladder the generator can
// read for this ability -- it survives as a single spell with no "Rank N" subtext and
// no ranked SkillLineAbility row -- so there is no data to build the spell from.
func (paladin *Paladin) registerHolyShock(rankConfig shared.SpellData) {
	panic("To be implemented")
}
