import * as PresetUtils from '@app/preset_utils';
import { Class, ConsumesSpec, Debuffs, Drums, IndividualBuffs, PartyBuffs, Profession, Race, RaidBuffs, TristateEffect } from '@generated/proto/common';
import {
	Hunter_Options as HunterOptions,
	HunterOptions_Ammo,
	HunterOptions_PetAttackSpeed,
	HunterOptions_PetType as PetType,
	HunterOptions_QuiverBonus,
} from '@generated/proto/hunter';
import { SavedTalents } from '@generated/proto/ui';
import { defaultExposeWeaknessSettings, defaultRaidBuffMajorDamageCooldowns } from '@sim/proto/utils';

import BeastMasteryAPL from './apls/bm.apl.json';
import MarksmanshipAPL from './apls/mm.apl.json';
import SurvivalAPL from './apls/sv.apl.json';

export const BeastMasteryRotation = PresetUtils.makePresetAPLRotation('Beast Mastery', BeastMasteryAPL);
export const MarksmanshipRotation = PresetUtils.makePresetAPLRotation('Marksmanship', MarksmanshipAPL);
export const SurvivalRotation = PresetUtils.makePresetAPLRotation('Survival', SurvivalAPL);
export const DefaultRotation = MarksmanshipRotation;

export const DefaultOptions = HunterOptions.create({
	classOptions: {
		ammo: HunterOptions_Ammo.Doomshot,
		quiverBonus: HunterOptions_QuiverBonus.Speed15,
		petType: PetType.Cat,
		petAttackSpeed: HunterOptions_PetAttackSpeed.OneTwo,
		petUptime: 1,
		petSingleAbility: false,
	},
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfMight: TristateEffect.TristateEffectImproved,
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
	unleashedRage: true,
});

export const DefaultPartyBuffs = PartyBuffs.create({
	battleShout: TristateEffect.TristateEffectImproved,
	braidedEterniumChain: true,
	ferociousInspiration: 1,
	graceOfAirTotem: TristateEffect.TristateEffectImproved,
	leaderOfThePack: TristateEffect.TristateEffectImproved,
	strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
	totemTwisting: true,
	windfuryTotem: TristateEffect.TristateEffectImproved,
	drums: Drums.LesserDrumsOfBattle,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	...defaultRaidBuffMajorDamageCooldowns(Class.ClassWarrior),
	arcaneBrilliance: true,
	divineSpirit: TristateEffect.TristateEffectImproved,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	powerWordFortitude: TristateEffect.TristateEffectImproved,
	shadowProtection: true,
});

export const DefaultDebuffs = Debuffs.create({
	bloodFrenzy: true,
	curseOfElements: TristateEffect.TristateEffectImproved,
	curseOfRecklessness: true,
	exposeArmor: TristateEffect.TristateEffectImproved,
	...defaultExposeWeaknessSettings(),
	faerieFire: TristateEffect.TristateEffectImproved,
	giftOfArthas: true,
	huntersMark: TristateEffect.TristateEffectImproved,
	improvedSealOfTheCrusader: TristateEffect.TristateEffectImproved,
	insectSwarm: true,
	judgementOfLight: true,
	judgementOfWisdom: true,
	mangle: true,
	misery: true,
	sunderArmor: true,
});

export const DefaultConsumables = ConsumesSpec.create({
	battleElixirId: 22831, // Elixir of Major Agility
	guardianElixirId: 22840, // Elixir of Major Mageblood
	foodId: 27659, // Warp Burger
	potId: 22838, // Haste Potion
	conjuredId: 12662,
	explosiveId: 30217,
	petFoodId: 33874, // Kibler's Bits
	petScrollAgi: true,
	petScrollStr: true,
	superSapper: true,
	goblinSapper: true,
	scrollAgi: true,
	scrollStr: true,
});

export const OtherDefaults = {
	distanceFromTarget: 7,
	iterationCount: 25000,
	profession1: Profession.Engineering,
	profession2: Profession.Blacksmithing,
	race: Race.RaceOrc,
};

// Talent presets, from master's ui/hunter spec.
export const TalentsP1 = PresetUtils.makePresetTalents('Marksmanship', SavedTalents.create({ talentsString: '5023000501-0050550501503051' }));
export const TalentsBeastMastery = PresetUtils.makePresetTalents('Beast Mastery 35/16/0', SavedTalents.create({ talentsString: '5520001505121251-0050551' }));
export const TalentsMarksmanship = PresetUtils.makePresetTalents('Marksmanship 0/39/12', SavedTalents.create({ talentsString: '-3050552301503151-50024001' }));
export const TalentsSurvival = PresetUtils.makePresetTalents('Survival 0/15/36', SavedTalents.create({ talentsString: '-005055-550230031051220151' }));
export const TalentPresets = [TalentsP1, TalentsBeastMastery, TalentsMarksmanship, TalentsSurvival];
