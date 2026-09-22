import * as PresetUtils from '@app/preset_utils';
import { ConsumesSpec, Debuffs, Drums, IndividualBuffs, PartyBuffs, Profession, RaidBuffs, Spec, TristateEffect } from '@generated/proto/common';
import { Mage_Options as MageOptions, Mage_Rotation, MageArmor } from '@generated/proto/mage';
import { SavedTalents } from '@generated/proto/ui';
import { defaultImprovedShadowBoltSettings } from '@sim/proto/utils';

import ArcaneApl from './apls/arcane.apl.json';
import BlankAPL from './apls/blank.apl.json';
import FireApl from './apls/fire.apl.json';
import FrostApl from './apls/frost.apl.json';
import LaunchGear from './gear_sets/launch.gear.json';
import P0BisGear from './gear_sets/p0.bis.gear.json';
import P1BisGear from './gear_sets/p1.bis.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const BLANK_APL = PresetUtils.makePresetAPLRotation('Blank', BlankAPL);

export const ROTATION_PRESET_ARCANE = PresetUtils.makePresetAPLRotation('Arcane', ArcaneApl);
export const ROTATION_PRESET_FIRE = PresetUtils.makePresetAPLRotation('Fire', FireApl);
export const ROTATION_PRESET_FROST = PresetUtils.makePresetAPLRotation('Frost', FrostApl);

export const ArcaneMageSimpleRotation = Mage_Rotation.create({
	conserveStart: 20,
	conserveEnd: 30,
	delayMajorCDs: 10,
});

export const APL_ARCANE_SIMPLE = PresetUtils.makePresetSimpleRotation('Simple', Spec.SpecMage, ArcaneMageSimpleRotation);

// The community builds our Forever sim ranks: Arcane 35/0/16, Fire 0/35/16 and Frost 14/0/37.
export const ArcaneTalents = PresetUtils.makePresetTalents('Arcane', SavedTalents.create({ talentsString: '055005023100311531--005500033' }));
export const FireTalents = PresetUtils.makePresetTalents('Fire', SavedTalents.create({ talentsString: '-03552020130133151-005500033' }));
export const FrostTalents = PresetUtils.makePresetTalents('Frost', SavedTalents.create({ talentsString: '050005013--0555003301001301251' }));

export const DefaultOptions = MageOptions.create({
	classOptions: {
		defaultMageArmor: MageArmor.MageArmorMageArmor,
	},
});

export const OtherDefaults = {
	distanceFromTarget: 20,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
};

export const DefaultConsumables = ConsumesSpec.create({
	guardianElixirId: 32067, // Elixir of Draenic Wisdom
	battleElixirId: 28103, // Adept's Elixir
	foodId: 27657, // Blackened Basilisk
	mhImbueId: 25122, // Brilliant Wizard Oil
	potId: 22832, // Super Mana Potion
});

export const DefaultRaidBuffs = RaidBuffs.create({
	bloodlust: true,
	divineSpirit: 2,
	arcaneBrilliance: true,
	giftOfTheWild: 2,
	powerWordFortitude: 2,
	shadowProtection: true,
});

export const DefaultPartyBuffs = PartyBuffs.create({
	manaSpringTotem: 2,
	manaTideTotems: 1,
	wrathOfAirTotem: 1,
	drums: Drums.LesserDrumsOfBattle,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfWisdom: 2,
	innervates: 1,
	blessingOfSalvation: true,
	shadowPriestDps: 1400,
});

export const DefaultDebuffs = Debuffs.create({
	misery: true,
	curseOfElements: 2,
	improvedSealOfTheCrusader: TristateEffect.TristateEffectImproved,
	judgementOfWisdom: true,
	...defaultImprovedShadowBoltSettings(),
});

// Our Forever sim's gear presets (master ui/<spec>/gear_sets).
export const GEAR_LAUNCH = PresetUtils.makePresetGear('Launch', LaunchGear);
export const GEAR_P0_BIS = PresetUtils.makePresetGear('Pre-BiS', P0BisGear);
export const GEAR_P1_BIS = PresetUtils.makePresetGear('P1 BiS', P1BisGear);
export const DEFAULT_GEAR = GEAR_P0_BIS;
export const GEAR_PRESETS = [GEAR_LAUNCH, GEAR_P0_BIS, GEAR_P1_BIS];
