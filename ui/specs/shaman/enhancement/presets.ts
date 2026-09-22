import * as PresetUtils from '@app/preset_utils';
import { ConsumesSpec, Debuffs, IndividualBuffs, PartyBuffs, Profession, Race, RaidBuffs, TristateEffect } from '@generated/proto/common';
import { EnhancementShaman_Options as EnhancementShamanOptions, ShamanImbue, ShamanSyncType } from '@generated/proto/shaman';
import { SavedTalents } from '@generated/proto/ui';

import ForeverApl from './apls/forever.apl.json';
import LaunchGear from './gear_sets/launch.gear.json';
import Phase1Gear from './gear_sets/phase_1.gear.json';
import Phase2Gear from './gear_sets/phase_2.gear.json';

// Our Forever APL (the parity tool's nextApl). The TBC default.apl.json it replaces never cast
// Earth Shock or Lightning Bolt (459 DPS at defaults, master 576).
export const ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Default', ForeverApl);

// Defaults below are what master's ui/enhancement_shaman opens with (currentSettings on a fresh
// profile); its raid-wide Battle Shout, Leader of the Pack and totems are party buffs here.
export const DefaultIndividualBuffs = IndividualBuffs.create({});

export const DefaultOptions = EnhancementShamanOptions.create({
	classOptions: {
		shieldProcrate: 0,
		imbueMh: ShamanImbue.WindfuryWeapon,
	},
	imbueOh: ShamanImbue.WindfuryWeapon,
	syncType: ShamanSyncType.Auto,
});

export const OtherDefaults = {
	reactionTime: 200, // master's default
	distanceFromTarget: 5,
	profession1: Profession.Alchemy,
	profession2: Profession.Enchanting,
	race: Race.RaceOrc,
};

// Master's consumables; Juju Power/Might, R.O.I.D.S., Dragonbreath Chili, Greater Arcane Elixir and
// Elixir of Firepower have no slot here (one battle elixir), Blessed Sunfruit, Mageblood and Major
// Mana Potion are not in this db.
export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 13512, // Flask of Supreme Power
	battleElixirId: 13452, // Elixir of the Mongoose
	guardianElixirId: 20007, // Mageblood Elixir
	foodId: 13810, // Blessed Sunfruit
	potId: 13444, // Major Mana Potion
	conjuredId: 12662, // Demonic Rune
});

export const DefaultPartyBuffs = PartyBuffs.create({
	battleShout: TristateEffect.TristateEffectImproved,
	fireResistanceTotem: true,
	leaderOfThePack: TristateEffect.TristateEffectRegular,
	manaSpringTotem: TristateEffect.TristateEffectRegular,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	divineSpirit: TristateEffect.TristateEffectRegular,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	powerWordFortitude: TristateEffect.TristateEffectImproved,
});

export const DefaultDebuffs = Debuffs.create({
	curseOfRecklessness: true,
	exposeArmor: TristateEffect.TristateEffectImproved,
	faerieFire: TristateEffect.TristateEffectRegular,
	sunderArmor: true,
});

// Talent presets, from master's ui/shaman spec.
export const TalentsLevel60 = PresetUtils.makePresetTalents('Level 60', SavedTalents.create({ talentsString: '5505301-053030031005112251' }));
export const TalentsEnhancement = PresetUtils.makePresetTalents('Enhancement 16/35/0', SavedTalents.create({ talentsString: '05023015-055030030205112251' }));
export const TalentPresets = [TalentsLevel60, TalentsEnhancement];

// Our Forever sim's gear presets (master ui/<spec>/gear_sets).
export const GEAR_LAUNCH = PresetUtils.makePresetGear('Launch', LaunchGear);
export const GEAR_PHASE_1 = PresetUtils.makePresetGear('Phase 1', Phase1Gear);
export const GEAR_PHASE_2 = PresetUtils.makePresetGear('Phase 2', Phase2Gear);
export const DEFAULT_GEAR = GEAR_LAUNCH;
export const GEAR_PRESETS = [GEAR_LAUNCH, GEAR_PHASE_1, GEAR_PHASE_2];
