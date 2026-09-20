import { Player } from '@generated/proto/api';
import { IndividualSimSettings } from '@generated/proto/ui';
import { describe, expect, it } from 'vitest';

import { migrateRetypedBuffFields, retypedBuffFields } from './buff_field_migration';

const v16Settings = () => ({
	apiVersion: 16,
	partyBuffs: { battleShout: 'TristateEffectImproved', manaSpringTotem: 'TristateEffectImproved' },
	debuffs: { faerieFire: 'TristateEffectMissing', misery: true },
	player: { buffs: { blessingOfMight: 2, blessingOfKings: true } },
});

describe('migrateRetypedBuffFields', () => {
	it('rewrites the retyped fields of a version-16 envelope so the parser accepts it', () => {
		const json = v16Settings();

		migrateRetypedBuffFields(json);
		const settings = IndividualSimSettings.fromJson(json as never);

		expect(settings.partyBuffs?.battleShout).toBe(true);
		expect(settings.debuffs?.faerieFire).toBe(false);
		expect(settings.player?.buffs?.blessingOfMight).toBe(true);
	});

	it('leaves the fields that are still a tristate, and every other field, alone', () => {
		const json = v16Settings();

		migrateRetypedBuffFields(json);

		expect(json.partyBuffs.manaSpringTotem).toBe('TristateEffectImproved');
		expect(json.debuffs.misery).toBe(true);
		expect(json.player.buffs.blessingOfKings).toBe(true);
	});

	it('leaves a payload stamped with the current version untouched', () => {
		const json = { apiVersion: 17, partyBuffs: { battleShout: 'TristateEffectImproved' } };

		migrateRetypedBuffFields(json);

		expect(json.partyBuffs.battleShout).toBe('TristateEffectImproved');
	});

	// SavedSettings and the raid messages carry no api_version, so they are migrated on sight.
	it('rewrites a message that carries no version, under either name for the individual buffs', () => {
		const savedSettings = { debuffs: { thunderClap: 1 }, playerBuffs: { blessingOfWisdom: 'TristateEffectRegular' } };
		const raid = { buffs: { thorns: 2 }, parties: [{ buffs: { devotionAura: 1 }, players: [{ buffs: { blessingOfMight: 0 } }] }] };

		migrateRetypedBuffFields(savedSettings);
		migrateRetypedBuffFields(raid);

		expect(savedSettings.debuffs.thunderClap).toBe(true);
		expect(savedSettings.playerBuffs.blessingOfWisdom).toBe(true);
		expect(raid.buffs.thorns).toBe(true);
		expect(raid.parties[0].buffs.devotionAura).toBe(true);
		expect(raid.parties[0].players[0].buffs.blessingOfMight).toBe(false);
	});

	it('rewrites a player handed over on its own', () => {
		const json = { buffs: { blessingOfWisdom: 'TristateEffectImproved' } };

		migrateRetypedBuffFields(json);

		expect(Player.fromJson(json as never).buffs?.blessingOfWisdom).toBe(true);
	});

	it('survives a blob that is not a message', () => {
		expect(() => migrateRetypedBuffFields(null)).not.toThrow();
		expect(() => migrateRetypedBuffFields('settings')).not.toThrow();
		expect(() => migrateRetypedBuffFields({ partyBuffs: 7 })).not.toThrow();
	});

	it('names 25 fields, the ones the proto retyped', () => {
		const fields = Object.values(retypedBuffFields).flat();

		expect(fields).toHaveLength(25);
		expect(new Set(fields).size).toBe(25);
	});
});
