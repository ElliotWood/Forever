package arenalib

import (
	googleProto "google.golang.org/protobuf/proto"

	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
)

// What every build in the arena drinks, decided here rather than by each spec's own test file.
//
// The file header claims every build meets the same fixed environment, and until this existed
// that was false. Fifteen specs brought four different consumable sets from four different
// content phases, and the gaps were not small ones: both paladins and the feral tank had their
// weapon imbue commented out entirely, while warrior, hunter, rogue, tank warrior and
// enhancement all carried Windfury. Measured on the warrior, stripping its imbues cost 14.2%
// of its damage - so the leaderboard was reporting a 19.6% gap between warrior and retribution
// while handing one of them a weapon buff and the other a bare weapon.
//
// Windfury deserves a note, because it is the biggest single item here and it is not really a
// consumable at all: WeaponImbue_Windfury applies spell 10610, which is Windfury Totem. The
// buffs proto has nowhere to put a totem, so the consumables slot is the only lever there is.
// Giving it to every melee matches what the arena already does everywhere else - every spec
// gets core.ForeverBuffs, which hands all four paladin blessings to Horde classes. The arena
// is a fixed environment, not a raid roster, and it has never modelled faction.
//
// Three sets rather than one, because one list for all fifteen would be a different lie.
// Elemental Sharpening Stone is +2% melee crit and -2% RANGED crit, so handing the hunter the
// melee list would equalise the shopping and quietly tax the one spec that shoots. Within a
// role every spec gets the identical list; what it is worth to you is your class's business,
// which is why Mighty Rage Potion stays in the melee list even though only warriors can spend
// it. Same shopping list, not same benefit - that is the part that makes two numbers comparable.
type Role int

const (
	Melee Role = iota
	Ranged
	Caster
)

var consumesMelee = core.ConsumesCombo{
	Label: "Arena-Melee",
	Consumes: &proto.Consumes{
		AgilityElixir:     proto.AgilityElixir_ElixirOfTheMongoose,
		AttackPowerBuff:   proto.AttackPowerBuff_JujuMight,
		StrengthBuff:      proto.StrengthBuff_JujuPower,
		DefaultPotion:     proto.Potions_MightyRagePotion,
		DragonBreathChili: true,
		Flask:             proto.Flask_FlaskOfTheTitans,
		Food:              proto.Food_FoodSmokedDesertDumpling,
		MainHandImbue:     proto.WeaponImbue_Windfury,
		OffHandImbue:      proto.WeaponImbue_ElementalSharpeningStone,
	},
}

// The melee list without the off-hand stone, whose ranged crit penalty is a real cost to the
// one spec that does its damage from thirty yards.
var consumesRanged = core.ConsumesCombo{
	Label: "Arena-Ranged",
	Consumes: &proto.Consumes{
		AgilityElixir:     proto.AgilityElixir_ElixirOfTheMongoose,
		AttackPowerBuff:   proto.AttackPowerBuff_JujuMight,
		StrengthBuff:      proto.StrengthBuff_JujuPower,
		DragonBreathChili: true,
		Flask:             proto.Flask_FlaskOfTheTitans,
		Food:              proto.Food_FoodSmokedDesertDumpling,
		MainHandImbue:     proto.WeaponImbue_Windfury,
	},
}

var consumesCaster = core.ConsumesCombo{
	Label: "Arena-Caster",
	Consumes: &proto.Consumes{
		DefaultPotion:  proto.Potions_MajorManaPotion,
		Flask:          proto.Flask_FlaskOfSupremePower,
		Food:           proto.Food_FoodSmokedDesertDumpling,
		MainHandImbue:  proto.WeaponImbue_BrilliantWizardOil,
		SpellPowerBuff: proto.SpellPowerBuff_GreaterArcaneElixir,
		FirePowerBuff:  proto.FirePowerBuff_ElixirOfFirepower,
	},
}

// A weapon imbue a class grants itself, which the role list must not overwrite.
//
// The role lists equalise what a character BUYS. Some imbues are not bought: Windfury Weapon
// is an enhancement shaman casting on their own weapons, and Instant Poison is a rogue's
// poison. Equalising those does not make a comparison fairer, it takes a class ability away -
// the first pass did exactly that and cost enhancement 23.6% of its damage, which is not a
// shaman being measured honestly, it is a shaman being disarmed.
type ClassImbues struct {
	MainHand proto.WeaponImbue
	OffHand  proto.WeaponImbue
}

func consumesFor(role Role, imbues ClassImbues) core.ConsumesCombo {
	var combo core.ConsumesCombo
	switch role {
	case Ranged:
		combo = consumesRanged
	case Caster:
		combo = consumesCaster
	default:
		combo = consumesMelee
	}
	if imbues.MainHand == proto.WeaponImbue_WeaponImbueUnknown && imbues.OffHand == proto.WeaponImbue_WeaponImbueUnknown {
		return combo
	}

	// Cloned, not mutated: the role lists are package level and shared by every spec, so
	// writing a shaman's Windfury into one would hand it to the next spec that ran. Cloned
	// rather than copied because a proto carries a mutex.
	consumes := googleProto.Clone(combo.Consumes).(*proto.Consumes)
	if imbues.MainHand != proto.WeaponImbue_WeaponImbueUnknown {
		consumes.MainHandImbue = imbues.MainHand
	}
	if imbues.OffHand != proto.WeaponImbue_WeaponImbueUnknown {
		consumes.OffHandImbue = imbues.OffHand
	}
	return core.ConsumesCombo{Label: combo.Label + "+class", Consumes: consumes}
}
