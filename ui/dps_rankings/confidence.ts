// Reads a finished sim run and reports what its damage was made of, against the evidence
// manifest. The tier vocabulary itself lives in core/spells/rests, because the arena needs
// the same names without a SimResult to read them from.

import { ActionId } from '../core/proto_utils/action_id';
import { UnitMetrics } from '../core/proto_utils/sim_result';
import { spellSource } from '../core/spells/index';
import { Composition, EMPTY, Tier, TIER_ORDER } from '../core/spells/rests';

export function composition(player: UnitMetrics): Composition {
	const damage: Composition = { ...EMPTY };
	let total = 0;

	for (const action of player.actions) {
		// A pet's action is the pet's, but it is damage this build produced, and the manifest
		// knows the pet's spells too.
		if (action.damage <= 0) continue;
		total += action.damage;
		damage[tierOf(action.actionId)] += action.damage;
	}

	if (total === 0) return { ...EMPTY };
	for (const tier of TIER_ORDER) damage[tier] /= total;
	return damage;
}

/**
 * A white swing is not an ability and has no manifest entry to have; it is weapon damage
 * times attack speed, the oldest and best-tested arithmetic in the sim. Counting it as
 * unclassified put a third of a fury warrior's damage in the doubt column and said nothing
 * true - but folding it into the settled pile would be its own quiet lie, so it gets its own
 * bucket and its own name.
 *
 * An id the manifest has genuinely never classified stays separate from a known guess. They
 * are different problems: one is a number somebody decided to estimate, the other is a
 * number nobody has looked at.
 */
function tierOf(actionId: ActionId): Tier {
	if (actionId.otherId) return 'core';
	if (!actionId.spellId) return 'unknown';
	const source = spellSource(actionId.spellId);
	if (!source) return 'unknown';
	if (source.measured) return 'measured';
	if (source.source === 'unreviewed') return 'unknown';
	return source.source;
}
