package priest

import (
	"github.com/wowsims/forever/sim/common/shared"
)

// Package-level state the commented-out implementations used:
// var VampiricTouchRankMap = spellData.VampiricTouch

// TODO: uncalled -- Forever drops the Vampiric Touch talent that gated this spell, so
// nothing calls this any more. Kept rather than deleted because "the talent is gone"
// and "the spell is gone" are not the same claim, and the client data does not
// distinguish them. Re-gate before wiring it back up.
// TODO: To be implemented. The Forever client ships this as a single unranked class spell:
// it has a SkillLineAbility row but no "Rank N" subtext, so no ladder can be built for it.
func (priest *Priest) registerVampiricTouchSpell(rank shared.SpellData) {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// tick := rank.Periodic.(shared.SpellDataPeriodic)
	//
	// manaMetrics := priest.NewManaMetrics(core.ActionID{SpellID: rank.SpellID}.WithTag(1))
	//
	// priest.RegisterSpell(core.SpellConfig{
	// 	ActionID:       core.ActionID{SpellID: rank.SpellID},
	// 	SpellSchool:    core.SpellSchoolShadow,
	// 	DefenseType:    core.DefenseTypeMagic,
	// 	ProcMask:       core.ProcMaskSpellDamage,
	// 	Flags:          core.SpellFlagAPL,
	// 	ClassSpellMask: PriestSpellVampiricTouch,
	// 	Rank:           rank.Rank,
	// 	MaxRange:       rank.MaxRange,
	//
	// 	ManaCost: core.ManaCostOptions{
	// 		FlatCost: rank.Cost,
	// 	},
	//
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			GCD:      rank.GCD,
	// 			CastTime: rank.CastTime,
	// 		},
	// 	},
	//
	// 	DamageMultiplier:         1,
	// 	DamageMultiplierAdditive: 1,
	// 	ThreatMultiplier:         1,
	//
	// 	Dot: core.DotConfig{
	// 		Aura: core.Aura{
	// 			Label: fmt.Sprintf("VampiricTouch-%d", rank.Rank),
	// 			OnInit: func(aura *core.Aura, sim *core.Simulation) {
	// 				aura.AttachProcTrigger(core.ProcTrigger{
	// 					Name:               "VampiricTouch-ManaReturn",
	// 					CanProcFromProcs:   true, // 34914, 34916, 34917 carry the bit.
	// 					Callback:           core.CallbackOnSpellHitTaken | core.CallbackOnPeriodicDamageTaken,
	// 					ClassSpellMask:     PriestShadowSpells,
	// 					RequireDamageDealt: true,
	// 					Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
	// 						priest.AddMana(sim, result.Damage*0.05, manaMetrics)
	// 					},
	// 				})
	// 			},
	// 		},
	// 		NumberOfTicks:       tick.NumberOfTicks,
	// 		TickLength:          tick.TickLength,
	// 		AffectedByCastSpeed: false,
	// 		BonusCoefficient:    tick.Coef,
	//
	// 		OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
	// 			dot.Snapshot(target, tick.Tick)
	// 		},
	// 		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
	// 			dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
	// 		},
	// 	},
	//
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)
	// 		if result.Landed() {
	// 			spell.Dot(target).Apply(sim)
	// 		}
	// 	},
	//
	// 	ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
	// 		if useSnapshot {
	// 			dot := spell.Dot(target)
	// 			return dot.CalcSnapshotDamage(sim, target, spell.OutcomeExpectedMagicHit)
	// 		}
	// 		return spell.CalcPeriodicDamage(sim, target, tick.Tick, spell.OutcomeExpectedMagicHit)
	// 	},
	// })
}
