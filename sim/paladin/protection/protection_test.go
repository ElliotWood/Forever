package protection

import (
	"github.com/wowsims/classic/sim/arenalib"
	"testing"

	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
)

func init() {
	RegisterProtectionPaladin()
}

func TestProtection(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class:      proto.Class_ClassPaladin,
			Phase:      1,
			Race:       proto.Race_RaceHuman,
			OtherRaces: []proto.Race{proto.Race_RaceDwarf},

			Talents:     Phase4ProtTalents,
			GearSet:     core.GetGearSet("../../../ui/protection_paladin/gear_sets", "launch"),
			Rotation:    core.GetAplRotation("../../../ui/protection_paladin/apls", "basic_prot"),
			Buffs:       core.ForeverBuffs,
			Consumes:    Phase4Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Basic Prot Rotation", SpecOptions: PlayerOptionsSealofRighteousness},

			ItemFilter: ItemFilters,
			// Without this the boss never attacks, so nothing the spec does in response to
			// being hit can fire and the damage taken metrics are all zero.
			IsTank:          true,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,

			Ruleset: proto.Ruleset_RulesetForever,
		},
	}))
}

var Phase4ProtTalents = "052003003-5530513321301501"

var Phase4Consumes = core.ConsumesCombo{
	Label: "P4-Consumes",
	Consumes: &proto.Consumes{
		DefaultPotion:     proto.Potions_MajorManaPotion,
		AgilityElixir:     proto.AgilityElixir_ElixirOfTheMongoose,
		AttackPowerBuff:   proto.AttackPowerBuff_JujuMight,
		Flask:             proto.Flask_FlaskOfSupremePower,
		SpellPowerBuff:    proto.SpellPowerBuff_GreaterArcaneElixir,
		DragonBreathChili: true,
		Food:              proto.Food_FoodSmokedDesertDumpling,
		//MainHandImbue:     proto.WeaponImbue_WildStrikes,
		//OffHandImbue:      proto.WeaponImbue_ConductiveShieldCoating,
		StrengthBuff: proto.StrengthBuff_JujuPower,
	},
}

var PlayerOptionsSealofCommand = &proto.Player_ProtectionPaladin{
	ProtectionPaladin: &proto.ProtectionPaladin{
		Options: optionsSealOfCommand,
	},
}

var PlayerOptionsSealofRighteousness = &proto.Player_ProtectionPaladin{
	ProtectionPaladin: &proto.ProtectionPaladin{
		Options: optionsSealOfRighteousness,
	},
}

var optionsSealOfCommand = &proto.PaladinOptions{
	PrimarySeal:   proto.PaladinSeal_Command,
	RighteousFury: true,
}

var optionsSealOfRighteousness = &proto.PaladinOptions{
	PrimarySeal:   proto.PaladinSeal_Righteousness,
	RighteousFury: true,
}

var ItemFilters = core.ItemFilter{
	WeaponTypes: []proto.WeaponType{
		proto.WeaponType_WeaponTypeAxe,
		proto.WeaponType_WeaponTypeSword,
		proto.WeaponType_WeaponTypeMace,
		proto.WeaponType_WeaponTypePolearm,
		proto.WeaponType_WeaponTypeShield,
	},
	RangedWeaponTypes: []proto.RangedWeaponType{
		proto.RangedWeaponType_RangedWeaponTypeLibram,
	},
}

var Stats = []proto.Stat{
	proto.Stat_StatHealth,
	proto.Stat_StatMana,
	proto.Stat_StatStrength,
	proto.Stat_StatStamina,
	proto.Stat_StatAgility,
	proto.Stat_StatIntellect,
	proto.Stat_StatAttackPower,
	proto.Stat_StatMeleeHit,
	proto.Stat_StatMeleeCrit,
	proto.Stat_StatMeleeHaste,
	proto.Stat_StatSpellHit,
	proto.Stat_StatSpellCrit,
	proto.Stat_StatSpellPower,
	proto.Stat_StatHolyPower,
	proto.Stat_StatHealingPower,
	proto.Stat_StatArmor,
	proto.Stat_StatBonusArmor,
	proto.Stat_StatDefense,
	proto.Stat_StatDodge,
	proto.Stat_StatParry,
	proto.Stat_StatBlock,
	proto.Stat_StatBlockValue,
	proto.Stat_StatFireResistance,
	proto.Stat_StatNatureResistance,
	proto.Stat_StatShadowResistance,
	proto.Stat_StatFrostResistance,
	proto.Stat_StatArcaneResistance,
}

// The arena entry for this spec. Skipped unless ARENA_OUT is set; see sim/arenalib.
func TestArena(t *testing.T) {
	arenalib.Run(t, arenalib.Spec{
		Dir:         "protection_paladin",
		Class:       proto.Class_ClassPaladin,
		Race:        proto.Race_RaceHuman,
		SpecOptions: PlayerOptionsSealofRighteousness,
		Consumes:    Phase4Consumes,
		Buffs:       core.ForeverBuffs,
		IsTank:      true,
	})
}
