import * as PresetUtils from '@app/preset_utils';
import { SavedTalents } from '@generated/proto/ui';
import {
	Class,
	ConsumesSpec,
	Debuffs,
	Drums,
	IndividualBuffs,
	PartyBuffs,
	Profession,
	Race,
	RaidBuffs,
	TristateEffect,
	UnitReference,
} from '@generated/proto/common';
import { BalanceDruid_Options as BalanceDruidOptions } from '@generated/proto/druid';
import { defaultRaidBuffMajorDamageCooldowns } from '@sim/proto/utils';

import DefaultAPL from './apls/default.apl.json';
import LaunchAPL from './apls/launch.apl.json';
import LaunchGear from './gear_sets/launch.gear.json';
import P0BisGear from './gear_sets/p0.bis.gear.json';
import P1BisGear from './gear_sets/p1.bis.gear.json';
import P2BisGear from './gear_sets/p2.bis.gear.json';

export const StandardRotation = PresetUtils.makePresetAPLRotation('Default', DefaultAPL);
// Master's Launch rotation (Wrath-led, Starfire on Eclipse), the one its arena ranks.
export const LaunchRotation = PresetUtils.makePresetAPLRotation('Launch', LaunchAPL);

export const BalanceTalents = PresetUtils.makePresetTalents('Balance', SavedTalents.create({ talentsString: '5532220115001351--505302' }));
export const MoonkinTalents = PresetUtils.makePresetTalents('Moonkin 38/0/13', SavedTalents.create({ talentsString: '5502220115501351--055003' }));

export const DefaultOptions = BalanceDruidOptions.create({
	classOptions: {
		innervateTarget: UnitReference.create(),
	},
});

export const DefaultRaidBuffs = RaidBuffs.create({
	...defaultRaidBuffMajorDamageCooldowns(Class.ClassShaman),
	arcaneBrilliance: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	powerWordFortitude: TristateEffect.TristateEffectImproved,
	divineSpirit: TristateEffect.TristateEffectImproved,
});

export const DefaultPartyBuffs = PartyBuffs.create({
	chainOfTheTwilightOwl: true,
	draeneiRacialCaster: true,
	drums: Drums.LesserDrumsOfBattle,
	eyeOfTheNight: true,
	totemOfWrath: 1,
	wrathOfAirTotem: TristateEffect.TristateEffectImproved,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
	shadowPriestDps: 800,
});

export const DefaultDebuffs = Debuffs.create({
	bloodFrenzy: true,
	curseOfElements: TristateEffect.TristateEffectImproved,
	curseOfRecklessness: true,
	exposeArmor: TristateEffect.TristateEffectImproved,
	giftOfArthas: true,
	huntersMark: TristateEffect.TristateEffectImproved,
	improvedSealOfTheCrusader: TristateEffect.TristateEffectImproved,
	judgementOfWisdom: true,
	mangle: true,
	misery: true,
	sunderArmor: true,
});

export const DefaultConsumables = ConsumesSpec.create({
	conjuredId: 12662, // Demonic Rune
	drumsId: Drums.LesserDrumsOfBattle,
	flaskId: 22861, // Flask of Blinding Light
	foodId: 27657, // Blackened Basilisk
	mhImbueId: 25122, // Brilliant Wizard Oil
	potId: 22832, // Super Mana Potion
});

export const OtherDefaults = {
	distanceFromTarget: 20,
	profession1: Profession.Enchanting,
	profession2: Profession.Tailoring,
	race: Race.RaceNightElf,
};

// Our Forever sim's gear presets (master ui/<spec>/gear_sets).
export const GEAR_LAUNCH = PresetUtils.makePresetGear('Launch', LaunchGear);
export const GEAR_P0_BIS = PresetUtils.makePresetGear('Pre-BiS', P0BisGear);
export const GEAR_P1_BIS = PresetUtils.makePresetGear('P1 BiS', P1BisGear);
export const GEAR_P2_BIS = PresetUtils.makePresetGear('P2 BiS', P2BisGear);
export const DEFAULT_GEAR = GEAR_P0_BIS;
export const GEAR_PRESETS = [GEAR_LAUNCH, GEAR_P0_BIS, GEAR_P1_BIS, GEAR_P2_BIS];
