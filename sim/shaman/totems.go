package shaman

import (
	"time"

	"github.com/wowsims/classic/sim/core"
)

// The shared buff values in sim/core/buffs.go are Classic's max ranks, 77 Strength and 77 Agility. The beta client
// gives Grace of Air 89, which is Classic's value with Enhancing Totems' 15% folded in (rounded up), but lowers
// Strength of Earth to 53. Each multiplier lands the core value exactly on the client's.
const graceOfAirMultiplier = 89.0 / 77
const strengthOfEarthMultiplier = 53.0 / 77

// The beta client puts every totem on a 1 sec global cooldown and gives the ones that last 1 or 2 min in Classic
// 5 min.
const totemGCD = time.Second
const totemDuration = time.Minute * 5

func (shaman *Shaman) newTotemSpellConfig(flatCost float64, spellID int32) core.SpellConfig {
	return core.SpellConfig{
		ActionID: core.ActionID{SpellID: spellID},
		Flags:    SpellFlagShaman | SpellFlagTotem | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			FlatCost:   flatCost,
			Multiplier: shaman.totemManaMultiplier(),
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: totemGCD,
			},
		},
	}
}
