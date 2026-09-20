import { Race } from '@generated/proto/common';
import { SimHostProvider } from '@sim/context/SimHostContext';
import { fakeHost } from '@sim/testing';
import { renderHook } from '@testing-library/react';
import type { ReactNode } from 'react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

import { useSavedSettings } from './useSavedSettings';

vi.mock('@sim/state/subscriptions', async () => (await import('@sim/testing')).mockSubscriptions());

let key = '';

const store = (entries: Record<string, unknown>) => window.localStorage.setItem(key, JSON.stringify(entries));

const load = () =>
	renderHook(() => useSavedSettings(), {
		wrapper: ({ children }: { children: ReactNode }) => (
			<SimHostProvider host={fakeHost({ getSavedSettingsStorageKey: () => key })}>{children}</SimHostProvider>
		),
	}).result.current;

beforeEach(() => {
	// A fresh key per test: `useTypedLocalStorage` caches the raw string per key across a file.
	key = `tbc-test-savedSettings-${Math.random()}`;
});

describe('useSavedSettings', () => {
	it('reads a current entry', () => {
		store({ Raid: { race: 'RaceOrc' } });

		const { entries } = load();
		expect(entries.map(entry => entry.name)).toEqual(['Raid']);
		expect(entries[0].data.race).toBe(Race.RaceOrc);
	});

	// A save written while the field was a TristateEffect holds an enum name, which
	// `SavedSettings.fromJson` throws on now that the field is a bool; `useSavedData` swallows that
	// and the entry would be gone from the panel without a word.
	it('keeps an entry whose improvedSealOfTheCrusader is a saved enum name', () => {
		store({ Legacy: { race: 'RaceOrc', debuffs: { improvedSealOfTheCrusader: 'TristateEffectRegular', huntersMark: true } } });

		const { entries } = load();
		expect(entries.map(entry => entry.name)).toEqual(['Legacy']);
		expect(entries[0].data.debuffs?.improvedSealOfTheCrusader).toBe(true);
		expect(entries[0].data.debuffs?.huntersMark).toBe(true);
	});

	// The same swallowing would hide an entry that names a field api version 17 retired.
	it('keeps an entry that still carries a retired buff', () => {
		store({ Legacy: { race: 'RaceOrc', debuffs: { misery: true, huntersMark: true }, partyBuffs: { drums: 'LesserDrumsOfBattle' } } });

		const { entries } = load();
		expect(entries.map(entry => entry.name)).toEqual(['Legacy']);
		expect(entries[0].data.debuffs?.huntersMark).toBe(true);
	});

	it('reads the missing state as off rather than as a set buff', () => {
		store({ Current: { debuffs: { improvedSealOfTheCrusader: 'TristateEffectMissing' } } });

		expect(load().entries[0].data.debuffs?.improvedSealOfTheCrusader).toBe(false);
	});

	it('does not drop the other entries of the slot', () => {
		store({ Legacy: { debuffs: { improvedSealOfTheCrusader: 'TristateEffectImproved' } }, Current: { race: 'RaceOrc' } });

		expect(load().entries.map(entry => entry.name)).toEqual(['Legacy', 'Current']);
	});
});
