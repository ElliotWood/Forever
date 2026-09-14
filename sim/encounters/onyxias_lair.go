package encounters

import (
	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
	"github.com/wowsims/classic/sim/core/stats"
)

// Onyxia's Lair is the one Forever tier 1 raid whose encounter is already known: it
// returns unchanged from Classic, at 40 players, alongside the two new ten and twenty
// player raids. Her stats here are Classic Era's.
func addOnyxia(bossPrefix string) {
	core.AddPresetTarget(&core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: &proto.Target{
			Id:        10184,
			Name:      "Onyxia's Lair Onyxia",
			Level:     63,
			MobType:   proto.MobType_MobTypeDragonkin,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      1_085_000,
				stats.Armor:       3731, // TODO:
				stats.AttackPower: 805,  // TODO: Unknown attack power
				// TODO: Resistances
			}.ToFloatArray(),

			SpellSchool:   proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:    2,
			MinBaseDamage: 3800,   // TODO: Estimated, she hits considerably softer than Vaelastrasz
			DamageSpread:  0.3333, // TODO:
			// The phase two air phase is not modelled, so this sims the ground phases
			// only - which is what a damage preset wants from her anyway.
			ParryHaste:       true,
			DualWield:        false,
			DualWieldPenalty: false,
			TargetInputs:     make([]*proto.TargetInput, 0),
		},
	})
	core.AddPresetEncounter("Onyxia's Lair Onyxia", []string{
		bossPrefix + "/Onyxia's Lair Onyxia",
	})
}
