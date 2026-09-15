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
const (
	holyStrikeWeaponDamage = 0.4
	holyStrikeMinDamage    = 36.0
	holyStrikeMaxDamage    = 46.0
	holyStrikeManaCost     = 20.0
	holyStrikeCooldown     = time.Second * 12
)

func (paladin *Paladin) registerHolyStrike() {
	// TODO: Only rank 1 of Improved Holy Strike was seen, the second second of cooldown is
	// assumed to scale linearly.
	cooldown := holyStrikeCooldown - time.Second*time.Duration(paladin.Talents.ImprovedHolyStrike)

	// Sacred Arbiter also refreshes the paladin's Judgement effects. Judgement of the Crusader
	// is the only Judgement that leaves anything behind, and it already refreshes off every
	// melee attack the paladin lands, so that half of the talent needs nothing here.
	damageMultiplier := paladin.getWeaponSpecializationModifier()
	if paladin.Talents.SacredArbiter {
		damageMultiplier *= 1.1
	}

	// TODO: Only rank 1 of Iron Creed's threat was seen at 5%, the 5% per rank the tree reads
	// comes from the community talent calculator rather than from a tooltip.
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

	// TODO: Only rank 1 of Iron Creed's damage reduction was seen at 2%, the 2% per rank the
	// tree reads comes from the community talent calculator rather than from a tooltip. The
	// 6 seconds is flat at every rank.
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
