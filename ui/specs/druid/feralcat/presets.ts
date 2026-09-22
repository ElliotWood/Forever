import * as PresetUtils from '@app/preset_utils';
import { SavedTalents } from '@generated/proto/ui';
import { ConsumesSpec, Drums, Profession, Race, Spec } from '@generated/proto/common';
import {
	FeralCatDruid_Options as FeralDruidOptions,
	FeralCatDruid_Rotation as FeralCatDruidRotation,
	FeralCatDruid_Rotation_FinishingMove as FinishingMove,
} from '@generated/proto/druid';

import DefaultApl from './apls/default.apl.json';
import LaunchGear from './gear_sets/launch.gear.json';
import P0BisGear from './gear_sets/p0.bis.gear.json';
import P2PreBisGear from './gear_sets/p2.pre-bis.gear.json';
import P2BisGear from './gear_sets/p2.bis.gear.json';

export const DefaultOptions = FeralDruidOptions.create({});

export const DefaultConsumables = ConsumesSpec.create({
	potId: 22838, // Haste Potion
	battleElixirId: 22831, // Elixir of Major Agility
	guardianElixirId: 32067, // Elixir of Draenic Wisdom
	foodId: 27664, // Grilled Mudfish (+20 Agility)
	mhImbueId: 34340, // Adamantite Weightstone
	conjuredId: 12662, // Demonic Rune
	drumsId: Drums.GreaterDrumsOfBattle,
	superSapper: true,
	goblinSapper: true,
	scrollAgi: true,
	scrollStr: true,
});

export const OtherDefaults = {
	distanceFromTarget: 0,
	profession1: Profession.Engineering,
	profession2: Profession.Enchanting,
	race: Race.RaceNightElf,
	reactionTime: 250,
};

export const DefaultRotation = FeralCatDruidRotation.create({
	finishingMove: FinishingMove.Rip,
	biteweave: true,
	ripMinComboPoints: 5,
	biteMinComboPoints: 5,
	mangleTrick: true,
	maintainFaerieFire: true,
});

export const SIMPLE = PresetUtils.makePresetSimpleRotation('Simple', Spec.SpecFeralCatDruid, DefaultRotation);

export const APL = PresetUtils.makePresetAPLRotation('APL', DefaultApl);

export const FeralTalents = PresetUtils.makePresetTalents('Feral', SavedTalents.create({ talentsString: '-5521002023132213051-05503' }));
export const FeralCatTalents = PresetUtils.makePresetTalents('Feral Cat 9/35/7', SavedTalents.create({ talentsString: '050022-5500002123032213051-052' }));

// Our Forever sim's gear presets (master ui/<spec>/gear_sets).
export const GEAR_LAUNCH = PresetUtils.makePresetGear('Launch', LaunchGear);
export const GEAR_P0_BIS = PresetUtils.makePresetGear('Pre-BiS', P0BisGear);
export const GEAR_P2_PRE_BIS = PresetUtils.makePresetGear('P2 Pre-BiS', P2PreBisGear);
export const GEAR_P2_BIS = PresetUtils.makePresetGear('P2 BiS', P2BisGear);
export const DEFAULT_GEAR = GEAR_P0_BIS;
export const GEAR_PRESETS = [GEAR_LAUNCH, GEAR_P0_BIS, GEAR_P2_PRE_BIS, GEAR_P2_BIS];
