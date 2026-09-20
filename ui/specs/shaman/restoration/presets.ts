import * as PresetUtils from '@app/preset_utils';
import { ConsumesSpec, Profession } from '@generated/proto/common';
import { RestorationShaman_Options as RestorationShamanOptions } from '@generated/proto/shaman';

import ElementalWardingTalentsJson from './presets/talents/elemental_warding.talents.json';
import EnhancingTotemsTalentsJson from './presets/talents/enhancing_totems.talents.json';

// Talent builds from wowhead's TBC guide. Uses the wowhead calculator format, make the talents on
// https://www.wowhead.com/forever/talent-calc and copy the numbers in the url.
export const ElementalWardingTalents = PresetUtils.makePresetTalentsFromJSON(ElementalWardingTalentsJson);
export const EnhancingTotemsTalents = PresetUtils.makePresetTalentsFromJSON(EnhancingTotemsTalentsJson);

export const DefaultOptions = RestorationShamanOptions.create({
	classOptions: {},
});

export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 22853, // Flask of Mighty Restoration
	foodId: 27666, // Golden Fish Sticks
	potId: 22832, // Super Mana Potion
});

export const OtherDefaults = {
	distanceFromTarget: 20,
	profession1: Profession.Leatherworking,
	profession2: Profession.Enchanting,
};
