package warrior

import (
	"github.com/wowsims/forever/sim/core"
)

var overpowerRank = spellData.Overpower.BySpellID(11585)
var overpowerBaseDamage, _ = overpowerRank.Direct.Range()

func (warrior *Warrior) registerOverpower() {
	actionID := core.ActionID{SpellID: overpowerRank.SpellID}
	overpowerCD := overpowerRank.Cooldown

	// TODO: Test in-game if OP activation only lasts 5 seconds
	warrior.OverpowerAura = warrior.RegisterAura(core.Aura{
		ActionID: actionID,
		Label:    "Overpower Aura",
		Duration: overpowerCD,
	})

	warrior.MakeProcTriggerAura(core.ProcTrigger{
		Name:               "Overpower - Trigger",
		TriggerImmediately: true,
		Outcome:            core.OutcomeDodge,
		Callback:           core.CallbackOnSpellHitDealt,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			warrior.OverpowerAura.Activate(sim)
		},
	})

	warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskOverpower,
		MaxRange:       core.MaxMeleeRange,

		RageCost: core.RageCostOptions{
			Cost:   overpowerRank.Cost,
			Refund: overpowerRank.MissRefund(),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: overpowerRank.GCD,
			},
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: overpowerCD,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		// TODO: Ingame validation needed
		// TODO: Manual review needed -- the threat coefficient is not in the client.
		ThreatMultiplier: 0.75,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.StanceMatches(BattleStance) && warrior.OverpowerAura.IsActive()
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := overpowerBaseDamage + spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower(target))
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialNoBlockDodgeParry)
			warrior.OverpowerAura.Duration = overpowerCD
			warrior.OverpowerAura.Deactivate(sim)

			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
}
