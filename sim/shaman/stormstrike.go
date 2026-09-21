package shaman

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/stats"
)

var stormstrikeRank = spellData.Stormstrike.HighestRank()
var StormstrikeActionID = core.ActionID{SpellID: stormstrikeRank.SpellID}

func (shaman *Shaman) StormstrikeDebuffAura(target *core.Unit) *core.Aura {
	aura := target.GetOrRegisterAura(core.Aura{
		Label:     "Stormstrike-" + shaman.Label,
		ActionID:  StormstrikeActionID,
		Duration:  stormstrikeRank.Duration,
		MaxStacks: stormstrikeRank.ProcCharges,
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.SpellSchool.Matches(core.SpellSchoolNature) {
				return
			}
			if !result.Landed() || result.Damage == 0 {
				return
			}
			aura.RemoveStack(sim)
		},
	})
	return aura.AttachMultiplicativePseudoStatBuff(
		&target.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexNature],
		1+stormstrikeRank.Effect(shared.A_MOD_SPELL_DAMAGE_FROM_CASTER, 0).Value/100,
	)
}

func (shaman *Shaman) newStormstrikeHitSpellConfig(spellID int32, isMH bool) core.SpellConfig {
	var procMask core.ProcMask
	var actionTag int32

	procMask = core.Ternary(isMH, core.ProcMaskMeleeMHSpecial, core.ProcMaskMeleeOHSpecial)
	actionTag = core.TernaryInt32(isMH, 1, 2)

	return core.SpellConfig{
		ActionID:         core.ActionID{SpellID: spellID}.WithTag(actionTag),
		SpellSchool:      core.SpellSchoolPhysical,
		DefenseType:      core.DefenseTypeMelee,
		ProcMask:         procMask,
		Flags:            core.SpellFlagMeleeMetrics,
		ClassSpellMask:   SpellMaskStormstrikeDamage,
		ThreatMultiplier: 1,
		DamageMultiplier: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			weaponDamage := core.Ternary(isMH, spell.Unit.MHWeaponDamage, spell.Unit.OHWeaponDamage)
			baseDamage := weaponDamage(sim, spell.MeleeAttackPower(target))
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialBlockAndCrit)
		},
	}
}

func (shaman *Shaman) newStormstrikeHitSpell(isMH bool) *core.Spell {
	return shaman.RegisterSpell(shaman.newStormstrikeHitSpellConfig(stormstrikeRank.SpellID, isMH))
}

func (shaman *Shaman) newStormstrikeSpellConfig(spellID int32, ssDebuffAuras *core.AuraArray, mhHit *core.Spell, ohHit *core.Spell) core.SpellConfig {
	stormstrikeSpellConfig := core.SpellConfig{
		ActionID:       core.ActionID{SpellID: spellID},
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskStormstrikeCast,
		ManaCost: core.ManaCostOptions{
			FlatCost: stormstrikeRank.Cost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: stormstrikeRank.Cooldown,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			shaman.StormstrikeCastResult = spell.CalcAndDealOutcome(sim, target, spell.OutcomeMeleeSpecialHitNoHitCounter)
			if shaman.StormstrikeCastResult.Landed() {
				ssDebuffAura := ssDebuffAuras.Get(target)
				ssDebuffAura.Activate(sim)
				ssDebuffAura.SetStacks(sim, ssDebuffAura.MaxStacks)

				if shaman.HasMHWeapon() {
					mhHit.Cast(sim, target)
				}

				if shaman.AutoAttacks.IsDualWielding && shaman.HasOHWeapon() {
					ohHit.Cast(sim, target)
				}
			}
			spell.DisposeResult(shaman.StormstrikeCastResult)
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return (shaman.HasMHWeapon() || shaman.HasOHWeapon())
		},
	}
	return stormstrikeSpellConfig
}

func (shaman *Shaman) registerStormstrikeSpell() {
	mhHit := shaman.newStormstrikeHitSpell(true)
	ohHit := shaman.newStormstrikeHitSpell(false)

	shaman.StormStrikeDebuffAuras = shaman.NewEnemyAuraArray(shaman.StormstrikeDebuffAura)

	shaman.Stormstrike = shaman.RegisterSpell(shaman.newStormstrikeSpellConfig(stormstrikeRank.SpellID, &shaman.StormStrikeDebuffAuras, mhHit, ohHit))
}
