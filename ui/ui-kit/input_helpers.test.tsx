import type { PartyBuffs } from '@generated/proto/buffs';
import type { Player } from '@sim/player/player';
import { ActionId } from '@sim/proto/action_id';
import { fireEvent, render, screen } from '@testing-library/react';
import { PickerShell } from '@ui-kit/PickerShell';
import { describe, expect, it } from 'vitest';

import {
	makeBooleanIconInput,
	makeClassOptionsBooleanIconInput,
	makeClassOptionsEnumIconInput,
	makeQuadstateIconInput,
	makeRotationEnumIconInput,
	makeSpecOptionsBooleanIconInput,
	makeSpecOptionsEnumIconInput,
} from './input_helpers';

// A quadstate input spreads four states over a numeric field and a second, boolean one. No buff row
// is shaped that way, so the fixture pairs two live party fields to exercise the helper.
const quadstate = (extra: { showWhen?: (modObj: PartyBuffs) => boolean } = {}) =>
	makeQuadstateIconInput<any, PartyBuffs, PartyBuffs>(
		{
			getModObject: (modObj: any) => modObj as PartyBuffs,
			getValue: (modObj: PartyBuffs) => modObj,
			setValue: (modObj: PartyBuffs, newVal: PartyBuffs) => Object.assign(modObj, newVal),
			storeField: 'raid:partyBuffs',
			...extra,
		},
		ActionId.fromSpellId(25289),
		ActionId.fromSpellId(12861),
		ActionId.fromItemId(30446),
		'battleShout',
		'totemTwisting',
	);

describe('makeQuadstateIconInput', () => {
	it('spreads its four states across the buff field and the second improved flag', () => {
		const buffs = { battleShout: 0, totemTwisting: false } as unknown as PartyBuffs;
		const input = quadstate();
		const player = buffs as unknown as Player<any>;

		expect(input.states).toBe(4);

		const roundTrip = [0, 1, 2, 3].map(value => {
			input.setValue(player, value);
			return [buffs.battleShout, buffs.totemTwisting, input.getValue(player)];
		});

		expect(roundTrip).toEqual([
			[0, false, 0],
			[1, false, 1],
			[2, false, 2],
			[2, true, 3],
		]);
	});

	// Every tristate, quadstate and multistate buff factory routes through makeNumberIconInput, so a
	// predicate it drops takes the faction gate on all of them with it.
	it('keeps showWhen, which the picker hides on', () => {
		const player = { battleShout: 0, totemTwisting: false } as unknown as Player<any>;
		const seen: PartyBuffs[] = [];

		expect(quadstate({ showWhen: modObj => (seen.push(modObj), false) }).showWhen!(player)).toBe(false);
		expect(seen).toEqual([player]);
		expect(quadstate({ showWhen: () => true }).showWhen!(player)).toBe(true);
		expect(quadstate().showWhen!(player)).toBe(true);
	});
});

// The picker is handed the player, while a buff config is written against the party or raid it
// reaches through, so every predicate has to be mapped on the way out or it can never fire.
describe('makeBooleanIconInput', () => {
	const partyInput = (enableWhen?: (party: { buffs: PartyBuffs; leaderPresent: boolean }) => boolean) =>
		makeBooleanIconInput<any, PartyBuffs, { buffs: PartyBuffs; leaderPresent: boolean }>(
			{
				getModObject: (player: Player<any>) => (player as unknown as { party: { buffs: PartyBuffs; leaderPresent: boolean } }).party,
				getValue: modObj => modObj.buffs,
				setValue: (modObj, newVal) => Object.assign(modObj.buffs, newVal),
				storeField: 'raid:partyBuffs',
				enableWhen,
			},
			ActionId.fromSpellId(25289),
			'battleShout',
		);

	it('maps enableWhen onto the mod object the config was written against', () => {
		const party = { buffs: { battleShout: false } as unknown as PartyBuffs, leaderPresent: false };
		const player = { party } as unknown as Player<any>;
		const input = partyInput(modObj => modObj.leaderPresent);

		expect(input.enableWhen!(player)).toBe(false);
		party.leaderPresent = true;
		expect(input.enableWhen!(player)).toBe(true);
	});

	it('leaves enableWhen unset when the config names none, so the picker stays enabled', () => {
		expect(partyInput().enableWhen).toBeUndefined();
	});
});

// A spec-level icon input's label and its tooltip are rendered only by PickerShell, so a factory that
// drops either leaves no type error behind — the picker simply comes out unlabelled.
describe('icon input factories', () => {
	const label = 'Maintain Judgement';
	const labelTooltip = 'Which Judgement debuff to keep active on the target.';
	const chrome = { label, labelTooltip };

	const factories: Record<string, () => { label?: string; labelTooltip?: unknown }> = {
		makeClassOptionsBooleanIconInput: () =>
			makeClassOptionsBooleanIconInput<any>({ ...chrome, fieldName: 'maintainJudgement' as never, id: ActionId.fromSpellId(27162) }),
		makeSpecOptionsBooleanIconInput: () =>
			makeSpecOptionsBooleanIconInput<any>({ ...chrome, fieldName: 'maintainJudgement' as never, id: ActionId.fromSpellId(27162) }),
		makeClassOptionsEnumIconInput: () => makeClassOptionsEnumIconInput<any, number>({ ...chrome, fieldName: 'maintainJudgement' as never, values: [] }),
		makeSpecOptionsEnumIconInput: () => makeSpecOptionsEnumIconInput<any, number>({ ...chrome, fieldName: 'maintainJudgement' as never, values: [] }),
		makeRotationEnumIconInput: () => makeRotationEnumIconInput<any, number>({ ...chrome, fieldName: 'maintainJudgement' as never, values: [] }),
	};

	// react-tooltip resolves anchors document-wide by id, so two pickers sharing one leaves every
	// labelled icon input on a tab showing all of its neighbours' tooltips at once.
	it('gives each picker a tooltip of its own', async () => {
		const other = makeClassOptionsEnumIconInput<any, number>({
			label: 'Aura',
			labelTooltip: 'Which paladin aura to activate in the prepull.',
			fieldName: 'aura' as never,
			values: [],
		});
		render(
			<>
				<PickerShell config={factories.makeRotationEnumIconInput() as any} className="ui-icon-field" hidden={false} disabled={false} />
				<PickerShell config={other as any} className="ui-icon-field" hidden={false} disabled={false} />
			</>,
		);

		fireEvent.mouseEnter(screen.getByText(label));
		expect(await screen.findByText(labelTooltip)).toBeTruthy();
		expect(screen.queryByText('Which paladin aura to activate in the prepull.')).toBeNull();
	});

	it.each(Object.keys(factories))('carries a label and its tooltip through %s into the shell', async name => {
		const config = factories[name]();
		render(<PickerShell config={{ ...config, id: 'icon-input' } as any} className="ui-icon-field" hidden={false} disabled={false} />);

		expect(screen.getByTestId('form-label').textContent).toBe(label);
		fireEvent.mouseEnter(screen.getByText(label));
		expect(await screen.findByText(labelTooltip)).toBeTruthy();
	});
});
