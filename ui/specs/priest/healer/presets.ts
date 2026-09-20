import { ConsumesSpec, Profession } from '@generated/proto/common';
import { HealerPriest_Options as HealerPriestOptions, PriestOptions_Armor } from '@generated/proto/priest';

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
