import { Class } from '@generated/proto/common';
import { describe, expect, it } from 'vitest';

import * as BuffDebuffInputs from './buffs_debuffs';
import type { RenderableStatOptions } from './stat_options';

// The three registries interleave the generated rows with the hand-written ones. These sequences are
// the settings tab's layout, so a diff here is a UI change and wants a reason, not a re-snapshot.
describe('the buff registries', () => {
	const labelsOf = (options: ReadonlyArray<RenderableStatOptions>) => options.map(option => option.config.label);

	it('shows the party buffs in their settled order', () => {
		expect(labelsOf(BuffDebuffInputs.PARTY_BUFFS_CONFIG)).toEqual([
			'Blood Pact',
			'Battle Shout',
			'Enhanced Battle Shout',
			'Devotion Aura',
			'Leader of the Pack',
			'Mana Spring',
			'Mana Tide Totem',
			'Moonkin Aura',
			'Retribution Aura',
			'Concentration Aura',
			'Trueshot Aura',
			'Atiesh - Mage',
			'Atiesh - Warlock',
			'Strength of Earth',
			'Grace of Air',
			'Windfury Totem',
		]);
	});

	it('shows the raid and individual buffs interleaved in their settled order', () => {
		expect(labelsOf(BuffDebuffInputs.BUFFS_CONFIG)).toEqual([
			'Arcane Brilliance',
			'Greater Blessing of Kings',
			'Prayer of Spirit',
			'Gift of the Wild',
			'Thorns',
			'Prayer of Fortitude',
			'Greater Blessing of Might',
			'Greater Blessing of Wisdom',
			'Greater Blessing of Salvation',
			'Prayer of Shadow Protection',
			'Fire Resistance Aura',
			'Frost Resistance Aura',
			'Shadow Resistance Aura',
			'Fire Resistance',
			'Frost Resistance',
			'Nature Resistance',
			'Aspect of the Wild',
			'Innervates',
			'Power Infusions',
		]);
	});

	it('shows the debuffs in their settled order', () => {
		expect(labelsOf(BuffDebuffInputs.DEBUFFS_CONFIG)).toEqual([
			"Hunter's Mark",
			'Seal of the Crusader',
			'Judgement of Light',
			'Judgement of Wisdom',
			'Mangle',
			'Curse of the Elements',
			'Curse of Recklessness',
			'Faerie Fire',
			'Expose Armor',
			'Sunder Armor',
			'Gift of Arthas',
			'Demoralizing Roar',
			'Demoralizing Shout',
			'Thunder Clap',
			'Insect Swarm',
			'Scorpid Sting',
		]);
	});

	it('has no miscellaneous debuff rows left', () => {
		expect(BuffDebuffInputs.DEBUFFS_MISC_CONFIG).toHaveLength(0);
	});

	it('puts the hand-written inputs themselves at those positions, not lookalikes', () => {
		expect(BuffDebuffInputs.PARTY_BUFFS_CONFIG[2].config).toBe(BuffDebuffInputs.EnhancedBattleShout);
		expect(BuffDebuffInputs.BUFFS_CONFIG[8].config).toBe(BuffDebuffInputs.GreaterBlessingOfSalvation);
	});

	it('gives the generated rows their owner class, so the settings tab can mark them external', () => {
		expect(BuffDebuffInputs.PARTY_BUFFS_CONFIG[1].config).toBe(BuffDebuffInputs.BattleShout);
		expect(BuffDebuffInputs.PARTY_BUFFS_CONFIG[1].ownerClass).toBe(Class.ClassWarrior);
		expect(BuffDebuffInputs.BUFFS_CONFIG[8].ownerClass).toBeUndefined();
	});
});
