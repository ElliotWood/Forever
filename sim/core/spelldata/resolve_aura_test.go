package spelldata

import (
	"testing"
	"time"

	"github.com/wowsims/forever/sim/core"
)

func TestAuraConfigFromARow(t *testing.T) {
	withResolverRows(t)
	aura := AuraConfig(Find(800))

	if aura.Label != "Bleed" {
		t.Errorf("label = %q, want the spell's name", aura.Label)
	}
	if aura.ActionID != (core.ActionID{SpellID: 800}) {
		t.Errorf("ActionID = %v, want spell 800", aura.ActionID)
	}
	if aura.Duration != time.Second*21 {
		t.Errorf("duration = %v, want 21s", aura.Duration)
	}
	if aura.MaxStacks != 5 {
		t.Errorf("max stacks = %d, want the row's 5", aura.MaxStacks)
	}
}

// The client states a permanent aura as -1, and its charges in the field the stacks share.
func TestAuraConfigPermanentRowAndCharges(t *testing.T) {
	withResolverRows(t)
	aura := AuraConfig(Find(900))

	if aura.Duration != core.NeverExpires {
		t.Errorf("duration = %v, want NeverExpires", aura.Duration)
	}
	if aura.MaxStacks != 3 {
		t.Errorf("max stacks = %d, want the row's 3 charges", aura.MaxStacks)
	}
}

func TestAuraConfigOptions(t *testing.T) {
	withResolverRows(t)
	aura := AuraConfig(Find(1000), Label("Tickless Aura - Debuff"), Permanent())

	if aura.Label != "Tickless Aura - Debuff" {
		t.Errorf("label = %q, want the caller's", aura.Label)
	}
	if aura.Duration != core.NeverExpires {
		t.Errorf("duration = %v, want NeverExpires", aura.Duration)
	}
	if aura.OnReset == nil {
		t.Error("a permanent aura got no OnReset to activate it")
	}
}

func TestDotConfigFromARow(t *testing.T) {
	withResolverRows(t)
	bleed := Find(800)
	dot := DotConfig(bleed, bleed.PeriodicEffect())

	if dot.TickLength != time.Second*3 {
		t.Errorf("tick length = %v, want 3s", dot.TickLength)
	}
	if dot.NumberOfTicks != 7 {
		t.Errorf("ticks = %d, want 21s over 3s ticks", dot.NumberOfTicks)
	}
	if dot.BonusCoefficient != 0.1 {
		t.Errorf("bonus coefficient = %v, want the effect's 0.1", dot.BonusCoefficient)
	}
	if dot.Aura.Duration != time.Second*21 {
		t.Errorf("aura duration = %v, want 21s", dot.Aura.Duration)
	}
	if dot.Aura.Label != "Bleed" {
		t.Errorf("aura label = %q, want the spell's name", dot.Aura.Label)
	}
	if dot.OnSnapshot == nil || dot.OnTick == nil {
		t.Error("a resolved dot has to come with its snapshot and tick")
	}
}

// The dot a caller would register, so the default snapshot has the unit and attack table it reads.
func dotForSnapshot(level int32) (*core.Dot, *core.Unit) {
	caster := &core.Unit{Level: level}
	target := &core.Unit{}
	caster.AttackTables = []*core.AttackTable{{Attacker: caster, Defender: target}}

	return &core.Dot{Spell: &core.Spell{Unit: caster, SpellSchool: core.SpellSchoolPhysical}}, target
}

// The default snapshot takes the effect's amount at the caster's level, so the same dot on a level 55
// caster snapshots less than on a level 60 one.
func TestDotConfigDefaultSnapshotReadsTheCastersLevel(t *testing.T) {
	withResolverRows(t)
	bleed := Find(800)
	config := DotConfig(bleed, bleed.PeriodicEffect())

	dot, target := dotForSnapshot(60)
	config.OnSnapshot(nil, target, dot)
	if dot.SnapshotBaseDamage != 80 {
		t.Errorf("snapshot at level 60 = %v, want 70 plus 10 levels of scaling", dot.SnapshotBaseDamage)
	}

	dot, target = dotForSnapshot(55)
	config.OnSnapshot(nil, target, dot)
	if dot.SnapshotBaseDamage != 75 {
		t.Errorf("snapshot at level 55 = %v, want 75", dot.SnapshotBaseDamage)
	}
}

// The callbacks are fields, so a caller that needs its own replaces them.
func TestDotConfigCallerReplacesTheCallbacks(t *testing.T) {
	withResolverRows(t)
	bleed := Find(800)
	config := DotConfig(bleed, bleed.PeriodicEffect())

	snapshots, ticks := 0, 0
	config.OnSnapshot = func(*core.Simulation, *core.Unit, *core.Dot) { snapshots++ }
	config.OnTick = func(*core.Simulation, *core.Unit, *core.Dot) { ticks++ }

	dot, target := dotForSnapshot(60)
	config.OnSnapshot(nil, target, dot)
	config.OnTick(nil, target, dot)

	if snapshots != 1 || ticks != 1 {
		t.Errorf("the caller's callbacks ran %d and %d times, want once each", snapshots, ticks)
	}
	if dot.SnapshotBaseDamage != 0 {
		t.Errorf("the default snapshot ran as well: base damage = %v", dot.SnapshotBaseDamage)
	}
}

func TestDotConfigTakesAuraOptions(t *testing.T) {
	withResolverRows(t)
	bleed := Find(800)

	if got := DotConfig(bleed, bleed.PeriodicEffect(), Label("Bleed - Off-hand")).Aura.Label; got != "Bleed - Off-hand" {
		t.Errorf("aura label = %q, want the caller's", got)
	}
}

// An effect that does not tick and a spell with no duration cannot answer a tick count, so they say
// so rather than registering a dot that never ticks.
func TestDotConfigPanicsWithoutAPeriodOrADuration(t *testing.T) {
	withResolverRows(t)

	tickless := Find(1000)
	requirePanic(t, "states no tick period", func() {
		DotConfig(tickless, tickless.EffectN(1))
	})

	permanent := Find(900)
	requirePanic(t, "states no duration to tick over", func() {
		DotConfig(permanent, Find(800).PeriodicEffect())
	})
}
