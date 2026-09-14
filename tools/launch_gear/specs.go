package main

import (
	"sort"

	"github.com/wowsims/classic/sim/core/proto"
	"github.com/wowsims/classic/sim/core/stats"
)

// The weights are stated, not derived. Every spec this tool runs for is one with no gear
// set at all, so there is nothing to compute stat weights from - a naked character's
// weights describe a naked character. They are rough on purpose: the job is to put a
// plausible launch set on a spec that currently fights in its underwear, not to find BiS.
//
// The tanks are weighted for threat rather than survival, because threat per second is
// what their sim measures.
type spec struct {
	class       proto.Class
	armor       proto.ArmorType
	weapons     []proto.WeaponType
	usesShield  bool
	allowTwoHnd bool
	dualWield   bool
	// Weapon damage is most of a melee weapon's value and none of a caster's.
	weaponDpsEP float64
	weights     stats.Stats
}

func (s spec) allows(item *proto.UIItem) bool {
	if len(item.ClassAllowlist) > 0 {
		found := false
		for _, c := range item.ClassAllowlist {
			if c == s.class {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	if item.Type == proto.ItemType_ItemTypeHead || item.Type == proto.ItemType_ItemTypeShoulder ||
		item.Type == proto.ItemType_ItemTypeChest || item.Type == proto.ItemType_ItemTypeWrist ||
		item.Type == proto.ItemType_ItemTypeHands || item.Type == proto.ItemType_ItemTypeWaist ||
		item.Type == proto.ItemType_ItemTypeLegs || item.Type == proto.ItemType_ItemTypeFeet {
		if item.ArmorType != s.armor && item.ArmorType != proto.ArmorType_ArmorTypeUnknown {
			return false
		}
	}
	return true
}

func (s spec) pickWeapons(pool []*proto.UIItem, weights stats.Stats, picked *[]map[string]any) {
	best := func(match func(*proto.UIItem) bool) *proto.UIItem {
		cands := []*proto.UIItem{}
		for _, item := range pool {
			if match(item) {
				cands = append(cands, item)
			}
		}
		if len(cands) == 0 {
			return nil
		}
		sort.SliceStable(cands, func(i, j int) bool { return s.ep(cands[i]) > s.ep(cands[j]) })
		return cands[0]
	}

	weaponAllowed := func(item *proto.UIItem) bool {
		if len(s.weapons) == 0 {
			return true
		}
		for _, w := range s.weapons {
			if item.WeaponType == w {
				return true
			}
		}
		return false
	}

	isWeapon := func(item *proto.UIItem) bool {
		return item.Type == proto.ItemType_ItemTypeWeapon && weaponAllowed(item)
	}
	oneHand := func(item *proto.UIItem) bool {
		return isWeapon(item) && item.HandType != proto.HandType_HandTypeTwoHand
	}

	add := func(item *proto.UIItem) {
		if item != nil {
			*picked = append(*picked, map[string]any{"id": item.Id})
		}
	}

	if s.usesShield {
		add(best(oneHand))
		add(best(func(item *proto.UIItem) bool {
			return item.WeaponType == proto.WeaponType_WeaponTypeShield
		}))
		return
	}

	mh := best(oneHand)
	if s.allowTwoHnd {
		th := best(func(item *proto.UIItem) bool {
			return isWeapon(item) && item.HandType == proto.HandType_HandTypeTwoHand
		})
		// A two hander gives up an off hand only if the spec had one to give up.
		rival := 1.0
		if s.dualWield {
			rival = 2.0
		}
		if th != nil && (mh == nil || s.ep(th) > rival*s.ep(mh)) {
			add(th)
			return
		}
	}
	add(mh)
	if s.dualWield {
		add(best(func(item *proto.UIItem) bool { return oneHand(item) && item != mh }))
	}
}

var (
	plate   = proto.ArmorType_ArmorTypePlate
	mail    = proto.ArmorType_ArmorTypeMail
	leather = proto.ArmorType_ArmorTypeLeather
	cloth   = proto.ArmorType_ArmorTypeCloth
)

func w(pairs map[stats.Stat]float64) stats.Stats {
	var out stats.Stats
	for k, v := range pairs {
		out[k] = v
	}
	return out
}

var specs = map[string]spec{
	// Melee damage with a holy component; Champion of the Light scales it off Intellect.
	"retribution_paladin": {
		class: proto.Class_ClassPaladin, armor: plate, allowTwoHnd: true,
		weapons:     []proto.WeaponType{proto.WeaponType_WeaponTypeAxe, proto.WeaponType_WeaponTypeMace, proto.WeaponType_WeaponTypeSword, proto.WeaponType_WeaponTypePolearm},
		weaponDpsEP: 12.0, weights: w(map[stats.Stat]float64{stats.Strength: 1, stats.AttackPower: 0.5, stats.Agility: 0.35,
			stats.MeleeCrit: 12, stats.MeleeHit: 14, stats.Intellect: 0.3, stats.SpellPower: 0.35, stats.Stamina: 0.1}),
	},
	// Threat, so damage stats lead and mitigation is a tiebreaker.
	"protection_paladin": {
		class: proto.Class_ClassPaladin, armor: plate, usesShield: true,
		weapons:     []proto.WeaponType{proto.WeaponType_WeaponTypeAxe, proto.WeaponType_WeaponTypeMace, proto.WeaponType_WeaponTypeSword},
		weaponDpsEP: 8.0, weights: w(map[stats.Stat]float64{stats.Strength: 1, stats.AttackPower: 0.5, stats.MeleeCrit: 10, stats.MeleeHit: 12,
			stats.SpellPower: 0.4, stats.Intellect: 0.2, stats.Stamina: 0.4, stats.Defense: 0.5, stats.BlockValue: 0.3, stats.Armor: 0.02}),
	},
	"warden_shaman": {
		class: proto.Class_ClassShaman, armor: mail, usesShield: true,
		weapons:     []proto.WeaponType{proto.WeaponType_WeaponTypeAxe, proto.WeaponType_WeaponTypeMace, proto.WeaponType_WeaponTypeDagger, proto.WeaponType_WeaponTypeFist},
		weaponDpsEP: 8.0, weights: w(map[stats.Stat]float64{stats.Strength: 1, stats.AttackPower: 0.5, stats.Agility: 0.6, stats.MeleeCrit: 10,
			stats.MeleeHit: 12, stats.SpellPower: 0.3, stats.Intellect: 0.2, stats.Stamina: 0.4, stats.Defense: 0.5, stats.Armor: 0.02}),
	},
	"tank_warrior": {
		class: proto.Class_ClassWarrior, armor: plate, usesShield: true,
		weapons:     []proto.WeaponType{proto.WeaponType_WeaponTypeAxe, proto.WeaponType_WeaponTypeMace, proto.WeaponType_WeaponTypeSword, proto.WeaponType_WeaponTypeDagger, proto.WeaponType_WeaponTypeFist},
		weaponDpsEP: 8.0, weights: w(map[stats.Stat]float64{stats.Strength: 1, stats.AttackPower: 0.5, stats.MeleeCrit: 10, stats.MeleeHit: 14,
			stats.Stamina: 0.4, stats.Defense: 0.5, stats.BlockValue: 0.3, stats.Armor: 0.02}),
	},
	// The two shaman specs have gear, but only Molten Core gear.
	"elemental_shaman": {
		class: proto.Class_ClassShaman, armor: mail,
		weapons:     []proto.WeaponType{proto.WeaponType_WeaponTypeAxe, proto.WeaponType_WeaponTypeMace, proto.WeaponType_WeaponTypeDagger, proto.WeaponType_WeaponTypeStaff},
		allowTwoHnd: true,
		weaponDpsEP: 0.0, weights: w(map[stats.Stat]float64{stats.SpellPower: 1, stats.SpellHit: 14, stats.SpellCrit: 10,
			stats.Intellect: 0.35, stats.MP5: 0.4, stats.Stamina: 0.05}),
	},
	"enhancement_shaman": {
		class: proto.Class_ClassShaman, armor: mail,
		weapons:     []proto.WeaponType{proto.WeaponType_WeaponTypeAxe, proto.WeaponType_WeaponTypeMace, proto.WeaponType_WeaponTypeDagger, proto.WeaponType_WeaponTypeFist},
		weaponDpsEP: 12.0, dualWield: true, weights: w(map[stats.Stat]float64{stats.AttackPower: 0.5, stats.Strength: 1, stats.Agility: 0.9, stats.MeleeCrit: 12,
			stats.MeleeHit: 14, stats.SpellPower: 0.2, stats.Intellect: 0.1, stats.Stamina: 0.1}),
	},
	// A holy caster. Spirit is worth more here than to the other casters, because
	// Forever's Spiritual Guidance turns it into spell damage as well as regen.
	"smite_priest": {
		class: proto.Class_ClassPriest, armor: cloth,
		weapons:     []proto.WeaponType{proto.WeaponType_WeaponTypeDagger, proto.WeaponType_WeaponTypeMace, proto.WeaponType_WeaponTypeStaff},
		allowTwoHnd: true,
		weaponDpsEP: 0.0, weights: w(map[stats.Stat]float64{stats.SpellPower: 1, stats.SpellHit: 14, stats.SpellCrit: 10,
			stats.Intellect: 0.35, stats.Spirit: 0.25, stats.MP5: 0.4, stats.Stamina: 0.05}),
	},
	"leather_placeholder": {class: proto.Class_ClassRogue, armor: leather},
}

// ep scores an item, adding weapon damage for specs that swing it.
func (s spec) ep(item *proto.UIItem) float64 {
	total := ep(item, s.weights)
	if s.weaponDpsEP > 0 && item.WeaponSpeed > 0 {
		total += (item.WeaponDamageMin + item.WeaponDamageMax) / 2 / item.WeaponSpeed * s.weaponDpsEP
	}
	return total
}
