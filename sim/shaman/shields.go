package shaman

import (
	"time"

	"github.com/wowsims/forever/sim/core"
)

// Package-level state the commented-out implementations used:
// var waterShieldRank = spellData.WaterShield.BySpellID(33736)

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

// TODO: To be implemented. The Forever client ships this as a single unranked class spell:
// it has a SkillLineAbility row but no "Rank N" subtext, so no ladder can be built for it.
func (shaman *Shaman) registerWaterShieldSpell() {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// bonusManaReturned := 0.0
	// mp5 := 50.0
	// if shaman.CouldHaveSetBonus(ItemSetTidefuryRaiment, 4) {
	// 	bonusManaReturned = 56
	// }
	//
	// actionID := core.ActionID{SpellID: waterShieldRank.SpellID}
	// waterShieldManaMetrics := shaman.NewManaMetrics(actionID)
	//
	// shaman.WaterShieldAura = shaman.RegisterAura(core.Aura{
	// 	Label:     "Water Shield",
	// 	ActionID:  actionID,
	// 	Duration:  10 * time.Minute,
	// 	MaxStacks: 3,
	// }).AttachProcTrigger(core.ProcTrigger{
	// 	Name:           "Water Shield Trigger",
	// 	Callback:       core.CallbackOnSpellHitTaken,
	// 	ICD:            3500 * time.Millisecond,
	// 	ClassSpellMask: SpellMaskShieldSelfProc,
	// 	Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
	// 		shaman.WaterShieldAura.RemoveStack(sim)
	// 		shaman.AddMana(sim, waterShieldRank.Direct.Damage(sim)+bonusManaReturned, waterShieldManaMetrics)
	// 	},
	// }).AttachStatBuff(stats.MP5, mp5)
	//
	// shaman.RegisterSpell(core.SpellConfig{
	// 	ActionID:    actionID,
	// 	SpellSchool: core.SpellSchoolNature,
	// 	DefenseType: core.DefenseTypeMagic,
	// 	Flags:       core.SpellFlagAPL | SpellFlagInstant,
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			GCD: waterShieldRank.GCD,
	// 		},
	// 	},
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		shaman.LightningShieldAura.Deactivate(sim)
	// 		shaman.WaterShieldAura.Activate(sim)
	// 		shaman.WaterShieldAura.SetStacks(sim, 3)
	// 	},
	// 	RelatedSelfBuff: shaman.WaterShieldAura,
	// })
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
