// What the fork changed to turn the Classic sim into a Forever sim, grouped by what a
// reader would look for rather than by the order it happened. Every entry names the pull
// requests that carried it and, where the change came from published Forever information,
// the place it was read from. Numbers marked demo were read off BlizzCon 2026 tooltips.

export type Source = {
	label: string;
	url: string;
};

export type Entry = {
	title: string;
	prs: Array<number>;
	// What changed in the simulator.
	changed: string;
	// What it does to the numbers, compared with the Classic Era sim this fork started from.
	effect: string;
	sources?: Array<Source>;
};

export type Section = {
	title: string;
	intro: string;
	entries: Array<Entry>;
};

const blizzardPanel: Source = {
	label: "Blizzard: World of Warcraft: Forever What's Next panel recap",
	url: 'https://news.blizzard.com/en-us/article/24303862/world-of-warcraft-forever-whats-next-panel-recap',
};
const blizzardAnnounce: Source = {
	label: "Blizzard: World of Warcraft at BlizzCon 2026, Discover What's Next",
	url: 'https://news.blizzard.com/en-us/article/24301145/world-of-warcraft-at-blizzcon-2026-discover-whats-next',
};
const wowheadOverview: Source = {
	label: 'Wowhead: World of Warcraft: Forever overview (features, zones, raids)',
	url: 'https://www.wowhead.com/forever/guide/overview-features-zones-raids',
};
const wowheadRoadmap: Source = {
	label: 'Wowhead: Forever content release roadmap',
	url: 'https://www.wowhead.com/forever/guide/content-release-roadmap',
};
const wowheadRacials: Source = {
	label: 'Wowhead: all racials and race-class combinations in Forever',
	url: 'https://www.wowhead.com/forever/guide/new-race-class-combinations',
};
const wowheadTalents: Source = {
	label: 'Wowhead: Forever talent calculator (built from the BlizzCon demo, to be refreshed from the beta client)',
	url: 'https://www.wowhead.com/forever/talent-calc',
};
const communityTalents: Source = {
	label: 'Community talent calculators rebuilt from the BlizzCon demo tooltips (wowforevertalents.com, zockify.com)',
	url: 'https://wowforevertalents.com/',
};
const worldBuffsReport: Source = {
	label: 'N_Tys, 13 September: world buffs will not work in Forever raids (from Savix)',
	url: 'https://x.com/N_Tys26/status/2099230107147931731',
};
const worldBuffsClip: Source = {
	label: 'Savix on Twitch: WoW: Forever will not have world buffs in raids',
	url: 'https://www.twitch.tv/savix/clip/NastySillyWerewolfCclamChamp-V04O2vDDSTTAnULX',
};
const talentDataset: Source = {
	label: "The wow-forever-talent-calc dataset behind wowforevertalents.com, read off the demo video (MIT; the icons and crops are Blizzard's)",
	url: 'https://github.com/Deradon/wow-forever-talent-calc',
};
const upstream: Source = {
	label: 'wowsims/classic, the Classic Era simulator this fork started from',
	url: 'https://github.com/wowsims/classic',
};

export const sections: Array<Section> = [
	{
		title: 'The rules of the game',
		intro: 'Engine rules that apply whenever the sim runs under Forever Rules. Classic Era Rules are still there under Sim Options and leave every one of these off.',
		entries: [
			{
				title: 'A Forever ruleset, on by default',
				prs: [1, 43],
				changed: 'The sim carries a ruleset setting, Forever or Classic Era. New sessions start on Forever; the header says which one is running.',
				effect: 'Every change below is gated on it. Flip to Classic Era Rules and the fork simulates the Classic sim it came from.',
				sources: [upstream],
			},
			{
				title: 'Periodic damage can crit',
				prs: [1, 16],
				changed:
					'Damage-over-time ticks and bleeds roll for critical strikes using the crit chance snapshotted when the effect was applied. Ignite is excluded, since it is already a share of a crit.',
				effect: 'Every dot-heavy build gains; the warlock builds and the shadow priest most of all. Read from the wording of the new talents (Pandemic, and the "non-periodic" qualifiers on others).',
				sources: [communityTalents, wowheadTalents],
			},
			{
				title: 'Hit and crit from gear count for every kind of attack',
				prs: [18],
				changed: "An item's melee and spell hit are summed and paid into both pools, and the same for crit. Attribute conversions are untouched.",
				effect: 'Hybrids and casters stop wasting the melee hit and crit on their gear, and hunters and enhancement stop wasting spell crit.',
				sources: [blizzardPanel],
			},
			{
				title: 'Bonus healing carries a damage component',
				prs: [19],
				changed: 'A third of the healing power on gear is added to spell damage.',
				effect: 'Healing gear becomes usable by damage casters; the smite priest and the paladins gain the most.',
				sources: [blizzardPanel],
			},
			{
				title: 'The old raid debuffs are personal now',
				prs: [15, 65, 66, 69],
				changed:
					"Improved Shadow Bolt, Shadow Weaving, Improved Scorch, Winter's Chill and Stormstrike's Nature vulnerability only raise the damage of the caster who applied them. Improved Shadow Bolt lasts a flat twelve seconds. The raid panel no longer offers them under Forever.",
				effect: 'Stacking casters loses its Classic payoff: four warlocks no longer share one Improved Shadow Bolt, and a fire mage is not buffing the raid with Scorch. Each caster keeps their own bonus.',
				sources: [blizzardPanel],
			},
			{
				title: 'World buffs do not work inside raids',
				prs: [89, 91],
				changed:
					"Rallying Cry, Songflower, Darkmoon Faire, Warchief's Blessing, the Dire Maul tribute buffs and Spirit of Zandalar are ignored under Forever Rules, off in every default, and the World Buffs section only shows under Classic Era Rules.",
				effect: 'Every number on the site fell by a quarter to two fifths against the world-buffed Classic sims, and the gap between physical and caster specs narrowed from 30 to 21 points, because those buffs paid out in crit and attack power. Reported from the demo, not yet in patch notes.',
				sources: [worldBuffsReport, worldBuffsClip],
			},
			{
				title: 'Racials reworked, the Skyborne, six new race and class pairings',
				prs: [20, 21, 59],
				changed:
					"Every race has two actives and two passives. The resistance racials are gone; weapon skill racials became crit while the weapon is held; Blood Fury, Expansive Mind, Elune's Light, Big Game Hunter and the Skyborne racials are modelled. Dwarf shamans, Undead paladins and the rest of the new pairings can be simulated.",
				effect: 'Race choice moves numbers by a few percent as before, but for different reasons; a Skyborne is available to Warrior, Hunter, Rogue, Druid, Horde Shaman and Alliance Mage.',
				sources: [wowheadRacials, blizzardAnnounce],
			},
			{
				title: 'Profession passives with published numbers',
				prs: [75],
				changed: 'Skinning is +5% damage against Beasts and Dragonkin; Mining is +5% health.',
				effect: 'Only shows against a target of those types; Molten Core bosses are not, so the presets are unaffected.',
				sources: [blizzardPanel],
			},
			{
				title: 'Baseline ability changes',
				prs: [77],
				changed:
					'Slam no longer resets the swing timer, Thunder Clap works in Defensive Stance, Improved Shield Wall shortens the cooldown, Tactical Mastery is baseline. Victory Rush is not modelled (it needs a killing blow).',
				effect: 'Warrior only; nothing concrete has been published for the other classes yet.',
				sources: [blizzardPanel],
			},
		],
	},
	{
		title: 'Talents',
		intro: 'The nine Forever talent trees, read from BlizzCon demo tooltips, and what the sim does with them.',
		entries: [
			{
				title: 'Nine classes on the Forever trees',
				prs: [2, 1, 3, 5, 6, 8, 10, 11, 12, 13],
				changed:
					"The trees were imported from tooltip data extracted at BlizzCon, then each class converted: new talents implemented where the tooltip gave enough to go on (Mangle, Berserk, Eclipse, Lava Burst, Lightning Overload, Maelstrom Weapon, Mutilate, Arcane Blast, Hot Streak, Ice Lance, Sniper Shot, Lone Wolf, Weaponmaster, Bloodthrill, Penance, Holy Nova, the paladin's Holy Strike and its dependents, and more), and Classic talents that Forever removed taken out.",
				effect: 'Every build on the site is a Forever build. Talent strings are positional, so a Classic talent string will not load.',
				sources: [communityTalents, wowheadTalents],
			},
			{
				title: 'The trees look and read like the published ones',
				prs: [44, 58, 60, 62, 48],
				changed:
					'New talents have icons and tooltips; positions match the published grids; Improved Fireball and the five-rank Shatter came back; talents the sim does not read are marked as such in the picker.',
				effect: 'What you click is what the sim runs, and the picker tells you when it is not.',
				sources: [wowheadTalents, communityTalents],
			},
			{
				title: 'Talent numbers audited against the tooltips',
				prs: [73, 47, 9, 17],
				changed:
					"Eight discrepancies fixed (Precision as spell hit, Holy Shield, Barrage on Aimed Shot, Mind Flay and Pyroblast base damage, Call of Flame on Lava Burst and Flame Shock, Savage Fury on Shred); the paladin's Forever talents implemented; the talents whose per-rank scaling was never shown are listed for the beta.",
				effect: 'Shadow priest +20% from Mind Flay alone, Elemental +5%, the rest one or two percent. Higher ranks of the demo-only talents are extrapolated from rank one until the beta shows them.',
				sources: [communityTalents],
			},
			{
				title: 'The community builds, everywhere',
				prs: [63, 71, 67, 78],
				changed:
					"Thirty-one community builds are talent presets in their specs, listed under each class on the homepage, and each spec's automatic rotation follows the build's tree (a Fire build casts the Fire rotation, a Mutilate rogue uses daggers).",
				effect: 'The rankings and the spec pages simulate the builds people are actually discussing.',
			},
			{
				title: 'Tests that pin the trees together',
				prs: [68],
				changed:
					'Go tests check that every tree, its proto message, the class tree sizes and every shipped build agree, and that every build link on the homepage names a real preset.',
				effect: 'Three separate incidents in this program would have been caught before deploy.',
			},
		],
	},
	{
		title: 'Content and items',
		intro: "Forever launches on 4 November 2026 with Onyxia's Lair and Molten Core as its first raids (9 December), and nothing past them.",
		entries: [
			{
				title: 'The item pool is launch content only',
				prs: [46, 25, 27, 29],
				changed:
					"Blackwing Lair, Zul'Gurub, Ahn'Qiraj and Naxxramas items, sets and effects are out of the database; the sim's phases are Forever's tiers; Onyxia is the one tier 1 encounter with a known fight.",
				effect: 'Gear pickers and the best in slot page only show what exists at launch. The two new raids (Barrow Deeps, Hyjal Summit) have no encounter until their bosses are announced.',
				sources: [wowheadOverview, wowheadRoadmap, blizzardAnnounce],
			},
			{
				title: 'Launch gear for every spec',
				prs: [26, 28],
				changed:
					'Pre-raid gear sets generated from the launch pool for the specs that had none (Retribution, Protection Paladin, Smite Priest and others), named as tiers in the UI.',
				effect: 'Every spec in the rankings wears comparable, launch-tier gear.',
			},
		],
	},
	{
		title: 'Specs and rotations',
		intro: 'What runs, and what was made to run properly.',
		entries: [
			{
				title: 'Smite Priest',
				prs: [31],
				changed: 'A holy damage priest, which Classic never had a reason to run.',
				effect: 'One more spec in the rankings.',
			},
			{
				title: 'Rotations that use the kit',
				prs: [35, 36, 37, 24, 49],
				changed:
					'The protection warrior stopped running the fury priority, the balance druid casts more than Wrath rank four, the arcane mage stopped going out of mana on Arcane Blast stacks, the feral got its bleeds and a rotation.',
				effect: 'Several specs moved by a large fraction of their DPS; these were bugs in the shipped presets, not Forever changes.',
			},
			{
				title: 'Warlock rotations by build',
				prs: [100],
				changed:
					'Demonic Pact runs with the Imp sacrificed, a Succubus out and Soul Link up; DS/Ruin summons its own Imp before sacrificing; every warlock keeps Immolate up, curses when nobody else has, and drinks when mana allows. Shadow and Flame has its own rotation. One raid preset per tree carries the pet setup each rotation expects, and the community builds are preset builds.',
				effect: 'Demonic Pact +41%, Deep Affliction +13%, DS/Ruin Pandemic +3%, Shadow and Flame +4% at two minutes; more at five, where Life Tap was a tenth of the fight.',
			},
			{
				title: "Feral cat: Tiger's Fury as the energy engine",
				prs: [105],
				changed:
					"Tiger's Fury takes Wrath's shape (no Energy cost, 30 sec cooldown) so King of the Jungle's 60 Energy is a cooldown rather than something to spam, the cat rotation casts it on cooldown, and Rake is dropped: it cost two fifths of the cat's Energy for a twentieth of its damage.",
				effect: "The ranked Feral Cat build goes from 393 to 478 DPS on a two minute fight, still the lowest of the melee, because Forever's Furor no longer hands out 40 Energy per shift.",
				sources: [communityTalents],
			},
			{
				title: 'Warrior: the off-hand misses again, and arms strikes',
				prs: [108],
				changed:
					"Queuing Heroic Strike lifted the dual wield miss penalty for the whole character instead of for the queued swing, so a Fury build's off-hand auto attacks almost never missed. The DPS priority list had no Mortal Strike line, so the ranked Arms build ran a Fury rotation, and Protection never cast Thunder Clap even though Forever allows it in Defensive Stance.",
				effect: 'Fury -1.5%, Arms +5.2%, Protection +1.8% DPS and +4.9% threat on a two minute fight.',
				sources: [communityTalents],
			},
			{
				title: 'Paladin: Holy Strike on its tooltip, Judgement paying its talents',
				prs: [109],
				changed:
					"Holy Strike was guessed at 110% weapon damage on a six second cooldown, and it was 40% of the retribution paladin's damage; the published tooltip is 40% weapon damage plus 36 to 46 Holy damage on a twelve second cooldown for 20 mana. Judgement no longer suppresses cast triggers, so Sanctified Judgement refunds mana and Swift Judgement's free cast stops being permanent. Retribution twists Seal of Command into Righteousness rather than the other way round, drinks its own potion and rune, casts Hammer of Wrath in the execute window, and fills spare globals with Consecration.",
				effect: 'Retribution goes from 780 to 679 DPS on a two minute fight and protection from 331 to 284, the rotation work giving back about half of what the ability correction took.',
			},
			{
				title: 'Hunter: the cat bites, Serpent Sting waits for mana',
				prs: [110],
				changed:
					"The pet's Bite was gated behind a focus income no pet can reach, so the cat only clawed; Serpent Sting was refreshed with a third of the dot still running and is the hunter's worst shot per point of mana, so it now waits until the dot is nearly gone and is skipped under thirty percent mana. Volley's crit bonus folded Mortal Shots in as a multiplier, and Lethal Attacks missed spell crit.",
				effect: 'Marksmanship +14 DPS at two minutes and +22 at five, Beast Mastery +2 and +7, Survival +2 and +10; the Marksmanship hunter spends 25 fewer seconds out of mana on a five minute fight.',
			},
			{
				title: 'Priest: Power Infusion returns, and Power in Light stays up',
				prs: [113],
				changed:
					"Power Infusion, the thirty-one point talent the Smite build spends its deepest point on, was commented out of the sim entirely. It is cast again. Holy Fire is now recast as its dot runs out rather than after it has dropped, so Power in Light no longer falls off for five seconds in every fifteen, Inner Focus goes on Smite instead of Penance, and Shadow Word: Pain and a downranked Smite fill the gaps a five minute fight opens. Devouring Plague registered only its first five ranks and was locked to the Undead, so the shadow rotation's Devouring Plague line did nothing at all for the ranked Dwarf.",
				effect: 'Smite +7.1% at two minutes and +5.6% at five; Shadow +0.3% and +0.4%.',
				sources: [communityTalents],
			},
			{
				title: 'Rogue: a rotation for the Hemorrhage build',
				prs: [111],
				changed:
					'The auto rotation only knew Mutilate, Backstab and Sinister Strike, so a Subtlety build that had spent fifteen points on Hemorrhage never cast it, and Rupture was in no rogue rotation at all, leaving Serrated Blades and Thousand Cuts reading off a debuff nobody applied. Hemorrhage gets its own list with Rupture and a Vanish-Premeditation-Ambush opener. Venom now reaches a Deadly Poison that is already ticking, and the Assassination list finishes at four combo points instead of letting nine a fight fall off the cap.',
				effect: 'Subtlety Hemo 582 to 709 DPS and Assassination Mutilate 638 to 696 on a two minute fight; Combat was measured against the same changes and kept its rotation.',
				sources: [communityTalents],
			},
			{
				title: 'Enhancement shaman: a rotation, and an off-hand that works',
				prs: [112],
				changed:
					"The enhancement priority list asked whether Strength of Earth Totem was up using the wrong rank's aura, so the shaman re-dropped it every global until it ran out of mana and never reached Stormstrike or Earth Shock; the totem lines now read how long each totem has left. Windfury Totem, which shared the air slot with Grace of Air and buffed nothing, is gone. The off-hand weapon imbue was never read at all, Windfury's extra attacks always came from the main hand, and Elemental Weapons was applied twice to Windfury's bonus attack power.",
				effect: 'The ranked Enhancement 16/35/0 build goes from 473 to 724 DPS at two minutes and from 458 to 699 at five, and is no longer the lowest non-tank.',
				sources: [communityTalents],
			},
			{
				title: 'Mage: Ignite paid the same crit twice',
				prs: [114],
				changed:
					"Ignite's ticks re-applied the Improved Scorch stacks and Curse of Elements that the critical strike had already carried, and a crit landing after the dot had ticked restarted it with the old damage still in the pool, so the talent paid out 62% of the crit it lit instead of 40%. Fingers of Frost handed the Shatter crit to the Frostbolt already half cast without spending the charge on it. Arcane Missiles stopped at rank 7 because its registration loop predated the AQ ranks, and the arcane rotation now holds three Arcane Blast stacks before spending them. Mana gems no longer spend the shared conjured cooldown on the smallest gem.",
				effect: 'Fire -7.6%, Frost -4.0%, Arcane +8.9% on a two minute fight; the three mage builds land within 90 DPS of each other instead of 300.',
				sources: [communityTalents],
			},
			{
				title: 'Season of Discovery leftovers removed',
				prs: [82, 83, 67],
				changed:
					'The Warden (tank) shaman, a Season of Discovery spec that came with the upstream code, is gone; the healing specs and the bear, which have no working sim, no longer carry Classic builds.',
				effect: 'The site only offers what Forever has and the sim can run.',
				sources: [upstream],
			},
		],
	},
	{
		title: 'The site',
		intro: 'Pages and tooling around the sim.',
		entries: [
			{
				title: 'DPS rankings from one raid',
				prs: [54, 74, 78, 81, 87, 88],
				changed:
					'One 25-player raid with a slot per community build, one run, one encounter, one set of buffs; both warlock curses up; rows named by spec and build.',
				effect: "The only page where every build's number comes from the same fight.",
			},
			{
				title: 'Rankings in launch gear',
				prs: [106],
				changed:
					'Every ranked build wears a Launch set: the best pre-raid gear in the launch item pool by its own stat weights, built by one tool for all sixteen specs, raid drops and faction-locked items left out. Before, rogues ranked in Pre-BiS, casters in thin Classic sets and the cat in eight Wildheart pieces and nine empty slots.',
				effect: 'The rankings compare specs rather than gear tiers; casters and the cat moved most.',
			},
			{
				title: 'Best in slot and stat weights',
				prs: [55, 56, 72, 79],
				changed: "A best in slot page from each spec's EP weights over the launch pool, and a page comparing stat weights across specs.",
				effect: 'Reference pages; neither runs the sim.',
			},
			{
				title: 'Landing page, feedback and the raid sim',
				prs: [38, 41, 45, 39, 40, 32, 33, 51, 52, 53],
				changed:
					"Forever's spec list and lockup on the homepage, a dropdown per class, in-sim feedback that opens a GitHub issue with a screenshot, the raid sim reachable and working on touch screens.",
				effect: 'Site plumbing.',
			},
			{
				title: 'Publishing and checks',
				prs: [34, 57, 76, 80, 85, 86],
				changed:
					"Published to GitHub Pages from master; every page opened in a headless browser before deploy and on every pull request; the sim pages generated by vite like upstream's newer sims.",
				effect: 'A build whose pages throw cannot go live; it has already stopped one.',
				sources: [upstream],
			},
			{
				title: 'Images served from the site',
				prs: [95, 101, 103],
				changed:
					"Every icon, tree background and pet icon mirrored into the repo and loaded from there instead of Wowhead's CDN; the talents Forever added drawn from the icon names the dataset matched, or from crops of the demo video where the icon is new.",
				effect: 'Icons load on networks that reject wow.zamimg.com, and no talent shows a question mark.',
				sources: [talentDataset],
			},
			{
				title: 'Written down for the beta and for upstream',
				prs: [70, 84, 90],
				changed:
					'A checklist of every demo-tooltip assumption the beta must confirm, a rules sheet of every Forever rule the sim models with its source, and a note on where upstream wowsims is heading.',
				effect: 'The beta pass is a checklist rather than an audit, and the official fork can take the data.',
				sources: [upstream],
			},
		],
	},
];
