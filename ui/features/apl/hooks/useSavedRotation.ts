import { SavedRotation } from '@generated/proto/ui';
import { useSimHost } from '@sim/context/SimHostContext';
import { dropRetiredRotationFields } from '@sim/proto/rotation_field_migration';
import { omitDeep } from '@sim/utils/collections';
import { useSavedData } from '@ui-kit/hooks/useSavedData';

// The live rotation is uuid-stripped, so a stored entry carrying uuids never compares equal to it.
// SavedRotation carries no api_version, so nothing migrates an entry on the way in: the simple
// rotation's JSON string decodes here whatever it holds, and only the spec's parser further on
// throws over a reserved key.
const CODEC = {
	toJson: SavedRotation.toJson.bind(SavedRotation),
	fromJson: (json: any) => {
		const saved = omitDeep(SavedRotation.fromJson(json), ['uuid']);
		if (saved.rotation?.simple) {
			saved.rotation.simple.specRotationJson = dropRetiredRotationFields(saved.rotation.simple.specRotationJson);
		}
		return saved;
	},
};

export const useSavedRotation = () => useSavedData(useSimHost().getSavedRotationStorageKey(), CODEC);
