package mage

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

// Each fire hit adds a stack of crit until the row's charges of fire crits (4 in Forever, 3 in
// Classic) are spent.
func (mage *Mage) registerCombustionSpell() {
	if !mage.Talents.Combustion {
		return
	}

	combustionRank := spellData.Combustion.HighestRank()
	critPerStack := spellData.CombustionTriggered.HighestRank().Effect(shared.A_ADD_FLAT_MODIFIER, shared.SPELLMOD_CRITICAL_CHANCE).Value
	maxCrits := combustionRank.ProcCharges

	actionID := core.ActionID{SpellID: combustionRank.SpellID}
	cd := core.Cooldown{
		Timer:    mage.NewTimer(),
		Duration: combustionRank.Cooldown,
	}

	critMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask: MageSpellsAll,
		School:    core.SpellSchoolFire,
		Kind:      core.SpellMod_BonusCrit_Percent,
	})

	numCrits := int32(0)
	combustAura := mage.RegisterAura(core.Aura{
		Label:     "Combustion",
		ActionID:  actionID,
		Duration:  core.NeverExpires,
		MaxStacks: 20,
		OnGain: func(_ *core.Aura, _ *core.Simulation) {
			numCrits = 0
			critMod.Activate()
		},
		OnExpire: func(_ *core.Aura, sim *core.Simulation) {
			critMod.Deactivate()
			cd.Use(sim)
			mage.UpdateMajorCooldowns()
		},
		OnStacksChange: func(_ *core.Aura, _ *core.Simulation, _ int32, newStacks int32) {
			critMod.UpdateFloatValue(critPerStack * float64(newStacks))
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// Ignite's ticks are not a cast of their own and never spend a charge.
			if !result.Landed() || !spell.Matches(MageSpellsAllDamaging) || !spell.SpellSchool.Matches(core.SpellSchoolFire) {
				return
			}

			aura.AddStack(sim)
			if result.DidCrit() {
				numCrits++
				if numCrits >= maxCrits {
					aura.Deactivate(sim)
				}
			}
		},
	})

	combustSpell := mage.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,
		ClassSpellMask: MageSpellCombustion,
		Cast: core.CastConfig{
			CD: cd,
		},
		ExtraCastCondition: func(_ *core.Simulation, _ *core.Unit) bool {
			return !combustAura.IsActive()
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			combustAura.Activate(sim)
			combustAura.AddStack(sim)
		},
		RelatedSelfBuff: combustAura,
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: combustSpell,
		Type:  core.CooldownTypeDPS,
	})
}
