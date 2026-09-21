import * as PresetUtils from '@app/preset_utils';
import { ConsumesSpec } from '@generated/proto/common';
import { Rogue_Options as RogueOptions } from '@generated/proto/rogue';
import { SavedTalents } from '@generated/proto/ui';

import CombatSinisterStrikeAPL from './apls/combat_sinister_strike.apl.json';
import HemorrhageAPL from './apls/forever_hemorrhage.apl.json';
import MutilateAPL from './apls/forever_mutilate.apl.json';

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
