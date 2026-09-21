import { Player } from '@generated/proto/api';
import { ConsumesSpec } from '@generated/proto/common';
import { IndividualSimSettings } from '@generated/proto/ui';
import { CURRENT_API_VERSION } from '@sim/constants/other';
import { describe, expect, it } from 'vitest';

const { updateIndividualProtoVersion } = await import('./proto_version');

const settings = (apiVersion: number) =>
	IndividualSimSettings.create({
		apiVersion,
		player: Player.create({ consumables: ConsumesSpec.create({ foodId: 27657 }) }),
	});

describe('updateIndividualProtoVersion', () => {
	it('stamps every migrated payload as current, so it is not migrated twice', () => {
		const proto = settings(6);

		updateIndividualProtoVersion(proto);

		expect(proto.apiVersion).toBe(CURRENT_API_VERSION);
	});

	// Both registered converters only stamp: version 15's oneof rename happens before parsing, and
	// version 17's retypes and retirements happen in the raw-JSON pre-pass.
	it('leaves the payload itself alone on the way up', () => {
		const proto = settings(14);

		updateIndividualProtoVersion(proto);

		expect(proto.apiVersion).toBe(CURRENT_API_VERSION);
		expect(proto.player?.consumables?.foodId).toBe(27657);
	});

	it('leaves a payload that is already current alone', () => {
		const proto = settings(CURRENT_API_VERSION);

		updateIndividualProtoVersion(proto);

		expect(proto.apiVersion).toBe(CURRENT_API_VERSION);
		expect(proto.player?.consumables?.foodId).toBe(27657);
	});
});
