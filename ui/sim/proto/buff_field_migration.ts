// Api version 17 types 25 buff fields bool where they were a TristateEffect, because the Improved
// talent behind each of them has no node in a Forever trait tree; the names are pinned by
// TestRetypedFieldsMatchTheMigration in tools/gen_buffs_proto.
//
// The version converters in `proto_migration` run on an already-parsed proto, and protobuf-ts
// `fromJson` throws on "TristateEffectImproved" in a bool field long before they are reached, so the
// rewrite happens here, on the JSON. Share links carry these fields as varints, where 1 and 2 both
// decode as true, so `fromBinary` needs nothing; the binary case in buff_field_migration.test.ts
// reads one off the wire.

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

type BuffScope = keyof typeof retypedBuffFields;

type JsonObject = Record<string, unknown>;

const asObject = (value: unknown): JsonObject | null => (typeof value === 'object' && value !== null && !Array.isArray(value) ? (value as JsonObject) : null);

// A message stamped with version 17 or later already holds bools. The messages that carry buffs
// without an api_version of their own (SavedSettings, Raid, Party) are always rewritten, which costs
// nothing: a value that is already a bool is left as it is.
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
};

const migratePlayer = (player: unknown) => {
	const message = asObject(player);
	if (!message) return;
	rewriteBuffs(message.buffs, 'individualBuffs');
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
 * Rewrites the retyped buff fields of a parsed-but-not-yet-decoded blob in place, so that `fromJson`
 * accepts a payload written before api version 17. The default shape covers the settings envelopes
 * -- IndividualSimSettings, SavedSettings and RaidSimSettings -- and the messages they nest; a Raid,
 * Party or Player passed on its own names itself.
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
