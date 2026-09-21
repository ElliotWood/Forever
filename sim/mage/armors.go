package mage

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

// The chosen armor is up for the whole fight. Forever has no Molten Armor, so that option, like None,
// applies nothing; Frost Armor picks the top rank of the line, Ice Armor.
func (mage *Mage) registerArmorSpells() {
	var armorRank shared.SpellData
	var buff stats.Stats
	var label string
	castingRegen := 0.0

	switch mage.Options.DefaultMageArmor {
	case proto.MageArmor_MageArmorFrostArmor:
		armorRank = spellData.IceArmor.HighestRank()
		label = "Ice Armor"
		buff = stats.Stats{
			stats.Armor:           armorRank.Effect(shared.A_MOD_RESISTANCE, 1).Value,
			stats.FrostResistance: armorRank.Effect(shared.A_MOD_RESISTANCE, 16).Value,
		}
	case proto.MageArmor_MageArmorMageArmor:
		armorRank = spellData.MageArmor.HighestRank()
		label = "Mage Armor"
		resist := armorRank.Effect(shared.A_MOD_RESISTANCE, 126).Value
		buff = stats.Stats{
			stats.ArcaneResistance: resist,
			stats.FireResistance:   resist,
			stats.FrostResistance:  resist,
			stats.NatureResistance: resist,
			stats.ShadowResistance: resist,
		}
		castingRegen = armorRank.Effect(shared.A_MOD_MANA_REGEN_INTERRUPT, 0).Value / 100
	default:
		return
	}

	core.MakePermanent(mage.RegisterAura(core.Aura{
		Label:    label,
		ActionID: core.ActionID{SpellID: armorRank.SpellID},
		OnGain: func(_ *core.Aura, _ *core.Simulation) {
			mage.PseudoStats.SpiritRegenRateCasting += castingRegen
			mage.UpdateManaRegenRates()
		},
		OnExpire: func(_ *core.Aura, _ *core.Simulation) {
			mage.PseudoStats.SpiritRegenRateCasting -= castingRegen
			mage.UpdateManaRegenRates()
		},
	}).AttachStatsBuff(buff))
}
