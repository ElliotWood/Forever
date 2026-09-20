import { Player } from '@generated/proto/api';
import { ConsumesSpec } from '@generated/proto/common';
import { IndividualSimSettings } from '@generated/proto/ui';
import { CURRENT_API_VERSION } from '@sim/constants/other';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const added = vi.hoisted(() => vi.fn());
vi.mock('@ui-kit/Toast', () => ({ toastManager: { add: added } }));
vi.mock('@i18n/config', () => ({ default: { t: (key: string) => key } }));

const { updateIndividualProtoVersion } = await import('./proto_version');

const GREATER_DRUMS_OF_BATTLE = 351355;

const settings = (apiVersion: number, drumsId: number) =>
	IndividualSimSettings.create({
		apiVersion,
		player: Player.create({ consumables: ConsumesSpec.create({ drumsId }) }),
	});

describe('updateIndividualProtoVersion', () => {
	beforeEach(() => added.mockClear());

	it('clears a pre-7 payload’s party drums, which version 17 left nowhere to put', () => {
		const proto = settings(6, GREATER_DRUMS_OF_BATTLE);

		updateIndividualProtoVersion(proto);

		expect(proto.player?.consumables?.drumsId).toBe(0);
		expect(added).toHaveBeenCalledTimes(1);
	});

	it('leaves a payload that is already past 7 alone', () => {
		const proto = settings(7, GREATER_DRUMS_OF_BATTLE);

		updateIndividualProtoVersion(proto);

		expect(proto.player?.consumables?.drumsId).toBe(GREATER_DRUMS_OF_BATTLE);
		expect(added).not.toHaveBeenCalled();
	});

	it('stamps every migrated payload as current, so it is not migrated twice', () => {
		const proto = settings(6, GREATER_DRUMS_OF_BATTLE);

		updateIndividualProtoVersion(proto);

		expect(proto.apiVersion).toBe(CURRENT_API_VERSION);
	});

	it('stamps a version-14 payload as current without touching it: the priest oneof rename happens before parsing', () => {
		const proto = settings(14, GREATER_DRUMS_OF_BATTLE);

		updateIndividualProtoVersion(proto);

		expect(proto.apiVersion).toBe(CURRENT_API_VERSION);
		expect(proto.player?.consumables?.drumsId).toBe(GREATER_DRUMS_OF_BATTLE);
		expect(added).not.toHaveBeenCalled();
	});

	it('says nothing when the old payload carried no drums', () => {
		const proto = settings(6, 0);

		updateIndividualProtoVersion(proto);

		expect(added).not.toHaveBeenCalled();
	});

	it('leaves the player’s own drums consumable alone, which is not the party buff', () => {
		const proto = settings(6, 29529);

		updateIndividualProtoVersion(proto);

		expect(proto.player?.consumables?.drumsId).toBe(29529);
		expect(added).not.toHaveBeenCalled();
	});
});
