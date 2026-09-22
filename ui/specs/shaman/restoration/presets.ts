import * as PresetUtils from '@app/preset_utils';
import { ConsumesSpec, Profession } from '@generated/proto/common';
import { RestorationShaman_Options as RestorationShamanOptions } from '@generated/proto/shaman';
import { SavedTalents } from '@generated/proto/ui';

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

// Talent presets, from master's ui/shaman spec.
export const TalentsTankHealing = PresetUtils.makePresetTalents('Tank Healing', SavedTalents.create({ talentsString: '--5533523315513151' }));
export const TalentsRaidHealing = PresetUtils.makePresetTalents('Raid Healing', SavedTalents.create({ talentsString: '-005102-5530500315513151' }));
export const TalentsRestoration = PresetUtils.makePresetTalents('Restoration 0/3/48', SavedTalents.create({ talentsString: '-003-5532503315513151' }));
export const TalentPresets = [TalentsTankHealing, TalentsRaidHealing, TalentsRestoration];
