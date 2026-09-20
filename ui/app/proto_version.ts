import { IndividualSimSettings } from '@generated/proto/ui';
import i18n from '@i18n/config';
import { migrateOldProto, type ProtoConversionMap } from '@sim/proto/proto_migration';
import { updateIndividualSimProtoVersion } from '@sim/state/serialization';
import { toastManager } from '@ui-kit/Toast';

// Party drums moved out of the player's consumables and into PartyBuffs in api version 7, and
// version 17 retired that party buff: a payload that named a Greater drum as its consumable has
// nowhere left to put it, so the converter drops the consumable and says so. It lives here rather
// than beside the shared migrations in `ui/sim` because it reports itself with a toast, and that
// layer may not reach `@ui-kit`.
// TODO: version 16 is the Forever talent rebuild -- every class's talent proto changed
// shape and StatSpellPenetration became StatSpellPiercing. No converter is registered for
// it on purpose: a saved TBC build's talent string has no meaning against a Forever tree,
// so there is nothing to migrate it to. Settings below version 16 fall through and the
// user re-picks their talents.
// Greater Drums of Battle, of War and of Restoration, which a pre-7 payload named as the player's
// own consumable to say the party had them.
const GREATER_DRUM_ITEM_IDS = [351355, 351360, 351358];

const TBC_CONVERSION_MAP: ProtoConversionMap<IndividualSimSettings> = new Map([
	[
		7,
		(oldProto: IndividualSimSettings) => {
			oldProto.apiVersion = 7;
			const oldPartyDrums = oldProto.player?.consumables?.drumsId as number;
			if (GREATER_DRUM_ITEM_IDS.includes(oldPartyDrums) && oldProto.player?.consumables) {
				oldProto.player.consumables.drumsId = 0;

				toastManager.add({
					variant: 'warning',
					delay: 8000,
					body: i18n.t('protoVersion.7.body', { ns: 'updates' }),
				});
			}

			return oldProto;
		},
	],
	// Version 15 renamed the shadow priest's Player oneof from `priest` to `dps_priest` (and the
	// Spec value to SpecDpsPriest): the sim covers any dps priest, not only shadow. Binary payloads
	// (share links) carry the field by number, and JSON ones (autosaved settings, the JSON importer)
	// are renamed before parsing by `parseLegacySettingsJson`, so there is nothing left to convert
	// here; the entry only stamps the version.
	[15, (oldProto: IndividualSimSettings) => oldProto],
	// Version 17 types the ghost-talent buff fields bool. JSON payloads are rewritten before they
	// are parsed (`migrateRetypedBuffFields`), and a binary one carries 1 or 2, which decode as
	// true, so this entry only stamps the version.
	[17, (oldProto: IndividualSimSettings) => oldProto],
]);

/**
 * TBC's own migrations, then the shared ones. The order matters: the shared pass stamps the proto
 * as current, after which nothing below would ever run again.
 */
export function updateIndividualProtoVersion(settingsProto: IndividualSimSettings) {
	migrateOldProto(settingsProto, settingsProto.apiVersion, TBC_CONVERSION_MAP);
	updateIndividualSimProtoVersion(settingsProto);
}
