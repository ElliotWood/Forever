import { PseudoStat, Spec } from '@generated/proto/common';
import { Stats } from '@sim/proto/stats';
import { defineSpec } from '@sim/spec_config';

import * as Presets from '../dps/presets';
import ShadowPage from '../dps/spec';

// Master ships Smite as its own sim page (ui/smite_priest). Here both pages sim the one dps priest
// spec; this one opens with master's Smite defaults and keeps its own saved settings.
export default defineSpec<Spec.SpecDpsPriest>({
	...ShadowPage,
	storageKeyPart: '_smite',
	pageTitle: 'Smite Priest',
	defaults: {
		...ShadowPage.defaults,
		gear: Presets.GEAR_SMITE_LAUNCH.gear,
		statCaps: new Stats().withPseudoStat(PseudoStat.PseudoStatSchoolHitPercentHoly, 16),
		consumables: Presets.SmiteConsumables,
		talents: Presets.TalentsLaunchSmite.data,
		other: Presets.SmiteOtherDefaults,
	},
	presets: {
		...ShadowPage.presets,
		talents: [Presets.TalentsLaunchSmite, Presets.SmiteTalents],
		rotations: [Presets.ROTATION_PRESET_SMITE],
		gear: [Presets.GEAR_SMITE_LAUNCH],
	},
});
