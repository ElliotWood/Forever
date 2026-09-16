import { Spec } from './proto/common.js';

// The community talent builds that the landing page lists under each spec, in the order
// it lists them. A name here is the name of one of that spec's talent presets: the
// ?build= link looks the preset up by name when the sim loads, so the two have to agree
// exactly.
//
// The list is written out here rather than read from each ui/<spec>/presets.ts because
// ui/core cannot import those. Every page bundles ui/core, so importing the presets of
// every spec would pull every sim's gear sets and rotations into every page. It could
// not be derived from them anyway: a spec's talent presets also hold its defaults and
// its per-phase builds, and which of them are community builds is a choice, not
// something written down in presets.ts.
//
// TestSimTitleDropdownBuildsMatchLandingPage in sim/talents_test.go fails if this and
// ui/index.html ever come to disagree, and checks that every name here is a preset the
// spec actually has.
export const communityBuilds: Record<Spec, string[]> = {
	[Spec.SpecBalanceDruid]: ['Moonkin 38/0/13'],
	[Spec.SpecFeralDruid]: ['Feral Cat 9/35/7'],
	[Spec.SpecFeralTankDruid]: ['Bear Tank 0/31/20'],
	[Spec.SpecRestorationDruid]: [],
	[Spec.SpecElementalShaman]: ['Elemental 31/6/14', 'Stormcaller 28/23/0'],
	[Spec.SpecEnhancementShaman]: ['Enhancement 16/35/0'],
	[Spec.SpecRestorationShaman]: [],
	[Spec.SpecHunter]: ['Beast Mastery 35/16/0', 'Marksmanship 0/39/12', 'Survival 0/15/36'],
	[Spec.SpecMage]: ['Fire 0/35/16', 'Frost 14/0/37', 'Arcane 35/0/16'],
	[Spec.SpecRogue]: ['Combat Dual-Wield 15/33/3', 'Assassination Mutilate 31/20/0', 'Subtlety Hemo 15/0/36'],
	[Spec.SpecHolyPaladin]: [],
	[Spec.SpecProtectionPaladin]: ['Protection 0/45/6'],
	[Spec.SpecRetributionPaladin]: ['Retribution 10/0/41'],
	[Spec.SpecHealingPriest]: [],
	[Spec.SpecShadowPriest]: ['Shadow 15/0/36'],
	[Spec.SpecSmitePriest]: ['Smite 31/17/3'],
	[Spec.SpecWarlock]: ['Demonic Pact 2/31/18', 'Deep Affliction 35/0/16', 'DS/Ruin Pandemic 24/11/16', 'Shadow and Flame 13/11/27'],
	[Spec.SpecWarrior]: ['Fury 17/34/0', 'Arms 39/12/0'],
	[Spec.SpecTankWarrior]: ['Protection 1/0/50'],
};

export function getCommunityBuilds(spec: Spec): string[] {
	return communityBuilds[spec] || [];
}
