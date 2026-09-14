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

// Neither the phase nor the zone check can see two groups of items. The Ahn'Qiraj and
// Naxxramas tier sets reach the database through the quests that turn their tokens in, so
// they carry a quest source and a phase of 1. Nefarian's head turn-ins carry no source
// record at all. Item level separates both from everything at launch: Sulfuras at 80 is
// the highest thing Molten Core or Onyxia gives up, and above it sit only the tier 2.5
// sets at 81, the Master Dragonslayer's rewards at 83 and the tier 3 sets from 86 to 92.
const maxLaunchIlvl = 80

// ObtainableAtLaunch reports whether an item exists in Forever at launch.
func ObtainableAtLaunch(item *proto.UIItem) bool {
	if item.Phase > lastLaunchPhase || item.Ilvl > maxLaunchIlvl {
		return false
	}
	for _, source := range item.Sources {
		if drop := source.GetDrop(); drop != nil && postLaunchZones[drop.ZoneId] {
			continue
		}
		return true
	}
	return len(item.Sources) == 0
}
