import { Drums, Race, Stat } from '@generated/proto/common';
import { Player } from '@sim/player/player';
import { ActionId } from '@sim/proto/action_id';
import { Party } from '@sim/raid/party';
import {
	makeBooleanDebuffInput,
	makeBooleanIndividualBuffInput,
	makeBooleanPartyBuffInput,
	makeBooleanRaidBuffInput,
	makeMultistateIndividualBuffInput,
	makeMultistatePartyBuffInput,
	makeTristatePartyBuffInput,
} from '@ui-kit/icon_inputs';
import * as InputHelpers from '@ui-kit/input_helpers';

import * as Generated from './buffs_debuffs_auto_gen';
import { DrumsBattle, DrumsRestoration, DrumsWar } from './consumables';
import { IconPickerStatOption, RenderableStatOptions } from './stat_options';

// Every buff the client database resolves to a spell has its input generated from the manifest;
// the rows below are the ones it cannot produce, and the registries at the end of this file
// interleave the two.
export * from './buffs_debuffs_auto_gen';

// The generated const takes its name from the proto field; the specs name the buff.
export const Innervate = Generated.Innervates;
export const PowerInfusion = Generated.PowerInfusions;
export const ManaTideTotem = Generated.ManaTideTotems;

///////////////////////////////////////////////////////////////////////////
//                                 RAID BUFFS
///////////////////////////////////////////////////////////////////////////

export const Bloodlust = makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(2825), fieldName: 'bloodlust', label: 'Bloodlust' });

///////////////////////////////////////////////////////////////////////////
//                                 PARTY BUFFS
///////////////////////////////////////////////////////////////////////////

export const BraidedEterniumChain = makeBooleanPartyBuffInput({
	actionId: ActionId.fromSpellId(31025),
	fieldName: 'braidedEterniumChain',
	label: 'Braided Eternium Chain',
});
export const ChainOfTheTwilightOwl = makeBooleanPartyBuffInput({
	actionId: ActionId.fromSpellId(31035),
	fieldName: 'chainOfTheTwilightOwl',
	label: 'Chain of the Twilight Owl',
});
export const DraeneiRacialCaster = makeBooleanPartyBuffInput({
	actionId: ActionId.fromSpellId(28878),
	fieldName: 'draeneiRacialCaster',
	label: 'Inspiring Presense - Caster',
	showWhen: (party: Party) => [Race.RaceDraenei, Race.RaceDwarf, Race.RaceGnome, Race.RaceHuman, Race.RaceNightElf].includes(party.getPlayer(0)!.getRace()),
});
export const DraeneiRacialMelee = makeBooleanPartyBuffInput({
	actionId: ActionId.fromSpellId(6562),
	fieldName: 'draeneiRacialMelee',
	label: 'Inspiring Presense - Melee',
	showWhen: (party: Party) => [Race.RaceDraenei, Race.RaceDwarf, Race.RaceGnome, Race.RaceHuman, Race.RaceNightElf].includes(party.getPlayer(0)!.getRace()),
});
export const EyeOfTheNight = makeBooleanPartyBuffInput({ actionId: ActionId.fromSpellId(31033), fieldName: 'eyeOfTheNight', label: 'Eye of the Night' });
export const FerociousInspiration = makeMultistatePartyBuffInput({
	actionId: ActionId.fromSpellId(34460),
	numStates: 5,
	fieldName: 'ferociousInspiration',
	label: 'Ferocious Inspiratation',
});
export const JadePendantOfBlasting = makeBooleanPartyBuffInput({
	actionId: ActionId.fromSpellId(25607),
	fieldName: 'jadePendantOfBlasting',
	label: 'Jade Pendant of Blasting',
});
export const SanctityAura = makeTristatePartyBuffInput({
	actionId: ActionId.fromSpellId(20218),
	impId: ActionId.fromSpellId(31870),
	fieldName: 'sanctityAura',
	label: 'Sanctity Aura',
});
export const TotemOfWrath = makeMultistatePartyBuffInput({
	actionId: ActionId.fromSpellId(30706),
	numStates: 5,
	fieldName: 'totemOfWrath',
	label: 'Totem of Wrath',
});
export const WrathOfAirTotem = makeTristatePartyBuffInput({
	actionId: ActionId.fromSpellId(3738),
	impId: ActionId.fromSpellId(37212),
	fieldName: 'wrathOfAirTotem',
	label: 'Wrath of Air Totem',
});

// The drums party buff is a mutually-exclusive swatch pick (none / battle / war /
// restoration) over the shared `PartyBuffs.drums` field, not a boolean toggle, so it
// is composed directly rather than through `makeEnumValuePartyBuffInput` (which only
// wraps a single on/off value).
export const DrumsBuff: InputHelpers.TypedIconEnumPickerConfig<Player<any>, Drums> = {
	type: 'iconEnum',
	label: 'Drums',
	values: [{ color: 'gray', value: Drums.DrumsUnknown }, DrumsBattle, DrumsWar, DrumsRestoration],
	zeroValue: Drums.DrumsUnknown,
	equals: (a: Drums, b: Drums) => a === b,
	storeField: 'raid:partyBuffs',
	getValue: (player: Player<any>) => player.getParty()!.getBuffs().drums,
	setValue: (player: Player<any>, newValue: Drums) => {
		const party = player.getParty()!;
		const newBuffs = party.getBuffs();
		newBuffs.drums = newValue;
		party.setBuffs(newBuffs);
	},
};

// Individual Buffs
export const BlessingOfSalvation = makeBooleanIndividualBuffInput({
	actionId: ActionId.fromSpellId(25895),
	fieldName: 'blessingOfSalvation',
	label: 'Blessing of Salvation',
	showWhen: player => !player.getPlayerSpec().isTankSpec && !player.getPlayerSpec().isHealingSpec,
});
export const BlessingOfSanctuary = makeBooleanIndividualBuffInput({
	actionId: ActionId.fromSpellId(27169),
	fieldName: 'blessingOfSanctuary',
	label: 'Blessing of Sanctuary',
});
export const UnleashedRage = makeBooleanIndividualBuffInput({
	actionId: ActionId.fromSpellId(30811),
	fieldName: 'unleashedRage',
	label: 'Unleashed Rage',
});
export const ShadowPriestDPS = makeMultistateIndividualBuffInput({
	actionId: ActionId.fromSpellId(34917),
	numStates: 1500,
	fieldName: 'shadowPriestDps',
	label: 'Vampiric Touch',
});

// Debuffs
export const BloodFrenzy = makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(29859), fieldName: 'bloodFrenzy', label: 'Blood Frenzy' });
export const ImprovedScorch = makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(12873), fieldName: 'improvedScorch', label: 'Improved Scorch' });
export const Misery = makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(33195), fieldName: 'misery', label: 'Misery' });
export const ShadowWeaving = makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(15334), fieldName: 'shadowWeaving', label: 'Shadow Weaving' });
export const WintersChill = makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(28595), fieldName: 'wintersChill', label: "Winter's Chill" });
export const Screech = makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(27051), fieldName: 'screech', label: 'Screech' });
export const ShadowEmbrace = makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(32394), fieldName: 'shadowEmbrace', label: 'Shadow Embrace' });

// A registry lists its rows in display order: a generated input config stands for the generated
// row that carries it, with the stat tags and owner class the client database produced, and a
// literal row is one of the hand-written buffs above.
const inDisplayOrder = (
	generated: Generated.GeneratedStatOption[],
	rows: Array<RenderableStatOptions | RenderableStatOptions['config']>,
): RenderableStatOptions[] =>
	rows.map(row => {
		if ('config' in row) return row;
		const option = generated.find(candidate => candidate.config === row);
		if (!option) throw new Error(`no generated row carries the buff input "${row.label}"`);
		return option;
	});

export const PARTY_BUFFS_CONFIG = inDisplayOrder(Generated.GENERATED_PARTY_BUFFS_CONFIG, [
	Generated.BloodPact,
	Generated.CommandingShout,
	Generated.BattleShout,
	Generated.DevotionAura,
	{ config: FerociousInspiration, stats: [] },
	Generated.LeaderOfThePack,
	Generated.ManaSpringTotem,
	Generated.ManaTideTotems,
	{ config: ShadowPriestDPS, stats: [Stat.StatMP5] },
	Generated.MoonkinAura,
	Generated.RetributionAura,
	Generated.ConcentrationAura,
	{ config: SanctityAura, stats: [] },
	{ config: TotemOfWrath, stats: [Stat.StatSpellCritRating, Stat.StatSpellHitRating] },
	Generated.TrueshotAura,
	{ config: WrathOfAirTotem, stats: [Stat.StatSpellDamage] },
	{ config: UnleashedRage, stats: [Stat.StatAttackPower] },
	Generated.AtieshMage,
	Generated.AtieshWarlock,
	{ config: BraidedEterniumChain, stats: [Stat.StatMeleeCritRating] },
	{ config: ChainOfTheTwilightOwl, stats: [Stat.StatSpellCritRating] },
	{ config: DraeneiRacialCaster, stats: [Stat.StatSpellHitRating] },
	{ config: DraeneiRacialMelee, stats: [Stat.StatMeleeHitRating] },
	{ config: EyeOfTheNight, stats: [Stat.StatSpellDamage] },
	{ config: JadePendantOfBlasting, stats: [Stat.StatSpellDamage] },
	Generated.StrengthOfEarthTotem,
	Generated.GraceOfAirTotem,
	Generated.WindfuryTotem,
	{ config: DrumsBuff, stats: [] },
	Generated.FrostResistanceTotem,
	Generated.NatureResistanceTotem,
	Generated.FireResistanceTotem,
	Generated.FrostResistanceAura,
	Generated.FireResistanceAura,
	Generated.ShadowResistanceAura,
	Generated.AspectOfTheWild,
]);

export const BUFFS_CONFIG = inDisplayOrder(
	[...Generated.GENERATED_RAID_BUFFS_CONFIG, ...Generated.GENERATED_INDIVIDUAL_BUFFS_CONFIG],
	[
		Generated.ArcaneBrilliance,
		Generated.BlessingOfKings,
		{ config: Bloodlust, stats: [] },
		Generated.DivineSpirit,
		Generated.GiftOfTheWild,
		Generated.Thorns,
		Generated.PowerWordFortitude,
		Generated.BlessingOfMight,
		Generated.BlessingOfWisdom,
		{ config: BlessingOfSanctuary, stats: [Stat.StatStamina, Stat.StatArmor] },
		{ config: BlessingOfSalvation, stats: [] },
		Generated.ShadowProtection,
		Generated.Innervates,
		Generated.PowerInfusions,
	],
);

export const DEBUFFS_CONFIG = inDisplayOrder(Generated.GENERATED_DEBUFFS_CONFIG, [
	{ config: BloodFrenzy, stats: [Stat.StatAttackPower] },
	Generated.HuntersMark,
	{ config: ImprovedScorch, stats: [Stat.StatFireDamage] },
	Generated.ImprovedSealOfTheCrusader,
	Generated.JudgementOfLight,
	Generated.JudgementOfWisdom,
	Generated.Mangle,
	{ config: Misery, stats: [] },
	{ config: ShadowWeaving, stats: [Stat.StatShadowDamage] },
	Generated.CurseOfElements,
	Generated.CurseOfRecklessness,
	Generated.FaerieFire,
	Generated.ExposeArmor,
	Generated.SunderArmor,
	{ config: WintersChill, stats: [Stat.StatFrostDamage] },
	Generated.GiftOfArthas,
	Generated.DemoralizingRoar,
	Generated.DemoralizingShout,
	{ config: Screech, stats: [Stat.StatStamina, Stat.StatResilienceRating] },
	Generated.ThunderClap,
	Generated.InsectSwarm,
	Generated.ScorpidSting,
	{ config: ShadowEmbrace, stats: [Stat.StatStamina, Stat.StatResilienceRating] },
]);

export const DEBUFFS_MISC_CONFIG = [] as IconPickerStatOption[];
