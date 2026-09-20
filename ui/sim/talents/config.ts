// Shape of a class's talent trees, plus the builder that validates a generated
// tree JSON. Pure data, kept out of the talents picker so ui/sim's import path
// stays free of a view module; the picker re-imports these from here.

export type TalentsConfig<TalentsProto> = Array<TalentTreeConfig<TalentsProto>>;

export type TalentTreeConfig<TalentsProto> = {
	name: string;
	backgroundUrl: string;
	talents: Array<TalentConfig<TalentsProto>>;
};

export type TalentLocation = {
	// 0-indexed row in the tree
	rowIdx: number;
	// 0-indexed column in the tree
	colIdx: number;
};

export type TalentConfig<TalentsProto> = {
	fieldName?: keyof TalentsProto | string;

	// Display name, as it appears in game.
	fancyName: string;

	location: TalentLocation;

	// Location of a prerequisite talent, if any
	prereqLocation?: TalentLocation;

	// Child talents depending on this talent. This is populated automatically.
	childLocations?: TalentLocation[];

	// Spell ID for each rank of this talent.
	// Omitted ranks will be inferred by incrementing from the last provided rank.
	spellIds: Array<number>;

	definitionId?: number;

	maxPoints: number;
};

export function newTalentsConfig<TalentsProto>(talents: TalentsConfig<TalentsProto>): TalentsConfig<TalentsProto> {
	talents.forEach(tree => {
		tree.talents.forEach((talent, i) => {
			talent.childLocations = [];
			// Validate that talents are given in the correct order (left-to-right top-to-bottom).
			if (i != 0) {
				const prevTalent = tree.talents[i - 1];
				if (
					talent.location.rowIdx < prevTalent.location.rowIdx ||
					(talent.location.rowIdx == prevTalent.location.rowIdx && talent.location.colIdx <= prevTalent.location.colIdx)
				) {
					throw new Error(`Out-of-order talent: ${String(talent.fieldName)}`);
				}
			}

			// Infer omitted spell IDs.
			if (talent.spellIds.length < talent.maxPoints) {
				let curSpellId = talent.spellIds[talent.spellIds.length - 1];
				for (let pointIdx = talent.spellIds.length; pointIdx < talent.maxPoints; pointIdx++) {
					curSpellId++;
					talent.spellIds.push(curSpellId);
				}
			}
		});
	});
	return talents;
}
