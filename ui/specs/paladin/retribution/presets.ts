import * as PresetUtils from '@app/preset_utils';
import { ConsumesSpec, Debuffs, IndividualBuffs, PartyBuffs, Profession, Race, RaidBuffs, TristateEffect } from '@generated/proto/common';
import { RetributionPaladin_Options as RetributionPaladinOptions } from '@generated/proto/paladin';
import { SavedTalents } from '@generated/proto/ui';

import DefaultApl from './apls/default.apl.json';
import LaunchGear from './gear_sets/launch.gear.json';

// Our Forever sim's Seal of Command / Seal of Righteousness twist rotation.
export const APL_PRESET = PresetUtils.makePresetAPLRotation('Basic Ret', DefaultApl);

// Our Forever sim's builds.
export const P4RetTalents = PresetUtils.makePresetTalents('P4/P5 Ret', SavedTalents.create({ talentsString: '0550030022001--052251310002330321' }));
export const TalentsRetribution = PresetUtils.makePresetTalents('Retribution 10/0/41', SavedTalents.create({ talentsString: '250003--552250312012331321' }));

export const TalentPresets = [P4RetTalents, TalentsRetribution];
export const DefaultTalents = P4RetTalents;

export const DefaultOptions = RetributionPaladinOptions.create({
	classOptions: {},
});

// Defaults below are what master's ui/retribution_paladin opens with (currentSettings on a fresh
// profile); its raid-wide Battle Shout, Leader of the Pack, Moonkin and Fire Resistance Aura are
// party buffs here. Juju Power/Might, R.O.I.D.S., Dragonbreath Chili and Greater Arcane Elixir have
// no slot here; Blessed Sunfruit and Major Mana Potion are not in this db.
export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 13512, // Flask of Supreme Power
	battleElixirId: 13452, // Elixir of the Mongoose
	foodId: 13810, // Blessed Sunfruit
	potId: 13444, // Major Mana Potion
	conjuredId: 12662, // Demonic Rune
});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	divineSpirit: TristateEffect.TristateEffectRegular,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
});

export const DefaultPartyBuffs = PartyBuffs.create({
	battleShout: TristateEffect.TristateEffectImproved,
	fireResistanceAura: true,
	leaderOfThePack: TristateEffect.TristateEffectRegular,
	moonkinAura: TristateEffect.TristateEffectRegular,
	sanctityAura: TristateEffect.TristateEffectMissing,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
	blessingOfMight: TristateEffect.TristateEffectImproved,
});

export const DefaultDebuffs = Debuffs.create({
	curseOfRecklessness: true,
	faerieFire: TristateEffect.TristateEffectRegular,
	giftOfArthas: true,
	improvedSealOfTheCrusader: TristateEffect.TristateEffectImproved,
	judgementOfWisdom: true,
	sunderArmor: true,
});

export const OtherDefaults = {
	reactionTime: 200, // master's default
	profession1: Profession.Blacksmithing,
	profession2: Profession.Enchanting,
	distanceFromTarget: 5,
	race: Race.RaceHuman,
};

// Our Forever sim's gear presets (master ui/<spec>/gear_sets).
export const GEAR_LAUNCH = PresetUtils.makePresetGear('Launch', LaunchGear);
export const DEFAULT_GEAR = GEAR_LAUNCH;
export const GEAR_PRESETS = [GEAR_LAUNCH];
