import { Debuffs, IndividualBuffs, PartyBuffs, RaidBuffs } from '@generated/proto/buffs';
import { ConsumesSpec } from '@generated/proto/common';

export const DefaultIndividualBuffs = IndividualBuffs.create({
	greaterBlessingOfKings: true,
	greaterBlessingOfMight: true,
});

export const DefaultPartyBuffs = PartyBuffs.create({
	graceOfAirTotem: true,
	strengthOfEarthTotem: true,
	windfuryTotem: true,
	leaderOfThePack: true,
	totemTwisting: true,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	prayerOfFortitude: true,
	giftOfTheWild: true,
});

export const DefaultDebuffs = Debuffs.create({
	improvedSealOfTheCrusader: true,
	giftOfArthas: true,
	mangle: true,
	exposeArmor: true,
	faerieFire: true,
	sunderArmor: true,
	curseOfRecklessness: true,
	huntersMark: true,
});

export const DefaultConsumables = ConsumesSpec.create({
	potId: 22838,
	flaskId: 22854,
	foodId: 27658,
	conjuredId: 22788,
	explosiveId: 30217,
	superSapper: true,
	goblinSapper: true,
	ohImbueId: 29453,
	scrollAgi: true,
	scrollStr: true,
});
