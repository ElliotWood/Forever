import { Drums } from '@generated/proto/common';
import { IndividualSimSettings } from '@generated/proto/ui';
import i18n from '@i18n/config';
import { migrateOldProto, type ProtoConversionMap } from '@sim/proto/proto_migration';
import { updateIndividualSimProtoVersion } from '@sim/state/serialization';
import { toastManager } from '@ui-kit/Toast';

// Party drums moved out of the player's consumables and into PartyBuffs in api version 7. The
// converter lives here rather than beside the shared migrations in `ui/sim` because it reports
// itself with a toast, and that layer may not reach `@ui-kit`.
// TODO: version 16 is the Forever talent rebuild -- every class's talent proto changed
// shape and StatSpellPenetration became StatSpellPiercing. No converter is registered for
// it on purpose: a saved TBC build's talent string has no meaning against a Forever tree,
// so there is nothing to migrate it to. Settings below version 16 fall through and the
// user re-picks their talents.
const TBC_CONVERSION_MAP: ProtoConversionMap<IndividualSimSettings> = new Map([
	[
		7,
		(oldProto: IndividualSimSettings) => {
			oldProto.apiVersion = 7;
			const oldPartyDrums = oldProto.player?.consumables?.drumsId as number;
			if (oldPartyDrums && oldProto.partyBuffs) {
				switch (oldPartyDrums) {
					case 351355: // Greater Drums of Battle
						oldProto.partyBuffs.drums = Drums.LesserDrumsOfBattle;
						break;
					case 351360: // Greater Drums of War
						oldProto.partyBuffs.drums = Drums.LesserDrumsOfWar;
						break;
					case 351358: // Greater Drums of Restoration
						oldProto.partyBuffs.drums = Drums.LesserDrumsOfRestoration;
						break;
				}

				if (oldProto.player?.consumables) oldProto.player.consumables.drumsId = 0;

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
]);

/**
 * TBC's own migrations, then the shared ones. The order matters: the shared pass stamps the proto
 * as current, after which nothing below would ever run again.
 */
export function updateIndividualProtoVersion(settingsProto: IndividualSimSettings) {
	migrateOldProto(settingsProto, settingsProto.apiVersion, TBC_CONVERSION_MAP);
	updateIndividualSimProtoVersion(settingsProto);
}
