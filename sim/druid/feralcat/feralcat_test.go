package feralcat

import (
	"testing"

	"github.com/wowsims/forever/sim/common"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
)

func init() {
	RegisterFeralCatDruid()
	common.RegisterAllEffects()
}

func TestFeralCat(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class:      proto.Class_ClassDruid,
			Race:       proto.Race_RaceNightElf,
			OtherRaces: []proto.Race{proto.Race_RaceTauren},

			// Naked: the generated item database does not carry the Forever gear our sim tests
			// with yet, and gives the rest TBC-shaped stats.
			GearSet: core.GearSetCombo{Label: "Naked", GearSet: &proto.EquipmentSpec{}},

			Talents: DefaultTalents,
			OtherTalentSets: []core.TalentsCombo{
				{Label: "FeralCat", Talents: FeralCatTalents},
			},

			SpecOptions: core.SpecOptionsCombo{Label: "Standard", SpecOptions: DefaultSpecOptions},

			Rotation: core.GetAplRotation("../../../ui/specs/druid/feralcat/apls", "default"),

			Consumables: DefaultConsumables,

			Profession1: proto.Profession_Engineering,
			Profession2: proto.Profession_Enchanting,

			ItemFilter: core.ItemFilter{
				ArmorType: proto.ArmorType_ArmorTypeLeather,
				WeaponTypes: []proto.WeaponType{
					proto.WeaponType_WeaponTypeDagger,
					proto.WeaponType_WeaponTypeFist,
					proto.WeaponType_WeaponTypeMace,
					proto.WeaponType_WeaponTypeStaff,
				},
				RangedWeaponTypes: []proto.RangedWeaponType{
					proto.RangedWeaponType_RangedWeaponTypeIdol,
				},
			},

			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh: []proto.Stat{
				proto.Stat_StatAgility,
				proto.Stat_StatStrength,
				proto.Stat_StatAttackPower,
				proto.Stat_StatFeralAttackPower,
				proto.Stat_StatMeleeHitRating,
				proto.Stat_StatExpertiseRating,
				proto.Stat_StatMeleeCritRating,
				proto.Stat_StatMeleeHasteRating,
				proto.Stat_StatArmorPenetration,
			},
		},
	}))
}

// Our Forever sim's feral builds.
const DefaultTalents = "-5521002023132213051-05503"
const FeralCatTalents = "050022-5500002123032213051-052"

var DefaultSpecOptions = &proto.Player_FeralCatDruid{
	FeralCatDruid: &proto.FeralCatDruid{
		Rotation: &proto.FeralCatDruid_Rotation{
			FinishingMove:      proto.FeralCatDruid_Rotation_Rip,
			Biteweave:          true,
			RipMinComboPoints:  5,
			BiteMinComboPoints: 5,
			MangleTrick:        true,
			MaintainFaerieFire: false,
		},
		Options: &proto.FeralCatDruid_Options{},
	},
}

var DefaultConsumables = &proto.ConsumesSpec{
	PotId:            22838, // Haste Potion
	BattleElixirId:   22831, // Elixir of Major Agility
	GuardianElixirId: 32067, // Elixir of Draenic Wisdom
	FoodId:           27664, // Grilled Mudfish
	MhImbueId:        34340, // Adamantite Weightstone
	ConjuredId:       12662, // Demonic Rune
	SuperSapper:      true,
	GoblinSapper:     true,
	ScrollAgi:        true,
	ScrollStr:        true,
}
