import * as PresetUtils from '@app/preset_utils';
import { ConsumesSpec, Profession } from '@generated/proto/common';
import { HolyPaladin_Options as HolyPaladinOptions } from '@generated/proto/paladin';
import { SavedTalents } from '@generated/proto/ui';

// Our Forever sim's builds.
export const StandardTalents = PresetUtils.makePresetTalents('Standard', SavedTalents.create({ talentsString: '005321013025131251-503210302' }));
export const TalentsHolyHealer = PresetUtils.makePresetTalents('Holy 38/13/0', SavedTalents.create({ talentsString: '205320213225131051-50323' }));

export const TalentPresets = [StandardTalents, TalentsHolyHealer];

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
