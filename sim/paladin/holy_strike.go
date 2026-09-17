package paladin

import (
	"time"

	"github.com/wowsims/classic/sim/core"
)

// Holy Strike is new in Forever and has no Classic ability behind it, but three talents hang off
// it - Improved Holy Strike shortens its cooldown, Iron Creed sharpens its threat and Sacred
// Arbiter its damage. The published tooltip reads 20 mana, melee range, instant, a 12 second
// cooldown, and 40% weapon damage plus 36 to 46 Holy damage. Classic's unused spell 13953 lends
// the name and the icon; Forever's own id for it is 17143, which nothing in the database knows
// about yet.
// TODO: assumed baseline, beta will confirm - only the level 60 rank is modelled, and the flat
// damage is taken from the published tooltip rather than from the game.
// TODO: beta will confirm - Holy damage on the melee hit table, so it rolls partial resists the
// way every other Holy ability here does. Whether a melee-table Holy strike actually partial
// resists is unknown; if it does not, it wants SpellFlagIgnoreResists.
const (
	holyStrikeWeaponDamage = 0.4
	holyStrikeMinDamage    = 36.0
	holyStrikeMaxDamage    = 46.0
	holyStrikeManaCost     = 20.0
	holyStrikeCooldown     = time.Second * 12
)

func (paladin *Paladin) registerHolyStrike() {
	// Rank 2 takes off 2 sec, so the linear reading was right. Confirmed on the beta.
	cooldown := holyStrikeCooldown - time.Second*time.Duration(paladin.Talents.ImprovedHolyStrike)

	// Sacred Arbiter also refreshes the paladin's Judgement effects. Judgement of the Crusader
	// is the only Judgement that leaves anything behind, and it already refreshes off every
	// melee attack the paladin lands, so that half of the talent needs nothing here.
	damageMultiplier := paladin.getWeaponSpecializationModifier()
	if paladin.Talents.SacredArbiter {
		damageMultiplier *= 1.1
	}

	// 5% per rank, confirmed on the beta at ranks 2, 3 and 4: 10%, 15% and 20%.
	threatMultiplier := 1 + 0.05*float64(paladin.Talents.IronCreed)

	ironCreedAura := paladin.registerIronCreedAura()

	paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 13953},
		SpellCode:   SpellCode_PaladinHolyStrike,
		SpellSchool: core.SpellSchoolHoly,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			FlatCost:   holyStrikeManaCost,
			Multiplier: paladin.benediction(),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: cooldown,
			},
		},

		DamageMultiplier: damageMultiplier,
		ThreatMultiplier: threatMultiplier,
		// Holy damage, so spell power feeds it on top of the weapon share and the flat
		// roll. 0.429 is the coefficient every other instant Holy paladin spell uses here
		// - Exorcism, Hammer of Wrath and Holy Shock.
		// TODO: beta will confirm. Reported by AdamRC as right "pretty sure", which the
		// other three agreeing with it supports but does not settle.
		BonusCoefficient: 0.429,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if ironCreedAura != nil {
				ironCreedAura.Activate(sim)
			}

			// A share of weapon damage, so it takes the normalized swing the way every other
			// percentage-of-weapon strike in the sim does.
			baseDamage := holyStrikeWeaponDamage*spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower(target)) +
				sim.Roll(holyStrikeMinDamage, holyStrikeMaxDamage)
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})
}

// The half of Iron Creed that is not threat: Holy Strike shaves the damage the paladin takes,
// but only while Righteous Fury is up. Ardent Defender's spell id stands in for the buff.
func (paladin *Paladin) registerIronCreedAura() *core.Aura {
	if paladin.Talents.IronCreed == 0 || !paladin.Options.RighteousFury {
		return nil
	}

	// 2% per rank, confirmed on the beta at ranks 2, 3 and 4: 4%, 6% and 8%. The 6 seconds
	// is flat at every rank, which those same tooltips show.
	damageTaken := 1 - 0.02*float64(paladin.Talents.IronCreed)

	return paladin.RegisterAura(core.Aura{
		Label:    "Iron Creed",
		ActionID: core.ActionID{SpellID: 31850},
		Duration: time.Second * 6,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			paladin.PseudoStats.DamageTakenMultiplier *= damageTaken
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			paladin.PseudoStats.DamageTakenMultiplier /= damageTaken
		},
	})
}
