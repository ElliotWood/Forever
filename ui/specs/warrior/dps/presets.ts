import * as PresetUtils from '@app/preset_utils';
import { makeSpecChangeWarningToast } from '@features/settings/utils/spec_change_warning_toast';
import { ConsumesSpec, HandType, ItemSlot, Profession, Race, Spec } from '@generated/proto/common';
import { SavedTalents } from '@generated/proto/ui';
import { DpsWarrior_Options as WarriorOptions, DpsWarrior_Rotation, DpsWarriorSpec, WarriorShout, WarriorStance, WarriorSunder } from '@generated/proto/warrior';
import { Player } from '@sim/player/player';

import * as WarriorPresets from '../shared/presets';
import DefaultArmsApl from './apls/arms.apl.json';
import DefaultFuryApl from './apls/fury.apl.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const isArmsSpec = (player: Player<Spec.SpecDpsWarrior>) =>
	player.getEquippedItem(ItemSlot.ItemSlotMainHand)?.item.handType === HandType.HandTypeTwoHand;

export const isArmsKebabSpec = (player: Player<Spec.SpecDpsWarrior>) => player.getTalents().mortalStrike && isFurySpec(player);

export const isFurySpec = (player: Player<Spec.SpecDpsWarrior>) =>
	player.getTalents().bloodthirst ||
	player.getEquippedItem(ItemSlot.ItemSlotMainHand)?.item.handType === HandType.HandTypeMainHand ||
	player.getEquippedItem(ItemSlot.ItemSlotMainHand)?.item.handType === HandType.HandTypeOneHand;

// Handlers for spec specific load checks
const FURY_PRESET_OPTIONS = {
	group: 'Fury',
	onLoad: (player: Player<Spec.SpecDpsWarrior>) => {
		makeSpecChangeWarningToast(
			[
				{
					condition: isArmsSpec,
					message: 'Check your gear: You have a two-handed weapon equipped, but the selected option is for dual wield.',
				},
				{
					condition: (player: Player<Spec.SpecDpsWarrior>) => !player.getTalents().dualWieldSpecialization,
					message: "Check your talents: You have selected a dual-wield spec but don't have [Dual Wield Specialization] talented.",
				},
			],
			player,
		);
	},
};
const ARMS_PRESET_OPTIONS = {
	group: 'Arms',
	onLoad: (player: Player<any>) => {
		makeSpecChangeWarningToast(
			[
				{
					condition: isFurySpec,
					message: 'Check your gear: You have a one-handed weapon equipped, but the selected option is for two-handed weapons.',
				},
			],
			player,
		);
	},
};

export const FURY_DEFAULT_ROTATION = PresetUtils.makePresetAPLRotation('Fury', DefaultFuryApl);
export const ARMS_DEFAULT_ROTATION = PresetUtils.makePresetAPLRotation('Arms', DefaultArmsApl);

export const SIMPLE_ROTATION = DpsWarrior_Rotation.create({
	spec: DpsWarriorSpec.DpsWarriorSpecFury,
	sunderArmor: WarriorSunder.WarriorSunderHelp,
	useOverpower: true,
	useRecklessness: false,
	bloodlustTiming: 5,
});
export const SIMPLE_DEFAULT_ROTATION = PresetUtils.makePresetSimpleRotation('Simple', Spec.SpecDpsWarrior, SIMPLE_ROTATION);
export const SIMPLE_ARMS_DEFAULT_ROTATION = PresetUtils.makePresetSimpleRotation('Simple', Spec.SpecDpsWarrior, {
	...SIMPLE_ROTATION,
	spec: DpsWarriorSpec.DpsWarriorSpecArms,
});

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/forever/talent-calc and copy the numbers in the url.
export const FuryTalents = {
	name: 'Fury',
	data: SavedTalents.create({
		talentsString: '3400502130201-05050005505012050115',
	}),
	...FURY_PRESET_OPTIONS,
};

export const ArmsTalents = {
	name: 'Arms',
	data: SavedTalents.create({
		talentsString: '32005011352010500221-0550000500521203',
	}),
	...ARMS_PRESET_OPTIONS,
};

export const ArmsKebabTalents = {
	name: 'Arms - Kebab',
	data: SavedTalents.create({
		talentsString: '34005021302010510321-0550000520501203',
	}),
	...FURY_PRESET_OPTIONS,
};

export const DefaultOptions = WarriorOptions.create({
	classOptions: {
		queueDelay: 250,
		startingRage: 50,
		defaultShout: WarriorShout.WarriorShoutBattle,
		defaultStance: WarriorStance.WarriorStanceBerserker,
		hasBsT2: true,
		stanceSnapshot: true,
	},
});

export const DefaultConsumables = ConsumesSpec.create({
	...WarriorPresets.DefaultConsumables,
});

export const OtherDefaults = {
	race: Race.RaceOrc,
	profession1: Profession.Engineering,
	profession2: Profession.Blacksmithing,
	distanceFromTarget: 25,
};
