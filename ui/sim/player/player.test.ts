import { Player as PlayerProto } from '@generated/proto/api';
import { APLRotation_Type } from '@generated/proto/apl';
import { Debuffs } from '@generated/proto/buffs';
import { PseudoStat, Stat } from '@generated/proto/common';
import { DpsWarrior_Rotation, DpsWarriorSpec, WarriorSunder } from '@generated/proto/warrior';
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
	// The crit the target loses comes from Improved Seal of the Crusader, which has no node in
	// any Forever trait tree, so sim/core's aura applies none of it and the sheet credits none.
	it('credits Seal of the Crusader with no crit', () => {
		const stats = debuffStats({ improvedSealOfTheCrusader: true });

		expect(stats.getPseudoStat(PseudoStat.PseudoStatMeleeCritPercent)).toBe(0);
		expect(stats.getPseudoStat(PseudoStat.PseudoStatRangedCritPercent)).toBe(0);
		expect(stats.getPseudoStat(PseudoStat.PseudoStatSpellCritPercent)).toBe(0);
	});

	it('credits nothing for a debuff that moves no stat on the sheet', () => {
		const stats = debuffStats({ faerieFire: true });

		expect(stats.getPseudoStat(PseudoStat.PseudoStatMeleeCritPercent)).toBe(0);
		expect(stats.getPseudoStat(PseudoStat.PseudoStatRangedCritPercent)).toBe(0);
		expect(stats.getPseudoStat(PseudoStat.PseudoStatSpellCritPercent)).toBe(0);
	});

	// The number has to be the one HuntersMarkValue states in sim/core/debuffs_auto_gen.go,
	// or the sheet and the simulation disagree about the same debuff.
	it("credits Hunter's Mark with the ranged attack power spell 14325 states", () => {
		expect(debuffStats({ huntersMark: true }).getStat(Stat.StatRangedAttackPower)).toBe(71);
		expect(debuffStats({ huntersMark: false }).getStat(Stat.StatRangedAttackPower)).toBe(0);
	});
});

// The simple rotation travels as an opaque JSON string, and `DpsWarrior_Rotation.fromJson` reads it
// without ignoreUnknownFields, so a reserved key left in a version-16 save costs the whole rotation.
describe('Player.updateProtoVersion', () => {
	const v16Warrior = (specRotationJson: string) =>
		PlayerProto.create({
			apiVersion: 16,
			spec: { oneofKind: 'dpsWarrior', dpsWarrior: {} },
			rotation: { type: APLRotation_Type.TypeSimple, simple: { specRotationJson } },
		});

	it('drops the warrior bloodlust timing a version-16 simple rotation still carries', () => {
		const proto = v16Warrior('{"spec":"DpsWarriorSpecArms","sunderArmor":"WarriorSunderMaintain","bloodlustTiming":5}');

		Player.updateProtoVersion(proto);
		const rotation = DpsWarrior_Rotation.fromJson(JSON.parse(proto.rotation!.simple!.specRotationJson));

		expect(rotation.spec).toBe(DpsWarriorSpec.DpsWarriorSpecArms);
		expect(rotation.sunderArmor).toBe(WarriorSunder.WarriorSunderMaintain);
	});

	it('drops the field under the name a proto3 JSON payload spells it', () => {
		const proto = v16Warrior('{"spec":"DpsWarriorSpecFury","bloodlust_timing":5}');

		Player.updateProtoVersion(proto);

		expect(() => DpsWarrior_Rotation.fromJson(JSON.parse(proto.rotation!.simple!.specRotationJson))).not.toThrow();
	});

	it('leaves a simple rotation it cannot parse alone', () => {
		const proto = v16Warrior('not json');

		Player.updateProtoVersion(proto);

		expect(proto.rotation!.simple!.specRotationJson).toBe('not json');
	});
});
