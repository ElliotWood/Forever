package warrior

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

// TODO: Manual review needed -- spell 25286 carries no threat effect, so the 194 is hand-supplied.
var heroicStrikeRank = shared.WithSpellDataFlatThreat(spellData.HeroicStrike, 194).HighestRank()
var heroicStrikeBaseDamage, _ = heroicStrikeRank.Direct.Range()
var cleaveBaseDamage, _ = cleaveRank.Direct.Range()

// TODO: Manual review needed -- spell 20569 carries no threat effect, so the 125 is hand-supplied.
var cleaveRank = shared.WithSpellDataFlatThreat(spellData.Cleave, 125).HighestRank()

func (warrior *Warrior) registerHeroicStrike() {
	spell := warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: heroicStrikeRank.SpellID},
		SpellSchool:    heroicStrikeRank.SpellSchool,
		DefenseType:    heroicStrikeRank.DefenseType,
		ProcMask:       core.ProcMaskMeleeMH,
		Flags:          core.SpellFlagMeleeMetrics,
		ClassSpellMask: SpellMaskHeroicStrike,
		MaxRange:       core.MaxMeleeRange,

		RageCost: core.RageCostOptions{
			Cost:   heroicStrikeRank.Cost,
			Refund: 0.8,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		FlatThreatBonus:  heroicStrikeRank.FlatThreatBonus,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := heroicStrikeBaseDamage + warrior.MHWeaponDamage(sim, spell.MeleeAttackPower(target))
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if !result.Landed() {
				spell.IssueRefund(sim)
			}

			if warrior.curQueueAura != nil {
				warrior.curQueueAura.Deactivate(sim)
			}
		},
	})
	warrior.makeQueueSpellsAndAura(spell)
}

func (warrior *Warrior) registerCleave() {
	const maxTargets int32 = 2

	spell := warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: cleaveRank.SpellID},
		SpellSchool:    cleaveRank.SpellSchool,
		DefenseType:    cleaveRank.DefenseType,
		ProcMask:       core.ProcMaskMeleeMH,
		Flags:          core.SpellFlagMeleeMetrics,
		ClassSpellMask: SpellMaskCleave,
		MaxRange:       core.MaxMeleeRange,

		RageCost: core.RageCostOptions{
			Cost: cleaveRank.Cost,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		FlatThreatBonus:  cleaveRank.FlatThreatBonus,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := cleaveBaseDamage + warrior.MHWeaponDamage(sim, spell.MeleeAttackPower(target))
			spell.CalcCleaveDamage(sim, target, maxTargets, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			spell.DealBatchedAoeDamage(sim)

			if warrior.curQueueAura != nil {
				warrior.curQueueAura.Deactivate(sim)
			}
		},
	})
	warrior.makeQueueSpellsAndAura(spell)
}

func (warrior *Warrior) makeQueueSpellsAndAura(srcSpell *core.Spell) *core.Spell {
	isQueueQueued := false

	queueAura := warrior.RegisterAura(core.Aura{
		Label:    "HS/Cleave Queue Aura-" + srcSpell.ActionID.String(),
		ActionID: srcSpell.ActionID.WithTag(1),
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			isQueueQueued = false
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if warrior.curQueueAura != nil {
				warrior.curQueueAura.Deactivate(sim)
			}
			warrior.PseudoStats.DisableDWMissPenalty = true
			warrior.curQueueAura = aura
			warrior.curQueuedAutoSpell = srcSpell
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.DisableDWMissPenalty = false
			warrior.curQueueAura = nil
			warrior.curQueuedAutoSpell = nil
		},
	})

	queueSpell := warrior.RegisterSpell(core.SpellConfig{
		ActionID:    srcSpell.ActionID.WithTag(1),
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: srcSpell.DefenseType,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL | core.SpellFlagNoMetrics,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.curQueueAura == nil &&
				!isQueueQueued &&
				warrior.CurrentRage() >= srcSpell.Cost.GetCurrentCost() &&
				warrior.queuedRealismICD.IsReady(sim)
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if warrior.queuedRealismICD.IsReady(sim) {
				isQueueQueued = true
				warrior.queuedRealismICD.Use(sim)
				sim.AddPendingAction(&core.PendingAction{
					NextActionAt: sim.CurrentTime + warrior.queuedRealismICD.Duration,
					OnAction: func(sim *core.Simulation) {
						queueAura.Activate(sim)
						isQueueQueued = false
					},
				})
			}
		},
	})

	return queueSpell
}

// Returns true if the regular melee swing should be used, false otherwise.
func (warrior *Warrior) TryHSOrCleave(sim *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
	if !warrior.curQueueAura.IsActive() || (mhSwingSpell.ActionID.Tag != 1 && mhSwingSpell.ActionID.Tag != 1290261) {
		warrior.PseudoStats.DisableDWMissPenalty = false
		return mhSwingSpell
	}

	if !warrior.curQueuedAutoSpell.CanCast(sim, warrior.CurrentTarget) {
		warrior.curQueueAura.Deactivate(sim)
		return mhSwingSpell
	}

	return warrior.curQueuedAutoSpell
}
