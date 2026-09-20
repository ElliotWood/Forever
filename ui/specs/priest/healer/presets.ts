import * as PresetUtils from '@app/preset_utils';
import { ConsumesSpec, Profession } from '@generated/proto/common';
import { HealerPriest_Options as HealerPriestOptions, PriestOptions_Armor } from '@generated/proto/priest';

import CircleOfHealingTalentsJson from './presets/talents/circle_of_healing.talents.json';

// Talent builds from wowhead's TBC guide. Uses the wowhead calculator format, make the talents on
// https://www.wowhead.com/forever/talent-calc and copy the numbers in the url.
export const CircleOfHealingTalents = PresetUtils.makePresetTalentsFromJSON(CircleOfHealingTalentsJson);

export const DefaultOptions = HealerPriestOptions.create({
	classOptions: {
		armor: PriestOptions_Armor.InnerFire,
	},
});

export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 22853, // Flask of Mighty Restoration
	foodId: 27666, // Golden Fish Sticks
	potId: 22832, // Super Mana Potion
});

export const OtherDefaults = {
	distanceFromTarget: 20,
	profession1: Profession.Tailoring,
	profession2: Profession.Enchanting,
};
