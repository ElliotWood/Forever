import * as PresetUtils from '@app/preset_utils';
import { ConsumesSpec, Debuffs, Drums, IndividualBuffs, PartyBuffs, Profession, Race, RaidBuffs, TristateEffect } from '@generated/proto/common';
import { RetributionPaladin_Options as RetributionPaladinOptions } from '@generated/proto/paladin';
import { SavedTalents } from '@generated/proto/ui';
import { defaultExposeWeaknessSettings } from '@sim/proto/utils';

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

export const DefaultConsumables = ConsumesSpec.create({
	potId: 22838,
	flaskId: 22854,
	foodId: 27658,
	conjuredId: 12662,
	superSapper: true,
	goblinSapper: true,
	scrollAgi: true,
	scrollStr: true,
	explosiveId: 30217,
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
	leaderOfThePack: TristateEffect.TristateEffectImproved,
	battleShout: TristateEffect.TristateEffectImproved,
	strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
	totemTwisting: true,
	windfuryTotem: TristateEffect.TristateEffectImproved,
	graceOfAirTotem: TristateEffect.TristateEffectImproved,
	drums: Drums.LesserDrumsOfBattle,
	sanctityAura: TristateEffect.TristateEffectMissing,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
	blessingOfMight: TristateEffect.TristateEffectImproved,
	unleashedRage: true,
});

export const DefaultDebuffs = Debuffs.create({
	misery: true,
	curseOfElements: TristateEffect.TristateEffectImproved,
	improvedSealOfTheCrusader: TristateEffect.TristateEffectImproved,
	jocRetribution2Pt4: true,
	judgementOfWisdom: true,
	bloodFrenzy: true,
	huntersMark: TristateEffect.TristateEffectImproved,
	curseOfRecklessness: true,
	sunderArmor: true,
	faerieFire: TristateEffect.TristateEffectImproved,
	exposeArmor: TristateEffect.TristateEffectImproved,
	...defaultExposeWeaknessSettings(),
});

export const OtherDefaults = {
	profession1: Profession.Engineering,
	profession2: Profession.Blacksmithing,
	distanceFromTarget: 5,
	iterationCount: 25000,
	race: Race.RaceBloodElf,
};

// Our Forever sim's gear presets (master ui/<spec>/gear_sets).
export const GEAR_LAUNCH = PresetUtils.makePresetGear('Launch', LaunchGear);
export const DEFAULT_GEAR = GEAR_LAUNCH;
export const GEAR_PRESETS = [GEAR_LAUNCH];
