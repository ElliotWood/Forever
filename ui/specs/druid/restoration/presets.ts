import * as PresetUtils from '@app/preset_utils';
import { ConsumesSpec, Profession, UnitReference } from '@generated/proto/common';
import { RestorationDruid_Options as RestorationDruidOptions } from '@generated/proto/druid';

import DreamstateTalentsJson from './presets/talents/dreamstate.talents.json';
import TreeOfLifeTalentsJson from './presets/talents/tree_of_life.talents.json';

// Talent builds from wowhead's TBC guide. Uses the wowhead calculator format, make the talents on
// https://www.wowhead.com/forever/talent-calc and copy the numbers in the url.
export const TreeOfLifeTalents = PresetUtils.makePresetTalentsFromJSON(TreeOfLifeTalentsJson);
export const DreamstateTalents = PresetUtils.makePresetTalentsFromJSON(DreamstateTalentsJson);

export const DefaultOptions = RestorationDruidOptions.create({
	classOptions: {
		innervateTarget: UnitReference.create(),
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
