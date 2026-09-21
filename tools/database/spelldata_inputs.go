package database

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/wowsims/forever/tools/database/dbc"
)

// The client rows sim/core/spelldata is built from, captured next to the generated store so that the
// store can be rebuilt - and checked - without the client database, which is gitignored and comes
// from a local WoW install.
//
// It holds the tables the loaders in spelldata_store.go read, cut down to the spells the store
// carries, plus the three things the store's shape depends on that are resolved elsewhere: the root
// ids the closure starts from (which come from the item, enchant and set-bonus tables), the ladder
// and tree spells the class files are generated from, and the talent tree's curve points. What is
// derived from these - the closure, the hand links, the tooltip hints, the overrides and the rows
// themselves - is left to be re-derived, which is what makes the regeneration a check on the
// generator rather than a copy of its answer.

// Beside the client extraction rather than in it: assets/db_inputs/dbc is gitignored, being tens of
// megabytes rebuilt by `make db`, and this file has to be committed for the check to run without a
// database. Gzipped like the extraction's files, and named the same way, since it is read back the
// same way.
const spellStoreInputsPath = "assets/db_inputs/spell_store_inputs.json"

type storeInputs struct {
	Names        map[int32]string
	Subtexts     map[int32]string
	Descriptions map[int32]string

	Misc         map[int32]miscRow
	Levels       map[int32]levelsRow
	Cooldowns    map[int32]cooldownRow
	Categories   map[int32]categoryRow
	AuraOptions  map[int32]auraOptionRow
	ClassOptions map[int32]storeClassFlags
	Interrupts   map[int32]interruptRow
	Shapeshift   map[int32]uint64
	Targets      map[int32]int16
	Equipped     map[int32]equippedRow

	Labels  map[int32][]int16
	Powers  map[int32][]storePower
	Effects map[int32][]storeEffect

	Roots []int32

	// The talent nodes of every class tree in the order storeCurves reads them, since the first
	// definition to price a spell is the one kept, and the points each definition states by the
	// client's EffectIndex and rank.
	TraitNodes  []traitNode
	TraitPoints map[int32]map[int32]map[int32]float64
}

// The tables the generator reads, from the captured rows.
//
// The effects are cloned because linkHandTriggers writes the client's missing trigger edges into
// them: sharing the slice would put those edges back into what is written out, and the regeneration
// would then read a hand link it is supposed to be re-deriving. Nothing else is written through.
func (in *storeInputs) tables() *spellTables {
	effects := make(map[int32][]storeEffect, len(in.Effects))
	for id, rows := range in.Effects {
		effects[id] = slices.Clone(rows)
	}

	return &spellTables{
		names:        in.Names,
		subtexts:     in.Subtexts,
		descriptions: in.Descriptions,
		misc:         in.Misc,
		levels:       in.Levels,
		cooldowns:    in.Cooldowns,
		categories:   in.Categories,
		auraOptions:  in.AuraOptions,
		classOptions: in.ClassOptions,
		interrupts:   in.Interrupts,
		shapeshift:   in.Shapeshift,
		targets:      in.Targets,
		equipped:     in.Equipped,
		labels:       in.Labels,
		powers:       in.Powers,
		effects:      effects,
	}
}

// Everything the store's ids reach, and nothing else. Restricting to those ids is what keeps the
// file to the size of the store rather than the size of the client: the closure only ever reads the
// rows of a spell it has already reached, and an edge to a spell outside the set is an edge the
// closure would have followed, so no such spell exists.
func captureStoreInputs(t *spellTables, roots []int32, ids []int32,
	nodes []traitNode, points map[int32]map[int32]map[int32]float64) *storeInputs {
	in := &storeInputs{
		Names:        map[int32]string{},
		Subtexts:     map[int32]string{},
		Descriptions: map[int32]string{},
		Misc:         map[int32]miscRow{},
		Levels:       map[int32]levelsRow{},
		Cooldowns:    map[int32]cooldownRow{},
		Categories:   map[int32]categoryRow{},
		AuraOptions:  map[int32]auraOptionRow{},
		ClassOptions: map[int32]storeClassFlags{},
		Interrupts:   map[int32]interruptRow{},
		Shapeshift:   map[int32]uint64{},
		Targets:      map[int32]int16{},
		Equipped:     map[int32]equippedRow{},
		Labels:       map[int32][]int16{},
		Powers:       map[int32][]storePower{},
		Effects:      map[int32][]storeEffect{},
		Roots:        roots,
		TraitNodes:   nodes,
		TraitPoints:  points,
	}

	for _, id := range ids {
		in.Names[id] = t.names[id]
		keepString(in.Subtexts, id, t.subtexts[id])
		keepString(in.Descriptions, id, t.descriptions[id])

		keepValue(in.Misc, id, t.misc)
		keepValue(in.Levels, id, t.levels)
		keepValue(in.Cooldowns, id, t.cooldowns)
		keepValue(in.Categories, id, t.categories)
		keepValue(in.AuraOptions, id, t.auraOptions)
		keepValue(in.ClassOptions, id, t.classOptions)
		keepValue(in.Interrupts, id, t.interrupts)
		keepValue(in.Shapeshift, id, t.shapeshift)
		keepValue(in.Targets, id, t.targets)
		keepValue(in.Equipped, id, t.equipped)

		keepSlice(in.Labels, id, t.labels)
		keepSlice(in.Powers, id, t.powers)
		keepSlice(in.Effects, id, t.effects)
	}
	return in
}

// A row is kept only where the client states one: the absence of a SpellLevels row is itself a
// reading - no level scaling - so writing a zero row in its place would change what the store says.
func keepValue[V comparable](into map[int32]V, id int32, from map[int32]V) {
	if v, ok := from[id]; ok {
		into[id] = v
	}
}

func keepSlice[V any](into map[int32][]V, id int32, from map[int32][]V) {
	if v := from[id]; len(v) > 0 {
		into[id] = v
	}
}

func keepString(into map[int32]string, id int32, s string) {
	if s != "" {
		into[id] = s
	}
}

func writeStoreInputs(in *storeInputs) error {
	out, err := json.Marshal(in)
	if err != nil {
		return fmt.Errorf("encoding the store's inputs: %w", err)
	}
	if err := dbc.WriteGzipFile(spellStoreInputsPath, out); err != nil {
		return fmt.Errorf("writing %s: %w", spellStoreInputsPath, err)
	}
	fmt.Fprintf(progress, "spelldata: wrote %s, %d spells\n", spellStoreInputsPath, len(in.Names))
	return nil
}

// The committed inputs, for a caller with no client database. The path is relative to the repository
// root, the way every other path the generator reads is.
func readStoreInputs(path string) (*storeInputs, error) {
	raw, err := dbc.ReadGzipFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}

	in := &storeInputs{}
	if err := json.Unmarshal(raw, in); err != nil {
		return nil, fmt.Errorf("decoding %s: %w", path, err)
	}
	return in, nil
}
