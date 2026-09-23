import { getLang } from '@i18n/locale_service';

import { CHARACTER_LEVEL } from '../constants/mechanics';
import { Database } from './database';

export type WowheadTooltipItemParams = {
	/**
	 * @description Item ID
	 * @see item - mapped value from wowhead
	 * */
	itemId: number;
	/**
	 * @description Item level
	 * @see ilvl - mapped value from wowhead
	 * */
	itemLevel?: number;
	/**
	 * @description Level
	 * @see lvl - mapped value from wowhead
	 * */
	level?: number;
	/**
	 * @description Enchant
	 * @see ench - mapped value from wowhead
	 * */
	enchantIds?: number[];
	/**
	 * @description Gems
	 * @see gems - mapped value from wowhead
	 * */
	gemIds?: number[];
	/**
	 * @description Extra Socket
	 * @see sock - mapped value from wowhead
	 * */
	hasExtraSocket?: boolean;
	/**
	 * @description Item Set Pieces
	 * @see pcs - mapped value from wowhead
	 * */
	setPieceIds?: number[];
	/**
	 * @description Random Enchantment
	 * @see rand - mapped value from wowhead
	 * */
	randomEnchantmentId?: number;
	/**
	 * @description Reforges
	 * @see forg - mapped value from wowhead
	 * */
	reforgeId?: number;
	/**
	 * @description Upgrades
	 * @see upgd - mapped value from wowhead
	 * */
	upgradeStep?: number;
	/**
	 * @description Transmogrified to
	 * @see transmog - mapped value from wowhead
	 * */
	transmogId?: number;
};

export type WowheadTooltipSpellParams = {
	/**
	 * @description Spell ID
	 * @see spell - mapped value from wowhead
	 * */
	spellId: number;
	/**
	 * @description Level
	 * @see lvl - mapped value from wowhead
	 * */
	level?: number;
	/**
	 * @description Buff
	 * @see buff - mapped value from wowhead
	 * */
	useBuffAura?: boolean;
	/**
	 * @description Difficulty
	 * @see dd - mapped value from wowhead
	 * */
	difficultyId?: 14 | 15 | 16;
};

// Wowhead serves each expansion under its own url segment and its own `dataEnv`
// id. Porting this file to a sibling sim repo means changing the one line below:
// the domain, and every url built from it, follows. The literal-union key means
// an unmapped id is a compile error rather than a `.../undefined/...` url.
const WOWHEAD_EXPANSIONS = {
	4: 'classic',
	5: 'tbc',
	15: 'mop-classic',
	16: 'forever',
} as const;
type WowheadExpansionEnv = keyof typeof WOWHEAD_EXPANSIONS;

export const WOWHEAD_EXPANSION_ENV: WowheadExpansionEnv = 16;
export const WOWHEAD_DOMAIN = WOWHEAD_EXPANSIONS[WOWHEAD_EXPANSION_ENV];

// Wowhead's Forever data only has the items and spells it has seen change: Classic items
// like Truestrike Shoulders 404 there. Master asks Classic Era (env 4) for everything, which
// 404s on what Forever added (items past 25000, spells past 100000, a few reused Classic ids).
// So Classic Era for what Classic had, Forever for the rest. Checked 2026-09-23 against every
// tooltip the 19 spec pages ask for on load: the only 404s left are TBC ids neither one has.
const FOREVER_ONLY_SPELLS = new Set([14084]); // Improved Distract
const wowheadEnvFor = (entity: WowheadEntity, id: number): WowheadExpansionEnv => {
	if (entity === 'item') return id < 25000 ? 4 : WOWHEAD_EXPANSION_ENV;
	if (entity === 'spell') return id < 100000 && !FOREVER_ONLY_SPELLS.has(id) ? 4 : WOWHEAD_EXPANSION_ENV;
	return WOWHEAD_EXPANSION_ENV;
};
export const wowheadTooltipDomain = (entity: WowheadEntity, id: number) => {
	const env = wowheadEnvFor(entity, id);
	return { env, domain: WOWHEAD_EXPANSIONS[env] };
};

export const buildWowheadTooltipDataset = async (options: WowheadTooltipItemParams | WowheadTooltipSpellParams) => {
	const lang = getLang();
	const params = new URLSearchParams();
	const langPrefix = lang && lang != 'en' ? lang + '.' : '';
	const { env, domain } = 'spellId' in options ? wowheadTooltipDomain('spell', options.spellId) : wowheadTooltipDomain('item', options.itemId);
	params.set('domain', `${langPrefix}${domain}`);
	params.set('dataEnv', String(env));

	params.set('lvl', String(options.level || CHARACTER_LEVEL));

	if ('spellId' in options) {
		if (options.spellId) {
			params.set('spell', String(options.spellId));
		}
		if (options.useBuffAura) {
			const data = await Database.getSpellIconData(options.spellId);
			if (data.hasBuff) params.set('buff', '1');
		}
	}

	if ('itemId' in options) {
		params.set('item', String(options.itemId));
		if (options.itemLevel) {
			params.set('ilvl', String(options.itemLevel));
		}
		if (options.gemIds?.length) {
			params.set('gems', options.gemIds.join(':'));
		}
		if (options.enchantIds) {
			params.set('ench', options.enchantIds.join(':'));
		}
		if (options.reforgeId) {
			params.set('forg', String(options.reforgeId));
		}
		if (options.randomEnchantmentId) {
			params.set('rand', String(options.randomEnchantmentId));
		}
		if (typeof options.upgradeStep === 'number') {
			params.set('upgd', String(options.upgradeStep));
		}
		if (options.setPieceIds?.length) {
			params.set('pcs', options.setPieceIds.join(':'));
		}
		if (options.hasExtraSocket) {
			params.set('sock', '');
		}
		if (options.transmogId) {
			params.set('transmog', String(options.transmogId));
		}
	}

	return decodeURIComponent(params.toString());
};

export function getWowheadLanguagePrefix(): string {
	const lang = getLang();
	return lang === 'en' ? '' : `${lang}/`;
}

// Every wowhead link this app builds hangs off one of these. The entity links
// keep the bare host they have always used (domain per id, see wowheadTooltipDomain); the gear planner page is the `www.`
// form we show to users verbatim in the importer.
export const WOWHEAD_GEAR_PLANNER_URL = `https://www.wowhead.com/${WOWHEAD_DOMAIN}/gear-planner`;
export const WOWHEAD_ICON_BASE_URL = 'https://wow.zamimg.com/images/wow/icons';

type WowheadEntity = 'item' | 'spell' | 'quest' | 'npc' | 'zone';

// `https://wowhead.com/<domain>/<lang>/<entity>=<id>` — the language segment
// is empty for English.
export function wowheadEntityUrl(entity: WowheadEntity, id: number, rank = 0, definitionId = 0): string {
	const url = `https://wowhead.com/${wowheadTooltipDomain(entity, id).domain}/${getWowheadLanguagePrefix()}${entity}=${id}`;
	const params = new URLSearchParams();
	if (definitionId > 0) params.set('def', String(definitionId));
	if (rank > 0) params.set('rank', String(rank));
	const query = params.toString();
	return query ? `${url}?${query}` : url;
}

export function wowheadIconUrl(iconLabel: string, size: 'large' | 'medium' | 'small' = 'large'): string {
	return `${WOWHEAD_ICON_BASE_URL}/${size}/${iconLabel}.jpg`;
}
