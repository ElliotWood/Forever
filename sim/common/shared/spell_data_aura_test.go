package shared

import (
	"math"
	"testing"
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
)

const (
	auraBoltID     int32 = 992001
	auraStrikeID   int32 = 992002
	auraManaBoltID int32 = 992003

	// Nature Aligned, Natural Alignment Crystal 19344's on-use: +20% spell damage (misc 126), +20%
	// healing and +20% mana cost (A_MOD_POWER_COST_SCHOOL_PCT, misc 126) on the wearer for 20 s. The
	// item's own row puts it on a 5 min cooldown in category 1141 for 20 s.
	natureAligned int32 = 23734
)

func init() {
	core.RegisterAgentFactory(
		proto.Player_RestorationShaman{},
		proto.Spec_SpecRestorationShaman,
		func(character *core.Character, _ *proto.Player, _ *proto.Raid) core.Agent {
			return &auraTester{Character: *character}
		},
		func(player *proto.Player, spec interface{}) {
			player.Spec = spec.(*proto.Player_RestorationShaman)
		},
	)
}

// A wearer with an arcane bolt, a physical strike and an arcane bolt costing 100 mana, each dealing
// 100 that no resistance or armor reduces.
type auraTester struct {
	core.Character
	bolt, strike, manaBolt *core.Spell
}

func (c *auraTester) GetCharacter() *core.Character       { return &c.Character }
func (c *auraTester) ApplyTalents()                       {}
func (c *auraTester) Reset(_ *core.Simulation)            {}
func (c *auraTester) OnGCDReady(_ *core.Simulation)       {}
func (c *auraTester) OnEncounterStart(_ *core.Simulation) {}

func (c *auraTester) Initialize() {
	c.bolt = c.RegisterSpell(testHit(auraBoltID, core.SpellSchoolArcane, core.ProcMaskSpellDamage))
	c.strike = c.RegisterSpell(testHit(auraStrikeID, core.SpellSchoolPhysical, core.ProcMaskMeleeMHSpecial))

	manaBolt := testHit(auraManaBoltID, core.SpellSchoolArcane, core.ProcMaskSpellDamage)
	manaBolt.ManaCost = core.ManaCostOptions{FlatCost: 100}
	c.manaBolt = c.RegisterSpell(manaBolt)
}

func testHit(id int32, school core.SpellSchool, mask core.ProcMask) core.SpellConfig {
	return core.SpellConfig{
		ActionID:         core.ActionID{SpellID: id},
		SpellSchool:      school,
		ProcMask:         mask,
		Flags:            core.SpellFlagIgnoreResists,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, 100, spell.OutcomeAlwaysHit)
		},
	}
}

func withAuraItem(id int32, itemType proto.ItemType, effects ...*proto.ItemEffect) {
	item := &proto.SimItem{Id: id, Name: "Test Aura Item", Type: itemType,
		ScalingOptions: map[int32]*proto.ScalingItemProperties{0: {}}, ItemEffects: effects}
	if itemType == proto.ItemType_ItemTypeWeapon {
		item = testOneHander(id)
		item.ItemEffects = effects
	}
	core.AddToDatabase(&proto.SimDatabase{Items: []*proto.SimItem{item}})
}

func auraPlayer(name string, class proto.Class, spec interface{}, slotted map[proto.ItemSlot]int32) *proto.Player {
	items := make([]*proto.ItemSpec, proto.ItemSlot_ItemSlotRanged+1)
	for i := range items {
		items[i] = &proto.ItemSpec{}
	}
	for slot, id := range slotted {
		items[slot] = &proto.ItemSpec{Id: id}
	}

	player := &proto.Player{
		Name:      name,
		Class:     class,
		Race:      proto.Race_RaceOrc,
		Buffs:     &proto.IndividualBuffs{},
		Equipment: &proto.EquipmentSpec{Items: items},
	}
	switch spec := spec.(type) {
	case *proto.Player_RestorationShaman:
		player.Spec = spec
	case *proto.Player_Hunter:
		player.Spec = spec
	}
	return player
}

func shaman(name string, slotted map[proto.ItemSlot]int32) *proto.Player {
	return auraPlayer(name, proto.Class_ClassShaman, &proto.Player_RestorationShaman{}, slotted)
}

func newAuraSim(players ...*proto.Player) *core.Simulation {
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid: &proto.Raid{Parties: []*proto.Party{{
			Players: players,
			Buffs:   &proto.PartyBuffs{},
		}}},
		Encounter: &proto.Encounter{
			Targets:  []*proto.Target{{Name: "target", Level: 60, MobType: proto.MobType_MobTypeDemon}},
			Duration: 180,
		},
	}, simsignals.CreateSignals())
	sim.Reset()
	return sim
}

func auraTesterAt(sim *core.Simulation, index int) *auraTester {
	return sim.Raid.Parties[0].Players[index].(*auraTester)
}

func dealt(sim *core.Simulation, spell *core.Spell, target *core.Unit) float64 {
	before := spell.SpellMetrics[target.UnitIndex].TotalDamage
	spell.Cast(sim, target)
	return spell.SpellMetrics[target.UnitIndex].TotalDamage - before
}

func near(got, want float64) bool {
	return math.Abs(got-want) < 1e-6
}

// Natural Alignment Crystal raises the wearer's spell damage and healing by 20% and its spells' mana
// cost by 20% for 20 s, and leaves a physical hit alone.
func TestNaturalAlignmentCrystal(t *testing.T) {
	const itemID int32 = 992101
	withAuraItem(itemID, proto.ItemType_ItemTypeTrinket, testOnUse(natureAligned, 300000, 1141, 20000))
	NewSpellDataAuraOnUse(itemID)

	sim := newAuraSim(shaman("Shaman", map[proto.ItemSlot]int32{proto.ItemSlot_ItemSlotTrinket1: itemID}))
	tester := auraTesterAt(sim, 0)
	target := sim.Encounter.ActiveTargetUnits[0]
	crystal := tester.GetSpell(core.ActionID{ItemID: itemID})
	if crystal == nil {
		t.Fatalf("the crystal registered no on-use")
	}

	check := func(when string, bolt, healing, cost float64) {
		t.Helper()
		if got := dealt(sim, tester.bolt, target); !near(got, bolt) {
			t.Errorf("%s: the arcane bolt dealt %v, want %v", when, got, bolt)
		}
		if got := dealt(sim, tester.strike, target); !near(got, 100) {
			t.Errorf("%s: the physical strike dealt %v, want 100", when, got)
		}
		if got := tester.PseudoStats.HealingDealtMultiplier; !near(got, healing) {
			t.Errorf("%s: healing dealt multiplier %v, want %v", when, got, healing)
		}
		if got := tester.manaBolt.Cost.GetCurrentCost(); !near(got, cost) {
			t.Errorf("%s: the mana bolt costs %v, want %v", when, got, cost)
		}
	}

	check("before use", 100, 1, 100)

	if !crystal.Cast(sim, target) {
		t.Fatalf("the crystal could not be used")
	}
	start := sim.CurrentTime
	check("while Nature Aligned is up", 120, 1.2, 120)

	if got := crystal.CD.ReadyAt(); got != start+5*time.Minute {
		t.Errorf("the crystal is ready again at %v, want 5 min after its use at %v", got, start)
	}

	stepPast(t, sim, start+20*time.Second+time.Millisecond)
	check("after its 20 s", 100, 1, 100)
}
