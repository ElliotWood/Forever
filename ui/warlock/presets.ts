import { Player } from '../core/player.js';
import * as PresetUtils from '../core/preset_utils.js';
import {
	Alcohol,
	Conjured,
	Consumes,
	Debuffs,
	FirePowerBuff,
	Flask,
	Food,
	IndividualBuffs,
	ManaRegenElixir,
	Potions,
	Profession,
	RaidBuffs,
	SaygesFortune,
	ShadowPowerBuff,
	SpellPowerBuff,
	TristateEffect,
	WeaponImbue,
	ZanzaBuff,
} from '../core/proto/common';
import { SavedTalents } from '../core/proto/ui.js';
import {
	WarlockOptions as WarlockOptions,
	WarlockOptions_Armor as Armor,
	WarlockOptions_Summon as Summon,
	WarlockOptions_WeaponImbue as WarlockWeaponImbue,
} from '../core/proto/warlock.js';
// apls
import AfflictionApl from './apls/forever_affliction.apl.json';
import DSRuinApl from './apls/forever_ds_ruin.apl.json';
import DemonicPactApl from './apls/forever_pact.apl.json';
import ShadowAndFlameApl from './apls/forever_shadow_and_flame.apl.json';
// gear
import LaunchGearJSON from './gear_sets/launch.gear.json';
import BlankGear from './gear_sets/blank.gear.json';
import MCGear from './gear_sets/mc.gear.json';
import PreBisGear from './gear_sets/prebis.gear.json';

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

export const GearLaunch = PresetUtils.makePresetGear('Launch', LaunchGearJSON);
export const GearBlank = PresetUtils.makePresetGear('Blank', BlankGear);
export const GearPreBis = PresetUtils.makePresetGear('Pre-BIS', PreBisGear);
export const GearMC = PresetUtils.makePresetGear('MC', MCGear);

export const GearPresets = [GearBlank, GearLaunch, GearPreBis, GearMC];

export const DefaultGear = GearPreBis;

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

// P1
export const RotationDemonicPact = PresetUtils.makePresetAPLRotation('Demonic Pact', DemonicPactApl);
export const RotationAffliction = PresetUtils.makePresetAPLRotation('Affliction', AfflictionApl);
export const RotationDSRuin = PresetUtils.makePresetAPLRotation('DS/Ruin', DSRuinApl);
export const RotationShadowAndFlame = PresetUtils.makePresetAPLRotation('Shadow and Flame', ShadowAndFlameApl);

export const APLPresets = [RotationDemonicPact, RotationAffliction, RotationDSRuin, RotationShadowAndFlame];

export const DefaultAPL = RotationDSRuin;

///////////////////////////////////////////////////////////////////////////
//                                 Talent Presets
///////////////////////////////////////////////////////////////////////////

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.

export const TalentsDemonicPact = {
	name: 'Demonic Pact',
	data: SavedTalents.create({ talentsString: '203-0055003221201001351-0550005' }),
};

export const TalentsAffliction = {
	name: 'Affliction',
	data: SavedTalents.create({ talentsString: '2435002013520135--0500055' }),
};

export const TalentsDSRuin = {
	name: 'DS/Ruin',
	data: SavedTalents.create({ talentsString: '233500201332-0340003001-0550105' }),
};

export const TalentsPactOptimised = PresetUtils.makePresetTalents(
	'Demonic Pact 2/31/18',
	SavedTalents.create({ talentsString: '113-0005003221220311351-0550005' }),
);
export const TalentsDeepAffliction = PresetUtils.makePresetTalents(
	'Deep Affliction 35/0/16',
	SavedTalents.create({ talentsString: '2535002013521105--05000551' }),
);
export const TalentsDSRuinPandemic = PresetUtils.makePresetTalents(
	'DS/Ruin Pandemic 24/11/16',
	SavedTalents.create({ talentsString: '25220010135201-0025003001-05500051' }),
);
export const TalentsShadowAndFlame = PresetUtils.makePresetTalents(
	'Shadow and Flame 13/11/27',
	SavedTalents.create({ talentsString: '25501-0025003001-055035510010002' }),
);

export const TalentPresets = [
	TalentsDemonicPact,
	TalentsAffliction,
	TalentsDSRuin,
	TalentsPactOptimised,
	TalentsDeepAffliction,
	TalentsDSRuinPandemic,
	TalentsShadowAndFlame,
];

export const DefaultTalents = TalentsDSRuin;

///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = WarlockOptions.create({
	armor: Armor.DemonArmor,
	summon: Summon.Imp,
	weaponImbue: WarlockWeaponImbue.NoWeaponImbue,
});

// Without pet talents the Succubus out-damages the Imp, so the Affliction builds run one; a
// sacrificing rotation summons and sacrifices its own Imp, which leaves Shadow damage in Forever.
export const AfflictionOptions = WarlockOptions.create({
	armor: Armor.DemonArmor,
	summon: Summon.Succubus,
	weaponImbue: WarlockWeaponImbue.NoWeaponImbue,
});

// Demonic Pact keeps the sacrifice when another demon is out. The Succubus stays out for Master
// Demonologist and Soul Link, and the Voidwalker is sacrificed for the mana, which Forever moved
// from the Felhunter to the Voidwalker.
export const DemonicPactOptions = WarlockOptions.create({
	armor: Armor.DemonArmor,
	summon: Summon.Succubus,
	sacrifice: Summon.Voidwalker,
	weaponImbue: WarlockWeaponImbue.NoWeaponImbue,
});

export const DefaultConsumes = Consumes.create({
	alcohol: Alcohol.AlcoholRumseyRumBlackLabel,
	defaultPotion: Potions.MajorManaPotion,
	defaultConjured: Conjured.ConjuredDemonicRune,
	flask: Flask.FlaskOfSupremePower,
	firePowerBuff: FirePowerBuff.ElixirOfFirepower,
	food: Food.FoodRunnTumTuberSurprise,
	// mainHandImbue: WeaponImbue.BrilliantWizardOil,
	manaRegenElixir: ManaRegenElixir.MagebloodPotion,
	spellPowerBuff: SpellPowerBuff.GreaterArcaneElixir,
	shadowPowerBuff: ShadowPowerBuff.ElixirOfShadowPower,
	zanzaBuff: ZanzaBuff.CerebralCortexCompound,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	divineSpirit: true,
	fireResistanceAura: true,
	fireResistanceTotem: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	manaSpringTotem: TristateEffect.TristateEffectRegular,
	moonkinAura: true,
	powerWordFortitude: TristateEffect.TristateEffectImproved,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
	// saygesFortune: SaygesFortune.SaygesDamage,
	// spiritOfZandalar: true,
});

export const DefaultDebuffs = Debuffs.create({
	exposeArmor: TristateEffect.TristateEffectImproved,
	faerieFire: true,
	judgementOfWisdom: true,
	sunderArmor: true,
});

///////////////////////////////////////////////////////////////////////////
//                                 Builds
///////////////////////////////////////////////////////////////////////////

// The community builds with the pet setup and rotation each one is measured with.
export const BuildDemonicPact = PresetUtils.makePresetBuild('Demonic Pact 2/31/18', {
	talents: TalentsPactOptimised,
	rotation: RotationDemonicPact,
	options: DemonicPactOptions,
});
export const BuildDeepAffliction = PresetUtils.makePresetBuild('Deep Affliction 35/0/16', {
	talents: TalentsDeepAffliction,
	rotation: RotationAffliction,
	options: AfflictionOptions,
});
export const BuildDSRuinPandemic = PresetUtils.makePresetBuild('DS/Ruin Pandemic 24/11/16', {
	talents: TalentsDSRuinPandemic,
	rotation: RotationDSRuin,
	options: DefaultOptions,
});
export const BuildShadowAndFlame = PresetUtils.makePresetBuild('Shadow and Flame 13/11/27', {
	talents: TalentsShadowAndFlame,
	rotation: RotationShadowAndFlame,
	options: DefaultOptions,
});

export const OtherDefaults = {
	distanceFromTarget: 25,
	profession1: Profession.Enchanting,
	profession2: Profession.Tailoring,
	channelClipDelay: 150,
};
