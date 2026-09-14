// Builds launch-tier gear sets for the specs that have none.
//
// Forever launches with no raids, so the pool is everything obtainable without one:
// dungeons, quests, crafting, reputation, world drops. Raid loot is excluded, by phase
// for the raid-release phases and by source for the rest. Within that pool each slot
// takes the item with the highest EP against the weights in specs.go, which are stated
// there rather than derived, because the specs this runs for have no gear to derive
// them from.
//
//	go run --tags=with_db ./tools/launch_gear -spec retribution_paladin
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/wowsims/classic/assets/database"
	"github.com/wowsims/classic/sim/core/proto"
	"github.com/wowsims/classic/sim/core/stats"
)

// Classic's phases 4, 5 and 6 were raid releases - Zul'Gurub, Ahn'Qiraj and
// Naxxramas - so nothing first available in them exists at Forever launch. Phases 1
// to 3 brought the dungeons, the quest and crafted gear, the reputation rewards and
// the PvP sets alongside their raids, so those are filtered by source instead.
const lastNonRaidPhase = 3

// The three raid zones inside phases 1-3.
var raidZones = map[int32]bool{
	2717: true, // Molten Core
	2677: true, // Blackwing Lair
	2159: true, // Onyxia's Lair
}

// Turn-ins gated behind those raids. Quest sources carry no zone, so they have to be
// named: these are the Onyxia head and Ragnaros/Nefarian drop turn-ins.
var raidQuests = map[int32]bool{
	7493: true, // The Great Masquerade (Onyxia, Alliance)
	7497: true, // For The Horde! (Onyxia, Horde)
	7783: true, // Thunderaan the Windseeker
	7848: true, // Examine the Vessel (Nefarian)
}

func obtainableAtLaunch(item *proto.UIItem) bool {
	if item.Phase > lastNonRaidPhase {
		return false
	}
	for _, source := range item.Sources {
		if drop := source.GetDrop(); drop != nil && raidZones[drop.ZoneId] {
			continue
		}
		if quest := source.GetQuest(); quest != nil && raidQuests[quest.Id] {
			continue
		}
		return true
	}
	// No source listed at all means a world drop or vendor item.
	return len(item.Sources) == 0
}

type slot struct {
	name  string
	types []proto.ItemType
	count int
}

var slots = []slot{
	{"head", []proto.ItemType{proto.ItemType_ItemTypeHead}, 1},
	{"neck", []proto.ItemType{proto.ItemType_ItemTypeNeck}, 1},
	{"shoulder", []proto.ItemType{proto.ItemType_ItemTypeShoulder}, 1},
	{"back", []proto.ItemType{proto.ItemType_ItemTypeBack}, 1},
	{"chest", []proto.ItemType{proto.ItemType_ItemTypeChest}, 1},
	{"wrist", []proto.ItemType{proto.ItemType_ItemTypeWrist}, 1},
	{"hands", []proto.ItemType{proto.ItemType_ItemTypeHands}, 1},
	{"waist", []proto.ItemType{proto.ItemType_ItemTypeWaist}, 1},
	{"legs", []proto.ItemType{proto.ItemType_ItemTypeLegs}, 1},
	{"feet", []proto.ItemType{proto.ItemType_ItemTypeFeet}, 1},
	{"finger", []proto.ItemType{proto.ItemType_ItemTypeFinger}, 2},
	{"trinket", []proto.ItemType{proto.ItemType_ItemTypeTrinket}, 2},
	{"ranged", []proto.ItemType{proto.ItemType_ItemTypeRanged}, 1},
}

func ep(item *proto.UIItem, weights stats.Stats) float64 {
	total := 0.0
	s := stats.FromFloatArray(item.Stats)
	for i := range weights {
		total += s[i] * weights[i]
	}
	return total
}

func main() {
	specName := flag.String("spec", "", "spec directory under ui/, e.g. retribution_paladin")
	out := flag.String("out", "launch", "gear set name to write")
	flag.Parse()

	spec, ok := specs[*specName]
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown spec %q\n", *specName)
		os.Exit(1)
	}

	db := database.Load()

	pool := []*proto.UIItem{}
	for _, item := range db.Items {
		if !spec.allows(item) || !obtainableAtLaunch(item) {
			continue
		}
		pool = append(pool, item)
	}

	picked := []map[string]any{}
	usedNames := map[string]bool{}

	pick := func(matches func(*proto.UIItem) bool, count int) {
		candidates := []*proto.UIItem{}
		for _, item := range pool {
			if matches(item) {
				candidates = append(candidates, item)
			}
		}
		sort.SliceStable(candidates, func(i, j int) bool {
			return ep(candidates[i], spec.weights) > ep(candidates[j], spec.weights)
		})
		taken := 0
		for _, item := range candidates {
			if taken >= count {
				break
			}
			// Two rings or trinkets of the same name cannot both be worn.
			if usedNames[item.Name] {
				continue
			}
			usedNames[item.Name] = true
			picked = append(picked, map[string]any{"id": item.Id})
			fmt.Printf("  %-9s %-42s ep %7.1f\n", item.Type.String()[8:], item.Name, ep(item, spec.weights))
			taken++
		}
		for ; taken < count; taken++ {
			fmt.Printf("  %-9s (nothing matched)\n", "?")
		}
	}

	for _, s := range slots {
		pick(func(item *proto.UIItem) bool {
			for _, t := range s.types {
				if item.Type == t {
					return true
				}
			}
			return false
		}, s.count)
	}

	spec.pickWeapons(pool, spec.weights, &picked)

	path := filepath.Join("ui", *specName, "gear_sets", *out+".gear.json")
	body, _ := json.Marshal(map[string]any{"items": picked})
	if err := os.WriteFile(path, append(body, '\n'), 0644); err != nil {
		panic(err)
	}
	fmt.Printf("wrote %s (%d items)\n", path, len(picked))
}
