import { Class, Stat } from '@generated/proto/common';
import { ActionId } from '@sim/proto/action_id';
import { makeBooleanIndividualBuffInput, makeBooleanPartyBuffInput } from '@ui-kit/icon_inputs';

import * as Generated from './buffs_debuffs_auto_gen';
import { IconPickerStatOption, inDisplayOrder } from './stat_options';

// Every buff the client database resolves to a spell has its input generated from the manifest;
// the rows below are the ones it cannot produce, and the registries at the end of this file
// interleave the two.
export * from './buffs_debuffs_auto_gen';

// The generated const takes its name from the proto field; the specs name the buff.
export const Innervate = Generated.Innervates;
export const PowerInfusion = Generated.PowerInfusions;
export const ManaTideTotem = Generated.ManaTideTotems;

// Party Buffs
// A sim toggle rather than a buff of its own: it says the warrior who shouted for the party wears
// three pieces of Battlegear of Wrath, which is worth 30 more attack power on the shout.
export const EnhancedBattleShout = makeBooleanPartyBuffInput({
	actionId: ActionId.fromSpellId(23563),
	fieldName: 'snapshotBsT2',
	label: 'Enhanced Battle Shout',
	enableWhen: party => party.getBuffs().battleShout,
});

// Individual Buffs
export const GreaterBlessingOfSalvation = makeBooleanIndividualBuffInput({
	actionId: ActionId.fromSpellId(25895),
	fieldName: 'greaterBlessingOfSalvation',
	label: 'Greater Blessing of Salvation',
	showWhen: player => !player.getPlayerSpec().isTankSpec && !player.getPlayerSpec().isHealingSpec,
});

export const PARTY_BUFFS_CONFIG = inDisplayOrder(Generated.GENERATED_PARTY_BUFFS_CONFIG, [
	Generated.BloodPact,
	Generated.BattleShout,
	{ config: EnhancedBattleShout, stats: [Stat.StatAttackPower], ownerClass: Class.ClassWarrior },
	Generated.DevotionAura,
	Generated.LeaderOfThePack,
	Generated.ManaSpringTotem,
	Generated.ManaTideTotems,
	Generated.MoonkinAura,
	Generated.RetributionAura,
	Generated.ConcentrationAura,
	Generated.TrueshotAura,
	Generated.AtieshMage,
	Generated.AtieshWarlock,
	Generated.StrengthOfEarthTotem,
	Generated.GraceOfAirTotem,
	Generated.WindfuryTotem,
]);

export const BUFFS_CONFIG = inDisplayOrder(
	[...Generated.GENERATED_RAID_BUFFS_CONFIG, ...Generated.GENERATED_INDIVIDUAL_BUFFS_CONFIG],
	[
		Generated.ArcaneBrilliance,
		Generated.GreaterBlessingOfKings,
		Generated.PrayerOfSpirit,
		Generated.GiftOfTheWild,
		Generated.Thorns,
		Generated.PrayerOfFortitude,
		Generated.GreaterBlessingOfMight,
		Generated.GreaterBlessingOfWisdom,
		{ config: GreaterBlessingOfSalvation, stats: [] },
		Generated.PrayerOfShadowProtection,
		Generated.FireResistanceAura,
		Generated.FrostResistanceAura,
		Generated.ShadowResistanceAura,
		Generated.FireResistanceTotem,
		Generated.FrostResistanceTotem,
		Generated.NatureResistanceTotem,
		Generated.AspectOfTheWild,
		Generated.Innervates,
		Generated.PowerInfusions,
	],
);

export const DEBUFFS_CONFIG = inDisplayOrder(Generated.GENERATED_DEBUFFS_CONFIG, [
	Generated.HuntersMark,
	Generated.ImprovedSealOfTheCrusader,
	Generated.JudgementOfLight,
	Generated.JudgementOfWisdom,
	Generated.Mangle,
	Generated.CurseOfElements,
	Generated.CurseOfRecklessness,
	Generated.FaerieFire,
	Generated.ExposeArmor,
	Generated.SunderArmor,
	Generated.GiftOfArthas,
	Generated.DemoralizingRoar,
	Generated.DemoralizingShout,
	Generated.ThunderClap,
	Generated.InsectSwarm,
	Generated.ScorpidSting,
]);

export const DEBUFFS_MISC_CONFIG = [] as IconPickerStatOption[];
