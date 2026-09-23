import * as PresetUtils from '@app/preset_utils';
import { ConsumesSpec, Debuffs, IndividualBuffs, PartyBuffs, Profession, Race, RaidBuffs, TristateEffect } from '@generated/proto/common';
import { HealerPriest_Options as HealerPriestOptions, PriestOptions_Armor } from '@generated/proto/priest';
import { SavedTalents } from '@generated/proto/ui';

// Defaults follow master's ui/healing_priest.
export const DefaultOptions = HealerPriestOptions.create({
	classOptions: {
		armor: PriestOptions_Armor.InnerFire,
	},
});

// Master sets no consumables for this spec (the old TBC flask, food and potion don't exist in the Forever client).
export const DefaultConsumables = ConsumesSpec.create({});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	divineSpirit: TristateEffect.TristateEffectRegular,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	powerWordFortitude: TristateEffect.TristateEffectImproved,
});

export const DefaultPartyBuffs = PartyBuffs.create({
	moonkinAura: TristateEffect.TristateEffectRegular,
	strengthOfEarthTotem: TristateEffect.TristateEffectRegular,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfWisdom: true,
});

export const DefaultDebuffs = Debuffs.create({});

export const OtherDefaults = {
	distanceFromTarget: 20,
	profession1: Profession.Tailoring,
	profession2: Profession.Enchanting,
	race: Race.RaceDwarf,
};

// Talent presets, from master's ui/priest spec.
export const TalentsHolyHealer = PresetUtils.makePresetTalents('Holy 19/32/0', SavedTalents.create({ talentsString: '005203031302-2350510323000053' }));
export const TalentsDisciplineHealer = PresetUtils.makePresetTalents(
	'Discipline 35/16/0',
	SavedTalents.create({ talentsString: '005203031325101531-03505003' }),
);
export const TalentPresets = [TalentsHolyHealer, TalentsDisciplineHealer];
