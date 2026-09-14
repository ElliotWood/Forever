import { ArmorType, HandType, ItemRandomSuffix, ItemSlot, ItemType, PseudoStat, Spec, Stat, WeaponType } from '../core/proto/common';
import { UIItem as Item } from '../core/proto/ui';
import { getWeaponDPS } from '../core/proto_utils/equipped_item';
import { Stats } from '../core/proto_utils/stats';
import { classToEligibleRangedWeaponTypes, classToEligibleWeaponTypes, classToMaxArmorType, isTankSpec, specToClass } from '../core/proto_utils/utils';

// How many items each slot keeps, so a reader can see what the pick beat.
const RANKED_PER_SLOT = 5;

export interface RankedItem {
	item: Item;
	ep: number;
}

export interface SlotRanking {
	slot: ItemSlot;
	// Best first, empty when the launch pool holds nothing the spec can wear here.
	ranked: Array<RankedItem>;
	// Set when the slot needs explaining, e.g. an off hand given up to a two hander.
	note?: string;
}

// The armour slots: the only ones where an item's armour type decides who may wear it.
const armorSlotTypes = [
	ItemType.ItemTypeHead,
	ItemType.ItemTypeShoulder,
	ItemType.ItemTypeChest,
	ItemType.ItemTypeWrist,
	ItemType.ItemTypeHands,
	ItemType.ItemTypeWaist,
	ItemType.ItemTypeLegs,
	ItemType.ItemTypeFeet,
];

const offHandOnlyWeaponTypes = [WeaponType.WeaponTypeShield, WeaponType.WeaponTypeOffHand];

const slotItemTypes: Partial<Record<ItemSlot, ItemType>> = {
	[ItemSlot.ItemSlotHead]: ItemType.ItemTypeHead,
	[ItemSlot.ItemSlotNeck]: ItemType.ItemTypeNeck,
	[ItemSlot.ItemSlotShoulder]: ItemType.ItemTypeShoulder,
	[ItemSlot.ItemSlotBack]: ItemType.ItemTypeBack,
	[ItemSlot.ItemSlotChest]: ItemType.ItemTypeChest,
	[ItemSlot.ItemSlotWrist]: ItemType.ItemTypeWrist,
	[ItemSlot.ItemSlotHands]: ItemType.ItemTypeHands,
	[ItemSlot.ItemSlotWaist]: ItemType.ItemTypeWaist,
	[ItemSlot.ItemSlotLegs]: ItemType.ItemTypeLegs,
	[ItemSlot.ItemSlotFeet]: ItemType.ItemTypeFeet,
};

// Follows canEquipItem() in core/proto_utils/utils, with two deliberate differences.
//
// An armour slot takes only the class's own armour type rather than anything up to it. A
// warrior can physically wear cloth, but no warrior raids in cloth, and EP weights that
// say nothing about armour would happily rank a cloth piece over a plate one.
// tools/launch_gear draws the same line for the same reason.
//
// Dual wield is asked of the spec rather than the class. canDualWield() is class-based
// and so excludes Enhancement, which learns dual wield from a talent; asking the spec
// gets Enhancement its off hand without lying to any other spec about theirs.
export function canSpecUseItem(spec: Spec, item: Item, slot: ItemSlot, dualWields: boolean): boolean {
	const playerClass = specToClass[spec];
	if (item.classAllowlist.length > 0 && !item.classAllowlist.includes(playerClass)) {
		return false;
	}

	if (item.type === ItemType.ItemTypeWeapon) {
		const eligible = classToEligibleWeaponTypes[playerClass].find(wt => wt.weaponType === item.weaponType);
		if (!eligible) {
			return false;
		}
		if (slot === ItemSlot.ItemSlotMainHand) {
			if (item.handType === HandType.HandTypeOffHand || offHandOnlyWeaponTypes.includes(item.weaponType)) {
				return false;
			}
			return item.handType !== HandType.HandTypeTwoHand || !!eligible.canUseTwoHand;
		}
		if (slot === ItemSlot.ItemSlotOffHand) {
			if (item.handType === HandType.HandTypeTwoHand || item.handType === HandType.HandTypeMainHand) {
				return false;
			}
			return dualWields || offHandOnlyWeaponTypes.includes(item.weaponType);
		}
		return false;
	}

	if (item.type === ItemType.ItemTypeRanged) {
		return slot === ItemSlot.ItemSlotRanged && classToEligibleRangedWeaponTypes[playerClass].includes(item.rangedWeaponType);
	}

	if (armorSlotTypes.includes(item.type)) {
		return item.armorType === classToMaxArmorType[playerClass] || item.armorType === ArmorType.ArmorTypeUnknown;
	}

	return true;
}

// core's Stats sizes its pseudo stat array by the number of PseudoStat enum entries, so a
// pseudo stat numbered past that count reads back undefined rather than zero.
// PseudoStatCastSpeedMultiplier is 34 of 32 entries and does that. No launched spec weighs
// it, but arithmetic on a weight still has to survive being handed one.
const weightFor = (epWeights: Stats, pseudoStat: PseudoStat): number => epWeights.getPseudoStat(pseudoStat) || 0;

// Mirrors Player.computeItemEP, so a pick here is the item the gear picker's EP column
// would put at the top of that slot. Random suffixes are scored at their best roll, as
// the picker does, because the picker is ranking the item rather than a given drop.
export function computeItemEP(item: Item, epWeights: Stats, slot: ItemSlot, randomSuffix: (id: number) => ItemRandomSuffix | undefined): number {
	let itemStats = new Stats(item.stats);

	if (item.weaponSpeed > 0) {
		const weaponDps = getWeaponDPS(item);
		if (slot === ItemSlot.ItemSlotMainHand) {
			itemStats = itemStats.withPseudoStat(PseudoStat.PseudoStatMainHandDps, weaponDps);
		} else if (slot === ItemSlot.ItemSlotOffHand) {
			itemStats = itemStats.withPseudoStat(PseudoStat.PseudoStatOffHandDps, weaponDps);
		} else if (slot === ItemSlot.ItemSlotRanged) {
			itemStats = itemStats.withPseudoStat(PseudoStat.PseudoStatRangedDps, weaponDps);
		}
	}
	itemStats = itemStats.addPseudoStat(PseudoStat.BonusPhysicalDamage, item.bonusPhysicalDamage);

	const suffixEPs = item.randomSuffixOptions.map(id => {
		const suffix = randomSuffix(id);
		return suffix ? new Stats(suffix.stats).computeEP(epWeights) : 0;
	});

	let ep = itemStats.computeEP(epWeights) + (suffixEPs.length ? Math.max(...suffixEPs) : 0);

	// Unique items are slightly worse than non-unique because you can have only one.
	if (item.unique) {
		ep -= 0.01;
	}
	if (item.stats[Stat.StatMeleeHaste] > 0) {
		ep += weightFor(epWeights, PseudoStat.PseudoStatMeleeSpeedMultiplier) * item.stats[Stat.StatMeleeHaste];
		ep += weightFor(epWeights, PseudoStat.PseudoStatRangedSpeedMultiplier) * item.stats[Stat.StatMeleeHaste];
	}
	if (item.stats[Stat.StatSpellHaste] > 0) {
		ep += weightFor(epWeights, PseudoStat.PseudoStatCastSpeedMultiplier) * item.stats[Stat.StatSpellHaste];
	}

	return ep;
}

// A spec dual wields if its own EP weights pay for an off hand's weapon damage. The weights
// are the sim's answer to "what is this spec holding", so they are asked rather than told.
export function specDualWields(epWeights: Stats): boolean {
	return weightFor(epWeights, PseudoStat.PseudoStatOffHandDps) > 0;
}

// Every slot a character fills, in the order the gear tab lists them.
export const bisSlots: Array<ItemSlot> = [
	ItemSlot.ItemSlotHead,
	ItemSlot.ItemSlotNeck,
	ItemSlot.ItemSlotShoulder,
	ItemSlot.ItemSlotBack,
	ItemSlot.ItemSlotChest,
	ItemSlot.ItemSlotWrist,
	ItemSlot.ItemSlotHands,
	ItemSlot.ItemSlotWaist,
	ItemSlot.ItemSlotLegs,
	ItemSlot.ItemSlotFeet,
	ItemSlot.ItemSlotFinger1,
	ItemSlot.ItemSlotFinger2,
	ItemSlot.ItemSlotTrinket1,
	ItemSlot.ItemSlotTrinket2,
	ItemSlot.ItemSlotMainHand,
	ItemSlot.ItemSlotOffHand,
	ItemSlot.ItemSlotRanged,
];

export interface SpecGear {
	slots: Array<SlotRanking>;
	// How many of the database's items this spec could wear anywhere, for the page to show
	// that a ranking came out of a real pool rather than a handful of leftovers.
	poolSize: number;
}

// Picks the highest-EP item for every slot out of the items the spec can wear.
//
// The two paired slots take the two best items of different names, because the rings and
// trinkets worth wearing are unique-equipped. The hands are decided together: a two hander
// is worth taking only if it beats a one hander and whatever the off hand would have held.
export function rankGear(spec: Spec, epWeights: Stats, items: Array<Item>, randomSuffix: (id: number) => ItemRandomSuffix | undefined): SpecGear {
	// A tank holds a shield, so its off hand is never a second weapon.
	const dualWields = specDualWields(epWeights) && !isTankSpec(spec);

	const rankFor = (slot: ItemSlot, pool: Array<Item>): Array<RankedItem> =>
		pool
			.filter(item => canSpecUseItem(spec, item, slot, dualWields))
			.map(item => ({ item, ep: computeItemEP(item, epWeights, slot, randomSuffix) }))
			.sort((a, b) => b.ep - a.ep);

	// Two rings or trinkets of the same name cannot both be worn.
	const distinctNames = (ranked: Array<RankedItem>): Array<RankedItem> => {
		const seen = new Set<string>();
		return ranked.filter(({ item }) => {
			if (seen.has(item.name)) {
				return false;
			}
			seen.add(item.name);
			return true;
		});
	};

	const slots: Partial<Record<ItemSlot, SlotRanking>> = {};

	for (const [slot, itemType] of Object.entries(slotItemTypes)) {
		const itemSlot = Number(slot) as ItemSlot;
		slots[itemSlot] = {
			slot: itemSlot,
			ranked: rankFor(
				itemSlot,
				items.filter(item => item.type === itemType),
			),
		};
	}

	for (const [first, second, itemType] of [
		[ItemSlot.ItemSlotFinger1, ItemSlot.ItemSlotFinger2, ItemType.ItemTypeFinger],
		[ItemSlot.ItemSlotTrinket1, ItemSlot.ItemSlotTrinket2, ItemType.ItemTypeTrinket],
	] as Array<[ItemSlot, ItemSlot, ItemType]>) {
		const ranked = distinctNames(
			rankFor(
				first,
				items.filter(item => item.type === itemType),
			),
		);
		slots[first] = { slot: first, ranked };
		slots[second] = { slot: second, ranked: ranked.slice(1) };
	}

	slots[ItemSlot.ItemSlotRanged] = {
		slot: ItemSlot.ItemSlotRanged,
		ranked: rankFor(
			ItemSlot.ItemSlotRanged,
			items.filter(item => item.type === ItemType.ItemTypeRanged),
		),
	};

	const weapons = items.filter(item => item.type === ItemType.ItemTypeWeapon);
	const mainHand = rankFor(ItemSlot.ItemSlotMainHand, weapons);
	const twoHanders = mainHand.filter(({ item }) => item.handType === HandType.HandTypeTwoHand);
	const oneHanders = mainHand.filter(({ item }) => item.handType !== HandType.HandTypeTwoHand);
	const offHand = rankFor(ItemSlot.ItemSlotOffHand, weapons).filter(({ item }) => item.name !== oneHanders[0]?.item.name);

	const twoHandEP = twoHanders[0]?.ep ?? -Infinity;
	const oneHandEP = oneHanders.length ? oneHanders[0].ep + (offHand[0]?.ep ?? 0) : -Infinity;
	if (twoHanders.length && twoHandEP > oneHandEP) {
		slots[ItemSlot.ItemSlotMainHand] = { slot: ItemSlot.ItemSlotMainHand, ranked: twoHanders };
		slots[ItemSlot.ItemSlotOffHand] = {
			slot: ItemSlot.ItemSlotOffHand,
			ranked: [],
			note: oneHanders.length
				? `Given up to a two hander, worth ${(twoHandEP - oneHandEP).toFixed(1)} EP more than the best one hander and off hand together.`
				: 'Given up to a two hander: this spec has no one handed weapon it can use.',
		};
	} else {
		slots[ItemSlot.ItemSlotMainHand] = { slot: ItemSlot.ItemSlotMainHand, ranked: oneHanders };
		slots[ItemSlot.ItemSlotOffHand] = { slot: ItemSlot.ItemSlotOffHand, ranked: offHand };
	}

	return {
		slots: bisSlots.map(slot => ({ ...slots[slot]!, ranked: slots[slot]!.ranked.slice(0, RANKED_PER_SLOT) })),
		poolSize: items.filter(item => bisSlots.some(slot => canSpecUseItem(spec, item, slot, dualWields))).length,
	};
}
