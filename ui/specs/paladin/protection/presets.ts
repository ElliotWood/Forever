import * as PresetUtils from '@app/preset_utils';
import { ConsumesSpec, Debuffs, Drums, HealingModel, IndividualBuffs, PartyBuffs, Profession, RaidBuffs, TristateEffect } from '@generated/proto/common';
import { ProtectionPaladin_Options as ProtectionPaladinOptions } from '@generated/proto/paladin';
import { SavedTalents } from '@generated/proto/ui';
import { defaultExposeWeaknessSettings } from '@sim/proto/utils';

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

export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 22861, // Flask of Blinding Light
	foodId: 27657, // Blackened Basilisk
	potId: 22849, // Ironshield Potion
	conjuredId: 12662, // Dark Rune
	mhImbueId: 28017,
	explosiveId: 30217,
	superSapper: true,
	goblinSapper: true,
	nightmareSeed: true,
	scrollStr: true,
	scrollAgi: true,
	scrollArm: true,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	bloodlust: true,
	divineSpirit: TristateEffect.TristateEffectImproved,
	arcaneBrilliance: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	powerWordFortitude: TristateEffect.TristateEffectImproved,
	shadowProtection: true,
	thorns: TristateEffect.TristateEffectImproved,
});

export const DefaultPartyBuffs = PartyBuffs.create({
	manaSpringTotem: TristateEffect.TristateEffectRegular,
	wrathOfAirTotem: TristateEffect.TristateEffectRegular,
	graceOfAirTotem: TristateEffect.TristateEffectMissing,
	strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
	windfuryTotem: TristateEffect.TristateEffectMissing,
	battleShout: TristateEffect.TristateEffectMissing,
	drums: Drums.LesserDrumsOfBattle,
	sanctityAura: TristateEffect.TristateEffectMissing,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
	blessingOfMight: TristateEffect.TristateEffectImproved,
	blessingOfSanctuary: true,
});

export const DefaultDebuffs = Debuffs.create({
	misery: true,
	curseOfElements: TristateEffect.TristateEffectImproved,
	improvedSealOfTheCrusader: TristateEffect.TristateEffectImproved,
	judgementOfWisdom: true,
	judgementOfLight: true,
	bloodFrenzy: true,
	huntersMark: TristateEffect.TristateEffectImproved,
	curseOfRecklessness: true,
	sunderArmor: true,
	faerieFire: TristateEffect.TristateEffectImproved,
	exposeArmor: TristateEffect.TristateEffectImproved,
	insectSwarm: true,
	...defaultExposeWeaknessSettings(),
});

export const OtherDefaults = {
	profession1: Profession.Engineering,
	profession2: Profession.Enchanting,
	distanceFromTarget: 5,
	iterationCount: 25000,
	healingModel: HealingModel.create({
		hps: 2200,
		cadenceSeconds: 0.4,
		cadenceVariation: 1.2,
		absorbFrac: 0.02,
		burstWindow: 6,
		inspirationUptime: 0.25,
	}),
};

// Our Forever sim's gear presets (master ui/<spec>/gear_sets).
export const GEAR_LAUNCH = PresetUtils.makePresetGear('Launch', LaunchGear);
export const DEFAULT_GEAR = GEAR_LAUNCH;
export const GEAR_PRESETS = [GEAR_LAUNCH];
