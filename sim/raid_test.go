package sim

import (
	"testing"

	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
	"github.com/wowsims/classic/sim/core/stats"
)

func init() {
	RegisterAll()
}

var SimOptions = &proto.SimOptions{
	Iterations: 1,
	IsTest:     true,
}

var StandardTarget = &proto.Target{
	Stats:   stats.Stats{stats.Armor: 7684}.ToFloatArray(),
	MobType: proto.MobType_MobTypeDemon,
}

var STEncounter = &proto.Encounter{
	Duration: 300,
	Targets: []*proto.Target{
		StandardTarget,
	},
}

var BasicRaid = &proto.Raid{
	Parties: []*proto.Party{
		{
			Players: []*proto.Player{},
		},
		{
			Players: []*proto.Player{},
		},
	},
}

// Tests that we don't crash with various combinations of empty parties / blank players.
// TODO: Classic
/*
func TestSparseRaid(t *testing.T) {
	sparseRaid := &proto.Raid{
		Parties: []*proto.Party{
			{},
			{
				Players: []*proto.Player{
					{},
					{},
				},
			},
			{
				Players: []*proto.Player{
					{},
					{},
				},
			},
		},
	}

	rsr := &proto.RaidSimRequest{
		Raid:       sparseRaid,
		Encounter:  STEncounter,
		SimOptions: SimOptions,
	}

	core.RunRaidSim(rsr)
	// Don't need to check results, as long as it doesn't crash we're fine.
}
*/

// A raid with a gap in it used to take the whole sim down. Empty party slots come back
// as a UnitMetrics with none of its distribution metrics set, and the concurrency
// combiner read straight through those nil pointers - so any raid that was not packed
// from slot zero crashed on combine, which is every raid anyone actually builds.
func TestCombineHandlesEmptyRaidSlots(t *testing.T) {
	// What the sim returns for a slot nobody is standing in: a name and nothing else.
	emptySlot := func() *proto.UnitMetrics { return &proto.UnitMetrics{} }

	player := func(dps float64) *proto.UnitMetrics {
		dist := func() *proto.DistributionMetrics {
			return &proto.DistributionMetrics{Avg: dps, Hist: map[int32]int32{}, AggregatorData: &proto.AggregatorData{N: 1}}
		}
		return &proto.UnitMetrics{
			Name: "Warrior", Dps: dist(), Dpasp: dist(), Threat: dist(),
			Dtps: dist(), Tmi: dist(), Hps: dist(), Tto: dist(),
		}
	}

	result := func(dps float64) *proto.RaidSimResult {
		return &proto.RaidSimResult{
			RaidMetrics: &proto.RaidMetrics{
				Dps: &proto.DistributionMetrics{Avg: dps * 2, Hist: map[int32]int32{}, AggregatorData: &proto.AggregatorData{N: 1}},
				Hps: &proto.DistributionMetrics{Hist: map[int32]int32{}, AggregatorData: &proto.AggregatorData{}},
				Parties: []*proto.PartyMetrics{
					{
						Dps:     &proto.DistributionMetrics{Avg: dps * 2, Hist: map[int32]int32{}, AggregatorData: &proto.AggregatorData{N: 1}},
						Hps:     &proto.DistributionMetrics{Hist: map[int32]int32{}, AggregatorData: &proto.AggregatorData{}},
						Players: []*proto.UnitMetrics{player(dps), emptySlot(), emptySlot(), player(dps)},
					},
				},
			},
			EncounterMetrics: &proto.EncounterMetrics{Targets: []*proto.UnitMetrics{}},
			IterationsDone:   10,
		}
	}

	combined := core.CombineConcurrentSimResults([]*proto.RaidSimResult{result(100), result(200)}, false)

	players := combined.RaidMetrics.Parties[0].Players
	if len(players) != 4 {
		t.Fatalf("expected 4 slots back, got %d", len(players))
	}
	if got := players[0].Dps.Avg; got != 150 {
		t.Errorf("expected the occupied slot to average 150 dps, got %v", got)
	}
	if players[1].Dps != nil {
		t.Errorf("expected the empty slot to stay empty, got %v", players[1].Dps)
	}
}

func TestBasicRaid(t *testing.T) {
	t.Skip()
	rsr := &proto.RaidSimRequest{
		Raid:       BasicRaid,
		Encounter:  STEncounter,
		SimOptions: SimOptions,
	}

	core.RaidSimTest("P1 ST", t, rsr, 6323.79)
}

// To quickly debug raid sim issues, uncomment this test and copy in a request string.
/*
func testRaidString(t *testing.T, raidString string) {
	rsr := &proto.RaidSimRequest{}

	data := []byte(raidString)
	if err := protojson.Unmarshal(data, rsr); err != nil {
		panic(err)
	}

	core.RunRaidSim(rsr)
	//core.RaidSimTest("Fixed Raid", t, rsr, 10000.00)
}

func TestFixedRaid(t *testing.T) {
 	testRaidString(t, `
 	`)
}
*/
