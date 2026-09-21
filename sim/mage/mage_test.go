package mage

import (
	"testing"

	"github.com/wowsims/forever/sim/common"
	_ "github.com/wowsims/forever/sim/common"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
)

func init() {
	RegisterMage()
	common.RegisterAllEffects()
}

func TestArcane(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{mageSuite("arcane", ArcaneTalents)}))
}

func TestFire(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{mageSuite("fire", FireTalents)}))
}

func TestFrost(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{mageSuite("frost", FrostTalents)}))
}

// The community builds our Forever sim ranks: Arcane 35/0/16, Fire 0/35/16 and Frost 14/0/37.
var ArcaneTalents = "055005023100311531--005500033"
var FireTalents = "-03552020130133151-005500033"
var FrostTalents = "050005013--0555003301001301251"

func mageSuite(apl string, talents string) core.CharacterSuiteConfig {
	return core.CharacterSuiteConfig{
		Class:      proto.Class_ClassMage,
		Race:       proto.Race_RaceGnome,
		OtherRaces: []proto.Race{proto.Race_RaceTroll},
		SpecOptions: core.SpecOptionsCombo{Label: "MageArmor", SpecOptions: &proto.Player_Mage{
			Mage: &proto.Mage{
				Options: &proto.Mage_Options{
					ClassOptions: &proto.MageOptions{
						DefaultMageArmor: proto.MageArmor_MageArmorMageArmor,
					},
				},
			},
		}},
		// Naked: the generated item database does not carry most of the pre-raid set our Forever sim
		// tests with yet, and gives the rest TBC-shaped stats.
		GearSet:  core.GearSetCombo{Label: "Naked", GearSet: &proto.EquipmentSpec{}},
		Talents:  talents,
		Rotation: core.GetAplRotation("../../ui/specs/mage/dps/apls", apl),
		ItemFilter: core.ItemFilter{
			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeDagger,
				proto.WeaponType_WeaponTypeSword,
				proto.WeaponType_WeaponTypeOffHand,
				proto.WeaponType_WeaponTypeStaff,
			},
			ArmorType: proto.ArmorType_ArmorTypeCloth,
			RangedWeaponTypes: []proto.RangedWeaponType{
				proto.RangedWeaponType_RangedWeaponTypeWand,
			},
			EnchantBlacklist: []int32{2673, 3225, 3273},
		},
	}
}
