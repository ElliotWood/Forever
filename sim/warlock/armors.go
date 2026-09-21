package warlock

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

// The client ships Demon Skin and Demon Armor and no Fel Armor, so the Fel Armor option buffs
// nothing. Demonic Aegis raises both halves by 15% a point (1235316).
func (warlock *Warlock) registerArmors() {
	rank := spellData.DemonArmor.HighestRank()
	aegis := spellData.DemonicAegis.MultiplierAt(warlock.Talents.DemonicAegis)

	armorBonus := spellData.DemonArmor.EffectAt(0).ValueAt(rank.Rank) * aegis
	shadowResBonus := spellData.DemonArmor.EffectAt(1).ValueAt(rank.Rank) * aegis

	warlock.DemonArmor = warlock.RegisterAura(core.Aura{
		Label:    "Demon Armor",
		ActionID: core.ActionID{SpellID: rank.SpellID},
		Duration: core.NeverExpires,
	}).AttachStatBuff(stats.Armor, armorBonus).AttachStatBuff(stats.ShadowResistance, shadowResBonus)

	if warlock.Options.Armor == proto.WarlockOptions_DemonArmor {
		core.MakePermanent(warlock.DemonArmor)
	}
}
