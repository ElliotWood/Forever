package dps

import (
	"github.com/wowsims/classic/sim/arenalib"
	"testing"

	_ "github.com/wowsims/classic/sim/common"
	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
)

func init() {
	RegisterDpsWarlock()
}

func TestWarlockDemonicPact(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class: proto.Class_ClassWarlock,
			Phase: 1,
			Race:  proto.Race_RaceOrc,

			Talents:     TalentsDemonicPact,
			GearSet:     core.GetGearSet("../../../ui/warlock/gear_sets", "prebis"),
			Rotation:    core.GetAplRotation("../../../ui/warlock/apls/", "forever_pact"),
			Buffs:       core.ForeverBuffs,
			Consumes:    Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Demonic Pact Warlock", SpecOptions: DefaultPactWarlock},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatSpellPower,
			StatsToWeigh:    Stats,

			Ruleset: proto.Ruleset_RulesetForever,
		},
	}))
}

func TestWarlockAffliction(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class: proto.Class_ClassWarlock,
			Phase: 1,
			Race:  proto.Race_RaceOrc,

			Talents:     TalentsAffliction,
			GearSet:     core.GetGearSet("../../../ui/warlock/gear_sets", "prebis"),
			Rotation:    core.GetAplRotation("../../../ui/warlock/apls/", "forever_affliction"),
			Buffs:       core.ForeverBuffs,
			Consumes:    Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Affliction Warlock", SpecOptions: DefaultDestroWarlock},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatSpellPower,
			StatsToWeigh:    Stats,

			Ruleset: proto.Ruleset_RulesetForever,
		},
	}))
}

func TestWarlockDSRuin(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class: proto.Class_ClassWarlock,
			Phase: 1,
			Race:  proto.Race_RaceOrc,

			Talents:     TalentsDSRuin,
			GearSet:     core.GetGearSet("../../../ui/warlock/gear_sets", "prebis"),
			Rotation:    core.GetAplRotation("../../../ui/warlock/apls/", "forever_ds_ruin"),
			Buffs:       core.ForeverBuffs,
			Consumes:    Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "DS/Ruin Warlock", SpecOptions: DefaultImpWarlock},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatSpellPower,
			StatsToWeigh:    Stats,

			Ruleset: proto.Ruleset_RulesetForever,
		},
	}))
}

// The community builds the rankings page runs, so the goldens measure what is ranked.
var TalentsDemonicPact = "113-0005003221220311351-0550005"
var TalentsAffliction = "2535002013521105--05500051"
var TalentsDSRuin = "25220010135201-0025003001-05500051"
var TalentsShadowAndFlame = "25501-0025003001-055035510010002"

func TestWarlockShadowAndFlame(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class: proto.Class_ClassWarlock,
			Phase: 1,
			Race:  proto.Race_RaceOrc,

			Talents:     TalentsShadowAndFlame,
			GearSet:     core.GetGearSet("../../../ui/warlock/gear_sets", "prebis"),
			Rotation:    core.GetAplRotation("../../../ui/warlock/apls/", "forever_shadow_and_flame"),
			Buffs:       core.ForeverBuffs,
			Consumes:    Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Shadow and Flame Warlock", SpecOptions: DefaultImpWarlock},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatSpellPower,
			StatsToWeigh:    Stats,

			Ruleset: proto.Ruleset_RulesetForever,
		},
	}))
}

var DefaultPactWarlock = &proto.Player_Warlock{
	Warlock: &proto.Warlock{
		Options: &proto.WarlockOptions{
			Armor:       proto.WarlockOptions_DemonArmor,
			Summon:      proto.WarlockOptions_Succubus,
			Sacrifice:   proto.WarlockOptions_Voidwalker,
			WeaponImbue: proto.WarlockOptions_NoWeaponImbue,
		},
	},
}

var DefaultImpWarlock = &proto.Player_Warlock{
	Warlock: &proto.Warlock{
		Options: &proto.WarlockOptions{
			Armor:       proto.WarlockOptions_DemonArmor,
			Summon:      proto.WarlockOptions_Imp,
			WeaponImbue: proto.WarlockOptions_NoWeaponImbue,
		},
	},
}

var DefaultDestroWarlock = &proto.Player_Warlock{
	Warlock: &proto.Warlock{
		Options: &proto.WarlockOptions{
			Armor:       proto.WarlockOptions_DemonArmor,
			Summon:      proto.WarlockOptions_Succubus,
			WeaponImbue: proto.WarlockOptions_NoWeaponImbue,
		},
	},
}

var Consumes = core.ConsumesCombo{
	Label: "Consumes",
	Consumes: &proto.Consumes{
		DefaultPotion:   proto.Potions_MajorManaPotion,
		Flask:           proto.Flask_FlaskOfSupremePower,
		FirePowerBuff:   proto.FirePowerBuff_ElixirOfGreaterFirepower,
		ShadowPowerBuff: proto.ShadowPowerBuff_ElixirOfShadowPower,
		Food:            proto.Food_FoodTenderWolfSteak,
		MainHandImbue:   proto.WeaponImbue_WizardOil,
		SpellPowerBuff:  proto.SpellPowerBuff_GreaterArcaneElixir,
	},
}

var ItemFilters = core.ItemFilter{
	WeaponTypes: []proto.WeaponType{
		proto.WeaponType_WeaponTypeSword,
		proto.WeaponType_WeaponTypeDagger,
	},
	HandTypes: []proto.HandType{
		proto.HandType_HandTypeOffHand,
	},
	ArmorType: proto.ArmorType_ArmorTypeCloth,
	RangedWeaponTypes: []proto.RangedWeaponType{
		proto.RangedWeaponType_RangedWeaponTypeWand,
	},
}

var Stats = []proto.Stat{
	proto.Stat_StatIntellect,
	proto.Stat_StatSpellPower,
	proto.Stat_StatSpellHit,
	proto.Stat_StatSpellCrit,
}

// The arena entry for this spec. Skipped unless ARENA_OUT is set; see sim/arenalib.
func TestArena(t *testing.T) {
	arenalib.Run(t, arenalib.Spec{
		Dir:                "warlock",
		Class:              proto.Class_ClassWarlock,
		Race:               proto.Race_RaceOrc,
		SpecOptions:        DefaultPactWarlock,
		Consumes:           Consumes,
		Buffs:              core.ForeverBuffs,
		DistanceFromTarget: 30,
	})
}
