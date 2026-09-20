import { Class, Spec } from '@generated/proto/common';

import { PlayerSpec } from '../player/player_spec';
import { specTypeFunctions } from '../proto/spec_functions';
import type { SpecTalents } from '../proto/spec_types';
import { TalentsConfig } from './config';
import { druidTalentsConfig } from './druid';
import { hunterTalentsConfig } from './hunter';
import { mageTalentsConfig } from './mage';
import { paladinTalentsConfig } from './paladin';
import { priestTalentsConfig } from './priest';
import { rogueTalentsConfig } from './rogue';
import { shamanTalentsConfig } from './shaman';
import { warlockTalentsConfig } from './warlock';
import { warriorTalentsConfig } from './warrior';

export const classTalentsConfig: Record<Class, TalentsConfig<any>> = {
	[Class.ClassUnknown]: [],
	[Class.ClassExtra1]: [],
	[Class.ClassExtra2]: [],
	[Class.ClassExtra3]: [],
	[Class.ClassExtra4]: [],
	[Class.ClassExtra5]: [],
	[Class.ClassExtra6]: [],
	[Class.ClassDruid]: druidTalentsConfig,
	[Class.ClassShaman]: shamanTalentsConfig,
	[Class.ClassHunter]: hunterTalentsConfig,
	[Class.ClassMage]: mageTalentsConfig,
	[Class.ClassRogue]: rogueTalentsConfig,
	[Class.ClassPaladin]: paladinTalentsConfig,
	[Class.ClassPriest]: priestTalentsConfig,
	[Class.ClassWarlock]: warlockTalentsConfig,
	[Class.ClassWarrior]: warriorTalentsConfig,
} as const;

export function talentSpellIdsToTalentString(playerClass: Class, talentIds: Array<number>): string {
	const talentsConfig = classTalentsConfig[playerClass];

	const talentsStr = talentsConfig
		.map(treeConfig => {
			const treeStr = treeConfig.talents
				.map(talentConfig => {
					const spellIdIndex = talentConfig.spellIds.findIndex(spellId => talentIds.includes(spellId));
					if (spellIdIndex == -1) {
						return '0';
					} else {
						return String(spellIdIndex + 1);
					}
				})
				.join('')
				.replace(/0+$/g, '');

			return treeStr;
		})
		.join('-')
		.replace(/-+$/g, '');

	return talentsStr;
}

export function playerTalentStringToProto<SpecType extends Spec>(playerSpec: PlayerSpec<SpecType>, talentString: string): SpecTalents<SpecType> {
	const specFunctions = specTypeFunctions[playerSpec.specID];
	const proto = specFunctions.talentsCreate() as SpecTalents<SpecType>;
	const talentsConfig = classTalentsConfig[playerSpec.classID] as TalentsConfig<SpecTalents<SpecType>>;

	return talentStringToProto(proto, talentString, talentsConfig);
}

export function talentStringToProto<TalentsProto>(proto: TalentsProto, talentString: string, talentsConfig: TalentsConfig<TalentsProto>): TalentsProto {
	const treeStrings = talentString.split('-');
	talentsConfig.forEach((treeConfig, treeIdx) => {
		const treeString = treeStrings[treeIdx] ?? '';
		treeConfig.talents.forEach((talentConfig, i) => {
			const points = parseInt(treeString.charAt(i));
			if (!isNaN(points) && talentConfig.fieldName) {
				if (talentConfig.maxPoints == 1) {
					(proto[talentConfig.fieldName as keyof TalentsProto] as unknown as boolean) = points == 1;
				} else {
					(proto[talentConfig.fieldName as keyof TalentsProto] as unknown as number) = points;
				}
			}
		});
	});

	return proto;
}

// Note that this function will fail if any of the talent names are not defined. TODO: Remove that condition
// once all talents are migrated to wrath and use all fields.
export function protoToTalentString<TalentsProto>(proto: TalentsProto, talentsConfig: TalentsConfig<TalentsProto>): string {
	return talentsConfig
		.map(treeConfig => {
			return treeConfig.talents
				.map(talentConfig => String(Number(proto[(talentConfig.fieldName as keyof TalentsProto)!])))
				.join('')
				.replace(/0+$/g, '');
		})
		.join('-')
		.replace(/-+$/g, '');
}
