import * as PresetUtils from '@app/preset_utils';
import { ConsumesSpec, Debuffs, Drums, IndividualBuffs, PartyBuffs, Profession, RaidBuffs, TristateEffect } from '@generated/proto/common';
import { DpsPriest_Options as Options } from '@generated/proto/priest';
import { SavedTalents } from '@generated/proto/ui';
import { defaultImprovedShadowBoltSettings, defaultRaidBuffMajorDamageCooldowns } from '@sim/proto/utils';

import DefaultApl from './apls/default.apl.json';

export const ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Default', DefaultApl);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://www.wowhead.com/forever/talent-calc/priest and copy the numbers in the url.
export const StandardTalents = {
	name: 'Shadow',
	data: SavedTalents.create({
		talentsString: '',
	}),
};

export const DefaultOptions = Options.create({
	classOptions: {
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
