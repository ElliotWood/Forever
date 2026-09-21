package warrior

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

func (warrior *Warrior) registerShieldBlock() {
	shieldBlockRank := spellData.ShieldBlock.HighestRank()

	actionId := core.ActionID{SpellID: shieldBlockRank.SpellID}

	aura := warrior.RegisterAura(core.Aura{
		Label:     "Shield Block",
		ActionID:  actionId,
		Duration:  shieldBlockRank.Duration,
		MaxStacks: shieldBlockRank.ProcCharges,
	}).AttachStatBuff(stats.BlockPercent, shieldBlockRank.Effect(shared.A_MOD_BLOCK_PERCENT, 0).Fraction())
	aura.AttachProcTrigger(core.ProcTrigger{
		Name:               "Shield Block - Consume",
		TriggerImmediately: true,
		Outcome:            core.OutcomeBlock,
		Callback:           core.CallbackOnSpellHitTaken,
		Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
			aura.RemoveStack(sim)
		},
	})

	warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionId,
		SpellSchool:    core.SpellSchoolPhysical,
		ClassSpellMask: SpellMaskShieldBlock,
		Flags:          core.SpellFlagAPL | core.SpellFlagHelpful,

		RageCost: core.RageCostOptions{
			Cost: shieldBlockRank.Cost,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: shieldBlockRank.Cooldown,
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.PseudoStats.CanBlock && warrior.StanceMatches(DefensiveStance)
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			aura.Activate(sim)
			aura.SetStacks(sim, aura.MaxStacks)
		},

		RelatedSelfBuff: aura,
	})

	warrior.deactivateWithoutShield(aura)
}

func (warrior *Warrior) deactivateWithoutShield(aura *core.Aura) {
	warrior.RegisterItemSwapCallback([]proto.ItemSlot{proto.ItemSlot_ItemSlotOffHand}, func(sim *core.Simulation, _ proto.ItemSlot) {
		if !warrior.PseudoStats.CanBlock {
			aura.Deactivate(sim)
		}
	})
}
