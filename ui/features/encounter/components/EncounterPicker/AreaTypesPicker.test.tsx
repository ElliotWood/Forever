import { AreaType } from '@generated/proto/common';
import type { Encounter } from '@sim/raid/encounter';
import { act, fireEvent, render, screen, within } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

import { AreaTypesPicker } from './AreaTypesPicker';

const source = vi.hoisted(() => {
	const listeners = new Set<() => void>();
	return {
		listeners,
		subscribe: (onChange: () => void) => {
			listeners.add(onChange);
			return () => listeners.delete(onChange);
		},
		notify: () => Array.from(listeners).forEach(listener => listener()),
	};
});
vi.mock('@sim/state/subscriptions', async () => (await import('@sim/testing')).mockSubscriptions(source.subscribe));
vi.mock('@i18n/config', () => ({ default: { t: (key: string) => key } }));
vi.mock('@i18n/localization', () => ({
	translateAreaType: (value: number) => `area-${value}`,
}));
const trackEvent = vi.hoisted(() => vi.fn());
vi.mock('../../../../tracking/utils', () => ({ trackEvent }));

class FakeEncounter {
	areaTypes: Array<AreaType> = [];
	getAreaTypes() {
		return this.areaTypes;
	}
	setInArea(areaType: AreaType, inArea: boolean) {
		this.areaTypes = inArea ? [...this.areaTypes, areaType].sort((a, b) => a - b) : this.areaTypes.filter(t => t !== areaType);
		source.notify();
	}
}

const mount = (encounter: FakeEncounter) => render(<AreaTypesPicker encounter={encounter as unknown as Encounter} />);
const input = () => screen.getByTestId('encounter-area-types-input') as HTMLInputElement;
const pills = () => within(screen.getByTestId('encounter-area-types-selected')).getAllByTestId('combo-box-selected-chip');

beforeEach(() => {
	source.listeners.clear();
	trackEvent.mockClear();
});

describe('AreaTypesPicker', () => {
	it('labels the field, starts with no pills, and shows the encounter areas as pills', () => {
		const encounter = new FakeEncounter();
		mount(encounter);

		expect(screen.getByText('settings_tab.encounter.area_types.label').closest('label')!.getAttribute('for')).toBe('encounter-area-types');
		expect(screen.queryByTestId('encounter-area-types-selected')).toBeNull();

		act(() => encounter.setInArea(AreaType.AreaTypeHaunted, true));
		act(() => encounter.setInArea(AreaType.AreaTypeForestGrassland, true));
		expect(pills().map(pill => pill.textContent)).toEqual([`area-${AreaType.AreaTypeForestGrassland}`, `area-${AreaType.AreaTypeHaunted}`]);
	});

	it('lists the areas not yet picked that match the query, and picks one on selection', () => {
		const encounter = new FakeEncounter();
		encounter.areaTypes = [AreaType.AreaTypeForestGrassland];
		mount(encounter);

		fireEvent.change(input(), { target: { value: `area-${AreaType.AreaTypeMountainous}` } });
		const list = screen.getByTestId('encounter-area-types-list');
		const options = within(list).getAllByRole('option');
		expect(options.map(option => option.textContent)).toEqual([`area-${AreaType.AreaTypeMountainous}`]);

		fireEvent.click(options[0]);
		expect(encounter.areaTypes).toEqual([AreaType.AreaTypeForestGrassland, AreaType.AreaTypeMountainous]);
		expect(input().value).toBe('');
		expect(trackEvent).toHaveBeenCalledWith(expect.objectContaining({ category: 'area', label: 'mountainous', value: true }));
	});

	it('drops an area through its pill', () => {
		const encounter = new FakeEncounter();
		encounter.areaTypes = [AreaType.AreaTypeSnowy, AreaType.AreaTypeVolcanic];
		mount(encounter);

		fireEvent.click(screen.getAllByRole('button', { name: 'settings_tab.encounter.area_types.remove' })[0]);
		expect(encounter.areaTypes).toEqual([AreaType.AreaTypeVolcanic]);
		expect(pills()).toHaveLength(1);
		expect(trackEvent).toHaveBeenCalledWith(expect.objectContaining({ category: 'area', label: 'snowy', value: false }));
	});
});
