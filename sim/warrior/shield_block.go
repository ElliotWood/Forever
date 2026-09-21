package warrior

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

var shieldBlockRank = spellData.ShieldBlock.Highest()

func (warrior *Warrior) registerShieldBlock() {
	actionId := core.ActionID{SpellID: shieldBlockRank.ID}

	var spell *core.Spell
	aura := warrior.RegisterAura(core.Aura{
		Label:     "Shield Block",
		ActionID:  actionId,
		Duration:  shieldBlockRank.Duration(),
		MaxStacks: int32(shieldBlockRank.ProcCharges),
	}).
		AttachStatBuff(stats.BlockPercent, shieldBlockRank.Effect(dbcenums.A_MOD_BLOCK_PERCENT, 0).Percent()).
		AttachProcTrigger(core.ProcTrigger{
			Name:               "Shield Block - Consume",
			TriggerImmediately: true,
			Outcome:            core.OutcomeBlock,
			Callback:           core.CallbackOnSpellHitTaken,
			Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
				spell.RelatedSelfBuff.RemoveStack(sim)
			},
		})

	spell = warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionId,
		SpellSchool:    core.SpellSchoolPhysical,
		ClassSpellMask: SpellMaskShieldBlock,
		ClassFlags:     SpellFlagsShieldBlock,
		Flags:          core.SpellFlagAPL | core.SpellFlagHelpful,

		RageCost: core.RageCostOptions{
			Cost: rageCost(shieldBlockRank),
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: cooldownOf(shieldBlockRank),
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.PseudoStats.CanBlock && warrior.StanceMatches(DefensiveStance)
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			spell.RelatedSelfBuff.Activate(sim)
			spell.RelatedSelfBuff.SetStacks(sim, spell.RelatedSelfBuff.MaxStacks)
		},

		RelatedSelfBuff: aura,
	})

	warrior.RegisterItemSwapCallback([]proto.ItemSlot{proto.ItemSlot_ItemSlotOffHand}, func(sim *core.Simulation, slot proto.ItemSlot) {
		if !warrior.PseudoStats.CanBlock {
			aura.Deactivate(sim)
		}
	})
}
