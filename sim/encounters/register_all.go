package encounters

import (
	"github.com/wowsims/classic/sim/core"
)

// Forever's tier 1 opens on 9 December 2026 with three raids: Barrow Deeps at ten
// players, Hyjal Summit at twenty, and Onyxia's Lair returning at forty. Only the last
// of those has an encounter to build - Onyxia is unchanged from Classic. BlizzCon gave
// the other two a name, a size and a date and nothing else: no boss lists, no stat
// blocks, no loot. They stay unregistered until there is something real to register.
func init() {
	// TODO: Classic encounters?
	// naxxramas.Register()
	addLevel60("Classic")
	addVaelastraszTheCorrupt("Classic")
	addOnyxia("Classic")
	// TODO: Barrow Deeps (10) and Hyjal Summit (20), once their bosses are announced.
}

func AddSingleTargetBossEncounter(presetTarget *core.PresetTarget) {
	core.AddPresetTarget(presetTarget)
	core.AddPresetEncounter(presetTarget.Config.Name, []string{
		presetTarget.Path(),
	})
}
