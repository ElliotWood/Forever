package paladin

func (paladin *Paladin) registerHolyTalents() {
	// Tier 1
	paladin.applyImprovedHolyStrike()
	paladin.applyDivineStrength()
	paladin.applyDivineIntellect()

	// Tier 2
	paladin.applyHealingLight()
	paladin.applySpiritualFocus()
	paladin.applyImprovedSeals()
	paladin.applyUnyieldingFaith()

	// Tier 3
	paladin.applyVoiceOfTruth()
	paladin.applyReverence()
	paladin.applyPurifyingPower()

	// Tier 4
	paladin.applyInfusionOfLight()
	paladin.applyIllumination()
	// Divine Favor registered in registerTalentSpells

	// Tier 5
	paladin.applyDivinePrecision()
	// Holy Shock registered in registerTalentSpells
	paladin.applyConsecratedGround()

	// Tier 6
	paladin.applyHolyPowerTalent()

	// Tier 7
	paladin.applyLightsVigil()
}

// TODO: To be implemented. TBC body below needs no porting; kept commented until this class's port is reviewed.
//
// Divine Strength - Increases your total Strength by 2/4/6/8/10%
func (paladin *Paladin) applyDivineStrength() {
	if paladin.Talents.DivineStrength == 0 {
		return
	}

	// The TBC implementation, kept for the port:
	// if paladin.Talents.DivineStrength == 0 {
	// 	return
	// }
	//
	// paladin.MultiplyStat(stats.Strength, 1+(float64(paladin.Talents.DivineStrength)*.02))
}

// TODO: To be implemented. TBC body below needs no porting; kept commented until this class's port is reviewed.
//
// Divine Intellect - Increases your total Intellect by 2/4/6/8/10%
func (paladin *Paladin) applyDivineIntellect() {
	if paladin.Talents.DivineIntellect == 0 {
		return
	}

	// The TBC implementation, kept for the port:
	// if paladin.Talents.DivineIntellect == 0 {
	// 	return
	// }
	//
	// bonus := 1.0 + 0.02*float64(paladin.Talents.DivineIntellect)
	// paladin.MultiplyStat(stats.Intellect, bonus)
}

// TODO: To be implemented. TBC body below needs no porting; kept commented until this class's port is reviewed.
//
// Healing Light - Increases the amount healed by your Holy Light and Flash of Light spells by 4/8/12%
func (paladin *Paladin) applyHealingLight() {
	if paladin.Talents.HealingLight == 0 {
		return
	}

	// The TBC implementation, kept for the port:
	// if paladin.Talents.HealingLight == 0 {
	// 	return
	// }
	//
	// paladin.AddStaticMod(core.SpellModConfig{
	// 	Kind:       core.SpellMod_DamageDone_Flat,
	// 	FloatValue: 0.04 * float64(paladin.Talents.HealingLight),
	// 	ClassMask:  SpellMaskHolyLight | SpellMaskFlashOfLight,
	// })
}

// Illumination - After getting a critical effect from your Flash of Light, Holy Light, or Holy Shock heal spell, you have a 20/40/60/80/100% chance to gain mana equal to 60% of the base cost of the spell
func (paladin *Paladin) applyIllumination() {
	if paladin.Talents.Illumination == 0 {
		return
	}

	// TODO: Implement mana return on crit
}

// TODO: To be implemented. TBC body below already accounts for Forever dropping Purifying Power's crit bonus (pinned to 0, per the TODO inside); kept commented until this class's port is reviewed.
//
// Purifying Power - Reduces the mana cost of your Cleanse and Consecration spells by 5/10%, and increases the critical strike chance of your Exorcism and Holy Wrath spells by 10/20%
func (paladin *Paladin) applyPurifyingPower() {
	if paladin.Talents.PurifyingPower == 0 {
		return
	}

	// The TBC implementation, kept for the port:
	// if paladin.Talents.PurifyingPower == 0 {
	// 	return
	// }
	//
	// paladin.AddStaticMod(core.SpellModConfig{
	// 	Kind:       core.SpellMod_PowerCost_Pct_Add,
	// 	FloatValue: -0.05 * float64(paladin.Talents.PurifyingPower),
	// 	ClassMask:  SpellMaskConsecration, // Cleanse not modeled
	// })
	// // TODO: Forever drops Purifying Power's crit bonus; the spell carries only a cost
	// // (-10% per rank) and a cooldown (-16.5% per rank) modifier, so the crit bonus is pinned
	// // to the untalented 0.
	// exorcismHolyWrathBonusCrit := 0.0
	// paladin.AddStaticMod(core.SpellModConfig{
	// 	Kind:       core.SpellMod_BonusCrit_Percent,
	// 	FloatValue: exorcismHolyWrathBonusCrit,
	// 	ClassMask:  SpellMaskExorcism | SpellMaskHolyWrath,
	// })
}

// TODO: To be implemented. TBC body below needs no porting; kept commented until this class's port is reviewed.
//
// Holy Power (talent) - Increases the critical effect chance of your Holy spells by 1/2/3/4/5%
func (paladin *Paladin) applyHolyPowerTalent() {
	if paladin.Talents.HolyPower == 0 {
		return
	}

	// The TBC implementation, kept for the port:
	// if paladin.Talents.HolyPower == 0 {
	// 	return
	// }
	//
	// paladin.AddStaticMod(core.SpellModConfig{
	// 	Kind:       core.SpellMod_BonusCrit_Percent,
	// 	FloatValue: float64(paladin.Talents.HolyPower),
	// 	School:     core.SpellSchoolHoly,
	// })
}

// applyImprovedHolyStrike implements Improved Holy Strike, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (paladin *Paladin) applyImprovedHolyStrike() {
	if paladin.Talents.ImprovedHolyStrike == 0 {
		return
	}
}

// applyImprovedSeals implements Improved Seals, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (paladin *Paladin) applyImprovedSeals() {
	if paladin.Talents.ImprovedSeals == 0 {
		return
	}
}

// applyInfusionOfLight implements Infusion of Light, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (paladin *Paladin) applyInfusionOfLight() {
	if paladin.Talents.InfusionOfLight == 0 {
		return
	}
}

// applyLightsVigil implements Light's Vigil, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (paladin *Paladin) applyLightsVigil() {
	if !paladin.Talents.LightsVigil {
		return
	}
}

// applyDivinePrecision implements Divine Precision, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (paladin *Paladin) applyDivinePrecision() {
	if paladin.Talents.DivinePrecision == 0 {
		return
	}
}

// applyConsecratedGround implements Consecrated Ground, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (paladin *Paladin) applyConsecratedGround() {
	if paladin.Talents.ConsecratedGround == 0 {
		return
	}
}

// applyReverence implements Reverence, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (paladin *Paladin) applyReverence() {
	if paladin.Talents.Reverence == 0 {
		return
	}
}

// applySpiritualFocus implements Spiritual Focus, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (paladin *Paladin) applySpiritualFocus() {
	if paladin.Talents.SpiritualFocus == 0 {
		return
	}
}

// applyUnyieldingFaith implements Unyielding Faith, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (paladin *Paladin) applyUnyieldingFaith() {
	if paladin.Talents.UnyieldingFaith == 0 {
		return
	}
}

// applyVoiceOfTruth implements Voice of Truth, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (paladin *Paladin) applyVoiceOfTruth() {
	if !paladin.Talents.VoiceOfTruth {
		return
	}
}
