import { ConsumesSpec, Profession } from '@generated/proto/common';
import { RestorationShaman_Options as RestorationShamanOptions } from '@generated/proto/shaman';

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
