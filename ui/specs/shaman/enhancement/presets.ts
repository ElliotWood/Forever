import * as PresetUtils from '@app/preset_utils';
import { Class, ConsumesSpec, Debuffs, Drums, IndividualBuffs, PartyBuffs, Profession, Race, RaidBuffs, TristateEffect } from '@generated/proto/common';
import { EnhancementShaman_Options as EnhancementShamanOptions, ShamanImbue, ShamanSyncType } from '@generated/proto/shaman';
import { defaultExposeWeaknessSettings, defaultRaidBuffMajorDamageCooldowns } from '@sim/proto/utils';

import DefaultApl from './apls/default.apl.json';
import { SavedTalents } from '@generated/proto/ui';
import LaunchGear from './gear_sets/launch.gear.json';
import Phase1Gear from './gear_sets/phase_1.gear.json';
import Phase2Gear from './gear_sets/phase_2.gear.json';

export const ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Default', DefaultApl);

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfMight: TristateEffect.TristateEffectImproved,
});

export const DefaultOptions = EnhancementShamanOptions.create({
	classOptions: {
		shieldProcrate: 0,
		imbueMh: ShamanImbue.WindfuryWeapon,
	},
	imbueOh: ShamanImbue.WindfuryWeapon,
	syncType: ShamanSyncType.DelayOffhandSwings,
});

export const OtherDefaults = {
	distanceFromTarget: 5,
	profession1: Profession.Engineering,
	profession2: Profession.Leatherworking,
	race: Race.RaceOrc,
};

export const DefaultConsumables = ConsumesSpec.create({
	potId: 22838, // Haste Potion
	flaskId: 22854, // Flask of Relentless Assault
	foodId: 27658, // Roasted Clefthoof
	drumsId: Drums.LesserDrumsOfBattle,
	conjuredId: 22788,
	explosiveId: 30217,
	superSapper: true,
	goblinSapper: true,
	scrollAgi: true,
	scrollStr: true,
});

export const DefaultPartyBuffs = PartyBuffs.create({
	ferociousInspiration: 2,
	braidedEterniumChain: true,
	leaderOfThePack: TristateEffect.TristateEffectRegular,
	battleShout: TristateEffect.TristateEffectImproved,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	...defaultRaidBuffMajorDamageCooldowns(Class.ClassShaman),
	powerWordFortitude: TristateEffect.TristateEffectImproved,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	arcaneBrilliance: true,
});

export const DefaultDebuffs = Debuffs.create({
	...defaultExposeWeaknessSettings(),
	improvedSealOfTheCrusader: TristateEffect.TristateEffectImproved,
	judgementOfWisdom: true,
	screech: true,
	misery: true,
	bloodFrenzy: true,
	giftOfArthas: true,
	mangle: true,
	exposeArmor: TristateEffect.TristateEffectImproved,
	faerieFire: TristateEffect.TristateEffectImproved,
	sunderArmor: true,
	curseOfElements: TristateEffect.TristateEffectImproved,
	curseOfRecklessness: true,
	huntersMark: TristateEffect.TristateEffectImproved,
});

// Talent presets, from master's ui/shaman spec.
export const TalentsLevel60 = PresetUtils.makePresetTalents('Level 60', SavedTalents.create({ talentsString: '5505301-053030031005112251' }));
export const TalentsEnhancement = PresetUtils.makePresetTalents('Enhancement 16/35/0', SavedTalents.create({ talentsString: '05023015-055030030205112251' }));
export const TalentPresets = [TalentsLevel60, TalentsEnhancement];

// Our Forever sim's gear presets (master ui/<spec>/gear_sets).
export const GEAR_LAUNCH = PresetUtils.makePresetGear('Launch', LaunchGear);
export const GEAR_PHASE_1 = PresetUtils.makePresetGear('Phase 1', Phase1Gear);
export const GEAR_PHASE_2 = PresetUtils.makePresetGear('Phase 2', Phase2Gear);
export const DEFAULT_GEAR = GEAR_LAUNCH;
export const GEAR_PRESETS = [GEAR_LAUNCH, GEAR_PHASE_1, GEAR_PHASE_2];
