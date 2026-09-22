// A simple rotation travels as the opaque JSON string `APLRotation.simple.spec_rotation_json`, and
// a spec's `rotationFromJson` reads it without `ignoreUnknownFields`: a key the proto has since
// reserved makes the parser throw, and the caller that catches hands back a blank rotation, so the
// whole rotation is lost rather than the one key. Every reader of a payload it did not just write
// drops the reserved keys first.
//
// v17 reserved DpsWarrior.Rotation's bloodlust_timing. Both spellings are listed because a payload
// may have been written by either JSON printer.
const RETIRED_ROTATION_FIELDS = ['bloodlustTiming', 'bloodlust_timing'];

// The string is handed back untouched unless a reserved key was actually there, so a rotation that
// needs nothing keeps the exact bytes it was stored with.
export const dropRetiredRotationFields = (specRotationJson: string): string => {
	if (!specRotationJson) return specRotationJson;

	try {
		const parsed = JSON.parse(specRotationJson);
		if (!parsed || typeof parsed !== 'object') return specRotationJson;

		const retired = RETIRED_ROTATION_FIELDS.filter(field => field in parsed);
		if (!retired.length) return specRotationJson;

		retired.forEach(field => delete parsed[field]);
		return JSON.stringify(parsed);
	} catch {
		return specRotationJson;
	}
};
