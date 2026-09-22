import * as PresetUtils from '@app/preset_utils';
import { ConsumesSpec, Profession } from '@generated/proto/common';
import { HealerPriest_Options as HealerPriestOptions, PriestOptions_Armor } from '@generated/proto/priest';
import { SavedTalents } from '@generated/proto/ui';

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

// Talent presets, from master's ui/priest spec.
export const TalentsHolyHealer = PresetUtils.makePresetTalents('Holy 19/32/0', SavedTalents.create({ talentsString: '005203031302-2350510323000053' }));
export const TalentsDisciplineHealer = PresetUtils.makePresetTalents('Discipline 35/16/0', SavedTalents.create({ talentsString: '005203031325101531-03505003' }));
export const TalentPresets = [TalentsHolyHealer, TalentsDisciplineHealer];
