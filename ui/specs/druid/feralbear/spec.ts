import * as OtherInputs from '@features/settings/model/other_inputs';
import { APLAction, APLListItem, APLRotation, APLRotation_Type as APLRotationType } from '@generated/proto/apl';
import { Cooldowns, Debuffs, IndividualBuffs, ItemSlot, PartyBuffs, PseudoStat, RaidBuffs, Spec, Stat, TristateEffect } from '@generated/proto/common';
import { FeralBearDruid_Rotation as DruidRotation } from '@generated/proto/druid';
import * as Mechanics from '@sim/constants/mechanics';
import { PlayerClasses } from '@sim/player/classes';
import { Player } from '@sim/player/player';
import { Stats, UnitStat } from '@sim/proto/stats';
import { defineSpec } from '@sim/spec_config';

import * as FeralBearInputs from './inputs';
import * as Presets from './presets';

export default defineSpec<Spec.SpecFeralBearDruid>({
	spec: Spec.SpecFeralBearDruid,
	enableHealing: true,

	className: 'feral-bear-druid-sim-ui',
	cssScheme: PlayerClasses.getCssScheme(PlayerClasses.Druid),
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [],
	warnings: [],

	epRatios: [0, 0, 0.6, 0, 1.0, 0],
	// All stats for which EP should be calculated.
	epStats: [
		Stat.StatStamina,
		Stat.StatAgility,
		Stat.StatStrength,
		Stat.StatAttackPower,
		Stat.StatFeralAttackPower,
		Stat.StatArmor,
		Stat.StatBonusArmor,
		Stat.StatDodgeRating,
		Stat.StatDefenseRating,
		Stat.StatMeleeHitRating,
		Stat.StatMeleeCritRating,
		Stat.StatMeleeHasteRating,
		Stat.StatExpertiseRating,
		Stat.StatResilienceRating,
		Stat.StatPhysicalDamage,
		Stat.StatArmorPenetration,
	],
	epPseudoStats: [],
	epReferenceStat: Stat.StatAgility,
	tankRefStat: Stat.StatStamina,
	displayStats: UnitStat.createDisplayStatArray(
		[
			Stat.StatHealth,
			Stat.StatStamina,
			Stat.StatAgility,
			Stat.StatStrength,
			Stat.StatAttackPower,
			Stat.StatArmor,
			Stat.StatBonusArmor,
			Stat.StatDodgeRating,
			Stat.StatDefenseRating,
			Stat.StatExpertiseRating,
			Stat.StatResilienceRating,
			Stat.StatNatureResistance,
			Stat.StatFireResistance,
			Stat.StatFrostResistance,
			Stat.StatArcaneResistance,
			Stat.StatShadowResistance,
		],
		[
			PseudoStat.PseudoStatMeleeHitPercent,
			PseudoStat.PseudoStatMeleeCritPercent,
			PseudoStat.PseudoStatMeleeHastePercent,
			PseudoStat.PseudoStatDodgePercent,
		],
	),

	defaults: {
		gear: Presets.DEFAULT_GEAR.gear,
		epWeights: new Stats(),
		statCaps: (() => {
			const hitCap = new Stats().withPseudoStat(PseudoStat.PseudoStatMeleeHitPercent, 9);
			const expCap = new Stats().withStat(Stat.StatExpertiseRating, 6.5 * 4 * Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION);
			const critImmunityCap = new Stats().withPseudoStat(PseudoStat.PseudoStatReducedCritTakenPercent, 5.6);
			return hitCap.add(expCap).add(critImmunityCap);
		})(),
		other: Presets.OtherDefaults,
		consumables: Presets.DefaultConsumables,
		rotationType: APLRotationType.TypeAPL,
		aplRotation: Presets.ROTATION_DEFAULT.rotation.rotation!,
		talents: Presets.BearTankTalents.data,
		specOptions: Presets.DefaultOptions,
		// Master's page (currentSettings on a fresh profile); its raid-wide totems and Battle Shout
		// are party buffs here, its Stoneskin Totem has no counterpart.
		raidBuffs: RaidBuffs.create({
			giftOfTheWild: TristateEffect.TristateEffectImproved,
			powerWordFortitude: TristateEffect.TristateEffectImproved,
		}),
		partyBuffs: PartyBuffs.create({
			battleShout: TristateEffect.TristateEffectImproved,
			fireResistanceTotem: true,
			graceOfAirTotem: TristateEffect.TristateEffectImproved,
			strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
		}),
		individualBuffs: IndividualBuffs.create({}),
		debuffs: Debuffs.create({
			curseOfRecklessness: true,
			exposeArmor: TristateEffect.TristateEffectImproved,
			faerieFire: TristateEffect.TristateEffectRegular,
			giftOfArthas: true,
			sunderArmor: true,
		}),
	},

	playerIconInputs: [],
	rotationInputs: FeralBearInputs.FeralBearRotationConfig,
	includeBuffDebuffInputs: [Stat.StatStamina, Stat.StatArmor],
	excludeBuffDebuffInputs: [Stat.StatParryRating],
	otherInputs: {
		inputs: [
			OtherInputs.TotemTwisting,
			FeralBearInputs.StartingRage,
			OtherInputs.InputDelay,
			OtherInputs.TankAssignment,
			OtherInputs.InspirationUptime,
			OtherInputs.IncomingHps,
			OtherInputs.HealingCadence,
			OtherInputs.HealingCadenceVariation,
			OtherInputs.AbsorbFrac,
			OtherInputs.BurstWindow,
			OtherInputs.HpPercentForDefensives,
			OtherInputs.InFrontOfTarget,
		],
	},
	itemSwapSlots: [ItemSlot.ItemSlotTrinket1, ItemSlot.ItemSlotTrinket2, ItemSlot.ItemSlotMainHand, ItemSlot.ItemSlotRanged],

	encounterPicker: {
		showExecuteProportion: false,
	},

	presets: {
		epWeights: [],
		talents: [Presets.BearTankTalents],
		// ROTATION_SIMPLE is kept in presets.ts for reference but omitted here —
		// the APL rotation is more user-friendly and handles CDs, re-shifting, and
		// on-use items more easily.
		rotations: [Presets.ROTATION_DEFAULT],
		gear: Presets.GEAR_PRESETS,
	},

	autoRotation: (_player: Player<Spec.SpecFeralBearDruid>): APLRotation => {
		return Presets.ROTATION_DEFAULT.rotation.rotation!;
	},

	simpleRotation: (_player: Player<Spec.SpecFeralBearDruid>, simple: DruidRotation, _cooldowns: Cooldowns): APLRotation => {
		const doRotation = APLAction.fromJsonString(
			`{"bearOptimalRotationAction":{"maintainFaerieFire":${simple.maintainFaerieFire},"maintainDemoralizingRoar":${simple.maintainDemoralizingRoar},"maulRageThreshold":${simple.maulRageThreshold},"swipeUsage":${simple.swipeUsage},"swipeApThreshold":${simple.swipeApThreshold}}}`,
		);
		return APLRotation.create({
			priorityList: [APLListItem.create({ action: doRotation })],
		});
	},

	reforge: {},
});
