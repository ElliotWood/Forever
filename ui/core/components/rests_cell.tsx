import clsx from 'clsx';

import { Composition, TIER_LABELS, TIER_ORDER, unsettledShare } from '../spells/rests';
import { formatToPercent } from '../utils';

/**
 * A stacked bar of where a build's damage came from, with the one number worth reading next
 * to it. Shared by the live rankings page and the precomputed arena so the two cannot start
 * saying the same thing differently.
 *
 * Titled rather than tooltipped: these are tables of hundreds of rows, and a tippy instance
 * per cell is hundreds of them built for a hover most people never make.
 */
export const restsCell = (rests: Composition): Element => {
	const unsettled = unsettledShare(rests);
	const breakdown = TIER_ORDER.filter(tier => rests[tier] > 0)
		.map(tier => `${formatToPercent(rests[tier] * 100, { maximumFractionDigits: 1 })} ${TIER_LABELS[tier]}`)
		.join('\n');

	return (
		<div className="rests" attributes={{ title: `Of this build's damage:\n${breakdown}` }}>
			<div className="rests-bar">
				{TIER_ORDER.map(tier => (
					<div className={`rests-part rests-${tier}`} style={{ '--percentage': formatToPercent(rests[tier] * 100) }} />
				))}
			</div>
			<span className={clsx('rests-percent', unsettled >= 0.05 && 'rests-high')}>{formatToPercent(unsettled * 100, { maximumFractionDigits: 1 })}</span>
		</div>
	) as Element;
};
