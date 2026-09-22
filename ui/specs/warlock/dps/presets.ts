import * as PresetUtils from '@app/preset_utils';
import { ConsumesSpec, Debuffs, Drums, IndividualBuffs, PartyBuffs, Profession, RaidBuffs, TristateEffect } from '@generated/proto/common';
import { SavedTalents } from '@generated/proto/ui';
import { Warlock_Options as WarlockOptions, WarlockOptions_Armor, WarlockOptions_CurseOptions, WarlockOptions_Summon } from '@generated/proto/warlock';
import { defaultExposeWeaknessSettings, defaultImprovedShadowBoltSettings, defaultRaidBuffMajorDamageCooldowns } from '@sim/proto/utils';

import AfflictionRot from './apls/affliction.apl.json';
import BlankAPL from './apls/blank.apl.json';
import DemoRot from './apls/demonology.apl.json';
import DestroFireRot from './apls/destro_fire.apl.json';
import DestroRot from './apls/destruction.apl.json';
import LaunchGear from './gear_sets/launch.gear.json';
import McGear from './gear_sets/mc.gear.json';
import PrebisGear from './gear_sets/prebis.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const BLANK_APL = PresetUtils.makePresetAPLRotation('Blank', BlankAPL);

// Rotations
export const AfflictionAPL = PresetUtils.makePresetAPLRotation('Affliction', AfflictionRot);
export const DemoAPL = PresetUtils.makePresetAPLRotation('Demonology', DemoRot);
export const DestroAPL = PresetUtils.makePresetAPLRotation('Destruction', DestroRot);
export const DestroFireAPL = PresetUtils.makePresetAPLRotation('Destruction (Fire)', DestroFireRot);

// Defaults
export const DefaultOptions = WarlockOptions.create({
	classOptions: {
		armor: WarlockOptions_Armor.FelArmor,
		curseOptions: WarlockOptions_CurseOptions.Recklessness,
		sacrificeSummon: true,
		summon: WarlockOptions_Summon.Succubus,
	},
});

export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 22866, // Flask of Pure Death
	foodId: 27657, // Blackened Basilisk
	conjuredId: 12662, // Demonic Rune
	mhImbueId: 25122, // Brilliant Wizard Oil
	potId: 22839, // Destruction Potion
	explosiveId: 30217,
	petScrollAgi: true,
	petScrollStr: true,
});

export const OtherDefaults = {
	distanceFromTarget: 20,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
};

export const DefaultRaidBuffs = RaidBuffs.create({
	...defaultRaidBuffMajorDamageCooldowns(),
	arcaneBrilliance: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	powerWordFortitude: TristateEffect.TristateEffectImproved,
	divineSpirit: TristateEffect.TristateEffectImproved,
});

export const DefaultPartyBuffs = PartyBuffs.create({
	manaSpringTotem: TristateEffect.TristateEffectRegular,
	moonkinAura: TristateEffect.TristateEffectRegular,
	totemOfWrath: 1,
	wrathOfAirTotem: TristateEffect.TristateEffectImproved,
	eyeOfTheNight: true,
	chainOfTheTwilightOwl: true,
	drums: Drums.LesserDrumsOfBattle,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
	shadowPriestDps: 0,
});

export const DefaultDebuffs = Debuffs.create({
	...defaultExposeWeaknessSettings(),
	...defaultImprovedShadowBoltSettings(),
	improvedSealOfTheCrusader: TristateEffect.TristateEffectImproved,
	judgementOfWisdom: true,
	misery: true,
	shadowWeaving: true,
	sunderArmor: true,
	screech: true,
	faerieFire: TristateEffect.TristateEffectImproved,
	curseOfRecklessness: true,
	shadowEmbrace: true,
	curseOfElements: TristateEffect.TristateEffectImproved,
	bloodFrenzy: true,
	giftOfArthas: true,
	mangle: true,
	exposeArmor: TristateEffect.TristateEffectImproved,
	huntersMark: TristateEffect.TristateEffectImproved,
});

// Talent presets, from master's ui/warlock spec.
export const TalentsDemonicPact = PresetUtils.makePresetTalents('Demonic Pact', SavedTalents.create({ talentsString: '203-0055003221201001351-0550005' }));
export const TalentsAffliction = PresetUtils.makePresetTalents('Affliction', SavedTalents.create({ talentsString: '2435002013520135--0500055' }));
export const TalentsDSRuin = PresetUtils.makePresetTalents('DS/Ruin', SavedTalents.create({ talentsString: '233500201332-0340003001-0550105' }));
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

// Our Forever sim's gear presets (master ui/<spec>/gear_sets).
export const GEAR_LAUNCH = PresetUtils.makePresetGear('Launch', LaunchGear);
export const GEAR_PREBIS = PresetUtils.makePresetGear('Pre-BiS', PrebisGear);
export const GEAR_MC = PresetUtils.makePresetGear('MC', McGear);
export const DEFAULT_GEAR = GEAR_PREBIS;
export const GEAR_PRESETS = [GEAR_LAUNCH, GEAR_PREBIS, GEAR_MC];
