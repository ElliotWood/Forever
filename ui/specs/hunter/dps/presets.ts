import * as PresetUtils from '@app/preset_utils';
import { ConsumesSpec, Debuffs, IndividualBuffs, PartyBuffs, Profession, Race, RaidBuffs, TristateEffect } from '@generated/proto/common';
import {
	Hunter_Options as HunterOptions,
	HunterOptions_Ammo,
	HunterOptions_PetAttackSpeed,
	HunterOptions_PetType as PetType,
	HunterOptions_QuiverBonus,
} from '@generated/proto/hunter';
import { SavedTalents } from '@generated/proto/ui';

import BeastMasteryAPL from './apls/bm.apl.json';
import MarksmanshipAPL from './apls/mm.apl.json';
import SurvivalAPL from './apls/sv.apl.json';
import LaunchGear from './gear_sets/launch.gear.json';
import P0BisGear from './gear_sets/p0.bis.gear.json';
import P1BisGear from './gear_sets/p1.bis.gear.json';

export const BeastMasteryRotation = PresetUtils.makePresetAPLRotation('Beast Mastery', BeastMasteryAPL);
export const MarksmanshipRotation = PresetUtils.makePresetAPLRotation('Marksmanship', MarksmanshipAPL);
export const SurvivalRotation = PresetUtils.makePresetAPLRotation('Survival', SurvivalAPL);
export const DefaultRotation = MarksmanshipRotation;

// Defaults below are what master's ui/hunter (the Forever site before the switch) opens with: its
// page drops the blessings, Fire Resistance Aura and Judgement of Wisdom its presets name.
export const DefaultOptions = HunterOptions.create({
	classOptions: {
		ammo: HunterOptions_Ammo.ThoriumHeadedArrow,
		quiverBonus: HunterOptions_QuiverBonus.Speed15,
		petType: PetType.Cat,
		petAttackSpeed: HunterOptions_PetAttackSpeed.OneTwo,
		petUptime: 1,
	},
});

export const DefaultIndividualBuffs = IndividualBuffs.create({});

export const DefaultPartyBuffs = PartyBuffs.create({
	battleShout: TristateEffect.TristateEffectImproved,
	fireResistanceTotem: true,
	graceOfAirTotem: TristateEffect.TristateEffectImproved,
	manaSpringTotem: TristateEffect.TristateEffectRegular,
	strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	divineSpirit: TristateEffect.TristateEffectRegular,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
});

export const DefaultDebuffs = Debuffs.create({
	curseOfRecklessness: true,
	exposeArmor: TristateEffect.TristateEffectImproved,
	faerieFire: TristateEffect.TristateEffectRegular,
	// Improved Hunter's Mark is gone from the Forever trees, assumed baseline rather than removed.
	huntersMark: TristateEffect.TristateEffectImproved,
	sunderArmor: true,
});

// Master's Juju Power/Might, Dragonbreath Chili, Ground Scorpok Assay and Windfury have no slot here.
export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 13512, // Flask of Supreme Power
	battleElixirId: 13452, // Elixir of the Mongoose
	guardianElixirId: 20007, // Mageblood Elixir
	foodId: 20452, // Smoked Desert Dumplings
	potId: 13444, // Major Mana Potion
	conjuredId: 12662, // Demonic Rune
	ohImbueId: 18262, // Elemental Sharpening Stone
	petScrollAgi: true,
	petScrollStr: true,
});

export const OtherDefaults = {
	reactionTime: 200, // master's default
	distanceFromTarget: 12,
	profession1: Profession.Enchanting,
	profession2: Profession.Engineering,
	race: Race.RaceTroll,
};

export const TalentsP1 = PresetUtils.makePresetTalents('Marksmanship', SavedTalents.create({ talentsString: '5023000501-0050550501503051' }));
export const TalentsBeastMastery = PresetUtils.makePresetTalents('Beast Mastery 35/16/0', SavedTalents.create({ talentsString: '5520001505121251-0050551' }));
export const TalentsMarksmanship = PresetUtils.makePresetTalents('Marksmanship 0/39/12', SavedTalents.create({ talentsString: '-3050552301503151-50024001' }));
export const TalentsSurvival = PresetUtils.makePresetTalents('Survival 0/15/36', SavedTalents.create({ talentsString: '-005055-550230031051220151' }));
export const TalentPresets = [TalentsP1, TalentsBeastMastery, TalentsMarksmanship, TalentsSurvival];

// Our Forever sim's gear presets (master ui/<spec>/gear_sets).
export const GEAR_LAUNCH = PresetUtils.makePresetGear('Launch', LaunchGear);
export const GEAR_P0_BIS = PresetUtils.makePresetGear('Pre-BiS', P0BisGear);
export const GEAR_P1_BIS = PresetUtils.makePresetGear('P1 BiS', P1BisGear);
export const DEFAULT_GEAR = GEAR_P0_BIS;
export const GEAR_PRESETS = [GEAR_LAUNCH, GEAR_P0_BIS, GEAR_P1_BIS];
