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
// gear
import BlankGear from './gear_sets/blank.gear.json';
import MCGear from './gear_sets/mc.gear.json';
import PreBisGear from './gear_sets/prebis.gear.json';

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

export const GearBlank = PresetUtils.makePresetGear('Blank', BlankGear);
export const GearPreBis = PresetUtils.makePresetGear('Pre-BIS', PreBisGear);
export const GearMC = PresetUtils.makePresetGear('MC', MCGear);

export const GearPresets = [
	GearBlank,
	GearPreBis,
	GearMC,
];

export const DefaultGear = GearPreBis;

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

// P1
export const RotationDemonicPact = PresetUtils.makePresetAPLRotation('Demonic Pact', DemonicPactApl);
export const RotationAffliction = PresetUtils.makePresetAPLRotation('Affliction', AfflictionApl);
export const RotationDSRuin = PresetUtils.makePresetAPLRotation('DS/Ruin', DSRuinApl);

export const APLPresets = [RotationDemonicPact, RotationAffliction, RotationDSRuin];

export const DefaultAPL = RotationDSRuin;

///////////////////////////////////////////////////////////////////////////
//                                 Talent Presets
///////////////////////////////////////////////////////////////////////////

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.

export const TalentsDemonicPact = {
	name: 'Demonic Pact',
	data: SavedTalents.create({ talentsString: '2-0055003231101001351-0550005003' }),
};

export const TalentsAffliction = {
	name: 'Affliction',
	data: SavedTalents.create({ talentsString: '2325002013500135--0550105003' }),
};

export const TalentsDSRuin = {
	name: 'DS/Ruin',
	data: SavedTalents.create({ talentsString: '23250020133-0320003201-0550105003' }),
};

export const TalentPresets = [TalentsDemonicPact, TalentsAffliction, TalentsDSRuin];

export const DefaultTalents = TalentsDSRuin;

///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = WarlockOptions.create({
	armor: Armor.DemonArmor,
	summon: Summon.Imp,
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
	fengusFerocity: true,
	moldarsMoxie: true,
	rallyingCryOfTheDragonslayer: true,
	// saygesFortune: SaygesFortune.SaygesDamage,
	slipkiksSavvy: true,
	songflowerSerenade: true,
	// spiritOfZandalar: true,
	warchiefsBlessing: true,
});

export const DefaultDebuffs = Debuffs.create({
	exposeArmor: TristateEffect.TristateEffectImproved,
	faerieFire: true,
	improvedScorch: true,
	judgementOfWisdom: true,
	shadowWeaving: true,
	sunderArmor: true,
});

export const OtherDefaults = {
	distanceFromTarget: 25,
	profession1: Profession.Enchanting,
	profession2: Profession.Tailoring,
	channelClipDelay: 150,
};
