import { SimSettingCategories } from '@sim/constants/sim_settings';
import { PlayerSpecs } from '@sim/player/specs';
import { tryParseUrlLocation } from '@sim/state/sim_links';
import { describe, expect, it } from 'vitest';

import { talentLink } from './ArenaPage';

describe('arena talent links', () => {
	it('open the spec page with only the talents', () => {
		const link = new URL(talentLink(PlayerSpecs.DpsWarrior, '30305001302-05050005525010051'), 'https://example.com');
		expect(link.pathname).toBe(PlayerSpecs.DpsWarrior.simLink);
		const parsed = tryParseUrlLocation(link);
		expect(parsed?.categories).toEqual([SimSettingCategories.Talents]);
		expect(parsed?.settings.player?.talentsString).toBe('30305001302-05050005525010051');
	});
});
