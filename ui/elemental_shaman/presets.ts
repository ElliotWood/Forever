import { ClassicPhase } from '../core/constants/other.js';
import * as PresetUtils from '../core/preset_utils.js';
import {
	AttackPowerBuff,
	Conjured,
	Consumes,
	Debuffs,
	FirePowerBuff,
	Flask,
	Food,
	FrostPowerBuff,
	IndividualBuffs,
	ManaRegenElixir,
	Potions,
	Profession,
	RaidBuffs,
	SaygesFortune,
	SpellPowerBuff,
	StrengthBuff,
	TristateEffect,
	WeaponImbue,
	ZanzaBuff,
} from '../core/proto/common.js';
import { ElementalShaman_Options as ElementalShamanOptions } from '../core/proto/shaman.js';
import { SavedTalents } from '../core/proto/ui.js';
import DefaultAPLJson from './apls/default.apl.json';
import Phase1GearJSON from './gear_sets/phase_1.gear.json';
import Phase2GearJSON from './gear_sets/phase_2.gear.json';
import Phase3GearJSON from './gear_sets/phase_3.gear.json';
import Phase4GearJSON from './gear_sets/phase_4.gear.json';
import Phase5GearJSON from './gear_sets/phase_5.gear.json';
import Phase6GearJSON from './gear_sets/phase_6.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

export const GearPhase1 = PresetUtils.makePresetGear('Phase 1', Phase1GearJSON);
export const GearPhase2 = PresetUtils.makePresetGear('Phase 2', Phase2GearJSON);
export const GearPhase3 = PresetUtils.makePresetGear('Phase 3', Phase3GearJSON);
export const GearPhase4 = PresetUtils.makePresetGear('Phase 4', Phase4GearJSON);
export const GearPhase5 = PresetUtils.makePresetGear('Phase 5', Phase5GearJSON);
export const GearPhase6 = PresetUtils.makePresetGear('Phase 6', Phase6GearJSON);

export const GearPresets = {
	[ClassicPhase.Phase1]: [GearPhase1],
	[ClassicPhase.Phase2]: [GearPhase2],
	[ClassicPhase.Phase3]: [GearPhase3],
	[ClassicPhase.Phase4]: [GearPhase4],
	[ClassicPhase.Phase5]: [GearPhase5],
	[ClassicPhase.Phase6]: [GearPhase6],

};

export const DefaultGear = GearPresets[ClassicPhase.Phase2][0];

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

export const APLDefault = PresetUtils.makePresetAPLRotation('Default', DefaultAPLJson);

export const APLPresets = {
	[ClassicPhase.Phase1]: [APLDefault],
	[ClassicPhase.Phase2]: [],
	[ClassicPhase.Phase3]: [],
	[ClassicPhase.Phase4]: [],
	[ClassicPhase.Phase5]: [],
	[ClassicPhase.Phase6]: [],

};

export const DefaultAPL = APLPresets[ClassicPhase.Phase1][0];

///////////////////////////////////////////////////////////////////////////
//                                 Talent Presets
///////////////////////////////////////////////////////////////////////////

export const TalentsLevel60 = PresetUtils.makePresetTalents('Level 60', SavedTalents.create({ talentsString: '5505301500103031--503352001' }));

export const TalentPresets = {
	[ClassicPhase.Phase1]: [TalentsLevel60],
	[ClassicPhase.Phase2]: [],
	[ClassicPhase.Phase3]: [],
	[ClassicPhase.Phase4]: [],
	[ClassicPhase.Phase5]: [],
	[ClassicPhase.Phase6]: [],
};

export const DefaultTalents = TalentPresets[ClassicPhase.Phase1][0];

///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = ElementalShamanOptions.create({});

export const DefaultConsumes = Consumes.create({
	attackPowerBuff: AttackPowerBuff.JujuMight,
	defaultConjured: Conjured.ConjuredDemonicRune,
	defaultPotion: Potions.MajorManaPotion,
	firePowerBuff: FirePowerBuff.ElixirOfFirepower,
	frostPowerBuff: FrostPowerBuff.ElixirOfFrostPower,
	flask: Flask.FlaskOfSupremePower,
	food: Food.FoodRunnTumTuberSurprise,
	// Not available until Phase 4
	// mainHandImbue: WeaponImbue.BrilliantWizardOil,
	manaRegenElixir: ManaRegenElixir.MagebloodPotion,
	spellPowerBuff: SpellPowerBuff.GreaterArcaneElixir,
	strengthBuff: StrengthBuff.JujuPower,
	zanzaBuff: ZanzaBuff.CerebralCortexCompound,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	divineSpirit: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	manaSpringTotem: TristateEffect.TristateEffectRegular,
	moonkinAura: true,
	powerWordFortitude: TristateEffect.TristateEffectImproved,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
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
	curseOfElements: true,
	stormstrike: true,
});

export const OtherDefaults = {
	distanceFromTarget: 15,
	profession2: Profession.Alchemy,
	profession1: Profession.Enchanting,
};
