import * as OtherInputs from '@features/settings/model/other_inputs';
import { APLAction, APLListItem, APLRotation, APLRotation_Type as APLRotationType } from '@generated/proto/apl';
import { Cooldowns, Debuffs, IndividualBuffs, ItemSlot, PartyBuffs, PseudoStat, RaidBuffs, Spec, Stat, TristateEffect } from '@generated/proto/common';
import { FeralCatDruid_Rotation as DruidRotation } from '@generated/proto/druid';
import * as Mechanics from '@sim/constants/mechanics';
import { PlayerClasses } from '@sim/player/classes';
import { Player } from '@sim/player/player';
import { masterEpWeights } from '@sim/proto/master_ep_weights';
import { Stats, UnitStat } from '@sim/proto/stats';
import { defineSpec } from '@sim/spec_config';

import * as FeralInputs from './inputs';
import * as Presets from './presets';

export default defineSpec<Spec.SpecFeralCatDruid>({
	spec: Spec.SpecFeralCatDruid,

	className: 'feral-druid-sim-ui',
	cssScheme: PlayerClasses.getCssScheme(PlayerClasses.Druid),
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [],
	warnings: [],

	// Extra stats shown in consumables picker (beyond epStats).
	consumableStats: [Stat.StatMeleeHasteRating, Stat.StatMana, Stat.StatSpirit, Stat.StatMP5],

	// All stats for which EP should be calculated.
	epStats: [
		Stat.StatStrength,
		Stat.StatAgility,
		Stat.StatAttackPower,
		Stat.StatFeralAttackPower,
		Stat.StatPhysicalDamage,
		Stat.StatMeleeHitRating,
		Stat.StatExpertiseRating,
		Stat.StatMeleeCritRating,
		Stat.StatMeleeHasteRating,
		Stat.StatArmorPenetration,
	],
	gemStats: [Stat.StatAgility, Stat.StatStrength],
	epPseudoStats: [],
	// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatAgility,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: UnitStat.createDisplayStatArray(
		[
			Stat.StatHealth,
			Stat.StatStrength,
			Stat.StatAgility,
			Stat.StatStamina,
			Stat.StatIntellect,
			Stat.StatSpirit,
			Stat.StatAttackPower,
			Stat.StatMana,
			Stat.StatExpertiseRating,
			Stat.StatArmorPenetration,
			Stat.StatArcaneResistance,
			Stat.StatFireResistance,
			Stat.StatFrostResistance,
			Stat.StatNatureResistance,
			Stat.StatShadowResistance,
		],
		[PseudoStat.PseudoStatMeleeHitPercent, PseudoStat.PseudoStatMeleeCritPercent, PseudoStat.PseudoStatMeleeHastePercent],
	),

	defaults: {
		// Default equipped gear.
		gear: Presets.DEFAULT_GEAR.gear,
		// Default EP weights for sorting gear in the gear picker.
		// Master's weights (percent stats per 1%), in our ratings.
		epWeights: masterEpWeights({
			Strength: 2.4,
			Agility: 2.43,
			Intellect: 0.61,
			Spirit: 0.38,
			MP5: 0.79,
			AttackPower: 1,
			MeleeHit: 26.59,
			MeleeCrit: 28.68,
			Expertise: 26.59,
			Mana: 0.03,
			FeralAttackPower: 1,
			BonusPhysicalDamage: 13.33,
			MeleeSpeedMultiplier: 16.5,
		}),
		statCaps: (() => {
			return new Stats()
				.withPseudoStat(PseudoStat.PseudoStatMeleeHitPercent, 9)
				.withStat(Stat.StatExpertiseRating, 6.5 * 4 * Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION);
		})(),
		other: Presets.OtherDefaults,
		// Default consumes settings.
		consumables: Presets.DefaultConsumables,
		// Default rotation settings.
		rotationType: APLRotationType.TypeAuto,
		// Default talents.
		talents: Presets.FeralTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultOptions,
		// Default raid/party buffs settings.
		// Master's page (currentSettings on a fresh profile); its raid-wide Strength of Earth,
		// Battle Shout, Leader of the Pack and Mana Spring are party buffs here.
		raidBuffs: RaidBuffs.create({
			arcaneBrilliance: true,
			divineSpirit: TristateEffect.TristateEffectRegular,
			giftOfTheWild: TristateEffect.TristateEffectImproved,
			powerWordFortitude: TristateEffect.TristateEffectImproved,
		}),
		partyBuffs: PartyBuffs.create({
			battleShout: TristateEffect.TristateEffectImproved,
			leaderOfThePack: TristateEffect.TristateEffectRegular,
			manaSpringTotem: TristateEffect.TristateEffectRegular,
			strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
		}),
		individualBuffs: IndividualBuffs.create({}),
		debuffs: Debuffs.create({
			exposeArmor: TristateEffect.TristateEffectImproved,
			curseOfRecklessness: true,
			sunderArmor: true,
		}),
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [],
	// Inputs to include in the 'Rotation' section on the settings tab.
	rotationInputs: FeralInputs.FeralDruidRotationConfig,
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [Stat.StatMP5, Stat.StatIntellect, Stat.StatStamina],
	excludeBuffDebuffInputs: [Stat.StatParryRating],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			OtherInputs.TotemTwisting,
			OtherInputs.InputDelay,
			OtherInputs.DistanceFromTarget,
			OtherInputs.TankAssignment,
			OtherInputs.InFrontOfTarget,
			FeralInputs.CannotShredTarget,
		],
	},
	itemSwapSlots: [ItemSlot.ItemSlotMainHand, ItemSlot.ItemSlotOffHand, ItemSlot.ItemSlotTrinket1, ItemSlot.ItemSlotTrinket2],
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: true,
	},

	presets: {
		epWeights: [],
		// Preset talents that the user can quickly select.
		talents: [Presets.FeralTalents, Presets.FeralCatTalents],
		rotations: [Presets.SIMPLE, Presets.APL],
		// Preset gear configurations that the user can quickly select.
		gear: Presets.GEAR_PRESETS,
	},

	autoRotation: (_player: Player<Spec.SpecFeralCatDruid>): APLRotation => {
		return Presets.APL.rotation.rotation!;
	},

	simpleRotation: (_player: Player<Spec.SpecFeralCatDruid>, simple: DruidRotation, _cooldowns: Cooldowns): APLRotation => {
		// All cooldowns (potions, sappers, runes) are fired by the Go rotation
		// during power shifts (ClearForm → fire MCDs → CatForm), not through the APL.
		const doRotation = APLAction.fromJsonString(
			`{"catOptimalRotationAction":{"finishingMove":${simple.finishingMove},"biteweave":${simple.biteweave},"ripweave":${simple.ripweave},"ripMinComboPoints":${simple.ripMinComboPoints},"biteMinComboPoints":${simple.biteMinComboPoints},"mangleTrick":${simple.mangleTrick},"rakeTrick":${simple.rakeTrick},"maintainFaerieFire":${simple.maintainFaerieFire}}}`,
		);

		return APLRotation.create({
			priorityList: [APLListItem.create({ action: doRotation })],
		});
	},

	hiddenMCDs: [],

	reforge: {},
});
