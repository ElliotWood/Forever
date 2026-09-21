package warrior

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/spelldata"
)

// TODO: Ingame research needed if HS/Cleave still allow for queueing
var heroicStrikeRank = spellData.HeroicStrike.Highest()
var heroicStrikeBaseDamage = heroicStrikeRank.DamageEffect().Average(core.CharacterLevel)

var cleaveRank = spellData.Cleave.Highest()
var cleaveBaseDamage = cleaveRank.DamageEffect().Average(core.CharacterLevel)

func (warrior *Warrior) registerHeroicStrike() {
	config := spelldata.SpellConfig(&warrior.Unit, heroicStrikeRank, spelldata.Flags(core.SpellFlagMeleeMetrics))
	config.ClassSpellMask = SpellMaskHeroicStrike
	config.ProcMask = core.ProcMaskMeleeMH
	config.DamageMultiplier = 1
	config.ThreatMultiplier = 1
	// TODO: Ingame research needed if this adds flat threat
	config.FlatThreatBonus = 0

	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		baseDamage := heroicStrikeBaseDamage + warrior.MHWeaponDamage(sim, spell.MeleeAttackPower(target))
		result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

		if !result.Landed() {
			spell.IssueRefund(sim)
		}

		if warrior.curQueueAura != nil {
			warrior.curQueueAura.Deactivate(sim)
		}
	}

	warrior.makeQueueSpellsAndAura(warrior.RegisterSpell(config))
}

func (warrior *Warrior) registerCleave() {
	const maxTargets int32 = 2

	config := spelldata.SpellConfig(&warrior.Unit, cleaveRank, spelldata.Flags(core.SpellFlagMeleeMetrics))
	config.ClassSpellMask = SpellMaskCleave
	config.ProcMask = core.ProcMaskMeleeMH
	config.DamageMultiplier = 1
	config.ThreatMultiplier = 1
	// TODO: Ingame research needed if this adds flat threat
	config.FlatThreatBonus = 0

	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		baseDamage := cleaveBaseDamage + warrior.MHWeaponDamage(sim, spell.MeleeAttackPower(target))
		results := spell.CalcCleaveDamage(sim, target, maxTargets, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
		spell.DealBatchedAoeDamage(sim)
		if !results[0].Landed() {
			spell.IssueRefund(sim)
		}

		if warrior.curQueueAura != nil {
			warrior.curQueueAura.Deactivate(sim)
		}
	}

	warrior.makeQueueSpellsAndAura(warrior.RegisterSpell(config))
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
