import { ConsumesSpec, Profession } from '@generated/proto/common';
import { HolyPaladin_Options as HolyPaladinOptions } from '@generated/proto/paladin';

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
