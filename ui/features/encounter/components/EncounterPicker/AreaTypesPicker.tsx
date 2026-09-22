import { AreaType } from '@generated/proto/common';
import i18n from '@i18n/config';
import { getAreaTypeI18nKey } from '@i18n/entity_mapping';
import { translateAreaType } from '@i18n/localization';
import { useStoreSubscribe } from '@sim/hooks/useStoreSubscribe';
import type { Encounter } from '@sim/raid/encounter';
import { subscribeEncounterField } from '@sim/state/subscriptions';
import { ComboBox } from '@ui-kit/ComboBox';
import { useMemo, useState } from 'react';

import { trackEvent } from '../../../../tracking/utils';

export const AREA_TYPES: ReadonlyArray<AreaType> = [
	AreaType.AreaTypeForestGrassland,
	AreaType.AreaTypeMountainous,
	AreaType.AreaTypeSnowy,
	AreaType.AreaTypeDesert,
	AreaType.AreaTypeSwamp,
	AreaType.AreaTypeWasteland,
	AreaType.AreaTypeHaunted,
	AreaType.AreaTypeCavernous,
	AreaType.AreaTypeVolcanic,
	AreaType.AreaTypeStrongholdsCities,
];

export interface AreaTypesPickerProps {
	encounter: Encounter;
}

export const AreaTypesPicker = ({ encounter }: AreaTypesPickerProps) => {
	const subscribe = useMemo(() => subscribeEncounterField(encounter, 'areaTypes'), [encounter]);
	const selected = useStoreSubscribe(subscribe, () => encounter.getAreaTypes());
	const [query, setQuery] = useState('');
	const [dismissed, setDismissed] = useState(true);

	const needle = query.trim().toLowerCase();
	const items = AREA_TYPES.filter(areaType => !selected.includes(areaType) && translateAreaType(areaType).toLowerCase().includes(needle));

	const setInArea = (areaType: AreaType, inArea: boolean) => {
		trackEvent({ action: 'settings', category: 'area', label: getAreaTypeI18nKey(areaType), value: inArea });
		encounter.setInArea(areaType, inArea);
	};

	return (
		<ComboBox
			id="encounter-area-types"
			label={i18n.t('settings_tab.encounter.area_types.label')}
			placeholder={i18n.t('settings_tab.encounter.area_types.placeholder')}
			value={query}
			onChange={next => {
				setQuery(next);
				setDismissed(false);
			}}
			open={!dismissed && items.length > 0}
			onOpenChange={next => setDismissed(!next)}
			items={items}
			itemKey={areaType => areaType}
			renderItem={areaType => <span className="p-2">{translateAreaType(areaType)}</span>}
			onItemSelect={areaType => {
				setInArea(areaType, true);
				setQuery('');
				setDismissed(true);
			}}
			closeOnSelect
			selected={selected}
			renderSelected={areaType => translateAreaType(areaType)}
			onSelectedRemove={areaType => setInArea(areaType, false)}
			removeLabel={i18n.t('settings_tab.encounter.area_types.remove')}
			footer={<div className="ui-combo-box-footer text-left">{i18n.t('settings_tab.encounter.area_types.tooltip')}</div>}
			className="w-full"
			inputTestId="encounter-area-types-input"
			listTestId="encounter-area-types-list"
			selectedTestId="encounter-area-types-selected"
		/>
	);
};
