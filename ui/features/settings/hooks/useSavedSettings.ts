import { SavedSettings } from '@generated/proto/ui';
import { useSimHost } from '@sim/context/SimHostContext';
import { migrateRetypedBuffFields } from '@sim/proto/buff_field_migration';
import type { SavedDataCodec } from '@ui-kit/hooks/useSavedData';
import { useSavedData } from '@ui-kit/hooks/useSavedData';

// SavedSettings carries no api_version, so every entry goes through the buff migration: a save
// holding "TristateEffectImproved" where the field is now a bool would throw in `fromJson`, and
// `useSavedData` drops any entry whose codec throws, leaving the panel short of it without a word.
const savedSettingsCodec: SavedDataCodec<SavedSettings> = {
	toJson: settings => SavedSettings.toJson(settings),
	fromJson: json => {
		migrateRetypedBuffFields(json);
		return SavedSettings.fromJson(json, { ignoreUnknownFields: true });
	},
};

export const useSavedSettings = () => useSavedData(useSimHost().getSavedSettingsStorageKey(), savedSettingsCodec);
