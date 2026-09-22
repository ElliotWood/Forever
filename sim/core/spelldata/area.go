package spelldata

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
)

// The client's terrain area groups, each named by the wording of a tooltip that requires it. Every
// other AreaGroup a spell requires is a single zone or instance, which no encounter setting names.
var areaTypeByGroup = map[int32]proto.AreaType{
	9161: proto.AreaType_AreaTypeForestGrassland, // "in Forest and Grassland areas", Grovewalker 1306076
	9165: proto.AreaType_AreaTypeMountainous,     // "in Mountainous areas", Might of Bedrock 1308764
	9164: proto.AreaType_AreaTypeSnowy,           // "in Snowy areas", Bitter Cold 1308511
	9097: proto.AreaType_AreaTypeDesert,          // "in Desert areas", The Desert Rose 1294646
	9162: proto.AreaType_AreaTypeSwamp,           // "in Swamp areas", Sporebloom 1306672
	9163: proto.AreaType_AreaTypeWasteland,       // "in Wasteland areas", Fallen Grove Walker 1306632
	9202: proto.AreaType_AreaTypeHaunted,         // "in Haunted areas", Phantomspeaker Mask 1307334
	9324: proto.AreaType_AreaTypeCavernous,       // "in Cavernous or Underground areas", Royal Seal of Eldre'Thalas 1318480
	9326: proto.AreaType_AreaTypeStrongholdsCities, // "in Strongholds and Cities", Black Horn Necklace 1320332
	// No spell requires 9203; it is the group holding Burning Steppes, Searing Gorge, Blackrock Mountain,
	// Spire and Depths, Molten Core, Blackwing Lair, Onyxia's Lair and Ironforge, which is what Molten
	// Fury's "in Volcanic areas" can only mean.
	9203: proto.AreaType_AreaTypeVolcanic,
}

func AreaTypeOfGroup(group int32) proto.AreaType {
	return areaTypeByGroup[group]
}

func (s *Spell) AreaType() proto.AreaType {
	return AreaTypeOfGroup(s.RequiredAreas)
}

// The factors on the spell's amounts and duration in the given encounter: its area bonus where the
// encounter is in any of the bonus's areas, 1 and 1 everywhere else.
func (s *Spell) AreaBonus(encounter *core.Encounter) (amount, duration float64) {
	for _, group := range s.AreaBonusGroups {
		if encounter.InArea(AreaTypeOfGroup(group)) {
			return float64(s.AreaMultiplier), float64(s.AreaDurationMultiplier)
		}
	}
	return 1, 1
}
