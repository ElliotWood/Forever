import { Class, Stat } from '@generated/proto/common';
import type { Player } from '@sim/player/player';
import { UnitStat } from '@sim/proto/stats';
import { fakeHost } from '@sim/testing';
import { describe, expect, it } from 'vitest';

import * as BuffDebuffInputs from './buffs_debuffs';
import { applyOwnerClassLabels, type PickerStatOptions, relevantStatOptions, type RenderableStatOptions } from './stat_options';

const option = (stats: Array<Stat>) => ({ config: { label: stats.join('/') }, stats }) as unknown as PickerStatOptions;

const manaSpring = option([Stat.StatMP5]);
const manaTide = option([Stat.StatMP5]);
const bloodlust = option([]);
const strengthOfEarth = option([Stat.StatStrength]);
const fortitude = option([Stat.StatStamina]);
const options = [manaSpring, manaTide, bloodlust, strengthOfEarth, fortitude];

const host = (parts: { epStats?: Array<Stat>; displayStats?: Array<UnitStat>; include?: Array<unknown>; exclude?: Array<unknown> }) =>
	fakeHost({
		individualConfig: {
			epStats: parts.epStats ?? [],
			displayStats: parts.displayStats ?? [],
			includeBuffDebuffInputs: parts.include ?? [],
			excludeBuffDebuffInputs: parts.exclude ?? [],
		},
	});

describe('relevantStatOptions', () => {
	it('keeps an option tagged with an EP stat or a displayed stat, and every untagged option', () => {
		const shown = relevantStatOptions(options, host({ epStats: [Stat.StatMP5], displayStats: [UnitStat.fromStat(Stat.StatStamina)] }));
		expect(shown).toEqual([manaSpring, manaTide, bloodlust, fortitude]);
	});

	it('includes by stat, the way every TBC spec lists them', () => {
		expect(relevantStatOptions(options, host({ include: [Stat.StatStrength] }))).toEqual([bloodlust, strengthOfEarth]);
	});

	it('excludes by stat, the sentinel way the feral specs drop Windfury', () => {
		expect(relevantStatOptions(options, host({ epStats: [Stat.StatMP5], exclude: [Stat.StatMP5] }))).toEqual([bloodlust]);
	});

	it('includes and excludes a single input by its config, so one MP5 buff can go while the other stays', () => {
		const shown = relevantStatOptions(
			options,
			host({ epStats: [Stat.StatMP5], include: [fortitude.config], exclude: [manaTide.config, bloodlust.config] }),
		);
		expect(shown).toEqual([manaSpring, fortitude]);
	});
});

const ownedOption = (label: string, ownerClass?: Class) => ({ config: { label }, stats: [], ownerClass }) as unknown as RenderableStatOptions;
const playerOf = (playerClass: Class) => ({ getClass: () => playerClass }) as unknown as Player<any>;

const battleShout = ownedOption('Battle Shout', Class.ClassWarrior);
const arcaneBrilliance = ownedOption('Arcane Brilliance', Class.ClassMage);
const giftOfArthas = ownedOption('Gift of Arthas');

describe('applyOwnerClassLabels', () => {
	it('marks the buffs the player casts itself as external, and returns every other row untouched', () => {
		const shown = applyOwnerClassLabels([battleShout, arcaneBrilliance, giftOfArthas], playerOf(Class.ClassWarrior));

		expect(shown.map(option => option.config.label)).toEqual(['Battle Shout (External)', 'Arcane Brilliance', 'Gift of Arthas']);
		expect(shown[1]).toBe(arcaneBrilliance);
		expect(shown[2]).toBe(giftOfArthas);
	});

	it('copies the relabelled row instead of renaming the shared config', () => {
		const shown = applyOwnerClassLabels([battleShout], playerOf(Class.ClassWarrior));

		expect(shown[0]).not.toBe(battleShout);
		expect(shown[0].config).not.toBe(battleShout.config);
		expect(battleShout.config.label).toBe('Battle Shout');
	});

	it('runs after relevantStatOptions, so a spec excluding a row by its config still drops it', () => {
		const shown = applyOwnerClassLabels(
			relevantStatOptions([battleShout, arcaneBrilliance], host({ exclude: [battleShout.config] })),
			playerOf(Class.ClassWarrior),
		);

		expect(shown.map(option => option.config.label)).toEqual(['Arcane Brilliance']);
	});
});

describe('the buff registries', () => {
	it('carries every generated row and every hand-written one, hand-written rows at their old positions', () => {
		expect(BuffDebuffInputs.PARTY_BUFFS_CONFIG).toHaveLength(36);
		expect(BuffDebuffInputs.BUFFS_CONFIG).toHaveLength(14);
		expect(BuffDebuffInputs.DEBUFFS_CONFIG).toHaveLength(23);
		expect(BuffDebuffInputs.DEBUFFS_MISC_CONFIG).toHaveLength(0);

		expect(BuffDebuffInputs.PARTY_BUFFS_CONFIG[4].config).toBe(BuffDebuffInputs.FerociousInspiration);
		expect(BuffDebuffInputs.PARTY_BUFFS_CONFIG[28].config).toBe(BuffDebuffInputs.DrumsBuff);
		expect(BuffDebuffInputs.BUFFS_CONFIG[2].config).toBe(BuffDebuffInputs.Bloodlust);
		expect(BuffDebuffInputs.DEBUFFS_CONFIG[0].config).toBe(BuffDebuffInputs.BloodFrenzy);
		expect(BuffDebuffInputs.DEBUFFS_CONFIG[22].config).toBe(BuffDebuffInputs.ShadowEmbrace);
	});

	it('gives the generated rows their owner class, so the settings tab can mark them external', () => {
		expect(BuffDebuffInputs.PARTY_BUFFS_CONFIG[2].config).toBe(BuffDebuffInputs.BattleShout);
		expect(BuffDebuffInputs.PARTY_BUFFS_CONFIG[2].ownerClass).toBe(Class.ClassWarrior);
		expect(BuffDebuffInputs.BUFFS_CONFIG[2].ownerClass).toBeUndefined();
	});
});
