import * as PresetUtils from '@app/preset_utils';
import { ConsumesSpec, Debuffs, Drums, IndividualBuffs, PartyBuffs, Profession, RaidBuffs, TristateEffect } from '@generated/proto/common';
import { DpsPriest_Options as Options, PriestOptions_Armor } from '@generated/proto/priest';
import { SavedTalents } from '@generated/proto/ui';
import { defaultImprovedShadowBoltSettings, defaultRaidBuffMajorDamageCooldowns } from '@sim/proto/utils';

import ShadowApl from './apls/shadow.apl.json';
import SmiteApl from './apls/smite.apl.json';
import LaunchGear from './gear_sets/launch.gear.json';
import P0BisGear from './gear_sets/p0.bis.gear.json';
import P1BisGear from './gear_sets/p1.bis.gear.json';
import SmiteLaunchGear from './gear_sets/smite_launch.gear.json';

export const ROTATION_PRESET_SHADOW = PresetUtils.makePresetAPLRotation('Shadow', ShadowApl);
export const ROTATION_PRESET_SMITE = PresetUtils.makePresetAPLRotation('Smite', SmiteApl);

// The community builds the rankings page runs on our Forever sim: Shadow 15/0/36 and Smite 31/17/3.
export const ShadowTalents = PresetUtils.makePresetTalents('Shadow 15/0/36', SavedTalents.create({ talentsString: '0253000311--550022501201302251' }));
export const SmiteTalents = PresetUtils.makePresetTalents('Smite 31/17/3', SavedTalents.create({ talentsString: '515030031305001031-00505023002-003' }));

export const DefaultOptions = Options.create({
	classOptions: {
		armor: PriestOptions_Armor.InnerFire,
		preShadowform: true,
	},
});

export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 22866, // Flask of Pure Death
	foodId: 27657, // Blackened Basilisk
	conjuredId: 12662, // Demonic Rune
	mhImbueId: 22522, // Superior Wizard Oil
	potId: 22839, // Destruction Potion
	explosiveId: 30217,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	...defaultRaidBuffMajorDamageCooldowns(),
	arcaneBrilliance: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	powerWordFortitude: TristateEffect.TristateEffectImproved,
	divineSpirit: TristateEffect.TristateEffectImproved,
});

export const DefaultPartyBuffs = PartyBuffs.create({
	manaSpringTotem: TristateEffect.TristateEffectRegular,
	wrathOfAirTotem: TristateEffect.TristateEffectImproved,
	eyeOfTheNight: true,
	chainOfTheTwilightOwl: true,
	drums: Drums.LesserDrumsOfBattle,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
	shadowPriestDps: 0,
});

export const DefaultDebuffs = Debuffs.create({
	improvedSealOfTheCrusader: TristateEffect.TristateEffectImproved,
	judgementOfWisdom: true,
	misery: false,
	shadowWeaving: false,
	faerieFire: TristateEffect.TristateEffectImproved,
	shadowEmbrace: true,
	curseOfElements: TristateEffect.TristateEffectImproved,
	exposeArmor: TristateEffect.TristateEffectImproved,
	...defaultImprovedShadowBoltSettings(),
});

export const OtherDefaults = {
	channelClipDelay: 100,
	distanceFromTarget: 28,
	profession1: Profession.Enchanting,
	profession2: Profession.Tailoring,
};

// Our Forever sim's gear presets (master ui/<spec>/gear_sets).
export const GEAR_LAUNCH = PresetUtils.makePresetGear('Launch', LaunchGear);
export const GEAR_P0_BIS = PresetUtils.makePresetGear('Pre-BiS', P0BisGear);
export const GEAR_P1_BIS = PresetUtils.makePresetGear('P1 BiS', P1BisGear);
export const GEAR_SMITE_LAUNCH = PresetUtils.makePresetGear('Smite Launch', SmiteLaunchGear);
export const DEFAULT_GEAR = GEAR_P0_BIS;
export const GEAR_PRESETS = [GEAR_LAUNCH, GEAR_P0_BIS, GEAR_P1_BIS, GEAR_SMITE_LAUNCH];
