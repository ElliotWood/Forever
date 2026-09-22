import * as PresetUtils from '@app/preset_utils';
import { ConsumesSpec, Debuffs, IndividualBuffs, PartyBuffs, Profession, Race, RaidBuffs, TristateEffect } from '@generated/proto/common';
import { RestorationShaman_Options as RestorationShamanOptions, ShamanImbue } from '@generated/proto/shaman';
import { SavedTalents } from '@generated/proto/ui';

// Defaults follow master's ui/restoration_shaman.
export const DefaultOptions = RestorationShamanOptions.create({
	classOptions: {
		imbueMh: ShamanImbue.RockbiterWeapon,
	},
	earthShieldPPM: 0,
});

// Master sets no flask or food for this spec (the old TBC flask, food and potion don't exist in the Forever client).
export const DefaultConsumables = ConsumesSpec.create({});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	divineSpirit: TristateEffect.TristateEffectRegular,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
});

export const DefaultPartyBuffs = PartyBuffs.create({
	moonkinAura: TristateEffect.TristateEffectRegular,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
});

export const DefaultDebuffs = Debuffs.create({
	faerieFire: TristateEffect.TristateEffectRegular,
	judgementOfWisdom: true,
});

export const OtherDefaults = {
	distanceFromTarget: 20,
	profession1: Profession.Leatherworking,
	profession2: Profession.Enchanting,
	race: Race.RaceOrc,
};

// Talent presets, from master's ui/shaman spec.
export const TalentsTankHealing = PresetUtils.makePresetTalents('Tank Healing', SavedTalents.create({ talentsString: '--5533523315513151' }));
export const TalentsRaidHealing = PresetUtils.makePresetTalents('Raid Healing', SavedTalents.create({ talentsString: '-005102-5530500315513151' }));
export const TalentsRestoration = PresetUtils.makePresetTalents('Restoration 0/3/48', SavedTalents.create({ talentsString: '-003-5532503315513151' }));
export const TalentPresets = [TalentsTankHealing, TalentsRaidHealing, TalentsRestoration];
