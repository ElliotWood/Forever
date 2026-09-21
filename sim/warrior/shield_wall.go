package warrior

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/spelldata"
)

var shieldWallRank = spellData.ShieldWall.Highest()

func (warrior *Warrior) registerShieldWall() {
	aura := warrior.RegisterAura(spelldata.AuraConfig(shieldWallRank)).AttachMultiplicativePseudoStatBuff(
		&warrior.PseudoStats.DamageTakenMultiplier,
		1+shieldWallRank.Effect(dbcenums.A_MOD_DAMAGE_PERCENT_TAKEN, 127).Percent(),
	)

	config := spelldata.SpellConfig(&warrior.Unit, shieldWallRank)

	config.ExtraCastCondition = func(sim *core.Simulation, target *core.Unit) bool {
		return warrior.StanceMatches(DefensiveStance) && warrior.PseudoStats.CanBlock
	}

	config.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
		aura.Activate(sim)
	}

	config.RelatedSelfBuff = aura

	spell := warrior.RegisterSpell(config)

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
