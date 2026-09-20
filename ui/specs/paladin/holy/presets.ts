import * as PresetUtils from '@app/preset_utils';
import { ConsumesSpec, Profession } from '@generated/proto/common';
import { HolyPaladin_Options as HolyPaladinOptions } from '@generated/proto/paladin';

import StandardTalentsJson from './presets/talents/standard.talents.json';

// Talent builds from wowhead's TBC guide. Uses the wowhead calculator format, make the talents on
// https://www.wowhead.com/forever/talent-calc and copy the numbers in the url.
export const StandardTalents = PresetUtils.makePresetTalentsFromJSON(StandardTalentsJson);

export const DefaultOptions = HolyPaladinOptions.create({
	classOptions: {},
});

export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 22853, // Flask of Mighty Restoration
	foodId: 27666, // Golden Fish Sticks
	potId: 22832, // Super Mana Potion
});

export const OtherDefaults = {
	distanceFromTarget: 20,
	profession1: Profession.Enchanting,
	profession2: Profession.Jewelcrafting,
};
