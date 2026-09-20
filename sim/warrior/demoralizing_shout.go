package warrior

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

var demoralizingShoutRank = shared.WithSpellDataFlatThreat(spellData.DemoralizingShout, 56).HighestRank()

// TODO: Manual review needed -- this was modelled during the Forever port, not carried
// over unchanged, so its numbers and shape want checking against the client.
func (war *Warrior) registerDemoralizingShout() {
	war.DemoralizingShoutAuras = war.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		// Nothing in the warrior tree prices this shout: Forever has no Improved
		// Demoralizing Shout, and Booming Voice modifies a field no buff reads,
		// so the aura is the client's -204 attack power for 45 seconds.
		return core.DemoralizingShoutAura(target, true, 0)
	})

	war.DemoralizingShout = war.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: demoralizingShoutRank.SpellID},
		SpellSchool:    demoralizingShoutRank.SpellSchool,
		DefenseType:    demoralizingShoutRank.DefenseType,
		ClassSpellMask: SpellMaskDemoralizingShout,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,

		RageCost: core.RageCostOptions{
			Cost: demoralizingShoutRank.Cost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: demoralizingShoutRank.GCD,
			},
			IgnoreHaste: true,
		},

		ThreatMultiplier: 1,
		FlatThreatBonus:  demoralizingShoutRank.FlatThreatBonus,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.ActiveTargetUnits {
				result := spell.CalcAndDealOutcome(sim, aoeTarget, spell.OutcomeMagicHit)
				if result.Landed() {
					war.DemoralizingShoutAuras.Get(aoeTarget).Activate(sim)
				}
			}
		},

		RelatedAuraArrays: war.DemoralizingShoutAuras.ToMap(),
	})
}
