package mage

import (
	"github.com/wowsims/forever/sim/core"
)

// Resets every frost spell with a cooldown.
func (mage *Mage) registerColdSnapSpell() {
	if !mage.Talents.ColdSnap {
		return
	}

	coldSnapRank := spellData.ColdSnap.HighestRank()

	spell := mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: coldSnapRank.SpellID},
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,
		ClassSpellMask: MageSpellColdSnap,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: coldSnapRank.Cooldown,
			},
		},
		ApplyEffects: func(_ *core.Simulation, _ *core.Unit, _ *core.Spell) {
			for _, frostSpell := range mage.Spellbook {
				if frostSpell.SpellSchool.Matches(core.SpellSchoolFrost) && frostSpell.CD.Timer != nil {
					frostSpell.CD.Reset()
				}
			}
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
	})
}
