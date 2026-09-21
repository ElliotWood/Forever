import { Player } from '@generated/proto/api';
import { Debuffs, IndividualBuffs, PartyBuffs, RaidBuffs } from '@generated/proto/buffs';
import { IndividualSimSettings } from '@generated/proto/ui';
import { ScalarType, UnknownFieldHandler } from '@protobuf-ts/runtime';
import { describe, expect, it } from 'vitest';

import { migrateRetypedBuffFields, retiredBuffFields, retiredFieldSpellings, retypedBuffFields } from './buff_field_migration';

const v16Settings = () => ({
	apiVersion: 16,
	partyBuffs: { battleShout: 'TristateEffectImproved', manaSpringTotem: 'TristateEffectImproved' },
	debuffs: { faerieFire: 'TristateEffectMissing', misery: true, jocRetribution2pt4: true },
	player: { buffs: { blessingOfMight: 2, blessingOfKings: true, unleashedRage: true } },
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
		const raidSettings = { raid: { buffs: { thorns: 2 }, parties: [{ buffs: { devotionAura: 1 }, players: [{ buffs: { blessingOfMight: 0 } }] }] } };

		migrateRetypedBuffFields(savedSettings);
		migrateRetypedBuffFields(raidSettings);

		expect(savedSettings.debuffs.thunderClap).toBe(true);
		expect(savedSettings.playerBuffs.blessingOfWisdom).toBe(true);
		expect(raidSettings.raid.buffs.thorns).toBe(true);
		expect(raidSettings.raid.parties[0].buffs.devotionAura).toBe(true);
		expect(raidSettings.raid.parties[0].players[0].buffs.blessingOfMight).toBe(false);
	});

	it('rewrites a raid, a party or a player handed over on its own', () => {
		const raid = { buffs: { thorns: 2 }, debuffs: { thunderClap: 1 } };
		const party = { buffs: { devotionAura: 1 } };
		const player = { buffs: { blessingOfWisdom: 'TristateEffectImproved' } };

		migrateRetypedBuffFields(raid, 'raid');
		migrateRetypedBuffFields(party, 'party');
		migrateRetypedBuffFields(player, 'player');

		expect(raid.buffs.thorns).toBe(true);
		expect(raid.debuffs.thunderClap).toBe(true);
		expect(party.buffs.devotionAura).toBe(true);
		expect(Player.fromJson(player as never).buffs?.blessingOfWisdom).toBe(true);
	});

	// `toJson` omits an empty repeated field, so a raid with no parties and a party with no players
	// reach the migration without the key that used to say which message they are.
	it('rewrites a raid with no parties and a party with no players', () => {
		const raid = { buffs: { giftOfTheWild: 'TristateEffectImproved' } };
		const party = { buffs: { battleShout: 'TristateEffectImproved' } };

		migrateRetypedBuffFields(raid, 'raid');
		migrateRetypedBuffFields(party, 'party');

		expect(raid.buffs.giftOfTheWild).toBe(true);
		expect(party.buffs.battleShout).toBe(true);
	});

	it('survives a blob that is not a message', () => {
		expect(() => migrateRetypedBuffFields(null)).not.toThrow();
		expect(() => migrateRetypedBuffFields('settings')).not.toThrow();
		expect(() => migrateRetypedBuffFields({ partyBuffs: 7 })).not.toThrow();
	});

	it('drops the retired fields of a version-16 envelope so the parser accepts it', () => {
		const json = {
			apiVersion: 16,
			partyBuffs: { drums: 'LesserDrumsOfBattle', snapshotBsT2: true, wrath_of_air_totem: 'TristateEffectImproved', battleShout: true },
			debuffs: { jocRetribution2pt4: true, misery: true },
			player: { buffs: { unleashedRage: true, blessingOfKings: true } },
		} as Record<string, any>;

		migrateRetypedBuffFields(json);

		expect(Object.keys(json.partyBuffs)).toEqual(['snapshotBsT2', 'battleShout']);
		expect(Object.keys(json.debuffs)).toEqual([]);
		expect(Object.keys(json.player.buffs)).toEqual(['blessingOfKings']);
		expect(() => IndividualSimSettings.fromJson(json as never)).not.toThrow();
	});

	// has_bs_solarian_sapphire was a warrior class option, not a buff field, so the pre-pass has
	// nothing to do for it: the loaders' `ignoreUnknownFields` skips it wherever it is nested.
	it('loads a version-16 payload that still carries the warrior class option the proto reserved', () => {
		const json = {
			apiVersion: 16,
			player: {
				dpsWarrior: { options: { classOptions: { hasBsSolarianSapphire: true, hasBsT2: true, stanceSnapshot: true } } },
			},
		} as Record<string, any>;

		migrateRetypedBuffFields(json);
		const settings = IndividualSimSettings.fromJson(json as never, { ignoreUnknownFields: true });

		const spec = settings.player?.spec;
		expect(spec?.oneofKind).toBe('dpsWarrior');
		expect(spec?.oneofKind === 'dpsWarrior' && spec.dpsWarrior.options?.classOptions?.hasBsT2).toBe(true);
		expect(spec?.oneofKind === 'dpsWarrior' && spec.dpsWarrior.options?.classOptions?.stanceSnapshot).toBe(true);
	});

	it('spells a retired field the three ways a payload may carry it', () => {
		expect(retiredFieldSpellings('joc_retribution_2pt4')).toEqual(['joc_retribution_2pt4', 'jocRetribution2pt4', 'jocRetribution2Pt4']);
		expect(retiredFieldSpellings('soe_enhancement_2pt4')).toEqual(['soe_enhancement_2pt4', 'soeEnhancement2pt4', 'soeEnhancement2Pt4']);
		expect(retiredFieldSpellings('totem_of_wrath')).toEqual(['totem_of_wrath', 'totemOfWrath']);
		expect(retiredFieldSpellings('drums')).toEqual(['drums']);
	});

	it('names 33 retired fields, the ones the proto reserved', () => {
		const fields = Object.values(retiredBuffFields).flat();

		expect(fields).toHaveLength(33);
		expect(new Set(fields).size).toBe(33);
	});

	it('names 25 fields, the ones the proto retyped', () => {
		const fields = Object.values(retypedBuffFields).flat();

		expect(fields).toHaveLength(25);
		expect(new Set(fields).size).toBe(25);
	});

	// A share link is binary, and nothing rewrites it: a tristate's varint 2 decodes as the bool's
	// true, and a retired field number is skipped as an unknown one. The bytes are written by hand
	// because no message in the tree encodes either shape: `[0x30, 0x02]` is RaidBuffs field 6
	// (thorns) carrying 2, and `[0x08, 0x01]` is the retired PartyBuffs field 1 carrying 1.
	it('reads a saved link with no pre-pass', () => {
		const raidBuffs = RaidBuffs.fromBinary(new Uint8Array([0x30, 0x02]));
		const partyBuffs = PartyBuffs.fromBinary(new Uint8Array([0x08, 0x01]));

		expect(raidBuffs.thorns).toBe(true);
		expect(UnknownFieldHandler.list(partyBuffs).map(field => field.no)).toEqual([1]);
		expect(PartyBuffs.toJson(partyBuffs)).toEqual({});
	});

	// A row that goes back to ProtoTristate in the manifest would have the rewrite write a bool into
	// an enum, which fromJson throws on. The list has to name bool fields, and only bool fields.
	it('names a bool field of the message it is grouped under', () => {
		const messages = {
			raidBuffs: RaidBuffs,
			partyBuffs: PartyBuffs,
			individualBuffs: IndividualBuffs,
			debuffs: Debuffs,
		};

		for (const [group, fields] of Object.entries(retypedBuffFields)) {
			for (const name of fields) {
				const field = messages[group as keyof typeof messages].fields.find(f => f.localName === name);
				expect(field, `${group}.${name}`).toBeDefined();
				expect(field!.kind === 'scalar' && field!.T === ScalarType.BOOL, `${group}.${name} is ${field!.kind}`).toBe(true);
			}
		}
	});
});
