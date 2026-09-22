import { fireEvent, render, screen, within } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

import { ComboBox, type ComboBoxProps } from './ComboBox';

const ITEMS = ['alpha', 'beta', 'gamma'];

const propsFor = (overrides: Partial<ComboBoxProps<string>> = {}): ComboBoxProps<string> => ({
	value: '',
	onChange: vi.fn(),
	items: ITEMS,
	itemKey: item => item,
	renderItem: item => <span>{item}</span>,
	onItemSelect: vi.fn(),
	open: false,
	onOpenChange: vi.fn(),
	inputTestId: 'combo-box-input',
	selectedTestId: 'combo-box-selected',
	...overrides,
});

describe('ComboBox multi-select', () => {
	it('renders nothing for the picks until there are some', () => {
		render(<ComboBox {...propsFor()} />);
		expect(screen.queryByTestId('combo-box-selected')).toBeNull();
	});

	it('shows every pick as an active pill under the input, in the order given', () => {
		render(<ComboBox {...propsFor({ selected: ['gamma', 'alpha'] })} />);

		const pills = within(screen.getByTestId('combo-box-selected')).getAllByTestId('combo-box-selected-chip');
		expect(pills.map(pill => pill.textContent)).toEqual(['gamma', 'alpha']);
		expect(pills.every(pill => pill.hasAttribute('data-active'))).toBe(true);
		expect(screen.getByTestId('combo-box-root').contains(pills[0])).toBe(true);
	});

	it('renders a pick through renderSelected when given, else through renderItem', () => {
		const { rerender } = render(<ComboBox {...propsFor({ selected: ['beta'], renderSelected: item => <b>{item.toUpperCase()}</b> })} />);
		expect(screen.getByTestId('combo-box-selected-chip').textContent).toBe('BETA');

		rerender(<ComboBox {...propsFor({ selected: ['beta'] })} />);
		expect(screen.getByTestId('combo-box-selected-chip').textContent).toBe('beta');
	});

	it('removes a pick with its pill button, without a confirmation step', () => {
		const onSelectedRemove = vi.fn();
		render(<ComboBox {...propsFor({ selected: ['alpha', 'beta'], onSelectedRemove, removeLabel: 'Remove' })} />);

		const remove = screen.getAllByRole('button', { name: 'Remove' });
		expect(remove).toHaveLength(2);
		fireEvent.click(remove[1]);
		expect(onSelectedRemove).toHaveBeenCalledWith('beta');
		expect(onSelectedRemove).toHaveBeenCalledTimes(1);
	});

	it('shows no remove button when nothing handles the removal', () => {
		render(<ComboBox {...propsFor({ selected: ['alpha'], removeLabel: 'Remove' })} />);
		expect(screen.queryByRole('button', { name: 'Remove' })).toBeNull();
	});
});
