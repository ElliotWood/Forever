import * as PresetUtils from '@app/preset_utils';
import { ConsumesSpec } from '@generated/proto/common';
import { Rogue_Options as RogueOptions } from '@generated/proto/rogue';
import { SavedTalents } from '@generated/proto/ui';

import SinisterAPL from './apls/swords.apl.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const SINSITER_APL = PresetUtils.makePresetAPLRotation('Rogue (Check Variables for Backstab/Shiv)', SinisterAPL);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.

export const Talents = {
	name: 'Combat Swords',
	data: SavedTalents.create({
		talentsString: '0053201252-023305200005015002321151',
	}),
};

export const DefaultOptions = RogueOptions.create({
	classOptions: {},
});

export const DefaultConsumables = ConsumesSpec.create({
	battleElixirId: 22831,
	guardianElixirId: 32062,
	foodId: 33872,
	potId: 22838,
	conjuredId: 7676,
	ohImbueId: 27186,
});

export const OtherDefaults = {
	distanceFromTarget: 5,
};
