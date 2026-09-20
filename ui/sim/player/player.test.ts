import { Debuffs } from '@generated/proto/buffs';
import { PseudoStat, Stat } from '@generated/proto/common';
import { describe, expect, it } from 'vitest';

import { Stats } from '../proto/stats';
import { Player } from './player';

// getDebuffStats reads nothing but the raid's debuffs, so the method is called against that one
// dependency rather than against a constructed Player, which needs a whole sim behind it.
const debuffStats = (debuffs: Partial<Debuffs>): Stats =>
	Player.prototype.getDebuffStats.call({
		sim: { raid: { getDebuffs: () => Debuffs.create(debuffs) } },
	} as unknown as Player<any>);

describe('Player.getDebuffStats', () => {
	it('credits Seal of the Crusader with 3% crit on every school', () => {
		const stats = debuffStats({ improvedSealOfTheCrusader: true });

		expect(stats.getPseudoStat(PseudoStat.PseudoStatMeleeCritPercent)).toBe(3);
		expect(stats.getPseudoStat(PseudoStat.PseudoStatRangedCritPercent)).toBe(3);
		expect(stats.getPseudoStat(PseudoStat.PseudoStatSpellCritPercent)).toBe(3);
	});

	it('credits nothing when the debuff is off', () => {
		const stats = debuffStats({ improvedSealOfTheCrusader: false, faerieFire: true });

		expect(stats.getPseudoStat(PseudoStat.PseudoStatMeleeCritPercent)).toBe(0);
		expect(stats.getPseudoStat(PseudoStat.PseudoStatRangedCritPercent)).toBe(0);
		expect(stats.getPseudoStat(PseudoStat.PseudoStatSpellCritPercent)).toBe(0);
	});

	it('credits Hunters Mark with its ranged attack power', () => {
		expect(debuffStats({ huntersMark: true }).getStat(Stat.StatRangedAttackPower)).toBe(440);
		expect(debuffStats({ huntersMark: false }).getStat(Stat.StatRangedAttackPower)).toBe(0);
	});
});
