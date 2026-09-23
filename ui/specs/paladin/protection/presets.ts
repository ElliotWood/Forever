import * as PresetUtils from '@app/preset_utils';
import { ConsumesSpec, Debuffs, IndividualBuffs, PartyBuffs, Profession, Race, RaidBuffs, TristateEffect } from '@generated/proto/common';
import { ProtectionPaladin_Options as ProtectionPaladinOptions } from '@generated/proto/paladin';
import { SavedTalents } from '@generated/proto/ui';

import DefaultApl from './apls/default.apl.json';
import P5Apl from './apls/p5.apl.json';
import LaunchGear from './gear_sets/launch.gear.json';

// Our Forever sim's rotations.
export const APL_PRESET = PresetUtils.makePresetAPLRotation('Basic Prot', DefaultApl);
export const APL_P5 = PresetUtils.makePresetAPLRotation('P5 Prot', P5Apl);

// Our Forever sim's builds.
export const P4ProtTalents = PresetUtils.makePresetTalents('P4 Prot', SavedTalents.create({ talentsString: '052003003-5530513321301501' }));
export const P5ProtTalents = PresetUtils.makePresetTalents('P5 Prot', SavedTalents.create({ talentsString: '055003-5530513321301501' }));
export const TalentsProtection = PresetUtils.makePresetTalents('Protection 0/45/6', SavedTalents.create({ talentsString: '-5532513321301551-15' }));

export const TalentPresets = [P5ProtTalents, P4ProtTalents, TalentsProtection];
export const DefaultTalents = P5ProtTalents;

export const DefaultOptions = ProtectionPaladinOptions.create({
	classOptions: {},
});

// Defaults below are what master's ui/protection_paladin opens with (currentSettings on a fresh
// profile); its raid-wide Battle Shout is a party buff here, its personal Blessing of Sanctuary
// the individual one. Juju Power/Might, R.O.I.D.S., Dragonbreath Chili, Rumsey Rum, Greater Arcane
// Elixir and Elixir of Fortitude have no slot here; Tender Wolf Steak, Greater Stoneshield, Superior
// Defense are not in this db.
export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 13510, // Flask of the Titans
	battleElixirId: 13452, // Elixir of the Mongoose
	guardianElixirId: 13445, // Elixir of Superior Defense
	foodId: 18045, // Tender Wolf Steak
	potId: 13455, // Greater Stoneshield Potion
	conjuredId: 12662, // Demonic Rune
	explosiveId: 18641, // Dense Dynamite
});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	divineSpirit: TristateEffect.TristateEffectRegular,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	powerWordFortitude: TristateEffect.TristateEffectImproved,
});

export const DefaultPartyBuffs = PartyBuffs.create({
	battleShout: TristateEffect.TristateEffectImproved,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfWisdom: true,
});

export const DefaultDebuffs = Debuffs.create({
	curseOfRecklessness: true,
	exposeArmor: TristateEffect.TristateEffectImproved,
	faerieFire: TristateEffect.TristateEffectRegular,
	giftOfArthas: true,
	judgementOfTheCrusader: true,
	judgementOfWisdom: true,
});

export const OtherDefaults = {
	reactionTime: 200, // master's default
	profession1: Profession.Blacksmithing,
	profession2: Profession.Engineering,
	distanceFromTarget: 5,
	race: Race.RaceHuman,
};

// Our Forever sim's gear presets (master ui/<spec>/gear_sets).
export const GEAR_LAUNCH = PresetUtils.makePresetGear('Launch', LaunchGear);
export const DEFAULT_GEAR = GEAR_LAUNCH;
export const GEAR_PRESETS = [GEAR_LAUNCH];
