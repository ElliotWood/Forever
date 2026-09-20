package priest

import (
	"time"

	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/stats"
)

func (priest *Priest) registerHolyTalents() {
	// Tier 1
	priest.applyTwilightFocus()
	priest.applyImprovedRenew()
	priest.applyHolySpecialization()

	// Tier 2
	priest.applySpellWarding()
	priest.applyDivineFury()

	// Tier 3
	priest.applyHolyNova()
	priest.applyBlessedRecovery()
	priest.applyInspiration()

	// Tier 4
	priest.applyHolyReach()
	priest.applyImprovedHealing()
	priest.applySearingLight()
	priest.applyBindingHeal()

	// Tier 5
	priest.applyLitanyOfLight()
	priest.applySpiritOfRedemption()
	priest.applySpiritualGuidance()

	// Tier 6
	priest.applySpiritualHealing()

	// Tier 7
	priest.applyPrayerOfMending()
}

// applyTwilightFocus implements Twilight Focus, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applyTwilightFocus() {
	if priest.Talents.TwilightFocus == 0 {
		return
	}

	panic("To be implemented")
}

// applyImprovedRenew implements Improved Renew, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applyImprovedRenew() {
	if priest.Talents.ImprovedRenew == 0 {
		return
	}

	panic("To be implemented")
}

// applyHolySpecialization implements Holy Specialization, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applyHolySpecialization() {
	if priest.Talents.HolySpecialization == 0 {
		return
	}

	panic("To be implemented")
}

// applySpellWarding implements Spell Warding, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applySpellWarding() {
	if priest.Talents.SpellWarding == 0 {
		return
	}

	panic("To be implemented")
}

func (priest *Priest) applyDivineFury() {
	if priest.Talents.DivineFury == 0 {
		return
	}
	// -0.1s per rank
	priest.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_CastTime_Flat,
		TimeValue: time.Millisecond * time.Duration(-100*priest.Talents.DivineFury),
		ClassMask: PriestSpellSmite | PriestSpellHolyFire,
	})
}

func (priest *Priest) applyHolyNova() {
	if !priest.Talents.HolyNova {
		return
	}
	HolyNovaRankMap.RegisterAll(priest.registerHolyNovaSpell)
}

var HolyNovaRankMap = spellData.HolyNova

func (priest *Priest) registerHolyNovaSpell(rank shared.SpellData) {
	priest.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.SpellID},
		SpellSchool:    core.SpellSchoolHoly,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: PriestSpellHolyNova,
		Rank:           rank.Rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: rank.Cost,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: rank.GCD,
			},
		},

		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		BonusCoefficient:         rank.Direct.BonusCoefficient(),
		ThreatMultiplier:         0,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := rank.Direct.Damage(sim)
			spell.CalcAndDealAoeDamage(sim, baseDamage, spell.OutcomeMagicHitAndCrit)

			baseHeal := rank.Direct.Damage(sim)
			spell.CalcAndDealHealing(sim, spell.Unit, baseHeal, spell.OutcomeHealing)
		},
	})
}

// applyBlessedRecovery implements Blessed Recovery, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applyBlessedRecovery() {
	if priest.Talents.BlessedRecovery == 0 {
		return
	}

	panic("To be implemented")
}

// applyInspiration implements Inspiration, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applyInspiration() {
	if priest.Talents.Inspiration == 0 {
		return
	}

	panic("To be implemented")
}

// applyHolyReach implements Holy Reach, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applyHolyReach() {
	if priest.Talents.HolyReach == 0 {
		return
	}

	panic("To be implemented")
}

// applyImprovedHealing implements Improved Healing, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applyImprovedHealing() {
	if priest.Talents.ImprovedHealing == 0 {
		return
	}

	panic("To be implemented")
}

func (priest *Priest) applySearingLight() {
	if priest.Talents.SearingLight == 0 {
		return
	}
	// +5% damage per rank
	priest.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: spellData.SearingLight.Effect(shared.A_MOD_DAMAGE_PERCENT_DONE, 0).FractionAt(priest.Talents.SearingLight),
		ClassMask:  PriestSpellSmite | PriestSpellHolyFire,
	})
}

// applyBindingHeal implements Binding Heal, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applyBindingHeal() {
	if !priest.Talents.BindingHeal {
		return
	}

	panic("To be implemented")
}

// applyLitanyOfLight implements Litany of Light, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applyLitanyOfLight() {
	if priest.Talents.LitanyOfLight == 0 {
		return
	}

	panic("To be implemented")
}

// Spirit of Redemption's passive half: +5% total Spirit. The on-death form is not modelled.
func (priest *Priest) applySpiritOfRedemption() {
	if !priest.Talents.SpiritOfRedemption {
		return
	}
	priest.MultiplyStat(stats.Spirit, 1.05)
}

func (priest *Priest) applySpiritualGuidance() {
	if priest.Talents.SpiritualGuidance == 0 {
		return
	}
	// 5% of Spirit added to spell damage per rank
	coeff := spellData.SpiritualGuidance.Effect(shared.A_MOD_SPELL_DAMAGE_OF_STAT_PERCENT, 126).FractionAt(priest.Talents.SpiritualGuidance)
	priest.AddStatDependency(stats.Spirit, stats.SpellDamage, coeff) // Only scaling damage for now since no healing sim....yet!
}

// applySpiritualHealing implements Spiritual Healing, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applySpiritualHealing() {
	if priest.Talents.SpiritualHealing == 0 {
		return
	}

	panic("To be implemented")
}

// applyPrayerOfMending implements Prayer of Mending, new in Forever.
//
// TODO: To be implemented. Needs the Forever tooltip and a spellData ladder before
// the effect can be modelled; there is no TBC equivalent to port.
func (priest *Priest) applyPrayerOfMending() {
	if !priest.Talents.PrayerOfMending {
		return
	}

	panic("To be implemented")
}
