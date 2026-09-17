package druid

import (
	"time"

	"github.com/wowsims/classic/sim/core"
)

func (druid *Druid) registerBerserkCD() {
	if !druid.Talents.Berserk {
		return
	}

	actionID := core.ActionID{SpellID: 50334}

	// The beta client's Berserk (417141) has a 3 min cooldown and lasts 15 sec.
	builders := []*DruidSpell{}
	druid.BerserkAura = druid.RegisterAura(core.Aura{
		Label:    "Berserk",
		ActionID: actionID,
		Duration: time.Second * 15,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			builders = core.FilterSlice(druid.DruidSpells, func(ds *DruidSpell) bool {
				return ds.Flags.Matches(SpellFlagBuilder)
			})
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range builders {
				spell.BonusCritRating += 100 * core.CritRatingPerCritChance
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range builders {
				spell.BonusCritRating -= 100 * core.CritRatingPerCritChance
			}
		},
	})

	druid.Berserk = druid.RegisterSpell(Cat|Bear, core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			druid.BerserkAura.Activate(sim)
		},
	})

	druid.AddMajorCooldown(core.MajorCooldown{
		Spell: druid.Berserk.Spell,
		Type:  core.CooldownTypeDPS,
	})
}
