package shadow

import (
	"testing"

	_ "github.com/wowsims/classic/sim/common" // imported to get caster sets included.
	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
	googleProto "google.golang.org/protobuf/proto"
)

func init() {
	RegisterShadowPriest()
}

func TestP1Shadow(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class:      proto.Class_ClassPriest,
			Phase:      1,
			Race:       proto.Race_RaceUndead,
			OtherRaces: []proto.Race{proto.Race_RaceDwarf},

			Talents:     P1Talents,
			GearSet:     core.GetGearSet("../../../ui/shadow_priest/gear_sets", "p0.bis"),
			Rotation:    core.GetAplRotation("../../../ui/shadow_priest/apls", "p1"),
			Buffs:       ForeverBuffs,
			Consumes:    P1Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsBasic},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatSpellPower,
			StatsToWeigh:    Stats,

			Ruleset: proto.Ruleset_RulesetForever,
		},
	}))
}

var P1Talents = "005300231303--505120501201300051"

// The raid debuff version of Shadow Weaving is a Classic mechanic, in Forever it buffs the priest instead.
var ForeverDebuffs = func() *proto.Debuffs {
	debuffs := googleProto.Clone(core.FullDebuffs).(*proto.Debuffs)
	debuffs.ShadowWeaving = false
	return debuffs
}()

var ForeverBuffs = core.BuffsCombo{
	Label:   "FullBuffs",
	Debuffs: ForeverDebuffs,
	Party:   core.FullPartyBuffs,
	Player:  core.FullIndividualBuffs,
	Raid:    core.FullRaidBuffs,
}

var P1Consumes = core.ConsumesCombo{
	Label: "P1-Consumes",
	Consumes: &proto.Consumes{
		DefaultPotion:   proto.Potions_MajorManaPotion,
		Flask:           proto.Flask_FlaskOfSupremePower,
		Food:            proto.Food_FoodRunnTumTuberSurprise,
		MainHandImbue:   proto.WeaponImbue_WizardOil,
		SpellPowerBuff:  proto.SpellPowerBuff_GreaterArcaneElixir,
		ShadowPowerBuff: proto.ShadowPowerBuff_ElixirOfShadowPower,
	},
}

var PlayerOptionsBasic = &proto.Player_ShadowPriest{
	ShadowPriest: &proto.ShadowPriest{
		Options: &proto.ShadowPriest_Options{
			Armor: proto.ShadowPriest_Options_InnerFire,
		},
	},
}

var ItemFilters = core.ItemFilter{
	WeaponTypes: []proto.WeaponType{
		proto.WeaponType_WeaponTypeDagger,
		proto.WeaponType_WeaponTypeMace,
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
	proto.Stat_StatShadowPower,
	proto.Stat_StatSpellHit,
	proto.Stat_StatSpellCrit,
}
