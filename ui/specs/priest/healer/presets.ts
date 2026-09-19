import * as PresetUtils from '@app/preset_utils';
import { ConsumesSpec, Profession } from '@generated/proto/common';
import { HealerPriest_Options as HealerPriestOptions, PriestOptions_Armor } from '@generated/proto/priest';

import P3Gear from './gear_sets/p3.gear.json';
import P3T5Gear from './gear_sets/p3_t5.gear.json';
import PreraidGear from './gear_sets/preraid.gear.json';
import DefaultEpJson from './presets/ep/default.ep.json';
import CircleOfHealingTalentsJson from './presets/talents/circle_of_healing.talents.json';

export const P3_PRESET = PresetUtils.makePresetGear('P3 BiS (T6 4P)', P3Gear);
export const P3_T5_PRESET = PresetUtils.makePresetGear('P3 BiS (T5 4P)', P3T5Gear);
export const PRERAID_PRESET = PresetUtils.makePresetGear('Pre-raid', PreraidGear);

// Stat weights with healing power = 1, ordered after wowhead's TBC stat priority for the spec
// (healing > haste > Spirit > Intellect > MP5 > crit). Not sim-derived: there is no healing sim, so these only sort the gear picker
// and drive the gem optimizer.
export const DEFAULT_EP_PRESET = PresetUtils.makePresetEpWeightsFromJSON(DefaultEpJson);

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
