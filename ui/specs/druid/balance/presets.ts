import * as PresetUtils from '@app/preset_utils';
import { Debuffs, IndividualBuffs, PartyBuffs, RaidBuffs } from '@generated/proto/buffs';
import { ConsumesSpec, Profession, Race, UnitReference } from '@generated/proto/common';
import { BalanceDruid_Options as BalanceDruidOptions } from '@generated/proto/druid';

import DefaultAPL from './apls/default.apl.json';

export const StandardRotation = PresetUtils.makePresetAPLRotation('Default', DefaultAPL);

export const DefaultOptions = BalanceDruidOptions.create({
	classOptions: {
		innervateTarget: UnitReference.create(),
	},
});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	giftOfTheWild: true,
	powerWordFortitude: true,
	divineSpirit: true,
});

export const DefaultPartyBuffs = PartyBuffs.create({});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfWisdom: true,
	shadowPriestDps: 800,
});

export const DefaultDebuffs = Debuffs.create({
	curseOfElements: true,
	curseOfRecklessness: true,
	exposeArmor: true,
	giftOfArthas: true,
	huntersMark: true,
	improvedSealOfTheCrusader: true,
	judgementOfWisdom: true,
	mangle: true,
	sunderArmor: true,
});

export const DefaultConsumables = ConsumesSpec.create({
	conjuredId: 12662, // Demonic Rune
	flaskId: 22861, // Flask of Blinding Light
	foodId: 27657, // Blackened Basilisk
	mhImbueId: 25122, // Brilliant Wizard Oil
	potId: 22832, // Super Mana Potion
});

export const OtherDefaults = {
	distanceFromTarget: 20,
	profession1: Profession.Enchanting,
	profession2: Profession.Tailoring,
	race: Race.RaceNightElf,
};
