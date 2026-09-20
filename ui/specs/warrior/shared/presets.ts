import { Debuffs, IndividualBuffs, PartyBuffs, RaidBuffs } from '@generated/proto/buffs';
import { Class, ConsumesSpec, Drums } from '@generated/proto/common';
import { defaultExposeWeaknessSettings, defaultRaidBuffMajorDamageCooldowns } from '@sim/proto/utils';

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfMight: true,
	unleashedRage: true,
});

export const DefaultPartyBuffs = PartyBuffs.create({
	ferociousInspiration: 2,
	braidedEterniumChain: true,
	graceOfAirTotem: true,
	strengthOfEarthTotem: true,
	windfuryTotem: true,
	leaderOfThePack: true,
	totemTwisting: true,
	drums: Drums.LesserDrumsOfBattle,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	...defaultRaidBuffMajorDamageCooldowns(Class.ClassWarrior),
	powerWordFortitude: true,
	giftOfTheWild: true,
});

export const DefaultDebuffs = Debuffs.create({
	...defaultExposeWeaknessSettings(),
	improvedSealOfTheCrusader: true,
	misery: true,
	bloodFrenzy: true,
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
