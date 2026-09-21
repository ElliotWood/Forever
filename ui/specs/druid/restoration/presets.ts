import * as PresetUtils from '@app/preset_utils';
import { ConsumesSpec, Profession, UnitReference } from '@generated/proto/common';
import { RestorationDruid_Options as RestorationDruidOptions } from '@generated/proto/druid';
import { SavedTalents } from '@generated/proto/ui';

export const RestorationTalents = PresetUtils.makePresetTalents('Restoration 10/0/41', SavedTalents.create({ talentsString: '05302--5053035153113051' }));

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
