// Builds launch-tier gear sets: the best pre-raid gear in the launch item pool.
//
// The item database is generated with everything past launch already filtered out (see
// tools/database/launch_content.go), but it keeps Molten Core and Onyxia's Lair because
// both open at launch. A launch set is what a character walks into those raids wearing,
// so their loot is excluded here as well, and the pool is every other item the spec can
// wear. Within it each slot takes the item with the highest EP against the weights in
// specs.go; those are each spec's own stat weights from its ui/<spec>/sim.ts, so the sets
// value the same things the spec pages do.
//
//	go run --tags=with_db ./tools/launch_gear -spec retribution_paladin [-dir /tmp/out]
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/wowsims/forever/assets/database"
	"github.com/wowsims/forever/sim/core/proto"
)

// The two raids open at launch. A launch set is worn on the way in, not out.
var launchRaidZones = map[int32]bool{
	2717: true, // Molten Core
	2159: true, // Onyxia's Lair
}

// Onyxia's drops and the head turn-ins for her and Ragnaros carry no source record at
// all, so the zone check cannot see them. Item level can: among sourceless items nothing
// sits between 66 and 70, and everything from 71 up is one of those - Onyxia's drops at 71
// and 72, her turn-in rewards at 74, the Benediction and Rhok'delar chains at 75, Sulfuras
// at 80. Below that it is world drops, PvP sets and dungeon loot, all obtainable at launch.
const maxSourcelessIlvl = 70

func preRaid(item *proto.UIItem) bool {
	for _, source := range item.Sources {
		if drop := source.GetDrop(); drop != nil && launchRaidZones[drop.ZoneId] {
			continue
		}
		return true
	}
	return len(item.Sources) == 0 && item.ScalingOptions[0].GetIlvl() <= maxSourcelessIlvl
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

func main() {
	specName := flag.String("spec", "", "entry in specs.go, e.g. retribution_paladin")
	out := flag.String("out", "", "gear set name to write; defaults to the spec entry's")
	dir := flag.String("dir", filepath.Join("ui", "specs"), "output root; the set goes to <dir>/<class>/<spec>/gear_sets/")
	flag.Parse()

	spec, ok := specs[*specName]
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown spec %q\n", *specName)
		os.Exit(1)
	}
	if *out != "" {
		spec.set = *out
	}
	if spec.set == "" {
		spec.set = "launch"
	}

	db := database.Load()

	pool := []*proto.UIItem{}
	for _, item := range db.Items {
		if !spec.allows(item) || !preRaid(item) {
			continue
		}
		pool = append(pool, item)
	}

	picked := []map[string]any{}
	used := map[string]bool{}

	pick := func(matches func(*proto.UIItem) bool, count int) {
		candidates := []*proto.UIItem{}
		for _, item := range pool {
			if matches(item) {
				candidates = append(candidates, item)
			}
		}
		sort.SliceStable(candidates, func(i, j int) bool {
			return spec.ep(candidates[i], ranged) > spec.ep(candidates[j], ranged)
		})
		taken := 0
		for _, item := range candidates {
			if taken >= count {
				break
			}
			// Two rings or trinkets of the same name cannot both be worn, and the two
			// factions' copies of one PvP ring are the same ring.
			if used[item.Name] || used[fingerprint(item)] {
				continue
			}
			used[item.Name], used[fingerprint(item)] = true, true
			picked = append(picked, map[string]any{"id": item.Id})
			fmt.Printf("  %-9s %-42s ep %7.1f\n", item.Type.String()[8:], item.Name, spec.ep(item, ranged))
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

	for _, held := range spec.pickWeapons(pool) {
		picked = append(picked, map[string]any{"id": held.item.Id})
		fmt.Printf("  %-9s %-42s ep %7.1f\n", held.item.WeaponType.String()[10:], held.item.Name, spec.ep(held.item, held.hand))
	}

	path := filepath.Join(*dir, spec.dir, "gear_sets", spec.set+".gear.json")
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		panic(err)
	}
	body, _ := json.Marshal(map[string]any{"items": picked})
	if err := os.WriteFile(path, append(body, '\n'), 0644); err != nil {
		panic(err)
	}
	fmt.Printf("wrote %s (%d items)\n", path, len(picked))
}
