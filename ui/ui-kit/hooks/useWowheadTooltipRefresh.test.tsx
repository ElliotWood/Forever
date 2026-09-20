import { render } from '@testing-library/react';
import { afterEach, describe, expect, it, vi } from 'vitest';

import { useWowheadTooltipRefresh } from './useWowheadTooltipRefresh';

const Probe = ({ token }: { token: unknown }) => {
	useWowheadTooltipRefresh(token);
	return null;
};

const withWowhead = () => {
	const refreshLinks = vi.fn();
	(window as unknown as { WH: unknown }).WH = { Tooltips: { refreshLinks } };
	return refreshLinks;
};

afterEach(() => {
	delete (window as unknown as { WH?: unknown }).WH;
	vi.restoreAllMocks();
});

describe('useWowheadTooltipRefresh', () => {
	it('refreshes once per change of the token, not on every render', () => {
		const refreshLinks = withWowhead();
		const { rerender } = render(<Probe token="2552-" />);
		expect(refreshLinks).toHaveBeenCalledTimes(1);

		rerender(<Probe token="2552-" />);
		expect(refreshLinks).toHaveBeenCalledTimes(1);

		rerender(<Probe token="2553-" />);
		expect(refreshLinks).toHaveBeenCalledTimes(2);
	});

	it('does nothing when the wowhead script is absent', () => {
		expect(() => render(<Probe token="2552-" />)).not.toThrow();
	});

	it('does nothing when the script loaded without the tooltip api', () => {
		(window as unknown as { WH: unknown }).WH = {};
		expect(() => render(<Probe token="2552-" />)).not.toThrow();
	});
});
