package shared

import (
	"math"
	"testing"
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/spelldata"
	"github.com/wowsims/forever/sim/core/stats"
)

// Recovery 8721's rows: the equip aura 1248761 and the 5% E_HEAL_PCT heal 1248759 it casts. The
// hint is the one its grant 1248760 states, "when you are Parried or Dodged".
const (
	recoveryTrigger int32 = 1248761
	recoveryHeal    int32 = 1248759
	recoveryHint          = core.ProcHintAttackDodged | core.ProcHintAttackParried
)

// A caster wearing a test weapon enchanted with a heal proc on the given rows, at half health so a
// heal shows on the health bar, and with the spell crit given.
func newHealProcSim(t *testing.T, weaponID, enchantID int32, healSpellID int32, spellCrit float64) (*core.Simulation, *testCaster, *core.Spell) {
	t.Helper()
	core.AddToDatabase(&proto.SimDatabase{
		Items: []*proto.SimItem{testOneHander(weaponID)},
		Enchants: []*proto.SimEnchant{{EffectId: enchantID, Name: "Test Heal Enchant",
			Type: proto.ItemType_ItemTypeWeapon}},
	})
	registerSpellDataHealProc(SpellDataProc{Name: "Test Heal Enchant", EnchantID: enchantID,
		TriggerSpellID: recoveryTrigger, BuffSpellID: healSpellID, ProcHint: recoveryHint})

	sim := newTestCasterSim(testHands(&proto.ItemSpec{Id: weaponID, Enchant: enchantID}, &proto.ItemSpec{}), nil)
	caster := sim.Raid.Parties[0].Players[0].(*testCaster)
	caster.AddStatsDynamic(sim, stats.Stats{stats.SpellCritPercent: spellCrit - caster.GetStat(stats.SpellCritPercent)})
	caster.RemoveHealth(sim, caster.MaxHealth()/2)

	heal := caster.GetSpell(core.ActionID{SpellID: healSpellID})
	if heal == nil {
		t.Fatalf("the heal %d the proc casts was not registered", healSpellID)
	}
	return sim, caster, heal
}

// What the caster heals for when the target meets its main-hand swing with the outcome.
func swingHealed(sim *core.Simulation, caster *testCaster, outcome core.HitOutcome, damage float64) float64 {
	before := caster.CurrentHealth()
	caster.OnSpellHitDealt(sim, &core.Spell{ProcMask: core.ProcMaskMeleeMHAuto},
		&core.SpellResult{Target: sim.Encounter.ActiveTargetUnits[0], Outcome: outcome, Damage: damage})
	return caster.CurrentHealth() - before
}

func stepPast(t *testing.T, sim *core.Simulation, until time.Duration) {
	t.Helper()
	sim.AddPendingAction(core.NewDelayedAction(core.DelayedActionOptions{DoAt: until, OnAction: func(*core.Simulation) {}}))
	for steps := 0; sim.CurrentTime < until; steps++ {
		if steps > 10000 {
			t.Fatalf("the sim did not reach %v", until)
		}
		sim.Step()
	}
}

// Recovery heals the wearer for 5% of its maximum health when the target dodges or parries the
// wearer's melee attack. A landed swing does not proc it, nor does the wearer dodging an attack it
// takes, and a second avoided swing inside the 10 s ProcCategoryRecovery heals nothing.
func TestRecoveryHealsWhenTheWearersAttackIsDodgedOrParried(t *testing.T) {
	sim, caster, heal := newHealProcSim(t, 990971, 990972, recoveryHeal, 0)
	want := 0.05 * caster.MaxHealth()

	if got := swingHealed(sim, caster, core.OutcomeHit, 100); got != 0 {
		t.Errorf("a landed swing healed %v, want nothing", got)
	}

	before := caster.CurrentHealth()
	enemy := sim.Encounter.ActiveTargetUnits[0]
	caster.OnSpellHitTaken(sim, &core.Spell{ProcMask: core.ProcMaskMeleeMHAuto, Unit: enemy},
		&core.SpellResult{Target: &caster.Unit, Outcome: core.OutcomeDodge})
	if got := caster.CurrentHealth() - before; got != 0 {
		t.Errorf("the wearer dodging an attack it takes healed %v, want nothing", got)
	}

	if got := swingHealed(sim, caster, core.OutcomeDodge, 0); math.Abs(got-want) > 1e-6 {
		t.Fatalf("a dodged swing healed %v, want 5%% of %v, %v", got, caster.MaxHealth(), want)
	}
	procTime := sim.CurrentTime

	if got := swingHealed(sim, caster, core.OutcomeParry, 0); got != 0 {
		t.Errorf("a parried swing inside the 10 s lockout healed %v, want nothing", got)
	}

	stepPast(t, sim, procTime+10*time.Second)
	if got := swingHealed(sim, caster, core.OutcomeParry, 0); math.Abs(got-want) > 1e-6 {
		t.Errorf("a parried swing after the lockout healed %v, want %v", got, want)
	}

	if got := heal.SpellMetrics[caster.UnitIndex].TotalHealing; math.Abs(got-2*want) > 1e-6 {
		t.Errorf("healing metrics = %v, want the two heals' %v", got, 2*want)
	}
}

// A proc's heal crits at the caster's spell crit unless its row carries ATTR_EX_2_CANT_CRIT. 1248759
// is not flagged, so a flagged copy stands in for the rows that are.
func TestProcHealCritsUnlessTheRowRulesItOut(t *testing.T) {
	sim, caster, heal := newHealProcSim(t, 990973, 990974, recoveryHeal, 100)
	want := 0.05 * caster.MaxHealth() * heal.CritDamageMultiplier(nil)
	if got := swingHealed(sim, caster, core.OutcomeDodge, 0); math.Abs(got-want) > 1e-6 {
		t.Errorf("at 100%% spell crit the heal was %v, want the critical %v", got, want)
	}

	row := spelldata.MustFind(recoveryHeal)
	original := *row
	row.Attr[dbcenums.ATTR_INDEX_EX_2] |= dbcenums.ATTR_EX_2_CANT_CRIT
	t.Cleanup(func() { *row = original })

	sim, caster, _ = newHealProcSim(t, 990975, 990976, recoveryHeal, 100)
	want = 0.05 * caster.MaxHealth()
	if got := swingHealed(sim, caster, core.OutcomeDodge, 0); math.Abs(got-want) > 1e-6 {
		t.Errorf("a heal whose row cannot crit healed %v at 100%% spell crit, want %v", got, want)
	}
}

// Darkmoon Card: Heroism's 23682 states its heal as E_HEAL 150 with a 0.4 spread, which the proc
// rolls rather than reading as a share of maximum health.
func TestProcHealRollsAStatedAmount(t *testing.T) {
	const fixedHeal int32 = 23682
	effect := spelldata.MustFind(fixedHeal).ProcHealEffect()
	if effect.Type != dbcenums.E_HEAL {
		t.Fatalf("%d heals through %v, want E_HEAL", fixedHeal, effect.Type)
	}

	sim, caster, _ := newHealProcSim(t, 990977, 990978, fixedHeal, 0)
	low, high := effect.Min(caster.Level), effect.Max(caster.Level)
	if got := swingHealed(sim, caster, core.OutcomeDodge, 0); got < low || got > high {
		t.Errorf("the heal was %v, want %v to %v", got, low, high)
	}
}
