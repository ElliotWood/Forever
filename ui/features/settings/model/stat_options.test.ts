import { Class, Stat } from '@generated/proto/common';
import type { Player } from '@sim/player/player';
import { UnitStat } from '@sim/proto/stats';
import { fakeHost } from '@sim/testing';
import { describe, expect, it } from 'vitest';

import * as BuffDebuffInputs from './buffs_debuffs';
import { applyOwnerClassLabels, inDisplayOrder, type PickerStatOptions, relevantStatOptions, type RenderableStatOptions } from './stat_options';

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
const commandingShout = ownedOption('Commanding Shout', Class.ClassWarrior);
const arcaneBrilliance = ownedOption('Arcane Brilliance', Class.ClassMage);
const giftOfArthas = ownedOption('Gift of Arthas');
const unlabelledWarriorRow = { config: {}, stats: [], ownerClass: Class.ClassWarrior } as unknown as RenderableStatOptions;

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

	it('marks every row the class owns, not just the first', () => {
		const shown = applyOwnerClassLabels([battleShout, arcaneBrilliance, commandingShout], playerOf(Class.ClassWarrior));

		expect(shown.map(option => option.config.label)).toEqual(['Battle Shout (External)', 'Arcane Brilliance', 'Commanding Shout (External)']);
		expect(shown[1]).toBe(arcaneBrilliance);
	});

	it('leaves a row with no label alone, since there is nothing to mark', () => {
		const shown = applyOwnerClassLabels([unlabelledWarriorRow], playerOf(Class.ClassWarrior));

		expect(shown[0]).toBe(unlabelledWarriorRow);
	});

	it('runs after relevantStatOptions, so a spec excluding a row by its config still drops it', () => {
		const shown = applyOwnerClassLabels(
			relevantStatOptions([battleShout, arcaneBrilliance], host({ exclude: [battleShout.config] })),
			playerOf(Class.ClassWarrior),
		);

		expect(shown.map(option => option.config.label)).toEqual(['Arcane Brilliance']);
	});
});

describe('inDisplayOrder', () => {
	const prebuiltBattleShout = ownedOption('Battle Shout', Class.ClassWarrior);
	const prebuiltSunderArmor = ownedOption('Sunder Armor', Class.ClassWarrior);
	const handWritten = ownedOption('Bloodlust');

	it('resolves a config to its prebuilt row and keeps the literal rows where they are written', () => {
		const composed = inDisplayOrder([prebuiltBattleShout, prebuiltSunderArmor], [prebuiltBattleShout.config, handWritten, prebuiltSunderArmor.config]);

		expect(composed).toEqual([prebuiltBattleShout, handWritten, prebuiltSunderArmor]);
		expect(composed[0]).toBe(prebuiltBattleShout);
	});

	it('reports a prebuilt row the display order never names, rather than dropping it off the tab', () => {
		expect(() => inDisplayOrder([prebuiltBattleShout, prebuiltSunderArmor], [prebuiltBattleShout.config])).toThrowError(
			'the display order leaves out Sunder Armor',
		);
	});

	it('reports a config with no prebuilt row', () => {
		expect(() => inDisplayOrder([prebuiltSunderArmor], [prebuiltBattleShout.config])).toThrowError(
			'the display order names Battle Shout, which no prebuilt row carries',
		);
	});
});

// The three registries interleave the generated rows with the hand-written ones. These sequences are
// the settings tab's layout, so a diff here is a UI change and wants a reason, not a re-snapshot.
describe('the buff registries', () => {
	const labelsOf = (options: ReadonlyArray<RenderableStatOptions>) => options.map(option => option.config.label);

	it('shows the party buffs in their settled order', () => {
		expect(labelsOf(BuffDebuffInputs.PARTY_BUFFS_CONFIG)).toEqual([
			'Blood Pact',
			'Commanding Shout',
			'Battle Shout',
			'Devotion Aura',
			'Ferocious Inspiratation',
			'Leader of the Pack',
			'Mana Spring',
			'Mana Tide Totem',
			'Vampiric Touch',
			'Moonkin Aura',
			'Retribution Aura',
			'Concentration Aura',
			'Sanctity Aura',
			'Totem of Wrath',
			'Trueshot Aura',
			'Wrath of Air Totem',
			'Unleashed Rage',
			'Atiesh - Mage',
			'Atiesh - Warlock',
			'Braided Eternium Chain',
			'Chain of the Twilight Owl',
			'Inspiring Presense - Caster',
			'Inspiring Presense - Melee',
			'Eye of the Night',
			'Jade Pendant of Blasting',
			'Strength of Earth',
			'Grace of Air',
			'Windfury Totem',
			'Drums',
			'Frost Resistance',
			'Nature Resistance',
			'Fire Resistance',
			'Frost Resistance Aura',
			'Fire Resistance Aura',
			'Shadow Resistance Aura',
			'Aspect of the Wild',
		]);
	});

	it('shows the raid and individual buffs interleaved in their settled order', () => {
		expect(labelsOf(BuffDebuffInputs.BUFFS_CONFIG)).toEqual([
			'Arcane Brilliance',
			'Blessing of Kings',
			'Bloodlust',
			'Divine Spirit',
			'Gift of the Wild',
			'Thorns',
			'Power Word: Fortitude',
			'Blessing of Might',
			'Blessing of Wisdom',
			'Blessing of Sanctuary',
			'Blessing of Salvation',
			'Shadow Protection',
			'Innervates',
			'Power Infusions',
		]);
	});

	it('shows the debuffs in their settled order', () => {
		expect(labelsOf(BuffDebuffInputs.DEBUFFS_CONFIG)).toEqual([
			'Blood Frenzy',
			"Hunter's Mark",
			'Improved Scorch',
			'Seal of the Crusader',
			'Judgement of Light',
			'Judgement of Wisdom',
			'Mangle',
			'Misery',
			'Shadow Weaving',
			'Curse of the Elements',
			'Curse of Recklessness',
			'Faerie Fire',
			'Expose Armor',
			'Sunder Armor',
			"Winter's Chill",
			'Gift of Arthas',
			'Demoralizing Roar',
			'Demoralizing Shout',
			'Screech',
			'Thunder Clap',
			'Insect Swarm',
			'Scorpid Sting',
			'Shadow Embrace',
		]);
		expect(BuffDebuffInputs.DEBUFFS_MISC_CONFIG).toHaveLength(0);
	});

	it('puts the hand-written inputs themselves at those positions, not lookalikes', () => {
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
