package paladin

import (
	"time"

	"github.com/wowsims/classic/sim/core"
)

// Holy Strike is new in Forever and has no Classic ability behind it, but three talents hang off
// it - Improved Holy Strike shortens its cooldown, Iron Creed sharpens its threat and Sacred
// Arbiter its damage - so none of them mean anything until it exists. The talents describe an
// instant weapon strike dealt as Holy damage on a short cooldown, and that is what it is modelled
// as. The only hard numbers are the 75 mana and the melee range on Classic's unused spell 13953,
// which is also where the name and the icon come from.
// TODO: assumed baseline, beta will confirm - the 110% weapon damage and the 6 second cooldown
// are both guesses.
const holyStrikeWeaponDamage = 1.1

func (paladin *Paladin) registerHolyStrike() {
	// TODO: Only rank 1 of Improved Holy Strike was seen, the second second of cooldown is
	// assumed to scale linearly.
	cooldown := time.Second*6 - time.Second*time.Duration(paladin.Talents.ImprovedHolyStrike)

	// Sacred Arbiter also refreshes the paladin's Judgement effects. Judgement of the Crusader
	// is the only Judgement that leaves anything behind, and it already refreshes off every
	// melee attack the paladin lands, so that half of the talent needs nothing here.
	damageMultiplier := holyStrikeWeaponDamage * paladin.getWeaponSpecializationModifier()
	if paladin.Talents.SacredArbiter {
		damageMultiplier *= 1.1
	}

	// TODO: Every rank of Iron Creed reads the same 5% threat, the rest are assumed to scale
	// linearly.
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
			FlatCost:   75,
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

			baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower(target))
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

	// TODO: Every rank of Iron Creed reads the same 2% for 6 sec, the rest are assumed to scale
	// linearly.
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
