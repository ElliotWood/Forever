import { useEffect } from 'react';

type WowheadTooltips = { refreshLinks?: () => void };

const wowheadTooltips = (): WowheadTooltips | undefined => (window as { WH?: { Tooltips?: WowheadTooltips } }).WH?.Tooltips;

export const useWowheadTooltipRefresh = (token: unknown) => {
	useEffect(() => {
		wowheadTooltips()?.refreshLinks?.();
	}, [token]);
};
