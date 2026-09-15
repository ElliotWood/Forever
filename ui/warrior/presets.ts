import { ClassicPhase } from '../core/constants/other.js';
import * as PresetUtils from '../core/preset_utils.js';
import {
	AgilityElixir,
	Alcohol,
	ArmorElixir,
	AttackPowerBuff,
	Consumes,
	Debuffs,
	Food,
	HealthElixir,
	IndividualBuffs,
	Potions,
	Profession,
	Race,
	RaidBuffs,
	SapperExplosive,
	StrengthBuff,
	TristateEffect,
	WeaponImbue,
	ZanzaBuff,
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';
import { Warrior_Options as WarriorOptions, WarriorShout, WarriorStance } from '../core/proto/warrior.js';
import APLNoReckJSON from './apls/dps_no_reck.apl.json';
import APLReckJSON from './apls/dps_reck.apl.json';
import LaunchGearJSON from './gear_sets/launch.gear.json';
import P0BISGear from './gear_sets/p0.bis.gear.json';
import Phase1Gear from './gear_sets/phase_1.gear.json';
import Phase2Gear from './gear_sets/phase_2.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

export const GearLaunch = PresetUtils.makePresetGear('Launch', LaunchGearJSON);
export const GearP0BIS = PresetUtils.makePresetGear('Pre-BiS', P0BISGear);
export const GearPhase1 = PresetUtils.makePresetGear('P1 BiS', Phase1Gear);
export const GearPhase2 = PresetUtils.makePresetGear('P2 BiS', Phase2Gear);

export const GearPresets = {
	[ClassicPhase.Phase1]: [GearLaunch, GearPhase1, GearP0BIS],
	[ClassicPhase.Phase2]: [GearPhase2],
};

export const DefaultGear = GearP0BIS;

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

export const AplReck = PresetUtils.makePresetAPLRotation('DPS (With Reck)', APLReckJSON);
export const APLNoReck = PresetUtils.makePresetAPLRotation('DPS (No Reck)', APLNoReckJSON);

export const APLPresets = {
	[ClassicPhase.Phase1]: [APLNoReck, AplReck],
};

export const DefaultAPLs = [APLPresets[ClassicPhase.Phase1][0]];

///////////////////////////////////////////////////////////////////////////
//                                 Talent Presets
///////////////////////////////////////////////////////////////////////////

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.

export const TalentsP1DPS = PresetUtils.makePresetTalents('DPS', SavedTalents.create({ talentsString: '303050130-050520035150310051' }));

export const TalentsFury = PresetUtils.makePresetTalents('Fury 17/34/0', SavedTalents.create({ talentsString: '30305213-550501015050010051' }));
export const TalentsArms = PresetUtils.makePresetTalents('Arms 39/12/0', SavedTalents.create({ talentsString: '32305213132515201-5502' }));

export const TalentPresets = {
	[ClassicPhase.Phase1]: [TalentsP1DPS, TalentsFury, TalentsArms],
};

export const DefaultTalents = TalentPresets[ClassicPhase.Phase1][0];

///////////////////////////////////////////////////////////////////////////
//                                 Options Presets
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = WarriorOptions.create({
	queueDelay: 250,
	startingRage: 0,
	shout: WarriorShout.WarriorShoutBattle,
	stance: WarriorStance.WarriorStanceBerserker,
});

export const DefaultConsumes = Consumes.create({
	agilityElixir: AgilityElixir.ElixirOfTheMongoose,
	alcohol: Alcohol.AlcoholRumseyRumBlackLabel,
	armorElixir: ArmorElixir.ElixirOfSuperiorDefense,
	attackPowerBuff: AttackPowerBuff.JujuMight,
	defaultPotion: Potions.MightyRagePotion,
	dragonBreathChili: true,
	food: Food.FoodSmokedDesertDumpling,
	healthElixir: HealthElixir.ElixirOfFortitude,
	mainHandImbue: WeaponImbue.Windfury,
	offHandImbue: WeaponImbue.ElementalSharpeningStone,
	sapperExplosive: SapperExplosive.SapperGoblinSapper,
	strengthBuff: StrengthBuff.JujuPower,
	zanzaBuff: ZanzaBuff.ROIDS,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	battleShout: TristateEffect.TristateEffectImproved,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	graceOfAirTotem: TristateEffect.TristateEffectImproved,
	leaderOfThePack: true,
	strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfMight: TristateEffect.TristateEffectImproved,
});

export const DefaultDebuffs = Debuffs.create({
	curseOfRecklessness: true,
	exposeArmor: TristateEffect.TristateEffectImproved,
	faerieFire: true,
	giftOfArthas: true,
	sunderArmor: true,
});

export const OtherDefaults = {
	profession1: Profession.Alchemy,
	profession2: Profession.Engineering,
	race: Race.RaceHuman,
};
