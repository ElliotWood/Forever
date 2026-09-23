package shared

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
)

// A proc's damage spell without a missile speed deals its damage without allocating: Fire Strike
// 7712 on its one target, and Fiery Blaze 6297 on each of three.
func TestProcDamageWithoutAMissileSpeedDoesNotAllocate(t *testing.T) {
	for _, tc := range []struct {
		neckID, row int32
	}{
		{991130, 7712},
		{991131, 6297},
	} {
		neck, _ := listeningNeck(tc.neckID, core.CallbackOnSpellHitDealt, tc.row)
		sim, caster := newOnUseSimAgainst(t, NewSpellDataDamageOnUse, map[int32]*proto.ItemEffect{}, neck, 3)
		spell := caster.GetSpell(core.ActionID{SpellID: tc.row})
		target := sim.Encounter.ActiveTargetUnits[0]

		if allocs := testing.AllocsPerRun(100, func() { spell.Cast(sim, target) }); allocs != 0 {
			t.Errorf("a cast of %d allocates %v times, want 0", tc.row, allocs)
		}
		if spell.SpellMetrics[target.UnitIndex].TotalDamage == 0 {
			t.Errorf("the casts of %d dealt nothing", tc.row)
		}
	}
}
