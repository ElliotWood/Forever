package warrior

import (
	"time"

	"github.com/wowsims/forever/sim/core"
)

func (warrior *Warrior) registerRetaliation() {
	actionID := core.ActionID{SpellID: 20230}

	attackSpell := warrior.RegisterSpell(core.SpellConfig{
		ClassSpellMask: SpellMaskRetaliationHit,
		ActionID:       core.ActionID{SpellID: 20240},
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMH,
		Flags:          core.SpellFlagMeleeMetrics,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := warrior.MHWeaponDamage(sim, spell.MeleeAttackPower(target))
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})

	aura := warrior.RegisterAura(core.Aura{
		ActionID:  actionID,
		Label:     "Retaliation",
		Duration:  time.Second * 15,
		MaxStacks: 30,
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask.Matches(core.ProcMaskMelee) && result.Landed() && result.Damage > 0 {
				attackSpell.Cast(sim, spell.Unit)
				aura.RemoveStack(sim)
			}
		},
	})

	spell := warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		DefenseType:    core.DefenseTypeMelee,
		ClassSpellMask: SpellMaskRetaliation,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Minute * 30,
			},
			SharedCD: core.Cooldown{
				Timer:    warrior.sharedMCD,
				Duration: time.Minute * 30,
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.StanceMatches(BattleStance)
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			aura.Activate(sim)
			aura.SetStacks(sim, 30)
		},

		RelatedSelfBuff: aura,
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
		// Require manual CD usage
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			return false
		},
	})
}
