package shaman

import (
	"slices"
	"time"

	"github.com/wowsims/classic/sim/core"
)

func (shaman *Shaman) registerStormstrikeSpell() {
	if !shaman.Talents.Stormstrike {
		return
	}

	// Forever's Stormstrike raises only this shaman's damage, so it gets an aura of its own rather than the
	// raid-wide debuff, which the Classic ruleset still uses. The beta client (17364) narrows it to the next
	// Lightning Bolt, Chain Lightning or Earth Shock that lands within 12 sec (one charge, +20%), puts the strike
	// on an 8 sec cooldown and charges a flat 125 mana.
	forever := shaman.Env.IsForever()
	stormstrikeSpellCodes := []int32{SpellCode_ShamanLightningBolt, SpellCode_ShamanChainLightning, SpellCode_ShamanEarthShock}

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
				if spell.Unit != &shaman.Unit || !slices.Contains(stormstrikeSpellCodes, spell.SpellCode) {
					return
				}
				// The charge goes on the first damage calculation, so an overload cast off the same bolt misses out.
				if aura := stormStrikeAuras.Get(result.Target); aura.IsActive() {
					result.Damage *= 1.20
					aura.Deactivate(sim)
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
			BaseCost: core.TernaryFloat64(forever, 0, .21),
			FlatCost: core.TernaryFloat64(forever, 125, 0),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: core.TernaryDuration(forever, time.Second*8, time.Second*20),
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
