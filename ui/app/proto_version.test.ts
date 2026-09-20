import { Player } from '@generated/proto/api';
import { PartyBuffs } from '@generated/proto/buffs';
import { ConsumesSpec, Drums } from '@generated/proto/common';
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
		partyBuffs: PartyBuffs.create(),
	});

describe('updateIndividualProtoVersion', () => {
	beforeEach(() => added.mockClear());

	it('moves a pre-7 payload’s party drums onto the party buffs and clears the consumable', () => {
		const proto = settings(6, GREATER_DRUMS_OF_BATTLE);

		updateIndividualProtoVersion(proto);

		expect(proto.partyBuffs?.drums).toBe(Drums.LesserDrumsOfBattle);
		expect(proto.player?.consumables?.drumsId).toBe(0);
		expect(added).toHaveBeenCalledTimes(1);
	});

	it('leaves a payload that is already past 7 alone', () => {
		const proto = settings(7, GREATER_DRUMS_OF_BATTLE);

		updateIndividualProtoVersion(proto);

		expect(proto.partyBuffs?.drums).toBe(Drums.DrumsUnknown);
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
});
