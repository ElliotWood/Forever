package mage

import (
	"testing"

	_ "github.com/wowsims/classic/sim/common"
	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
	googleProto "google.golang.org/protobuf/proto"
)

func init() {
	RegisterMage()
}

func TestP1Mage(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class:      proto.Class_ClassMage,
			Phase:      1,
			Race:       proto.Race_RaceTroll,
			OtherRaces: []proto.Race{proto.Race_RaceGnome},

			Talents:     P1FrostTalents,
			GearSet:     core.GetGearSet("../../ui/mage/gear_sets", "p0.bis"),
			Rotation:    core.GetAplRotation("../../ui/mage/apls", "forever_frost"),
			Buffs:       ForeverBuffs,
			Consumes:    P1Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "DPS", SpecOptions: PlayerOptions},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatSpellPower,
			StatsToWeigh:    Stats,

			Ruleset: proto.Ruleset_RulesetForever,
		},
	}))
}

func TestP1MageArcane(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class: proto.Class_ClassMage,
			Phase: 1,
			Race:  proto.Race_RaceTroll,

			Talents:     P1ArcaneTalents,
			GearSet:     core.GetGearSet("../../ui/mage/gear_sets", "p0.bis"),
			Rotation:    core.GetAplRotation("../../ui/mage/apls", "forever_arcane"),
			Buffs:       ForeverBuffs,
			Consumes:    P1Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "DPS", SpecOptions: PlayerOptions},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatSpellPower,
			StatsToWeigh:    Stats,

			Ruleset: proto.Ruleset_RulesetForever,
		},
	}))
}

func TestP1MageFire(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class: proto.Class_ClassMage,
			Phase: 1,
			Race:  proto.Race_RaceTroll,

			Talents:     P1FireTalents,
			GearSet:     core.GetGearSet("../../ui/mage/gear_sets", "p0.bis"),
			Rotation:    core.GetAplRotation("../../ui/mage/apls", "forever_fire"),
			Buffs:       ForeverBuffs,
			Consumes:    P1Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "DPS", SpecOptions: PlayerOptions},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatSpellPower,
			StatsToWeigh:    Stats,

			Ruleset: proto.Ruleset_RulesetForever,
		},
	}))
}

var P1FrostTalents = "0502050030003--0555000331000301241"
var P1ArcaneTalents = "050215003100311531-235003202003-"
var P1FireTalents = "0502252000003-2350031130133151-"

// Improved Scorch and Winter's Chill only help the mage that applied them in Forever, so the
// raid no longer supplies either debuff.
var ForeverDebuffs = func() *proto.Debuffs {
	debuffs := googleProto.Clone(core.FullDebuffs).(*proto.Debuffs)
	debuffs.ImprovedScorch = false
	debuffs.WintersChill = false
	return debuffs
}()

var ForeverBuffs = core.BuffsCombo{
	Label:   "FullBuffs",
	Debuffs: ForeverDebuffs,
	Party:   core.FullPartyBuffs,
	Player:  core.FullIndividualBuffs,
	Raid:    core.FullRaidBuffs,
}

var PlayerOptions = &proto.Player_Mage{
	Mage: &proto.Mage{
		Options: &proto.Mage_Options{
			Armor: proto.Mage_Options_MoltenArmor,
		},
	},
}

var P1Consumes = core.ConsumesCombo{
	Label: "P1-Consumes",
	Consumes: &proto.Consumes{
		DefaultPotion:  proto.Potions_MajorManaPotion,
		Flask:          proto.Flask_FlaskOfSupremePower,
		FirePowerBuff:  proto.FirePowerBuff_ElixirOfGreaterFirepower,
		FrostPowerBuff: proto.FrostPowerBuff_ElixirOfFrostPower,
		Food:           proto.Food_FoodRunnTumTuberSurprise,
		MainHandImbue:  proto.WeaponImbue_BrilliantWizardOil,
		SpellPowerBuff: proto.SpellPowerBuff_GreaterArcaneElixir,
	},
}

var ItemFilters = core.ItemFilter{
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
}

var Stats = []proto.Stat{
	proto.Stat_StatIntellect,
	proto.Stat_StatSpellPower,
	proto.Stat_StatArcanePower,
	proto.Stat_StatFirePower,
	proto.Stat_StatFrostPower,
	proto.Stat_StatSpellHit,
	proto.Stat_StatSpellCrit,
}
