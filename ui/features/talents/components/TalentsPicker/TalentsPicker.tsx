import i18n from '@i18n/config';
import { usePlayer } from '@sim/context/SimHostContext';
import type { Player } from '@sim/player/player';
import type { TalentsConfig } from '@sim/talents/config';
import { Button } from '@ui-kit/Button';
import { useCopyToClipboard } from '@ui-kit/hooks/useCopyToClipboard';
import { useInput } from '@ui-kit/hooks/useInput';
import { useWowheadTooltipRefresh } from '@ui-kit/hooks/useWowheadTooltipRefresh';
import { Icon } from '@ui-kit/Icon';
import type { InputConfig } from '@ui-kit/input';
import { PickerShell } from '@ui-kit/PickerShell';
import { Tooltip, tooltipAnchorProps } from '@ui-kit/Tooltip';
import type { CSSProperties } from 'react';
import { useCallback, useId, useMemo, useState } from 'react';

import { MAX_POINTS_PLAYER } from '../../model/can_set_points';
import type { TalentPoints } from '../../model/talents_string';
import { parseTalentsString, serializeTalentsString, totalPointsSpent } from '../../model/talents_string';
import { buildTalentGraph } from '../../model/tree_graph';
import { TalentTreePicker } from './TalentTreePicker';

export interface TalentsPickerConfig<ModObject, TalentsProto> extends InputConfig<ModObject, string> {
	trees: TalentsConfig<TalentsProto>;
	pointsPerRow: number;
}

export interface TalentsPickerProps<TalentsProto> {
	config: TalentsPickerConfig<Player<any>, TalentsProto>;
}

// TBC shows the middle tree first on mobile, and the carousel slides one tree either side of it.
const INITIAL_TREE_IDX = 1;

export const TalentsPicker = <TalentsProto,>({ config }: TalentsPickerProps<TalentsProto>) => {
	const player = usePlayer();
	const fallbackId = useId();
	const { value, setValue, hidden, disabled } = useInput(player, config);

	const copyTooltipId = useId();
	const { copy, copied } = useCopyToClipboard(() => player.getTalentsString());

	const [activeTreeIdx, setActiveTreeIdx] = useState(INITIAL_TREE_IDX);

	const graph = useMemo(() => buildTalentGraph(config.trees), [config.trees]);
	const points = useMemo(() => parseTalentsString(config.trees, value ?? ''), [config.trees, value]);
	const limits = useMemo(() => ({ maxPoints: MAX_POINTS_PLAYER, pointsPerRow: config.pointsPerRow }), [config.pointsPerRow]);

	const onChange = useCallback((next: TalentPoints) => setValue(serializeTalentsString(next)), [setValue]);

	useWowheadTooltipRefresh(value);

	return (
		<PickerShell
			config={{ ...config, id: config.id ?? fallbackId }}
			className="col-span-full flex w-fit flex-row gap-section max-fhd:flex-col max-xl:m-auto max-md:w-full"
			hidden={hidden}
			disabled={disabled}>
			<div className="flex flex-col max-lg:w-full">
				<div className="mb-1 flex w-full items-center">
					<label>
						{i18n.t('talents_tab.points_remaining', { defaultValue: 'Points Remaining' })}:{' '}
						<span data-testid="talents-picker-points-remaining">{limits.maxPoints - totalPointsSpent(points)}</span>
					</label>
					<div className="ml-auto" data-testid="talents-picker-actions">
						<Button
							variant="outline-primary"
							size="sm"
							className="w-24"
							data-testid="copy-button"
							onClick={copy}
							{...tooltipAnchorProps(copyTooltipId)}>
							<Icon name={copied ? 'check' : 'copy'} className="mr-1" />
							{copied ? i18n.t('common.copy_button.copied') : i18n.t('talents_tab.copy_button.label')}
						</Button>
						<Tooltip id={copyTooltipId} content={i18n.t('talents_tab.copy_button.tooltip')} />
					</div>
				</div>
				<div className="relative flex-1 max-lg:-mx-page max-lg:flex max-lg:justify-center max-lg:overflow-x-hidden">
					<div
						className="ui-talents-picker-trees"
						data-testid="talents-picker-trees"
						style={{ '--talents-carousel-offset': `${(33.3 * (INITIAL_TREE_IDX - activeTreeIdx)).toFixed(1)}%` } as CSSProperties}>
						{config.trees.map((_, treeIdx) => (
							<TalentTreePicker
								key={treeIdx}
								config={config.trees}
								graph={graph}
								points={points}
								treeIdx={treeIdx}
								limits={limits}
								active={treeIdx === activeTreeIdx}
								onChange={onChange}
							/>
						))}
					</div>
					<Button
						variant={null}
						className="absolute top-0 bottom-0 left-0 z-2 px-2 text-white lg:hidden"
						data-testid="talents-carousel-prev"
						disabled={activeTreeIdx === 0}
						onClick={() => setActiveTreeIdx(idx => Math.max(0, idx - 1))}>
						<Icon name="chevron-left" />
						<span className="sr-only">{i18n.t('talents_tab.carousel.previous', { defaultValue: 'Previous' })}</span>
					</Button>
					<Button
						variant={null}
						className="absolute top-0 right-0 bottom-0 z-2 px-2 text-white lg:hidden"
						data-testid="talents-carousel-next"
						disabled={activeTreeIdx === config.trees.length - 1}
						onClick={() => setActiveTreeIdx(idx => Math.min(config.trees.length - 1, idx + 1))}>
						<Icon name="chevron-right" />
						<span className="sr-only">{i18n.t('talents_tab.carousel.next', { defaultValue: 'Next' })}</span>
					</Button>
				</div>
			</div>
		</PickerShell>
	);
};
