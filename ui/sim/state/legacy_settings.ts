import { IndividualSimSettings } from '@generated/proto/ui';

import { convertClassicSettingsJson, isClassicSettingsJson } from './classic_links';

// Settings JSON written before a proto rename. Share links are binary, so field numbers carry
// them across; only the JSON forms (autosaved settings in localStorage and the JSON importer) name
// fields and need a hand.

// The shadow priest's Player oneof was `priest` until the spec became `dps_priest` (the sim covers
// any dps priest, not only shadow). Renames the field in place so an old blob parses.
export const migrateLegacySettingsJson = (json: unknown): unknown => {
	if (!json || typeof json !== 'object') return json;
	const player = (json as { player?: unknown }).player;
	if (!player || typeof player !== 'object') return json;
	const spec = player as Record<string, unknown>;
	if ('priest' in spec && !('dpsPriest' in spec)) {
		const { priest, ...rest } = spec;
		return { ...(json as object), player: { ...rest, dpsPriest: priest } };
	}
	return json;
};

// JSON text in, migrated object out; a parse failure surfaces to the caller as before.
// Settings from the classic-engine site (master before the switch) are another proto altogether.
export const parseLegacySettingsJson = (text: string): unknown => {
	const json = JSON.parse(text);
	return isClassicSettingsJson(json) ? IndividualSimSettings.toJson(convertClassicSettingsJson(json)) : migrateLegacySettingsJson(json);
};
