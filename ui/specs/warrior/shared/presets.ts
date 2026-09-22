import { ConsumesSpec, Debuffs, IndividualBuffs, PartyBuffs, TristateEffect } from '@generated/proto/common';

// Defaults follow master's ui/warrior and ui/tank_warrior.
export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfMight: TristateEffect.TristateEffectImproved,
});

// Master's page (currentSettings on a fresh profile) has Battle Shout and Leader of the Pack as
// raid buffs; they are party buffs here. No totems.
export const DefaultPartyBuffs = PartyBuffs.create({
	battleShout: TristateEffect.TristateEffectImproved,
	leaderOfThePack: TristateEffect.TristateEffectRegular,
});

export const DefaultDebuffs = Debuffs.create({
	curseOfRecklessness: true,
	exposeArmor: TristateEffect.TristateEffectImproved,
	faerieFire: TristateEffect.TristateEffectRegular,
	giftOfArthas: true,
	sunderArmor: true,
});

// Master's consumables. Juju Power/Might, R.O.I.D.S., Dragonbreath Chili and Rumsey Rum have no
// field here; Smoked Desert Dumplings and Elixir of Fortitude are not in this db (no effect).
export const DefaultConsumables = ConsumesSpec.create({
	battleElixirId: 13452, // Elixir of the Mongoose
	guardianElixirId: 3825, // Elixir of Fortitude
	foodId: 20452, // Smoked Desert Dumplings
	potId: 13442, // Mighty Rage Potion
	ohImbueId: 18262, // Elemental Sharpening Stone
	goblinSapper: true,
});
