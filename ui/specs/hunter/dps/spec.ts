import * as other_inputs from '@features/settings/model/other_inputs';
import { StatCapType } from '@generated/proto/api';
import { APLRotation, APLRotation_Type } from '@generated/proto/apl';
import { ItemSlot, PseudoStat, Spec, Stat } from '@generated/proto/common';
import { SavedTalents } from '@generated/proto/ui';
import { PlayerClasses } from '@sim/player/classes';
import { Player } from '@sim/player/player';
import { StatCap, Stats, UnitStat } from '@sim/proto/stats';
import { defineSpec } from '@sim/spec_config';

import * as HunterInputs from './inputs';
import * as Presets from './presets';

export default defineSpec<Spec.SpecHunter>({
	spec: Spec.SpecHunter,

	className: 'hunter-sim-ui',
	cssScheme: PlayerClasses.getCssScheme(PlayerClasses.Hunter),
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [],
	warnings: [],
	// All stats for which EP should be calculated.
	epStats: [
		Stat.StatAgility,
		Stat.StatStrength,
		Stat.StatIntellect,
		Stat.StatMP5,
		Stat.StatAttackPower,
		Stat.StatRangedAttackPower,
		Stat.StatArmorPenetration,
		Stat.StatMeleeHitRating,
		Stat.StatMeleeHasteRating,
		Stat.StatMeleeCritRating,
		Stat.StatExpertiseRating,
		Stat.StatPhysicalDamage,
	],
	gemStats: [Stat.StatStamina, Stat.StatAgility],
	epPseudoStats: [PseudoStat.PseudoStatRangedHitPercent, PseudoStat.PseudoStatRangedCritPercent, PseudoStat.PseudoStatRangedDps],
	consumableStats: [Stat.StatStamina, Stat.StatHealth, Stat.StatMana],
	// Reference stat against which to calculate EP.
	epReferenceStat: Stat.StatAgility,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: UnitStat.createDisplayStatArray(
		[
			Stat.StatHealth,
			Stat.StatMana,
			Stat.StatStamina,
			Stat.StatStrength,
			Stat.StatAgility,
			Stat.StatIntellect,
			Stat.StatMP5,
			Stat.StatAttackPower,
			Stat.StatRangedAttackPower,
			Stat.StatExpertiseRating,
			Stat.StatArmorPenetration,
			Stat.StatArcaneResistance,
			Stat.StatFireResistance,
			Stat.StatFrostResistance,
			Stat.StatNatureResistance,
			Stat.StatShadowResistance,
		],
		[
			PseudoStat.PseudoStatMeleeHitPercent,
			PseudoStat.PseudoStatMeleeCritPercent,
			PseudoStat.PseudoStatMeleeHastePercent,
			PseudoStat.PseudoStatRangedHitPercent,
			PseudoStat.PseudoStatRangedCritPercent,
			PseudoStat.PseudoStatRangedHastePercent,
		],
	),
	itemSwapSlots: [ItemSlot.ItemSlotMainHand, ItemSlot.ItemSlotOffHand, ItemSlot.ItemSlotRanged, ItemSlot.ItemSlotTrinket1, ItemSlot.ItemSlotTrinket2],
	defaults: {
		// Default equipped gear.
		gear: Presets.DEFAULT_GEAR.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: new Stats(),
		softCapBreakpoints: [
			StatCap.fromPseudoStat(PseudoStat.PseudoStatRangedHitPercent, {
				breakpoints: [9],
				capType: StatCapType.TypeSoftCap,
				postCapEPs: [0],
			}),
		],
		rotationType: APLRotation_Type.TypeAPL,
		other: Presets.OtherDefaults,
		// Default consumes settings.
		consumables: Presets.DefaultConsumables,
		// Default talents.
		talents: SavedTalents.create({
			talentsString: '-3050552301503151-50024001',
		}),
		// Default spec-specific settings.
		specOptions: Presets.DefaultOptions,
		// Default raid/party buffs settings.
		raidBuffs: Presets.DefaultRaidBuffs,
		partyBuffs: Presets.DefaultPartyBuffs,
		individualBuffs: Presets.DefaultIndividualBuffs,
		debuffs: Presets.DefaultDebuffs,
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [HunterInputs.PetTypeInput(), HunterInputs.QuiverInput(), HunterInputs.AmmoInput()],
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [Stat.StatSpirit, Stat.StatSpellCritRating, Stat.StatSpellDamage],
	excludeBuffDebuffInputs: [],
	rotationInputs: HunterInputs.RotationInputs,
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			other_inputs.TotemTwisting,
			HunterInputs.PetUptime(),
			HunterInputs.PetSingleAbility(),
			HunterInputs.PetAttackSpeedInput(),
			other_inputs.InputDelay,
			other_inputs.DistanceFromTarget,
			other_inputs.TankAssignment,
		],
	},
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		epWeights: [],
		// Preset talents that the user can quickly select.
		talents: Presets.TalentPresets,
		// Preset rotations that the user can quickly select.
		rotations: [Presets.BeastMasteryRotation, Presets.MarksmanshipRotation, Presets.SurvivalRotation],
		// Preset gear configurations that the user can quickly select.
		gear: Presets.GEAR_PRESETS,
	},

	// The build decides the rotation: Sniper Shot is Marksmanship's capstone and Summon Hawk is
	// Beast Mastery's, so each tree gets the list that uses its own.
	autoRotation: (player: Player<Spec.SpecHunter>): APLRotation => {
		const talents = player.getTalents();
		if (talents.sniperShot) return APLRotation.clone(Presets.MarksmanshipRotation.rotation.rotation!);
		if (talents.summonHawk) return APLRotation.clone(Presets.BeastMasteryRotation.rotation.rotation!);
		return APLRotation.clone(Presets.SurvivalRotation.rotation.rotation!);
	},

	reforge: {},
});
