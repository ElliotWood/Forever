import * as PresetUtils from '@app/preset_utils';
import { ConsumesSpec, Debuffs, IndividualBuffs, PartyBuffs, Profession, Race, RaidBuffs, TristateEffect } from '@generated/proto/common';
import { HolyPaladin_Options as HolyPaladinOptions } from '@generated/proto/paladin';
import { SavedTalents } from '@generated/proto/ui';

// Defaults follow master's ui/holy_paladin.
// Our Forever sim's builds.
export const StandardTalents = PresetUtils.makePresetTalents('Standard', SavedTalents.create({ talentsString: '005321013025131251-503210302' }));
export const TalentsHolyHealer = PresetUtils.makePresetTalents('Holy 38/13/0', SavedTalents.create({ talentsString: '205320213225131051-50323' }));

export const TalentPresets = [StandardTalents, TalentsHolyHealer];

export const DefaultOptions = HolyPaladinOptions.create({
	classOptions: {},
});

// Master sets no consumables for this spec (the old TBC flask, food and potion don't exist in the Forever client).
export const DefaultConsumables = ConsumesSpec.create({});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	powerWordFortitude: TristateEffect.TristateEffectImproved,
	shadowProtection: true,
	thorns: TristateEffect.TristateEffectImproved,
});

// Master's DevotionAura spec option is the party Devotion Aura here.
export const DefaultPartyBuffs = PartyBuffs.create({
	devotionAura: true,
	manaSpringTotem: TristateEffect.TristateEffectRegular,
	moonkinAura: TristateEffect.TristateEffectRegular,
	strengthOfEarthTotem: TristateEffect.TristateEffectRegular,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfMight: true,
	blessingOfWisdom: true,
});

export const DefaultDebuffs = Debuffs.create({
	exposeArmor: TristateEffect.TristateEffectImproved,
	faerieFire: TristateEffect.TristateEffectRegular,
	insectSwarm: true,
	judgementOfLight: true,
	judgementOfWisdom: true,
	sunderArmor: true,
	thunderClap: TristateEffect.TristateEffectImproved,
});

export const OtherDefaults = {
	distanceFromTarget: 20,
	profession1: Profession.Enchanting,
	profession2: Profession.Jewelcrafting,
	race: Race.RaceHuman,
};
