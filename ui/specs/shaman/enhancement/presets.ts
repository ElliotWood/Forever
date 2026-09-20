import * as PresetUtils from '@app/preset_utils';
import { Class, ConsumesSpec, Debuffs, Drums, IndividualBuffs, PartyBuffs, Profession, Race, RaidBuffs, TristateEffect } from '@generated/proto/common';
import { EnhancementShaman_Options as EnhancementShamanOptions, ShamanImbue, ShamanSyncType } from '@generated/proto/shaman';
import { SavedTalents } from '@generated/proto/ui';
import { defaultExposeWeaknessSettings, defaultRaidBuffMajorDamageCooldowns } from '@sim/proto/utils';

import DefaultApl from './apls/default.apl.json';

export const ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Default', DefaultApl);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/forever/talent-calc and copy the numbers in the url.
export const SubRestoIWT = {
	name: 'Sub-Restoration IWT',
	data: SavedTalents.create({
		talentsString: '03-500502210501133531151-50005301',
	}),
};

export const SubRestoILS = {
	name: 'Sub-Restoration ILS',
	data: SavedTalents.create({
		talentsString: '03-500503210500133531151-50005301',
	}),
};

export const SubEle = {
	name: 'Sub-Elemental',
	data: SavedTalents.create({
		talentsString: '250031501-500503210500133531151',
	}),
};

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
