package warrior

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
)

var shieldWallRank = spellData.ShieldWall.HighestRank()

func (warrior *Warrior) registerShieldWall() {
	actionID := core.ActionID{SpellID: shieldWallRank.SpellID}
	aura := warrior.RegisterAura(core.Aura{
		Label:    "Shield Wall",
		ActionID: actionID,
		Duration: shieldWallRank.Duration,
	}).AttachMultiplicativePseudoStatBuff(
		&warrior.PseudoStats.DamageTakenMultiplier,
		1+shieldWallRank.Effect(shared.A_MOD_DAMAGE_PERCENT_TAKEN, 127).Value/100,
	)

	spell := warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		DefenseType:    core.DefenseTypeMelee,
		ClassSpellMask: SpellMaskShieldWall,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: shieldWallRank.GCD,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: shieldWallRank.Cooldown,
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
