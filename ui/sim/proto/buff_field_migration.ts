// Api version 17 does two things to the buff messages, and one to the player's consumables. It
// types 25 fields bool where they were a TristateEffect, because the Improved talent behind each of
// them has no node in a Forever trait tree; the names are pinned by
// TestRetypedFieldsMatchTheMigration in tools/gen_buffs_proto. It retires 33 fields outright,
// because the Forever client describes no spell for them; the numbers they gave up are `reserved`
// in proto/buffs.proto and listed in buffmanifest.Retired. And it retires the drums consumable for
// the same reason, which proto/common.proto reserves on ConsumesSpec.
//
// The version converters in `proto_migration` run on an already-parsed proto, and protobuf-ts
// `fromJson` throws on "TristateEffectImproved" in a bool field - and, without
// `ignoreUnknownFields`, on a retired key - long before they are reached, so both happen here, on
// the JSON. Share links carry these fields as varints: 1 and 2 both decode as true, and a retired
// field number is skipped as an unknown one, so `fromBinary` needs nothing; the binary case in
// buff_field_migration.test.ts reads both off the wire.

const RETYPED_API_VERSION = 17;

export const retypedBuffFields = {
	raidBuffs: ['powerWordFortitude', 'divineSpirit', 'giftOfTheWild', 'thorns'],
	partyBuffs: [
		'bloodPact',
		'moonkinAura',
		'leaderOfThePack',
		'devotionAura',
		'retributionAura',
		'concentrationAura',
		'graceOfAirTotem',
		'strengthOfEarthTotem',
		'windfuryTotem',
		'battleShout',
		'commandingShout',
	],
	individualBuffs: ['blessingOfWisdom', 'blessingOfMight'],
	debuffs: [
		'improvedSealOfTheCrusader',
		'curseOfElements',
		'exposeArmor',
		'faerieFire',
		'huntersMark',
		'demoralizingRoar',
		'demoralizingShout',
		'thunderClap',
	],
} as const;

// The proto names of the 33 retired fields, which is how buffmanifest.Retired and
// proto/buffs.proto's `reserved` lines spell them.
export const retiredBuffFields = {
	raidBuffs: ['bloodlust'],
	partyBuffs: [
		'braided_eternium_chain',
		'bs_solarian_sapphire',
		'chain_of_the_twilight_owl',
		'draenei_racial_caster',
		'draenei_racial_melee',
		'drums',
		'eye_of_the_night',
		'ferocious_inspiration',
		'jade_pendant_of_blasting',
		'sanctity_aura',
		'snapshot_bs_booming_voice_rank',
		'snapshot_bs_solarian_sapphire',
		'snapshot_improved_strength_of_earth_totem',
		'snapshot_improved_wrath_of_air_totem',
		'soe_enhancement_2pt4',
		'totem_of_wrath',
		'tranquil_air_totem',
		'wrath_of_air_totem',
	],
	individualBuffs: ['blessing_of_sanctuary', 'unleashed_rage'],
	debuffs: [
		'blood_frenzy',
		'expose_weakness_hunter_agility',
		'expose_weakness_uptime',
		'hemorrhage_uptime',
		'improved_scorch',
		'isb_uptime',
		'joc_retribution_2pt4',
		'misery',
		'screech',
		'shadow_embrace',
		'shadow_weaving',
		'winters_chill',
	],
} as const;

// A payload may spell a field three ways: the proto name, protoc's json name, and the property
// protobuf-ts generates, which differs from the json name by capitalising after a digit as well
// ("joc_retribution_2pt4" is "jocRetribution2pt4" and "jocRetribution2Pt4"). `fromJson` accepts the
// first two and `toJson` writes the second, and a hand-built blob may hold the third.
export const retiredFieldSpellings = (protoName: string): string[] => {
	const jsonName = protoName.replace(/_([a-z0-9])/g, (_, character: string) => character.toUpperCase());
	const propertyName = jsonName.replace(/(\d)([a-z])/g, (_, digit: string, character: string) => digit + character.toUpperCase());
	return [...new Set([protoName, jsonName, propertyName])];
};

// Api version 17 also retires the player's own drums: ConsumesSpec reserves field 12 and the name
// `drums_id`, because the Forever client describes no drum item at all. A settings blob written
// before the bump still names it, and it sits on the player rather than on a buff message.
export const retiredConsumesFields = ['drums_id'] as const;

type BuffScope = keyof typeof retypedBuffFields;

type JsonObject = Record<string, unknown>;

const asObject = (value: unknown): JsonObject | null => (typeof value === 'object' && value !== null && !Array.isArray(value) ? (value as JsonObject) : null);

// A message stamped with version 17 or later already holds bools and has no retired fields. The
// messages that carry buffs without an api_version of their own (SavedSettings, Raid, Party) are
// always rewritten, which costs nothing: a value that is already a bool is left as it is.
const isOutdated = (message: JsonObject): boolean => {
	const version = message.apiVersion;
	return typeof version !== 'number' || version < RETYPED_API_VERSION;
};

// `toJson` writes an enum as its name and a `fromJson` payload may also hold the number, so both
// spellings of every state are accepted. Anything else - a bool, a missing field, a typo - is left
// for the parser to judge.
const rewriteBuffs = (buffs: unknown, scope: BuffScope) => {
	const message = asObject(buffs);
	if (!message) return;

	for (const field of retypedBuffFields[scope]) {
		const value = message[field];
		if (value === 'TristateEffectMissing' || value === 0) {
			message[field] = false;
		} else if (value === 'TristateEffectRegular' || value === 'TristateEffectImproved' || value === 1 || value === 2) {
			message[field] = true;
		}
	}

	for (const field of retiredBuffFields[scope]) {
		for (const spelling of retiredFieldSpellings(field)) {
			delete message[spelling];
		}
	}
};

const dropRetiredConsumes = (consumables: unknown) => {
	const message = asObject(consumables);
	if (!message) return;

	for (const field of retiredConsumesFields) {
		for (const spelling of retiredFieldSpellings(field)) {
			delete message[spelling];
		}
	}
};

const migratePlayer = (player: unknown) => {
	const message = asObject(player);
	if (!message) return;
	rewriteBuffs(message.buffs, 'individualBuffs');
	dropRetiredConsumes(message.consumables);
};

const migrateParty = (party: unknown) => {
	const message = asObject(party);
	if (!message) return;
	rewriteBuffs(message.buffs, 'partyBuffs');
	if (Array.isArray(message.players)) message.players.forEach(migratePlayer);
};

const migrateRaid = (raid: unknown) => {
	const message = asObject(raid);
	if (!message) return;
	rewriteBuffs(message.buffs, 'raidBuffs');
	rewriteBuffs(message.debuffs, 'debuffs');
	if (Array.isArray(message.parties)) message.parties.forEach(migrateParty);
};

// Raid, Party and Player all call their buffs `buffs`, and an empty repeated field is absent from
// the JSON, so a blob cannot be told apart by its shape: the caller names the message it holds.
export type BuffMessageShape = 'settings' | 'raid' | 'party' | 'player';

/**
 * Rewrites the retyped buff fields of a parsed-but-not-yet-decoded blob in place and drops the
 * retired ones, so that `fromJson` accepts a payload written before api version 17. The default
 * shape covers the settings envelopes -- IndividualSimSettings, SavedSettings and RaidSimSettings --
 * and the messages they nest; a Raid, Party or Player passed on its own names itself.
 */
export function migrateRetypedBuffFields(json: unknown, shape: BuffMessageShape = 'settings'): void {
	const message = asObject(json);
	if (!message || !isOutdated(message)) return;

	switch (shape) {
		case 'raid':
			migrateRaid(message);
			return;
		case 'party':
			migrateParty(message);
			return;
		case 'player':
			migratePlayer(message);
			return;
	}

	// IndividualSimSettings, and SavedSettings, which names the individual buffs `playerBuffs`.
	rewriteBuffs(message.raidBuffs, 'raidBuffs');
	rewriteBuffs(message.partyBuffs, 'partyBuffs');
	rewriteBuffs(message.debuffs, 'debuffs');
	rewriteBuffs(message.playerBuffs, 'individualBuffs');
	migratePlayer(message.player);

	// RaidSimSettings.
	migrateRaid(message.raid);
}
