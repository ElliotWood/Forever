package shaman

import (
	"time"

	"github.com/wowsims/forever/sim/core"
)

func (shaman *Shaman) ApplyRestorationTalents() {
	shaman.applyNaturesSwiftness()
	shaman.applyRestorativeTotems()
	shaman.applyTidalMastery()
	shaman.applyTotemicFocus()

	// Forever additions, not yet implemented.
	shaman.applyImprovedHealingWave()
	shaman.applyMindfulness()
	shaman.applyNaturalGrace()
	shaman.applyImprovedReincarnation()
	shaman.applyAncestralHealing()
	shaman.applyHealingFocus()
	shaman.applyWaterShield()
	shaman.applyManaTideTotem()
	shaman.applyHealingWay()
	shaman.applyPurification()
	shaman.applyRiptide()
}

func (shaman *Shaman) applyNaturesSwiftness() {
	if !shaman.Talents.NaturesSwiftness {
		return
	}
	nsAura := shaman.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 16188},
		Label:    "Nature's Swiftness",
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if !spell.Matches(SpellMaskChainLightning | SpellMaskLightningBolt) {
				return
			}
			aura.Deactivate(sim)
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_CastTime_Pct,
		FloatValue: -100,
		ClassMask:  SpellMaskChainLightning | SpellMaskLightningBolt,
	})

	shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 16188},
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMagic,
		Flags:       core.SpellFlagAPL | core.SpellFlagNoOnCastComplete | SpellFlagInstant,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Second * 180,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			nsAura.Activate(sim)
		},
	})
}

func (shaman *Shaman) applyRestorativeTotems() {
	if shaman.Talents.RestorativeTotems == 0 {
		return
	}
	// In totems.go
}

func (shaman *Shaman) applyTidalMastery() {
	if shaman.Talents.TidalMastery == 0 {
		return
	}
	shaman.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_BonusCrit_Percent,
		FloatValue: spellData.TidalMastery.ValueAt(shaman.Talents.TidalMastery),
		ClassMask:  SpellMaskChainLightning | SpellMaskLightningBolt | SpellMaskLightningShield | SpellMaskOverload,
	})
}

func (shaman *Shaman) applyTotemicFocus() {
	if shaman.Talents.TotemicFocus == 0 {
		return
	}
	shaman.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_PowerCost_Pct_Add,
		FloatValue: -0.05 * float64(shaman.Talents.TotemicFocus),
		ClassMask:  SpellMaskTotem,
	})
}

// applyImprovedHealingWave implements Improved Healing Wave, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (shaman *Shaman) applyImprovedHealingWave() {
	if shaman.Talents.ImprovedHealingWave == 0 {
		return
	}

	panic("To be implemented")
}

// applyMindfulness implements Mindfulness, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (shaman *Shaman) applyMindfulness() {
	if shaman.Talents.Mindfulness == 0 {
		return
	}

	panic("To be implemented")
}

// applyNaturalGrace implements Natural Grace, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (shaman *Shaman) applyNaturalGrace() {
	if shaman.Talents.NaturalGrace == 0 {
		return
	}

	panic("To be implemented")
}

// applyImprovedReincarnation implements Improved Reincarnation, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (shaman *Shaman) applyImprovedReincarnation() {
	if shaman.Talents.ImprovedReincarnation == 0 {
		return
	}

	panic("To be implemented")
}

// applyAncestralHealing implements Ancestral Healing, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (shaman *Shaman) applyAncestralHealing() {
	if shaman.Talents.AncestralHealing == 0 {
		return
	}

	panic("To be implemented")
}

// applyHealingFocus implements Healing Focus, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (shaman *Shaman) applyHealingFocus() {
	if shaman.Talents.HealingFocus == 0 {
		return
	}

	panic("To be implemented")
}

// applyWaterShield implements Water Shield, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (shaman *Shaman) applyWaterShield() {
	if !shaman.Talents.WaterShield {
		return
	}

	panic("To be implemented")
}

// applyManaTideTotem implements Mana Tide Totem, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (shaman *Shaman) applyManaTideTotem() {
	if !shaman.Talents.ManaTideTotem {
		return
	}

	panic("To be implemented")
}

// applyHealingWay implements Healing Way, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (shaman *Shaman) applyHealingWay() {
	if shaman.Talents.HealingWay == 0 {
		return
	}

	panic("To be implemented")
}

// applyPurification implements Purification, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (shaman *Shaman) applyPurification() {
	if shaman.Talents.Purification == 0 {
		return
	}

	panic("To be implemented")
}

// applyRiptide implements Riptide, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (shaman *Shaman) applyRiptide() {
	if !shaman.Talents.Riptide {
		return
	}

	panic("To be implemented")
}
