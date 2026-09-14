import { Stat } from '../proto/common';

// Forever's content tiers. Forever launches with no raids at all - the original 1-60
// content plus three new zones and nine new dungeons - and opens its first two raids five
// weeks later. This is what a sim's tier means.
export enum Phase {
	// 4 November 2026. No raids.
	Launch = 1,
	// 9 December 2026. Barrow Deeps (10 player), Hyjal Summit (20 player).
	Tier1,
	// Spring 2027. Two more raids, a 10 and a 20.
	Tier2,
	// Summer 2027. A returning raid alongside another new one.
	Tier3,
}

export const CURRENT_PHASE = Phase.Launch;

// Classic Era's six content tiers, Molten Core through Naxxramas. Separate from the
// Forever tiers above because it means something different: it is when an item became
// available in Classic, which is the only availability the item database carries. Gear
// set presets are tagged with it and the gear picker filters on it, and it stays until
// there is a Forever item database to replace it with.
export enum ClassicPhase {
	Phase1 = 1,
	Phase2,
	Phase3,
	Phase4,
	Phase5,
	Phase6,
}

export const CURRENT_CLASSIC_PHASE = ClassicPhase.Phase6;

// Github pages serves our site under the /classic directory
export const REPO_NAME = 'classic';

// Get 'elemental_shaman', the pathname part after the repo name
const pathnameParts = window.location.pathname.split('/');
const repoPartIdx = pathnameParts.findIndex(part => part == REPO_NAME);
export const SPEC_DIRECTORY = repoPartIdx == -1 ? '' : pathnameParts[repoPartIdx + 1];

export const GLOBAL_DISPLAY_STATS = [Stat.StatHealth, Stat.StatStamina, Stat.StatFireResistance, Stat.StatFrostResistance, Stat.StatNatureResistance];

export const GLOBAL_DISPLAY_PSEUDO_STATS = [];

export const GLOBAL_EP_STATS = [Stat.StatFireResistance, Stat.StatFrostResistance, Stat.StatNatureResistance];

export enum SortDirection {
	ASC,
	DESC,
}
