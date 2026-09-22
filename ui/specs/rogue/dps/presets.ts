import * as PresetUtils from '@app/preset_utils';
import { ConsumesSpec } from '@generated/proto/common';
import { Rogue_Options as RogueOptions } from '@generated/proto/rogue';
import { SavedTalents } from '@generated/proto/ui';

import CombatSinisterStrikeAPL from './apls/combat_sinister_strike.apl.json';
import HemorrhageAPL from './apls/forever_hemorrhage.apl.json';
import MutilateAPL from './apls/forever_mutilate.apl.json';
import BackstabLaunchGear from './gear_sets/backstab_launch.gear.json';
import SinisterStrikeLaunchGear from './gear_sets/sinister_strike_launch.gear.json';
import CombatBackstabPrebisGear from './gear_sets/combat_backstab_prebis.gear.json';
import CombatSinisterStrikePrebisGear from './gear_sets/combat_sinister_strike_prebis.gear.json';
import CombatBackstabP1BisGear from './gear_sets/combat_backstab_p1_bis.gear.json';
import CombatSinisterStrikeP1BisGear from './gear_sets/combat_sinister_strike_p1_bis.gear.json';
import CombatBackstabP2BisGear from './gear_sets/combat_backstab_p2_bis.gear.json';
import CombatSinisterStrikeP2BisGear from './gear_sets/combat_sinister_strike_p2_bis.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const ROTATION_PRESET_COMBAT = PresetUtils.makePresetAPLRotation('Combat (Sinister Strike)', CombatSinisterStrikeAPL);
export const ROTATION_PRESET_MUTILATE = PresetUtils.makePresetAPLRotation('Assassination (Mutilate)', MutilateAPL);
export const ROTATION_PRESET_HEMORRHAGE = PresetUtils.makePresetAPLRotation('Subtlety (Hemorrhage)', HemorrhageAPL);

// The three builds our Forever sim ranks.
export const CombatTalents = PresetUtils.makePresetTalents('Combat (Sinister Strike)', SavedTalents.create({ talentsString: '00530310501-32003311201515231' }));
export const AssassinationTalents = PresetUtils.makePresetTalents(
	'Assassination (Mutilate)',
	SavedTalents.create({ talentsString: '00530310551021051-302303202004' }),
);
export const SubtletyTalents = PresetUtils.makePresetTalents('Subtlety (Hemorrhage)', SavedTalents.create({ talentsString: '125320101--5320003310013211551' }));

export const DefaultOptions = RogueOptions.create({
	classOptions: {},
});

export const DefaultConsumables = ConsumesSpec.create({
	battleElixirId: 22831,
	guardianElixirId: 32062,
	foodId: 33872,
	potId: 22838,
	conjuredId: 7676,
	mhImbueId: 26891, // Instant Poison
	ohImbueId: 27186, // Deadly Poison
});

export const OtherDefaults = {
	distanceFromTarget: 5,
};

// Our Forever sim's gear presets (master ui/<spec>/gear_sets).
export const GEAR_BACKSTAB_LAUNCH = PresetUtils.makePresetGear('Backstab Launch', BackstabLaunchGear);
export const GEAR_SINISTER_STRIKE_LAUNCH = PresetUtils.makePresetGear('Sinister Strike Launch', SinisterStrikeLaunchGear);
export const GEAR_COMBAT_BACKSTAB_PREBIS = PresetUtils.makePresetGear('Backstab Pre-BiS', CombatBackstabPrebisGear);
export const GEAR_COMBAT_SINISTER_STRIKE_PREBIS = PresetUtils.makePresetGear('Sinister Strike Pre-BiS', CombatSinisterStrikePrebisGear);
export const GEAR_COMBAT_BACKSTAB_P1_BIS = PresetUtils.makePresetGear('Backstab P1 BiS', CombatBackstabP1BisGear);
export const GEAR_COMBAT_SINISTER_STRIKE_P1_BIS = PresetUtils.makePresetGear('Sinister Strike P1 BiS', CombatSinisterStrikeP1BisGear);
export const GEAR_COMBAT_BACKSTAB_P2_BIS = PresetUtils.makePresetGear('Backstab P2 BiS', CombatBackstabP2BisGear);
export const GEAR_COMBAT_SINISTER_STRIKE_P2_BIS = PresetUtils.makePresetGear('Sinister Strike P2 BiS', CombatSinisterStrikeP2BisGear);
export const DEFAULT_GEAR = GEAR_COMBAT_SINISTER_STRIKE_PREBIS;
export const GEAR_PRESETS = [GEAR_BACKSTAB_LAUNCH, GEAR_SINISTER_STRIKE_LAUNCH, GEAR_COMBAT_BACKSTAB_PREBIS, GEAR_COMBAT_SINISTER_STRIKE_PREBIS, GEAR_COMBAT_BACKSTAB_P1_BIS, GEAR_COMBAT_SINISTER_STRIKE_P1_BIS, GEAR_COMBAT_BACKSTAB_P2_BIS, GEAR_COMBAT_SINISTER_STRIKE_P2_BIS];
