// Whether a proposed point change would leave a legal tree.
//
// A transcription of ui/core/talents/talents_picker.tsx:476-521, which asked the DOM for every
// talent's current points. Nothing here reads state: the caller passes the parsed `TalentPoints`,
// so the same rules serve the click handler and the "can this still take a point" highlight.
//
// The two branches are deliberately asymmetric and must stay so: the add branch gates on the whole
// tree total against rowIdx * pointsPerRow, while the remove branch gates on per-row totals
// accumulated down to the row above. Collapsing them into one rule breaks tier gating silently.

import type { TalentsConfig } from '@sim/talents/config';

import type { TalentPoints } from './talents_string';
import { CHARACTER_LEVEL } from '@sim/constants/mechanics';

import { totalPointsSpent, treePointTotal } from './talents_string';
import type { TalentGraph } from './tree_graph';

// A character earns one talent point per level from 10 onwards, so the cap follows
// CHARACTER_LEVEL: 51 at level 60, where TBC's 61 was the level 70 figure.
export const MAX_POINTS_PLAYER = CHARACTER_LEVEL - 9;
export const POINTS_PER_ROW = 5;

export interface TalentLimits {
	maxPoints: number;
	pointsPerRow: number;
}

export const DEFAULT_TALENT_LIMITS: TalentLimits = { maxPoints: MAX_POINTS_PLAYER, pointsPerRow: POINTS_PER_ROW };

export const canSetPoints = <TalentsProto>(
	config: TalentsConfig<TalentsProto>,
	graph: TalentGraph,
	points: TalentPoints,
	treeIdx: number,
	talentIdx: number,
	newPoints: number,
	limits: TalentLimits = DEFAULT_TALENT_LIMITS,
): boolean => {
	const talents = config[treeIdx].talents;
	const treeGraph = graph.trees[treeIdx];
	const oldPoints = points[treeIdx][talentIdx];
	const rowIdx = talents[talentIdx].location.rowIdx;

	if (newPoints > oldPoints) {
		if (totalPointsSpent(points) + (newPoints - oldPoints) > limits.maxPoints) return false;
		if (treePointTotal(points, treeIdx) < rowIdx * limits.pointsPerRow) return false;

		const prereqIdx = treeGraph.prereqIdx[talentIdx];
		if (prereqIdx >= 0 && points[treeIdx][prereqIdx] < talents[prereqIdx].maxPoints) return false;
		return true;
	}

	const pointTotalsByRow = Array.from({ length: graph.numRows }, () => 0);
	talents.forEach((talent, idx) => (pointTotalsByRow[talent.location.rowIdx] += points[treeIdx][idx]));
	pointTotalsByRow[rowIdx] -= oldPoints - newPoints;

	let running = 0;
	const cumulativeTotalsByRow = pointTotalsByRow.map(total => (running += total));

	const strands = talents.some((talent, idx) => {
		const row = talent.location.rowIdx;
		return points[treeIdx][idx] > 0 && row > 0 && cumulativeTotalsByRow[row - 1] < row * limits.pointsPerRow;
	});
	if (strands) return false;

	return !treeGraph.childIdxs[talentIdx].some(childIdx => points[treeIdx][childIdx] > 0);
};
