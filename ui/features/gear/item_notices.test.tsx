import { Spec } from '@generated/proto/common';
import type { Database } from '@sim/proto/database';
import { render } from '@testing-library/react';
import { renderToStaticMarkup } from 'react-dom/server';
import { afterEach, describe, expect, it, vi } from 'vitest';

import { ITEM_NOTICES, MISSING_RANDOM_SUFFIX_WARNING, registerAreaStatsNotices, registerSetBonusNotices, SET_BONUS_NOTICES } from './item_notices';

vi.mock('@i18n/localization', () => ({
	translateAreaType: (value: number) => `area-${value}`,
	translateStat: (value: number) => `stat-${value}`,
}));

const markup = (itemId: number, spec: Spec = Spec.SpecUnknown) => renderToStaticMarkup(ITEM_NOTICES.get(itemId)?.[spec]);
const noticeContainer = (itemId: number, spec: Spec = Spec.SpecUnknown) => render(<>{ITEM_NOTICES.get(itemId)?.[spec]}</>).container;

describe('the item notice table', () => {
	it('renders the tentative-implementation notice', () => {
		const container = noticeContainer(95346);
		expect([...container.children].map(child => child.tagName.toLowerCase())).toEqual(['p', 'p']);
		const paragraphs = container.querySelectorAll('p');
		expect(paragraphs[0].textContent).toBe('This item is implemented, but detailed proc behavior will be confirmed on PTR.');
		expect(paragraphs[0].querySelectorAll('span')).toHaveLength(1);
		expect(paragraphs[0].querySelector('span')!.className).toBe('font-bold');
		expect(paragraphs[0].querySelector('span')!.textContent).toBe('is');
		expect(paragraphs[1].className).toBe('mb-0');
		expect(paragraphs[1].textContent).toBe('Want to help out by providing additional information? Contact us on our Discord!');
	});

	it('lists the tooltips a missing item effect carries', () => {
		const container = noticeContainer(17182);
		expect([...container.children].map(child => child.tagName.toLowerCase())).toEqual(['p', 'ul']);
		const heading = container.querySelector('p')!;
		expect(heading.className).toBe('font-bold');
		expect(heading.textContent).toBe('The following item effect (on-use or proc) is not implemented!');
		const items = container.querySelectorAll('ul > li');
		expect(items).toHaveLength(2);
		expect(items[0].textContent).toBe('Hurls a fiery ball that causes 303 Fire damage and an additional 75 damage over 10s.');
		expect(items[1].textContent).toBe('Deals 5 Fire damage to anyone who strikes you with a melee attack.');
	});

	it('renders the hand-written trinket notice', () => {
		const container = noticeContainer(94523);
		expect([...container.children].map(child => child.tagName.toLowerCase())).toEqual(['p', 'p', 'p']);
		const paragraphs = container.querySelectorAll('p');
		expect(paragraphs[0].textContent).toBe('The Agility proc on this trinket has been implemented, but the Voodoo Gnomes are not currently implemented!');
		expect(paragraphs[1].textContent).toBe('PTR testing is required in order to fit out accurate damage parameters for the Voodoo Gnomes.');
		expect(paragraphs[2].textContent).toBe('Want to help out by providing additional information? Contact us on our Discord!');
		expect(paragraphs[0].querySelectorAll('span')).toHaveLength(1);
		expect(paragraphs[0].querySelector('span')!.className).toBe('font-bold');
		expect(paragraphs[0].querySelector('span')!.textContent).toBe('not');
	});

	it('renders the random suffix warning', () => {
		const container = render(<>{MISSING_RANDOM_SUFFIX_WARNING}</>).container;
		expect([...container.children].map(child => child.tagName.toLowerCase())).toEqual(['p']);
		const p = container.querySelector('p')!;
		expect(p.className).toBe('mb-0');
		expect(p.textContent).toBe('Please select a random suffix');
	});
});

describe('registerSetBonusNotices', () => {
	const SET_ID = 9999;
	const ITEM_IDS = [90001, 90002];

	afterEach(() => {
		SET_BONUS_NOTICES.delete(SET_ID);
		ITEM_IDS.forEach(id => ITEM_NOTICES.delete(id));
	});

	// The set-bonus notices are written into the same map the pickers read at runtime, which is why
	// there is one table rather than a second copy.
	it('writes a notice into the shared table for every item in the set', () => {
		SET_BONUS_NOTICES.set(SET_ID, null);
		registerSetBonusNotices({ getItemIdsForSet: (setId: number) => (setId === SET_ID ? ITEM_IDS : []) } as unknown as Database);

		for (const id of ITEM_IDS) {
			expect(markup(id)).toBe(
				'<p class="mb-1"> This item set has the following warnings:</p>' +
					'<ul class="mb-0"><li>2-piece: Not yet implemented</li><li>4-piece: Not yet implemented</li></ul>',
			);
		}
	});
});

describe('registerAreaStatsNotices', () => {
	const RUNE = 90010;
	const PLAIN = 90011;
	const db = {
		getAllItems: () => [
			{ id: RUNE, scalingOptions: { 0: { stats: { 17: 15 }, areaStats: [{ areaType: 1, stats: { 17: 29, 18: 29 } }] } } },
			{ id: PLAIN, scalingOptions: { 0: { stats: { 17: 15 }, areaStats: [] } } },
		],
	} as unknown as Database;

	afterEach(() => {
		ITEM_NOTICES.delete(RUNE);
		ITEM_NOTICES.delete(PLAIN);
	});

	it('writes a notice naming each area and the stats it adds, and none for an item without any', () => {
		registerAreaStatsNotices(db);

		expect(markup(RUNE)).toBe(
			'<p class="mb-1">Only while the encounter is in one of these areas:</p>' + '<ul class="mb-0"><li>area-1: +29 stat-17, +29 stat-18</li></ul>',
		);
		expect(ITEM_NOTICES.has(PLAIN)).toBe(false);
	});

	it('keeps a notice the item already carries in front of the area lines', () => {
		ITEM_NOTICES.set(RUNE, { [Spec.SpecUnknown]: <p>existing</p> });
		registerAreaStatsNotices(db);

		expect(markup(RUNE).startsWith('<p>existing</p><p class="mb-1">Only while')).toBe(true);
	});
});
