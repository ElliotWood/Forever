package warlock

import (
	"github.com/wowsims/forever/sim/core"
)

func (warlock *Warlock) registerLifeTap() {
	actionID := core.ActionID{SpellID: 27222}
	manaMetrics := warlock.NewManaMetrics(actionID)
	healthCost := 582.0
	baseRestore := healthCost * spellData.ImprovedLifeTap.MultiplierAt(warlock.Talents.ImprovedLifeTap)

	// TODO: Forever drops Mana Feed; no pet mana restore until we know whether the effect
	// moved onto another talent.

	warlock.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolShadow,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellLifeTap,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// Life tap adds 0.8*sp to mana restore
			restore := baseRestore + (warlock.GetSpellDamageValue(spell, nil) * 0.8)
			warlock.RemoveHealth(sim, healthCost)
			warlock.AddMana(sim, restore, manaMetrics)
		},
	})
}
