import { BalanceDruidSimUI } from '../balance_druid/sim.js';
import { IndividualSimUI, IndividualSimUIConfig, RaidSimPreset } from '../core/individual_sim_ui.js';
import { getSpecConfig, Player } from '../core/player.js';
import { Raid as RaidProto } from '../core/proto/api.js';
import { Class, EquipmentSpec, Profession, Spec, TristateEffect } from '../core/proto/common.js';
import { BalanceDruid_Options as BalanceDruidOptions } from '../core/proto/druid.js';
import { Blessings } from '../core/proto/paladin.js';
import { BlessingsAssignments } from '../core/proto/ui.js';
import { isTankSpec, naturalSpecOrder, newUnitReference, playerToSpec } from '../core/proto_utils/utils.js';
import { Raid } from '../core/raid.js';
import { Sim } from '../core/sim.js';
import { EventID } from '../core/typed_event.js';
import { ElementalShamanSimUI } from '../elemental_shaman/sim.js';
import { EnhancementShamanSimUI } from '../enhancement_shaman/sim.js';
import { FeralDruidSimUI } from '../feral_druid/sim.js';
import { FeralTankDruidSimUI } from '../feral_tank_druid/sim.js';
import { HealingPriestSimUI } from '../healing_priest/sim.js';
import { HolyPaladinSimUI } from '../holy_paladin/sim.js';
import { HunterSimUI } from '../hunter/sim.js';
import { MageSimUI } from '../mage/sim.js';
import { ProtectionPaladinSimUI } from '../protection_paladin/sim.js';
import { RestorationDruidSimUI } from '../restoration_druid/sim.js';
import { RestorationShamanSimUI } from '../restoration_shaman/sim.js';
import { RetributionPaladinSimUI } from '../retribution_paladin/sim.js';
import { RogueSimUI } from '../rogue/sim.js';
import { ShadowPriestSimUI } from '../shadow_priest/sim.js';
import { SmitePriestSimUI } from '../smite_priest/sim.js';
import { TankWarriorSimUI } from '../tank_warrior/sim.js';
import { WardenShamanSimUI } from '../warden_shaman/sim.js';
import { WarlockSimUI } from '../warlock/sim.js';
import { WarriorSimUI } from '../warrior/sim.js';

export const specSimFactories: Record<Spec, (parentElem: HTMLElement, player: Player<any>) => IndividualSimUI<any>> = {
	[Spec.SpecBalanceDruid]: (parentElem: HTMLElement, player: Player<any>) => new BalanceDruidSimUI(parentElem, player),
	[Spec.SpecFeralDruid]: (parentElem: HTMLElement, player: Player<any>) => new FeralDruidSimUI(parentElem, player),
	[Spec.SpecFeralTankDruid]: (parentElem: HTMLElement, player: Player<any>) => new FeralTankDruidSimUI(parentElem, player),
	[Spec.SpecRestorationDruid]: (parentElem: HTMLElement, player: Player<any>) => new RestorationDruidSimUI(parentElem, player),
	[Spec.SpecElementalShaman]: (parentElem: HTMLElement, player: Player<any>) => new ElementalShamanSimUI(parentElem, player),
	[Spec.SpecEnhancementShaman]: (parentElem: HTMLElement, player: Player<any>) => new EnhancementShamanSimUI(parentElem, player),
	[Spec.SpecRestorationShaman]: (parentElem: HTMLElement, player: Player<any>) => new RestorationShamanSimUI(parentElem, player),
	[Spec.SpecWardenShaman]: (parentElem: HTMLElement, player: Player<any>) => new WardenShamanSimUI(parentElem, player),
	[Spec.SpecHunter]: (parentElem: HTMLElement, player: Player<any>) => new HunterSimUI(parentElem, player),
	[Spec.SpecMage]: (parentElem: HTMLElement, player: Player<any>) => new MageSimUI(parentElem, player),
	[Spec.SpecRogue]: (parentElem: HTMLElement, player: Player<any>) => new RogueSimUI(parentElem, player),
	[Spec.SpecHolyPaladin]: (parentElem: HTMLElement, player: Player<any>) => new HolyPaladinSimUI(parentElem, player),
	[Spec.SpecProtectionPaladin]: (parentElem: HTMLElement, player: Player<any>) => new ProtectionPaladinSimUI(parentElem, player),
	[Spec.SpecRetributionPaladin]: (parentElem: HTMLElement, player: Player<any>) => new RetributionPaladinSimUI(parentElem, player),
	[Spec.SpecHealingPriest]: (parentElem: HTMLElement, player: Player<any>) => new HealingPriestSimUI(parentElem, player),
	[Spec.SpecShadowPriest]: (parentElem: HTMLElement, player: Player<any>) => new ShadowPriestSimUI(parentElem, player),
	[Spec.SpecSmitePriest]: (parentElem: HTMLElement, player: Player<any>) => new SmitePriestSimUI(parentElem, player),
	[Spec.SpecWarrior]: (parentElem: HTMLElement, player: Player<any>) => new WarriorSimUI(parentElem, player),
	[Spec.SpecTankWarrior]: (parentElem: HTMLElement, player: Player<any>) => new TankWarriorSimUI(parentElem, player),
	[Spec.SpecWarlock]: (parentElem: HTMLElement, player: Player<any>) => new WarlockSimUI(parentElem, player),
};

export const playerPresets: Array<RaidSimPreset<any>> = naturalSpecOrder
	.map(getSpecConfig)
	.map(config => {
		const indSimUiConfig = config as IndividualSimUIConfig<any>;
		return indSimUiConfig.raidSimPresets;
	})
	.flat();

export const implementedSpecs: Array<Spec> = [...new Set(playerPresets.map(preset => preset.spec))];

// Every preset's gear is keyed by the phase it was authored for, and most only carry
// one. Asking for a phase a preset has no entry for used to hand lookupEquipmentSpec an
// undefined spec, which throws inside the promise and leaves the player naked - the raid
// sim then runs a raid of unequipped characters. Take the best phase at or below the one
// asked for instead, and failing that the earliest the preset has.
function gearForPhase(byPhase: Record<number, EquipmentSpec>, phase: number): EquipmentSpec | null {
	const phases = Object.keys(byPhase)
		.map(k => parseInt(k))
		.sort((a, b) => a - b);
	if (!phases.length) return null;
	const atOrBelow = phases.filter(p => p <= phase);
	return byPhase[atOrBelow.length ? atOrBelow[atOrBelow.length - 1] : phases[0]];
}

// Builds the player the raid picker drops into a slot. Anything that fills a raid from
// presets goes through here, so every page gets the same build for a given spec.
export function newPlayerFromPreset(eventID: EventID, sim: Sim, preset: RaidSimPreset<any>): Player<any> {
	const newPlayer = new Player(preset.spec, sim);
	newPlayer.applySharedDefaults(eventID);
	newPlayer.setRace(eventID, preset.defaultFactionRaces[sim.getFaction()]);
	newPlayer.setTalentsString(eventID, preset.talents.talentsString);
	newPlayer.setSpecOptions(eventID, preset.specOptions);
	newPlayer.setConsumes(eventID, preset.consumes);
	newPlayer.setName(eventID, preset.defaultName);
	newPlayer.setProfession1(eventID, preset.otherDefaults?.profession1 || Profession.Engineering);
	newPlayer.setProfession2(eventID, preset.otherDefaults?.profession2 || Profession.Enchanting);
	newPlayer.setDistanceFromTarget(eventID, preset.otherDefaults?.distanceFromTarget || 0);

	// Need to wait because the gear might not be loaded yet.
	sim.waitForInit().then(() => {
		const gear = gearForPhase(preset.defaultGear[sim.getFaction()], sim.getPhase());
		if (gear) {
			newPlayer.setGear(eventID, sim.db.lookupEquipmentSpec(gear));
		}
	});

	return newPlayer;
}

// Assignments that only make sense once the player is in the raid, so this has to run
// after the player has been placed in a party.
export function applyNewPlayerAssignments(eventID: EventID, newPlayer: Player<any>, raid: Raid) {
	if (isTankSpec(newPlayer.spec)) {
		const tanks = raid.getTanks();
		const emptyIdx = tanks.findIndex(tank => raid.getPlayerFromUnitReference(tank) == null);
		if (emptyIdx == -1) {
			if (tanks.length < 3) {
				raid.setTanks(eventID, tanks.concat([newPlayer.makeUnitReference()]));
			}
		} else {
			tanks[emptyIdx] = newPlayer.makeUnitReference();
			raid.setTanks(eventID, tanks);
		}
	}

	// Spec-specific assignments. For most cases, default to buffing self.
	if (newPlayer.spec == Spec.SpecBalanceDruid) {
		const newOptions = newPlayer.getSpecOptions() as BalanceDruidOptions;
		newOptions.innervateTarget = newUnitReference(newPlayer.getRaidIndex());
		newPlayer.setSpecOptions(eventID, newOptions);
	}
}

// Stamps the paladins' blessing assignments onto the player protos they were assigned
// to. Blessings aren't a player setting, so they can only be applied once the whole raid
// has been serialized - every page that runs a raid does this from its modifyRaidProto.
export function applyBlessings(raidProto: RaidProto, assignments: BlessingsAssignments, numPaladins: number) {
	implementedSpecs.forEach(spec => {
		const playerProtos = raidProto.parties
			.map(party => party.players.filter(player => player.class != Class.ClassUnknown && playerToSpec(player) == spec))
			.flat();

		assignments.paladins.forEach((paladin, i) => {
			if (i >= numPaladins) {
				return;
			}

			if (paladin.blessings[spec] == Blessings.BlessingOfKings) {
				playerProtos.forEach(playerProto => (playerProto.buffs!.blessingOfKings = true));
			} else if (paladin.blessings[spec] == Blessings.BlessingOfMight) {
				playerProtos.forEach(playerProto => (playerProto.buffs!.blessingOfMight = TristateEffect.TristateEffectImproved));
			} else if (paladin.blessings[spec] == Blessings.BlessingOfWisdom) {
				playerProtos.forEach(playerProto => (playerProto.buffs!.blessingOfWisdom = TristateEffect.TristateEffectImproved));
			} else if (paladin.blessings[spec] == Blessings.BlessingOfSanctuary) {
				playerProtos.forEach(playerProto => (playerProto.buffs!.blessingOfSanctuary = true));
			}
		});
	});
}
