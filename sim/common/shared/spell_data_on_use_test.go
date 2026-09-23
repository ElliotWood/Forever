package shared

import (
	"math"
	"testing"
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/spelldata"
	"github.com/wowsims/forever/sim/core/stats"
)

// Infernal Lasso 219345's on-use: ItemEffect 183857 casts 443265, A_PERIODIC_DAMAGE 84 every 2 s for
// 12 s on the target, on a 5 min cooldown in category 1141 for 12 s.
const (
	lassoSpell    int32 = 443265
	lassoCategory int32 = 1141
)

func lassoOnUse() *proto.ItemEffect {
	return testOnUse(lassoSpell, 300000, lassoCategory, 12000)
}

func testOnUse(spellID int32, cooldownMs int32, categoryID int32, categoryCooldownMs int32) *proto.ItemEffect {
	return &proto.ItemEffect{BuffId: spellID, Effect: &proto.ItemEffect_OnUse{OnUse: &proto.OnUseEffect{
		CooldownMs: cooldownMs, CategoryId: categoryID, CategoryCooldownMs: categoryCooldownMs}}}
}

// A caster wearing test trinkets carrying the on-use effects, registered through the constructor, at
// half health, with its spells certain to hit and never to crit.
func newOnUseSim(t *testing.T, register func(int32), trinkets map[int32]*proto.ItemEffect) (*core.Simulation, *testCaster) {
	t.Helper()
	items := make([]*proto.ItemSpec, proto.ItemSlot_ItemSlotTrinket2+1)
	for i := range items {
		items[i] = &proto.ItemSpec{}
	}

	slot := proto.ItemSlot_ItemSlotTrinket1
	for id, effect := range trinkets {
		core.AddToDatabase(&proto.SimDatabase{Items: []*proto.SimItem{{Id: id, Name: "Test Trinket",
			Type: proto.ItemType_ItemTypeTrinket, ScalingOptions: map[int32]*proto.ScalingItemProperties{0: {}},
			ItemEffects: []*proto.ItemEffect{effect}}}})
		register(id)
		items[slot] = &proto.ItemSpec{Id: id}
		slot++
	}

	sim := newTestCasterSim(items, nil)
	caster := sim.Raid.Parties[0].Players[0].(*testCaster)
	caster.AddStatsDynamic(sim, stats.Stats{
		stats.SpellHitPercent:  100,
		stats.SpellCritPercent: -caster.GetStat(stats.SpellCritPercent),
	})
	caster.RemoveHealth(sim, caster.MaxHealth()/2)
	return sim, caster
}

func onUseSpell(t *testing.T, caster *testCaster, itemID int32) *core.Spell {
	t.Helper()
	spell := caster.GetSpell(core.ActionID{ItemID: itemID})
	if spell == nil {
		t.Fatalf("item %d registered no on-use spell", itemID)
	}
	return spell
}

// What the spell deals to the target between now and the given time after start, stepping just past
// it so a tick landing on it is counted.
func dealtUntil(t *testing.T, sim *core.Simulation, spell *core.Spell, start, after time.Duration) float64 {
	t.Helper()
	target := sim.Encounter.ActiveTargetUnits[0]
	before := spell.SpellMetrics[target.UnitIndex].TotalDamage
	stepPast(t, sim, start+after+time.Millisecond)
	return spell.SpellMetrics[target.UnitIndex].TotalDamage - before
}

// The lasso lands nothing when it is used, then deals 84 every 2 s until its 12 s run out: six ticks,
// 504 in all, on the target it was used on.
func TestOnUseDotTicksItsAmountOnTheTarget(t *testing.T) {
	const itemID int32 = 991001
	sim, caster := newOnUseSim(t, NewSpellDataDamageOnUse, map[int32]*proto.ItemEffect{itemID: lassoOnUse()})
	lasso := onUseSpell(t, caster, itemID)
	target := sim.Encounter.ActiveTargetUnits[0]

	if !lasso.Cast(sim, target) {
		t.Fatalf("the lasso could not be used")
	}
	start := sim.CurrentTime
	if !lasso.Dot(target).IsActive() {
		t.Fatalf("using the lasso applied no damage over time to the target")
	}
	if got := lasso.SpellMetrics[target.UnitIndex].TotalDamage; got != 0 {
		t.Errorf("using the lasso dealt %v at once, want nothing until the first tick", got)
	}

	for tick := 1; tick <= 6; tick++ {
		if got := dealtUntil(t, sim, lasso, start, time.Duration(tick)*2*time.Second); got != 84 {
			t.Errorf("tick %d dealt %v, want 84", tick, got)
		}
	}
	if lasso.Dot(target).IsActive() {
		t.Errorf("the damage over time is still up after its 12 s")
	}
	if got := dealtUntil(t, sim, lasso, start, 14*time.Second); got != 0 {
		t.Errorf("the damage over time dealt %v after it ran out, want nothing", got)
	}

	metrics := lasso.SpellMetrics[target.UnitIndex]
	if metrics.Casts != 1 || metrics.Ticks != 6 || metrics.TotalDamage != 504 {
		t.Errorf("metrics = %d casts, %d ticks, %v damage; want 1, 6, 504", metrics.Casts, metrics.Ticks, metrics.TotalDamage)
	}
}

// Each tick reads the spell power the wearer has when it lands. 443265 states no spell power share, so
// a copy stating 1 stands in for the rows that do.
func TestOnUseDotTicksOnTheSpellPowerOfTheTick(t *testing.T) {
	const itemID int32 = 991002
	editRow(t, lassoSpell, func(s *spelldata.Spell) { s.Effects[0].SPCoef = 1 })
	sim, caster := newOnUseSim(t, NewSpellDataDamageOnUse, map[int32]*proto.ItemEffect{itemID: lassoOnUse()})
	caster.AddStatsDynamic(sim, stats.Stats{stats.SpellDamage: -caster.GetStat(stats.SpellDamage)})
	lasso := onUseSpell(t, caster, itemID)

	lasso.Cast(sim, sim.Encounter.ActiveTargetUnits[0])
	start := sim.CurrentTime
	if got := dealtUntil(t, sim, lasso, start, 2*time.Second); got != 84 {
		t.Errorf("the first tick at no spell power dealt %v, want 84", got)
	}

	caster.AddStatsDynamic(sim, stats.Stats{stats.SpellDamage: 100})
	if got := dealtUntil(t, sim, lasso, start, 4*time.Second); got != 184 {
		t.Errorf("the tick after gaining 100 spell power dealt %v, want 184", got)
	}
}

// The lasso's own 5 min cooldown, and its category's 12 s that a second item in category 1141 waits
// out as well.
func TestOnUseRunsOnTheItemsCooldownAndCategory(t *testing.T) {
	const lassoID, otherID int32 = 991003, 991004
	sim, caster := newOnUseSim(t, NewSpellDataDamageOnUse, map[int32]*proto.ItemEffect{
		lassoID: lassoOnUse(),
		otherID: testOnUse(lassoSpell, 60000, lassoCategory, 12000),
	})
	lasso, other := onUseSpell(t, caster, lassoID), onUseSpell(t, caster, otherID)
	target := sim.Encounter.ActiveTargetUnits[0]

	for _, mcd := range []*core.Spell{lasso, other} {
		if got := caster.GetInitialMajorCooldown(mcd.ActionID); !got.Type.Matches(core.CooldownTypeDPS) {
			t.Errorf("%v is a major cooldown of type %v, want a DPS one", mcd.ActionID, got.Type)
		}
	}

	lasso.Cast(sim, target)
	start := sim.CurrentTime
	if offCooldown(sim, other) {
		t.Errorf("an item sharing category %d could be used while the lasso's category cooldown runs", lassoCategory)
	}

	stepPast(t, sim, start+12*time.Second)
	if !offCooldown(sim, other) {
		t.Errorf("an item sharing category %d could not be used once the 12 s ran out", lassoCategory)
	}
	if offCooldown(sim, lasso) {
		t.Errorf("the lasso could be used again after 12 s, inside its 5 min cooldown")
	}
	if got := lasso.CD.ReadyAt(); got != start+5*time.Minute {
		t.Errorf("the lasso is ready again at %v, want 5 min after its use at %v", got, start)
	}
}

func offCooldown(sim *core.Simulation, spell *core.Spell) bool {
	return core.BothTimersReady(spell.CD.Timer, spell.SharedCD.Timer, sim)
}

// Gem-studded Leather Belt 4262's on-use, 9163, heals the wearer for 300 with a 0.5 spread.
func TestOnUseHealsTheWearer(t *testing.T) {
	const itemID int32 = 991005
	const heal int32 = 9163
	sim, caster := newOnUseSim(t, NewSpellDataHealOnUse, map[int32]*proto.ItemEffect{itemID: testOnUse(heal, 300000, 0, 0)})
	spell := onUseSpell(t, caster, itemID)

	if got := caster.GetInitialMajorCooldown(spell.ActionID); !got.Type.Matches(core.CooldownTypeSurvival) || !spell.Flags.Matches(core.SpellFlagHelpful) {
		t.Errorf("the heal is a major cooldown of type %v, helpful %v; want a survival one cast on the wearer",
			got.Type, spell.Flags.Matches(core.SpellFlagHelpful))
	}

	effect := spelldata.MustFind(heal).ProcHealEffect()
	before := caster.CurrentHealth()
	spell.Cast(sim, &caster.Unit)
	if got, low, high := caster.CurrentHealth()-before, effect.Min(caster.Level), effect.Max(caster.Level); got < low || got > high {
		t.Errorf("the heal was %v, want %v to %v", got, low, high)
	}
	if offCooldown(sim, spell) {
		t.Errorf("the heal could be used again inside its 5 min cooldown")
	}
}

// Furbolg Medicine Pouch 16768's on-use, 20631, heals the wearer 100 every 1 s for 10 s.
func TestOnUseHotHealsTheWearerOverTime(t *testing.T) {
	const itemID int32 = 991006
	sim, caster := newOnUseSim(t, NewSpellDataHealOnUse, map[int32]*proto.ItemEffect{itemID: testOnUse(20631, 1200000, 0, 0)})
	spell := onUseSpell(t, caster, itemID)

	spell.Cast(sim, &caster.Unit)
	start := sim.CurrentTime
	for tick := 1; tick <= 10; tick++ {
		if got := healedUntil(t, sim, caster, start, time.Duration(tick)*time.Second); got != 100 {
			t.Errorf("tick %d healed %v, want 100", tick, got)
		}
	}
	if spell.SelfHot().IsActive() {
		t.Errorf("the heal over time is still up after its 10 s")
	}
	if got := spell.SpellMetrics[caster.UnitIndex].TotalHealing; got != 1000 {
		t.Errorf("healing metrics = %v, want the ten ticks' 1000", got)
	}
}

// Helm of Fire 8348's on-use, 10578, deals 331 with a spread and applies 33 every 2 s for 8 s with it.
func TestOnUseDealsItsDirectDamageAndItsDamageOverTime(t *testing.T) {
	const itemID int32 = 991007
	const fireball int32 = 10578
	sim, caster := newOnUseSim(t, NewSpellDataDamageOnUse, map[int32]*proto.ItemEffect{itemID: testOnUse(fireball, 300000, lassoCategory, 10000)})
	spell := onUseSpell(t, caster, itemID)
	target := sim.Encounter.ActiveTargetUnits[0]

	direct := spelldata.MustFind(fireball).DamageEffect()
	spell.Cast(sim, target)
	start := sim.CurrentTime
	if got, low, high := spell.SpellMetrics[target.UnitIndex].TotalDamage, direct.Min(caster.Level), direct.Max(caster.Level); got < low || got > high {
		t.Errorf("the direct damage was %v, want %v to %v", got, low, high)
	}
	if got := dealtUntil(t, sim, spell, start, 8*time.Second); math.Abs(got-4*33) > 1e-6 {
		t.Errorf("the damage over time dealt %v over its 8 s, want four ticks of 33", got)
	}
}
