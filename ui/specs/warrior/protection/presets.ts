import * as PresetUtils from '@app/preset_utils';
import { ConsumesSpec, HealingModel, Profession, Race } from '@generated/proto/common';
import { SavedTalents } from '@generated/proto/ui';
import { ProtectionWarrior_Options as ProtectionWarriorOptions, WarriorShout, WarriorStance } from '@generated/proto/warrior';
import { OtherDefaults as SimUIOtherDefaults } from '@sim/spec_config';

import * as WarriorPresets from '../shared/presets';
import GenericApl from './apls/default.apl.json';
import ForeverProtectionApl from './apls/protection.apl.json';
import LaunchGear from './gear_sets/launch.gear.json';
import P0BisGear from './gear_sets/p0.bis.gear.json';
import P1BisGear from './gear_sets/p1.bis.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const ROTATION_DEFAULT = PresetUtils.makePresetAPLRotation('Generic', GenericApl);
export const ROTATION_PRESET_PROTECTION = PresetUtils.makePresetAPLRotation('Protection', ForeverProtectionApl);

// The two builds our Forever sim ships for the tank.
export const ProtectionTalents = PresetUtils.makePresetTalents('Protection', SavedTalents.create({ talentsString: '31--552531233330012531' }));
export const DeepProtectionTalents = PresetUtils.makePresetTalents('Protection 1/0/50', SavedTalents.create({ talentsString: '1--552531233331212531' }));

export const DefaultOptions = ProtectionWarriorOptions.create({
	classOptions: {
		queueDelay: 250,
		startingRage: 100,
		defaultShout: WarriorShout.WarriorShoutBattle,
		defaultStance: WarriorStance.WarriorStanceDefensive,
		hasBsT2: true,
		stanceSnapshot: true,
	},
});

export const DefaultConsumables = ConsumesSpec.create({
	...WarriorPresets.DefaultConsumables,
	conjuredId: 22105,
	foodId: 27667,
	flaskId: undefined,
	battleElixirId: 22831,
	guardianElixirId: 9088,
	potId: 22849,
	nightmareSeed: true,
	scrollStr: true,
	scrollAgi: true,
	scrollArm: true,
});

export const OtherDefaults: Partial<SimUIOtherDefaults> = {
	profession1: Profession.Engineering,
	profession2: Profession.Blacksmithing,
	race: Race.RaceOrc,
	distanceFromTarget: 0,
	healingModel: HealingModel.create({
		hps: 2200,
		cadenceSeconds: 0.4,
		cadenceVariation: 1.2,
		absorbFrac: 0.02,
		burstWindow: 6,
		inspirationUptime: 0.25,
	}),
	// Morogrim
	// healingModel: HealingModel.create({
	// 	hps: 3300,
	// 	cadenceSeconds: 1.5,
	// 	cadenceVariation: 1.0,
	// 	absorbFrac: 0.02,
	// 	burstWindow: 6,
	// 	inspirationUptime: 0.12,
	// }),
};

// Our Forever sim's gear presets (master ui/<spec>/gear_sets).
export const GEAR_LAUNCH = PresetUtils.makePresetGear('Launch', LaunchGear);
export const GEAR_P0_BIS = PresetUtils.makePresetGear('Pre-BiS', P0BisGear);
export const GEAR_P1_BIS = PresetUtils.makePresetGear('P1 BiS', P1BisGear);
export const DEFAULT_GEAR = GEAR_LAUNCH;
export const GEAR_PRESETS = [GEAR_LAUNCH, GEAR_P0_BIS, GEAR_P1_BIS];
