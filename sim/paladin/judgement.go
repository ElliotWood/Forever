package paladin

import (
	"github.com/wowsims/forever/sim/core"
)

var JudgementRankMap = spellData.Judgement

// Judgement
// https://www.wowhead.com/forever/spell=20271
//
// Unleash the energy of a Seal spell upon an enemy. Does not consume the Seal. Refer to individual
// Seals for Judgement effect.
//
// The spell itself is a dummy that rolls on the spell hit table; the seal's own judgement spell
// carries the effect and, for the damaging ones, the crit roll.
func (paladin *Paladin) registerJudgement() {
	row := JudgementRankMap.HighestRank()

	paladin.Judgement = paladin.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: row.SpellID},
		SpellSchool:    core.SpellSchoolHoly,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL | core.SpellFlagPassiveSpell | core.SpellFlagNoOnCastComplete,
		ClassSpellMask: SpellMaskJudgement,
		MaxRange:       row.MaxRange,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ManaCost: manaCost(row),
		Cast: core.CastConfig{
			// Off the global cooldown, as the client states.
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    paladin.sharedTimer(&paladin.judgementTimer),
				Duration: row.Cooldown,
			},
		},

		ExtraCastCondition: func(_ *core.Simulation, _ *core.Unit) bool {
			return paladin.activeSeal() != nil
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			seal := paladin.activeSeal()
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				seal.judgement.Cast(sim, target)
			}

			// Let the rotation react when the cooldown ends, since nothing on the GCD marks it.
			pa := sim.GetConsumedPendingActionFromPool()
			pa.NextActionAt = sim.CurrentTime + spell.TimeToReady(sim) + core.SpellBatchWindow
			pa.OnAction = func(sim *core.Simulation) {
				paladin.ReactToEvent(sim, false, false)
			}
			sim.AddPendingAction(pa)
		},
	})

	// Every melee strike that lands refreshes the judgement debuffs on its target.
	paladin.MakeProcTriggerAura(core.ProcTrigger{
		Name:               "Judgement Refresh" + paladin.Label,
		Callback:           core.CallbackOnSpellHitDealt,
		ProcMask:           core.ProcMaskMelee,
		Outcome:            core.OutcomeLanded,
		TriggerImmediately: true,
		Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
			paladin.refreshJudgements(sim, result.Target)
		},
	})
}
