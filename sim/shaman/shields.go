package shaman

import (
	"time"

	"github.com/wowsims/forever/sim/core"
)

func (shaman *Shaman) registerShieldsSpells() {
	shaman.registerWaterShieldSpell()
	shaman.registerLightningShieldSpell()
	shaman.registerShieldEffectTriggerSpell()
}

func (shaman *Shaman) registerShieldEffectTriggerSpell() {
	shaman.ShieldSelfProcSpell = shaman.RegisterSpell(core.SpellConfig{
		Flags:          core.SpellFlagNoMetrics | core.SpellFlagNoLogs,
		ClassSpellMask: SpellMaskShieldSelfProc,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, 0, spell.OutcomeAlwaysHit)
		},
	})
}

func (shaman *Shaman) startShieldProcPeriodicAction(sim *core.Simulation) {
	if shaman.SelfBuffs.ShieldProcrate > 0 {
		core.StartPeriodicAction(sim, core.PeriodicActionOptions{
			Period:   60 * time.Second / time.Duration(shaman.SelfBuffs.ShieldProcrate),
			Priority: core.ActionPriorityGCD,
			OnAction: func(sim *core.Simulation) {
				shaman.ShieldSelfProcSpell.Cast(sim, &shaman.Unit)
			},
		})
	}
}

// TODO: To be implemented. The Forever client ships no rank ladder the generator can
// read for this ability -- it survives as a single spell with no "Rank N" subtext and
// no ranked SkillLineAbility row -- so there is no data to build the spell from.
func (shaman *Shaman) registerWaterShieldSpell() {
	// Registered unconditionally, so this returns instead of panicking -- a panic
	// here would stop the sim from starting at all rather than flagging one ability.
}

var lightningShieldRank = spellData.LightningShield.HighestRank()

func (shaman *Shaman) registerLightningShieldSpell() {
	actionID := core.ActionID{SpellID: lightningShieldRank.SpellID}

	lsDamage := shaman.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: lightningShieldRank.SpellID},
		SpellSchool:      lightningShieldRank.SpellSchool,
		DefenseType:      lightningShieldRank.DefenseType,
		ProcMask:         core.ProcMaskEmpty,
		Flags:            SpellFlagShamanSpell,
		ClassSpellMask:   SpellMaskLightningShield,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: lightningShieldRank.Direct.BonusCoefficient(),
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := lightningShieldRank.Direct.Damage(sim)
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})

	shaman.LightningShieldAura = shaman.RegisterAura(core.Aura{
		Label:     "Lightning Shield",
		ActionID:  actionID,
		Duration:  10 * time.Minute,
		MaxStacks: 3,
	}).AttachProcTrigger(core.ProcTrigger{
		Name:           "Lightning Shield Trigger",
		Callback:       core.CallbackOnSpellHitTaken,
		ICD:            3500 * time.Millisecond,
		ClassSpellMask: SpellMaskShieldSelfProc,
		Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
			shaman.LightningShieldAura.RemoveStack(sim)
			lsDamage.Cast(sim, shaman.CurrentTarget)
		},
	})

	shaman.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolNature,
		DefenseType: core.DefenseTypeMagic,
		Flags:       core.SpellFlagAPL | SpellFlagInstant,
		ManaCost: core.ManaCostOptions{
			FlatCost: lightningShieldRank.Cost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: lightningShieldRank.GCD,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			shaman.WaterShieldAura.Deactivate(sim)
			shaman.LightningShieldAura.Activate(sim)
			shaman.LightningShieldAura.SetStacks(sim, 3)
		},
		RelatedSelfBuff: shaman.LightningShieldAura,
	})
}
