package mage

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/stats"
)

// Forever's Arcane Blast buff (400573). Since client 70425 (16 Oct beta notes) each stack gives 15% of
// the mage's Intellect as spell damage and raises Arcane Blast's cost; the stacks last 15 sec from
// the first one and no other spell spends them.
func (mage *Mage) registerArcaneCharges() {
	if !mage.Talents.ArcaneBlast {
		return
	}

	buffRank := spellData.ArcaneBlastTriggered.Highest()
	intellectPerStack := buffRank.Effect(dbcenums.A_MOD_SPELL_DAMAGE_OF_STAT_PERCENT, 126).Average(core.CharacterLevel) / 100
	costPerStack := buffRank.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_COST)).Average(core.CharacterLevel) / 100

	spellDamageDep := mage.NewDynamicStatDependency(stats.Intellect, stats.SpellDamage, 0)
	costMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask: MageSpellArcaneBlast,
		Kind:      core.SpellMod_PowerCost_Pct_Add,
	})

	mage.ArcaneBlastAura = mage.RegisterAura(core.Aura{
		Label:     "Arcane Blast",
		ActionID:  core.ActionID{SpellID: buffRank.ID},
		Duration:  buffRank.Duration(),
		MaxStacks: int32(buffRank.MaxStack),
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.EnableDynamicStatDep(sim, spellDamageDep)
			costMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.DisableDynamicStatDep(sim, spellDamageDep)
			costMod.Deactivate()
		},
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, _ int32, newStacks int32) {
			aura.Unit.UpdateDynamicStatDep(sim, spellDamageDep, intellectPerStack*float64(newStacks))
			costMod.UpdateFloatValue(costPerStack * float64(newStacks))
		},
	})
}
