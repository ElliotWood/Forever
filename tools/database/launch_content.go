package database

import (
	"github.com/wowsims/classic/sim/core/proto"
)

// Forever reorders the raid release: Onyxia's Lair opens first, then Molten Core, and
// nothing beyond that exists at launch. Classic's phases 4, 5 and 6 were the later raid
// releases - Zul'Gurub, Ahn'Qiraj and Naxxramas - so cutting by phase covers those. Phases
// 1 to 3 brought the dungeons, the quest and crafted gear, the reputation rewards and the
// PvP sets alongside their raids, so within them it is only Blackwing Lair that has to go,
// which is a check on the source rather than the phase.
const lastLaunchPhase = 3

// The one raid zone inside phases 1-3 that Forever does not open at launch. Molten Core
// (2717) and Onyxia's Lair (2159) both do, so neither is listed here.
var postLaunchZones = map[int32]bool{
	2677: true, // Blackwing Lair
}

// Nefarian's head turn-ins carry no source record at all, so the zone check cannot see
// them. Item level can: among items with no source and a phase of 3 or lower, the only
// ones above 80 are the Master Dragonslayer's Orb, Medallion and Ring, all rewards for
// handing in his head. Everything at or below it is launch content - Onyxia's own drops at
// 71 and 72, her turn-in rewards at 74, the Benediction and Rhok'delar chains at 75 and
// Sulfuras at 80, none of which have a source record either.
const maxSourcelessIlvl = 80

// ObtainableAtLaunch reports whether an item exists in Forever at launch.
func ObtainableAtLaunch(item *proto.UIItem) bool {
	if item.Phase > lastLaunchPhase {
		return false
	}
	for _, source := range item.Sources {
		if drop := source.GetDrop(); drop != nil && postLaunchZones[drop.ZoneId] {
			continue
		}
		return true
	}
	return len(item.Sources) == 0 && item.Ilvl <= maxSourcelessIlvl
}
