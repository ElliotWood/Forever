import type { TalentsConfig } from '@sim/talents/config';
import { newTalentsConfig } from '@sim/talents/config';
import { createEvent, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, describe, expect, it, vi } from 'vitest';

import { DEFAULT_TALENT_LIMITS } from '../../model/can_set_points';
import { parseTalentsString, serializeTalentsString } from '../../model/talents_string';
import { buildTalentGraph } from '../../model/tree_graph';
import { TalentTreePicker } from './TalentTreePicker';

// The icon and the wowhead href come from the database, which no unit test has; echoing the spell id
// back is what lets the rank-aware-icon case see which rank was asked for.
vi.mock('@ui-kit/hooks/useActionId', () => ({
	useActionId: (actionId: { spellId?: number }) => ({ iconUrl: `icon-${actionId.spellId}`, name: '', href: `href-${actionId.spellId}`, ready: true }),
}));

type Fake = Record<string, never>;

const COLS = 4;

const mkTalent = (index: number, overrides: Record<string, unknown> = {}) => ({
	fieldName: `talent${index}`,
	fancyName: `Talent ${index}`,
	location: { rowIdx: Math.floor(index / COLS), colIdx: index % COLS },
	spellIds: [100 * index + 1, 100 * index + 2, 100 * index + 3, 100 * index + 4, 100 * index + 5],
	maxPoints: 5,
	...overrides,
});

const plainTree = (name: string) => ({
	name,
	backgroundUrl: `${name}.jpg`,
	talents: Array.from({ length: 5 * COLS }, (_, index) => mkTalent(index)),
});

const config: TalentsConfig<Fake> = newTalentsConfig<Fake>([
	{
		name: 'First',
		backgroundUrl: 'first.jpg',
		talents: [
			mkTalent(0),
			mkTalent(1, { maxPoints: 1, spellIds: [999] }),
			mkTalent(2),
			mkTalent(3),
			mkTalent(4),
			mkTalent(5, { prereqLocation: { rowIdx: 0, colIdx: 1 } }),
			mkTalent(6),
			mkTalent(7),
			mkTalent(8),
			mkTalent(9),
			mkTalent(10),
			mkTalent(11),
		],
	},
	plainTree('Second'),
	plainTree('Third'),
]);

const graph = buildTalentGraph(config);

const tree = (talentsString: string, treeIdx = 0) => {
	const onChange = vi.fn();
	render(
		<TalentTreePicker
			config={config}
			graph={graph}
			points={parseTalentsString(config, talentsString)}
			treeIdx={treeIdx}
			limits={DEFAULT_TALENT_LIMITS}
			active
			onChange={onChange}
		/>,
	);
	return { onChange, written: () => serializeTalentsString(onChange.mock.calls.at(-1)![0]) };
};

const talents = () => Array.from(document.querySelectorAll<HTMLAnchorElement>('[data-testid="talent-picker-root"]'));
const arrows = () => Array.from(document.querySelectorAll<HTMLElement>('[data-testid="talent-req-arrow"]'));

afterEach(() => vi.useRealTimers());

describe('TalentTreePicker layout', () => {
	it('lays every talent of the tree into ONE grid, placed by row and column', () => {
		tree('');
		const main = screen.getByTestId('talent-tree-main');
		expect(main.style.gridTemplateRows).toBe(`repeat(${graph.numRows}, 1fr)`);
		expect(main.style.gridTemplateColumns).toBe(`repeat(${graph.numCols}, 1fr)`);
		expect(talents()).toHaveLength(config[0].talents.length);
		expect(talents()[5].style.gridRow).toBe('2');
		expect(talents()[5].style.gridColumn).toBe('2');
	});

	it('shows each talent as N/M rather than a name label', () => {
		tree('35');
		const badges = Array.from(document.querySelectorAll('[data-testid="talent-picker-points"]')).map(el => el.textContent);
		expect(badges.slice(0, 3)).toEqual(['3/5', '1/1', '0/5']);
	});

	it('reports the tree total against the 61-point budget', () => {
		tree('55');
		expect(screen.getByTestId('talent-tree-points').textContent).toBe('6 / 51');
	});

	it('names the tree from the config, not from the spec', () => {
		tree('', 2);
		expect(screen.getByTestId('talent-tree-title').textContent).toBe('Third');
	});
});

describe('TalentTreePicker icons', () => {
	it('asks for a different spell id per rank', () => {
		tree('');
		expect(talents()[0].style.backgroundImage).toBe('url("icon-1")');
	});

	it('follows the rank up as points are spent', () => {
		tree('3');
		expect(talents()[0].style.backgroundImage).toBe('url("icon-3")');
	});
});

describe('TalentTreePicker tier gating', () => {
	it('refuses a row-1 talent until five points sit in the tree', () => {
		tree('4');
		expect(talents()[4].dataset.canAdd).toBe('false');
	});

	it('allows it once the tier is met', () => {
		tree('5');
		expect(talents()[4].dataset.canAdd).toBe('true');
	});

	it('gates on points spent anywhere in the tree, not on character level', () => {
		tree('00005');
		expect(talents()[4].dataset.canAdd).toBe('true');
	});

	it('leaves the string untouched when a gated click is rejected', () => {
		const { onChange, written } = tree('4');
		fireEvent.mouseDown(talents()[4], { button: 0 });
		expect(onChange).toHaveBeenCalled();
		expect(written()).toBe('4');
	});
});

describe('TalentTreePicker prerequisites', () => {
	it('refuses a talent whose prerequisite is unfilled', () => {
		tree('50');
		expect(talents()[5].dataset.canAdd).toBe('false');
	});

	it('allows it once the prerequisite is full', () => {
		tree('51');
		expect(talents()[5].dataset.canAdd).toBe('true');
	});

	it('refuses to free a prerequisite that a spent child depends on', () => {
		const { written } = tree('510001');
		fireEvent.mouseDown(talents()[1], { button: 2 });
		expect(written()).toBe('510001');
	});

	it('frees it once the child is empty', () => {
		const { written } = tree('51');
		fireEvent.mouseDown(talents()[1], { button: 2 });
		expect(written()).toBe('5');
	});
});

describe('TalentTreePicker arrows', () => {
	it('draws one arrow per prerequisite edge, spanning the grid cells between them', () => {
		tree('');
		expect(arrows()).toHaveLength(1);
		const [arrow] = arrows();
		expect(arrow.dataset.reqDir).toBe('down');
		expect(arrow.dataset.reqArrowRowSize).toBe('1');
		expect(arrow.style.gridRow).toBe('1 / 3');
		expect(arrow.style.gridColumn).toBe('2 / 2');
	});

	it('lights the arrow only when the child can take a point or is already full', () => {
		tree('50');
		expect(arrows()[0].dataset.reqActive).toBeUndefined();
		document.body.innerHTML = '';

		tree('51');
		expect(arrows()[0].dataset.reqActive).toBe('true');
	});
});

describe('TalentTreePicker mouse', () => {
	it('spends a point on left mousedown', () => {
		const { written } = tree('');
		fireEvent.mouseDown(talents()[0], { button: 0 });
		expect(written()).toBe('1');
	});

	it('unspends a point on right mousedown', () => {
		const { written } = tree('3');
		fireEvent.mouseDown(talents()[0], { button: 2 });
		expect(written()).toBe('2');
	});

	it('stops at the talent maximum rather than wrapping', () => {
		const { written } = tree('5');
		fireEvent.mouseDown(talents()[0], { button: 0 });
		expect(written()).toBe('5');
	});

	it('does not navigate the wowhead link a talent carries', () => {
		tree('');
		const click = new MouseEvent('click', { bubbles: true, cancelable: true });
		talents()[0].dispatchEvent(click);
		expect(click.defaultPrevented).toBe(true);
	});

	it('suppresses the context menu so a right click can unspend', () => {
		tree('');
		const menu = new MouseEvent('contextmenu', { bubbles: true, cancelable: true });
		talents()[0].dispatchEvent(menu);
		expect(menu.defaultPrevented).toBe(true);
	});
});

describe('TalentTreePicker touch', () => {
	// The touch handlers are native listeners rather than React props precisely so this holds: React
	// attaches touchstart passively, where preventDefault is a no-op.
	it('prevents the default on touchstart, which a passive React listener could not', () => {
		tree('');
		const start = createEvent.touchStart(talents()[0]);
		fireEvent(talents()[0], start);
		expect(start.defaultPrevented).toBe(true);
	});

	it('spends a point on a short tap', () => {
		const { written } = tree('');
		fireEvent.touchStart(talents()[0]);
		fireEvent.touchEnd(talents()[0]);
		expect(written()).toBe('1');
	});

	it('wraps a tap on a full talent back to nothing, where a click would have stuck', () => {
		const { written } = tree('5');
		fireEvent.touchStart(talents()[0]);
		fireEvent.touchEnd(talents()[0]);
		expect(written()).toBe('');
	});

	it('clears the talent on a long press, and the release that follows changes nothing', () => {
		vi.useFakeTimers();
		const { onChange, written } = tree('3');
		fireEvent.touchStart(talents()[0]);
		vi.advanceTimersByTime(750);
		expect(written()).toBe('');

		onChange.mockClear();
		fireEvent.touchEnd(talents()[0]);
		expect(onChange).not.toHaveBeenCalled();
	});

	it('cancels the long press when the finger moves', () => {
		vi.useFakeTimers();
		const { onChange } = tree('3');
		fireEvent.touchStart(talents()[0]);
		fireEvent.touchMove(talents()[0]);
		vi.advanceTimersByTime(750);
		expect(onChange).not.toHaveBeenCalled();
	});
});

describe('TalentTreePicker reset', () => {
	it('zeroes its own tree and leaves the others alone', () => {
		const onChange = vi.fn();
		render(
			<TalentTreePicker
				config={config}
				graph={graph}
				points={parseTalentsString(config, '51-5-5')}
				treeIdx={0}
				limits={DEFAULT_TALENT_LIMITS}
				active
				onChange={onChange}
			/>,
		);
		fireEvent.click(screen.getByTestId('talent-tree-reset'));
		expect(serializeTalentsString(onChange.mock.calls.at(-1)![0])).toBe('-5-5');
	});
});
