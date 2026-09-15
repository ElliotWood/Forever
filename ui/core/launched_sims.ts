import { Phase } from './constants/other';
import { Class, Spec } from './proto/common';
import { specToClass } from './proto_utils/utils';

// This file is for anything related to launching a new sim. DO NOT touch this
// file until your sim is ready to launch!

export enum LaunchStatus {
	Unlaunched,
	Alpha,
	Beta,
	Launched,
}

export type SimStatus = {
	phase: Phase;
	status: LaunchStatus;
};

export const raidSimStatus: SimStatus = {
	phase: Phase.Launch,
	status: LaunchStatus.Alpha,
};

// This list controls which links are shown in the top-left dropdown menu.
export const simLaunchStatuses: Record<Spec, SimStatus> = {
	[Spec.SpecBalanceDruid]: {
		phase: Phase.Launch,
		status: LaunchStatus.Beta,
	},
	[Spec.SpecFeralDruid]: {
		phase: Phase.Launch,
		status: LaunchStatus.Alpha,
	},
	[Spec.SpecFeralTankDruid]: {
		phase: Phase.Launch,
		status: LaunchStatus.Alpha,
	},
	[Spec.SpecRestorationDruid]: {
		phase: Phase.Launch,
		status: LaunchStatus.Unlaunched,
	},
	[Spec.SpecElementalShaman]: {
		phase: Phase.Launch,
		status: LaunchStatus.Launched,
	},
	[Spec.SpecEnhancementShaman]: {
		phase: Phase.Launch,
		status: LaunchStatus.Launched,
	},
	[Spec.SpecRestorationShaman]: {
		phase: Phase.Launch,
		status: LaunchStatus.Unlaunched,
	},
	[Spec.SpecHunter]: {
		phase: Phase.Launch,
		status: LaunchStatus.Alpha,
	},
	[Spec.SpecMage]: {
		phase: Phase.Launch,
		status: LaunchStatus.Alpha,
	},
	[Spec.SpecRogue]: {
		phase: Phase.Launch,
		status: LaunchStatus.Alpha,
	},
	[Spec.SpecHolyPaladin]: {
		phase: Phase.Launch,
		status: LaunchStatus.Unlaunched,
	},
	[Spec.SpecProtectionPaladin]: {
		phase: Phase.Launch,
		status: LaunchStatus.Alpha,
	},
	[Spec.SpecRetributionPaladin]: {
		phase: Phase.Launch,
		status: LaunchStatus.Alpha,
	},
	[Spec.SpecHealingPriest]: {
		phase: Phase.Launch,
		status: LaunchStatus.Unlaunched,
	},
	[Spec.SpecShadowPriest]: {
		phase: Phase.Launch,
		status: LaunchStatus.Alpha,
	},
	[Spec.SpecSmitePriest]: {
		phase: Phase.Launch,
		status: LaunchStatus.Alpha,
	},
	[Spec.SpecWarlock]: {
		phase: Phase.Launch,
		status: LaunchStatus.Alpha,
	},
	[Spec.SpecWarrior]: {
		phase: Phase.Launch,
		status: LaunchStatus.Alpha,
	},
	[Spec.SpecTankWarrior]: {
		phase: Phase.Launch,
		status: LaunchStatus.Alpha,
	},
};

export function getLaunchedSims(): Array<Spec> {
	return Object.keys(simLaunchStatuses)
		.map(specStr => parseInt(specStr) as Spec)
		.filter(spec => simLaunchStatuses[spec].status > LaunchStatus.Unlaunched);
}

export function getLaunchedSimsForClass(klass: Class): Array<Spec> {
	return Object.keys(specToClass)
		.map(specStr => parseInt(specStr) as Spec)
		.filter(spec => specToClass[spec] == klass && isSimLaunched(spec));
}

export function isSimLaunched(specIndex: Spec): boolean {
	return simLaunchStatuses[specIndex].status > LaunchStatus.Unlaunched;
}
