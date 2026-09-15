package shaman

import (
	"time"

	"github.com/wowsims/classic/sim/core"
)

func (shaman *Shaman) registerStormstrikeSpell() {
	if !shaman.Talents.Stormstrike {
		return
	}

	// Forever's Stormstrike raises only the Nature damage this shaman deals to the
	// target, so it gets an aura of its own rather than the raid-wide debuff, which
	// the Classic ruleset still uses.
	forever := shaman.Env.IsForever()

	stormStrikeAuras := shaman.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		if !forever {
			return core.StormstrikeAura(target)
		}

		return target.RegisterAura(core.Aura{
			Label:    "Stormstrike-" + shaman.Label,
			ActionID: core.ActionID{SpellID: 17364},
			Duration: time.Second * 12,
		})
	})

	if forever {
		for _, target := range shaman.Env.Encounter.TargetUnits {
			target.AddDynamicDamageTakenModifier(func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.Unit == &shaman.Unit && spell.SpellSchool.Matches(core.SpellSchoolNature) && stormStrikeAuras.Get(result.Target).IsActive() {
					result.Damage *= 1.20
				}
			})
		}
	}

	shaman.Stormstrike = shaman.RegisterSpell(core.SpellConfig{
		SpellCode:   SpellCode_ShamanStormstrike,
		ActionID:    core.ActionID{SpellID: 17364},
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       SpellFlagShaman | core.SpellFlagAPL | core.SpellFlagMeleeMetrics,

		ManaCost: core.ManaCostOptions{
			BaseCost: .21,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Second * 20,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := shaman.MHWeaponDamage(sim, spell.MeleeAttackPower(target))
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if result.Landed() {
				aura := stormStrikeAuras.Get(target)
				aura.Activate(sim)
				if !forever {
					aura.SetStacks(sim, aura.MaxStacks)
				}
			}
		},
	})
}
