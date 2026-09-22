import { ItemSlot, PseudoStat, Stat } from '@generated/proto/common';
import { UIEnchant as Enchant, UIItem as Item } from '@generated/proto/ui';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@i18n/localization', () => ({
	translateStat: (stat: unknown) => String(stat),
	translatePseudoStat: (pseudoStat: unknown) => String(pseudoStat),
}));

import { Stats } from '../proto/stats';
import { Player } from './player';

const weights = Stats.fromMap({ [Stat.StatStrength]: 2 }, { [PseudoStat.PseudoStatMainHandDps]: 10, [PseudoStat.PseudoStatOffHandDps]: 4 });

const computeEnchantEP = (enchant: Enchant, slot?: ItemSlot, weapon?: Item | null) =>
	Player.prototype.computeEnchantEP.call(
		{ enchantEPCache: new Map<number, number>(), computeStatsEP: (stats: Stats) => stats.computeEP(weights) } as unknown as Player<any>,
		enchant,
		slot,
		weapon,
	);

const striking = Enchant.create({ effectId: 1897, weaponDamage: 5 });
const sword = Item.create({ id: 1, weaponSpeed: 2.5 });

describe('Player.computeEnchantEP', () => {
	it('values flat weapon damage as the DPS it adds at the weapon speed of the slot it enchants', () => {
		expect(computeEnchantEP(striking, ItemSlot.ItemSlotMainHand, sword)).toBeCloseTo((5 / 2.5) * 10);
		expect(computeEnchantEP(striking, ItemSlot.ItemSlotOffHand, sword)).toBeCloseTo((5 / 2.5) * 4);
	});

	it('adds the weapon damage on top of the stats', () => {
		const enchant = Enchant.create({ effectId: 2, weaponDamage: 5, stats: new Stats().withStat(Stat.StatStrength, 3).asProtoArray() });
		expect(computeEnchantEP(enchant, ItemSlot.ItemSlotMainHand, sword)).toBeCloseTo(3 * 2 + (5 / 2.5) * 10);
	});

	it('gives weapon damage nothing without a weapon to read the speed off', () => {
		expect(computeEnchantEP(striking, ItemSlot.ItemSlotMainHand, null)).toBe(0);
		expect(computeEnchantEP(striking)).toBe(0);
	});
});
