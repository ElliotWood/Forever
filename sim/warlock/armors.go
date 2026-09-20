package warlock

// TODO: To be implemented. Port the TBC Armors implementation below; not yet verified against the Forever client.
func (warlock *Warlock) registerArmors() {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	//
	// felArmorBonus := 100.0
	// felArmorHealingBonus := 1.2
	//
	// demonArmorBonus := 660.0
	// demonArmorSRBonus := 18.0
	//
	// if warlock.Talents.DemonicAegis > 0 {
	// 	bonusMultiplier := spellData.DemonicAegis.MultiplierAt(warlock.Talents.DemonicAegis)
	//
	// 	felArmorBonus *= bonusMultiplier
	// 	felArmorHealingBonus *= bonusMultiplier
	//
	// 	demonArmorBonus *= bonusMultiplier
	// 	demonArmorSRBonus *= bonusMultiplier
	// }
	//
	// warlock.FelArmor = warlock.RegisterAura(core.Aura{
	// 	Label:    "Fel Armor",
	// 	ActionID: core.ActionID{SpellID: 28176},
	// 	Duration: time.Minute * 30,
	// }).AttachMultiplicativePseudoStatBuff(&warlock.PseudoStats.SelfHealingMultiplier, felArmorHealingBonus).AttachStatBuff(stats.SpellDamage, felArmorBonus)
	//
	// warlock.DemonArmor = warlock.RegisterAura(core.Aura{
	// 	Label:    "Demon Armor",
	// 	ActionID: core.ActionID{SpellID: 27260},
	// 	Duration: time.Minute * 30,
	// }).AttachStatBuff(stats.Armor, demonArmorBonus).AttachStatBuff(stats.ShadowResistance, demonArmorSRBonus)
	//
	// // Armor selection
	// switch warlock.Options.Armor {
	//
	// case proto.WarlockOptions_FelArmor:
	// 	core.MakePermanent(warlock.FelArmor)
	//
	// case proto.WarlockOptions_DemonArmor:
	// 	core.MakePermanent(warlock.DemonArmor)
	// }
	//
}
