// The individual sims register their configs as a side effect of being imported, and the
// raid presets module is the one place that already imports every one of them. Without
// this, getSpecConfig() below has nothing to look up.
import '../raid/presets';

import { IndividualSimUIConfig } from '../core/individual_sim_ui';
import { getSpecConfig, Player } from '../core/player';
import { StatWeightsResult } from '../core/proto/api';
import { Profession, Spec, Stat } from '../core/proto/common';
import { isHealingSpec, isTankSpec, specToEligibleRaces } from '../core/proto_utils/utils';
import { Sim } from '../core/sim';
import { TypedEvent } from '../core/typed_event';
import { WorkerProgressCallback } from '../core/worker_pool';

// Everything this page needs to know about a spec before it has been simmed: which stats
// the spec's own sim asks to be weighed, and what those weights get normalised against.
export interface SpecEpConfig {
	epStats: Array<Stat>;
	epReferenceStat: Stat;
}

export function getSpecEpConfig(spec: Spec): SpecEpConfig {
	const config = getSpecConfig(spec) as IndividualSimUIConfig<any>;
	return {
		epStats: config.epStats,
		epReferenceStat: config.epReferenceStat,
	};
}

/**
 * Sets the sim up as if the player had just opened that spec's individual sim and touched
 * nothing, then runs the same stat weights request its 'Stat Weights' button issues.
 *
 * This mirrors IndividualSimUI.applyDefaults(), which cannot be called directly because it
 * belongs to a full sim UI and building fifteen of those off-screen would be absurd. The
 * player is installed in raid slot 0 so the request has a party and a raid to read buffs
 * from, exactly as the individual sims do.
 */
export async function runSpecStatWeights(sim: Sim, spec: Spec, iterations: number, onProgress: WorkerProgressCallback): Promise<StatWeightsResult> {
	await sim.waitForInit();

	const config = getSpecConfig(spec) as IndividualSimUIConfig<any>;
	const player = new Player<any>(spec, sim);
	const tankSpec = isTankSpec(spec);
	const healingSpec = isHealingSpec(spec);

	const eventID = TypedEvent.nextEventID();
	TypedEvent.freezeAllAndDo(() => {
		sim.raid.setPlayer(eventID, 0, player);

		player.applySharedDefaults(eventID);
		player.setName(eventID, 'Player');
		player.setRace(eventID, config.defaults.race ?? specToEligibleRaces[spec][0]);
		player.setGear(eventID, sim.db.lookupEquipmentSpec(config.defaults.gear));
		player.setConsumes(eventID, config.defaults.consumes);
		player.setTalentsString(eventID, config.defaults.talents.talentsString);
		player.setSpecOptions(eventID, config.defaults.specOptions);
		player.setBuffs(eventID, config.defaults.individualBuffs);
		player.getParty()!.setBuffs(eventID, config.defaults.partyBuffs);
		player.getRaid()!.setBuffs(eventID, config.defaults.raidBuffs);
		player.setProfession1(eventID, config.defaults.other?.profession1 || Profession.Engineering);
		player.setProfession2(eventID, config.defaults.other?.profession2 || Profession.ProfessionUnknown);
		player.setDistanceFromTarget(eventID, config.defaults.other?.distanceFromTarget || 0);
		player.setChannelClipDelay(eventID, config.defaults.other?.channelClipDelay || 0);

		sim.raid.setTargetDummies(eventID, healingSpec ? 9 : 0);
		sim.encounter.applyDefaults(eventID);
		sim.raid.setDebuffs(eventID, config.defaults.debuffs);
		sim.applyDefaults(eventID, tankSpec, healingSpec);
		sim.raid.setTanks(eventID, tankSpec ? [player.makeUnitReference()] : []);

		// sim.applyDefaults() resets the iteration count, so this has to come after it.
		sim.setIterations(eventID, iterations);
	});

	// The rotation is left on 'auto', which is what the individual sims default to, and is
	// resolved into a real APL by Player.toProto() when the request is built.
	//
	// Only the spec's own epStats are weighed. The individual sims also append the three
	// resistances from GLOBAL_EP_STATS and weigh pseudo-stats such as main hand DPS, which
	// between them would cost up to ten more sims per spec for columns this page does not
	// show: none of them is a stat you compare across specs.
	return sim.statWeights(player, config.epStats, [], config.epReferenceStat, onProgress);
}
