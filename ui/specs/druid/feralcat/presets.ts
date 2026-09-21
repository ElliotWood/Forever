import * as PresetUtils from '@app/preset_utils';
import { SavedTalents } from '@generated/proto/ui';
import { ConsumesSpec, Drums, Profession, Race, Spec } from '@generated/proto/common';
import {
	FeralCatDruid_Options as FeralDruidOptions,
	FeralCatDruid_Rotation as FeralCatDruidRotation,
	FeralCatDruid_Rotation_FinishingMove as FinishingMove,
} from '@generated/proto/druid';

import DefaultApl from './apls/default.apl.json';

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
