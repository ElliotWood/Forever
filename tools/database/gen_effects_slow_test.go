package database

import (
	"slices"
	"testing"

	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/spelldata"
	"github.com/wowsims/forever/tools/database/dbc"
)

// Frostguard 12797's Chilled 16927 slows movement beside its melee slow, and routes to the slow shape.
// The Lobotomizer 19324's Brain Damage 1290950 deals damage beside its cast slow, and routes to the
// damage shape, whose spell carries the slow. Neither states a rate for its chance on hit, so both
// stay commented out. Mug O' Hurt 4090's Dazed 13496 slows movement alone and stays skipped.
func TestChanceOnHitSlowsRouteAndStateNoRate(t *testing.T) {
	inRepositoryRoot(t)
	instance := dbc.GetDBC()

	if ItemEffectIsSupported(instance, 13496) {
		t.Errorf("Dazed 13496, a movement slow alone, is not skipped")
	}

	for _, tc := range []struct {
		itemID  int
		spellID int32
		damage  bool
		slows   []int32
	}{
		{12797, 16927, false, []int32{1}},
		{19324, 1290950, true, []int32{2}},
	} {
		if !ItemEffectIsSupported(instance, int(tc.spellID)) {
			t.Errorf("item %d: %d is skipped", tc.itemID, tc.spellID)
			continue
		}
		if got := spelldata.Find(tc.spellID).SlowEffects(); !slices.Equal(got, tc.slows) {
			t.Errorf("%d slows the target through effects %v, want %v", tc.spellID, got, tc.slows)
		}

		item := instance.Items[tc.itemID]
		parsed := item.ToUIItem()
		parsed.ItemEffects = dbc.MergeItemEffectsForAllStates(parsed)
		i := slices.IndexFunc(parsed.ItemEffects, func(e *proto.ItemEffect) bool { return e.BuffId == tc.spellID })
		if i < 0 {
			t.Fatalf("item %d carries no effect on %d", tc.itemID, tc.spellID)
		}

		groups := map[string]Group{}
		if got := TryParseProcEffect(parsed, parsed.ItemEffects[i], instance, groups); got != EffectParseResultRefused {
			t.Errorf("item %d parsed as %v, want refused with a reason", tc.itemID, got)
			continue
		}
		entry := groups["Procs"].Entries[0]
		r := entry.Proc
		if r.Damage != tc.damage || r.Slow == tc.damage || entry.Slows == tc.damage || !r.IsWeaponProc ||
			r.TriggerSpellID != int(tc.spellID) {
			t.Errorf("item %d: damage %v, slow %v/%v, weapon proc %v, trigger %d; want a weapon proc on %d, damage %v",
				tc.itemID, r.Damage, entry.Slows, r.Slow, r.IsWeaponProc, r.TriggerSpellID, tc.spellID, tc.damage)
		}
		if want := []string{spelldata.ReasonStatesNoRate}; !slices.Equal(r.Unsupported, want) {
			t.Errorf("item %d refused for %q, want %q", tc.itemID, r.Unsupported, want)
		}
	}
}
