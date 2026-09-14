import { IndividualSimUIConfig } from '../core/individual_sim_ui';
import { getLaunchedSims } from '../core/launched_sims';
import { getSpecConfig } from '../core/player';
import { Spec } from '../core/proto/common';
import { naturalSpecOrder } from '../core/proto_utils/utils';

// Where each spec's sim UI lives. Importing that module is what puts the spec's config,
// and with it the EP weights this page ranks by, into the player registry.
const specModules: Partial<Record<Spec, () => Promise<unknown>>> = {
	[Spec.SpecBalanceDruid]: () => import('../balance_druid/sim'),
	[Spec.SpecFeralDruid]: () => import('../feral_druid/sim'),
	[Spec.SpecElementalShaman]: () => import('../elemental_shaman/sim'),
	[Spec.SpecEnhancementShaman]: () => import('../enhancement_shaman/sim'),
	[Spec.SpecWardenShaman]: () => import('../warden_shaman/sim'),
	[Spec.SpecHunter]: () => import('../hunter/sim'),
	[Spec.SpecMage]: () => import('../mage/sim'),
	[Spec.SpecRogue]: () => import('../rogue/sim'),
	[Spec.SpecProtectionPaladin]: () => import('../protection_paladin/sim'),
	[Spec.SpecRetributionPaladin]: () => import('../retribution_paladin/sim'),
	[Spec.SpecShadowPriest]: () => import('../shadow_priest/sim'),
	[Spec.SpecSmitePriest]: () => import('../smite_priest/sim'),
	[Spec.SpecWarlock]: () => import('../warlock/sim'),
	[Spec.SpecWarrior]: () => import('../warrior/sim'),
	[Spec.SpecTankWarrior]: () => import('../tank_warrior/sim'),
};

// The launched specs this page can rank, in the order the sim menu lists them. Anything
// launched without a module here would be a spec whose sim exists but whose EP weights
// this page cannot reach, so it is left out rather than guessed at.
export const bisSpecs: Array<Spec> = naturalSpecOrder.filter(spec => getLaunchedSims().includes(spec) && !!specModules[spec]);

const configs: Partial<Record<Spec, IndividualSimUIConfig<any>>> = {};

// Loads a spec's sim config, which carries the EP weights its gear picker sorts by.
//
// registerSpecConfig() hands the whole IndividualSimUIConfig to the player registry, which
// declares only the rotation part of it. Nothing else needs the rest, so nothing else casts
// it back; this page does, because the weights are the one number it is built around.
export async function loadSpecConfig(spec: Spec): Promise<IndividualSimUIConfig<any>> {
	if (!configs[spec]) {
		await specModules[spec]!();
		configs[spec] = getSpecConfig(spec) as unknown as IndividualSimUIConfig<any>;
	}
	return configs[spec]!;
}
