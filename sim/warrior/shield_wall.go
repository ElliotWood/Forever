package warrior

import (
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
)

func (warrior *Warrior) registerShieldWall() {
	actionID := core.ActionID{SpellID: 871}
	aura := warrior.RegisterAura(core.Aura{
		Label:    "Shield Wall",
		ActionID: actionID,
		// TODO: Manual review needed -- spell 871 states a 12 second duration.
		Duration: time.Second * 12,
	}).AttachMultiplicativePseudoStatBuff(
		// TODO: Manual review needed -- spell 871 states 60% reduced damage taken.
		&warrior.PseudoStats.DamageTakenMultiplier, 0.4,
	)

	spell := warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		DefenseType:    core.DefenseTypeMelee,
		ClassSpellMask: SpellMaskShieldWall,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer: warrior.NewTimer(),
				// TODO: Manual review needed -- spell 871 states a 15 minute cooldown.
				Duration: time.Minute * 15,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.StanceMatches(DefensiveStance) && warrior.PseudoStats.CanBlock
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			aura.Activate(sim)
		},
		RelatedSelfBuff: aura,
	})

	warrior.RegisterItemSwapCallback([]proto.ItemSlot{proto.ItemSlot_ItemSlotOffHand}, func(sim *core.Simulation, slot proto.ItemSlot) {
		if !warrior.PseudoStats.CanBlock {
			aura.Deactivate(sim)
		}
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeSurvival,
		ShouldActivate: func(s *core.Simulation, c *core.Character) bool {
			if warrior.Spec == proto.Spec_SpecDpsWarrior {
				return false
			}
			return warrior.CurrentHealthPercent() < 0.4
		},
	})
}
