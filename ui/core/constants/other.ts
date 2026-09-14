import { Stat } from '../proto/common';

// Forever's content tiers. The game is set after Warcraft III Reforged: Forsaken Kingdom
// and *before* Molten Core, so it launches with no raids at all - the original 1-60 content
// plus three new zones and nine new dungeons - and opens its first raids five weeks later.
// That placement is why the launch gear sets contain no raid loot: at this point in the
// timeline there is none to have.
export enum Phase {
	// 4 November 2026. No raids.
	Launch = 1,
	// 9 December 2026. Barrow Deeps (10 player), Hyjal Summit (20 player) and Onyxia's
	// Lair returning at 40 - Forever raids at all three sizes, not only the small ones.
	Tier1,
	// Spring 2027. Two more raids, a 10 and a 20.
	Tier2,
	// Summer 2027. An iconic raid returns alongside another new one. Which one has not
	// been named; Forever's timeline sits before Molten Core, so that is the direction of
	// travel rather than a confirmed raid.
	Tier3,
}

export const CURRENT_PHASE = Phase.Launch;

// The tiers are named, not numbered, so they need labels of their own rather than
// 'Phase ' plus the enum value.
export const phaseNames: Record<Phase, string> = {
	[Phase.Launch]: 'Launch',
	[Phase.Tier1]: 'Tier 1',
	[Phase.Tier2]: 'Tier 2',
	[Phase.Tier3]: 'Tier 3',
};

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

// Where the site is rooted, with a trailing slash. Vite fills this from its `base`, so
// it follows wherever the site is actually published: '/classic/' by default, and
// something like '/Forever/classic/' when it is a project page hanging off a repo name.
// Everything that used to hard-code '/classic/' goes through this.
export const SITE_BASE = import.meta.env.BASE_URL;

// Which build of the sim this is. Vite fills it from `git describe` so nobody has to
// remember to bump a number, which means it is a tag once the fork starts cutting them
// and a commit until then. The sidebar shows it so a bug report can name the build.
export const SITE_VERSION = __SITE_VERSION__;

// The Github repository this build points at, as '<owner>/<repo>'. Defaults to upstream so
// a wowsims build is unchanged, and a fork sets SITE_REPO at build time. Everything that
// links to source, issues, crash reports or releases goes through these two.
export const SITE_REPO = __SITE_REPO__;
export const SITE_REPO_URL = `https://github.com/${SITE_REPO}`;

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
