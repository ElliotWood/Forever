package shared

import (
	"math"
	"slices"
	"testing"
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/spelldata"
	"github.com/wowsims/forever/sim/core/stats"
)

const (
	frostguardChilled         int32 = 16927
	lobotomizerBrainDamage    int32 = 1290950
	speedTolerance                  = 1e-9
	frostguardChilledDuration       = 5 * time.Second
	brainDamageDuration             = 30 * time.Second
)

// Neither row states a rate for its chance on hit, so both stay unregistered; the tests below give
// each a 100% chance of their own.
func TestSlowProcRowsStateNoRate(t *testing.T) {
	for _, id := range []int32{frostguardChilled, lobotomizerBrainDamage} {
		got := spelldata.ItemProcUnsupported(spelldata.MustFind(id), true)
		if !slices.Equal(got, []string{spelldata.ReasonStatesNoRate}) {
			t.Errorf("%d: unsupported = %v, want only the missing rate", id, got)
		}
	}
}

func alwaysProcs(s *spelldata.Spell) {
	s.ProcChance = 100
	s.ProcChanceSource = spelldata.ProcChanceColumn
}

func newSlowProcSim(t *testing.T, weaponID int32, register func(SpellDataProc), row int32) (*core.Simulation, *testCaster, *core.Unit) {
	t.Helper()
	core.AddToDatabase(&proto.SimDatabase{Items: []*proto.SimItem{testOneHander(weaponID)}})
	register(SpellDataProc{Name: "Test Slow Weapon", ItemID: weaponID, TriggerSpellID: row, BuffSpellID: row,
		IsWeaponProc: true})

	sim := newTestCasterSim(testHands(&proto.ItemSpec{Id: weaponID}, &proto.ItemSpec{}), nil)
	caster := sim.Raid.Parties[0].Players[0].(*testCaster)
	caster.AddStatsDynamic(sim, stats.Stats{stats.SpellHitPercent: 100})
	return sim, caster, sim.Encounter.ActiveTargetUnits[0]
}

func landMainHandSwing(sim *core.Simulation, caster *testCaster, target *core.Unit) {
	caster.OnSpellHitDealt(sim, &core.Spell{ProcMask: core.ProcMaskMeleeMHAuto},
		&core.SpellResult{Target: target, Outcome: core.OutcomeHit, Damage: 100})
}

func expectSpeeds(t *testing.T, when string, target *core.Unit, melee float64, cast float64) {
	t.Helper()
	if got := target.TotalMeleeHasteMultiplier(); math.Abs(got-melee) > speedTolerance {
		t.Errorf("%s: the target's melee speed is %v×, want %v×", when, got, melee)
	}
	if got := target.TotalSpellHasteMultiplier(); math.Abs(got-cast) > speedTolerance {
		t.Errorf("%s: the target's cast speed is %v×, want %v×", when, got, cast)
	}
}

// Frostguard 12797's Chilled 16927 states A_MOD_MELEE_HASTE_3 -25 on the target for 5 s. A hit slows
// the target's melee attacks to 0.75×; a second hit 3 s later runs the slow to 5 s from then rather
// than slowing it twice, and it lifts when that runs out.
func TestFrostguardSlowsTheTargetsMeleeAttacks(t *testing.T) {
	editRow(t, frostguardChilled, alwaysProcs)
	sim, caster, target := newSlowProcSim(t, 991400, registerSpellDataSlowProc, frostguardChilled)
	expectSpeeds(t, "before the proc", target, 1, 1)

	start := sim.CurrentTime
	landMainHandSwing(sim, caster, target)
	expectSpeeds(t, "after the proc", target, 0.75, 1)

	stepPast(t, sim, start+3*time.Second)
	landMainHandSwing(sim, caster, target)
	expectSpeeds(t, "after the second proc", target, 0.75, 1)

	refreshed := sim.CurrentTime
	stepPast(t, sim, start+frostguardChilledDuration+time.Millisecond)
	expectSpeeds(t, "past the first proc's 5 s", target, 0.75, 1)

	stepPast(t, sim, refreshed+frostguardChilledDuration+time.Millisecond)
	expectSpeeds(t, "past the second proc's 5 s", target, 1, 1)
}

// The Lobotomizer 19324's Brain Damage 1290950 deals its damage and states A_MOD_CASTING_SPEED_NOT_STACK
// -25 on the target for 30 s. A landed hit slows the target's casts to 0.75× and deals the damage.
func TestLobotomizerDamagesAndSlowsTheTargetsCasts(t *testing.T) {
	editRow(t, lobotomizerBrainDamage, alwaysProcs)
	sim, caster, target := newSlowProcSim(t, 991401, registerSpellDataDamageProc, lobotomizerBrainDamage)

	start := sim.CurrentTime
	landMainHandSwing(sim, caster, target)
	expectSpeeds(t, "after the proc", target, 1, 0.75)

	brainDamage := caster.GetSpell(core.ActionID{SpellID: lobotomizerBrainDamage})
	if brainDamage == nil || brainDamage.SpellMetrics[target.UnitIndex].TotalDamage == 0 {
		t.Errorf("the proc dealt no damage alongside the slow")
	}

	stepPast(t, sim, start+brainDamageDuration-time.Millisecond)
	expectSpeeds(t, "inside the 30 s", target, 1, 0.75)

	stepPast(t, sim, start+brainDamageDuration+time.Millisecond)
	expectSpeeds(t, "past the 30 s", target, 1, 1)
}
