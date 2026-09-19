package druid

import (
	"github.com/wowsims/forever/sim/core"
)

var thornsRank = spellData.Thorns.BySpellID(26992)

// Self-cast Thorns (rank 7). Reuses the core raid-buff aura, passing the
// druid's own Brambles talent points. If the Thorns raid buff is selected it
// is already registered (buffs apply before Initialize) and wins; otherwise
// the self-cast version uses the druid's actual talent.
func (druid *Druid) registerThornsSpell() {
	thornsAura := druid.GetAura("Thorns")
	if thornsAura == nil {
		thornsAura = core.ThornsAura(druid.GetCharacter(), druid.Talents.Brambles)
	}

	druid.RegisterSpell(Humanoid, core.SpellConfig{
		ActionID:       core.ActionID{SpellID: thornsRank.SpellID},
		SpellSchool:    thornsRank.SpellSchool,
		DefenseType:    thornsRank.DefenseType,
		Flags:          core.SpellFlagAPL | core.SpellFlagHelpful,
		ClassSpellMask: DruidSpellThorns,
		ProcMask:       core.ProcMaskEmpty,
		MaxRange:       thornsRank.MaxRange,

		ManaCost: core.ManaCostOptions{
			FlatCost: thornsRank.Cost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: thornsRank.GCD,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			thornsAura.Activate(sim)
		},

		RelatedSelfBuff: thornsAura,
	})
}
